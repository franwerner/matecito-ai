# EDR — Un ítem reaparecido se reconoce por igualdad exacta de string sobre su identificador

- **Status:** Accepted
- **Date:** 2026-08-14

## Contexto
Un ítem ya ratificado en un gate anterior del mismo cambio puede volver a aparecer en un gate posterior — por ejemplo, si una fase distinta lo vuelve a proponer. Reconocer esa reaparición necesita un criterio de comparación único y determinístico, no una heurística de similitud.

## Decisión
La reaparición se detecta comparando el identificador del ítem entrante contra el registro, con igualdad de string exacta. No hay comparación por similitud de resumen ni un criterio auxiliar sobre el ancla.

## Reglas verificables
- **[auto]** La comparación de reaparición es igualdad exacta de string sobre el identificador del ítem, sin normalización ni resolución difusa.

## Alternativas consideradas
Comparar por identificador más una pista adicional del ancla — descartada porque esa pista termina decidiendo por su cuenta (convirtiéndose en una segunda regla difusa) o no aporta nada si nunca decide. Comparar por similitud del resumen — descartada porque un resumen puede cambiar de redacción sin cambiar de decisión, y la similitud no es determinística.

## Consecuencias
La detección de reaparición es exacta y predecible — a costa de que un mismo ítem con un identificador distinto (por ejemplo, por un error de tipeo en una fase) no se reconoce como el mismo.

## Relacionados
- `relacionado-con` → [ratification-ledger-row.md](ratification-ledger-row.md) — la fila que esta comparación consulta.
- `relacionado-con` → [re-emergence-short-form.md](re-emergence-short-form.md) — qué se muestra cuando la comparación encuentra una coincidencia.
