<!-- matecito-ai: template canónico del retorno de sdd-apply.
     Existe por lo mismo que el de sdd-design, y acá el desorden era mayor: el blocker tenía TRES
     destinos posibles (`### Issues Found`, `### Status`, `risks`) y ninguno de los tres lo leía un
     guard, y `partial` significaba una cosa en la Sección D.1 y otra en la skill. Un formato con
     dos definiciones no es un formato.
     Este archivo es LA fuente del formato. La skill y el módulo strict-tdd.md lo referencian; no lo
     copian. -->

# Return template — `sdd-apply`

The exact shape of what this phase hands back to the orchestrator. The orchestrator validates the
return against this file: it matches section titles literally, so a title that differs in wording,
casing or heading level is a section it will not find.

- **Fields of the envelope**: Section D of `~/.claude/skills/_shared/sdd-phase-common.md`. Not repeated here.
- **This file**: the `detailed_report` block — which sections, in which order, with which titles, per status.

<!-- matecito-ai: parallel-batch note — added when the isolated/consolidation split landed. -->
**A batch with two or more independence-marked tasks runs a different shape for HALF of it.** An
**isolated run** does not return this block at all — it returns a **Task Run Report**, defined once in
`~/.claude/references/phase-returns/sdd-apply/parallel-batch.md`, and persists nothing. Everything below
this point — the `## Implementation Progress` block, its sections, `status` resolution, the Merge
Protocol, the artifact format — is what the **consolidation run** produces, exactly as a serial batch
always has. Read `parallel-batch.md` first if you are dispatched as, or are consolidating, a parallel
batch; this file does not repeat its content.

## Sections of this phase

| Section | Emitted | Read by |
| --- | --- | --- |
| `### Completed Tasks` | always | the orchestrator, as context |
| `### Files Changed` | always | the orchestrator, as context — and `sdd-verify` reads it from the artifact (Strict TDD coverage, and the "backing" check of `decision-gaps`, which joins by the row's `Task` column) |
| `### Decisions Materialized` | only when this run materialized ≥1 ratified proposal | the orchestrator, as context; `sdd-verify`'s `decision-gaps` group reads it from the artifact — see `~/.claude/references/decision-capture/in-flow-capture.md` |
| `### Rejected Proposals Checked` | only when the dispatch prompt forwarded ≥1 rejected proposal | Unresolved Decisions Guard — classifies each `design-conflict` verdict; a `conflicts` verdict requires `status: blocked` — see `~/.claude/matecito-ai/domains/development.md` → "Forwarding a proposal's resolution to `sdd-apply`" |
| `### TDD Cycle Evidence` | only in **Strict TDD Mode** | the orchestrator, as context; `sdd-verify` reads it from the artifact |
| `### Test Summary` | only in **Strict TDD Mode** | the orchestrator, as context |
| `### Unmandated Forks` | always | Unresolved Decisions Guard — **Tier 1** |
| `### Mandated Departures` | always | Unresolved Decisions Guard — **Tier 2** |
| `### Contract Shapes Proposed` | conditional — only when `has_contract_proposals` is true and status is `blocked` | Unresolved Decisions Guard — **Tier 1** |
| `### Blocker` | only when a blocker stopped you | the orchestrator: it puts the question to the user |
| `### Issues Found` | always | the orchestrator, as context — **never** the blocker |
| `### Remaining Tasks` | always | the orchestrator, to route: continuation batch vs verify |
| `### Workload / PR Boundary` | always | the orchestrator, to follow through on the workload decision |
| `### Status` | always | the orchestrator, to route |

Titles are fixed and this phase declares no variants of them.

**`### Files Changed` carries a `Task` column — every row is attributed to the task that produced
it.** With a single task this reads as trivially uniform; with two or more it is what makes the
table usable at all. `sdd-verify`'s `decision-gaps` group joins on it (see
`~/.claude/references/decision-capture/in-flow-capture.md`, "Backing"): it filters this table down
to the implementing task's own rows before checking whether any of them lands outside
`.matecito-ai/edr/`. A table with no `Task` column cannot support that filter once more than one
task is in play — do not drop the column, even when this batch has only one task.

The two conditional TDD sections are gated on the `**Mode**` line of the block, which is why that
line is part of the header: `Strict TDD` means both sections must be there, `Standard` means both
are legitimately absent. That is the only condition — a batch cut short still ships the evidence
table (see `~/.claude/skills/sdd-apply/strict-tdd.md`).

