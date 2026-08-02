package deploy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	matecitoai "github.com/franwerner/matecito-ai"
)

type Mode int

const (
	ModeFile Mode = iota
	ModeDir
	ModeGrouped
)

type Mapping struct {
	SourceRel string
	TargetRel string
	Mode      Mode
}

// domainComponents are the optional subtrees a domain may ship, mapped to
// their flat destination under ~/.claude. A domain that lacks one is skipped.
var domainComponents = []Mapping{
	{SourceRel: "agents", TargetRel: "agents", Mode: ModeDir},
	{SourceRel: "skills", TargetRel: "skills", Mode: ModeGrouped},
	{SourceRel: "references", TargetRel: "references", Mode: ModeDir},
	// Executable helpers a domain's agents invoke. They share one flat destination like agents and
	// references do, so a script name is global across domains, not scoped to the one that ships it.
	{SourceRel: "scripts", TargetRel: "scripts", Mode: ModeDir},
}

// domainActive reports whether a domain id is in the active set. An empty active
// set means "all domains present in the payload" (compat shim).
func domainActive(active []string, id string) bool {
	if len(active) == 0 {
		return true
	}
	for _, a := range active {
		if a == id {
			return true
		}
	}
	return false
}

// buildMappings discovers the per-domain component plan from the payload
// layout: every ACTIVE domain under domains/ contributes its agents/skills/
// references. Shared components under shared/ are appended unconditionally
// (no active-domain gate). The composed matecito-ai.md (core + domain CLAUDE.md
// fragments) is handled separately in Plan via composeClaudeMdOp.
func buildMappings(payloadFS fs.FS, active []string) ([]Mapping, error) {
	var mappings []Mapping

	domains, err := fs.ReadDir(payloadFS, "domains")
	if err != nil {
		return nil, fmt.Errorf("payload sin domains/: %w", err)
	}

	for _, d := range domains {
		if !d.IsDir() || !domainActive(active, d.Name()) {
			continue
		}
		base := path.Join("domains", d.Name())
		for _, c := range domainComponents {
			src := path.Join(base, c.SourceRel)
			if _, err := fs.Stat(payloadFS, src); err != nil {
				continue
			}
			mappings = append(mappings, Mapping{SourceRel: src, TargetRel: c.TargetRel, Mode: c.Mode})
		}
	}

	// Shared components deploy unconditionally, independent of the active-domain set.
	for _, c := range domainComponents {
		src := path.Join("shared", c.SourceRel)
		if _, err := fs.Stat(payloadFS, src); err != nil {
			continue
		}
		mappings = append(mappings, Mapping{SourceRel: src, TargetRel: c.TargetRel, Mode: c.Mode})
	}

	return mappings, nil
}

type FileStatus int

const (
	StatusNew FileStatus = iota
	StatusChanged
	StatusSame
	// StatusRemoved marks a destination the registry (or, on migration, the
	// embedded legacy catalog) has but today's plan no longer produces —
	// renames, moves and deactivated domains all resolve to this, with no
	// special-casing: diff of sets, nothing more.
	StatusRemoved
)

// FileOp describe la escritura (o el borrado) de un archivo en el filesystem
// real del usuario (Target). Para New/Changed/Same el contenido sale de una de
// dos fuentes: Source (ruta interna al fs.FS del payload, copiada byte a byte)
// o Inline (bytes generados en memoria, usado para matecito-ai.md = core +
// índice de dominios); si Inline != nil tiene prioridad y Source se ignora.
// Todos los archivos se escriben con permisos 0o644.
//
// Para StatusRemoved, Source/Inline no aplican: Target es lo único que se
// borra, y el criterio de borrado depende de la FUENTE de la ruta, no de
// Origin (que sigue siendo puramente descriptivo): RegisteredHash presente
// (el op viene del registro) se borra SIEMPRE — pertenecer al registro ya es
// la prueba de propiedad, lo escribimos y lo anotamos nosotros — y el hash en
// disco sólo decide qué se reporta (borrado limpio si coincide, edición
// pisada si no). LegacyHashes presente (el op viene del catálogo legacy
// embebido, sin registro) conserva el criterio de hash: se borra sólo si el
// contenido coincide con alguno de los hashes históricos, si no se preserva
// sin tocarlo.
type FileOp struct {
	Source string
	Inline []byte
	Target string
	Status FileStatus
	Origin Origin

	// RegisteredHash es el hash que el registro tenía para Target, sólo
	// presente en un op StatusRemoved producido por diff contra el registro.
	RegisteredHash string
	// LegacyHashes es el conjunto de hashes históricos de un op StatusRemoved
	// producido por el catálogo legacy embebido (migración, sin registro).
	LegacyHashes []string
}

