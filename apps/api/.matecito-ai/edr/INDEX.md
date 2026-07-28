# Development Decision Records — Índice raíz

Las decisiones están organizadas por **dominio**. Este índice te dice qué dominio mirar; el detalle de cada decisión está en el índice de su dominio.

Contexto del proyecto: **broker** de Matecito UI (cockpit) — daemon backend local en Go bajo `apps/api`. Recibe JSON estructurado del MCP por HTTP, lo persiste en una SQLite (event-log) y lo sirve a una UI React por WebSocket. Modelo de despliegue: instancia única global (un daemon y una SQLite compartida en `~/.matecito-ai/` para todos los proyectos, scopeado por proyecto).

## Cómo usar este índice

1. Identificá qué tipo de tarea estás por hacer.
2. Encontrá el dominio correspondiente abajo y abrí su `INDEX.md`.
3. Leé los EDRs relevantes antes de escribir código.
4. Si hay contradicción entre tu plan y un EDR: pará y preguntale al usuario.

## Dominios de este proyecto

(Solo se listan los dominios que tienen al menos un EDR-archivo.)

| Dominio | Qué agrupa | Índice |
|---|---|---|
| `structure` | Estilo de arquitectura, layout de carpetas y convenciones de código | [structure/INDEX.md](structure/INDEX.md) |
| `runtime` | Manejo de errores, concurrencia y resiliencia en ejecución | [runtime/INDEX.md](runtime/INDEX.md) |
| `data` | Acceso a datos y modelado del event-log | [data/INDEX.md](data/INDEX.md) |
| `contracts` | Superficies de contrato (HTTP-in del MCP, WS-out a la UI) | [contracts/INDEX.md](contracts/INDEX.md) |
| `observability` | Logging estructurado y health checks | [observability/INDEX.md](observability/INDEX.md) |
| `delivery` | Configuración, testing, empaquetado y despliegue | [delivery/INDEX.md](delivery/INDEX.md) |
| `tech` | Tecnologías concretas elegidas | [tech/INDEX.md](tech/INDEX.md) — **consultá siempre antes de instalar algo nuevo** |

**Leyenda de status:** `Accepted` = vigente · `Pending` = decidir más adelante · `Not Applicable` = decidido que no aplica · `Deferred` = postergado con condición.

> Para EDRs `Pending`/`Deferred`, leé la sección "Razón de omisión / aplazamiento" del archivo; para los `Not Applicable`, la razón está en la sección "No aplican" del INDEX del dominio (o "Dominios sin uso" del raíz). **No asumas que la falta de decisión es un olvido** — está documentada.

## Dominios sin uso en este proyecto

(Dominios cuyas fases quedaron todas `Not Applicable` — no tienen carpeta. Se listan acá para dejar constancia de que se consideraron.)

| Dominio | Razón |
|---|---|
| `security` | Local-first, single-user. `auth` / `authorization` / `rate-limiting` / `secrets-management`: sin superficie de autenticación/autorización ni secretos. `cors`: la UI se sirve desde el broker (mismo origen), sin request cross-origin de browser; el handshake de WebSocket valida Origin==Host por default, que en mismo-origen pasa — decisión consciente, no olvido. |
| `quality` | Single-user local. `scalability`: sin escala horizontal. `i18n`: sin internacionalización. |
| `frontend` | `accessibility`: fuera del alcance del broker (es backend); la a11y de la UI se trata en el dominio de diseño, no acá. |

## Estado y mantenimiento

- Última actualización: 2026-07-23
- **Actualizar una decisión:** editá el EDR en el lugar, sea cambio menor o de fondo. El historial lo lleva git.
- **Decisión nueva:** creá el EDR en su dominio y sumá la fila al índice de ese dominio (y, si el dominio es nuevo en el proyecto, a este índice raíz).
