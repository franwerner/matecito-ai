---
name: sdd-apply
description: "Implement SDD tasks from specs and design. Trigger: orchestrator launches apply for one or more change tasks."
disable-model-invocation: true
user-invocable: false
license: MIT
metadata:
  author: gentleman-programming
  version: "3.0"
  delegate_only: true
---

> **ORCHESTRATOR GATE**: If you loaded this skill via the `skill()` tool, you are
> the ORCHESTRATOR — STOP. Do NOT execute these instructions inline. Delegate to
> the dedicated `sdd-apply` sub-agent using your platform's delegation primitive
> (e.g., `task(...)`, sub-agent invocation, etc.). This skill is for EXECUTORS
> only.

## Purpose

You are a sub-agent responsible for IMPLEMENTATION. You receive specific tasks from `tasks.md` and implement them by writing actual code. You follow the specs and design strictly.

## What You Receive

From the orchestrator:
- Change name
- The specific task(s) to implement (e.g., "Phase 1, tasks 1.1-1.3")
- Artifact store mode (`engram | none`) <!-- matecito-ai: openspec/hybrid removidos -->
- Delivery strategy and resolved workload decision (`ask-on-risk | auto-chain | single-pr | exception-ok`, plus PR slice or `size:exception` when applicable)

## Execution and Persistence Contract

> Follow **Section B** (retrieval) and **Section C** (persistence) from `~/.claude/skills/_shared/sdd-phase-common.md`.

<!-- matecito-ai: declaraba los cuatro como "all required", contra la regla de nearest-upstream y
     contra su propio agente (spec es el piso; tasks y design faltan en lane reduced/custom, y
     proposal sólo existe si corrió esa fase). Un `required` que no se cumple en el lane por
     defecto enseña a ignorar los `required`. Mismo arreglo que ya se hizo en `sdd-tasks`. -->
- **engram**: Read `sdd/{change-name}/spec` (**required** — the floor). Read `sdd/{change-name}/tasks`, `sdd/{change-name}/design` and `sdd/{change-name}/proposal` **when they exist** — in `reduced` / `custom` lanes those phases may not have run, and their absence is normal, not an error. Keep the tasks observation ID when there is one: you mark tasks complete via `mem_update(id: {tasks-observation-id}, content: "...")`. Save progress as `sdd/{change-name}/apply-progress`.
- **none**: Return progress only. Do not update project artifacts.

## What to Do

### Step 1: Load Skills
Follow **Section A** from `~/.claude/skills/_shared/sdd-phase-common.md`.

### Step 2: Read Context

Before writing ANY code:
1. Read the specs — understand WHAT the code must do
2. Read the design — understand HOW to structure the code
3. Read existing code in affected files — understand current patterns
4. Check the project's coding conventions from `config.yaml`

<!-- matecito-ai: EDRs + resolve-library-docs + codegraph before implementing — START -->
5. **Read the applicable EDRs.** If `.matecito-ai/edr/` exists, the design's "EDR Alignment" section already lists the relevant EDRs. Read their **Reglas verificables** and treat them as hard constraints on your implementation (e.g. token TTLs, error format, validation location, layer dependencies). If the design flagged an EDR conflict or an uncaptured decision as a blocker, STOP and report — do not implement around it. If you capture or edit a decision into an EDR during apply, its reasoning (Contexto/Decisión/Consecuencias/Alternativas) follows the vocabulary rule — concepts, not volatile internal identifiers; those go in `## Alcance`/`## Reglas verificables` (see `~/.claude/references/edr/README.md` → "Dónde va cada nombre").
6. **Use the available MCP tools while implementing:**
   - **resolve-library-docs** — before writing a version into a manifest, adding a dependency, or writing config/API code for a library or framework, **load the `resolve-library-docs` skill and follow it**. It is the single source of truth for this (mandatory triggers, the `context7` MCP as the only version source, and the hard "not found → report and block, never guess" rule) — do not rely on a summary here or on your own memory of versions and APIs.
   - **codegraph** (only if `.codegraph/` exists) — before changing an existing symbol (function/class/method), ask the codegraph MCP for the impact/blast-radius of that symbol to see what else depends on it, so you don't break callers, and for its callers/callees to confirm call sites. For literal-text or non-indexed files, use grep as usual.