// FindPayloadDir busca payload/ en el filesystem real desde start hacia
// arriba. Solo se usa cuando el binario se ejecuta desde un clone del repo
// (modo dev). En binarios distribuidos se prefiere el embedded FS.
func FindPayloadDir(start string) (string, error) {
	dir := start
	for {
		p := filepath.Join(dir, "payload")
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return p, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("no se encontró payload/ en cwd ni en ningún directorio padre")
		}
		dir = parent
	}
}

func ClaudeHome() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".claude"), nil
}

func BackupRoot() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".matecito-ai", "backups"), nil
}

// BackupDir devuelve la carpeta de backup para la corrida actual:
// ~/.matecito-ai/backups/<timestamp>/. NO crea el directorio; el caller
// lo crea solo si llega a haber algo que respaldar (lazy).
func BackupDir() (string, error) {
	root, err := BackupRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, time.Now().Format("20060102-150405")), nil
}

// composeClaudeMdOp builds matecito-ai.md = the core CLAUDE.md plus a generated
// index of the ACTIVE domains. The domain CLAUDE.md bodies are NOT concatenated
// anymore — each is deployed as a standalone file (domainFragmentOps) and loaded
// on demand. The composed bytes are carried Inline (generated, not a payload file).
func composeClaudeMdOp(payloadFS fs.FS, claudeHome string, active []string) (FileOp, error) {
	core, err := fs.ReadFile(payloadFS, "core/CLAUDE.md")
	if err != nil {
		return FileOp{}, fmt.Errorf("payload sin core/CLAUDE.md: %w", err)
	}
	index, err := domainsIndex(payloadFS, active)
	if err != nil {
		return FileOp{}, err
	}
	content := append(bytes.TrimRight(core, "\n"), []byte("\n\n")...)
	content = append(content, index...)
	return FileOp{Inline: content, Target: filepath.Join(claudeHome, "matecito-ai.md")}, nil
}

// domainMeta is the slice of a domain manifest needed to render the index. It is
// parsed locally (not via the manifest package) to avoid an import cycle
// (manifest already imports this deploy package).
type domainMeta struct {
	Label   string `json:"label"`
	Summary string `json:"summary"`
}

func readDomainMeta(payloadFS fs.FS, id string) domainMeta {
	var m domainMeta
	if data, err := fs.ReadFile(payloadFS, path.Join("domains", id, "manifest.json")); err == nil {
		_ = json.Unmarshal(data, &m)
	}
	if m.Label == "" {
		m.Label = id
	}
	return m
}

