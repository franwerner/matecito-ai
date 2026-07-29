# Dominio: `delivery` — Decisiones

Configuración, empaquetado/despliegue, gates de calidad, documentación y feature flags de la UI.

## EDRs en este dominio

| EDR | Status | Consultá cuando... |
|---|---|---|
| [configuration.md](configuration.md) | Accepted | agregás una variable de entorno, tocás la validación de config al boot o el `.env`. |
| [deployment-topology.md](deployment-topology.md) | Accepted | tocás el build estático, el embebido en el binario del broker o el orden de build en la release. |
| [ci-quality-gates.md](ci-quality-gates.md) | Accepted | tocás los gates de merge, el pre-commit o el sync de tipos generados. |
| [documentation.md](documentation.md) | Accepted | escribís el README o te preguntás dónde está la doc de API. |
| [feature-flags.md](feature-flags.md) | Pending | aparece una necesidad de flag de larga duración. |

## No aplican en este dominio

(Fases recomendadas para este tipo de proyecto que se descartaron. No generan EDR-archivo; su razón queda acá.)

| Concern | Razón |
|---|---|
| testing-strategy | Decisión consciente de no formalizar una suite de tests en esta etapa (read-first, un dev). |

**Leyenda de status:** `Accepted` · `Pending` · `Not Applicable` · `Deferred`.
