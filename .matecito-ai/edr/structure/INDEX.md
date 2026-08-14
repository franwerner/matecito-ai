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
| [change-level-worktree-isolation.md](change-level-worktree-isolation.md) | Accepted | Vas a tocar cómo se aísla un cambio completo del resto del repo, o cuándo se abre y se integra ese aislamiento. |
| [phase-fanout-two-cases.md](phase-fanout-two-cases.md) | Accepted | Vas a tocar la prosa de fan-out del pipeline o a agregar un tercer caso de despacho concurrente. |
| [prose-register-single-home.md](prose-register-single-home.md) | Accepted | Vas a agregar o editar una instrucción de registro de prosa gate-facing y te tienta copiarla en más de un contrato de retorno. |
| [verdict-classified-by-the-orchestrator.md](verdict-classified-by-the-orchestrator.md) | Accepted | Vas a tocar quién clasifica un token declarado por una fase, o a agregar un tercer token del mecanismo de captura in-flow. |
| [pr-base-branch-rule-reach.md](pr-base-branch-rule-reach.md) | Accepted | Vas a tocar la descripción de una estrategia de cadena de PRs (stacked-to-main, feature-branch-chain) o a agregar una nueva. |
| [pr-base-explicit-argument-form.md](pr-base-explicit-argument-form.md) | Accepted | Vas a documentar o revisar cómo se abre un PR (por el flujo o a mano) y te tienta enunciarlo como norma en vez de argumento obligatorio. |
| [change-isolation-activation-flag.md](change-isolation-activation-flag.md) | Accepted | Vas a tocar cómo se activa o se recomienda una elección junto con el lane en el fork o el INTAKE GATE. |
| [change-workspace-identity.md](change-workspace-identity.md) | Accepted | Vas a tocar la identidad (rama, directorio) del espacio de trabajo aislado de un cambio, o cómo se lo mantiene fuera de la vista de git. |
| [change-level-integration-act.md](change-level-integration-act.md) | Accepted | Vas a tocar cómo se cierra el ciclo de un cambio con aislamiento activo, o cómo se maneja un conflicto al integrarlo. |
| [change-workspace-cleanup.md](change-workspace-cleanup.md) | Accepted | Vas a tocar la limpieza del espacio de trabajo de un cambio, en cualquiera de sus dos resultados (integración limpia o fallida). |
| [change-workspace-prose-homes.md](change-workspace-prose-homes.md) | Accepted | Vas a agregar o mover prosa de un mecanismo que toca el kernel, un dominio y una referencia de fase a la vez, y no sabés dónde va cada parte. |
| [item-shaping-helper-seam.md](item-shaping-helper-seam.md) | Accepted | Vas a agregar una sección de retorno con ítems que renderiza como tabla, o a tocar cómo se imprime el adorno de un ítem. |
| [validate-side-needs-no-mirror.md](validate-side-needs-no-mirror.md) | Accepted | Te tienta agregar al validador de retornos un helper equivalente al del lado de renderizado. |
| [conditional-section-with-status-filter.md](conditional-section-with-status-filter.md) | Accepted | Vas a agregar una sección condicional de retorno que además tenga que filtrar por el status de la fase. |

## No aplican en este dominio

(Fases recomendadas para este tipo de proyecto que se descartaron. No generan EDR-archivo; su razón queda acá.)

| Concern | Razón |
|---|---|
| — | Ninguna descartada todavía en este dominio. |

**Leyenda de status:** `Accepted` · `Pending` · `Not Applicable` · `Deferred`.
