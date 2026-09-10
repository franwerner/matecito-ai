# EDR — `turn claim` performs the queue-and-retry mechanics, within one invocation; the caller only reads the report and asks

- **Status:** Accepted
- **Date:** 2026-08-15

## Contexto
`structure/merge-turn-wait-and-stop.md` requires a failed claim to queue the pending request, pause, and re-attempt once before stopping. That obligation splits into a mechanical half a program can perform deterministically, and a conversational half that needs a channel to a human. It has now been assessed twice for a performer with no such channel of its own — first a `PreToolUse` hook, now `turn claim`, an ordinary command a caller invokes and reads the result of directly — and the split re-derives to the same place both times.

## Decisión
`turn claim` performs the mechanical half, inside one invocation: it writes the queue entry, pauses briefly, re-attempts the claim exactly once, and deletes its own queue entry by compare-and-swap the moment the retry succeeds. If still held, it reports the holder, the age, and the queue count, and exits non-zero. The conversational half stays with the caller, which now has a channel to the user by construction — it is the one that invoked the command — rather than by reading a message off a blocked command: it puts those three facts to the user with a recommendation, choosing none of the three itself; a headless consolidation run returns `status: blocked` instead. Only the performer's name changes — first the agent, then the guard, now `turn claim` — the split itself does not.

## Reglas verificables
- **[auto]** `contend` (package `internal/turn`) writes the queue entry, pauses via the `sleep` seam, and re-attempts the claim exactly once before reporting the holder.
- **[auto]** A successful retry deletes the queue entry it wrote; a still-held retry leaves the queue entry in place — and since there is no release hook anymore, nothing later sweeps an entry left this way (the accepted cost recorded in `contracts/turn-claimed-by-explicit-call.md`).
- **[manual]** The report carries the holder, how long it has been held, and how many sessions are queued behind it — the three facts the skill tells the reader to put to the user.
- **[auto]** A failed queue write degrades the report (it is appended, not swallowed) but never the verdict — the call still exits non-zero.

## Alternativas consideradas
The command exits non-zero and the caller performs the whole queue-and-retry sequence by hand — rejected: it puts the deterministic part back into prose a caller has to remember, the exact failure mode this change exists to remove. The command asking the user directly — impossible, a CLI invocation has no channel of its own to a person; only the caller that invoked it does.

## Consecuencias
The skill's `## When a command comes back blocked` section carries no queue-and-retry steps of its own — it only tells the reader `turn claim` already performed them and describes what to do with the resulting facts.
