# Tier compartido

`payload/shared/` entrega **componentes transversales** que se despliegan a **todos** los dominios activos, sin importar cuáles tengas instalados. Son single-source: viven una sola vez acá y se aplanan dentro de los árboles compartidos `~/.claude/...` en el deploy — no se duplican por dominio.

El **mecanismo** de deploy (aplanamiento, reglas de colisión, hooks siempre activos vía `hook.SharedDomain`) está documentado en [`../docs/build-a-domain.md`](../docs/build-a-domain.md) → sección "## Shared tier". Acá no lo re-explicamos: este README es el **catálogo** de QUÉ entrega el tier.

## Componentes

`skills/` y `agents/` están reservados como placeholders: todavía no entregan ningún componente. Cuando
aparezca una skill o un agente genuinamente cross-domain —que valga para todos los dominios activos, sin
importar cuáles tengas instalados— se cataloga acá.

`references/` entrega:

- [`gate-presentation.md`](references/gate-presentation.md) — el recorrido único (índice, item por
  item, "confirmar el resto", retomar) y la plantilla de huecos fijos que usa cualquier gate que
  ratifique items. Se despliega a `~/.claude/references/gate-presentation.md`.

## Ver también

- [Contrato de área](../docs/build-a-domain.md) — incluye el mecanismo de deploy del tier compartido ("## Shared tier").
- [README raíz del ecosistema](../../README.md) — visión general de matecito-ai.