// domainsIndex renders the markdown index of the active domains that ship a
// CLAUDE.md fragment, in directory order (deterministic). It tells the agent to
// load a domain's fragment on demand (at the latest when intake classifies it).
func domainsIndex(payloadFS fs.FS, active []string) ([]byte, error) {
	domains, err := fs.ReadDir(payloadFS, "domains")
	if err != nil {
		return nil, fmt.Errorf("payload sin domains/: %w", err)
	}
	var b strings.Builder
	b.WriteString("## Active domains — load on demand\n\n")
	b.WriteString("Each active work domain ships a behavior fragment that is NOT loaded with this file. ")
	b.WriteString("When you determine a request belongs to a domain (at the latest when intake classifies it), ")
	b.WriteString("**READ that domain's fragment before applying its rules or dispatching its intake**. ")
	b.WriteString("For cross-domain work, load each that applies. Conceptual questions that execute no domain work need only the kernel.\n\n")
	b.WriteString("| Domain | When it applies | Fragment to load |\n")
	b.WriteString("| --- | --- | --- |\n")
	for _, d := range domains {
		if !d.IsDir() || !domainActive(active, d.Name()) {
			continue
		}
		if _, err := fs.Stat(payloadFS, path.Join("domains", d.Name(), "CLAUDE.md")); err != nil {
			continue
		}
		meta := readDomainMeta(payloadFS, d.Name())
		fmt.Fprintf(&b, "| %s (`%s`) | %s | `~/.claude/matecito-ai/domains/%s.md` |\n", meta.Label, d.Name(), meta.Summary, d.Name())
	}
	return []byte(b.String()), nil
}

// domainFragmentOps deploys each active domain's CLAUDE.md as a standalone file
// under ~/.claude/matecito-ai/domains/<id>.md, loaded on demand per the index.
func domainFragmentOps(payloadFS fs.FS, claudeHome string, active []string) ([]FileOp, error) {
	domains, err := fs.ReadDir(payloadFS, "domains")
	if err != nil {
		return nil, fmt.Errorf("payload sin domains/: %w", err)
	}
	var ops []FileOp
	for _, d := range domains {
		if !d.IsDir() || !domainActive(active, d.Name()) {
			continue
		}
		frag := path.Join("domains", d.Name(), "CLAUDE.md")
		if _, err := fs.Stat(payloadFS, frag); err != nil {
			continue
		}
		ops = append(ops, FileOp{
			Source: frag,
			Target: filepath.Join(claudeHome, "matecito-ai", "domains", d.Name()+".md"),
		})
	}
	return ops, nil
}

// opSource returns the bytes a FileOp would write: the generated Inline bytes
// when set, otherwise the single Source file from the payload.
func opSource(payloadFS fs.FS, op FileOp) ([]byte, error) {
	if op.Inline != nil {
		return op.Inline, nil
	}
	return fs.ReadFile(payloadFS, op.Source)
}

// Plan arma la lista de operaciones de copia: la composición de matecito-ai.md
// (core + fragments de los dominios activos) más los componentes de cada
// dominio activo. active es la lista de dominios habilitados; nil/empty = todos
// los presentes en el payload (shim de compatibilidad).
// payloadFS es la raíz del payload (sea os.DirFS sobre un payload/ local, o
// un sub-FS de PayloadFS embebido). claudeHome es la carpeta real de destino.
// stateDir es la raíz de estado (~/.matecito-ai) donde vive el registro de
// despliegue; vacío resuelve StateDir().
//
// Además de las operaciones de escritura, Plan agrega StatusRemoved por
// diferencia de conjuntos: si hay registro, contra sus entradas; si no lo hay
// (primera corrida, migración), contra el catálogo legacy embebido — nunca
// ambos, y nunca se borra por ausencia de registro.
func Plan(payloadFS fs.FS, claudeHome, stateDir string, active []string) ([]FileOp, error) {
	if stateDir == "" {
		sd, err := StateDir()
		if err != nil {
			return nil, err
		}
		stateDir = sd
	}

	mappings, err := buildMappings(payloadFS, active)
	if err != nil {
		return nil, err
	}

	composed, err := composeClaudeMdOp(payloadFS, claudeHome, active)
	if err != nil {
		return nil, err
	}
	composed.Origin = OriginComposed
	ops := []FileOp{composed}

	fragOps, err := domainFragmentOps(payloadFS, claudeHome, active)
	if err != nil {
		return nil, err
	}
	for i := range fragOps {
		fragOps[i].Origin = originForSource(fragOps[i].Source)
	}
	ops = append(ops, fragOps...)

	for _, m := range mappings {
		more, err := expandMapping(payloadFS, claudeHome, m)
		if err != nil {
			return nil, err
		}
		for i := range more {
			more[i].Origin = originForSource(more[i].Source)
		}
		ops = append(ops, more...)
	}

	seen := map[string]string{}
	for _, op := range ops {
		if prev, dup := seen[op.Target]; dup {
			return nil, clashError(prev, op.Source, op.Target, claudeHome)
		}
		seen[op.Target] = op.Source
	}

	for i := range ops {
		ops[i].Status = computeStatus(payloadFS, ops[i])
	}

	removed, err := removedOps(ops, claudeHome, stateDir)
	if err != nil {
		return nil, err
	}
	ops = append(ops, removed...)

	return ops, nil
}

