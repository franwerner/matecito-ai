# EDR — La lista de gaps del mine gate no lleva ningún hint de `## Alcance` del artefacto de tareas

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
El alcance que se pasaba al ejecutor de minería incluía "cualquier hint de `## Alcance` del artefacto de tareas". Dos problemas encadenados: `## Alcance` es una sección del TEMPLATE del EDR, y el artefacto `tasks` nunca la tuvo. Y aunque la tuviera, no podría ayudar — un gap es por definición un EDR que no existe todavía, así que no hay ningún `## Alcance` real con el cual dar la pista.

## Decisión
La lista de gaps que se pasa como `scope` al ejecutor de minería lleva sólo tres cosas por ítem: `domain/slug`, la tarea que lo implementó, y la raíz del repo — minado contra el trabajo ya entregado. Ningún hint de `## Alcance` del artefacto de tareas se incluye.

## Reglas verificables
- **[manual]** `payload/core/CLAUDE.md`, sección `### Decision-Gap Capture (mine gate)`, construye el scope del ejecutor con exactamente `domain/slug` + tarea implementadora + raíz del repo, sin ningún hint de `## Alcance`.

## Alternativas consideradas
Mantener el hint de `## Alcance` del artefacto de tareas — descartada: la sección no existe en ese artefacto, y aunque existiera, un gap por definición no tiene un EDR (y por lo tanto ningún `## Alcance`) con qué darla.

## Consecuencias
El ejecutor de minería tiene menos contexto de entrada por diseño, pero el que tenía era una promesa vacía — nunca hubo datos reales detrás de ese campo.
