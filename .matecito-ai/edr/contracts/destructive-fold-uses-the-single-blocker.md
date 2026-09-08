# EDR — A destructive delta-spec fold stops through sdd-apply's single Blocker section, as a fourth cause, never a section of its own

- **Status:** Accepted
- **Date:** 2026-09-04

## Contexto
sdd-apply's Step 5b (the delta-spec fold, relocated from sdd-archive) can find a merge that would drop scenarios or sections the delta never mentions. sdd-apply already declares one `### Blocker` section that carries every reason the phase stops, and the return template's blocked-cause list already has three entries: an unpinned contract, a `design-conflict: conflicts` verdict, and a fork or failure that stops the rest of the batch. A destructive fold needed a place to report from.

## Decisión
No new return section. The destructive-fold stop adds one more cause — a fourth — to the existing blocked-cause list in `sdd-apply.md`'s "Which status", and is reported through the existing `### Blocker` section like the other three: same lead/options/what-unblocks shape, no dedicated conditional section. `contracts/verdict-in-its-own-conditional-section.md` fixes the criterion for when a new section IS warranted — a verdict that must be reported even when it does NOT block. The destructive-fold verdict has no such non-blocking half: it either stops the phase (the fold does not run) or produces nothing to report at all (the fold ran cleanly and folded normally, reported under `### Capability-Specs Materialized` instead). It also carries no token any guard classifies. Because cause 4 can only arise after every task of the batch is already complete, `### Remaining Tasks` on this blocked variant reads `None — every task completed; what remains is the fold.` and `### Status` gains a non-task stop line: `{N}/{N} tasks complete. Progress persisted: yes. Stopped at the durable-spec fold — see Blocker.`

## Reglas verificables
- **[manual]** A destructive delta-spec fold is reported ONLY through `### Blocker`, as a fourth listed cause in `sdd-apply.md`'s blocked-cause ladder — never a section of its own.
- **[manual]** The `blocked` variant of `sdd-apply`'s return admits `### Remaining Tasks: None — every task completed; what remains is the fold.` and a non-task `### Status` line naming the fold stop — both documented explicitly in the template.
- **[manual]** The destructive-fold cause names the capability-spec and exactly which scenarios or sections the fold would drop, and states the same two options every other Blocker cause states: apply as-is, or a corrected input (here: a corrected delta spec).

## Alternativas consideradas
A dedicated `### Destructive Fold` conditional section, mirroring `### Rejected Proposals Checked` / `### Contract Shapes Proposed` — rejected. Those two sections exist because their verdict must be reported even on a non-blocking outcome (a rejection checked clean, a contract shape ratified), which gives them content on statuses other than `blocked`. The destructive-fold verdict has no such half: either it stops the phase, in which case `### Blocker` already carries everything a reader needs (the lead, the options, what unblocks it), or it never happens, in which case there is nothing to report beyond the ordinary `### Capability-Specs Materialized` success row. A section with no non-blocking half would just restate the blocker in a second place, which "One blocker, one place" already forbids for every other cause.

## Consecuencias
The blocked-cause list grows to four items, ordered so cause 4 is listed last: it is the only one that can arise after every task in the batch is already complete, which is also why `### Remaining Tasks` and `### Status` needed new admitted forms on the `blocked` variant — a case that used to assume "everything left, including the task you stopped on" now has to admit "nothing left, the stop is the fold itself".
