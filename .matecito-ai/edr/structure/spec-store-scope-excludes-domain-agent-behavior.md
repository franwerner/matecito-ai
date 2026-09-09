# EDR — The capability-spec store's scope excludes domain-agent behavior

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
Retiring `flow/materialize-mined-artifacts.md` removes the only capability-spec that described a post-gate materialization step, and `design`'s own mine-materialization step (its `design-decisions-mine` executor, gated and materialized the same way) is not separately documented anywhere in `.matecito-ai/development-specs/`. Deleting the retired spec could therefore read as also deleting coverage `design` depended on. It does not: the retired spec's `## Actores` names only the two `development` executors (`development-decisions-mine`, `development-spec-mine`); the flow it describes invokes `render-artifact.js`, which lives under `payload/domains/development/scripts/` and is referenced nowhere under `payload/domains/design/`; its INDEX paths are all `.matecito-ai/edr/**`; and across all 49 capability-specs in the store, zero describe `design`-domain agent behavior. `design`'s equivalent step is already documented where `design`'s behavior lives, `payload/domains/design/CLAUDE.md`'s "Decision-Gap Capture — design specifics" section — not in this store.

## Decisión
The capability-spec store's scope excludes `design`-domain agent behavior, matching its own declared scope at `.matecito-ai/development-specs/INDEX.md:5` ("el store cubre el **CLI** de matecito-ai"). No narrowed successor to `flow/materialize-mined-artifacts.md` scoped to `design` is written. `design`'s mine-materialization step stays documented exactly where it already is, `payload/domains/design/CLAUDE.md`, and this record exists so a future sweep over the retired spec's deletion does not re-raise the question of whether `design` lost coverage.

## Reglas verificables
- **[manual]** No capability-spec under `.matecito-ai/development-specs/` describes `design`-domain agent behavior; a spec whose `## Actores` names only `design` executors, or whose flow calls tooling that lives under `payload/domains/design/`, does not belong in this store.
- **[manual]** `design`'s equivalent to a `development` mechanism is documented in `payload/domains/design/CLAUDE.md` (or a design-domain artifact), never as a narrowed capability-spec in `.matecito-ai/development-specs/` written to preserve coverage for a retired `development`-only spec.
- **[manual]** `.matecito-ai/development-specs/INDEX.md:5`'s declared scope ("el CLI de matecito-ai") is the test consulted before writing any new capability-spec whose actor is a domain other than `development`'s CLI surface.

## Alternativas consideradas
Author a narrowed `flow/materialize-mined-artifacts` scoped to `design`, so the retirement leaves a successor naming `design-decisions-mine` in place of the two deleted `development` executors. Rejected: it would be the first record in this store to describe `design`-domain agent behavior, expanding the store's own declared scope as a side effect of a pruning change the user ratified as `development`-only.

## Consecuencias
`flow/materialize-mined-artifacts.md` is deleted with no successor of any scope; `design`'s mine-materialization behavior remains sole-sourced to `payload/domains/design/CLAUDE.md`. A future capability-spec proposal for `design`-domain agent behavior is either rejected on this record's grounds, or is itself the ratified decision to widen the store's scope — a decision this record does not make.
