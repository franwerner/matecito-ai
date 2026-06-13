---
name: design-verify
description: >
  Method for the design VERIFY phase. Use when produce reports done (or partial) and the work must be
  validated against the brief, the locked system, the DDRs, and accessibility before archive — reads
  the real Figma colors, type, and hierarchy and flags CRITICAL / WARNING / SUGGESTION findings.
---

# Design verify — method

The recipe the `design-verify` executor follows. This phase is the quality gate: it proves the
produced work matches its contract, runs the guards against the REAL Figma values, and reports
issues (it does not fix them).

## Reads / writes (design Phase Read/Write)

- **Reads:** `design/{change-name}/brief` (required — the floor) + `design/{change-name}/system`
  (when present) + `design/{change-name}/produce-progress` (required) + **DDRs touched** (when
  active).
- **Writes:** `design/{change-name}/verify-report`.

## Capability skills this phase invokes

- `consistency-audit` — the `brand-consistency` guard: check each piece against the brand guide and
  the accepted DDRs.
- `visual-accessibility` — the WCAG guard: contrast ratios, minimum sizes, and hierarchy against the
  real colors and type.
- `design-review` — check each brief requirement against the produced work.

This phase ORCHESTRATES those skills; the technique lives in them — do not duplicate it here.

## Steps

1. Read brief artifact (required — the floor): `mem_search("design/{change-name}/brief")` →
   `mem_get_observation`.
2. Read system artifact if present: `mem_search("design/{change-name}/system")` → if found,
   `mem_get_observation`; if absent (reduced / custom lane), verify against the brief alone.
3. Read produce-progress (required): `mem_search("design/{change-name}/produce-progress")` →
   `mem_get_observation`.
4. If a Figma file is connected, READ it (`mcp__figma__*`) to inspect the real colors, type scale,
   and hierarchy of the produced pieces.
5. Run the `visual-accessibility` guard: WCAG contrast ratios, minimum sizes, hierarchy against the
   real values — flag anything below AA as CRITICAL.
6. Run the `brand-consistency` guard (`consistency-audit`): check each piece against the brand guide
   and the accepted DDRs — flag a piece that contradicts a decision record.
7. Check each brief requirement against the produced work (`design-review`) — flag CRITICAL /
   WARNING / SUGGESTION.
7b. DDR activation gate: if `.matecito-ai/ddr/` is absent or empty, DDRs are inactive — skip the
    DDR-compliance check silently. If active: for the DDRs touched by this change (or
    `.matecito-ai/ddr/<surface>/` for touched surfaces), confirm the produced work honors their
    concrete rules. Any violation → CRITICAL `DDR-VIOLATION` (cite the DDR).
8. Confirm tasks are marked complete and match the produced state.
9. Persist the verify report to the active backend.

## Mentor mode

Explain the WHY behind each finding in 1-2 lines — the design principle violated or honored (e.g. why
a 3:1 contrast fails AA for body text, why this spacing breaks the rhythm). Cite the
`design-principles` catalog; when a concept the person may not know surfaces, derive to
`explain-concept`. Do not repeat the full mentor rule (it lives in the domain CLAUDE.md).
