# Tier compartido

`payload/shared/` entrega **componentes transversales** que se despliegan a **todos** los dominios activos, sin importar cuáles tengas instalados. Son single-source: viven una sola vez acá y se aplanan dentro de los árboles compartidos `~/.claude/...` en el deploy — no se duplican por dominio.

El **mecanismo** de deploy (aplanamiento, reglas de colisión, hooks siempre activos vía `hook.SharedDomain`) está documentado en [`../domains/README.md`](../domains/README.md) → sección "## Shared tier". Acá no lo re-explicamos: este README es el **catálogo** de QUÉ entrega el tier.

## Componentes

| Capa | Pieza | Rol |
| --- | --- | --- |
| **Skill** | `roadmap` | Capa de planificación y continuidad por encima del flujo SDD/design: organiza QUÉ hay que hacer en fases, mientras el flujo ejecuta CÓMO. |

### `roadmap`

Guía al usuario para definir un roadmap multi-fase de forma conversacional, rastrea el avance paso a paso, y emite un "next context prompt" al cerrar cada sesión para que la siguiente arranque sin perder continuidad.

- **Invocación:** `/roadmap new <titulo>` · `continue` · `next`, más disparadores en lenguaje natural ("armar un roadmap", "plan de implementación por fases", "qué sigue en mi roadmap", …).
- **Layout en runtime:** los artefactos viven en el **proyecto del usuario**, en `.matecito-ai/roadmaps/<titulo>/` (hermano de `adr/` y `ddr/`): un `INDEX.md` por roadmap (objetivo, scope, dominio(s), rollup de progreso machine-readable) más un `STEP-N.md` por paso.
- **Forma de un `STEP-N.md`:** header de `status` (`pending | in-progress | done`), un checklist de tareas, una sección `## Pendientes` (loose ends tickables, con carry-forward al step siguiente cuando el step se marca `done`), y un `## Next context prompt` para retomar en una sesión nueva.
- **Handoff a flujos:** cuando un step mapea a desarrollo o diseño, `/roadmap next` PROPONE el comando de flujo (`/sdd-new` · `/design-new`) con el scope pre-llenado; nunca ejecuta autónomamente y siempre pasa por el INTAKE GATE del flujo destino.

Esto es una entrada de catálogo, no la spec completa: la fuente de verdad es el [`SKILL.md`](skills/matecito-ai/roadmap/SKILL.md) de la skill (con sus plantillas `INDEX-TEMPLATE.md` y `STEP-TEMPLATE.md`).

## Placeholders

El tier reserva además `agents/` y `references/` para componentes transversales, pero **hoy no entregan ninguno** — son placeholders. Cuando aparezca un agente o una referencia cross-domain, se cataloga acá.

## Ver también

- [Contrato de área](../domains/README.md) — incluye el mecanismo de deploy del tier compartido ("## Shared tier").
- [README raíz del ecosistema](../../README.md) — visión general de matecito-ai.
