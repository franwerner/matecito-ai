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

> Follow **Section B** (retrieval) and **Section C** (persistence) from `~/.claude/skills/_shared/sdd-phase-common.md`.

<!-- matecito-ai: declaraba `proposal` como required, pero esta fase corre en lanes donde `propose` no
     existe — su propio agente hace fallback a spec y después al intake brief. Tercer archivo con el
     mismo defecto (los otros dos: `sdd-tasks`, `sdd-apply`); barrido: no quedan más. -->
- **engram**: Read the nearest available upstream — `sdd/{change-name}/proposal` when it exists, else `sdd/{change-name}/spec`, else `sdd/{change-name}/intake`. In `reduced` / `custom` lanes `propose` may not have run, and its absence is normal, not an error. `sdd/{change-name}/spec` may also be absent when this phase runs in parallel with `sdd-spec`. Save as `sdd/{change-name}/design`.
<!-- matecito-ai: the brief is the LAST option of the chain above, and in the lanes where this phase runs
     there is always a proposal or a spec: read only as a fallback, the `diagram` flag never arrived.
     It is read as well, always, as its own retrieval (same pattern as `ui-test` in `sdd-verify`). -->
- **engram, additionally and unconditionally**: read `sdd/{change-name}/intake` on top of whatever the chain above resolved — it carries the `diagram` flag this phase acts on (Step 3-bis). It is a separate retrieval, never a fallback: `intake` is a base phase, so the brief always exists.
- **none**: Return result only. Never create or modify project files.

## What to Do

### Step 1: Load Skills
Follow **Section A** from `~/.claude/skills/_shared/sdd-phase-common.md`.

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
   - If the change requires a decision **no EDR covers** (a genuinely new architectural choice) → flag it explicitly under "New Decisions (not yet in EDRs)" in your design, and recommend capturing it via `development-decisions-bootstrap` before/with implementation. Document your proposed choice, but mark it as pending EDR capture. When your new decision IS a canonical pattern, name it (the catalog is available at `~/.claude/references/design-patterns/`) so the future EDR can record `Applied pattern: <Nombre>`.

If `.matecito-ai/edr/` does NOT exist, proceed normally (note that the project has no captured decisions) — but keep emitting `## New Decisions`, per Step 2a-bis.

<!-- matecito-ai: la activation gate fusionaba dos cosas distintas — la maquinaria de EDRs (leer,
     citar, alinear) y el RECONOCIMIENTO de que algo es una decisión que le toca al usuario. Sin
     store se apagaban las dos, así que en el repo más desprotegido (cero decisiones capturadas)
     el buzón Tier 1 nunca se emitía y el Unresolved Decisions Guard nunca disparaba desde design. -->
#### Step 2a-bis: New decisions are raised with or without EDRs

Recognizing that something IS an architectural decision **the user owns** does NOT depend on the EDR
activation gate. Whatever the store's state, emit it in BOTH places — `## New Decisions` in the
artifact AND `### New Decisions` in your return to the orchestrator (the guard only ever reads the
return; an artifact-only section is invisible to it) — with the architectural choices
this change requires — approach, contracts, dependencies, boundaries between components, where a
responsibility lives. That section is **Tier 1** for the orchestrator's Unresolved Decisions Guard:
the user has not agreed to those choices yet.

<!-- matecito-ai: sin este puntero, un ejecutor que lee este step y no llega a ## Rules archiva acá
     una decisión que tenía que bloquear. El criterio NO se duplica: vive una sola vez en ## Rules. -->
Not every such choice belongs in this section. Apply the **blocking test** in `## Rules` first: when
the alternatives differ in new infrastructure, in the public contract, or in the data model, the
decision returns `blocked` instead — `## New Decisions` carries the ones the test clears.

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
<!-- matecito-ai: simetría con EDR Conflicts — sin store no se menciona NADA de EDRs, así que esta
     sección también se omite entera (no se emite con "None", que ya sería nombrarlos). -->
## EDR Alignment

