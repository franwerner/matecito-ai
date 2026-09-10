# EDR — A session finding the shared-branch turn taken re-attempts once, then stops and asks — never waits unbounded, never breaks it

- **Status:** Accepted
- **Date:** 2026-08-15

## Contexto
A session whose claim fails needs a bounded reaction: waiting forever is indistinguishable from a hang, and breaking another session's turn contradicts the orphan policy already ratified for this mechanism.

## Decisión
A failed claim queues the session's pending request, pauses briefly, and re-attempts the claim exactly once. If still taken, the session stops trying. An orchestrator reads the state back from the refs (`git for-each-ref refs/matecito-ai/` plus `git cat-file -p`) to report who holds the turn, since when, and what is queued behind it, then asks the user how to proceed — wait and retry, leave it, or release it under the user's explicit instruction — with a recommendation, choosing none of the three itself. A headless `sdd-apply` consolidation run has no channel to ask through: it integrates nothing and returns `status: blocked` with the same three facts, leaving its Task Run Reports, commits and branches untouched so a later re-dispatch integrates them unchanged.

## Reglas verificables
- **[manual]** A failed claim is re-attempted exactly once, after a brief pause, before the session stops trying.
- **[manual]** No session breaks another session's turn on its own initiative, whatever the evidence — not after a wait, not on a stale timestamp, not on an absent holder; only the user's explicit instruction releases a turn that is not the releasing session's own.
- **[manual]** The stop-and-ask report is built by reading the refs themselves (`for-each-ref` plus `cat-file -p`), never from a second store.
- **[manual]** A headless consolidation run that cannot get the turn returns `status: blocked` with the holder, since-when, and queue facts, integrating nothing and leaving every Task Run Report untouched for a later re-dispatch.

## Alternativas consideradas
Waiting until the turn frees — rejected: an unbounded wait inside an agent run is indistinguishable from a hang, with no channel to say why it is waiting. A timeout that breaks the lock — rejected outright by the ratified orphan policy: neither elapsed time nor an absent holder authorizes one session to break another's turn.

## Consecuencias
Two attempts cover the common case — a turn held for the seconds an integration takes — and anything longer becomes a human's call, exactly as the ratified orphan policy already says. The report needs no second store to be accurate, since it is built from the refs themselves. Accepted cost: a crashed holder blocks the queue until a human intervenes.
