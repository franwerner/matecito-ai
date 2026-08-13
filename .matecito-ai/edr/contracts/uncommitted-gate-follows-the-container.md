# EDR — El Uncommitted-Work Gate inspecciona el contenedor inmediato, no siempre el repo principal

- **Status:** Accepted
- **Date:** 2026-08-13

## Contexto
El Uncommitted-Work Gate existía para un solo nivel de aislamiento: inspeccionaba el repo principal porque ahí era donde vivía el trabajo sucio que podía filtrarse a una corrida aislada. Con el aislamiento anidado por cambio, ese trabajo puede vivir en el espacio de trabajo del cambio en vez de en el repo principal — inspeccionar siempre el repo principal avisaría sobre suciedad que la corrida nunca toca, o se perdería suciedad real dentro del espacio de trabajo.

## Decisión
El gate inspecciona el árbol del **contenedor inmediato** de la ronda: el espacio de trabajo del cambio cuando el aislamiento anidado está activo, el repo principal cuando no — extendiendo el mismo principio que `contracts/worktree-base-handshake.md` ya fija para contra qué se verifica la base.

## Reglas verificables
- **[manual]** El Uncommitted-Work Gate corre `git status --porcelain --untracked-files=all` desde la raíz del contenedor inmediato de la ronda — el espacio de trabajo del cambio cuando el aislamiento anidado está activo, el repo principal cuando no.
- **[manual]** El outcome 'trabajar en la rama, sin worktree' degrada la ronda a la rama del contenedor inmediato, no siempre a la rama de trabajo original.

## Alternativas consideradas
Mantener siempre el repo principal. Descartado: avisa sobre trabajo que la ronda no puede tocar, y se pierde trabajo sucio real dentro del espacio de trabajo. Inspeccionar los dos. Descartado: ruido — dos avisos por una sola ronda cuando sólo uno de los dos árboles es relevante.

## Consecuencias
El gate necesita saber cuál es el contenedor inmediato de cada ronda antes de correr, el mismo dato que ya resuelve el handshake de base. Ningún cambio de comportamiento cuando el aislamiento por cambio está inactivo: el contenedor inmediato sigue siendo el repo principal, exactamente como antes.

## Relacionados
- `relacionado-con` → [worktree-base-handshake.md](worktree-base-handshake.md) — el mismo principio de 'contenedor inmediato' que esta decisión extiende del handshake de base al gate.
- `relacionado-con` → [../structure/change-level-worktree-isolation.md](../structure/change-level-worktree-isolation.md) — el aislamiento anidado que introduce el segundo contenedor posible.