// originForSource clasifica el Origin de un FileOp a partir de su Source: el
// dominio dueño ("domains/<id>/...") o "shared" ("shared/..."). Usado para
// todo op que no sea el matecito-ai.md compuesto (ese se marca aparte).
func originForSource(source string) Origin {
	if id := domainOf(source); id != "" {
		return domainOrigin(id)
	}
	if strings.HasPrefix(source, "shared/") {
		return OriginShared
	}
	return ""
}

// removedOps calcula los ops StatusRemoved: por diferencia de conjuntos contra
// el registro cuando existe, o barriendo el catálogo legacy embebido cuando no
// (migración). currentOps son las operaciones de escritura ya resueltas por
// Plan (New/Changed/Same) — su conjunto de Target define qué ya NO está "a
// borrar", sea cual sea la fuente de comparación.
//
// Una entrada registrada cuyo Target ya no existe en disco (la persona lo
// borró a mano) NO se emite como op: Apply ya la ignora sin error llegado el
// caso, pero el preview de Plan la mostraba igual — inflando el conteo y la
// lista de "a borrar" con rutas que nada va a tocar. targetExists es el único
// costo nuevo (un stat por entrada registrada).
func removedOps(currentOps []FileOp, claudeHome, stateDir string) ([]FileOp, error) {
	currentTargets := make(map[string]bool, len(currentOps))
	for _, op := range currentOps {
		currentTargets[op.Target] = true
	}

	reg, found, err := LoadRegistry(stateDir)
	if err != nil {
		return nil, err
	}
	if !found {
		return legacyRemovedOps(currentTargets, claudeHome, stateDir), nil
	}

	var ops []FileOp
	for _, e := range reg.Entries {
		target := filepath.Join(claudeHome, filepath.FromSlash(e.Path))
		if currentTargets[target] {
			continue
		}
		if !targetExists(target) {
			continue
		}
		ops = append(ops, FileOp{
			Target:         target,
			Status:         StatusRemoved,
			Origin:         e.Origin,
			RegisteredHash: e.Hash,
		})
	}
	return ops, nil
}

// targetExists reports whether path is currently present on disk. Used to
// keep a StatusRemoved op out of the plan entirely when there is nothing
// left to remove — see removedOps and legacyRemovedOpsFrom.
func targetExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func expandMapping(payloadFS fs.FS, claudeHome string, m Mapping) ([]FileOp, error) {
	source := m.SourceRel
	info, err := fs.Stat(payloadFS, source)
	if err != nil {
		return nil, fmt.Errorf("payload no tiene %q: %w", m.SourceRel, err)
	}

	switch m.Mode {
	case ModeFile:
		if info.IsDir() {
			return nil, fmt.Errorf("se esperaba archivo en %q, es directorio", source)
		}
		return []FileOp{{Source: source, Target: filepath.Join(claudeHome, m.TargetRel)}}, nil

	case ModeDir:
		return walkDir(payloadFS, source, filepath.Join(claudeHome, m.TargetRel))

	case ModeGrouped:
		groups, err := fs.ReadDir(payloadFS, source)
		if err != nil {
			return nil, err
		}
		var ops []FileOp
		for _, g := range groups {
			if !g.IsDir() {
				continue
			}
			groupPath := path.Join(source, g.Name())
			entries, err := fs.ReadDir(payloadFS, groupPath)
			if err != nil {
				return nil, err
			}
			targetBase := filepath.Join(claudeHome, m.TargetRel)
			for _, e := range entries {
				childSrc := path.Join(groupPath, e.Name())
				childDst := filepath.Join(targetBase, e.Name())
				if e.IsDir() {
					more, err := walkDir(payloadFS, childSrc, childDst)
					if err != nil {
						return nil, err
					}
					ops = append(ops, more...)
				} else {
					ops = append(ops, FileOp{Source: childSrc, Target: childDst})
				}
			}
		}
		return ops, nil
	}
	return nil, fmt.Errorf("modo desconocido: %d", m.Mode)
}

