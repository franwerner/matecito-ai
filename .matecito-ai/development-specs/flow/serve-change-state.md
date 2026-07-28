# Capability — Servir el estado de un change a la UI

- **Status:** Accepted
- **Date:** 2026-07-27

## Propósito

Servirle a la UI el estado de un change —snapshot inicial del event-log y sus artefactos, y push en vivo de los eventos nuevos— para que el cockpit muestre el flujo SDD en tiempo real sin polling.

## Actores

- **La UI (el cockpit)** — abre un change para verlo, en nombre del usuario.

## Precondiciones

- El proyecto está registrado.
- El change existe (en cualquier estado: `active` o `closed` — un change cerrado también se sirve, para historial).
- No hay permisos ni auth: el daemon es local (localhost, single-user); cualquier change de cualquier proyecto registrado es visible.

## Flujo principal

1. La UI pide el snapshot de un change.
2. El sistema devuelve el estado del change (activo/cerrado, metadata) y el listado ordenado de su event-log —cada evento con su forma según su clase: los de fase con su artefacto en forma canónica; los mecánicos (arranque/fin de agente) con su tipo, agente y timestamp, sin artefacto— junto con la posición del log hasta la que cubre el snapshot.
3. La UI se suscribe por WebSocket al change, declarando la posición que trajo el snapshot.
4. El sistema empuja por la suscripción todo evento del log posterior a esa posición (incluidos los que ocurrieron entre el snapshot y la suscripción), y de ahí en adelante cada evento nuevo que se persiste.
5. La UI aplica los eventos incrementales sobre el snapshot y queda al día.

## Ramas / flujos alternativos

- **La conexión WS se corta** → la UI se re-suscribe declarando la última posición que aplicó; el sistema le empuja lo que ocurrió en el medio, sin necesidad de re-pedir el snapshot.
- **El change está `closed`** → se sirve igual (snapshot + suscripción); si el change se reabre y llegan eventos nuevos, la suscripción los recibe.

## Casos borde

- **Ventana entre snapshot y suscripción** → no se pierden eventos: la suscripción arranca desde la posición declarada, no desde "ahora".
- **Varias suscripciones al mismo change** (dos ventanas, dos vistas) → cada una recibe todos los eventos; el broker no limita a una conexión.
- **Posición de resume inválida** (una posición que el log no reconoce) → el sistema la rechaza indicándolo y la UI debe re-pedir el snapshot; nunca se sirve un stream con huecos silenciosos.

## Reglas de negocio

- La suscripción es **por change**: la UI se suscribe al change que está mirando y se desuscribe al salir; el broker no empuja eventos de otros changes por esa suscripción.
- El resume es por **posición del log**: tanto la suscripción inicial (posición del snapshot) como la reconexión (última posición aplicada) garantizan entrega sin huecos ni duplicados desde la posición declarada.
- La lectura no muta nada: servir snapshot o suscripción no agrega eventos al log ni cambia el estado del change.
- Snapshot y push sirven **todas las clases de eventos** del log (de fase y mecánicos), cada uno con su forma; la UI distingue por la clase del evento.

## Errores de cara al actor

- **Change inexistente** → `change_not_found`.
- **Proyecto no registrado** → `project_not_registered`.
- **Posición de resume inválida** → `invalid_position`; la UI re-pide el snapshot.
- La forma de error es `{error, code}`, sin stack ni internals, como en el resto del sistema.

## Escenarios

### Scenario: snapshot completo al abrir

- **GIVEN** un proyecto registrado con un change que tiene eventos persistidos
- **WHEN** la UI pide el snapshot del change
- **THEN** recibe el estado del change y el listado ordenado de eventos con sus artefactos en forma canónica, junto con la posición del log cubierta

### Scenario: push en vivo tras suscribirse

- **GIVEN** una UI suscrita a un change con la posición del snapshot declarada
- **WHEN** se persiste un evento nuevo en el event-log de ese change
- **THEN** el evento se empuja por la suscripción y la UI queda al día sin polling

### Scenario: sin huecos entre snapshot y suscripción

