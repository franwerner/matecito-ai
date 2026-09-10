package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/franwerner/matecito-ai/internal/turn"
)

// validMoments is the closed set `--moment` accepts — the two write-moments
// the mechanism exists for (contracts/turn-claimed-by-explicit-call.md).
var validMoments = map[string]bool{
	"change-integration":  true,
	"batch-consolidation": true,
}

// NewTurnCmd returns the visible "turn" command group: claim, release,
// status. Unlike the hidden "hook" group — generated from the handler
// registry because the host invokes it — this one is a call an agent makes
// on purpose, and a person can run directly to see what is happening.
func NewTurnCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "turn",
		GroupID: "status",
		Short:   "Claim, release or inspect the shared-branch turn",
		Long: `turn serializes a write to a shared branch between two concurrent sessions
with one git ref — no interception, no identity. Claim it before one unit
of work, release it with the receipt claim returned, on every path out.`,
	}
	cmd.AddCommand(newTurnClaimCmd())
	cmd.AddCommand(newTurnReleaseCmd())
	cmd.AddCommand(newTurnStatusCmd())
	return cmd
}

func newTurnClaimCmd() *cobra.Command {
	var change, destination, moment, tree string

	cmd := &cobra.Command{
		Use:   "claim",
		Short: "Claim the shared-branch turn for one unit of work",
		Example: `  matecito-ai turn claim --change my-change --destination main --moment change-integration \
    && git merge --ff-only matecito-ai/my-change

  # then release with the receipt the claim above printed (its "token: <sha>" line):
  matecito-ai turn release --token <sha>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !validMoments[moment] {
				return fmt.Errorf("--moment must be change-integration or batch-consolidation, got %q", moment)
			}
			if tree == "" {
				if wd, err := os.Getwd(); err == nil {
					tree = wd
				}
			}

			result := turn.Claim(turn.ClaimOptions{
				Change:      change,
				Tree:        tree,
				Destination: destination,
				Moment:      moment,
			})
			switch result.Outcome {
			case turn.Claimed:
				fmt.Fprintln(cmd.OutOrStdout(), "turn claimed")
				fmt.Fprintln(cmd.OutOrStdout(), "token:", result.Token)
				return nil
			case turn.Held:
				fmt.Fprintln(cmd.OutOrStdout(), result.Report)
				return errTurnHeld
			default: // turn.NotArbitrated
				fmt.Fprintln(cmd.ErrOrStderr(), "turn not arbitrated:", result.Reason)
				return nil
			}
		},
	}

	cmd.Flags().StringVar(&change, "change", "", "The change name — names the queue ref and appears in the blocked report")
	cmd.Flags().StringVar(&destination, "destination", "", "The branch being written to (report material)")
	cmd.Flags().StringVar(&moment, "moment", "", "change-integration or batch-consolidation (report material)")
	cmd.Flags().StringVar(&tree, "tree", "", "The working tree this write comes from (default: cwd)")
	for _, name := range []string{"change", "destination", "moment"} {
		_ = cmd.MarkFlagRequired(name)
	}
	return cmd
}

// errTurnHeld is the sentinel `claim` returns when the turn is held —
// non-nil is all main.go needs to exit non-zero; the report itself was
// already printed above, so the error carries no duplicate text.
var errTurnHeld = errors.New("turn held")

func newTurnReleaseCmd() *cobra.Command {
	var token string

	cmd := &cobra.Command{
		Use:     "release",
		Short:   "Release the shared-branch turn with the receipt claim returned",
		Example: `  matecito-ai turn release --token <sha>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := turn.Release(turn.ReleaseOptions{Token: token})
			if err != nil {
				return err
			}
			switch result.Outcome {
			case turn.Released:
				fmt.Fprintln(cmd.OutOrStdout(), "turn released")
				return nil
			default: // turn.Mismatch
				fmt.Fprintln(cmd.ErrOrStderr(), result.Message)
				return nil
			}
		},
	}

	cmd.Flags().StringVar(&token, "token", "", "The receipt claim returned — required, no fallback")
	_ = cmd.MarkFlagRequired("token")
	return cmd
}

func newTurnStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Report the shared-branch turn without claiming or releasing it",
		RunE: func(cmd *cobra.Command, args []string) error {
			result := turn.Status("")
			fmt.Fprintln(cmd.OutOrStdout(), result.Report)
			return nil
		},
	}
}
