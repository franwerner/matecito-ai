# EDR — The shared-branch turn's waiting list is a set of git refs, not an entry in a memory store

- **Status:** Accepted
- **Date:** 2026-08-15

## Contexto
A session that finds the turn taken needs somewhere to record that it is waiting, in the order requests arrived, without that record itself being mistaken for the lock.

## Decisión
A waiting session creates `refs/matecito-ai/merge-queue/<change-name>`, pointing at its own record blob (same shape the holder writes — see `contracts/merge-turn-holder-token.md`). It deletes that ref, by compare-and-swap on its own value, the moment its claim succeeds. Reading the queue is one `git for-each-ref refs/matecito-ai/`; reading an entry is `git cat-file -p <ref>`.

## Reglas verificables
- **[manual]** A queue entry lives at `refs/matecito-ai/merge-queue/<change-name>`, created and deleted only by its own owning session.
- **[manual]** The queue is never read to decide who holds the turn — only `refs/matecito-ai/merge-lock` decides that; an entry existing is evidence of a pending request, never evidence of holding the turn.
- **[manual]** Reading the whole queue is `git for-each-ref refs/matecito-ai/` in one call; reading one entry's fields is `git cat-file -p <ref>` — no second store, no per-entry round trip.

## Alternativas consideradas
A single shared key in a memory store — rejected outright and always: recording under a key is an unconditional overwrite, so two sessions writing at the same moment lose an entry, which is the same property that makes such a store incapable of being the lock either. One memory key per session — rejected: it worked, but split the mechanism across two stores and made 'what is queued' a search whose completeness had to be verified, instead of one deterministic call.

## Consecuencias
The whole mechanism collapses into one store — the one that already provides the atomic primitive the lock depends on. Every worktree of the repository sees the same refs, so a queue entry is visible from wherever the waiting session runs. Accepted cost, ratified at the gate: refs belonging to dead sessions accumulate, because no session may delete a ref that belongs to another — the orphan policy applies to the queue exactly as it applies to the lock, and nothing here cleans up after another session. Verified: `for-each-ref` listed the lock and two queue entries together, with no truncation and no second call per entry, and `git status --porcelain` stayed empty throughout.
