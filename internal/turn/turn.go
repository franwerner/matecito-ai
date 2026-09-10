// Package turn is the mechanics behind `matecito-ai turn`: one git ref
// (refs/matecito-ai/merge-lock) that serializes a write to a shared branch
// between two concurrent sessions. There is no interception and no
// identity — a caller claims the turn, does one unit of work, and releases
// it with the receipt the claim returned. Ownership is the receipt alone;
// nothing here ever compares a record's fields against anything.
package turn

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// lockRef is the single git ref every worktree of the repository shares to
// represent the turn (contracts/merge-turn-holder-token.md).
const lockRef = "refs/matecito-ai/merge-lock"

// queueRefPrefix is the namespace of pending-request refs
// (structure/merge-queue-lives-in-refs.md).
const queueRefPrefix = "refs/matecito-ai/merge-queue/"

func queueRef(change string) string {
	return queueRefPrefix + change
}

// contentionPause is how long a claim that found the turn held waits before
// its single retry (structure/guard-queues-and-retries-agent-asks.md).
const contentionPause = 300 * time.Millisecond

// runGit runs `git -C dir <args...>`, optionally piping stdin, and returns
// trimmed stdout. On failure the returned error's message is git's stderr
// (falling back to the Go error) so callers can classify failures (e.g.
// "already exists") without re-running the command.
var runGit = func(dir string, stdin []byte, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	if stdin != nil {
		cmd.Stdin = bytes.NewReader(stdin)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	trimmed := strings.TrimRight(stdout.String(), "\n")
	if err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return trimmed, errors.New(msg)
	}
	return trimmed, nil
}

var nowUTC = func() string {
	return time.Now().UTC().Format(time.RFC3339)
}

var sleep = func(d time.Duration) {
	time.Sleep(d)
}

var getwd = os.Getwd

// Record is the holder/queue record shape pinned by
// contracts/merge-turn-holder-token.md: six fields, every one of them
// descriptive. None is ever compared against anything — see
// contracts/turn-release-keyed-to-the-claim-token.md.
type Record struct {
	Change      string
	Tree        string
	Destination string
	Moment      string
	Base        string
	At          string
}

func (r Record) String() string {
	return fmt.Sprintf(
		"change: %s\ntree: %s\ndestination: %s\nmoment: %s\nbase: %s\nat: %s\n",
		r.Change, r.Tree, r.Destination, r.Moment, r.Base, r.At,
	)
}

var recordFields = []string{"change", "tree", "destination", "moment", "base", "at"}

func parseRecord(body string) (Record, error) {
	fields := map[string]string{}
	for _, line := range strings.Split(strings.TrimRight(body, "\n"), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) != 2 {
			return Record{}, fmt.Errorf("malformed record line: %q", line)
		}
		fields[parts[0]] = parts[1]
	}
	for _, f := range recordFields {
		if _, ok := fields[f]; !ok {
			return Record{}, fmt.Errorf("record missing field %q", f)
		}
	}
	return Record{
		Change:      fields["change"],
		Tree:        fields["tree"],
		Destination: fields["destination"],
		Moment:      fields["moment"],
		Base:        fields["base"],
		At:          fields["at"],
	}, nil
}

func resolveCommonDir(dir string) (string, error) {
	out, err := runGit(dir, nil, "rev-parse", "--git-common-dir")
	if err != nil {
		return "", err
	}
	if out == "" {
		return "", errors.New("empty --git-common-dir output")
	}
	if filepath.IsAbs(out) {
		return out, nil
	}
	return filepath.Join(dir, out), nil
}

func bestEffort(out string, err error) string {
	if err != nil || out == "" {
		return "unknown"
	}
	return out
}

func resolveDir(dir string) (string, error) {
	if dir != "" {
		return dir, nil
	}
	return getwd()
}

func readLock(commonDir string) (rec Record, sha string, held bool, err error) {
	out, gerr := runGit(commonDir, nil, "rev-parse", "--verify", "--quiet", lockRef)
	if gerr != nil || out == "" {
		return Record{}, "", false, nil
	}
	sha = out
	body, cerr := runGit(commonDir, nil, "cat-file", "-p", sha)
	if cerr != nil {
		return Record{}, sha, true, fmt.Errorf("holder record %s is unreadable: %w", sha, cerr)
	}
	rec, perr := parseRecord(body)
	if perr != nil {
		return Record{}, sha, true, fmt.Errorf("holder record %s is malformed: %w", sha, perr)
	}
	return rec, sha, true, nil
}

