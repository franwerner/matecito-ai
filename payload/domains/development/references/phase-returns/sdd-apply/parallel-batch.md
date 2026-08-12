<!-- matecito-ai: sole home of the sdd-apply parallel-batch mechanism — eligibility, isolation, the
     uncommitted-work gate, repositioning onto the local base, the base handshake, the commit
     convention, the Task Run Report shape, integration mechanics and conflict handling. This is the
     SECOND declared case of Phase fan-out (see
     `~/.claude/matecito-ai/domains/development.md` → "Phase fan-out"); the first is
     `subverifier-groups.md`. No `.yaml` pair: the Task Run Report is a bespoke block nothing renders
     or validates mechanically, exactly like `sdd-verify`'s Sub-Report — `sdd-apply.yaml` stays the
     shape of the CONSOLIDATED `## Implementation Progress` block and is unchanged by this file. -->

# `sdd-apply` parallel batch

Applies to a tasks artifact whenever at least one `· parallel-group: <id>` id has **two or more**
members — eligibility is stated and evaluated **per group**, never over the artifact as a whole (four
marked tasks split into two groups of two are two eligible batches run in successive rounds, not one
eligible batch of four; see "One group is one batch" below). A group with fewer than two members, an
unmarked task, or no tasks artifact at all (`reduced`/`custom` lanes) runs **serial mode**: today's
single-agent path, unchanged by this file. Serial mode is not documented here because nothing about it
changed.

A dirty `git status --porcelain` at dispatch time no longer degrades a round to serial by itself — see
"Uncommitted-Work Gate" below. Its three outcomes decide the round's fate; "work on the branch, no
worktree" is the only one of them that reaches serial mode, and it is reached only by the user's choice,
never automatically.

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
below, not a separate step elsewhere. For each eligible round, immediately after that validation and
before capturing `base`, the orchestrator runs the **Uncommitted-Work Gate** (below). Once the gate
clears — silently, or through a chosen outcome — the orchestrator reads `HEAD` of the working branch,
records it as `base`, and dispatches every task of that group in **one message** — N `Task` tool calls,
each with its own task, `base`, and `isolation: "worktree"`. It waits for the whole batch to return
before dispatching the consolidation run (batch-bound dispatch — background per-completion dispatch is
rejected; see the EDR). Then it moves to the next group's round, in ascending order of each group's
lowest task id. The worktree shares the repository's object store and history: isolation materializes
only the versioned files, never a second copy of history.

## Uncommitted-Work Gate

Before dispatching a round that will use worktree isolation — after `validate-parallel-marks.js`, before
the orchestrator captures `HEAD` as that round's `base` — the orchestrator inspects the main repo's
uncommitted changes. It runs once per **eligible round**, not once per phase run (`base` is captured per
round, so a single phase-wide pass would miss a round that got dirtied mid-phase), and it does not
re-prompt within the same phase run when the dirty set is unchanged from its last check. A serial
dispatch and the consolidation run never trigger it — no worktree is in play, so there is nothing to
warn about.

**Legal inputs, closed list: the dirty file list, `repo.components[].paths`
(`~/.claude/references/repo-components/README.md`), and the brief's `Components` line.** Task titles are
never parsed for paths — they are prose, not a field, and cross-referencing them against dirty files is
explicitly out of scope.

**Entry.** `git status --porcelain --untracked-files=all`, run from the main repo's root.
`--untracked-files=all` is deliberate over the default: an untracked (`??`) file is exactly the kind of
change that never travels into a worktree (see "Repositioning onto `base`" below), so it is exactly what
this gate must not miss. A rename takes its new path. When the dirty (and especially untracked) list is
large, the warning and the notice below print the first 20 entries plus a count of the rest — this caps
what prints, never the mapping or the decision to warn.

**Mapping (only with `repo.components` declared).** A dirty file belongs to **every** component whose
`paths` is a segment-prefix of its path — multivalued, no arbitration between matches. A file matching no
component is an **orphan**.

**Decision:**
- **No `repo.components` declared** → always warn when there are dirty files, listing them with no
  attempt to map them to anything, and present the three outcomes below.
- **`repo.components` declared, and the mapped components intersect the brief's `Components` line** →
  warn, naming the dirty files and the components that intersected, and present the three outcomes.
- **`repo.components` declared, no intersection, but there is an orphan file, or the brief's
  `Components` is `unassigned`, or it names something `repo.components` does not declare** → a **soft
  notice** mentioning it. It neither blocks nor forces a choice among the three outcomes — the decision
  stays the user's, made freely, not funneled through this gate.
