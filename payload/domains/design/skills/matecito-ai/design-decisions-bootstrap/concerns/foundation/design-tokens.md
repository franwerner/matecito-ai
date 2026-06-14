---
name: design-tokens
depth: light
domain: foundation
type: decision
source: W3C Design Tokens Community Group Format · Material Design (design tokens)
---

# Fase: Design tokens

## Qué decide

Si las decisiones visuales (color, tipografía, spacing, radios, elevación) se expresan como **tokens nombrados** (Figma Variables / styles) y con qué taxonomía de nombres. Es lo que convierte valores sueltos en un sistema verificable y portable a código.

## Preguntas

Una o dos, según haga falta.

### 1. Estrategia de tokens

> Tokens nombrados hacen las decisiones chequeables (mine/verify) y portables a código; valores crudos no.

- **Tokens nombrados para todo (Figma Variables + styles)** — *default para `app-ui` y `brand-system`; única forma de chequear drift contra Figma.*
- Solo styles (color/text/effect), sin variables — parcial; sirve para color y tipografía, no para spacing/radios.
- Sin tokens, valores aplicados a mano — *no recomendado para un sistema; rompe verificabilidad.*
- No sé, recomendame.

### 2. Taxonomía de nombres

> Una convención de nombres consistente hace los tokens predecibles y mapeables a código. **Solo si en la 1 eligió tokenizar.**

- **Jerárquica por rol (`color.primary.500`, `space.md`, `radius.lg`)** — *default; alineada con W3C Design Tokens.*
- Plana semántica (`bg-default`, `text-muted`) — buena para mapear directo a CSS vars.
- Mix (tokens base + alias semánticos) — más robusto, más trabajo.

## Notas de lógica (para el motor)

- Si en la pregunta 1 eligió "sin tokens", no hagas la pregunta 2. Materializá el DDR con `Status: Pending` y razón ("pieza sin sistema tokenizado todavía"), no como decisión hueca.
- Este concern es transversal: `color-palette`, `type-scale` y `spacing-grid` referencian la taxonomía decidida acá en lugar de redefinirla.

## Qué materializar

DDR `design-tokens` materializado según el template `~/.claude/references/ddr/templates/ddr.md`. La **Decisión** captura: la estrategia de tokens elegida (variables+styles / solo styles / ninguno) y la taxonomía de nombres concreta con ejemplos.

**Reglas verificables** (cada una con su mecanismo al inicio):

- **[tool: figma]** las decisiones visuales del sistema (color, tipografía, spacing, radios) están expresadas como Figma Variables o styles nombrados, no como valores crudos repetidos.
- **[tool: figma]** los nombres de los tokens siguen la taxonomía declarada (ej: `color.primary.500`, `space.md`); no hay nombres fuera del patrón.
- **[manual]** cada token tiene un único valor de verdad; no hay dos tokens con el mismo rol y distinto valor.

**Alcance:** el conjunto de colecciones de Figma Variables y styles que conforman el set de tokens — el ancla de drift contra Figma.

Si se eligió "sin tokens", el DDR va con `Status: Pending` y la razón concreta; en ese caso no lleva Reglas verificables.
