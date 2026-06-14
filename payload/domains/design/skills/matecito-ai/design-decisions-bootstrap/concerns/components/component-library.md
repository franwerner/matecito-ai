---
name: component-library
depth: deep
domain: components
type: decision
source: Atomic Design (Brad Frost) · Material Design (components) · Figma component/variant model
---

# Fase: Biblioteca de componentes

## Qué decide

Qué piezas son reutilizables, en qué nivel de atomicidad, y cómo se modelan sus variantes y estados como componentes nombrados en Figma. Es lo que evita que cada pantalla redibuje un botón a mano y lo que hace la UI consistente y mantenible.

## Preguntas

Una por turno. Para cada una: línea de "por qué importa", opciones con default, y siempre la opción "no sé, recomendame".

### 1. Nivel de atomicidad

> Define qué cuenta como componente reutilizable. Demasiado fino genera ruido; demasiado grueso, duplicación.

- **Atomic Design (átomos → moléculas → organismos)** — *default para `brand-system`; vocabulario compartido y escalable.*
- Solo componentes de UI comunes (botón, input, card, modal) sin jerarquía formal — *default para `app-ui` mediano.*
- Ad-hoc por pantalla — *no recomendado; lleva a duplicación.*
- No sé, recomendame.

### 2. Modelado de variantes

> Las variantes (tamaño, jerarquía, ícono) deben vivir como propiedades de un componente, no como copias sueltas, para que sean chequeables y mantenibles.

- **Component set con variant properties (`size`, `variant`, `state`)** — *default; el modelo nativo de Figma, verificable.*
- Componentes separados por variante (`Button/Primary`, `Button/Secondary`) — más simple, menos potente.
- Mix — base como set, casos raros sueltos.
- No sé, recomendame.

### 3. Estados interactivos

> Una pieza interactiva sin estados definidos (hover, focus, disabled) deja huecos que cada implementador inventa distinto.

- **Estados explícitos como variant property (`state: default|hover|focus|disabled`)** — *default para `app-ui`.*
- Solo default + disabled — para piezas simples.
- Sin estados modelados — *solo para `marketing-asset` estático.*

### 4. Tokenización del componente

> Un componente que usa valores sueltos en vez de tokens rompe la trazabilidad con foundation.

- **El componente referencia color/text/spacing styles, no valores crudos** — *default; condición para chequear drift.*
- Parcialmente tokenizado — aceptable transitorio, anotar deuda.
- No sé, recomendame.

## Notas de lógica (para el motor)

- **Default según tipo de pieza:** `brand-system` → proponé Atomic Design en la pregunta 1. `app-ui` → proponé set de componentes comunes. `marketing-asset` → este concern suele ser `Not Applicable`.
- Depende de `color-palette`, `type-scale`, `spacing-grid`: los componentes consumen esos styles, no redefinen valores.

## Qué materializar

DDR `component-library` materializado según el template `~/.claude/references/ddr/templates/ddr.md`. La **Decisión** captura: el nivel de atomicidad, la lista de componentes reutilizables nombrados, cómo se modelan variantes (variant properties concretas) y estados.

**Reglas verificables** (cada una con su mecanismo al inicio):

- **[tool: figma]** cada pieza reutilizable enumerada existe como un componente (o component set) nombrado en el archivo; no hay piezas reutilizables como grupos sueltos copiados.
- **[tool: figma]** las variantes se modelan como variant properties del component set (ej: `Button` con `variant`, `size`, `state`), no como componentes duplicados ad-hoc.
- **[tool: figma]** los estados interactivos decididos existen como valores de la property `state`; no falta ningún estado declarado en la Decisión.
- **[tool: figma]** los fills/textos/spacings del componente referencian styles nombrados de `foundation`, no valores crudos inline.
- **[manual]** una instancia del componente no tiene overrides que contradigan la decisión (ej: un Button con color fuera de la paleta).

**Alcance:** la lista de componentes y component sets nombrados que forman la biblioteca (`Button`, `Input`, `Card`, …) — el ancla de drift contra Figma.
