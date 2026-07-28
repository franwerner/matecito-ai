# Capability — Ingerir eventos mecánicos de sesión

- **Status:** Accepted
- **Date:** 2026-07-27

## Propósito

Ingerir eventos mecánicos de actividad de sesión —un agente arrancó, un agente terminó— reportados por una fuente de actividad, y sumarlos al event-log del change correcto, para que la UI muestre en vivo qué está corriendo sin intervención del orquestador.

## Actores

- **Fuente de actividad** — el mecanismo del runtime del coder que reporta la actividad de forma automática y determinista. Hoy: los hooks de arranque/fin de sub-agente de Claude Code, vía un cliente mínimo que reporta a la ingesta del broker. Otros runtimes pueden alimentar el mismo contrato con su propio adapter (posibilidad futura, no especificada).

## Precondiciones

- La fuente de actividad está configurada en el runtime del coder (la gestiona el installer de matecito).
- Para que un evento se ingeste: el proyecto está registrado y hay un change **`active`** para el contexto. Si falta cualquiera de las dos —incluido un change que existe pero está `closed`—, el evento se descarta en silencio (ver casos borde).

## Flujo principal

1. El runtime del coder lanza o termina un sub-agente; la fuente de actividad lo reporta a la ingesta del broker con el contexto de la sesión (directorio del proyecto, identificación del agente, timestamp del evento).
2. El broker resuelve el (proyecto, change) destino con la regla de scoping de eventos.
3. Appendea el evento al event-log del change (append-only).
4. El evento sale empujado en vivo a la UI vía la suscripción del change.

## Casos borde

- **El broker está caído o inalcanzable** → aplica el **spool de ingesta offline** ([`../rule/ingestion-spool.md`](../rule/ingestion-spool.md)): el evento se guarda localmente con su timestamp original y se reconcilia en el próximo contacto, sin que la sesión se entere.
- **Proyecto no registrado** (la fuente es global: dispara en cualquier repo) → el broker descarta el evento en silencio; no es un error.
- **Proyecto registrado pero sin change activo para el contexto** (trabajo suelto, sin flow SDD andando) → el broker descarta el evento en silencio; solo se ingesta actividad de un change activo.
- **El change del contexto existe pero está `closed`** → cuenta como "sin change activo": descarte silencioso. Un change cerrado no acumula más actividad mecánica; el canal deliberado (las tools) es el que rechaza con `change_closed` y dispara la pregunta de reapertura — una vez reabierto, la actividad mecánica vuelve a entrar sola.
- **Stop sin start previo, o start duplicado** (reportes que se pierden o repiten) → ingesta tolerante: el broker appendea lo que llega sin exigir pareja; la UI interpreta con lo que hay.

## Reglas de negocio

- La resiliencia de la ingesta —fire-and-forget que jamás afecta la sesión, spool local ante falla, timestamps originales, reconciliación idempotente y en orden— la gobierna el **spool de ingesta offline** ([`../rule/ingestion-spool.md`](../rule/ingestion-spool.md)); este proceso lo aplica, no lo define.
- La captura es **automática y determinista**: la dispara el runtime, no un agente que "se acuerda" de reportar; el orquestador no interviene.
- El alcance es **solo los bordes**: arranque y fin de sub-agente. La actividad intra-agente (qué tool usa, progreso) queda fuera de esta capacidad.
- Los eventos mecánicos entran por la **ingesta HTTP del broker**, no por la superficie MCP: la fuente ejecuta comandos, no tools.

## Escenarios

### Scenario: arranque de agente visible en vivo

- **GIVEN** un proyecto registrado con un change activo y una UI suscrita a ese change
- **WHEN** el runtime lanza un sub-agente y la fuente de actividad lo reporta
- **THEN** el evento de arranque queda appendeado al event-log del change y la UI lo recibe empujado en vivo

### Scenario: fin de agente cierra la actividad

- **GIVEN** un agente cuyo arranque fue ingestado
- **WHEN** el agente termina y la fuente reporta el fin
- **THEN** el evento de fin queda en el event-log y la UI puede mostrar la fase como terminada, con su duración

### Scenario: proyecto no registrado se descarta

- **GIVEN** un repo que no está registrado en el broker
- **WHEN** llega un evento mecánico originado en ese repo
- **THEN** el broker lo descarta en silencio, sin error y sin persistir nada

### Scenario: sin change activo se descarta

- **GIVEN** un proyecto registrado pero sin change activo para el contexto
- **WHEN** llega un evento mecánico de ese contexto
- **THEN** el broker lo descarta en silencio; el event-log no registra actividad suelta

### Scenario: change cerrado se descarta

- **GIVEN** un proyecto registrado cuya rama mapea a un change en estado `closed`
- **WHEN** llega un evento mecánico de ese contexto
- **THEN** el broker lo descarta en silencio, igual que si no hubiera change; si el change se reabre, la actividad mecánica posterior vuelve a ingestarse

### Scenario: stop huérfano tolerado

- **GIVEN** un change activo cuyo event-log no tiene el arranque de un agente
- **WHEN** llega el evento de fin de ese agente
- **THEN** el broker lo appendea igual, sin exigir el arranque previo

## Referencias

- **EDR** → [`../../../apps/api/.matecito-ai/edr/contracts/api-contract.md`](../../../apps/api/.matecito-ai/edr/contracts/api-contract.md) — la ingesta HTTP como superficie separada de la MCP (escritura deliberada) y de la lectura de la UI.
- **EDR** → [`../../../apps/api/.matecito-ai/edr/data/data-modeling.md`](../../../apps/api/.matecito-ai/edr/data/data-modeling.md) — cómo se modela el evento mecánico en el event-log.
- **Rule** → [`../rule/ingestion-spool.md`](../rule/ingestion-spool.md) — la resiliencia de la ingesta: fire-and-forget, spool local y reconciliación.
- **Rule** → [`../rule/event-scoping.md`](../rule/event-scoping.md) — cómo se resuelve el (proyecto, change) del contexto reportado.
- **Flow** → [`../flow/serve-change-state.md`](../flow/serve-change-state.md) — la suscripción por la que estos eventos llegan en vivo a la UI.
- **Lifecycle** → [`../lifecycle/change.md`](../lifecycle/change.md) — los estados `active`/`closed` que gobiernan el descarte, y la reapertura por la vía deliberada.
