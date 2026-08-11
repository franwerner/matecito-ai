# Dominio: `contracts` — Decisiones

Decisiones sobre los contratos entre piezas del ecosistema: qué forma tiene lo que entra y sale de una herramienta, y cómo se comunica un resultado.

## EDRs en este dominio

| EDR | Status | Consultá cuando... |
|---|---|---|
| [three-level-checker-severity.md](three-level-checker-severity.md) | Accepted | Vas a agregar un chequeo a un validador, decidir si algo rompe el build, o interpretar la salida de un validador sobre material heredado. |
| [data-contract-derived-and-producer-neutral.md](data-contract-derived-and-producer-neutral.md) | Accepted | Vas a agregar un productor de artefactos durables, extender el esquema de entrada de un renderizador, o te tienta aceptar un campo que el renderizador ya sabe calcular. |
| [one-commit-per-isolated-run.md](one-commit-per-isolated-run.md) | Accepted | Vas a tocar cómo una corrida aislada de implementación persiste su trabajo. |
| [worktree-base-handshake.md](worktree-base-handshake.md) | Accepted | Vas a tocar cómo se verifica la base de un espacio de trabajo aislado antes de implementar o integrar. |
| [single-writer-per-batch.md](single-writer-per-batch.md) | Accepted | Vas a tocar quién escribe el artefacto de progreso o el estado de las tareas de un cambio. |

## No aplican en este dominio

(Fases recomendadas para este tipo de proyecto que se descartaron. No generan EDR-archivo; su razón queda acá.)

| Concern | Razón |
|---|---|
| — | Ninguno descartado todavía. |

**Leyenda de status:** `Accepted` · `Pending` · `Not Applicable` · `Deferred`.