{Omit this whole section when the EDR store is inactive (absent or empty) — with no store there is
nothing to align against and EDRs are not mentioned at all. When active, one row per applicable EDR;
if none applies, state "None."}

| Applicable EDR | Status | How this design respects it |
|----------------|--------|------------------------------|
| `<domain>/<slug>.md` | Accepted | {how the design complies} |

<!-- matecito-ai: título y coletilla condicionales — la SECCIÓN se emite siempre (Step 2a-bis) -->
## New Decisions

{Title this `## New Decisions (not yet in EDRs)` when the EDR store is active; plain
`## New Decisions` when it is absent or empty — see Step 2a-bis.
The architectural choices this change requires AND that pass the blocking test in `## Rules`
(alternatives that differ in new infrastructure, public contract or data model do NOT go here —
they return `blocked`). Store active: those no existing EDR covers,
each with the domain it would belong to + a note to capture it via development-decisions-bootstrap.
Store inactive: the choices themselves, with no mention of EDRs.
Emitted in BOTH cases — this is Tier 1 for the orchestrator's Unresolved Decisions Guard.
Every item carries its `· blocking-test:` token, per `## Rules`.
If there are genuinely none, state "None."}

- {the choice}: {what you chose} — {alternatives weighed, and why this one}
  · blocking-test: none

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

<!-- matecito-ai: rol acotado — Open Questions y New Decisions se solapaban y el ejecutor duplicaba
     (razonamiento en una, pregunta en la otra) => fatiga de confirmación. Acá NO van decisiones. -->
## Open Questions

{Only what does NOT pin a decision: implementation doubts and things to validate during apply.
Anything that fixes the approach, a contract, a dependency, a boundary or where a responsibility
lives belongs in `## New Decisions` — or returns `blocked`, per the blocking test in `## Rules`.
If none, state "None."}

- [ ] {An implementation detail to confirm while applying}
- [ ] {Something to measure or check against the running system}
```

<!-- matecito-ai: `sdd-intake` decides the `diagram` flag and writes it into the brief, and this phase is
     its only reader. It was not mentioned anywhere in this skill: it lived only in the agent file,
     which is a copy with no authority over content. The full doctrine (when to draw, and that the
     diagram is EPHEMERAL) lives in the domain fragment; here, only what this phase does with the flag. -->
### Step 3-bis: Architecture diagram flag

Read the flag at its literal location in the intake brief: the line `- Diagram: {needed|not-needed}`
under `### Classification`. Do not look for it anywhere else and do not re-derive it — `sdd-intake`
decided it and the user confirmed it at the INTAKE GATE.

- `needed` → add ONE clause to your `executive_summary` recommending a live diagram of the chosen
  architecture. That is the whole action.
- `not-needed`, or the line is absent → skip silently, no mention.

You do NOT generate, export or persist any diagram: you have no drawio tooling and diagrams are
ephemeral in this ecosystem (live preview rendered by the main thread, zero `.drawio` files in the
repo). The recommendation goes in `executive_summary` — never in `risks`, which Section D.4 reserves
for risks and assumptions to validate.

### Step 4: Persist Artifact

**This step is MANDATORY — do NOT skip it.**

Follow **Section C** from `~/.claude/skills/_shared/sdd-phase-common.md`.
- artifact: `design`
- topic_key: `sdd/{change-name}/design`
- type: `architecture`

### Step 5: Return Summary

<!-- matecito-ai: la plantilla literal salió de acá. Mantenerla inline Y en el template crea una
     copia más para desincronizar — que es exactamente el defecto que el template vino a cerrar. -->
The shape of your return lives in **`~/.claude/references/phase-returns/sdd-design.md`**. Read it and
follow it literally: it declares the `## Design Created` block — its sections, their titles, their
order, which ones are unconditional, and what changes when you return `blocked`. The orchestrator
validates your return against that same file, matching titles literally, so a section you drop,
rename or re-level is a gate that never fires.

Three things that file makes explicit and that this phase gets wrong most often:

