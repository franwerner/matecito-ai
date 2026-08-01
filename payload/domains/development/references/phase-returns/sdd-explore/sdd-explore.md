<!-- matecito-ai: template canónico del retorno de sdd-explore.
     Existe por la misma razón que el de sdd-design: el formato vivía inline en la skill, sin
     variantes por status. La skill decía "si el pedido es demasiado vago, aclará qué falta" sin
     designar NINGUNA sección donde ponerlo, así que cada ejecutor que se topaba con ese caso
     inventaba la suya y ninguna gate del orquestador la encontraba.
     Este archivo es LA fuente del formato. La skill y el agente lo referencian; no lo copian. -->

# Return template — `sdd-explore`

The exact shape of what this phase hands back to the orchestrator. The orchestrator validates the
return against this file: it matches section titles literally, so a title that differs in wording,
casing or heading level is a section it will not find.

- **Fields of the envelope**: Section D of `~/.claude/skills/_shared/sdd-phase-common.md`. Not repeated here.
- **This file**: the `detailed_report` block — which sections, in which order, with which titles, per status.

## Sections of this phase

| Section | Emitted | Read by |
| --- | --- | --- |
| `### Current State` | always | the orchestrator, as context; the next phase, as its starting point |
| `### Affected Areas` | always | the orchestrator, as context |
| `### Approaches` | always | the orchestrator: it shows them to the user |
| `### Recommendation` | always | the orchestrator: it shows it to the user |
| `### Risks` | always | the orchestrator, surfaced but never blocking |
| `### Blocker` | only on `status: blocked` | the orchestrator: it puts the question to the user |
| `### Ready for Proposal` | always | the orchestrator, to route |

Titles are fixed. `### Ready for Proposal` keeps that title in every lane, including the ones where
the `propose` add-on is off — the **body** names the phase that actually comes next; the title does
not change with the lane.

This phase has **no Tier-1 mailbox** (see the mailbox table, Section D.3): nothing it returns opens
the Unresolved Decisions Guard. It investigates and reports; it fixes nothing.

`### Risks` here is the analysis: what could go wrong with the approaches you compared. It is part of
the block, and it does not replace the envelope's `risks` field (D.4) — the same restriction applies
to both, so a decision the user owns never goes in either.

Two statuses this phase does **not** use. `partial`: an investigation that hit a wall is still
`done`, with the gap stated plainly in `### Current State` and `### Ready for Proposal` — say what
you could not find out, do not dress the gap up as a finding. `needs-input`: D.1 reserves it for
phases whose own skill designates it, and this one does not — a question only the user can answer
makes you return `blocked`.

## `status: done` — the exploration landed

```markdown
## Exploration: {topic}

### Current State
{How the system works today, in the part relevant to this topic. What you actually read, not what
you assume is there. If you could not establish something, say so here rather than filling it in.}

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
{Yes/No — and what the orchestrator should tell the user. On "No", say what is missing.
In a lane where `propose` is off, still title the section this way and name the phase that comes
next instead.}
```

## `status: blocked` — the topic cannot be explored as given

For the case the skill flags as "too vague to explore" or "not enough information to answer": the
question is not yours to settle, and guessing at what the user meant produces an exploration nobody
asked for. Same block, with the differences below. Everything else is emitted as usual — a blocked
exploration still carries whatever you did establish.

```markdown
## Exploration: {topic}

### Current State
{Whatever you DID establish before hitting the wall. If you got nowhere: "Nothing established —
see Blocker."}

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

**What I could not establish**: {what is missing and why — an ambiguity in the topic as given, or
something the codebase does not answer.}

**What unblocks it**: {the answer, or the access, you need.}

### Ready for Proposal
No — blocked. This phase cannot continue until the blocker is resolved.
```

`### Blocker` is where the question goes. Not `risks` — Section D.4 forbids routing through that
field a question the user owns — and not `### Recommendation`, which is for the choice you propose
once you have something to choose between.

## Artifact vs return — the same content here

Unlike the other phases, this one's **artifact** is this same block: the skill persists the identical
content to Engram (`sdd/{change-name}/explore`, or `sdd/explore/{topic-slug}` when standalone). So
whatever you trim from the return you also trim from what the next phase reads. In `none` mode there
is no artifact at all and this block is the entire output.