func writeLock(commonDir string, rec Record) (string, error) {
	blobSha, err := runGit(commonDir, []byte(rec.String()), "hash-object", "-w", "--stdin")
	if err != nil {
		return "", err
	}
	if _, err := runGit(commonDir, nil, "update-ref", lockRef, blobSha, ""); err != nil {
		return "", err
	}
	return blobSha, nil
}

func looksLikeAlreadyExists(err error) bool {
	return err != nil && strings.Contains(strings.ToLower(err.Error()), "already exists")
}

// Outcome is the three-way result `claim` MUST produce
// (see spec: "The three outcomes are told apart").
type Outcome int

const (
	// Claimed: the turn is now the caller's. Token is the receipt.
	Claimed Outcome = iota
	// Held: another claim (or this same one, double-claiming) holds the
	// turn. Report names it.
	Held
	// NotArbitrated: the mechanism could not decide at all. The work MUST
	// be allowed to proceed; Reason says why the check could not run.
	NotArbitrated
)

// ClaimOptions is everything a caller supplies to Claim. Every field is
// descriptive — none of it is ever compared against anything; see
// contracts/turn-release-keyed-to-the-claim-token.md.
type ClaimOptions struct {
	Change      string
	Tree        string
	Destination string
	Moment      string
	// Dir is where the repository is resolved from. Empty uses the
	// process's own working directory.
	Dir string
}

// ClaimResult is what Claim returns. Exactly one of Token, Report, Reason
// is meaningful, selected by Outcome.
type ClaimResult struct {
	Outcome Outcome
	Token   string // set on Claimed
	Report  string // set on Held
	Reason  string // set on NotArbitrated
}

// Claim is the whole PreToolUse-shaped decision, now an ordinary function
// call: free → write the record, CAS-create the ref, return its sha as the
// receipt. Held → queue, pause, retry once, then either claim (as above)
// or report the holder. Unable to decide at all → NotArbitrated, so the
// caller's work proceeds and the failure is never silent.
func Claim(opts ClaimOptions) ClaimResult {
	dir, err := resolveDir(opts.Dir)
	if err != nil {
		return ClaimResult{Outcome: NotArbitrated, Reason: "could not resolve the working directory: " + err.Error()}
	}
	commonDir, err := resolveCommonDir(dir)
	if err != nil {
		return ClaimResult{Outcome: NotArbitrated, Reason: "could not resolve the repository: " + err.Error()}
	}

	rec := Record{
		Change:      opts.Change,
		Tree:        opts.Tree,
		Destination: opts.Destination,
		Moment:      opts.Moment,
		Base:        bestEffort(runGit(dir, nil, "rev-parse", "HEAD")),
		At:          nowUTC(),
	}

	return evaluateClaim(commonDir, rec)
}

func evaluateClaim(commonDir string, rec Record) ClaimResult {
	holder, _, held, err := readLock(commonDir)
	if err != nil {
		return ClaimResult{Outcome: NotArbitrated, Reason: err.Error()}
	}
	if !held {
		sha, claimErr := writeLock(commonDir, rec)
		if claimErr == nil {
			return ClaimResult{Outcome: Claimed, Token: sha}
		}
		if !looksLikeAlreadyExists(claimErr) {
			return ClaimResult{Outcome: NotArbitrated, Reason: "could not claim the turn: " + claimErr.Error()}
		}
		// Raced with another claimant since the read above — re-read and fall through to contention.
		holder, _, held, err = readLock(commonDir)
		if err != nil {
			return ClaimResult{Outcome: NotArbitrated, Reason: err.Error()}
		}
		if !held {
			return ClaimResult{Outcome: NotArbitrated, Reason: "the turn's lock disappeared during the claim attempt"}
		}
	}
	return contend(commonDir, rec, holder)
}

// contend is the queue-and-retry procedure
// (structure/guard-queues-and-retries-agent-asks.md), performed within this
// one call: it queues the pending request, pauses, retries the claim
// exactly once, and either succeeds (deleting its own queue entry) or
// reports the holder. A failed queue write degrades the report, never the
// verdict.
func contend(commonDir string, rec, held Record) ClaimResult {
	qSha, qErr := writeQueueEntry(commonDir, rec)
	sleep(contentionPause)

	sha, claimErr := writeLock(commonDir, rec)
	if claimErr == nil {
		if qErr == nil {
			deleteQueueEntry(commonDir, queueRef(rec.Change), qSha)
		}
		return ClaimResult{Outcome: Claimed, Token: sha}
	}

	current, _, currentHeld, readErr := readLock(commonDir)
	if readErr != nil || !currentHeld {
		current = held
	}
	report := buildReport(current, commonDir)
	if qErr != nil {
		report += " Your wait could not be recorded in the queue: " + qErr.Error() + "."
	}
	return ClaimResult{Outcome: Held, Report: report}
}

