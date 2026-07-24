# EDR — Resiliencia

- **Status:** Accepted
- **Type:** decision
- **Date:** 2026-07-23

## Contexto

La UI vive de un único broker local con dos canales (WS y HTTP snapshot). El WS puede caerse; hace falta que la UI reconecte sin martillar al broker, muestre estado degradado y resincronice al volver, apoyándose en lo que ya proveen las libs.

## Decisión

Timeout y retry con backoff **provistos por las libs**, sin política custom:
- `partysocket` para el WS: reconnect con backoff, buffer y heartbeat;
- TanStack Query para el HTTP snapshot: retry con backoff y abort.

**Sin circuit breaker** (un solo broker local). **Resync del snapshot al reconectar** (invalidate/refetch vía Query) más un indicador stale/reconectando mientras el WS está caído. Los valores concretos son los defaults de las libs, ajustables en config. Complementa la resiliencia del lado del broker.

## Reglas verificables

- **[manual]** el WS reconecta con backoff (partysocket), sin reconnect fijo que martille al broker.
- **[manual]** con el WS caído la UI muestra estado stale/reconectando, no se cuelga.
- **[manual]** al reconectar, el snapshot se resincroniza.
- **[manual]** las llamadas del snapshot tienen retry con backoff, sin espera indefinida.

## Relacionados

- `relacionado-con` → [../frontend/data-fetching.md](../frontend/data-fetching.md) — el cliente único y el bridge WS → Query/Zustand.
