# EDR — `sdd-spec` lee el intake brief siempre, no sólo como upstream de fallback

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
`sdd-spec` pasó a leer el intake brief SIEMPRE, no sólo como upstream de fallback: es el único lugar que
lleva el flag `ui-test`, y la proposal no lo transporta. Bajo el modelo de cuatro lanes
(direct/reduced/full/custom) — retirado por `two-lanes-fixed-flow`, que lo reemplazó por dos lanes fijos
sin fork ni gate de confirmación — sin esa lectura la producción de `ui-scenarios` habría funcionado en
lane `reduced` y fallado en `full`.

## Decisión
La fila de `sdd-spec` en `### SDD Phase Read/Write` declara la lectura del intake brief como
incondicional (`always, for the ui-test flag`), no como un fallback que sólo se consulta cuando la
proposal falta.

## Reglas verificables
- **[manual]** `payload/domains/development/CLAUDE.md`, sección `### SDD Phase Read/Write`, marca la lectura del intake brief por `sdd-spec` como siempre, para el flag `ui-test`.

## Alternativas consideradas
Dejar la lectura del intake brief como fallback condicional a la ausencia de proposal — descartada: bajo
el modelo de lanes anterior producía `ui-scenarios` en una lane y no en otra, un comportamiento que
dependía silenciosamente de cuál upstream existiera.

## Consecuencias
`sdd-spec` produce `ui-scenarios` de forma consistente independientemente de qué upstream exista, porque
el flag que lo decide siempre se lee.