- **`repo.components` declared, no intersection, no orphan, and `Components` names something that does
  exist** → total silence — not even a mention that the gate ran.

**Three outcomes — nothing dispatches until the user picks one, not even in Automatic mode:**

| Outcome | What happens |
|---|---|
| **Commit first** | The user commits. The orchestrator re-reads the working branch's `HEAD` and uses that new sha as this round's `base`; the tree is clean and the round dispatches isolated, normally. If the tree is still dirty and still intersects after the commit, the gate re-runs and presents the three outcomes again. |
| **Continue anyway** | The round dispatches isolated as-is, `base` already captured, the warning understood. The consolidation run records the notice (below). |
| **Work on the branch, no worktree** | The round degrades to the serial path on the working branch — no worktree, no fan-out, uncommitted work available to it. No isolated run is dispatched for this round. |

Silence from the user is not consent: with no answer, nothing dispatches — not isolated, not serial, not
even in Automatic mode.

**Notice.** Only "continue anyway" leaves a trace — commit-first leaves its commit as the record, and
"no worktree" removes the risk instead of accepting it. The orchestrator carries what triggered the
notice (the round, the dirty files, the intersecting components) in the consolidation run's launch
prompt; the consolidation run alone writes it (single-writer rule, see the EDR), artifact-only — never
the rendered return, never `sdd-apply.yaml` — as a cumulative section of `apply-progress`:

```markdown
### Uncommitted-Work Notice
| Round | Dirty files | Components | Choice |
|-------|-------------|------------|--------|
| <group id> | <paths, comma-separated> | <components that intersected, or —> | continue-anyway |
```

Only "continue anyway" adds a row — a clean round, one resolved by commit-first, or one degraded to
serial, adds none.

## Repositioning onto `base` (before the handshake)

The harness's worktree starting point is not the working base — it can be `origin/<branch>`, behind the
working branch's local `HEAD` whenever there are unpushed local commits, which is the common case. An
isolated run MUST reposition its worktree onto the `base` sha it received **before** running the Level 1
handshake and before writing anything.

**`<branch>` is the worktree's own branch — never the working branch.** Read it from inside the
worktree, before repositioning: `git rev-parse --abbrev-ref HEAD` (the harness already checked the
worktree out onto an ephemeral branch of its own, typically `worktree-agent-<hex>`). The mechanism:
`git checkout -B <branch> <base>`, using that same name — it works because the worktree shares the
repository's object store, so the `base` sha is reachable from inside it. Repositioning does not replace
the handshake; it precedes it — the handshake still runs afterward exactly as it always has.

**Using the working branch's name instead (e.g. `main`) is a bug, not a variant — worktrees share the
repository's ref store.** `git checkout -B <working-branch> <base>` run inside a worktree resets that
ref **globally**: every worktree checked out on that same branch name (including the main working
directory) moves with it, and any commit that was only reachable through the branch's old tip becomes
orphaned — no ref points at it any longer, even though `git worktree list` still shows the other
worktrees "on" that name, now silently repositioned to `base` alongside this one. This is not a
theoretical risk: it is exactly what a run that skips "read your own branch name first" produces.
Determine `<branch>` from the worktree itself, every time — never assume it, never reuse the name of the
branch you are trying to reach.

Repositioning onto a base the worktree already sits on is inert: it neither fails nor changes the
worktree's content, and the handshake still passes.

**A failed repositioning implements nothing.** Base unreachable, a git error, any reason — the run MUST
NOT implement anything, MUST NOT continue on the harness's starting point, and MUST NOT retry on an
unverified base. Report `Result: not-implemented`, `Why not: base-not-established` — the same value the
handshake itself produces, no new vocabulary — with the concrete failure reason left legible on that
line.

**This obligation belongs to every isolated run, not only the parallel batch.** The parallel batch is
today's only invoker; that is a fact about the moment, not the scope of the rule. Any future mechanism
that adopts worktree isolation without touching this file inherits the same repositioning,
unconditionally.

**Uncommitted work in the main repo never reaches an isolated run — before or after repositioning.** A
worktree is a checkout of a commit: an uncommitted new file is simply absent inside it, an uncommitted
modification to a versioned file reads as its last committed content, and `git status --porcelain` inside
the isolated run's worktree comes back empty regardless of what the main repo's tree looks like. This is
what makes the handshake's "clean tree" half satisfiable by construction, and it is the premise the
Uncommitted-Work Gate (above) acts on.

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