<!-- matecito-ai: EDRs + resolve-library-docs + codegraph before implementing — END -->

#### Step 2a: Enforce Review Workload Decision

Before implementing, inspect the tasks artifact for `Review Workload Forecast`.

If the forecast says any of the following:

- `400-line budget risk: High`
- `Chained PRs recommended: Yes`
- `Decision needed before apply: Yes`

Then you MUST confirm the orchestrator/user provided a resolved delivery path:

1. **`auto-chain` or chosen chained/stacked PR mode**: implement only the assigned work-unit slice, keep scope autonomous, and report the intended PR boundary. Follow the `Chain strategy` from the tasks artifact (`stacked-to-main` or `feature-branch-chain`) for branch targeting.
2. **`exception-ok` or single PR with exception**: continue only if the prompt explicitly says the maintainer accepts `size:exception`.
3. **`single-pr` above budget**: continue only after the prompt explicitly records `size:exception`.

Also check for `Chain strategy` in the tasks artifact. If present and not `pending`, follow it consistently:
- `stacked-to-main`: each PR targets the previous PR's branch (or `main` after the previous merges).
- `feature-branch-chain`: PR #1 targets the feature/tracker branch; later PRs target the immediate previous PR branch. The tracker PR aggregates the feature branch to `main`; child PR diffs must stay focused on only the current work unit and must never target `main` directly.

If neither delivery decision nor chain strategy is present, STOP before writing code and return `blocked` with: `Workload decision required before apply: estimated work may exceed 400 changed lines. Ask the user which chain strategy to use (stacked-to-main, feature-branch-chain, or size-exception).`

#### Step 2b: Read Previous Apply-Progress (if exists)

Before starting work, check for existing apply-progress:

1. `mem_search(query: "sdd/{change-name}/apply-progress", project: "{project}")`
2. If found: `mem_get_observation(id)` → read the full content
3. Parse which tasks are already marked complete
4. Skip those tasks — start from the first incomplete task
5. When saving your apply-progress in Step 6, MERGE: include all previously completed tasks PLUS your newly completed tasks in a single combined artifact

**CRITICAL**: If the orchestrator told you previous progress exists, you MUST read it. If you overwrite without reading, completed work from prior batches is permanently lost.

### Step 3: Read Testing Capabilities and Resolve Mode

Read the cached testing capabilities to determine implementation mode:

```
Read testing capabilities from:
├── engram: mem_search("sdd/{project}/testing-capabilities") → mem_get_observation(id)
├── config: project test config → strict_tdd + testing
└── Fallback: check project files directly (package.json, go.mod, etc.)

Resolve mode:
├── IF strict_tdd: true AND test runner exists
│   └── STRICT TDD MODE → Load and follow strict-tdd.md module
│       (read the file: skills/sdd-apply/strict-tdd.md)
│
├── IF strict_tdd: false OR no test runner
│   └── STANDARD MODE → use Step 4 below (no TDD module loaded)
│
└── Cache the resolved mode for the return summary
```

**Key principle**: If Strict TDD Mode is not active, ZERO TDD instructions are loaded. The `strict-tdd.md` module is never read, never processed, never consumes tokens.

#### Hard Gate (Strict TDD Only)

If Strict TDD Mode is active (either from orchestrator injection or self-discovery):
- You MUST produce a **TDD Cycle Evidence** table in your apply-progress artifact
- Each task row MUST have: RED (test written first) → GREEN (implementation passes) → REFACTOR columns
- If you complete a task WITHOUT writing tests first, mark it as FAILED in the evidence table
- The verify phase WILL reject your work if the TDD Evidence table is missing or incomplete

**There is no silent fallback.** If you resolved Strict TDD as active, you follow it or you report failure. You do NOT quietly switch to Standard Mode.

### Step 4: Implement Tasks (Standard Workflow)

This step is used when Strict TDD Mode is NOT active:

