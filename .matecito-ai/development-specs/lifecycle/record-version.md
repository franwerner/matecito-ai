# Capability — Ciclo de vida de la versión de un record

- **Status:** Accepted
- **Date:** 2026-07-28
- **Components:** api

## Propósito

Definir el ciclo de vida de una versión de record (EDR/spec) —sus estados y qué la mueve entre ellos— para que las capacidades que la tocan (el pin al persistir un evento, el indexado ante un cambio del `.md`) compartan una sola definición de cuándo una versión se congela y qué garantiza una versión congelada.

## Entidades y estados

- **Versión de un record** — una foto del contenido y la proyección de un `.md` de record (EDR/spec) en un momento dado. Estados:
  - **`actual`** — la versión vigente del record en su rama; es la que resuelve cualquier referencia nueva y la que se actualiza en el lugar cuando el `.md` cambia sin estar congelada.
  - **`congelada`** — inmutable para siempre; una versión congelada nunca vuelve a `actual` ni se modifica.

  Transición: `actual → congelada`, disparada por **la primera vez que un evento la referencia** (el pin). No hay transición inversa.

## Reglas de negocio

- La referencia de un evento resuelve siempre a la **versión actual** del record, de forma implícita al momento de persistir — el agente nunca manda una versión. Al quedar referenciada, esa versión se **congela**.
- Al cambiar el `.md`, hay **una sola razón** por la que la versión actual no se puede tocar: que ya esté **congelada** (referenciada por al menos un evento). Si lo está, se hace **fork**: se crea una **nueva versión actual** desde el `.md` cambiado, para la rama que cambió, y la congelada queda **intacta** para los eventos que la referencian. Si no lo está, se hace **update-in-place**: se actualiza el contenido de esa misma versión, sin bump.
- **Las versiones no se comparten entre ramas.** Cada `(record, rama)` tiene su propia versión, aun cuando dos ramas tengan contenido idéntico: lo que se comparte y se deduplica es el **contenido**, no la versión. Por eso una edición en una rama nunca puede alterar lo que resuelve otra, y ninguna rama depende del estado de versión de otra.
- Las versiones **congeladas sobreviven a todo**: a cambios posteriores del `.md` y a su borrado; los eventos que las referencian las siguen resolviendo siempre.
- El **contenido** de las versiones es **content-addressable y compartido entre ramas por hash** (deduplicado); la **versión** en cambio es propia de cada rama. El estado del record (`active`/`deleted`) se scopea por rama, y una versión congelada no pertenece a ninguna: sobrevive por sí misma para los eventos que la referencian.

## Escenarios

### Scenario: el pin congela la versión actual

- **GIVEN** un record cuya versión actual no está referenciada por ningún evento
- **WHEN** un evento la referencia al persistirse
- **THEN** esa versión pasa a `congelada` y queda inmutable; ninguna edición posterior la altera

### Scenario: update-in-place sin referencia

- **GIVEN** un record cuya versión actual no está referenciada por ningún evento
- **WHEN** su `.md` cambia
- **THEN** el contenido de esa misma versión se actualiza en el lugar, sin crear una versión nueva

### Scenario: fork al cambiar una versión congelada

- **GIVEN** un record cuya versión actual está congelada (referenciada por al menos un evento)
- **WHEN** su `.md` cambia
- **THEN** se crea una nueva versión `actual` desde el `.md` cambiado, la congelada queda intacta y las referencias existentes siguen resolviendo a la congelada

### Scenario: dos ramas con contenido idéntico no comparten versión

- **GIVEN** un record cuyo `.md` es idéntico en dos ramas
- **WHEN** el `.md` cambia en una de las dos
- **THEN** cada rama tiene su propia versión `actual` apuntando al mismo contenido deduplicado, el cambio afecta solo a la versión de la rama que cambió, y la otra sigue resolviendo la suya con su contenido intacto

### Scenario: la congelada sobrevive al borrado

- **GIVEN** una versión congelada de un record referenciada por un evento
- **WHEN** el `.md` del record se borra
- **THEN** la versión congelada sobrevive y el evento la sigue resolviendo

### Scenario: sin vuelta atrás

- **GIVEN** una versión en estado `congelada`
- **WHEN** ocurre cualquier operación posterior (ediciones, borrados, nuevos eventos)
- **THEN** la versión nunca vuelve a `actual` ni cambia su contenido

## Referencias

- **Flow** → [`../flow/submit-phase-artifact.md`](../flow/submit-phase-artifact.md) — el flujo cuyo pin dispara la transición `actual → congelada`.
- **Process** → [`../process/index-decision-records.md`](../process/index-decision-records.md) — el proceso que aplica update-in-place o fork según el estado de la versión, al (re)indexar un `.md`.
- **EDR** → [`../../../apps/api/.matecito-ai/edr/data/storage-sync-model.md`](../../../apps/api/.matecito-ai/edr/data/storage-sync-model.md) — el modelo de versionado content-addressable en la base y su dedup por hash.
- **EDR** → [`../../../apps/api/.matecito-ai/edr/data/data-modeling.md`](../../../apps/api/.matecito-ai/edr/data/data-modeling.md) — el store de versiones y el pin.
