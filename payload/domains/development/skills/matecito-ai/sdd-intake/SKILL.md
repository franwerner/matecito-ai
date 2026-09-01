---
name: sdd-intake
description: "Intake and structure a raw user request before the SDD flow. Trigger: orchestrator launches intake, or a user describes a feature/bug/change in natural language that needs structuring before exploration."
disable-model-invocation: true
user-invocable: false
license: MIT
metadata:
  author: matecito-ai
  version: "2.0"
  delegate_only: true
---

> **ORCHESTRATOR GATE**: If you loaded this skill via the `skill()` tool, you are the
> ORCHESTRATOR — STOP. Do NOT execute these instructions inline. Delegate to the dedicated
> `sdd-intake` sub-agent. This skill is for EXECUTORS only.

## Purpose

You are a sub-agent responsible for **INTAKE** — the first phase of the SDD flow. You take a raw,
natural-language request from the user (as typed in the chat) and turn it into a **structured brief**
that the rest of the flow can act on, in a single pass. You also catch EDR conflicts or undecided
architectural questions *before* exploration burns effort.

You do NOT explore the codebase in depth (that is `sdd-explore`). You do NOT design or implement, and
you do NOT ask the user anything — you run headless. This phase always runs the `full` pipeline;
`direct` is decided outside the flow, before any phase is dispatched, and this phase is never invoked
for it.

## What You Receive

- A raw request from the user, in natural language (e.g. "quiero que se puedan exportar los reportes a CSV").
- Artifact store mode (`engram | none`).
- On a correction re-dispatch from the Brief Confirmation Gate: the original request plus the user's
  correction, verbatim. This is still a single, fresh pass — you produce the brief exactly as you
  would from a first dispatch, and it upserts the same `topic_key` (Step 4), overwriting the prior
  version.

## Execution and Persistence Contract

> Follow **Section B** (retrieval) and **Section C** (persistence) from `~/.claude/skills/_shared/sdd-phase-common.md`.

- **engram**: Save artifact as `sdd/{change-name}/intake`.
- **none**: Return result inline only.

## What to Do

### Step 1: Load Skills
Follow **Section A** from `~/.claude/skills/_shared/sdd-phase-common.md`.

### Step 2: Classify the Change

From the raw request, classify:
- **Type:** `feature` | `bug` | `refactor` | `chore`
- **Domains touched:** map to the canonical EDR domains (e.g. an export endpoint touches `contracts`, `security`, `runtime`, maybe `data`). This is a rough mapping to help routing — NOT a deep analysis.

<!-- matecito-ai: neither flag appeared anywhere in this skill, which declares itself the authority on
     CONTENT ("This skill defines WHAT goes in each section") and enumerated Classification as "type,
     domains touched". They lived only in the agent file. This phase is the ONLY one that decides
     them and two downstream phases read them from the brief: drop them from the classification and they
     drop from the brief, and their readers find absence — which both gates read as "does not apply",
     silently. -->
Plus the three **downstream flags** this phase is the only one to decide. They are part of the
classification and they travel in the brief (`### Classification`). None is executed here — this
phase decides, others act:

- **`diagram`:** `needed` | `not-needed`, with a one-line reason. Apply the **diagram inference test**
  in `~/.claude/matecito-ai/domains/development.md` (`## Architecture diagrams (drawio)`) — that is
  the single source of truth for when a diagram is warranted; do not restate or re-derive its
  criteria. Read by `sdd-design`, which only *recommends* the live render. You never generate one.
- **`ui-test`:** `needed` | `not-needed`, with a one-line reason. Scan the request's description and
  scenarios for `browser`, `page`, `form`, `screen`, `visual`, `click`, `render` — any hit → `needed`.
  An explicit override written in the request (`ui-test: needed` / `ui-test: not-needed`) beats the
  keyword inference; no hit and no override → `not-needed`. Read by `sdd-spec` (it authors the
  `ui-scenarios` block only when this says `needed`) and by `sdd-verify` (its UI gate). You never run
  proofshot.
