# EDR — The delta-spec fold fires on the final apply dispatch, keyed off the tasks artifact, never the return

- **Status:** Accepted
- **Date:** 2026-09-04

## Contexto
The change that moves the delta-spec fold from sdd-archive into sdd-apply's Step 5b (between Step 5, marking tasks [x], and Step 6, persisting apply-progress) needed a firing rule: WHEN inside a batch that may span several dispatches does the fold run. sdd-apply already resolves `status` and builds `### Remaining Tasks` in Step 7, downstream of both Step 5 and the new Step 5b — so a rule that reads either of those would depend on a value that does not exist yet at the point the fold has to decide.

## Decisión
Step 5b evaluates one mechanical predicate immediately after Step 5: does the tasks artifact, as just updated, contain no `- [ ]`? True means this is the dispatch that closes the change's implementation, and the fold runs, covering the change's complete delta across every capability it touches, in this one run. False means tasks remain, and Step 5b is skipped entirely — no section, no mention. The predicate is read from the tasks artifact alone, never from `status` or from `### Remaining Tasks`. A re-dispatch that follows a destructive-fold stop re-arms the predicate for free: apply-progress already shows every task `[x]`, so the re-dispatched batch implements nothing, reaches Step 5 with nothing to mark, and Step 5b evaluates the same predicate again against the corrected delta.

## Reglas verificables
- **[manual]** Step 5b's firing predicate reads the tasks artifact's own `- [ ]` marks, never `status` and never `### Remaining Tasks` — both are built in Step 7, downstream of Step 5b.
- **[manual]** Step 5b runs in Consolidation/Serial Mode only, immediately after Step 5 and before Step 6; an isolated run never reaches it, because it already skips Steps 5 and 6 under the single-writer rule.
- **[manual]** A continuation batch that leaves any task unchecked does not fire the fold; the batch that leaves none unchecked fires it once, covering the whole change's delta.
- **[manual]** A re-dispatch that follows a destructive-fold stop re-arms the predicate without any special-casing: apply-progress already marks every task done, so the re-dispatched batch reaches Step 5b again with nothing left to implement.

## Alternativas consideradas
(a) Key off `status: done` — rejected: status is resolved in Step 7, downstream of the fold, so the step would depend on a value that does not exist yet. (b) Key off `### Remaining Tasks` being empty — rejected for the same reason: it is a return section built in Step 7. (c) Fold inside Step 6 as a sub-step of persistence — rejected: a destructive-fold stop must still run Step 6 ("Stopping Mid-Batch" makes Steps 5 and 6 mandatory on every exit path), and a fold nested inside Step 6 would have to re-enter its own caller. (d) Per-task fold mirroring Step 4b (decision materialization) — rejected at the propose gate, and independently unworkable: nothing tags a task with the capability it implements, so there is no way to know which capability-spec a given task's code belongs to until the whole change's delta is read as a unit.

## Consecuencias
Placement between Steps 5 and 6 makes a fold that stops an ordinary "Stopping Mid-Batch" exit: Steps 5 and 6 already run on every exit path, so a destructive-fold stop does not need a new persistence carve-out. The fold happens exactly once per change, on whichever dispatch turns out to be the last one — which may not be knowable in advance, so every Consolidation/Serial dispatch pays the cost of checking the predicate, even the ones that turn out not to fire it. An isolated run's exemption falls out of the single-writer rule for free, with no separate case to write.
