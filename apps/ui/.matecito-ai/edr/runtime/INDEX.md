# Dominio: `runtime` — Decisiones

Manejo de errores y resiliencia de la UI en ejecución.

## EDRs en este dominio

| EDR | Status | Type | Consultá cuando... |
|---|---|---|---|
| [error-handling.md](error-handling.md) | Accepted | decision | tocás error boundaries, errores de data, el mapeo del error del broker o el estado de conexión. |
| [resilience.md](resilience.md) | Accepted | decision | tocás reconexión del WS, retry del snapshot o el resync tras la caída. |

## No aplican en este dominio

(Fases recomendadas para este tipo de proyecto que se descartaron. No generan EDR-archivo; su razón queda acá.)

| Concern | Razón |
|---|---|
| concurrency-async | El async lo dan Query y partysocket; sin política de concurrencia propia. |
| background-jobs | La UI no corre jobs en background. |
| caching | Sin capa de cache propia; la cache es la de TanStack Query. |

**Leyenda de status:** `Accepted` · `Pending` · `Not Applicable` · `Deferred`.
