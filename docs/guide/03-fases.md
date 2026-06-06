# 03 — Las fases

[← 02 El flujo SDD](02-flujo-sdd.md) · [Índice](README.md) · Siguiente: [04 — Decisiones y ADRs →](04-decisiones-adr.md)

Cada fase es un **sub-agente de contexto fresco** que el orquestador despacha. Lee sus artefactos de Engram (por topic key), hace su trabajo, y escribe el suyo. Acá, fase por fase: qué hace, qué herramienta engancha, y qué deja para la siguiente.

## intake _(base, entrada)_

Estructura el pedido crudo. Hace 2-4 preguntas de descubrimiento, clasifica el cambio (tipo/tamaño), recomienda la **lane**, y decide si amerita **diagrama** y/o **verificación de UI**. Si los ADRs están activos, corre una **guardia temprana**: frena (`blocked`) si el pedido choca con un ADR `Accepted`, o deriva a `bootstrap` (`needs-decision`) si hay una pregunta arquitectónica sin decidir.
→ Produce el **brief**. El orquestador lo muestra en el **INTAKE GATE** y espera confirmación.

## explore _(add-on)_

Investiga el código antes de comprometer un diseño: arquitectura actual, comparación de enfoques. Engancha **codegraph** para explorar por estructura (no archivo por archivo).
→ Produce `explore`.

## propose _(add-on)_

Formaliza la idea en una propuesta (intención, alcance, enfoque).
→ Produce `proposal`.

## spec _(base)_

Escribe los requisitos con escenarios verificables (Given/When/Then). Es el contrato observable del cambio. Si no hubo `proposal`, arranca del brief de intake.
→ Produce `spec`.

## design _(add-on)_

Decide el enfoque técnico. **Lee los ADRs vigentes** y alinea el diseño a los `Accepted` (es un consumidor de ADRs). Si el cambio amerita **diagrama** (marcado en intake), el **thread principal** lo renderiza en vivo (efímero, `localhost:6002`); **no se exporta ningún archivo** — el sub-agente design solo lo señala (es headless, no previsualiza).
→ Produce `design`.

## tasks _(add-on)_

Descompone el cambio en tareas atómicas y ordenadas. Cada tarea lleva:
- un **`criteria:`** verificable (obligatorio) — lo que `verify` chequea después;
- un **`· adr: <dominio>/<slug>`** opcional — solo si la tarea toca una decisión.

Si está activo el auto-mine, una `· adr:` cuyo ADR **no existe** queda como **gap** (ver [05](05-auto-mine.md)).
→ Produce `tasks`.

## apply _(base)_

Implementa el código siguiendo los patrones existentes; marca cada tarea `[x]` a medida que avanza. Engancha **context7** (docs de librerías al día) y **codegraph** (estructura). Si está activo **Strict TDD**, sigue el ciclo test-first.
→ Produce `apply-progress` (mergea con batches previos).

## verify _(base)_

Valida la implementación contra spec/design/tasks. Corre tests/build/coverage cuando están disponibles. Sobre ADRs hace **dos chequeos distintos**:
1. **Cumplimiento** de los ADRs `Accepted` que el cambio tocó → violación = `ADR-VIOLATION` (CRITICAL).
2. **Confirmación de gaps** (si auto-mine activo): por cada `· adr:` dangling de tasks, marca `implemented: yes/no` (tarea completa **y** su `criteria:` pasa) en una sección `## Decision Gaps`.

Si el cambio toca UI y **proofshot** está disponible, conduce el browser y valida los escenarios.
→ Produce `verify-report`.

## archive _(base)_

Cierra el cambio: persiste el reporte final con las observation IDs para trazabilidad, marca el estado archivado. **No registra ADRs** (viven solo en `.md`).
→ Produce `archive-report`.

---

Cada fase corre con su **modelo configurable** (ver [07](07-configuracion.md)). La resolución de modelo / Strict TDD / auto-mine la hace el **orquestador** antes de despachar, y se la pasa al sub-agente ya resuelta — el ejecutor no lee config.
