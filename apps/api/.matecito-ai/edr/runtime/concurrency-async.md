# EDR — Concurrencia y asincronía

- **Status:** Accepted
- **Type:** decision
- **Date:** 2026-07-23

## Contexto

El broker corre varias actividades concurrentes de larga vida: un escritor de la persistencia, un hub que hace fan-out de eventos a la UI por WebSocket y uno o más tailers incrementales del transcript. Todo esto debe arrancar, esperarse y cancelarse de forma coordinada, y las escrituras a la persistencia deben serializarse para evitar contención. Escribir a mano el boilerplate de coordinación (esperas y cancelación) es propenso a errores.

## Decisión

Adoptamos el modelo de concurrencia idiomático de Go: goroutines, channels y context para la cancelación. El lifecycle (arranque, espera y cancelación coordinada del escritor, el hub y los tailers) se orquesta con un grupo de goroutines con cancelación propagada, en lugar de coordinar esperas y cancelación a mano. Hay un único goroutine escritor de la persistencia, alimentado por channel (single-writer), un hub para el fan-out a la UI, y el tail del transcript es incremental por offset.

## Reglas verificables

- **[manual]** hay un único goroutine escritor de la persistencia, alimentado por channel.
- **[manual]** el lifecycle (arranque, espera y cancelación) se coordina con el grupo de goroutines + context, no con esperas manuales.
- **[manual]** el tail del transcript es incremental (avanza por offset), no re-lee el archivo entero.

## Alternativas consideradas

- **Coordinación manual de esperas y cancelación:** descartada por reescribir a mano boilerplate que la librería de sincronización ya resuelve, con más riesgo de fugas de goroutines.
- **Múltiples escritores concurrentes a la persistencia:** descartada; el single-writer evita contención y simplifica la frontera transaccional.

## Consecuencias

- El apagado es limpio y coordinado; cancelar el context detiene todas las actividades.
- El single-writer elimina condiciones de carrera sobre la persistencia.
- Trade-off: todo el volumen de escritura pasa por un solo goroutine, que es el punto de serialización.

## Relacionados

- `relacionado-con` → [../data/data-access-entity-framework.md](../data/data-access-entity-framework.md) — la frontera transaccional se apoya en este único escritor.
