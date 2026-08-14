<!-- matecito-ai: template canónico del retorno de sdd-propose.
     Existe porque el formato vivía inline en la skill y declaraba una sola forma —la del camino
     feliz—, sin decir cómo se devuelve un `blocked`: la fase que ELIGE el approach era justamente
     la que no tenía dónde poner la pregunta que no le toca contestar, y el ejecutor terminaba
     eligiendo para poder emitir el bloque que sí conocía.
     Este archivo es LA fuente del formato. La skill y el agente lo referencian; no lo copian. -->

# Return template — `sdd-propose`

The exact shape of what this phase hands back to the orchestrator. The orchestrator validates the
return against this file: it matches section titles literally, so a title that differs in wording,
casing or heading level is a section it will not find.

- **Fields of the envelope**: Section D of `~/.claude/skills/_shared/sdd-phase-common.md`. Not repeated here.
- **This file**: the `detailed_report` block — which sections, in which order, with which titles, per status.

## Sections of this phase

| Section | Emitted | Read by |
| --- | --- | --- |
| `### Summary` | always | the orchestrator, as context |
| `### Scope and approach (unconfirmed)` | always | Unresolved Decisions Guard — **Tier 1** |
| `### Contract Shapes Proposed` | conditional — only when `has_contract_proposals` is true and status is `blocked` | Unresolved Decisions Guard — **Tier 1** |
| `### Blocker` | only on `status: blocked` | the orchestrator: it puts the question to the user |
| `### Next Step` | always | the orchestrator, to route |

Titles are fixed. This phase declares no accepted variants of them.

<!-- matecito-ai: esta fase NO tenía buzón, y es donde se fijan el approach y el mapeo de
     capabilities que `sdd-spec` consume como contrato. En Interactive lo salvaba el checkpoint entre
     fases; en Automatic la propuesta pasaba derecho, y un mapeo equivocado recién asomaba una fase
     después, cuando spec derivaba el suyo. Propose es un add-on opcional: si lo encendiste, es
     justamente porque querías confirmar el alcance antes de especificar. -->
**`### Scope and approach (unconfirmed)` is this phase's Tier-1 mailbox**, and it is emitted always.
It carries the two things this phase fixes on the user's behalf: the **approach** chosen (with the
alternatives you weighed, if any) and the **capability mapping** — which capabilities this change
adds or modifies, since `sdd-spec` consumes that mapping as a contract and turns it into delta specs.
Neither has been agreed to yet: the guard raises the batch and the user confirms or corrects before
spec runs. It is never a `None` case — a proposal always fixes an approach and always touches some
capability; if you genuinely cannot name them, you do not have a proposal, you have a `blocked`.

**Split into summary/rationale.** Each item declares two parts, `summary` and `rationale`, in the
`scope_unconfirmed` JSON: `summary` is what the gate prints, `rationale` is the full reasoning behind
it — always emitted into this block, never printed by default. Both are non-empty, single-line
strings; a missing one, or one with an embedded newline, fails the render naming the item and the
part, and nothing reaches stdout. `summary`'s register is fixed once in Section D.3 of
`sdd-phase-common.md` — not restated here. `summary` also carries a **250-character cap**, enforced by
`render-return.js`; an over-cap summary fails the render before anything reaches stdout.

Every item also carries the `· anchor:` token, declared first so it prints directly under the
summary. What counts as a legitimate anchor — `<repo-path>[:line]` or `<engram-key>`, and what a
not-yet-written target anchors to — is fixed once in Section D.3 of `sdd-phase-common.md`; not
restated here.

**`### Contract Shapes Proposed`** is the dedicated home for an unpinned contract or definition — a
domain entity, a DB model, a DTO, a public type — the shape "Contract & definition shapes — never
inferred" (`~/.claude/matecito-ai/domains/development.md`) forbids guessing. It never appears as free
prose inside `### Blocker`. Emitted only when this phase declares `has_contract_proposals: true` **and**
`status: blocked` — both conditions must hold, and the gate is supplied by the phase, never derived
from the list being non-empty: an empty list under a `true` gate is legitimate (every proposal, once
checked, needed no further ratification) and reads differently from `false` ("nothing to propose").
Absence is never a violation on its own; only the section appearing on a status outside `[blocked]` is.

Each item is ONE compound entry, never split across items and never one item per field: a one-line
`summary` (the contract, what needs it, why it is unpinned), one `· anchor:`, then one
`· field: {name} — {type} — {description}` continuation line per field the contract needs — every
field proposed, none dropped, none summarized — then `· rationale:`. `summary` keeps the ordinary
250-character cap; each field's own `description` carries its own 160-character cap — the **field
count is never capped**. One anchor per contract, never per field: it is also the identity the forward
uses when the ratified shape comes back.

The item carries no field for where the shape will be stored or persisted — that omission is
deliberate, not missing: this section proposes the shape, it does not decide its home. Once ratified
(or adjusted) at the gate, the field list travels back only in this phase's re-dispatch prompt,
identified by the item's own `anchor` and `summary`, and is written into this phase's own artifact —
this phase MUST NOT read it from any stored artifact.

Only two statuses have a shape here: `done` and `blocked`. This phase's skill does not designate
`needs-input`, and a proposal is written whole or not at all, so `partial` does not arise.

## `status: done` — the proposal was written