func writeQueueEntry(commonDir string, rec Record) (string, error) {
	blobSha, err := runGit(commonDir, []byte(rec.String()), "hash-object", "-w", "--stdin")
	if err != nil {
		return "", err
	}
	if _, err := runGit(commonDir, nil, "update-ref", queueRef(rec.Change), blobSha, ""); err != nil {
		return "", err
	}
	return blobSha, nil
}

func deleteQueueEntry(commonDir, ref, sha string) {
	_, _ = runGit(commonDir, nil, "update-ref", "-d", ref, sha)
}

func buildReport(holder Record, commonDir string) string {
	age := "an unknown time ago"
	if t, err := time.Parse(time.RFC3339, holder.At); err == nil {
		age = time.Since(t).Round(time.Second).String() + " ago"
	}
	return fmt.Sprintf(
		"turn held by change %q (tree %s, destination %s, moment %s), claimed %s. %s",
		holder.Change, holder.Tree, holder.Destination, holder.Moment, age, describeQueue(commonDir),
	)
}

func describeQueue(commonDir string) string {
	out, err := runGit(commonDir, nil, "for-each-ref", "--format=%(refname)", queueRefPrefix)
	if err != nil || out == "" {
		return "No other session is queued."
	}
	n := len(strings.Split(out, "\n"))
	if n == 1 {
		return "1 session is queued behind it."
	}
	return fmt.Sprintf("%d sessions are queued behind it.", n)
}

// ReleaseOutcome is release's two-way result.
type ReleaseOutcome int

const (
	Released ReleaseOutcome = iota
	Mismatch
)

// ReleaseOptions is what Release needs. Token is mandatory — there is
// deliberately no fallback that identifies the caller any other way (see
// contracts/turn-release-keyed-to-the-claim-token.md).
type ReleaseOptions struct {
	Token string
	Dir   string
}

type ReleaseResult struct {
	Outcome ReleaseOutcome
	Message string // set on Mismatch
}

// Release deletes the lock ref conditioned on it still carrying exactly
// opts.Token — the receipt Claim returned. Ownership is the token alone:
// no field of the holder record is ever read or compared here.
func Release(opts ReleaseOptions) (ReleaseResult, error) {
	if opts.Token == "" {
		return ReleaseResult{}, errors.New("release requires a token")
	}
	dir, err := resolveDir(opts.Dir)
	if err != nil {
		return ReleaseResult{}, fmt.Errorf("could not resolve the working directory: %w", err)
	}
	commonDir, err := resolveCommonDir(dir)
	if err != nil {
		return ReleaseResult{}, fmt.Errorf("could not resolve the repository: %w", err)
	}
	if _, err := runGit(commonDir, nil, "update-ref", "-d", lockRef, opts.Token); err != nil {
		return ReleaseResult{
			Outcome: Mismatch,
			Message: "the turn no longer carries this receipt (already released, or moved): " + err.Error(),
		}, nil
	}
	return ReleaseResult{Outcome: Released}, nil
}

// StatusResult is status's read-only report. Status MUST always succeed —
// there is no failure outcome, only what it found.
type StatusResult struct {
	Held   bool
	Report string
}

// Status reads the turn without claiming or releasing anything. It always
// returns a result; when the repository itself cannot be resolved, that
// fact becomes part of the report instead of a returned error, so a caller
// that only wants to look is never blocked by the same checks Claim runs.
func Status(dir string) StatusResult {
	resolvedDir, err := resolveDir(dir)
	if err != nil {
		return StatusResult{Report: "could not resolve the working directory: " + err.Error()}
	}
	commonDir, err := resolveCommonDir(resolvedDir)
	if err != nil {
		return StatusResult{Report: "could not resolve the repository: " + err.Error()}
	}
	holder, _, held, err := readLock(commonDir)
	if err != nil {
		return StatusResult{Report: err.Error()}
	}
	if !held {
		return StatusResult{Held: false, Report: "turn free. " + describeQueue(commonDir)}
	}
	return StatusResult{Held: true, Report: buildReport(holder, commonDir)}
}
