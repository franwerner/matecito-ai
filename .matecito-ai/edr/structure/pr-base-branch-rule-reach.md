# EDR — The PR base rule is branch-agnostic — it reaches every chain's last hop

- **Status:** Accepted
- **Date:** 2026-08-12

## Contexto
The base-branch resolution rule for pull requests is stated once, agnostic of any particular branch name — it never hardcodes which branch a PR targets. `feature-branch-chain`'s tracker PR was the case named in the brief, but `stacked-to-main`'s final hop faces the identical resolution question: what its last PR targets is also the working branch, not a fixed name. Fixing only the named case would leave the two chain strategies inconsistent — one honoring the rule, the other still naming a literal branch at its last hop.

## Decisión
The rule reaches every chain strategy that has a final hop, not only `feature-branch-chain`'s tracker PR. Both contingent sites this change scoped — `sdd-apply/SKILL.md:117` (`stacked-to-main`'s per-PR description) and `sdd-tasks/SKILL.md:178` (the forecast's description of the same strategy) — are updated together, in the same pass, so no chain strategy is left naming a literal branch as its destination.

## Reglas verificables
- **[manual]** Every chain strategy with a final hop (`feature-branch-chain`'s tracker PR, `stacked-to-main`'s last PR) targets the working branch, never a literal branch name.
- **[manual]** A change to the base-branch resolution rule updates every contingent site it reaches (`sdd-apply/SKILL.md`, `sdd-tasks/SKILL.md`) in the same pass, so the chain strategies never diverge from each other.

## Alternativas consideradas
Fix only `feature-branch-chain`'s tracker line, since it was the one named in the brief. Rejected: the resolution rule is written branch-agnostic on purpose, so leaving `stacked-to-main`'s last hop naming a literal branch would let one chain strategy honor the rule while the other silently didn't — an inconsistency between two things meant to behave the same way.

## Consecuencias
`sdd-apply/SKILL.md:117` and `sdd-tasks/SKILL.md:178` change together in this pass. A future change to the base-branch rule must check both contingent sites, not just the one a request happens to name.

## Relacionados
- `relacionado-con` → [pr-base-explicit-argument-form.md](pr-base-explicit-argument-form.md) — la otra decisión ratificada en el mismo cambio: cómo se enuncia la regla que este alcance extiende.
