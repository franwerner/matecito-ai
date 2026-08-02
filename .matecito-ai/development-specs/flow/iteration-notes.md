# Capability — Notas de iteración desde la UI

- **Status:** Accepted
- **Date:** 2026-07-27
- **Components:** api, ui

## Propósito

Permitir que el usuario marque en el cockpit notas de iteración —sobre una parte de un artefacto o del código, ancladas a líneas o generales— que el broker guarda y señala al chat en el próximo mensaje de esa conversación, para que la iteración se converse donde se trabaja: en el chat. Es la única escritura que produce la UI: anota; nunca resuelve gates ni conversa.

## Actores

- **El usuario vía la UI** — crea la nota mientras mira un artefacto o un diff en el cockpit.
- **La fuente de señal** — el mecanismo del runtime del coder que, en cada mensaje del usuario en el chat, consulta al broker e inyecta la señal de notas pendientes como contexto. Hoy: el hook de prompt de Claude Code, vía el mismo cliente mínimo de los otros canales.
- **El AI (orquestador/agente)** — recibe la señal, pregunta, y lee y resuelve las notas vía tools del MCP.

## Precondiciones

- El proyecto está registrado y el change existe.
- La fuente de señal está configurada en el runtime (la gestiona el installer de matecito).

## Flujo principal

1. El usuario, mirando el contenido de un evento en el cockpit (un artefacto de fase, el diff de un batch), crea una nota: su comentario de qué rever o rehacer, con su anclaje.
2. El broker la guarda como `pendiente`, anclada al change y, según cómo se creó, al evento y al archivo + rango de líneas.
3. En el **próximo mensaje del usuario** en una sesión de chat cuyo contexto resuelve a ese (proyecto, change), la fuente de señal inyecta **solo la señal**: cuántas notas pendientes hay — nunca su contenido.
4. El AI le pregunta al usuario si quiere que las lea e itere sobre ellas.
5. Con la confirmación del usuario, el AI lee las notas vía la tool de consulta del MCP → pasan a `entregada`.
6. El AI atiende cada nota conversando/iterando en el chat; al atenderla la marca vía la tool de resolución → pasa a `resuelta`, y el vínculo queda trazado en el event-log (qué nota disparó qué trabajo).

## Ramas / flujos alternativos

- **El usuario responde que no** (paso 4) → las notas siguen `pendientes`; el AI no las lee ni insiste en ese turno. La señal vuelve a aparecer en el próximo mensaje.
- **El AI no puede atender una nota** (no aplica, requiere decisión del usuario) → queda `entregada` sin resolver; lo conversa en el chat y la nota sigue visible en el cockpit como no resuelta.
- **Broker caído al momento de la señal** → la consulta se saltea en silencio: el mensaje del usuario llega sin señal y la sesión no se bloquea ni ve un error (las notas siguen `pendientes` para la próxima).

## Casos borde

- **Nota sin líneas** → nota general sobre el evento anclado (todo el artefacto, todo el batch).
- **Nota sin evento** → nota general del change.
- **El contenido anclado cambió después de la nota** (el archivo siguió evolucionando, el artefacto se re-submiteó) → la nota no se rompe ni se desplaza: ancla al evento, y el evento pinea la versión exacta que el usuario vio al escribirla (la versión del artefacto, las fotos del batch).
- **Varias sesiones de chat sobre el mismo change** → cualquiera de ellas recibe la señal en su próximo mensaje; la primera que lee las notas las pasa a `entregada` (las demás ya no las verán como pendientes).

## Reglas de negocio

- **Anclaje**: toda nota ancla al **change**; opcionalmente a un **evento del event-log** (el submit de fase o el batch de apply que el usuario estaba mirando); opcionalmente a **archivo + rango de líneas (from,to)** dentro del contenido de ese evento. Anclar al evento = anclar al contenido pineado por ese evento; y como cada evento porta su fase SDD, la nota hereda a qué parte del flujo refiere.
- **La señal nunca lleva contenido**: solo la existencia y cantidad de notas pendientes. El AI no lee ninguna nota sin la confirmación explícita del usuario en el chat.
- **La resolución es explícita y del AI**: solo la tool de resolución pasa una nota a `resuelta`, y eso queda trazado en el event-log. La UI no resuelve notas; el usuario las resuelve conversando en el chat.
- La creación y la resolución de una nota quedan reflejadas como eventos del change — llegan en vivo al cockpit por la suscripción, como todo evento.
- La UI solo escribe lo suyo: notas. Todo lo demás de la UI sigue siendo lectura.

