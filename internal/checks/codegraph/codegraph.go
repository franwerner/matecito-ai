package codegraph

import (
	"github.com/franwerner/matecito-ai/internal/check"
	"github.com/franwerner/matecito-ai/internal/mcp"
)

func All() []check.Result {
	return []check.Result{
		detectBinary(),
		detectMCP(),
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
