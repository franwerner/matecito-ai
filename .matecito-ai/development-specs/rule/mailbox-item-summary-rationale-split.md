# Capability — Split de items buzón en summary/rationale

- **Status:** Accepted
- **Date:** 2026-08-04
- **Components:** cli

## Propósito

Permitir que items en secciones buzón (5 secciones de fase-returns que imprimen gates) declaren un split opt-in de dos partes: un `summary` breve que se imprime en el gate, y una `rationale` completa que siempre se emite pero no se imprime por defecto. Solo la presentación se divide; la emisión es total.

## Actores

- **Agente SDD**: dispara la renderización de secciones buzón en su retorno de fase
- **Orquestador**: imprime el summary de cada item en un gate; en caso de pedido, reproduce la rationale del bloque en contexto
- **Validador**: rechaza items incompletos y detecta drift entre contrato y template

## Precondiciones

- La sección declara `items.rationale: rationale` en su bloque `items` del contrato `.yaml`
- Cada item en la sección lleva ambas partes: `summary` y `rationale`

## Flujo principal

1. El contrato de fase-returns declara `items.rationale` en una sección buzón
2. El agente redacta los items con dos campos: `summary` y `rationale`
3. El renderizador construye líneas `- {summary}` seguidas de `· rationale: {texto}`
4. El orquestador imprime solo la línea `- {summary}` al mostrar el gate
5. El bloque completo (ambas partes) se persiste en el `detailed_report`

## Ramas / flujos alternativos

- **Sección sin declaración de split** → el renderizador trata la sección como hoy, con un item por fila de `text`
- **Item con token de bloqueo** → el token se emite entre el summary y la rationale, en línea separada
- **Sección sin items** → se emite el sentinel `None.` (comportamiento sin cambios)

## Casos borde

- **Item sin rationale o summary vacío** → renderizador falla nombrando el item y el campo; exit 1, stdout vacío
- **Rationale que cruza varias líneas** → renderizador falla nombrando el item y el campo; exit 1
- **Summary que cruza varias líneas** → renderizador falla nombrando el item y el campo; exit 1
- **Item con campo `text` multi-línea en sección no-declarante** → renderizador falla nombrando el item y el campo; exit 1 (tightening: antes salía markdown malformado, exit 0)

## Reglas de negocio

- El split es opt-in: declarar `items.rationale` en el contrato es la única forma de habilitar el comportamiento
- Impresión vs emisión son conceptos distintos: la emisión es total (ambas partes siempre viajan en el bloque persistido); la presentación (lo que imprime el gate) es declarada por el contrato, nunca juzgada por brevedad
- Toda sección que renderiza items es un `render: items` con o sin split; el validador rechaza items multi-línea en cualquier items-section
- La rationale DEBE ser una línea no-vacía; la ausencia o vaciado causa fallo de renderización
- El marker `· rationale:` es literal en el template `.md` de retorno si la sección declara el split; su ausencia/presencia es audit-trail de contrato-template coherencia

## Entidades y estados

- **Sección buzón** — sección de gate cuya salida es una lista de items. Puede declarar split o no. Declarante: secciones de sdd-propose (Scope and approach), sdd-spec (Derived capabilities), sdd-design (New Decisions, Open Questions), sdd-tasks (Tasks not traceable), sdd-apply (Unmandated Forks, Mandated Departures).
- **Item** — entrada de una lista. Pre-split: un solo campo `text`. Post-split: dos campos (`summary`, `rationale`).

## Errores de cara al actor

- **Item incompleto al renderizar**: no se emite bloque, exit 1, stderr nombra la sección, el item y el campo faltante
- **Campo multi-línea al renderizar**: no se emite bloque, exit 1, stderr nombra la sección, el item y el campo
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

## Referencias

- **Contexto de negocio** → Intake Brief `sdd/mailbox-summary-rationale-split/intake`; Spec `sdd/mailbox-summary-rationale-split/spec`
