---
name: development-decisions-validate
description: Validador de coherencia, completitud y verificabilidad de las decisiones de ingeniería (EDRs) de un proyecto, organizadas por dominio en .matecito-ai/edr/<dominio>/. Usá esta skill cuando el usuario pida "validar la arquitectura", "revisar los EDRs", "chequear coherencia", "¿mis decisiones se contradicen?", después de editar EDRs a mano o de correr development-decisions-bootstrap. Lee `.matecito-ai/edr/` (recursivo por dominio) y reporta hallazgos con severidad. NO modifica nada — es consultiva.
---

# Development Decisions Validate

Lee los EDRs producidos por `development-decisions-bootstrap` (o editados a mano) y los chequea contra una rúbrica: **completitud**, **coherencia entre decisiones**, **verificabilidad**, e **integridad de la taxonomía**. Reporta hallazgos con severidad. No modifica archivos.

Los EDRs están organizados por dominio (`.matecito-ai/edr/<dominio>/<slug>.md`), con un índice raíz (`.matecito-ai/edr/INDEX.md`) y un índice por dominio (`.matecito-ai/edr/<dominio>/INDEX.md`).

## Por qué contexto fresco

Esta validación SOLO sirve si es adversarial: leé únicamente lo que está ESCRITO en los EDRs. No asumas la intención del autor ni el contexto de cómo se tomó cada decisión. Lo que no está en el archivo, no existe. Por eso esta skill corre con contexto limpio (standalone, o lanzada por el bootstrap como sub-agente).

## Dominios canónicos

La taxonomía es fija (la misma que impone `development-decisions-bootstrap`):

**Activos:** `structure` · `runtime` · `data` · `observability` · `security` · `contracts` · `delivery` · `frontend` · `quality`
**Reservados:** `lifecycle` · `integration` · `privacy` · `release` · `domain-logic` · `compliance` · `ux-product`

Cualquier carpeta bajo `.matecito-ai/edr/` que no sea uno de estos dominios (ni `tech/`) es un hallazgo de integridad de taxonomía — lo detecta el motor mecánico (`folder-taxonomy` en `checks.yaml`), no vos; ver "Motor mecánico" más abajo.

## Pre-flight

Leé `.matecito-ai/edr/INDEX.md`. Si no existe, no hay nada que validar → sugerí correr `development-decisions-bootstrap` y frená.

## Motor mecánico (`validate-artifact.js`)

La mitad mecánica de la rúbrica —presencia/vacío de sección, sincronía de índices, taxonomía de carpetas, links colgados, drift de `## Alcance`, status deprecado, ubicación por slug— vive en `~/.claude/references/artifact-checks/checks.yaml` y la evalúa el motor. **No la re-derives**: invocá el script e ingerí su salida tal cual.

```
node ~/.claude/scripts/validate-artifact.js --type edr --store .matecito-ai/edr --root <raíz del repo>
```

`--root` no tiene default — pasalo siempre con la raíz real del repo (no `.matecito-ai/edr`). Si algún chequeo resuelto lo necesita y falta, el script sale con exit 2 nombrando cuáles; en ese caso no hay JSON que ingerir, solo el mensaje en stderr.

La salida es un único JSON (`findings[]`, `counts`, `skipped[]`). Cada `finding` ya trae `severity` **y** `display` (`CRITICAL`/`WARNING`/`SUGGESTION`) — **imprimí el `display` que recibís, no lo traduzcas ni lo re-deduzcas**. `skipped[]` puede traer chequeos no evaluados por falta de `.matecito-ai/config.json`/`repo.components` bajo `--root`; mencionalos en el reporte, no los silencies.

## Proceso

1. **Inventariá la estructura.** Listá todo: `find .matecito-ai/edr -name '*.md'`. Identificá el índice raíz, los índices de dominio (`<dominio>/INDEX.md`), los EDRs (`<dominio>/<slug>.md`) y `tech/INDEX.md` + `tech/*.md`.
2. **Identificá el tipo de proyecto** desde el parámetro que te pasa el invocador (el bootstrap lo pasa al lanzarte). Si corrés standalone, inferilo del catálogo `tech/` y la estructura del repo, o preguntalo.
3. **Para el chequeo de completitud** necesitás saber qué fases son relevantes a ese tipo:
   - Si el bootstrap te lanzó, usá la lista de fases relevantes que te pasó.
   - Si corrés standalone y podés acceder al catálogo `concerns/INDEX.md` de `development-decisions-bootstrap`, usalo.
   - Si no tenés ninguna de las dos, marcá completitud como "no verificable" y seguí con el resto (que solo necesita los EDRs).
   - **EDRs con `Status: Inferred` NO son decisiones cerradas:** no los contés en el total de decisiones tomadas para el chequeo de completitud (no satisfacen la preocupación), y no reportes como defecto las secciones WHY/Consecuencias/Reglas vacías (se espera que estén vacías). Sí considerá `## Alcance` como ancla de drift (verificá que los globs sigan matcheando) — este último ya lo cubre el motor mecánico (`glob-drift`); no lo repitas a mano.