**`### Rejected Proposals Checked` — one item per forwarded rejection whose governed task this run
reached.** The gate field `has_forwarded_rejections` is a fact about the dispatch prompt, not the
list: a prompt that forwarded a rejection but reached no item to report renders the gate `true` with
an empty list, which the renderer shows as `None.` — that IS the declaration, not an omission. Each
item names the point the rejected proposal governs, its `record:` token, and a `design-conflict:`
verdict — `none` when the design's approach and the rejected proposal describe the same
implementation, `conflicts` when they describe different ones for the same point. A `conflicts` item
is never reported alongside `status: done` or `status: partial` with the governed task marked
complete: it means the task could not be implemented without picking a version nobody ratified, so
the return MUST be `blocked`, with `### Blocker` pointing at the conflicting item rather than
restating it (see "One blocker, one place" below). Both `none` and `conflicts` are `passing` values
for the `design-conflict` token — unlike `mandate`/`verify-checks` below, there is no non-passing
value here, because a value the renderer refuses to render would make an honest conflict
unreportable.

**Split into summary/rationale.** `### Unmandated Forks`, `### Mandated Departures`, and
`### Rejected Proposals Checked` all declare it. Each item carries two parts, `summary` and
`rationale`, in the `unmandated_forks` / `mandated_departures` / `rejected_proposals` JSON: `summary`
is what the gate prints, `rationale` is the full reasoning — always emitted into this block, never
printed by default. Both are non-empty, single-line strings; a missing one, or one with an embedded
newline, fails the render naming the item and the part, and nothing reaches stdout. `summary` also
carries a **250-character cap**, enforced by `render-return.js`. Each item also carries its declared
tokens, each on its own `· ` line, in fixed order, with `· rationale:` always last — every one of the
three sections leads with `· anchor:` (free-form, per Section D.3 of `sdd-phase-common.md`), then:
for `### Unmandated Forks` / `### Mandated Departures`, `mandate: covered|forced|chosen`, then
`verify-checks: yes|no`; for `### Rejected Proposals Checked`, `record: <domain>/<slug>`, then
`design-conflict: none|conflicts`. `summary`'s register is fixed once in Section D.3 of
`sdd-phase-common.md` — not restated here.

