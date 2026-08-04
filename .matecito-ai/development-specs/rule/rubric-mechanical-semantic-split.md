# Capability — Reparto mecánico/semántico de las rúbricas de validación

- **Status:** Accepted
- **Date:** 2026-08-03
- **Components:** cli

## Propósito

Garantizar que los chequeos de coherencia entre semántica (decisiones de diseño, lenguaje sin ambigüedad) y mecánica (presencia de secciones, orden, sincronía de índices, taxonomía de carpetas) tengan dueños distintos y no se dupliquen. La mitad mecánica vive **solo** en un contrato de datos evaluado por script; la rúbrica retiene **solo** lo irreductiblemente semántico.

## Reglas de negocio

- Un chequeo decidible por inspección estructural —presencia, orden, enums de cabecera, sincronía índice↔archivo, links, sets declarados— MUST vivir SOLO en el contrato de datos (`artifact-checks/checks.yaml`) y MUST NOT repetirse en la rúbrica (`coherence-rules.md`). La evaluación la hace el script; la rúbrica no re-deriva.
- Un chequeo que requiere juicio semántico —contradicción de significado, lenguaje vago, calco de identificadores volátiles, clasificación de dominio— vive SOLO en la rúbrica. El script no lo puede evaluar.
- Un hallazgo mecánico **nunca lleva dos severidades** — una sola severidad por un mapeo declarado UNA vez en el contrato. El script traduce `severity` (crudo: `error|warning|nota`) a `display` (público: `CRITICAL|WARNING|SUGGESTION`) una sola vez.
- El **ratchet de la mitad mecánica** es un dato simple: agregar un chequeo significa una fila más del `checks.yaml`. El ratchet de la semántica también es simple: una línea más en la rúbrica. Ambos crecen sin tocar código.

## Escenarios

### Scenario: La skill no re-deriva lo mecánico

- **GIVEN** un store con un índice desincronizado (una carpeta sin INDEX.md) y una sección faltante (presente/orden)
- **WHEN** la skill valida ese store
- **THEN** ambos hallazgos vienen del script (`validate-artifact.js`); la skill los imprime sin re-derivar; nunca re-calcula sync ni orden

### Scenario: Un hallazgo, una severidad

- **GIVEN** un hallazgo mecánico emitido por el script (e.g., `SECTION-MISSING` con `severity: error`)
- **WHEN** la skill lo incorpora a su reporte
- **THEN** su severidad en el reporte es exactamente la que el script declaró y tradujo via `severity_map`; nunca aparece con dos valores, nunca se endurece ni se suaviza

### Scenario: Ratchet de la mitad que lo conserve

- **GIVEN** el contrato `checks.yaml` con un `kind` que la mitad declaró como ratchet-able (ya tiene al menos una fila)
- **WHEN** el usuario quiere agregar un chequeo de ese `kind`
- **THEN** alcanza con agregar una fila al contrato; no se toca la lógica del script; el script existente evalúa la nueva fila

## Referencias

- **Capability** → [`../process/validate-artifact-structure.md`](../process/validate-artifact-structure.md) — El validador que ejecuta el contrato de datos.
- **Rule** → [`./single-rooted-spec-store.md`](./single-rooted-spec-store.md) — Dónde viven los specs y las rúbricas.
