# Dominio: `structure` — Decisiones

Decisiones sobre cómo se organiza el payload del repo (dónde vive cada concepto, cómo se citan entre sí).

## EDRs en este dominio

| EDR | Status | Consultá cuando... |
|---|---|---|
| [repo-components-home.md](repo-components-home.md) | Accepted | Vas a documentar o citar el concepto repo-level de `components` desde cualquier consumidor. |
| [components-concept-bridge.md](components-concept-bridge.md) | Accepted | Vas a mudar la mitad de un concepto documentado a un hogar propio y necesitás decidir qué queda en el lugar de origen. |
| [claude-md-components-split.md](claude-md-components-split.md) | Accepted | Vas a resumir en la guía del dominio un concepto que tiene una declaración repo-level y una o más proyecciones. |
| [two-scripts-render-and-validate.md](two-scripts-render-and-validate.md) | Accepted | Vas a agregar una herramienta que produzca o audite artefactos durables, o te tienta unificar producción y auditoría en un solo ejecutable. |
| [contract-pair-in-templates.md](contract-pair-in-templates.md) | Accepted | Vas a crear el contrato máquina de una familia de artefactos, o a decidir dónde vive un esquema. |
| [root-index-cardinality-per-domain-type.md](root-index-cardinality-per-domain-type.md) | Accepted | Vas a materializar una tanda de artefactos durables y tenés que actualizar los índices del store. |
| [dispatch-batch-bound-integration.md](dispatch-batch-bound-integration.md) | Accepted | Vas a tocar cómo se despacha o integra un batch de implementación paralelo. |
| [consolidation-run-is-the-integrator.md](consolidation-run-is-the-integrator.md) | Accepted | Vas a decidir quién ejecuta una integración, en cualquiera de sus dos niveles. |
| [phase-fanout-two-cases.md](phase-fanout-two-cases.md) | Accepted | Vas a tocar la prosa de fan-out del pipeline o a agregar un tercer caso de despacho concurrente. |
| [prose-register-single-home.md](prose-register-single-home.md) | Accepted | Vas a agregar o editar una instrucción de registro de prosa gate-facing y te tienta copiarla en más de un contrato de retorno. |
| [pr-base-branch-rule-reach.md](pr-base-branch-rule-reach.md) | Accepted | Vas a tocar la descripción de una estrategia de cadena de PRs (stacked-to-main, feature-branch-chain) o a agregar una nueva. |
| [pr-base-explicit-argument-form.md](pr-base-explicit-argument-form.md) | Accepted | Vas a documentar o revisar cómo se abre un PR (por el flujo o a mano) y te tienta enunciarlo como norma en vez de argumento obligatorio. |
| [item-shaping-helper-seam.md](item-shaping-helper-seam.md) | Accepted | Vas a agregar una sección de retorno con ítems que renderiza como tabla, o a tocar cómo se imprime el adorno de un ítem. |
| [footer-survives-the-sentinel.md](footer-survives-the-sentinel.md) | Accepted | Vas a agregar un footer de tabla que registra evidencia de un chequeo distinto al de las filas, y necesitás que sobreviva al retorno `None.` de una tabla vacía. |
| [render-constraints-restated-for-subverifiers.md](render-constraints-restated-for-subverifiers.md) | Accepted | Vas a documentar una restricción que `render-return.js` impone sobre un campo que un sub-verificador de sdd-verify construye a mano. |
| [script-change-ships-with-dev-tests.md](script-change-ships-with-dev-tests.md) | Accepted | Vas a editar un script de `payload/domains/development/scripts/` que no tiene cobertura previa para el comportamiento que estás por cambiar. |
| [validate-side-needs-no-mirror.md](validate-side-needs-no-mirror.md) | Accepted | Te tienta agregar al validador de retornos un helper equivalente al del lado de renderizado. |
| [conditional-section-with-status-filter.md](conditional-section-with-status-filter.md) | Accepted | Vas a agregar una sección condicional de retorno que además tenga que filtrar por el status de la fase. |
| [side-discussion-prose-homes.md](side-discussion-prose-homes.md) | Accepted | Vas a decidir dónde vive la prosa de un mecanismo del orquestador que tiene un lector fuera de cualquier dominio. |
| [side-discussion-topic-key-namespace.md](side-discussion-topic-key-namespace.md) | Accepted | Vas a agregar una clave de Engram fuera del namespace `sdd/` para un intercambio entre dos actores. |
| [gating-vocabulary-sweep-scope.md](gating-vocabulary-sweep-scope.md) | Accepted | Vas a barrer el vocabulario de gating retirado (Tier 1/Tier 2) fuera del Scope original ratificado de un cambio. |
| [retired-vocabulary-in-record-stores.md](retired-vocabulary-in-record-stores.md) | Accepted | Vas a decidir si el sweep final de un cambio de vocabulario debe extender su scope a un record store pre-existente cuyo sustantivo queda desactualizado. |
| [lane-vocabulary-sweep-in-record-stores.md](lane-vocabulary-sweep-in-record-stores.md) | Accepted | Vas a barrer vocabulario de lane retirado y necesitás decidir si el sweep alcanza una fila de navegación de un índice de store durable, o un archivo cuyo sustantivo incidental sobrevive como residuo aceptado. |
| [sdd-propose-contract-joins-lane-sweep.md](sdd-propose-contract-joins-lane-sweep.md) | Accepted | Vas a barrer vocabulario de lane retirado y encontrás un archivo fuera del Scope original que describe un mecanismo (add-on opcional, contraste de modos) que el cambio retira. |
| [serial-apply-dispatch-budget.md](serial-apply-dispatch-budget.md) | Accepted | Vas a tocar el guard que corta un despacho serial de sdd-apply por costo estimado, o a decidir dónde vive un mecanismo de este tipo. |
| [edr-mark-is-the-implementing-task.md](edr-mark-is-the-implementing-task.md) | Accepted | Vas a escribir un checklist de tareas (`sdd-tasks`) que incluye una decisión de arquitectura, o a decidir si la materialización de un EDR va en su propia tarea. |
| [falsified-clause-joins-the-sweep.md](falsified-clause-joins-the-sweep.md) | Accepted | Vas a barrer el vocabulario que un cambio retira y encontrás una cláusula falsificada fuera del Scope original del spec. |
| [release-resolution-helper-home.md](release-resolution-helper-home.md) | Accepted | Vas a resolver un release de GitHub para instalar algo, y te tienta agregarle un modo a releasedl. |
| [release-asset-by-stable-name.md](release-asset-by-stable-name.md) | Accepted | Vas a resolver qué asset de una release de GitHub instalar por nombre. |
| [mcp-step-guards-both-artifacts.md](mcp-step-guards-both-artifacts.md) | Accepted | Vas a tocar el paso de instalación del MCP de qmd, o a decidir cómo un paso de MCP con binario propio detecta trabajo pendiente. |
| [installer-step-ships-with-tests.md](installer-step-ships-with-tests.md) | Accepted | Vas a escribir o revisar un paso de instalación nuevo y te tienta omitir sus tests porque el Scope del spec no los menciona. |
| [mcp-connectivity-read-from-host.md](mcp-connectivity-read-from-host.md) | Accepted | Vas a leer `.Connection` sobre un `Found` de `mcp.Find` —o a tratar sus tres estados, incluido el "no preguntado"—, o a escribir la lógica de `Describe()` para una integración registrada por MCP. |
| [qmd-run-fails-on-shadowing-executable.md](qmd-run-fails-on-shadowing-executable.md) | Accepted | Vas a tocar `qmdMCPStep.Run`, o a decidir qué hacer con un ejecutable que un paso de instalación no instaló pero sigue ganando en PATH. |
| [binary-presence-probe-shared-by-both-surfaces.md](binary-presence-probe-shared-by-both-surfaces.md) | Accepted | Vas a agregar un segundo caller de "probar si un binario está instalado en su ubicación canónica" a check.ProbeAt, o a decidir si un check de presencia va en internal/check o en el paquete que lo consume. |
| [binary-detection-at-canonical-path.md](binary-detection-at-canonical-path.md) | Accepted | Vas a tocar cómo sync.Detect resuelve la presencia de engram, codegraph o proofshot, o a decidir si conviene reemplazar los tres bloques por una tabla. |
| [step-install-location-accessors.md](step-install-location-accessors.md) | Accepted | Vas a necesitar la ubicación canónica donde un paso de instalación deja su binario, desde fuera de internal/setup/install. |
| [npm-bin-dir-single-resolver.md](npm-bin-dir-single-resolver.md) | Accepted | Vas a resolver dónde npm deja los ejecutables globales, en vez de volver a derivarlo con tu propio `npm config get prefix`. |
| [check-resolver-seam.md](check-resolver-seam.md) | Accepted | Vas a agregar o tocar un check de verify que prueba la presencia de un binario que otro paso instala. |
| [spec-fold-fires-on-the-final-dispatch.md](spec-fold-fires-on-the-final-dispatch.md) | Accepted | Vas a tocar cuándo o con qué predicado se dispara el pliegue del delta de comportamiento en sdd-apply. |
| [spec-merge-relocation-sweep-scope.md](spec-merge-relocation-sweep-scope.md) | Accepted | Vas a reubicar un mecanismo del flujo y necesitás decidir si el barrido de prosa desactualizada alcanza un archivo del repo que nunca se despliega. |
| [single-root-record-store-after-cockpit-removal.md](single-root-record-store-after-cockpit-removal.md) | Accepted | Vas a documentar o citar dónde viven las decisiones o el comportamiento del sistema (EDRs o capability-specs), o a decidir si un componente nuevo necesita su propio store. |
| [closed-front-recorded-not-deleted.md](closed-front-recorded-not-deleted.md) | Accepted | Vas a cerrar (no borrar) una sección de un documento de trabajo cuya contraparte fue retirada, y necesitás decidir si el registro queda o se va con ella. |
| [spec-deletion-over-deprecation.md](spec-deletion-over-deprecation.md) | Accepted | Vas a retirar una capability-spec sin reemplazo y necesitás decidir si se borra o se marca Deprecated. |
| [sweep-reach-outside-the-ratified-scope.md](sweep-reach-outside-the-ratified-scope.md) | Accepted | Vas a barrer vocabulario retirado fuera del Scope original de un cambio y encontrás un sitio adicional que el cambio falsifica pero el Scope no nombra. |
| [single-conditional-skill-reached-by-read.md](single-conditional-skill-reached-by-read.md) | Accepted | Antes de dar a un agente de fase el mecanismo para alcanzar exactamente una skill condicional. |
| [moment-count-counts-gates-not-bullets.md](moment-count-counts-gates-not-bullets.md) | Accepted | Vas a agregar o quitar un momento que cita el walkthrough compartido de gate-presentation.md, o a tocar cómo se cuenta. |
| [retire-by-deletion-when-the-premise-is-gone.md](retire-by-deletion-when-the-premise-is-gone.md) | Accepted | Vas a retirar un record sin reemplazo y necesitás decidir si se borra o se marca Deprecated — en cualquier cambio, no solo en el que originó el precedente. |
| [spec-store-scope-excludes-domain-agent-behavior.md](spec-store-scope-excludes-domain-agent-behavior.md) | Accepted | Vas a proponer una capability-spec en .matecito-ai/development-specs/ cuyo actor es un dominio distinto de la superficie CLI de development. |
| [mirror-sites-join-the-prune-sweep.md](mirror-sites-join-the-prune-sweep.md) | Accepted | Antes de tocar un archivo que instruye/enlaza algo que un cambio retira |
| [apply-progress-continuity-binds-the-writer-role.md](apply-progress-continuity-binds-the-writer-role.md) | Accepted | Vas a tocar la continuidad de apply-progress cuando el fragmento de dominio declara más de un rol de despacho para la misma fase. |
| [brief-flags-declared-where-orchestrator-reads.md](brief-flags-declared-where-orchestrator-reads.md) | Accepted | Vas a tocar dónde se declara la tabla de flags de decisión del brief, o quién tiene que reportarlos. |
| [decision-record-statuses-checked-on-return.md](decision-record-statuses-checked-on-return.md) | Accepted | Vas a tocar dónde se chequean los status `blocked`/`needs-decision` que dependen de los decision records. |
| [discovery-runs-after-code-is-read.md](discovery-runs-after-code-is-read.md) | Accepted | Vas a tocar en qué fase corre el Discovery Gate, o cuándo se formula el formulario de discovery. |
| [domain-fragment-trigger-is-an-act.md](domain-fragment-trigger-is-an-act.md) | Accepted | Vas a tocar cuándo se carga el fragmento de dominio, o a decidir si un lane queda exento de cargarlo. |
| [guard-references-canonical-mailbox-list.md](guard-references-canonical-mailbox-list.md) | Accepted | Vas a tocar la lista de secciones gateables del Unresolved Decisions Guard, o te tienta agregarle una copia propia. |
| [guard-turns-mailboxes-into-a-trigger.md](guard-turns-mailboxes-into-a-trigger.md) | Accepted | Vas a tocar el propósito o el disparo del Unresolved Decisions Guard. |
| [init-guard-key-is-declared-not-derived.md](init-guard-key-is-declared-not-derived.md) | Accepted | Vas a tocar cómo el Init Guard construye la clave de búsqueda, o a agregar un dominio nuevo al flujo. |
| [kernel-does-not-enumerate-domain-shared-files.md](kernel-does-not-enumerate-domain-shared-files.md) | Accepted | Vas a agregar o sacar un archivo compartido de un dominio y te tienta nombrarlo en el kernel. |
| [kernel-fanout-pointer-not-a-pattern.md](kernel-fanout-pointer-not-a-pattern.md) | Accepted | Vas a tocar la mención de sdd-verify como excepción de despacho, o te tienta ofrecer fan-out a otra fase. |
| [kernel-states-discovery-invariant-not-mechanism.md](kernel-states-discovery-invariant-not-mechanism.md) | Accepted | Vas a tocar el invariante de discovery del kernel, o a decidir si nombra un mecanismo concreto. |
| [mine-gate-override-by-presence.md](mine-gate-override-by-presence.md) | Accepted | Vas a tocar el mine gate del kernel, o a decidir si un dominio con mecanismo propio de captura queda exento. |
| [no-slash-commands-for-phases.md](no-slash-commands-for-phases.md) | Accepted | Vas a documentar cómo se invoca una fase del flujo. |
| [open-questions-is-not-a-decision-mailbox.md](open-questions-is-not-a-decision-mailbox.md) | Accepted | Vas a decidir si un item pendiente va a `Open Questions` o a `New Decisions`. |
| [parallel-batch-role-not-a-table-row.md](parallel-batch-role-not-a-table-row.md) | Accepted | Vas a tocar la tabla `SDD Phase Read/Write` para el caso de batch paralelo de sdd-apply. |
| [reasoning-record-granularity-per-decision.md](reasoning-record-granularity-per-decision.md) | Accepted | Vas a reubicar un comentario de razonamiento del payload y necesitás decidir la granularidad del record. |
| [short-guard-points-to-single-mechanism-file.md](short-guard-points-to-single-mechanism-file.md) | Accepted | Vas a tocar Parallel-Mark Validation o Uncommitted-Work Gate, o te tienta duplicar su mecanismo en el fragmento. |
| [spec-always-reads-intake-for-ui-test.md](spec-always-reads-intake-for-ui-test.md) | Accepted | Vas a tocar si sdd-spec lee el intake brief siempre o sólo como fallback. |
| [change-scoped-requirements-excluded-from-the-fold.md](change-scoped-requirements-excluded-from-the-fold.md) | Accepted | Vas a tocar cómo el paso que cierra un cambio decide qué requisitos del spec se pliegan en un capability-spec durable. |
| [merge-turn-granularity.md](merge-turn-granularity.md) | Accepted | Vas a tocar en qué punto del loop de consolidación se reclama o se libera el turno de rama compartida. |
| [merge-turn-wait-and-stop.md](merge-turn-wait-and-stop.md) | Accepted | Vas a tocar qué hace una corrida cuando encuentra el turno de rama compartida tomado. |
| [merge-queue-lives-in-refs.md](merge-queue-lives-in-refs.md) | Accepted | Vas a tocar dónde vive la lista de espera del turno de rama compartida. |
| [guard-queues-and-retries-agent-asks.md](guard-queues-and-retries-agent-asks.md) | Accepted | Vas a tocar la mecánica de cola y reintento de `turn claim`, o a decidir qué ve el que llama. |
| [turn-mechanism-home-and-wiring-register.md](turn-mechanism-home-and-wiring-register.md) | Accepted | Vas a tocar dónde vive el mecanismo del turno de rama compartida, o dónde se cita cada punto de wiring. |
| [manual-turn-procedure-retained-for-unguarded-sessions.md](manual-turn-procedure-retained-for-unguarded-sessions.md) | Accepted | Vas a tocar el procedimiento manual de claim/release para una sesión sin `matecito-ai` instalado. |
| [ported-record-adapted-to-the-target-store.md](ported-record-adapted-to-the-target-store.md) | Accepted | Vas a portar un record a otro store y su premisa (un write-moment, una relación citada) deja de existir en el destino. |
| [turn-port-sweep-reaches-five-more-sites.md](turn-port-sweep-reaches-five-more-sites.md) | Accepted | Vas a barrer vocabulario retirado y encontrás un sitio fuera del Scope original ratificado del cambio. |

## No aplican en este dominio

(Fases recomendadas para este tipo de proyecto que se descartaron. No generan EDR-archivo; su razón queda acá.)

| Concern | Razón |
|---|---|
| — | Ninguna descartada todavía en este dominio. |

**Leyenda de status:** `Accepted` · `Pending` · `Not Applicable` · `Deferred`.
