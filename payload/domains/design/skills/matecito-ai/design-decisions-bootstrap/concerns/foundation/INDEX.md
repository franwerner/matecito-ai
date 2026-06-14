# Surface: `foundation`

Los tokens base sobre los que se construye todo lo demás: paleta de color, escala tipográfica, espaciado/grid y la estrategia de design tokens. Define el *vocabulario* visual del sistema.

## Criterio de pertenencia

Un concern nuevo va en `foundation` si trata sobre un **token o escala base reutilizable** (color, tipografía, espaciado, elevación, radios). Si trata sobre cómo se ensamblan en piezas reutilizables, va en `components`; si trata sobre la composición de la página, va en `layout`.

## Concerns en esta surface

| Concern | Prof. | Type | Qué decide |
|---|---|---|---|
| [color-palette](color-palette.md) | deep | decision | La paleta de color del sistema: primarios, secundarios, semánticos y neutros, expresada como color styles nombrados. |
| [design-tokens](design-tokens.md) | light | decision | Si las decisiones visuales se expresan como tokens nombrados (Figma Variables / styles) y con qué taxonomía de nombres. |
| [spacing-grid](spacing-grid.md) | light | convention | La escala de espaciado base (step y unidad) que gobierna paddings, gaps y márgenes en todo el sistema. |
| [type-scale](type-scale.md) | deep | decision | La escala tipográfica: familias, ratio modular, pesos y los text styles nombrados que la materializan. |
