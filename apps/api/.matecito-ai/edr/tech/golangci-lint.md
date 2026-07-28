# golangci-lint

- **Category:** Other
- **Version:** v2.12.2
- **Status:** Accepted
- **Decided in phase:** structure
- **Date:** 2026-07-23

## Por qué la elegimos

Enforcer del estilo por encima del piso del formateador y el vet estándar: agrupa múltiples linters de estilo y correctness en una sola corrida, sin costo de configuración relevante.

## Alternativas descartadas

- **Solo el formateador y el vet estándar:** dejan fuera checks de estilo/correctness que el agregador cubre.

## Notas

Usada en: structure/code-conventions, delivery/ci-quality-gates.

La versión está pineada: sin pin, una release nueva puede romper el gate sin que nadie haya tocado el código, que es lo contrario de lo que se le pide a una verificación determinista.

Se compila desde fuente en vez de bajar el binario publicado. Los binarios de release se construyen con una versión del lenguaje anterior a la que este módulo declara, y la herramienta se niega a analizar un módulo cuya versión de lenguaje es más nueva que la suya.

