# TanStack Router

- **Category:** Other
- **Version:** sin pinear
- **Status:** Accepted
- **Decided in phase:** frontend
- **Date:** 2026-07-23

## Por qué la elegimos

Router file-based con search params tipados para el estado remontable y `errorComponent` por ruta; los loaders prefetchean el snapshot vía Query. Paquetes npm: `@tanstack/react-router` + `@tanstack/router-plugin` (Vite).

## Alternativas descartadas

- Ninguna evaluada formalmente.

## Notas

Usada en: frontend/routing (rutas file-based, search params tipados, loaders) y runtime/error-handling (boundary por ruta).