**`### Contract Shapes Proposed`** is the dedicated home for an unpinned contract or definition — the
shape "Contract & definition shapes — never inferred" (`~/.claude/matecito-ai/domains/development.md`)
forbids guessing. It never appears as free prose inside `### Blocker`, and it is not an
`### Unmandated Forks` item: a fork admits more than one valid resolution, an unspecified contract
admits none until the user fixes it. Emitted only when this run declares `has_contract_proposals: true`
**and** `status: blocked` — both conditions must hold, and the gate is supplied by the run, never
derived from the list being non-empty: an empty list under a `true` gate is legitimate (every proposal,
once checked, needed no further ratification) and reads differently from `false` ("nothing to
propose"). Absence is never a violation on its own; only the section appearing on a status outside
`[blocked]` is.

Each item is ONE compound entry, never split across items and never one item per field: a one-line
`summary` (the contract, what needs it, why it is unpinned), one `· anchor:`, then one
`· field: {name} — {type} — {description}` continuation line per field the contract needs — every
field proposed, none dropped, none summarized — then `· rationale:`. `summary` keeps the ordinary
250-character cap; each field's own `description` carries its own 160-character cap — the **field
count is never capped**. One anchor per contract, never per field: it is also the identity the forward
uses when the ratified shape comes back. The item carries no field for where the shape will be stored
or persisted — this is the one proposing phase where that omission matters most, because "where it
lives" for `sdd-apply` IS the code, and this section is not where that code is written.

Once ratified (or adjusted) at the gate, the field list travels back only in this batch's re-dispatch
prompt, identified by the item's own `anchor` and `summary` — **never re-read from `apply-progress` or
any other stored artifact**. When the prompt carries a ratified shape, implement it as the code it
governs: no separate return section records it once written, because the code IS the record, per the
domain rule. A batch that reaches the governed point with nothing in its prompt MUST stop and name the
missing contract, exactly as "Consulting an Unmandated Fork" governs any other point no artifact fixes.

## Which status

Resolve in this order, top down, and stop at the first that fits:

1. **`blocked`** — you cannot continue without a resolution that is not yours to make. Emit `### Blocker`.
2. **`partial`** — the phase is not finished: tasks of this change remain. This is the normal
   continuation batch as much as it is the stop-with-a-blocker case; `partial` does NOT claim
   anything about blockers either way, and `### Blocker` is emitted only if one actually stopped you.
3. **`done`** — nothing remains: `### Remaining Tasks` is `None`.

<!-- matecito-ai: `partial` = "la fase no terminó", el sentido de D.1. NO significa "hay un blocker
     que no frenó todo": esa lectura dejaba sin status al batch de continuación normal, que es el
     caso más frecuente de esta fase. Si hay blocker, se reporta en su sección, con cualquier status. -->

Whatever the status, the completed work is persisted BEFORE you return (Steps 5 and 6 of the skill).
There is no exit path from this phase that skips them.

## `status: done` — the batch finished and nothing remains

```markdown
## Implementation Progress

**Change**: {change-name}
**Mode**: {Strict TDD | Standard}

### Completed Tasks
- [x] {task id and description}
- [x] {task id and description}

### Files Changed
| File | Task | Action | What Was Done |
|------|------|--------|---------------|
| `path/to/file.ext` | {task id} | Created | {brief description} |
| `path/to/other.ext` | {task id} | Modified | {brief description} |

{CONDITIONAL — only when this run materialized at least one ratified decision proposal; omit the
whole section otherwise, do not even print "None." Full mechanism:
`~/.claude/references/decision-capture/in-flow-capture.md`.}

### Decisions Materialized
| Record | Task | Result |
|--------|------|--------|
| `<domain>/<slug>` | {task id} | materialized \| failed: {reason} |

{The next two sections are Strict TDD Mode ONLY — in Standard Mode omit both, entirely.
What each column means, and what counts as evidence, is in `~/.claude/skills/sdd-apply/strict-tdd.md`.}

### TDD Cycle Evidence
| Task | Test File | Layer | Safety Net | RED | GREEN | TRIANGULATE | REFACTOR |
|------|-----------|-------|------------|-----|-------|-------------|----------|
| 1.1 | `path/test.ext` | Unit | ✅ 5/5 | ✅ Written | ✅ Passed | ✅ 3 cases | ✅ Clean |
| 1.2 | `path/test.ext` | Integration | N/A (new) | ✅ Written | ✅ Passed | ➖ Single | ✅ Clean |
| 1.3 | `path/test.ext` | Unit | ✅ 2/2 | ✅ Written | ✅ Passed | ✅ 2 cases | ➖ None needed |

### Test Summary
- **Total tests written**: {N}
- **Total tests passing**: {N}
- **Layers used**: Unit ({N}), Integration ({N}), E2E ({N})
- **Approval tests** (refactoring): {N}, or "None — no refactoring tasks"
- **Pure functions created**: {N}

<!-- matecito-ai: two Tier sections replace the old single mailbox — the guard classifies by section,
     so an unresolved fork routes to Tier 1 (consult before applying) and everything you did apply
     routes to Tier 2 (surface, don't block). Each item's tokens are literal, not prose — a phrase the
     reader has to interpret is how "no declaration" and "declared clean" become indistinguishable. -->
### Unmandated Forks
{One entry per point the confirmed artifacts do not fix, where more than one resolution was valid.
You applied NONE of them — the point stays untouched, its task stays open in `### Remaining Tasks`,
and the fork travels back as a question before the next dispatch:

- {the fork, phrased so the user can answer it without reading the rest — naming the task it blocks
  — plus your recommended resolution, if you have one}
  · anchor: {the concrete source this fork is about — a `<repo-path>[:line]` or `<engram-key>`}
  · mandate: chosen
  · verify-checks: yes|no
  · rationale: {the resolutions that were valid, and why you recommend the one you do}

`mandate: chosen` is the only legal value in this section — an item whose mandate is `covered` or
`forced` belongs in `### Mandated Departures` instead.
If none: "None."}

### Mandated Departures
{One entry per point where the implementation departs from the design and you DID apply it without
consulting — because a confirmed artifact already fixed the point (`mandate: covered`), or no
alternative was valid and you can name the concrete constraint that closed the others
(`mandate: forced`):

- {what you applied — for `forced`, name the constraint: an upstream artifact, an Accepted decision
  record, or what leaving the point untouched broke}
  · anchor: {the concrete source this deviation is about — a `<repo-path>[:line]` or `<engram-key>`}
  · mandate: covered|forced
  · verify-checks: yes|no
  · rationale: {one line: the full reasoning behind the deviation — always emitted here, never printed at the gate by default}

`verify-checks:` is a **literal token, not prose** — `sdd-verify` matches it to classify the
deviation. Answer `yes` when the deviation touches a spec requirement or a task `criteria:`: the
design artifact is now stale and verify will read it as a mismatch. Answer `no` when it is pure
execution detail that nothing verifiable covers. You are the only one with the context to say
which; the orchestrator is not, and verify will not guess.
**A deviation with no `verify-checks:` token, or a hedged one, is treated as `yes`** — the strict
default, so forgetting to declare is never the cheap way past the gate. **A missing or hedged
`mandate:` is treated as `chosen`** — an item you cannot back with a named `covered` or `forced`
constraint was never legally absorbed; it belongs in `### Unmandated Forks` instead.
If none: "None."}

