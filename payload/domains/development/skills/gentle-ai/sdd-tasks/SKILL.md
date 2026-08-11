---
name: sdd-tasks
description: "Break an SDD change into implementation tasks. Trigger: orchestrator launches task planning for a change."
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
> the dedicated `sdd-tasks` sub-agent using your platform's delegation primitive
> (e.g., `task(...)`, sub-agent invocation, etc.). This skill is for EXECUTORS
> only.

## Purpose

You are a sub-agent responsible for creating the TASK BREAKDOWN. You take the proposal, specs, and design, then produce a `tasks.md` with concrete, actionable implementation steps organized by phase.

## What You Receive

From the orchestrator:
- Change name
- Artifact store mode (`engram | none`) <!-- matecito-ai: openspec/hybrid removidos -->
- Delivery strategy (`ask-on-risk | auto-chain | single-pr | exception-ok`)

## Execution and Persistence Contract

> Follow **Section B** (retrieval) and **Section C** (persistence) from `~/.claude/skills/_shared/sdd-phase-common.md`.

<!-- matecito-ai: declaraba los tres como required, contra la regla de nearest-upstream del fragmento:
     en lane `reduced`/`custom` no hay proposal ni design, y el agente ya trata design como opcional.
     Un `required` que no se cumple en el lane por defecto enseña a ignorar los `required`. -->
- **engram**: Read `sdd/{change-name}/spec` (**required** — the floor). Read `sdd/{change-name}/design` and `sdd/{change-name}/proposal` **when they exist**: in `reduced` and `custom` lanes those phases may not have run, and their absence is normal, not an error — decompose from what you do have. Save as `sdd/{change-name}/tasks`.
<!-- matecito-ai: also read the durable capability-specs of the capabilities this change touches — `.matecito-ai/development-specs/<type>/<capability>.md` (type ∈ flow|rule|lifecycle|process; concept at ~/.claude/references/spec/README.md), when present. They are the accumulated behavior contract the tasks must uphold, alongside the change spec and the design. -->
- **none**: Return result only. Never create or modify project files.

## What to Do

### Step 1: Load Skills
Follow **Section A** from `~/.claude/skills/_shared/sdd-phase-common.md`.

### Step 2: Analyze the Design

From the design document, identify:
- All files that need to be created/modified/deleted
- The dependency order (what must come first)
- Testing requirements per component

### Step 3: Write tasks.md

<!-- matecito-ai: engram-only; no se crean archivos. -->
Compose the tasks content in memory — you will persist it in Step 4 (Engram).

#### Task File Format

```markdown
# Tasks: {Change Title}

## Review Workload Forecast

| Field | Value |
|-------|-------|
| Estimated changed lines | <rough estimate or range> |
| 400-line budget risk | Low / Medium / High |
| Chained PRs recommended | Yes / No |
| Suggested split | <single PR or PR 1 → PR 2 → PR 3> |
| Delivery strategy | <ask-on-risk / auto-chain / single-pr / exception-ok> |
| Chain strategy | <stacked-to-main / feature-branch-chain / size-exception / pending> |

Decision needed before apply: <Yes|No>
Chained PRs recommended: <Yes|No>
Chain strategy: <stacked-to-main|feature-branch-chain|size-exception|pending>
400-line budget risk: <Low|Medium|High>

### Suggested Work Units

| Unit | Goal | Likely PR | Notes |
|------|------|-----------|-------|
| 1 | <standalone deliverable> | PR 1 | <base branch; tests/docs included> |
| 2 | <standalone deliverable> | PR 2 | <immediate parent/base branch boundary; depends on PR 1 or independent> |

## Phase 1: {Phase Name} (e.g., Infrastructure / Foundation)

<!-- matecito-ai: each task carries an indented sub-line with `criteria:` (MANDATORY) and, only when EDRs are active and the task touches a decision, `· edr: <domain>/<slug>` (OPTIONAL). Do NOT touch the `- [ ]`: apply marks progress by flipping `- [ ]` → `- [x]` on the task line. The EDR ref is slug-based (`structure/layering`), never numeric. A task MAY also carry a further indented `· parallel-group: <id>` sub-line — its own line, never folded into the `criteria:` line and never on the `- [ ]` line — declaring which tasks run concurrently; see "Parallel-group mark" below. -->
- [ ] 1.1 {Concrete action — what file, what change}
      criteria: {observable check — input → result}  · edr: {<domain>/<slug> | omit}
- [ ] 1.2 {Concrete action}
      criteria: {observable check}
      · parallel-group: {<id> | omit}

## Phase 2: {Phase Name} (e.g., Core Implementation)

- [ ] 2.1 {Concrete action}
      criteria: {observable check}
- [ ] 2.2 {Concrete action}
      criteria: {observable check}  · edr: {<dominio>/<slug>}

## Phase 3: {Phase Name} (e.g., Testing / Verification)

- [ ] 3.1 {Write tests for ...}
      criteria: {observable check}
- [ ] 3.2 {Verify integration between ...}
      criteria: {observable check}

## Phase 4: {Phase Name} (e.g., Cleanup / Documentation)

- [ ] 4.1 {Update docs/comments}
      criteria: {observable check}
```

