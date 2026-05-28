package context7

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
		Name:     "context7 MCP",
		Required: true,
		FixHint:  "Registrá context7: claude mcp add --scope user context7 -- npx -y @upstash/context7-mcp@latest",
	}
	if f, ok := mcp.Find("context7"); ok {
		r.Status = check.StatusOK
		r.Detail = f.Describe()
		return r
	}
	r.Status = check.StatusMissing
	r.Detail = "no registrado"
	return r
}
