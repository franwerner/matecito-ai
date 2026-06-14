# Surface: `accessibility`

Los objetivos verificables de accesibilidad visual: contraste, tamaños mínimos, y cómo se chequean contra el archivo Figma. Define el *piso de accesibilidad* que toda pieza debe cumplir.

## Criterio de pertenencia

Un concern nuevo va en `accessibility` si trata sobre una **regla de accesibilidad visual verificable** (contraste, tamaño táctil/tipográfico, foco visible) chequeable contra el diseño. Las reglas de implementación técnica de a11y en código viven en el dominio development (`frontend/accessibility`), no acá.

## Concerns en esta surface

| Concern | Prof. | Type | Qué decide |
|---|---|---|---|
| [contrast-target](contrast-target.md) | light | policy | El nivel de contraste WCAG objetivo (AA/AAA) y los ratios mínimos verificables para texto y elementos de UI. |
