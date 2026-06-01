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
catch ADR conflicts or undecided architectural questions *before* exploration burns effort.

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
Follow **Section A** from `sdd-phase-common.md`.

### Step 2: Ask the Discovery Form (2-4 questions)

The raw request is almost always underspecified. Before structuring anything, ask **2-4 targeted
questions** to lock down what is ambiguous. Keep it short — this is a 30-second form, not an interrogation.

Pick the questions that actually matter for *this* request. Typical axes:
- **Scope:** what exactly is in and out? (e.g. "¿solo export CSV, o también otros formatos?")
- **Trigger / surface:** where does the user invoke it? (endpoint, button, CLI, job)
- **Constraints:** size, performance, limits (e.g. "¿reportes chicos o pueden ser de cientos de miles de filas?")
- **Behavior:** sync vs async, what happens on failure, edge cases.

Ask only what's genuinely unclear. If the user already answered something in the raw request, don't re-ask it.

### Step 3: Classify the Change

From the request + answers, classify:
- **Type:** `feature` | `bug` | `refactor` | `chore`
- **Domains touched:** map to the canonical ADR domains (e.g. an export endpoint touches `contracts`, `security`, `runtime`, maybe `data`). This is a rough mapping to help routing — NOT a deep analysis.
- **Rough size:** `trivial` | `small` | `medium` | `large`.

### Step 4: Triage — recommend a lane

Recommend a lane. You only **recommend**; the orchestrator surfaces it and the user confirms or adjusts — never apply a lane unilaterally.

**Default bias — minimum viable lane.** Recommend the *lightest* lane that still covers the change. Escalate ONLY for a concrete, named trigger; absent a trigger, the recommendation is `reduced`, not `full`. Resolve top-down and stop at the first that fits:

1. **`direct`** (no SDD) — `trivial` / `small` change with no architectural impact. Route to `direct-implementation`.
2. **`reduced`** = base, no add-ons (`intake → spec → apply → verify → archive`, always; `sdd-spec` starts from THIS brief when no proposal exists). **This is the default for substantial work** — any `small` / `medium` change with no escalation trigger lands here. Expect to recommend this most of the time; it is the norm, not an edge case.
3. **`custom`** = base + only the add-ons a trigger requires. Use this for the common middle ground instead of jumping to `full`. Add each add-on only for its trigger: `design` when there's an architectural decision, `tasks` when the work has many pieces, `explore` when the codebase area is unclear, `propose` when scope/approach needs sign-off.
4. **`full`** = base + all add-ons. Reserved for `large` changes, or work touching architecture across **multiple** domains. Requires a named trigger — do NOT recommend it as the generic "this is important" choice; that's what `custom` is for.

Escalation triggers (the ONLY reasons to go above `reduced`): an architectural decision is needed, multiple domains are touched, the surface is `large`, or the codebase area is unclear. One isolated trigger → `custom` with the matching add-on; several triggers or `large` size → `full`.

Emit the lane as the base plus the list of enabled add-ons. Be honest about size: over-routing to `full` wastes effort and is the more common failure mode — under-routing skips rigor, but the default is to trust `reduced` until a trigger says otherwise.

### Step 5: Early Guard — ADR conflicts and undecided questions

First apply the **ADR activation gate** (single source of truth in `matecito-ai:behavior`): if `.matecito-ai/adr/` is absent or empty, ADRs are inactive — **skip this entire step silently** (`status: done`, no mention of ADRs in the brief). Only when active, continue.

Read `.matecito-ai/adr/INDEX.md` and the indexes of the domains this request touches.
This is a **shallow** check — you are looking for early blockers, not doing design:

- **Conflict:** does the request contradict an `Accepted` ADR? (e.g. "endpoint público sin login" vs an auth ADR that requires protection.) → set `status: blocked`, name the ADR, and recommend resolving via `project-decisions-bootstrap` (update) or adjusting the request. Do NOT proceed to recommend the flow.
- **Undecided question:** does the request require an architectural decision that NO ADR covers? (e.g. export of huge files — sync or background job? no ADR says.) → set `status: needs-decision`, name the gap, and recommend `project-decisions-bootstrap` to capture it *before* the flow runs.
- **All clear** → `status: done`, proceed with the routing from Step 4.

The point: catch the blocker now, at intake, instead of letting the flow discover it at the design phase after wasting explore/propose/spec.

### Step 6: Persist Artifact

Follow **Section C** from `sdd-phase-common.md`.
- artifact: `intake`
- topic_key: `sdd/{change-name}/intake`
- type: `architecture`

### Step 7: Return the Structured Brief

Return EXACTLY this format (and persist the same content):

```markdown
## Intake Brief: {short title}

### Request (structured)
{1-2 sentences: what the user wants, restated clearly after the discovery form}

### Classification
- Type: {feature|bug|refactor|chore}
- Domains touched: {list of canonical ADR domains}
- Size: {trivial|small|medium|large}

### Discovery answers
- {question}: {answer}
- ...

### Triage
Lane: {direct | reduced | full | custom} — add-ons: [{explore? propose? design? tasks?} or none] — {one line why}

### Early guard (ADRs)
{One of:
- "Clear — no conflict with existing ADRs, no undecided question."
- "⛔ BLOCKED: conflicts with `<domain>/<slug>.md` — {what}. Resolve before proceeding."
- "🟡 NEEDS DECISION: `<domain>` has no ADR for {what}. Capture via project-decisions-bootstrap first."}

### Next
{direct-implementation | project-decisions-bootstrap | the first phase the chosen lane runs — `sdd-explore` if `explore` is on, else `sdd-propose` if `propose` is on, else `sdd-spec`}
```

This brief is the entry artifact for the flow. The next phase reads it as its starting point — `sdd-explore` in the full lane, `sdd-spec` in the reduced lane — so the flow doesn't start from a vague one-liner.

<!-- matecito-ai: GATE de confirmación -->
**Confirmation gate (handled by the orchestrator):** after you return this brief, the orchestrator MUST show it to the user and wait for **confirm / adjust / cancel** before launching any next phase — always, even for `trivial` changes. Do NOT assume the flow proceeds automatically. If the user adjusts the scope, the brief is updated and re-shown. See `~/.claude/skills/_shared/orchestration.md` (GATE de confirmación del alcance).

## Rules

- ALWAYS ask the discovery form first (2-4 questions) — never structure a request you haven't clarified.
- Ask ONLY what's genuinely ambiguous; don't re-ask what the user already stated.
- Do NOT explore the codebase in depth — that's `sdd-explore`. Your domain mapping is a rough routing aid, not analysis.
- Do NOT design or implement.
- The ADR check is SHALLOW — catch obvious early blockers, don't do design-level analysis (that's `sdd-design`).
- If the request conflicts with an Accepted ADR → `blocked`, don't route to the flow.
- If the request needs an undecided architectural choice → `needs-decision`, route to bootstrap first.
- Be honest in triage: trivial changes should skip the full flow.
- Return envelope per **Section D** from `sdd-phase-common.md`.
<!-- matecito-ai: el brief siempre pasa por el gate de confirmación del orquestador antes de la fase siguiente -->
- The brief ALWAYS goes through the orchestrator's confirmation gate (show to user → confirm/adjust/cancel) before any next phase runs — never assume auto-proceed. See `_shared/orchestration.md`.
