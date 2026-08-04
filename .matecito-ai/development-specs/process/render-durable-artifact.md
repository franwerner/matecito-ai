# Capability — Renderizar artefacto durable a partir de datos

- **Status:** Accepted
- **Date:** 2026-08-03
- **Components:** cli

## Propósito

Construir de manera determinista el cuerpo completo de un artefacto durable (EDR o capability-spec) a partir de datos estructurados y del contrato declarado del tipo. Cada encabezado, su orden y su presencia provienen del contrato, garantizando conformidad por construcción.

## Actores

- El paso de materialización post-gate de los executors de minería (development-decisions-mine, development-spec-mine)
- Bootstrap del árbol de specs (fuera de alcance en este ciclo)

## Precondiciones

- El contrato del tipo de artefacto existe (`.yaml` con `artifact`, `header`, `sections`, `index_entries`)
- Los datos del artefacto cumplen con los campos requeridos del contrato

## Flujo principal

1. Leer el contrato declarado para el tipo de artefacto
2. Validar que los datos contienen todos los campos requeridos; fallar nombando el campo faltante si no
3. Validar que los valores de cabecera están dentro de sus enums declarados; fallar nombre el valor ilegal si no
4. Construir el encabezado del archivo a partir de `header` del contrato y los datos
5. Para cada sección en el contrato, decidir presencia según `emitted` y `Status` del artefacto
6. Renderizar las secciones presentes en el orden declarado
7. Emitir el cuerpo completo a stdout

## Ramas / flujos alternativos

- **Campo requerido faltante** → Fallar a stderr nombrando el campo; no emitir texto; salir 1
- **Valor de cabecera fuera de enum** → Fallar a stderr nombrando el valor ilegal; no emitir texto; salir 1
- **Sección condicionada a Status no aplica** → Omitir la sección entera (no dejar vacía); continuar con la siguiente
- **Sección con `emitted: when-present` y sin datos** → Omitir la sección entera

## Casos borde

- **EDR `Status: Inferred` con `contexto` / `decision` / `reglas` / `alternativas` / `consecuencias` suministrados** → Las secciones Contexto y Decision se renderizan vacías; las secciones Reglas, Alternativas y Consecuencias se omiten entera (no presentes)
- **Capability-spec sin componentes** → La línea `- **Components:** ...` se omite entera (no se escribe `N/A`)
- **Sección omitible sin datos** → Se omite; no se deja `N/A` ni encabezado vacío
- **Línea de format en --schema con placeholders repetidos** → Se deduplica via `Set` antes de emitir

## Reglas de negocio

- La presencia de una sección está gobernada por su `emitted` y el `Status` del artefacto; nunca por el invocador
- Dos mecanismos de presencia distintos: `empty_on` (sección se emite con cuerpo vacío en ciertos Status) e `emitted: when-present` (sección se omite entera si no hay datos)
- El contrato es la única fuente de verdad para la forma; el rendidero aplica solo lo que declara
- Artefacto y entradas de INDEX se derivan de los mismos datos y no pueden diverger

## Entidades y estados

- **Artefacto durable** — Un archivo `.md` (EDR o capability-spec) con estructura declarada. Estados: no renderizado → renderizado → persisted (por el thread principal tras el gate, nunca por el executor)

## Errores de cara al actor

- **Campo requerido faltante** → `stderr: field <name> is required; stdout: empty; exit 1`
- **Valor de cabecera ilegal** → `stderr: value <value> is not in enum for <field>; stdout: empty; exit 1`
- **Tipo de artefacto no soportado** → `stderr: unknown artifact type <type>; exit 1`

## Escenarios

### Scenario: Cuerpo construido a partir de datos válidos

- **GIVEN** datos válidos para un tipo de artefacto
- **WHEN** se renderiza el cuerpo con `render-artifact.js --type <type> --data <json>`
- **THEN** stdout contiene el cuerpo con las secciones declaradas en el contrato, en orden, el encabezado relleno de datos

### Scenario: Datos inválidos (campo faltante) emiten error

- **GIVEN** datos que faltan un campo requerido
- **WHEN** se intenta renderizar
- **THEN** stderr nombra el campo faltante; stdout vacío; salida 1

### Scenario: Datos inválidos (valor de cabecera fuera de enum) emiten error

- **GIVEN** datos con un valor de cabecera (`Status`, `applied_pattern`, etc.) fuera de su enum declarado
- **WHEN** se intenta renderizar
- **THEN** stderr nombra el valor ilegal y el enum; stdout vacío; salida 1

### Scenario: EDR Inferred con contenido adversarial omite secciones

- **GIVEN** datos EDR con `Status: Inferred` y `contexto`, `decision`, `reglas`, `alternativas`, `consecuencias` suministrados no-vacíos
- **WHEN** se renderiza
- **THEN** Contexto y Decision se renderizan vacíos; Reglas, Alternativas y Consecuencias se omiten entera; sección Evidencia presente (conforme `empty_on: [Inferred]`)

### Scenario: Capability-spec sin componentes omite la línea

- **GIVEN** datos de capability-spec sin `components` o vacío
- **WHEN** se renderiza
- **THEN** la línea `- **Components:** ...` se omite entera; no se escribe `N/A`

### Scenario: Sección omitible sin datos se omite

- **GIVEN** capability-spec con `emitted: when-present` en una sección (e.g., Precondiciones) y sin datos para ella
- **WHEN** se renderiza
- **THEN** la sección se omite entera, no presente

### Scenario: Entradas de INDEX desde los mismos datos

- **GIVEN** datos válidos para un artefacto
- **WHEN** se solicitan entradas de INDEX con `render-artifact.js --type <type> --data <json> --index-entries`
- **THEN** se emiten líneas `row:` del contrato, con `path` y campos derivados de los datos, coincidiendo con el cuerpo

### Scenario: Forma reportada y agnóstica del productor

- **GIVEN** un tipo de artefacto
- **WHEN** se solicita su forma esperada con `render-artifact.js --type <type> --schema`
- **THEN** se reportan required, optional y derived sin mencionar JSON de minería ni productor específico

## Referencias

- **EDR** → [`../../edr/structure/contract-pair-in-templates.md`](../../edr/structure/contract-pair-in-templates.md) — Decisión de ubicación y granularidad del contrato
- **EDR** → [`../../edr/contracts/data-contract-derived-and-producer-neutral.md`](../../edr/contracts/data-contract-derived-and-producer-neutral.md) — Decisión de que --data es agnóstica del productor
- **EDR** → [`../../edr/structure/two-scripts-render-and-validate.md`](../../edr/structure/two-scripts-render-and-validate.md) — Decisión de dos scripts distintos