### Task Writing Rules

Each task MUST be:

| Criteria | Example ✅ | Anti-example ❌ |
|----------|-----------|----------------|
| **Specific** | "Create `internal/auth/middleware.go` with JWT validation" | "Add auth" |
| **Actionable** | "Add `ValidateToken()` method to `AuthService`" | "Handle tokens" |
| **Verifiable** | "Test: `POST /login` returns 401 without token" | "Make sure it works" |
| **Small** | One file or one logical unit of work | "Implement the feature" |

<!-- matecito-ai: per-task contract (criteria + edr)
- `criteria:` is MANDATORY on every task: an observable, checkable condition (input → result), not a vague claim. It is what verify consumes later; a "works fine" is useless.
- `· edr: <domain>/<slug>` is OPTIONAL: add it ONLY when the task implements or touches a registered decision. Mechanical tasks (add a field, move a file) carry NO edr — a false ref is worse than none.
- What counts as a "decision" (and what does not) is canonical in `~/.claude/references/edr/README.md` — a task carries `· edr:` only if it touches an EDR in that sense, not for every change.
- EDR activation gate (with `flagDecisionGaps` OFF, the default): if `.matecito-ai/edr/` does not exist or is empty, do NOT add the `edr` sub-part on any task and do not mention it. `criteria` is still required. matecito-ai NEVER requires an EDR. (With the flag ON the behavior changes — see the decision-gap hook below.) -->
<!-- matecito-ai: decision-gap detection hook
Active ONLY when flagDecisionGaps=true (does NOT depend on EDRs existing). When active: for EACH task that touches a decision, emit `· edr: <domain>/<slug>` mapped to a catalog concern, whether or not the EDR exists — this overrides the flag-off rule of "omit edr if there is no `.matecito-ai/edr/`". Then, for each `· edr:`, check whether `.matecito-ai/edr/<domain>/<slug>.md` exists: if NOT, the ref is a dangling decision gap — leave it as-is (do not modify or mark it). With zero EDRs, every decision is a gap (bootstrap the first ones). The dangling refs in the artifact are the gap list that sdd-verify consumes. When flag off: byte-identical behavior to before, no mention. -->
- (Decision-gap detection — flag-gated) When `flagDecisionGaps=true` (regardless of EDR presence): emit a concern-mapped `· edr:` for each decision-touching task even if no EDR exists yet; a `· edr:` whose target file is absent under `.matecito-ai/edr/` is a decision gap — leave it verbatim (the dangling ref IS the signal; with zero EDRs every decision is a gap → bootstrap). Silent when flag off.

### Parallel-group mark

<!-- matecito-ai: this is the canonical definition `sdd-apply`'s `parallel-batch.md` cites instead of
     assuming a form — see its `## Eligibility` section. Defined once, here; do not restate it there. -->

A task MAY carry **exactly one** further indented sub-line, of the same family as `criteria:` and
`· edr:`: `· parallel-group: <id>`, on its own line, never on the `- [ ]` line. `<id>` is a short
token naming a group; it carries no meaning beyond identity — two tasks are in the same group if and
only if their ids are equal.

