# SDD Phase — Common Protocol

Boilerplate identical across all SDD phase skills. Sub-agents MUST load this alongside their phase-specific SKILL.md.

Executor boundary: every SDD phase agent is an EXECUTOR, not an orchestrator. Do the phase work yourself. Do NOT launch sub-agents, do NOT call `delegate`/`task`, and do NOT bounce work back unless the phase skill explicitly says to stop and report a blocker.

<!-- matecito-ai: the two mining executors load this file like any phase, and inherited three things they
     are forbidden or unable to honor: Section C ORDERS them to persist their artifact, when the kernel's
     invariant is that a mining executor NEVER writes (the gate and the materialization happen in the main
     thread, under explicit user confirmation); Section D.2 assigns them a `detailed_report` block that
     does not exist, since the table enumerates the ten flow phases and no `*-mine` is among them, so a
     literal reader goes looking for a template in `~/.claude/references/phase-returns/` that was never
     written; and the D.1 status enum describes flow states (`partial`, `needs-decision`, routing via
     `next_recommended`) that a scope → `candidates[]` executor has no use for. The scoping note lives
     here rather than as a patch in each agent: this is the file all ten executors read, and a note here
     cannot drift out of sync with copies that do not exist. -->
**Scope — mining executors load Sections A and B only.** If you are `development-decisions-mine` or
`development-spec-mine`, Sections **C** and **D** do NOT apply to you. You write nothing: you return
`candidates[]` and the main thread decides, at a user-confirmed gate, what gets materialized — so
there is no artifact for Section C to persist, and no `detailed_report` block for Section D to shape.
The form of your return lives in your own skill. Section A (skill loading, including the domain
fragment) and Section B (artifact retrieval) bind you exactly as they bind a phase, and so does the
no-self-invented-defaults rule below.

No self-invented defaults (absolute): if you hit a genuine decision or an open question that your inputs (brief / spec / design / tasks / confirmed scope + this phase's skill) do NOT resolve, do NOT pick a "most likely" default to keep going. Return `status: blocked` with the exact question so the orchestrator can put it to the user (or `needs-input`, when your phase skill designates that status for this situation — see Section D). A missing or unanswered question is NOT permission and NOT a default. This holds even in Automatic mode. (Executor-side of the kernel's "Open question = blocked, not permission" rule.)

## A. Skill Loading

<!-- matecito-ai: el mecanismo de inyección (skill-registry / skill-resolver / Project Standards) fue removido de este ecosistema. Las fases cargan su propia skill directamente. -->

1. Load this phase's own `SKILL.md` (and any module it explicitly tells you to load, e.g. `strict-tdd.md`).
2. If the orchestrator's launch prompt includes explicit `SKILL: Load` instructions, load those exact skill files too.
3. Project-level conventions live in the project's own files — read `.matecito-ai/edr/` (architecture decisions), `CLAUDE.md`, and any `config.yaml` the phase references. Those are the project's standards in matecito-ai; there is no separate registry to consult.
<!-- matecito-ai: el fragmento del dominio se carga on-demand en el main thread, así que un ejecutor de contexto fresco NO lo tiene. Sin este punto, reglas como "Contract & definition shapes" nunca llegan a quien escribe el código. -->
4. **Read the domain fragment: `~/.claude/matecito-ai/domains/development.md`.** It binds this flow's vocabulary AND carries rules the orchestrator cannot enforce on your behalf — notably **"Contract & definition shapes — never inferred"**, which governs every entity, DB model/migration/schema, DTO, public type/interface/enum, event payload, and config schema you touch. Those rules bind you exactly as this file does. Not optional, and a summary does not count.

Proceed with the phase skill as your authority. Loading a skill file is NOT delegation.

## B. Artifact Retrieval (Engram Mode)

**CRITICAL**: `mem_search` returns 300-char PREVIEWS, not full content. You MUST call `mem_get_observation(id)` for EVERY artifact. **Skipping this produces wrong output.**

**Run all searches in parallel** — do NOT search sequentially.

```
mem_search(query: "sdd/{change-name}/{artifact-type}", project: "{project}") → save ID
```

Then **run all retrievals in parallel**:

```
mem_get_observation(id: {saved_id}) → full content (REQUIRED)
```

Do NOT use search previews as source material.

## C. Artifact Persistence

Every phase that produces an artifact MUST persist it. Skipping this BREAKS the pipeline — downstream phases will not find your output.

### Engram mode

```
mem_save(
  title: "sdd/{change-name}/{artifact-type}",
  topic_key: "sdd/{change-name}/{artifact-type}",
  type: "architecture",
  project: "{project}",
  capture_prompt: false,
  content: "{your full artifact markdown}"
)
```

