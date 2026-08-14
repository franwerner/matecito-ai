# EDR — Los campos de un contrato se imprimen como líneas de continuación, nunca como viñetas anidadas

- **Status:** Accepted
- **Date:** 2026-08-14

## Contexto
Un ítem de contrato con varios campos necesita imprimir cada campo de alguna forma legible. El parser que ya lee un ítem de gate interpreta cualquier línea que empiece con guión o asterisco, sin importar cuánto esté indentada, como el inicio de un ítem nuevo — así que anidar los campos como viñetas los haría leer como ítems propios en vez de partes del mismo ítem; una tabla anidada, en cambio, no la reconoce en absoluto y la salta en silencio.

## Decisión
Cada campo se imprime como una línea de continuación repetida, en la misma familia de líneas que ya usan los demás tokens de un ítem — nunca como una viñeta anidada ni como una tabla dentro del ítem.

## Reglas verificables
- **[auto]** Cada campo de un contrato propuesto se imprime como una línea de continuación del mismo ítem, nunca como una viñeta con guión o asterisco.
- **[manual]** La forma elegida se valida contra el comportamiento real del parser existente, no se asume compatible.

## Alternativas consideradas
Viñetas anidadas por campo — descartada porque el parser ya interpreta cualquier viñeta, en cualquier indentación, como un ítem nuevo, así que un campo anidado se leería como una propuesta separada y fallaría por token faltante. Una tabla anidada por campo — descartada porque el parser no la reconoce en absoluto, y el campo se perdería en silencio en vez de fallar de forma visible.

## Consecuencias
Un contrato con varios campos se imprime y se valida con la misma gramática que ya existe para tokens, sin gramática nueva que aprender — a costa de que la lectura visual de muchos campos es una lista plana de líneas, no una tabla anidada más legible a simple vista.

## Relacionados
- `relacionado-con` → [contract-item-anchors-once.md](contract-item-anchors-once.md) — el ancla que acompaña a estos campos.
- `relacionado-con` → [per-field-description-cap.md](per-field-description-cap.md) — el tope que se aplica sobre el contenido de cada línea de campo.
