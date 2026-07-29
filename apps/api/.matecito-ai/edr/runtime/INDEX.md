# Dominio: `runtime` — Decisiones

Cómo se comporta el broker en ejecución: manejo de errores, concurrencia y resiliencia.

## EDRs en este dominio

| EDR | Status | Consultá cuando... |
|---|---|---|
| [error-handling.md](error-handling.md) | Accepted | propagás, traducís o loggeás un error. |
| [concurrency-async.md](concurrency-async.md) | Accepted | agregás una goroutine, un channel o coordinás el lifecycle. |
| [resilience.md](resilience.md) | Accepted | tocás presupuestos de hooks sincrónicos, reconexión del WS o la política del spool local. |

## No aplican en este dominio

(Fases recomendadas para este tipo de proyecto que se descartaron. No generan EDR-archivo; su razón queda acá.)

| Concern | Razón |
|---|---|
| caching | El estado derivado sale de queries sobre el event-log; no hay capa de cache. |

**Leyenda de status:** `Accepted` · `Pending` · `Not Applicable` · `Deferred`.
