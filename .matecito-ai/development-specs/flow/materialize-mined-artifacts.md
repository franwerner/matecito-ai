# Capability — Materializar artefactos minados post-confirmación

- **Status:** Accepted
- **Date:** 2026-08-03
- **Components:** cli

## Propósito

Después del gate de confirmación de una ejecución de minería, materializar cada candidato confirmado a archivo — obteniendo el cuerpo y entradas de INDEX del renderer, no componiéndolas a mano. El executor sigue sin escribir; toda la escritura ocurre en el thread principal tras la confirmación.

## Actores

- El thread principal del orquestador (post gate de confirmación, ambos executors: development-decisions-mine y development-spec-mine)
- El usuario que confirmó candidatos en el gate

## Precondiciones

- Una ejecución de minería completó y produjo candidatos
- El usuario confirmó al menos un candidato en el gate de confirmación
- El contrato (`.yaml`) para el tipo de artefacto existe

## Flujo principal

1. Para cada candidato confirmado:
   - Construir un objeto `--data` a partir del candidato, usando el adapter específico de la skill
   - Invocar `render-artifact.js --type <type> --data <json>` para obtener el cuerpo
   - Escribir el cuerpo al archivo en `.matecito-ai/<path>`
   - Invocar `render-artifact.js --type <type> --data <json> --index-entries` para obtener las entradas de INDEX
   - Acumular la entrada (agrupada por domain/type)
2. Después de todos los candidatos materializados:
   - Leer el `INDEX.md` raíz (o crear si no existe usando el scaffold del contrato)
   - Mergedear las entradas acumuladas en el INDEX: una fila por domain/type, deduplicada por clave
   - Escribir el `INDEX.md` raíz exactamente una vez
3. Salir 0 si todo se escribió; salir 1 si alguna invocación del renderer falló

## Ramas / flujos alternativos

- **Renderer falla en datos inválidos** → No escribir el archivo; fallar la materialización; salida 1
- **No hay candidatos confirmados** → No escribir nada; salida 0 (no es error)
- **INDEX raíz no existe** → Usar el scaffold del contrato para crearlo; sumar las entradas
- **Gate de confirmación presenta candidatos** → El gate camina los candidatos a través del template compartido de items (índice único, uno a uno, "confirmar el resto" disponible a mitad)

## Casos borde

- **Varios candidatos en el mismo domain/type** → Una sola fila en INDEX raíz, con la clave domain/type; no duplicar
- **INDEX raíz bajo un store vacío** → Crear el INDEX raíz si el contrato declara un `scaffold:`
- **Candidato cuya `path` ya existe (retry)** → Sobrescribir sin preguntar; la materialización es idempotente tras confirmación
- **Entradas de INDEX del tipo per-domain y per-root** → Ambas se aplican: una fila en `.matecito-ai/edr/<domain>/INDEX.md` (per-domain) y una en `.matecito-ai/edr/INDEX.md` (root)

## Reglas de negocio

- El executor NUNCA escribe; todo ocurre en el thread principal tras la confirmación
- El contrato es la única fuente de verdad de forma y entradas
- El INDEX raíz se escribe exactamente una vez, después de todos los candidatos, nunca por candidato
- El renderer es el único que produce cuerpo e INDEX; el thread principal solo orquesta y escribe
- La materialización es idempotente: re-ejecutarla tras confirmación (retry) no crea duplicados

## Entidades y estados

- **Candidato minado confirmado** — Un candidato aceptado por el usuario en el gate. Estados: sin materializar → materializado → persistido en git
- **INDEX** — Registro de artefactos por domain/type (per-domain) y raíz. Estados: no existe → existe con entradas acumuladas

## Errores de cara al actor

- **Renderer falla** → La materialización se detiene; el usuario re-confirma tras corregir los datos
- **INDEX raíz no puede escribirse** → Fallar a stderr; salida 1; el usuario invoca manualmente tras diagnosticar
- **No hay candidatos confirmados** → No es error; materialización silenciosa (nada que hacer)

## Escenarios

### Scenario: Candidato confirmado se renderiza

- **GIVEN** un candidato confirmado en el gate
- **WHEN** se materializa en el thread principal
- **THEN** se invoca al renderer con los datos, se obtiene cuerpo, se escribe a `.matecito-ai/<path>`

### Scenario: Nada se escribe antes de confirmación

- **GIVEN** una ejecución de minería que produjo candidatos
- **WHEN** no se ha confirmado ninguno
- **THEN** no se escribieron artefacto ni INDEX, en ningún modo de ejecución (Interactive ni Automatic)

### Scenario: INDEX raíz se escribe una sola vez

- **GIVEN** varios candidatos confirmados en un lote
- **WHEN** se materializan todos
- **THEN** el `.matecito-ai/edr/INDEX.md` (raíz) se escribió exactamente una vez, después del último candidato; una fila por domain/type; no duplicados

### Scenario: Materialización es idempotente

- **GIVEN** candidatos materializados y confirmados
- **WHEN** se vuelve a ejecutar la materialización (retry tras un fallo parcial, con la misma confirmación)
- **THEN** los archivos se sobrescriben; el INDEX no duplica entradas; salida 0

### Scenario: INDEX por dominio y raíz

- **GIVEN** candidatos EDR de varios dominios (structure, contracts)
- **WHEN** se materializan
- **THEN** se escriben filas en `.matecito-ai/edr/structure/INDEX.md` y `.matecito-ai/edr/contracts/INDEX.md` (per-domain) Y una fila en `.matecito-ai/edr/INDEX.md` (raíz) que resume ambos por clave

### Scenario: Renderer falla en datos inválidos

- **GIVEN** un candidato adaptado con datos que faltan un campo requerido
- **WHEN** se invoca al renderer
- **THEN** el renderer falla a stderr nombrando el campo; no se escribe el archivo; la materialización se detiene; salida 1

### Scenario: El gate de confirmación presenta candidatos por el template compartido

- **GIVEN** una ejecución de minería que produjo candidatos
- **WHEN** su gate de confirmación abre
- **THEN** los candidatos se presentan a través del template compartido: un índice único, luego uno a uno, con "confirmar el resto" disponible antes del primero y a mitad del walk

### Scenario: Confirmar el resto a mitad del walk de candidatos

- **GIVEN** un walk de candidatos donde algunos ya fueron decididos
- **WHEN** el usuario pide confirmar el resto
- **THEN** los candidatos undecided se registran como confirmados, los decididos mantienen sus resultados, y exactamente el conjunto confirmado se materializa

### Scenario: Un candidato rechazado en el walk no se materializa

- **GIVEN** un candidato que el usuario rechaza mientras camina la lista
- **WHEN** se ejecuta la materialización
- **THEN** no se escriben archivo ni entrada de INDEX para él, y el resto del conjunto confirmado se materializa normalmente

## Referencias

- **EDR** → [`../../edr/structure/root-index-cardinality-per-domain-type.md`](../../edr/structure/root-index-cardinality-per-domain-type.md) — Decisión de que el INDEX raíz tiene una fila per domain/type, deduplicada por el caller
- **EDR** → [`../../edr/structure/contract-pair-in-templates.md`](../../edr/structure/contract-pair-in-templates.md) — Decisión de contrato ubicado en templates/
- **Contrato compartido** → [`../../shared/references/gate-presentation.md`](../../shared/references/gate-presentation.md) — Template compartido de presentación (índice, uno a uno, "confirmar el resto")
- **Contexto de negocio** → Spec de Intake `mine-materialization-renderer` — El paso 2 de descubrimiento motiva esta materialización en vez de prosa manual
