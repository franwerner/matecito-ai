# EDR — El despacho paralelo es batch-bound: se integra recién cuando vuelve todo el batch

- **Status:** Accepted
- **Date:** 2026-08-10

## Contexto

Cuando un batch de implementación trae tareas marcadas independientes, hay dos formas de encadenar el
despacho con la integración: esperar a que cierre cada corrida y consolidarla al toque, o esperar a que
cierre el batch entero y consolidar todo junto. Medido sobre un batch real, el reloj de pared lo domina
la corrida más lenta (112s), y una consolidación —cherry-pick más escritura de artefacto— toma segundos:
integrar apenas cierra cada corrida no ahorra tiempo perceptible frente a integrar todas al final.

Además, el despacho por finalización individual devuelve el control de la fase tantas veces como tareas
tenga el batch, lo que choca con el modelo de "una fase, una vuelta de control" que sostiene el resto del
pipeline.

## Decisión

El despacho de un batch paralelo ocurre en **un solo mensaje** del hilo principal, con las N tareas
marcadas; la integración arranca recién **cuando el batch entero devolvió el control**, y procesa las
corridas en orden ascendente de id de tarea.

## Reglas verificables

- **[manual]** El despacho de un batch con tareas marcadas independientes ocurre en un único mensaje, nunca uno por tarea.
- **[manual]** Ninguna integración arranca antes de que el batch entero haya devuelto el control.
- **[manual]** El orden de integración es ascendente por id de tarea y queda legible en el artefacto de progreso del cambio sin tener que re-derivarlo de la rama.

## Alternativas consideradas

- **Despacho background por finalización individual.** Descartado: devuelve el control de la fase una vez por tarea en vez de una vez por batch, contradiciendo el modelo de retorno único; y el ahorro de reloj de pared frente a integrar todo al cierre es despreciable, porque el batch está dominado por su miembro más lento, no por el momento de cada merge.

## Consecuencias

- La atribución de un conflicto sigue siendo legible por tarea, porque cada una aterriza en su propio commit — integrar diferido no la debilita.
- El orden de integración queda determinado de antemano (ascendente por id), así que consolidar no requiere ninguna decisión adicional en el momento.
- Una corrida que termina rápido espera igual a que termine la más lenta antes de que su trabajo llegue a la rama — la ganancia de paralelismo está en la implementación, no en el aterrizaje.

## Relacionados

- `relacionado-con` → [consolidation-run-is-the-integrator.md](consolidation-run-is-the-integrator.md) — quién ejecuta la integración que esta decisión ordena diferir.
- `relacionado-con` → [phase-fanout-two-cases.md](phase-fanout-two-cases.md) — la prosa que nombra este mecanismo como el segundo caso de fan-out.
- `relacionado-con` → [../contracts/worktree-base-handshake.md](../contracts/worktree-base-handshake.md) — la corrida de consolidación re-chequea la base antes de confiar en cada commit que integra.
