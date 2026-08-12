# EDR — El veredicto design-conflict lo clasifica el orquestador, con una tabla de tres filas espejo de blocking-test

- **Status:** Accepted
- **Date:** 2026-08-12

## Contexto
sdd-apply ahora declara un veredicto `design-conflict` por cada propuesta rechazada que chequeó, pero declarar no es lo mismo que actuar: alguien tiene que leer ese token y decidir qué hace el flujo con él. El patrón ya existente para el token `blocking-test` de sdd-design (clasificación por el orquestador, nunca por el ejecutor re-derivando la prueba) es el precedente directo — dos tokens del mismo mecanismo de captura in-flow, dos lecturas hechas por el mismo participante.

## Decisión
El orquestador clasifica el token en el Unresolved Decisions Guard de `CLAUDE.md`, con una tabla de tres filas espejo de la de `blocking-test`: `none` procede sin gate ni mención; `conflicts` exige que el retorno sea `status: blocked` — un ítem `conflicts` sobre un retorno no-`blocked` es una violación de contrato, no un caso a rutear por el flujo ordinario del guard de mailboxes; y la ausencia de un ítem para un `record:` que el orquestador mismo forwardeó como rechazado, cuya tarea gobernante la corrida alcanzó, no se lee como `none` — "sin pronunciamiento" y "no hay conflicto" son afirmaciones distintas, y confundirlas dejaría pasar una corrida que nunca chequeó lo que debía.

## Reglas verificables
- **[manual]** `domains/development.md` declara la tabla de tres filas ("Reading the `design-conflict` token") inmediatamente después del párrafo de forwarding, mismo patrón que "Reading the `blocking-test` token".
- **[manual]** Un ítem `conflicts` sobre un retorno no `blocked` se trata como violación de contrato — se detiene y se surfacea, no se rutea por el guard ordinario de mailboxes Tier 1/2.
- **[manual]** La fila "sin ítem" es una tercera clasificación explícita, distinta de `none`: la ausencia nunca se lee como veredicto limpio.

## Alternativas consideradas
Que `sdd-verify` clasifique el token en vez del orquestador (como el "Return Contract Check" mecánico): descartado — el token decide si el batch se detiene ANTES de llegar a verify; para cuando verify corre, el código ya estaría escrito sobre una base no ratificada si el orquestador no paró antes. Que `decision-gaps` (el grupo de `sdd-verify`) lo chequee: descartado por el propio diseño — `decision-gaps` corre demasiado tarde para frenar código que ya aterrizó, y el spec excluye explícitamente esa ruta.

## Consecuencias
El mismo participante (el orquestador) es responsable de leer los dos tokens hermanos de este mecanismo de captura in-flow (`blocking-test`, `design-conflict`), con el mismo patrón de lectura — menos superficie nueva que aprender. Un veredicto `conflicts` que se cuela en un retorno no-`blocked` no pasa desapercibido: el orquestador lo trata como una violación de contrato, con la misma severidad que cualquier otra ruptura del Return Contract Check.

## Relacionados
- `relacionado-con` → [../contracts/rejected-proposal-verdict-token.md](../contracts/rejected-proposal-verdict-token.md) — el token que esta decisión clasifica.
- `relacionado-con` → [../contracts/expected-set-outside-the-validator.md](../contracts/expected-set-outside-the-validator.md) — la mitad de completitud que esta misma lectura del orquestador también resuelve.