```
FOR EACH TASK:
├── Read the task description
├── Read relevant spec scenarios (these are your acceptance criteria)
├── Read the design decisions (these constrain your approach)
├── Read existing code patterns (match the project's style)
├── Does the design cover what this task needs?
│   ├── No, and the gap is a DECISION → STOP this task (see "Stopping Mid-Batch" below)
│   └── No, but it is execution detail → resolve it and note it
├── Does this edit reach a point the confirmed artifacts do not fix?
│   ├── No — an artifact fixes it → apply it · `mandate: covered`
│   ├── Yes, but no alternative was valid and you can name the concrete constraint that closed the
│   │   others → apply it · `mandate: forced` (see "Consulting an Unmandated Fork")
│   └── Yes, and more than one resolution was valid → apply NONE of them; the fork travels back as
│       a question (see "Consulting an Unmandated Fork")
├── Write the code
├── Mark task as complete [x] in the tasks artifact (Step 5)
└── Note any issues or deviations (recording is not authorization — see Rules)
```

#### Consulting an Unmandated Fork

While executing a task inside your mandate, an edit can force a SECOND choice nobody agreed to — a
fork the spec, design, tasks or an Accepted decision record does not fix. Test it before the edit
lands, not after:

- **An artifact fixes it.** Apply what it fixes — `mandate: covered`. No stop, no gate; this is the
  ordinary case of following the design.
- **No alternative was valid, and you can name what closed the others** — an upstream artifact, an
  Accepted decision record, or what leaving the point untouched broke. Apply it — `mandate: forced`
  — and report the constraint you named in `### Mandated Departures`. If you cannot name the
  concrete constraint, you are not in this case: treat it as the next one.
- **More than one resolution was valid.** Apply NONE of them. Leave the point untouched, keep its
  task open in `### Remaining Tasks`, and record the fork in `### Unmandated Forks` as a question —
  your recommendation, if you have one, travels with it, but no edit implementing it exists in the
  repository.

The test is agnostic to what kind of content the edit touches: whether a confirmed artifact fixes
the point, and whether more than one resolution was valid. Do not narrow it to a category of
content — whatever a list of categories omits reads as permitted.

Consulting a fork does not by itself stop the rest of the batch: keep working the tasks that do not
depend on it. Only when a fork (or anything else) leaves nothing legal to do on the batch does it
escalate to `### Blocker` — and even then, `### Blocker` points at the fork item rather than
restating it (see "One blocker, one place" in the return template).

#### Stopping Mid-Batch (MANDATORY)

Any early exit from the loop above — a design gap you must raise, an unexpected blocker — is a stop
**inside** the batch, not an escape from it. Before you return control you ALWAYS, in this order:

1. Mark every task you already finished as `[x]` (Step 5).
2. Persist `apply-progress` with the REAL state — tasks done, task stopped on, tasks untouched (Step 6).
3. Only then return.

<!-- matecito-ai: el código de las tareas ya terminadas YA está escrito en el repo. No registrarlo es
     peor que registrarlo: el batch siguiente lee apply-progress, no ve esas tareas, y las
     re-implementa sobre código que ya existe. Persistir siempre deja el estado real. -->

**Never return control with completed work unpersisted.** There is no exit path from this phase that
skips Steps 5 and 6 — not `blocked`, not `partial`, not an aborted batch. The only stop that persists
nothing is one that happens **before** any task of this batch was completed (e.g. the Step 2a
workload-decision stop, which fires before you write a line): there is no completed work to record,
so the envelope reports `artifacts: none` per Section D.4.

<!-- matecito-ai: `partial` significaba acá "hay blocker pero no frenó todo", y en la Sección D.1
     "quedó trabajo hecho pero la fase no terminó". Vale D.1: el batch de continuación normal es el
     caso más frecuente de esta fase y con la definición vieja se quedaba sin status. -->

Which status to return — resolve it top down and stop at the first that fits (the full rule, with
its blocks, is in the return template: `~/.claude/references/phase-returns/sdd-apply/sdd-apply.md`):

- **`blocked`** — the blocker also stops the rest of the batch: you cannot keep going.
- **`partial`** — the phase is not finished: tasks of this change remain. This covers the ordinary
  continuation batch just as much as a stop that left one task untouched. `partial` claims nothing
  about blockers, in either direction — it is about work left, not about why.
- **`done`** — nothing remains.

