package context7

import (
	"os"

	"github.com/franwerner/matecito-ai/internal/check"
	"github.com/franwerner/matecito-ai/internal/mcp"
)

func All() []check.Result {
	return []check.Result{
		detectMCP(),
		detectAPIKey(),
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

func detectAPIKey() check.Result {
	r := check.Result{
		Name:     "context7 API key",
		Required: false,
		FixHint:  "Si tu plan la requiere: export CONTEXT7_API_KEY=... (o registralo con --api-key).",
	}
	if v := os.Getenv("CONTEXT7_API_KEY"); v != "" {
		r.Status = check.StatusOK
		r.Detail = "CONTEXT7_API_KEY definida"
		return r
	}
	r.Status = check.StatusMissing
	r.Detail = "CONTEXT7_API_KEY no definida"
	return r
}
