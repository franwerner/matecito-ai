# EDR — El traspaso de discusión lateral lleva cinco secciones fijas más una línea de tipo

- **Status:** Accepted
- **Date:** 2026-08-14

## Contexto
El traspaso es lo único que le da contexto a una sesión lateral que arranca en blanco. Prosa libre fue rechazada en intake porque cada traspaso saldría distinto y nada garantizaría que las referencias necesarias viajen. Por separado, el principal necesita saber si la discusión es bloqueante o consultiva incluso después de una compactación o al retomar la sesión.

## Decisión
El traspaso siempre lleva cinco secciones — `## Topic` / `## What I already read` / `## References` / `## Open question` / `## Return` — todas llenadas por el orquestador en cada traspaso; una vacía se declara vacía, nunca se omite. Encima de esas cinco va una línea de encabezado, `- **Type:** blocking | consultive`, que no es una sexta sección: gobierna qué hace el principal, no lo que la sesión lateral tiene que leer.

## Reglas verificables
- **[manual]** El traspaso lleva siempre las cinco secciones `## Topic`, `## What I already read`, `## References`, `## Open question` y `## Return`, llenadas por el orquestador en cada traspaso; una sección sin contenido se declara vacía explícitamente, nunca se omite.
- **[manual]** Encima de `## Topic` va una línea de encabezado `- **Type:** blocking | consultive`; no cuenta como una sexta sección y no se pliega dentro de `## Topic` ni de `## Return`.

## Alternativas consideradas
Prosa libre — descartada en intake. Las cuatro secciones del brief, sin `## Return` — descartada en diseño: la sesión lateral arranca sin contexto y necesita que se le diga adónde devolver la conclusión. Para el tipo específicamente: una sexta sección — descartada, una sección entera para una palabra es el instrumento equivocado; plegarlo dentro de `## Topic` o `## Return` — descartado, mezclaría un dato que gobierna al principal dentro de una sección dirigida a la sesión lateral.

## Consecuencias
La sesión lateral siempre sabe adónde escribir su conclusión, sin ambigüedad. El tipo queda registrado de forma durable en el propio traspaso, así que una sesión retomada o compactada lo recupera sin tener que preguntarle de nuevo al usuario.

## Relacionados
- `relacionado-con` → [ratification-ledger-row.md](ratification-ledger-row.md) — mismo criterio — campos fijos en vez de prosa libre, porque un lector posterior tiene que poder encontrar cada parte sin interpretar.
- `relacionado-con` → [side-discussion-conclusion-shape.md](side-discussion-conclusion-shape.md) — la otra mitad del intercambio, con el mismo criterio de secciones fijas.
