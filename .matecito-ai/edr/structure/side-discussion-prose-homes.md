# EDR — La prosa del mecanismo de discusión lateral se reparte entre el kernel y una referencia compartida

- **Status:** Accepted
- **Date:** 2026-08-14

## Contexto
El mecanismo de discusión lateral tiene dos lectores — el orquestador y la sesión lateral misma — y uno de ellos no pertenece a ningún dominio: la sesión lateral no es una fase ni un ejecutor de dominio, es una segunda sesión de Claude Code cualquiera. Meter todo en un único lugar obliga a elegir un hogar que no sirve a los dos lectores por igual.

## Decisión
La política domain-neutral (qué es una discusión lateral, que sólo el usuario la abre, que el comportamiento del principal depende del tipo declarado, que la conclusión vuelve por el artifact store) vive en `payload/core/CLAUDE.md`, en la zona del orquestador, como `### Side Discussion (opt-in)`. El mecanismo concreto (el traspaso, el comando de lanzamiento, la conclusión, el formato de clave, el pickup, el abandono) vive en una referencia nueva del tier compartido, `payload/shared/references/side-discussion.md`, desplegada a `~/.claude/references/side-discussion.md`.

## Reglas verificables
- **[manual]** La política general del mecanismo de discusión lateral (qué es, que sólo el usuario la abre, que el comportamiento del principal depende del tipo declarado, que la conclusión vuelve por el artifact store) vive en `payload/core/CLAUDE.md`, sección `### Side Discussion (opt-in)`, sin nombrar el formato del traspaso ni la clave de Engram.
- **[manual]** El mecanismo concreto (traspaso, comando de lanzamiento, conclusión, formato de clave, pickup, abandono) vive en `payload/shared/references/side-discussion.md`, desplegado a `~/.claude/references/side-discussion.md` — nunca en un fragmento de dominio.

## Alternativas consideradas
Todo en el kernel, con la forma de `### Change Workspace (opt-in)` — descartado: un formato que se lee sólo cuando se abre una discusión no merece residencia permanente en un archivo que está siempre en contexto. Todo en `payload/domains/development/references/` — descartado: la sesión lateral no pertenece a ningún dominio, y esa referencia no se despliega cuando development está inactivo.

## Consecuencias
Dos archivos en vez de uno, pero ninguno de los dos duplica lo que el otro ya dice: el kernel dice qué es y cuándo, la referencia dice cómo. El kernel no crece con el detalle de un mecanismo de bajo uso.

## Relacionados
- `relacionado-con` → [change-workspace-prose-homes.md](change-workspace-prose-homes.md) — mismo criterio de reparto (política neutral en el kernel, mecanismo al lado), aplicado antes a otro mecanismo del orquestador.
