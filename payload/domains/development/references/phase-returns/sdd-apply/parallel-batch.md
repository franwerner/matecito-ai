<!-- matecito-ai: sole home of the sdd-apply parallel-batch mechanism — eligibility, isolation, the
     base handshake, the commit convention, the Task Run Report shape, integration mechanics and
     conflict handling. This is the SECOND declared case of Phase fan-out (see
     `~/.claude/matecito-ai/domains/development.md` → "Phase fan-out"); the first is
     `subverifier-groups.md`. No `.yaml` pair: the Task Run Report is a bespoke block nothing renders
     or validates mechanically, exactly like `sdd-verify`'s Sub-Report — `sdd-apply.yaml` stays the
     shape of the CONSOLIDATED `## Implementation Progress` block and is unchanged by this file. -->

# `sdd-apply` parallel batch

Applies to a tasks artifact whenever at least one `· parallel-group: <id>` id has **two or more**
members — eligibility is stated and evaluated **per group**, never over the artifact as a whole (four
marked tasks split into two groups of two are two eligible batches run in successive rounds, not one
eligible batch of four; see "One group is one batch" below). A group with fewer than two members, an
unmarked task, no tasks artifact at all (`reduced`/`custom` lanes), or a dirty `git status --porcelain`
at dispatch time all run **serial mode**: today's single-agent path, unchanged by this file. Serial mode
is not documented here because nothing about it changed.

## Eligibility (the mark, not an inference)

The mark's form and semantics are defined once, in `sdd-tasks`'s task-format definition —
`~/.claude/skills/sdd-tasks/SKILL.md` → "Parallel-group mark". This file cites that definition instead
of restating it; read it there for the sub-line's exact shape and for what makes two tasks share an id.

The eligibility mark comes from the tasks artifact **as emitted** — this file, the SKILL and the agent
MUST NOT infer independence from file overlap, task description, or anything else, and MUST NOT promote
an unmarked task to parallel because it looks safe. Eligibility is scoped to **one batch**: it does not
cross batch boundaries and does not cross a chained/stacked PR slice — two slices of a delivery chain
never run at the same time because of this mechanism.

**One group is one batch; groups never mix.** Tasks sharing an id are dispatched together as one
isolated-run batch; tasks with different ids MUST NOT appear in the same batch — each group is its own
round, and rounds run one after another, in ascending order of each group's lowest task id (tasks are
already ordered by dependency). Two ids are never merged on the grounds that their tasks look
independent of each other.

**A group of one is legal, not malformed.** It yields no fan-out: that task runs through the unchanged
serial path in its own round, and nothing reports the lone member as an error.

**Malformed ⇒ unmarked.** A mark is malformed when its sub-line carries no id, when a task carries more
than one such sub-line, or when it is written on the task's `- [ ]` line instead of its own indented
sub-line. A malformed mark is never repaired by guessing the intended group: the task is handled as
unmarked — serial — and the malformation is reported, visibly, without halting the phase or the rest of
its group.

**Validated mechanically, before a batch is formed.** Before forming any batch from a tasks artifact's
marks, the orchestrator runs `~/.claude/scripts/validate-parallel-marks.js --file <tasks artifact>`
against it — same point and manner as the Return Contract Check's `validate-return.js`. The check is
read-only: it reports the three malformed shapes above (naming the offending task) and never rewrites,
normalizes or repairs a mark, and it never judges whether tasks sharing an id are genuinely independent
— that determination stays the producer's. A clean artifact produces no output; a finding degrades only
its own task to serial and is surfaced visibly, never halting the phase. If the check cannot run at all,
the orchestrator surfaces that fact, never records the artifact as clean, and forms no batch from the
marks — every task takes the serial path.

## Two runtime modes, plus serial

A dispatch is one of three shapes; the launch prompt tells a run which one it is running as — a run
never infers its own mode:

- **Isolated run.** One task, `isolation: "worktree"` set on the Task-tool launch, and a `base` sha in
  the prompt. It implements exactly that one task, produces exactly one commit, and returns a **Task
  Run Report** (below). It persists nothing: no `apply-progress`, no task marking.
- **Consolidation run.** No isolation. The prompt carries the **N Task Run Reports** the batch's
  isolated runs returned, verbatim. It integrates them one at a time, writes the single `apply-progress`
  artifact and the tasks artifact, and produces the batch's one `## Implementation Progress` return —
  same shape as today, unchanged sections (`~/.claude/references/phase-returns/sdd-apply/sdd-apply.md`).
- **Serial run.** Today's path. Not affected by anything in this file.

## Dispatch (orchestrator side, for context)

