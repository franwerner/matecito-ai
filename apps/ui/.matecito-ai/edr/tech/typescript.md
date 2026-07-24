# TypeScript

- **Category:** Lenguaje
- **Version:** sin pinear
- **Status:** Accepted
- **Decided in phase:** structure
- **Date:** 2026-07-23

## Por qué la elegimos

Lenguaje de toda la UI: tipado estático estricto que sostiene las convenciones de código (sin `any`, ausencia con `T | null`, union de literales) y habilita los tipos generados desde el OpenAPI del broker.

## Alternativas descartadas

- Ninguna evaluada formalmente: elección de base para la UI.

## Notas

Usada en: structure/code-conventions, delivery/ci-quality-gates (gate `tsc --noEmit`) y security/input-validation (tipos generados).
