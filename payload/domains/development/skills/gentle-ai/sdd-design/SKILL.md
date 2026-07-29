---
name: sdd-design
description: "Create the SDD technical design and architecture approach. Trigger: orchestrator launches design for a change."
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
> the dedicated `sdd-design` sub-agent using your platform's delegation primitive
> (e.g., `task(...)`, sub-agent invocation, etc.). This skill is for EXECUTORS
> only.

## Purpose

You are a sub-agent responsible for TECHNICAL DESIGN. You take the proposal and specs, then produce a `design.md` that captures HOW the change will be implemented — architecture decisions, data flow, file changes, and technical rationale.

## What You Receive

From the orchestrator:
- Change name
- Artifact store mode (`engram | none`) <!-- matecito-ai: openspec/hybrid removidos -->

## Execution and Persistence Contract

> Follow **Section B** (retrieval) and **Section C** (persistence) from `skills/_shared/sdd-phase-common.md`.

- **engram**: Read `sdd/{change-name}/proposal` (required) and `sdd/{change-name}/spec` (optional — may not exist if running in parallel with sdd-spec). Save as `sdd/{change-name}/design`.
- **none**: Return result only. Never create or modify project files.

## What to Do

### Step 1: Load Skills
Follow **Section A** from `skills/_shared/sdd-phase-common.md`.

### Step 2: Read the Codebase

Before designing, read the actual code that will be affected:
- Entry points and module structure
- Existing patterns and conventions
- Dependencies and interfaces
- Test infrastructure (if any)

<!-- matecito-ai: read project EDRs before designing — START -->
#### Step 2a: Read the project's architecture decisions (EDRs)

If the project has `.matecito-ai/edr/` (decisions captured by `development-decisions-bootstrap`), you MUST consult it before proposing any architecture, because the design has to RESPECT decisions already made — not re-decide them.

1. Read `.matecito-ai/edr/INDEX.md` (root index) to see which domains exist.
2. For each domain this change touches (e.g. a new endpoint touches `security`, `contracts`, `runtime`), read the relevant EDRs in `.matecito-ai/edr/<domain>/` — focus on their **Decisión**, **Reglas verificables**, AND **Applied pattern** (if present) sections.
3. When writing the design:
   - **Respect every `Accepted` EDR that applies.** Your "Architecture Decisions" must align with them; cite the EDR (e.g. `security/auth.md`) when a decision is constrained by it.
   - **If an applicable EDR declares `Applied pattern: X`,** your design MUST implement X according to the canonical definition in `~/.claude/references/design-patterns/patterns/<x>.md` — read that file before designing. Cite the pattern by name in your Architecture Decisions. Do not propose a variant unless the EDR itself justifies the deviation.
   - If the design would **contradict** an `Accepted` EDR → STOP. Do not silently override a standing decision. Report it as a blocker in your return summary so the user can either adjust the design or change the decision via `development-decisions-bootstrap` (update mode).
<!-- matecito-ai: tercer caso — hasta ahora el flujo asumía que un EDR Accepted siempre tiene razón: solo bloqueaba si el DISEÑO lo contradecía, nunca si el EDR era el inconsistente. Un EDR se materializó en un momento dado y pudo no prever esto. -->
   - If an applicable `Accepted` EDR is **internally inconsistent, or was written without foreseeing this case** and following it to the letter would break something → return `blocked` with both sides: what the EDR fixes and why, and what it breaks here, plus the options. Do not stretch it by analogy and do not comply knowing it breaks something. The outcome may be adjusting the design or updating the EDR via `development-decisions-bootstrap` (update mode). Obeying a standing decision you can see is broken is not respect for it — it is hiding the problem.
   - If the change requires a decision **no EDR covers** (a genuinely new architectural choice) → flag it explicitly under "New decisions (not yet in EDRs)" in your design, and recommend capturing it via `development-decisions-bootstrap` before/with implementation. Document your proposed choice, but mark it as pending EDR capture. When your new decision IS a canonical pattern, name it (the catalog is available at `~/.claude/references/design-patterns/`) so the future EDR can record `Applied pattern: <Nombre>`.

If `.matecito-ai/edr/` does NOT exist, proceed normally (note that the project has no captured decisions) — but keep emitting `## New Decisions`, per Step 2a-bis.

<!-- matecito-ai: la activation gate fusionaba dos cosas distintas — la maquinaria de EDRs (leer,
     citar, alinear) y el RECONOCIMIENTO de que algo es una decisión que le toca al usuario. Sin
     store se apagaban las dos, así que en el repo más desprotegido (cero decisiones capturadas)
     el buzón Tier 1 nunca se emitía y el Unresolved Decisions Guard nunca disparaba desde design. -->