Before forming any batch, the orchestrator runs `validate-parallel-marks.js` against the tasks artifact
(see "Validated mechanically, before a batch is formed" above) — this is a precondition of everything
below, not a separate step elsewhere. For the round in play, the orchestrator reads `HEAD` of the
working branch once, records it as `base`, and dispatches every task of that group in **one message** —
N `Task` tool calls, each with its own task, `base`, and `isolation: "worktree"`. It waits for the whole
batch to return before dispatching the consolidation run (batch-bound dispatch — background
per-completion dispatch is rejected; see the EDR). Then it moves to the next group's round, in ascending
order of each group's lowest task id. The worktree shares the repository's object store and history:
isolation materializes only the versioned files, never a second copy of history.

## The base handshake (two levels)

A worktree the harness hands out is not guaranteed to start from the branch's current state — a stale or
reused directory is a base nobody can trust. Two checks close that gap, one per side:

**Level 1 — the isolated run, before writing anything:**
1. `git rev-parse HEAD` must equal the `base` sha received in the prompt.
2. `git status --porcelain` must be empty.

Either check failing → implement **nothing**. Report `Result: not-implemented`, `Why not:
base-not-established`, and return — do not attempt the task on an unverified base.

**Level 2 — the consolidation run, before trusting a report's commit:** `git rev-parse <commit>^` must
equal the same `base` the orchestrator handed to that run. A report that fails this check is handled
exactly like a cherry-pick conflict for integration purposes (see "Conflict handling" below), with the
Integration Log reason recorded as `base-mismatch` instead of `conflict` — its lineage cannot be trusted,
so it is not cherry-picked, and its worktree/branch (if the harness kept them) stay intact for
inspection.

## Isolated run: task, then one commit

1. Run the base handshake (Level 1). Fails → report `not-implemented / base-not-established`, done.
2. Implement the assigned task per Steps 2-4 of `~/.claude/skills/sdd-apply/SKILL.md` — same reading,
   same fork test, same UI-counterpart obligation. **Do not** run Steps 5/6 (mark tasks, persist
   `apply-progress`) — those belong to the consolidation run only (single-writer rule, see the EDR).
3. Commit **exactly once**, before returning control. The message follows the format, types, scope,
   subject rules and hard attribution rules of `~/.claude/skills/git/SKILL.md` **cited by reference,
   not invoked** — this run does not call that skill and does not run its interactive atomicity
   STOP-and-ask loop (there is no human to ask inside an isolated run) or its "no commit without
   compiling" gate (that check is deferred to the end-of-batch `sdd-verify`).
4. Return the **Task Run Report** (below). Never the `## Implementation Progress` block — that shape
   belongs to the consolidation run alone.

An isolated run that needs work outside its own task **reports it**; it does not dispatch anything — the
Executor boundary holds inside isolation exactly as it does outside it.

### Task Run Report

```markdown
## Task Run Report: <task-id>
**Task**: <id> — <description>
**Base**: <sha> · **Result**: committed | not-implemented · **Why not**: <line> | —
**Branch**: <name> · **Commit**: <sha> · **Worktree**: <abs path>

### Files Changed
| File | Action | What Was Done |
|------|--------|---------------|

### Unmandated Forks
{Same shape and tokens as `sdd-apply.md`'s section of the same name. If none: "None."}

### Mandated Departures
{Same shape and tokens as `sdd-apply.md`'s section of the same name. If none: "None."}

### UI Scenario Counterparts
{Only when the spec artifact carries a `ui-scenarios:` block. Absent otherwise.}

### TDD Cycle Evidence
{Only in Strict TDD Mode. Absent otherwise.}
```

Every element traces to a spec requirement — none is decorative: `Result`/`Why not` (a run that did not
persist is not registered as done), `Base`/`Commit` (the base handshake and the audit trail), `Branch`/
`Worktree` (what survives isolation and where to inspect a conflict), `Files Changed` (what the
consolidation run's artifact accumulates), the two mailboxes (deviations must survive concurrency
intact), the two conditional sections (same conditions as the non-parallel path).

For `Result: not-implemented`, `Branch`/`Commit` are `—` and `Worktree` is whatever the harness left (see
"Isolation and the worktree directory" below).

## Isolation and the worktree directory

`isolation: "worktree"` on the Task-tool launch is the confirmed harness parameter for per-run
isolation — this is a harness fact, not a design choice this file makes. The harness cleans the worktree
directory automatically **when the agent leaves no changes**. A committed run always has changes (a
commit is one), so its directory is expected to survive; a `not-implemented` run leaves none, so its
directory is the one case the harness is expected to reclaim on its own.

**The branch is the reliable survivor, not the directory.** `git worktree remove` detaches a directory
from its branch; the branch ref itself is untouched unless explicitly deleted. Report `Worktree: <abs
path>` when it is still there to report; when it is not (harness already reclaimed it), the branch is
still the thing to point someone at for inspection.

