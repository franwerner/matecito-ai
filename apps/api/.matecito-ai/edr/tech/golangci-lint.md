# golangci-lint

- **Category:** Other
- **Version:** sin pinear
- **Status:** Accepted
- **Decided in phase:** structure
- **Date:** 2026-07-23

## Por qué la elegimos

Enforcer del estilo por encima del piso de gofmt + go vet: agrupa múltiples linters de estilo y correctness en una sola corrida, sin costo de configuración relevante.

## Alternativas descartadas

- **Solo gofmt + go vet:** dejan fuera checks de estilo/correctness que el agregador cubre.

## Notas

Usada en: structure/code-conventions.
