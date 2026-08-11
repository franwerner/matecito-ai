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
| `structure` | Decisiones sobre cómo se organiza el payload del repo (dónde vive cada concepto, cómo se citan entre sí). | [structure/INDEX.md](structure/INDEX.md) |
| `contracts` | Decisiones sobre los contratos entre piezas del ecosistema: qué forma tiene lo que entra y sale de una herramienta, y cómo se comunica un resultado. | [contracts/INDEX.md](contracts/INDEX.md) |
| `tech` | Tecnologías concretas elegidas | [tech/INDEX.md](tech/INDEX.md) — **consultá siempre antes de instalar algo nuevo** |

**Leyenda de status:** `Accepted` = vigente · `Pending` = decidir más adelante · `Not Applicable` = decidido que no aplica · `Deferred` = postergado con condición.

> Para EDRs `Pending`/`Deferred`, leé la sección "Razón de omisión / aplazamiento" del archivo; para los `Not Applicable`, la razón está en la sección "No aplican" del INDEX del dominio (o "Dominios sin uso" del raíz). **No asumas que la falta de decisión es un olvido** — está documentada.

## Dominios sin uso en este proyecto

(Dominios cuyas fases quedaron todas `Not Applicable` — no tienen carpeta. Se listan acá para dejar constancia de que se consideraron.)

| Dominio | Razón |
|---|---|
| — | Ninguno considerado todavía fuera de `structure`/`tech`. |

## Estado y mantenimiento

- Última actualización: 2026-08-10
- **Actualizar una decisión:** editá el EDR en el lugar, sea cambio menor o de fondo. El historial lo lleva git.
- **Decisión nueva:** creá el EDR en su dominio y sumá la fila al índice de ese dominio (y, si el dominio es nuevo en el proyecto, a este índice raíz).
