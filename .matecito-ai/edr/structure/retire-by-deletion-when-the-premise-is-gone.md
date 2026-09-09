# EDR — Retire by deletion only when the record's premise is gone

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
`structure/spec-deletion-over-deprecation.md` ratified deletion instead of `Deprecated`-with-link for the fourteen cockpit-only capability-specs `sdd/drop-apps-subtree` retired, framed explicitly as a one-off exception for that change — it left `.matecito-ai/development-specs/INDEX.md:48` ("marcá el spec Deprecated con link a su reemplazo; no borres el archivo") untouched, and its own Alternativas section rejected relying on `sdd-archive`'s maintenance line for the same reason a future record needed to exist. This change retires two more capability-specs with no replacement — `flow/materialize-mined-artifacts.md` and `rule/decision-re-emergence-reconfirmation.md` — and faces the identical question a third time. Writing a third per-change exception record repeats the precedent's shape without ever generalizing the criterion a per-change record can't state.

## Decisión
A reusable criterion decides the case, not a per-change exception record: a record whose **premise** disappeared is deleted as a file; a capability **replaced** by another way of doing the same thing keeps the store's `Deprecated`-with-link convention, still stated exactly as written at `.matecito-ai/development-specs/INDEX.md:48`. Both capability-specs this change retires are premise-gone — `materialize-mined-artifacts` documented a post-gate write path whose only actors (`development-decisions-mine`, `development-spec-mine`) are deleted in the same change, and `decision-re-emergence-reconfirmation` documented reconfirming a per-change ratification ledger this change also deletes — so both are deleted as files, not marked `Deprecated`.

## Reglas verificables
- **[manual]** A retired record is deleted as a file when its premise (the mechanism, actor, or artifact it documents) no longer exists; it is marked `Deprecated` with a link to its replacement only when a successor covers the same capability a different way. This is the criterion consulted for every future retirement — not a per-change record copied from precedent.
- **[manual]** `.matecito-ai/development-specs/INDEX.md:48` keeps its unconditional "marcá el spec Deprecated... no borres el archivo" wording unedited by this criterion; the criterion decides which case a retirement falls into, it does not rewrite the convention's text.
- **[manual]** `flow/materialize-mined-artifacts.md` and `rule/decision-re-emergence-reconfirmation.md` are absent as files after this change and neither is left marked `Deprecated`, consistent with both being premise-gone rather than replaced.

## Alternativas consideradas
(a) Copy `structure/spec-deletion-over-deprecation.md`'s shape and write a third per-change exception record enumerating this change's two specs. Rejected: it accumulates one record per change and answers nothing for the next one, and its stated rationale ("no replacement to link") does not even hold cleanly here — `flow/materialize-mined-artifacts` has a partial successor in `flow/materialize-records-straight-through`, which a per-change record would have to explain away each time. (b) Rewrite `.matecito-ai/development-specs/INDEX.md:48` into an explicit two-case convention. Rejected: `structure/spec-deletion-over-deprecation.md` explicitly ratified that line as untouched, and rewriting the store's general convention would put a convention edit inside a pruning change's scope, which this change was never ratified to do.

## Consecuencias
Future retirements — in this change and any later one — consult this criterion instead of writing a new per-change exception record or re-litigating whether `INDEX.md:48`'s wording still holds. The criterion agrees with `structure/spec-deletion-over-deprecation.md`'s own carve-out and supersedes nothing it decided; that record remains the account of why the fourteen-spec batch was itself premise-gone.