- **GIVEN** un evento persistido después de que el sistema armó el snapshot pero antes de que la suscripción quedó activa
- **WHEN** la UI se suscribe declarando la posición del snapshot
- **THEN** el sistema le empuja ese evento intermedio; la UI no se pierde nada

### Scenario: reconexión sin re-snapshot

- **GIVEN** una suscripción que se cortó, con la UI habiendo aplicado hasta una posición del log
- **WHEN** la UI se re-suscribe declarando esa posición
- **THEN** el sistema empuja los eventos posteriores ocurridos durante el corte y el stream continúa sin huecos ni re-pedir snapshot

### Scenario: change cerrado se sirve para historial

- **GIVEN** un change en estado `closed`
- **WHEN** la UI pide su snapshot y se suscribe
- **THEN** el sistema sirve el snapshot completo y acepta la suscripción; si el change se reabre y llegan eventos, la suscripción los recibe

### Scenario: varias vistas del mismo change

- **GIVEN** dos suscripciones activas al mismo change
- **WHEN** se persiste un evento nuevo
- **THEN** ambas suscripciones reciben el evento

### Scenario: change inexistente

- **GIVEN** un identificador de change que no existe
- **WHEN** la UI pide su snapshot o se suscribe
- **THEN** el sistema responde `change_not_found` con la forma `{error, code}`

### Scenario: posición de resume inválida

- **GIVEN** una UI que se suscribe declarando una posición que el log no reconoce
- **WHEN** el sistema procesa la suscripción
- **THEN** la rechaza con `invalid_position` y la UI re-pide el snapshot; no se sirve un stream con huecos silenciosos

### Scenario: proyecto no registrado

- **GIVEN** un contexto cuyo proyecto no fue registrado
- **WHEN** la UI pide un snapshot o se suscribe
- **THEN** el sistema responde `project_not_registered` con la forma `{error, code}`

### Scenario: la lectura no muta

- **GIVEN** un change con su event-log y su estado
- **WHEN** la UI pide el snapshot y se suscribe
- **THEN** el event-log y el estado del change quedan idénticos: ningún evento nuevo, ninguna transición

### Scenario: eventos mecánicos servidos por su clase

- **GIVEN** un change cuyo event-log mezcla eventos de fase y eventos mecánicos
- **WHEN** la UI recibe el snapshot o el push
- **THEN** los eventos de fase llegan con su artefacto en forma canónica y los mecánicos con su tipo, agente y timestamp, sin artefacto

### Scenario: la suscripción solo trae el change suscrito

- **GIVEN** una UI suscrita al change A y eventos persistiéndose en el change B
- **WHEN** llegan los eventos de B
- **THEN** la suscripción de A no los recibe

## Referencias

- **EDR** → [`../../../apps/api/.matecito-ai/edr/contracts/api-contract.md`](../../../apps/api/.matecito-ai/edr/contracts/api-contract.md) — la superficie de lectura de la UI (OpenAPI Go-first + tipos generados) separada de la superficie de escritura (MCP).
- **EDR** → [`../../../apps/api/.matecito-ai/edr/data/data-modeling.md`](../../../apps/api/.matecito-ai/edr/data/data-modeling.md) — cómo se modela el event-log del que se sirven snapshot y push.
- **EDR** → [`../../../apps/api/.matecito-ai/edr/delivery/deployment-topology.md`](../../../apps/api/.matecito-ai/edr/delivery/deployment-topology.md) — el daemon local único en localhost; por qué no hay permisos ni auth.
- **EDR** → [`../../../apps/api/.matecito-ai/edr/runtime/error-handling.md`](../../../apps/api/.matecito-ai/edr/runtime/error-handling.md) — la forma de error `{error, code}` sin internals.
- **Flow** → [`submit-phase-artifact.md`](submit-phase-artifact.md) — la superficie de escritura que produce los eventos que este flujo sirve.
- **Lifecycle** → [`../lifecycle/change.md`](../lifecycle/change.md) — los estados `active`/`closed` del change servido.
