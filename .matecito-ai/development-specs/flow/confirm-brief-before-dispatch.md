# Capability — Confirmar el brief de intake antes de despachar

- **Status:** Accepted
- **Date:** 2026-09-01
- **Components:** cli

## Propósito

Nada se sienta entre que `sdd-intake` retorna el brief y la fase siguiente se despacha, pero todo lo que viene leyendo ese brief como alcance *confirmado* sin que nadie lo haya confirmado. Un stop obligatorio — una pregunta sobre el brief completo, dos respuestas — cierra ese fallo: nada se despacha hasta que el usuario responde.

## Actores

- **Orquestador** — presenta el brief, espera, y no despacha nada hasta que la respuesta llega
- **Usuario** — acepta el brief tal como está, o escribe una corrección
- **Fase de intake** — el único escritor del brief; corre de nuevo sobre el pedido original más la corrección

## Precondiciones

- La fase de intake ha retornado y su brief está persistido bajo su propia clave de artefacto
- Ninguna fase después de intake ha sido despachada

## Flujo principal

1. La fase de intake retorna el brief.
2. El orquestador lo presenta como **un solo** item, a través del walkthrough de presentación compartida, y pregunta **una sola** pregunta.
3. El usuario responde: **aceptar** o **corregir**.
4. On accept, the flow proceeds — the next phase is dispatched.
5. On correct, the intake phase is re-dispatched with the original request plus the user's words verbatim; it overwrites its own brief, and the gate runs again over what came back.

## Ramas

- **Correction, repeated** → the loop has no bound; it closes when the user accepts.

## Casos borde

- **A correction that changes a decision flag** → the re-run brief carries the new value, and nothing built on the old one exists yet to undo.
- **A brief with nothing the user wants to change** → one answer closes the gate; no per-flag question is ever asked.

## Reglas de negocio

- Nothing is dispatched until the question is answered.
- The brief is exactly one ratifiable item and gets exactly one question.
- Exactly two answers exist: accept and correct. There is no cancel.
- The orchestrator never writes the brief; only the intake phase does.
- The gate is not the removed intake gate: it confirms no lane, walks no flag, offers no cancel, and keys off no execution mode.

## Entidades y estados

- **Brief** — carries the request and decision flags; states: emitted by intake → presented at gate → accepted or corrected → emitted again if corrected
- **Gate state** — states: open (waiting for answer) → closed (on accept or after re-dispatch returns)

## Errores de cara al actor

- **An answer that is neither an acceptance nor a correction** → the question stands; nothing is dispatched and nothing is assumed.

## Requisitos

### Requisito: A mandatory stop sits between the brief and the next dispatch

The orchestrator MUST NOT dispatch any phase after intake until the user has answered the brief-confirmation question. The governing text MUST carry this stop as its own named section, positioned immediately before the section that states how phases run, and that section's "phases run back-to-back, no between-phase checkpoint" claim MUST be scoped so it governs the phases *after* this stop rather than contradicting it.

#### Scenario: Nothing moves before the answer

- GIVEN a brief that has just returned
- WHEN the orchestrator holds it
- THEN no later phase is dispatched
- AND it stays so until the user answers

#### Scenario: The back-to-back claim is scoped, not deleted

- GIVEN the section stating how phases run
- WHEN it is read after the change
- THEN it still states that phases run back-to-back with no between-phase checkpoint
- AND it names this one stop as preceding that run, so the two claims do not contradict

#### Scenario: The lane text keeps its own no-confirmation claim

- GIVEN the text fixing the two lanes
- WHEN it is read after the change
- THEN it still states that no lane is forked, recommended or confirmed
- AND that claim is scoped to the lane, not to the brief

#### Scenario: Nothing left stating a workspace opens at this gate

- GIVEN the text governing this gate, after the change
- WHEN it is read end to end
- THEN it does not mention opening any change workspace, neither before nor after the answer

### Requisito: The brief is offered as one item, with one question

The orchestrator MUST present the brief as **exactly one** ratifiable item and ask **exactly one** question over it. It MUST delegate the presentation to the shared walkthrough rather than stating a presentation of its own; with one item, that walkthrough resolves to its fixed item template alone — no index, no "confirm the rest". The item's anchor MUST be the brief's own artifact key. Its summary MUST be one line carrying intake's reading of the request **alone**; no decided flag value MUST be folded into it. The brief MUST be modelled as a **compound item**: the brief's change type and each flag intake decided print as their own field line beneath the summary and above the actions, in the order intake writes them into the brief's classification, carrying the brief's own value verbatim. `Domains touched` MUST NOT print — it is the only line of the classification the item leaves out, and it stays reachable through the item's detail retrieval. The governing text MUST enumerate which lines print rather than describing them by category: an orchestrator holding only that text, with no memory of this change, MUST be able to produce the right set from it. A flag whose axis is not declared for the project MUST NOT print a field line, and MUST NOT be stood in for by a placeholder. No per-flag question MUST exist anywhere, and printing a flag as a field line MUST NOT be read as asking about it. The governing text MUST NOT state anywhere that the flag values are folded into the summary. (Previously: the same obligation, with `Worktree isolation` among the enumerated field lines and with four flags decided.)

#### Scenario: One item, no index

- GIVEN a brief presented at this gate
- WHEN the presentation is inspected
- THEN it is the fixed item template alone, with no index and no bulk shortcut
- AND the count rule of the shared walkthrough produced that form, not an exception written for this gate

