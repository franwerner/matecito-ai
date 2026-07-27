# Capability — Enviar el artefacto de una fase

- **Status:** Accepted
- **Date:** 2026-07-24

## Propósito

Recibir el artefacto estructurado que produce una fase del flujo SDD, validarlo contra el contrato de esa fase, persistirlo en el event-log y asociarlo al (proyecto, change) correcto, para que la salida de cada fase quede estructurada, consultable y disponible para la UI en vivo.

## Actores

- **Claude Code** — el agente de fase (o el orquestador) que envía el artefacto vía la tool `submit_<fase>` del MCP (una por fase: intake, propuesta, spec, diseño, tareas, apply, verify, archivo).

## Precondiciones

- El proyecto está registrado.
- Hay un change activo para el contexto (ya corrió `start_change` / la rama existe).
- El artefacto corresponde a una fase válida y trae la forma que esa fase espera.
- El change target no está `closed`.

## Flujo principal

1. El agente envía el artefacto de la fase mediante la tool `submit_<fase>`.
2. El sistema identifica el (proyecto, change) a partir del contexto de la request (gobernado por la regla de scoping de eventos).
3. Valida el artefacto contra el contrato de su fase.
4. Persiste el artefacto como un evento en el event-log del change (append-only).
5. Renderiza el artefacto a su forma canónica y se la devuelve al agente.
6. El evento queda disponible para la UI en vivo.

## Ramas / flujos alternativos

- **La validación del artefacto falla** → se rechaza con `validation_failed` y no se persiste nada.
- **El proyecto no está registrado** → se rechaza con `project_not_registered`.
- **No hay change iniciado para el contexto** → se rechaza con `change_not_started`.
- **El change target está `closed`** → se rechaza con `change_closed` (guard del ciclo de vida del change); el actor le avisa al usuario y le pregunta si lo reactiva.

## Casos borde

- **Reintento o envío duplicado del mismo artefacto** → la operación es idempotente: produce el mismo efecto y no duplica el evento en el log.
- **Fase enviada fuera de orden** → se acepta igual. El broker NO valida el orden de las fases; el log es append-only y el orden lo gobierna el orquestador, no la superficie de escritura.

## Reglas de negocio

- El event-log es append-only: un envío válido agrega un evento, nunca reescribe ni borra los anteriores.
- El broker no impone orden entre fases; acepta cualquier fase válida sobre un change activo.
- Un envío rechazado por validación no deja rastro en el event-log (todo-o-nada respecto de la persistencia).
- Un evento que referencia un EDR/spec lo ancla a su versión vigente al momento de persistir; esa versión queda congelada (inmutable), de modo que cambios posteriores del `.md` no la alteren. El agente no manda la versión: se resuelve la actual de forma implícita.

## Errores de cara al actor

- **Artefacto inválido para su fase** → `validation_failed`.
- **Proyecto no registrado** → `project_not_registered`.
- **Change no iniciado** → `change_not_started`.
- **Change cerrado** → `change_closed`.
- La forma de error de cara al actor es `{error, code}`, sin stack ni internals; los detalles internos van al log, no a la respuesta. Un envío exitoso devuelve el artefacto renderizado a su forma canónica.

## Escenarios

### Scenario: artefacto válido persistido y disponible para la UI

- **GIVEN** un proyecto registrado y un change activo, y un artefacto válido para su fase
- **WHEN** el agente lo envía vía `submit_<fase>`
- **THEN** el sistema lo valida, lo persiste como evento en el event-log del change, devuelve el artefacto renderizado al agente y el evento queda disponible para la UI en vivo

### Scenario: artefacto inválido no persiste nada

- **GIVEN** un proyecto registrado y un change activo, y un artefacto que no cumple el contrato de su fase
- **WHEN** el agente lo envía
- **THEN** el sistema responde `validation_failed` y no persiste ningún evento

### Scenario: proyecto no registrado

- **GIVEN** un contexto cuyo proyecto no fue registrado previamente
- **WHEN** el agente envía un artefacto
- **THEN** el sistema responde `project_not_registered` y no persiste nada

### Scenario: sin change iniciado

- **GIVEN** un proyecto registrado pero sin un change iniciado para el contexto
- **WHEN** el agente envía un artefacto
- **THEN** el sistema responde `change_not_started` y no persiste nada

### Scenario: reintento idempotente

- **GIVEN** un artefacto que ya fue enviado y persistido con éxito
- **WHEN** el agente reenvía el mismo artefacto
- **THEN** el efecto es el mismo y el evento no se duplica en el event-log

### Scenario: fase fuera de orden aceptada

- **GIVEN** un change activo cuyo orden de fases esperado aún no llegó a esta fase
- **WHEN** el agente envía el artefacto de una fase posterior
- **THEN** el sistema lo acepta y lo persiste igual, sin validar el orden

### Scenario: error sin internals

- **GIVEN** cualquier envío que falla
- **WHEN** el sistema responde el error al actor
- **THEN** la respuesta es de la forma `{error, code}`, sin stack ni detalles internos

### Scenario: change cerrado

- **GIVEN** un change en estado `closed`
- **WHEN** el agente envía un artefacto sobre ese change
- **THEN** el sistema responde `change_closed`, no persiste nada, y el actor le avisa al usuario y le pregunta si desea reactivar el change

### Scenario: pin de la versión del EDR/spec referenciado

- **GIVEN** un artefacto que referencia un EDR/spec en su versión vigente
- **WHEN** se persiste el evento
- **THEN** el evento queda anclado a esa versión y esa versión se congela; un cambio posterior del `.md` crea una versión nueva sin tocar la pineada

## Referencias

- **EDR** → [`../../../apps/api/.matecito-ai/edr/contracts/mcp-server.md`](../../../apps/api/.matecito-ai/edr/contracts/mcp-server.md) — la superficie MCP, las tools de submit por fase y la identidad por request.
- **EDR** → [`../../../apps/api/.matecito-ai/edr/contracts/api-contract.md`](../../../apps/api/.matecito-ai/edr/contracts/api-contract.md) — separa esta superficie de escritura (MCP) de la de lectura de la UI.
- **EDR** → [`../../../apps/api/.matecito-ai/edr/data/data-modeling.md`](../../../apps/api/.matecito-ai/edr/data/data-modeling.md) — cómo se modela y persiste el evento en el event-log.
- **EDR** → [`../../../apps/api/.matecito-ai/edr/runtime/error-handling.md`](../../../apps/api/.matecito-ai/edr/runtime/error-handling.md) — por qué la respuesta de error es `{error, code}` sin internals y dónde van los detalles internos.
- **Lifecycle** → [`../lifecycle/change.md`](../lifecycle/change.md) — el guard `change_closed`.
- **Rule** → [`../rule/event-scoping.md`](../rule/event-scoping.md) — cómo se resuelve el (proyecto, change) del contexto.
- **Process** → [`../process/index-decision-records.md`](../process/index-decision-records.md) — el indexado y versionado del EDR/spec, y cómo se resuelve y congela la versión vigente que este flujo pinea.
