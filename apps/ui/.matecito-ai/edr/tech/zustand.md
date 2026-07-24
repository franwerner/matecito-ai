# Zustand

- **Category:** Other
- **Version:** sin pinear
- **Status:** Accepted
- **Decided in phase:** structure
- **Date:** 2026-07-23

## Por qué la elegimos

Store del estado efímero de UI (viewport del canvas, nodo seleccionado, posición del scrubber, tema, command palette, lente/filtro activo), separado del server-state que vive en TanStack Query.

## Alternativas descartadas

- Ninguna evaluada formalmente.

## Notas

Usada en: structure/architecture-style (split server-state vs UI-state). El WS empuja al store vía el borde de datos.