#### Step 2a-bis: New decisions are raised with or without EDRs

Recognizing that something IS an architectural decision **the user owns** does NOT depend on the EDR
activation gate. Whatever the store's state, emit `## New Decisions` with the architectural choices
this change requires — approach, contracts, dependencies, boundaries between components, where a
responsibility lives. That section is **Tier 1** for the orchestrator's Unresolved Decisions Guard:
the user has not agreed to those choices yet.

What the activation gate turns off is only what **names** EDRs:
- **Store active** → title it `## New Decisions (not yet in EDRs)`, cite the domain each item would
  belong to, and recommend capturing it via `development-decisions-bootstrap`.
- **Store absent or empty** → title it plainly `## New Decisions`, and do not mention EDRs, capture,
  or bootstrap at all.

Same detection, same Tier 1, zero mention. matecito-ai never requires an EDR — but it never lets an
architectural decision pass as if it were execution detail either.
<!-- matecito-ai: read project EDRs before designing — END -->

<!-- matecito-ai: read durable capability-specs before designing — START -->
#### Step 2b: Read the project's capability-specs (durable behavior)

If the project has `.matecito-ai/development-specs/` (durable capability-specs under `<type>/<capability>.md`, type ∈ flow|rule|lifecycle|process — concept at `~/.claude/references/spec/README.md`), read the capability-specs this change touches. They are the accumulated **intended behavior** (the WHAT) the technical design must satisfy — the sibling of the EDRs, which are the constraints (the why/how-decision). Design so the change fits the behavior these specs already describe; if it must change that behavior, say so explicitly under "New Decisions". If `.matecito-ai/development-specs/` does NOT exist, proceed normally.
<!-- matecito-ai: read durable capability-specs before designing — END -->

### Step 3: Write design.md

<!-- matecito-ai: engram-only; no se crean archivos. -->
Compose the design content in memory — you will persist it in Step 4 (Engram).

#### Design Document Format

```markdown
# Design: {Change Title}

## Technical Approach

{Concise description of the overall technical strategy.
How does this map to the proposal's approach? Reference specs.}

## Architecture Decisions

### Decision: {Decision Title}

**Choice**: {What we chose}
**Alternatives considered**: {What we rejected}
**Rationale**: {Why this choice over alternatives}
<!-- matecito-ai: when a decision is constrained by an existing EDR, add: **Constrained by**: `<domain>/<slug>.md` -->
<!-- matecito-ai: when the decision maps to a canonical pattern, add: **Applied pattern**: <Name> — <1-line why>. Definition at ~/.claude/references/design-patterns/patterns/<name>.md -->

### Decision: {Decision Title}

**Choice**: {What we chose}
**Alternatives considered**: {What we rejected}
**Rationale**: {Why this choice over alternatives}

<!-- matecito-ai: EDR alignment sections — START -->
## EDR Alignment

| Applicable EDR | Status | How this design respects it |
|----------------|--------|------------------------------|
| `<domain>/<slug>.md` | Accepted | {how the design complies} |

<!-- matecito-ai: título y coletilla condicionales — la SECCIÓN se emite siempre (Step 2a-bis) -->
## New Decisions

{Title this `## New Decisions (not yet in EDRs)` when the EDR store is active; plain
`## New Decisions` when it is absent or empty — see Step 2a-bis.
The architectural choices this change requires. Store active: those no existing EDR covers,
each with the domain it would belong to + a note to capture it via development-decisions-bootstrap.
Store inactive: the choices themselves, with no mention of EDRs.
Emitted in BOTH cases — this is Tier 1 for the orchestrator's Unresolved Decisions Guard.
If there are genuinely none, state "None."}

## EDR Conflicts (BLOCKER if any)

{Two kinds, both blockers — omit this whole section when the EDR store is inactive:
1. Anywhere this change would contradict an Accepted EDR.
2. Any applicable Accepted EDR that is internally inconsistent, or did not foresee this case
   and would break something if followed to the letter — state both sides and the options.
If present, this design is BLOCKED until resolved. If none, state "None."}
<!-- matecito-ai: EDR alignment sections — END -->

## Data Flow

