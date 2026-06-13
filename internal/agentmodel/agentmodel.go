package agentmodel

import (
	"bytes"
	"strings"
)

// ValidModels lists the model identifiers the tool supports.
// These are Claude Code aliases resolved at runtime by Claude Code itself; the
// binary never probes which are available on the running install. "fable" is a
// newer alias — a config that selects it on a Claude that lacks it degrades to
// the agent's frontmatter default at the forwarding layer (see
// payload/core/CLAUDE.md "Model resolution"), not here.
var ValidModels = []string{"fable", "opus", "sonnet", "haiku"}

// Agents is the canonical 10-agent list in declaration order, matching payload/domains/development/agents/ filenames.
// sdd-init leads: it is the setup phase that bootstraps a project before the flow phases run.
var Agents = []string{
	"sdd-init",
	"sdd-intake",
	"sdd-explore",
	"sdd-propose",
	"sdd-spec",
	"sdd-design",
	"sdd-tasks",
	"sdd-apply",
	"sdd-verify",
	"sdd-archive",
}

// IsValidModel reports whether m is one of the three valid model identifiers.
func IsValidModel(m string) bool {
	for _, v := range ValidModels {
		if v == m {
			return true
		}
	}
	return false
}

// IsValidAgent reports whether name is one of the canonical 10 agents.
func IsValidAgent(name string) bool {
	for _, a := range Agents {
		if a == name {
			return true
		}
	}
	return false
}

// ReadModel returns the value of the first non-indented `model:` key found inside
// the opening `---`…`---` YAML frontmatter block. Returns "" when the key is absent
// or there is no frontmatter. Never returns an error today; kept for interface symmetry.
func ReadModel(content []byte) (string, error) {
	lines := bytes.Split(content, []byte("\n"))

	// the first line must be exactly "---" to open a frontmatter block
	if len(lines) == 0 || string(lines[0]) != "---" {
		return "", nil
	}

	for i := 1; i < len(lines); i++ {
		line := lines[i]
		// closing delimiter ends the block
		if string(line) == "---" {
			break
		}
		// only non-indented lines are valid YAML keys
		if len(line) > 0 && line[0] == ' ' || len(line) > 0 && line[0] == '\t' {
			continue
		}
		s := string(line)
		if strings.HasPrefix(s, "model:") {
			val := strings.TrimSpace(strings.TrimPrefix(s, "model:"))
			return val, nil
		}
	}
	return "", nil
}
