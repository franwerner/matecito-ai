package cli

import (
	"os/exec"
	"strings"
	"testing"
)

// execTurn runs one turn subcommand end to end through cobra (as
// install_test.go exercises installRunner), capturing stdout/stderr
// separately so a test can tell "printed the report" from "surfaced an
// error" — the same distinction `claim`'s three outcomes depend on.
func execTurn(args ...string) (stdout, stderr string, err error) {
	cmd := NewTurnCmd()
	var out, errBuf strings.Builder
	cmd.SetOut(&out)
	cmd.SetErr(&errBuf)
	cmd.SetArgs(args)
	err = cmd.Execute()
	return out.String(), errBuf.String(), err
}

func requiredFlagErr(t *testing.T, args ...string) {
	t.Helper()
	_, _, err := execTurn(args...)
	if err == nil {
		t.Fatalf("expected a required-flag error for args %v, got nil", args)
	}
}

func TestTurnClaimCmd_RequiredFlags(t *testing.T) {
	t.Run("missing --change", func(t *testing.T) {
		requiredFlagErr(t, "claim", "--destination", "main", "--moment", "change-integration")
	})
	t.Run("missing --destination", func(t *testing.T) {
		requiredFlagErr(t, "claim", "--change", "c1", "--moment", "change-integration")
	})
	t.Run("missing --moment", func(t *testing.T) {
		requiredFlagErr(t, "claim", "--change", "c1", "--destination", "main")
	})
}

func TestTurnClaimCmd_InvalidMoment(t *testing.T) {
	_, _, err := execTurn("claim", "--change", "c1", "--destination", "main", "--moment", "nonsense")
	if err == nil {
		t.Fatal("expected an error for an out-of-set --moment value")
	}
}

func TestTurnReleaseCmd_RequiredToken(t *testing.T) {
	requiredFlagErr(t, "release")
}

// gitScratchRepo builds a real git repository in t.TempDir() and moves the
// test process into it for the duration of the test, restoring the
// original working directory on cleanup. `claim`/`status` resolve their
// repository from the process's own cwd when --tree/--dir is not overridden.
func gitScratchRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if out, err := exec.Command("git", "-C", dir, "init").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	t.Chdir(dir)
	return dir
}

func TestTurnClaimAndReleaseCmd_RoundTrip(t *testing.T) {
	gitScratchRepo(t)

	stdout, stderr, err := execTurn("claim", "--change", "c1", "--destination", "main", "--moment", "change-integration")
	if err != nil {
		t.Fatalf("claim: unexpected error: %v (stderr: %s)", err, stderr)
	}
	if !strings.Contains(stdout, "turn claimed") {
		t.Errorf("stdout %q does not confirm the claim", stdout)
	}
	var token string
	for _, line := range strings.Split(stdout, "\n") {
		if strings.HasPrefix(line, "token:") {
			token = strings.TrimSpace(strings.TrimPrefix(line, "token:"))
		}
	}
	if token == "" {
		t.Fatalf("no token printed in stdout: %q", stdout)
	}

	relOut, relErr, err := execTurn("release", "--token", token)
	if err != nil {
		t.Fatalf("release: unexpected error: %v (stderr: %s)", err, relErr)
	}
	if !strings.Contains(relOut, "turn released") {
		t.Errorf("release stdout %q does not confirm the release", relOut)
	}
}

func TestTurnClaimCmd_HeldExitsNonZero(t *testing.T) {
	gitScratchRepo(t)

	if _, _, err := execTurn("claim", "--change", "c1", "--destination", "main", "--moment", "change-integration"); err != nil {
		t.Fatalf("first claim: %v", err)
	}

	stdout, _, err := execTurn("claim", "--change", "c2", "--destination", "main", "--moment", "change-integration")
	if err == nil {
		t.Fatal("second claim: expected a non-nil error (exit non-zero) while the turn is held")
	}
	if !strings.Contains(stdout, "c1") {
		t.Errorf("held report %q does not name the holder", stdout)
	}
}

func TestTurnReleaseCmd_MismatchReportsWithoutError(t *testing.T) {
	gitScratchRepo(t)

	if _, _, err := execTurn("claim", "--change", "c1", "--destination", "main", "--moment", "change-integration"); err != nil {
		t.Fatalf("claim: %v", err)
	}

	_, stderr, err := execTurn("release", "--token", "not-the-real-token")
	if err != nil {
		t.Fatalf("a mismatch is not a command error, want nil, got: %v", err)
	}
	if stderr == "" {
		t.Error("expected the mismatch to be reported on stderr")
	}
}

func TestTurnStatusCmd_ReportsFreeThenHeld(t *testing.T) {
	gitScratchRepo(t)

	free, _, err := execTurn("status")
	if err != nil {
		t.Fatalf("status: %v", err)
	}
	if !strings.Contains(free, "free") {
		t.Errorf("expected the free turn to be reported, got %q", free)
	}

	if _, _, err := execTurn("claim", "--change", "c1", "--destination", "main", "--moment", "change-integration"); err != nil {
		t.Fatalf("claim: %v", err)
	}

	held, _, err := execTurn("status")
	if err != nil {
		t.Fatalf("status: %v", err)
	}
	if !strings.Contains(held, "c1") {
		t.Errorf("expected the held turn's report to name c1, got %q", held)
	}
}

// TestNewTurnCmd_IsVisible guards the one property the design's intro
// paragraph is explicit about: unlike the hidden "hook" group, "turn" is a
// command an agent calls on purpose and a person can run to see what is
// happening — it must not be Hidden and must not disable its own help.
func TestNewTurnCmd_IsVisible(t *testing.T) {
	cmd := NewTurnCmd()
	if cmd.Hidden {
		t.Error("the turn command group must be visible, not hidden")
	}
	for _, sub := range cmd.Commands() {
		if sub.Hidden {
			t.Errorf("subcommand %q must be visible, not hidden", sub.Use)
		}
	}
}

// TestNewTurnCmd_WiredIntoRoot guards that NewRootCmd actually registers the
// group, so a person running `matecito-ai turn ...` finds it.
func TestNewTurnCmd_WiredIntoRoot(t *testing.T) {
	root := NewRootCmd()
	found := false
	for _, c := range root.Commands() {
		if c.Name() == "turn" {
			found = true
		}
	}
	if !found {
		t.Fatal("NewRootCmd does not register the turn command group")
	}
}

// TestTurnClaimCmd_ExampleDoesNotChainASecondClaim guards against the exact
// bug the shipped --help text had: substituting a second `claim` call into
// `--token` returns the second claim's blocked report (the first claim
// still holds the turn), not a token, guaranteeing release's CAS mismatches.
func TestTurnClaimCmd_ExampleDoesNotChainASecondClaim(t *testing.T) {
	example := newTurnClaimCmd().Example
	if strings.Contains(example, "$(matecito-ai turn claim") {
		t.Error("the example must not substitute a second claim into --token")
	}
	if !strings.Contains(example, "matecito-ai turn release --token") {
		t.Error("the example must show release with the first claim's own token")
	}
}
