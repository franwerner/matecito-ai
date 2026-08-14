<!-- matecito-ai: template canónico del retorno de sdd-intake.
     Existe por la misma razón que el de sdd-design: el formato vivía inline en la skill, y acá el
     costo era mayor porque esta fase tiene DOS pasadas con bloques distintos y cuatro status
     posibles. La plantilla inline cubría bien la pasada feliz y dejaba a interpretación el resto:
     dos ejecutores independientes inventaron secciones distintas para el mismo caso.
     Ojo con el caso que ya costó caro: la lista de preguntas VACÍA en Pass 1 es un retorno legítimo
     y completo, no un atajo al brief. Está contemplado abajo, explícito.
     Este archivo es LA fuente del formato. La skill y el agente lo referencian; no lo copian. -->

# Return template — `sdd-intake`

The exact shape of what this phase hands back to the orchestrator. The orchestrator validates the
return against this file: it matches section titles literally, so a title that differs in wording,
casing or heading level is a section it will not find.

- **Fields of the envelope**: Section D of `~/.claude/skills/_shared/sdd-phase-common.md`. Not repeated here.
- **This file**: the `detailed_report` block — which sections, in which order, with which titles, per status.

## Two passes, two different blocks

This is the only phase that carries a different block depending on the pass (D.2 assigns it both):

| Pass | Condition | Status | Block |
| --- | --- | --- | --- |
| Pass 1 | your launch prompt carries NO answers | `needs-input`, always | `## Discovery Form: {short title}` |
| Pass 2 | your launch prompt carries the user's answers | `done` · `needs-decision` · `blocked` | `## Intake Brief: {short title}` |

**Pass 1 always returns `needs-input`.** There is no path from Pass 1 to the brief — not even when
you find nothing to ask. See the empty-list case below; it is the failure this phase has actually
had, not a hypothetical.

`needs-decision` and `blocked` come from the early guard and therefore **exist only when the EDR
store is active** (per the activation gate). With the store absent or empty, the guard is skipped
silently and neither status is reachable for EDR reasons.

## Pass 1 — sections of `## Discovery Form`

| Section | Emitted | Read by |
| --- | --- | --- |
| `### Request (as received)` | always | the orchestrator, to confirm it dispatched what the user said |
| `### Questions (unanswered — for the orchestrator to ask)` | always, **may be empty** | the orchestrator: it puts them to the user verbatim |
| `### Next` | always | the orchestrator, to route back into this same phase |

Nothing else. No classification, no triage, no early guard, no brief, and no `mem_save` — Pass 1
persists nothing, so the envelope's `artifacts` is `none`.

## `status: needs-input` — Pass 1, questions formulated

```markdown
## Discovery Form: {short title}

**Status**: needs-input

### Request (as received)
{the raw request, verbatim}

### Questions (unanswered — for the orchestrator to ask)
1. {question} — {why it matters: what changes downstream depending on the answer}
2. {question} — {why it matters}
3. {question} — {why it matters}

### Next
Re-dispatch `sdd-intake` with these answers (or with the confirmation) to produce the brief.
```

## `status: needs-input` — Pass 1, nothing to ask

Same block. The question list is **empty**, and that is a complete, legitimate return — not a
degenerate one, and not licence to produce the brief. You do not get to decide on the user's behalf
that they had nothing to add: the empty list travels with your one-line reading of the request, and
the user confirms or corrects it. A one-word confirmation is cheap; a brief built on an ambiguity
you failed to notice is not.

```markdown
## Discovery Form: {short title}

**Status**: needs-input

### Request (as received)
{the raw request, verbatim}

### Questions (unanswered — for the orchestrator to ask)
No questions — I read the request as: {one line}. Confirm or correct before I build the brief.

### Next
Re-dispatch `sdd-intake` with the confirmation to produce the brief.
```

Never pad the list with invented questions to make the return look fuller, and never answer a
question yourself. Returning `needs-input` is the successful outcome of Pass 1.

## Pass 2 — sections of `## Intake Brief`

| Section | Emitted | Read by |
| --- | --- | --- |
| `### Request (structured)` | always | the orchestrator, and the first phase of the lane as its starting point |
| `### Classification` | always | the orchestrator, as context |
| `### Discovery answers` | always | the orchestrator: the user checks their own answers came back intact |
| `### Triage` | always | the orchestrator: the lane the user confirms or adjusts at the INTAKE GATE |
| `### Early guard (EDRs)` | only when the EDR store is active | the orchestrator: it routes `blocked` / `needs-decision` from here |
| `### Next` | always | the orchestrator, to route |

`### Early guard (EDRs)` is the one **conditional** section: when `.matecito-ai/edr/` is absent or
empty the guard is skipped silently and the section is absent with it — no mention of EDRs anywhere
in the brief. That absence is a legitimate return, not a dropped section. When the store IS active
the section is emitted, including in the all-clear case.

