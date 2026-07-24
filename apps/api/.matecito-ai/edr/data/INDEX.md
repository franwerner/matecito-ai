# Dominio: `data` — Decisiones

Cómo se persiste y accede al estado del broker: el borde de acceso a datos y el modelado del event-log.

## EDRs en este dominio

| EDR | Status | Type | Consultá cuando... |
|---|---|---|---|
| [data-access.md](data-access.md) | Accepted | decision | tocás el store, agregás una query o una migración. |
| [data-modeling.md](data-modeling.md) | Pending | decision | vas a definir el schema concreto de tablas del store. |

## No aplican en este dominio

(Fases recomendadas para este tipo de proyecto que se descartaron. No generan EDR-archivo; su razón queda acá.)

| Concern | Razón |
|---|---|

**Leyenda de status:** `Accepted` · `Pending` · `Not Applicable` · `Deferred` · `Superseded`.