func walkDir(payloadFS fs.FS, sourceRoot, targetRoot string) ([]FileOp, error) {
	var ops []FileOp
	err := fs.WalkDir(payloadFS, sourceRoot, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(sourceRoot, p)
		if err != nil {
			return err
		}
		// .gitkeep files exist only to keep empty directories in the repository;
		// they must never be deployed to the user's ~/.claude/ tree.
		if filepath.Base(rel) == ".gitkeep" {
			return nil
		}
		ops = append(ops, FileOp{Source: p, Target: filepath.Join(targetRoot, rel)})
		return nil
	})
	return ops, err
}

func computeStatus(payloadFS fs.FS, op FileOp) FileStatus {
	dst, err := os.ReadFile(op.Target)
	if err != nil {
		return StatusNew
	}
	src, err := opSource(payloadFS, op)
	if err != nil {
		return StatusNew
	}
	if bytes.Equal(src, dst) {
		return StatusSame
	}
	return StatusChanged
}

// Summary resume el plan por categoría. Removals lista qué se va a borrar,
// relativo a claudeHome cuando el destino cae bajo claudeHome (todo caso salvo
// el único legacy entry con root "state"; ese se muestra con su ruta absoluta
// en su lugar, ya que no hay una forma relativa sensata de mostrarlo).
type Summary struct {
	New     int
	Changed int
	Same    int
	Removed int

	Removals []string
}

// Summarize cuenta el plan por categoría. claudeHome relativiza Removals para
// que el desglose se lea igual que el resto del plan (rutas cortas, no
// absolutas) — necesario porque los ops StatusRemoved cargan Target absoluto,
// igual que cualquier otro FileOp.
func Summarize(ops []FileOp, claudeHome string) Summary {
	var s Summary
	for _, op := range ops {
		switch op.Status {
		case StatusNew:
			s.New++
		case StatusChanged:
			s.Changed++
		case StatusSame:
			s.Same++
		case StatusRemoved:
			s.Removed++
			s.Removals = append(s.Removals, displayRemoval(op.Target, claudeHome))
		}
	}
	return s
}

// displayRemoval relativiza target a claudeHome cuando cae bajo él; si no
// (el único caso hoy es el manifiesto abandonado, bajo ~/.matecito-ai), cae a
// la ruta absoluta en vez de un ".." confuso.
func displayRemoval(target, claudeHome string) string {
	rel, err := filepath.Rel(claudeHome, target)
	if err != nil || strings.HasPrefix(rel, "..") {
		return target
	}
	return filepath.ToSlash(rel)
}

