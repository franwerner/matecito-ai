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

> Follow **Section B** (retrieval) and **Section C** (persistence) from `~/.claude/skills/_shared/sdd-phase-common.md`.

<!-- matecito-ai: dos lecturas con propósitos distintos. La de requisitos es "el upstream más cercano"
     (proposal, o el brief en lane reduced). La del brief es INCONDICIONAL y sólo por el flag
     `ui-test`, que ninguna otra fase upstream transporta. No las colapses en una. -->
- **engram**: Read the nearest available upstream for requirements — `sdd/{change-name}/proposal`, falling back to `sdd/{change-name}/intake` when no proposal exists (`reduced`/`custom` lanes). Read `sdd/{change-name}/intake` **as well, always**, for the `ui-test` flag (Step 4b) — even when the proposal was your requirements source. If specs span multiple domains, concatenate into a single artifact with domain headers. Save as `sdd/{change-name}/spec`.
- **none**: Return result only. Never create or modify project files.

## What to Do

### Step 1: Load Skills
Follow **Section A** from `~/.claude/skills/_shared/sdd-phase-common.md`.

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

### Step 4b: Write the `ui-scenarios` Block (conditional — only when `ui-test: needed`)

<!-- matecito-ai: esta fase PRODUCE `ui-scenarios`. Antes nadie lo escribía: sdd-verify lo iteraba y
     su compuerta de UI cerraba en silencio al no encontrarlo, así que la verificación de UI entera
     era inalcanzable sin ningún error visible. El productor es spec porque el bloque es la
     contraparte ejecutable de los escenarios Given/When/Then que esta fase ya escribió. -->

<!-- matecito-ai: el flag se lee SIEMPRE del intake brief, NO del upstream de requisitos. En lane
     `full` el upstream de esta fase es la proposal, que no lleva el flag; sólo en lane `reduced` el
     upstream ES el brief. Como intake es fase BASE, el brief existe en todos los lanes: es el único
     lugar donde el flag está garantizado. Leerlo "del upstream" lo perdería en el lane completo. -->

**Read the intake brief for this flag — ALWAYS, whatever your requirements upstream was.**
`mem_search("sdd/{change-name}/intake")` → `mem_get_observation`, and look under `### Classification`
for the line `- UI test: {needed|not-needed}`. Do this even when you took your requirements from the
proposal: the proposal does not carry the flag, and intake is a base phase, so the brief always
exists. The flag is decided by `sdd-intake` and confirmed by the user at the INTAKE GATE — you never
decide it and never override it.

```
IF the brief says `UI test: needed`
├── AUTHOR the block (below) into the spec artifact
└── report the scenario count in your return's `### Coverage`

ELSE (not-needed, or the line is absent)
└── Skip SILENTLY — no block, no placeholder, no mention anywhere
```

**Before authoring, read `~/.claude/references/ui-scenarios-schema.md` in full.** It is the contract:
per-scenario fields, the legal step primitives, the `wait` primitive, the target rules, the assertion
classes, and the validation rules your output must already satisfy. Do not author from memory, and do
not invent a primitive it does not define — `sdd-verify` executes the block literally.

**Derive, do not duplicate.** These are not a second set of scenarios: take the Given/When/Then
scenarios you already wrote and, for the capabilities with a visual surface, express each as a
behavioral UI scenario — the GIVEN becomes `given`, the WHEN becomes `when`, the THEN becomes the
`then` conditions. A capability with no visual surface produces no `ui-scenarios` entry. Never
introduce behavior in this block that the requirements above it do not already state.

<!-- matecito-ai: this step used to author the FULLY EXECUTABLE block — `url` plus steps with concrete
     locators. It was unworkable and wrong at once. Unworkable: the schema demanded a route and role+name
     targets, the brief carries neither, and for a new feature the control does not exist yet, so the only
     legal move was `blocked` — every `ui-test: needed` change stopped here. Wrong: a route and an
     accessible name are volatile implementation identifiers, which the capability-spec vocabulary rule
     forbids in a spec. You own the WHAT; `sdd-apply` authors the executable counterpart because it knows
     the real locators — it just wrote them. -->
**You write the WHAT, not the HOW.** No routes, no CSS selectors, no accessible names, no component
names. If you are about to write `/login` or `role=button name="Submit"`, stop: that is the executable
counterpart, and `sdd-apply` authors it in its `apply-progress` (Part 2 of the schema). You never need
to know a route or a control's name to say what must be true — and you must not pin them, because they
are the code's business.

This is also why you no longer block here. Nothing in this block requires information that does not
exist yet, so there is no case where the flag says `needed` and you cannot author it.

Name what the user sees and does, in the vocabulary the requirement already uses. If a scenario does
not name the control it acts on in user-visible terms, that is a gap in the requirement — fix the
requirement, do not invent a selector to paper over it.

**Where it goes**: at the END of the spec artifact you persist in Step 5 — never in your return
block. The orchestrator never reads the artifact, and `sdd-verify` never reads your return.

````markdown
## UI Scenarios

```yaml
ui-scenarios:
  - name: {the binding key — sdd-apply reuses it verbatim, sdd-verify maps it 1:1 to a STATE row}
    given: {the starting situation, in domain language}
    when: {the interaction the user performs}
    then:
      - {an observable condition that must hold afterwards}
