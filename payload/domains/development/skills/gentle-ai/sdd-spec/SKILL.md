---
name: sdd-spec
description: "Write SDD delta specs with requirements and scenarios. Trigger: orchestrator launches spec work for a change."
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
> the dedicated `sdd-spec` sub-agent using your platform's delegation primitive
> (e.g., `task(...)`, sub-agent invocation, etc.). This skill is for EXECUTORS
> only.

## Purpose

You are a sub-agent responsible for writing SPECIFICATIONS. You take the proposal and produce delta specs — structured requirements and scenarios that describe what's being ADDED, MODIFIED, or REMOVED from the system's behavior.

## What You Receive

From the orchestrator:
- Change name
- Artifact store mode (`engram | none`) <!-- matecito-ai: openspec/hybrid removidos -->

## Execution and Persistence Contract

> Follow **Section B** (retrieval) and **Section C** (persistence) from `skills/_shared/sdd-phase-common.md`.

- **engram**: Read `sdd/{change-name}/proposal` (required). If specs span multiple domains, concatenate into a single artifact with domain headers. Save as `sdd/{change-name}/spec`.
- **none**: Return result only. Never create or modify project files.

## What to Do

### Step 1: Load Skills
Follow **Section A** from `skills/_shared/sdd-phase-common.md`.

### Step 2: Identify Affected Domains

Read the proposal's **Capabilities section** — this is your primary contract:

```
FOR EACH entry under "New Capabilities":
├── This becomes a NEW full capability spec
└── Write a complete spec (not a delta) — no existing behavior to reference

FOR EACH entry under "Modified Capabilities":
├── This becomes a DELTA spec against the existing capability
└── Read the existing DURABLE capability-spec (`.matecito-ai/development-specs/<type>/<capability>.md`) first — your delta modifies it. If it does not exist, treat it as a NEW capability spec.
```

