# EDR — La discusión lateral vive en su propio namespace de Engram, fuera de `sdd/`, con una clave por mitad

- **Status:** Accepted
- **Date:** 2026-08-14

## Contexto
El mecanismo puede pasar con o sin un cambio activo del flujo, así que la clave no puede depender de un `{change-name}`. El traspaso y la conclusión, además, los escribe cada uno un actor distinto — el principal escribe el traspaso, la sesión lateral escribe la conclusión.

## Decisión
Dos claves, un namespace propio fuera de `sdd/`: `side-discussion/{slug}/handoff` y `side-discussion/{slug}/conclusion`, donde `{slug}` es un kebab-case del tema que el orquestador fija al componer el traspaso. Declarado una sola vez, en la referencia compartida — no se agrega a la tabla de namespace de `development`, porque ninguna fase SDD produce o lee ninguna de las dos claves.

## Reglas verificables
- **[manual]** El traspaso se guarda en Engram bajo `side-discussion/{slug}/handoff`; la conclusión bajo `side-discussion/{slug}/conclusion`.
- **[manual]** El slug es un kebab-case del tema y lo fija el orquestador al componer el traspaso — nunca lo elige la sesión lateral.

## Alternativas consideradas
Anidar bajo `sdd/{change-name}/side-discussion/{slug}` — descartado: se rompe en cuanto la discusión ocurre sin un cambio activo (trabajo directo o ad-hoc). Una sola clave para ambas mitades — descartado: hay dos escritores distintos, y un upsert de Engram sobre una clave compartida forzaría a la sesión lateral a reescribir el texto del orquestador para poder agregar el suyo.

## Consecuencias
El mecanismo funciona igual con o sin cambio activo. Es territorio nuevo — ningún namespace fuera de `sdd/` había tenido antes un intercambio determinístico actor-A-escribe / actor-B-lee — así que la convención queda declarada explícitamente en vez de inferida del patrón `sdd/`.

## Relacionados
- `relacionado-con` → [side-discussion-prose-homes.md](side-discussion-prose-homes.md) — fija dónde vive la prosa que declara este namespace.
