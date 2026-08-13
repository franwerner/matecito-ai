# EDR — Cada nivel de integración tiene su propio ejecutor

- **Status:** Accepted
- **Date:** 2026-08-13

## Contexto

Un batch de implementación paralelo necesita que algo cherry-pickee el trabajo de cada corrida aislada
sobre el destino del batch, escriba el artefacto de progreso del cambio y marque las tareas. Ese trabajo
—mutar el repositorio y persistir memoria— es exactamente lo que ya hace la fase de implementación al
cierre de cualquier batch, y reutiliza su Merge Protocol existente tal cual. El orquestador, en cambio,
es un coordinador: no muta el repositorio ni escribe memoria por su cuenta.

Con el aislamiento anidado por cambio aparece un segundo nivel de integración, de otra naturaleza: llevar
el cambio entero, ya verificado, del espacio de trabajo que lo aisló a la rama original. Ese nivel no
pertenece a ninguna fase —arranca cuando el pipeline terminó— y no tiene reportes de corrida que
consolidar ni artefacto de progreso que escribir: es un único movimiento sobre el espacio de trabajo que
el orquestador abrió y cuyo ciclo de vida ya administra.

## Decisión

El ejecutor se asigna **por nivel**:

- **Nivel de batch:** quien integra un batch paralelo es una **segunda invocación de la misma fase de
  implementación**, sin aislamiento — un segundo modo del mismo agente, nunca el orquestador ni un agente
  nuevo. Su destino es el espacio de trabajo del cambio cuando ese nivel está activo, y la rama de trabajo
  cuando no.
- **Nivel de cambio:** quien integra el espacio de trabajo del cambio a la rama original, al cerrarse el
  ciclo, es el **orquestador**, que es el que abrió ese espacio y administra su ciclo de vida.

## Reglas verificables

- **[manual]** La corrida que integra un batch paralelo es una invocación de la fase de implementación, sin aislamiento, nunca el orquestador ni un agente distinto.
- **[manual]** La corrida de consolidación reutiliza el Merge Protocol existente de la fase tal cual, sin un contrato de persistencia paralelo para el caso batch.
- **[manual]** La integración de nivel de cambio la ejecuta el orquestador, una sola vez, al cerrarse el ciclo; ninguna fase la ejecuta ni la anticipa.

## Alternativas consideradas

- **Que el orquestador integre también el nivel de batch.** Descartado: para ese nivel el orquestador es coordinador, no ejecutor — asumir cherry-pick por tarea, commits y escritura de memoria ahí rompe el reparto de responsabilidades que sostiene el resto del pipeline, y duplicaría un Merge Protocol que la fase ya define y mantiene.
- **Que una fase ejecute también el nivel de cambio.** Descartado: ese nivel arranca cuando el pipeline ya terminó, sobre un espacio de trabajo cuyo ciclo de vida administra el orquestador; asignárselo a una fase la obligaría a integrar el contenedor que la contiene.
- **Un agente nuevo dedicado a integrar.** Descartado: para el nivel de batch duplicaría el Merge Protocol y el formato del artefacto que la fase de implementación ya define y mantiene; un agente nuevo es un segundo lugar donde ese contrato puede desalinearse del original. Para el nivel de cambio sería un despacho entero por un único movimiento sobre un espacio de trabajo que el orquestador ya tiene abierto.

## Consecuencias

- No hace falta un contrato de persistencia nuevo: la consolidación es la misma fase, con el mismo formato de artefacto y el mismo return template, aplicados sobre N reportes en vez de sobre su propia implementación.
- El agente de implementación gana un segundo modo de operación (aislado vs. consolidación) que su propio método tiene que distinguir explícitamente al arrancar.
- El reparto "el orquestador no muta el repositorio" deja de ser absoluto y pasa a estar acotado al nivel de batch, que es donde la mutación es por tarea y necesita el protocolo de una fase.
- Los dos niveles pueden fallar por separado: un conflicto al integrar el batch se atribuye a una tarea; uno al integrar el cambio se atribuye al cambio entero contra la rama original.

## Relacionados

- `relacionado-con` → [dispatch-batch-bound-integration.md](dispatch-batch-bound-integration.md) — cuándo arranca la integración de nivel de batch que esta decisión asigna.
- `relacionado-con` → [change-level-worktree-isolation.md](change-level-worktree-isolation.md) — el aislamiento anidado que introduce el segundo nivel.
- `relacionado-con` → [../contracts/single-writer-per-batch.md](../contracts/single-writer-per-batch.md) — la corrida de consolidación es, además, el único escritor del batch.
