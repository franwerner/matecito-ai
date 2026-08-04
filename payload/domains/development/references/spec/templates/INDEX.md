# Plantillas de estructura de capability-spec — Índice

Plantillas canónicas de la **estructura** de los artefactos de capability-spec en
`.matecito-ai/development-specs/`. Referencia consultable y **agnóstica de flujo**: cualquier
productor de capability-specs (la entrevista de `development-spec-bootstrap`, o la minería desde
código de `development-spec-mine`) las lee desde acá antes de materializar; el validador puede
compararlas contra la salida real para detectar drift.

Separadas en archivos individuales (uno por plantilla) para que sean auditables y consultables de
forma aislada.

| Plantilla | Salida que genera | Cuándo se usa |
|---|---|---|
| [capability.md](capability.md) | `.matecito-ai/development-specs/<type>/<capability>.md` | Un capability-spec (`Inferred`/`Draft`/`Accepted`/`Deprecated`), cualquiera de los cuatro tipos. |
| [capability.yaml](capability.yaml) | — (no es un artefacto) | La **forma** machine-checkable de `capability.md`: qué campos lleva el `--data`, qué secciones son opcionales (todas — presence-based, ninguna gateada por `Status` salvo `Reemplazado por`), y las filas de INDEX que le corresponden. La lee `~/.claude/scripts/render-artifact.js` (`--type capability-spec`) y `~/.claude/scripts/validate-artifact.js` (`--type capability-spec`) — nunca un humano a mano. `capability.md` es el significado; `capability.yaml` es la forma; deben mantenerse fieles entre sí. |
| [index-root.md](index-root.md) | `.matecito-ai/development-specs/INDEX.md` | Índice raíz: enruta por tipo + tipos sin uso. |
| [index-type.md](index-type.md) | `.matecito-ai/development-specs/<type>/INDEX.md` | Índice de cada tipo con al menos un capability-spec. |
| [../../artifact-checks/checks.yaml](../../artifact-checks/checks.yaml) | — (no es un artefacto) | Chequeos mecánicos entre archivos (índice↔archivo, links colgados, taxonomía de carpetas, y más) que `validate-artifact.js --store` evalúa sobre el store de capability-specs. Ver [`../../artifact-checks/README.md`](../../artifact-checks/README.md). |

**Leyenda de placeholders:** `<...>` = valor a completar al materializar. Los bloques `<!-- ... -->`
son instrucciones del contrato (no van en el archivo generado).
