package codegraph

import (
	"github.com/franwerner/matecito-ai/internal/check"
	"github.com/franwerner/matecito-ai/internal/mcp"
	"github.com/franwerner/matecito-ai/internal/setup/install"
)

// resolveBinary is the install-location seam so tests can exercise every
// branch (including a resolver error) without depending on the host's actual
// install state — mirrors internal/checks/debugger/debugger.go's `var find`.
var resolveBinary = install.CodegraphBinaryPath

func All() []check.Result {
	return []check.Result{
		detectBinary(),
		detectMCP(),
	}
}

func detectBinary() check.Result {
	return check.ProbeAt("codegraph", resolveBinary, []string{"--version"}, true,
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