## Entidades y estados

- **Nota de iteración** — el comentario del usuario con su anclaje. Estados: **`pendiente`** (creada, aún no leída por el AI) → **`entregada`** (el AI la leyó vía la tool de consulta, con confirmación del usuario) → **`resuelta`** (el AI la atendió y la marcó vía la tool de resolución). Sin vuelta atrás.

## Errores de cara al actor

- **Change inexistente al crear la nota** → `change_not_found`.
- **Proyecto no registrado** → `project_not_registered`.
- La forma de error es `{error, code}`, sin internals, como en el resto del sistema.

## Escenarios

### Scenario: nota anclada a líneas de un diff

- **GIVEN** el usuario mirando en el cockpit el diff de un batch de apply
- **WHEN** crea una nota sobre las líneas (from,to) de un archivo con su comentario
- **THEN** el broker la guarda `pendiente`, anclada al change, a ese evento de apply y a ese archivo + rango

### Scenario: señal sin contenido en el próximo mensaje

- **GIVEN** notas `pendientes` de un change y una sesión de chat cuyo contexto resuelve a ese change
- **WHEN** el usuario escribe su próximo mensaje
- **THEN** la señal inyectada indica solo cuántas notas pendientes hay, sin el contenido, y el AI le pregunta al usuario si quiere que las lea e itere

### Scenario: el AI no lee sin confirmación

- **GIVEN** la señal de notas pendientes inyectada
- **WHEN** el usuario no confirma (o responde que no)
- **THEN** el AI no lee las notas, siguen `pendientes` y la señal reaparece en el próximo mensaje

### Scenario: lectura confirmada entrega las notas

- **GIVEN** el usuario confirma que quiere iterar sobre las notas
- **WHEN** el AI las lee vía la tool de consulta
- **THEN** las notas pasan a `entregada` y el AI las atiende en la conversación

### Scenario: resolución trazada

- **GIVEN** una nota `entregada` que el AI atendió
- **WHEN** el AI la marca vía la tool de resolución
- **THEN** la nota pasa a `resuelta` y el vínculo queda en el event-log del change, visible en vivo en el cockpit

### Scenario: el anclaje sobrevive a cambios posteriores

- **GIVEN** una nota anclada a un evento cuyo archivo siguió cambiando en batches posteriores
- **WHEN** el AI (o el cockpit) resuelve el anclaje de la nota
- **THEN** la nota apunta a la versión exacta que pineaba ese evento (las fotos de ese batch), no al estado actual

### Scenario: broker caído no bloquea el chat

- **GIVEN** el broker caído al momento del mensaje del usuario
- **WHEN** la fuente de señal intenta consultar
- **THEN** se saltea en silencio: el mensaje llega sin señal, la sesión no se bloquea y las notas siguen `pendientes`

### Scenario: nota general del change

- **GIVEN** el usuario crea una nota sin evento ni líneas
- **WHEN** el AI la lee tras confirmar
- **THEN** la trata como una nota general del change, sin anclaje fino

## Referencias

- **Flow** → [`submit-phase-artifact.md`](submit-phase-artifact.md) — los eventos de fase a los que anclan las notas, y el pin de versiones que las mantiene estables.
- **Process** → [`../process/capture-code-snapshots.md`](../process/capture-code-snapshots.md) — las fotos que resuelven el anclaje de una nota sobre código.
- **Flow** → [`serve-change-state.md`](serve-change-state.md) — la suscripción por la que la creación/resolución de notas llega en vivo al cockpit.
- **Rule** → [`../rule/event-scoping.md`](../rule/event-scoping.md) — cómo la sesión de chat resuelve a qué (proyecto, change) pertenece la señal.
- **EDR** → [`../../../apps/api/.matecito-ai/edr/contracts/mcp-server.md`](../../../apps/api/.matecito-ai/edr/contracts/mcp-server.md) — la superficie MCP donde viven las tools de consulta y resolución de notas.
