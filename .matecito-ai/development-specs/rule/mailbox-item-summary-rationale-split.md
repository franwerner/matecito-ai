# Capability — Split de items buzón en summary/rationale

- **Status:** Accepted
- **Date:** 2026-08-04
- **Components:** cli

## Propósito

Permitir que items en secciones buzón (12 secciones de fase-returns que imprimen gates, en cualquier forma de renderización) declaren un split opt-in de tres partes: un `summary` breve (con longitud máxima declarada) que se imprime en el gate, un `anchor` (fuente concreta del item, suministrado por la fase) y una `rationale` completa que siempre se emite pero no se imprime por defecto. Solo la presentación se divide; la emisión es total. La sección elige declarar el split, no su forma de renderización.

## Actores

- **Agente SDD**: dispara la renderización de secciones buzón en su retorno de fase
- **Orquestador**: imprime el summary de cada item en un gate; en caso de pedido, reproduce la rationale del bloque en contexto
- **Validador**: rechaza items incompletos y detecta drift entre contrato y template

## Precondiciones

- La sección declara `items.rationale: rationale` en su bloque `items` del contrato `.yaml`
- La sección declara `items.summary_max: 250` (o el límite específico) en caracteres
- Cada item en la sección lleva las tres partes: `summary`, `anchor` y `rationale`
- El `anchor` es suministrado por la fase, nunca derivado

## Flujo principal

1. El contrato de fase-returns declara `items.rationale`, `items.summary_max` y token `anchor` en una sección buzón
2. El agente redacta los items con tres campos: `summary` (máximo 250 caracteres), `anchor` y `rationale`
3. El renderizador valida que summary no supere el límite; si lo hace, falla con exit 1
4. El renderizador construye líneas `- {summary}`, luego `· anchor: {valor}`, luego opcionalmente tokens, luego `· rationale: {texto}`
5. El orquestador imprime solo las líneas `- {summary}` y `· anchor:` al mostrar el gate
6. El bloque completo (las tres partes + tokens) se persiste en el `detailed_report`

## Ramas / flujos alternativos

- **Sección sin declaración de split** → el renderizador trata la sección como hoy, con un item por fila de `text`
- **Item con token de bloqueo** → el token se emite entre el summary y la rationale, en línea separada
- **Sección sin items** → se emite el sentinel `None.` (comportamiento sin cambios)

## Casos borde

- **Item sin rationale, summary vacío o anchor faltante** → renderizador falla nombrando el item y el campo; exit 1, stdout vacío
- **Rationale que cruza varias líneas** → renderizador falla nombrando el item y el campo; exit 1
- **Summary que cruza varias líneas** → renderizador falla nombrando el item y el campo; exit 1
- **Summary que supera `summary_max`** → renderizador falla durante renderización; exit 1, stdout vacío
- **Anchor que cruza varias líneas** → renderizador falla nombrando el item y el campo; exit 1
- **Item con campo `text` multi-línea en sección no-declarante** → renderizador falla nombrando el item y el campo; exit 1 (tightening: antes salía markdown malformado, exit 0)
- **Bloque escrito a mano en lugar de renderizado** → summary nunca se mide en longitud (la hora de aplicar el límite es durante renderización; una vez cerrada la fase, el re-envío sería un re-despacho completo)

## Reglas de negocio

