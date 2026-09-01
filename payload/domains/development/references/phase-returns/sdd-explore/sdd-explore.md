<!-- matecito-ai: template canónico del retorno de sdd-explore.
     Existe por la misma razón que el de sdd-design: el formato vivía inline en la skill, sin
     variantes por status. La skill decía "si el pedido es demasiado vago, aclará qué falta" sin
     designar NINGUNA sección donde ponerlo, así que cada ejecutor que se topaba con ese caso
     inventaba la suya y ninguna gate del orquestador la encontraba.
     Now two passes, two different blocks — the discovery cycle moved here from `sdd-intake`: this
     phase reads code FIRST, then formulates its questions, then re-runs with the answers.
     Este archivo es LA fuente del formato. La skill y el agente lo referencian; no lo copian. -->

# Return template — `sdd-explore`

The exact shape of what this phase hands back to the orchestrator. The orchestrator validates the
return against this file: it matches section titles literally, so a title that differs in wording,
casing or heading level is a section it will not find.

- **Fields of the envelope**: Section D of `~/.claude/skills/_shared/sdd-phase-common.md`. Not repeated here.
- **This file**: the `detailed_report` block — which sections, in which order, with which titles, per status.

## Two passes, two different blocks

This phase carries a different block depending on what it found and on the pass:

| Pass | Condition | Status | Block |
| --- | --- | --- | --- |
| Pass 1 | your launch prompt carries NO answers, and you found real questions | `needs-input` | `## Discovery Form: {topic}` |
| Pass 1 | your launch prompt carries NO answers, and you found no questions | `done` | `## Exploration: {topic}` |
| Pass 2 | your launch prompt carries the user's answers to real Pass-1 questions | `done` · `blocked` | `## Exploration: {topic}` |

**Pass 1 returns `needs-input` only when it found real questions.** You run headless — you have no
channel to the user — so a question you cannot resolve from the code or the request itself does not
get invented, it gets returned for the orchestrator to ask. Finding nothing to ask is NOT a Pass-1
return: it proceeds straight into the exploration block, in the same pass, with status `done` — the
request's own reading was already put to the user once, at the Brief Confirmation Gate, before this
phase ran.

## Pass 1 — sections of `## Discovery Form`

| Section | Emitted | Read by |
| --- | --- | --- |
| `### Request (as received)` | always | the orchestrator, to confirm it dispatched what the user said |
| `### Questions (unanswered — for the orchestrator to ask)` | always, **never empty** | the orchestrator: it puts them to the user verbatim |
| `### Next` | always | the orchestrator, to route back into this same phase |

Nothing else. No current state, no approaches, no recommendation, no exploration artifact, and no
`mem_save` — Pass 1 persists nothing, so the envelope's `artifacts` is `none`.

Each question carries an `anchor`, under the ordinary anchor criterion — no exception exists for this
gate. A question grounded in something you read anchors to the repo path it came from (with a start
line when the source is a specific place); a question about the request's own intent anchors to the
intake brief's artifact key.

## `status: needs-input` — Pass 1, questions formulated

```markdown
## Discovery Form: {topic}

**Status**: needs-input

### Request (as received)
{the raw request, verbatim}

### Questions (unanswered — for the orchestrator to ask)
1. {question} — {why it matters: what changes downstream depending on the answer}
   · anchor: {repo-path[:line] the question is grounded in, or the intake brief's artifact key for a
     question about the request's own intent}
2. {question} — {why it matters}
   · anchor: {...}

### Next
Re-dispatch `sdd-explore` with these answers (or with the confirmation) to produce the exploration.
```

Never pad the list with invented questions to make the return look fuller, and never answer a question
yourself. Returning `needs-input` is the successful outcome of a Pass 1 that found real questions —
with none, Pass 1 does not return this block at all; it proceeds to `## Exploration` below, in the same
pass, with status `done`.

## Pass 2 — sections of `## Exploration`

| Section | Emitted | Read by |
| --- | --- | --- |
| `### Current State` | always | the orchestrator, as context; the next phase, as its starting point |
| `### Discovery answers` | always | the orchestrator: the user checks their own answers came back intact |
| `### Affected Areas` | always | the orchestrator, as context |
| `### Approaches` | always | the orchestrator: it shows them to the user |
| `### Recommendation` | always | the orchestrator: it shows it to the user |
| `### Risks` | always | the orchestrator, surfaced but never blocking |
| `### Blocker` | only on `status: blocked` | the orchestrator: it puts the question to the user |
| `### Ready for Proposal` | always | the orchestrator, to route |

