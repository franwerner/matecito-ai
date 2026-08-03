// Command gen-legacy-orphans reconstructs
// internal/setup/deploy/legacyorphans.json: the embedded catalog of legacy
// deploy destinations — targets a past version of the payload wrote that
// today's deploy.Plan() no longer produces at all (as opposed to StatusRemoved
// entries derived from the registry, which cover destinations a *recent* run
// tracked). Each entry carries every content hash that destination is known to
// have had historically, since a host being migrated to the new registry can
// be on any intermediate version.
//
// Run manually from the repo root:
//
//	go run ./scripts/gen-legacy-orphans
//
// The output is committed, not generated at build time — there is no runtime
// dependency on git, and the embedded catalog ships inside the binary.
//
// Re-run it whenever a release removes or renames payload files. A host that
// upgrades late deployed those releases with its pre-registry binary, so the
// orphans they left are only ever cleaned by the catalog embedded in the
// version it finally jumps to. The release workflow enforces this.
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/franwerner/matecito-ai/internal/setup/deploy"
)

// dirComponents mirrors deploy.go's domainComponents that use ModeDir: the
// source subtree copies through unchanged below its component name.
var dirComponents = map[string]bool{"agents": true, "references": true, "scripts": true}

// legacyEntry mirrors deploy.LegacyEntry (kept local — the generator is a
// separate binary and must not depend on deploy's unexported catalog loader).
type legacyEntry struct {
	Root   string   `json:"root"`
	Path   string   `json:"path"`
	Hashes []string `json:"hashes"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "gen-legacy-orphans:", err)
		os.Exit(1)
	}
}

func run() error {
	repoRoot, err := repoRoot()
	if err != nil {
		return fmt.Errorf("resolving repo root: %w", err)
	}

	everPaths, err := everPayloadPaths(repoRoot)
	if err != nil {
		return fmt.Errorf("listing historical payload/ paths: %w", err)
	}
	nowPaths, err := nowPayloadPaths(repoRoot)
	if err != nil {
		return fmt.Errorf("listing current payload/ paths: %w", err)
	}

	// Only paths that vanished from the payload source tree are candidates —
	// a path still present today (even if the current deploy code no longer
	// deploys it, e.g. a .gitkeep under a component that now filters them) is
	// not a legacy orphan by this catalog's contract: it was never removed
	// from the payload, so there is nothing to reconstruct from git history
	// for it. This mirrors the oracle reconstruction validated in Engram #1109.
	gone := make([]string, 0, len(everPaths))
	for p := range everPaths {
		if !nowPaths[p] {
			gone = append(gone, p)
		}
	}
	sort.Strings(gone)

	// Group every gone source path by its mapped target, replicating
	// ModeFile/ModeDir/ModeGrouped for both payload layout eras.
	byTarget := make(map[string][]string)
	for _, rel := range gone {
		t := targetFor(rel)
		if t == "" {
			continue
		}
		byTarget[t] = append(byTarget[t], rel)
	}

	currentTargets, err := currentPlanTargets(repoRoot)
	if err != nil {
		return fmt.Errorf("resolving current plan targets: %w", err)
	}

	var targets []string
	for t := range byTarget {
		if !currentTargets[t] {
			targets = append(targets, t)
		}
	}
	sort.Strings(targets)

	entries := make([]legacyEntry, 0, len(targets)+1)
	for _, t := range targets {
		hashes, err := historicalHashes(repoRoot, byTarget[t])
		if err != nil {
			return fmt.Errorf("hashing history for %q: %w", t, err)
		}
		entries = append(entries, legacyEntry{Root: "claude", Path: t, Hashes: hashes})
	}

	// The abandoned ~/.matecito-ai/deployed-manifest.json: no code writes it
	// anymore and it predates any commit, so it cannot be reconstructed from git.
	entries = append(entries, abandonedManifestEntry())

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Root != entries[j].Root {
			return entries[i].Root < entries[j].Root
		}
		return entries[i].Path < entries[j].Path
	})

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal catalog: %w", err)
	}
	data = append(data, '\n')

	out := filepath.Join(repoRoot, "internal", "setup", "deploy", "legacyorphans.json")
	if err := os.WriteFile(out, data, 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", out, err)
	}

	fmt.Printf("gen-legacy-orphans: wrote %d entries to %s\n", len(entries), out)
	return nil
}

// targetFor mirrors deploy's buildMappings/expandMapping decisions but works
// off a historical path STRING relative to payload/ (the file may no longer
// exist on disk), for both eras of the payload layout: domains/<id>/... +
// shared/... (current), and a flat agents/skills/references/scripts root
// (legacy, pre-domains). Returns "" for a path that never deployed under any
// era (domain README/GUIDE, hook manifests reconciled into settings.json, a
// file loose at a skill-group root, etc.).
func targetFor(payloadRel string) string {
	parts := strings.Split(payloadRel, "/")
	if len(parts) == 0 {
		return ""
	}

	// Era 2: domains/<id>/<component>/... and shared/<component>/...
	if parts[0] == "domains" && len(parts) >= 3 {
		return componentTarget(parts[2:])
	}
	if parts[0] == "shared" && len(parts) >= 2 {
		return componentTarget(parts[1:])
	}

	// Era 1: component at the payload root.
	if dirComponents[parts[0]] || parts[0] == "skills" {
		return componentTarget(parts)
	}

	// The payload's own CLAUDE.md composes into matecito-ai.md in memory — it
	// is never copied byte-for-byte, and that target is always present in the
	// current plan, so it is never itself flagged as an orphan. Mapped anyway
	// for completeness of the byTarget grouping.
	if payloadRel == "CLAUDE.md" {
		return "matecito-ai.md"
	}
	return ""
}

func componentTarget(rest []string) string {
	if len(rest) == 0 {
		return ""
	}
	comp := rest[0]
	switch {
	case dirComponents[comp]:
		if len(rest) < 2 {
			return ""
		}
		return comp + "/" + strings.Join(rest[1:], "/")
	case comp == "skills":
		// ModeGrouped drops the group level (vendor/domain folder); a file
		// loose at the group level (no skill folder under it) never deployed.
		if len(rest) < 3 {
			return ""
		}
		return "skills/" + strings.Join(rest[2:], "/")
	}
	return ""
}

// everPayloadPaths lists every path under payload/ that has ever existed in
// any commit reachable from any ref, relative to payload/.
func everPayloadPaths(repoRoot string) (map[string]bool, error) {
	out, err := gitOutput(repoRoot, "log", "--all", "--no-renames", "--name-only", "--pretty=format:", "--", "payload/")
	if err != nil {
		return nil, err
	}
	paths := make(map[string]bool)
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "payload/") {
			continue
		}
		paths[strings.TrimPrefix(line, "payload/")] = true
	}
	return paths, nil
}

// nowPayloadPaths walks the current payload/ directory on disk.
func nowPayloadPaths(repoRoot string) (map[string]bool, error) {
	root := filepath.Join(repoRoot, "payload")
	paths := make(map[string]bool)
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		paths[filepath.ToSlash(rel)] = true
		return nil
	})
	return paths, err
}

// currentPlanTargets calls the real deploy.Plan() against the local payload/
// tree and returns the set of relative-to-claudeHome targets it produces
// (excluding StatusRemoved — those are not "current", they are what Plan
// itself would already flag for deletion). A throwaway stateDir guarantees no
// registry is found, so Plan never mixes in this host's real removed-ops.
func currentPlanTargets(repoRoot string) (map[string]bool, error) {
	payloadFS := os.DirFS(filepath.Join(repoRoot, "payload"))
	const claudeHome = "/gen-legacy-orphans/claudehome"
	stateDir, err := os.MkdirTemp("", "gen-legacy-orphans-state-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(stateDir)

	ops, err := deploy.Plan(payloadFS, claudeHome, stateDir, nil)
	if err != nil {
		return nil, err
	}
	targets := make(map[string]bool, len(ops))
	for _, op := range ops {
		if op.Status == deploy.StatusRemoved {
			continue
		}
		rel, err := filepath.Rel(claudeHome, op.Target)
		if err != nil {
			return nil, err
		}
		targets[filepath.ToSlash(rel)] = true
	}
	return targets, nil
}

// historicalHashes returns every distinct content hash any of sourcePaths
// (relative to payload/) is known to have had at any commit in its history —
// the set a host migrating from any intermediate version could be showing.
func historicalHashes(repoRoot string, sourcePaths []string) ([]string, error) {
	seen := make(map[string]bool)
	for _, rel := range sourcePaths {
		full := "payload/" + rel
		shas, err := commitsTouching(repoRoot, full)
		if err != nil {
			return nil, err
		}
		for _, sha := range shas {
			content, err := gitShow(repoRoot, sha, full)
			if err != nil {
				// Deleted-at-this-commit or otherwise unreadable at this
				// revision — not a content version to track.
				continue
			}
			sum := sha256.Sum256(content)
			seen[hex.EncodeToString(sum[:])] = true
		}
	}
	hashes := make([]string, 0, len(seen))
	for h := range seen {
		hashes = append(hashes, h)
	}
	sort.Strings(hashes)
	return hashes, nil
}

func commitsTouching(repoRoot, path string) ([]string, error) {
	out, err := gitOutput(repoRoot, "log", "--all", "--no-renames", "--format=%H", "--", path)
	if err != nil {
		return nil, err
	}
	var shas []string
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			shas = append(shas, line)
		}
	}
	return shas, nil
}

func gitShow(repoRoot, sha, path string) ([]byte, error) {
	cmd := exec.Command("git", "show", sha+":"+path)
	cmd.Dir = repoRoot
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%s: %w", errOut.String(), err)
	}
	return out.Bytes(), nil
}

func gitOutput(repoRoot string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = repoRoot
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, errOut.String())
	}
	return out.String(), nil
}

func repoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// abandonedManifestHash is the sha256 of ~/.matecito-ai/deployed-manifest.json,
// pinned as a constant: no commit ever produced that file, so it cannot be
// reconstructed from git. Reading it from the running machine's home directory
// made this generator emit a different catalog depending on where it ran — the
// entry appeared on a host that still had the file and vanished on CI, which
// never has it. That is unverifiable by definition, so the value is fixed here.
const abandonedManifestHash = "81f1faa8c205985ebe8aefad26ba93ea57b589921b298b3056b6a868b3b7435f"

func abandonedManifestEntry() legacyEntry {
	return legacyEntry{
		Root:   "state",
		Path:   "deployed-manifest.json",
		Hashes: []string{abandonedManifestHash},
	}
}
