package hook

import (
	"strings"
	"testing"
)

func TestValidateGitCommit(t *testing.T) {
	// buildPayload wraps a Bash command into a minimal PreToolUse JSON payload.
	buildPayload := func(bashCmd string) []byte {
		return []byte(`{"tool_input":{"command":` + `"` + strings.ReplaceAll(bashCmd, `"`, `\"`) + `"}}`)
	}

	tests := []struct {
		name        string
		payload     []byte
		wantCode    int
		wantMsgFrag string // substring expected in Message; "" means no message
		wantNoMsg   bool   // true when Message must be empty
	}{
		{
			name:      "valid conventional message",
			payload:   buildPayload(`git commit -m "feat(hooks): add validator"`),
			wantCode:  0,
			wantNoMsg: true,
		},
		{
			name:      "valid conventional no scope",
			payload:   buildPayload(`git commit -m "fix: handle nil pointer"`),
			wantCode:  0,
			wantNoMsg: true,
		},
		{
			name:        "Co-Authored-By blocks",
			payload:     buildPayload(`git commit -m "feat: add feature\n\nCo-Authored-By: Claude"`),
			wantCode:    2,
			wantMsgFrag: "BLOCKED",
		},
		{
			name:        "Claude in message blocks",
			payload:     buildPayload(`git commit -m "fix: claude helped write this"`),
			wantCode:    2,
			wantMsgFrag: "BLOCKED",
		},
		{
			name:        "robot emoji blocks",
			payload:     []byte(`{"tool_input":{"command":"git commit -m \"feat: done ` + "\U0001F916" + `\""}}`),
			wantCode:    2,
			wantMsgFrag: "BLOCKED",
		},
		{
			name:        "Generated with blocks",
			payload:     buildPayload(`git commit -m "feat: Generated with Claude Code"`),
			wantCode:    2,
			wantMsgFrag: "BLOCKED",
		},
		{
			name:        "bad format warns only",
			payload:     buildPayload(`git commit -m "Fixed stuff."`),
			wantCode:    0,
			wantMsgFrag: "WARN",
		},
		{
			name:        "capitalized bad format warns",
			payload:     buildPayload(`git commit -m "Added the new feature"`),
			wantCode:    0,
			wantMsgFrag: "WARN",
		},
		{
			name:      "missing -m fails open",
			payload:   buildPayload(`git commit -F msg.txt`),
			wantCode:  0,
			wantNoMsg: true,
		},
		{
			name:      "interactive commit fails open",
			payload:   buildPayload(`git commit`),
			wantCode:  0,
			wantNoMsg: true,
		},
		{
			name:      "invalid JSON fails open",
			payload:   []byte(`{not valid json`),
			wantCode:  0,
			wantNoMsg: true,
		},
		{
			name:      "empty payload fails open",
			payload:   []byte{},
			wantCode:  0,
			wantNoMsg: true,
		},
		{
			name:      "missing tool_input.command fails open",
			payload:   []byte(`{"tool_input":{}}`),
			wantCode:  0,
			wantNoMsg: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateGitCommit(tt.payload)
			if got.Code != tt.wantCode {
				t.Errorf("Code = %d, want %d (message: %q)", got.Code, tt.wantCode, got.Message)
			}
			if tt.wantNoMsg && got.Message != "" {
				t.Errorf("Message = %q, want empty", got.Message)
			}
			if tt.wantMsgFrag != "" && !strings.Contains(got.Message, tt.wantMsgFrag) {
				t.Errorf("Message = %q, want substring %q", got.Message, tt.wantMsgFrag)
			}
		})
	}
}

func TestExtractCommitMessage(t *testing.T) {
	tests := []struct {
		name  string
		cmd   string
		want  string
	}{
		{"double quoted", `git commit -m "feat: do stuff"`, "feat: do stuff"},
		{"single quoted", `git commit -m 'fix: typo'`, "fix: typo"},
		{"no -m", `git commit -F file.txt`, ""},
		{"empty command", ``, ""},
		{"unquoted value", `git commit -m feat`, "feat"},
		{"--amend no -m", `git commit --amend`, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractCommitMessage(tt.cmd)
			if got != tt.want {
				t.Errorf("extractCommitMessage(%q) = %q, want %q", tt.cmd, got, tt.want)
			}
		})
	}
}
