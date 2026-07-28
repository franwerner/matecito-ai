# Capability specs — `lifecycle`

Máquinas de estado de una entidad (su ciclo de vida y transiciones).

**Cuándo consultar este tipo:** antes de tocar los estados de una entidad o las tools que los transicionan — para respetar sus transiciones válidas y sus guards.

## Capacidades

| Capacidad | Qué hace | Status | Spec |
|---|---|---|---|
| `change` | Ciclo de vida del change (`active` ↔ `closed`) y el guard `change_closed` | Accepted | [`change.md`](change.md) |
| `record-version` | Ciclo de vida de la versión de un record (`actual` → `congelada`): el pin congela, update-in-place vs fork, inmutabilidad | Accepted | [`record-version.md`](record-version.md) |
