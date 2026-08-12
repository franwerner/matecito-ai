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
| [consolidation-run-is-the-integrator.md](consolidation-run-is-the-integrator.md) | Accepted | Vas a decidir quién ejecuta la integración de un batch paralelo. |
| [phase-fanout-two-cases.md](phase-fanout-two-cases.md) | Accepted | Vas a tocar la prosa de fan-out del pipeline o a agregar un tercer caso de despacho concurrente. |
| [prose-register-single-home.md](prose-register-single-home.md) | Accepted | Vas a agregar o editar una instrucción de registro de prosa gate-facing y te tienta copiarla en más de un contrato de retorno. |
| [verdict-classified-by-the-orchestrator.md](verdict-classified-by-the-orchestrator.md) | Accepted | Vas a tocar quién clasifica un token declarado por una fase, o a agregar un tercer token del mecanismo de captura in-flow. |

## No aplican en este dominio

(Fases recomendadas para este tipo de proyecto que se descartaron. No generan EDR-archivo; su razón queda acá.)

| Concern | Razón |
|---|---|
| — | Ninguna descartada todavía en este dominio. |

**Leyenda de status:** `Accepted` · `Pending` · `Not Applicable` · `Deferred`.
