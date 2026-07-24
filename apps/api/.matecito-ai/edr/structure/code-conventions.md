# EDR — Convenciones de código

- **Status:** Accepted
- **Type:** convention
- **Date:** 2026-07-23

## Contexto

El broker está escrito en Go, que trae de base un formateador y un analizador estático canónicos, lo que elimina gran parte de la discusión de estilo. Para un proyecto single-dev conviene apoyarse en ese set idiomático y sumar un enforcer que agrupe checks de estilo adicionales sin ceremonia.

## Decisión

Adoptamos el set idiomático de Go (formateo y análisis estático estándar ya de base) y usamos un linter agregador como enforcer del estilo por encima de ese piso.

## Reglas verificables

- **[tool: gofmt]** el código está formateado con gofmt.
- **[tool: go vet]** el código pasa go vet.
- **[tool: golangci-lint]** el código pasa golangci-lint.

## Alternativas consideradas

- **Solo gofmt + go vet sin linter agregador:** descartada por dejar fuera checks de estilo y correctness que golangci-lint agrupa sin costo de configuración relevante.

## Consecuencias

- Estilo consistente y verificable por herramienta, sin debate manual.
- Trade-off: el linter agregador introduce una dependencia de tooling que hay que mantener actualizada.
