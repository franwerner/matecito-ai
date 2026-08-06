# Capability — Consultar un fork sin mandato antes de aplicarlo

- **Status:** Accepted
- **Date:** 2026-08-06
- **Components:** cli

## Propósito

Mientras ejecuta trabajo dentro de su mandato, `sdd-apply` puede encontrar un fork que los artefactos confirmados no fijan. El punto de control se desplaza desde el retorno (clasificar después de la edición) a antes de la edición: un fork con más de una resolución válida vuelve atrás como una pregunta; un fork con ninguna alternativa válida se absorbe, y solo cuando el ejecutor puede nombrar qué cerró las alternativas.

## Actores

- **sdd-apply**: evalúa si una edición abre un fork no-fijado, y consulta antes de aplicar si hay más de una resolución válida
- **Ejecutor**: nombra la restricción cuando absorbe una desviación forzada
- **Orquestador**: imprime los forks `chosen` (sin mandato) en una sección Tier 1; los forks `forced` y `covered` en Tier 2

## Precondiciones

- `sdd-apply` tiene un contexto de edición que evaluar (las tareas, el diseño, el artefacto que se especificó)
- Existe un mecanismo de stop-mid-batch para enrutar preguntas antes de que la edición aterrice

## Flujo principal

1. Una edición mandatada abre, como efecto colateral, un punto que los artefactos no fijan
2. `sdd-apply` evalúa si el punto tiene más de una resolución válida
3. Si hay más de una opción válida → parar; no aplica ninguna. La pregunta viaja en `### Unmandated Forks` (Tier 1)
4. Si no hay alternativa válida → aplicar y reportar la desviación; nombrar la restricción que la forzó. La desviación viaja en `### Mandated Departures` (Tier 2)
5. Si un fork `chosen` se reporta (la pregunta llegó sin resolverse) → el Unresolved Decisions Guard del orquestador lo detiene antes de la siguiente fase

## Ramas / flujos alternativos

- **La restricción no puede ser nombrada**: no era forzada. Reclasificar como `chosen` y consultar antes de aplicar
- **El artefacto (spec, design, task, Accepted EDR) ya fija el punto**: aplicar sin parar; el token es `mandate: covered`

## Casos borde

- **El ejecutor elige absorber un fork `chosen` sin poder nombrar la restricción**: inválido. La sección Tier 1 de `### Unmandated Forks` es la consulta previa, no un vertedero de elecciones ya hechas
- **Dos alternativas parecen igual de válidas**: ambas son válidas — es un fork `chosen`, no se aplica, se consulta
- **El fork no es contenido nuevo, sino una edición sobre contenido existente**: el criterio es agnóstico al tipo de contenido — lo que importa es si la edición abre un punto y si hay más de una resolución

## Reglas de negocio

- El test es agnóstico al tipo de contenido: se aplica a cualquier edición, en cualquier archivo, independientemente de qué describe ese archivo
- `mandate:` responde «¿quién lo decidió?». Los valores: `covered` (un artefacto lo fija), `forced` (ninguna alternativa, restricción nombrada), `chosen` (había más de una opción)
- `verify-checks:` responde «¿sdd-verify lo audita contra el design?» — y sigue siendo independiente de `mandate:`. Un fork `chosen` no aplicado puede llevar `verify-checks: no` (nada llegó a auditar)
- Nombrar la restricción es obligatorio para `forced`. La imposibilidad de nombrarla significa que no era forzada
- El mecanismo de consulta es la sección Tier 1 del buzón de desviaciones, no un mecanismo nuevo — `sdd-apply` ya tiene stop-mid-batch para otros puntos de parada

## Entidades y estados

- **Fork**: un punto en la edición donde más de una resolución es válida
- **Desviación cubierta** (`mandate: covered`): un artefacto confirmado (spec, design, task, Accepted decision record) fija el punto
- **Desviación forzada** (`mandate: forced`): no hay alternativa válida; una restricción concreta compelió la elección
- **Desviación no-mandatada / elegida** (`mandate: chosen`): más de una resolución era válida y el ejecutor eligió una
- **Buzón `### Unmandated Forks`** (Tier 1): desviaciones `chosen` — no aplicadas, consultadas antes de aterrizaje; la pregunta sin resolver detiene la siguiente fase
- **Buzón `### Mandated Departures`** (Tier 2): desviaciones `covered` o `forced` — reportadas, no bloqueantes

## Errores de cara al actor

- **Ejecutor reporta una desviación como `forced` pero no puede nombrar la restricción**: inválido; reclasificar como `chosen` y no aplicar
- **Un `chosen` aterrizó en el código sin pasar por consulta**: CRITICAL — el gate debe haber parado la edición

## Escenarios

### Scenario: Un punto se abre como efecto colateral de una edición mandatada

- **GIVEN** una tarea dentro del mandato cuya edición fuerza una segunda elección que nadie acordó
- **AND** los artefactos confirmados no fijan esa elección y más de una resolución es válida
- **WHEN** el ejecutor llega a ese punto
- **THEN** no aplica ninguna de ellas, y el fork viaja atrás como una pregunta antes del siguiente despacho

### Scenario: Los artefactos ya lo fijan

- **GIVEN** un punto que la spec, el design, una tarea, o un Accepted decision record fija
- **WHEN** el ejecutor llega a él
- **THEN** aplica lo que el artefacto fija, sin parar y sin abrir una consulta

### Scenario: La pregunta lleva una propuesta

- **GIVEN** un fork consultado para el cual el ejecutor tiene una resolución preferida
- **WHEN** devuelve la pregunta
- **THEN** la preferencia se declara como una recomendación
- **AND** ninguna edición la implementa en el repositorio

### Scenario: Absorber una desviación requiere nombrar qué la forzó

- **GIVEN** una edición cuya única alternativa válida dejaría el cambio roto
- **WHEN** el ejecutor la absorbe
- **THEN** la aplica, la reporta con la restricción nombrada, y el flujo no se detiene

### Scenario: Forzada pero sin nombre

- **GIVEN** una edición absorbida para la cual no se puede nombrar una restricción concreta
- **WHEN** el ejecutor prepara su retorno
- **THEN** no es elegible para absorción y es consultada en su lugar

### Scenario: Toda desviación reportada declara quién la decidió

- **GIVEN** una desviación reportada
- **WHEN** el retorno y el artefacto se producen
- **THEN** el ítem lleva `mandate: covered | forced | chosen`, junto a `verify-checks:`, en el orden declarado

### Scenario: Token `mandate:` ausente o hedgeado se lee como `chosen`

- **GIVEN** una desviación con no `mandate:` token, o uno impreciso
- **WHEN** el orquestador lee el retorno
- **THEN** se trata como `chosen` y es gateada como tal

## Referencias

- **Conceptualmente relacionado**: `rule/mailbox-item-summary-rationale-split.md` — ambas secciones del buzón de desviaciones declaran `items.rationale` y honran el split
- **Artefactos de diseño**: Spec `sdd/sdd-apply-deviation-consult/spec`, Design `sdd/sdd-apply-deviation-consult/design`