// clashError builds the deploy clash message. When both colliding sources come
// from distinct domains it is domain-aware (names the two domains and, for a
// skill target, the skill folder) so the convention violation is obvious:
// each domain must name its skill folders uniquely (no auto-prefix — M4 Option 1).
func clashError(prevSource, source, target, claudeHome string) error {
	d1, d2 := domainOf(prevSource), domainOf(source)
	if d1 != "" && d2 != "" && d1 != d2 {
		if name, ok := skillName(target, claudeHome); ok {
			return fmt.Errorf("clash de skills entre dominios: %q y %q exponen la skill %q (destino %s). Convención: cada dominio nombra único sus carpetas de skills — renombrá la de uno de los dominios", d1, d2, name, target)
		}
		return fmt.Errorf("clash entre dominios %q y %q en %s", d1, d2, target)
	}
	return fmt.Errorf("clash de nombres en deploy: %q y %q apuntan al mismo destino %q", prevSource, source, target)
}

// domainOf extracts the domain id from a payload source path "domains/<id>/...".
func domainOf(source string) string {
	parts := strings.Split(source, "/")
	if len(parts) >= 2 && parts[0] == "domains" {
		return parts[1]
	}
	return ""
}

// skillName returns the skill folder name when target lives under
// <claudeHome>/skills/, i.e. the first path segment below it.
func skillName(target, claudeHome string) (string, bool) {
	rel, err := filepath.Rel(filepath.Join(claudeHome, "skills"), target)
	if err != nil || rel == "." || strings.HasPrefix(rel, "..") {
		return "", false
	}
	if i := strings.IndexRune(rel, filepath.Separator); i >= 0 {
		return rel[:i], true
	}
	return rel, true
}

// ResolvePayloadFS prefiere un payload/ local (modo dev/source) si existe en
// el cwd o algún directorio padre. Si no, cae al payload embebido en el
// binario via go:embed.
func ResolvePayloadFS() (payloadFS fs.FS, source string, err error) {
	if cwd, cwdErr := os.Getwd(); cwdErr == nil {
		if local, findErr := FindPayloadDir(cwd); findErr == nil {
			return os.DirFS(local), local, nil
		}
	}
	sub, err := fs.Sub(matecitoai.PayloadFS, "payload")
	if err != nil {
		return nil, "", err
	}
	return sub, "embedded", nil
}

// ApplyResult resume lo que pasó en una corrida de Apply más allá del error:
// si se creó carpeta de respaldo; qué huérfanos legacy StatusRemoved se
// preservaron porque su contenido en disco no coincidía con ningún hash
// histórico conocido (drift); y qué entradas del registro se borraron pese a
// no coincidir su hash — una edición de la persona que quedó pisada. El
// invocador (sync.go) muestra los dos reportes; ninguno es motivo de fallo.
type ApplyResult struct {
	BackupCreated    bool
	Preserved        []string
	OverwrittenEdits []OverwrittenEdit
}

// OverwrittenEdit describe una entrada del registro que Apply borró aunque su
// contenido en disco ya no coincidía con el hash registrado — la pertenencia
// al registro decide el borrado, el hash sólo decide este reporte. BackupPath
// es dónde quedó la versión de la persona antes de borrarla.
type OverwrittenEdit struct {
	Target     string
	BackupPath string
}

