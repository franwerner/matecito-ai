# Capability — Aislamiento de workspace en un único nivel, sin workspace de cambio

- **Status:** Accepted
- **Date:** 2026-08-13
- **Components:** cli

## Propósito

Concurrent sessions on one repository stop sharing a single working tree. Isolation exists at a single level: the per-task isolated runs an implementation batch opens on top of the working branch. There is no change-level workspace above it — nothing opens one, forwards its location in a dispatch prompt, integrates it, or cleans it up.

## Actores

- **Implementation batch** — opens per-task isolated runs on top of the working branch
- **Isolated run** — repositions onto the round's base, verifies the pre-write handshake, and writes inside its own per-task worktree
- **Consolidation run** — a second invocation of the implementation phase, without isolation; integrates the batch's commits into the working branch one at a time, holding the shared-branch turn while it does

## Flujo principal

1. The Uncommitted-Work Gate inspects the main repository before a round is dispatched.
2. The round's base is captured as the current head of the working branch.
3. Each isolated run repositions onto that base and opens its own per-task worktree.
4. Each isolated run verifies its head matches the base and its tree is clean before writing anything.
5. The consolidation run claims the shared-branch turn, integrates the batch's commits into the working branch one at a time in ascending task order, and releases the turn before the cleanup pass.

## Ramas / flujos alternativos

- **Base mismatch in pre-write check** → an isolated run detects its head does not match the received base or the tree is not clean; it reports `not-implemented / base-not-established` and does not proceed.
- **Commit parent mismatch in integration** → the consolidation run detects a reported commit's parent is not the expected base; the commit is not integrated, the report is recorded as `base-mismatch`, and the branch remains intact.
- **Uncommitted work found in the main repo** → the gate offers to work on the working branch directly, without a worktree; choosing that degrades the round onto that same branch.

## Casos borde

- **Cherry-pick conflict during integration** → the consolidation run's cherry-pick fails for one task's commit; it reports the conflict against that task, does not force it through, and leaves the branch and the remaining isolated runs' work intact.
- **Shared-branch turn already held** → a consolidation run whose claim is refused integrates nothing that round and reports blocked with the facts the claim reported (holder, branch, write moment, how long ago, and the queue).

## Reglas de negocio

- Isolation exists at one level only: the per-task isolated runs an implementation batch opens on the working branch. No change-level workspace exists — nothing opens one, forwards its location in a dispatch prompt, integrates it, or cleans it up.
- A round's base is always the current head of the working branch; the consolidation run's integration destination is that same branch; and the pre-write check compares against that base. None of the three is a choice between two containers.
- The Uncommitted-Work Gate inspects the main repository, always, before a round is dispatched; its *work on the branch without a worktree* outcome degrades the round onto that same branch, and the text governing it retains no conditional form over which container it inspects.
- An isolated run must verify its head equals the received base and its tree is clean before writing anything; failing either, it writes nothing and reports `not-implemented / base-not-established`.
- The consolidation run must verify each reported commit's parent equals the base it received; commits with mismatched parents are not integrated.
- The consolidation run — a second invocation of the implementation phase, without isolation, reusing its merge protocol — is the sole integrator of a batch's commits into the working branch. No phase integrates or prepares an integration toward any other branch, and the orchestrator does not mutate the repository.
- The consolidation run holds the shared-branch turn while it integrates: it claims the turn once before the first cherry-pick and releases it once after the last, regardless of the loop's outcome.
- An integration failure is attributed to the task whose commit failed; there is no second level of failure to attribute it to.

## Entidades y estados

- **Isolated run** — a git worktree opened on its own per-task branch. States: repositioned onto base (open) → committed → integrated → cleaned up.
- **Base of a round** — the current head of the working branch; used to verify each isolated run's head in the pre-write check, and each reported commit's parent in the integration check.
- **Shared-branch turn** — held by the consolidation run for the duration of its integration loop; claimed once before the first cherry-pick, released once after the last.

## Errores de cara al actor

- **Base not established** → reported as `not-implemented / base-not-established`; the run wrote nothing.
- **Base mismatch during integration** → reported as `base-mismatch`; the commit is not integrated, the branch is left intact.
- **Turn refused** → the consolidation run reports blocked with the holder, its branch, its write moment, how long ago, and the queue; nothing is integrated that round.

## Escenarios

### Scenario: the round's base is the working branch

- **GIVEN** a round eligible for parallel dispatch
- **WHEN** its base is captured
- **THEN** it is the working branch's head
- **AND** there is no other container option to evaluate

### Scenario: consolidation integrates into the working branch

