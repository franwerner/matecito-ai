# Capability specs — `rule`

Reglas de negocio transversales, sin flujo (scoping, políticas, invariantes).

**Cuándo consultar este tipo:** antes de tocar cómo un evento se asocia a un (proyecto, change), cualquier lógica que dependa de la identidad por request (proyecto, change, sesión), o la resiliencia de cualquier productor que ingesta contenido al broker.

## Capacidades

| Capacidad | Qué hace | Status | Spec |
|---|---|---|---|
| `event-scoping` | Asocia cada evento (de una tool MCP o de un hook) a un (proyecto, change) | Accepted | [`event-scoping.md`](event-scoping.md) |
| `ingestion-spool` | Resiliencia transversal de la ingesta: fire-and-forget, spool local completo ante falla y reconciliación idempotente en orden | Accepted | [`ingestion-spool.md`](ingestion-spool.md) |
