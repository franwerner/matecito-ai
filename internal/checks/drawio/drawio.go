package drawio

import (
	"github.com/franwerner/matecito-ai/internal/check"
	"github.com/franwerner/matecito-ai/internal/mcp"
)

func All() []check.Result {
	return []check.Result{
		detectMCP(),
	}
}

func detectMCP() check.Result {
	r := check.Result{
		Name:     "drawio MCP",
		Required: true,
		FixHint:  "Registrá drawio: claude mcp add --scope user drawio -- npx -y @next-ai-drawio/mcp-server@latest",
	}
	if f, ok := mcp.Find("drawio"); ok {
		r.Status = check.StatusOK
		r.Detail = f.Describe()
		return r
	}
	r.Status = check.StatusMissing
	r.Detail = "no registrado"
	return r
}
