<!-- matecito-ai: template canónico del retorno de sdd-design.
     Existe porque el formato de salida vivía disperso en la skill, el agente y el guard, y cada
     edición desincronizaba las copias: dos ejecutores independientes llegaron a inventar secciones
     distintas para devolver un bloqueo, porque ninguna plantilla contemplaba ese caso.
     Este archivo es LA fuente del formato. La skill y el agente lo referencian; no lo copian. -->

# Return template — `sdd-design`

The exact shape of what this phase hands back to the orchestrator. The orchestrator validates the
return against this file: it matches section titles literally, so a title that differs in wording,
casing or heading level is a section it will not find.

- **Fields of the envelope**: Section D of `~/.claude/skills/_shared/sdd-phase-common.md`. Not repeated here.
- **This file**: the `detailed_report` block — which sections, in which order, with which titles, per status.

## Sections of this phase

| Section | Emitted | Read by |
| --- | --- | --- |
| `### Summary` | always | the orchestrator, as context |
| `### Blocker` | only on `status: blocked` | the orchestrator: it puts the question to the user |
| `### New Decisions` | always | Unresolved Decisions Guard — **Tier 1** |
| `### Open Questions` | always | Unresolved Decisions Guard — Tier 2, informative |
| `### Next Step` | always | the orchestrator, to route |

Titles are fixed. `### New Decisions` becomes `### New Decisions (not yet in EDRs)` **only** when the
decision store is active — both forms are valid and the orchestrator accepts either.

**Split into summary/rationale.** Both `### New Decisions` and `### Open Questions` declare it. Each
item carries two parts, `summary` and `rationale`, in the `new_decisions` / `open_questions` JSON:
`summary` is what the gate prints, `rationale` is the full reasoning — always emitted into this
block, never printed by default. Both are non-empty, single-line strings; a missing one, or one with
an embedded newline, fails the render naming the item and the part, and nothing reaches stdout. In
`### New Decisions`, the `· rationale:` line sits directly below the item's `· blocking-test:` line —
same item, same section, no separate channel.

## `status: done` — the design was produced

```markdown
## Design Created

**Change**: {change-name}
**Location**: Engram `sdd/{change-name}/design`

### Summary
- **Approach**: {one line}
- **Key Decisions**: {N} documented
- **Files Affected**: {N} new, {M} modified, {K} deleted
- **Testing Strategy**: {coverage planned}

### New Decisions
{The architectural choices this change requires, one per item, each with the choice, the
alternatives weighed and why this one. Only those that PASS the blocking test in the skill's
`## Rules`: the ones it catches are not filed here, they make you return `blocked`.
Emitted whether or not the decision store exists — only the "(not yet in EDRs)" suffix, the
domain citation and the capture recommendation depend on the store.
Every item carries the `· blocking-test:` token below it — see "The blocking-test token".
If there are genuinely none: "None."}

- {the choice}: {what you chose} — {alternatives weighed, and why this one}
  · blocking-test: none
  · rationale: {one line: the full reasoning — why this choice, restated for the record even though the gate only prints the line above}

### Open Questions
{What does NOT fix a decision: implementation doubts, things to validate during apply.
Anything that fixes a decision belongs above, or makes you return `blocked`.
If there are none: "None."}

- {the open question, as one line}
  · rationale: {one line: why it is open — what makes it not yours to settle here}

### Next Step
Ready for tasks (sdd-tasks).
```

## `status: blocked` — a decision the user owns is in the way

Same block, with three differences. Everything else is emitted as usual — a blocked design still
carries the work you did complete.

```markdown
## Design Created

**Change**: {change-name}
**Location**: not persisted — or the Engram key, if you did persist. State which.

### Summary
- **Approach**: undefined — depends on the blocker below. {What IS settled, if anything.}
- **Key Decisions**: {N} proposed + {M} blocking, unresolved
- **Files Affected**: {N} / {M} / {K}, or "not determinable — see Blocker"
- **Testing Strategy**: {or "not fixed; depends on the resolution"}

### Blocker
{The question, in one line, phrased so the user can answer it without reading the rest.}

**Options weighed**: {each one, with what it costs and what it buys.}

**Why this is not mine to settle**: {name WHICH of the blocking test's three axes the alternatives
differ in — new infrastructure, public contract, data model — and HOW. If the blocker is an
Accepted decision record that did not foresee this case, state both sides instead: what the record
fixes and why, and what it breaks here.}

**What unblocks it**: {the answer you need.}

### New Decisions
{The ones that DO pass the blocking test and hold regardless of how the blocker resolves, same shape
as `done` — each item still needs both `summary` and `· rationale:`.
If every decision is subordinate to the blocker: "None — all pending decisions depend on the blocker."}

### Open Questions
{As above — same shape as `done`, each item with its `· rationale:`.}

### Next Step
None. This phase cannot continue until the blocker is resolved.
```

`### Blocker` is where the question goes. Not `risks` — Section D.4 forbids routing a decision the
user owns through that field — and not `### New Decisions`, which is for the choices the blocking
test lets through.

<!-- matecito-ai: the blocking test used to be pure self-assessment: the executor ran it in its head and
     published only the verdict. In the functional test a decision whose OWN text said "this needs a
     queue and a worker the project does not have today" — axis 1, literally — arrived filed under
     `New Decisions` with `status: done`, and the orchestrator had no mechanical way to notice: catching
     it meant reading the decision's prose and re-running the test itself, which is the interpretation
     these guards exist not to do. Same fix as `verify-checks:` for design deviations one level up:
     the phase DECLARES the axis, and the reader classifies without re-deriving anything. -->
## The blocking-test token

Every item under `### New Decisions` carries one token line directly beneath it:

```
· blocking-test: none | infra | contract | data-model
```

The token is a **claim about the alternatives**, not about your confidence:

- `none` — "I put the alternatives side by side and they differ in NONE of the three axes; that is why
  this item is here and not in `### Blocker`." This is the only value consistent with the item's
  location, and therefore the normal one.
- `infra` / `contract` / `data-model` — the alternatives DO differ in that axis. Writing this under
  `### New Decisions` contradicts the item's own destination: an axis that differs makes the decision
  `blocked`, so the item is in the wrong mailbox.

The orchestrator reads it mechanically, without reopening the decision:

| Token | What it asserts | What the orchestrator does |
| --- | --- | --- |
| `none` | the test ran and came back negative | ordinary Tier 1 — the user confirms |
| an axis named | the item is in the wrong mailbox | stop and surface, as it would for a `blocked` |
| absent, or hedged | the test did not run, or the answer is being withheld | Tier 1 under the strict reading — same default as an undeclared deviation |

The token is what makes the skill's own requirement satisfiable: it demands the test be *auditable
from outside* — a reader who sees only the alternatives reaching the same verdict without
reconstructing your reasoning — and until now there was nowhere for that verdict to be seen. Do not
hedge it and do not omit it: one line per decision, and the strict reading is what an omission buys.

On `status: blocked`, the `### Blocker` section already names the axis in prose ("Why this is not mine
to settle"), so the blocking item needs no token. The items that remain under `### New Decisions` in a
blocked return still carry theirs.

## Artifact vs return — do not confuse them

The design **artifact** (persisted to Engram, format in the skill) is the full document: technical
approach, data flow, file changes, interfaces, testing, migration. The **return** is this block: what
the orchestrator needs in order to route and to gate. The orchestrator never reads the artifact.

A section that exists only in the artifact is invisible to every gate. That is why `### New Decisions`
appears in both: `##` in the artifact, `###` here.
