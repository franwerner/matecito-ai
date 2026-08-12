<!-- matecito-ai: template canónico del retorno de sdd-tasks.
     Existe por lo mismo que el de sdd-design: el formato vivía inline en la skill, incompleto —
     sólo contemplaba el caso feliz — y nada decía qué hacer con el retorno cuando la fase se
     bloquea. Dos copias del mismo formato (skill + guard) se desincronizan en la primera edición.
     Este archivo es LA fuente del formato. La skill y el agente lo referencian; no lo copian. -->

# Return template — `sdd-tasks`

The exact shape of what this phase hands back to the orchestrator. The orchestrator validates the
return against this file: it matches section titles literally, so a title that differs in wording,
casing or heading level is a section it will not find.

- **Fields of the envelope**: Section D of `~/.claude/skills/_shared/sdd-phase-common.md`. Not repeated here.
- **This file**: the `detailed_report` block — which sections, in which order, with which titles, per status.

## Sections of this phase

| Section | Emitted | Read by |
| --- | --- | --- |
| `### Breakdown` | always | the orchestrator, as context |
| `### Implementation Order` | always | the orchestrator, as context |
| `### Blocker` | only on `status: blocked` | the orchestrator: it puts the question to the user |
| `### Tasks not traceable to spec/design` | always | Unresolved Decisions Guard — **Tier 1** |
| `### Review Workload Forecast` | always | Review Workload Guard — matched **line by line** |
| `### Next Step` | always | the orchestrator, to route |

Titles are fixed and this phase declares no variants of them.

**Split into summary/rationale.** `### Tasks not traceable to spec/design` declares it. Each item
carries two parts, `summary` and `rationale`, in the `untraceable_tasks` JSON: `summary` is what the
gate prints, `rationale` is the full reasoning — always emitted into this block, never printed by
default. Both are non-empty, single-line strings; a missing one, or one with an embedded newline,
fails the render naming the item and the part, and nothing reaches stdout. `summary`'s register is
fixed once in Section D.3 of `sdd-phase-common.md` — not restated here.

## `status: done` — the breakdown was produced

```markdown
## Tasks Created

**Change**: {change-name}
**Location**: Engram `sdd/{change-name}/tasks`

### Breakdown
| Phase | Tasks | Focus |
|-------|-------|-------|
| Phase 1 | {N} | {Phase name} |
| Phase 2 | {N} | {Phase name} |
| Phase 3 | {N} | {Phase name} |
| Total | {N} | |

### Implementation Order
{The recommended order and why — what has to exist before what. One short paragraph.}

### Tasks not traceable to spec/design
{Every task that links to NO spec requirement and to nothing the design establishes. This is work
the user did not ask for: it is where silent scope creep enters — it lands in the checklist, gets
implemented and gets verified without anyone having agreed to it. Do NOT drop those tasks and do
NOT fold them in as if they came from the spec; list them and let the user decide. One item per task:

- {the task, as one line}
  · rationale: {one line: what motivated it — a gap you found, an implied prerequisite, a project convention}

If every task traces: "None — every task links to spec or design."}

### Review Workload Forecast
- Estimated changed lines: {estimate or range}
- 400-line budget risk: {Low | Medium | High}
- Chained PRs recommended: {Yes | No}
- Delivery strategy: {ask-on-risk | auto-chain | single-pr | exception-ok}
- Decision needed before apply: {Yes | No}
- Suggested work-unit PR split: {brief list, or "Not needed"}

<!-- matecito-ai: decía "ask the user which chain strategy" — pero esta fase es headless y el Review
     Workload Guard del orquestador ya hace esa pregunta entre esta fase y apply. Acá se SEÑALA que
     hace falta la decisión; no se pide. -->
### Next Step
{Ready for implementation (sdd-apply) — or, when `Decision needed before apply: Yes`, say so plainly:
the chain strategy is still `pending` and the orchestrator's Review Workload Guard resolves it with
the user before apply. You are flagging that the decision is due, not asking for it.}
```

## `status: blocked` — a decision the user owns is in the way

Same block, with three differences. Everything else is emitted as usual — a blocked breakdown still
carries the work you did complete.

```markdown
## Tasks Created

**Change**: {change-name}
**Location**: not persisted — or the Engram key, if you did persist. State which.

### Breakdown
{The phases you did manage to break down, same table. If you got nowhere:
"Not produced — see Blocker."}

### Implementation Order
{What IS settled, if anything. Otherwise: "undefined — depends on the blocker below."}

### Blocker
{The question, in one line, phrased so the user can answer it without reading the rest.}

**Options weighed**: {each one, with what it costs and what it buys.}

**Why this is not mine to settle**: {what makes this a decision and not execution detail — a
contract whose shape is not pinned, a spec requirement the design does not cover, two orderings
that are not equivalent. If an upstream artifact pins it but is inconsistent or does not cover
this case, state both sides: what it fixes and why, and what it breaks here.}

**What unblocks it**: {the answer you need.}

### Tasks not traceable to spec/design
{As above, same shape, over the tasks you did write — each item still needs both `summary` and
`· rationale:`. If you wrote none: "None — no tasks were written."}

### Review Workload Forecast
- Estimated changed lines: {your best estimate over the scope as briefed, or "not estimable — the breakdown is incomplete"}
- 400-line budget risk: {Low | Medium | High}
- Chained PRs recommended: {Yes | No}
- Delivery strategy: {the one you received}
- Decision needed before apply: Yes
- Suggested work-unit PR split: {brief list, or "Not needed"}

### Next Step
None. This phase cannot continue until the blocker is resolved.
```

`### Blocker` is where the question goes. Not `risks` — Section D.4 forbids routing a decision the
user owns through that field — and not `### Tasks not traceable to spec/design`, which is for tasks
you DID write and whose origin the user has to sanction.

The three guard lines below still carry a legal value on `blocked`: the section is unconditional and
the Review Workload Guard matches it literally, so it cannot be replaced by prose. `Decision needed
before apply` is `Yes` because the blocker IS that decision.

## The Review Workload Forecast lines are a contract

<!-- matecito-ai: estas líneas no son prosa — el Review Workload Guard las matchea literales antes de
     dispatchear sdd-apply. Reescribir una etiqueta apaga el guard en silencio. -->

Three of the lines in that section are matched **literally**, by label, by the orchestrator's Review
Workload Guard. Reproduce the labels verbatim — do not reword, re-case or turn them into prose:

```text
400-line budget risk: {Low | Medium | High}
Chained PRs recommended: {Yes | No}
Decision needed before apply: {Yes | No}
```

Their values are the ones enumerated and nothing else. Any of `400-line budget risk: High`,
`Chained PRs recommended: Yes` or `Decision needed before apply: Yes` stops the flow before
`sdd-apply` so the user can choose a chain strategy — which is exactly what a reworded label
silently prevents.

## Artifact vs return — do not confuse them

The tasks **artifact** (persisted to Engram, format in the skill) is the checklist itself: the
`- [ ]` lines with their `criteria:` sub-line and, when they apply, their `· edr:` refs, plus the
forecast table and the suggested work units. The **return** is this block: what the orchestrator
needs in order to gate and to route. The orchestrator never reads the artifact.

A section that exists only in the artifact is invisible to every gate. The forecast therefore
appears in both — its plain-text lines in the artifact for `sdd-apply` to enforce (Step 2a), and the
same lines here for the Review Workload Guard, which reads only this block.
