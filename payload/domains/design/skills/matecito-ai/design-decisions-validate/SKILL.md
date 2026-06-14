---
name: design-decisions-validate
description: Validador de coherencia, completitud y verificabilidad de las decisiones de DISEÑO (DDRs) de un proyecto, organizadas por surface en .matecito-ai/ddr/<surface>/. Usá esta skill cuando el usuario pida "validar el sistema de diseño", "revisar los DDRs", "chequear coherencia de diseño", "¿mis decisiones de diseño se contradicen?", después de editar DDRs a mano o de correr design-decisions-bootstrap. Lee `.matecito-ai/ddr/` (recursivo por surface) y reporta hallazgos con severidad. NO modifica nada — es consultiva.
---

# Design Decisions Validate

Lee los DDRs producidos por `design-decisions-bootstrap` (o editados a mano, o drafteados como `Inferred` por `design-decisions-mine`) y los chequea contra una rúbrica: **completitud**, **coherencia entre decisiones**, **verificabilidad**, e **integridad de la taxonomía**. Reporta hallazgos con severidad. No modifica archivos.

Los DDRs están organizados por surface (`.matecito-ai/ddr/<surface>/<slug>.md`), con un índice raíz (`.matecito-ai/ddr/INDEX.md`) y un índice por surface (`.matecito-ai/ddr/<surface>/INDEX.md`).

## Por qué contexto fresco

Esta validación SOLO sirve si es adversarial: leé únicamente lo que está ESCRITO en los DDRs. No asumas la intención del autor ni el contexto de cómo se tomó cada decisión. Lo que no está en el archivo, no existe. Por eso esta skill corre con contexto limpio (standalone, o lanzada por el bootstrap como sub-agente).

## Surfaces canónicas

La taxonomía es fija (la misma que impone `design-decisions-bootstrap`):

**Activas:** `foundation` · `components` · `layout` · `brand` · `accessibility`

Cualquier carpeta bajo `.matecito-ai/ddr/` que no sea una de estas surfaces es un hallazgo de integridad de taxonomía.

## Pre-flight

Leé `.matecito-ai/ddr/INDEX.md`. Si no existe, no hay nada que validar → sugerí correr `design-decisions-bootstrap` y frená.

## Proceso

1. **Inventariá la estructura.** Listá todo: `find .matecito-ai/ddr -name '*.md'`. Identificá el índice raíz, los índices de surface (`<surface>/INDEX.md`) y los DDRs (`<surface>/<slug>.md`).
2. **Identificá el tipo de pieza** desde el contexto disponible (el INDEX raíz suele anotarlo; si no, miralo en los DDRs de `brand`/`foundation`).
3. **Para el chequeo de completitud** necesitás saber qué fases (concerns) son relevantes a ese tipo de pieza:
   - Si el bootstrap te lanzó, usá la lista de fases relevantes que te pasó.
   - Si corrés standalone y podés acceder al catálogo `concerns/INDEX.md` de `design-decisions-bootstrap` (la matriz de aplicabilidad por tipo de pieza), usalo.
   - Si no tenés ninguna de las dos, marcá completitud como "no verificable" y seguí con el resto (que solo necesita los DDRs).
   - **DDRs con `Status: Inferred` NO son decisiones cerradas:** no los contés en el total de decisiones tomadas para el chequeo de completitud (no satisfacen la preocupación), y no reportes como defecto las secciones Contexto/Decisión/Consecuencias/Reglas verificables vacías ni el porqué vacío (se espera que estén vacías hasta que el modo update los ratifique a `Accepted`). Sí considerá `## Alcance` como ancla de drift (verificá que el locator siga existiendo en Figma).
4. **Leé `coherence-rules.md`** (en esta misma skill) y aplicá cada chequeo. Cada regla indica la/las surface(s) donde viven los DDRs involucrados, así sabés qué archivos abrir.
5. **Emití el reporte** agrupado por surface y, dentro de cada surface, por severidad.

## Resolución de archivos por surface

La rúbrica nombra los DDRs por su slug (`color-palette`, `type-scale`, etc.). Para abrir el archivo, la surface está en la tabla de mapeo de `coherence-rules.md` o en la ruta del propio DDR (`.matecito-ai/ddr/<surface>/<slug>.md`). Ej: `color-palette` → `.matecito-ai/ddr/foundation/color-palette.md`.

Las contradicciones **entre surfaces** (ej: `foundation` vs `accessibility`) requieren abrir DDRs de carpetas distintas — usá el mapeo para localizarlos.

## Formato del reporte

Agrupado **por surface**, y dentro de cada surface **por severidad**. Cerrá con una sección de hallazgos **cross-surface** (contradicciones que involucran DDRs de más de una surface) y un veredicto final.

```
## Surface: foundation
🔴 CRITICAL — <qué> · DDRs: <cuáles> · <por qué> · <sugerencia>
🟡 WARNING  — ...
🔵 SUGGESTION — ...
(si un nivel no tiene hallazgos en la surface, omitilo)

## Cross-surface
🔴/🟡/🔵 — hallazgos que cruzan surfaces (ej: foundation ↔ accessibility)

## Veredicto
<una línea: ej "2 CRITICAL, 1 WARNING — resolver los CRITICAL antes de producir">
```

Leyenda de severidad:

```
🔴 CRITICAL — contradicen el sistema/decisiones de diseño; hay que resolverlas.
🟡 WARNING  — inconsistencias o riesgo de pudrición.
🔵 SUGGESTION — mejoras de claridad o robustez.
```

Si una surface no tiene ningún hallazgo, no la listes. Si NO hay hallazgos en ningún lado, decilo explícitamente y dá un veredicto verde.

## Después del reporte

- **No modifiques DDRs.** Si el usuario quiere resolver un hallazgo, derivá a `design-decisions-bootstrap` en modo update para el DDR afectado.
- **Ratchet:** si detectaste una contradicción real que NO está en `coherence-rules.md`, ofrecé agregarla a la rúbrica para que se atrape en el futuro (con su severidad, surface(s) y mensaje qué/por qué/sugerencia).

## Anti-patterns

- ❌ Inferir intención no escrita para "salvar" una contradicción → si no está en el DDR, es un hallazgo.
- ❌ Modificar o arreglar DDRs vos mismo → solo reportás; el usuario decide y resuelve vía bootstrap modo update.
- ❌ Reportar todo como CRITICAL → reservá CRITICAL para lo que rompe el sistema de diseño; usá WARNING/SUGGESTION para el resto.
- ❌ Buscar DDRs con glob plano (`.matecito-ai/ddr/*.md`) → los DDRs están en subcarpetas de surface; recorré recursivo.
- ❌ Tratar un `Inferred` como decisión tomada, o reportar su porqué/reglas vacías como defecto → es un borrador, no cuenta para completitud y se espera vacío.
- ❌ Reportar como defecto verificabilidad ausente en un `Inferred` → las "Reglas verificables" solo aplican a `Accepted`.
- ❌ Ignorar carpetas no canónicas (surfaces fuera de la taxonomía fija) → son hallazgos de integridad de taxonomía, repórtalos.
- ❌ Confundir drift DDR-vs-Figma (trabajo de `design-decisions-mine`) con coherencia interna entre DDRs (lo de esta skill) → acá chequeás que los DDRs no se contradigan entre sí y que sean verificables; el `## Alcance` solo se usa como ancla de completitud/drift cuando tenés acceso al archivo Figma.
