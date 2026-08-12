# EDR — The PR base rule is a mandatory argument, not a norm to repeat

- **Status:** Accepted
- **Date:** 2026-08-12

## Contexto
The base-branch resolution rule could be written two ways: as a norm ("aim the PR at the working branch") trusted to whoever opens it, or as a mandatory, explicit argument the PR-opening command itself must carry. Session evidence (obs #1410, a controlled three-run comparison of the same conflict scenario) showed a prose rule — even duplicated across multiple files and steps — let the executor reason around it and proceed past the exact case it was meant to catch; only when the same information was demanded as an obligatory field did the executor confront it directly, because omitting it produced a visibly malformed output instead of a silent shortcut.

## Decisión
The rule is stated as a mandatory argument: every PR-opening command MUST carry an explicit `--base <working-branch>`. Omitting it is a malformed command even when the tool's own fallback would have produced the correct branch — the omission is the defect, not whichever branch the fallback happens to pick. The web-UI equivalent follows the same shape: the base selector is set deliberately, never accepted as pre-filled.

## Reglas verificables
- **[manual]** `gh pr create` (or any PR-opening command) always carries an explicit `--base` naming the working branch, even when the tool's default would resolve to the same branch.
- **[manual]** A PR opened without an explicit base argument is read as a malformed command, not an acceptable shortcut, regardless of which branch the tool's fallback would have produced.

## Alternativas consideradas
State it as a norm ("open the PR against the working branch") and rely on phase prose or the agent's own judgment to apply it consistently. Rejected on session evidence: a prose rule, however many places it is repeated, lets the executor reason past the case it exists to catch; a mandatory field forces the value to be declared, and the declaration itself is what the rest of the ecosystem can act on.

## Consecuencias
Every PR-opening command in the payload — flow-driven or documented for a person — names `--base` explicitly; a command missing it is a defect to fix, not a stylistic choice. This deliberately stops short of a mechanical hook that blocks a bare `gh pr create`: that remains a possible follow-up if the explicit-argument form alone proves insufficient.

## Relacionados
- `relacionado-con` → [pr-base-branch-rule-reach.md](pr-base-branch-rule-reach.md) — la otra decisión ratificada en el mismo cambio: hasta dónde alcanza la regla que este enunciado fija.
