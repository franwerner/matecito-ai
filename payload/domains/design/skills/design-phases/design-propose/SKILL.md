---
name: design-propose
description: >
  Method for the design PROPOSE phase. Use when exploration is complete and the idea is ready to be
  formalized into a small set of distinct directions to choose from, each with its intent, scope, and
  visual rationale.
---

# Design propose — method

The recipe the `design-propose` executor follows. This phase turns exploration into a small set of
distinct directions. It does NOT lock the brief or produce final assets — it proposes the directions,
nothing more.

## Reads / writes (design Phase Read/Write)

- **Reads:** `design/{change-name}/explore` (optional).
- **Writes:** `design/{change-name}/proposal` (the directions).

## Capability skills this phase invokes

- `explore-variations` — generate a small set of genuinely distinct visual directions, not minor
  variants of one idea.

This phase ORCHESTRATES that skill; the technique lives in it — do not duplicate it here.

## Steps

1. Read exploration artifact (optional): `mem_search("design/{change-name}/explore")` →
   `mem_get_observation`.
2. Define intent (what the design must convey, for whom, what success looks like).
3. Define scope (in-scope / out-of-scope surfaces and pieces, explicit).
4. Produce a small set of distinct directions (`explore-variations`) — each with its visual rationale.
5. Persist the directions to the active backend.

## Mentor mode

Explain the WHY behind each direction in 1-2 lines — the design principle it leans on (e.g. why one
direction is type-led and another image-led). Cite the `design-principles` catalog; when a concept
the person may not know surfaces, derive to `explain-concept`. Do not repeat the full mentor rule
(it lives in the domain CLAUDE.md).
