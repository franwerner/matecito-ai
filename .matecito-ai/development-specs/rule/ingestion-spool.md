# Capability — Spool de ingesta offline

- **Status:** Accepted
- **Date:** 2026-07-27

## Propósito

Garantizar que nada de lo que un productor local ingesta al broker se pierda ni bloquee la sesión cuando el broker no está disponible: ante falla se guarda localmente completo y se reconcilia al volver el contacto. Es la regla transversal de resiliencia de toda la ingesta.

## Reglas de negocio

- **Aplica a todo productor local que ingesta contenido al broker** — hoy: los eventos mecánicos de sesión y las fotos de código; cualquier productor futuro adopta la misma regla.
- **La ingesta jamás bloquea ni rompe la sesión del coder**: reportar es fire-and-forget; ninguna falla del broker (caído, lento, rechazando) produce un error visible ni una espera del lado de la sesión.
- **Ante falla de contacto, el productor materializa el ítem COMPLETO en un spool local** — incluido el contenido (por ejemplo, la foto de un archivo el productor la lee y la guarda él mismo) y el timestamp original. Nada queda "por referencia" que pueda perderse mientras el broker no está.
- **En el próximo contacto exitoso, el spool se descarga ANTES que el ítem nuevo**: primero lo pendiente en orden, después lo actual. El orden de reconciliación preserva el orden original.
- **Los ítems viajan con su timestamp original**: un ítem reconciliado tarde se ordena por cuándo ocurrió, no por cuándo llegó.
- **La reconciliación es idempotente**: un flush repetido (o un ítem reportado dos veces) no duplica nada del lado del broker.

## Escenarios

### Scenario: falla de contacto no pierde el ítem

- **GIVEN** el broker caído en el momento en que un productor debe ingestar un ítem con contenido
- **WHEN** el intento de contacto falla
- **THEN** el productor guarda el ítem completo (contenido + timestamp original) en su spool local y la sesión continúa sin error ni espera

### Scenario: reconciliación en orden, lo pendiente primero

- **GIVEN** un spool con ítems pendientes y el broker nuevamente disponible
- **WHEN** el productor logra el próximo contacto
- **THEN** descarga primero el spool en orden y después el ítem nuevo; el broker los ordena por su timestamp original

### Scenario: flush repetido no duplica

- **GIVEN** un ítem del spool ya reconciliado con éxito
- **WHEN** el mismo ítem se reporta de nuevo
- **THEN** el broker no lo duplica; el efecto es el mismo

### Scenario: la sesión nunca se bloquea

- **GIVEN** cualquier estado del broker (caído, lento, rechazando)
- **WHEN** un productor ingesta
- **THEN** la sesión del coder no se bloquea ni ve un error, con o sin spool de por medio

## Referencias

- **Process** → [`../process/ingest-mechanical-events.md`](../process/ingest-mechanical-events.md) — productor: los eventos mecánicos de sesión.
- **Process** → [`../process/capture-code-snapshots.md`](../process/capture-code-snapshots.md) — productor: las fotos de código (el caso que exige spool con contenido).
- **EDR** → [`../../../apps/api/.matecito-ai/edr/contracts/api-contract.md`](../../../apps/api/.matecito-ai/edr/contracts/api-contract.md) — la ingesta HTTP del broker como superficie separada.
