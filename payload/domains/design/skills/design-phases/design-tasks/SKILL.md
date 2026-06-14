---
name: design-tasks
description: >
  Method for the design TASKS phase. Use when the brief and the visual system are both ready and the
  work needs to be sliced into actionable, ordered production pieces (screens, assets, states) — each
  with a verifiable criteria, ready to produce.
---

# Design tasks — method

The recipe the `design-tasks` executor follows. This phase decomposes the change into a production
checklist. It does NOT produce assets — it produces the checklist only.

## Reads / writes (design Phase Read/Write)

- **Reads:** `design/{change-name}/brief` (required — the floor) + `design/{change-name}/system`
  (when present).
- **Writes:** `design/{change-name}/tasks`.

## Steps

1. Read brief artifact (required — the floor): `mem_search("design/{change-name}/brief")` →
   `mem_get_observation`.
2. Read system artifact if present: `mem_search("design/{change-name}/system")` → if found,
   `mem_get_observation`; if absent (custom lane without system), decompose from the brief alone.
3. Decompose work into ordered pieces (each small enough to produce in isolation — a screen, an
   asset, a state).
4. Link each task to the brief requirement it satisfies.
5. Per task, add a verifiable `criteria:` sub-line (mandatory — consumed by `design-verify`). When
   DDRs are active and the task touches a brand decision, add an optional `· ddr: <surface>/<slug>`
   ref (slug-based). If `.matecito-ai/ddr/` is absent or empty, omit the `ddr` part; `criteria:` is
   always required. Keep the `- [ ]` checkbox so `design-produce` can mark progress.
<!-- matecito-ai: decision-gap detection hook
Active ONLY when flagDecisionGaps=true (does NOT depend on DDRs existing). When active: for EACH task that touches a brand decision, emit `· ddr: <surface>/<slug>` mapped to a design concern, whether or not the DDR exists — this overrides step 5's rule of "omit ddr when `.matecito-ai/ddr/` is empty". Then, for each `· ddr:`, check whether `.matecito-ai/ddr/<surface>/<slug>.md` exists: if NOT, the ref is a dangling decision gap — leave it verbatim (do not modify or mark it). With zero DDRs, every brand decision is a gap (bootstrap the first ones). The dangling refs in the artifact are the gap list that `design-verify` consumes. When flag off: byte-identical behavior to step 5, no mention. -->
5b. (Decision-gap detection — flag-gated) When `flagDecisionGaps=true` (regardless of DDR presence): emit a concern-mapped `· ddr: <surface>/<slug>` for each decision-touching task even if no DDR exists yet; a `· ddr:` whose target file is absent under `.matecito-ai/ddr/` is a decision gap — leave it verbatim (the dangling ref IS the signal; with zero DDRs every brand decision is a gap → bootstrap). Silent when flag off.
6. Mark which tasks can run in parallel vs sequential.
7. Persist tasks to the active backend.

## Mentor mode

Explain the WHY behind the slicing in 1-2 lines — the design principle (e.g. why the system-level
components are produced before the screens that compose them). Cite the `design-principles` catalog;
when a concept the person may not know surfaces, derive to `explain-concept`. Do not repeat the full
mentor rule (it lives in the domain CLAUDE.md).