```
````

<!-- matecito-ai: acá había un marcador `ui-test: needed` copiado dentro del spec, de cuando verify
     leía el flag del artefacto spec. Verify ahora lo lee del intake brief, que es donde se decide y
     se confirma: copiarlo acá era duplicar estado, y dos copias del mismo dato divergen. -->
**Do NOT copy the `ui-test` flag into the spec.** It lives in the intake brief, which is where it is
decided and confirmed, and that is where `sdd-verify` reads it. The spec carries the scenarios; the
brief carries the flag. One fact, one home.

### Step 4c: Recognize decisions — only when no `design` add-on is active

<!-- matecito-ai: in-flow decision capture (development-specifics). Full mechanism, the ratification
     gate table, and why the title matches sdd-design's byte-for-byte: in-flow-capture.md. -->

**Read `decisions_gate_here` first.** `mem_search("sdd/{change-name}/intake")` → `mem_get_observation`
(you already read this artifact in Step 4b for the `ui-test` flag — reuse it), and look at
`### Triage`: `Lane: {...} — add-ons: [...]`. `decisions_gate_here` is `true` **iff** `design` is NOT
one of the listed add-ons. `false` → skip this whole step silently: `sdd-design` will run later and its
own `### New Decisions` is the single ratification gate for this change; a second one here would ask
the user to ratify the same decision twice.

When `true`: while writing the delta spec, you MAY notice a genuine architectural decision — the
approach, a contract, a dependency, a boundary between components, where a responsibility lives — that
no confirmed upstream artifact already fixes. You are NOT designing (you still do not choose HOW the
code is structured), but writing WHAT sometimes surfaces exactly this kind of choice, and in this lane
there is no later phase to catch it. Apply the **same blocking test** `sdd-design`'s `## Rules` defines
(new infrastructure / public contract / data model):

- **The alternatives differ in at least one axis** → return `blocked`, same shape as any other blocker
  in this phase's `### Blocker` section — name which axis and how, exactly as `sdd-design` would.
- **They differ in none** → propose it under `### New Decisions` in your return (see
  `~/.claude/references/phase-returns/sdd-spec/sdd-spec.md`), with the same item shape `sdd-design`
  uses: `summary`, `· blocking-test: none`, `· record: <domain>/<slug>`, `· rationale:`.

If you notice nothing decision-shaped while writing the spec — the ordinary case for most changes —
emit the section anyway, with `None.`; do not invent a decision to fill it.

### Step 5: Persist Artifact

**This step is MANDATORY — do NOT skip it.**

Follow **Section C** from `~/.claude/skills/_shared/sdd-phase-common.md`.
- artifact: `spec`
- topic_key: `sdd/{change-name}/spec`
- type: `architecture`

### Step 6: Return Summary

<!-- matecito-ai: la plantilla literal salió de acá. Mantenerla inline Y en el template crea una
     copia más para desincronizar — que es exactamente el defecto que el template vino a cerrar. -->
The shape of your return lives in **`~/.claude/references/phase-returns/sdd-spec/sdd-spec.md`**. Read it and
follow it literally: it declares the `## Specs Created` block — its sections, their titles, their
order, which ones are unconditional, and what changes when you return `blocked`. The orchestrator
validates your return against that same file, matching titles literally, so a section you drop,
rename or re-level is a gate that never fires.

Two things that file makes explicit and that this phase gets wrong most often:

- `### Derived capabilities (unconfirmed)` is this phase's **Tier-1 mailbox**. It carries the
  mapping you derived in Step 2 and it is emitted **always** — with the `None — mapping was
  explicit.` sentinel when the upstream artifact had its own Capabilities section. A derived mapping
  is not contract until the main thread confirms it.
- The ambiguous-derivation stop from Step 2 returns `blocked`, and the possible readings go in the
  `### Blocker` section that file designates — never in `risks`, never as a mapping you picked.
- `### New Decisions` (Step 4c) is a **second, conditional** Tier-1 mailbox — emitted only when
  `decisions_gate_here` is true, omitted entirely otherwise. Do not confuse it with `### Derived
  capabilities (unconfirmed)`: a capability mapping and an architecture-decision proposal are
  different things, each in its own section.
<!-- matecito-ai: el bloque `ui-scenarios` va en el ARTEFACTO, no en el retorno; del retorno sólo
     cuelga el conteo, y en una sección que ya existe — no se abre un buzón nuevo por esto. -->
