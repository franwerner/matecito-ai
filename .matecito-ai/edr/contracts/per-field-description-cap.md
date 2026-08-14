# EDR — La descripción de cada campo de un contrato tiene su propio tope de caracteres

- **Status:** Accepted
- **Date:** 2026-08-14

## Contexto
Un ítem de contrato propuesto puede declarar más de un campo, cada uno con su propia descripción. El tope que ya existía sobre el resumen del ítem no dice nada sobre cuánto puede crecer la descripción de cada campo individual.

## Decisión
El tope del resumen del ítem sigue midiendo sólo el resumen. Cada descripción de campo tiene su propio tope, de 160 caracteres — más chico que el del resumen porque describe una sola pieza, no el ítem entero. La cantidad de campos que un contrato propuesto puede declarar nunca se topea.

## Reglas verificables
- **[auto]** La descripción de cada campo de un contrato propuesto se mide contra su propio tope de 160 caracteres, independiente del tope del resumen del ítem.
- **[manual]** Ningún límite se aplica a la cantidad de campos que un ítem de contrato puede declarar.

## Alternativas consideradas
Medir la descripción de campo contra el mismo tope que el resumen del ítem — descartada porque una descripción de campo cubre una sola pieza del contrato, no el ítem entero, y compartir el mismo tope resultaba o demasiado laxo o demasiado estricto según el caso.

## Consecuencias
Un contrato con muchos campos sigue siendo legible campo por campo, sin que la cantidad de campos misma se vuelva un motivo de rechazo — a costa de dos topes distintos que hay que recordar mantener alineados dentro del mismo contrato.

## Relacionados
- `relacionado-con` → [gate-summary-cap-in-characters.md](gate-summary-cap-in-characters.md) — mismo criterio de tope por caracteres, aplicado al resumen del ítem en vez de a la descripción de un campo.
- `relacionado-con` → [nested-field-continuation-line.md](nested-field-continuation-line.md) — la forma en la que cada campo, con su descripción topeada, se imprime.
