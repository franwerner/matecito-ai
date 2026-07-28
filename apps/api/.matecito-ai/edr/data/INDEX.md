# Dominio: `data` — Decisiones

Cómo se persiste y accede al estado del broker: el borde de acceso a datos y el modelado del event-log.

## EDRs en este dominio

| EDR | Status | Type | Consultá cuando... |
|---|---|---|---|
| [data-access-entity-framework.md](data-access-entity-framework.md) | Accepted | decision | tocás el store: definís schema, agregás una lectura/escritura o generás una migración. |
| [data-modeling.md](data-modeling.md) | Accepted | decision | tocás el schema del store: IDs, borrado, timestamps, idempotencia, envelope de eventos o scoping por proyecto. |
| [storage-sync-model.md](storage-sync-model.md) | Accepted | decision | tocás la sincronización archivo↔base, el versionado content-addressable, o el deploy/compartir. |

## No aplican en este dominio

(Fases recomendadas para este tipo de proyecto que se descartaron. No generan EDR-archivo; su razón queda acá.)

| Concern | Razón |
|---|---|

**Leyenda de status:** `Accepted` · `Pending` · `Not Applicable` · `Deferred` · `Superseded`.
