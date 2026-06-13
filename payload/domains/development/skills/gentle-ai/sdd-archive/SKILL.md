---
name: sdd-archive
description: "Archive a completed SDD change by syncing delta specs. Trigger: orchestrator launches archive after implementation and verification."
disable-model-invocation: true
user-invocable: false
license: MIT
metadata:
  author: gentleman-programming
  version: "2.0"
  delegate_only: true
---

> **ORCHESTRATOR GATE**: If you loaded this skill via the `skill()` tool, you are
> the ORCHESTRATOR — STOP. Do NOT execute these instructions inline. Delegate to
> the dedicated `sdd-archive` sub-agent using your platform's delegation primitive
> (e.g., `task(...)`, sub-agent invocation, etc.). This skill is for EXECUTORS
> only.

## Purpose

You are a sub-agent responsible for ARCHIVING. You merge delta specs into the main specs (source of truth), then move the change folder to the archive. You complete the SDD cycle.

## What You Receive

From the orchestrator:
- Change name
- Artifact store mode (`engram | none`) <!-- matecito-ai: openspec/hybrid removidos -->

## Execution and Persistence Contract

> Follow **Section B** (retrieval) and **Section C** (persistence) from `skills/_shared/sdd-phase-common.md`.

- **engram**: Read `sdd/{change-name}/proposal`, `sdd/{change-name}/spec`, `sdd/{change-name}/design`, `sdd/{change-name}/tasks`, `sdd/{change-name}/verify-report` (all required). Record all observation IDs in the archive report for traceability. Save as `sdd/{change-name}/archive-report`.
- **none**: Return closure summary only. Do not perform archive file operations.
<!-- matecito-ai: ADRs (any status, incl. Inferred) live ONLY in their `.md` under `.matecito-ai/adr/` — never in Engram or the archive-report. This step MUST NOT add an Inferred-ADR listing. Guard prevents regeneration from re-introducing an inclusion hook. -->

## What to Do

### Step 1: Load Skills
Follow **Section A** from `skills/_shared/sdd-phase-common.md`.

### Step 2: Sync Delta Specs to Main Specs

<!-- matecito-ai: engram-only — no hay filesystem sync de specs. -->
Artifacts live in Engram only. There is no filesystem spec tree to sync. The archive report (Step 5) records all observation IDs for traceability. In `none` mode, there's nothing to sync.

### Step 3: Move to Archive

<!-- matecito-ai: engram-only — no hay directorios openspec/ que mover. -->
There are no project directories to move. The archive report saved to Engram serves as the audit trail. Mark the change as archived in its Engram state.

### Step 4: Verify Archive

Confirm:
- [ ] Archive report saved to Engram with all artifact observation IDs
- [ ] Change state marked as archived
- [ ] Active changes directory no longer has this change

**IF mode is `engram`:** Confirm all artifact observation IDs are recorded in the archive report.

**IF mode is `none`:** Skip verification — no persisted artifacts.

### Step 5: Persist Archive Report

**This step is MANDATORY — do NOT skip it.**

Follow **Section C** from `skills/_shared/sdd-phase-common.md`.
- artifact: `archive-report`
- topic_key: `sdd/{change-name}/archive-report`
- type: `architecture`

### Step 6: Return Summary

Return to the orchestrator:

```markdown
## Change Archived

**Change**: {change-name}
**Archived to**: Engram archive report (engram) | inline (none)

### Specs Synced
| Domain | Action | Details |
|--------|--------|---------|
| {domain} | Created/Updated | {N added, M modified, K removed requirements} |

### Archive Contents
- proposal.md ✅
- specs/ ✅
- design.md ✅
- tasks.md ✅ ({N}/{N} tasks complete)

### Source of Truth Updated
The following specs now reflect the new behavior:

### SDD Cycle Complete
The change has been fully planned, implemented, verified, and archived.
Ready for the next change.
```

## Rules

- NEVER archive a change that has CRITICAL issues in its verification report
- ALWAYS sync delta specs BEFORE moving to archive
- When merging into existing specs, PRESERVE requirements not mentioned in the delta
- Use ISO date format (YYYY-MM-DD) for archive folder prefix
- If the merge would be destructive (removing large sections), WARN the orchestrator and ask for confirmation
- The archive is an AUDIT TRAIL — never delete or modify archived changes
- Return envelope per **Section D** from `skills/_shared/sdd-phase-common.md`.
