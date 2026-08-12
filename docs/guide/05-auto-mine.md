# 05 — Captura de decisiones (siempre activa)

[← 04 Decisiones y EDRs](04-decisiones-edr.md) · [Índice](README.md) · Siguiente: [06 — Auto-mine de specs →](06-auto-mine-spec.md)

> Esta página describe cómo cada dominio se asegura de que una decisión que el código termina
> implementando no se quede sin capturar. No es una opción de configuración: es comportamiento por
> defecto, sin interruptor, en todo cambio — cada dominio conserva su **propio mecanismo**.

No hay un único mecanismo compartido: cada dominio resuelve esto a su manera, según en qué punto del
flujo le conviene detectar el hueco.

- **development** — **in-flow**: la fase que llega a la decisión la propone en su propio retorno, el
  gate del lane la ratifica **una sola vez**, y `sdd-apply` la materializa como EDR `Accepted` en el
  mismo paso que implementa el código que la gobierna. No hay pase de minería después de verify.
- **design** — **mine gate post-verify**: `design-tasks` marca cada decisión de marca con
  `· ddr: <surface>/<slug>` (exista o no el DDR), `design-verify` confirma cuáles se implementaron
  (`## Decision Gaps`), y al cerrar el orquestador dispara `design-decisions-mine` — que lee el
  archivo Figma vivo — y ofrece materializar las decisiones como DDRs `Inferred`, siempre con tu
  confirmación.

## Por qué dos mecanismos distintos

`development` puede proponer la decisión en el momento en que la fase la encuentra, porque sus fases
son texto y estructura: no hace falta que el código exista todavía para articular el porqué. `design`
depende de un artefacto externo (el archivo Figma) cuya evidencia fuerte —los styles y components
nombrados— solo es confiable **después** de que la pieza se produjo; por eso su mecanismo detecta el
hueco post-verify, sobre lo ya embarcado, en vez de proponerlo por adelantado.

## Mecanismo de development: in-flow, fase por fase

```
fase que decide → propone la decisión en su propio retorno (spec, design, …)
gate del lane    → ratifica la decisión una sola vez (no se re-pregunta después)
sdd-apply        → materializa la decisión como EDR Accepted, en el mismo paso que el código
sdd-verify       → grupo decision-gaps: corre siempre, valida estructura + que el código corresponda
```

Mecanismo completo: [`in-flow-capture.md`](../../payload/domains/development/references/decision-capture/in-flow-capture.md)
(ruta desplegada: `~/.claude/references/decision-capture/in-flow-capture.md`).

El ejecutor `development-decisions-mine` sigue existiendo, pero solo para su **Mode A** — un scan
brownfield que invocás vos mismo sobre un repo con código y EDRs escasos o ausentes. No hay Mode B: nada
en el flujo de `development` arma una gap list post-verify ni lo despacha.

## Mecanismo de design: mine gate post-verify, fase por fase

```
tasks      → marca los gaps         (a qué DDR pertenece + cuáles no existen)
produce    → produce la pieza
verify     → confirma cuáles se implementaron (## Decision Gaps: implemented yes/no)
boundary   → el ORQUESTADOR lanza design-decisions-mine sobre los gaps implemented:yes
gate       → mostrar candidatos → confirm/edit/skip (thread principal)
materialize→ escribe los .md Inferred (+ INDEX) — solo lo confirmado
archive    → cierra (no registra DDRs)
```

1. **tasks** asigna `· ddr: <surface>/<slug>` (mapeado a un concern) a las tareas que tocan una
   decisión de marca, y chequea si el archivo existe. Una `· ddr:` cuyo DDR **no existe** = **gap
   dangling**. El conjunto de dangling = la gap list.
2. **verify** confirma, por cada gap, si `implemented: yes` — la tarea está `[x]` **y** su `criteria:`
   pasa contra el archivo Figma embarcado. Lo vuelca en `## Decision Gaps`.
3. **boundary verify→archive**: si hay ≥1 gap `implemented: yes`, el orquestador arma la gap list
   (scope) y **lanza `design-decisions-mine`**.
4. **fan-out**: si son muchos gaps, el orquestador parte la lista en batches y despacha varios miners
   en paralelo (cada uno mode-agnóstico: `scope → candidates[]`). Mergea y **deduplica por
   `surface/slug`**.
5. **gate** (thread principal): presenta candidatos ordenados por confianza, agrupados por surface, con
   acciones bulk (accept-all / por-surface / por-ítem). Nada se escribe sin confirmar.
6. **materialize**: escribe los `.md` Inferred confirmados y actualiza el INDEX (una vez).

## Detectar temprano, materializar tarde (design)

tasks marca el gap **sin la pieza todavía** (evidencia débil) → barato. El miner arma la evidencia
**con la pieza ya embarcada en Figma** (evidencia fuerte) → recién en el boundary. Por eso solo se
minan los `implemented: yes`: si una tarea no se completó o su `criteria:` no pasa, no hay pieza de
donde sacar evidencia.

## Quién hace qué (design)

- **tasks / verify**: producen los datos del gap (a qué DDR, y si se implementó). **No lanzan miners.**
- **orquestador** (hilo principal): el único que lanza sub-agentes. Lee `## Decision Gaps`, arma el
  scope, dispara el miner, mergea, presenta el gate, materializa.
- **miner (ejecutor)**: mode-agnóstico, recibe un scope, devuelve candidatos. No decide si correr — eso
  vive en quien lo invoca — ni escribe DDRs (eso es el thread principal post-confirm).

## No recomienda "siempre"

`design-decisions-mine` solo recomienda cuando hay un hueco **real**: la tarea tocó una decisión de
marca, esa decisión no tiene DDR, y verify confirmó que se implementó. Un cambio mecánico (sin
decisiones de marca) → cero gaps → silencio. Y aun con candidatos, **ofrece**: confirmás o saltás.

## Después: ratificar

Lo que cualquiera de los dos mecanismos materializa son **borradores `Inferred`** (sin porqué, en el
caso de design; ya con porqué en el caso de development, porque su ratificación ocurre en el gate del
lane, antes de materializar). Para convertir un `Inferred` de design en una decisión plena (`Accepted`),
corrés `design-decisions-bootstrap` en modo update, que te entrevista el porqué. Ver
[04 — ciclo de vida](04-decisiones-edr.md#el-ciclo-de-vida-de-una-decisión) (el ciclo aplica igual a
DDRs, con el término del dominio).
