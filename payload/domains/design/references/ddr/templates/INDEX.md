# Plantillas de estructura de DDR — Índice

Plantillas canónicas de la **estructura** de los artefactos de DDR en `.matecito-ai/ddr/`. Referencia consultable y **agnóstica de flujo**: cualquier productor de DDRs (la entrevista de decisiones de diseño `design-decisions-bootstrap`, o la minería desde Figma `design-decisions-mine`) las lee desde acá antes de materializar; el validador puede compararlas contra la salida real para detectar drift.

Separadas en archivos individuales (uno por plantilla) para que sean auditables y consultables de forma aislada.

| Plantilla | Salida que genera | Cuándo se usa |
|---|---|---|
| [ddr.md](ddr.md) | `.matecito-ai/ddr/<surface>/<slug>.md` | Un DDR (`Accepted`/`Pending`/`Deferred`/`Inferred`). |
| [index-root.md](index-root.md) | `.matecito-ai/ddr/INDEX.md` | Índice raíz: enruta por surface + surfaces sin uso. |
| [index-surface.md](index-surface.md) | `.matecito-ai/ddr/<surface>/INDEX.md` | Índice de cada surface con al menos un DDR. |

> A diferencia del dominio development, design NO tiene catálogo de tecnologías (`tech/`) ni escribe un `CLAUDE.md` raíz: los DDR viven solo como `.md` bajo `.matecito-ai/ddr/`.

**Leyenda de placeholders:** `<...>` = valor a completar al materializar. Los bloques `<!-- ... -->` son instrucciones del contrato (no van en el archivo generado).
