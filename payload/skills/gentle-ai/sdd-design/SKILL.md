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

<!-- matecito-ai: read project ADRs before designing — START -->
#### Step 2a: Read the project's architecture decisions (ADRs)

If the project has `.matecito-ai/adr/` (decisions captured by `project-decisions-bootstrap`), you MUST consult it before proposing any architecture, because the design has to RESPECT decisions already made — not re-decide them.

1. Read `.matecito-ai/adr/INDEX.md` (root index) to see which domains exist.
2. For each domain this change touches (e.g. a new endpoint touches `security`, `contracts`, `runtime`), read the relevant ADRs in `.matecito-ai/adr/<domain>/` — focus on their **Decisión**, **Reglas concretas**, AND **Patrón aplicado** (if present) sections.
3. When writing the design:
   - **Respect every `Accepted` ADR that applies.** Your "Architecture Decisions" must align with them; cite the ADR (e.g. `security/auth.md`) when a decision is constrained by it.
   - **If an applicable ADR declares `Patrón aplicado: X`,** your design MUST implement X according to the canonical definition in `~/.claude/references/design-patterns/patterns/<x>.md` — read that file before designing. Cite the pattern by name in your Architecture Decisions. Do not propose a variant unless the ADR itself justifies the deviation.
   - If the design would **contradict** an `Accepted` ADR → STOP. Do not silently override a standing decision. Report it as a blocker in your return summary so the user can either adjust the design or change the decision via `project-decisions-bootstrap` (update mode).
   - If the change requires a decision **no ADR covers** (a genuinely new architectural choice) → flag it explicitly under "New decisions (not yet in ADRs)" in your design, and recommend capturing it via `project-decisions-bootstrap` before/with implementation. Document your proposed choice, but mark it as pending ADR capture. When your new decision IS a canonical pattern, name it (the catalog is available at `~/.claude/references/design-patterns/`) so the future ADR can record `Patrón aplicado: <Nombre>`.

If `.matecito-ai/adr/` does NOT exist, proceed normally (note that the project has no captured decisions).
<!-- matecito-ai: read project ADRs before designing — END -->

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
<!-- matecito-ai: when a decision is constrained by an existing ADR, add: **Constrained by**: `<domain>/<slug>.md` -->
<!-- matecito-ai: when the decision maps to a canonical pattern, add: **Patrón aplicado**: <Name> — <1-line why>. Definition at ~/.claude/references/design-patterns/patterns/<name>.md -->

### Decision: {Decision Title}

**Choice**: {What we chose}
**Alternatives considered**: {What we rejected}
**Rationale**: {Why this choice over alternatives}

<!-- matecito-ai: ADR alignment sections — START -->
## ADR Alignment

| Applicable ADR | Status | How this design respects it |
|----------------|--------|------------------------------|
| `<domain>/<slug>.md` | Accepted | {how the design complies} |

## New Decisions (not yet in ADRs)

{Architectural choices this change requires that NO existing ADR covers.
For each: the proposed choice + a note to capture it via project-decisions-bootstrap.
If none, state "None — all decisions are covered by existing ADRs."}

## ADR Conflicts (BLOCKER if any)

{Any place where this change would contradict an Accepted ADR.
If present, this design is BLOCKED until resolved. If none, state "None."}
<!-- matecito-ai: ADR alignment sections — END -->

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

### Open Questions
{List any unresolved questions, or "None"}

### Next Step
Ready for tasks (sdd-tasks).
```

## Rules

- ALWAYS read the actual codebase before designing — never guess
<!-- matecito-ai: ADRs are binding — respect Accepted ADRs in .matecito-ai/adr/; never contradict one silently (report as blocker); flag uncovered decisions for capture via project-decisions-bootstrap (see Step 2a) -->
- ALWAYS read `.matecito-ai/adr/` (if present) before designing; treat Accepted ADRs as binding constraints, surface conflicts as blockers, and flag new uncovered decisions
- Every decision MUST have a rationale (the "why")
- Include concrete file paths, not abstract descriptions
- Use the project's ACTUAL patterns and conventions, not generic best practices
- If you find the codebase uses a pattern different from what you'd recommend, note it but FOLLOW the existing pattern unless the change specifically addresses it
- Keep ASCII diagrams simple — clarity over beauty
- If you have open questions that BLOCK the design, say so clearly — don't guess
- **Size budget**: Design artifact MUST be under 800 words. Architecture decisions as tables (option | tradeoff | decision). Code snippets only for non-obvious patterns.
- Return envelope per **Section D** from `skills/_shared/sdd-phase-common.md`.
