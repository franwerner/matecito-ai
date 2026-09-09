# EDR — No hay `/sdd-*` slash-commands; una fase se despacha, no se invoca por comando

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
La sección listaba `/sdd-init`, `/sdd-intake`, … como comandos de usuario, y no existen: no hay
`~/.claude/commands/` en el deploy y las once skills `sdd-*` llevan `disable-model-invocation: true`
junto con `user-invocable: false`, que cierran las dos puertas. Un usuario que tipeaba `/sdd-verify` no
encontraba nada, y la lista además omitía cuatro fases (propose, spec, design, tasks), siendo `spec`
obligatoria.

## Decisión
`### How a phase runs` documenta la vía que realmente existe: no hay slash-commands, las skills de fase
no son auto-invocables por el modelo, y una fase corre porque el orquestador despacha su agente. Pedir
el trabajo en lenguaje llano dispara el pipeline; nombrar una fase es un pedido de despacho, no un
comando que el harness resuelve.

## Reglas verificables
- **[manual]** `payload/domains/development/CLAUDE.md`, sección `### How a phase runs`, no lista ningún `/sdd-*` como comando de usuario y declara que una fase corre por despacho del orquestador.

## Alternativas consideradas
Documentar los `/sdd-*` como comandos aspiracionales, a implementar después — descartada: describe algo
que no existe en el deploy actual, y un usuario que los tipea no encuentra nada.

## Consecuencias
La documentación coincide con lo que el deploy realmente ofrece: pedir el trabajo en lenguaje llano.
