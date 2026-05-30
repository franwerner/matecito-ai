package install

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	matecitoai "github.com/franwerner/matecito-ai"
	"github.com/franwerner/matecito-ai/internal/deploy"
	"github.com/franwerner/matecito-ai/internal/mcp"
	"github.com/franwerner/matecito-ai/internal/platform"
	"github.com/franwerner/matecito-ai/internal/releasedl"
	"github.com/franwerner/matecito-ai/internal/settings"
)

type Step struct {
	Name  string
	Plan  string
	Check func() bool
	Run   func() error
}

type Options struct {
	DryRun bool
	Yes    bool
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	// BackupDir es la carpeta donde se respaldan todos los archivos que
	// la corrida de install modifique. La setea Run una sola vez para
	// que todas las funciones de backup compartan el mismo timestamp.
	BackupDir string
}

func Run(opts Options) error {
	if opts.Stdin == nil {
		opts.Stdin = os.Stdin
	}
	if opts.Stdout == nil {
		opts.Stdout = os.Stdout
	}
	if opts.Stderr == nil {
		opts.Stderr = os.Stderr
	}
	if opts.BackupDir == "" {
		bd, err := deploy.BackupDir()
		if err != nil {
			return fmt.Errorf("resolviendo backup dir: %w", err)
		}
		opts.BackupDir = bd
	}

	steps := AllSteps(opts)
	plan := make([]Step, 0, len(steps))
	for _, s := range steps {
		if s.Check() {
			plan = append(plan, s)
		}
	}

	if len(plan) == 0 {
		fmt.Fprintln(opts.Stdout, "Nada para hacer — todo está instalado y registrado.")
		return nil
	}

	fmt.Fprintln(opts.Stdout, "Plan:")
	for i, s := range plan {
		fmt.Fprintf(opts.Stdout, "  %d. %s\n     %s\n", i+1, s.Name, s.Plan)
	}

	if opts.DryRun {
		fmt.Fprintln(opts.Stdout, "\n(dry-run) no se ejecutó nada.")
		return nil
	}

	if !opts.Yes {
		if !confirm(opts.Stdin, opts.Stdout, "\n¿Ejecutar? [y/N]: ") {
			fmt.Fprintln(opts.Stdout, "Cancelado.")
			return nil
		}
	}

	if err := backupClaudeJSON(opts.BackupDir); err != nil {
		return fmt.Errorf("falló el backup de ~/.claude.json: %w", err)
	}

	for i, s := range plan {
		fmt.Fprintf(opts.Stdout, "\n[%d/%d] %s\n", i+1, len(plan), s.Name)
		if err := s.Run(); err != nil {
			fmt.Fprintf(opts.Stderr, "✗ falló: %v\n", err)
			if hasBackup(opts.BackupDir) {
				fmt.Fprintf(opts.Stderr, "Backup intacto en %s\n", opts.BackupDir)
			}
			return fmt.Errorf("install detenido en %q", s.Name)
		}
		fmt.Fprintln(opts.Stdout, "✓ OK")
	}

	if hasBackup(opts.BackupDir) {
		fmt.Fprintf(opts.Stdout, "\nBackup: %s\n", opts.BackupDir)
	}
	fmt.Fprintln(opts.Stdout, "\nListo. Verificá con: matecito-ai verify")
	return nil
}

// hasBackup devuelve true si la carpeta de backup existe en disco — i.e. si
// alguno de los pasos de la corrida tuvo algo que respaldar.
func hasBackup(backupDir string) bool {
	if backupDir == "" {
		return false
	}
	info, err := os.Stat(backupDir)
	return err == nil && info.IsDir()
}

func AllSteps(opts Options) []Step {
	return []Step{
		engramBinaryStep(opts),
		codegraphBinaryStep(opts),
		engramMCPStep(opts),
		codegraphMCPStep(opts),
		context7MCPStep(opts),
		deployStep(opts),
		claudeMdReferenceStep(opts),
		mcpPermissionsStep(opts),
	}
}