#### Scenario: The flags are not walked one by one

- GIVEN a brief carrying every decision flag intake decided
- WHEN the gate runs
- THEN one question is asked over the whole brief
- AND no question is asked about any individual flag, at this gate or anywhere else

#### Scenario: The whole brief is one retrieval away

- GIVEN an item whose summary is one line, whose field lines carry the decided flags, and whose anchor is the brief's artifact key
- WHEN the user asks to see the detail
- THEN the whole brief is retrieved through that anchor
- AND nothing of the brief is hidden by the summary being one line

#### Scenario: Every decided flag prints as its own field line

- GIVEN a brief carrying the three decided flags
- WHEN the item is presented
- THEN each one prints as its own field line beneath the summary and above the actions
- AND the summary carries intake's reading of the request with no flag value in it

#### Scenario: One item, one anchor, one outcome

- GIVEN the brief presented as a compound item
- WHEN the gate counts what it is holding
- THEN it is one item with one anchor, and it takes exactly one outcome
- AND no field line takes an outcome of its own

#### Scenario: A project without the components axis prints one field line fewer

- GIVEN a project that declares no components axis
- WHEN the item is presented
- THEN the change type and the two remaining decided flags print, one field line each
- AND no placeholder or "n/a" line stands in for the absent one

#### Scenario: The domains touched stay out of the field lines

- GIVEN a brief whose classification also carries the change type and the domains touched
- WHEN the item is presented
- THEN the change type prints as a field line and the domains touched does not
- AND the domains touched stays reachable through the item's detail retrieval

#### Scenario: The enumerated list no longer names isolation

- GIVEN the text enumerating which field lines print, after the change
- WHEN it is read
- THEN it names the change type and the three remaining flags
- AND it names no workspace-isolation field line

### Requisito: The two actions go through the host's question widget

The two actions MUST be offered through the host's question control, per the shared walkthrough's own discrete-options rule, with the item's summary as the question and one line of cost per action. The correction the user then writes MUST be prose, per that same rule's second half.

#### Scenario: A closed choice takes the widget

- GIVEN a gate offering accept and correct
- WHEN the question is put
- THEN it is put through the host's question control, with each action carrying one line of what it costs
- AND the anchor travels in the question's own text

#### Scenario: The correction itself is written, not chosen

- GIVEN a user who picks the correct action
- WHEN they state the correction
- THEN they write it as prose, with no options enumerated for them

### Requisito: Exactly two answers exist

The gate MUST offer exactly two answers: accept as-is, and correct. A cancel path MUST NOT exist.

#### Scenario: No cancel is offered

- GIVEN the gate as presented
- WHEN its actions are counted
- THEN there are two, and neither discards the change

#### Scenario: The removed gate's third path is not reinstated

- GIVEN the governing text after the change
- WHEN it is searched for a way to cancel the change at this point
- THEN nothing offers one

### Requisito: A correction re-runs the intake phase, which overwrites its own brief

A correction MUST be resolved by re-dispatching the intake phase with the original request plus the user's words verbatim. That phase MUST overwrite its own brief under the same artifact key. The orchestrator MUST NOT edit the brief itself. The gate MUST then run again over what came back, and the loop MUST close only when the user accepts. No new plumbing MUST be introduced to propagate the accepted version: downstream phases retrieve the key at their own dispatch time.

#### Scenario: The correction goes back to intake, verbatim

- GIVEN a user's written correction
- WHEN it is resolved
- THEN the intake phase is re-dispatched with the original request plus those words, unchanged
- AND the orchestrator writes nothing into the brief itself

#### Scenario: One writer, one key

- GIVEN a change whose brief was corrected two times
- WHEN the artifact key is inspected
- THEN every version of it was written by the intake phase
- AND the key holds the accepted version, having been overwritten rather than duplicated

#### Scenario: The loop closes on acceptance

- GIVEN a corrected brief presented again
- WHEN the user corrects it again
- THEN it is re-dispatched again, with no bound on the number of rounds
- AND the flow proceeds only after an acceptance

#### Scenario: Downstream sees the accepted version with no extra step

- GIVEN a phase dispatched after the gate closed
- WHEN it retrieves the brief
- THEN it gets the accepted version
- AND nothing forwarded it there by a second channel

### Requisito: This is not the removed intake gate

The gate MUST NOT be named after the removed intake gate. It MUST NOT confirm a lane, MUST NOT walk the decision flags one by one, MUST NOT offer a cancel, and MUST NOT key off an execution mode. Only its position in the flow and its "nothing dispatches until it is answered" character carry over.

#### Scenario: The name does not re-import the old semantics

- GIVEN the new section
- WHEN its name is read
- THEN it names what it confirms — the brief — and does not reuse the removed gate's name

#### Scenario: No lane is confirmed at it

- GIVEN the gate running over a brief
- WHEN what it asks about is inspected
- THEN it asks about the brief and never about which lane runs

## Referencias

- **Conceptualmente relacionado**: `flow/two-fixed-lanes.md` — define el modelo de dos lanes fijos, y donde el gate no tiene rol
- **Conceptualmente relacionado**: `flow/ratify-gate-items.md` — define cómo un gate presenta una vez que disparó
- **Conceptualmente relacionado**: `rule/intake-passthrough-contract.md` — define el contrato de intake que produce el brief