```markdown
## Proposal Created

**Change**: {change-name}
**Location**: Engram `sdd/{change-name}/proposal` (engram) | inline (none)

### Summary
- **Intent**: {one line: the problem being solved}
- **Scope**: {N deliverables in, M items deferred}
- **Approach**: {one line: the technical approach chosen}
- **Risk Level**: {Low/Medium/High}

### Scope and approach (unconfirmed)
- **Approach**: {the approach this proposal fixes, and what you weighed against it and discarded. If the alternatives differ in new infrastructure, in the public contract or in the data model, this is not yours to fix — return `blocked` instead.}
  · anchor: {the concrete source this item is about — a `<repo-path>[:line]` or `<engram-key>`; not-yet-written work anchors to what surfaced it (e.g. the intake brief's Engram key)}
  · rationale: {one line: the full reasoning behind this approach and why it beat the alternatives}

- **Capability mapping**: {which capabilities this change touches, as `New` or `Modified`, exactly as the artifact's `## Capabilities` section states them. `sdd-spec` consumes this as its contract: a New becomes a full spec, a Modified becomes a delta against the durable capability-spec of that name. A wrong name here silently writes the delta against the wrong capability.}
  · anchor: {e.g. the durable capability-spec path for a Modified capability, or the intake brief's Engram key for a New one}
  · rationale: {one line: what in the artifact's Affected Areas or Request grounds this mapping}

- **In scope / deferred**: {what this proposal deliberately leaves out, if the boundary is not obvious.}
  · anchor: {what surfaced this boundary — the intake brief's Engram key, or a repo path if code already draws the line}
  · rationale: {one line: why this boundary and not another}

### Next Step
Ready for specs (sdd-spec) or design (sdd-design).
```

The Capabilities section of the artifact is the contract with `sdd-spec` — it does not travel in
this return, and it must be filled in (`None` under both sub-sections when nothing changes at spec
level, never left as a placeholder). A proposal that ships placeholder capabilities is what makes
`sdd-spec` derive its own mapping downstream and raise it as unconfirmed.

## `status: blocked` — something that fixes scope or approach is not yours to settle

Same block, with three differences, plus the conditional `### Contract Shapes Proposed` when the
blocker is an unpinned contract. Everything else is emitted as usual — a blocked proposal still
carries what you did settle.

```markdown
## Proposal Created

**Change**: {change-name}
**Location**: not persisted — or the Engram key, if you did persist. State which.

### Summary
- **Intent**: {what IS settled — the intent usually survives the blocker}
- **Scope**: {what is in/out so far, or "undefined — depends on the blocker below"}
- **Approach**: undefined — depends on the blocker below. {What IS settled, if anything.}
- **Risk Level**: {or "not assessable until the blocker is resolved"}

{CONDITIONAL — only when `has_contract_proposals` is true; omit the whole section otherwise, do not
even print "None." This is the shape a stop over an unspecified contract takes — the reason this
proposal is `blocked`:}

### Contract Shapes Proposed
- {the contract, what needs it, why it is unpinned — one line, ≤250}
  · anchor: {the concrete source this contract is about — a `<repo-path>[:line]` or `<engram-key>`}
  · field: {name} — {type} — {the field's description; every field proposed, none dropped}
  · field: {name} — {type} — {a second field, same shape; as many lines as the contract needs}
  · rationale: {one line: what a wrong guess here propagates to}

### Blocker
{The question, in one line, phrased so the user can answer it without reading the rest.}

**Options weighed**: {each one, with what it costs and what it buys. If they imply different
Capabilities entries, say which — that is the contract `sdd-spec` will read.}

**Why this is not mine to settle**: {what the alternatives change that is the user's call — the
problem being solved, what is in and out of scope, an approach whose alternatives differ in
infrastructure the project does not have, in the public contract, or in the data model, or a
contract shape covered by "Contract & definition shapes — never inferred". Name it concretely; "I
was not sure" is not a reason to block.}

**What unblocks it**: {the answer you need.}

### Scope and approach (unconfirmed)
{Emitted here too — it is unconditional. Carry whatever IS settled: the parts of the approach the
blocker does not touch, and the capability mapping if you could name it. Each item still needs both
`summary` and `· rationale:`, same as `done`. When the blocker is precisely the approach or the
mapping, say so and point at `### Blocker` instead of repeating it: "Depends on the blocker above."}

### Next Step
None. This phase cannot continue until the blocker is resolved.
```

`### Blocker` is where the question goes. Not `risks` — Section D.4 forbids routing a decision the
user owns through that field, and the artifact's own `## Risks` table is for risks you accepted and
mitigated, not for a choice you did not make.

## Artifact vs return — do not confuse them

The proposal **artifact** (persisted to Engram, format in the skill) is the full document: intent,
scope, capabilities, approach, affected areas, risks, rollback, dependencies, success criteria. The
**return** is this block: what the orchestrator needs in order to route and to gate. The
orchestrator never reads the artifact.

That asymmetry is sharper here than in any other phase: this is where scope and approach get fixed,
and none of it is in the return beyond one line each. In Interactive mode the between-phase
checkpoint is what puts the artifact in front of the user; in Automatic mode nothing does. That is
the contract as it stands — respect it, and raise anything you cannot fit into it as a blocker
rather than smuggling it into a section of your own.
