# EDR — El chequeo de backing de decision-gaps se declara no-aplicable sobre una fila `modified`, en vez de pasarla o fallarla

- **Status:** Accepted
- **Date:** 2026-08-31

## Contexto
El chequeo de backing de sdd-verify's decision-gaps es una diferencia de conjuntos sobre el changed-file set de una sola tarea: si esa tarea tocó algún archivo fuera de .matecito-ai/edr/, backing es OK; si sólo tocó el registro y su INDEX, es CRITICAL. Con record-mode: modify ahora posible, una tarea de sólo-gobernanza edita únicamente .matecito-ai/edr/ — el chequeo, corrido tal cual, fallaría todo edit-in-place por construcción; y si el grupo no reconoce el nuevo valor de result, la fila se salta en silencio, que es peor que cualquiera de los dos veredictos.

## Decisión
En una fila `modified` de `### Decisions Materialized`, el chequeo de structure corre sin cambios (mismo comando, misma lectura OK/CRITICAL, contra el archivo editado). El chequeo de backing NO corre: se declara no-aplicable, y la celda queda `n/a — record edited in place, governance-only` — la fila no es ni OK ni CRITICAL en esa columna. Las columnas de la sección (record | task | structure | backing) no cambian; las celdas son texto libre, así que no hace falta editar el contrato de retorno.

## Reglas verificables
- **[manual]** in-flow-capture.md documenta la tercera rama de `result` (`modified`) en el grupo decision-gaps: structure sin cambios, backing con la celda `n/a — record edited in place, governance-only`.
- **[manual]** Las columnas `record | task | structure | backing` de `## Decision Gaps` no cambian de forma — el `n/a` es texto libre en la celda existente, no un valor nuevo que el contrato de retorno tenga que declarar.

## Alternativas consideradas
(a) dejar el grupo sin enterarse de `modified` — su lista de chequeos se indexa por `result`, así que un valor no reconocido saltearía la fila en silencio, peor que cualquiera de los dos veredictos y exactamente la falla de "la ausencia no es un veredicto limpio" que este ecosistema evita en otros lados. (b) correr backing tal cual está escrito — una tarea de sólo-gobernanza que no toca nada fuera de .matecito-ai/edr/ es CRITICAL por construcción, así que toda edición en el lugar fallaría la verificación. (c) forzar un archivo no relacionado dentro de la tarea para satisfacer la diferencia de conjuntos — es jugarle trampa al chequeo, no responderlo.

## Consecuencias
Costo aceptado, enunciado y NO mitigado: una fila `modified` cuya edición realmente no tiene nada detrás en ningún lado del cambio ahora pasa desapercibida. Ningún chequeo de este diseño la atrapa, y no se propone ninguno.
