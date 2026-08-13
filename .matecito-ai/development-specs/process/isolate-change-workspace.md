# Capability — Aislar el workspace de cambio en un nivel propio

- **Status:** Accepted
- **Date:** 2026-08-13
- **Components:** cli

## Propósito

Concurrent sessions on one repository stop sharing a single working tree. Isolation nests in two levels: a per-change workspace the orchestrator opens and later integrates, and the per-task workspaces an implementation batch already opens — from now on, inside the change's workspace rather than on the original branch. The mechanism is opt-in; with it off, every existing path behaves exactly as before.

## Actores

- **Orchestrator** — opens the change-level workspace when isolation is active, and performs the final integration into the original branch once the cycle closes
- **Phase agents** — run and write files inside the change workspace when it is open; they receive its path in their dispatch prompt
- **Implementation batch** — opens per-task isolation inside the change workspace instead of the original branch
- **Consolidation run** — integrates the batch's commits into the change workspace, not the original branch

## Flujo principal

1. User chooses isolation at the lane fork or confirms it with the change's scope
2. Orchestrator opens the change workspace on its own branch, on first file write (for direct/ad-hoc work) or after scope confirmation (for flow work)
3. Every subsequent phase writes files inside the change workspace
4. When a parallel implementation batch is eligible, its isolated runs open on top of the change workspace
5. The consolidation run integrates commits into the change workspace
6. After the pipeline closes, the orchestrator merges the change workspace into the original branch

## Ramas / flujos alternativos

- **Isolation is inactive** → the system behaves exactly as before this change; the base and integration destination remain the working branch, and no workspace is opened.
- **Opening moment differs by lane** → in the flow, the workspace opens after the INTAKE GATE once the scope is confirmed; in direct/ad-hoc work, it opens right before the first file is written.
- **Base mismatch in pre-write check** → an isolated run detects its head does not match the received base or the tree is not clean; it reports `not-implemented / base-not-established` and does not proceed.
- **Commit parent mismatch in integration** → the consolidation run detects a reported commit's parent is not the expected base; the commit is not integrated, the report is recorded as `base-mismatch`, and the branch remains intact.

## Casos borde

- **Merge conflict when closing the cycle** → the orchestrator attempts to merge the change workspace and encounters a conflict against the original branch; it reports the failure and leaves both branches intact for inspection, never forced.
- **Workspace already open for the same change** → no second workspace is opened; the existing one is reused.
- **A trivial change without isolation** → no workspace is opened and nothing about it is mentioned.
- **Change workspace opening failure** → reported as a system error; the change is marked as unable to proceed, and the original branch remains untouched.

## Reglas de negocio

- Isolation is opt-in and never activated by default; the user must choose it explicitly at the lane fork.
- Once confirmed, the isolation choice applies to the entire change; no later phase re-asks it.
- The base handshake checks are performed against the immediate container's head — the change workspace when nested isolation is active, the original branch when it is not.
- An isolated run must verify its head equals the received base and its tree is clean before writing anything; failing either, it writes nothing and reports `not-implemented / base-not-established`.
- The consolidation run must verify each reported commit's parent equals the base it received; commits with mismatched parents are not integrated.
- The orchestrator alone performs the final merge of the change workspace into the original branch, exactly once, after the pipeline finishes.
- With isolation inactive, the behavior is unchanged; no new prompts, gates, or fields appear.

## Entidades y estados

- **Change workspace** — a git worktree opened on its own branch (`matecito-ai/<change-name>`), stored in `.matecito-ai/workspaces/<change-name>`. States: created (open) → integrated (merge complete) → cleaned up (branch deleted). Transitions triggered by orchestrator actions.
- **Isolation state** — for a change: `active` or `inactive`, confirmed at scope confirmation and persisted for the entire cycle.
- **Base of an isolated run** — the current head of its immediate container (the change workspace when nested, the original branch when not); used to verify the run's head in the pre-write check.

## Errores de cara al actor

- **Workspace opening failed** → reported as a system-level failure; the change cannot proceed with isolation.
- **Merge conflict closing the cycle** → reported with the conflict details and the branches left intact for manual inspection.
- **Base not established** → reported as `not-implemented / base-not-established`; the run wrote nothing.
- **Base mismatch during integration** → reported as `base-mismatch`; the commit is not integrated, the branch is left intact.

## Escenarios

### Scenario: the choice is offered together with the lane

- **GIVEN** a request the system reads as substantial
- **WHEN** the lane options are surfaced
- **THEN** change-level isolation is surfaced with them as an explicit choice, with a recommendation stated and the decision left to the user

### Scenario: nothing activates it silently

