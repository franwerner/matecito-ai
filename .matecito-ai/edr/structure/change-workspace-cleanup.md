# EDR — La limpieza del espacio de trabajo de un cambio espeja la del batch, sólo tras integrar limpio

- **Status:** Accepted
- **Date:** 2026-08-13

## Contexto
El batch de implementación paralelo ya limpia sus worktrees por tarea con una secuencia verificada: unlock → remove → branch -D, nunca forzada. El espacio de trabajo de un cambio es el mismo mecanismo un nivel más arriba — dejarlo sin definir su propia limpieza dejaría worktrees y ramas de cambio acumulándose indefinidamente, o forzaría a inventar una segunda secuencia.

## Decisión
Tras una integración limpia, el espacio de trabajo del cambio se limpia con la misma secuencia que el batch ya usa: `git worktree unlock` → `git worktree remove` → `git branch -D`, nunca `remove -f -f`; una falla en cualquier paso se registra, no se fuerza. Tras una integración fallida, tanto el espacio de trabajo como su rama quedan intactos para inspección.

## Reglas verificables
- **[manual]** Tras una integración limpia, el espacio de trabajo del cambio se limpia con `git worktree unlock` → `git worktree remove` → `git branch -D`, en ese orden, nunca `remove -f -f`.
- **[manual]** Una falla en cualquier paso de la limpieza se registra, no se fuerza; el trabajo ya está integrado independientemente de si su worktree se limpia.
- **[manual]** Tras una integración fallida, el espacio de trabajo y su rama quedan intactos, sin tocar, para inspección.

## Alternativas consideradas
Conservar la rama del cambio como rastro de auditoría después de integrar limpio. Descartado: los commits ya están en la rama original, así que conservar la rama del cambio es una copia sin usar, no un rastro adicional.

## Consecuencias
El repo no acumula worktrees ni ramas de cambios ya integrados. Que una integración fallida deje todo intacto es, a la vez, lo que permite la inspección y lo que impide que la limpieza borre evidencia de un conflicto no resuelto.

## Relacionados
- `relacionado-con` → [change-level-worktree-isolation.md](change-level-worktree-isolation.md) — el aislamiento cuyo espacio de trabajo esta decisión limpia.
- `relacionado-con` → [change-level-integration-act.md](change-level-integration-act.md) — el resultado de la integración que decide si la limpieza corre o no.
