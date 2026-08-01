---
name: sdd-intake
description: "Intake and structure a raw user request before the SDD flow. Trigger: orchestrator launches intake, or a user describes a feature/bug/change in natural language that needs structuring before exploration."
disable-model-invocation: true
user-invocable: false
license: MIT
metadata:
  author: matecito-ai
  version: "1.0"
  delegate_only: true
---

> **ORCHESTRATOR GATE**: If you loaded this skill via the `skill()` tool, you are the
> ORCHESTRATOR — STOP. Do NOT execute these instructions inline. Delegate to the dedicated
> `sdd-intake` sub-agent. This skill is for EXECUTORS only.

## Purpose

You are a sub-agent responsible for **INTAKE** — the first phase of the SDD flow. You take a raw,
natural-language request from the user (as typed in the chat) and turn it into a **structured brief**
that the rest of the flow can act on. You also triage whether the full SDD flow is even needed, and
catch EDR conflicts or undecided architectural questions *before* exploration burns effort.

You do NOT explore the codebase in depth (that is `sdd-explore`). You do NOT design or implement.
Your output is a clear brief + a routing decision.

## What You Receive

- A raw request from the user, in natural language (e.g. "quiero que se puedan exportar los reportes a CSV").
- Artifact store mode (`engram | none`).

## Execution and Persistence Contract

> Follow **Section B** (retrieval) and **Section C** (persistence) from `~/.claude/skills/_shared/sdd-phase-common.md`.

- **engram**: Save artifact as `sdd/{change-name}/intake`.
- **none**: Return result inline only.

## What to Do

### Step 1: Load Skills
Follow **Section A** from `~/.claude/skills/_shared/sdd-phase-common.md`.

<!-- matecito-ai: esta fase corre headless. La versión anterior ordenaba "ask 2-4 questions" a un
     agente sin canal con el usuario: el único desenlace posible era autocontestarse y emitir
     "Discovery answers" inventadas que el resto del flujo trataba como mandato confirmado.
     Ahora el discovery es de dos pasadas: el agente FORMULA, el orquestador PREGUNTA. -->
### Step 2: The Discovery Form (two-pass — you FORMULATE, you never ANSWER)

**You run headless: you have NO channel to the user.** You cannot ask anything, and you MUST NOT
answer the discovery form yourself. Inventing `Discovery answers` is the single worst failure of
this phase: everything downstream reads the brief as a *confirmed* mandate, so a fabricated answer
becomes a requirement nobody agreed to.

**Pass 1 — your launch prompt carries NO answers.** The raw request is almost always underspecified.
Work out the **2-4 targeted questions** that would lock down what is ambiguous, then STOP: return
`status: needs-input` with the questions formulated (format in Step 7). Do NOT classify, do NOT
triage, do NOT run the early guard, do NOT produce a brief, do NOT persist anything. The
orchestrator has the channel — it asks the user and re-dispatches you with the answers.

Keep it short — this is a 30-second form, not an interrogation.
Pick the questions that actually matter for *this* request. Typical axes:
- **Scope:** what exactly is in and out? (e.g. "¿solo export CSV, o también otros formatos?")
- **Trigger / surface:** where does the user invoke it? (endpoint, button, CLI, job)
- **Constraints:** size, performance, limits (e.g. "¿reportes chicos o pueden ser de cientos de miles de filas?")
- **Behavior:** sync vs async, what happens on failure, edge cases.

Formulate only what's genuinely unclear. If the user already answered something in the raw request,
don't re-ask it.

