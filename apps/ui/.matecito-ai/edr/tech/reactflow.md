# ReactFlow

- **Category:** Other
- **Version:** sin pinear
- **Status:** Accepted
- **Decided in phase:** structure
- **Date:** 2026-07-23

## Por qué la elegimos

Motor del canvas del cockpit (canvas-first): renderiza nodos y aristas derivados de la proyección del event-log, con culling de viewport como mitigación de performance con muchos nodos.

## Alternativas descartadas

- Ninguna evaluada formalmente.

## Notas

Usada en: structure/architecture-style (el canvas deriva del event-log) y quality/nfr-performance (culling de viewport).
