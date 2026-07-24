# Vite

- **Category:** Other
- **Version:** sin pinear
- **Status:** Accepted
- **Decided in phase:** delivery
- **Date:** 2026-07-23

## Por qué la elegimos

Build tool y dev server de la UI: produce el build estático que se embebe en el binario del broker y provee el dev server con proxy al broker (same-origin) y la superficie `import.meta.env` para la config.

## Alternativas descartadas

- Ninguna evaluada formalmente.

## Notas

Usada en: delivery/deployment-topology (build estático a `dist`), delivery/configuration (`import.meta.env`, proxy de dev) y frontend/routing (plugin de router file-based).
