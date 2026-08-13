# Capability — Split de items buzón en summary/rationale

- **Status:** Accepted
- **Date:** 2026-08-04
- **Components:** cli

## Propósito

Permitir que items en secciones buzón (6+ secciones de fase-returns que imprimen gates) declaren un split opt-in de tres partes: un `summary` breve (con longitud máxima declarada) que se imprime en el gate, un `anchor` (fuente concreta del item, suministrado por la fase) y una `rationale` completa que siempre se emite pero no se imprime por defecto. Solo la presentación se divide; la emisión es total.

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
- Impresión vs emisión son conceptos distintos: la emisión es total (las tres partes siempre viajan en el bloque persistido); la presentación (lo que imprime el gate) es declarada por el contrato, nunca juzgada por brevedad
- El `summary_max` es declarado per-sección y se aplica en tiempo de renderización; un bloque escrito a mano nunca se mide
- El `anchor` es siempre suministrado por la fase, nunca derivado por el renderizador
- Toda sección que renderiza items es un `render: items` con o sin split; el validador rechaza items multi-línea en cualquier items-section
- La rationale DEBE ser una línea no-vacía; la ausencia o vaciado causa fallo de renderización
- El `anchor` DEBE ser presente (una línea no-vacía); la ausencia causa fallo de renderización
- Los markers `· anchor:` y `· rationale:` son literales en el template `.md` de retorno si la sección declara el split; su ausencia/presencia es audit-trail de contrato-template coherencia
- El orden de líneas es: `- {summary}`, luego `· anchor:`, luego tokens opcionales, luego `· rationale:` (anchor PRIMERO, antes de otros tokens)

## Entidades y estados

- **Sección buzón** — sección de gate cuya salida es una lista de items. Puede declarar split o no. Declarante (con split): secciones de sdd-propose (Scope and approach), sdd-spec (Derived capabilities, New Decisions), sdd-design (New Decisions, Open Questions), sdd-tasks (Tasks not traceable), sdd-apply (Rejected Proposals Checked, Unmandated Forks, Mandated Departures). Todas llevan `items.summary_max: 250` y token `anchor` si declaran split.
- **Item** — entrada de una lista. Pre-split: un solo campo `text`. Post-split: tres campos (`summary` ≤ límite, `anchor` libre-form, `rationale`).

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

## Referencias

- **Contexto de negocio** → Intake Brief `sdd/mailbox-summary-rationale-split/intake`; Spec `sdd/mailbox-summary-rationale-split/spec`; Cambio `sdd/sequential-decision-gate`
- **Criterio de anchor** → Sección D.3 de `~/.claude/skills/_shared/sdd-phase-common.md` — Formas legales (`<repo-path>[:line]` | `<engram-key>`), línea inicial solamente, regla target-no-existe-aún
- **Tope al summary** → Contrato de fase-returns declara `items.summary_max` per-sección; aplicado en render-return.js, NO en validate-return.js
