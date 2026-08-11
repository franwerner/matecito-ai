# EDR — Una corrida aislada cierra con exactamente un commit, formato heredado por referencia

- **Status:** Accepted
- **Date:** 2026-08-10

## Contexto

El trabajo de una corrida aislada tiene que sobrevivir a la liberación de su espacio de trabajo, y la
única forma de que eso pase de manera confiable es que quede persistido en el repositorio antes de que
esa corrida devuelva el control. La convención de commits del proyecto ya existe y ya define formato,
tipos, scope y reglas duras de atribución — no hace falta inventar una segunda convención para este caso.

Esa convención, sin embargo, tiene como premisa que la IA commitea lo que una persona ya escribió y
revisó, con un loop de autorización explícita ante cualquier cambio que no cierre atómico. Una corrida
aislada de este mecanismo es headless: no tiene a quién preguntarle, y produce exactamente un cambio
deliberado (una tarea), así que la premisa de esa convención está invertida acá — no hay ni mezcla que
atomizar ni una persona esperando del otro lado del loop.

## Decisión

Cada corrida aislada produce **exactamente un commit** antes de devolver el control. El mensaje sigue el
formato, los tipos, el scope y las reglas duras de atribución de la convención de commits del proyecto,
**citada por referencia, nunca invocada como skill**. No se hereda el loop interactivo de autorización
ante no-atomicidad (no hay canal para preguntar dentro del aislamiento) ni el gate de "no commit sin
compilar" (ese chequeo se desplaza al verify de fin de batch).

## Reglas verificables

- **[manual]** Cada corrida aislada produce exactamente un commit antes de devolver el control.
- **[manual]** El mensaje del commit sigue el formato, los tipos, el scope y las reglas de atribución de la convención de commits del proyecto, sin invocar esa convención como skill.
- **[manual]** Ninguna corrida aislada corre el loop interactivo de atomicidad ni bloquea su commit a que el código compile; ese chequeo queda a cargo del verify de fin de batch.

## Alternativas consideradas

- **Invocar la convención de commits tal cual, con su loop interactivo.** Descartado: el loop exige autorización explícita de una persona ante cualquier cambio no atomizable, y una corrida aislada no tiene ese canal — invocarla headless la deja bloqueada sin salida.
- **Permitir varios commits por corrida.** Descartado: rompe la correspondencia 1 tarea = 1 commit que sostiene la atribución de conflictos y el cherry-pick por tarea de la integración.

## Consecuencias

- Se hereda el formato, los tipos, el scope y las reglas duras de atribución tal cual — nada nuevo que aprender para quien ya sigue esa convención.
- No se hereda el loop de autorización ni el gate de compilación — quedan explícitamente fuera, y el gate de compilación se recupera en el verify de fin de batch en vez de perderse.
- La correspondencia 1 tarea = 1 commit es lo que hace posible que la integración procese cada tarea con un solo cherry-pick.

## Relacionados

- `relacionado-con` → [worktree-base-handshake.md](worktree-base-handshake.md) — el chequeo que corre antes de este commit, para no escribir sobre una base desconocida.
- `relacionado-con` → [../structure/dispatch-batch-bound-integration.md](../structure/dispatch-batch-bound-integration.md) — el commit es lo que la integración diferida cherry-pickea.
