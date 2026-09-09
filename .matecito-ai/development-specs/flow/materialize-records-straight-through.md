# Capability — Materializar un record de decisión directo, en el paso que lo implementa

- **Status:** Accepted
- **Date:** 2026-09-09
- **Components:** cli

## Propósito

Escribir un record de decisión durable sin canal de forwarding, sin ledger por-cambio y sin confirmación intermedia: `sdd-apply` lee `## New Decisions` del artefacto de diseño que ya lee como input ordinario, y lo escribe como record `Accepted` en el mismo despacho que implementa el trabajo que ese record rige.

## Actores

- **`sdd-apply`**: lee `## New Decisions` del artefacto `sdd/{change-name}/design`; para cada entrada cuya tarea gobernada implementa en este despacho, rama por `record-mode` y escribe el record
- **`sdd-design`**: produce `### New Decisions`, declarando `gates: reported` — nunca abre un gate de confirmación para esta sección
- **Orquestador**: no interviene en este camino — no arma un bloque de decisiones ratificadas para el prompt de despacho de `sdd-apply`, no lee ni escribe un ledger por-cambio

## Precondiciones

- El artefacto de diseño existe en Engram bajo `sdd/{change-name}/design` y `sdd-apply` lo lee como parte de su lectura ordinaria de upstream
- Cada entrada de `## New Decisions` lleva `record:`, `record-mode:` (`create` | `modify`) y `blocking-test:`

## Flujo principal

1. `sdd-apply` lee el artefacto de diseño completo, incluyendo `## New Decisions`
2. Para cada entrada cuya tarea gobernada esta corrida implementa, rama por `record-mode`: `create` renderiza el record completo y agrega filas de INDEX; `modify` edita solo las cláusulas nombradas y no agrega fila de INDEX
3. El record se escribe con `Status: Accepted`, incondicionalmente
4. La escritura ocurre en el mismo paso que implementa el trabajo que el record rige — nunca en un paso separado ni diferido

## Ramas / flujos alternativos

- **Bajo un fan-out (batch paralelo)** → solo el rol único que la fragmento de dominio nombra como writer materializa los records; los demás workers no escriben ninguno, así la misma lista `## New Decisions` produce un solo conjunto de archivos, no uno por worker. Una corrida aislada lleva las filas de INDEX sin aplicar en su Task Run Report; la corrida de consolidación las aplica una sola vez
- **Una entrada cuya tarea gobernada este despacho no alcanza** → no se materializa en este paso; queda para el despacho que sí implemente esa tarea

## Casos borde

- **Una declaración de `record-mode` que no coincide con lo que hay en disco** → sigue siendo una falla, nunca un cambio de rama silencioso
- **Ningún gate abrió para la entrada** → no cambia nada de este flujo: la entrada se materializa igual, porque nunca hubo un gate del que depender

## Reglas de negocio

- No existe canal de forwarding: nada arma un bloque de "decisiones ratificadas" para el prompt de despacho de `sdd-apply`
- No existe ledger por-cambio (`sdd/{change-name}/ratified-decisions`): no se lee, no se escribe, no se consulta para re-emergence
- El rechazo no tiene canal: una proposal de `## New Decisions` no es rechazable en este camino — no hay `### Rejected Proposals Checked`, no hay token `design-conflict`
- La rama `create` / `modify` de `record-mode` sobrevive sin cambios: `create` renderiza, escribe y aplica filas de INDEX; `modify` edita solo las cláusulas nombradas y no agrega fila de INDEX
- Todo record escrito por este camino es `Accepted` incondicionalmente — el status no depende de si un gate abrió

## Entidades y estados

- **Entrada de `## New Decisions`** — una proposal de decisión con `record`, `record-mode`, `blocking-test` y su propio `summary`/`rationale`. Estados: reportada (en el artefacto de diseño) → materializada (record `Accepted` en disco)

## Errores de cara al actor

- **Una tarea gobernada por una entrada que el prompt de despacho no menciona**: `sdd-apply` retorna `blocked` nombrando el item y la resolución faltante

## Requisitos

### Requisito: El trigger lee el artefacto de diseño y nada más

El Step 4b de `sdd-apply` DEBE nombrar la clave del artefacto de diseño como la fuente de `## New Decisions`, y DEBE NO nombrar ningún bloque del prompt de despacho, ningún ledger de decisiones ratificadas ni ningún resultado de confirmación como precondición.

#### Scenario: El trigger nombra el artefacto, no el prompt

- **GIVEN** el Step 4b de `sdd-apply` después del cambio
- **WHEN** se lee su condición de disparo
- **THEN** nombra la clave del artefacto de diseño como fuente de `## New Decisions`
- **AND** no nombra ningún bloque del prompt de despacho, ningún ledger de decisiones ratificadas ni ningún resultado de confirmación como precondición

### Requisito: El canal de forwarding desapareció de las instrucciones del orquestador

`payload/domains/development/CLAUDE.md` DEBE NO contener el párrafo que arma un bloque de decisiones ratificadas para el prompt de despacho de `sdd-apply`, ni instrucción alguna de leer o escribir `sdd/{change-name}/ratified-decisions`.

#### Scenario: El párrafo de forwarding no existe

- **GIVEN** `payload/domains/development/CLAUDE.md` después del cambio
- **WHEN** se busca el párrafo que arma un bloque de decisiones ratificadas para el prompt de `sdd-apply`
- **THEN** no existe tal párrafo
- **AND** ningún texto instruye al orquestador a escribir o leer `sdd/{change-name}/ratified-decisions`

### Requisito: Ambas ramas de `record-mode` sobreviven

#### Scenario: Ambas ramas siguen andando

- **GIVEN** el Step 4b después del cambio
- **WHEN** se lee la rama `record-mode`
- **THEN** `create` sigue renderizando, escribiendo y aplicando filas de INDEX, y `modify` sigue editando solo las cláusulas nombradas y sin agregar fila de INDEX
- **AND** una declaración que no coincide con lo que hay en disco sigue siendo una falla, nunca un cambio de rama silencioso

### Requisito: Todo record escrito por este camino es `Accepted`

#### Scenario: El status no depende del gate

- **GIVEN** el Step 4b después del cambio
- **WHEN** se lee el status que escribe
- **THEN** es `Accepted` incondicionalmente
- **AND** ningún status depende de si un gate abrió

### Requisito: Un solo writer materializa, incluso bajo un fan-out

#### Scenario: Un solo conjunto de archivos, no uno por worker

- **GIVEN** un batch de implementación paralelo donde varios workers leen el mismo artefacto de diseño
- **WHEN** los records se materializan
- **THEN** solo el rol único de despacho que el fragmento de dominio nombra como writer los materializa
- **AND** los demás workers no escriben ningún record, así la misma lista `## New Decisions` produce un solo conjunto de archivos, no uno por worker

## Escenarios

El spec no declara escenarios adicionales más allá de los definidos en sus cinco requisitos anteriormente.

## Referencias

- **Contrato de fase** → [`../../../payload/domains/development/references/phase-returns/sdd-design/sdd-design.md`](../../../payload/domains/development/references/phase-returns/sdd-design/sdd-design.md) — `### New Decisions` declara `gates: reported`, sin token `contested`
- **Rule** → [`../rule/gate-firing-triggers.md`](../rule/gate-firing-triggers.md) — cuándo un gate dispara y por qué esta sección nunca lo hace
- **Single-writer** → `contracts/single-writer-per-batch.md`, `structure/consolidation-run-is-the-integrator.md` — el rol único que materializa bajo un fan-out
