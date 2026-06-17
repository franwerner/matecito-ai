package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/franwerner/matecito-ai/internal/hook"
)

// NewHookCmd returns the hidden "hook" command group. It is invoked by Claude
// Code hook handlers (not by users), so the group and its subcommands are
// marked Hidden to keep the help output clean.
func NewHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "hook",
		Short:  "Hook handlers invoked by Claude Code (not user-facing)",
		Hidden: true,
	}
	cmd.AddCommand(newGitCommitValidateCmd())
	return cmd
}

// newGitCommitValidateCmd returns the "hook git-commit-validate" subcommand.
// It reads the PreToolUse hook payload from stdin, validates the git commit
// message, and exits with the appropriate code (0 allow, 2 block).
func newGitCommitValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "git-commit-validate",
		Short:  "Validate git commit messages from a Claude Code PreToolUse hook payload",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			payload, err := io.ReadAll(os.Stdin)
			if err != nil {
				// stdin read error — fail open (do not block)
				return nil
			}
			result := hook.ValidateGitCommit(payload)
			if result.Message != "" {
				fmt.Fprintln(os.Stderr, result.Message)
			}
			if result.Code != 0 {
				os.Exit(result.Code)
			}
			return nil
		},
	}
}
