# EDR — La fila del registro de ratificación lleva exactamente cuatro campos

- **Status:** Accepted
- **Date:** 2026-08-14

## Contexto
Cuando un gate ratifica un ítem, algo tiene que quedar registrado para que un gate posterior del mismo cambio pueda reconocer si ese mismo ítem vuelve a aparecer. Ese registro necesita ser comparable mecánicamente, no prosa libre que cada gate interprete distinto.

## Decisión
Cada fila del registro de ratificación lleva exactamente cuatro campos: el identificador del ítem, el resumen ratificado (el texto final ajustado por el usuario gana sobre el propuesto originalmente), el ancla y el gate en el que se ratificó. Ningún campo de resolución adicional: el registro cubre sólo lo ratificado, no un mecanismo de rechazo que nadie pidió para este registro.

## Reglas verificables
- **[auto]** Toda fila del registro de ratificación declara los cuatro campos: identificador, resumen ratificado, ancla, gate.
- **[manual]** El resumen que queda en la fila es el texto final ajustado por el usuario, nunca el propuesto originalmente si el usuario lo cambió.

## Alternativas consideradas
Agregar una columna de resolución — descartada porque cubrir rechazos es una capacidad que nadie pidió para este registro. Prosa libre en vez de campos fijos — descartada porque un gate posterior necesita comparar mecánicamente, y la prosa libre no es comparable entre sí.

## Consecuencias
El registro queda angosto y mecánicamente comparable — a costa de no poder representar, todavía, un ítem que fue rechazado en vez de ratificado.

## Relacionados
- `relacionado-con` → [re-emergence-match-on-record.md](re-emergence-match-on-record.md) — usa el identificador de esta fila como clave de comparación.
- `relacionado-con` → [re-emergence-short-form.md](re-emergence-short-form.md) — usa el resumen ratificado y el ancla de esta fila para el aviso corto.
