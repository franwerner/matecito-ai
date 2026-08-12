# EDR — El veredicto de una propuesta rechazada sobrevive el fan-out, en el envelope del Task Run Report

- **Status:** Accepted
- **Date:** 2026-08-12

## Contexto
Un batch paralelo de sdd-apply fanea la implementación en N corridas aisladas, cada una devolviendo un Task Run Report en vez del bloque `## Implementation Progress`; la corrida de consolidación es la única que escribe el retorno final. Si el veredicto `design-conflict` de una propuesta rechazada sólo viajara en el retorno serial, una corrida aislada no tendría dónde declararlo, y el veredicto moriría en el límite del worktree — la consolidación no tendría manera de saber que una tarea quedó sin implementar por un conflicto de contenido en vez de, digamos, un handshake de base fallido.

## Decisión
`### Rejected Proposals Checked` se agrega al envelope del Task Run Report (mismo formato de ítem que en el retorno serial), y la corrida de consolidación lo copia hacia el `## Implementation Progress` consolidado, verbatim, sin re-juzgar el veredicto. Un veredicto `conflicts` obliga a la corrida aislada a reportar `Result: not-implemented` para esa tarea — nunca `committed` — así que la tarea nunca llega a `[x]` y el batch nunca cierra `done` mientras ese ítem exista sin resolver en un batch posterior.

## Reglas verificables
- **[manual]** `parallel-batch.md` declara `### Rejected Proposals Checked` como sección condicional del Task Run Report, misma forma que el retorno serial.
- **[manual]** La corrida de consolidación copia cada ítem de `### Rejected Proposals Checked` de cada reporte al retorno consolidado, sin volver a correr la guardia de contenido — "never re-judging" es una prohibición explícita del texto, no sólo una descripción.
- **[manual]** Un veredicto `conflicts` fuerza `Result: not-implemented` en el Task Run Report de esa tarea; la consolidación nunca reporta `status: done` mientras un ítem `conflicts` exista sin una tarea `[x]` que lo resuelva.

## Alternativas consideradas
Que la corrida de consolidación re-derive el veredicto releyendo el diseño y la propuesta rechazada: descartado — la corrida de consolidación no implementa nada, no tiene el contexto de la tarea individual que la corrida aislada sí tiene, y re-juzgar duplicaría trabajo y podría producir un veredicto distinto del que realmente decidió si commitear o no.

## Consecuencias
El veredicto de una propuesta rechazada no se pierde en el límite del worktree: llega al retorno consolidado exactamente igual que en modo serial. La consolidación es un mero copista para esta sección — no reevalúa nada, lo que mantiene la invariante de "single writer, no re-judging" del resto del mecanismo de batch paralelo.

## Relacionados
- `relacionado-con` → [rejected-proposal-verdict-token.md](rejected-proposal-verdict-token.md) — el token que este ítem transporta a través del fan-out.
- `relacionado-con` → [verdict-in-its-own-conditional-section.md](verdict-in-its-own-conditional-section.md) — la sección, en el retorno serial, que este mecanismo replica en el Task Run Report.
