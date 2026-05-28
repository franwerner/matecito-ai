package prereqs

import (
	"fmt"

	"github.com/franwerner/matecito-ai/internal/check"
)

func All() []check.Result {
	return []check.Result{
		detectClaude(),
		detectNode(),
		detectNpm(),
		detectNpx(),
		detectGit(),
		detectGo(),
	}
}

func detectClaude() check.Result {
	return check.RunVersion("claude", "claude", []string{"--version"}, true,
		"Instalá Claude Code: https://docs.claude.com/")
}

func detectNode() check.Result {
	r := check.RunVersion("node", "node", []string{"--version"}, true,
		"Instalá Node.js 18+: https://nodejs.org/")
	if r.Status != check.StatusOK {
		return r
	}
	major, ok := check.ParseMajor(r.Version)
	if !ok {
		r.Status = check.StatusOutdated
		r.Detail = "no se pudo parsear la versión"
		return r
	}
	if major < 18 {
		r.Status = check.StatusOutdated
		r.Detail = fmt.Sprintf("v%d.x (se requiere ≥ 18)", major)
		r.FixHint = "Actualizá Node a 18+: https://nodejs.org/"
	}
	return r
}

func detectNpm() check.Result {
	return check.RunVersion("npm", "npm", []string{"--version"}, true,
		"npm viene con Node.js. Si tenés Node 18+, deberías tener npm.")
}

func detectNpx() check.Result {
	return check.RunVersion("npx", "npx", []string{"--version"}, true,
		"npx viene con Node.js. Si tenés Node 18+, deberías tener npx.")
}

func detectGit() check.Result {
	return check.RunVersion("git", "git", []string{"--version"}, true,
		"Instalá git: https://git-scm.com/")
}

func detectGo() check.Result {
	r := check.RunVersion("go", "go", []string{"version"}, false,
		"Opcional. Sólo si instalás Engram vía `go install`: https://go.dev/")
	if r.Status != check.StatusOK {
		return r
	}
	major, minor, ok := check.ParseMajorMinor(r.Version)
	if !ok {
		return r
	}
	if major < 1 || (major == 1 && minor < 18) {
		r.Status = check.StatusOutdated
		r.Detail = fmt.Sprintf("go%d.%d (se requiere ≥ 1.18)", major, minor)
		r.FixHint = "Actualizá Go a 1.18+: https://go.dev/"
	}
	return r
}
