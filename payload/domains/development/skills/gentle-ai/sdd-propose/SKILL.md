---
name: sdd-propose
description: "Create an SDD change proposal with intent, scope, and approach. Trigger: orchestrator launches proposal work for a change."
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
> the dedicated `sdd-propose` sub-agent using your platform's delegation primitive
> (e.g., `task(...)`, sub-agent invocation, etc.). This skill is for EXECUTORS
> only.

## Purpose

<!-- matecito-ai: acá se describía un `proposal.md` "inside the change folder". No hay change folder:
     los artefactos del pipeline viven en Engram (regla del dominio: "Never write pipeline artifacts
     to the filesystem"). Solo el conocimiento durable —EDRs, capability-specs— es archivo. -->
You are a sub-agent responsible for creating PROPOSALS. You take the exploration analysis (or direct user input) and produce a structured `proposal` artifact, persisted in Engram (Step 5). You write NO files.

## What You Receive

From the orchestrator:
- Change name (e.g., "add-dark-mode")
- Exploration analysis (from sdd-explore) OR direct user description
- Artifact store mode (`engram | none`) <!-- matecito-ai: openspec/hybrid removidos -->

## Execution and Persistence Contract

> Follow **Section B** (retrieval) and **Section C** (persistence) from `~/.claude/skills/_shared/sdd-phase-common.md`.

<!-- matecito-ai: `sdd-init/{project}` had no declared format, so reading it meant guessing at its
     structure. Its shape is now fixed in init-details.md; read it by section. -->
- **engram**: Read `sdd/{change-name}/explore` (optional) and `sdd-init/{project}` (optional — its shape is fixed by `## Project Context Format` in `~/.claude/skills/sdd-init/references/init-details.md`: `### Stack`, `### Architecture`, `### Conventions`, read by those titles; an axis marked `— not detected` is a gap, not a value to assume). Save artifact as `sdd/{change-name}/proposal`.
- **none**: Return result only. Never create or modify project files.

## What to Do

### Step 1: Load Skills
Follow **Section A** from `~/.claude/skills/_shared/sdd-phase-common.md`.

### Step 2: Read Existing Context

<!-- matecito-ai: openspec/hybrid removidos; engram-only. No se crean directorios openspec/. -->
Existing context was already retrieved from Engram in the Persistence Contract (or passed inline in `none` mode). Do NOT create any project files or directories. Skip filesystem reads.

### Step 3: Compose the proposal

<!-- matecito-ai: esto es el CONTENIDO del artefacto que se persiste en Engram en el Step 5,
     no un archivo a escribir en el repo. -->
This is the body of the `proposal` artifact you persist in Step 5 — compose it, do not write it to disk.

```markdown
# Proposal: {Change Title}

## Intent

{What problem are we solving? Why does this change need to happen?
Be specific about the user need or technical debt being addressed.}

## Scope

### In Scope
- {Concrete deliverable 1}
- {Concrete deliverable 2}
- {Concrete deliverable 3}

### Out of Scope
- {What we're explicitly NOT doing}
- {Future work that's related but deferred}

## Capabilities

> This section is the CONTRACT between proposal and specs phases.
> The sdd-spec agent reads this to know exactly which spec files to create or update.


### New Capabilities
<!-- Capabilities being introduced.
     Use kebab-case names (e.g., user-auth, data-export, api-rate-limiting).
     Leave empty if no new capabilities. -->
- `<capability-name>`: <brief description of what this capability covers>

### Modified Capabilities
<!-- Existing capabilities whose REQUIREMENTS are changing (not just implementation).
     Only list here if spec-level behavior changes. Each needs a delta spec.
     Leave empty if none. -->
- `<existing-capability-name>`: <what requirement is changing>

## Approach

{High-level technical approach. How will we solve this?
Reference the recommended approach from exploration if available.}

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `path/to/area` | New/Modified/Removed | {What changes} |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| {Risk description} | Low/Med/High | {How we mitigate} |

## Rollback Plan

{How to revert if something goes wrong. Be specific.}

## Dependencies

- {External dependency or prerequisite, if any}

## Success Criteria

- [ ] {How do we know this change succeeded?}
- [ ] {Measurable outcome}
```

### Step 5: Persist Artifact

**This step is MANDATORY — do NOT skip it.**

Follow **Section C** from `~/.claude/skills/_shared/sdd-phase-common.md`.
- artifact: `proposal`
- topic_key: `sdd/{change-name}/proposal`
- type: `architecture`

### Step 6: Return Summary

<!-- matecito-ai: la plantilla literal salió de acá. Mantenerla inline Y en el template crea una
     copia más para desincronizar — que es exactamente el defecto que el template vino a cerrar. -->
The shape of your return lives in **`~/.claude/references/phase-returns/sdd-propose/sdd-propose.md`**. Read it
and follow it literally: it declares the `## Proposal Created` block — its sections, their titles,
their order, which ones are unconditional, and what changes when you return `blocked`. The
orchestrator validates your return against that same file, matching titles literally, so a section
you drop, rename or re-level is a gate that never fires.

Two things that file makes explicit and that this phase gets wrong most often:

<!-- matecito-ai: esta fase pasó a tener buzón. Fija el approach y el mapeo de capabilities —
     lo que `sdd-spec` consume como contrato— y hasta ahora nada de eso llegaba a un gate: en
     Automatic la propuesta pasaba derecho y el error asomaba una fase después. -->
- `### Scope and approach (unconfirmed)` is this phase's **Tier-1 mailbox**, emitted always. Carries
  the approach you fixed and the capability mapping (`New` / `Modified`), which is what `sdd-spec`
  turns into full specs or deltas — a wrong name there writes the delta against the wrong capability.
  It is never `None`: a proposal always fixes an approach and always touches some capability. The
  artifact's `## Capabilities` section must still be filled in and must match what you report here;
  leaving it as a placeholder is what forces `sdd-spec` to derive its own mapping downstream.
- Anything that would fix scope or approach on the user's behalf returns `blocked`, with the
  question and the options in the `### Blocker` section that file designates — never in `risks`, and
  never resolved by picking the option you prefer so the happy-path block can be emitted.
- When that `blocked` stop is over an unspecified **contract or definition** — per the domain
  fragment's "Contract & definition shapes — never inferred" — it does NOT stay as free `### Blocker`
  prose: emit `### Contract Shapes Proposed` (`has_contract_proposals: true`), one compound item per
  contract, per that same return template. The ratified (or user-adjusted) shape comes back to you
  only in a later re-dispatch's own instructions, never re-read from a stored artifact — see
  `~/.claude/matecito-ai/domains/development.md`, "Forwarding a ratified contract shape to the
  proposing phase."
- Every item under `### Scope and approach (unconfirmed)` also carries its own `anchor` — the concrete
  source it traces to. Free-form, per the anchor criterion in `~/.claude/skills/_shared/sdd-phase-common.md`,
  Section D.3: `<repo-path>[:line]` or `<engram-key>`, start line only (say the range in words), and a
  point nothing has written yet anchors to what surfaced the need, never to a file that would only
  exist if the item is confirmed. You supply it — nothing derives it for you.

## Rules

<!-- matecito-ai: "change directory" no existe — la continuidad se resuelve leyendo el artefacto
     previo en Engram (`sdd/{change-name}/proposal`), no un directorio del repo. -->
- If a `proposal` artifact for this change already exists in Engram, READ it first and UPDATE it
- Keep the proposal CONCISE - it's a thinking tool, not a novel
- Every proposal MUST have a rollback plan
- Every proposal MUST have success criteria
- Use concrete file paths in "Affected Areas" when possible
- **ALWAYS fill in the Capabilities section** — this is the contract with sdd-spec.
- New Capabilities → each becomes a new capability spec for sdd-spec
- Modified Capabilities → each will become a delta spec inside `sdd-spec`'s `spec` artifact
- If nothing changes at the spec level (pure refactor, config change), explicitly write "None" under both sub-sections — don't leave them as template placeholders
- **Size budget**: Proposal artifact MUST be under 450 words. Use bullet points and tables over prose. Headers organize, not explain.
- Return envelope per **Section D** from `~/.claude/skills/_shared/sdd-phase-common.md`.