1. Reposition the worktree onto `base` — reading `<branch>` from inside the worktree first, never the
   working branch's name (see "Repositioning onto `base` (before the handshake)" above). Fails → report
   `not-implemented / base-not-established`, done.
2. Run the base handshake (Level 1). Fails → report `not-implemented / base-not-established`, done.
3. Implement the assigned task per Steps 2-4 of `~/.claude/skills/sdd-apply/SKILL.md` — same reading,
   same fork test, same UI-counterpart obligation. If this task carries a ratified decision proposal
   (forwarded in the launch prompt), materialize it in this same step per "Materializing decision
   records" below — the record's `.md` body ships **inside this run's one commit**, but neither INDEX
   file is touched here. **Do not** run Steps 5/6 (mark tasks, persist `apply-progress`) — those belong
   to the consolidation run only (single-writer rule, see the EDR).
4. Commit **exactly once**, before returning control. The message follows the format, types, scope,
   subject rules and hard attribution rules of `~/.claude/skills/git/SKILL.md` **cited by reference,
   not invoked** — this run does not call that skill and does not run its interactive atomicity
   STOP-and-ask loop (there is no human to ask inside an isolated run) or its "no commit without
   compiling" gate (that check is deferred to the end-of-batch `sdd-verify`).
5. Return the **Task Run Report** (below). Never the `## Implementation Progress` block — that shape
   belongs to the consolidation run alone.

An isolated run that needs work outside its own task **reports it**; it does not dispatch anything — the
Executor boundary holds inside isolation exactly as it does outside it.

## Materializing decision records (single-writer split, per `structure/root-index-cardinality-per-domain-type.md`)

Full mechanism (proposal shape, ratification gate, the two `render-artifact.js` invocations, `sdd-verify`'s
checks): `~/.claude/references/decision-capture/in-flow-capture.md`. This section fixes only the split
between the two roles, per the design decision `structure/root-index-cardinality-per-domain-type.md`
governs — **not** `contracts/single-writer-per-batch.md`, whose scope is `apply-progress` and task state:

- **Isolated run** — writes only the record's `.md` body to `.matecito-ai/edr/<domain>/<slug>.md`,
  inside its one commit, beside the code the decision governs (`contracts/one-commit-per-isolated-run.md`).
  It does **not** touch `.matecito-ai/edr/<domain>/INDEX.md` or the root `.matecito-ai/edr/INDEX.md` —
  it computes the `--index-entries` rows (second `render-artifact.js` invocation, no write either) and
  carries them, unapplied, in its Task Run Report's `### Decisions Materialized` section.
- **Consolidation run** — after the cherry-pick loop and the cleanup pass (below), applies every
  carried-forward `--index-entries` row from every cleanly-integrated report **exactly once**, deduping
  the root INDEX row by `domain` (two isolated tasks materializing in the same domain contribute one
  root row, not two — the domain INDEX itself gets one row per record, never deduped, since two
  different `slug`s are two different rows). Scaffolds an absent `INDEX.md` (domain or root) from
  `references/edr/templates/index-domain.md` / `index-root.md` on first write. A report whose
  `Result: not-implemented`, whose Level-2 base check failed, or whose cherry-pick conflicted
  contributes **no** INDEX row — its record body never reached the working branch, so indexing it would
  create a dangling entry.
- **Serial mode** — does both in the same step: writes the body and applies the INDEX rows, since there
  is no isolation split to observe.

Nothing is ever lost or duplicated across a parallel batch: every cleanly-integrated task's record body
already landed via its cherry-picked commit; the one remaining write (the INDEX rows) happens exactly
once, by the run this split already names as the sole writer of everything else this phase persists.

### Task Run Report

```markdown
## Task Run Report: <task-id>
**Task**: <id> — <description>
**Base**: <sha> · **Result**: committed | not-implemented · **Why not**: <line> | —
**Branch**: <name> · **Commit**: <sha> · **Worktree**: <abs path>

### Files Changed
| File | Task | Action | What Was Done |
|------|------|--------|---------------|
{Task is this report's own `<task-id>` on every row — trivial within a single-task report, but the
consolidation run copies rows straight into `apply-progress`'s cumulative `### Files Changed` without
re-deriving attribution, and that copy is what `sdd-verify`'s `decision-gaps` group joins on. Fill it
in here rather than leaving it for the consolidation run to guess.}

