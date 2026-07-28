# Development Decision Records — Índice raíz

Las decisiones están organizadas por **dominio**. Este índice te dice qué dominio mirar; el detalle de cada decisión está en el índice de su dominio.

Contexto del proyecto: **Matecito UI** — cockpit React (`web-spa`) bajo `apps/ui` que consume el broker (`apps/api`, Go) por WebSocket (push en vivo) y HTTP (snapshot). Canvas-first (ReactFlow), inspector drawer, timeline scrubber, keyboard-first, tema dual (default light), read-first. Single-dev, greenfield, local-first; se sirve como build estático embebido en el binario del broker (same-origin).

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
| `frontend` | Accesibilidad, acceso a datos, estilos/theming y ruteo | [frontend/INDEX.md](frontend/INDEX.md) |
| `runtime` | Manejo de errores y resiliencia en ejecución | [runtime/INDEX.md](runtime/INDEX.md) |
| `security` | Validación de entrada y escaneo de dependencias | [security/INDEX.md](security/INDEX.md) |
| `delivery` | Configuración, empaquetado/despliegue, gates de calidad y documentación | [delivery/INDEX.md](delivery/INDEX.md) |
| `quality` | Performance e internacionalización | [quality/INDEX.md](quality/INDEX.md) |
| `tech` | Tecnologías concretas elegidas | [tech/INDEX.md](tech/INDEX.md) — **consultá siempre antes de instalar algo nuevo** |

**Leyenda de status:** `Accepted` = vigente · `Pending` = decidir más adelante · `Not Applicable` = decidido que no aplica · `Deferred` = postergado con condición.

> Para EDRs `Pending`/`Deferred`, leé la sección "Razón de omisión / aplazamiento" del archivo; para los `Not Applicable`, la razón está en la sección "No aplican" del INDEX del dominio (o "Dominios sin uso" del raíz). **No asumas que la falta de decisión es un olvido** — está documentada.

## Dominios sin uso en este proyecto

(Dominios cuyas fases quedaron todas `Not Applicable` — no tienen carpeta. Se listan acá para dejar constancia de que se consideraron.)

| Dominio | Razón |
|---|---|
| `observability` | La UI es un cliente: `logging` es una decisión consciente de no hacer logging formal del cliente en esta etapa (los errores se surfacean por UI/consola); `metrics`, `tracing` y `health-checks` son concerns de backend, no aplican a un cliente. |
| `data` | La UI no persiste datos propios: `data-access` y `data-modeling` no aplican; su estado es efímero + remoto del broker. |
| `contracts` | La UI consume el contrato del broker, no expone ninguno: `api-contract`, `cli-contract`, `library-contract` y `event-contract` no aplican. |

## Estado y mantenimiento

- Última actualización: 2026-07-23
- **Actualizar una decisión:** editá el EDR en el lugar, sea cambio menor o de fondo. El historial lo lleva git.
- **Decisión nueva:** creá el EDR en su dominio y sumá la fila al índice de ese dominio (y, si el dominio es nuevo en el proyecto, a este índice raíz).
