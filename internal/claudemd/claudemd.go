package claudemd

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/franwerner/matecito-ai/internal/check"
)

// marker es la línea de import que el install prependa en ~/.claude/CLAUDE.md
// para que Claude Code cargue el matecito-ai.md.
const marker = "@matecito-ai.md"

func All() []check.Result {
	return []check.Result{
		detectMatecitoMd(),
		detectClaudeMd(),
		detectReference(),
	}
}

func detectMatecitoMd() check.Result {
	r := check.Result{
		Name:     "matecito-ai.md",
		Required: true,
		FixHint:  "Corré `matecito-ai install` para deployar el payload",
	}
	path, err := claudeFile("matecito-ai.md")
	if err != nil {
		r.Status = check.StatusMissing
		r.Detail = "no se pudo resolver $HOME"
		return r
	}
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		r.Status = check.StatusMissing
		r.Detail = "no existe ~/.claude/matecito-ai.md"
		return r
	}
	if err != nil || info.IsDir() {
		r.Status = check.StatusMissing
		r.Detail = "no se pudo leer ~/.claude/matecito-ai.md"
		return r
	}
	r.Status = check.StatusOK
	r.Detail = "~/.claude/matecito-ai.md"
	return r
}

func detectClaudeMd() check.Result {
	r := check.Result{
		Name:     "CLAUDE.md",
		Required: true,
		FixHint:  "Corré `matecito-ai install` (lo crea si falta)",
	}
	path, err := claudeFile("CLAUDE.md")
	if err != nil {
		r.Status = check.StatusMissing
		r.Detail = "no se pudo resolver $HOME"
		return r
	}
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		r.Status = check.StatusMissing
		r.Detail = "no existe ~/.claude/CLAUDE.md"
		return r
	}
	if err != nil || info.IsDir() {
		r.Status = check.StatusMissing
		r.Detail = "no se pudo leer ~/.claude/CLAUDE.md"
		return r
	}
	r.Status = check.StatusOK
	r.Detail = "~/.claude/CLAUDE.md"
	return r
}

func detectReference() check.Result {
	r := check.Result{
		Name:     "referencia matecito-ai.md",
		Required: true,
		FixHint:  "Corré `matecito-ai install` (prependa `@matecito-ai.md` al CLAUDE.md)",
	}
	path, err := claudeFile("CLAUDE.md")
	if err != nil {
		r.Status = check.StatusMissing
		r.Detail = "no se pudo resolver $HOME"
		return r
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		r.Status = check.StatusMissing
		r.Detail = "~/.claude/CLAUDE.md no existe"
		return r
	}
	if err != nil {
		r.Status = check.StatusMissing
		r.Detail = "no se pudo leer ~/.claude/CLAUDE.md"
		return r
	}
	if !strings.Contains(string(data), marker) {
		r.Status = check.StatusMissing
		r.Detail = "CLAUDE.md no contiene `" + marker + "`"
		return r
	}
	r.Status = check.StatusOK
	r.Detail = "CLAUDE.md importa `" + marker + "`"
	return r
}

func claudeFile(name string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".claude", name), nil
}
