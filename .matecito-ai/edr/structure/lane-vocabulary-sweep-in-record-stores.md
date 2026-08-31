# EDR — El sweep de vocabulario de lane llega a dos entradas de índice, no al tercer residuo ya aceptado

- **Status:** Accepted
- **Date:** 2026-08-31

## Contexto
two-lanes-fixed-flow retira el modelo de cuatro lanes (direct/reduced/full/custom) y el fork/INTAKE GATE que los seleccionaba, reemplazándolos por dos lanes fijos sin gate de confirmación de scope. Dos filas de navegación en los índices de los stores durables seguían nombrando el mecanismo retirado por su nombre: `.matecito-ai/edr/structure/INDEX.md:23`, la celda "Consultá cuando..." de `change-isolation-activation-flag.md`, y `.matecito-ai/development-specs/rule/INDEX.md:34`, la celda de descripción de `intake-brief-components-line.md`. La tarea 6.14 necesitaba decidir si el sweep final del cambio alcanza esas dos celdas de navegación, y si alcanza además a una tercera candidata, `.matecito-ai/edr/contracts/side-discussion-conclusion-is-not-a-record.md:10`, que sigue leyendo "se ratifica una vez en el gate del lane".

## Decisión
El sweep alcanza las dos celdas de navegación y deja la tercera intacta. IN: `.matecito-ai/edr/structure/INDEX.md:23` y `.matecito-ai/development-specs/rule/INDEX.md:34` son filas de navegación de records que este mismo cambio reescribe — su celda "leé esto cuando..." apuntaba a un fork y a un gate que ya no existen, así que quedarían apuntando a algo que no resuelve si no se reescribían. Ambas se reescriben para no nombrar el mecanismo retirado. OUT: `.matecito-ai/edr/contracts/side-discussion-conclusion-is-not-a-record.md:10` queda tal cual — la regla que enuncia (una conclusión de discusión lateral que zanja una pregunta arquitectónica entra al camino existente del flujo, se ratifica una vez y `sdd-apply` materializa el EDR) sigue siendo correcta palabra por palabra salvo por el sustantivo "el gate del lane", que ahora es una etiqueta vieja para el mismo mecanismo (un único gate por fase, sin fork de lane). Ambos criterios ya estaban ratificados antes de esta tarea, no se re-deciden acá: el criterio residuo-vs-defecto viene de `structure/retired-vocabulary-in-record-stores.md`, y el criterio de cuándo un puntero stale sí entra al sweep viene de `structure/gating-vocabulary-sweep-scope.md`.

## Alcance
- `.matecito-ai/edr/structure/INDEX.md` — Celda de navegación reescrita.
- `.matecito-ai/development-specs/rule/INDEX.md` — Celda de navegación reescrita.

## Reglas verificables
- **[manual]** `.matecito-ai/edr/structure/INDEX.md:23` no nombra "fork de lane" ni construcciones equivalentes; describe la activación del aislamiento sin mecanismo de lane.
- **[manual]** `.matecito-ai/development-specs/rule/INDEX.md:34` no nombra "INTAKE GATE"; describe `intake-brief-components-line.md` como decidido y reportado por `sdd-intake`, nunca ratificado.
- **[manual]** `.matecito-ai/edr/contracts/side-discussion-conclusion-is-not-a-record.md:10` queda sin editar por este cambio — su sustantivo "el gate del lane" es residuo aceptado, no un defecto que el sweep final deba señalar.

## Alternativas consideradas
Barrer también el tercer archivo, reescribiendo "el gate del lane" a algo neutral. Descartada por el mismo criterio que `structure/retired-vocabulary-in-record-stores.md` ya fijó para casos análogos: es un EDR de razonamiento cuyo sustantivo incidental no cambia lo que manda, así que editarlo acá sería scope creep sin beneficio funcional — la regla que enuncia sigue siendo correcta con o sin la palabra "lane". Extender el sweep sólo a una de las dos celdas de navegación y dejar la otra. Descartada: ambas apuntan a records que este mismo cambio reescribe, así que dejar una sin tocar deja al lector siguiendo un puntero muerto exactamente en el mismo lugar donde el cambio ya está reescribiendo el resto de la fila.

## Consecuencias
Después de este cambio, las dos celdas de navegación de los stores durables leen correctamente contra el mecanismo post-cambio (dos lanes fijos, sin fork, sin INTAKE GATE). Un tercer archivo, `side-discussion-conclusion-is-not-a-record.md:10`, queda como residuo aceptado explícito — nombrado acá para que no se vuelva a señalar como hallazgo en un sweep futuro. Ningún chequeo mecánico impide que ese archivo divergiera más adelante; un futuro cambio sobre `New Decisions` o el mecanismo de discusiones laterales es el punto natural para reescribirlo, no éste.

## Relacionados
- `edr` → [structure/retired-vocabulary-in-record-stores.md](structure/retired-vocabulary-in-record-stores.md) — Fija el criterio residuo-vs-defecto que esta decisión aplica sin re-decidir.
- `edr` → [structure/gating-vocabulary-sweep-scope.md](structure/gating-vocabulary-sweep-scope.md) — Fija el criterio de cuándo un puntero stale en un archivo leído en runtime entra al sweep.