- El split es opt-in: declarar `items.rationale` en el contrato es la única forma de habilitar el comportamiento (las tres partes son obligatorias si se declara)
- El split es declarado por la sección en su contrato, nunca por su forma de renderización: una tabla puede declararlo, una labeled-list puede declararlo, una lista de items puede declararlo; la forma no restringe ni habilita la declaración
- Una única estrategia produce y valida las partes (summary, anchor, rationale) sean cuales sean la forma y render — nunca un mecanismo diferente por render-form
- Impresión vs emisión son conceptos distintos: la emisión es total (las tres partes siempre viajan en el bloque persistido); la presentación (lo que imprime el gate) es declarada por el contrato, nunca juzgada por brevedad
- El `summary_max` es declarado per-sección y se aplica en tiempo de renderización; un bloque escrito a mano nunca se mide
- El `anchor` es siempre suministrado por la fase, nunca derivado por el renderizador
- Toda sección que declara `items.rationale` lleva las tres partes en cada item, cualquiera que sea su forma de renderización; el validador rechaza items que falten cualquier parte
- La rationale DEBE ser una línea no-vacía; la ausencia o vaciado causa fallo de renderización
- El `anchor` DEBE ser presente (una línea no-vacía); la ausencia causa fallo de renderización
- Los markers `· anchor:` y `· rationale:` son literales en el template `.md` de retorno si la sección declara el split; su ausencia/presencia es audit-trail de contrato-template coherencia
- El orden de líneas es: `- {summary}` / `· anchor:` / tokens opcionales / `· rationale:` para listas; para tablas/labeled-lists el orden está definido en el template, pero mantiene anchor antes de rationale
- Las secciones de verify que declaran el split (`Decision Gaps`, `UI Verdict`, `Issues Found`) mantienen todas sus columnas/estructura hoy; el split es aditivo, nunca remodela los datos existentes
- Una sección MAY declarar un repeated typed field bajo `items.fields`; cada field MUST emitirse como una línea de continuación — después de los tokens, antes de la línea `· rationale:` — llevando las partes declaradas separadas por el separador declarado, así una descripción que contiene el separador sobrevive intacta
- Un field MUST NOT ser un sub-bullet ni una sub-tabla: las líneas emitidas MUST ser del mismo tipo que la verificación que lee el bloque ya camina (continúo-lineas con patrón `· name: value`)
- Una sección que no declara fields está sin cambios, byte por byte
- El `summary_max` MUST medir solo la línea del item summary, manteniéndose en 250 caracteres. Cada descripción de field MUST llevar su propio máximo de 160 caracteres. La cantidad de fields MUST NOT limitarse: todo field propuesto se emite
- El checker MUST reconocer las field lines como fueron emitidas, nunca saltarlas. Una field line malformada (part vacía) y un item declarante con cero field lines MUST cada uno producir un hallazgo nombrando el item ofensor. Un bloque emitido con field lines MUST leer de vuelta limpio, sin hallazgo de sección no-parseable

## Entidades y estados

- **Sección buzón** — sección de gate cuya salida puede ser una lista de items, una tabla tipada, o labeled-lists. Puede declarar split o no, independientemente de su forma. Puede opcionalmente declarar `fields` (lista tipada de sub-campos bajo el summary de cada item). Declarante de split (con split): dieciséis secciones (cuatro de contrato-shape + doce de mailbox ordinarias). Todas llevan `items.summary_max: 250` y token `anchor` si declaran split. Las cuatro declarantes de fields también llevan `items.fields: {key, parts: [...], separator, field_max}`.
- **Item/Fila/Entrada** — unidad de información en una sección. Pre-split: un solo campo `text` (en listas) o datos tipados (en tablas/labeled-lists). Post-split: tres partes adicionales (`summary` ≤ límite, `anchor` libre-form, `rationale`) que se emiten en todas las formas de renderización. Si la sección declara `fields`, el item también emite líneas `· field: {name} — {type} — {description}` (cada field capped a `field_max`).

## Errores de cara al actor

- **Item incompleto al renderizar**: no se emite bloque, exit 1, stderr nombra la sección, el item y el campo faltante (summary, anchor o rationale)
- **Summary sobre límite al renderizar**: no se emite bloque, exit 1, stderr nombra la sección, el item, el campo y ambas longitudes
- **Campo multi-línea al renderizar**: no se emite bloque, exit 1, stderr nombra la sección, el item y el campo (summary, anchor o rationale)
- **Marker en template sin declaración en contrato**: `--self-check` reporta DRIFT, exit 1
- **Declaración en contrato sin marker en template**: `--self-check` reporta DRIFT, exit 1

## Escenarios

### Scenario: Sección que declara el split

- **GIVEN** una sección cuyo contrato declara `items.rationale`
- **WHEN** se renderiza desde data con dos campos por item
- **THEN** cada línea emitida es `- {summary}` seguida de `  · rationale: {rationale}` (y opcionalmente `  · {token}: {valor}` si la sección tiene token)

### Scenario: Item sin rationale

- **GIVEN** una sección declarante con un item que carece de campo `rationale`
- **WHEN** se intenta renderizar
- **THEN** el renderizador falla nombrando la sección, el item y el campo `rationale`; exit 1, stdout vacío

### Scenario: Summary multi-línea en sección declarante

- **GIVEN** un item de sección declarante cuyo `summary` contiene un salto de línea
- **WHEN** se renderiza
- **THEN** el renderizador falla nombrando el item y el campo; exit 1, stdout vacío

### Scenario: Rationale multi-línea en sección declarante

- **GIVEN** un item de sección declarante cuyo `rationale` contiene un salto de línea
- **WHEN** se renderiza
- **THEN** el renderizador falla nombrando el item y el campo; exit 1, stdout vacío

### Scenario: Sección sin declaración mantiene comportamiento actual

- **GIVEN** una sección items sin declaración de split, con data de campo único `text`
- **WHEN** se renderiza
- **THEN** la salida es byte-idéntica a la pre-split, salvo que multi-línea en `text` ahora falla (antes: markdown malformado)

### Scenario: Item con token y rationale

