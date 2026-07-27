# Capability specs — `rule`

Reglas de negocio transversales, sin flujo (scoping, políticas, invariantes).

**Cuándo consultar este tipo:** antes de tocar cómo un evento se asocia a un (proyecto, change), o cualquier lógica que dependa de la identidad por request (proyecto, change, sesión).

## Capacidades

| Capacidad | Qué hace | Status | Spec |
|---|---|---|---|
| `event-scoping` | Asocia cada evento (de una tool MCP o de un hook) a un (proyecto, change) | Accepted | [`event-scoping.md`](event-scoping.md) |
