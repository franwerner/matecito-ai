package turn

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"
)

func withRunGit(t *testing.T, fn func(dir string, stdin []byte, args ...string) (string, error)) {
	t.Helper()
	orig := runGit
	runGit = fn
	t.Cleanup(func() { runGit = orig })
}

func withSleep(t *testing.T, fn func(time.Duration)) {
	t.Helper()
	orig := sleep
	sleep = fn
	t.Cleanup(func() { sleep = orig })
}

// memGit is an in-memory git backend simulating exactly the plumbing this
// package uses (hash-object, cat-file, update-ref, rev-parse, for-each-ref),
// so tests can exercise the real decision tree end to end without invoking
// git.
type memGit struct {
	blobs map[string]string
	refs  map[string]string
	seq   int
}

func newMemGit() *memGit {
	return &memGit{blobs: map[string]string{}, refs: map[string]string{}}
}

func (g *memGit) run(dir string, stdin []byte, args ...string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("no args")
	}
	switch args[0] {
	case "rev-parse":
		return g.revParse(args[1:])
	case "hash-object":
		g.seq++
		sha := "blob" + itoa(g.seq)
		g.blobs[sha] = string(stdin)
		return sha, nil
	case "cat-file":
		sha := args[len(args)-1]
		body, ok := g.blobs[sha]
		if !ok {
			return "", errors.New("fatal: Not a valid object name " + sha)
		}
		return strings.TrimRight(body, "\n"), nil
	case "update-ref":
		return g.updateRef(args[1:])
	case "for-each-ref":
		return g.forEachRef(args[1:])
	}
	return "", nil
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	digits := ""
	for n > 0 {
		digits = string(rune('0'+n%10)) + digits
		n /= 10
	}
	return digits
}

func (g *memGit) revParse(args []string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("rev-parse: no args")
	}
	switch args[0] {
	case "--git-common-dir":
		return "/repo/.git", nil
	case "HEAD":
		return "deadbeef", nil
	case "--verify":
		ref := args[len(args)-1]
		sha, ok := g.refs[ref]
		if !ok {
			return "", errors.New("fatal: needed a single revision")
		}
		return sha, nil
	}
	return "", nil
}

func (g *memGit) updateRef(args []string) (string, error) {
	if len(args) >= 1 && args[0] == "-d" {
		ref, oldSha := args[1], args[2]
		cur, ok := g.refs[ref]
		if !ok || cur != oldSha {
			return "", errors.New("fatal: cannot lock ref '" + ref + "'")
		}
		delete(g.refs, ref)
		return "", nil
	}
	ref, newSha := args[0], args[1]
	if len(args) >= 3 {
		old := args[2]
		cur, exists := g.refs[ref]
		if old == "" {
			if exists {
				return "", errors.New("fatal: cannot lock ref '" + ref + "': reference already exists")
			}
		} else if cur != old {
			return "", errors.New("fatal: cannot lock ref '" + ref + "': is at " + cur + " but expected " + old)
		}
	}
	g.refs[ref] = newSha
	return "", nil
}

func (g *memGit) forEachRef(args []string) (string, error) {
	format, prefix := "", ""
	for _, a := range args {
		if strings.HasPrefix(a, "--format=") {
			format = strings.TrimPrefix(a, "--format=")
		} else {
			prefix = a
		}
	}
	var lines []string
	for ref, sha := range g.refs {
		if !strings.HasPrefix(ref, prefix) {
			continue
		}
		line := strings.ReplaceAll(format, "%(refname)", ref)
		line = strings.ReplaceAll(line, "%(objectname)", sha)
		lines = append(lines, line)
	}
	sort.Strings(lines)
	return strings.Join(lines, "\n"), nil
}

func opts(change string) ClaimOptions {
	return ClaimOptions{
		Change: change, Tree: "/t/" + change, Destination: "main", Moment: "change-integration", Dir: "/repo",
	}
}

func TestParseRecord(t *testing.T) {
	full := Record{Change: "c", Tree: "/t", Destination: "main", Moment: "change-integration", Base: "sha1", At: "2026-08-15T00:00:00Z"}
	rec, err := parseRecord(full.String())
	if err != nil {
		t.Fatalf("parseRecord(valid) error: %v", err)
	}
	if rec != full {
		t.Errorf("round-trip = %+v, want %+v", rec, full)
	}

	t.Run("missing required field is an error", func(t *testing.T) {
		body := "change: c\ntree: /t\ndestination: main\nmoment: m\nbase: b\n" // no "at"
		if _, err := parseRecord(body); err == nil {
			t.Error("expected an error for a record missing a required field")
		}
	})

	t.Run("malformed line is an error", func(t *testing.T) {
		if _, err := parseRecord("change: c\nnot-a-valid-line\n"); err == nil {
			t.Error("expected an error for a malformed line")
		}
	})

	t.Run("no session or agent field exists", func(t *testing.T) {
		if strings.Contains(full.String(), "session") || strings.Contains(full.String(), "agent") {
			t.Error("the record must carry no identity field at all")
		}
	})
}

