# Tier compartido

`payload/shared/` entrega **componentes transversales** que se despliegan a **todos** los dominios activos, sin importar cuáles tengas instalados. Son single-source: viven una sola vez acá y se aplanan dentro de los árboles compartidos `~/.claude/...` en el deploy — no se duplican por dominio.

El **mecanismo** de deploy (aplanamiento, reglas de colisión, hooks siempre activos vía `hook.SharedDomain`) está documentado en [`../domains/README.md`](../domains/README.md) → sección "## Shared tier". Acá no lo re-explicamos: este README es el **catálogo** de QUÉ entrega el tier.

## Componentes

Hoy el tier compartido **no entrega ningún componente**: `skills/`, `agents/` y `references/` están reservados como placeholders. Cuando aparezca una skill, un agente o una referencia genuinamente cross-domain —que valga para todos los dominios activos, sin importar cuáles tengas instalados— se cataloga acá.

## Ver también

- [Contrato de área](../domains/README.md) — incluye el mecanismo de deploy del tier compartido ("## Shared tier").
- [README raíz del ecosistema](../../README.md) — visión general de matecito-ai.
