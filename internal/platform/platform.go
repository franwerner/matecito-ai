package platform

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Adapter interface {
	Name() string
	EnsurePathInShell(binDir string, out io.Writer) (modified bool, err error)
}

func Detect() Adapter {
	switch runtime.GOOS {
	case "windows":
		return &windowsAdapter{}
	default:
		return &posixAdapter{}
	}
}

type posixAdapter struct{}

func (p *posixAdapter) Name() string { return runtime.GOOS }

func (p *posixAdapter) EnsurePathInShell(binDir string, out io.Writer) (bool, error) {
	rc, err := shellRC()
	if err != nil {
		return false, err
	}
	data, err := os.ReadFile(rc)
	if err != nil && !os.IsNotExist(err) {
		return false, err
	}
	if strings.Contains(string(data), binDir) {
		return false, nil
	}
	var line string
	if strings.HasSuffix(rc, "config.fish") {
		line = fmt.Sprintf("set -gx PATH %s $PATH\n", binDir)
	} else {
		line = fmt.Sprintf("export PATH=\"%s:$PATH\"\n", binDir)
	}
	if err := os.MkdirAll(filepath.Dir(rc), 0o755); err != nil {
		return false, err
	}
	f, err := os.OpenFile(rc, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return false, err
	}
	defer f.Close()
	if _, err := f.WriteString("\n# matecito-ai\n" + line); err != nil {
		return false, err
	}
	fmt.Fprintf(out, "  PATH añadido a %s — `source %s` o abrí shell nueva\n", rc, rc)
	return true, nil
}

func shellRC() (string, error) {
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

type windowsAdapter struct{}

func (w *windowsAdapter) Name() string { return "windows" }

func (w *windowsAdapter) EnsurePathInShell(binDir string, out io.Writer) (bool, error) {
	fmt.Fprintf(out, "  Windows: agregá %s al PATH del usuario manualmente:\n", binDir)
	fmt.Fprintf(out, "    PowerShell:  [Environment]::SetEnvironmentVariable(\"Path\", \"$env:Path;%s\", \"User\")\n", binDir)
	fmt.Fprintf(out, "    cmd:         setx PATH \"%%PATH%%;%s\"\n", binDir)
	fmt.Fprintln(out, "  Después abrí una shell nueva para que tome el nuevo PATH.")
	return false, nil
}
