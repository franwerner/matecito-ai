# EDR — Quien integra un batch paralelo es la fase de implementación, y sostiene el turno mientras lo hace

- **Status:** Accepted
- **Date:** 2026-08-13

## Contexto

Un batch de implementación paralelo necesita que algo cherry-pickee el trabajo de cada corrida aislada
sobre el destino del batch, escriba el artefacto de progreso del cambio y marque las tareas. Ese trabajo
—mutar el repositorio y persistir memoria— es exactamente lo que ya hace la fase de implementación al
cierre de cualquier batch, y reutiliza su Merge Protocol existente tal cual. El orquestador, en cambio,
es un coordinador: no muta el repositorio ni escribe memoria por su cuenta.

## Decisión

Quien integra un batch paralelo es una **segunda invocación de la misma fase de implementación**, sin
aislamiento — un segundo modo del mismo agente, nunca el orquestador ni un agente nuevo. Su destino es
la rama de trabajo.

La corrida de consolidación sostiene el turno de rama compartida mientras integra: lo reclama antes del
primer cherry-pick del loop y lo libera recién después de la última iteración, cualquiera sea el
resultado.

## Reglas verificables

- **[manual]** La corrida que integra un batch paralelo es una invocación de la fase de implementación, sin aislamiento, nunca el orquestador ni un agente distinto.
- **[manual]** La corrida de consolidación reutiliza el Merge Protocol existente de la fase tal cual, sin un contrato de persistencia paralelo para el caso batch.

## Alternativas consideradas

- **Que el orquestador integre también el nivel de batch.** Descartado: para ese nivel el orquestador es coordinador, no ejecutor — asumir cherry-pick por tarea, commits y escritura de memoria ahí rompe el reparto de responsabilidades que sostiene el resto del pipeline, y duplicaría un Merge Protocol que la fase ya define y mantiene.
- **Un agente nuevo dedicado a integrar.** Descartado: duplicaría el Merge Protocol y el formato del artefacto que la fase de implementación ya define y mantiene; un agente nuevo es un segundo lugar donde ese contrato puede desalinearse del original.

## Consecuencias

- No hace falta un contrato de persistencia nuevo: la consolidación es la misma fase, con el mismo formato de artefacto y el mismo return template, aplicados sobre N reportes en vez de sobre su propia implementación.
- El agente de implementación gana un segundo modo de operación (aislado vs. consolidación) que su propio método tiene que distinguir explícitamente al arrancar.
- Un conflicto al integrar el batch se atribuye a la tarea que lo generó, nunca al cambio entero.

## Relacionados

- `relacionado-con` → [dispatch-batch-bound-integration.md](dispatch-batch-bound-integration.md) — cuándo arranca la integración de nivel de batch que esta decisión asigna.
- `relacionado-con` → [../contracts/single-writer-per-batch.md](../contracts/single-writer-per-batch.md) — la corrida de consolidación es, además, el único escritor del batch.
