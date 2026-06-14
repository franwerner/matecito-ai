---
name: contrast-target
depth: light
domain: accessibility
type: policy
source: WCAG 2.2 (W3C — contrast minimum 1.4.3 / 1.4.11)
---

# Fase: Objetivo de contraste

## Qué decide

El nivel de contraste WCAG objetivo (AA/AAA) y los ratios mínimos verificables para texto y elementos de UI. Es el piso de accesibilidad visual que toda pieza debe cumplir, chequeable directamente contra los colores del archivo Figma.

## Preguntas

Una o dos, según haga falta.

### 1. Nivel de contraste objetivo

> El nivel define la exigencia legal y de UX. WCAG 2.2 separa nivel AA (estándar legal en la mayoría de los países) de AAA (máximo).

- **WCAG 2.2 AA** — *default recomendado: texto normal ≥ 4.5:1, texto grande ≥ 3:1, componentes/UI no textuales ≥ 3:1.*
- WCAG 2.2 AAA — texto normal ≥ 7:1, texto grande ≥ 4.5:1; solo si se requiere explícitamente.
- Sin objetivo formal por ahora — *solo piezas internas sin obligación.*
- No sé, recomendame.

### 2. Alcance del chequeo

> Una regla de contraste sin definir sobre qué pares aplica deja huecos. **Solo si se eligió un nivel.**

- **Todos los pares texto/fondo + estados de componentes (placeholder, disabled) + íconos informativos** — *default; cobertura completa.*
- Solo texto sobre fondos principales — mínimo, deja huecos en estados.

## Notas de lógica (para el motor)

- **Default según tipo de pieza:** Requerido para `app-ui`/`landing`. Para `marketing-asset` se recomienda AA en el texto legible aunque la pieza sea expresiva.
- Provee el ratio que `color-palette` debe respetar en sus pares texto/fondo.
- Si se eligió "sin objetivo formal", materializá con `Status: Pending` y el trigger esperado (ej: "antes del lanzamiento público").

## Qué materializar

DDR `contrast-target` materializado según el template `~/.claude/references/ddr/templates/ddr.md`. La **Decisión** captura: el nivel WCAG objetivo (AA/AAA) y los ratios mínimos concretos por tipo de contenido (texto normal, texto grande, UI no textual).

**Reglas verificables** (cada una con su mecanismo al inicio):

- **[tool: contrast]** todo texto normal contra su fondo cumple el ratio del nivel decidido (AA → ≥ 4.5:1; AAA → ≥ 7:1).
- **[tool: contrast]** todo texto grande (≥ 18.66px bold o ≥ 24px) cumple el ratio reducido (AA → ≥ 3:1; AAA → ≥ 4.5:1).
- **[tool: contrast]** los elementos de UI no textuales y bordes de foco cumplen ≥ 3:1 contra su entorno (WCAG 1.4.11).
- **[tool: contrast]** los estados de componente que comunican información (placeholder, disabled, error) cumplen el ratio aplicable a su rol.

Si se eligió "sin objetivo formal por ahora", el DDR va con `Status: Pending` indicando el trigger esperado; en ese caso no lleva Reglas verificables exigibles.
