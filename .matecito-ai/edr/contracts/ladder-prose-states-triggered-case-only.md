# EDR — The gates-ladder prose states the triggered case only, no untriggered-item comparison

- **Status:** Accepted
- **Date:** 2026-08-31

## Contexto
Section D.3 of `_shared/sdd-phase-common.md` introduces the `gates:` value ladder with a two-clause paragraph: the first names what an item declaring no trigger gets, the second names what an item declaring one gets. The second clause used to claim a triggered item 'gates in every value except `reported`, where it only ever surfaces' — which is wrong for `muted`: the table beneath it (and the Accepted capability-spec `rule/gate-firing-triggers.md`, its ladder table and Requisito) both say a `muted` section never gates at all, triggered or not. A proposal to fix this drafted a replacement clause ending '...same as an untriggered item', which is itself false for `muted` (triggered surfaces, untriggered reaches the user nowhere) and fails the proposal's own Success Criterion that the paragraph agree, clause by clause, with all four table rows.

## Decisión
The paragraph's opening clause (what an item declaring no trigger gets) stays byte-identical. The second clause is rewritten to state only the triggered case, with no comparison to the untriggered one: 'an item that does declare one gates only in `always` and `contested`, and merely surfaces in `reported` and `muted`.' This is narrower than the proposal's drafted text, which this record departs from because that draft contradicts the same table it is meant to agree with.

## Reglas verificables
- **[manual]** The gates-ladder paragraph in `_shared/sdd-phase-common.md` (D.3) states only what an item declaring a trigger does — gates only in `always` and `contested`, merely surfaces in `reported` and `muted` — and makes no equivalence claim to the untriggered case.
- **[manual]** The paragraph's opening clause, which introduces the table's untriggered-item column, stays byte-identical across edits; only the trailing clause about a triggered item may change.
- **[manual]** Any future edit to this paragraph is re-checked clause by clause against all four rows of the table at D.3 and against `rule/gate-firing-triggers.md`'s ladder table and Requisito before it ships.

## Alternativas consideradas
Adopting the original proposal's drafted tail verbatim ('...same as an untriggered item') was rejected: it is false for the `muted` row, where a triggered item surfaces but an untriggered one reaches the user nowhere — it would have replaced one contradiction against the table with a subtler one. A comparison clause scoped to `muted` alone (true, but restating what the paragraph three lines below already says at length) was also rejected on the file's own single-statement register rule.

## Consecuencias
The paragraph and the table it introduces can no longer be read as disagreeing for the `muted` case; a reader classifying a `muted`-plus-trigger item gets the same answer from both. The paragraph now makes a narrower claim than before (no equivalence to the untriggered case), so a future edit that reintroduces a comparison clause must be re-verified against all four table rows rather than assumed correct by symmetry.