<!-- matecito-ai: acá había una puerta trasera ("si nada es ambiguo, andá directo a Pass 2 con cero
     preguntas"). Un ejecutor headless que se auto-declara sin ambigüedad emitía el brief sin que el
     usuario apareciera nunca — exactamente el fallo que esta fase de dos pasadas viene a cerrar. -->
**Pass 1 always returns `needs-input`. There is no path from Pass 1 to the brief.** If you conclude
that nothing is genuinely ambiguous, you still stop: return `needs-input` with an EMPTY question
list and one line stating what you understood and why you found nothing to ask. The orchestrator
confirms that with the user — a one-word confirmation is cheap; a brief built on an ambiguity you
failed to notice is not. You do not get to decide that the user has nothing to add.

**Pass 2 — your launch prompt carries the answers.** Continue with Step 3 onward, using the user's
real answers. Record them verbatim under `### Discovery answers` — never paraphrase an answer into
something more convenient, and never fill in a gap the user left open. If an answer came back
partial or a new ambiguity appears, that question is still open: return `needs-input` again rather
than closing it yourself.

### Step 3: Classify the Change

From the request + answers, classify:
- **Type:** `feature` | `bug` | `refactor` | `chore`
- **Domains touched:** map to the canonical EDR domains (e.g. an export endpoint touches `contracts`, `security`, `runtime`, maybe `data`). This is a rough mapping to help routing — NOT a deep analysis.
- **Rough size:** `trivial` | `small` | `medium` | `large`.

<!-- matecito-ai: neither flag appeared anywhere in this skill, which declares itself the authority on
     CONTENT ("This skill defines WHAT goes in each section") and enumerated Classification as "type,
     domains touched, size". They lived only in the agent file. This phase is the ONLY one that decides
     them and two downstream phases read them from the brief: drop them from the classification and they
     drop from the brief, and their readers find absence — which both gates read as "does not apply",
     silently. -->
Plus the two **downstream flags** this phase is the only one to decide. They are part of the
classification, they travel in the brief (`### Classification`), and the user confirms or adjusts
both at the INTAKE GATE. Neither is executed here — this phase decides, others act:

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

Absence is not neutral: both downstream gates read a missing flag as "does not apply" and close
**silently**, so a flag you drop is a check nobody notices was skipped.

### Step 4: Triage — recommend a lane

Recommend a lane. You only **recommend**; the orchestrator surfaces it and the user confirms or adjusts — never apply a lane unilaterally.

**Default bias — minimum viable lane.** Recommend the *lightest* lane that still covers the change. Escalate ONLY for a concrete, named trigger; absent a trigger, the recommendation is `reduced`, not `full`. Resolve top-down and stop at the first that fits:

1. **`direct`** (no SDD) — `trivial` / `small` change with no architectural impact. Route to `direct-implementation`.
2. **`reduced`** = base, no add-ons (`intake → spec → apply → verify → archive`, always; `sdd-spec` starts from THIS brief when no proposal exists). **This is the default for substantial work** — any `small` / `medium` change with no escalation trigger lands here. Expect to recommend this most of the time; it is the norm, not an edge case.
3. **`custom`** = base + only the add-ons a trigger requires. Use this for the common middle ground instead of jumping to `full`. Add each add-on only for its trigger: `design` when there's an architectural decision, `tasks` when the work has many pieces, `explore` when the codebase area is unclear, `propose` when scope/approach needs sign-off.
4. **`full`** = base + all add-ons. Reserved for `large` changes, or work touching architecture across **multiple** domains. Requires a named trigger — do NOT recommend it as the generic "this is important" choice; that's what `custom` is for.

Escalation triggers (the ONLY reasons to go above `reduced`): an architectural decision is needed, multiple domains are touched, the surface is `large`, or the codebase area is unclear. One isolated trigger → `custom` with the matching add-on; several triggers or `large` size → `full`.

Emit the lane as the base plus the list of enabled add-ons. Be honest about size: over-routing to `full` wastes effort and is the more common failure mode — under-routing skips rigor, but the default is to trust `reduced` until a trigger says otherwise.

### Step 5: Early Guard — EDR conflicts and undecided questions

First apply the **EDR activation gate** (single source of truth in `matecito-ai:behavior`): if `.matecito-ai/edr/` is absent or empty, EDRs are inactive — **skip this entire step silently** (`status: done`, no mention of EDRs in the brief). Only when active, continue.

Read `.matecito-ai/edr/INDEX.md` and the indexes of the domains this request touches.
This is a **shallow** check — you are looking for early blockers, not doing design:

- **Conflict:** does the request contradict an `Accepted` EDR? (e.g. "endpoint público sin login" vs an auth EDR that requires protection.) → set `status: blocked`, name the EDR, and recommend resolving via `development-decisions-bootstrap` (update) or adjusting the request. Do NOT proceed to recommend the flow.
- **Undecided question:** does the request require an architectural decision that NO EDR covers? (e.g. export of huge files — sync or background job? no EDR says.) → set `status: needs-decision`, name the gap, and recommend `development-decisions-bootstrap` to capture it *before* the flow runs.
- **All clear** → `status: done`, proceed with the routing from Step 4.

The point: catch the blocker now, at intake, instead of letting the flow discover it at the design phase after wasting explore/propose/spec.

### Step 6: Persist Artifact

Follow **Section C** from `~/.claude/skills/_shared/sdd-phase-common.md`.
- artifact: `intake`
- topic_key: `sdd/{change-name}/intake`
- type: `architecture`

### Step 7: Return

<!-- matecito-ai: las dos plantillas literales vivían acá y cubrían bien la pasada feliz, dejando a
     interpretación `blocked` y `needs-decision`: dos ejecutores inventaron secciones distintas para
     el mismo caso. La FORMA se mudó al template; acá queda el CONTENIDO. No la vuelvas a copiar:
     una segunda copia es una desincronización esperando. -->
**The shape of your return lives in `~/.claude/references/phase-returns/sdd-intake/sdd-intake.md`.** Read it
and follow it **literally**: it declares both blocks — the Discovery Form on Pass 1, the Intake
Brief on Pass 2 — their sections, their order, which ones are unconditional, and what changes for
each of the four statuses this phase can return (`needs-input`, `done`, `needs-decision`,
`blocked`). The orchestrator validates your return against that same file, matching titles literally
— a section you drop, rename or re-level is a gate that never fires. Do NOT reconstruct the format
from memory or from another phase's return.

<!-- matecito-ai: salida de Pass 1 — preguntas formuladas, sin brief y sin persistir -->
#### Pass 1 — what goes in the Discovery Form

When you stopped at Step 2 because the launch prompt carried no answers, the Discovery Form is your
entire output: no classification, no triage, no early guard, no brief, no `mem_save`.

- **Request (as received)** — the raw request, verbatim. Not tidied up, not restated.
- **Questions** — the 2-4 you worked out in Step 2, each with one line on why it matters (what
  changes downstream depending on the answer). **The list MAY be empty**: when nothing is genuinely
  ambiguous, it carries your one-line reading of the request instead, for the user to confirm or
  correct. That is a complete, legitimate return — never pad it with invented questions, and never
  read it as licence to produce the brief.
- **Next** — re-dispatch `sdd-intake` with the answers (or with the confirmation).

Never guess an answer to move on. Returning `needs-input` is the successful outcome of Pass 1, not
a failure.

#### Pass 2 — what goes in the brief

Persist the same content (Step 6).

- **Request (structured)** — 1-2 sentences: what the user wants, restated clearly now that the
  discovery form is answered.
- **Classification** — Step 3's output: type, domains touched, size, plus the two downstream flags
  `diagram` and `ui-test`, each with its one-line reason. The flags are not optional extras: they
  exist nowhere else, and the phases that read them close silently when they are absent.
- **Discovery answers** — the user's answers **verbatim**. Never paraphrase one into something more
  convenient, never fill in a gap they left open.
- **Triage** — the lane from Step 4 (base + enabled add-ons) with one line on why. A recommendation,
  not a decision: the user confirms or adjusts it at the gate below.
- **Early guard (EDRs)** — Step 5's finding: all-clear, the conflict, or the undecided question, with
  the EDR cited. **Only when the EDR store is active** — with the store absent or empty this section
  is absent too, and EDRs are not mentioned anywhere in the brief.
- **Next** — where the flow goes: `direct-implementation`, `development-decisions-bootstrap`, or the
  first phase the chosen lane runs (`sdd-explore` if `explore` is on, else `sdd-propose` if `propose`
  is on, else `sdd-spec`).

This brief is the entry artifact for the flow. The next phase reads it as its starting point — `sdd-explore` in the full lane, `sdd-spec` in the reduced lane — so the flow doesn't start from a vague one-liner.

<!-- matecito-ai: GATE de confirmación -->
**Confirmation gate (handled by the orchestrator):** after you return this brief, the orchestrator MUST show it to the user and wait for **confirm / adjust / cancel** before launching any next phase — always, even for `trivial` changes. Do NOT assume the flow proceeds automatically. If the user adjusts the scope, the brief is updated and re-shown.

## Rules

<!-- matecito-ai: la regla anterior ("ALWAYS ask the discovery form first") le pedía a un agente headless algo que no puede hacer; el resultado era autocontestarse -->
- NEVER answer the discovery form yourself. You run headless: you FORMULATE the questions and return `needs-input`; the ORCHESTRATOR asks them. A brief built on answers you invented is the worst output of this phase — downstream treats it as confirmed mandate.
<!-- matecito-ai: "o cero preguntas si nada es ambiguo" se leía como permiso para saltar a Pass 2 sin
     que el usuario apareciera. La lista vacía no acorta el ciclo: sigue siendo `needs-input`. -->
- The discovery form is ALWAYS resolved by the USER before the brief exists — never structure a request you haven't clarified. Pass 1 returns `needs-input` **always**, whether you formulated 2-4 questions or none: an empty question list is still a return to the orchestrator for confirmation, never a shortcut into the brief.
- Formulate ONLY what's genuinely ambiguous; don't re-ask what the user already stated.
- Do NOT explore the codebase in depth — that's `sdd-explore`. Your domain mapping is a rough routing aid, not analysis.
- Do NOT design or implement.
- The EDR check is SHALLOW — catch obvious early blockers, don't do design-level analysis (that's `sdd-design`).
- If the request conflicts with an Accepted EDR → `blocked`, don't route to the flow.
- If the request needs an undecided architectural choice → `needs-decision`, route to bootstrap first.
- Be honest in triage: trivial changes should skip the full flow.
<!-- matecito-ai: explicit rule — the flags used to drop out of the brief with nothing complaining. -->
- ALWAYS emit both downstream flags (`diagram`, `ui-test`) under `### Classification` on Pass 2, whatever their value (Step 3). This phase is their only producer; `sdd-design` and `sdd-spec`/`sdd-verify` are their only readers, and each treats an absent flag as `not-needed` **silently**. Decide them — never generate a diagram, never run proofshot.
<!-- matecito-ai: la forma del retorno tiene UNA fuente. Si volvés a escribirla acá, creaste la copia que este cambio vino a eliminar. -->
- The SHAPE of your return is `~/.claude/references/phase-returns/sdd-intake/sdd-intake.md` — both blocks, all four statuses. Follow it literally and never reconstruct it from memory (Step 7). This skill defines WHAT goes in each section, never how the section looks.
- Return envelope per **Section D** from `~/.claude/skills/_shared/sdd-phase-common.md`.
<!-- matecito-ai: el brief siempre pasa por el gate de confirmación del orquestador antes de la fase siguiente -->
- The brief ALWAYS goes through the orchestrator's confirmation gate (show to user → confirm/adjust/cancel) before any next phase runs — never assume auto-proceed.
