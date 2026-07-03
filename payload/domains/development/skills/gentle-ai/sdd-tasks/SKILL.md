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

> Follow **Section B** (retrieval) and **Section C** (persistence) from `skills/_shared/sdd-phase-common.md`.

- **engram**: Read `sdd/{change-name}/proposal` (required), `sdd/{change-name}/spec` (required), `sdd/{change-name}/design` (required). Save as `sdd/{change-name}/tasks`.
<!-- matecito-ai: also read the durable capability-specs of the capabilities this change touches — `.matecito-ai/development-specs/<type>/<capability>.md` (type ∈ flow|rule|lifecycle|process; concept at ~/.claude/references/spec/README.md), when present. They are the accumulated behavior contract the tasks must uphold, alongside the change spec and the design. -->
- **none**: Return result only. Never create or modify project files.

## What to Do

### Step 1: Load Skills
Follow **Section A** from `skills/_shared/sdd-phase-common.md`.

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

<!-- matecito-ai: each task carries an indented sub-line with `criteria:` (MANDATORY) and, only when ADRs are active and the task touches a decision, `· adr: <domain>/<slug>` (OPTIONAL). Do NOT touch the `- [ ]`: apply marks progress by flipping `- [ ]` → `- [x]` on the task line. The ADR ref is slug-based (`structure/layering`), never numeric. -->
- [ ] 1.1 {Concrete action — what file, what change}
      criteria: {observable check — input → result}  · adr: {<domain>/<slug> | omit}
- [ ] 1.2 {Concrete action}
      criteria: {observable check}

## Phase 2: {Phase Name} (e.g., Core Implementation)

- [ ] 2.1 {Concrete action}
      criteria: {observable check}
- [ ] 2.2 {Concrete action}
      criteria: {observable check}  · adr: {<dominio>/<slug>}

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

<!-- matecito-ai: per-task contract (criteria + adr)
- `criteria:` is MANDATORY on every task: an observable, checkable condition (input → result), not a vague claim. It is what verify consumes later; a "works fine" is useless.
- `· adr: <domain>/<slug>` is OPTIONAL: add it ONLY when the task implements or touches a registered decision. Mechanical tasks (add a field, move a file) carry NO adr — a false ref is worse than none.
- What counts as a "decision" (and what does not) is canonical in `~/.claude/references/adr/README.md` — a task carries `· adr:` only if it touches an ADR in that sense, not for every change.
- ADR activation gate (with `flagDecisionGaps` OFF, the default): if `.matecito-ai/adr/` does not exist or is empty, do NOT add the `adr` sub-part on any task and do not mention it. `criteria` is still required. matecito-ai NEVER requires an ADR. (With the flag ON the behavior changes — see the decision-gap hook below.) -->
<!-- matecito-ai: decision-gap detection hook
Active ONLY when flagDecisionGaps=true (does NOT depend on ADRs existing). When active: for EACH task that touches a decision, emit `· adr: <domain>/<slug>` mapped to a catalog concern, whether or not the ADR exists — this overrides the flag-off rule of "omit adr if there is no `.matecito-ai/adr/`". Then, for each `· adr:`, check whether `.matecito-ai/adr/<domain>/<slug>.md` exists: if NOT, the ref is a dangling decision gap — leave it as-is (do not modify or mark it). With zero ADRs, every decision is a gap (bootstrap the first ones). The dangling refs in the artifact are the gap list that sdd-verify consumes. When flag off: byte-identical behavior to before, no mention. -->
- (Decision-gap detection — flag-gated) When `flagDecisionGaps=true` (regardless of ADR presence): emit a concern-mapped `· adr:` for each decision-touching task even if no ADR exists yet; a `· adr:` whose target file is absent under `.matecito-ai/adr/` is a decision gap — leave it verbatim (the dangling ref IS the signal; with zero ADRs every decision is a gap → bootstrap). Silent when flag off.


### Review Workload Forecast Rules

Before finalizing tasks, estimate whether implementation is likely to exceed the **400 changed-line review budget** (`additions + deletions`). This is a planning guard, not an exact diff count.

Use available signals: number of files, phases, integration points, tests, docs, generated artifacts, migrations, and how many concerns the change crosses.

If the estimate is **High** or likely above 400 lines:

1. Mark `Chained PRs recommended` as `Yes`.
2. Split tasks into **work units** that can become chained or stacked PRs.
3. Each suggested PR must have a clear start, clear finish, verification, and autonomous scope.
4. **Ask the user which chain strategy to use** (this is a team decision):
   - **Stacked PRs to main** — each PR merges to main in order. Fast iteration, fix on the go. Best for speed-first teams and independent slices.
   - **Feature Branch Chain** — the feature/tracker branch accumulates the final integration; PR #1 targets the tracker branch, later PRs target the immediate previous PR branch so each child diff stays focused. Only the tracker merges to main. Best for rollback control and coordinated releases.
   - **size:exception** — keep it as a single PR with maintainer approval. Best for generated code, migrations, or vendor diffs.
5. Cache the user's choice and set `Decision needed before apply` from delivery strategy:
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

Follow **Section C** from `skills/_shared/sdd-phase-common.md`.
- artifact: `tasks`
- topic_key: `sdd/{change-name}/tasks`
- type: `architecture`

### Step 5: Return Summary

Return to the orchestrator:

```markdown
## Tasks Created

**Change**: {change-name}
**Location**: Engram `sdd/{change-name}/tasks` (engram) | inline (none)

### Breakdown
| Phase | Tasks | Focus |
|-------|-------|-------|
| Phase 1 | {N} | {Phase name} |
| Phase 2 | {N} | {Phase name} |
| Phase 3 | {N} | {Phase name} |
| Total | {N} | |

### Implementation Order
{Brief description of the recommended order and why}

### Review Workload Forecast
- Estimated changed lines: {estimate or range}
- 400-line budget risk: {Low | Medium | High}
- Chained PRs recommended: {Yes | No}
- Delivery strategy: {ask-on-risk | auto-chain | single-pr | exception-ok}
- Decision needed before apply: {Yes | No}
- Suggested work-unit PR split: {brief list or "Not needed"}

### Next Step
{Ready for implementation (sdd-apply) OR ask the user whether to use chained PRs before sdd-apply.}
```

## Rules

- ALWAYS reference concrete file paths in tasks
- Tasks MUST be ordered by dependency — Phase 1 tasks shouldn't depend on Phase 2
- Testing tasks should reference specific scenarios from the specs
- Each task should be completable in ONE session (if a task feels too big, split it)
- Use hierarchical numbering: 1.1, 1.2, 2.1, 2.2, etc.
- NEVER include vague tasks like "implement feature" or "add tests"
- If the project uses TDD, integrate test-first tasks: RED task (write failing test) → GREEN task (make it pass) → REFACTOR task (clean up)
<!-- matecito-ai: budget subido de 530 → 800 palabras para absorber la sub-línea `criteria:` (criteria + adr) por tarea. -->
- **Size budget**: Tasks artifact MUST be under 800 words. Each task: the `- [ ]` line + one indented `criteria:` sub-line (max 2 lines total). Use checklist format, not paragraphs.
- **Review workload guard**: ALWAYS include the Review Workload Forecast. If likely above 400 changed lines, recommend chained PRs and honor the received delivery strategy for whether a decision/exception is needed before apply.
- Return envelope per **Section D** from `skills/_shared/sdd-phase-common.md`.
