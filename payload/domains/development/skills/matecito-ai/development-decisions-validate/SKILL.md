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

**Activos:** `context` · `structure` · `runtime` · `data` · `observability` · `security` · `contracts` · `delivery` · `frontend` · `quality`
**Reservados:** `lifecycle` · `integration` · `privacy` · `release` · `domain-logic` · `compliance` · `ux-product`

Cualquier carpeta bajo `.matecito-ai/edr/` que no sea uno de estos dominios (ni `tech/`) es un hallazgo de integridad de taxonomía.

## Pre-flight

Leé `.matecito-ai/edr/INDEX.md`. Si no existe, no hay nada que validar → sugerí correr `development-decisions-bootstrap` y frená.

## Proceso

1. **Inventariá la estructura.** Listá todo: `find .matecito-ai/edr -name '*.md'`. Identificá el índice raíz, los índices de dominio (`<dominio>/INDEX.md`), los EDRs (`<dominio>/<slug>.md`) y `tech/INDEX.md` + `tech/*.md`.
2. **Identificá el tipo de proyecto** desde el EDR `context` (en `.matecito-ai/edr/context/context.md`).
3. **Para el chequeo de completitud** necesitás saber qué fases son relevantes a ese tipo:
   - Si el bootstrap te lanzó, usá la lista de fases relevantes que te pasó.
   - Si corrés standalone y podés acceder al catálogo `concerns/INDEX.md` de `development-decisions-bootstrap`, usalo.
   - Si no tenés ninguna de las dos, marcá completitud como "no verificable" y seguí con el resto (que solo necesita los EDRs).
   - **EDRs con `Status: Inferred` NO son decisiones cerradas:** no los contés en el total de decisiones tomadas para el chequeo de completitud (no satisfacen la preocupación), y no reportes como defecto las secciones WHY/Consecuencias/Reglas vacías (se espera que estén vacías). Sí considerá `## Alcance` como ancla de drift (verificá que los globs sigan matcheando).
4. **Leé `coherence-rules.md`** (en esta misma skill) y aplicá cada chequeo. Cada regla indica el/los dominio(s) donde viven los EDRs involucrados, así sabés qué archivos abrir.
5. **Emití el reporte** agrupado por dominio y, dentro de cada dominio, por severidad.

## Resolución de archivos por dominio

La rúbrica nombra los EDRs por su slug (`auth`, `layers-and-dependencies`, etc.). Para abrir el archivo, el dominio está en la tabla de mapeo de `coherence-rules.md` o en el campo `Dominio:` del encabezado del propio EDR. Ej: `auth` → `.matecito-ai/edr/security/auth.md`.

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
- ❌ Ignorar carpetas no canónicas o mismatches `domain`/carpeta → son hallazgos de integridad de taxonomía, repórtalos.
