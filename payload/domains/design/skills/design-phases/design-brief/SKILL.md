---
name: design-brief
description: >
  Method for the design BRIEF phase (design's "spec"). Use when a direction is chosen and the work
  needs a formal alignment artifact — what MUST be true of the finished design, with acceptance
  criteria — before the system is locked or any asset is produced.
---

# Design brief — method

The recipe the `design-brief` executor follows. The brief is design's alignment artifact (the
equivalent of development's spec): it describes WHAT the finished design must achieve, not HOW it is
built. It does NOT lock the visual system or produce assets.

## Reads / writes (design Phase Read/Write)

- **Reads:** `design/{change-name}/proposal` (required); falls back to the nearest available upstream
  — the intake brief — in reduced / custom lanes where no proposal ran.
- **Writes:** `design/{change-name}/brief`.

## Steps

1. Read the upstream artifact — directions / proposal if present, else fall back to the intake brief:
   `mem_search("design/{change-name}/proposal")`; if it has no result,
   `mem_search("design/{change-name}/intake")` → `mem_get_observation`. The nearest available
   upstream is the source of requirements.
2. Extract requirements from that upstream artifact (the chosen direction, or the intake brief in
   reduced / custom lanes).
3. Write the brief — what MUST be true of the finished design (surfaces, tone, constraints,
   deliverables).
4. Add acceptance criteria (how we will know the design is "done" and on-brief).
5. Persist the brief to the active backend.

## Mentor mode

Explain the WHY behind each requirement in 1-2 lines — the design principle it enforces (e.g. why a
minimum contrast or a fixed type scale is a hard criterion). Cite the `design-principles` catalog;
when a concept the person may not know surfaces, derive to `explain-concept`. Do not repeat the full
mentor rule (it lives in the domain CLAUDE.md).