- **GIVEN** that round's isolated runs having returned their commits
- **WHEN** the consolidation run integrates them one at a time, in ascending task order
- **THEN** each lands on the working branch

### Scenario: the pre-write check compares against the received base

- **GIVEN** an isolated run on the working branch
- **WHEN** its pre-write check runs
- **THEN** it compares its head against the base it received in its dispatch prompt
- **AND** there is no second container it could be compared against

### Scenario: the Uncommitted-Work Gate inspects the main repository

- **GIVEN** a round eligible for dispatch, about to be dispatched
- **WHEN** the uncommitted work that might intersect it is inspected
- **THEN** the main repository is inspected
- **AND** the text governing that gate retains no conditional form over which container it inspects

### Scenario: the degrade outcome names one destination

- **GIVEN** the gate offering to work on the branch without a worktree
- **WHEN** that outcome is chosen
- **THEN** the round runs on the working branch
- **AND** that outcome's description names no other possible container

### Scenario: no phase receives a workspace's location

- **GIVEN** any phase dispatch prompt after the change
- **WHEN** it is read
- **THEN** it carries no line with a change workspace's location
- **AND** the phase resolves its paths against the tree it was launched from

### Scenario: both handshake checks still stated

- **GIVEN** the durable record governing the base handshake, after the change
- **WHEN** it is read
- **THEN** it carries the pre-write check and the pre-integrate check, both complete
- **AND** each is named by the end of the flow it sits at

### Scenario: the end-names are not deleted as nesting vocabulary

- **GIVEN** that same record, rewritten by this change
- **WHEN** it is compared with its prior version
- **THEN** both end-names — pre-write and pre-integrate — are still present
- **AND** the only thing retired is the base being able to come from two different containers

### Scenario: one base, no condition

- **GIVEN** the statement of what each handshake check compares against
- **WHEN** it is read
- **THEN** it names a single base, the working branch's
- **AND** no clause makes it depend on whether any isolation is active

### Scenario: the integrator is the implementation phase, not the orchestrator

- **GIVEN** a round whose isolated runs have returned their commits
- **WHEN** it is integrated
- **THEN** a second invocation of the implementation phase, without isolation, performs it
- **AND** neither the orchestrator nor a new agent integrates anything

### Scenario: the orchestrator does not mutate the repository, again without exception

- **GIVEN** the statement of what the orchestrator mutates, after the change
- **WHEN** it is read
- **THEN** it is absolute, with no change-level exception
- **AND** no phase integrates or prepares an integration toward any other branch

### Scenario: the integrator holds the turn

- **GIVEN** a consolidation run about to integrate
- **WHEN** its bracket runs
- **THEN** it holds the shared-branch turn for the whole integration
- **AND** it releases it after the last cherry-pick, before the cleanup pass

### Scenario: an integration failure is attributed to its task

- **GIVEN** a batch whose integration of one commit fails
- **WHEN** the failure is reported
- **THEN** it names that task
- **AND** there is no second level of failure to attribute it to

### Scenario: every retired record is fully gone

- **GIVEN** each durable record this change retires
- **WHEN** it is looked for after the change
- **THEN** its file does not exist
- **AND** its row in its domain's index does not either

### Scenario: nothing cites a retired record as a live rule

- **GIVEN** the durable record store after the change
- **WHEN** it is searched for references to each retired record
- **THEN** no match is a live rule depending on it

### Scenario: historical citers carry their note, not a repaired citation

- **GIVEN** the records that cited a retired one as historical precedent
- **WHEN** they are read after the change
- **THEN** each carries a note saying the target is retired and why that is expected
- **AND** none was re-pointed to another record

### Scenario: rewritten records carry no link to the retired ones

- **GIVEN** every durable record this change rewrites
- **WHEN** it is read after the change
- **THEN** it links to none of the retired ones

### Scenario: the behavior store's index describes the surviving capability

- **GIVEN** the index of the type this capability belongs to
- **WHEN** its row is read
- **THEN** it describes per-task isolation on the working branch
- **AND** it no longer describes a two-level nested isolation or an opt-in mechanism

### Scenario: an unestablished base implements nothing

- **GIVEN** an isolated run whose head does not match the base it received, or whose tree is not clean
- **WHEN** it reaches the point of writing
- **THEN** it writes nothing, produces no commit, and reports `not-implemented / base-not-established`

### Scenario: a commit whose parent is not the base is not integrated

- **GIVEN** a reported commit whose parent differs from the base the consolidation run was handed
- **WHEN** that report is processed
- **THEN** the commit is not integrated, the report is recorded as `base-mismatch`, and its branch is left intact for inspection
