# Capability — Serial sdd-apply dispatch bounded by cost

- **Status:** Accepted
- **Date:** 2026-08-31
- **Components:** cli

## Propósito

Una descomposición serial de `sdd-apply` crece en costo a medida que la ejecución avanza (el prompt aumenta desde ~50k a ~540k tokens en una ronda de 487 solicitudes), así que el orquestador DEBE cortar la descomposición en slices más cortos ANTES de enviarla, basado en un signal de tamaño que la fase de tareas declara por Fase. Obligatorio y mecánico: sin turno del usuario, sin gate.

## Actores

- **Orquestador**: lee el `Est. lines` de cada Fase del retorno de `sdd-tasks`, acumula en orden de implementación, corta en límites de Fase cuando la acumulación superaría 600
- **Fase de tareas**: declara un `Est. lines` por Fase (estimación de adiciones + supresiones sobre la lista de tareas completa de esa Fase)
- **Ejecutor de tareas aplicadas**: adapta cada slice como un call de `sdd-apply`, el primero nuevo y los siguientes como continuation batches

## Precondiciones

- El retorno de `sdd-tasks` declara `### Breakdown` como una tabla con columna `Est. lines`, uno por Fase
- El total acumulado de `Est. lines` sobre todas las Fases es conocido antes de enviar
- Cada Fase es un átomo: nunca se divide entre slices
- El umbral es 600 líneas estimadas de cambio por slice, calibrado de un ciclo, no de una curva medida

## Flujo principal

1. El orquestador lee `### Breakdown` del retorno de `sdd-tasks` sostenido en memoria (nunca del artifact)
2. Camina las Fases en el orden declarado (orden de implementación), acumulando `Est. lines`
3. Abre el primer slice con la primera Fase
4. Para cada Fase siguiente: si el total actual + `Est. lines` de esta Fase > 600, cierra el slice actual y abre uno nuevo con esta Fase; si no, suma la Fase al slice actual
5. Despacha cada slice como su propio call de `sdd-apply`, llamado por su rango de Fases e ids de tareas ("Fase 1-2, tareas 1.1-2.4")
6. Los slices después del primero son continuation batches ordinarios: el kernel's Apply-Progress Continuity se aplica sin cambios — el prompt dice a la ejecución leer `apply-progress` primero y MERGE en él

## Ramas / flujos alternativos

- **Breakdown ausente o imparseable** → degradar silenciosamente a un dispatch, comportamiento de hoy
- **Todas las Fases caben en un slice** → un dispatch, idéntico a hoy
- **Una Fase individual supera 600** → se despacha sola, nunca dividida; la siguiente Fase abre un nuevo slice

## Casos borde

- **Fase estimada en 900** → su propio slice, todas sus tareas, sin cut a nivel de tarea
- **Cambio pequeño de 380 líneas totales** → un slice, sin corte
- **Error de degradación (columna faltante)** → silencio completo, un dispatch, sin error ni warning ni gate
- **El segundo slice se envía como continuation batch** → lee `apply-progress`, merge automático, sin re-lectura manual de artifact

## Reglas de negocio

- El corte es SOLO para dispatch serial — no se aplica dentro de una ronda concurrente, a través de rondas, ni a consolidation runs
- La acumulación es por Fase en orden de implementación, nunca rearreglada
- El 600 es calibrado de un ciclo medido (nueve batches manuales de ~190k–470k tokens; una ronda de 487 solicitudes cuyo per-quarter price se duplicó), no derivado de una curva lines-to-cost
- La Fase es el átomo — nunca se divide; si una Fase entera supera 600, se despacha sola
- El guard no pregunta — es mecánica de orquestación, nunca un gate, nunca un turno
- Para una Fase cuyos tasks están en parte marcados (parallel-group), se acumula el `Est. lines` completo como-está — deliberadamente sobre-contando, sesgando hacia un corte extra, porque hacer el estimate depender de los marks acopla dos números que `sdd-tasks` redacta independientemente
- La degradación silenciosa es por diseño — si la columna falta o es imparseable, un dispatch, nada reportado, nada fallido
- La cantidad del slice es reportada en una sola línea de aviso nombrando el conteo; un conteo de uno no se reporta en absoluto

## Entidades y estados

