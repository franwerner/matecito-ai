# TanStack Query

- **Category:** Other
- **Version:** sin pinear
- **Status:** Accepted
- **Decided in phase:** frontend
- **Date:** 2026-07-23

## Por qué la elegimos

Cache del server/remote state del broker (snapshot HTTP + eventos WS que la alimentan/invalidan), con retry+backoff y estados de error de fetch de fábrica. Paquete npm: `@tanstack/react-query`.

## Alternativas descartadas

- Ninguna evaluada formalmente.

## Notas

Usada en: structure/architecture-style, frontend/data-fetching, frontend/routing (loaders con `ensureQueryData`), runtime/error-handling (error-state), runtime/resilience (retry/refetch) y security/input-validation (los schemas parsean antes de entrar a la cache).
