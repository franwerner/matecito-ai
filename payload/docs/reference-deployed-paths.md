# Referenciar rutas desplegadas

Lo que vive en `payload/` es **fuente**; lo que se ejecuta vive en `~/.claude/`.
Un agente o una skill corre en el repo **del usuario**, donde `payload/` no existe.

> **Regla:** ninguna instrucción que se lea en runtime puede nombrar una ruta
> `payload/…`. Si el texto le dice a alguien "leé este archivo", esa ruta tiene
> que ser la **desplegada**.

## Mapeo fuente → desplegado

| En el repo | En la máquina del usuario |
| --- | --- |
| `payload/domains/<id>/agents/<x>.md` | `~/.claude/agents/<x>.md` |
| `payload/domains/<id>/skills/<group>/<x>/SKILL.md` | `~/.claude/skills/<x>/SKILL.md` |
| `payload/domains/<id>/references/<...>` | `~/.claude/references/<...>` |
| `payload/domains/<id>/CLAUDE.md` | `~/.claude/matecito-ai/domains/<id>.md` |
| `payload/shared/<component>/<...>` | igual que su equivalente de dominio |

La capa `<group>` bajo `skills/` es organizativa y se descarta en el deploy
(detalle completo en [`../domains/README.md`](../domains/README.md)). Por eso la
ruta desplegada de una skill **nunca** incluye ni el dominio ni el grupo.

`payload/docs/` —donde vive este archivo— es documentación del repo: **no se
despliega** y nadie la lee en runtime.

## Cómo referenciar cada cosa

**Skills — por nombre, sin ruta, cuando se pueda.** El frontmatter de un agente
acepta el campo `skills:`, que toma nombres y precarga el contenido completo al
arrancar:

```yaml
tools: Read, Grep, Glob
skills:
  - development-spec-mine
```

En el cuerpo alcanza con "tu skill `X` viene precargada — seguila exactamente".

**Excepción:** una skill con `disable-model-invocation: true` **no se puede
precargar** (solo la invoca el usuario). Esos agentes leen su ruta desplegada:
`~/.claude/skills/<x>/SKILL.md`. Verificá el flag antes de usar `skills:`.

Para skills **condicionales** —las que se cargan solo si aplican— agregá `Skill`
a `tools:` y nombralas en la prosa; no hace falta ruta.

**References y archivos compartidos — ruta desplegada.** Templates y
`_shared/sdd-phase-common.md` no son skills, así que se referencian por su ruta
bajo `~/.claude/…`, que es portable.

**Agentes — por nombre.** Un agente se nombra (`development-spec-mine`), nunca
por archivo.

## La excepción legítima

Nombrar `payload/…` está bien cuando el texto habla del **repo matecito-ai**, no
del proyecto del usuario — por ejemplo, explicar que para sumar una entrada al
catálogo compartido hay que editar el repo. Decilo explícito ("en el repo
matecito-ai") y usá la ruta completa con `domains/<id>/`.

## Checklist antes de mergear

1. ¿Alguna instrucción de runtime dice `payload/`? → cambiala por la desplegada
   o por el nombre.
2. ¿Agregaste `skills:`? → confirmá que esa skill no tiene
   `disable-model-invocation`.
3. ¿La ruta de skill que escribiste incluye el dominio o el grupo? → sobra: es
   `~/.claude/skills/<x>/SKILL.md`.

Ojo con probar dentro de este repo: acá `payload/` existe, así que una ruta rota
igual se encuentra y el error queda oculto hasta que alguien la corre afuera.
