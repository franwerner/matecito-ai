# EDR — El espacio de trabajo del cambio viaja como línea explícita en cada despacho de fase

- **Status:** Accepted
- **Date:** 2026-08-13

## Contexto
El directorio de trabajo de una sesión se fija al arrancar y ningún participante se reubica durante el ciclo de un cambio: cada fase sigue corriendo desde donde fue lanzada. Con el aislamiento por cambio activo, el trabajo tiene que pasar dentro del espacio de trabajo en vez de en ese directorio fijo — pero ninguna fase puede descubrir por su cuenta dónde está ese espacio sin un estado que nadie escribe, y confundirlo con los worktrees por tarea que el batch paralelo ya abre es un riesgo real dado que ambos son worktrees.

## Decisión
Mientras el espacio de trabajo de un cambio está abierto, cada instrucción de despacho de fase — incluidas las que sólo leen — lleva una línea `workspace: <ruta absoluta>`. Las fases resuelven rutas del repo bajo esa ruta y corren git con `git -C <workspace> ...`; un comando de proyecto (una corrida de tests, un build) corre con esa misma ruta como directorio de trabajo del comando.

## Reglas verificables
- **[manual]** Cada despacho de fase, mientras el espacio de trabajo del cambio está abierto, lleva una línea `workspace: <ruta absoluta>` — incluidas las fases que sólo leen.
- **[manual]** Una fase resuelve rutas del repo bajo esa ruta y corre git con `git -C <workspace> ...` en vez de asumir el directorio de trabajo de la sesión.
- **[manual]** Un comando de proyecto que una fase ejecuta (tests, build, lint) corre con esa misma ruta como directorio de trabajo del comando.

## Alternativas consideradas
Que cada fase descubra el espacio de trabajo por su cuenta. Descartado: necesita un estado que hoy nadie escribe, y se presta a confundirse con los worktrees por tarea que el batch paralelo ya abre — ambos son worktrees, y sólo uno de los dos es el que la fase debe usar.

## Consecuencias
Cada prompt de despacho gana una línea más mientras el espacio de trabajo está abierto — costo bajo, pagado una vez por despacho. Ninguna fase necesita lógica de descubrimiento ni relocalización: recibe la ruta, la usa, listo.

## Relacionados
- `relacionado-con` → [../structure/change-level-worktree-isolation.md](../structure/change-level-worktree-isolation.md) — el aislamiento cuyo espacio de trabajo esta decisión hace localizable para cada fase.
- `relacionado-con` → [uncommitted-gate-follows-the-container.md](uncommitted-gate-follows-the-container.md) — otro chequeo que depende de saber cuál es el contenedor inmediato de la ronda.
