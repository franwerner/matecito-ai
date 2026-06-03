package proofshot

import (
	"github.com/franwerner/matecito-ai/internal/check"
)

func All() []check.Result {
	return []check.Result{
		detectCLI(),
	}
}

func detectCLI() check.Result {
	return check.RunVersion("proofshot", "proofshot", []string{"--version"}, false,
		"Instalá proofshot: npm install -g proofshot (https://github.com/AmElmo/proofshot)")
}
