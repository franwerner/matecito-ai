# EDR — One shared-branch turn covers a whole consolidation round, not one turn per cherry-pick

- **Status:** Accepted
- **Date:** 2026-08-15

## Contexto
The `sdd-apply` consolidation run cherry-picks a batch's commits onto the round's container one at a time, in ascending task-id order. The turn could be claimed once for the whole loop, or claimed and released around each individual cherry-pick.

## Decisión
The consolidation run claims the turn once, before the first cherry-pick of its loop, and releases it once, after the loop ends — before the cleanup pass, which touches only per-task worktrees and branches, never the round's container.

## Reglas verificables
- **[manual]** The turn is claimed before the loop's first cherry-pick and released after the loop's last one, never per individual cherry-pick.
- **[manual]** The release happens before the cleanup pass, so cleanup — which never touches the round's container — never runs while the turn is still held.

## Alternativas consideradas
Claiming and releasing around each cherry-pick — rejected: it would let another session's integration land in the middle of a round, so the branch would end up carrying half of one change interleaved with another, and a conflict would stop being attributable to a single task — the exact property `structure/dispatch-batch-bound-integration.md` was written to preserve.

## Consecuencias
The round is one integration act, the same way the cycle-close merge is one. The loop already runs strictly one cherry-pick at a time inside a single non-isolated run, so a per-round turn adds no serialization the run did not already have — it only keeps another session's round from interleaving with this one.
