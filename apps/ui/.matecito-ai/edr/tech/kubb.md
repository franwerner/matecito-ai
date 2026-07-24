# Kubb

- **Category:** Other
- **Version:** sin pinear
- **Status:** Accepted
- **Decided in phase:** security
- **Date:** 2026-07-23

## Por qué la elegimos

Genera los tipos TS y los schemas Zod desde el OpenAPI del broker, de modo que el contrato de la UI espeja el schema-first del broker sin escribirse a mano.

## Alternativas descartadas

- Ninguna evaluada formalmente.

## Notas

Usada en: security/input-validation (tipos/schemas generados) y delivery/ci-quality-gates (gate de sync generado vs OpenAPI).
