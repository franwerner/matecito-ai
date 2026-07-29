# EDR — Acceso a datos y fetching

- **Status:** Accepted
- **Date:** 2026-07-23

## Contexto

La UI consume el broker por dos canales: HTTP para el snapshot inicial y WebSocket para el push en vivo. El estado del cockpit se deriva de un event-log remoto; hace falta un único borde de datos que centralice ambos canales y alimente la cache y el estado de UI.

## Decisión

**TanStack Query** para el server-state y el snapshot (HTTP al broker) más **`partysocket`** como cliente WebSocket único, ubicado en el borde de datos compartido. El cliente WS provee auto-reconnect con backoff, buffer mientras está caído y heartbeat, y **bridgea** los eventos con envelope `{type, payload}` a la cache de Query (`setQueryData` / invalidate) y al store Zustand. La proyección del event-log deriva de ahí los nodos y aristas del canvas. Sin polling y sin fetch a mano desperdigado.

## Alcance

- `apps/ui/src/shared/api/**` — el cliente único del broker (WS + setup de Query).

## Reglas verificables

- **[manual]** todo acceso al broker pasa por el cliente único del borde de datos.
- **[manual]** el WS alimenta Query y Zustand; no hay fetch imperativo suelto.
- **[manual]** sin polling.

## Relacionados

- `relacionado-con` → [../structure/architecture-style.md](../structure/architecture-style.md) — el split server-state / UI-state y el rol del WS como puente.
- `relacionado-con` → [../runtime/resilience.md](../runtime/resilience.md) — reconexión, backoff y resync del snapshot.
- `relacionado-con` → [../security/input-validation.md](../security/input-validation.md) — el parseo del envelope y el snapshot en el borde.
