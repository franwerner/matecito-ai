---
name: spacing-grid
depth: light
domain: foundation
type: convention
source: Material Design (8dp grid) · W3C Design Tokens (dimension)
---

# Fase: Escala de espaciado

## Qué decide

La escala de espaciado base (unidad y step) que gobierna paddings, gaps y márgenes en todo el sistema. Sin una escala, cada frame elige spacings sueltos y el ritmo visual se rompe.

## Preguntas

Una o dos, según haga falta.

### 1. Unidad y step base

> Una escala consistente hace que el spacing se sienta intencional y se pueda tokenizar; valores arbitrarios (13px, 17px) no.

- **Base 8 (8, 16, 24, 32, 40…)** — *default; estándar de industria (Material 8dp grid), divisible y predecible.*
- Base 4 (4, 8, 12, 16…) — más granular, para UI densa.
- Base 10 — *menos común; solo si la marca lo pide.*
- No sé, recomendame.

### 2. Materialización como tokens de spacing

> Un gap suelto repetido a mano no es chequeable; un token de spacing nombrado sí.

- **Cada step es un token de spacing nombrado (`Spacing/sm`=8, `Spacing/md`=16…)** — *default; condición para chequear drift.*
- Escala documentada pero aplicada a mano — parcial, no verificable automáticamente.
- No sé, recomendame.

## Notas de lógica (para el motor)

- Si la pieza es un `marketing-asset` único sin sistema, materializá con `Status: Pending` y razón ("pieza única sin sistema reutilizable"), no como decisión hueca.
- Depende de `design-tokens` si el equipo expresa la escala como Figma Variables.

## Qué materializar

DDR `spacing-grid` materializado según el template `~/.claude/references/ddr/templates/ddr.md`. La **Decisión** captura: la unidad base, el step, y la lista de tokens de spacing con su valor px concreto.

**Reglas verificables** (cada una con su mecanismo al inicio):

- **[tool: figma]** todo padding/gap/margin de los frames del sistema es un múltiplo de la unidad base decidida (ej: múltiplo de 8); no hay valores fuera de la escala.
- **[tool: figma]** los tokens de spacing nombrados existen con el valor px declarado (`Spacing/md` = 16); no hay spacings de sistema aplicados como valor suelto.

**Alcance:** la lista de tokens de spacing nombrados (`Spacing/*`) y los frames del sistema que la escala gobierna — el ancla de drift contra Figma.

Si la pieza es un asset único sin sistema, el DDR va con `Status: Pending` y la razón concreta; en ese caso no lleva Reglas verificables.
