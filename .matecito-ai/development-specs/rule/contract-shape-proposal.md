# Capability — Proponer forma de contrato no especificada

- **Status:** Accepted
- **Date:** 2026-08-14
- **Components:** cli

## Propósito

Una forma de datos sin especificar — la única cosa que el sistema está prohibido de adivinar — se propone hoy como prosa libre en un slot genérico de stop: sin límite, multiline, sin campos declarados, así que nada verifica que los campos y sus tipos fueron mostrados alguna vez. Y ninguna forma de gate presenta un item que es un conjunto de campos tipados en lugar de una oración. Tras este cambio una forma de contrato se propone como un item compuesto único, declarado, verificado mecánicamente; un gate la presenta como un item a un ritmo que el usuario fija; y la forma que el usuario ratifica vuelve a la fase que la pidió.

## Actores

- **Fase SDD**: propone un contrato cuya forma no está especificada, emitiendo `### Contract Shapes Proposed` en lugar de solo prosa libre en `### Blocker`
- **Gate / Orquestador**: recibe la forma propuesta como item compuesto, la presenta uno o varios a ritmo que el usuario elige (uno-a-uno o todos-a-la-vez)
- **Usuario**: ratifica, ajusta o rechaza la forma; la forma ratificada se devuelve en las instrucciones de re-despacho de la fase
- **Fase re-despachada**: recibe la forma ratificada en sus instrucciones, nunca re-lee de un almacenamiento

## Precondiciones

- La fase declara un stop por contrato no especificado, no por otra razón (ambigüedad, conflicto, fork)
- El item compuesto lleva: un summary de una línea, un anchor (fuente de la necesidad), y una línea por campo (nombre, tipo, descripción)
- La forma propuesta no lleva slot de persistencia declarado

## Flujo principal

1. Fase descubre que necesita un contrato pero su forma no está especificada
2. Emite `### Contract Shapes Proposed` con un item por contrato, cada item compuesto
3. Gate presenta el contrato (summary + campos) como un único item
4. Si hay múltiples contratos, ofrece ritmo: uno-a-uno o todos-a-la-vez, antes de mostrar el primero
5. Usuario ratifica, ajusta o rechaza
6. Orquestador devuelve la forma ratificada en las instrucciones de re-despacho de la fase (anchor + summary + lista de campos)
7. Fase re-despachada recibe la forma en sus instrucciones y continúa; nunca re-lee de almacenamiento

## Ramas / flujos alternativos

- **Stop sin contrato** → La sección está ausente; no es un defecto
- **Contrato propuesto de nuevo en el mismo cambio** → Se presenta en su totalidad de nuevo; no hay forma corta de re-confirmación
- **Ritmo = uno-a-uno (default)** → Cada contrato se presenta, decide, y continúa al siguiente
- **Ritmo = todos-a-la-vez** → Se presentan juntos pero cada uno toma un resultado
- **Un solo contrato** → No se ofrece ritmo

## Casos borde

- **Contrato con muchos campos** → Se emite completo, ninguno se elide ni se abrevia
- **Campo con descripción larga** → Tiene su propio límite (160 caracteres); excederse falla la renderización
- **Contrato mezcla campos nuevos con campos de estructura existente** → Un anchor único; la distinción se declara en las descripciones de campo, no como marcador por-campo

## Reglas de negocio

- El test de alcance: una forma que persiste, cruza límite o es pública está gobernada. *Cruzar límite* ahora significa alcanzar base de datos, red, o superficie fuera de esta herramienta. Una forma leída solo en coordinación intra-flujo no cruza y se maneja como decisión arquitectónica ordinaria
- Cada contrato es UN item: uno índice, un turno, un resultado. Campos nunca se dividen en items
- Un item compuesto lleva un summary, un anchor, y líneas de campo (nombre — tipo — descripción); nada más en slots declarados
- Anchor: uno por contrato (no por-campo); se usa para identificar la forma cuando vuelve
- La sección `### Contract Shapes Proposed` aparece solo cuando fase declara `has_contract_proposals: true` Y el status es `blocked`; ambas condiciones se requieren
- Una forma ratificada viaja **solo en las instrucciones de re-despacho**, nunca se re-lee de almacenamiento. Fase re-despachada sin forma en sus instrucciones: retorna `blocked` nombrando la forma faltante, nunca adivina
- Todo destino después difiere: `sdd-spec` y `sdd-design` escriben en su propio artefacto; `sdd-apply` escribe en el código (per regla de dominio)
- Ningún contrato propuesto lleva field indicando dónde vivirá la forma almacenada

## Entidades y estados

- **Item compuesto** — Propuesta de contrato: summary (≤250 chars), anchor, lista de campos. Estados: pendiente → ratificado → ajustado → rechazado
- **Campo de contrato** — Nombre, tipo, descripción (≤160 chars cada descripción). Emitido como línea `· field: {nombre} — {tipo} — {descripción}`
- **Oferta de ritmo** — Cuando hay ≥2 contratos: ofrecer uno-a-uno (default) o todos-a-la-vez antes de mostrar el primero

## Errores de cara al actor

- **Field missing its type**: la renderización falla nombrando item + field + parte; exit 1
- **Description over 160 chars**: falla durante renderización; exit 1
- **Item con cero fields**: falla nombrando el item; exit 1
- **Re-despacho sin forma en instrucciones**: fase retorna `blocked` nombrando el contrato faltante

## Escenarios

### Scenario: Una forma que sale de la herramienta

- **GIVEN** un payload devuelto por la red o una fila escrita a base de datos
- **WHEN** su forma no está especificada
- **THEN** está gobernada, y se propone como contrato antes de que nada la fije

