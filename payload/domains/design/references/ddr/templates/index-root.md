<!-- Canonical template: índice RAÍZ (`.matecito-ai/ddr/INDEX.md`). Enruta por surface. Consumido por la fase de Materialización de design-decisions-bootstrap/SKILL.md. -->

# Design Decision Records — Índice raíz

Las decisiones de diseño están organizadas por **surface**. Este índice te dice qué surface mirar; el detalle de cada decisión está en el índice de su surface.

## Cómo usar este índice

1. Identificá qué tipo de trabajo de diseño estás por hacer.
2. Encontrá la surface correspondiente abajo y abrí su `INDEX.md`.
3. Leé los DDRs relevantes antes de tocar el sistema visual.
4. Si hay contradicción entre tu plan y un DDR: pará y preguntale al usuario.

## Surfaces de esta pieza / sistema

(Solo se listan las surfaces que tienen al menos un DDR-archivo.)

| Surface | Qué agrupa | Índice |
|---|---|---|
| `<surface>` | <una línea> | [<surface>/INDEX.md](<surface>/INDEX.md) |
| ... | | |

**Leyenda de status:** `Accepted` = vigente · `Pending` = decidir más adelante · `Not Applicable` = decidido que no aplica · `Deferred` = postergado con condición · `Superseded` = reemplazado por otro DDR · `Inferred` = borrador minado de Figma, sin porqué (ratificar vía `design-decisions-bootstrap` modo update).

> Para DDRs `Pending`/`Deferred`, leé la sección "Razón de omisión / aplazamiento" del archivo; para los `Not Applicable`, la razón está en la sección "No aplican" del INDEX de la surface (o "Surfaces sin uso" del raíz). **No asumas que la falta de decisión es un olvido** — está documentada.

## Surfaces sin uso en esta pieza / sistema

(Surfaces cuyas fases quedaron todas `Not Applicable` — no tienen carpeta. Se listan acá para dejar constancia de que se consideraron.)

| Surface | Razón |
|---|---|
| `<surface>` | <1 línea: por qué no aplica a este tipo de pieza> |
| ... | |

## Estado y mantenimiento

- Última actualización: <YYYY-MM-DD>
- **Actualizar una decisión (cambio menor):** editá el DDR. El historial lo lleva git.
- **Cambiar una decisión (cambio de fondo):** creá un DDR nuevo en la misma surface, marcá el viejo `Superseded` con link al nuevo. No edites la decisión vieja en el lugar.
- **Decisión nueva:** creá el DDR en su surface y sumá la fila al índice de esa surface (y, si la surface es nueva en la pieza, a este índice raíz).
- **Ratificar un `Inferred`:** vía `design-decisions-bootstrap` modo update — entrevistá el porqué, llená Contexto/Decisión/Alternativas/Consecuencias, descartá `## Evidencia (inferida)`, pasá a `Accepted`.