// Apply ejecuta las copias (y los borrados) en disco. Lee de payloadFS y
// escribe en claudeHome. Si hay cambios sobre archivos existentes, o un
// borrado, los respalda en backupDir/<rel-path>. backupDir es la carpeta
// concreta de esta corrida (típicamente de deploy.BackupDir()) — Apply lo crea
// on-demand solo si llega a haber al menos un archivo cambiado o borrado.
// Los archivos de payload se copian byte a byte sin ninguna transformación.
//
// Al terminar, Apply reescribe stateDir/deployed.json con una entrada por cada
// destino que sigue vigente (New/Changed/Same) — ninguna para lo borrado
// (coincida o no el hash cuando el origen es el registro), y ninguna para lo
// preservado por drift (sólo puede pasar con un huérfano legacy: eso no es
// "lo que escribió", dejar de rastrearlo es correcto, no un olvido).
func Apply(payloadFS fs.FS, ops []FileOp, claudeHome, stateDir, backupDir string) (ApplyResult, error) {
	if stateDir == "" {
		sd, err := StateDir()
		if err != nil {
			return ApplyResult{}, err
		}
		stateDir = sd
	}

	var result ApplyResult
	entries := make([]Entry, 0, len(ops))

	for _, op := range ops {
		if op.Status == StatusRemoved {
			deleted, preserved, edit, err := applyRemoval(op, claudeHome, stateDir, backupDir)
			if err != nil {
				return result, err
			}
			if deleted {
				result.BackupCreated = true
			}
			if preserved {
				result.Preserved = append(result.Preserved, op.Target)
			}
			if edit != nil {
				result.OverwrittenEdits = append(result.OverwrittenEdits, *edit)
			}
			continue
		}

		if op.Status != StatusSame {
			if err := os.MkdirAll(filepath.Dir(op.Target), 0o755); err != nil {
				return result, err
			}
		}

		// backup solo para archivos que ya existían y cambiaron
		if op.Status == StatusChanged {
			if _, err := backupFile(op.Target, claudeHome, stateDir, backupDir); err != nil {
				return result, err
			}
			result.BackupCreated = true
		}

		if op.Status != StatusSame {
			if op.Inline != nil {
				if err := os.WriteFile(op.Target, op.Inline, 0o644); err != nil {
					return result, err
				}
			} else if err := copyFromFS(payloadFS, op.Source, op.Target, 0o644); err != nil {
				return result, err
			}
		}

		content, err := opSource(payloadFS, op)
		if err != nil {
			return result, err
		}
		rel, err := relToClaudeHome(op.Target, claudeHome)
		if err != nil {
			return result, err
		}
		entries = append(entries, Entry{Path: rel, Hash: hashBytes(content), Origin: op.Origin})
	}

	sort.Slice(entries, func(i, j int) bool { return entries[i].Path < entries[j].Path })
	reg := Registry{Version: 1, UpdatedAt: time.Now().UTC(), Entries: entries}
	if err := SaveRegistry(stateDir, reg); err != nil {
		return result, err
	}

	return result, nil
}

// applyRemoval decide el destino de un op StatusRemoved, ramificado por su
// FUENTE (ver FileOp): si no hay nada en disco, sale del registro sin error
// ni reporte, para las dos fuentes por igual.
//
// RegisteredHash presente (viene del registro): se respalda y se borra
// SIEMPRE — la pertenencia al registro ya prueba que ese archivo lo
// escribimos nosotros. El hash en disco sólo decide el `edit` devuelto: nil
// si coincidía (borrado limpio), o la entrada de "edición pisada" —con dónde
// quedó el respaldo— si no.
//
// LegacyHashes presente (viene del catálogo legacy embebido, sin registro):
// conserva el criterio viejo, sin cambios — se borra sólo si el contenido
// coincide con alguno de los hashes históricos; si no, se preserva intacto y
// se reporta, nunca como fallo.
func applyRemoval(op FileOp, claudeHome, stateDir, backupDir string) (deleted, preserved bool, edit *OverwrittenEdit, err error) {
	diskHash, err := hashDiskFile(op.Target)
	if err != nil {
		// Ya no está en disco: sale del registro sin error, sin reporte.
		return false, false, nil, nil
	}

	if op.RegisteredHash != "" {
		bk, err := backupFile(op.Target, claudeHome, stateDir, backupDir)
		if err != nil {
			return false, false, nil, err
		}
		if err := os.Remove(op.Target); err != nil {
			return false, false, nil, err
		}
		pruneEmptyDirsAfterDelete(op.Target, claudeHome)
		if diskHash != op.RegisteredHash {
			return true, false, &OverwrittenEdit{Target: op.Target, BackupPath: bk}, nil
		}
		return true, false, nil, nil
	}

	if len(op.LegacyHashes) > 0 && slices.Contains(op.LegacyHashes, diskHash) {
		if _, err := backupFile(op.Target, claudeHome, stateDir, backupDir); err != nil {
			return false, false, nil, err
		}
		if err := os.Remove(op.Target); err != nil {
			return false, false, nil, err
		}
		pruneEmptyDirsAfterDelete(op.Target, claudeHome)
		return true, false, nil, nil
	}

	return false, true, nil, nil
}