### Scenario: Una forma leída solo entre pasos del flujo

- **GIVEN** un schema cuya única lectura es en pasos de esta herramienta coordinándose entre sí
- **WHEN** se aplica el test angostado
- **THEN** no cruza límite por ese test, y se maneja como decisión arquitectónica ordinaria

### Scenario: El test angostado no gana segunda cláusula

- **GIVEN** la declaración del test de alcance tras este cambio
- **WHEN** se lee
- **THEN** solo el significado de "cruzar" cambió, y ninguna excepción de "ya ratificado en otro lado" aparece en él

### Scenario: Una forma transitoria dentro de una unidad de trabajo

- **GIVEN** una forma usada dentro de una sola función y nunca saliendo de ella
- **WHEN** se aplica el test
- **THEN** no está gobernada, exactamente como antes de este cambio

### Scenario: Un contrato con muchos fields

- **GIVEN** una propuesta de contrato que lleva catorce fields
- **WHEN** se emite
- **THEN** es un item llevando catorce líneas de field, ninguno se cae ni se abrevia

### Scenario: Field missing su type

- **GIVEN** un contrato propuesto con un field cuyo type está ausente o vacío
- **WHEN** se emite
- **THEN** la emisión falla nombrando item, field y parte faltante; nada se emite

### Scenario: Item que declara fields pero no lleva ninguno

- **GIVEN** un item declarando que lleva forma de contrato pero listando cero fields
- **WHEN** se emite
- **THEN** la emisión falla nombrando el item; nada se emite

### Scenario: Nada dice dónde vivirá la forma

- **GIVEN** una propuesta de contrato presentada al usuario
- **WHEN** se muestra
- **THEN** ninguna declaración de dónde se almacenará la forma aparece, en slot ni como prosa

### Scenario: El mismo contrato propuesto de nuevo en el cambio

- **GIVEN** un contrato propuesto una segunda vez, más adelante en el mismo cambio
- **WHEN** llega a un gate
- **THEN** se presenta en su totalidad; no existe forma corta de re-confirmación para una forma de contrato

### Scenario: Un contrato mezcla fields nuevos con de estructura existente

- **GIVEN** una forma propuesta donde algunos fields son nuevos y otros vienen de una estructura que existe
- **WHEN** se emite
- **THEN** lleva un anchor para el contrato completo, y la distinción entre los dos tipos de field se declara en las descripciones de field, nunca como marcador por-field

### Scenario: Un stop por razón que no es contrato

- **GIVEN** una fase que para por derivación ambigua, conflicto con decisión aceptada, o fork que no puede resolver
- **WHEN** su retorno se chequea
- **THEN** la sección está ausente, y su ausencia no se reporta como defecto

### Scenario: La sección en status que la excluye

- **GIVEN** un retorno que no es stop pero lleva la sección
- **WHEN** se chequea
- **THEN** se reporta como violación

### Scenario: Un stop sobre un contrato no especificado

- **GIVEN** una fase que para porque un contrato que necesita no está especificado
- **WHEN** su retorno se emite
- **THEN** la sección está presente y lleva un item por contrato

### Scenario: Una forma ratificada viaja en instrucciones de re-despacho

- **GIVEN** un contrato que el usuario ratifica o ajusta en el gate
- **WHEN** la fase se re-despacha
- **THEN** la lista de fields ratificada es lo que viaja (o la lista ajustada si el usuario la modificó)
- **AND** viaja en el campo instructions del re-despacho, nunca en un artefacto almacenado

### Scenario: El usuario ajusta la forma antes de ratificar

- **GIVEN** un contrato que el usuario corrige en el gate
- **WHEN** la fase se re-despacha
- **THEN** la lista de fields ajustada es lo que viaja, y la originalmente ofrecida no

### Scenario: La fase implementadora llega sin forma en instrucciones

- **GIVEN** una fase implementadora llegando a trabajo gobernado por contrato cuyas instrucciones no llevan la forma
- **WHEN** busca la forma
- **THEN** para y nombra el contrato y la forma faltante, sin adivinar ni leer de `apply-progress`

### Scenario: Una fase que especifica la recibe

- **GIVEN** una fase que especifica re-despachada con una forma ratificada
- **WHEN** continúa su trabajo
- **THEN** la forma se escribe en su propio artefacto, del cual fases más adelante ya leen

### Scenario: Tres y dos campos en el mismo gate

- **GIVEN** un gate presentando dos contratos y dos items ordinarios
- **WHEN** abre
- **THEN** un índice cuenta cuatro y el walkthrough es plano, cada contrato mostrado con sus propios fields

## Referencias

- **Cómo se propone hoy** → Exploration `sdd/contract-shape-gate/explore` — Current State #1-2: dos o tres fields libres sin verificación en `### Blocker`
- **Precedente de sección dedicada** → Decision Record `contracts/verdict-in-its-own-conditional-section.md` (Accepted) — una sección nueva, condicional, con puerta no-derivable
- **Precisión de alcance** → Decision Record `contracts/narrowed-boundary-scope-test.md` (Accepted) — qué significa "cruzar límite"
- **Cómo vuelve la forma** → Decision Record `contracts/ratified-contract-travels-in-the-dispatch-prompt.md` (Accepted) — nunca desde almacenamiento, solo en instrucciones
- **Shape del item** → Decision Record `contracts/contract-item-anchors-once.md` (Accepted) — un anchor, no por-field
- **Líneas de field** → Decision Record `contracts/nested-field-continuation-line.md` (Accepted) — emitidas como `· field:` lineas
- **Tope por descripción** → Decision Record `contracts/per-field-description-cap.md` (Accepted) — 160 chars per descripción, el summary sigue en 250
