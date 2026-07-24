# EDR — Health checks

- **Status:** Accepted
- **Type:** decision
- **Date:** 2026-07-23

## Contexto

El broker es una dependencia dura del MCP (que hace un chequeo de arranque contra él) y de la UI. Necesita un chequeo que confirme que el proceso está vivo y que la persistencia está lista. Al ser un daemon local sin orquestador ni pool de instancias, no hay valor en separar liveness de readiness.

## Decisión

Exponemos un único endpoint de liveness. Responde 200 cuando el broker está arriba y la persistencia está abierta y migrada (chequeo local, sin tocar dependencias externas); el body incluye lo mínimo para diagnóstico (estado del proceso y estado del store). No hay split liveness/readiness porque no hay orquestador ni pool. Sirve al chequeo de arranque del MCP y a la UI.

## Reglas verificables

- **[manual]** el endpoint de health responde 200 solo con el broker vivo y la persistencia abierta y migrada, sin tocar dependencias externas.
- **[tool: test]** existe un test que verifica que el endpoint de health refleja proceso vivo más estado de la persistencia.

## Alternativas consideradas

- **Split liveness/readiness:** descartado; sin orquestador ni pool no aporta, un solo endpoint alcanza.

## Consecuencias

- El MCP y la UI tienen una señal simple y local de disponibilidad.
- Trade-off: al no separar readiness, no hay una señal distinta para "vivo pero aún no listo"; en un daemon local sin orquestación no hace falta.
