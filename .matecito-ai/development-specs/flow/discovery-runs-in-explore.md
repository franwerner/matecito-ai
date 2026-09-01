# Capability — Discovery es formulada y resuelta dentro de exploración

- **Status:** Accepted
- **Date:** 2026-08-31
- **Components:** cli

## Propósito

Las preguntas de descubrimiento de un cambio NO se formulan sin haber visto el código. El orquestador pide descubrimiento ANTES de que intake corra, guardando los archivos del cambio sin leerlos. El cambio de esta spec: discovery vive dentro de `sdd-explore`, que SIEMPRE corre (parte de `full`) y ha leído el código **antes** de preguntar.

## Actores

- **Explorador de código** (`sdd-explore`): inspecciona los archivos, formula preguntas, vuelve al usuario si tiene alguna, y retorna un artefacto de exploración
- **Orquestador**: pone las preguntas al usuario si explora las retorna, y re-dispara el mismo explorador con las respuestas
- **Usuario**: responde, o cancela el change
- **Intake**: NUNCA formula discovery

## Precondiciones

- El change tiene archivos identificados para leer
- El orquestador puede re-despachar la misma fase con argumentos nuevos (las respuestas del usuario)

## Flujo principal

1. `sdd-explore` inspecciona el código del cambio
2. Formula 0-N preguntas de descubrimiento ancladas a lo que encontró
3. Si hay preguntas: retorna `needs-input`, el orquestador pregunta, re-despacha con respuestas
4. Si no hay preguntas: retorna `done` con el artefacto de exploración
5. Las respuestas viajan en un bloque `### Discovery answers` del artefacto de exploración

## Ramas / flujos alternativos

- **Exploración sin preguntas** → retorna el artefacto de exploración con `done`; no interrumpe al usuario
- **Usuario rechaza responder** → le pide que confirme lo que entiende; si el usuario no entiende más que antes, queda `blocked` y deja la pregunta sin respuesta (no inventa una)
- **Pregunta sobre el pedido mismo** (ambigüedad del request) → es una pregunta de discovery que también lleva anchor (el brief que intake produjo)

## Casos borde

- **Petición tan vaga que ni ver código la aclara** → sigue siendo `needs-input` pero la pregunta ahora está anclada en lo que explora leyó
- **Usuario pide cambio trivial** → exploración aún corre (parte de `full`), podría no tener preguntas, retorna `done`
- **Las preguntas del usuario contradicen lo que explora leyó** → explora reporta la contradicción como hallazgo, no como pregunta

## Reglas de negocio

- Discovery NUNCA ocurre antes de que se lea el código
- Las preguntas SIEMPRE están ancladas a algo (código leído o intención del pedido)
- La forma es un ciclo de dos pasos: Pass 1 (lee + formula questions) retorna `needs-input` con un bloque discovery form; Pass 2 (con respuestas) retorna `done` con artefacto de exploración
- El mismo explorador se re-despacha, no un second agent or re-invoked intake
- Intake NUNCA toca discovery
- Ninguna otra fase corre hasta que la forma de discovery esté resuelta CON el usuario (no a priori, no por defecto)

## Entidades y estados

- **Discovery question** — una pregunta anclada a código o a intención; estados: formulada → respondida (o sin respuesta si usuario rechaza)
- **Discovery form** — el bloque de preguntas sin responder; estados: no existe → mostrado al usuario → respondido
- **Discovery answers** — las respuestas verbatim del usuario; estados: ausentes → proporcionadas → viajando en el artefacto de exploración

## Errores de cara al actor

- **Pregunta ambigua**: explora DEBE poder explicarla en una línea, nunca como jerga
- **Usuario no puede contestar**: la pregunta sigue abierta, o explora propone una interpretación alternativa
- **Hallazgo incompatible con el pedido**: reportado como hallazgo/blocker, no como una pregunta a responder

## Requisitos

### Requisito: Discovery is formulated and resolved inside the exploration phase

