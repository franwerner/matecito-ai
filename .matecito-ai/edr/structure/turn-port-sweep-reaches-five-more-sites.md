# EDR — El sweep de este cambio llega a cinco sitios que el Scope ratificado no nombraba

- **Status:** Accepted
- **Date:** 2026-09-10

## Contexto
El Scope de este cambio fija qué archivos toca el port del turno de rama compartida y el retiro del aislamiento a nivel de cambio. Durante la ejecución aparecieron cinco sitios adicionales, ninguno nombrado en esa tabla, cada uno falsificado por el mismo motivo que el Scope original: siguen describiendo un mecanismo de dos contenedores, o un flag/campo que este cambio retira, después de que el resto del payload ya dejó de tenerlo. Aparecieron también tres candidatos que, examinados, resultan residuo aceptado y no defecto — nombrados acá para que un sweep futuro no los vuelva a señalar.

## Decisión
El sweep alcanza los cinco sitios IN con un veredicto cada uno, y deja los tres OUT explícitamente intactos:

IN:
1. `payload/domains/development/README.md:33` — decía "los cuatro decision flags (`diagram`, `ui-test`, `components`, `worktree-isolation`)"; falsificado dos veces por el requisito de tres flags. Reescrito a "tres", nombrando los tres flags sobrevivientes.
2. `payload/domains/development/references/phase-returns/sdd-intake/sdd-intake.yaml:14,27` — **declaraba** la línea de campo `Worktree isolation` en el contrato de retorno. La proposal original sólo nombraba el `.md` de la plantilla; dejar el `.yaml` declarando un campo que intake ya no emite rompe el render, no sólo la prosa. Declaración y comentario eliminados.
3. `payload/domains/development/skills/gentle-ai/sdd-apply/SKILL.md:184,644` — dos formulaciones "contenedor inmediato ... o el espacio de trabajo propio del cambio `matecito-ai/<change-name>`" — exactamente la forma que el spec prohíbe en cualquier parte del payload. Ambas colapsadas a la rama de trabajo, una cláusula por sitio.
4. `payload/shared/references/side-discussion.md:71-72` — una quinta formulación de dos contenedores, en el árbol que se despliega a toda instalación. Colapsada a un único árbol de trabajo.
5. `.matecito-ai/edr/contracts/side-discussion-launcher-test.md:10` — registro Accepted cuya propiedad (4) es esa misma formulación de dos contenedores. Una cuarta reescritura (además de las tres del Scope original), colapsada al árbol de trabajo ordinario.

Además: `structure/side-discussion-prose-homes.md:23` es el único link markdown resoluble del store que apunta a un record retirado (`change-workspace-prose-homes.md`). Por el tratamiento ya ratificado de `structure/mirror-sites-join-the-prune-sweep.md` para exactamente este caso, la fila `relacionado-con` se **elimina**, no se convierte en nota — su criterio ya está enunciado en la línea 20 del propio registro, así que no se pierde nada.

OUT (residuo aceptado, nombrado para que no se re-señale):
- `.matecito-ai/development-specs/INDEX.md:38,40,45` — un ledger cronológico de cambios pasados.
- `.gitignore:16` — mantenido deliberadamente.
- `.matecito-ai/.markdown_vault_mcp/state.json` — un índice generado.

## Reglas verificables
- **[manual]** `payload/domains/development/README.md:33` nombra tres flags, no cuatro, y no nombra `worktree-isolation`.
- **[manual]** `payload/domains/development/references/phase-returns/sdd-intake/sdd-intake.yaml` no declara ninguna línea de campo `Worktree isolation` ni su comentario asociado.
- **[manual]** `payload/domains/development/skills/gentle-ai/sdd-apply/SKILL.md` no contiene, en ninguna línea, la formulación de dos contenedores ("o el espacio de trabajo propio del cambio").
- **[manual]** `payload/shared/references/side-discussion.md` abre su propiedad (5) sobre un único árbol de trabajo, sin segundo contenedor.
- **[manual]** `.matecito-ai/edr/contracts/side-discussion-launcher-test.md`, propiedad (4), colapsa al árbol de trabajo ordinario, sin cláusula condicional.
- **[manual]** `.matecito-ai/edr/structure/side-discussion-prose-homes.md` no contiene ningún link `relacionado-con` hacia un record retirado.
- **[manual]** Los tres sitios OUT (`development-specs/INDEX.md:38,40,45`, `.gitignore:16`, `.markdown_vault_mcp/state.json`) quedan sin editar por este cambio.

## Alternativas consideradas
Limitar el sweep estrictamente al Scope ratificado y dejar los cinco sitios adicionales para un cambio futuro. Descartada: cuatro de los cinco están en árboles que se despliegan a toda instalación o en un contrato de retorno que ya se está editando en este mismo cambio — dejarlos divergentes es entregar el cambio con el mismo defecto que el cambio existe para eliminar, visible en el primer grep que alguien corra. Barrer también los tres candidatos OUT. Descartada: son residuo aceptado por los mismos criterios ya ratificados en `structure/sweep-reach-outside-the-ratified-scope.md` y `structure/falsified-clause-joins-the-sweep.md` — un ledger cronológico, un `.gitignore` deliberado y un índice generado no son defectos que este sweep deba corregir.

## Consecuencias
El diff de este cambio crece en cinco archivos más allá del Scope original, todos justificados por el mismo criterio de falsificación que motivó el Scope. Un sweep futuro que vuelva a correr sobre este vocabulario encuentra los tres sitios OUT ya nombrados y no los reporta como hallazgo nuevo.
