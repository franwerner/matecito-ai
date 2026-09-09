# EDR — Moment count counts gates, not bullets

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
`gate-presentation.md`'s "Where this governs" section declared a count of interruption moments in prose (previously nine), but its own enumeration split one moment into two separate rows of its "Seven more moments" table: `A validator's findings` and `risks` fire through the identical mechanism (one `gates:` value, `reported`, per `rule/gate-firing-triggers.md:298`), so they are one moment shown with two anchor sources, not two moments. The table's arithmetic already overcounted by one before this change touched it. This change also retires the spec-mine confirmation gate along with its executor, which would take the declared count from nine to eight without ever correcting the underlying bug — leaving the corrected count and the enumeration permanently one apart.

## Decisión
A moment is counted as a **gate** — an actual mechanism instance a user can be stopped by — never as a bullet in a presentation table. `gate-presentation.md`'s "Seven more moments" table is corrected to six rows: `A validator's findings` and `risks` merge into one row naming both anchor sources. The "Both mining confirmation gates" bullet narrows to name only the decision-mine confirmation gate, since the spec-mine confirmation gate retires with its executor; the ratifying-gates group ("The two gates that ratify a batch of items") becomes a literal two instead of a bullet packing two sub-gates into one. Total: two phase gates (pending-decisions, decision-mine) plus six orchestrator moments (Brief Confirmation Gate, Discovery Gate, Uncommitted-Work Gate, Review Workload Guard, a phase's `blocked` return, and the merged validator-findings/`risks` row) = eight.

## Reglas verificables
- **[manual]** `gate-presentation.md`'s "Where this governs" section states the moment count as eight, and its own enumeration (two phase gates + a six-row table) sums to eight.
- **[manual]** The orchestrator-moments table carries exactly six rows; `A validator's findings` and `risks` are merged into one row naming both anchor sources, never listed as two separate rows.
- **[manual]** Any later addition or removal of a moment recounts by counting gates (mechanism instances), never by counting table rows or bullets — a row that packs one shared trigger behind two anchor sources is always one moment.

## Alternativas consideradas
(a) Leave the presentation table's grouping alone and only drop the spec-mine bullet — rejected: the declared count would still disagree with the table's own enumeration (which would still overcount by one), failing the requirement that the declared count and the enumeration agree. (b) Withdraw the recount entirely and keep declaring nine, on the reading that this file never counted spec-mine separately from decision-mine — rejected: the user ratified eight, and the table's own overcount bug (risks/validator-findings split into two rows) predates and is independent of the spec-mine removal.

## Consecuencias
`flow/ratify-gate-items.md`'s References row (folded into the durable capability-spec store at this change's final dispatch, not by this record) must also read "eight (two phase gates + six orchestrator moments)" to stay consistent with this corrected arithmetic. Any future audit of the moment count reads `gate-presentation.md`'s "Where this governs" section as the single source of truth, and verifies it by counting distinct gate mechanisms, not table rows.
