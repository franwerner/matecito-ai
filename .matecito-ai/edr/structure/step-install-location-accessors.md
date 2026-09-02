# EDR — One exported install-location accessor per step

- **Status:** Accepted
- **Date:** 2026-09-02

## Contexto
The install surface (sync.Detect) and the verify surface (internal/checks/engram, /codegraph, /proofshot) both need to know where a step's binary lands, so both can probe it there instead of asking the running session's PATH. Before this decision neither surface exposed that location as a reusable value: sync.go resolved it inline for qmd (via the already-ratified resolveUserNpmBinDir) and the verify checks did not resolve it at all.

## Decisión
internal/setup/install exports one accessor per step — EngramBinaryPath, CodegraphBinaryPath, ProofshotBinaryPath, each func() (string, error) — matching exactly the resolve parameter shape check.ProbeAt takes, so every one of the six call sites (three in sync.Detect, three in the verify checks) passes it by name, with no closure and no filepath.Join of its own. EngramBinaryPath returns releasedl.DefaultBinaryPath(releasedl.EngramRepo); CodegraphBinaryPath and ProofshotBinaryPath each join UserNpmBinDir() with their own binary name.

## Reglas verificables
- **[auto]** each *BinaryPath accessor returns the exact path its own step's Run already writes to — internal/setup/sync/binary_detection_test.go's real-resolver subtests for codegraph and proofshot, run through the real npm-backed resolver.
- **[manual]** install.EngramBinaryPath() matches the actual location `engram` is installed at on a real machine (verified by running the installed binary during this change's own apply).

## Alternativas consideradas
(a) Each of the six call sites joins dir + name itself — six copies of "where codegraph/proofshot lives", and requirement 4 would then rest on six literals continuing to agree with each other and with InstallCodegraph/InstallProofshot's own destination. (b) A new shared package for install locations — churn a bug fix does not justify; already rejected in the design's first pass for UserNpmBinDir on the same grounds. (c) The verify checks calling releasedl.DefaultBinaryPath directly for engram instead of going through install.EngramBinaryPath — two owners of "engram's place", and it puts release-download plumbing (releasedl) inside a reporting package that has no other reason to import it. (d) Folding the version args ("version" vs "--version") into the same accessor as well — that is the shared binary registry the user already discarded; the accessor answers only "where", never "how to probe it".

## Consecuencias
internal/setup/install grows three small exported functions, each named after the step it serves, so "where does X live" has exactly one answer in the codebase. internal/checks/* gains a new import of internal/setup/install (verified acyclic: install imports check, hook, manifest, mcp, platform, setup/deploy, setup/releasedl, setup/settings — none of them imports internal/checks), and the verify binary already links install through internal/cli/verify.go. Neither package has an init(), so no import-time side effect is added by this new edge.