- **`worktree-isolation`:** `active` | `inactive`, with a one-line reason, per
  `structure/change-isolation-activation-flag.md`. `active` only when the request explicitly asks for
  isolated work; otherwise `inactive`. Do not derive it from anything else about the request — size,
  complexity, or concurrency are not signals for this flag; an explicit ask is the only one. Read by
  the orchestrator (kernel's `### Change Workspace (opt-in)`), never by a later phase agent — you
  never open a workspace yourself.

Absence is not neutral: all three downstream readers read a missing flag as "does not apply" /
inactive and close **silently**, so a flag you drop is a check nobody notices was skipped.

<!-- matecito-ai: a THIRD classification value, but not a third "downstream flag" — it has no phase
     reader. Presence-based on `repo.components`, same gate family as EDRs and capability-specs. -->
Plus **`components`**, multivalued — unlike `diagram` and `ui-test`, **no phase reads it**: it is
metadata reported alongside the rest of the brief's flags, nothing more. It exists only when the
project's `repo.components` set is declared (**presence-based gate**, per
`~/.claude/references/repo-components/README.md`): with no set declared in the project config, do
not emit the bullet and do not mention components anywhere in the brief. With the set declared:

- **Infer** it by mapping the request's scope against `repo.components[].paths`. Emit the line
  **always** when the axis is active — an absent `Components` line with the set declared is an
  anomaly, the same weight as a missing `Diagram`, never "the axis doesn't apply here".
- **Multivalued**: name every component whose `paths` cover the request's scope, comma-separated,
  using each component's exact `name`. Never a path, a folder, a path prefix, or an invented label.
- **No match → `unassigned`**: if no declared `paths` cover the request's scope, emit the line with
  the single token `unassigned`. Do NOT omit the line to mean "nothing matched" — that omission means
  the axis is off, and conflating the two erases the difference between "this repo doesn't use the
  axis" and "this change didn't land on any declared surface".
- **No later phase re-infers or re-decides it.** It is reported once, with the rest of the brief's
  flags, and nothing waits on it.

### Step 3: Early Guard — EDR conflicts and undecided questions

First apply the **EDR activation gate** (single source of truth in `matecito-ai:behavior`): if `.matecito-ai/edr/` is absent or empty, EDRs are inactive — **skip this entire step silently** (`status: done`, no mention of EDRs in the brief). Only when active, continue.

Read `.matecito-ai/edr/INDEX.md` and the indexes of the domains this request touches.
This is a **shallow** check — you are looking for early blockers, not doing design:

- **Conflict:** does the request contradict an `Accepted` EDR? (e.g. "endpoint público sin login" vs an auth EDR that requires protection.) → set `status: blocked`, name the EDR, and recommend resolving via `development-decisions-bootstrap` (update) or adjusting the request. Do NOT proceed to recommend the flow.
- **Undecided question:** does the request require an architectural decision that NO EDR covers? (e.g. export of huge files — sync or background job? no EDR says.) → set `status: needs-decision`, name the gap, and recommend `development-decisions-bootstrap` to capture it *before* the flow runs.
- **All clear** → `status: done`, next is `sdd-explore`.

The point: catch the blocker now, at intake, instead of letting the flow discover it at the design phase after wasting explore/propose/spec.

### Step 4: Persist Artifact

Follow **Section C** from `~/.claude/skills/_shared/sdd-phase-common.md`.
- artifact: `intake`
- topic_key: `sdd/{change-name}/intake`
- type: `architecture`

### Step 5: Return

<!-- matecito-ai: la plantilla literal vivía acá y cubría bien la pasada feliz, dejando a
     interpretación `blocked` y `needs-decision`: dos ejecutores inventaron secciones distintas para
     el mismo caso. La FORMA se mudó al template; acá queda el CONTENIDO. No la vuelvas a copiar:
     una segunda copia es una desincronización esperando. -->
**The shape of your return lives in `~/.claude/references/phase-returns/sdd-intake/sdd-intake.md`.** Read it
and follow it **literally**: it declares the block, its sections, their order, which ones are
unconditional, and what changes for each of the three statuses this phase can return (`done`,
`needs-decision`, `blocked`). The orchestrator validates your return against that same file, matching
titles literally — a section you drop, rename or re-level is a gate that never fires. Do NOT
reconstruct the format from memory or from another phase's return.

