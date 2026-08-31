# EDR — Las respuestas de discovery viajan verbatim en el retorno de sdd-explore, no en el brief de sdd-intake

- **Status:** Accepted
- **Date:** 2026-08-31

## Contexto
Antes de este cambio, el ciclo de discovery vivía en sdd-intake y sus respuestas se persistían en la sección `### Discovery answers` del brief. Con el ciclo movido a sdd-explore (porque puede leer el código antes de preguntar), sdd-intake pasa a ser una sola pasada que retorna antes de que exploración corra — no hay una segunda pasada de intake donde esas respuestas puedan aterrizar.

## Decisión
`### Discovery answers` se declara en la variante `[done, blocked]` de sdd-explore.yaml (`emitted: always`, `sentinel: true`, `render: items`), y se elimina del contrato de sdd-intake junto con el resto de su rol de discovery. El artefacto que persiste sdd-explore (`sdd/{change-name}/explore`) es el mismo bloque que retorna, así que las respuestas quedan autocontenidas ahí.

## Reglas verificables
- **[auto]** Un retorno `done`/`blocked` de sdd-explore sin el campo `discovery_answers` falla el render — la sección es `emitted: always`; la lista vacía renderiza el sentinel `None.`, nunca se omite el campo.
- **[manual]** sdd-intake.yaml ya no declara ninguna sección `### Discovery answers` — verificado leyendo el contrato colapsado; nada mecánico impide reagregarla ahí en el futuro.

## Alternativas consideradas
(a) no persistir las respuestas en ningún artefacto, dejarlas sólo en la conversación del orquestador — descartado porque rompe la garantía de nunca parafrasear una respuesta que la fase de discovery ya sostenía. (b) mantener la sección en el brief de intake, completada en una pasada posterior — descartado porque intake ahora es una sola pasada que retorna antes de que exploración corra.

## Consecuencias
El brief de intake pierde `### Discovery answers`; quien busque las respuestas del usuario a partir de ahora las encuentra en el artefacto de exploración, no en el de intake.
