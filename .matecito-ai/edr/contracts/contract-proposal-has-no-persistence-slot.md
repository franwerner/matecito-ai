# EDR — Una propuesta de contrato no declara dónde va a persistir

- **Status:** Accepted
- **Date:** 2026-08-14

## Contexto
Cuando una fase propone la forma de un contrato todavía no fijado (una entidad, un DTO, un esquema), esa propuesta tiene que viajar de vuelta a la fase que la va a implementar, sin decidir todavía dónde va a persistir esa forma en el repo.

## Decisión
El ítem que propone un contrato no lleva ningún campo que indique dónde vivirá una vez ratificado. La ubicación la decide la fase que lo implementa, siguiendo la regla del dominio que ya dice dónde vive cada tipo de contrato — nunca la propuesta misma.

## Reglas verificables
- **[manual]** Un ítem de contrato propuesto no declara un campo de destino o ubicación de persistencia.

## Alternativas consideradas
Agregar un campo que indique dónde persistirá el contrato — descartada porque esa ubicación ya la fija la regla del dominio sobre dónde vive cada tipo de contrato; declararla de nuevo en la propuesta sería un segundo lugar diciendo lo mismo, con riesgo de desalinearse.

## Consecuencias
La propuesta se mantiene angosta —sólo la forma del contrato, nunca su destino— a costa de que la fase que la implementa tiene que resolver la ubicación por su cuenta, apoyándose en la regla del dominio y no en la propuesta.

## Relacionados
- `relacionado-con` → [ratified-contract-travels-in-the-dispatch-prompt.md](ratified-contract-travels-in-the-dispatch-prompt.md) — cómo la forma ratificada llega de vuelta a la fase que la implementa.