### Issues Found
{Problems discovered while implementing that did NOT stop you: a fragile test, a pre-existing
defect you did not touch, a slow path worth revisiting. Context, not a gate.
If none: "None."}

### Remaining Tasks
None — all tasks complete.

### Workload / PR Boundary
- Mode: {single PR | chained PR slice | stacked PR slice | size:exception}
- Current work unit: {unit name, or "N/A"}
- Boundary: {what this batch starts from and ends with}
- Estimated review budget impact: {brief note}

### Status
{N}/{N} tasks complete. Progress persisted: yes.
Ready for verify.
```

## `status: partial` — real work landed, the phase is not finished

Same block. `### Remaining Tasks` is what makes it `partial`, and `### Blocker` appears only if a
blocker is what stopped you — a batch that simply ran out of assigned tasks has none.

```markdown
## Implementation Progress

**Change**: {change-name}
**Mode**: {Strict TDD | Standard}

### Completed Tasks
- [x] {every task you actually finished — this is what the next batch will skip}

### Files Changed
{As above, over what this batch touched.}

{Strict TDD Mode only → `### TDD Cycle Evidence` and `### Test Summary`, exactly as in the `done`
block. On an early exit the table still ships: one row per completed cycle, plus a row for the task
you stopped on showing the last stage it reached (e.g. RED written, GREEN never run).}

### Unmandated Forks
{As above.}

### Mandated Departures
{As above.}

### Blocker
{ONLY if a blocker stopped you. Omit the whole section otherwise.}

{The gap, in one line, phrased so the user can answer it without reading the rest.}

**Options weighed**: {each one, with what it costs and what it buys.}

**Why this is not mine to settle**: {it is a decision, not execution detail — a contract, a new
dependency, which layer the logic lives in, a change to the design's approach. If an upstream
artifact pins it but is inconsistent or does not cover this case, state both sides: what it fixes
and why, and what it breaks here.}

**What unblocks it**: {the answer you need.}

**Task stopped on**: {task id}{, and the stage it reached, in Strict TDD Mode}

### Issues Found
{As above — problems that did NOT stop you. The blocker is not one of them.}

### Remaining Tasks
- [ ] {task id and description}
- [ ] {task id and description}

### Workload / PR Boundary
{As above.}

### Status
{N}/{total} tasks complete. Progress persisted: yes.
{Ready for next batch — OR: stopped at task {id}, see Blocker.}
```

## `status: blocked` — you cannot continue at all

Same block as `partial`, with `### Blocker` unconditional and no route out (`next_recommended: none`,
per Section D.4). Use it when the blocker stops the rest of the batch too, not only the task you were
on. Everything you did finish is still reported and still persisted — `blocked` is not "nothing
happened".

```markdown
## Implementation Progress

**Change**: {change-name}
**Mode**: {Strict TDD | Standard}

### Completed Tasks
{What you finished before the stop. If nothing: "None — the stop happened before any task completed."}

### Files Changed
{What you touched, or "None."}

{Strict TDD Mode only → `### TDD Cycle Evidence` and `### Test Summary`, as in `partial`.}

### Unmandated Forks
{As above.}

### Mandated Departures
{As above.}

