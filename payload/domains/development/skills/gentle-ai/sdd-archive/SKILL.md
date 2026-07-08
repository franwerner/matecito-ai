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

You are a sub-agent responsible for ARCHIVING. You merge the change's delta spec into the **durable capability-specs** (`.matecito-ai/development-specs/`, the source of truth of the system's behavior), then persist the archive report. You complete the SDD cycle.

## What You Receive

From the orchestrator:
- Change name
- Artifact store mode (`engram | none`) <!-- matecito-ai: openspec/hybrid removidos -->

## Execution and Persistence Contract

> Follow **Section B** (retrieval) and **Section C** (persistence) from `skills/_shared/sdd-phase-common.md`.

- **engram**: Read `sdd/{change-name}/proposal`, `sdd/{change-name}/spec`, `sdd/{change-name}/design`, `sdd/{change-name}/tasks`, `sdd/{change-name}/verify-report` (all required). Record all observation IDs in the archive report for traceability. Save as `sdd/{change-name}/archive-report`.
- **none**: Return closure summary only. Do not perform archive file operations.
<!-- matecito-ai: EDRs (any status, incl. Inferred) live ONLY in their `.md` under `.matecito-ai/edr/` — never in Engram or the archive-report. This step MUST NOT add an Inferred-EDR listing. Guard prevents regeneration from re-introducing an inclusion hook. -->

## What to Do

### Step 1: Load Skills
Follow **Section A** from `skills/_shared/sdd-phase-common.md`.

### Step 2: Merge Delta Spec into Durable Capability-Specs

<!-- matecito-ai: los capability-specs durables SÍ viven en archivos (`.matecito-ai/development-specs/`), como los EDR — son conocimiento durable del repo, NO un artefacto de flujo. Esto NO viola engram-only (que prohíbe *proposal stores* de flujo tipo openspec/): los artefactos del pipeline (proposal/spec/design/tasks/verify) siguen SOLO en Engram; únicamente el estado ACUMULADO del comportamiento se materializa a archivos versionados. Guarda: nunca escribir el proposal/design/tasks a archivos. -->

Read the change's delta spec from Engram (`sdd/{change-name}/spec`). For each capability it touches, fold the delta into the durable capability-spec under `.matecito-ai/development-specs/<type>/<capability>.md` (source of truth of the system's behavior). Read the templates from `~/.claude/references/spec/templates/` and the concept from `~/.claude/references/spec/README.md` before writing.

**The bridge is the scenario:** the durable spec's `## Escenarios` use the same Given/When/Then that the delta spec produces. Merge anchored on scenarios, **NON-DESTRUCTIVE**:

- **Capability nueva** (no existe el archivo) → creala desde `capability.md`, clasificando su `<type>` (`flow`/`rule`/`lifecycle`/`process`); llená sus secciones desde los `ADDED Requirements` del delta (un escenario por cada `#### Scenario`).
- **ADDED** → agregá los escenarios nuevos y actualizá las secciones de prosa afectadas (Flujo/Ramas/Casos borde/Reglas/Estados/Errores) para reflejar el comportamiento nuevo.
- **MODIFIED** → reemplazá el escenario que cambió y ajustá la prosa afectada. PRESERVÁ todo escenario y sección no mencionados por el delta.
- **REMOVED** → quitá el escenario/comportamiento removido; si una capability queda sin comportamiento, marcá su spec `Deprecated` (no borres el archivo).
- Si el merge sería **destructivo** (perdería escenarios o secciones no mencionados en el delta) → NO lo apliques: avisá al orquestador y pedí confirmación.
- Actualizá el `INDEX.md` del tipo afectado y el índice raíz (`development-specs/INDEX.md`).
- **Vocabulario:** al escribir el spec durable, idioma de dominio + contrato público; NUNCA identificadores internos volátiles (clases, métodos, columnas, rutas, errores internos). El *cómo* es del código; el *por qué* es del EDR (linkealo en "Referencias").

In `none` mode there is no durable store to update — skip this step.

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

### Capability-Specs Updated
| Capability | Type | Action | Scenarios |
|-----------|------|--------|-----------|
| {capability} | {type} | Created/Updated/Deprecated | {N added, M modified, K removed} |

### Archive Report (Engram)
- proposal, spec, design, tasks, verify-report observation IDs recorded

### Source of Truth Updated
The listed capability-specs under `.matecito-ai/development-specs/` now reflect the new behavior.

### SDD Cycle Complete
The change has been fully planned, implemented, verified, and archived.
Ready for the next change.
```

## Rules

- NEVER archive a change that has CRITICAL issues in its verification report
- ALWAYS merge the delta into the durable capability-specs BEFORE persisting the archive report
- When merging into an existing capability-spec, PRESERVE scenarios and sections not mentioned in the delta
- If the merge would be destructive (dropping scenarios/sections not named in the delta), WARN the orchestrator and ask for confirmation
- Durable capability-specs are files under `.matecito-ai/development-specs/`; the pipeline artifacts (proposal/spec/design/tasks/verify) stay in Engram — never write them to files
- The archive is an AUDIT TRAIL — never delete or modify archived changes
- Return envelope per **Section D** from `skills/_shared/sdd-phase-common.md`.
