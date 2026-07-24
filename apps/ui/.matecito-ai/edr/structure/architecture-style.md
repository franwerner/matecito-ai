# EDR — Estilo de arquitectura

- **Status:** Accepted
- **Type:** decision
- **Date:** 2026-07-23

## Contexto

El cockpit es una SPA React que consume el broker por WebSocket (push en vivo) y HTTP (snapshot). Es greenfield, single-dev y local-first, canvas-first (ReactFlow) con inspector, timeline y command palette. Hace falta fijar cómo se organiza el código y cómo se separa el estado remoto del efímero de UI.

## Decisión

Organización **feature-based (vertical slice)** con acoplamiento **pragmático**: la abstracción se concentra solo en el borde de datos hacia el broker (WS/HTTP), no en cada capa.

El **estado se parte** por naturaleza:
- **Server/remote state** en TanStack Query — los eventos del WS alimentan e invalidan la cache.
- **UI/efímero** en Zustand — viewport del canvas, nodo seleccionado, posición del scrubber, tema, command palette y lente/filtro activo.
- El **canvas** deriva su render (nodos/aristas) de la proyección del event-log; no mantiene una copia de verdad propia.

El WebSocket es el puente que empuja tanto a la cache de Query como al store de UI.

## Alcance

- `apps/ui/src/features/**` — módulos organizados por feature (vertical slice).

## Reglas verificables

- **[manual]** el server-state vive en TanStack Query, no duplicado en el store de UI.
- **[manual]** el estado efímero de UI vive en Zustand.
- **[manual]** no hay fetch imperativo desperdigado: todo acceso al broker pasa por el borde de datos.
- **[manual]** el canvas deriva nodos y aristas de la proyección del event-log, no de un estado paralelo.

## Relacionados

- `relacionado-con` → [folder-structure.md](folder-structure.md) — el layout de carpetas es consecuencia de esta organización por slice.
- `relacionado-con` → [../frontend/data-fetching.md](../frontend/data-fetching.md) — el borde de datos y el bridge WS → Query/Zustand.