- `### New Decisions` is this phase's **only Tier-1 mailbox**, and it is emitted **always** — with a
  `None.` sentinel when there is nothing. It carries only the choices that PASS the blocking test in
  `## Rules`; the ones the test catches make you return `blocked` instead, into `### Blocker`. Every
  item carries its `· blocking-test:` token: the orchestrator reads it mechanically, and an omission
  buys the strict reading, not the benefit of the doubt.
- `### Open Questions` is **not** a decision mailbox and **not** Tier 1: it carries what fixes
  nothing. Anything that pins the approach, a contract, a dependency, a boundary or where a
  responsibility lives belongs in `### New Decisions` or in `blocked`.
- On `blocked`, the question and the options go in `### Blocker` — never in `risks`, which Section D.4
  forbids for decisions the user owns.

## Rules

- ALWAYS read the actual codebase before designing — never guess
<!-- matecito-ai: EDRs are binding — respect Accepted EDRs in .matecito-ai/edr/; never contradict one silently (report as blocker); flag uncovered decisions for capture via development-decisions-bootstrap (see Step 2a) -->
- ALWAYS read `.matecito-ai/edr/` (if present) before designing; treat Accepted EDRs as binding constraints and surface conflicts as blockers — including the case where the EDR itself is the inconsistent one, or did not foresee this case (Step 2a)
- ALWAYS emit New Decisions in BOTH places — `## New Decisions` in the artifact and `### New Decisions` **in your return** — whether or not `.matecito-ai/edr/` exists (Step 2a-bis). The guard reads the return only. Only the EDR-naming parts are gated on the store; the detection is not. A repo with no captured decisions is the one that most needs its architectural choices raised, not the one where they pass silently. **What lands in this section is filtered by the blocking test below** — a decision the test catches returns `blocked` instead of being filed here
<!-- matecito-ai: recordatorio apuntado — la doctrina completa vive en el fragmento del dominio (cargado en Step 1), no se duplica acá. Aplica sobre todo a la sección "Interfaces / Contracts". -->
- **Before pinning ANY contract or definition** in the design — domain entity, DB model/migration/schema, DTO, public/exported type, interface or enum, event payload, or config schema (this is what the "Interfaces / Contracts" section materializes) — apply **"Contract & definition shapes — never inferred"** from the domain fragment (`~/.claude/matecito-ai/domains/development.md`, read in Step 1). Never infer which fields it has nor their types. Pinned-and-coherent upstream → carry it through; unspecified, or pinned by something that conflicts or does not cover this case → return `blocked` proposing the FULL contract as one reviewable unit
- Every decision MUST have a rationale (the "why")
- Include concrete file paths, not abstract descriptions
- Use the project's ACTUAL patterns and conventions, not generic best practices
- If you find the codebase uses a pattern different from what you'd recommend, note it but FOLLOW the existing pattern unless the change specifically addresses it
- Keep ASCII diagrams simple — clarity over beauty
<!-- matecito-ai: calificar con "que BLOQUEAN" dejaba fuera todo lo demás, y como el rol del design es ELEGIR el approach, casi ninguna pregunta se auto-califica como bloqueante. El criterio pasa a ser "¿fija una decisión?", no "¿me impide seguir?". -->
<!-- matecito-ai: ARBITRAJE. El criterio anterior ("¿podés proponer una respuesta fundada?") era
     AUTOEVALUACIÓN: sin umbral, cualquier racionalización se aprueba a sí misma, y en una prueba
     funcional el MISMO input (export CSV de 500k filas) admitió las DOS salidas cumpliendo el texto
     al pie de la letra. El corte pasa a ser una propiedad OBSERVABLE de las alternativas, auditable
     por un tercero sin conocer el razonamiento interno del ejecutor. -->
