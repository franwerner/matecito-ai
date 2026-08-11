# EDR — La corrida de consolidación es el único escritor del artefacto y del estado de las tareas

- **Status:** Accepted
- **Date:** 2026-08-10

## Contexto

Un batch paralelo produce N corridas concurrentes, y el artefacto de progreso del cambio junto con el
estado de las tareas necesitan quedar en un estado consistente al cierre del batch. Engram resuelve por
`topic_key` con la última escritura ganando — si dos corridas escribieran el mismo artefacto, una pisaría
a la otra sin que quede rastro de cuál prevaleció.

## Decisión

El artefacto de progreso del cambio y el estado de las tareas tienen **exactamente un escritor por
batch**: la corrida de consolidación. Las corridas aisladas pueden leer ese artefacto si lo necesitan,
pero nunca lo escriben ni marcan una tarea por su cuenta.

## Reglas verificables

- **[manual]** Solo la corrida de consolidación llama a las operaciones de escritura sobre el artefacto de progreso y el artefacto de tareas del cambio.
- **[manual]** Ninguna corrida aislada persiste su propio artefacto de progreso ni marca una tarea como completa.

## Alternativas consideradas

- **Una clave de Engram por tarea.** Descartado: multiplica namespaces que después hay que deduplicar y limpiar al consolidar, para resolver un problema que un único escritor con una única clave ya resuelve sin coordinación adicional.

## Consecuencias

- No hace falta ningún mecanismo de coordinación entre escrituras concurrentes — no las hay, por construcción.
- Todo el contenido que una corrida aislada produce (archivos tocados, forks, desvíos, contrapartes de UI) viaja hacia el escritor único a través de su Task Run Report, nunca directamente a Engram.
- Si la corrida de consolidación falla antes de escribir, el batch entero queda sin persistir — el costo de tener un solo punto de escritura.

## Relacionados

- `relacionado-con` → [../structure/consolidation-run-is-the-integrator.md](../structure/consolidation-run-is-the-integrator.md) — quién es ese escritor único.
- `relacionado-con` → [worktree-base-handshake.md](worktree-base-handshake.md) — lo que el escritor único re-verifica antes de integrar cada reporte.
