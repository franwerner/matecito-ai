# Plantillas de estructura de EDR — Índice

Plantillas canónicas de la **estructura** de los artefactos de EDR en `.matecito-ai/edr/`. Referencia consultable y **agnóstica de flujo**: cualquier productor de EDRs (la entrevista de decisiones, o la minería desde código) las lee desde acá antes de materializar; el validador puede compararlas contra la salida real para detectar drift.

Separadas en archivos individuales (uno por plantilla) para que sean auditables y consultables de forma aislada.

| Plantilla | Salida que genera | Cuándo se usa |
|---|---|---|
| [edr.md](edr.md) | `.matecito-ai/edr/<dominio>/<slug>.md` | Un EDR (`Accepted`/`Pending`/`Deferred`/`Inferred`). |
| [edr.yaml](edr.yaml) | — (no es un artefacto) | La **forma** machine-checkable de `edr.md`: qué campos lleva el `--data`, qué secciones son siempre/condicionales/opcionales por `Status`, y las filas de INDEX que le corresponden. La lee `~/.claude/scripts/render-artifact.js` (`--type edr`) y `~/.claude/scripts/validate-artifact.js` (`--type edr`) — nunca un humano a mano. `edr.md` es el significado; `edr.yaml` es la forma; deben mantenerse fieles entre sí. |
| [index-root.md](index-root.md) | `.matecito-ai/edr/INDEX.md` | Índice raíz: enruta por dominio + dominios sin uso. |
| [index-domain.md](index-domain.md) | `.matecito-ai/edr/<dominio>/INDEX.md` | Índice de cada dominio con al menos un EDR. |
| [tech-edr.md](tech-edr.md) | `.matecito-ai/edr/tech/<nombre>.md` | Mini-EDR por tecnología concreta elegida. |
| [tech-index.md](tech-index.md) | `.matecito-ai/edr/tech/INDEX.md` | Catálogo de tecnologías por categoría. |
| [../../artifact-checks/checks.yaml](../../artifact-checks/checks.yaml) | — (no es un artefacto) | Chequeos mecánicos entre archivos (índice↔archivo, links colgados, taxonomía de carpetas, y más) que `validate-artifact.js --store` evalúa sobre el store de EDRs. Ver [`../../artifact-checks/README.md`](../../artifact-checks/README.md). |

`tech-edr.md` / `tech-index.md` no tienen contrato `.yaml`: mine nunca produce EDRs de tecnología, así que no hay nada que renderizar ni validar ahí.

> El template del `CLAUDE.md` raíz del proyecto NO vive acá: es propio de `development-decisions-bootstrap` (`templates/claude-md.md`), porque solo bootstrap escribe ese archivo.

**Leyenda de placeholders:** `<...>` = valor a completar al materializar. Los bloques `<!-- ... -->` son instrucciones del contrato (no van en el archivo generado).
