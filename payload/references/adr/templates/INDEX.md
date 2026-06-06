# Plantillas de estructura de ADR — Índice

Plantillas canónicas de la **estructura** de los artefactos de ADR en `.matecito-ai/adr/`. Referencia consultable y **agnóstica de flujo**: cualquier productor de ADRs (la entrevista de decisiones, o la minería desde código) las lee desde acá antes de materializar; el validador puede compararlas contra la salida real para detectar drift.

Separadas en archivos individuales (uno por plantilla) para que sean auditables y consultables de forma aislada.

| Plantilla | Salida que genera | Cuándo se usa |
|---|---|---|
| [adr.md](adr.md) | `.matecito-ai/adr/<dominio>/<slug>.md` | Un ADR (`Accepted`/`Pending`/`Deferred`/`Inferred`). |
| [index-root.md](index-root.md) | `.matecito-ai/adr/INDEX.md` | Índice raíz: enruta por dominio + dominios sin uso. |
| [index-domain.md](index-domain.md) | `.matecito-ai/adr/<dominio>/INDEX.md` | Índice de cada dominio con al menos un ADR. |
| [tech-adr.md](tech-adr.md) | `.matecito-ai/adr/tech/<nombre>.md` | Mini-ADR por tecnología concreta elegida. |
| [tech-index.md](tech-index.md) | `.matecito-ai/adr/tech/INDEX.md` | Catálogo de tecnologías por categoría. |

> El template del `CLAUDE.md` raíz del proyecto NO vive acá: es propio de `project-decisions-bootstrap` (`templates/claude-md.md`), porque solo bootstrap escribe ese archivo.

**Leyenda de placeholders:** `<...>` = valor a completar al materializar. Los bloques `<!-- ... -->` son instrucciones del contrato (no van en el archivo generado).