Titles are fixed. `### Ready for Proposal` keeps that title always — the **body** names the phase that
actually comes next; the title does not change.

This phase has **no gating mailbox** (see the mailbox table, Section D.3): nothing it returns opens
the Unresolved Decisions Guard. It investigates and reports; it fixes nothing.

`### Risks` here is the analysis: what could go wrong with the approaches you compared. It is part of
the block, and it does not replace the envelope's `risks` field (D.4) — the same restriction applies
to both, so a decision the user owns never goes in either.

Two statuses this phase does **not** use. `partial`: an investigation that hit a wall is still
`done`, with the gap stated plainly in `### Current State` and `### Ready for Proposal` — say what
you could not find out, do not dress the gap up as a finding. `blocked` narrows: it is reserved for a
stop a re-dispatch **with answers** cannot clear — missing access, or a contradiction between the
ratified inputs you were given. A question about the request itself — ambiguity in what was asked, the
"too vague to explore" case — is a discovery-form question, and you return `needs-input` for it instead.

## `status: done` — the exploration landed

```markdown
## Exploration: {topic}

### Current State
{How the system works today, in the part relevant to this topic. What you actually read, not what
you assume is there. If you could not establish something, say so here rather than filling it in.}

### Discovery answers
- {question}: {answer verbatim — never paraphrased into something more convenient}
- ...

{If you reached this block directly from Pass 1, with no questions formulated, leave this list empty —
it renders the "None." sentinel. Nothing here claims a confirmation this phase did not run.}

### Affected Areas
- `path/to/file.ext` — {why it's affected}
- `path/to/other.ext` — {why it's affected}

{If nothing in the codebase is touched — a greenfield topic — write "None — nothing exists yet for
this topic." Do not drop the section.}

### Approaches
1. **{Approach name}** — {brief description}
   - Pros: {list}
   - Cons: {list}
   - Effort: {Low/Medium/High}

2. **{Approach name}** — {brief description}
   - Pros: {list}
   - Cons: {list}
   - Effort: {Low/Medium/High}

{If the investigation turned up only one viable approach, emit it as a single entry and say why the
alternatives were discarded. One entry is a legitimate answer; a dropped section is not.}

### Recommendation
{Which approach and why. Not a decision anyone has agreed to — the user has not seen this yet.}

### Risks
- {Risk 1}
- {Risk 2}

{If there are none: "None."}

### Ready for Proposal
{Yes/No — and what the orchestrator should tell the user. On "No", say what is missing.}
```

## `status: blocked` — a re-dispatch with answers cannot resolve this

For a stop that answers cannot clear: missing access, or a contradiction between the ratified inputs
you were given. Same block, with the differences below. Everything else is emitted as usual — a
blocked exploration still carries whatever you did establish.

```markdown
## Exploration: {topic}

### Current State
{Whatever you DID establish before hitting the wall. If you got nowhere: "Nothing established —
see Blocker."}

### Discovery answers
{As above — whatever answers you already had when you hit the blocker.}

### Affected Areas
{The ones you could identify, or "Not determinable — see Blocker."}

### Approaches
{The ones that survive regardless of how the blocker resolves.
If every approach depends on it: "None — the alternatives depend on the blocker."}

### Recommendation
None — depends on the blocker below.

### Risks
{As above.}

### Blocker
{The question, in one line, phrased so the user can answer it without reading the rest.}

**What I could not establish**: {what is missing and why — missing access, or a contradiction between
the ratified inputs, never an ambiguity in the request itself — that is a discovery-form question.}

**What unblocks it**: {the answer, or the access, you need.}

### Ready for Proposal
No — blocked. This phase cannot continue until the blocker is resolved.
```

`### Blocker` is where the question goes. Not `risks` — Section D.4 forbids routing through that
field a question the user owns — and not `### Recommendation`, which is for the choice you propose
once you have something to choose between.

## Artifact vs return — the same content here

Unlike most other phases, this one's `## Exploration` **artifact** is this same block: the skill
persists the identical content to Engram (`sdd/{change-name}/explore`, or `sdd/explore/{topic-slug}`
when standalone) — whether it was reached via Pass 2 (real questions, now answered) or directly from
Pass 1 (no questions found). So whatever you trim from that return you also trim from what the next
phase reads. In `none` mode there is no artifact at all and this block is the entire output. A Pass 1
that stops with real questions and returns `## Discovery Form` persists nothing — there is no
exploration yet, only a form.
