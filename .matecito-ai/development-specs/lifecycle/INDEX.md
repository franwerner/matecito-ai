# Capability specs — `lifecycle`

Máquinas de estado de una entidad (su ciclo de vida y transiciones).

**Cuándo consultar este tipo:** antes de tocar los estados de una entidad o las tools que los transicionan — para respetar sus transiciones válidas y sus guards.

## Capacidades

| Capacidad | Qué hace | Status | Spec |
|---|---|---|---|
| `component-check-status` | Los tres estados de un resultado de chequeo (`OK` / `Missing` / `Outdated`), sus reglas de derivación y la remediación obligatoria del `Missing` | Inferred | [`component-check-status.md`](component-check-status.md) |