func TestClaim_FreeTurnIsClaimed(t *testing.T) {
	g := newMemGit()
	withRunGit(t, g.run)
	got := Claim(opts("c1"))
	if got.Outcome != Claimed {
		t.Fatalf("Outcome = %v, want Claimed", got.Outcome)
	}
	if got.Token == "" {
		t.Error("Token is empty — the receipt must be non-empty")
	}
	if got.Report != "" || got.Reason != "" {
		t.Errorf("Claimed result must carry no Report/Reason, got %+v", got)
	}
}

func TestClaim_ClaimingLeavesNoOtherRefBehindOnSuccess(t *testing.T) {
	g := newMemGit()
	withRunGit(t, g.run)
	Claim(opts("c1"))
	if len(g.refs) != 1 {
		t.Errorf("expected exactly one ref (the lock) after a clean claim, got %v", g.refs)
	}
}

func TestClaim_SecondClaimInsideTheBracketIsRefused(t *testing.T) {
	// No identity anywhere: a caller that claims twice cannot be told apart
	// from anyone else. The mechanism must not try — it reports "held",
	// exactly as it would for a different caller.
	g := newMemGit()
	withRunGit(t, g.run)
	withSleep(t, func(time.Duration) {}) // stay held — nobody releases mid-wait

	first := Claim(opts("c1"))
	if first.Outcome != Claimed {
		t.Fatalf("first claim: Outcome = %v, want Claimed", first.Outcome)
	}

	second := Claim(opts("c1"))
	if second.Outcome != Held {
		t.Fatalf("second claim: Outcome = %v, want Held — a double claim must be refused, not recognized", second.Outcome)
	}
	// Legible: the report names the caller's OWN change, making the mistake
	// evident to a reader even though the mechanism decided nothing from it.
	if !strings.Contains(second.Report, "c1") {
		t.Errorf("report %q does not name the caller's own change — self-blocking must be legible", second.Report)
	}
}

func TestClaim_HeldByAnother_BlocksAfterOneRetry(t *testing.T) {
	g := newMemGit()
	withRunGit(t, g.run)
	Claim(opts("holder-change")) // someone else's claim, still held

	slept := false
	withSleep(t, func(time.Duration) { slept = true })

	got := Claim(opts("waiter-change"))
	if got.Outcome != Held {
		t.Fatalf("Outcome = %v, want Held", got.Outcome)
	}
	if !strings.Contains(got.Report, "holder-change") {
		t.Errorf("report %q does not name the holder", got.Report)
	}
	if !slept {
		t.Error("the claim did not pause before its retry")
	}
	// The waiter's queue entry was written and, since the retry still found
	// it held, was never cleaned up — abandoned entries accumulate by design.
	if _, ok := g.refs[queueRef("waiter-change")]; !ok {
		t.Error("expected the queue entry to remain after a refused retry")
	}
}

func TestClaim_HeldByAnother_SucceedsOnRetry(t *testing.T) {
	g := newMemGit()
	withRunGit(t, g.run)
	first := Claim(opts("holder-change"))

	// The holder releases mid-wait — simulated as a side effect of the pause
	// seam, since the in-memory double has no real concurrency.
	withSleep(t, func(time.Duration) {
		Release(ReleaseOptions{Token: first.Token, Dir: "/repo"})
	})

	got := Claim(opts("waiter-change"))
	if got.Outcome != Claimed {
		t.Fatalf("Outcome = %v, want Claimed after the retry succeeds", got.Outcome)
	}
	if got.Token == "" {
		t.Error("expected a receipt for the successful retry")
	}
	if _, ok := g.refs[queueRef("waiter-change")]; ok {
		t.Error("the queue entry must be removed once the retry succeeds")
	}
}

func TestClaim_NotArbitratedPathsAreNeverSilent(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T)
	}{
		{
			"repository unresolvable",
			func(t *testing.T) {
				withRunGit(t, func(dir string, stdin []byte, args ...string) (string, error) {
					return "", errors.New("fatal: not a git repository")
				})
			},
		},
		{
			"holder record unreadable (blob missing)",
			func(t *testing.T) {
				g := newMemGit()
				g.refs[lockRef] = "ghost-sha"
				withRunGit(t, g.run)
			},
		},
		{
			"holder record malformed",
			func(t *testing.T) {
				g := newMemGit()
				g.blobs["bad1"] = "not: a: valid\nrecord"
				g.refs[lockRef] = "bad1"
				withRunGit(t, g.run)
			},
		},
		{
			"claim fails for a reason other than already-exists",
			func(t *testing.T) {
				withRunGit(t, func(dir string, stdin []byte, args ...string) (string, error) {
					switch args[0] {
					case "rev-parse":
						return "", errors.New("fatal: needed a single revision") // absent
					case "hash-object":
						return "", errors.New("fatal: unable to write sha1 file")
					}
					return "", nil
				})
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(t)
			got := Claim(opts("c1"))
			if got.Outcome != NotArbitrated {
				t.Fatalf("Outcome = %v, want NotArbitrated", got.Outcome)
			}
			if got.Reason == "" {
				t.Error("Reason is empty — an unarbitrated path must never be silent")
			}
			if got.Token != "" || got.Report != "" {
				t.Errorf("NotArbitrated result must carry no Token/Report, got %+v", got)
			}
		})
	}
}

