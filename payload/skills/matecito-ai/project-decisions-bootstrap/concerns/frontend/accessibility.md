---
name: accessibility
depth: light
domain: frontend
type: decision
source: WCAG 2.2 (W3C Web Content Accessibility Guidelines)
---

# Fase: Accesibilidad

## Qué decide

El nivel de conformidad WCAG objetivo y cómo se verifica que se cumple durante el desarrollo.

## Preguntas

### 1. Nivel WCAG objetivo

> El nivel define la exigencia legal y de UX. WCAG 2.2 tiene tres niveles: A (mínimo), AA (estándar de industria y legal en muchos países), AAA (máximo; raramente alcanzable en su totalidad).

- **WCAG 2.2 nivel AA** — *default recomendado; cubre los requisitos legales más comunes (ADA, EN 301 549, RGAA) y las expectativas de usuarios con discapacidad.*
- WCAG 2.2 nivel A — mínimo; solo para proyectos internos sin obligación legal.
- WCAG 2.2 nivel AAA — solo si el proyecto lo requiere explícitamente (ej: app gubernamental accesible).
- Sin objetivo formal por ahora — se atiende caso por caso.
- No sé, recomendame.

### 2. Cómo se verifica

> Una regla de accesibilidad sin forma de verificarla es una intención, no una decisión. El objetivo es que al menos parte de los errores se detecten automáticamente.

- **Linter estático** (ej: `eslint-plugin-jsx-a11y` para React, `axe-linter`) — *default; detecta errores obvios en tiempo de desarrollo sin costo de ejecución.*
- Testing automatizado con `axe-core` / `@axe-core/react` en suite de tests.
- Ambos — linter en desarrollo + axe en tests (mayor cobertura).
- Revisión manual solamente (lector de pantalla, teclado).
- Sin verificación formal por ahora.

## Tech a registrar

Si se elige un linter o librería de testing de a11y (ej: `eslint-plugin-jsx-a11y`, `axe-core`), registrarlos en `tech/`.

## Qué materializar

ADR `accessibility` materializado según el template `~/.claude/references/adr/templates/adr.md`.

- **Contexto:** por qué importa la accesibilidad en este proyecto (obligación legal aplicable, tipo de usuarios) y el nivel WCAG objetivo elegido (2.2 A | AA | AAA, o "sin objetivo formal por ahora").
- **Decisión:** nivel de conformidad objetivo y el mecanismo de verificación elegido (linter estático, testing con axe-core, ambos, revisión manual).
- **Reglas verificables:** reformulá la exigencia de a11y como aserciones chequeables, cada una con su mecanismo según lo elegido. Ejemplos:
  - **[tool: eslint-plugin-jsx-a11y]** todo elemento interactivo expone nombre accesible (`aria-label` o texto visible); el lint falla el build si falta.
  - **[tool: axe-core]** la suite de tests no reporta violaciones de nivel ≥ al objetivo (A/AA) en las vistas cubiertas.
  - **[manual]** navegación completa por teclado y verificación con lector de pantalla en los flujos críticos.

  Si se eligió "sin verificación formal por ahora", marcar la regla mínima como `[manual]` y materializar el ADR con `Status: Pending` indicando el trigger esperado.