**Verified: the consolidation run CAN remove a harness-created worktree — but only via `unlock` first.**
A bare `git worktree remove <path>` on a worktree the harness created fails with exit 128: `fatal:
cannot remove a locked working tree, lock reason: claude agent agent-<id> (pid <n>) — use 'remove -f -f'
to override or unlock first`. The harness leaves every worktree it creates **locked**, not
permission-restricted — that is the actual cause, not a guess. `git worktree unlock <path>` followed by
`git worktree remove <path>` works: both exit 0, and the worktree disappears from `git worktree list`.
**`remove -f -f` is never used** — the override is not needed once `unlock` runs first, and this run has
no standing authorization for a forced removal. The branch itself survives the directory removal —
`worktree-agent-<hex>` still exists after `remove` — so deleting it is the separate, explicit `git
branch -D` step. Where and when this sequence runs: "Cleanup (after the loop, once)", below.

## Consolidation run: integrate, then write once

Receives the batch's N Task Run Reports verbatim, in the same message. For each, ascending by task id:

1. **Level-2 base check** (above). Fails → treat as a conflict for this task; do not cherry-pick;
   `base-mismatch` in the Integration Log; continue to the next report.
2. `Result: not-implemented` → nothing to integrate; the task stays incomplete with its `Why not`
   carried into `### Remaining Tasks` / `### Issues Found`; continue.
3. `Result: committed` → `git cherry-pick <sha>`.
   - **Clean.** The task lands on the working branch. Record it as integrated and move on — its
     worktree and branch are **not** touched here; removal is deferred to the single cleanup pass
     below, which runs only after the whole loop has finished.
   - **Conflict.** `git cherry-pick --abort` — the only recovery. Never `git reset --hard`, `git
     stash`, or `git checkout --`: those are destructive operations this run does not have standing
     authorization to run. Keep the branch and whatever worktree the harness left, for inspection.
     Continue to the next report — a conflict never aborts the batch and never reverts what is already
     on the branch.

Integrations run strictly one at a time, in this loop, in ascending task-id order — this is automatic
inside a single non-isolated run and needs no extra guard.

### Cleanup (after the loop, once)

Once every report has been processed — not per task as it integrates — run one cleanup pass over
exactly the tasks that landed **Clean** above. A task that conflicted, failed the Level-2 base check, or
was `not-implemented` is never touched here: its worktree and branch stay in place for inspection, same
as always. This is deliberate, not incidental: if a later cherry-pick in the loop still ends up in
conflict, every earlier worktree is still there to inspect or redo, and there is one pass to report
instead of N interleaved ones.

For each cleanly integrated task, in the same ascending task-id order:

1. `git worktree unlock <path>` — the harness leaves every worktree it creates **locked**, not
   permission-restricted; a bare `remove` on it fails with exit 128. `unlock` is the verified way past
   that lock (see "Isolation and the worktree directory" above).
2. `git worktree remove <path>` — after `unlock`, this exits 0 and the worktree disappears from `git
   worktree list`.
3. `git branch -D <branch>` — the branch survives step 2 untouched (verified: it is still listed after
   `remove`), so deleting it is its own explicit step.

**Never `git worktree remove -f -f`.** The forced override is not needed once `unlock` runs first, and
this run has no standing authorization for it. A failure at any of the three steps is not blocking — the
task is already integrated regardless of whether its worktree gets cleaned up — so record what actually
happened (`removed` / `left` / `remove failed` / `branch kept`) in the Integration Log's `Worktree kept`
column and continue to the next task.

After cleanup: run Steps 5 and 6 of `~/.claude/skills/sdd-apply/SKILL.md` **once**, over the union of
every report — mark every completed, integrated task `[x]`; merge `Files Changed`, `Unmandated Forks`,
`Mandated Departures`, the UI counterparts and the TDD evidence from every report into the
`apply-progress` artifact exactly as the Merge Protocol already describes for a normal batch, plus the
`### Integration Log` (artifact-only, cumulative):

```markdown
### Integration Log
| Order | Task | Result | Commit | Worktree kept |
|-------|------|--------|--------|----------------|
```

`Order` is the ascending task-id order this run integrated in — the requirement that the order stay
legible without re-deriving it from the branch. A task whose result was `not-implemented`,
`base-mismatch`, or a cherry-pick `conflict` gets a row here too, with `Commit: —`.

Produce exactly **one** `## Implementation Progress` return, same sections as always
(`~/.claude/references/phase-returns/sdd-apply/sdd-apply.md`): a non-integrated task is absent from
`### Completed Tasks`, stays in `### Remaining Tasks`, its conflict or mismatch is an `### Issues Found`
entry, and the batch is `partial` whenever one remains — never `done` with an unintegrated task.

## Scope

This is the **second** declared case of Phase fan-out — see `~/.claude/matecito-ai/domains/development.md`
→ "Phase fan-out (the two declared cases)". It is scoped to `sdd-apply` batches with the independence
mark; it does not extend to any other phase, and a third fan-out case is a change of its own, not an
extension of this file.