- **GIVEN** una sección declarante cuyo contrato tiene `items.token` y la data lleva `token` e `rationale`
- **WHEN** se renderiza
- **THEN** se emite `- {summary}`, luego `  · {token}: {valor}`, luego `  · rationale: {rationale}`

### Scenario: Sección declarante vacía

- **GIVEN** una sección declarante cuya lista está vacía
- **WHEN** se renderiza
- **THEN** se emite el sentinel `None.`; el gate sigue leyéndolo como vacío (comportamiento sin cambios)

### Scenario: Template y contrato difieren

- **GIVEN** un contrato que declara la sección como split y un template cuya sección NO documenta la línea `· rationale:`
- **WHEN** corre `--self-check`
- **THEN** reporta DRIFT; exit 1

### Scenario: El gate imprime solo summary

- **GIVEN** un gate que expone una sección declarante
- **WHEN** el orquestador imprime la salida
- **THEN** cada item aparece como `- {summary}` solamente; la línea `· rationale:` no se imprime
- **AND** la línea `· rationale:` vive en el bloque persistido del `detailed_report`

### Scenario: Multi-línea en sección no-declarante es rechazada

- **GIVEN** una sección items sin split (e.g., `### Remaining Tasks`) cuyo `text` contiene un salto de línea
- **WHEN** se renderiza
- **THEN** el renderizador falla nombrando el item y el campo `text`; exit 1, stdout vacío

### Scenario: Un item split ahora declara tres partes, no dos

- **GIVEN** una sección cuyo contrato declara `items.rationale` y token `anchor`
- **WHEN** se renderiza desde data con los tres campos por item
- **THEN** cada línea emitida es `- {summary}`, luego `· anchor: {valor}`, luego opcionalmente `· {token}: {valor}`, luego `· rationale: {rationale}`

### Scenario: Item sin anchor

- **GIVEN** una sección declarante con un item que carece de campo `anchor`
- **WHEN** se intenta renderizar
- **THEN** el renderizador falla nombrando la sección, el item y el campo `anchor`; exit 1, stdout vacío

### Scenario: Anchor cuando el target no existe aún

- **GIVEN** un item que propone algo no escrito, así que no hay archivo para apuntar
- **WHEN** se redacta su anchor
- **THEN** nombra la fuente que surfaceó la necesidad, nunca una ubicación que no existe

### Scenario: Anchor cubriendo más de una línea

- **GIVEN** un item cuya fuente abarca varias líneas
- **WHEN** se redacta su anchor
- **THEN** nombra la ubicación y la línea inicial (sola), y la extensión se declara en palabras no como rango

### Scenario: Item informativo también ancla

- **GIVEN** un item informativo en una sección declarante (Tier 2, como Open Questions)
- **WHEN** se renderiza
- **THEN** lleva su anchor como cualquier otro item, y sigue sin bloquear

### Scenario: Summary dentro del límite declarado

- **GIVEN** una sección declarante cuyos items están todos dentro del máximo declarado
- **WHEN** se renderiza
- **THEN** la salida es exactamente como antes de este cambio (solo difiere en la presencia de anchor y su validación)

### Scenario: Summary sobre el límite declarado

- **GIVEN** un item cuyo summary supera el máximo de su sección
- **WHEN** se intenta renderizar
- **THEN** el renderizador falla nombrando la sección, el item y el campo; exit 1, stdout vacío

### Scenario: El gate imprime solo summary y anchor

- **GIVEN** un gate que expone una sección declarante
- **WHEN** el orquestador imprime la salida
- **THEN** cada item aparece como `- {summary}` seguido de `· anchor: {valor}` solamente; las líneas `· rationale:` y tokens no se imprimen
- **AND** la línea `· rationale:` y los tokens viven en el bloque persistido del `detailed_report`

### Scenario: Una sección tabla declara el split

- **GIVEN** una sección renderizada como tabla de filas tipadas, cuyo contrato declara `items.rationale`
- **WHEN** se renderiza desde data con filas que llevan `summary`, `anchor` y `rationale`
- **THEN** cada columna declarada se emite como hoy
- **AND** cada fila también emite sus partes (summary, anchor, rationale) en un bloque de detalle keyed por `items.key`
- **AND** se produce mediante el mismo mecanismo que una sección lista usa

### Scenario: Una sección labeled-list declara el split

- **GIVEN** una sección de labeled-lists, cuyo contrato declara `items.rationale`
- **WHEN** se renderiza
- **THEN** cada entry emite sus partes (summary, anchor, rationale) bajo su propia label
- **AND** se produce mediante el mismo mecanismo que una sección lista usa

### Scenario: Una entrada sin rationale en tabla/labeled-list declarante

