# EDR — El cierre del ciclo integra por merge, con un intento de rebase antes de reportar el conflicto

- **Status:** Accepted
- **Date:** 2026-08-13

## Contexto
Al cerrar el ciclo de un cambio con aislamiento activo, el espacio de trabajo aislado tiene que volver a la rama original en un único movimiento — nadie más lo hace, y postergarlo a un paso manual deja el trabajo verificado fuera de la rama por tiempo indefinido. Un merge puede fallar; forzarlo pierde información sobre qué versión gana en cada conflicto, y abortar sin más deja al usuario sin ninguna pista de qué pasó.

## Decisión
El cierre es un único `git merge --no-ff`, corrido por el orquestador desde el repo principal, después de `sdd-archive` en flujo, o al reportarse terminado el trabajo directo/ad-hoc. Si el merge falla, el orquestador no aborta de inmediato: primero intenta un `git rebase <rama original>` dentro del espacio de trabajo. Si el rebase sale limpio, reintenta el merge — ahora un fast-forward. Si el rebase también choca, el orquestador aborta el rebase y el merge, y reporta al usuario qué pasó — qué archivo, en qué commit, de qué lado viene cada versión — pidiéndole elegir, con una recomendación según el caso (resolver a mano en el espacio de trabajo · dejarlo abierto · otra salida). Nada se fuerza en ningún punto, y el espacio de trabajo queda intacto en cualquier camino de falla.

## Reglas verificables
- **[manual]** El cierre del ciclo es un único `git merge --no-ff` desde el repo principal, ejecutado por el orquestador, una sola vez.
- **[manual]** Ante un merge fallido, el orquestador intenta primero `git rebase <rama original>` dentro del espacio de trabajo antes de abortar cualquier cosa.
- **[manual]** Si el rebase sale limpio, el orquestador reintenta el merge, ahora como fast-forward.
- **[manual]** Si el rebase también choca, el orquestador aborta rebase y merge, y reporta al usuario el conflicto (archivo, commit, versión de cada lado) con una recomendación, sin forzar nada; el espacio de trabajo queda intacto.

## Alternativas consideradas
Reproducir los commits para una historia lineal. Descartado frente a `--no-ff`: pierde la agrupación de un cambio como una unidad. Dejar la integración a una persona a mano. Descartado ya en `structure/change-level-worktree-isolation.md`: posterga el cierre a un paso manual sin fecha. Abortar directamente ante el primer conflicto de merge, sin intentar el rebase. Descartado en el gate: el usuario aceptó que un rebase no esquiva un conflicto de contenido real, pero sí compra una atribución más fina, porque reaplica commit por commit — vale la pena intentarlo antes de rendirse.

## Consecuencias
El orquestador ejecuta una mutación irreversible del repositorio al cerrar el ciclo — la misma acotación que ya reconoce `structure/change-level-worktree-isolation.md`. Un conflicto de integración ahora tiene dos intentos automáticos (merge, luego rebase-y-reintentar) antes de involucrar al usuario, y cuando lo involucra, lo hace con la información más fina que el rebase ya reunió, no con un mensaje genérico de conflicto.

## Relacionados
- `relacionado-con` → [change-level-worktree-isolation.md](change-level-worktree-isolation.md) — el aislamiento cuyo cierre esta decisión define.
- `relacionado-con` → [change-workspace-cleanup.md](change-workspace-cleanup.md) — qué pasa con el espacio de trabajo después de que la integración termina, en cualquiera de sus dos resultados.
