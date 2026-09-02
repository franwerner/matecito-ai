# Tier compartido

`payload/shared/` entrega **componentes transversales** que se despliegan a **todos** los dominios activos, sin importar cuáles tengas instalados. Son single-source: viven una sola vez acá y se aplanan dentro de los árboles compartidos `~/.claude/...` en el deploy — no se duplican por dominio.

El **mecanismo** de deploy (aplanamiento, reglas de colisión, hooks siempre activos vía `hook.SharedDomain`) está documentado en [`../docs/build-a-domain.md`](../docs/build-a-domain.md) → sección "## Shared tier". Acá no lo re-explicamos: este README es el **catálogo** de QUÉ entrega el tier.

## Componentes

`skills/` entrega:

- **`find-records`** — localiza los records durables que gobiernan una parte del proyecto (EDR, DDR,
  capability-specs) antes de escribir código, proponer un diseño o verificar contra ellos, combinando
  el índice del store, la búsqueda literal y la búsqueda por significado.
- **`setup-record-search`** — configura la búsqueda por significado sobre esos records: crea una
  colección por store, declara sus exclusiones y las credenciales de embeddings, e indexa.

`agents/` está reservado como placeholder: todavía no entrega ningún componente. Cuando aparezca un
agente genuinamente cross-domain —que valga para todos los dominios activos, sin importar cuáles tengas
instalados— se cataloga acá.

`references/` entrega:

- [`gate-presentation.md`](references/gate-presentation.md) — el recorrido único (índice, item por
  item, "confirmar el resto", retomar) y la plantilla de huecos fijos que usa cualquier gate que
  ratifique items. Se despliega a `~/.claude/references/gate-presentation.md`.
- [`side-discussion.md`](references/side-discussion.md) — el mecanismo de la discusión lateral: la
  forma del traspaso y de la conclusión, el requisito de lanzamiento (cinco propiedades, sin nombrar
  herramienta), y cómo se recoge la conclusión según el tipo declarado. Lo leen el orquestador y la
  sesión lateral misma, que no pertenece a ningún dominio. Se despliega a
  `~/.claude/references/side-discussion.md`.

## Ver también

- [Contrato de área](../docs/build-a-domain.md) — incluye el mecanismo de deploy del tier compartido ("## Shared tier").
- [README raíz del ecosistema](../../README.md) — visión general de matecito-ai.
