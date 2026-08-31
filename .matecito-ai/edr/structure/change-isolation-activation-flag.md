# EDR — La activación del aislamiento por cambio viaja como cuarta línea del brief

- **Status:** Accepted
- **Date:** 2026-08-13

## Contexto
El aislamiento por espacio de trabajo a nivel de cambio (`structure/change-level-worktree-isolation.md`) se activa por elección explícita en el fork de lane, nunca por default — pero el fork de lane y el brief que produce `sdd-intake` ya tienen tres flags de decisión existentes (`diagram`, `ui-test`, `components`), confirmados en un único punto: el INTAKE GATE. Sin un lugar fijo donde viajar, la elección de aislamiento se resolvería con una pregunta ad-hoc en el gate, sin artefacto que la sostuviera después — exactamente el problema que el mecanismo de flags ya resuelve para las otras tres.

## Decisión
La elección de aislamiento viaja como una cuarta línea del brief, `- Isolation: {active|inactive}`, bajo `### Classification`, junto a `Diagram`, `UI test` y `Components`. `sdd-intake` la decide por cuenta del usuario, per `structure/change-isolation-activation-flag.md`: activa solo cuando el pedido la pide explícitamente. El orquestador la reporta junto con el resto de los flags del brief en una única línea de aviso; ningún gate la confirma. Para trabajo `direct`/ad-hoc — que nunca produce un brief — la elección se resuelve en el pedido mismo: un pedido que nunca la pidió implica aislamiento inactivo.

## Reglas verificables
- **[manual]** La elección de aislamiento se agrega como una cuarta línea `- Isolation: {active|inactive}` bajo `### Classification` del brief, junto a `Diagram`, `UI test` y `Components`.
- **[manual]** La elección se decide por `sdd-intake` y se reporta junto con el resto de los flags del brief; ningún gate la confirma y ninguna fase posterior la vuelve a preguntar.
- **[manual]** Para trabajo `direct`/ad-hoc, la elección se resuelve en el pedido mismo; un pedido que nunca la pidió explícitamente implica aislamiento inactivo.

## Alternativas consideradas
Una pregunta ad-hoc en el gate, separada del fork de lane. Descartado: ningún artefacto la sostendría después de ese punto, y el mecanismo de flags de decisión ya resuelve exactamente este problema para `diagram`, `ui-test` y `components` — agregar una cuarta flag reusa el mecanismo en vez de duplicarlo.

## Consecuencias
El brief gana una cuarta línea de clasificación, y el INTAKE GATE gana un flag más que surfacear por nombre con su valor y su razón. El trabajo `direct`/ad-hoc, que no tiene brief, resuelve la misma elección en el fork de lane — un único mecanismo con dos puntos de confirmación, no dos mecanismos.

## Relacionados
- `relacionado-con` → [change-level-worktree-isolation.md](change-level-worktree-isolation.md) — el aislamiento por cambio que esta elección activa.
