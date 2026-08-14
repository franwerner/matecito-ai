# EDR — El test de cruce hacia territorio de contrato se angosta, sin ganar una segunda cláusula

- **Status:** Accepted
- **Date:** 2026-08-14

## Contexto
La regla que obliga a proponer la forma completa de un contrato en vez de inferirla ya tenía un test para decidir cuándo una edición "cruza" hacia territorio de contrato. Extender esa regla a un lugar nuevo del pipeline tentaba a agregarle una segunda cláusula específica para ese caso.

## Decisión
El test de cruce se angosta —se precisa mejor qué cuenta como cruzar— pero se mantiene como un único test, sin ganar una segunda cláusula paralela para el caso nuevo.

## Reglas verificables
- **[manual]** El test que decide si una edición cruza hacia territorio de contrato sigue siendo un único test; ningún caso nuevo le agrega una cláusula paralela.

## Alternativas consideradas
Agregar una segunda cláusula específica para el caso nuevo — descartada porque dos cláusulas independientes decidiendo lo mismo (cuándo algo es un contrato que no se puede inferir) son dos lugares donde esa definición puede desalinearse entre sí.

## Consecuencias
La regla de "nunca inferido" se mantiene como una única definición aplicable a todos los lugares que la usan, con una precisión mayor de qué cuenta como cruce — a costa de que cualquier ajuste futuro a esa definición tiene que seguir siendo compatible con todos esos usos a la vez.

## Relacionados
- `relacionado-con` → [contract-proposal-has-no-persistence-slot.md](contract-proposal-has-no-persistence-slot.md) — un caso que usa este test angostado.
