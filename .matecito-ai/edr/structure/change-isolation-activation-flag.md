# EDR — La activación del aislamiento por cambio viaja como cuarta línea del brief

- **Status:** Accepted
- **Date:** 2026-08-13

## Contexto
El aislamiento por espacio de trabajo a nivel de cambio (`structure/change-level-worktree-isolation.md`) se activa solo cuando el pedido lo pide explícitamente, nunca por default — pero el brief que produce `sdd-intake` ya tiene tres flags de decisión existentes (`diagram`, `ui-test`, `components`), que `sdd-intake` decide por cuenta del usuario y el orquestador reporta en una única línea de aviso, sin que ningún gate las confirme. Sin un lugar fijo donde viajar, la elección de aislamiento se resolvería con una pregunta ad-hoc, sin artefacto que la sostuviera después — exactamente el problema que el mecanismo de flags ya resuelve para las otras tres.

## Decisión
La elección de aislamiento viaja como una cuarta línea del brief, `- Isolation: {active|inactive}`, bajo `### Classification`, junto a `Diagram`, `UI test` y `Components`. `sdd-intake` la decide por cuenta del usuario, per `structure/change-level-worktree-isolation.md`: activa solo cuando el pedido la pide explícitamente. El orquestador la reporta junto con el resto de los flags del brief en una única línea de aviso; ningún gate la confirma como ítem propio — viaja dentro del brief que el Brief Confirmation Gate confirma entero. Para trabajo `direct`/ad-hoc — que nunca produce un brief — la elección se resuelve en el pedido mismo: un pedido que nunca la pidió implica aislamiento inactivo.

## Reglas verificables
- **[manual]** La elección de aislamiento se agrega como una cuarta línea `- Isolation: {active|inactive}` bajo `### Classification` del brief, junto a `Diagram`, `UI test` y `Components`.
- **[manual]** La elección se decide por `sdd-intake` y se reporta junto con el resto de los flags del brief; ningún gate la confirma como ítem propio, y ninguna fase posterior la vuelve a preguntar.
- **[manual]** Para trabajo `direct`/ad-hoc, la elección se resuelve en el pedido mismo; un pedido que nunca la pidió explícitamente implica aislamiento inactivo.

## Alternativas consideradas
Una pregunta ad-hoc en el gate, separada del fork de lane. Descartado: ningún artefacto la sostendría después de ese punto, y el mecanismo de flags de decisión ya resuelve exactamente este problema para `diagram`, `ui-test` y `components` — agregar una cuarta flag reusa el mecanismo en vez de duplicarlo.

## Consecuencias
El brief gana una cuarta línea de clasificación, y el orquestador gana un flag más que reportar por nombre, con su valor y su razón, en la misma línea de aviso que ya cubre los otros tres — sin que nadie la confirme. El trabajo `direct`/ad-hoc, que no tiene brief, resuelve la misma elección en el pedido mismo — un único mecanismo con dos puntos de resolución, no dos mecanismos. Un flag mal decidido llega a `sdd-apply` sin que nadie lo haya revisado, el mismo costo que `structure/change-level-worktree-isolation.md` ya asume para la activación que esta elección dispara.

## Relacionados
- `relacionado-con` → [change-level-worktree-isolation.md](change-level-worktree-isolation.md) — el aislamiento por cambio que esta elección activa.
