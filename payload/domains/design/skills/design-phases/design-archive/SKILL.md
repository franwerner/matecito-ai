---
name: design-archive
description: >
  Method for the design ARCHIVE phase. Use when verification has passed and the change must be
  closed — records the final deliverables, the brief and system, and persists the archive report with
  all observation IDs for traceability. Completes the design cycle.
---

# Design archive — method

The recipe the `design-archive` executor follows. This phase closes a verified change: it records the
final deliverables and all artifact observation IDs as an audit trail, then marks the change archived.

## Reads / writes (design Phase Read/Write)

- **Reads:** all change artifacts (`proposal`, `brief`, `system`, `tasks`, `verify-report`).
- **Writes:** `design/{change-name}/archive-report`.

## Steps

1. Read all change artifacts (required):
   - `mem_search("design/{change-name}/proposal")` → `mem_get_observation`
   - `mem_search("design/{change-name}/brief")` → `mem_get_observation`
   - `mem_search("design/{change-name}/system")` → `mem_get_observation`
   - `mem_search("design/{change-name}/tasks")` → `mem_get_observation`
   - `mem_search("design/{change-name}/verify-report")` → `mem_get_observation`
2. Write the final archive report with all observation IDs and produced deliverables for traceability.
3. Mark the change state as archived in Engram.
4. Persist the archive report to the active backend.

NEVER archive a change that has CRITICAL issues in its verify-report. DDRs are NOT recorded here —
DDRs (any status) live ONLY as `.md` under `.matecito-ai/ddr/`, never duplicated into Engram or the
archive-report.

## Mentor mode

Keep it brief: when the archive surfaces a recurring pattern worth naming for next time, note it in
1-2 lines and cite the `design-principles` catalog; derive to `explain-concept` only if asked. Do not
repeat the full mentor rule (it lives in the domain CLAUDE.md).
