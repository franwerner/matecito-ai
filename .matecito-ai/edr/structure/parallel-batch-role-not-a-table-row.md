# EDR — La fila de `sdd-apply` en Phase Read/Write es el ideal single-dispatch; el batch paralelo no gana una fila propia

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
La fila de `sdd-apply` en `### SDD Phase Read/Write` describe el ideal full-lane, single-dispatch. Un
batch con tareas marcadas por independencia lo divide entre los dos roles de fan-out declarados
(isolated run, consolidation run) en vez de cambiar esa fila — cada rol ya tiene su propia fila
inmediatamente debajo.

## Decisión
La fila de `sdd-apply` en la tabla no se reescribe para reflejar el batch paralelo; ese caso se documenta
en las dos filas dedicadas (`↳ isolated run`, `↳ consolidation run`) que ya existen debajo, sin duplicar
la explicación en una nota aparte sobre la fila original.

## Reglas verificables
- **[manual]** `payload/domains/development/CLAUDE.md`, sección `### SDD Phase Read/Write`, mantiene la fila de `sdd-apply` como el ideal single-dispatch y documenta el batch paralelo únicamente en sus dos filas `↳` dedicadas.

## Alternativas consideradas
Reescribir la fila de `sdd-apply` para cubrir también el caso de batch paralelo — descartada: las dos
filas `↳` ya existen para eso, y duplicar la explicación en la fila original es la misma clase de
duplicación que otras notas de esta tabla ya corrigieron.

## Consecuencias
Quien lee la fila de `sdd-apply` ve el caso ideal; quien necesita el caso de batch paralelo baja a las
dos filas `↳` que ya lo cubren.
