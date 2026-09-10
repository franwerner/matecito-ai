# EDR — Holding the receipt is what permits a release — not proving who you are

- **Status:** Accepted
- **Date:** 2026-08-15

## Contexto
The holder record carries descriptive fields (change, tree, destination, moment, base, a timestamp) so a person can diagnose a stuck turn. With identity gone from the mechanism, the temptation to let a caller free a turn by matching those fields against something it supplies remains live — the fields loosely identify the write the turn is for, which looks like enough. This gets its own record rather than living inside `contracts/merge-turn-holder-token.md`, because it carries a boundary that one does not: the stored fields are informative and never decisive.

## Decisión
Holding the receipt is what lets a caller give the turn back — not proving who it is. No code path may ever decide anything by comparing a field of the holder record against anything else; the fields exist to be printed in a report a person reads. Ownership is the receipt, and only the receipt.

## Reglas verificables
- **[auto]** No code path compares a field of the holder record against a caller-supplied value or against another record's field — `release` takes only a token, no other identifying option.
- **[auto]** A release with every field of the record different from the releasing caller's own still succeeds on the correct receipt.
- **[auto]** A release with every field identical to the releasing caller's own still fails on the wrong receipt.
- **[manual]** There is no route to a release that does not hold the exact value the claim returned — giving back the turn without the receipt is a usage error with no fallback path.

## Alternativas consideradas
Matching on `change` or `tree` as a fallback for a lost receipt — rejected: it would reopen every hole the three prior identity-based rounds closed, in a form that looks reasonable precisely because it resembles a legitimate diagnostic use of those fields.

## Consecuencias
The orphan policy — nobody frees another's turn on their own initiative — stops being a rule anyone has to remember and becomes a property of the code: there is no path that can delete a ref whose value the caller does not already hold. This is the rule most likely to erode over time, because the moment someone wants "release my own turn without the receipt", matching on a descriptive field looks like the obvious fix; this record exists so a later reader does not drift it back into a matching rule.
