# EDR — Accesibilidad

- **Status:** Accepted
- **Type:** decision
- **Date:** 2026-07-23

## Contexto

El cockpit es keyboard-first, con tema dual (default light) y flujos densos de datos. La accesibilidad debe sostenerse en ambos temas sin una suite de tests automatizados.

## Decisión

**WCAG 2.2 AA** en ambos temas. Verificación combinada:
- `eslint-plugin-jsx-a11y` en tiempo de desarrollo (falla el build si algo interactivo no expone nombre accesible),
- revisión manual: navegación por teclado, foco visible, lector de pantalla en los flujos críticos y validación de contraste sobre los pares reales texto/fondo de cada tema.

Los primitives de Radix (vía shadcn) aportan la base accesible: manejo de foco, ARIA y teclado. Sin axe en tests, porque no hay suite de tests.

## Reglas verificables

- **[tool: eslint-plugin-jsx-a11y]** todo elemento interactivo expone un nombre accesible; el lint falla el build si falta.
- **[manual]** el contraste AA se valida sobre los pares reales texto/fondo de cada tema.
- **[manual]** el color nunca es el único portador de significado: las severidades se acompañan de texto o ícono.
- **[manual]** navegación completa por teclado con foco visible en los flujos críticos.
