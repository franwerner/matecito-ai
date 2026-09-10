# EDR — Taking and giving back the turn are two explicit commands, with nothing watching and nothing identifying anyone

- **Status:** Accepted
- **Date:** 2026-08-15

## Contexto
The mechanism went through three rounds of hook-based enforcement, and each round required identifying which session or sub-agent held the turn — first via a session identifier, then via a (session, agent) pair — chasing edge cases the identity model kept producing. Measured against real transcripts, every merge live sessions actually ran contained `&&`, so a `PreToolUse` classifier would have been bypassed regardless of how precisely it matched a command string; a quote-aware rewrite still let `$(git merge x)` execute unguarded while blocking an unrelated quoted string. The direct observation that ended the chase: identity is only needed to answer "is the turn already mine", and that question only exists because a caller could check its own turn mid-work.

Three records that fixed pieces of that retired hook-based model are superseded outright by this decision, since their decisions are now false: `merge-guard-registration-and-matching` fixed how the `PreToolUse` handler registered and classified commands, and no handler exists anymore to register; `turn-released-on-both-session-end-events` fixed the asymmetric `Stop`/`SubagentStop` release rule, and neither event is registered anymore; `wiring-paragraph-defers-to-the-guard` fixed that the turn is taken and released by a guard automatically, and that claim is no longer true.

## Decisión
Taking and giving back the turn are two explicit commands (`matecito-ai turn claim`, `matecito-ai turn release`) around one unit of work, with nothing watching and nothing identifying anyone. Every automatic hook is deleted, and with them the whole question of who holds what. Enforcement is `claim`'s exit code composing with `&&`: a chained write never runs behind a refused claim. If taking and giving back bracket one unit of work, nobody ever has to ask whether the turn is already theirs — the one write-moment the mechanism exists for is a single unit: the consolidation round takes the turn once before its loop and gives it back once after.

## Reglas verificables
- **[auto]** No PreToolUse/Stop/SubagentStop handler is registered for this mechanism — `hook.ForDomains(["development"])` returns only the commit-message hook.
- **[auto]** `claim` returns exactly three distinguishable outcomes: claimed (receipt, exit 0, silent), held (report, exit 1), not arbitrated (message on stderr, exit 0).
- **[manual]** The one write-moment the mechanism guards claims once and releases once around a single unit of work — the consolidation round around its whole cherry-pick loop — never a second claim inside that bracket.
- **[auto]** A caller that claims twice inside its own bracket is refused exactly as any other holder would refuse it — the mechanism does not and cannot detect that the holder is itself, and the refusal report names the caller's own change, tree and moment, which makes the mistake legible without the mechanism deciding anything from it.

## Alternativas consideradas
Identity-based interception, three rounds in a row (session id, then the (session, agent) pair, then a quote-aware command classifier) — rejected outright: every live-transcript merge contained `&&` and would have been missed by a `PreToolUse` classifier regardless of how precisely it matched, and the quote-aware rewrite still let a subshell invocation execute unguarded while blocking an unrelated quoted string. Command-string inspection cannot be made both complete and precise at once.

## Consecuencias
The whole question of "is this session's own turn" disappears — self-blocking becomes legible (the report names the caller's own change) rather than something the mechanism has to get right. Accepted cost: no automatic release — a session that dies holding the turn leaves it held until a human clears it, which `structure/merge-turn-wait-and-stop.md`'s orphan policy already required for exactly that case. Surface kept deliberately small: three subcommands, no machine-readable flag, since the only value a caller carries is the receipt.