func mcpPermissionsStep(opts Options) Step {
	return Step{
		Name: "Auto-aprobación de tools del ecosistema (settings.json)",
		Plan: "agregar patrones del ecosistema (MCP + Skill) a permissions.allow en ~/.claude/settings.json (no toca defaultMode ni Bash/Write/Edit)",
		Check: func() bool {
			doc, err := settings.Load()
			if err != nil {
				return false
			}
			return len(settings.MissingPatterns(settings.AllowList(doc))) > 0
		},
		Run: func() error {
			path, err := settings.Path()
			if err != nil {
				return err
			}
			doc, err := settings.Load()
			if err != nil {
				return err
			}
			if !settings.Merge(doc) {
				return nil
			}
			if err := backupSettings(path, opts.BackupDir); err != nil {
				return err
			}
			out, err := json.MarshalIndent(doc, "", "  ")
			if err != nil {
				return err
			}
			out = append(out, '\n')
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return err
			}
			return os.WriteFile(path, out, 0o644)
		},
	}
}

func backupSettings(settingsPath, backupDir string) error {
	data, err := os.ReadFile(settingsPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if err := os.MkdirAll(backupDir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(backupDir, "settings.json"), data, 0o644)
}

func deployStep(opts Options) Step {
	prep := prepareDeploy()

	return Step{
		Name:  "Deploy del fork (payload → ~/.claude/)",
		Plan:  prep.summary,
		Check: func() bool { return prep.shouldRun },
		Run: func() error {
			if prep.err != nil {
				return prep.err
			}
			claudeHome, err := deploy.ClaudeHome()
			if err != nil {
				return err
			}
			_, err = deploy.Apply(prep.payloadFS, prep.ops, claudeHome, opts.BackupDir)
			if err != nil {
				return err
			}
			s := deploy.Summarize(prep.ops)
			fmt.Fprintf(opts.Stdout, "  %d nuevos, %d cambiados, %d sin cambio\n", s.New, s.Changed, s.Same)
			return nil
		},
	}
}

type deployPrep struct {
	payloadFS fs.FS
	source    string
	ops       []deploy.FileOp
	summary   string
	shouldRun bool
	err       error
}

// resolvePayloadFS prefiere un payload/ local (modo dev/source) si existe en
// el cwd o algún directorio padre. Si no, cae al payload embebido en el
// binario via go:embed.
func resolvePayloadFS() (payloadFS fs.FS, source string, err error) {
	if cwd, cwdErr := os.Getwd(); cwdErr == nil {
		if local, findErr := deploy.FindPayloadDir(cwd); findErr == nil {
			return os.DirFS(local), local, nil
		}
	}
	sub, err := fs.Sub(matecitoai.PayloadFS, "payload")
	if err != nil {
		return nil, "", err
	}
	return sub, "embedded", nil
}

func prepareDeploy() deployPrep {
	var p deployPrep
	payloadFS, source, err := resolvePayloadFS()
	if err != nil {
		p.summary = "no se pudo resolver payload: " + err.Error()
		p.err = err
		return p
	}
	p.payloadFS = payloadFS
	p.source = source

	claudeHome, err := deploy.ClaudeHome()
	if err != nil {
		p.err = err
		return p
	}
	p.ops, err = deploy.Plan(p.payloadFS, claudeHome)
	if err != nil {
		p.err = err
		p.summary = "error planeando deploy: " + err.Error()
		return p
	}
	s := deploy.Summarize(p.ops)
	p.summary = fmt.Sprintf("desde %s: %d nuevos, %d cambiados (%d sin cambio)",
		source, s.New, s.Changed, s.Same)
	p.shouldRun = s.New+s.Changed > 0
	return p
}

func engramBinaryStep(opts Options) Step {
	return Step{
		Name: "Engram (binario)",
		Plan: "descargar última release de GitHub → ~/.local/bin/engram (SHA256 verificado)",
		Check: func() bool {
			_, err := exec.LookPath("engram")
			return err != nil
		},
		Run: func() error {
			return InstallEngram(opts)
		},
	}
}

// InstallSelf descarga la última release de matecito-ai desde GitHub,
// verifica SHA256 y reemplaza el binario actualmente en ejecución por la
// nueva versión. Usa os.Executable para resolver dónde escribir.
//
// En Linux y macOS reemplazar el binario en uso es seguro: el kernel mantiene
// la inode abierta y el proceso en curso sigue usando la versión vieja hasta
// que termine; la siguiente invocación ya usa la nueva.
func InstallSelf(opts Options) error {
	plat, err := releasedl.Detect()
	if err != nil {
		return err
	}
	rel, err := releasedl.LatestRelease(releasedl.MatecitoRepo, plat)
	if err != nil {
		return err
	}
	dest, err := os.Executable()
	if err != nil {
		return err
	}
	return releasedl.Download(releasedl.MatecitoRepo, rel, dest, opts.Stdout)
}

// InstallEngram descarga la última release de Engram desde GitHub, verifica
// el checksum SHA256, instala el binario y asegura que la carpeta destino
// esté en PATH. Es reutilizada por el comando `matecito-ai update`.
func InstallEngram(opts Options) error {
	plat, err := releasedl.Detect()
	if err != nil {
		return err
	}
	rel, err := releasedl.LatestRelease(releasedl.EngramRepo, plat)
	if err != nil {
		return err
	}
	dest, err := releasedl.DefaultBinaryPath(releasedl.EngramRepo)
	if err != nil {
		return err
	}
	if err := releasedl.Download(releasedl.EngramRepo, rel, dest, opts.Stdout); err != nil {
		return err
	}
	_, err = platform.Detect().EnsurePathInShell(filepath.Dir(dest), opts.Stdout)
	return err
}

func codegraphBinaryStep(opts Options) Step {
	return Step{
		Name: "CodeGraph (binario)",
		Plan: "npm install -g @colbymchenry/codegraph (configura ~/.npm-global si hace falta)",
		Check: func() bool {
			_, err := exec.LookPath("codegraph")
			return err != nil
		},
		Run: func() error {
			if _, err := exec.LookPath("npm"); err != nil {
				return errors.New("npm no está instalado")
			}
			if err := ensureUserNpmPrefix(opts); err != nil {
				return err
			}
			return runIO(opts, "npm", "install", "-g", "@colbymchenry/codegraph")
		},
	}
}

func isSystemPath(p string) bool {
	switch p {
	case "/usr", "/usr/local", "/opt":
		return true
	}
	return strings.HasPrefix(p, "/usr/") || strings.HasPrefix(p, "/opt/")
}

func ensureUserNpmPrefix(opts Options) error {
	out, err := exec.Command("npm", "config", "get", "prefix").CombinedOutput()
	if err != nil {
		return err
	}
	prefix := strings.TrimSpace(string(out))
	if !isSystemPath(prefix) {
		return nil
	}
	fmt.Fprintf(opts.Stdout, "  npm prefix actual = %q (system-owned) → reconfigurando a ~/.npm-global\n", prefix)
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	userPrefix := filepath.Join(home, ".npm-global")
	if err := os.MkdirAll(userPrefix, 0o755); err != nil {
		return err
	}
	if err := runIO(opts, "npm", "config", "set", "prefix", userPrefix); err != nil {
		return err
	}
	binDir := filepath.Join(userPrefix, "bin")
	os.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	_, err = platform.Detect().EnsurePathInShell(binDir, opts.Stdout)
	return err
}

func engramMCPStep(opts Options) Step {
	return Step{
		Name: "Engram MCP (plugin)",
		Plan: "claude plugin marketplace add Gentleman-Programming/engram && claude plugin install engram",
		Check: func() bool {
			_, ok := mcp.Find("engram")
			return !ok
		},
		Run: func() error {
			if _, err := exec.LookPath("claude"); err != nil {
				return errors.New("claude no está en PATH")
			}
			if err := runIO(opts, "claude", "plugin", "marketplace", "add", "Gentleman-Programming/engram"); err != nil {
				return err
			}
			return runIO(opts, "claude", "plugin", "install", "engram")
		},
	}
}

func codegraphMCPStep(opts Options) Step {
	return Step{
		Name: "CodeGraph MCP",
		Plan: "claude mcp add --scope user codegraph -- codegraph serve --mcp",
		Check: func() bool {
			_, ok := mcp.Find("codegraph")
			return !ok
		},
		Run: func() error {
			if _, err := exec.LookPath("claude"); err != nil {
				return errors.New("claude no está en PATH")
			}
			return runIO(opts, "claude", "mcp", "add", "--scope", "user", "codegraph", "--", "codegraph", "serve", "--mcp")
		},
	}
}

func context7MCPStep(opts Options) Step {
	return Step{
		Name: "context7 MCP",
		Plan: "claude mcp add --scope user context7 -- npx -y @upstash/context7-mcp@latest",
		Check: func() bool {
			_, ok := mcp.Find("context7")
			return !ok
		},
		Run: func() error {
			if _, err := exec.LookPath("claude"); err != nil {
				return errors.New("claude no está en PATH")
			}
			return runIO(opts, "claude", "mcp", "add", "--scope", "user", "context7", "--", "npx", "-y", "@upstash/context7-mcp@latest")
		},
	}
}

// claudeMdMarker es la línea que importa el matecito-ai.md en el CLAUDE.md
// del usuario. Claude Code lee @<archivo>.md como import nativo.
const claudeMdMarker = "@matecito-ai.md"

func claudeMdReferenceStep(opts Options) Step {
	home, _ := os.UserHomeDir()
	claudeMdPath := filepath.Join(home, ".claude", "CLAUDE.md")

	return Step{
		Name: "Referencia en ~/.claude/CLAUDE.md",
		Plan: fmt.Sprintf("prependear `%s` (crea el archivo si no existe; backup automático si ya hay contenido)", claudeMdMarker),
		Check: func() bool {
			data, err := os.ReadFile(claudeMdPath)
			if errors.Is(err, os.ErrNotExist) {
				return true
			}
			if err != nil {
				return false
			}
			return !strings.Contains(string(data), claudeMdMarker)
		},
		Run: func() error {
			existing, err := os.ReadFile(claudeMdPath)
			if errors.Is(err, os.ErrNotExist) {
				if err := os.MkdirAll(filepath.Dir(claudeMdPath), 0o755); err != nil {
					return err
				}
				return os.WriteFile(claudeMdPath, []byte(claudeMdMarker+"\n"), 0o644)
			}
			if err != nil {
				return err
			}
			if err := backupClaudeMd(claudeMdPath, opts.BackupDir); err != nil {
				return err
			}
			newContent := claudeMdMarker + "\n\n" + string(existing)
			return os.WriteFile(claudeMdPath, []byte(newContent), 0o644)
		},
	}
}

func backupClaudeMd(claudeMdPath, backupDir string) error {
	data, err := os.ReadFile(claudeMdPath)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(backupDir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(backupDir, "CLAUDE.md"), data, 0o644)
}

func runIO(opts Options, bin string, args ...string) error {
	c := exec.Command(bin, args...)
	c.Stdin = opts.Stdin
	c.Stdout = opts.Stdout
	c.Stderr = opts.Stderr
	return c.Run()
}

func backupClaudeJSON(backupDir string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	src := filepath.Join(home, ".claude.json")
	info, err := os.Stat(src)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("%s es un directorio", src)
	}
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(backupDir, 0o755); err != nil {
		return err
	}
	dst := filepath.Join(backupDir, ".claude.json")
	return os.WriteFile(dst, data, 0600)
}

func confirm(in io.Reader, out io.Writer, prompt string) bool {
	fmt.Fprint(out, prompt)
	sc := bufio.NewScanner(in)
	if !sc.Scan() {
		return false
	}
	ans := strings.ToLower(strings.TrimSpace(sc.Text()))
	return ans == "y" || ans == "yes" || ans == "s" || ans == "si" || ans == "sí"
}
