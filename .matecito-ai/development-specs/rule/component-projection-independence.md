# Capability — Independencia de proyecciones del set de componentes

- **Status:** Accepted
- **Date:** 2026-08-06
- **Components:** cli

## Propósito

El set declarado en `repo.components` es vocabulario **de repo**, compartido por varios consumidores. Esta capability fija cómo conviven esos consumidores: cada uno **proyecta** el set a su propia granularidad y con su propia vida, y ninguna proyección escribe en otra.

## Reglas de negocio

### Regla: La declaración es independiente de cualquier proyección

La definición repo-level del concepto —qué es un componente, cómo se declara el set, y el gate presence-based— MUST tener un hogar canónico propio, no subordinado a ningún consumidor. Un consumidor nuevo MUST poder leerla sin abrir la documentación del concepto de capability-spec. Todo consumidor que nombre componentes MUST tomar sus valores del set declarado y MUST NOT declarar un set propio.

### Regla: Cada proyección tiene granularidad, vida y momento de ratificación propios

Una **proyección** es el uso que un consumidor hace del set declarado, con su propia granularidad, su propia vida y su propio momento de ratificación. Hoy hay dos: la línea `Components:` de cada capability-spec (por-capability, durable, ratificada spec por spec) y el valor del cambio (por-cambio, muere con el cambio, ratificado una sola vez en la INTAKE GATE). Toda proyección nueva MUST declarar esas tres cosas.

### Regla: Las proyecciones comparten el set y el gate, nunca el valor

Dos proyecciones MUST compartir exactamente dos cosas: el **set declarado** y el **gate presence-based**. MUST NOT compartir el **valor**: una ratificación hecha en una granularidad MUST NOT escribirse, copiarse ni propagarse a otra. En particular, propagar automáticamente el valor de una granularidad más gruesa a una más fina (batch) está prohibido.

### Regla: Un cambio que toca capabilities de componentes distintos

Cuando un cambio toca capabilities cuyos componentes difieren entre sí, el valor por-cambio MUST NOT decidir los componentes de cada capability. El valor por-cambio describe **el cambio**; cada capability conserva el suyo, ratificado spec por spec.

### Regla: Los contratos existentes no cambian

Este cambio MUST NOT alterar la forma de `repo.components` (bloque top-level del config de proyecto, array de objetos `{name, paths}`) ni el contrato de la línea `Components:`. El permiso para tocarlos existe —lo dio el intake— pero no se ejerce.

## Escenarios

### Scenario: un consumidor nuevo lee la definición sin pasar por el store de specs

- **GIVEN** un consumidor que necesita el concepto repo-level (qué es un componente, cómo se declara el set, el gate)
- **WHEN** busca la definición
- **THEN** la encuentra completa en su hogar canónico, sin abrir la documentación del concepto de capability-spec

### Scenario: la documentación del capability-spec conserva sólo su proyección

- **GIVEN** la documentación del concepto de capability-spec después del cambio
- **WHEN** se busca en ella la definición repo-level
- **THEN** sólo encuentra la proyección spec —la línea del header y por qué el store nunca se parte por app— más una cita al hogar canónico

### Scenario: ninguna cita del concepto queda colgada

- **GIVEN** cada sitio del payload que hoy cita el concepto de components (4 citas del ancla, más las citas conceptuales de las capability-specs `Accepted`)
- **WHEN** se sigue cada cita después de la mudanza
- **THEN** resuelve a una sección existente que contiene efectivamente el contenido citado
- **AND** ninguna cita queda apuntando a un ancla que ya no existe, y ninguna nombra una ruta `payload/…` en texto que se lea en runtime

### Scenario: las dos proyecciones coexisten con valores distintos

- **verification:** `deferred → intake-components-flag`
- **GIVEN** un cambio cuyo valor por-cambio nombra `cli` y una capability tocada cuyo header nombra `api`
- **WHEN** el cambio se ejecuta y se cierra
- **THEN** ambos valores conviven sin modificarse: el header sigue diciendo `api` y el valor del cambio sigue diciendo `cli`

### Scenario: el gate apagado no deja línea en ningún capability-spec

- **GIVEN** un repo sin `repo.components` declarado
- **WHEN** se autoriza o se valida cualquier capability-spec
- **THEN** ningún spec lleva la línea `Components:` y ninguna validación se dispara por su ausencia

### Scenario: el gate apagado tampoco habilita el valor por-cambio

- **verification:** `deferred → intake-components-flag`
- **GIVEN** un repo sin `repo.components` declarado
- **WHEN** un cambio pasa por su ratificación de alcance
- **THEN** no existe valor por-cambio de componentes que ratificar

### Scenario: una ratificación por-cambio no escribe en la proyección durable

- **verification:** `deferred → intake-components-flag`
- **GIVEN** un valor de componentes ratificado una sola vez para el cambio
- **WHEN** el cambio se archiva
- **THEN** ninguna línea `Components:` de las capabilities tocadas se crea, modifica ni ensancha

### Scenario: capabilities de componentes distintos en un mismo cambio

- **verification:** `deferred → intake-components-flag`
- **GIVEN** un cambio que toca una capability de `api` y otra de `ui`, con un único valor por-cambio ratificado
- **WHEN** el cambio se cierra
- **THEN** cada capability conserva su propio valor y ninguna recibe el del cambio

### Scenario: un config y un spec válidos antes del cambio siguen siéndolo

- **GIVEN** un `repo.components` y un capability-spec con su línea `Components:`, ambos válidos antes del cambio
- **WHEN** se cargan y se validan después del cambio
- **THEN** se cargan y validan igual que antes: sin errores nuevos, sin findings nuevos y sin migración

## Referencias

- **Rule** → [`./repo-component-declaration.md`](./repo-component-declaration.md) — la declaración repo-level del set de componentes.
- **Rule** → [`./capability-spec-components-axis.md`](./capability-spec-components-axis.md) — la proyección por-capability.