// componentRoots are the flat directories deploy ever writes into under
// claudeHome (see domainComponents, plus "matecito-ai" for the per-domain
// fragment tree) — the only directories a deletion is ever allowed to prune
// upward into. They are shared with the person's own files and other
// plugins', so pruning MUST NOT ever remove one of them, empty or not, and
// MUST NOT touch anything outside them at all.
var componentRoots = []string{"agents", "skills", "references", "scripts", "matecito-ai"}

// componentRootFor returns the componentRoots entry target lives under (as an
// absolute path under claudeHome), or ok=false when target isn't under any of
// them — the case for the one legacy entry rooted at stateDir (the abandoned
// deployed-manifest.json, loose under ~/.matecito-ai): that deletion must
// never trigger directory cleanup at all.
func componentRootFor(target, claudeHome string) (root string, ok bool) {
	rel, err := filepath.Rel(claudeHome, target)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", false
	}
	first := rel
	if i := strings.IndexRune(rel, filepath.Separator); i >= 0 {
		first = rel[:i]
	}
	if !slices.Contains(componentRoots, first) {
		return "", false
	}
	return filepath.Join(claudeHome, first), true
}

// pruneEmptyDirsAfterDelete removes the directories a just-deleted target
// left empty, climbing from its parent up to — but never including — its
// component root (see componentRootFor). That bound is the entire difference
// between pruning and destroying a directory shared with the person's own
// files or another plugin's, so it is checked on every step of the climb, not
// assumed from how target/root were derived.
//
// os.Remove fails with ENOTEMPTY on a directory that still has something in
// it — ours or not — and that failure IS the stop condition: trying it
// avoids a separate "is it empty" check that would leave a race window
// between looking empty and still being empty once we act on it.
func pruneEmptyDirsAfterDelete(target, claudeHome string) {
	root, ok := componentRootFor(target, claudeHome)
	if !ok {
		return
	}

	for dir := filepath.Dir(target); ; dir = filepath.Dir(dir) {
		rel, err := filepath.Rel(root, dir)
		if err != nil || strings.HasPrefix(rel, "..") {
			return
		}
		if rel == "." {
			return // dir IS the component root: never removed, even empty.
		}
		if err := os.Remove(dir); err != nil {
			return
		}
	}
}

// backupFile respalda target bajo backupDir antes de sobreescribirlo o
// borrarlo, preservando su ruta relativa a claudeHome (el caso normal) o, para
// el único destino bajo stateDir (el manifiesto abandonado), bajo un
// subdirectorio "state/" — así nunca se produce un ".." confuso en la carpeta
// de respaldo. Devuelve la ruta absoluta donde quedó el respaldo, para que el
// caller pueda reportarla (ver OverwrittenEdit).
func backupFile(target, claudeHome, stateDir, backupDir string) (string, error) {
	bk := filepath.Join(backupDir, backupRelPath(target, claudeHome, stateDir))
	if err := os.MkdirAll(filepath.Dir(bk), 0o755); err != nil {
		return "", err
	}
	if err := copyDiskFile(target, bk); err != nil {
		return "", err
	}
	return bk, nil
}

func backupRelPath(target, claudeHome, stateDir string) string {
	if rel, err := filepath.Rel(claudeHome, target); err == nil && !strings.HasPrefix(rel, "..") {
		return rel
	}
	if rel, err := filepath.Rel(stateDir, target); err == nil && !strings.HasPrefix(rel, "..") {
		return filepath.Join("state", rel)
	}
	return filepath.Base(target)
}

func copyFromFS(payloadFS fs.FS, src, dst string, mode os.FileMode) error {
	in, err := payloadFS.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

func copyDiskFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}
