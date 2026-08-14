# EDR — El ancla de un ítem de gate es de forma libre, sólo se chequea su presencia

- **Status:** Accepted
- **Date:** 2026-08-14

## Contexto
El ancla que acompaña a un ítem de gate señala la fuente concreta de la que ese ítem surgió: un archivo, una línea, una clave de artefacto. Las formas legítimas son abiertas por construcción — una ruta de archivo, una clave de memoria persistente, la línea de un artefacto que todavía no existe — y la existencia real de lo que el ancla nombra no es algo que un chequeo mecánico pueda confirmar.

## Decisión
El ancla se acepta en la forma en que la fase la escriba; el único chequeo es que el token esté presente, igual que cualquier otro token declarado sin un conjunto fijo de valores legales. No se valida la forma del valor.

## Reglas verificables
- **[auto]** Un ítem sin el token de ancla falla la validación por ausencia, igual que cualquier otro token requerido.
- **[manual]** Ningún patrón de forma se aplica sobre el valor del ancla una vez que está presente.

## Alternativas consideradas
Validar la forma del ancla con una expresión regular — descartado porque las formas legítimas son abiertas por construcción; un patrón lo bastante estricto para rechazar basura también rechaza anclas reales, y chequear la sintaxis da la falsa impresión de que también se chequeó la existencia.

## Consecuencias
El ancla queda tan flexible como la fuente que señala, sin falsos rechazos por forma — a costa de que un ancla presente pero inválida (apuntando a algo que no existe) no se detecta mecánicamente; sólo su ausencia se atrapa.

## Relacionados
- `relacionado-con` → [gate-summary-cap-in-characters.md](gate-summary-cap-in-characters.md) — mismo mecanismo de ítems de gate, mismo cambio que lo introdujo.
- `relacionado-con` → [contract-item-anchors-once.md](contract-item-anchors-once.md) — extiende este mismo token a un ítem con más de un campo.
