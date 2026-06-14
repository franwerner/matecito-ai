# Catálogo de design concerns — Índice raíz

Menú de fases (concerns) de diseño que la skill puede recorrer, **organizado por surface**. El motor (`SKILL.md`) lee este índice para (a) entender el mapa de surfaces y (b) armar la lista de fases relevantes según el tipo de pieza. Recién después lee el archivo individual de cada fase que va a tratar, para no cargar al contexto lo que no aplica.

## Cómo lo usa el motor

1. Fase 0 detecta el tipo de pieza → lo mapea a un token de abajo.
2. El motor recorre la matriz de aplicabilidad y arma dos grupos:
   - **Requerido** → se incluye por default (el usuario puede marcarlo `Not Applicable` con razón).
   - **Recomendado** → se sugiere; el usuario decide.
3. Muestra el set; el usuario elige qué definir ahora, y puede agregar una **fase custom**.
4. Por cada fase elegida, el motor lee `<surface>/<slug>.md` y la trata con las reglas del motor.
5. Las relevantes NO elegidas → DDR `Not Applicable` / `Pending` + razón. Nunca hueco silencioso.

## Surfaces canónicas (fijas)

La taxonomía de surfaces es **cerrada y la impone el motor** — la misma para el catálogo interno y para la salida `.matecito-ai/ddr/`, así todas las piezas del equipo se ven igual. Cada surface tiene su propio `INDEX.md` con el detalle de sus concerns y el **criterio de pertenencia** (cuándo un concern nuevo va ahí).

### Activas (con concerns)

| Surface | Qué agrupa | Índice |
|---|---|---|
| `foundation` | Tokens base: paleta de color, escala tipográfica, espaciado/grid, design tokens | [foundation/INDEX.md](foundation/INDEX.md) |
| `components` | Biblioteca de componentes, variantes, estados | [components/INDEX.md](components/INDEX.md) |
| `layout` | Grilla de layout, breakpoints, composición responsive | [layout/INDEX.md](layout/INDEX.md) |
| `brand` | Identidad: tono de voz, uso del logo, expresión de marca | [brand/INDEX.md](brand/INDEX.md) |
| `accessibility` | Objetivos verificables de accesibilidad visual (contraste, tamaños) | [accessibility/INDEX.md](accessibility/INDEX.md) |

## Tokens de tipo de pieza

`landing` · `app-ui` · `brand-system` · `marketing-asset`

`todos` = cualquier tipo.

## Fase custom

Si el usuario tiene un tema fuera de este catálogo, el motor le hace las preguntas genéricas, determina a qué **surface canónica** pertenece, crea `<surface>/<slug>.md` con el formato estándar y suma la fila al índice de esa surface + a la matriz de abajo. Antes de guardarlo pregunta: **¿reusable (queda en el catálogo) o solo para esta pieza (solo genera el DDR)?**

## Leyenda

- **Prof.** = profundidad: `deep` (cuestionario propio) · `light` (1-2 preguntas).
- **Tipo** = `decisión` (alternativas y trade-offs) · `convención` (acuerdo de estilo) · `política` (regla verificable).

---

## Matriz de aplicabilidad

Cada fila apunta a `<surface>/<slug>.md`. La columna **Surface** es la carpeta canónica (interna y de salida).

### foundation
| Fase | Surface | Prof. | Requerido para | Recomendado para |
|---|---|---|---|---|
| [color-palette](foundation/color-palette.md) | foundation | deep | todos | — |
| [type-scale](foundation/type-scale.md) | foundation | deep | todos | — |
| [spacing-grid](foundation/spacing-grid.md) | foundation | light | `app-ui`, `landing`, `brand-system` | `marketing-asset` |
| [design-tokens](foundation/design-tokens.md) | foundation | light | `brand-system`, `app-ui` | `landing` |

### components
| Fase | Surface | Prof. | Requerido para | Recomendado para |
|---|---|---|---|---|
| [component-library](components/component-library.md) | components | deep | `app-ui`, `brand-system` | `landing` |

### layout
| Fase | Surface | Prof. | Requerido para | Recomendado para |
|---|---|---|---|---|
| [grid-breakpoints](layout/grid-breakpoints.md) | layout | light | `app-ui`, `landing` | `brand-system` |

### brand
| Fase | Surface | Prof. | Requerido para | Recomendado para |
|---|---|---|---|---|
| [tone-of-voice](brand/tone-of-voice.md) | brand | light | `brand-system`, `marketing-asset` | `landing` |
| [logo-usage](brand/logo-usage.md) | brand | light | `brand-system` | `landing`, `marketing-asset` |

### accessibility
| Fase | Surface | Prof. | Requerido para | Recomendado para |
|---|---|---|---|---|
| [contrast-target](accessibility/contrast-target.md) | accessibility | light | `app-ui`, `landing` | `brand-system`, `marketing-asset` |

---

## Mantenimiento (ratchet)

- **Agregar una fase:** decidí a qué **surface canónica** pertenece (mirá el criterio de pertenencia en `<surface>/INDEX.md`). Creá `<surface>/<slug>.md` con el formato estándar (ver `foundation/color-palette.md` para `deep`, `foundation/spacing-grid.md` para `light`), sumá la fila al `<surface>/INDEX.md` y a la matriz de arriba.
- **No crear surfaces nuevas por pieza:** la taxonomía es fija. Si de verdad falta una surface, es una decisión de catálogo (agregarla acá y en el motor), no algo improvisado en un repo.
- **Origen del catálogo:** sembrado de W3C Design Tokens, Atomic Design, WCAG 2.x y Material Design. Nace casi completo y solo crece.
