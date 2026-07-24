# Dominio: `contracts` — Decisiones

Las superficies de contrato del broker: qué recibe del MCP y qué emite a la UI.

## EDRs en este dominio

| EDR | Status | Type | Consultá cuando... |
|---|---|---|---|
| [api-contract.md](api-contract.md) | Accepted | decision | tocás HTTP-in, WS-out, el envelope, la idempotencia o el registro de proyecto. |

## No aplican en este dominio

(Fases recomendadas para este tipo de proyecto que se descartaron. No generan EDR-archivo; su razón queda acá.)

| Concern | Razón |
|---|---|
| cli | El broker no expone una CLI pública; su contrato son las dos superficies HTTP/WS. |
| library | El broker no es una librería consumida como API; se embebe como binario. |

**Leyenda de status:** `Accepted` · `Pending` · `Not Applicable` · `Deferred` · `Superseded`.
