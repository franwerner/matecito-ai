package engram

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/franwerner/matecito-ai/internal/check"
	"github.com/franwerner/matecito-ai/internal/mcp"
)

func All() []check.Result {
	return []check.Result{
		detectBinary(),
		detectDB(),
		detectMCP(),
	}
}

func detectBinary() check.Result {
	return check.RunVersion("engram", "engram", []string{"version"}, true,
		"Instalá Engram: go install github.com/Gentleman-Programming/engram/cmd/engram@latest")
}

func detectDB() check.Result {
	r := check.Result{
		Name:     "engram DB",
		Required: false,
		FixHint:  "Se crea automáticamente al iniciar el MCP por primera vez.",
	}
	home, err := os.UserHomeDir()
	if err != nil {
		r.Status = check.StatusMissing
		r.Detail = "no se pudo resolver $HOME"
		return r
	}
	p := filepath.Join(home, ".engram", "engram.db")
	fi, err := os.Stat(p)
	if err != nil || fi.IsDir() {
		r.Status = check.StatusMissing
		r.Detail = "no existe ~/.engram/engram.db"
		return r
	}
	r.Status = check.StatusOK
	r.Detail = fmt.Sprintf("~/.engram/engram.db (%s)", humanSize(fi.Size()))
	return r
}

func detectMCP() check.Result {
	r := check.Result{
		Name:     "engram MCP",
		Required: true,
		FixHint:  "Registrá Engram: claude plugin marketplace add Gentleman-Programming/engram && claude plugin install engram",
	}
	if f, ok := mcp.Find("engram"); ok {
		r.Status = check.StatusOK
		r.Detail = f.Describe()
		return r
	}
	r.Status = check.StatusMissing
	r.Detail = "no registrado"
	return r
}

func humanSize(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for x := n / unit; x >= unit; x /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(n)/float64(div), "KMGTPE"[exp])
}
