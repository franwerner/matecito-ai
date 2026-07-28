# Dominio: `delivery` — Decisiones

Cómo se configura, testea, empaqueta y despliega el broker.

## EDRs en este dominio

| EDR | Status | Type | Consultá cuando... |
|---|---|---|---|
| [configuration.md](configuration.md) | Accepted | convention | agregás un parámetro de config, un flag o una env var. |
| [testing-strategy.md](testing-strategy.md) | Accepted | decision | escribís tests o decidís qué nivel/recurso usar. |
| [deployment-topology.md](deployment-topology.md) | Accepted | decision | tocás el empaquetado, el modelo de instancia o el scope por proyecto. |
| [ci-quality-gates.md](ci-quality-gates.md) | Accepted | policy | armás el CI, un gate de PR o el orden de build de la release. |

## No aplican en este dominio

(Fases recomendadas para este tipo de proyecto que se descartaron. No generan EDR-archivo; su razón queda acá.)

| Concern | Razón |
|---|---|

**Leyenda de status:** `Accepted` · `Pending` · `Not Applicable` · `Deferred` · `Superseded`.
