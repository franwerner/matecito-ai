# EDR — El espacio de trabajo de un cambio es un git worktree en una ruta y rama fijas

- **Status:** Accepted
- **Date:** 2026-08-13

## Contexto
Con el aislamiento por cambio activo, el orquestador necesita abrir una copia separada del proyecto por cada cambio. Git ofrece dos convenciones para eso: un directorio hermano fuera del repo (`../<repo>-<change>`), o un directorio dentro del propio repo. La sesión ya trabaja dentro de un límite de permisos fijado al arrancar; un directorio fuera de ese límite haría que cada escritura de cada fase dispare una confirmación.

## Decisión
El espacio de trabajo de un cambio es un git worktree en la rama `matecito-ai/<change-name>`, con directorio `<repo>/.matecito-ai/workspaces/<change-name>`. Al abrirse, se agrega `.matecito-ai/workspaces/` al `.gitignore` del repo si no está ya, antes de que cualquier fase escriba adentro — así el Uncommitted-Work Gate nunca lee el propio directorio del espacio de trabajo como suciedad del repo principal.

## Reglas verificables
- **[manual]** El espacio de trabajo de un cambio es un git worktree en la rama `matecito-ai/<change-name>`, con directorio `<repo>/.matecito-ai/workspaces/<change-name>`.
- **[manual]** Al abrirse, `.matecito-ai/workspaces/` se agrega al `.gitignore` del repo si no está ya presente, antes de que cualquier fase escriba dentro del espacio de trabajo.

## Alternativas consideradas
Un directorio hermano fuera del repo, la convención habitual de git worktree. Descartado: cada escritura de cada fase caería entonces fuera del límite de permisos de la sesión, disparando una confirmación por archivo.

## Consecuencias
El repo gana una carpeta ignorada por git donde viven los espacios de trabajo activos; el chequeo de `.gitignore` es parte de la apertura, no un paso aparte que alguien pueda saltear. Cada cambio con aislamiento activo cuesta un worktree y una rama propios, liberados al cerrar el ciclo.

## Relacionados
- `relacionado-con` → [change-level-worktree-isolation.md](change-level-worktree-isolation.md) — el aislamiento por cambio cuya identidad esta decisión fija.
- `relacionado-con` → [../contracts/uncommitted-gate-follows-the-container.md](../contracts/uncommitted-gate-follows-the-container.md) — por qué el directorio necesita quedar gitignored antes de que el gate corra.