- Do NOT guess. Every open question that **pins a decision** — the approach, a contract, a dependency, a boundary between components, where a responsibility lives — must be raised, whether or not it blocks the rest of the design. Which of the two destinations it takes is NOT decided by how well you can argue for an option. It is decided by the **blocking test**: an objective comparison of the alternatives themselves.

  **The blocking test.** Put the alternatives side by side and check whether they differ in any of these three:
  1. **New infrastructure** — one alternative needs a piece of runtime infrastructure the project does not have today (a queue, a worker or background process, a storage, a cache, a scheduler, a notification channel) and the other does not.
  2. **Public contract** — they expose different API surfaces, event payloads, exported file formats, CLI or config schemas, or a different user-visible interaction (e.g. a synchronous download vs. "we'll notify you when it's ready").
  3. **Data model** — they imply different entities, a different persisted schema, or a migration one of them does not need.

  - **They differ in at least one** → return `blocked` with the question and the options you weighed — **even if you have a founded preference, and even if one option looks obvious**. State explicitly WHICH of the three they differ in and HOW (e.g. "async job needs a queue + worker that do not exist today (1) and turns the export from a download into a notification (2)"). Do not pick one to keep moving; a preference is not a mandate.
  - **They differ in none of the three** → propose the founded choice, file it under `## New Decisions` / `### New Decisions` and keep designing. The guard raises it as Tier 1 and the user confirms or corrects. This is the normal case.

  The test is auditable from outside: a reader who sees only the alternatives must be able to reach the same verdict without reconstructing your reasoning. If you cannot name which of the three differs, you do not have a blocker — you have a preference, and it goes to `## New Decisions`.

<!-- matecito-ai: the test above was pure self-assessment — run in the executor's head, with only the
     verdict published. In the functional test a decision whose OWN text said "this needs a queue and a
     worker the project does not have today" (axis 1, literally) came back filed under `New Decisions`
     with `status: done`, and nothing caught it: catching it would have meant the orchestrator reading
     the decision's prose and re-running the test itself, which is the interpretation the guards exist
     not to do. Same fix as `verify-checks:` for design deviations, one level up: DECLARE the axis, so
     the reader classifies without re-deriving. Also: the requirement right above ("auditable from
     outside") was unsatisfiable until this token gave the verdict somewhere to be seen. -->
- **Declare the verdict — the `blocking-test` token.** Running the test is not enough; it has to be visible. Every item you file under `## New Decisions` / `### New Decisions` carries one token line directly beneath it:

  ```
  · blocking-test: none | infra | contract | data-model
  ```

  `none` means "I put the alternatives side by side and they differ in NONE of the three axes; that is why this item is here and not in `blocked`" — the only value consistent with the item's location, and therefore the normal one. Naming an axis instead contradicts the item's own destination: an axis that differs makes the decision `blocked`. The orchestrator reads the token mechanically and never reopens your reasoning: `none` → ordinary Tier 1; an axis named → it stops, because the item is in the wrong mailbox; absent or hedged → Tier 1 under the strict reading, the same default an undeclared deviation gets in `sdd-apply`. Do not hedge it, do not omit it, and do not write `none` for a decision you did not actually put side by side — one line per decision, and the token IS the audit trail the paragraph above asks for. Shape and the reader's table: `~/.claude/references/phase-returns/sdd-design.md`, section "The blocking-test token".
<!-- matecito-ai: Open Questions dejó de ser buzón de decisiones (se solapaba con New Decisions y el
     ejecutor duplicaba: razonamiento en una, pregunta en la otra => fatiga de confirmación). -->
- `## Open Questions` is NOT a decision mailbox. A question that pins a decision belongs to `## New Decisions` / `### New Decisions`, or makes you return `blocked` per the test above — never here. What stays here is what fixes nothing: implementation doubts and things to validate during apply. And a question you write there and answer yourself in the same delivery is not an open question — it is a decision you took without asking. Whatever remains genuinely open at the end MUST still appear under `### Open Questions` **in your return** (not only in the artifact), so the orchestrator can carry it forward — it is no longer Tier 1
- **Size budget**: Design artifact MUST be under 800 words. Architecture decisions as tables (option | tradeoff | decision). Code snippets only for non-obvious patterns.
- Return envelope per **Section D** from `~/.claude/skills/_shared/sdd-phase-common.md`.
