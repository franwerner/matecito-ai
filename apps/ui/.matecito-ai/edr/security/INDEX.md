# Dominio: `security` — Decisiones

Validación de la entrada del broker y escaneo de dependencias de la UI.

## EDRs en este dominio

| EDR | Status | Consultá cuando... |
|---|---|---|
| [input-validation.md](input-validation.md) | Accepted | parseás lo entrante del broker (envelope WS / snapshot) o generás schemas desde el OpenAPI. |
| [dependency-scanning.md](dependency-scanning.md) | Accepted | tocás el escaneo de deps npm o el gate de vulnerabilidades. |

## No aplican en este dominio

(Fases recomendadas para este tipo de proyecto que se descartaron. No generan EDR-archivo; su razón queda acá.)

| Concern | Razón |
|---|---|
| auth | Local single-user; sin superficie de autenticación. |
| authorization | Local single-user; sin autorización. |
| rate-limiting | Cliente local; no aplica limitar tasa. |
| cors | La UI se sirve same-origin desde el broker; sin request cross-origin. |
| secrets-management | Sin secretos del lado cliente. |

**Leyenda de status:** `Accepted` · `Pending` · `Not Applicable` · `Deferred`.