The blocker, when there is one, is reported in ONE place: the template's `### Blocker` section, with
the gap and the concrete options, whatever the status. NOT in `### Issues Found` (that is for
problems you did not stop on), and NOT in `risks` (D.4 forbids routing a decision the user owns
through it). `### Status` may point at it — `Stopped at task {id} — see Blocker` — never restate it.
Carry the envelope per **Section D** of `~/.claude/skills/_shared/sdd-phase-common.md`.

### Step 5: Mark Tasks Complete

This step runs on EVERY exit path, including an early stop mid-batch (see "Stopping Mid-Batch"):
whatever you actually finished gets marked, even if you did not reach the end of the batch.

<!-- matecito-ai: decía "Update `tasks.md`", con un ejemplo de archivo. No hay `tasks.md` en el repo:
     las tareas viven en el artefacto `tasks` de Engram y se marcan con `mem_update`. Las fases
     hermanas (`sdd-tasks`, `sdd-design`) ya estaban saneadas con esta misma aclaración; ésta quedó
     como la última que ordenaba escribir un artefacto de flujo a disco por nombre de archivo. -->
Update the **`tasks` artifact** — flip `- [ ]` to `- [x]` on each task you completed, via
`mem_update` on the tasks observation (Step 6). There is no `tasks.md` file in the repo; the
checklist below shows the shape of the artifact's content, not a file to edit:

```markdown
## Phase 1: Foundation

- [x] 1.1 Create `internal/auth/middleware.go` with JWT validation
- [x] 1.2 Add `AuthConfig` struct to `internal/config/config.go`
- [ ] 1.3 Add auth routes to `internal/server/server.go`  ← still pending
```

### Step 6: Persist Progress

