# EDR — The worktrees skill retains its hand-run claim/release procedure, gated on whether `matecito-ai` is installed

- **Status:** Accepted
- **Date:** 2026-08-15

## Contexto
`matecito-ai turn claim`/`release` only exist where the tool is installed. A machine without it, for any write the scope criterion covers, still needs a way to take and release the turn.

## Decisión
The skill keeps its manual claim/guarded-write/release procedure — raw git commands, no CLI — for the one case that still needs it: a machine where `matecito-ai` is not installed. It lives in a final section, `## When matecito-ai is not installed`, that opens by naming its own condition before a single instruction. The section states plainly why the procedure is second-class: a turn claimed by hand carries no receipt the tool issued, so nothing distinguishes a hand-written record from a valid claim except that a human wrote it — the same trust the mechanism withdrew everywhere else once it stopped depending on anyone's word.

## Reglas verificables
- **[manual]** The manual claim/guarded-write/release procedure lives in its own gated section (`## When matecito-ai is not installed`), opening on the condition under which it applies.
- **[auto]** No manual claim or release heading survives outside that gated section — covered by the documentary test on the skill file.
- **[manual]** The gated section states that a turn claimed by hand carries no receipt the tool issued, and that this is why the procedure is second-class.

## Alternativas consideradas
Deleting the procedure outright — rejected: it would leave a machine without `matecito-ai` installed with no way to coordinate a shared-branch write at all.

## Consecuencias
The skill now has one primary procedure — `matecito-ai turn claim` / `turn release`, used whenever the tool is available — and one gated fallback for when it is not. The distinction that used to separate two audiences (a session with the guard registered vs. one without, further narrowed by what the guard's classifier happened to match) collapses to a single, simpler test: has the tool, or does raw git by hand.
