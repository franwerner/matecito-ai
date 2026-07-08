# Plantillas de estructura de EDR — Índice

Plantillas canónicas de la **estructura** de los artefactos de EDR en `.matecito-ai/edr/`. Referencia consultable y **agnóstica de flujo**: cualquier productor de EDRs (la entrevista de decisiones, o la minería desde código) las lee desde acá antes de materializar; el validador puede compararlas contra la salida real para detectar drift.

Separadas en archivos individuales (uno por plantilla) para que sean auditables y consultables de forma aislada.

| Plantilla | Salida que genera | Cuándo se usa |
|---|---|---|
| [edr.md](edr.md) | `.matecito-ai/edr/<dominio>/<slug>.md` | Un EDR (`Accepted`/`Pending`/`Deferred`/`Inferred`). |
| [index-root.md](index-root.md) | `.matecito-ai/edr/INDEX.md` | Índice raíz: enruta por dominio + dominios sin uso. |
| [index-domain.md](index-domain.md) | `.matecito-ai/edr/<dominio>/INDEX.md` | Índice de cada dominio con al menos un EDR. |
| [tech-edr.md](tech-edr.md) | `.matecito-ai/edr/tech/<nombre>.md` | Mini-EDR por tecnología concreta elegida. |
| [tech-index.md](tech-index.md) | `.matecito-ai/edr/tech/INDEX.md` | Catálogo de tecnologías por categoría. |

> El template del `CLAUDE.md` raíz del proyecto NO vive acá: es propio de `development-decisions-bootstrap` (`templates/claude-md.md`), porque solo bootstrap escribe ese archivo.

**Leyenda de placeholders:** `<...>` = valor a completar al materializar. Los bloques `<!-- ... -->` son instrucciones del contrato (no van en el archivo generado).
