# EDR — Un ítem reaparecido se confirma con un sí/no corto, no con el walkthrough completo

- **Status:** Accepted
- **Date:** 2026-08-14

## Contexto
Cuando un ítem reaparece y coincide con uno ya ratificado, el usuario no necesita repasar el ítem completo otra vez — pero tampoco conviene ocultarle si la premisa cambió desde la última vez que lo confirmó.

## Decisión
Un ítem reaparecido se muestra como una sola confirmación sí/no, con el resumen y el ancla ya ratificados. Sólo cuando el resumen entrante difiere del ratificado, se muestra también el resumen entrante, para que el usuario vea qué cambió antes de confirmar de nuevo.

## Reglas verificables
- **[auto]** Un ítem reaparecido se presenta como una única confirmación sí/no, nunca como el walkthrough completo del ítem.
- **[manual]** El resumen entrante se muestra junto al ratificado únicamente cuando difieren; si son iguales, no se duplica.

## Alternativas consideradas
Mostrar sólo el texto ya ratificado — descartada porque oculta que la premisa pudo haber cambiado, dejando al usuario confirmar a ciegas algo distinto de lo que ratificó antes.

## Consecuencias
Confirmar una reaparición es más rápido que la primera vez, sin perder visibilidad sobre un cambio de premisa — a costa de una comparación de texto adicional en cada reaparición, para decidir si mostrar el resumen entrante.

## Relacionados
- `relacionado-con` → [ratification-ledger-row.md](ratification-ledger-row.md) — fuente del resumen y el ancla ya ratificados.
- `relacionado-con` → [re-emergence-match-on-record.md](re-emergence-match-on-record.md) — la comparación que decide si un ítem reaparece.
