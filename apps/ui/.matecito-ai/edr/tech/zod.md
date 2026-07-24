# Zod

- **Category:** Other
- **Version:** sin pinear
- **Status:** Accepted
- **Decided in phase:** security
- **Date:** 2026-07-23

## Por qué la elegimos

Validación por schema en runtime: parsea lo entrante del broker en el borde de datos (los schemas los genera Kubb desde el OpenAPI) y valida la config al startup. A futuro, validación de forms en la fase de escritura.

## Alternativas descartadas

- Ninguna evaluada formalmente.

## Notas

Usada en: security/input-validation (parseo del envelope/snapshot) y delivery/configuration (validación de config al boot).
