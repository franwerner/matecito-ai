# Capability — Dos lanes fijos sin fork de usuario

- **Status:** Accepted
- **Date:** 2026-08-31
- **Components:** cli

## Propósito

El fork de lane pregunta al usuario por tamaño de cambio **antes** de que nadie haya leído el código, y lo pregunta de nuevo en el INTAKE GATE para confirmar un lane más cuatro flags. Ambos momentos cuestan un turno y ninguno está respaldado por evidencia. Este spec declara la vía de reemplazo: dos lanes fijos (`full` siempre; `direct` solo si el usuario lo pide explícitamente), sin fork, sin recomendación, sin confirmación.

## Actores

- **Orquestador**: ofrece dos lanes (ninguno de los cuales es elegido por el usuario, a menos que pida `direct` explícitamente)
- **Usuario**: pide `direct` explícitamente cuando necesita hacerlo; `full` corre sin pregunta
- **Intake**: decide sin investigar, sin clasificar tamaño, sin recomendar lane

## Precondiciones

- La solicitud del usuario llega al orquestador sin mediación de fork previo
- No hay clasificación de tamaño ni triage de lane guardada de antes

## Flujo principal

1. Usuario hace una solicitud
2. Si la solicitud nombra explícitamente `directo` o `direct` (variante), intake decide `direct` → no flujo
3. Si no, intake procesa con `full` → todos los fases, siempre
4. El orquestador no ofrece opción de lane

## Ramas / flujos alternativos

- **Solicitud imperativa sin mención de `direct`** → se trata como `full`, no se infiere `direct`
- **Solicitud trivial obvio** → sigue siendo `full`; el tamaño no decide
- **Solicitud que pide explícitamente directo** → `direct` corre, no flujo SDD

## Casos borde

- **El usuario dice "arreglá X directo"** → `direct`
- **El usuario dice "arreglá X"** → `full`
- **El usuario dice "agregá Y, corto"** → `full` (el tamaño no activa `direct`)
- **"¿podés arreglar Z?"** → forma interrogativa, pero misma semántica: pide trabajo, no ofrece lane, es `full`

## Reglas de negocio

- `full` es el default y NUNCA es cuestionado
- `direct` SOLAMENTE cuando el usuario lo pide por nombre explícitamente
- El orquestador no muestra un fork, no pide confirmación de lane, no hace recomendación
- Las únicas dos opciones son `direct` (explícitamente pedido) y `full` (todo lo demás)
- El vocabulario de `reduced`, `custom`, y "add-on" NO APARECE en ningún lado como mecanismo vivo

## Entidades y estados

- **Lane** — `direct` o `full`; estado: decidido (nunca confirmado)
- **Solicitud** — contiene intención del usuario; estados: ingresa → decision de lane → flujo dispatched

## Errores de cara al actor

- Si el usuario pide `direct` pero la solicitud está ambigua: se interpreta como `full` y se sigue adelante
- La ambigüedad no frena el flujo

## Requisitos

### Requisito: Two fixed lanes, neither chosen at a fork

The system MUST offer exactly two lanes. `full` MUST be the lane for every request by default and always. `direct` MUST run only when the user explicitly asks for it. The system MUST NOT surface a lane fork, MUST NOT produce a lane recommendation, and MUST NOT take a lane confirmation.

#### Scenario: A substantial request runs the full pipeline with no question asked

- GIVEN a request the system would previously have sized as substantial
- WHEN the work starts
- THEN the full pipeline runs and no lane choice is surfaced at any point
- AND no artifact records a lane recommendation

#### Scenario: `direct` runs only on an explicit request

- GIVEN a request the user words as an explicit ask for direct work
- WHEN the work starts
- THEN no flow phase runs
- AND nothing is asked to confirm that reading

#### Scenario: An imperative or an obviously trivial request is still `full`

- GIVEN a trivially small request phrased as an imperative, with no explicit ask for direct work
- WHEN the work starts
- THEN it runs `full`
- AND the system does not infer `direct` from the request's size or its grammatical form

### Requisito: Every phase of the pipeline always runs, and the add-on vocabulary is gone

Under `full`, every phase of the domain's pipeline MUST run; no phase is optionally present within a run. The vocabulary of "add-on", "optional phase", and the lanes `reduced` and `custom` MUST NOT appear anywhere as a live mechanism. The pipeline order MUST have a single declared source.

#### Scenario: No live add-on mechanism survives

- GIVEN the payload after the change
- WHEN it is searched for `reduced`, `custom`, and "add-on" as a lane or phase mechanism
- THEN nothing matches, apart from a governance record recounting its own history
- AND a match that is not about lanes at all — a word like "reduced" used in ordinary prose — is left alone

#### Scenario: One row declares the phase order

- GIVEN the domain's vocabulary table
- WHEN it is read
- THEN it carries one row listing all its phases in order, and no base/add-on split
- AND that row is the only declaration of the order in the domain fragment

#### Scenario: Nothing maps add-ons into slots

- GIVEN the kernel after the change
- WHEN an insertion map for optional phases is looked for
- THEN none exists, and nothing replaces it

### Requisito: Every phase reads its full upstream

Each phase MUST read every upstream artifact its Read/Write row names. No phase MUST carry a fallback to a nearer upstream, and no phase's `next_recommended` MUST branch on a lane. The presence-based gate for durable capability-specs and EDRs MUST survive verbatim — it is not a lane fallback.

#### Scenario: The specification phase reads the proposal unconditionally

- GIVEN a change entering the specification phase
- WHEN that phase resolves its inputs
- THEN it reads the proposal, and carries no branch for the proposal being absent

#### Scenario: No upstream-fallback prose survives

- GIVEN the payload after the change
- WHEN it is searched for a rule phrased as reading the nearest available upstream
- THEN nothing matches

#### Scenario: The presence-based gate is untouched

- GIVEN a repo with no durable capability-spec store
- WHEN a phase that reads capability-specs runs
- THEN it skips them silently, exactly as before this change

### Requisito: Every domain that delegates to the kernel stays consistent with it

The lane model lives in the kernel and every domain inherits it. A domain whose own text states which phases run MUST agree with the kernel after this change: its vocabulary table MUST carry no base/add-on split, its phases MUST NOT branch on a lane, and no text of it MUST delegate to a kernel section this change deletes. This reaches every domain the kernel serves, not only the one this change was scoped around. A domain's own phases, its decision-record mechanism, its own gates and its own discovery placement MUST NOT change: this requirement keeps a domain consistent with a kernel that moved, it does not redesign the domain.

#### Scenario: No domain is left pointing at a deleted kernel section

- GIVEN every domain fragment after the change
- WHEN each is read for how it resolves which phases run
- THEN none defers to a lane fork, and each states its own pipeline as always running
- AND no fragment cites a kernel section this change deleted

#### Scenario: A second domain's phases stop branching on a lane

- GIVEN a domain whose phase contracts carried per-lane fallbacks and lane-branching next steps
- WHEN each is read after the change
- THEN each reads its full upstream and names the next pipeline phase mechanically

#### Scenario: The rest of a second domain is untouched

- GIVEN a domain swept only for consistency with the kernel
- WHEN its diff is inspected
- THEN it changes only where it named the removed lanes, the removed gate, or the deleted kernel section
- AND its own phases, its decision-record mechanism, its own gates and where it asks its discovery questions are all unchanged

#### Scenario: A word that is not a lane is not swept

- GIVEN a file matching the search only because it uses a lane word in ordinary prose
- WHEN the sweep runs
- THEN that file is left untouched, and the fact that it is not a match is written down so it is not re-flagged
