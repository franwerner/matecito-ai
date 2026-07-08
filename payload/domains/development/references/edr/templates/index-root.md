<!-- Canonical template: índice RAÍZ (`.matecito-ai/edr/INDEX.md`). Enruta por dominio. Consumido por la fase de Materialización de SKILL.md. -->

# Development Decision Records — Índice raíz

Las decisiones están organizadas por **dominio**. Este índice te dice qué dominio mirar; el detalle de cada decisión está en el índice de su dominio.

## Cómo usar este índice

1. Identificá qué tipo de tarea estás por hacer.
2. Encontrá el dominio correspondiente abajo y abrí su `INDEX.md`.
3. Leé los EDRs relevantes antes de escribir código.
4. Si hay contradicción entre tu plan y un EDR: pará y preguntale al usuario.

## Dominios de este proyecto

(Solo se listan los dominios que tienen al menos un EDR-archivo.)

| Dominio | Qué agrupa | Índice |
|---|---|---|
| `<dominio>` | <una línea> | [<dominio>/INDEX.md](<dominio>/INDEX.md) |
| ... | | |
| `tech` | Tecnologías concretas elegidas | [tech/INDEX.md](tech/INDEX.md) — **consultá siempre antes de instalar algo nuevo** |

**Leyenda de status:** `Accepted` = vigente · `Pending` = decidir más adelante · `Not Applicable` = decidido que no aplica · `Deferred` = postergado con condición · `Superseded` = reemplazado por otro EDR.

> Para EDRs `Pending`/`Deferred`, leé la sección "Razón de omisión / aplazamiento" del archivo; para los `Not Applicable`, la razón está en la sección "No aplican" del INDEX del dominio (o "Dominios sin uso" del raíz). **No asumas que la falta de decisión es un olvido** — está documentada.

## Dominios sin uso en este proyecto

(Dominios cuyas fases quedaron todas `Not Applicable` — no tienen carpeta. Se listan acá para dejar constancia de que se consideraron.)

| Dominio | Razón |
|---|---|
| `<dominio>` | <1 línea: por qué no aplica al proyecto> |
| ... | |

## Estado y mantenimiento

- Última actualización: <YYYY-MM-DD>
- **Actualizar una decisión (cambio menor):** editá el EDR. El historial lo lleva git.
- **Cambiar una decisión (cambio de fondo):** creá un EDR nuevo en el mismo dominio, marcá el viejo `Superseded` con link al nuevo. No edites la decisión vieja en el lugar.
- **Decisión nueva:** creá el EDR en su dominio y sumá la fila al índice de ese dominio (y, si el dominio es nuevo en el proyecto, a este índice raíz).
