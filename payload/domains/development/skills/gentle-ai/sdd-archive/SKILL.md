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

> Follow **Section B** (retrieval) and **Section C** (persistence) from `~/.claude/skills/_shared/sdd-phase-common.md`.

- **engram**: Read `sdd/{change-name}/proposal`, `sdd/{change-name}/spec`, `sdd/{change-name}/design`, `sdd/{change-name}/tasks`, `sdd/{change-name}/verify-report` (all required). Record all observation IDs in the archive report for traceability. Save as `sdd/{change-name}/archive-report`.
- **none**: Return closure summary only. Do not perform archive file operations.
<!-- matecito-ai: EDRs (any status, incl. Inferred) live ONLY in their `.md` under `.matecito-ai/edr/` — never in Engram or the archive-report. This step MUST NOT add an Inferred-EDR listing. Guard prevents regeneration from re-introducing an inclusion hook. -->

## What to Do

### Step 1: Load Skills
Follow **Section A** from `~/.claude/skills/_shared/sdd-phase-common.md`.

### Step 2: Merge Delta Spec into Durable Capability-Specs

<!-- matecito-ai: los capability-specs durables SÍ viven en archivos (`.matecito-ai/development-specs/`), como los EDR — son conocimiento durable del repo, NO un artefacto de flujo. Esto NO viola engram-only (que prohíbe *proposal stores* de flujo tipo openspec/): los artefactos del pipeline (proposal/spec/design/tasks/verify) siguen SOLO en Engram; únicamente el estado ACUMULADO del comportamiento se materializa a archivos versionados. Guarda: nunca escribir el proposal/design/tasks a archivos. -->

Read the change's delta spec from Engram (`sdd/{change-name}/spec`). For each capability it touches, fold the delta into the durable capability-spec under `.matecito-ai/development-specs/<type>/<capability>.md` (source of truth of the system's behavior). Read the templates from `~/.claude/references/spec/templates/` and the concept from `~/.claude/references/spec/README.md` before writing.

**The bridge is the scenario:** the durable spec's `## Escenarios` use the same Given/When/Then that the delta spec produces. Merge anchored on scenarios, **NON-DESTRUCTIVE**:

<!-- matecito-ai: these seven bullets were in Spanish inside an English body. Translated without changing
     the instruction. The template section names (`## Escenarios`, Flujo/Ramas/…) are deliberately left
     in Spanish: they are literal titles of files that are written in Spanish, not prose. -->
<!-- matecito-ai: verification-scope-token, fourth closing pass. Three rounds in a row fixed the drift the
     previous round found and hardened the wording to protect exactly that spot — the token's backticks,
     then a path citation and an accent in scenario bodies — and each time the NEXT round found a new drift
     in a part the wording never named (this time: the Purpose paragraph's own prose, "capability" retyped
     as "capacidad"). Enumerating protected parts was the recurring defect, not which parts got named. The
     fix is at the mechanism: the default is COPY for everything that crosses the delta → durable boundary,
     and only what archive itself legitimately generates (because the delta never wrote it) is named as the
     exception. -->
- **Default: copy, don't recompose.** Any prose the delta already writes out in full, that crosses into
  the durable spec, is carried by reproducing the delta's characters — never by re-typing, re-describing,
  translating, re-flowing, or "correcting" it (an accent, a term, a citation, a parenthetical: none of
  those are yours to touch while merging). This is the default for the **whole document** on that
  boundary, not a named set of protected parts — it covers the `## Purpose`/`## Propósito` paragraph,
  every scenario's `GIVEN`/`WHEN`/`THEN` body (including a `verification:` token line and its inline
  Markdown markup — backticks, bold, italics — around the token's value), and any other sentence the
  delta already wrote out in full. The verb is **copy**, not reproduce or preserve: you are moving
  characters, not restating an idea in your own words.
  The only content this default does NOT bind is what archive itself legitimately generates or
  transforms because the delta never wrote it out in the first place: the file's `Status`/`Date`/
  `Components` header, a **new** capability's `Actores`/`Reglas de negocio`/`Entidades y estados`
  sections synthesized from the delta's `Requirements` prose (there is no literal counterpart in the
  delta to copy from), the `INDEX.md` entries, and — when the delta's `ADDED` or `MODIFIED`
  requirement lands on an **existing** capability — the durable prose sections (`Flujo`/`Ramas`/
  `Casos borde`/`Reglas`/`Estados`/`Errores`) that were **already in the file before this merge** and
  that the new or changed scenario now makes stale: reading them, deciding whether they still hold,
  and rewriting the ones that don't is synthesis, the same kind as the New-capability case, because
  the delta carries no literal counterpart in that section's shape to copy from — only a scenario body
  and/or `Requirements` prose to synthesize from. This fourth exception is scoped to that pre-existing
  prose alone: anything the delta itself writes out in full — a new scenario's `GIVEN`/`WHEN`/`THEN`
  body, a new sentence, a rewritten `Purpose` paragraph — is still copied by the general default, even
  on the ADDED/MODIFIED path, exactly as in the New-capability case. Outside those four, copy.
