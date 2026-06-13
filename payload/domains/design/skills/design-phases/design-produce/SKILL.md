---
name: design-produce
description: >
  Method for the design PRODUCE phase (design's "apply"). Use when tasks are ready and production
  should begin — generate the deliverables (Figma frames, brand guides, exported assets) following
  the locked visual system, marking tasks complete as they finish.
---

# Design produce — method

The recipe the `design-produce` executor follows. This phase generates the actual deliverables and
records WHAT was produced and where as markdown progress; the visual work itself lives in Figma /
exported files.

## Reads / writes (design Phase Read/Write)

- **Reads:** `design/{change-name}/tasks` + `design/{change-name}/brief` +
  `design/{change-name}/system` + `design/{change-name}/produce-progress`. The brief is the floor;
  tasks / system are optional (absent in reduced / custom lanes) + **DDRs touched** (when active).
- **Writes:** `design/{change-name}/produce-progress`.

## Capability skills this phase invokes

- `generate-assets` — produce the assigned pieces consistent with the locked system and brand guide.

This phase ORCHESTRATES that skill; the technique lives in it — do not duplicate it here.

## Steps

1. Read brief artifact (required — the floor): `mem_search("design/{change-name}/brief")` →
   `mem_get_observation`.
2. Read tasks artifact if present: `mem_search("design/{change-name}/tasks")` → if found,
   `mem_get_observation`; if absent (reduced / custom lane), produce directly from the brief.
3. Read system artifact if present: `mem_search("design/{change-name}/system")` → if found,
   `mem_get_observation`; if absent, there is no locked system to follow.
3a. DDR activation gate: if `.matecito-ai/ddr/` is absent or empty, DDRs are inactive — skip this
    step silently. If active: read the applicable DDRs in `.matecito-ai/ddr/` — when a system exists,
    use the ones listed in its DDR Alignment; without a system (reduced / custom lane), read
    `.matecito-ai/ddr/INDEX.md` for the touched surfaces. Treat their concrete rules as hard
    constraints. If the system flagged a DDR conflict / uncaptured decision as a blocker → return
    `blocked`.
3b. Read previous produce-progress (if it exists):
    `mem_search("design/{change-name}/produce-progress")` → if found, `mem_get_observation` → read
    and merge (skip completed tasks, merge when saving).
4. Produce the assigned pieces (`generate-assets`) — every piece consistent with the locked system
   and brand guide.
5. Mark each task `[x]` complete as you finish it (update the tasks artifact via `mem_update`).
6. Persist progress to the active backend.

## Mentor mode

Explain the WHY behind any non-obvious production choice in 1-2 lines — the design principle applied
(e.g. why optical alignment over metric here). Cite the `design-principles` catalog; when a concept
the person may not know surfaces, derive to `explain-concept`. Do not repeat the full mentor rule
(it lives in the domain CLAUDE.md).
