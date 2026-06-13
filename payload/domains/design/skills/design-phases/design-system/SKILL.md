---
name: design-system
description: >
  Method for the design SYSTEM phase (design's "design" phase). Use when a direction is chosen and
  the visual system must be fixed before assets are produced — palette, type scale, grid, spacing,
  components — with the rationale for each choice. Reads and writes DDRs (Design Decision Records).
---

# Design system — method

The recipe the `design-system` executor follows. This phase locks the visual system — the HOW at the
visual-foundation level — and is where DDRs are read and captured. It does NOT slice tasks or produce
final assets.

## Reads / writes (design Phase Read/Write)

- **Reads:** `design/{change-name}/brief` (required) + **DDRs** (when active).
- **Writes:** `design/{change-name}/system`, plus any new / updated DDRs under `.matecito-ai/ddr/`.

## Capability skills this phase invokes

- `brand-guide` — produce the locked brand guide (palette, type scale, grid, spacing, key components)
  with rationale.

This phase ORCHESTRATES that skill; the technique lives in it — do not duplicate it here.

## Steps

1. Read the brief (required): `mem_search("design/{change-name}/brief")` → `mem_get_observation`.
1b. DDR activation gate: if `.matecito-ai/ddr/` is absent or empty, DDRs are inactive — skip all DDR
    steps (1b and 4b) silently, no mention. If active: read root `INDEX.md` + the DDRs of the
    surfaces this work touches. Accepted DDRs are binding constraints.
2. If a Figma file is connected, READ its existing styles / components (`mcp__figma__*`) as the
   starting point.
3. Lock the system — palette, type scale, grid, spacing, key components — and the brand guide
   (`brand-guide`).
4. Capture DDR-style decisions with rationale and rejected alternatives ("the palette is X because…").
4b. Align decisions with existing DDRs (cite them). If the system contradicts an Accepted DDR →
    return `blocked`. If it needs a brand decision no DDR covers → flag it for capture via
    project-decisions-bootstrap.
5. Persist the system (and any new / updated DDRs under `.matecito-ai/ddr/`) to the active backend.

DDRs live ONLY as `.md` under `.matecito-ai/ddr/` (with an `INDEX.md`) — never duplicated into Engram.

## Mentor mode

Explain the WHY behind each system choice in 1-2 lines — the design principle (e.g. why this type
scale ratio, why this neutral, why an 8pt grid). Cite the `design-principles` catalog; when a concept
the person may not know surfaces, derive to `explain-concept`. Do not repeat the full mentor rule
(it lives in the domain CLAUDE.md).