4. **Corré el motor mecánico** (ver arriba) e ingerí su JSON.
5. **Leé `coherence-rules.md`** (en esta misma skill, ya recortada a lo semántico) y aplicá cada chequeo restante sobre los EDRs. Cada regla indica el/los dominio(s) donde viven los EDRs involucrados, así sabés qué archivos abrir.
6. **Emití el reporte** combinando los `findings` del motor (con su `display` tal cual) y los hallazgos semánticos, agrupado por dominio y, dentro de cada dominio, por severidad.

## Resolución de archivos por dominio

La rúbrica nombra los EDRs por su slug (`auth`, `layers-and-dependencies`, etc.). El dominio de un EDR sale de la carpeta donde vive (`<dominio>/<slug>.md`) — el template canónico no tiene campo `Dominio:` en el header. Para localizar el archivo de un slug nombrado en `coherence-rules.md`, usá su tabla de mapeo. Ej: `auth` → `.matecito-ai/edr/security/auth.md`. Un EDR cuya carpeta no corresponde a ese mapeo es un hallazgo del motor mecánico (`location-by-slug`), no algo que detectes vos.

Las contradicciones **entre dominios** (ej: `privacy` vs `lifecycle`) requieren abrir EDRs de carpetas distintas — usá el mapeo para localizarlos.

## Formato del reporte

Agrupado **por dominio**, y dentro de cada dominio **por severidad**. Cerrá con una sección de hallazgos **cross-dominio** (contradicciones que involucran EDRs de más de un dominio) y un veredicto final.

```
## Dominio: security
🔴 CRITICAL — <qué> · EDRs: <cuáles> · <por qué> · <sugerencia>
🟡 WARNING  — ...
🔵 SUGGESTION — ...
(si un nivel no tiene hallazgos en el dominio, omitilo)

## Cross-dominio
🔴/🟡/🔵 — hallazgos que cruzan dominios (ej: privacy ↔ lifecycle)

## Veredicto
<una línea: ej "2 CRITICAL, 1 WARNING — resolver los CRITICAL antes de codear">
```

Leyenda de severidad:

```
🔴 CRITICAL — contradicen la arquitectura/decisiones; hay que resolverlas.
🟡 WARNING  — inconsistencias o riesgo de pudrición.
🔵 SUGGESTION — mejoras de claridad o robustez.
```

Si un dominio no tiene ningún hallazgo, no lo listes. Si NO hay hallazgos en ningún lado, decilo explícitamente y dá un veredicto verde.

## Después del reporte

- **No modifiques EDRs.** Si el usuario quiere resolver un hallazgo, derivá a `development-decisions-bootstrap` en modo update para el EDR afectado.
- **Ratchet:** si detectaste una contradicción real que NO está en `coherence-rules.md`, ofrecé agregarla a la rúbrica para que se atrape en el futuro (con su severidad, dominio(s) y mensaje qué/por qué/sugerencia).

## Anti-patterns

- ❌ Inferir intención no escrita para "salvar" una contradicción → si no está en el EDR, es un hallazgo.
- ❌ Modificar o arreglar EDRs vos mismo → solo reportás; el usuario decide y resuelve vía update.
- ❌ Reportar todo como CRITICAL → reservá CRITICAL para lo que rompe la arquitectura; usá WARNING/SUGGESTION para el resto.
- ❌ Buscar EDRs con glob plano (`.matecito-ai/edr/*.md`) → los EDRs están en subcarpetas de dominio; recorré recursivo.
- ❌ Ignorar carpetas no canónicas o EDRs en el dominio equivocado según el catálogo slug→dominio → son hallazgos de integridad/ubicación que el motor mecánico ya detecta (`folder-taxonomy`, `location-by-slug`); repórtalos, no los descartes.
- ❌ Re-derivar a mano un chequeo que ya cubre el motor mecánico (sincronía de índices, taxonomía de carpetas, links colgados, drift de `## Alcance`, secciones vacías, status deprecado, ubicación por slug) → invocá `validate-artifact.js` e ingerí su JSON; no lo dupliques leyendo los archivos vos mismo.
- ❌ Traducir o re-clasificar la `severity`/`display` que devuelve el motor → imprimí el `display` recibido tal cual.
