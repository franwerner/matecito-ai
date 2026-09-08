# EDR — REFACTOR.md §2 se registra como frente cerrado, no se borra

- **Status:** Accepted
- **Date:** 2026-09-08

## Contexto
`REFACTOR.md` §2 argumentaba por adoptar qmd como buscador y tenía una pregunta abierta (“El punto a resolver antes de tocar nada”) sobre la relación de qmd con `process/index-decision-records`. Desde entonces pasaron dos cosas: qmd quedó adoptado (`process/configure-record-search` y `flow/find-durable-records`, ambos `Accepted` el 2026-09-02), y `index-decision-records` — la contraparte de esa pregunta, junto con los cuatro EDRs de `apps/api` de los que dependía — se retira con el resto del cockpit en este mismo cambio (`sdd/drop-apps-subtree`).

## Decisión
§2 se reescribe en el lugar para registrar que su punto se cerró —no que se respondió— porque su contraparte dejó de existir; el ítem asociado (§5 ítem 1) sale de “Decisiones abiertas” en vez de borrarse o volver a plantearse como pregunta nueva. La sección queda, con su historia, en vez de irse junto con el código sobre el que argumentaba.

## Reglas verificables
- **[manual]** `REFACTOR.md` §2 afirma que qmd está adoptado (no propuesto), citando `process/configure-record-search` y `flow/find-durable-records`.
- **[manual]** El cierre de `REFACTOR.md` §2 dice que la pregunta se cerró porque su contraparte (`process/index-decision-records`) se retiró junto con el cockpit, no que se resolvió por sus méritos.
- **[manual]** `REFACTOR.md` §5 (“Decisiones abiertas”) lista exactamente dos ítems, renumerados desde los tres previos al cambio.

## Alternativas consideradas
Borrar §2 y su ítem de decisión abierta directamente. Descartada: pierde el registro de por qué se cerró el frente, que es justamente para lo que existe un documento de trabajo como `REFACTOR.md`. También se consideró refundar §2 como una pregunta abierta nueva (qmd contra la navegación curada a mano de `INDEX.md`) — descartada: inventa una bifurcación que nadie pidió.

## Consecuencias
`REFACTOR.md` conserva un rastro legible de un debate resuelto en vez de perderlo en silencio; quién lea §5 después ve dos ítems, no tres, sin ninguna referencia colgante a una capability retirada.