- **GIVEN** una sección tabla o labeled-list declarante, con una entrada que carece de `rationale`
- **WHEN** se intenta renderizar
- **THEN** el renderizador falla nombrando la sección, la entrada y el campo faltante; exit 1, stdout vacío

### Scenario: Una entrada con summary sobre límite en tabla/labeled-list

- **GIVEN** una sección tabla o labeled-list declarante con `items.summary_max: 250`, con una entrada cuyo summary supera ese límite
- **WHEN** se intenta renderizar
- **THEN** el renderizador falla nombrando la sección, la entrada y el campo; exit 1, stdout vacío

### Scenario: Las secciones de verify adoptan el split manteniendo columnas

- **GIVEN** las secciones `## Decision Gaps`, `## UI Verdict` y `### Issues Found` del retorno de verify, todas declarando `items.rationale` en su contrato
- **WHEN** se renderizan desde data con las partes del split
- **THEN** cada sección emite todas las columnas que declara hoy (e.g., `structure`, `backing` en Decision Gaps; `scenario`, `counterpart`, `covers`, `state`, `failure_reason` en UI Verdict)
- **AND** además emite las partes del split (summary capped, anchor, rationale) para cada fila/entry
- **AND** las emisiones son totales — las tres partes + columnas siempre viajan en el bloque persistido

### Scenario: Dieciséis secciones declaran el split ahora

- **GIVEN** la definición canónica de cuántas secciones buzón declaran el split
- **WHEN** se lee después de este cambio
- **THEN** declara dieciséis secciones: doce de fases ordinarias (propose 1, spec 2, design 2, tasks 1, apply 3, más cuatro de contrato-shape en propose/spec/design/apply) y tres de verify (Decision Gaps, UI Verdict, Issues Found como una sección)
- **AND** declara seis pares de contrato-template (cinco de fases ordinarias, uno de verify) — la cuarta forma (compuesto) aterriza dentro de los cuatro pares existentes, no crea pares nuevos

### Scenario: Una sección declarante renderiza un item con fields

- **GIVEN** una sección declarante y un item llevando tres fields
- **WHEN** se renderiza
- **THEN** las líneas emitidas son el summary, luego el anchor, luego cualesquier tokens, luego las tres líneas de field, luego la línea de rationale

### Scenario: Una descripción conteniendo el separador

- **GIVEN** un field cuya descripción contiene el separador declarado
- **WHEN** se renderiza y se lee de vuelta
- **THEN** la descripción sobrevive completa, porque solo las primeras ocurrencias del separador separan las partes

### Scenario: Líneas de field no se cuentan como nuevos items

- **GIVEN** un item renderizado llevando líneas de field
- **WHEN** el bloque se camina por el checker que cuenta items
- **THEN** cada línea de field se lee como continuación de su item, nunca como nuevo item

### Scenario: Una sección que no declara fields

- **GIVEN** una sección cuyo contrato no declara fields
- **WHEN** se renderiza
- **THEN** su salida es byte-idéntica a antes de este cambio

### Scenario: Una descripción de field sobre su máximo

- **GIVEN** un item cuya descripción de field excede 160 caracteres
- **WHEN** se renderiza
- **THEN** la renderización falla nombrando el item, field y parte; nada se emite

### Scenario: Summary sobre el límite declarado

- **GIVEN** un item cuyo summary excede el máximo de su sección
- **WHEN** se renderiza
- **THEN** la renderización falla nombrando la sección, item y campo; exit 1, stdout vacío

### Scenario: Muchos fields, todos dentro de sus límites

- **GIVEN** un item llevando muchos fields, cada uno dentro de su máximo
- **WHEN** se renderiza
- **THEN** cada field se emite, ninguno se cae ni se abrevia

### Scenario: Una línea de field malformada

- **GIVEN** un bloque emitido con una línea de field faltando una de sus partes
- **WHEN** el bloque se chequea
- **THEN** un hallazgo nombra el item y el field, y el check no pasa

### Scenario: Un bloque va derecho de emisión a verificación

- **GIVEN** un item compuesto renderizado por el emisor
- **WHEN** su salida se alimenta directamente al checker
- **THEN** pasa, sin hallazgo de que la sección no es parseable y sin ninguna línea de field saltada en silencio

## Referencias

- **Contexto de negocio** → Intake Brief `sdd/mailbox-summary-rationale-split/intake`; Spec `sdd/mailbox-summary-rationale-split/spec`; Cambio `sdd/sequential-decision-gate`
- **Criterio de anchor** → Sección D.3 de `~/.claude/skills/_shared/sdd-phase-common.md` — Formas legales (`<repo-path>[:line]` | `<engram-key>`), línea inicial solamente, regla target-no-existe-aún
- **Tope al summary** → Contrato de fase-returns declara `items.summary_max` per-sección; aplicado en render-return.js, NO en validate-return.js
