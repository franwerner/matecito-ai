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

## Sections of this phase

| Section | Emitted | Read by |
| --- | --- | --- |
| `### Completed Tasks` | always | the orchestrator, as context |
| `### Files Changed` | always | the orchestrator, as context — and `sdd-verify` reads it from the artifact (Strict TDD coverage) |
| `### TDD Cycle Evidence` | only in **Strict TDD Mode** | the orchestrator, as context; `sdd-verify` reads it from the artifact |
| `### Test Summary` | only in **Strict TDD Mode** | the orchestrator, as context |
| `### Deviations from Design` | always | Unresolved Decisions Guard — **Tier 2** |
| `### Blocker` | only when a blocker stopped you | the orchestrator: it puts the question to the user |
| `### Issues Found` | always | the orchestrator, as context — **never** the blocker |
| `### Remaining Tasks` | always | the orchestrator, to route: continuation batch vs verify |
| `### Workload / PR Boundary` | always | the orchestrator, to follow through on the workload decision |
| `### Status` | always | the orchestrator, to route |

Titles are fixed and this phase declares no variants of them.

The two conditional TDD sections are gated on the `**Mode**` line of the block, which is why that
line is part of the header: `Strict TDD` means both sections must be there, `Standard` means both
are legitimately absent. That is the only condition — a batch cut short still ships the evidence
table (see `~/.claude/skills/sdd-apply/strict-tdd.md`).

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
| File | Action | What Was Done |
|------|--------|---------------|
| `path/to/file.ext` | Created | {brief description} |
| `path/to/other.ext` | Modified | {brief description} |

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

<!-- matecito-ai: la declaración era prosa libre y `sdd-verify` tenía que interpretarla — una frase
     ambigua caía en "no declarado" y escalaba a CRITICAL. El token la vuelve determinista: el
     ejecutor tiene el contexto para decidirlo, el lector no tiene que adivinarlo. -->
### Deviations from Design
{One entry per deviation, each with this exact shape:

- {what departed from the design, and why} — `verify-checks: yes|no`

`verify-checks:` is a **literal token, not prose** — `sdd-verify` matches it to classify the
deviation. Answer `yes` when the deviation touches a spec requirement or a task `criteria:`: the
design artifact is now stale and verify will read it as a mismatch. Answer `no` when it is pure
execution detail that nothing verifiable covers. You are the only one with the context to say
which; the orchestrator is not, and verify will not guess.
**A deviation with no token, or with a hedged one, is treated as `yes`** — the strict default, so
that forgetting to declare is never the cheap way past the gate.
If none: "None — implementation matches design."}

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

### Deviations from Design
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

### Deviations from Design
{As above.}

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

Three sections are **mandatory** in the persisted artifact, and all three are cumulative — each batch
merges its own into what previous batches left:

```markdown
## Apply Progress: {change-name}

### Tasks
{Every task of the change, from every batch, with its state: `[x]` done, `[ ]` pending.
Never drop a task a previous batch completed — that is what makes the next batch skip it.}

### Files Changed
{Every file touched, from every batch, with what was done to it. `sdd-verify` reads THIS copy
to scope its Strict TDD coverage check — it is not only orchestrator context.}

### Deviations from Design
{Every deviation from every batch, each in the exact shape this template declares above,
including its `verify-checks: yes|no` token. `sdd-verify` reads THIS copy — not the return —
to classify each deviation. A deviation that reaches the orchestrator and not the artifact
is a deviation `sdd-verify` will never see.}

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
```

**`### Deviations from Design` goes in BOTH places**, with the same content: in the return, because
the Unresolved Decisions Guard reads it there as Tier 2; in the artifact, because `sdd-verify` reads
it there to apply the `verify-checks` classification. Same reason the TDD evidence lives in both.
Dropping either copy breaks a different consumer, and neither failure is loud.
