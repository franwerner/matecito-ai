---
name: grid-breakpoints
depth: light
domain: layout
type: convention
source: Material Design (responsive layout grid) · Bootstrap/Tailwind breakpoint conventions
---

# Fase: Grilla y breakpoints

## Qué decide

La grilla de columnas y los breakpoints responsive que gobiernan la composición en cada tamaño de pantalla. Sin una grilla acordada, cada frame compone a ojo y el layout pierde consistencia entre pantallas.

## Preguntas

Una o dos, según haga falta.

### 1. Sistema de grilla

> La grilla define el esqueleto de la composición. Una grilla compartida hace que todas las pantallas se alineen.

- **12 columnas con gutter y margin fijos** — *default para `app-ui`/`landing` web; flexible y estándar.*
- 4/8 columnas (mobile-first) — para apps primariamente mobile.
- Sin grilla formal, layout libre — *solo para `marketing-asset` único.*
- No sé, recomendame.

### 2. Breakpoints

> Los breakpoints definen dónde cambia la composición. **Solo si la pieza es responsive.**

- **Mobile / Tablet / Desktop (3 puntos, ej: 360 / 768 / 1280)** — *default; cubre la mayoría de los casos.*
- Mobile / Desktop (2 puntos) — para piezas simples.
- Escala completa (xs…xl, 5 puntos) — para `app-ui` complejas.

## Notas de lógica (para el motor)

- Si la pieza es de un solo tamaño (`marketing-asset` fijo), no hagas la pregunta 2 y materializá con `Status: Pending` y razón ("pieza de tamaño único, sin responsive").
- Los anchos de gutter/margin deben ser múltiplos de la unidad de `spacing-grid`.

## Qué materializar

DDR `grid-breakpoints` materializado según el template `~/.claude/references/ddr/templates/ddr.md`. La **Decisión** captura: el número de columnas, gutter y margin concretos, y la lista de breakpoints con su ancho px.

**Reglas verificables** (cada una con su mecanismo al inicio):

- **[tool: figma]** los frames principales usan layout grids de Figma con el número de columnas, gutter y margin decididos; no hay columnas arbitrarias por frame.
- **[tool: figma]** existe un frame/variante por cada breakpoint declarado (ej: 360, 768, 1280); no falta ningún breakpoint de la Decisión.
- **[tool: figma]** gutter y margin son múltiplos de la unidad base de `spacing-grid`.

**Alcance:** los frames de página/pantalla y sus layout grids que la decisión gobierna, por breakpoint — el ancla de drift contra Figma.

Si la pieza es de tamaño único, el DDR va con `Status: Pending` y la razón concreta; en ese caso no lleva Reglas verificables.
