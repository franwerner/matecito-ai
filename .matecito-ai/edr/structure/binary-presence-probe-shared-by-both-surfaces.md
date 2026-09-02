# EDR — One presence predicate shared by both surfaces

- **Status:** Accepted
- **Date:** 2026-09-02

## Contexto
Both matecito-ai install and matecito-ai verify decide the same question — is this step's binary installed? — and before this change both answered it by handing a bare name to check.RunVersion, which resolves through exec.LookPath against the invoking session's PATH: sync.Detect did it at sync.go:235/:252/:269, and the three verify checks did it at internal/checks/engram/engram.go:21, internal/checks/codegraph/codegraph.go:16 and internal/checks/proofshot/proofshot.go:14. One question, answered twice, wrongly, in two packages — and the two answers could disagree about the same machine.

## Decisión
Add check.ProbeAt(name string, resolve func() (string, error), args []string, required bool, fixHint string) Result to internal/check, beside RunVersion — the leaf package both surfaces already import, so the shared predicate creates zero new dependency edges. It takes the binary's canonical location as a resolve closure (it cannot import internal/setup/install: that would cycle, since install imports check), resolves it, and delegates to RunVersion with the absolute path. A resolver error is StatusMissing without calling RunVersion — a step that cannot establish where its own binary would live has no honest way to call it installed. RunVersion's "no encontrado en PATH" detail literal is hoisted into a package const so ProbeAt can substitute a detail naming the resolved path on the healthy branch; RunVersion's own signature and behavior come out byte-identical. sync.probeInstalledBinary becomes a three-line adapter over check.ProbeAt (mapping present = Status != StatusMissing, version = Version); the three verify checks call check.ProbeAt directly through their own resolver seam.

## Reglas verificables
- **[auto]** given a resolve closure returning an absolute path to a binary that runs and prints a version, check.ProbeAt returns StatusOK with that version and a Detail naming the resolved path — internal/check/check_test.go::TestProbeAt.
- **[auto]** given a resolve closure that errors, check.ProbeAt returns StatusMissing without ever invoking RunVersion (Version stays empty) — internal/check/check_test.go::TestProbeAt/resolver_error_reports_Missing_without_calling_RunVersion.
- **[auto]** RunVersion's own bare-name behavior (Status/Detail on a missing binary) stays byte-identical after the const hoist — internal/check/check_test.go::TestRunVersion_Unaffected.
- **[auto]** sync.probeInstalledBinary's (present, version) tuple matches check.ProbeAt's Status/Version 1:1 across the same five scenarios — internal/setup/sync/binary_detection_test.go.

## Alternativas consideradas
(a) Export sync.ProbeInstalledBinary and have internal/checks/* import internal/setup/sync — a layering inversion (the environment reporter depending on the install planner) and it drags Detect/decide/Sync into three check packages that need none of it. (b) Each surface keeps its own probe — two implementations of the predicate the spec's new requirement ("Instalar y reportar el estado no se contradicen") says must never disagree, i.e. the drift itself. (c) Change RunVersion to take a resolver — breaks the ratified deferral on RunVersion's bare-name meaning and touches five internal/checks/prereqs call sites (claude, node, npm, npx, git) for which PATH is the correct question. (d) ProbeAt re-implementing the exec instead of delegating to RunVersion — duplicates the existence/health logic and the version parse. (e) ProbeAt doing its own os.Stat to avoid the const hoist — one extra syscall and a second definition of "present", and it still leaves the stale "no encontrado en PATH" detail for a present-but-not-executable file.

## Consecuencias
internal/check gains its first test file (check_test.go) and a second public entry point beside RunVersion, but no new dependency edges: both surfaces already imported this package. sync.probeInstalledBinary's exported signature is unchanged from before this decision, so its callers (the three Detect blocks) needed no further changes beyond passing an install.<X>BinaryPath resolver by name. A future third caller of "is this binary installed at its canonical location" reaches for check.ProbeAt instead of writing a sixth ad hoc probe.