The `Components` bullet **inside** `### Classification` has the same kind of gate, one level deeper:
it is a **conditional bullet**, present only when `repo.components` is declared in the project config
(`components_axis_active`). Resolved false, the bullet is legitimately absent from `### Classification`
— no mention of components anywhere in the brief. Resolved true, it is required like any other bullet:
supply it or the render fails naming the field. Unlike `Diagram` and `UI test`, no phase reads this
value — it is metadata for the person confirming the gate.

## `status: done` — Pass 2, the brief

```markdown
## Intake Brief: {short title}

### Request (structured)
{1-2 sentences: what the user wants, restated clearly after the discovery form}

### Classification
- Type: {feature|bug|refactor|chore}
- Domains touched: {list of canonical EDR domains}
- Size: {trivial|small|medium|large}
<!-- matecito-ai: estos flags los DECIDE intake y los LEEN fases/consumidores posteriores (`sdd-design`
     lee `diagram` del brief; `sdd-verify` lee `ui-test`; el orquestador lee `worktree-isolation`). Vivían sólo
     en el envelope, que es efímero: al persistir el brief se perdían, y el lector downstream no los
     encontraba. Van en el brief. -->
- Diagram: {needed|not-needed} — {one line why}
- UI test: {needed|not-needed} — {one line why}
- Worktree isolation: {active|inactive} — {one line why}
{If `repo.components` is declared: "- Components: {name[, name...] | unassigned}" — omit this bullet
entirely when the project config declares no `repo.components`.}

### Discovery answers
- {question}: {answer verbatim — never paraphrased into something more convenient}
- ...

{If Pass 1 returned an empty question list: "No questions were asked. The user confirmed the reading:
{the one line}."}

### Triage
Lane: {direct | reduced | full | custom} — add-ons: [{explore? propose? design? tasks?} or none] — {one line why}

### Early guard (EDRs)
Clear — no conflict with existing EDRs, no undecided question.

### Next
{direct-implementation | the first phase the chosen lane runs — `sdd-explore` if `explore` is on,
else `sdd-propose` if `propose` is on, else `sdd-spec`}
```

## `status: needs-decision` — an architectural question no EDR covers

Same block. The guard found a decision the flow would otherwise take by itself, so the flow does not
start until it is captured. Everything before the guard is emitted as usual — the work you completed
still travels.

```markdown
## Intake Brief: {short title}

### Request (structured)
{as above}

### Classification
{as above}

### Discovery answers
{as above}

### Triage
Lane: {the lane you would recommend} — add-ons: [...] — {one line why}. Held: the decision below is
captured first.

### Early guard (EDRs)
🟡 NEEDS DECISION: `<domain>` has no EDR for {what}. Capture via development-decisions-bootstrap first.

**What is undecided**: {the question, in one line, phrased so the user can answer it without reading the rest.}

**Why it cannot wait for design**: {what the later phases would have to assume in order to proceed.}

### Next
development-decisions-bootstrap — then re-run the lane from `### Triage`.
```

## `status: blocked` — the request contradicts an Accepted EDR

Same block. Do NOT recommend the flow: the conflict is resolved first, either by adjusting the
request or by updating the EDR. Naming which EDR, and both sides of the conflict, is the whole point
— the orchestrator carries a conversation it cannot have without them.

```markdown
## Intake Brief: {short title}

### Request (structured)
{as above}

### Classification
{as above}

### Discovery answers
{as above}

### Triage
None — blocked. No lane is recommended until the conflict below is resolved.

### Early guard (EDRs)
⛔ BLOCKED: conflicts with `<domain>/<slug>.md` — {what}. Resolve before proceeding.

**What the EDR fixes, and why**: {its side.}

**What the request asks for**: {the other side.}

**Options**: {adjust the request | update the EDR via development-decisions-bootstrap (update mode)}
— {what each one costs.}

### Next
None. This phase cannot recommend a lane until the conflict is resolved.
```

`### Early guard (EDRs)` is where the blocker goes — it is this phase's blocker mailbox, so there is
no separate `### Blocker` here. Not `risks` either: Section D.4 forbids routing through that field a
decision the user owns. And an ambiguity in the request is neither — it is a discovery question, on
Pass 1.

## Artifact vs return — same content, Pass 2 only

On Pass 2 the brief is persisted to Engram (`sdd/{change-name}/intake`) with this same content, so
what you trim from the return you also trim from what the next phase reads. On Pass 1 nothing is
persisted at all: there is no brief yet, only a form.

The brief always passes through the orchestrator's **INTAKE GATE** — shown to the user, confirmed /
adjusted / cancelled — before any next phase runs, including for `trivial` changes and in Automatic
mode. Never write the return as if the flow proceeds on its own.
