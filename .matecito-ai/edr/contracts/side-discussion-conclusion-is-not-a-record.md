# EDR — La conclusión de una discusión lateral es material de trabajo, nunca un registro de decisión

- **Status:** Accepted
- **Date:** 2026-08-14

## Contexto
El dominio de desarrollo tiene una regla firme: los EDRs y los capability-specs nunca se graban en Engram, viven sólo como archivo. El traspaso y la conclusión de una discusión lateral, en cambio, viven en Engram por diseño. Hacía falta fijar qué pasa cuando una conclusión de discusión lateral zanja algo arquitectónico, para no dejar la impresión de que esa regla se dobla.

## Decisión
La conclusión es material de trabajo, de la misma categoría que el artefacto `sdd/{change-name}/design` — razonamiento que lleva a una decisión, no la decisión ya ratificada. Cuando una conclusión zanja una pregunta arquitectónica, entra al camino existente del flujo: se vuelve un ítem bajo `New Decisions` de la fase que la propone, se ratifica una vez en el gate del lane, y `sdd-apply` materializa el archivo de EDR. La sesión lateral nunca escribe un EDR ni un capability-spec, y el principal tampoco copia una conclusión a un EDR sin pasar por ese gate.

## Reglas verificables
- **[manual]** La sesión lateral nunca escribe un archivo de EDR ni de capability-spec — su único output es la conclusión en Engram.
- **[manual]** El principal nunca copia una conclusión directamente a un EDR sin que pase antes por el gate de ratificación de la fase que la propone (`New Decisions`).

## Alternativas consideradas
La sesión lateral escribe el EDR directamente — descartado: rompe el límite "sólo discute" y materializa un registro que nadie ratificó. El principal copia la conclusión a un EDR sin gate — descartado por el mismo motivo: materializa un registro sin ratificación.

## Consecuencias
Ninguna regla existente se dobla — EDRs y capability-specs siguen sin grabarse nunca en Engram. El costo es un paso extra: una conclusión que zanja algo arquitectónico no se vuelve EDR de inmediato, tiene que pasar por proponer y ratificar como cualquier otra decisión del flujo.

## Relacionados
- `relacionado-con` → [side-discussion-conclusion-shape.md](side-discussion-conclusion-shape.md) — la forma que hace posible que la conclusión entre sin transformación al camino de captura.
