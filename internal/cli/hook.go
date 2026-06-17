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
// marked Hidden to keep the help output clean. Subcommands are built from the
// compiled-in hook registry.
func NewHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "hook",
		Short:  "Hook handlers invoked by Claude Code (not user-facing)",
		Hidden: true,
	}
	for _, h := range hook.Registered() {
		cmd.AddCommand(newHookSubcommand(h))
	}
	return cmd
}

// newHookSubcommand wraps one registered hook as a cobra subcommand: it reads
// the hook payload from stdin, runs the hook, writes any message to stderr, and
// exits with the hook's code (0 allow, 2 block).
func newHookSubcommand(h hook.Hook) *cobra.Command {
	return &cobra.Command{
		Use:    h.Subcommand,
		Short:  "Hook handler " + h.Id,
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			payload, err := io.ReadAll(os.Stdin)
			if err != nil {
				// stdin read error — fail open (do not block)
				return nil
			}
			result := h.Run(payload)
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
