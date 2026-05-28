---
name: accessibility
depth: light
domain: frontend
tipo: decisión
adr-output: accessibility
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

ADR `accessibility` con: nivel WCAG objetivo, mecanismo de verificación elegido, y las reglas concretas (ej: "todo componente interactivo debe tener `aria-label` o texto visible; el linter falla el build si no").
