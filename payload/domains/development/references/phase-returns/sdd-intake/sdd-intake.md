<!-- matecito-ai: template canónico del retorno de sdd-intake.
     Existe por la misma razón que el de sdd-design: el formato vivía inline en la skill. Ahora es un
     passthrough de una sola pasada: recibe el pedido crudo, decide los cuatro flags, corre el early
     guard cuando aplica, y devuelve el brief — sin discovery, sin triage, sin size. El discovery de
     dos pasadas se mudó a `sdd-explore`, que sí puede leer el código antes de preguntar.
     Este archivo es LA fuente del formato. La skill y el agente lo referencian; no lo copian. -->

# Return template — `sdd-intake`

The exact shape of what this phase hands back to the orchestrator. The orchestrator validates the
return against this file: it matches section titles literally, so a title that differs in wording,
casing or heading level is a section it will not find.

- **Fields of the envelope**: Section D of `~/.claude/skills/_shared/sdd-phase-common.md`. Not repeated here.
- **This file**: the `detailed_report` block — which sections, in which order, with which titles, per status.

## One pass, one block

This phase returns once, with `## Intake Brief: {short title}`. There is no discovery form here — it
reads the raw request, decides the four downstream flags, runs the early guard when the EDR store is
active, and emits the brief. Nothing about this phase asks the user anything.

`needs-decision` and `blocked` come from the early guard and therefore **exist only when the EDR
store is active** (per the activation gate). With the store absent or empty, the guard is skipped
silently and neither status is reachable for EDR reasons.

## Sections of `## Intake Brief`

| Section | Emitted | Read by |
| --- | --- | --- |
| `### Request (structured)` | always | the orchestrator, and the next phase as its starting point |
| `### Classification` | always | the orchestrator, and each flag's own reader |
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
value — it is metadata, decided by this phase and reported once when the brief returns.

## `status: done` — the brief

```markdown
## Intake Brief: {short title}

### Request (structured)
{1-2 sentences: what the user wants, restated clearly}

### Classification
- Type: {feature|bug|refactor|chore}
- Domains touched: {list of canonical EDR domains}
<!-- matecito-ai: estos flags los DECIDE intake y los LEEN fases/consumidores posteriores (`sdd-design`
     lee `diagram` del brief; `sdd-verify` lee `ui-test`; el orquestador lee `worktree-isolation`). Vivían sólo
     en el envelope, que es efímero: al persistir el brief se perdían, y el lector downstream no los
     encontraba. Van en el brief. -->
- Diagram: {needed|not-needed} — {one line why}
- UI test: {needed|not-needed} — {one line why}
- Worktree isolation: {active|inactive} — {one line why}
{If `repo.components` is declared: "- Components: {name[, name...] | unassigned}" — omit this bullet
entirely when the project config declares no `repo.components`.}

### Early guard (EDRs)
Clear — no conflict with existing EDRs, no undecided question.

### Next
sdd-explore
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

### Early guard (EDRs)
🟡 NEEDS DECISION: `<domain>` has no EDR for {what}. Capture via development-decisions-bootstrap first.

**What is undecided**: {the question, in one line, phrased so the user can answer it without reading the rest.}

**Why it cannot wait for design**: {what the later phases would have to assume in order to proceed.}

### Next
development-decisions-bootstrap — then re-run from `sdd-explore`.
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

### Early guard (EDRs)
⛔ BLOCKED: conflicts with `<domain>/<slug>.md` — {what}. Resolve before proceeding.

**What the EDR fixes, and why**: {its side.}

**What the request asks for**: {the other side.}

**Options**: {adjust the request | update the EDR via development-decisions-bootstrap (update mode)}
— {what each one costs.}

### Next
None. This phase cannot proceed until the conflict below is resolved.
```

`### Early guard (EDRs)` is where the blocker goes — it is this phase's blocker mailbox, so there is
no separate `### Blocker` here. Not `risks` either: Section D.4 forbids routing through that field a
decision the user owns.

## Artifact vs return — the same content

The brief is persisted to Engram (`sdd/{change-name}/intake`) with this same content, so what you
trim from the return you also trim from what the next phase reads.

The brief always goes straight to the next phase — there is no scope-confirmation gate. The
orchestrator reports the flags in one notice line and dispatches; nothing waits on it.