- **GIVEN** a change whose lane was resolved without the isolation choice ever being made
- **WHEN** any later step needs to know whether isolation is active
- **THEN** it counts as inactive, no workspace is opened, and nothing assumes otherwise

### Scenario: confirmed once, with the scope

- **GIVEN** an in-flow change whose brief is presented for confirmation
- **WHEN** the user confirms or adjusts it
- **THEN** the isolation choice is settled in that same act, and no later phase asks about it again

### Scenario: a trivial change is untouched by the mechanism

- **GIVEN** a trivial change that goes direct without the lane fork being surfaced at all
- **WHEN** it runs
- **THEN** isolation is inactive and nothing about it is mentioned anywhere

### Scenario: in the flow, once the scope is confirmed

- **GIVEN** an in-flow change with isolation active whose scope has just been confirmed
- **WHEN** the next phase is about to be dispatched
- **THEN** the change's workspace already exists, and every phase that writes files works inside it

### Scenario: direct work, right before the first file

- **GIVEN** direct or ad-hoc work with isolation active and no scope to confirm
- **WHEN** the first file of the work is about to be created or modified
- **THEN** the change's workspace is opened first, and that file is written inside it

### Scenario: opened once per change

- **GIVEN** a change whose workspace is already open
- **WHEN** a later phase runs, or the session resumes
- **THEN** that same workspace is used and no second one is opened for the same change

### Scenario: never opened before there is a scope to isolate

- **GIVEN** an in-flow change still resolving its discovery questions
- **WHEN** the scope has not been confirmed yet
- **THEN** no workspace exists yet and no file of the change has been written

### Scenario: the round's base is the change's workspace

- **GIVEN** change-level isolation active and a round eligible for parallel dispatch
- **WHEN** the base for that round is captured
- **THEN** it is the current head of the change's workspace, not of the original branch

### Scenario: consolidation integrates into the change's workspace

- **GIVEN** the isolated runs of that round returned their commits
- **WHEN** the consolidation run integrates them one at a time, in ascending task order
- **THEN** each lands on the change's workspace and the original branch receives nothing from the batch

### Scenario: with isolation off, the batch behaves as before

- **GIVEN** a change without isolation and a round eligible for parallel dispatch
- **WHEN** the round is dispatched and then consolidated
- **THEN** the base and the integration destination are the working branch, exactly as before this change

### Scenario: the nested base is what the pre-write check compares against

- **GIVEN** an isolated run nested on an active change workspace
- **WHEN** it runs its pre-write check
- **THEN** it compares its head against the change workspace's base it received, never against the original branch

### Scenario: an unestablished base implements nothing

- **GIVEN** an isolated run whose head does not match the base it received, or whose tree is not clean
- **WHEN** it reaches the point of writing
- **THEN** it writes nothing, produces no commit, and reports `not-implemented / base-not-established`

### Scenario: a commit whose parent is not the base is not integrated

- **GIVEN** a reported commit whose parent differs from the base the consolidation run was handed
- **WHEN** that report is processed
- **THEN** the commit is not integrated, the report is recorded as `base-mismatch`, and its branch is left intact for inspection

### Scenario: the orchestrator closes the cycle

- **GIVEN** a change with isolation active whose pipeline has finished
- **WHEN** the cycle closes
- **THEN** the orchestrator integrates the change's workspace into the original branch, once

### Scenario: no phase integrates the change level

- **GIVEN** any phase of a change with isolation active
- **WHEN** it finishes its work
- **THEN** it has neither integrated the change's workspace into the original branch nor prepared that integration

### Scenario: a conflict against the original branch is reported, not forced

- **GIVEN** a change-level integration that cannot complete cleanly
- **WHEN** the orchestrator attempts it
- **THEN** it reports the failure as the change's, against the original branch, and leaves both the original branch and the workspace as they were

### Scenario: the two levels are attributed separately

- **GIVEN** a batch whose integration of one task's commit fails while the change level is untouched
- **WHEN** the failure is reported
- **THEN** it names that task, and says nothing about the change against the original branch

### Scenario: a full change with the option off

- **GIVEN** a change taken end to end with isolation inactive
- **WHEN** every phase runs, including a parallel implementation batch
- **THEN** the observable behavior is identical to before this change, with no extra step and no extra prompt

### Scenario: no new gate appears for work that did not opt in

- **GIVEN** direct or ad-hoc work with isolation inactive
- **WHEN** files are written
- **THEN** they are written on the working tree as before, and nothing asks about or mentions a workspace

## Referencias

- **`deploy-payload-to-host.md`** → [`./deploy-payload-to-host.md`](./deploy-payload-to-host.md) — the change workspace must be excluded from the deployed payload at the target host.
