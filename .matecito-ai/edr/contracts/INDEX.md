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
| [uncommitted-gate-follows-the-container.md](uncommitted-gate-follows-the-container.md) | Accepted | Vas a tocar qué árbol inspecciona un chequeo cuando el aislamiento anidado puede desplazar cuál es el contenedor relevante. |
| [workspace-forwarded-in-dispatch.md](workspace-forwarded-in-dispatch.md) | Accepted | Vas a tocar cómo una fase se entera de dónde está el espacio de trabajo en el que tiene que operar. |
| [gate-summary-cap-in-characters.md](gate-summary-cap-in-characters.md) | Accepted | Vas a tocar el límite de longitud de un resumen de ítem en un gate, o en qué momento del ciclo de vida del retorno se aplica ese chequeo. |
| [anchor-token-free-form.md](anchor-token-free-form.md) | Accepted | Vas a validar (o te tienta validar) la forma del token de ancla de un ítem de gate. |
| [ratification-ledger-row.md](ratification-ledger-row.md) | Accepted | Vas a tocar los campos que lleva una fila del registro de ratificación de un cambio. |
| [re-emergence-match-on-record.md](re-emergence-match-on-record.md) | Accepted | Vas a tocar cómo se detecta que un ítem de gate ya fue ratificado antes en el mismo cambio. |
| [re-emergence-short-form.md](re-emergence-short-form.md) | Accepted | Vas a tocar cómo se presenta un ítem de gate que ya reapareció y fue ratificado antes. |
| [subverifier-item-shape-single-declaration.md](subverifier-item-shape-single-declaration.md) | Accepted | Vas a tocar la forma de un hallazgo de sub-verificador, o el prompt de despacho de uno de los siete. |
| [contract-proposal-has-no-persistence-slot.md](contract-proposal-has-no-persistence-slot.md) | Accepted | Te tienta agregar a una propuesta de contrato un campo que diga dónde va a persistir. |
| [per-field-description-cap.md](per-field-description-cap.md) | Accepted | Vas a tocar el límite de longitud de la descripción de un campo dentro de un contrato propuesto. |
| [contract-item-anchors-once.md](contract-item-anchors-once.md) | Accepted | Te tienta agregar un ancla por campo a un contrato propuesto con varios campos. |
| [nested-field-continuation-line.md](nested-field-continuation-line.md) | Accepted | Vas a tocar cómo se imprime un campo dentro de un ítem de contrato propuesto. |
| [narrowed-boundary-scope-test.md](narrowed-boundary-scope-test.md) | Accepted | Vas a tocar el test que decide si una edición cruza hacia territorio de contrato no inferible. |
| [ratified-contract-travels-in-the-dispatch-prompt.md](ratified-contract-travels-in-the-dispatch-prompt.md) | Accepted | Vas a tocar cómo una fase recibe la forma ratificada de un contrato que ella misma propuso. |
| [side-discussion-handoff-shape.md](side-discussion-handoff-shape.md) | Accepted | Vas a tocar la forma del traspaso de una discusión lateral o a agregarle un campo nuevo. |
| [side-discussion-conclusion-shape.md](side-discussion-conclusion-shape.md) | Accepted | Vas a tocar la forma de la conclusión de una discusión lateral o cómo entra al camino de captura in-flow. |
| [side-discussion-conclusion-is-not-a-record.md](side-discussion-conclusion-is-not-a-record.md) | Accepted | Te tienta que una sesión lateral escriba un EDR directamente, o copiar una conclusión a un EDR sin gate. |
| [side-discussion-launch-boundary.md](side-discussion-launch-boundary.md) | Accepted | Vas a tocar el comando de lanzamiento de la sesión lateral o el límite de qué puede hacer. |
| [side-discussion-pickup-on-consult.md](side-discussion-pickup-on-consult.md) | Accepted | Vas a tocar cómo el principal decide si esperar o seguir trabajando durante una discusión lateral, o cómo recoge su conclusión. |
| [side-discussion-launcher-test.md](side-discussion-launcher-test.md) | Accepted | Vas a tocar cómo se lanza la sesión lateral, qué cuenta como un lanzamiento válido, o el mensaje cuando no hay ninguno disponible. |

## No aplican en este dominio

(Fases recomendadas para este tipo de proyecto que se descartaron. No generan EDR-archivo; su razón queda acá.)

| Concern | Razón |
|---|---|
| — | Ninguno descartado todavía. |

**Leyenda de status:** `Accepted` · `Pending` · `Not Applicable` · `Deferred`.