func TestClaim_OutOfScopeIsIndistinguishableFromNothingToHide(t *testing.T) {
	// The claimed path's zero value is the discriminator: a Claimed result
	// with no Report/Reason is what "the turn is yours, say nothing" means.
	g := newMemGit()
	withRunGit(t, g.run)
	got := Claim(opts("c1"))
	if got != (ClaimResult{Outcome: Claimed, Token: got.Token}) {
		t.Errorf("got %+v, want only Outcome and Token set", got)
	}
}

func TestRelease_MatchingTokenReleases(t *testing.T) {
	g := newMemGit()
	withRunGit(t, g.run)
	claimed := Claim(opts("c1"))

	got, err := Release(ReleaseOptions{Token: claimed.Token, Dir: "/repo"})
	if err != nil {
		t.Fatalf("Release: %v", err)
	}
	if got.Outcome != Released {
		t.Fatalf("Outcome = %v, want Released", got.Outcome)
	}
	if _, held := g.refs[lockRef]; held {
		t.Error("the lock is still held after a matching release")
	}
}

func TestRelease_WrongTokenIsAMismatch(t *testing.T) {
	g := newMemGit()
	withRunGit(t, g.run)
	Claim(opts("c1"))

	got, err := Release(ReleaseOptions{Token: "not-the-real-token", Dir: "/repo"})
	if err != nil {
		t.Fatalf("Release: %v", err)
	}
	if got.Outcome != Mismatch {
		t.Fatalf("Outcome = %v, want Mismatch", got.Outcome)
	}
	if got.Message == "" {
		t.Error("Message is empty — a mismatch must be reported, not silent")
	}
	if _, held := g.refs[lockRef]; !held {
		t.Error("a mismatched release must delete nothing")
	}
}

func TestRelease_NoTokenIsAUsageError(t *testing.T) {
	_, err := Release(ReleaseOptions{Dir: "/repo"})
	if err == nil {
		t.Fatal("expected an error when no token is supplied — there is deliberately no fallback")
	}
}

// TestRelease_EveryFieldDiffersAndReleaseStillWorks and its sibling below
// pin the ownership invariant: contracts/turn-release-keyed-to-the-claim-token.md.
func TestRelease_EveryFieldDiffersAndReleaseStillWorks(t *testing.T) {
	g := newMemGit()
	withRunGit(t, g.run)
	claimed := Claim(ClaimOptions{Change: "a", Tree: "/a", Destination: "x", Moment: "change-integration", Dir: "/repo"})

	// The releasing call supplies nothing that matches the record's fields —
	// Release takes no such options at all. Only the token can ever matter.
	got, err := Release(ReleaseOptions{Token: claimed.Token, Dir: "/repo"})
	if err != nil || got.Outcome != Released {
		t.Fatalf("got %+v, err=%v — release must succeed on the right receipt regardless of record fields", got, err)
	}
}

func TestRelease_EveryFieldMatchesButWrongTokenStillFails(t *testing.T) {
	g := newMemGit()
	withRunGit(t, g.run)
	rec := Record{Change: "a", Tree: "/a", Destination: "x", Moment: "change-integration", Base: "deadbeef", At: nowUTC()}
	realSha, _ := g.run("/repo/.git", []byte(rec.String()), "hash-object", "-w", "--stdin")
	_, _ = g.run("/repo/.git", nil, "update-ref", lockRef, realSha, "")

	// A record byte-for-byte identical to the real one, but written as a
	// DIFFERENT blob — same fields, different sha.
	forgedSha, _ := g.run("/repo/.git", []byte(rec.String()+" "), "hash-object", "-w", "--stdin")

	got, err := Release(ReleaseOptions{Token: forgedSha, Dir: "/repo"})
	if err != nil {
		t.Fatalf("Release: %v", err)
	}
	if got.Outcome != Mismatch {
		t.Fatalf("Outcome = %v, want Mismatch — identical fields must not substitute for the receipt", got.Outcome)
	}
	if _, held := g.refs[lockRef]; !held {
		t.Error("the real turn must remain held when released with the wrong token")
	}
}

