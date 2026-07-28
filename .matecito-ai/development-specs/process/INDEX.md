# Capability specs — `process`

Comportamiento reactivo/de fondo, disparado por un evento o por el sistema (no por un actor): file-watchers, jobs, reconciliación, indexado.

**Cuándo consultar este tipo:** antes de tocar cualquier proceso de fondo del daemon —el file-watch e indexado de records, su versionado o su proyección— leé el spec para conocer su disparador, su flujo, sus casos borde y sus reglas de versionado.

## Capacidades

| Capacidad | Qué hace | Status | Spec |
|---|---|---|---|
| `index-decision-records` | Indexa y versiona (copy-on-write lazy) los `.md` de EDR/spec del proyecto ante file-watch, sosteniendo el pin de la versión que aplicó cada evento | Accepted | [`index-decision-records.md`](index-decision-records.md) |
| `ingest-mechanical-events` | Ingesta automática de arranque/fin de sub-agentes (hooks del runtime) al event-log del change activo | Accepted | [`ingest-mechanical-events.md`](ingest-mechanical-events.md) |
| `capture-code-snapshots` | Fotos del contenido de cada archivo tocado por el apply, en los bordes de batch (antes sincrónico pre-edición, después al submit), por change | Accepted | [`capture-code-snapshots.md`](capture-code-snapshots.md) |
