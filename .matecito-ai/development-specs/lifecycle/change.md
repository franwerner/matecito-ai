# Capability — Ciclo de vida del change

- **Status:** Accepted
- **Date:** 2026-07-24

## Propósito

Definir el ciclo de vida de un change —sus estados y qué lo mueve entre ellos— para que la superficie de escritura sepa cuándo un change acepta eventos y cuándo los rechaza.

## Entidades y estados

- **`change`** — la unidad de trabajo del flujo SDD sobre una rama; contra ella se appendean los eventos de cada fase. Estados:
  - **`active`** — en progreso; se appendean eventos.
  - **`closed`** — terminado; no se esperan más eventos.

  Transiciones: `active → closed` (cerrar) y `closed → active` (reabrir), ambas SOLO vía la tool `change_status(status)`.

## Reglas de negocio

- La transición de estado ocurre únicamente vía la tool `change_status(status)`; es explícita, nunca un efecto lateral de otra tool.
- No hay auto-close: archivar el change no lo cierra por sí solo.
- No hay auto-reopen: un change `closed` solo vuelve a `active` con una llamada explícita a `change_status(active)`.
- **Guard `change_closed`:** cualquier tool que opere sobre un change `closed` (por ejemplo `start_change` o `submit_<fase>`) se rechaza con `change_closed`, informando que está cerrado. La tool `change_status` es la excepción: es la única que puede operar sobre un change `closed`, para reabrirlo.
- Ante un `change_closed`, el actor (Claude) le avisa al usuario y le pregunta si lo reactiva; la reactivación es explícita vía `change_status(active)`.

## Errores de cara al actor

- **Tool sobre un change cerrado** → `change_closed`.

## Escenarios

### Scenario: cerrar

- **GIVEN** un change en estado `active`
- **WHEN** se llama `change_status(closed)`
- **THEN** el change pasa a `closed` y deja de aceptar eventos de otras tools

### Scenario: reabrir

- **GIVEN** un change en estado `closed`
- **WHEN** se llama `change_status(active)`
- **THEN** el change pasa a `active` y vuelve a aceptar eventos

### Scenario: tool sobre change cerrado

- **GIVEN** un change en estado `closed`
- **WHEN** una tool distinta de `change_status` opera sobre ese change
- **THEN** el sistema responde `change_closed` y el actor le avisa al usuario y le pregunta si desea reactivarlo

### Scenario: reactivación tras prompt

- **GIVEN** un change `closed` sobre el que una operación fue rechazada con `change_closed`
- **WHEN** el usuario confirma la reactivación y se llama `change_status(active)`
- **THEN** el change pasa a `active` y la operación puede reintentarse

### Scenario: archive no cierra el change

- **GIVEN** un change en estado `active`
- **WHEN** se persiste el artefacto de la fase archive (submit de archive)
- **THEN** el change sigue `active` (solo `change_status(closed)` lo cierra; ninguna fase lo hace como efecto lateral)

### Scenario: no hay auto-reopen

- **GIVEN** un change en estado `closed`
- **WHEN** llega cualquier tool que no sea `change_status`
- **THEN** el change sigue `closed` y la operación se rechaza con `change_closed` (solo `change_status(active)` lo reabre)

## Referencias

- **EDR** → [`../../../apps/api/.matecito-ai/edr/delivery/deployment-topology.md`](../../../apps/api/.matecito-ai/edr/delivery/deployment-topology.md) — el daemon global único que aloja el store donde vive el estado del change.
- **EDR** → [`../../../apps/api/.matecito-ai/edr/contracts/mcp-server.md`](../../../apps/api/.matecito-ai/edr/contracts/mcp-server.md) — la tool `change_status` y el error `change_closed`.
- **Rule** → [`../rule/event-scoping.md`](../rule/event-scoping.md) — cómo se identifica el change sobre el que se aplica el guard.
