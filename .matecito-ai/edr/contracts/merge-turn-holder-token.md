# EDR — The shared-branch turn's holder value is a written git blob, not a reused commit sha

- **Status:** Accepted
- **Date:** 2026-08-15

## Contexto
`git update-ref` takes an object id as its value, so the turn ref cannot point directly at arbitrary text. Something has to prove who holds the turn well enough that a later release can be conditioned on it, and well enough that a stop-and-ask report can name the holder.

## Decisión
Claiming the turn writes a plain-text record (change, tree, destination, moment, base, an ISO-8601 UTC timestamp) as a git blob via `git hash-object -w --stdin`, then points `refs/matecito-ai/merge-lock` at that blob's sha with a create-only-if-absent `git update-ref`. **That sha is returned to the caller as its receipt.** Releasing deletes the ref conditioned on that same sha. The same record shape serves a queue entry.

## Reglas verificables
- **[manual]** The turn's value is always a blob written by `git hash-object -w --stdin`, never a commit sha, never `HEAD`, never an annotated tag.
- **[auto]** A release always conditions the delete on the sha the claim returned to its caller —
  `git update-ref -d refs/matecito-ai/merge-lock <sha>` — never an unconditional delete.
- **[auto]** The record carries exactly these fields, one per line: change, tree, destination,
  moment, base, at (UTC) — read with `git cat-file -p` for the blocked report.

## Alternativas consideradas
Pointing the ref at the current `HEAD` — rejected: two sessions standing on the same head would write the identical value, so a release could delete a turn that is not the releasing session's own, which the spec forbids outright. An annotated tag object — rejected: same result as a blob, with more machinery and no added property.

## Consecuencias
One object shape covers both the holder value and a queue entry, so the mechanism has no second shape to keep aligned. The blob is readable (`git cat-file -p`), which is what the stop-and-ask report needs to name the holder, and unique to its claim, which is what makes the conditional delete a real ownership proof. Verified in a scratch repository: creation with compare-and-swap, refusal of a second claim, release with the correct old value, and a clean working tree throughout.

Rewritten when the turn stopped depending on identity: the record was previously amended to carry a
session and a run identifier, and that amendment is withdrawn — it is no longer necessary, since no
caller ever asks whether a turn is already its own. The process id it originally carried is also gone,
for an unrelated reason: the claiming process now exits the instant its claim returns, so `pid` named a
dead process in every report it appeared in. Ownership is proved by holding the value's identifier —
the sha, returned to the caller as its receipt — not by matching any of the record's contents. This
file's pre-amendment wording, from before the session/agent amendment, was never committed and is not
recoverable.