### Unmandated Forks
{Same shape and tokens as `sdd-apply.md`'s section of the same name. If none: "None."}

### Mandated Departures
{Same shape and tokens as `sdd-apply.md`'s section of the same name. If none: "None."}

### UI Scenario Counterparts
{Only when the spec artifact carries a `ui-scenarios:` block. Absent otherwise.}

### TDD Cycle Evidence
{Only in Strict TDD Mode. Absent otherwise.}

### Decisions Materialized
{Only when this task materialized a ratified decision proposal — see
`~/.claude/references/decision-capture/in-flow-capture.md`. Absent otherwise. One row per proposal this
task's step materialized:}
| Record | Result |
|--------|--------|
| <domain>/<slug> | materialized \| failed: <reason> |
{Carries the record's rendered body already written to disk (inside this run's one commit) PLUS the
`--index-entries` JSON rows for it, held here rather than applied — the isolated run never touches
either INDEX file (see "Materializing decision records" below).}
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

**Verified: the consolidation run CAN remove a harness-created worktree — but only via `unlock` first,
and `unlock` succeeding does not guarantee `remove` will.** A bare `git worktree remove <path>` on a
worktree the harness created fails with exit 128: `fatal: cannot remove a locked working tree, lock
reason: claude agent agent-<id> (pid <n>) — use 'remove -f -f' to override or unlock first`. The harness
leaves every worktree it creates **locked**, not permission-restricted — that is the actual cause, not a
guess. `git worktree unlock <path>` clears that lock. The following `git worktree remove <path>` is a
**second, independent** check: verified, it can still fail on its own, exit 128, with a **different**
message — `fatal: '<path>' contains modified or untracked files, use --force to delete it` — when the
worktree carries modified or untracked files at the moment of removal (a run's own build output, an edit
left behind outside the one commit it made). This is not the lock failure recurring; `unlock` already
succeeded, and repeating it does not help. When neither failure fires, `unlock` then `remove` both exit
0 and the worktree disappears from `git worktree list`.
**`remove -f -f` is never used, for either failure** — the override is not needed once `unlock` runs
first, and this run has no standing authorization for a forced removal; a `remove` that still fails is
recorded, not forced past. The branch itself survives the directory removal — `worktree-agent-<hex>`
still exists after a successful `remove` — so deleting it is the separate, explicit `git branch -D` step.
Where and when this sequence runs: "Cleanup (after the loop, once)", below.

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
2. `git worktree remove <path>` — after a successful `unlock`, this **usually** exits 0 and the worktree
   disappears from `git worktree list`. It can still fail on its own, exit 128, `fatal: '<path>' contains
   modified or untracked files, use --force to delete it` — a **second, distinct** failure mode from the
   lock (see "Isolation and the worktree directory" above), reachable even though `unlock` succeeded.
3. `git branch -D <branch>` — the branch survives step 2 untouched whether step 2 succeeded or failed
   (verified: it is still listed after a successful `remove`), so deleting it is its own explicit step.

**Never `git worktree remove -f -f`, for either failure mode.** The forced override is not needed once
`unlock` runs first, and this run has no standing authorization for it — a `remove` that still fails
after a successful `unlock` (lock persisting, or leftover modified/untracked files) is recorded, not
forced past. A failure at any of the three steps is not blocking — the task is already integrated
regardless of whether its worktree gets cleaned up — so record what actually happened (`removed` /
`left` / `remove failed` / `branch kept`) in the Integration Log's `Worktree kept` column and continue to
the next task.

**Immediately after cleanup, apply the carried-forward decision-record INDEX rows — once.** For every
cleanly-integrated report that carried a `### Decisions Materialized` section, apply its
`--index-entries` rows to the domain and root `.matecito-ai/edr/INDEX.md` files exactly as described in
"Materializing decision records" above (dedup the root row by `domain`; scaffold an absent INDEX from
the templates). This is a single pass, not one write per report.

After cleanup: run Steps 5 and 6 of `~/.claude/skills/sdd-apply/SKILL.md` **once**, over the union of
every report — mark every completed, integrated task `[x]`; merge `Files Changed`, `Unmandated Forks`,
`Mandated Departures`, the UI counterparts, the TDD evidence, and `### Decisions Materialized` from every
report into the `apply-progress` artifact exactly as the Merge Protocol already describes for a normal
batch, plus the `### Integration Log` (artifact-only, cumulative):

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
