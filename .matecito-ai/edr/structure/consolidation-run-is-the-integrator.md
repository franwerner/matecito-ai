# EDR — La integración de un batch paralelo la hace una segunda corrida de la misma fase

- **Status:** Accepted
- **Date:** 2026-08-10

## Contexto

Un batch de implementación paralelo necesita que algo cherry-pickee el trabajo de cada corrida aislada
sobre la rama de trabajo, escriba el artefacto de progreso del cambio y marque las tareas. Ese trabajo
—mutar el repositorio y persistir memoria— es exactamente lo que ya hace la fase de implementación al
cierre de cualquier batch, y reutiliza su Merge Protocol existente tal cual. El orquestador, en cambio,
es un coordinador: no muta el repositorio ni escribe memoria por su cuenta.

## Decisión

Quien integra un batch paralelo es una **segunda invocación de la misma fase de implementación**, sin
aislamiento — un segundo modo del mismo agente, nunca el orquestador ni un agente nuevo.

## Reglas verificables

- **[manual]** La corrida que integra un batch paralelo es una invocación de la fase de implementación, sin aislamiento, nunca el orquestador ni un agente distinto.
- **[manual]** La corrida de consolidación reutiliza el Merge Protocol existente de la fase tal cual, sin un contrato de persistencia paralelo para el caso batch.

## Alternativas consideradas

- **Que integre el orquestador.** Descartado: el orquestador es coordinador, no ejecutor — asumir cherry-pick, commits y escritura de memoria ahí rompe el reparto de responsabilidades que sostiene el resto del pipeline (mutaciones irreversibles del repositorio son trabajo de fase, no de orquestación).
- **Un agente nuevo dedicado a integrar.** Descartado: duplicaría el Merge Protocol y el formato del artefacto que la fase de implementación ya define y mantiene; un agente nuevo es un segundo lugar donde ese contrato puede desalinearse del original.

## Consecuencias

- No hace falta un contrato de persistencia nuevo: la consolidación es la misma fase, con el mismo formato de artefacto y el mismo return template, aplicados sobre N reportes en vez de sobre su propia implementación.
- El agente de implementación gana un segundo modo de operación (aislado vs. consolidación) que su propio método tiene que distinguir explícitamente al arrancar.

## Relacionados

- `relacionado-con` → [dispatch-batch-bound-integration.md](dispatch-batch-bound-integration.md) — cuándo arranca la integración que esta decisión asigna.
- `relacionado-con` → [../contracts/single-writer-per-batch.md](../contracts/single-writer-per-batch.md) — esta corrida es, además, el único escritor del batch.
