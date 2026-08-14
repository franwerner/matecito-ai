# EDR — Un ítem de contrato lleva un solo ancla, nunca una por campo

- **Status:** Accepted
- **Date:** 2026-08-14

## Contexto
Un contrato propuesto puede declarar varios campos. Si cada campo llevara su propia ancla, un mismo contrato tendría tantas anclas como campos, y la identidad de la propuesta —qué usar para reconocerla cuando vuelve ratificada— quedaría repartida en vez de ser una sola cosa.

## Decisión
El ancla se declara una sola vez por ítem de contrato, nunca una por campo. Esa misma ancla, junto con el resumen del ítem, es también la identidad que se usa para reconocer la propuesta cuando su forma ratificada vuelve a la fase que la implementa.

## Reglas verificables
- **[auto]** Un ítem de contrato propuesto declara exactamente un token de ancla, sin importar cuántos campos tenga.
- **[manual]** El ancla y el resumen del ítem son la identidad que identifica la propuesta cuando su forma ratificada regresa a la fase que la implementa.

## Alternativas consideradas
Un ancla por campo — descartada porque fragmentaría la identidad de una misma propuesta en tantas piezas como campos, sin que ninguna sea la identidad de la propuesta completa.

## Consecuencias
Reconocer una propuesta ratificada es siempre una sola comparación, contra una sola ancla — a costa de que el ancla no distingue de qué campo vino un cambio puntual dentro del contrato, si sólo un campo cambió.

## Relacionados
- `relacionado-con` → [nested-field-continuation-line.md](nested-field-continuation-line.md) — los campos que este ancla único acompaña.
- `relacionado-con` → [ratified-contract-travels-in-the-dispatch-prompt.md](ratified-contract-travels-in-the-dispatch-prompt.md) — usa esta identidad para reconocer la propuesta al reenviarla.