- **Optional; absence means serial** — today's behavior. Do NOT infer a group for an unmarked task,
  and do NOT require re-marking any existing tasks artifact.
- **Same id ⇒ concurrently safe, and ONLY when it is.** Give two tasks the same id ONLY when they can
  run at the same time without conflicting. Two tasks that would collide (e.g. both rewrite the same
  file) MUST NOT share an id — each is still markable, just in a group of its own.
- **One group is one batch.** Tasks sharing an id run together; tasks with different ids never mix in
  the same batch — each group is its own round. Never merge two ids because their tasks look
  independent of each other; that is the inference this mark exists to remove.
- **A group of one is legal**, not malformed — it simply yields no fan-out for that task.
- **The mark is a form, not prose.** It is the only way to declare parallel eligibility; the emitted
  artifact carries no free-text note about which tasks are parallel.

### The per-task line budget with the mark

A marked task's cap is **three lines**: the `- [ ]` line, the `criteria:` sub-line, and the
`· parallel-group:` sub-line. An unmarked task keeps the existing two-line cap — see "Size budget"
below.

### Review Workload Forecast Rules

Before finalizing tasks, estimate whether implementation is likely to exceed the **400 changed-line review budget** (`additions + deletions`). This is a planning guard, not an exact diff count.

Use available signals: number of files, phases, integration points, tests, docs, generated artifacts, migrations, and how many concerns the change crosses.

If the estimate is **High** or likely above 400 lines:

1. Mark `Chained PRs recommended` as `Yes`.
2. Split tasks into **work units** that can become chained or stacked PRs.
3. Each suggested PR must have a clear start, clear finish, verification, and autonomous scope.
<!-- matecito-ai: acá decía "Ask the user which chain strategy to use" y "Cache the user's choice", a
     un ejecutor headless sin canal y sin estado de sesión — el mismo defecto que se arregló en
     sdd-intake. Y encima era redundante: el Review Workload Guard del orquestador YA pregunta esto
     después de esta fase y antes de apply, que es donde corresponde. Vos documentás las opciones
     para que la decisión se pueda tomar; no la tomás ni la pedís. -->
4. **Document the three chain strategies in the forecast** so the user can choose when the guard asks. You do NOT ask and you do NOT choose: you have no channel to the user, and the orchestrator's Review Workload Guard puts this question after your phase and before apply. Lay them out:
   - **Stacked PRs to main** — each PR merges to main in order. Fast iteration, fix on the go. Best for speed-first teams and independent slices.
   - **Feature Branch Chain** — the feature/tracker branch accumulates the final integration; PR #1 targets the tracker branch, later PRs target the immediate previous PR branch so each child diff stays focused. Only the tracker merges to main. Best for rollback control and coordinated releases.
   - **size:exception** — keep it as a single PR with maintainer approval. Best for generated code, migrations, or vendor diffs.

   Set `Chain strategy: pending` unless the orchestrator handed you an already-resolved one in the launch prompt. `pending` is the correct, expected value here — not a gap you should fill in.
5. Set `Decision needed before apply` from the `delivery_strategy` the orchestrator passed you (you receive it; you never cache or infer it):
   - `ask-on-risk`: `Yes` — orchestrator asks before apply.
   - `auto-chain`: `No` — orchestrator proceeds with the first slice using the chosen chain strategy.
   - `single-pr`: `Yes` — orchestrator must require `size:exception` before apply.
   - `exception-ok`: `No` — maintainer has accepted `size:exception`.

Do not bury this in prose. Put the forecast near the top of the tasks artifact so the user sees it before implementation starts.

The forecast MUST include these exact plain-text lines so downstream guards can match them literally:

```text
Decision needed before apply: Yes|No
Chained PRs recommended: Yes|No
Chain strategy: stacked-to-main|feature-branch-chain|size-exception|pending
400-line budget risk: Low|Medium|High
```

You may keep the table for readability, but the plain-text lines are the guard contract.

For `feature-branch-chain`, suggested work units SHOULD name the intended base boundary: PR #1 base = feature/tracker branch; PR #2 base = PR #1 branch; PR #3 base = PR #2 branch. If a child PR would show previous PR changes, the base is wrong and must be retargeted/rebased before review.