In the development domain, the exploration phase MUST always run, MUST read the code before it asks anything, and MUST own the discovery cycle: when it has questions it formulates the form and returns `needs-input`; the orchestrator puts the questions to the user and re-dispatches the same phase with the answers verbatim. When it has **no** questions it MUST return `done` with its exploration artifact and MUST NOT interrupt the user to confirm a no-questions reading — the brief-confirmation gate has already put the request's reading to the user, once, before this phase ran. The development intake phase MUST NOT formulate a discovery form. This placement is this domain's own; another domain's placement of its discovery MUST NOT change.

**Accepted cost, recorded and not mitigated:** nothing now checks what the exploration phase understood *after reading the code*. The brief-confirmation gate covers the request, not the code reading.

#### Scenario: Questions come out after the code was read

- GIVEN a change whose exploration phase has inspected the affected code
- WHEN it formulates its discovery questions
- THEN each question is grounded in something it found or in the request's own intent
- AND no discovery question was asked before any code was read

#### Scenario: The two-pass cycle belongs to exploration

- GIVEN an exploration run that has questions
- WHEN it returns
- THEN its status is `needs-input` and its return carries the discovery form beside its exploration block
- AND the same phase is re-dispatched with the answers to produce the exploration artifact

#### Scenario: An empty question list still goes to the user

- GIVEN an exploration run that formulated no questions
- WHEN it returns
- THEN its status is `done`, and this phase puts no no-questions reading to the user
- AND what the user confirms about the request reached them earlier, once, at the brief-confirmation gate

#### Scenario: Nobody answers the form on the user's behalf

- GIVEN a headless phase holding an unanswered discovery question
- WHEN it cannot reach the user
- THEN it hands the question back unanswered
- AND it never supplies an answer of its own, whatever the execution mode

#### Scenario: The form is resolved before scope is fixed

- GIVEN a change whose discovery form is still open
- WHEN the phase that fixes the change's scope would be dispatched
- THEN it is not dispatched until the form is resolved with the user

#### Scenario: Another domain keeps its discovery where it is

- GIVEN a domain whose intake phase asks its own discovery questions
- WHEN this change ships
- THEN that phase still asks them, and the relocated invariant is satisfied by it being resolved with the user before scope is fixed

### Requisito: A discovery question carries an anchor like any other gated item

Every discovery question MUST carry an anchor under the ordinary anchor criterion. No exception from the anchor requirement MUST exist for the discovery gate, and none MUST replace it. A question grounded in something exploration found anchors to a repo path with an optional start line; a question about the request's own intent anchors to the intake brief's artifact key, under the existing not-yet-written-target clause.

#### Scenario: A code-grounded question points at the code

- GIVEN a discovery question raised by something exploration read
- WHEN it is presented
- THEN it carries the repo path, and a start line when the source is a specific place

#### Scenario: An intent question points at the brief

- GIVEN a discovery question about what the request meant rather than about the code
- WHEN it is presented
- THEN it carries the intake brief's artifact key as its anchor

#### Scenario: No anchor exception is declared anywhere

- GIVEN the shared presentation contract after the change
- WHEN it is read
- THEN it declares no exception to the anchor requirement, and every moment it governs requires a real anchor

### Requisito: The exploration phase can read the brief itself

The exploration phase's agent MUST be granted the two artifact-store read capabilities the artifact retrieval protocol requires — search and full retrieval — and MUST NOT be granted any additional write capability. It MUST NOT depend on the brief being restated in its dispatch prompt.

#### Scenario: A fresh dispatch reads the brief on its own

- GIVEN an exploration dispatch whose prompt does not restate the brief
- WHEN the phase runs
- THEN it retrieves the brief from the artifact store and works from it

#### Scenario: The grant adds no write capability

- GIVEN the exploration agent's capability grant after the change
- WHEN it is compared with the one before
- THEN it gained exactly the two read capabilities, and its only write remains persisting its own artifact
