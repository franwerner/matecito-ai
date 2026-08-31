# EDR — Whether the vocabulary sweep reaches the record stores

- **Status:** Accepted
- **Date:** 2026-08-31

## Contexto
narrow-gating-triggers retires the Tier 1/Tier 2 gating nouns across the payload. Two Accepted EDRs (`structure/verdict-classified-by-the-orchestrator.md:14`, `contracts/verdict-in-its-own-conditional-section.md:14`) and one Accepted capability-spec scenario (`rule/mailbox-item-summary-rationale-split.md:171`, a GIVEN reading "un item informativo en una sección declarante (Tier 2, como Open Questions)") keep the retired nouns in incidental positions. The spec's Scope names three capability-spec files under `.matecito-ai/` and zero EDRs — none of these three files was in the ratified 23-file Scope, and the final sweep (task 3.4) needs to know whether they count as residue or as a gap.

## Decisión
Leave all three record-store files exactly as they are. In each, the rule's substance survives intact and only an incidental noun goes stale — the rule they state is still correct and still enforced; only the WORD used to describe one example inside it ("Tier 2", "Tier 1/2") is now stale vocabulary. They are named explicitly in task 3.4's final sweep as known, accepted residue, not as a defect the sweep should flag.

## Reglas verificables
- **[manual]** `.matecito-ai/edr/structure/verdict-classified-by-the-orchestrator.md:14` keeps its `Tier 1/2` wording; the rule it states (a `conflicts` verdict on a non-`blocked` return is a contract violation, not a mailbox-guard routing) is unaffected by the wording.
- **[manual]** `.matecito-ai/edr/contracts/verdict-in-its-own-conditional-section.md:14` keeps its `Tier 1 o Tier 2` wording; the rule it states (no orchestrator guard classifies `### Rejected Proposals Checked` into either tier — its only reader is the `conflicts → status: blocked` content check) is unaffected.
- **[manual]** `.matecito-ai/development-specs/rule/mailbox-item-summary-rationale-split.md:171` keeps its `(Tier 2, como Open Questions)` GIVEN clause; the scenario's business rule (an informative item in a declaring section) is unaffected.
- **[manual]** The final sweep (task 3.4) treats these three files as allow-listed residue, not as violations of the retired-vocabulary criterion.

## Alternativas consideradas
Extend this change's scope now to touch the three files. Rejected: the two EDRs are reasoning documents whose incidental noun does not change what they mandate, so editing them here would be scope creep with no functional benefit. The capability-spec scenario is worse — `sdd-archive` merges capability-spec deltas scenario-anchored and never accepts a hand edit mid-flow, so fixing its wording would require re-running `sdd-spec` to author a formal delta, which this change's ratified scope does not include.

## Consecuencias
After this change, three files in the record stores describe the retired classification by name, in positions where only the label is stale and the rule survives. A future reader of `mailbox-item-summary-rationale-split.md`'s GIVEN clause sees "Tier 2" without the corresponding `gates: muted` vocabulary existing anywhere else — the scenario still reads correctly (Open Questions is still the example it names), but the noun no longer resolves to an active classification. No mechanical check enforces this file staying in sync with the retired vocabulary; a future `sdd-spec` delta on `rule/gate-firing-triggers` or `rule/mailbox-item-summary-rationale-split` is the natural point to reword it, not this change.