`topic_key` enables upserts — saving again updates, not duplicates.
`capture_prompt: false` is mandatory for SDD artifacts because they are automated pipeline outputs, not human/proactive memory saves. Set it when the Engram tool schema supports it; if an older schema rejects or does not expose the field, omit it rather than failing.

<!-- matecito-ai: acá vivían los modos `OpenSpec` e `Hybrid`, removidos del ecosistema hace tiempo.
     Sobrevivían en el archivo que leen las DIEZ fases, y eran el único lugar que le decía a un
     ejecutor que podía escribir artefactos de flujo a disco — contra el kernel ("matecito-ai is
     engram-only"), contra `persistence-contract.md` y contra el `engram | none` que declara cada
     skill. Los modos son dos: `engram` y `none`. -->
### None mode

Return result inline only. Do not write any files or call `mem_save`.

## D. Return Envelope — CANONICAL CONTRACT

<!-- matecito-ai: FUENTE ÚNICA. Este contrato estaba duplicado en este archivo y en los nueve
     Result Contracts de agents/, y cada edición desalineaba las otras copias (tres rondas
     seguidas de defectos por esta causa). Ahora se define UNA vez acá y los agentes REFERENCIAN.
     Un agente NO re-declara `status`, `detailed_report` ni sus buzones: los hereda de esta sección
     y sólo agrega lo específico de su fase. Si algo de acá tiene que cambiar, cambia acá. -->

This section is the **single source of truth** for what every phase returns. An agent file does NOT
redeclare these fields — it points here and adds only what is specific to its phase. Never encode
this contract by a literal phrase to match ("when your skill defines a … block"); the phases are
enumerated below on purpose.

### D.1 `status`

`done` · `partial` · `blocked` · `needs-input` · `needs-decision`

- `done` — the phase finished its work. (Historical note: some older text says `success`; `done` is the value, they are the same thing.)
- `partial` — real work landed but the phase is not finished (e.g. an apply batch with tasks left).
- `blocked` — you cannot continue without a resolution that is not yours to make. Carry the exact question and the options you weighed.
- `needs-input` — you need an answer only the user can give, and the flow resumes by **re-dispatching this same phase** with it. NOT a failure. Used by `sdd-intake` Pass 1; any phase may use it if its own skill says so.
- `needs-decision` — an architectural decision must be captured before the flow proceeds (`sdd-intake` early guard, when decision records are active). The orchestrator routes to the decision-capture skill.

### D.2 `detailed_report` — MANDATORY, by phase

<!-- matecito-ai: el formato literal de cada retorno vive en un template propio, no acá ni disperso en
     las skills. La tabla sólo dice CUÁL bloque te toca; el template dice cómo se ve, incluidas las
     variantes por status y qué secciones son condicionales. -->
<!-- matecito-ai: the block used to be typed by hand against a markdown template, and the failures were
     always the same three: a dropped section, a re-levelled heading, a forgotten sentinel — each one a
     gate that never fires, silently. The shape now lives as data in `<phase>.yaml` and a script builds
     the block from your content, so those three are not mistakes you can make. What you still have to
     get right is the CONTENT, which is what the `.md` is for. -->
**You do not write the block by hand. You render it:**

1. `node ~/.claude/scripts/render-return.js --phase <your-phase> --schema` — the exact data shape.
   Read it instead of reconstructing it: it also states which fields are **derived** (never supply
   those — the renderer computes them) and that an empty list is `[]`, which is a different thing from
   an omitted field.
2. Write your content as JSON to a temporary file.
3. `node ~/.claude/scripts/render-return.js --phase <your-phase> --data <file>` — its stdout IS your
   `detailed_report`. Paste it verbatim.

A non-zero exit names the offending field: fix the data and re-run. **Never hand-write the block to
get past a failure** — the failure is the contract telling you the content is wrong, not the tool
being in the way. If the script cannot run at all, say so in your return rather than substituting a
hand-made block that nothing checked.

**The summary/rationale split.** A section whose contract declares `items.rationale` splits each of
its items into two required parts: `summary` (what gates print) and `rationale` (the full reasoning
— always emitted into this block, never printed by default). Both are non-empty, single-line strings
— supply them per item as `{ summary: "...", rationale: "..." }` in your JSON. A missing one, or one
with an embedded newline, fails the render naming the item and the part, and nothing reaches stdout —
the same contract as any other missing field. This does not change the principle above: you still
never hand-write the block, and `--schema` tells you, section by section, whether it wants this
two-field shape or the plain one.

**`~/.claude/references/phase-returns/<your-phase>/<your-phase>.md` is what you read to know WHAT
belongs in each section** — the rationale, the examples, and the reason each rule exists. It is the
authority on meaning; the `.yaml` beside it is the authority on shape. The orchestrator validates what
you return against that same contract.

Sections whose body is only a `None…` sentinel are emitted by the renderer when their list is empty —
you get that for free, and it is why "nothing to report" and "the phase dropped the section" stay
distinguishable.

| Phase | Block to carry |
| --- | --- |
| `sdd-intake` | `## Intake Brief: {title}` (Pass 2) · `## Discovery Form: {title}` (Pass 1) |
| `sdd-explore` | `## Exploration: {topic}` |
| `sdd-propose` | `## Proposal Created` |
| `sdd-spec` | `## Specs Created` |
| `sdd-design` | `## Design Created` |
| `sdd-tasks` | `## Tasks Created` |
| `sdd-apply` | `## Implementation Progress` |
| `sdd-verify` | `## Verification Report` |
| `sdd-archive` | `## Change Archived` |
| `sdd-init` | `## SDD Initialized` |

### D.3 Mailboxes the orchestrator's guards read

These live **inside `detailed_report`**, never in `risks`. Tier definitions and gate behavior are in
the domain fragment (`~/.claude/matecito-ai/domains/development.md`, `## Guards`); this table only
fixes which section belongs to which phase.

| Phase | Section | Tier | Emitted |
| --- | --- | --- | --- |
| `sdd-propose` | `### Scope and approach (unconfirmed)` | 1 | always |
| `sdd-spec` | `### Derived capabilities (unconfirmed)` | 1 | always |
| `sdd-spec` | `### New Decisions` — a proposal's ratification gate for a lane with no `design` add-on active (see `~/.claude/references/decision-capture/in-flow-capture.md`) | 1 | conditional — only when the lane running has no `design` add-on |
| `sdd-design` | `### New Decisions` — or `### New Decisions (not yet in EDRs)` when the decision store is active; **both titles are valid and the orchestrator accepts either** | 1 | always |
| `sdd-design` | `### Open Questions` | 2 | always |
| `sdd-tasks` | `### Tasks not traceable to spec/design` | 1 | always |
| `sdd-apply` | `### Unmandated Forks` | 1 | always |
| `sdd-apply` | `### Mandated Departures` | 2 | always |
| `sdd-verify` | `## Decision Gaps` | — | only when the change materialized at least one decision record (`### Decisions Materialized` in `apply-progress` carries ≥1 row) — no flag, see `in-flow-capture.md` |
| `sdd-verify` | `## UI Verdict` | — | only when the UI check applies |

**Tier 1** stops the flow and asks the user; **Tier 2** is surfaced but does not block. `sdd-apply`'s
two rows are not a Tier-1/Tier-2 pair over the same content: `### Unmandated Forks` carries a fork
the confirmed artifacts do not fix and that `sdd-apply` did NOT apply (its item's `mandate:` token is
always `chosen`); `### Mandated Departures` carries what it DID apply, either because an artifact
already fixed the point (`mandate: covered`) or because no alternative was valid and the constraint
is named (`mandate: forced`). A missing or hedged `mandate:` is read as `chosen`, so an absorbed
deviation nobody can back with a named constraint routes to the Tier-1 section by default, never the
cheap way past the gate. `sdd-spec`'s `### New Decisions` row is the SAME mailbox concept as
`sdd-design`'s — a decision-proposal ratification gate — surfacing conditionally, one lane earlier;
it is not a third kind of thing. `sdd-verify`'s `## Decision Gaps` is not a tier section: for
`development` it feeds nothing kernel-side (the domain declares its own in-flow decision-capture
mechanism, materialized during `sdd-apply` — see `~/.claude/references/decision-capture/in-flow-capture.md`);
a domain with no such mechanism (e.g. `design`) may still feed a post-verify mine gate under its own
flag. Gate behavior lives in the domain fragment (`~/.claude/matecito-ai/domains/development.md`, `##
Guards`) — it reads this table and keeps no parallel copy of it.

**Eight of these mailboxes split each item into `summary`/`rationale`**: `sdd-propose`'s `Scope and
approach`, `sdd-spec`'s `Derived capabilities` and its conditional `New Decisions`, both `sdd-design`
rows (`New Decisions` and `Open Questions`), `sdd-tasks`'s `Tasks not traceable`, and both `sdd-apply`
rows (`Unmandated Forks` and `Mandated Departures`) — their `.yaml` contract declares
`items.rationale`. Emission stays total:
both parts always land in `detailed_report`. Printing only the `summary` at the gate is that
section's own **declared** presentation, fixed by its contract — never an assistant judging what
counts as brief. The `rationale` is reproduced verbatim, from the block already in context, the
moment it is asked for.

**Register.** Write an item's `summary` in plain language: state what was decided and what follows
from it, so a reader who does not know this flow's vocabulary understands the outcome. The
ecosystem's internal vocabulary — phase names, token names, section titles, the tier/mailbox nouns —
belongs in `rationale`, not in `summary`. This binds the item's prose only; a section's own title
(which may legitimately carry that vocabulary, e.g. `### Derived capabilities (unconfirmed)`) is
untouched. Where a check keys off literal words in the envelope's `Summary` — `validate-return.js`
check 4 matches a section's `summary_claims` pattern together with a non-zero digit — keep those
exact words and state the count as a digit whenever the summary claims that section has content: the
check is a literal match, and dropping the words would turn it off silently, not on purpose.

Everything marked **always** is emitted even when there is nothing to report, with a `None…`
sentinel — that is what lets the orchestrator tell "nothing to raise" from "the phase dropped the
section". A section that some other rule of your own skill declares **conditional**
(`sdd-intake`'s `### Early guard (EDRs)` when the decision store is inactive, the UI Verdict when
the UI gate is off, the TDD evidence table outside Strict TDD) is legitimately absent when its
condition does not hold: that is not a broken return.
<!-- matecito-ai: the first example used to be `## EDR Conflicts`, which lives only in the design ARTIFACT
     and never in a return. Illustrating a rule ABOUT RETURNS with a section that is not a return
     section is the same artifact/return confusion this section exists to close. -->

The `None…` sentinel is canonically written in English. The orchestrator's guard also recognizes its
Spanish forms (`Ninguna`, `Ninguno`, `Nada`) because this ecosystem converses in Spanish — but you
write `None`: relying on the tolerance is how a sentinel ends up in a form nobody matches.

### D.4 The other fields

<!-- matecito-ai: nobody cross-checked the `Summary`, so a return could announce "2 decisions documented"
     over a body carrying the empty sentinel and dispatch silently. The Return Contract Check now puts
     it against the body (check 4). -->
- `executive_summary`: 1-3 sentence summary of what was done, in the same plain register as a mailbox item's `summary` — see Section D.3's "Register" above, not restated here. **It is checked against your own `detailed_report`**: a section you closed with the empty sentinel cannot be summarized as having content, and a section carrying rows cannot be summarized as empty. The orchestrator compares the two claims mechanically and stops on a mismatch — so write the summary from the block you actually emitted, not from the work you intended to emit
- `artifacts`: artifact keys/paths written — `none` only when the phase wrote nothing at all: `needs-input`, or a stop that happened **before any work landed**. Stopping partway does NOT excuse you from persisting: if work of yours completed, it is persisted first and listed here, whatever the status (see `sdd-apply`'s mid-batch rule)
- `next_recommended`: the next SDD phase, or `none` — always legal, and the only correct value on `blocked` / `needs-input`
- `risks`: risks and assumptions to validate. **Never** a decision the user owns, an ambiguity you resolved by assuming, or content that belongs to a mailbox in D.3
- `skill_resolution`: `phase-skill` (loaded this phase's own SKILL.md), `fallback-path` (loaded via an explicit SKILL: Load path), or `none`

Example:

```markdown
**Status**: done
**Summary**: Proposal created for `{change-name}`. Defined scope, approach, and rollback plan.
**Artifacts**: Engram `sdd/{change-name}/proposal`
**Next**: sdd-spec or sdd-design
**Risks**: None
**Skill Resolution**: phase-skill — loaded sdd-{phase}/SKILL.md
(other values: `fallback-path`, or `none`)
```

## E. Review Workload Guard

SDD must protect reviewer cognitive load, not only generate tasks.

- The default PR review budget is **400 changed lines** (`additions + deletions`).
- The orchestrator MUST cache a delivery strategy at session start: `ask-on-risk` (default), `auto-chain`, `single-pr`, or `exception-ok`.
- The orchestrator MUST pass `delivery_strategy` to `sdd-tasks` and the resolved decision to `sdd-apply`.
- `sdd-tasks` MUST forecast whether the planned work may exceed that budget.
- The forecast MUST include exact plain-text guard lines: `Decision needed before apply: Yes|No`, `Chained PRs recommended: Yes|No`, and `400-line budget risk: Low|Medium|High`.
- If the forecast is high, `sdd-tasks` MUST recommend chained or stacked PRs using deliverable work units.
- `sdd-apply` MUST NOT start oversized work unless the delivery strategy resolves to chained/stacked PR slices or explicitly accepted `size:exception`.
- Each chained PR slice must have a clear start, clear finish, autonomous scope, verification, and reasonable rollback.
- In a Feature Branch Chain, PR #1 targets the feature/tracker branch and later child PRs target the immediate previous PR branch; if GitHub shows previous slices in a child diff, retarget/rebase until the diff is clean.

This guard exists to reduce reviewer burnout and keep implementation delivery safe. Do not treat it as optional process noise.