- **Slice** — un conjunto de Fases cuyo total `Est. lines` ≤ 600, formado precomputing antes de despachar; cada slice es su propio call de `sdd-apply`, el primero nuevo y los siguientes como continuation batches
- **Fase** — unidad indivisible de la descomposición; nunca se divide a pesar del tamaño
- **Umbral** — 600 líneas, el budget de cambio estimado por slice
- **Acumulación** — suma corriente de `Est. lines` a lo largo de Fases en orden de implementación; se reinicia en cada nuevo slice

## Errores de cara al actor

- **Breakdown sin `Est. lines`**: el guard degrada silenciosamente, nada reporta (escenario de artifact antiguo)
- **Breakdown carrieando el sentinel `None.`**: el guard degrada, un dispatch
- **Un valor `Est. lines` no es un entero**: el guard degrada
- **Nada falla** — el guard nunca abre un gate, nunca pregunta, nunca reporta error si se degrada

## Escenarios

### Scenario: Cuatro Fases se cortan en dos dispatches

- **GIVEN** Fases estimadas en 250, 200, 300, 150 líneas, en orden de implementación
- **WHEN** el guard forma la ronda serial
- **THEN** Fases 1-2 se despachan juntas (450) y Fases 3-4 juntas (450) — dos calls de `sdd-apply`
- **AND** el segundo call es un continuation batch que lee el `apply-progress` existente y merges en él

### Scenario: Una Fase individual sobre-presupuesto se despacha sola, nunca dividida

- **GIVEN** una Fase estimada en 900 líneas
- **WHEN** el guard forma la ronda
- **THEN** esa Fase es su propio dispatch, con todas sus tareas, y ningún corte a nivel de tarea se hace

### Scenario: Un cambio pequeño no se ve afectado

- **GIVEN** Fases totalizando 380 líneas estimadas
- **WHEN** el guard forma la ronda
- **THEN** exactamente un dispatch de `sdd-apply` se forma, idéntico al comportamiento de hoy

### Scenario: El corte nunca pregunta

- **GIVEN** cualquier breakdown que el guard corta, corriendo desatendido
- **WHEN** los slices se despachan
- **THEN** ningún gate se abre y nada espera al usuario; a lo sumo una línea de aviso reporta el conteo de slice

### Scenario: Una ronda paralela se deja intacta

- **GIVEN** un grupo de cuatro tareas marcadas con el mismo `· parallel-group:` id
- **WHEN** el orquestador forma esa ronda
- **THEN** el budget no se aplica a ella, y uno isolated run por tarea se despacha como hoy

### Scenario: Una Fase mixta cuenta completa

- **GIVEN** una Fase de 500 líneas estimadas cuyas tareas son mitad marcadas para un grupo paralelo
- **WHEN** el guard acumula
- **THEN** suma 500, no la cuota de tasks no-marcadas

### Scenario: La consolidation run no se corta

- **GIVEN** una consolidation run siguiendo a una ronda de muchos reports de tareas
- **WHEN** se despacha
- **THEN** el budget no se aplica a ella, y esto se declara como residual exposure conocida

### Scenario: Estimates ausentes o no-parseables degradan silenciosamente

- **GIVEN** un tasks return cuyo `### Breakdown` carece de columna `Est. lines` (artifact antiguo)
- **WHEN** el guard corre
- **THEN** no dispara, un solo dispatch de `sdd-apply` se forma, nada se reporta como error

### Scenario: Breakdown que nunca fue producido

- **GIVEN** un `### Breakdown` que lleva solo el sentinel `None.`
- **WHEN** el guard corre
- **THEN** no dispara y forma un dispatch

## Referencias

- **Signal** → La columna `Est. lines` de `### Breakdown` en el retorno de `sdd-tasks` (nuevo)
- **Orden de implementación** → El orden en que las Fases se listan en `### Breakdown`
- **Continuation batch** → Kernel's Apply-Progress Continuity — el prompt dice a `sdd-apply` leer `apply-progress` y MERGE
- **Threshold provenance** → `structure/serial-apply-dispatch-budget.md` en el repo del proyecto — uno-ciclo calibración, no una curva; rutas de recalibración; ambas direcciones de fallo nombradas
- **Mechanical shape precedent** → `Parallel-Mark Validation` en el fragment de desarrollo — mandatory, silent, degrada sin error