- The `ui-scenarios` you wrote in Step 4b are reported in the return **only as a count**, on the
  `### Coverage` line that file declares for them. The block itself lives in the persisted artifact.

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
- **In a lane with no `design` add-on, recognize architecture decisions you notice while writing the delta spec** (Step 4c) — apply `sdd-design`'s blocking test; a decision that fails it returns `blocked`, one that passes it goes to `### New Decisions` with the same `· blocking-test:` / `· record: <domain>/<slug>` tokens `sdd-design` uses. In a lane WITH `design`, skip this entirely — `sdd-design`'s own mailbox is the single ratification gate. Full mechanism: `~/.claude/references/decision-capture/in-flow-capture.md`
- **Before a scenario or requirement pins ANY contract or definition** — domain entity, DB model/migration/schema, DTO, public/exported type, interface or enum, event payload, or config schema — apply **"Contract & definition shapes — never inferred"** from the domain fragment (`~/.claude/matecito-ai/domains/development.md`, read in Step 1). Never invent which fields it has nor their types to make a scenario concrete. Pinned-and-coherent upstream → use it; unspecified, or pinned by something that conflicts or does not cover this case → return `blocked` proposing the FULL contract as one reviewable unit
<!-- matecito-ai: reglas del bloque `ui-scenarios` — el detalle operativo está en el Step 4b; acá van
     las que un ejecutor apurado rompe si sólo lee las Rules. -->
- **The BEHAVIORAL half of `ui-scenarios` is produced HERE, and only when the intake brief says `ui-test: needed`** — read that flag from `sdd/{change-name}/intake` **always**, whatever your requirements upstream was (the proposal does not carry it; the brief always exists). Flag `not-needed` or absent → skip silently: no block, no mention. The **executable counterpart** — real route, real locators — is `sdd-apply`'s, in its `apply-progress`: you never author it, and you never need a route or a control's name to say what must be true
- **`ui-scenarios` targets MUST be role+name or CSS — NEVER `@eN` runtime snapshot refs.** `sdd-verify` rejects a `@e\d+` target as CRITICAL and the scenario fails static validation
- **`ui-scenarios` entries are DERIVED from the Given/When/Then scenarios you already wrote** (for the capabilities with a visual surface), never a parallel set of behavior. Author them against `~/.claude/references/ui-scenarios-schema.md`, read in full first
- **MODIFIED requirements MUST be the FULL block** — copy entire requirement + all scenarios from main spec, then edit. Partial MODIFIED blocks lose content at archive time.
- If adding new behavior without changing existing behavior → use ADDED, not MODIFIED
<!-- matecito-ai: el presupuesto acota la PROSA de requisitos. Contar el YAML de `ui-scenarios` contra
     él haría que un ejecutor recorte escenarios de UI para "entrar" — degradando en silencio lo que
     sdd-verify va a ejecutar. El bloque ya está acotado por su origen: sale de los escenarios que
     escribiste, para capabilities con superficie visual. -->
<!-- matecito-ai: `## Scope` se exceptúa por la MISMA razón, y no en abstracto: pasó. Un ejecutor con
     el presupuesto ajustado la fue plegando y terminó dropeándola entera para entrar — y `sdd-apply`
     NO lee el intake brief, así que ese spec era su única fuente de qué archivos tocar. Se tapó
     pasándole las rutas en el prompt de despacho, que es el orquestador cubriendo un agujero del
     contrato. La lista de archivos no es prosa que se pueda comprimir: o está completa o miente.
     Y el conteo: `wc -w` cuenta cada guion de viñeta y cada flecha como palabra — reportó 724 sobre
     643 reales, así que ese ejecutor plegó contenido por ~80 palabras que no existían. -->
<!-- matecito-ai: el cap numérico (650 palabras) se quitó. Estaba mal calibrado y el síntoma lo probó: una
     spec que el usuario ratificó salía en ~665 ANTES de cualquier enmienda, así que el tope no medía
     "spec inflada", medía "spec normal" — y lo único que producía era plegar contenido para entrar.
     Eso es exactamente la falla que las dos excepciones de abajo existen para frenar, sólo que
     llegando por la puerta del número en vez de por la del ejecutor apurado. La presión hacia tablas
     y escenarios cortos se queda; el número no. -->
- **Size budget**: Keep the artifact tight — prefer requirement tables over narrative descriptions, and hold each scenario to 3-5 lines. **There is no word cap**: the length follows the change's real surface, and folding content to hit a number is the failure this budget exists to prevent, not its purpose. **Never drop or trim the `## UI Scenarios` block or `## Scope`** to make the artifact shorter — the file list is `sdd-apply`'s only source of what to touch, since it does not read the intake brief.
- Return envelope per **Section D** from `~/.claude/skills/_shared/sdd-phase-common.md`.

## RFC 2119 Keywords Quick Reference

| Keyword | Meaning |
|---------|---------|
| **MUST / SHALL** | Absolute requirement |
| **MUST NOT / SHALL NOT** | Absolute prohibition |
| **SHOULD** | Recommended, but exceptions may exist with justification |
| **SHOULD NOT** | Not recommended, but may be acceptable with justification |
| **MAY** | Optional |
