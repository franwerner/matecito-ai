# Surface: `layout`

Cómo se compone la página: la grilla de layout, los breakpoints responsive y las reglas de composición. Define la *estructura espacial* en la que viven los componentes.

## Criterio de pertenencia

Un concern nuevo va en `layout` si trata sobre la **composición de la página o pantalla** (grilla, columnas, breakpoints, regiones). Si trata sobre la escala de espaciado como token base, va en `foundation` (`spacing-grid`); si trata sobre una pieza reutilizable, va en `components`.

## Concerns en esta surface

| Concern | Prof. | Type | Qué decide |
|---|---|---|---|
| [grid-breakpoints](grid-breakpoints.md) | light | convention | La grilla de columnas y los breakpoints responsive que gobiernan la composición en cada tamaño de pantalla. |
