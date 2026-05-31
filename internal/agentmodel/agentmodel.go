package agentmodel

import (
	"bytes"
	"fmt"
	"strings"
)

// ValidModels lists the model identifiers the tool supports.
var ValidModels = []string{"opus", "sonnet", "haiku"}

// Agents is the canonical 9-agent list in declaration order, matching payload/agents/ filenames.
var Agents = []string{
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

// IsValidAgent reports whether name is one of the canonical 9 agents.
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

// ApplyModelOverride performs a line-targeted replacement of the `model:` line inside
// the first frontmatter block. All other bytes are left byte-identical.
// Returns an error when model is not a ValidModel. Returns original bytes unchanged
// when no `model:` line is found (no-op). Preserves the trailing-newline state of content.
func ApplyModelOverride(content []byte, model string) ([]byte, error) {
	if !IsValidModel(model) {
		return nil, fmt.Errorf("agentmodel: invalid model %q (valid: %v)", model, ValidModels)
	}

	// split while keeping track of EOL style per line
	// We work line-by-line preserving the raw line endings.
	trailingNewline := len(content) > 0 && content[len(content)-1] == '\n'

	// split on '\n' to get lines; each element excludes the '\n'
	rawLines := bytes.Split(content, []byte("\n"))
	// bytes.Split on "\n" adds an empty element at the end when content ends with '\n'

	if len(rawLines) == 0 || string(rawLines[0]) != "---" {
		// no frontmatter; no-op
		out := make([]byte, len(content))
		copy(out, content)
		return out, nil
	}

	inFrontmatter := true
	replaced := false

	for i := 1; i < len(rawLines); i++ {
		if !inFrontmatter {
			break
		}
		line := rawLines[i]
		if string(line) == "---" {
			inFrontmatter = false
			break
		}
		// only non-indented lines are valid keys
		if len(line) > 0 && (line[0] == ' ' || line[0] == '\t') {
			continue
		}
		s := string(line)
		if strings.HasPrefix(s, "model:") {
			rawLines[i] = []byte("model: " + model)
			replaced = true
			break
		}
	}

	if !replaced {
		// no model line found; return original bytes unchanged
		out := make([]byte, len(content))
		copy(out, content)
		return out, nil
	}

	result := bytes.Join(rawLines, []byte("\n"))

	// restore trailing newline state
	hasNewline := len(result) > 0 && result[len(result)-1] == '\n'
	if trailingNewline && !hasNewline {
		result = append(result, '\n')
	} else if !trailingNewline && hasNewline {
		result = result[:len(result)-1]
	}

	return result, nil
}