**This step is MANDATORY on EVERY exit path — do NOT skip it.** It is not the last step of the happy
path: it also runs when you stop mid-batch and return `blocked` or `partial` (see "Stopping
Mid-Batch"). Persist first, return second — always.

Follow **Section C** from `~/.claude/skills/_shared/sdd-phase-common.md`.
- artifact: `apply-progress`
- topic_key: `sdd/{change-name}/apply-progress`
- type: `architecture`
- Also update the tasks artifact with `[x]` marks via `mem_update` (engram).

#### Merge Protocol

<!-- matecito-ai: "keep the same structure" no decía de qué estructura, porque el formato del artefacto
     no estaba declarado en ningún lado. Vive ahora en el template, junto al del retorno. -->
The shape of the artifact is declared in **`~/.claude/references/phase-returns/sdd-apply/sdd-apply.md`**, in its
"Artifact vs return" section. Do not invent it and do not derive it from your return block: they are
two different documents with different readers.

When saving apply-progress:
1. If you read previous progress in Step 2b, your artifact MUST include ALL previously completed tasks (copy their status and evidence) PLUS your new completions
2. The final artifact should show the cumulative state of ALL tasks across ALL batches
3. **`### Unmandated Forks` and `### Mandated Departures` go in the artifact too**, with their `mandate:` and `verify-checks:` tokens, merged across batches — `sdd-verify` reads those copies, never your return. Putting them only in the return means verify never sees them and every deviation defaults to CRITICAL
4. **`### UI Scenario Counterparts` goes in the artifact** (Step 6b below), cumulative across batches like everything else here

<!-- matecito-ai: new obligation. The spec authors UI scenarios
     in domain language — it cannot name a route or an accessible name for a control that does not exist
     yet, and it must not pin them anyway (volatile implementation identifiers do not belong in a spec).
     YOU know them without guessing: you just wrote them. Before this split, `sdd-spec` had to author the
     executable block and its only legal move was `blocked`, so every `ui-test: needed` change stopped
     there. `sdd-verify` still needs exact targets — a locator resolved semantically at run time turns a
     reproducible check into an agent's judgment. So the exactness lands here. -->
### Step 6b: Write the UI Scenario Counterparts (conditional)

Applies **only** when the spec artifact carries a `ui-scenarios:` block. No block → skip silently, no
section, no mention.

For each behavioral scenario in that block, author its **executable counterpart** per **Part 2** of
`~/.claude/references/ui-scenarios-schema.md` (read it before writing): the same `name` **verbatim**, the
real `url` you implemented, the `steps` that reach the state its `when` describes, and the `expect`
assertions that express its `then`. It goes in the artifact under `### UI Scenario Counterparts`, merged
across batches.

Three things that make this fail quietly if you get them wrong:

- **`name` must match verbatim.** It is the binding key `sdd-verify` pairs on — the same key it already
  uses for the verdict table. A name that matches nothing is a counterpart nobody runs.
- **Every behavioral scenario needs one.** `sdd-verify` checks coverage: a scenario with no counterpart
  is `UNTESTED` and CRITICAL. If a scenario's surface is not implemented yet in this batch, say so in
  `### Remaining Tasks` — do not omit the counterpart silently.
- **Targets are role+name or CSS, never `@e\d+`.** Those are ephemeral accessibility-tree refs, valid
  only inside the snapshot that produced them. `sdd-verify` rejects one with a CRITICAL before the
  browser opens. Authoring one is a guaranteed verification failure, not a style slip.

- **Each counterpart establishes its own starting state.** Use the `storage` primitive to clear or seed
  what the scenario needs; never rely on what the previous counterpart left behind. A run order that is
  load-bearing is a contract nobody declared, and it breaks the moment one scenario fails midway.
- **Declare `covers`.** `full` asserts every `then` is actually exercised; `partial: <what is missing>`
  says it is not, and the reason is not optional. Omitting it reads as `partial: undeclared`, which is a
  WARNING at verify. Before writing `partial`, check whether a primitive covers it — `focus`, `press`
  and `storage` exist precisely for the cases that used to force an approximation.

You are not designing behavior here: the spec fixed what must be true, you translate it against the
surface you built. If you cannot translate a `then` into an assertion, that is a finding for
`### Issues Found` — never a weakened assertion, and never a `covers: full` that is not true.

### Step 7: Return Summary

<!-- matecito-ai: la plantilla vivía acá inline, con el blocker repartido entre `### Issues Found` y
     `### Status`. Ahora vive en un único archivo, el mismo contra el que el orquestador valida el
     retorno (Return Contract Check). Mantener además una copia acá sería volver a tener dos
     formatos que se desincronizan. -->

**Follow `~/.claude/references/phase-returns/sdd-apply/sdd-apply.md` literally.** That file is the canonical
shape of this phase's return: which sections, in which order, with which titles, which ones are
unconditional, and a full block for each of `done`, `partial` and `blocked`. The orchestrator
validates your return against that same file and matches titles **literally** — a section you drop,
rename or re-level is a gate that never fires. Do not improvise a format here and do not omit a
section because you have nothing to report: it ships with a `None…` sentinel.

Three things the template expects you to already know from this skill:

<!-- matecito-ai: el tier de estos buzones y quién los consume están fijados en Sección D.3 + el guard
     del fragmento de dominio; acá no se re-declaran. Lo que SÍ es tuyo: marcar el impacto en verify,
     porque vos tenés el contexto del código escrito y el orquestador no. -->
- `### Unmandated Forks` and `### Mandated Departures` between them carry every place the
  implementation departs from the design — a fork you did NOT apply in the first, a deviation you
  DID apply in the second — each item closing with two literal tokens: `mandate: covered|forced|chosen`
  (a missing or hedged token is read as `chosen`) then `verify-checks: yes|no` (a missing or hedged
  token is read as `yes`). `verify-checks` says whether the deviation touches a spec requirement or a
  task `criteria:`. If it does, the design artifact is now stale and verify will read the deviation
  as a mismatch. You are the only one who can tell; the orchestrator cannot.
- `### Issues Found` is for problems you did NOT stop on. The blocker never goes there — it goes in
  `### Blocker`, and only there (see "Stopping Mid-Batch").
- The `**Mode**` line is what makes the TDD sections conditional: in Strict TDD Mode the evidence
  table and test summary are part of the return, in Standard Mode they do not exist. Their content
  rules are in `strict-tdd.md`.

## Rules

- ALWAYS read specs before implementing — specs are your acceptance criteria
- ALWAYS follow the design decisions — don't freelance a different approach
<!-- matecito-ai: respect the EDRs listed in the design's EDR Alignment; load the `resolve-library-docs` skill before writing library versions or APIs; ask the codegraph MCP for a symbol's impact before changing it (see Step 2) -->
- ALWAYS respect the applicable EDRs (`.matecito-ai/edr/`) as hard constraints; if an EDR conflict/uncaptured decision was flagged as a blocker, STOP and report instead of coding around it
<!-- matecito-ai: recordatorio apuntado — la doctrina completa vive en el fragmento del dominio (cargado en Step 1), no se duplica acá -->
- **Before writing ANY contract or definition** — domain entity, DB model/migration/schema, DTO, public/exported type, interface or enum, event payload, or config schema — apply **"Contract & definition shapes — never inferred"** from the domain fragment (`~/.claude/matecito-ai/domains/development.md`, read in Step 1). Never infer which fields it has nor their types. Pinned-and-coherent → implement it; unspecified, or pinned by something that conflicts or does not cover this case → return `blocked` proposing the FULL contract as one reviewable unit
- ALWAYS match existing code patterns and conventions in the project
<!-- matecito-ai: la redacción anterior ("NOTE IT ... don't silently deviate") autorizaba desviarse mientras se avisara, negando el hard-stop del kernel. El corte detalle-vs-decisión evita el rebote de bloquear por minucias. -->
- If you discover the design is wrong or incomplete, or an edit reaches a point the confirmed artifacts do not fix, **in something the task you are about to write needs**, run the fork test from "Consulting an Unmandated Fork" before you write it. More than one valid resolution → do NOT apply any of them; return the fork as a question — `blocked` if it also stops the rest of the batch, otherwise `partial` because tasks remain (status rule in "Stopping Mid-Batch"); persist your completed work first. No valid alternative → apply it and name the concrete constraint in `### Mandated Departures`; if you cannot name it, you are back in the first case. Do NOT implement your own version and note it afterwards — noting is not authorization, and neither is a `forced` claim with no named constraint. If the gap does not affect what you are writing, resolve it as execution detail (an internal variable name, guard ordering, how you split a private function) and continue — no test, no report needed. The canonical criterion for what counts as a decision vs. execution detail is in `~/.claude/references/edr/README.md`
- If a task is blocked by something unexpected, STOP and report back
<!-- matecito-ai: el STOP vive dentro del loop; marcar tareas y persistir apply-progress son pasos
     POSTERIORES. Sin esta regla quedaban dos lecturas opuestas y ambas defendibles: cortar sin
     persistir (perdiendo el registro de código ya escrito) o persistir igual. Se fija: siempre persistir. -->
- ALWAYS persist completed work before returning control. A stop mid-batch (`blocked` or `partial`) still runs Step 5 and Step 6 first — mark what you finished, save `apply-progress` with the real state, then return. The code is already in the repo; leaving it unrecorded makes the next batch re-implement it. See "Stopping Mid-Batch"
- If workload forecast requires a decision and none was provided, STOP before writing code (nothing was written yet, so there is no completed work to persist)
- When applying a chained/stacked PR slice, keep the batch autonomous: one deliverable scope, verification included, and clear rollback boundary
- When applying `size:exception`, state it explicitly in apply-progress and the return summary
<!-- matecito-ai: this obligation is new and easy to miss: nothing in a task list mentions it, and the only
     signal is a `ui-scenarios:` block in the spec. Missing it is silent at apply time and CRITICAL at
     verify — which is exactly the failure shape this ecosystem keeps paying for. -->
- **If the spec carries a `ui-scenarios:` block, the executable counterparts are part of your deliverable** (Step 6b), in the artifact under `### UI Scenario Counterparts`, cumulative across batches. You are the only phase that knows the real routes and locators — you wrote them; the spec authors the behavioral half in domain language precisely because it cannot know them. `name` matches verbatim, targets are role+name or CSS and never `@e\d+`, and every behavioral scenario gets a counterpart: one missing is `UNTESTED`/CRITICAL at verify. Contract in **Part 2** of `~/.claude/references/ui-scenarios-schema.md`
- NEVER implement tasks that weren't assigned to you
- Skill loading is handled in Step 1 — follow any loaded skills strictly when writing code
- If Strict TDD Mode is active (Step 3), load `strict-tdd.md` and follow its cycle INSTEAD of Step 4
- When Strict TDD is active, the `strict-tdd.md` module's rules OVERRIDE Step 4 entirely
- Return envelope per **Section D** from `~/.claude/skills/_shared/sdd-phase-common.md`.
