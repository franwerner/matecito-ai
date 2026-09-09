# EDR — Los estados que dependen de los decision records se chequean sobre el retorno del brief, no en un gate propio

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
Los estados `blocked`/`needs-decision` que dependen de si los decision records están activos solían enrutarse a través del INTAKE GATE, que fue eliminado. El brief mismo sigue cargando esos estados, así que necesitaban un lugar nuevo donde chequearse.

## Decisión
El chequeo de esos estados pasa a ser parte de leer el retorno del brief en la sección `### Execution`: si el brief vuelve `status: blocked` (conflicto con un decision record Accepted) o `status: needs-decision` (decisión arquitectónica sin resolver), el orquestador actúa ahí mismo, sin un gate propio intermedio.

## Reglas verificables
- **[manual]** `payload/core/CLAUDE.md`, sección `### Execution`, chequea `status: blocked` / `status: needs-decision` directamente sobre el retorno de `sdd-intake`, sin pasar por un gate propio.

## Alternativas consideradas
Reintroducir un gate dedicado sólo para estos dos estados — descartada: duplica el trabajo que la Sección de Execution ya hace al leer el retorno del brief, por un caso que ya está cubierto ahí.

## Consecuencias
Un solo lugar del kernel lee el status del brief para todo — incluidos estos dos casos condicionados a que los decision records estén activos.
