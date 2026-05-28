package install

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/franwerner/matecito-ai/internal/deploy"
	"github.com/franwerner/matecito-ai/internal/mcp"
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

	backup, err := backupClaudeJSON()
	if err != nil {
		return fmt.Errorf("falló el backup de ~/.claude.json: %w", err)
	}
	if backup != "" {
		fmt.Fprintf(opts.Stdout, "\nBackup: %s\n", backup)
	}

	for i, s := range plan {
		fmt.Fprintf(opts.Stdout, "\n[%d/%d] %s\n", i+1, len(plan), s.Name)
		if err := s.Run(); err != nil {
			fmt.Fprintf(opts.Stderr, "✗ falló: %v\n", err)
			if backup != "" {
				fmt.Fprintf(opts.Stderr, "Backup intacto en %s\n", backup)
			}
			return fmt.Errorf("install detenido en %q", s.Name)
		}
		fmt.Fprintln(opts.Stdout, "✓ OK")
	}

	fmt.Fprintln(opts.Stdout, "\nListo. Verificá con: matecito-ai verify")
	return nil
}

func AllSteps(opts Options) []Step {
	return []Step{
		engramBinaryStep(opts),
		codegraphBinaryStep(opts),
		engramMCPStep(opts),
		codegraphMCPStep(opts),
		context7MCPStep(opts),
		deployStep(opts),
	}
}

func deployStep(opts Options) Step {
	prep := prepareDeploy()

	return Step{
		Name:  "Deploy del fork (payload/ → ~/.claude/)",
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
			backupRoot, err := deploy.BackupRoot()
			if err != nil {
				return err
			}
			backup, err := deploy.Apply(prep.ops, claudeHome, backupRoot)
			if err != nil {
				return err
			}
			s := deploy.Summarize(prep.ops)
			fmt.Fprintf(opts.Stdout, "  %d nuevos, %d cambiados, %d sin cambio\n", s.New, s.Changed, s.Same)
			if backup != "" {
				fmt.Fprintf(opts.Stdout, "  Backup: %s\n", backup)
			}
			return nil
		},
	}
}

type deployPrep struct {
	payloadDir string
	ops        []deploy.FileOp
	summary    string
	shouldRun  bool
	err        error
}

func prepareDeploy() deployPrep {
	var p deployPrep
	cwd, err := os.Getwd()
	if err != nil {
		p.summary = "no se pudo resolver cwd"
		return p
	}
	p.payloadDir, err = deploy.FindPayload(cwd)
	if err != nil {
		p.summary = "payload/ no encontrado (skip)"
		return p
	}
	claudeHome, err := deploy.ClaudeHome()
	if err != nil {
		p.err = err
		return p
	}
	p.ops, err = deploy.Plan(p.payloadDir, claudeHome)
	if err != nil {
		p.err = err
		p.summary = "error planeando deploy: " + err.Error()
		return p
	}
	s := deploy.Summarize(p.ops)
	p.summary = fmt.Sprintf("desde %s: %d nuevos, %d cambiados (%d sin cambio)",
		p.payloadDir, s.New, s.Changed, s.Same)
	p.shouldRun = s.New+s.Changed > 0
	return p
}

func engramBinaryStep(opts Options) Step {
	method := pickEngramMethod()
	plan := "brew install gentleman-programming/tap/engram"
	if method == "go" {
		plan = "go install github.com/Gentleman-Programming/engram/cmd/engram@latest"
	} else if method == "none" {
		plan = "(requiere go ≥ 1.21 o brew — ninguno disponible)"
	}
	return Step{
		Name: "Engram (binario)",
		Plan: plan,
		Check: func() bool {
			_, err := exec.LookPath("engram")
			return err != nil
		},
		Run: func() error {
			switch method {
			case "go":
				if err := runIO(opts, "go", "install", "github.com/Gentleman-Programming/engram/cmd/engram@latest"); err != nil {
					return err
				}
				return ensureGoBinPath(opts)
			case "brew":
				return runIO(opts, "brew", "install", "gentleman-programming/tap/engram")
			default:
				return errors.New("no se puede instalar Engram: no hay go ≥ 1.21 ni brew")
			}
		},
	}
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

func pickEngramMethod() string {
	if hasGoAtLeast(1, 21) {
		return "go"
	}
	if _, err := exec.LookPath("brew"); err == nil {
		return "brew"
	}
	return "none"
}

func hasGoAtLeast(majorReq, minorReq int) bool {
	if _, err := exec.LookPath("go"); err != nil {
		return false
	}
	out, err := exec.Command("go", "version").CombinedOutput()
	if err != nil {
		return false
	}
	fields := strings.Fields(string(out))
	if len(fields) < 3 {
		return false
	}
	v := strings.TrimPrefix(fields[2], "go")
	parts := strings.Split(v, ".")
	if len(parts) < 2 {
		return false
	}
	major, err1 := strconv.Atoi(parts[0])
	minor, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return false
	}
	if major > majorReq {
		return true
	}
	return major == majorReq && minor >= minorReq
}

func ensureGoBinPath(opts Options) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	binDir := filepath.Join(home, "go", "bin")
	return ensurePathManaged(opts, binDir, "go install")
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
	return ensurePathManaged(opts, binDir, "npm prefix")
}

func ensurePathManaged(opts Options, binDir, source string) error {
	rc, err := userShellRC()
	if err != nil {
		return err
	}
	added, err := appendPathToRC(binDir, rc)
	if err != nil {
		return err
	}
	if added {
		fmt.Fprintf(opts.Stdout, "  PATH (%s) añadido a %s — abrí shell nueva o `source %s`\n", source, rc, rc)
	}
	return nil
}

func userShellRC() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	shell := os.Getenv("SHELL")
	switch {
	case strings.Contains(shell, "zsh"):
		return filepath.Join(home, ".zshrc"), nil
	case strings.Contains(shell, "fish"):
		return filepath.Join(home, ".config", "fish", "config.fish"), nil
	default:
		return filepath.Join(home, ".bashrc"), nil
	}
}

func appendPathToRC(binDir, rcPath string) (bool, error) {
	data, err := os.ReadFile(rcPath)
	if err != nil && !os.IsNotExist(err) {
		return false, err
	}
	if strings.Contains(string(data), binDir) {
		return false, nil
	}
	var line string
	if strings.HasSuffix(rcPath, "config.fish") {
		line = fmt.Sprintf("set -gx PATH %s $PATH\n", binDir)
	} else {
		line = fmt.Sprintf("export PATH=\"%s:$PATH\"\n", binDir)
	}
	if err := os.MkdirAll(filepath.Dir(rcPath), 0o755); err != nil {
		return false, err
	}
	f, err := os.OpenFile(rcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return false, err
	}
	defer f.Close()
	if _, err := f.WriteString("\n# matecito-ai\n" + line); err != nil {
		return false, err
	}
	return true, nil
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

func runIO(opts Options, bin string, args ...string) error {
	c := exec.Command(bin, args...)
	c.Stdin = opts.Stdin
	c.Stdout = opts.Stdout
	c.Stderr = opts.Stderr
	return c.Run()
}

func backupClaudeJSON() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	src := filepath.Join(home, ".claude.json")
	info, err := os.Stat(src)
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		return "", fmt.Errorf("%s es un directorio", src)
	}
	data, err := os.ReadFile(src)
	if err != nil {
		return "", err
	}
	dst := fmt.Sprintf("%s.bak.%s", src, time.Now().Format("20060102-150405"))
	if err := os.WriteFile(dst, data, 0600); err != nil {
		return "", err
	}
	return dst, nil
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
