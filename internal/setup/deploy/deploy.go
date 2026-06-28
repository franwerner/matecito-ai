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
)

// FileOp describe la escritura de un archivo en el filesystem real del usuario
// (Target). El contenido sale de una de dos fuentes: Source (ruta interna al
// fs.FS del payload, copiada byte a byte) o Inline (bytes generados en memoria,
// usado para matecito-ai.md = core + índice de dominios). Si Inline != nil tiene
// prioridad y Source se ignora. Todos los archivos se escriben con permisos 0o644.
type FileOp struct {
	Source string
	Inline []byte
	Target string
	Status FileStatus
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
func Plan(payloadFS fs.FS, claudeHome string, active []string) ([]FileOp, error) {
	mappings, err := buildMappings(payloadFS, active)
	if err != nil {
		return nil, err
	}

	composed, err := composeClaudeMdOp(payloadFS, claudeHome, active)
	if err != nil {
		return nil, err
	}
	ops := []FileOp{composed}

	fragOps, err := domainFragmentOps(payloadFS, claudeHome, active)
	if err != nil {
		return nil, err
	}
	ops = append(ops, fragOps...)

	for _, m := range mappings {
		more, err := expandMapping(payloadFS, claudeHome, m)
		if err != nil {
			return nil, err
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
	return ops, nil
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

type Summary struct {
	New     int
	Changed int
	Same    int
}

func Summarize(ops []FileOp) Summary {
	var s Summary
	for _, op := range ops {
		switch op.Status {
		case StatusNew:
			s.New++
		case StatusChanged:
			s.Changed++
		case StatusSame:
			s.Same++
		}
	}
	return s
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

// Apply ejecuta las copias en disco. Lee de payloadFS y escribe en
// claudeHome. Si hay cambios sobre archivos existentes, los respalda en
// backupDir/<rel-path>. backupDir es la carpeta concreta de esta corrida
// (típicamente de deploy.BackupDir()) — Apply lo crea on-demand solo si
// llega a haber al menos un archivo cambiado.
// Los archivos de payload se copian byte a byte sin ninguna transformación.
func Apply(payloadFS fs.FS, ops []FileOp, claudeHome, backupDir string) (bool, error) {
	backupCreated := false
	for _, op := range ops {
		if op.Status == StatusSame {
			continue
		}

		if err := os.MkdirAll(filepath.Dir(op.Target), 0o755); err != nil {
			return backupCreated, err
		}

		// backup solo para archivos que ya existían y cambiaron
		if op.Status == StatusChanged {
			rel, err := filepath.Rel(claudeHome, op.Target)
			if err != nil {
				return backupCreated, err
			}
			bk := filepath.Join(backupDir, rel)
			if err := os.MkdirAll(filepath.Dir(bk), 0o755); err != nil {
				return backupCreated, err
			}
			if err := copyDiskFile(op.Target, bk); err != nil {
				return backupCreated, err
			}
			backupCreated = true
		}

		if op.Inline != nil {
			if err := os.WriteFile(op.Target, op.Inline, 0o644); err != nil {
				return backupCreated, err
			}
		} else if err := copyFromFS(payloadFS, op.Source, op.Target, 0o644); err != nil {
			return backupCreated, err
		}
	}
	return backupCreated, nil
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
