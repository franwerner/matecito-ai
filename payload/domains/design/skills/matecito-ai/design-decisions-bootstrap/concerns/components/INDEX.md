# Surface: `components`

Las piezas reutilizables ensambladas a partir de los tokens de `foundation`: la biblioteca de componentes, sus variantes y estados. Define las *unidades de construcción* de la UI.

## Criterio de pertenencia

Un concern nuevo va en `components` si trata sobre una **pieza reutilizable** (componente, set de variantes, estado) y su contrato visual. Si trata sobre un token base, va en `foundation`; si trata sobre cómo se disponen los componentes en la página, va en `layout`.

## Concerns en esta surface

| Concern | Prof. | Type | Qué decide |
|---|---|---|---|
| [component-library](component-library.md) | deep | decision | Qué componentes son reutilizables, su nivel de atomicidad y cómo se modelan variantes y estados como componentes nombrados en Figma. |
