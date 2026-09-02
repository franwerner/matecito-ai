# EDR — The three verify checks take their resolver through a package-level seam

- **Status:** Accepted
- **Date:** 2026-09-02

## Contexto
The three verify-surface checks (internal/checks/engram, internal/checks/codegraph, internal/checks/proofshot) decided presence by handing a bare binary name to check.RunVersion, resolved through the running session's PATH, matching the same wrong pattern install.Detect had before this change's Phase 1-2 (structure/binary-detection-at-canonical-path). Fixing the install surface without also fixing these three checks would leave the two commands able to contradict each other about the same machine, which is exactly what spec process/binary-install-step-idempotency's requirement "Instalar y reportar el estado no se contradicen" forbids.

## Decisión
Each of the three verify-surface check packages (internal/checks/engram, internal/checks/codegraph, internal/checks/proofshot) declares its own package-level `var resolveBinary = install.<X>BinaryPath`, and its `detectBinary`/`detectCLI` calls `check.ProbeAt(name, resolveBinary, args, required, fixHint)` instead of `check.RunVersion` with a bare name. Branch tests live in a `*_internal_test.go` in the same package and override the seam directly, following internal/checks/debugger/debugger.go's `var find = mcp.Find` convention (overridden in debugger_internal_test.go, with an optional black-box debugger_test.go asserting shape invariants).

## Reglas verificables
- **[auto]** each check package exposes its own `var resolveBinary = install.<X>BinaryPath` seam; overriding it to a closure that errors makes detectBinary/detectCLI return StatusMissing — internal/checks/engram/engram_internal_test.go, internal/checks/codegraph/codegraph_internal_test.go and internal/checks/proofshot/proofshot_internal_test.go, subtest "resolver-fails".
- **[auto]** with the binary present only at its own install.<X>BinaryPath() (absent from the test's PATH), detectBinary/detectCLI returns StatusOK — same three files, subtest "canonical-present" and "path-not-yet-on-session-PATH".
- **[manual]** case names in the three *_internal_test.go files match verbatim the twin cases in internal/setup/sync/binary_detection_test.go — canonical-present, canonical-absent, canonical-but-broken, path-not-yet-on-session-PATH, resolver-fails — so a divergence between the install and verify surfaces fails on exactly one side. No automated cross-file assertion enforces the name match; it is sustained by convention alone.

## Alternativas consideradas
(a) drive the real resolver with a fake npm on an isolated PATH plus a doctored HOME for every case — works for the codegraph/proofshot pair but engram's path comes from os.UserHomeDir(), and the resolver-error branch cannot be produced that way for any of the three; (b) export detectBinary/detectCLI for tests — widens the package's public API for testing only, a convention this repo does not follow elsewhere; (c) ship the three checks with no unit tests, manual verification only — contradicts structure/installer-step-ships-with-tests.md's Consecuencias.

## Consecuencias
The three check packages gain one new package-level var each and one new import (internal/setup/install), mirroring debugger's existing shape exactly. All five spec scenarios of process/binary-install-step-idempotency become reachable in tests, including the resolver-error branch that a real-npm/real-HOME setup cannot reach for engram. Any future check added to this family follows the same seam convention by precedent.
