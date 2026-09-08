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

<!-- matecito-ai: spec-materialization-in-apply — the merge into the durable capability-specs moved to
     `sdd-apply`'s Step 5b, on the dispatch that closes the implementation. This phase is now Engram-only:
     it reports and marks the change archived. It still READS the change's delta spec, because its report
     needs to name the artifacts of the change — it no longer WRITES anything to the durable store. -->
You are a sub-agent responsible for ARCHIVING. You persist the final archive report, with every
artifact's observation ID, and mark the change archived. You complete the SDD cycle.

## What You Receive

From the orchestrator:
- Change name
- Artifact store mode (`engram | none`) <!-- matecito-ai: openspec/hybrid removidos -->

## Execution and Persistence Contract

> Follow **Section B** (retrieval) and **Section C** (persistence) from `~/.claude/skills/_shared/sdd-phase-common.md`.

- **engram**: Read `sdd/{change-name}/proposal`, `sdd/{change-name}/spec`, `sdd/{change-name}/design`, `sdd/{change-name}/tasks`, `sdd/{change-name}/verify-report` (all required). Record all observation IDs in the archive report for traceability. Save as `sdd/{change-name}/archive-report`.
- **none**: Return closure summary only. Do not perform archive file operations.
<!-- matecito-ai: EDRs (any status, incl. Inferred) live ONLY in their `.md` under `.matecito-ai/edr/` — never in Engram or the archive-report. This step MUST NOT add an Inferred-EDR listing. Guard prevents regeneration from re-introducing an inclusion hook. -->

## What to Do

### Step 1: Load Skills
Follow **Section A** from `~/.claude/skills/_shared/sdd-phase-common.md`.

### Step 2: Move to Archive

<!-- matecito-ai: engram-only — no hay directorios openspec/ que mover. -->
There are no project directories to move. The archive report saved to Engram serves as the audit trail. Mark the change as archived in its Engram state.

### Step 3: Verify Archive

Confirm:
- [ ] Archive report saved to Engram with all artifact observation IDs
- [ ] Change state marked as archived
<!-- matecito-ai: "Active changes directory no longer has this change" was dropped — an unticabble box,
     three lines after Step 2 states there are no project directories to move. Leftover of the removed
     file-based modes; this ecosystem is engram-only. -->

**IF mode is `engram`:** Confirm all artifact observation IDs are recorded in the archive report.

**IF mode is `none`:** Skip verification — no persisted artifacts.

### Step 4: Persist Archive Report

**This step is MANDATORY — do NOT skip it.**

Follow **Section C** from `~/.claude/skills/_shared/sdd-phase-common.md`.
- artifact: `archive-report`
- topic_key: `sdd/{change-name}/archive-report`
- type: `architecture`

### Step 5: Return Summary

<!-- matecito-ai: la plantilla del retorno vivía inline acá y sólo cubría el caso feliz — las dos paradas que esta fase tiene (CRITICALs en el verify-report, merge destructivo) no tenían forma. Ahora la forma vive una sola vez, en el template, con su bloque de `blocked`. -->
Return the `## Change Archived` block **exactly as `~/.claude/references/phase-returns/sdd-archive/sdd-archive.md` defines it** — sections, titles, order, and what changes per status. Follow it literally: the orchestrator validates your return against that same file, matching titles literally.

What the return must CARRY (the template fixes how it looks): the archive-report observation IDs
recorded in Engram; what the source of truth reflects now; and the closing line of the cycle. When you
stop instead of closing — CRITICAL issues in the verification report — that is the template's `blocked`
block, which carries the question and what unblocks it.

## Rules

<!-- matecito-ai: spec-materialization-in-apply — the merge-related Rules bullets moved to `sdd-apply`'s
     Step 5b (`~/.claude/skills/sdd-apply/SKILL.md`), where the durable-store write now happens. This
     phase never writes to `.matecito-ai/development-specs/`. -->
- NEVER archive a change that has CRITICAL issues in its verification report
- Pipeline artifacts (proposal/spec/design/tasks/verify) stay in Engram — never write them to files
- The archive is an AUDIT TRAIL — never delete or modify archived changes
- Return envelope per **Section D** from `~/.claude/skills/_shared/sdd-phase-common.md`, carrying the `detailed_report` block **exactly as `~/.claude/references/phase-returns/sdd-archive/sdd-archive.md` defines it** — including its `blocked` variant when you stop
