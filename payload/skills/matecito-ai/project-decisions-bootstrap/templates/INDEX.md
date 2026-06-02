# Templates de salida — Índice

Templates canónicos de lo que la skill **materializa** en el repo objetivo (`.matecito-ai/adr/` + `CLAUDE.md`). Cada archivo es el contrato de un artefacto de salida: la fase de Materialización de `SKILL.md` los lee desde acá, y el validador `project-decisions-validate` puede compararlos contra la salida real para detectar drift.

Separados en archivos individuales (uno por template) para que sean auditables y consultables de forma aislada, sin parsear todo el `SKILL.md`.

| Template | Salida que genera | Cuándo se usa |
|---|---|---|
| [adr.md](adr.md) | `.matecito-ai/adr/<dominio>/<slug>.md` | Un ADR por fase tratada (`Accepted`/`Pending`/`Deferred`). |
| [index-root.md](index-root.md) | `.matecito-ai/adr/INDEX.md` | Índice raíz: enruta por dominio + dominios sin uso. |
| [index-domain.md](index-domain.md) | `.matecito-ai/adr/<dominio>/INDEX.md` | Índice de cada dominio con al menos un ADR. |
| [tech-adr.md](tech-adr.md) | `.matecito-ai/adr/tech/<nombre>.md` | Mini-ADR por tecnología concreta elegida (intercalado). |
| [tech-index.md](tech-index.md) | `.matecito-ai/adr/tech/INDEX.md` | Catálogo de tecnologías por categoría. |
| [claude-md.md](claude-md.md) | `CLAUDE.md` (raíz del proyecto) | Puntero mínimo al índice de ADRs. NO sobrescribir uno existente sin permiso. |

**Leyenda de placeholders:** `<...>` = valor a completar al materializar. Los bloques `<!-- ... -->` son instrucciones del contrato (no van en el archivo generado).