{CONDITIONAL — only when the dispatch prompt forwarded ≥1 rejected proposal; omit the whole section
otherwise, do not even print "None." Full mechanism:
`~/.claude/matecito-ai/domains/development.md` → "Forwarding a proposal's resolution to
`sdd-apply`". This is the shape a `conflicts` verdict takes — the reason this batch is `blocked`:}

### Rejected Proposals Checked
- {the point the rejected proposal governs, and the verdict, in one line}
  · anchor: {the concrete source this item is about — a `<repo-path>[:line]` or `<engram-key>`}
  · record: <domain>/<slug>
  · design-conflict: conflicts
  · rationale: {the design's approach vs. what the rejected proposal proposed, one line}

{CONDITIONAL — only when `has_contract_proposals` is true; omit the whole section otherwise, do not
even print "None." This is the shape a stop over an unspecified contract takes — the reason this batch
is `blocked`:}

### Contract Shapes Proposed
- {the contract, what needs it, why it is unpinned — one line, ≤250}
  · anchor: {the concrete source this contract is about — a `<repo-path>[:line]` or `<engram-key>`}
  · field: {name} — {type} — {the field's description; every field proposed, none dropped}
  · field: {name} — {type} — {a second field, same shape; as many lines as the contract needs}
  · rationale: {one line: what a wrong guess here propagates to}

### Blocker
{Same shape as in `partial`. Always present here.}

### Issues Found
{As above.}

### Remaining Tasks
- [ ] {everything left, including the task you stopped on}

### Workload / PR Boundary
{As above.}

### Status
{N}/{total} tasks complete. Progress persisted: {yes | no — nothing had been completed yet}.
Stopped at task {id} — see Blocker.
```

## One blocker, one place

<!-- matecito-ai: el blocker tenía tres destinos y ningún consumidor. Se fija UNO, igual que en
     sdd-design: la sección dedicada. Los otros dos lugares quedan explícitamente prohibidos, si no
     la copia vuelve a aparecer sola. -->

`### Blocker` is where the blocker is stated, and nowhere else:

- **not `### Issues Found`** — that section is for problems you worked around or left standing; a
  gate that has to guess which of its items is the blocker is not a gate;
- **not `### Status`** — it may point at the blocker (`Stopped at task 2.3 — see Blocker`) but never
  restates it;
- **not `risks`** — Section D.4 forbids routing a decision the user owns through that field.

Stating the same blocker twice is the same failure as stating it nowhere: the two copies drift and
then nobody knows which one is current.

## Artifact vs return — do not confuse them

The `apply-progress` **artifact** (persisted to Engram) is the cumulative state of the change across
ALL batches — every completed task from every batch, merged. The **return** is this block: what THIS
batch did and what the orchestrator needs in order to gate and to route. The orchestrator never reads
the artifact; the next `sdd-apply` batch and `sdd-verify` read only the artifact.

<!-- matecito-ai: el formato del artefacto no estaba declarado en NINGÚN lado — este archivo decía
     "format in the skill" y la skill no lo define; su Merge Protocol pide "keep the same structure"
     de una estructura inexistente. Peor: `### Deviations from Design` figuraba sólo como sección del
     RETORNO, y `sdd-verify` la lee del ARTEFACTO. O sea que el token `verify-checks:` —el único
     mecanismo que baja un desvío de CRITICAL a WARNING— viajaba por un canal que su lector no
     recibe. Se declara acá porque este archivo es la fuente del formato de esta fase. -->
### Format of the `apply-progress` artifact

Four sections are **mandatory** in the persisted artifact, and all four are cumulative — each batch
merges its own into what previous batches left:

```markdown
## Apply Progress: {change-name}

### Tasks
{Every task of the change, from every batch, with its state: `[x]` done, `[ ]` pending.
Never drop a task a previous batch completed — that is what makes the next batch skip it.}

### Files Changed
{Every file touched, from every batch, with what was done to it — File, Task, Action, What Was
Done, same shape as the return's version of this section. `sdd-verify` reads THIS copy to scope its
Strict TDD coverage check, and its `decision-gaps` group joins on the `Task` column to find the
implementing task's own rows for the "backing" check (see
`~/.claude/references/decision-capture/in-flow-capture.md`) — it is not only orchestrator context.
A batch that merges rows from more than one task (a consolidation run integrating N Task Run
Reports) MUST keep each row's `Task` column intact through the merge — never collapse it away.}

### Unmandated Forks
{Every fork from every batch that reached this point unresolved, each in the exact shape this
template declares above, including its `· anchor:`, `mandate: chosen` and `verify-checks: yes|no`
tokens and its `· rationale:` line. `sdd-verify` reads THIS copy — not the return — alongside
`### Mandated Departures` below, to find every declared deviation.}

### Mandated Departures
{Every deviation from every batch that you applied without consulting, each in the exact shape this
template declares above, including its `· anchor:`, `mandate: covered|forced` and
`verify-checks: yes|no` tokens and its `· rationale:` line. `sdd-verify` reads THIS copy — not the
return — to classify each
deviation. A deviation that reaches the orchestrator and not the artifact is a deviation `sdd-verify`
will never see; a rationale that never left the return is a rationale nobody persisted.}

<!-- matecito-ai: artifact-only, and only present when a parallel batch actually ran — it is the
     integration record `parallel-batch.md` promises, not something a serial batch produces or needs.
     Absence here is legitimate for a serial batch, the same way the UI/TDD sections below are
     legitimate absences outside their own conditions. -->
### Integration Log
{CONDITIONAL — only when this artifact accumulates the result of at least one parallel batch; absent
otherwise. Cumulative across every parallel batch this change has run, in the shape
`~/.claude/references/phase-returns/sdd-apply/parallel-batch.md` fixes:

| Order | Task | Result | Commit | Worktree kept |
|-------|------|--------|--------|----------------|

`Order` is the ascending task-id order the consolidation run actually integrated in — what makes the
order legible without re-deriving it from the branch. A task whose result was `not-implemented`,
`base-mismatch`, or a cherry-pick conflict gets a row too, with `Commit: —`.}

<!-- matecito-ai: artifact-only section, and deliberately so — it is INPUT for sdd-verify's browser run,
     not something the orchestrator gates on. It lives here because the spec authors UI scenarios in
     domain language (it cannot know a route or an accessible name before the code exists, and must not
     pin them anyway) while sdd-verify needs exact targets to stay deterministic. This phase is the only
     one that knows them without guessing. -->
### UI Scenario Counterparts
{CONDITIONAL — only when the spec artifact carries a `ui-scenarios:` block; absent otherwise, and that
absence is legitimate. One executable counterpart per behavioral scenario, cumulative across batches,
in the shape Part 2 of `~/.claude/references/ui-scenarios-schema.md` fixes: the same `name` verbatim,
the real `url`, the `steps`, the `expect`. `sdd-verify` pairs them by `name` and checks coverage — a
behavioral scenario with no counterpart is `UNTESTED` and CRITICAL, and a missing section while the
spec declares scenarios is a CRITICAL finding, never a silent skip. Artifact only: it is input for the
browser run, not a section the orchestrator gates on.}

### TDD Cycle Evidence
<!-- matecito-ai: estaba sólo como prosa dentro de las llaves, no como sección declarada. Su
     consumidor la busca por título y marca CRITICAL si no la encuentra: mismo patrón que
     `### Files Changed`, un escalón más abajo. -->
{CONDITIONAL — only when Strict TDD is active; absent otherwise, and that absence is legitimate.
One row per task: RED (test written first) → GREEN (implementation passes) → REFACTOR.
Cumulative across batches, like the sections above. `sdd-verify` enforces its hard gate against
THIS copy and marks CRITICAL when it is missing while Strict TDD was on — see
`~/.claude/skills/sdd-verify/strict-tdd-verify.md`.}

<!-- matecito-ai: in-flow decision capture (development-specifics). Artifact-only in the sense that
     this is the copy `sdd-verify` reads — it is ALSO emitted in the return, conditionally, unlike
     UI Scenario Counterparts which lives only here. Full mechanism: in-flow-capture.md. -->
### Decisions Materialized
{CONDITIONAL — only when at least one batch of this change materialized a ratified decision proposal;
absent otherwise, and that absence is legitimate (a change that opened no decision materializes
nothing). Cumulative across every batch, same `record | task | result` shape as the return:

| Record | Task | Result |
|--------|------|--------|
| `<domain>/<slug>` | {task id} | materialized \| failed: {reason} |

`sdd-verify`'s `decision-gaps` group reads THIS copy — not the return — to build its check list. A
`failed` row is never dropped: it is exactly the "propuesta ratificada sin registro materializado"
case that group flags CRITICAL.}
```

**`### Unmandated Forks` and `### Mandated Departures` go in BOTH places**, with the same content: in
the return, because the Unresolved Decisions Guard reads them there — Tier 1 and Tier 2
respectively; in the artifact, because `sdd-verify` reads both there to find every declared deviation
and apply the `verify-checks` classification. Same reason the TDD evidence lives in both. Dropping
either copy breaks a different consumer, and neither failure is loud.

**`### Rejected Proposals Checked` is return-only — it does NOT go in the artifact.** Its one
consumer is the orchestrator's Unresolved Decisions Guard, at the moment this batch's return is
gated: a `conflicts` verdict blocks the governing task, so that task never reaches `[x]` and is never
part of the cumulative task list the artifact carries. There is nothing durable this section would
add to `apply-progress` that the artifact's own `### Tasks` state does not already say by the
governed task staying `[ ]`.
