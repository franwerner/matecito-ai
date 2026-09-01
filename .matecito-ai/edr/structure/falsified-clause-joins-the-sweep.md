# EDR — Three further clauses this change falsifies join the sweep, beyond the spec's declared-closed Scope

- **Status:** Accepted
- **Date:** 2026-09-01

## Contexto
This change's spec declares a closed `## Scope` of every site it touches. While implementing it, three more clauses turned out to be falsified by the new Brief Confirmation Gate but were not named in that Scope: `sdd-intake/SKILL.md`'s Step 5 prose ("there is no confirmation gate"), its closing rule ("no scope-confirmation gate exists in this flow"), and `sdd-explore.yaml`'s header/inline comments calling the phase two-pass unconditionally.

## Decisión
All three join the sweep. The two `sdd-intake/SKILL.md` clauses are corrected to state that the brief is confirmed at the Brief Confirmation Gate before the next phase is dispatched. The `sdd-explore.yaml` comments are corrected to describe the phase's now-conditional two-pass cycle (a Pass 1 with no questions proceeds straight to `done`). Only comments and prose move — `sdd-explore.yaml`'s schema (`statuses`, `variants`, `sections`) is untouched; `render-return.js --phase sdd-explore --schema` is byte-identical before and after.

## Reglas verificables
- **[manual]** `payload/domains/development/skills/matecito-ai/sdd-intake/SKILL.md` does not state, anywhere, that no confirmation gate exists for the brief — both its Step 5 prose and its closing Rules-section clause name the Brief Confirmation Gate instead.
- **[manual]** `payload/domains/development/references/phase-returns/sdd-explore/sdd-explore.yaml`'s comments describe the phase's two-pass cycle as conditional on finding real questions, not as unconditional — with no change to its `statuses`, `variants`, or `sections` keys.

## Alternativas consideradas
Leaving all three as accepted residue for a later change to find. Rejected: the first two make a payload file the flow reads at runtime disagree with the kernel about whether a confirmation gate exists — a direct contradiction, not stale flavor text — and the third is the identical shape `structure/gating-vocabulary-sweep-scope.md` already ruled gets swept: a comment restating a retired rule, with no behavior change once corrected.

## Consecuencias
The sweep for this change now covers three sites beyond the spec's original closed Scope, all corrected in the same commit as the rest of the falsified-clause work (Phase 6 of `sdd-tasks`'s checklist) so nothing mid-change states the two contradictory claims side by side. `sdd-explore.yaml`'s schema is confirmed unchanged by running its `--schema` output before and after and diffing byte-for-byte.

## Relacionados
- `relacionado-con` → [gating-vocabulary-sweep-scope.md](gating-vocabulary-sweep-scope.md) — precedente directo: un comentario que restablece una regla retirada se barre en el mismo cambio, sin cambio de comportamiento.
- `relacionado-con` → [retired-vocabulary-in-record-stores.md](retired-vocabulary-in-record-stores.md) — precedente: si el sweep final de un cambio de vocabulario debe extender su scope a un sitio fuera del Scope original.
- `relacionado-con` → [lane-vocabulary-sweep-in-record-stores.md](lane-vocabulary-sweep-in-record-stores.md) — precedente: mismo tipo de decisión sobre el ancho del sweep, aplicado a vocabulario de lane retirado.
- `relacionado-con` → [sdd-propose-contract-joins-lane-sweep.md](sdd-propose-contract-joins-lane-sweep.md) — precedente: un archivo fuera del Scope original que describe un mecanismo que el cambio retira, sumado al sweep.