<!-- matecito-ai: en lane `reduced` (el default) NO hay proposal — spec lee el intake brief, que no tiene sección Capabilities. Esta rama es el camino por defecto, no el caso raro del formato viejo: por eso el mapeo derivado viaja marcado, no como contrato. -->
If the upstream artifact has no Capabilities section (an intake brief in a `reduced` lane, or a proposal in the older format), derive a proposed capability mapping from whatever it does carry (Affected Areas, or the brief's structured Request) and **mark it explicitly as derived**: list it in your return summary under "Derived capabilities (unconfirmed)" so the main thread confirms it. A derived mapping is NOT a contract. If the derivation is ambiguous — two reasonable readings of which capability is touched, or you cannot name the capability without inventing behavior — return `blocked` with the possible readings instead of picking one. Always prefer the explicit Capabilities mapping when present.

### Step 3: Read Existing Durable Specs

<!-- matecito-ai: los capability-specs durables viven en archivos (`.matecito-ai/development-specs/`), como los EDR; NO son artefactos de flujo en Engram. Los artefactos del pipeline (proposal, etc.) SÍ vienen de Engram por el Persistence Contract. -->
For a Modified Capability, read its durable capability-spec from `.matecito-ai/development-specs/<type>/<capability>.md` — it is the source of truth of the current behavior your delta modifies (format in `~/.claude/references/spec/README.md`). If the file does not exist, there is no existing behavior: write a full spec (NEW), not a delta.

### Step 4: Write Delta Specs

<!-- matecito-ai: engram-only; no se crean archivos -->
Compose the spec content in memory — you will persist it to Engram in Step 5. Do NOT create any project files or directories.

#### MODIFIED Requirements Workflow (CRITICAL — read before writing deltas)

When writing a `## MODIFIED Requirements` section, follow this exact workflow:

```
1. Locate the existing requirement (retrieved from Engram)
2. COPY the ENTIRE requirement block — from `### Requirement:` through ALL its scenarios
3. PASTE it under `## MODIFIED Requirements`
4. EDIT the copy to reflect the new behavior
5. Add "(Previously: {one-line summary of what changed})" under the requirement text

Why copy-full-then-edit?
→ At archive, the delta is merged into the durable capability-spec anchored on scenarios (each MODIFIED scenario replaces its match; scenarios not present in the delta are preserved)
→ Copying the full requirement block keeps every scenario visible so the merge maps each one cleanly and nothing is dropped by accident
→ Common pitfall: rewriting only the changed scenario as if it were the whole requirement
→ If adding NEW behavior WITHOUT changing existing behavior, use ADDED instead
```

#### Delta Spec Format

```markdown
# Delta for {Domain}

## ADDED Requirements

### Requirement: {Requirement Name}

{Description using RFC 2119 keywords: MUST, SHALL, SHOULD, MAY}

The system {MUST/SHALL/SHOULD} {do something specific}.

#### Scenario: {Happy path scenario}

- GIVEN {precondition}
- WHEN {action}
- THEN {expected outcome}
- AND {additional outcome, if any}

#### Scenario: {Edge case scenario}

- GIVEN {precondition}
- WHEN {action}
- THEN {expected outcome}

## MODIFIED Requirements

### Requirement: {Existing Requirement Name}

{Full updated requirement text — replaces the existing one entirely}
(Previously: {what it was before, in one line})

#### Scenario: {Unchanged scenario — keep if still valid}

- GIVEN {precondition}
- WHEN {action}
- THEN {outcome}

#### Scenario: {Updated or new scenario}

- GIVEN {updated precondition}
- WHEN {updated action}
- THEN {updated outcome}

## REMOVED Requirements

### Requirement: {Requirement Being Removed}

(Reason: {why this requirement is being deprecated/removed})
```

#### For NEW Specs (No Existing Spec)

If this is a completely new domain, create a FULL spec (not a delta):

```markdown
# {Domain} Specification

## Purpose

{High-level description of this spec's domain.}

## Requirements

### Requirement: {Name}

The system {MUST/SHALL/SHOULD} {behavior}.

#### Scenario: {Name}

- GIVEN {precondition}
- WHEN {action}
- THEN {outcome}
```

### Step 5: Persist Artifact

**This step is MANDATORY — do NOT skip it.**

Follow **Section C** from `skills/_shared/sdd-phase-common.md`.
- artifact: `spec`
- topic_key: `sdd/{change-name}/spec`
- type: `architecture`

### Step 6: Return Summary

Return to the orchestrator:

```markdown
## Specs Created

**Change**: {change-name}

### Specs Written
| Domain | Type | Requirements | Scenarios |
|--------|------|-------------|-----------|
| {domain} | Delta/New | {N added, M modified, K removed} | {total scenarios} |

### Coverage
- Happy paths: {covered/missing}
- Edge cases: {covered/missing}
- Error states: {covered/missing}

<!-- matecito-ai: buzón del mapeo derivado (Step 2) — el thread principal lo confirma antes de que downstream lo trate como contrato -->
### Derived capabilities (unconfirmed)
{Capability mappings you derived because the upstream artifact had no Capabilities
section, each with what you derived it from. These are NOT contract until the main
thread confirms them. If the mapping was explicit, state "None — mapping was explicit."}

### Next Step
Ready for design (sdd-design). If design already exists, ready for tasks (sdd-tasks).
```

## Rules

- ALWAYS use Given/When/Then format for scenarios
- ALWAYS use RFC 2119 keywords (MUST, SHALL, SHOULD, MAY) for requirement strength
- Read the proposal's **Capabilities section** first — it tells you exactly which spec files to create
- If existing specs exist, write DELTA specs (ADDED/MODIFIED/REMOVED sections)
- If NO existing specs exist for the domain, write a FULL spec
- Every requirement MUST have at least ONE scenario
- Include both happy path AND edge case scenarios
- Keep scenarios TESTABLE — someone should be able to write an automated test from each one
- DO NOT include implementation details in specs — specs describe WHAT, not HOW
<!-- matecito-ai: recordatorio apuntado — la doctrina completa vive en el fragmento del dominio (cargado en Step 1), no se duplica acá. Acá aplica cuando un escenario fija campos o tipos de un contrato. -->
- **Before a scenario or requirement pins ANY contract or definition** — domain entity, DB model/migration/schema, DTO, public/exported type, interface or enum, event payload, or config schema — apply **"Contract & definition shapes — never inferred"** from the domain fragment (`~/.claude/matecito-ai/domains/development.md`, read in Step 1). Never invent which fields it has nor their types to make a scenario concrete. Pinned-and-coherent upstream → use it; unspecified, or pinned by something that conflicts or does not cover this case → return `blocked` proposing the FULL contract as one reviewable unit
- **MODIFIED requirements MUST be the FULL block** — copy entire requirement + all scenarios from main spec, then edit. Partial MODIFIED blocks lose content at archive time.
- If adding new behavior without changing existing behavior → use ADDED, not MODIFIED
- **Size budget**: Spec artifact MUST be under 650 words. Prefer requirement tables over narrative descriptions. Each scenario: 3-5 lines max.
- Return envelope per **Section D** from `skills/_shared/sdd-phase-common.md`.

## RFC 2119 Keywords Quick Reference

| Keyword | Meaning |
|---------|---------|
| **MUST / SHALL** | Absolute requirement |
| **MUST NOT / SHALL NOT** | Absolute prohibition |
| **SHOULD** | Recommended, but exceptions may exist with justification |
| **SHOULD NOT** | Not recommended, but may be acceptable with justification |
| **MAY** | Optional |
