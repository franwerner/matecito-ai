package proofshot

import (
	"github.com/franwerner/matecito-ai/internal/check"
	"github.com/franwerner/matecito-ai/internal/setup/install"
)

// resolveBinary is the install-location seam so tests can exercise every
// branch (including a resolver error) without depending on the host's actual
// install state — mirrors internal/checks/debugger/debugger.go's `var find`.
var resolveBinary = install.ProofshotBinaryPath

func All() []check.Result {
	return []check.Result{
		detectCLI(),
	}
}

func detectCLI() check.Result {
	return check.ProbeAt("proofshot", resolveBinary, []string{"--version"}, false,
		"Instalá proofshot: npm install -g proofshot (https://github.com/AmElmo/proofshot)")
}
