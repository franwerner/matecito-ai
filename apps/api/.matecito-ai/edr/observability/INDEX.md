# Dominio: `observability` — Decisiones

Cómo se observa el broker: logging estructurado y health checks.

## EDRs en este dominio

| EDR | Status | Consultá cuando... |
|---|---|---|
| [logging.md](logging.md) | Accepted | agregás logs, definís niveles o adjuntás contexto de correlación. |
| [health-checks.md](health-checks.md) | Accepted | tocás el endpoint de liveness o su semántica. |

## No aplican en este dominio

(Fases recomendadas para este tipo de proyecto que se descartaron. No generan EDR-archivo; su razón queda acá.)

| Concern | Razón |
|---|---|
| metrics | Daemon local sin agregador de métricas. |
| tracing | Daemon local, sin tracing distribuido. |

**Leyenda de status:** `Accepted` · `Pending` · `Not Applicable` · `Deferred`.
