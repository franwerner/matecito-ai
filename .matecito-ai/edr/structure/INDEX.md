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

## No aplican en este dominio

(Fases recomendadas para este tipo de proyecto que se descartaron. No generan EDR-archivo; su razón queda acá.)

| Concern | Razón |
|---|---|
| — | Ninguna descartada todavía en este dominio. |

**Leyenda de status:** `Accepted` · `Pending` · `Not Applicable` · `Deferred`.
