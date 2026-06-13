# 05 — Auto-mine de ADRs (`flagDecisionGaps`)

[← 04 Decisiones y ADRs](04-decisiones-adr.md) · [Índice](README.md) · Siguiente: [06 — Herramientas →](06-herramientas.md)

> Esta página describe el auto-mine de **ADRs**, el decision record del dominio **development**. El mecanismo (detectar gaps in-flow, opt-in vía `flagDecisionGaps`, materializar `Inferred` con confirmación) es del núcleo; otro dominio lo aplica a su propio tipo de record (p. ej. auto-mine de DDR en design).

`mine` tiene dos modos, un mismo motor:

- **Mode A — scan brownfield**: lo invocás explícitamente sobre un repo. Escanea todo y propone ADRs `Inferred`.
- **Mode B — in-flow**: corre dentro del flujo SDD, opt-in vía el flag **`flagDecisionGaps`**. Detecta decisiones que el cambio implementó y no tienen ADR.

Esta página cubre el Mode B. (El concepto de mine como proponente está en [04](04-decisiones-adr.md).)

## El gate es la INTENCIÓN, no la presencia de ADRs

`flagDecisionGaps` es **opt-in, off por default**. Con el flag off: silencio total, el flujo se comporta como siempre. Con el flag on, mine funciona **aunque no exista `.matecito-ai/adr/`**: si no hay ADRs, toda decisión implementada es un hueco, y mine **bootstrapea los primeros**. Su única dependencia dura es el **catálogo de concerns** (viene con la skill).

Esto preserva el invariante de matecito-ai: **los ADRs nunca son obligatorios ni molestan** — detecta y ofrece; nunca impone.

## El flujo, fase por fase

```
tasks      → marca los gaps         (a qué ADR pertenece + cuáles no existen)
apply      → implementa
verify     → confirma cuáles se implementaron (## Decision Gaps: implemented yes/no)
boundary   → el ORQUESTADOR lanza N miners sobre los gaps implemented:yes
gate       → mostrar candidatos → confirm/edit/skip (thread principal)
materialize→ escribe los .md Inferred (+ INDEX) — solo lo confirmado
archive    → cierra (no registra ADRs)
```

1. **tasks** asigna `· adr: <dominio>/<slug>` (mapeado a un concern) a las tareas que tocan una decisión, y chequea si el archivo existe. Una `· adr:` cuyo ADR **no existe** = **gap dangling**. El conjunto de dangling = la gap list (sin campo extra). Solo con el flag on.
2. **verify** confirma, por cada gap, si `implemented: yes` — la tarea está `[x]` **y** su `criteria:` pasa en el código embarcado. Lo vuelca en `## Decision Gaps`.
3. **boundary verify→archive**: el **orquestador** evalúa el *mine gate*. Si `flagDecisionGaps` on **y** hay ≥1 gap `implemented: yes`, arma la gap list (scope) y **lanza los miners**.
4. **fan-out**: si son muchos gaps, el orquestador parte la lista en batches y despacha **N miners en paralelo** (cada uno mode-agnóstico: `scope → candidates[]`). Mergea y **deduplica por `dominio/slug`**.
5. **gate** (thread principal): presenta candidatos ordenados por confianza, agrupados por dominio, con acciones bulk (accept-all / por-dominio / por-ítem). Nada se escribe sin confirmar.
6. **materialize**: escribe los `.md` Inferred confirmados y actualiza el INDEX (una vez).

## Detectar temprano, materializar tarde

tasks marca el gap **sin código todavía** (evidencia débil) → barato. El miner arma la evidencia **con el código ya embarcado** (evidencia fuerte) → recién en el boundary. Por eso solo se minan los `implemented: yes`: si una tarea no se completó o su `criteria:` no pasa, no hay código de donde sacar evidencia.

## Quién hace qué

- **tasks / verify**: producen los datos del gap (a qué ADR, y si se implementó). **No lanzan miners.**
- **orquestador** (hilo principal): el único que lanza sub-agentes. Lee `## Decision Gaps`, arma el scope, dispara los N miners, mergea, presenta el gate, materializa. La regla vive en `CLAUDE.md` ("Decision-Gap Capture").
- **miner (ejecutor)**: mode-agnóstico, recibe un scope, devuelve candidatos. No lee el flag, no se ramifica por modo, no escribe ADRs (eso es el thread principal post-confirm).

## No recomienda "siempre"

Con el flag on, mine solo recomienda cuando hay un hueco **real**: la tarea tocó una decisión, esa decisión no tiene ADR, y verify confirmó que se implementó. Un cambio mecánico (sin decisiones) → cero gaps → silencio. Y aun con candidatos, **ofrece**: confirmás o saltás.

## Después: ratificar

Lo que mine materializa son **borradores `Inferred`** (sin porqué). Para convertirlos en decisiones plenas (`Accepted`), corrés **`bootstrap` en modo update**, que te entrevista el porqué. Ver [04 — ciclo de vida](04-decisiones-adr.md#el-ciclo-de-vida-de-una-decisión).
