package codegraph

import (
	"os"
	"path/filepath"

	"github.com/franwerner/matecito-ai/internal/check"
	"github.com/franwerner/matecito-ai/internal/mcp"
)

func All() []check.Result {
	return []check.Result{
		detectBinary(),
		detectMCP(),
		detectProjectInit(),
	}
}

func detectBinary() check.Result {
	return check.RunVersion("codegraph", "codegraph", []string{"--version"}, true,
		"Instalá CodeGraph: npm install -g @colbymchenry/codegraph")
}

func detectMCP() check.Result {
	r := check.Result{
		Name:     "codegraph MCP",
		Required: true,
		FixHint:  "Registrá CodeGraph: claude mcp add --scope user codegraph -- codegraph serve --mcp",
	}
	if f, ok := mcp.Find("codegraph"); ok {
		r.Status = check.StatusOK
		r.Detail = f.Describe()
		return r
	}
	r.Status = check.StatusMissing
	r.Detail = "no registrado"
	return r
}

func detectProjectInit() check.Result {
	r := check.Result{
		Name:     ".codegraph/",
		Required: false,
		FixHint:  "Inicializá en este proyecto: matecito-ai init (o `codegraph init -i`)",
	}
	cwd, err := os.Getwd()
	if err != nil {
		r.Status = check.StatusMissing
		r.Detail = "no se pudo resolver cwd"
		return r
	}
	p := filepath.Join(cwd, ".codegraph")
	fi, err := os.Stat(p)
	if err != nil || !fi.IsDir() {
		r.Status = check.StatusMissing
		r.Detail = "no existe en cwd"
		return r
	}
	r.Status = check.StatusOK
	r.Detail = "presente en cwd"
	return r
}