func TestStatus_Free(t *testing.T) {
	withRunGit(t, newMemGit().run)
	got := Status("/repo")
	if got.Held {
		t.Error("expected Held = false")
	}
	if got.Report == "" {
		t.Error("expected a non-empty report")
	}
}

func TestStatus_Held(t *testing.T) {
	g := newMemGit()
	withRunGit(t, g.run)
	Claim(opts("c1"))

	got := Status("/repo")
	if !got.Held {
		t.Error("expected Held = true")
	}
	if !strings.Contains(got.Report, "c1") {
		t.Errorf("report %q does not name the holder", got.Report)
	}
}

func TestStatus_NeverFails(t *testing.T) {
	withRunGit(t, func(dir string, stdin []byte, args ...string) (string, error) {
		return "", errors.New("fatal: not a git repository")
	})
	got := Status("/repo") // no panic, no error return possible at all — Status has none
	if got.Report == "" {
		t.Error("even an unresolvable repository must produce a report, not nothing")
	}
}

func TestStatus_NeverMutatesState(t *testing.T) {
	g := newMemGit()
	withRunGit(t, g.run)
	Claim(opts("c1"))
	before := len(g.refs)

	Status("/repo")

	if len(g.refs) != before {
		t.Errorf("status wrote or deleted a ref: before=%d after=%d refs=%v", before, len(g.refs), g.refs)
	}
}

func TestWorktreesSkillDocumentation(t *testing.T) {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		t.Skip("not inside a git repository — cannot locate the skill file")
	}
	root := strings.TrimSpace(string(out))
	path := filepath.Join(root, "payload", "domains", "development", "skills", "matecito-ai", "worktrees", "SKILL.md")
	body, err := os.ReadFile(path)
	if err != nil {
		t.Skip("worktrees skill not found at " + path)
	}
	text := string(body)

	if !strings.Contains(text, "One claim, one release, per unit of work") {
		t.Error("the bracket invariant is not stated in the skill")
	}
	if !strings.Contains(text, "matecito-ai turn claim") || !strings.Contains(text, "matecito-ai turn release") {
		t.Error("the skill does not name both commands")
	}
	if strings.Contains(text, "PreToolUse") || strings.Contains(text, "merge-turn-guard") {
		t.Error("the skill still mentions the deleted PreToolUse guard")
	}
	if !strings.Contains(text, "## When matecito-ai is not installed") {
		t.Error("the retained hand-run procedure's gated section is missing or misnamed")
	}
}

func TestIntegration_RefEffectsInAScratchRepository(t *testing.T) {
	dir := t.TempDir()
	if out, err := exec.Command("git", "-C", dir, "init").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}

	first := Claim(ClaimOptions{Change: "merge-lock-hook", Tree: dir, Destination: "main", Moment: "change-integration", Dir: dir})
	if first.Outcome != Claimed || first.Token == "" {
		t.Fatalf("first claim: %+v", first)
	}

	// A second claim finds it held, queues, retries, and gives up (nobody
	// releases mid-wait in this test) — the queue entry then stays.
	withSleep(t, func(time.Duration) {})
	second := Claim(ClaimOptions{Change: "other-change", Tree: dir, Destination: "main", Moment: "batch-consolidation", Dir: dir})
	if second.Outcome != Held {
		t.Fatalf("second claim: %+v, want Held", second)
	}

	out, err := exec.Command("git", "-C", dir, "for-each-ref", queueRefPrefix).Output()
	if err != nil {
		t.Fatalf("for-each-ref: %v", err)
	}
	if strings.TrimSpace(string(out)) == "" {
		t.Error("expected the abandoned queue entry to remain — nothing sweeps it")
	}

	// Release with the correct token succeeds.
	rel, err := Release(ReleaseOptions{Token: first.Token, Dir: dir})
	if err != nil {
		t.Fatalf("Release: %v", err)
	}
	if rel.Outcome != Released {
		t.Fatalf("Release outcome = %v, want Released", rel.Outcome)
	}

	// Releasing again (already gone) reports the mismatch instead of forcing.
	rel2, err := Release(ReleaseOptions{Token: first.Token, Dir: dir})
	if err != nil {
		t.Fatalf("Release (again): %v", err)
	}
	if rel2.Outcome != Mismatch {
		t.Fatalf("second release outcome = %v, want Mismatch", rel2.Outcome)
	}

	status := Status(dir)
	if status.Held {
		t.Error("status should report the turn free after release")
	}

	out, err = exec.Command("git", "-C", dir, "status", "--porcelain").Output()
	if err != nil {
		t.Fatalf("git status: %v", err)
	}
	if strings.TrimSpace(string(out)) != "" {
		t.Errorf("working tree is not clean: %q", out)
	}
}
