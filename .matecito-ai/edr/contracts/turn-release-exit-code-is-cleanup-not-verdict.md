# EDR — Release exits zero on both Released and Mismatch — only a usage error is non-zero

- **Status:** Accepted
- **Date:** 2026-08-15

## Contexto
`contracts/turn-claimed-by-explicit-call.md` pins `claim`'s three outcomes to exit codes explicitly. `release` has two outcomes of its own — Released and Mismatch — and nothing had pinned what either one exits with. In code both have always returned nil from `RunE`, so both exit 0, with only a stderr line distinguishing a clean release from a mismatched one. The gap was never a code defect: the code already behaved this way. It was that nobody had decided whether it should.

## Decisión
`release` exits 0 on both outcomes — a matching token that deletes the turn, and a token that no longer matches (already released, or moved) — with the mismatch always reported on stderr, never silently. Only a genuine usage error (no `--token` supplied, the repository cannot be resolved) exits non-zero. This is deliberately asymmetric with `claim`: `release` runs as cleanup at the end of a chain whose write already happened, and its own exit code is not what the chain composes on.

## Reglas verificables
- **[auto]** `release` exits 0 when the token matches and the turn is deleted.
- **[auto]** `release` exits 0 when the token does not match, and reports the mismatch on stderr rather than staying silent.
- **[auto]** `release` exits non-zero only for a usage error — no token supplied — never for a mismatch.
- **[manual]** No instruction anywhere chains a write, or any other command, on `release`'s exit code — `release` is cleanup, not a gate.

## Alternativas consideradas
Making `release` exit non-zero on a mismatch, mirroring `claim`'s asymmetry — rejected: a chain like `claim && git merge && release` would report a failure for a merge that succeeded, the moment the release step alone hit a mismatch (a slow release racing a queued caller's retry, for instance). The gap was not a missing feature; it was an undecided default that happened to already be correct.

## Consecuencias
A caller cannot tell "released" from "mismatch" by exit code alone — it has to read stderr, or call `turn status` to check what actually happened. That is accepted: release's job is to try the deletion and say what happened, not to gate anything downstream of it, since by the time release runs the unit of work is already over.
