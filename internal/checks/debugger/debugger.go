package debugger

import (
	"github.com/franwerner/matecito-ai/internal/check"
	"github.com/franwerner/matecito-ai/internal/mcp"
)

// find is the MCP lookup seam so tests can exercise the registered branch
// without depending on the host's actual MCP configuration.
var find = mcp.Find

func All() []check.Result {
	return []check.Result{
		detectMCP(),
	}
}

func detectMCP() check.Result {
	r := check.Result{
		Name:     "debugger MCP",
		Required: true,
		FixHint:  "Registrá debugger: claude mcp add --scope user debugger -- npx -y @debugmcp/mcp-debugger@latest stdio",
	}
	if f, ok := find("debugger"); ok {
		r.Status = check.StatusOK
		r.Detail = f.Describe()
		return r
	}
	r.Status = check.StatusMissing
	r.Detail = "no registrado"
	return r
}
