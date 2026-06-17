// Package development self-registers the development domain's hooks into the
// hook registry via init(). It is imported for side effects only (blank import
// in main); nothing reads its exported symbols.
package development

import (
	"encoding/json"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/franwerner/matecito-ai/internal/hook"
)

func init() {
	hook.Register(hook.Hook{
		Domain:     "development",
		Subcommand: "git-commit-validate",
		Event:      "PreToolUse",
		Matcher:    "Bash",
		If:         "Bash(git commit *)",
		Run:        run,
	})
}

// preToolUsePayload is the minimal shape of the Claude Code PreToolUse hook
// payload we need to extract the Bash command.
type preToolUsePayload struct {
	ToolInput struct {
		Command string `json:"command"`
	} `json:"tool_input"`
}

// blockPatterns are the case-insensitive substrings that hard-block a commit.
// The robot emoji is matched via its UTF-8 bytes directly.
var blockPatterns = []string{
	"co-authored-by",
	"claude",
	"generated with",
	"\U0001F916", // 🤖
}

// conventionalCommit is a permissive approximation of the Conventional Commits
// spec: type[(scope)]: lowercase-subject, no trailing period.
var conventionalCommit = regexp.MustCompile(
	`^(feat|fix|refactor|docs|test|chore|build|perf|style|revert)(\([a-z0-9-]+\))?: [a-z].*[^.]$`,
)

// run parses the PreToolUse hook payload from payloadJSON, extracts the git
// commit message from a -m "..." / -m '...' argument, and applies the
// attribution block and format-warn rules.
//
// Fail-open rules (returns Code=0, Message=""):
//   - payloadJSON is empty, malformed, or has no tool_input.command
//   - the Bash command has no -m <message> argument (interactive/file-based commit)
func run(payloadJSON []byte) hook.Result {
	if len(payloadJSON) == 0 {
		return hook.Result{}
	}

	var p preToolUsePayload
	if err := json.Unmarshal(payloadJSON, &p); err != nil {
		return hook.Result{} // fail open
	}
	bashCmd := p.ToolInput.Command
	if bashCmd == "" {
		return hook.Result{} // fail open
	}

	msg := extractCommitMessage(bashCmd)
	if msg == "" {
		return hook.Result{} // no -m argument — fail open
	}

	// Hard block: AI/Claude attribution (case-insensitive).
	lower := strings.ToLower(msg)
	for _, pat := range blockPatterns {
		// The robot emoji is multi-byte; check the original message too.
		patLower := strings.ToLower(pat)
		if strings.Contains(lower, patLower) || (utf8.RuneCountInString(pat) > 1 && strings.Contains(msg, pat)) {
			return hook.Result{
				Code:    2,
				Message: "BLOCKED: commit message contains AI/Claude attribution (" + pat + "). Remove AI attribution before committing.",
			}
		}
	}

	// Format warn: not Conventional Commits (exit 0, note to stderr).
	if !conventionalCommit.MatchString(msg) {
		return hook.Result{
			Code:    0,
			Message: "WARN: commit message does not follow Conventional Commits (e.g. feat(scope): lowercase subject). Consider revising.",
		}
	}

	return hook.Result{}
}

// extractCommitMessage extracts the value of the first -m <message> argument
// from a git commit command string. It handles both double-quoted and
// single-quoted values. Returns "" when no -m argument is found.
func extractCommitMessage(cmd string) string {
	// Match -m "..." (double quotes) or -m '...' (single quotes).
	// We use a simple state-machine approach rather than a full shell parser
	// because the common case is a single -m flag.
	for i := 0; i < len(cmd); i++ {
		// Look for " -m " or beginning "-m " patterns.
		if cmd[i] != '-' {
			continue
		}
		if i+2 >= len(cmd) {
			break
		}
		if cmd[i+1] != 'm' {
			continue
		}
		// Ensure the char before '-' (if any) is a space or the string start,
		// to avoid matching "--some-option".
		if i > 0 && cmd[i-1] != ' ' {
			continue
		}
		// Skip " -m" and look for the value.
		j := i + 2
		// skip optional spaces after -m
		for j < len(cmd) && cmd[j] == ' ' {
			j++
		}
		if j >= len(cmd) {
			break
		}
		quote := cmd[j]
		if quote == '"' || quote == '\'' {
			end := strings.IndexByte(cmd[j+1:], quote)
			if end < 0 {
				break
			}
			return cmd[j+1 : j+1+end]
		}
		// unquoted value: take until next space
		end := strings.IndexByte(cmd[j:], ' ')
		if end < 0 {
			return cmd[j:]
		}
		return cmd[j : j+end]
	}
	return ""
}
