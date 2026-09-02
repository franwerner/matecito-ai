# EDR — Each install step keeps its own detection block, sharing one probe routine

- **Status:** Accepted
- **Date:** 2026-09-02

## Contexto
sync.Detect's engram, codegraph and proofshot blocks (sync.go:234-282) each decided presence by handing a bare name to check.RunVersion, resolved through the running session's PATH — so a binary a prior run installed and persisted to a shell rc file, but that this session has not yet sourced, was reported missing and reinstalled.

## Decisión
Each of the three steps keeps its own detection block (sync.go's engram/codegraph/proofshot sections in Detect) and its own way of locating its binary; they share one small routine, probeInstalledBinary(name string, resolve func() (string, error), versionArgs []string) (present bool, version string), that checks a binary at a given absolute location — so that rule is written once instead of three times, but the per-step wiring stays explicit and readable rather than folded into a loop or a table. probeInstalledBinary is a three-line adapter over check.ProbeAt: present = Status != StatusMissing, version = Version. Each block calls it with its own install.<X>BinaryPath accessor, passed by name.

## Reglas verificables
- **[auto]** with a healthy engram/codegraph/proofshot binary present ONLY at its own install.<X>BinaryPath() (absent from the test's PATH), Detect's ComponentState.Present for that binary is true — internal/setup/sync/binary_detection_test.go::TestProbeInstalledBinary_Engram|Codegraph|Proofshot.
- **[auto]** one binary's unresolvable canonical path leaves the other two Detect verdicts intact — internal/setup/sync/binary_detection_test.go::TestProbeInstalledBinary_Independence.
- **[auto]** two consecutive Detect passes over the same healthy canonical binary both report present, with no state mutation between calls — internal/setup/sync/binary_detection_test.go::TestProbeInstalledBinary_SecondRunIdempotent.
- **[auto]** a present, up-to-date engram never returns ActionInstall from decide() — internal/setup/sync/sync_test.go::TestDecide_MatecitoAIDevBuild.

## Alternativas consideradas
(a) Fold all three binaries into a single loop over a (name, resolver, args) table — the shared binary registry approach, discarded: it centralizes the (name, resolver, version-args) triple that requirement 4 ("El lugar de instalación es propio de cada paso") wants owned per step, and it is exactly the Approach B the user already declined at the tasks gate. (b) Have each block call check.RunVersion with a manually resolved absolute path inline, with no shared adapter — three copies of the resolve-then-probe glue instead of one, reopening the drift the new spec requirement forbids. (c) Move the three blocks into the verify-surface checks' own package and have sync import them — a layering inversion, the install planner depending on the environment reporter.

## Consecuencias
sync.go's Detect keeps its familiar three-block shape; the only material change per block is which function it calls and what resolver it passes. probeInstalledBinary's signature is unchanged from before this decision, so nothing outside sync.go needed to change to accommodate it. Any future fourth binary is another explicit block plus its own install.<X>BinaryPath accessor — deliberately not a table entry.