- **Request (structured)** — 1-2 sentences: what the user wants, restated clearly.
- **Classification** — Step 2's output: type, domains touched, plus the three downstream flags
  `diagram`, `ui-test` and `worktree-isolation`, each with its one-line reason. The flags are not optional
  extras: they exist nowhere else, and the readers that act on them close silently when they are
  absent. Plus **`Components`**, multivalued and reader-less (Step 2) — emit it only when
  `repo.components` is declared for this project; when you render the return block (Step 5's tool), supply the boolean
  that gates this bullet (`components_axis_active`) explicitly: `true` with the line rendered, `false`
  when the set is not declared, never omitted — an omitted gate boolean fails the render, it does not
  read as "off".
- **Early guard (EDRs)** — Step 3's finding: all-clear, the conflict, or the undecided question, with
  the EDR cited. **Only when the EDR store is active** — with the store absent or empty this section
  is absent too, and EDRs are not mentioned anywhere in the brief.
- **Next** — where the flow goes: `sdd-explore` on the all-clear path, or
  `development-decisions-bootstrap` when the early guard raised an undecided question.

This brief is the entry artifact for the flow. `sdd-explore` reads it as its starting point — it is
the only phase that can read the raw request in more depth than a rough classification, so the flow
doesn't start from a vague one-liner without ever getting a closer look.

<!-- matecito-ai: sin gate de confirmación — el brief pasa directo a la fase siguiente -->
The next phase is dispatched immediately once this brief returns — there is no confirmation gate. The
orchestrator reports the four decided flags in a single notice line; nothing waits on them.

## Rules

- Do NOT ask the user anything and do NOT wait for an answer — you run headless, and this phase never
  stops to clarify. A request that is genuinely ambiguous is `sdd-explore`'s discovery cycle to raise,
  after it has read the affected code — not something you resolve here by inventing an answer or by
  refusing to classify.
- Do NOT explore the codebase in depth — that's `sdd-explore`. Your domain mapping is a rough routing aid, not analysis.
- Do NOT design or implement.
- The EDR check is SHALLOW — catch obvious early blockers, don't do design-level analysis (that's `sdd-design`).
- If the request conflicts with an Accepted EDR → `blocked`, don't route to the flow.
- If the request needs an undecided architectural choice → `needs-decision`, route to bootstrap first.
<!-- matecito-ai: explicit rule — the flags used to drop out of the brief with nothing complaining. -->
- ALWAYS emit all three downstream flags (`diagram`, `ui-test`, `worktree-isolation`) under `### Classification`, whatever their value (Step 2). This phase is their only producer; `sdd-design`, `sdd-spec`/`sdd-verify`, and the orchestrator are their only readers, and each treats an absent flag as `not-needed`/inactive **silently**. Decide them — never generate a diagram, never run proofshot, never open a workspace. `worktree-isolation` is decided by an explicit ask alone, nothing else, per Step 2.
<!-- matecito-ai: presence-based, reader-less on purpose — never treat it like diagram/ui-test's "absent = not-needed" silence, because a set-declared repo with a missing line is the anomaly, not the default. -->
- `Components` is presence-based, not absence-tolerant: with `repo.components` declared, emit the line on EVERY brief — a missing line is an anomaly, not "the axis doesn't apply". With no set declared, never emit it and never mention components. It has NO phase reader — do not invent one, do not use it to scope any later phase's work. When rendering the return (Step 5), the gate boolean (`components_axis_active`) is REQUIRED and explicit — never omit it hoping it defaults to "off".
<!-- matecito-ai: la forma del retorno tiene UNA fuente. Si volvés a escribirla acá, creaste la copia que este cambio vino a eliminar. -->
- The SHAPE of your return is `~/.claude/references/phase-returns/sdd-intake/sdd-intake.md` — one block, all three statuses. Follow it literally and never reconstruct it from memory (Step 5). This skill defines WHAT goes in each section, never how the section looks.
- Return envelope per **Section D** from `~/.claude/skills/_shared/sdd-phase-common.md`.
<!-- matecito-ai: no scope-confirmation gate — la fase siguiente se despacha apenas vuelve este brief -->
- The brief is dispatched onward immediately — no scope-confirmation gate exists in this flow.