### Phase Organization Guidelines

```
Phase 1: Foundation / Infrastructure
  └─ New types, interfaces, database changes, config
  └─ Things other tasks depend on

Phase 2: Core Implementation
  └─ Main logic, business rules, core behavior
  └─ The meat of the change

Phase 3: Integration / Wiring
  └─ Connect components, routes, UI wiring
  └─ Make everything work together

Phase 4: Testing
  └─ Unit tests, integration tests, e2e tests
  └─ Verify against spec scenarios

Phase 5: Cleanup (if needed)
  └─ Documentation, remove dead code, polish
```

### Step 4: Persist Artifact

**This step is MANDATORY — do NOT skip it.**

Follow **Section C** from `~/.claude/skills/_shared/sdd-phase-common.md`.
- artifact: `tasks`
- topic_key: `sdd/{change-name}/tasks`
- type: `architecture`

### Step 5: Return Summary

<!-- matecito-ai: la plantilla vivía acá inline y sólo cubría el caso feliz. Ahora vive en un único
     archivo, el mismo contra el que el orquestador valida el retorno (Return Contract Check).
     Mantener además una copia acá sería volver a tener dos formatos que se desincronizan. -->

**Follow `~/.claude/references/phase-returns/sdd-tasks/sdd-tasks.md` literally.** That file is the canonical
shape of this phase's return: which sections, in which order, with which titles, which ones are
unconditional and what changes on `blocked`. The orchestrator validates your return against that
same file and matches titles **literally** — a section you drop, rename or re-level is a gate that
never fires. Do not improvise a format here and do not omit a section because you have nothing to
report: it ships with a `None…` sentinel.

Three things the template expects you to already know from this skill:

- `### Tasks not traceable to spec/design` carries every task that links to NO spec requirement and
  to nothing the design establishes, each with what motivated it (a gap you found, an implied
  prerequisite, a project convention). This is where silent scope creep enters — work the user did
  not ask for that lands in the checklist, gets implemented and gets verified. Do NOT drop those
  tasks and do NOT fold them in as if they came from the spec.
- `### Review Workload Forecast` repeats the guard lines of the artifact (see the forecast rules
  above). Their labels are matched literally downstream — reproduce them verbatim.
- A decision you cannot settle yourself does not go into either of those. It makes you return
  `blocked`, with the question in the template's `### Blocker` section — never in `risks` (D.4).

## Rules

- ALWAYS reference concrete file paths in tasks
- Tasks MUST be ordered by dependency — Phase 1 tasks shouldn't depend on Phase 2
- Testing tasks should reference specific scenarios from the specs
- Each task should be completable in ONE session (if a task feels too big, split it)
- Use hierarchical numbering: 1.1, 1.2, 2.1, 2.2, etc.
- NEVER include vague tasks like "implement feature" or "add tests"
<!-- matecito-ai: una tarea sin origen no está prohibida — pero se declara, no se cuela -->
- Every task MUST link to a spec requirement or to something the design establishes. A task that links to neither is not forbidden (a real gap or an implied prerequisite may warrant it) but it MUST be declared under `### Tasks not traceable to spec/design` in your return — never folded in silently
- If the project uses TDD, integrate test-first tasks: RED task (write failing test) → GREEN task (make it pass) → REFACTOR task (clean up)
<!-- matecito-ai: budget subido de 530 → 800 palabras para absorber la sub-línea `criteria:` (criteria + edr) por tarea; el cap por tarea sube de 2 a 3 líneas cuando lleva `· parallel-group:`. -->
- **Size budget**: Tasks artifact MUST be under 800 words. Each task: the `- [ ]` line + one indented `criteria:` sub-line (2 lines total) — 3 lines total for a task that also carries `· parallel-group:`. Use checklist format, not paragraphs.
- **Review workload guard**: ALWAYS include the Review Workload Forecast. If likely above 400 changed lines, recommend chained PRs and honor the received delivery strategy for whether a decision/exception is needed before apply.
- Return envelope per **Section D** from `~/.claude/skills/_shared/sdd-phase-common.md`.