- **New capability** (the file does not exist) → create it from `capability.md`, classifying its `<type>` (`flow`/`rule`/`lifecycle`/`process`); fill its `## Escenarios` from the delta's `ADDED Requirements` (one scenario per `#### Scenario`, copied per the default above) and its surrounding `Actores`/`Reglas de negocio`/`Entidades y estados` by synthesis, per the exception above.
- **ADDED** → add the new scenarios (copied per the default) and update the affected prose sections (Flujo/Ramas/Casos borde/Reglas/Estados/Errores) to reflect the new behavior — that update is the fourth exception above: only the pre-existing prose the new scenario makes stale, synthesized, never the scenario itself.
- **MODIFIED** → replace the scenario that changed (copied per the default) and adjust the affected prose — same fourth exception: the adjustment is bound to pre-existing prose the change makes stale, not to the scenario, which is copied. PRESERVE every scenario and section the delta does not mention.
- **REMOVED** → drop the removed scenario/behavior; if a capability is left with no behavior at all, mark its spec `Deprecated` (do not delete the file).
<!-- matecito-ai: verification-scope-token. Canonical definition lives at `~/.claude/references/spec/README.md`
     (deployed path) — this bullet cites it, it does not redefine the token's form or semantics. -->
- **`verification:` token** (canonical definition: `~/.claude/references/spec/README.md`) → its line is bound by the copy default above like the rest of the scenario body it belongs to; what is specific to it is a precedence rule, not a text rule: do NOT reconcile the delta's token against the durable copy's — the delta is the source, so on a MODIFIED scenario the delta's token wins outright, unmerged, same as the rest of its content. A scenario with no token merges without one.
- If the merge would be **destructive** (losing scenarios or sections the delta does not mention) → do NOT apply it: tell the orchestrator and ask for confirmation.
- Update the `INDEX.md` of the affected type and the root index (`development-specs/INDEX.md`).
- **Vocabulary:** write the durable spec in domain language + public contract; NEVER volatile internal identifiers (classes, methods, columns, routes, internal errors). The *how* belongs to the code; the *why* belongs to the EDR (link it under "Referencias").

In `none` mode there is no durable store to update — skip this step.

### Step 3: Move to Archive

<!-- matecito-ai: engram-only — no hay directorios openspec/ que mover. -->
There are no project directories to move. The archive report saved to Engram serves as the audit trail. Mark the change as archived in its Engram state.

### Step 4: Verify Archive

Confirm:
- [ ] Archive report saved to Engram with all artifact observation IDs
- [ ] Change state marked as archived
<!-- matecito-ai: "Active changes directory no longer has this change" was dropped — an unticabble box,
     three lines after Step 3 states there are no project directories to move. Leftover of the removed
     file-based modes; this ecosystem is engram-only. -->

**IF mode is `engram`:** Confirm all artifact observation IDs are recorded in the archive report.

**IF mode is `none`:** Skip verification — no persisted artifacts.

### Step 5: Persist Archive Report

**This step is MANDATORY — do NOT skip it.**

Follow **Section C** from `~/.claude/skills/_shared/sdd-phase-common.md`.
- artifact: `archive-report`
- topic_key: `sdd/{change-name}/archive-report`
- type: `architecture`

### Step 6: Return Summary

<!-- matecito-ai: la plantilla del retorno vivía inline acá y sólo cubría el caso feliz — las dos paradas que esta fase tiene (CRITICALs en el verify-report, merge destructivo) no tenían forma. Ahora la forma vive una sola vez, en el template, con su bloque de `blocked`. -->
Return the `## Change Archived` block **exactly as `~/.claude/references/phase-returns/sdd-archive/sdd-archive.md` defines it** — sections, titles, order, and what changes per status. Follow it literally: the orchestrator validates your return against that same file, matching titles literally.

What the return must CARRY (the template fixes how it looks): which durable capability-specs were created, updated or deprecated and with how many scenarios; the archive-report observation IDs recorded in Engram; what the source of truth reflects now; and the closing line of the cycle. When you stop instead of closing — CRITICAL issues in the verification report, or a merge that would be destructive — that is the template's `blocked` block, which carries the question and what unblocks it.

## Rules

- NEVER archive a change that has CRITICAL issues in its verification report
- ALWAYS merge the delta into the durable capability-specs BEFORE persisting the archive report
- When merging into an existing capability-spec, PRESERVE scenarios and sections not mentioned in the delta
- ALWAYS copy — never retype, re-describe, translate, or "correct" — any prose the delta already writes out in full that crosses into the durable spec: this is the default for the whole document (the `Purpose`/`Propósito` paragraph, every scenario's `GIVEN`/`WHEN`/`THEN` body, any `verification:` token line and its inline Markdown markup), not a named set of protected parts. The only exceptions are content archive itself legitimately generates because the delta never wrote it — the `Status`/`Date`/`Components` header, a new capability's synthesized `Actores`/`Reglas de negocio`/`Entidades y estados`, the `INDEX.md` entries, and — on an ADDED/MODIFIED merge onto an **existing** capability — the pre-existing `Flujo`/`Ramas`/`Casos borde`/`Reglas`/`Estados`/`Errores` prose that the new or changed scenario makes stale (synthesized, the same way as the new-capability case). That fourth exception never extends to anything the delta itself writes out in full — a new or changed scenario body is still copied, never synthesized
- A scenario's `verification:` token (definition: `~/.claude/references/spec/README.md`) is bound by the copy rule above like the rest of its body; its one extra rule is precedence, not text: never reconcile the delta's token against the durable copy's — the delta is the source, so on a MODIFIED scenario the delta's token wins outright, unmerged
- If the merge would be destructive (dropping scenarios/sections not named in the delta), WARN the orchestrator and ask for confirmation
- Durable capability-specs are files under `.matecito-ai/development-specs/`; the pipeline artifacts (proposal/spec/design/tasks/verify) stay in Engram — never write them to files
- The archive is an AUDIT TRAIL — never delete or modify archived changes
- Return envelope per **Section D** from `~/.claude/skills/_shared/sdd-phase-common.md`, carrying the `detailed_report` block **exactly as `~/.claude/references/phase-returns/sdd-archive/sdd-archive.md` defines it** — including its `blocked` variant when you stop
