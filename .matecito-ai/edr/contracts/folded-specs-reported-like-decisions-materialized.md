# EDR — The delta-spec fold is reported both in the return and in the cumulative apply-progress artifact, mirroring Decisions Materialized

- **Status:** Accepted
- **Date:** 2026-09-04

## Contexto
Moving the delta-spec fold from `sdd-archive` into `sdd-apply`'s Step 5b left a gap: `sdd-archive` used to report what it merged in its own `### Capability-Specs Updated` return section; that reporting duty had to move somewhere too, or the fold would happen with nothing telling the orchestrator or a later reader that it ran. `sdd-apply` already carries a section with the exact shape this needed — `### Decisions Materialized` — reported in both its return AND the cumulative `apply-progress` artifact, conditional on at least one materialization this run, no status filter.

## Decisión
New section `### Capability-Specs Materialized` in `sdd-apply.yaml`, `emitted: conditional`, `when: has_capability_specs_materialized`, no `statuses` filter, `render: table`, columns `Capability-spec | Action | Scenarios` — the exact column shape `sdd-archive.yaml`'s retired `### Capability-Specs Updated` carried, so the moved content keeps its shape unchanged. It goes into the cumulative `apply-progress` artifact too, per the Merge Protocol, exactly as `### Decisions Materialized` does — cumulative across batches, present only when at least one batch's Step 5b actually wrote a durable capability-spec. `sdd-archive` gains NO new read to keep reporting this: its own `### Source of Truth Updated` line becomes phase-generic instead of enumerating merges it no longer performs.

## Reglas verificables
- **[auto]** `sdd-apply.yaml` declares `### Capability-Specs Materialized` as `emitted: conditional`, gated on `has_capability_specs_materialized`, with no `statuses` filter — checked by `render-return.js --phase sdd-apply --schema`.
- **[manual]** The section's columns are `Capability-spec | Action | Scenarios`, matching the column shape `sdd-archive.yaml`'s retired `### Capability-Specs Updated` carried.
- **[manual]** The section is also written into the cumulative `apply-progress` artifact, per the Merge Protocol, cumulative across batches, no status filter — a partial fold that then hit the destructive stop still reports what it folded.
- **[manual]** `sdd-archive` reads no new artifact to keep reporting the fold; its `### Source of Truth Updated` line is phase-generic, naming no specific merge mechanism.

## Alternativas consideradas
(a) Return-only, mirroring `### Rejected Proposals Checked` — rejected: that section's return-only asymmetry has a specific cause (a task under a `conflicts` verdict never reaches `[x]`, so the cumulative task list has nothing to carry for it), which does not apply here — a folded capability-spec IS a durable, persisted fact that later readers need. (b) Artifact-only — rejected: the orchestrator never reads the artifact, so the fold would be invisible at the point it happened, with no line in the return to notify anyone. (c) Keep the reporting inside `sdd-archive` by having it read `apply-progress` for the fold's outcome — rejected: it adds a new read to archive's contract and a Read/Write row this change did not ratify, purely to report something archive no longer does.

## Consecuencias
`sdd-apply`'s return and artifact both gain a fourth conditional table section, alongside `### Decisions Materialized`, `### TDD Cycle Evidence` and `### Test Summary` — all following the same emitted-conditional, no-status-filter shape. `sdd-archive`'s own return loses a section and a whole reporting duty, simplifying its contract to report+mark-archived only.
