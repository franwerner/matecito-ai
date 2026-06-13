---
name: design-explore
description: >
  Method for the design EXPLORE phase. Use when a brand or screen must be thought through before
  anything is locked: extract a palette from references, audit an existing Figma file, compare visual
  approaches, or clarify the brief — investigation only, before any direction or brief is committed.
---

# Design explore — method

The recipe the `design-explore` executor follows. This phase investigates visual directions and
returns a structured analysis with a recommendation. It is investigation ONLY — it does not lock a
brief or produce assets.

## Reads / writes (design Phase Read/Write)

- **Reads:** `design/{change-name}/intake`.
- **Writes:** `design/{change-name}/explore`.

## Capability skills this phase invokes

- `brand-from-references` — extract brand signals (palette, mood, type feel) from supplied references.
- `figma-hygiene` — audit the connected Figma file's styles, components, and structure.

This phase ORCHESTRATES those skills; the technique lives in them — do not duplicate it here.

## Steps

1. Read the intake brief: `mem_search("design/{change-name}/intake")` → `mem_get_observation`.
2. Understand the topic to investigate (brand feel, target audience, medium, references).
3. If a Figma file is connected, READ it (`mcp__figma__*`) — review existing styles, components,
   structure. The Figma index is active only when a file is connected; absent it, work from
   references and the brief alone.
4. Extract brand signals from references (`brand-from-references`); audit Figma hygiene
   (`figma-hygiene`).
5. Compare visual approaches with a pros / cons / effort table.
6. Return structured analysis with a recommendation.

## Mentor mode

Explain the WHY behind each finding in 1-2 lines — the design principle (e.g. why this contrast
direction reads as "premium", why this grid carries the hierarchy). Cite the `design-principles`
catalog; when a concept the person may not know surfaces, derive to `explain-concept`. Do not repeat
the full mentor rule (it lives in the domain CLAUDE.md).
