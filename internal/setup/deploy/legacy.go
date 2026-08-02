package deploy

import (
	_ "embed"
	"encoding/json"
	"path/filepath"
	"sync"
)

//go:embed legacyorphans.json
var legacyOrphansJSON []byte

// LegacyEntry is one destination a past version of the payload deployed that
// the current Plan() no longer produces at all — a host may be migrating from
// any intermediate version, so Hashes carries every content version that
// destination is known to have had historically, not just one. Root picks
// which real root the entry lives under: "claude" (~/.claude, the normal
// deploy destination) or "state" (~/.matecito-ai, for the one entry that is
// the abandoned deployed-manifest.json itself). Hashes empty means "never
// reconstructable from git" (e.g. a dirty working tree at deploy time) — such
// an entry can never match on hash, so it is never deleted, only ever
// reported as preserved when something happens to sit at that path.
type LegacyEntry struct {
	Root   string   `json:"root"`
	Path   string   `json:"path"`
	Hashes []string `json:"hashes"`
}

var (
	legacyCatalogOnce sync.Once
	legacyCatalogData []LegacyEntry
)

// legacyCatalog parses the embedded legacy orphan catalog once per process.
// The catalog is generated offline (scripts/gen-legacy-orphans) and committed;
// nothing here ever writes it.
func legacyCatalog() []LegacyEntry {
	legacyCatalogOnce.Do(func() {
		var entries []LegacyEntry
		// A malformed embed would be a build-time defect (the file ships with the
		// binary and is never user input) — decode failures are silently treated
		// as an empty catalog rather than panicking a live deploy run.
		_ = json.Unmarshal(legacyOrphansJSON, &entries)
		legacyCatalogData = entries
	})
	return legacyCatalogData
}

// legacyRemovedOps builds StatusRemoved ops from the embedded legacy catalog —
// used only during migration (no registry found yet), since a fresh registry
// diff would otherwise have nothing to compare against.
func legacyRemovedOps(currentTargets map[string]bool, claudeHome, stateDir string) []FileOp {
	return legacyRemovedOpsFrom(legacyCatalog(), currentTargets, claudeHome, stateDir)
}

// legacyRemovedOpsFrom is legacyRemovedOps with an injectable catalog, so unit
// tests can exercise the Root routing ("claude" vs "state") and the
// already-current skip without depending on the real embedded data. entries is
// the set of absolute targets the current Plan() actually produces; an entry
// that coincides with one (e.g. a reactivated domain reproducing an old path)
// is skipped defensively.
//
// An entry whose target is absent from this host (the overwhelming common
// case: the 103-entry catalog covers every version's orphans, and any given
// host only ever had a handful of them) is skipped too — Apply already treats
// an absent target as a silent no-op, but emitting the op anyway made the
// migration preview list — and count — paths nothing will actually touch.
func legacyRemovedOpsFrom(catalog []LegacyEntry, currentTargets map[string]bool, claudeHome, stateDir string) []FileOp {
	var ops []FileOp
	for _, e := range catalog {
		root := claudeHome
		if e.Root == "state" {
			root = stateDir
		}
		target := filepath.Join(root, filepath.FromSlash(e.Path))
		if currentTargets[target] {
			continue
		}
		if !targetExists(target) {
			continue
		}
		ops = append(ops, FileOp{
			Target:       target,
			Status:       StatusRemoved,
			LegacyHashes: e.Hashes,
		})
	}
	return ops
}