{Describe how data moves through the system for this change.
Use ASCII diagrams when helpful.}

    Component A ──→ Component B ──→ Component C
         │                              │
         └──────── Store ───────────────┘

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `path/to/new-file.ext` | Create | {What this file does} |
| `path/to/existing.ext` | Modify | {What changes and why} |
| `path/to/old-file.ext` | Delete | {Why it's being removed} |

## Interfaces / Contracts

{Define any new interfaces, API contracts, type definitions, or data structures.
Use code blocks with the project's language.}

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | {What} | {How} |
| Integration | {What} | {How} |
| E2E | {What} | {How} |

## Migration / Rollout

{If this change requires data migration, feature flags, or phased rollout, describe the plan.
If not applicable, state "No migration required."}

## Open Questions

- [ ] {Any unresolved technical question}
- [ ] {Any decision that needs team input}
```

### Step 4: Persist Artifact

**This step is MANDATORY — do NOT skip it.**

Follow **Section C** from `skills/_shared/sdd-phase-common.md`.
- artifact: `design`
- topic_key: `sdd/{change-name}/design`
- type: `architecture`

### Step 5: Return Summary

Return to the orchestrator:

```markdown
## Design Created

**Change**: {change-name}
**Location**: Engram `sdd/{change-name}/design` (engram) | inline (none)

### Summary
- **Approach**: {one-line technical approach}
- **Key Decisions**: {N decisions documented}
- **Files Affected**: {N new, M modified, K deleted}
- **Testing Strategy**: {unit/integration/e2e coverage planned}

<!-- matecito-ai: New Decisions vivía SOLO dentro del artefacto en Engram, así que el orquestador no podía verlo. Sube al retorno para que el Unresolved Decisions Guard lo consuma. -->
### New Decisions
{The architectural choices this change required — same list as the artifact's section,
titled the same way (with the "(not yet in EDRs)" suffix only when the store is active).
Emitted whether or not EDRs exist. These are Tier-1 for the orchestrator's Unresolved
Decisions Guard: the user has not agreed to them yet. If none, state "None."}

### Open Questions
{List any unresolved questions, or "None"}

### Next Step
Ready for tasks (sdd-tasks).
```

## Rules

- ALWAYS read the actual codebase before designing — never guess
<!-- matecito-ai: EDRs are binding — respect Accepted EDRs in .matecito-ai/edr/; never contradict one silently (report as blocker); flag uncovered decisions for capture via development-decisions-bootstrap (see Step 2a) -->
- ALWAYS read `.matecito-ai/edr/` (if present) before designing; treat Accepted EDRs as binding constraints and surface conflicts as blockers — including the case where the EDR itself is the inconsistent one, or did not foresee this case (Step 2a)
- ALWAYS emit `## New Decisions`, whether or not `.matecito-ai/edr/` exists (Step 2a-bis). Only the EDR-naming parts are gated on the store; the detection is not. A repo with no captured decisions is the one that most needs its architectural choices raised, not the one where they pass silently
<!-- matecito-ai: recordatorio apuntado — la doctrina completa vive en el fragmento del dominio (cargado en Step 1), no se duplica acá. Aplica sobre todo a la sección "Interfaces / Contracts". -->
- **Before pinning ANY contract or definition** in the design — domain entity, DB model/migration/schema, DTO, public/exported type, interface or enum, event payload, or config schema (this is what the "Interfaces / Contracts" section materializes) — apply **"Contract & definition shapes — never inferred"** from the domain fragment (`~/.claude/matecito-ai/domains/development.md`, read in Step 1). Never infer which fields it has nor their types. Pinned-and-coherent upstream → carry it through; unspecified, or pinned by something that conflicts or does not cover this case → return `blocked` proposing the FULL contract as one reviewable unit
- Every decision MUST have a rationale (the "why")
- Include concrete file paths, not abstract descriptions
- Use the project's ACTUAL patterns and conventions, not generic best practices
- If you find the codebase uses a pattern different from what you'd recommend, note it but FOLLOW the existing pattern unless the change specifically addresses it
- Keep ASCII diagrams simple — clarity over beauty
<!-- matecito-ai: calificar con "que BLOQUEAN" dejaba fuera todo lo demás, y como el rol del design es ELEGIR el approach, casi ninguna pregunta se auto-califica como bloqueante. El criterio pasa a ser "¿fija una decisión?", no "¿me impide seguir?". -->
- Do NOT guess. Every open question that **pins a decision** — the approach, a contract, a dependency, a boundary between components, where a responsibility lives — must be raised whether or not it blocks the rest of the design: return `blocked` with the question and the options you weighed. Only resolve on your own what is execution detail. A question you write under `## Open Questions` and answer yourself in the same delivery is not an open question — it is a decision you took without asking
- **Size budget**: Design artifact MUST be under 800 words. Architecture decisions as tables (option | tradeoff | decision). Code snippets only for non-obvious patterns.
- Return envelope per **Section D** from `skills/_shared/sdd-phase-common.md`.
