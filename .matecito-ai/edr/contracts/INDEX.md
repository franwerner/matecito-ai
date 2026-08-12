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
| [tone-prohibitions-kept-as-is.md](tone-prohibitions-kept-as-is.md) | Accepted | Te tienta aflojar o eliminar la prohibición de emojis o de frases motivacionales en cualquier zona de tono. |
| [acknowledgement-opener-test.md](acknowledgement-opener-test.md) | Accepted | Vas a escribir o revisar una línea de apertura de una respuesta y no estás seguro si aporta información. |
| [closing-line-only-for-a-pending-decision.md](closing-line-only-for-a-pending-decision.md) | Accepted | Te tienta cerrar una respuesta ofreciendo ayuda, o vas a decidir si una pregunta de cierre está permitida. |
| [plain-register-over-technical-neutral.md](plain-register-over-technical-neutral.md) | Accepted | Vas a describir el registro de una zona de tono y te tienta usar una etiqueta técnica interna en vez de una instrucción positiva. |
| [emit-in-english-present-translated.md](emit-in-english-present-translated.md) | Accepted | Vas a decidir si un aviso interno de una línea se traduce igual que el material gate-facing. |
| [rejected-proposal-verdict-token.md](rejected-proposal-verdict-token.md) | Accepted | Vas a agregar un token a un mailbox de retorno y te tienta usar un valor no-passing para un caso que el ejecutor necesita poder reportar honestamente. |
| [verdict-in-its-own-conditional-section.md](verdict-in-its-own-conditional-section.md) | Accepted | Vas a decidir dónde vive un veredicto que no encaja en ningún mailbox existente de sdd-apply. |
| [verdict-survives-the-fan-out.md](verdict-survives-the-fan-out.md) | Accepted | Vas a tocar cómo un veredicto declarado por una corrida aislada llega al retorno consolidado de un batch paralelo. |
| [expected-set-outside-the-validator.md](expected-set-outside-the-validator.md) | Accepted | Te tienta agregar un flag de completitud a un validador de forma de retorno. |

## No aplican en este dominio

(Fases recomendadas para este tipo de proyecto que se descartaron. No generan EDR-archivo; su razón queda acá.)

| Concern | Razón |
|---|---|
| — | Ninguno descartado todavía. |

**Leyenda de status:** `Accepted` · `Pending` · `Not Applicable` · `Deferred`.
