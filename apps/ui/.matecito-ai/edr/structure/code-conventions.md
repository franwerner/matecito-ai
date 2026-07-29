# EDR — Convenciones de código

- **Status:** Accepted
- **Date:** 2026-07-23

## Contexto

Proyecto TS/React single-dev; hace falta un set idiomático cerrado que mantenga el código consistente y verificable por herramientas, con una regla propia del usuario sobre conjuntos cerrados.

## Decisión

Set idiomático TS/React:
- conjuntos cerrados con **union de literales** por sobre enum;
- ausencia con `T | null` / `undefined` bajo strict;
- inmutable por default (`const` / `readonly`);
- prohibir `any` (tipos explícitos en los bordes, inferencia adentro);
- literales mágicos siempre a constante nombrada;
- casing: componentes y tipos en PascalCase, funciones y variables en camelCase, constantes en UPPER_SNAKE, booleanos con prefijo `is/has/can`;
- guard clauses / early return, ≤ 3-4 parámetros o se pasa un objeto de opciones, evitar parámetros booleanos;
- comparación estricta (`===`) siempre;
- iteración declarativa (`map` / `filter` / `reduce`);
- imports ordenados (simple-import-sort), alias absoluto `@/`, sin deep-imports a internos de otra feature (se importa por su índice público).

**Regla especial del usuario:** los valores de un conjunto cerrado se definen una sola vez como constante nombrada (objeto `as const` o enum) y se comparan y usan siempre por esa referencia, **nunca** por el string literal inline.

Enforcer: ESLint (flat config) con typescript-eslint, eslint-plugin-react-hooks y eslint-plugin-jsx-a11y, más Prettier.

## Reglas verificables

- **[auto]** union sobre enum, no-magic-numbers, eqeqeq, orden de imports y casing.
- **[auto]** sin `any`; ausencia modelada con `T | null`.
- **[auto] + [manual]** los valores de un conjunto cerrado se usan por su constante nombrada, nunca por string literal inline.
- **[auto] + [manual]** forma de función: pocos parámetros o objeto de opciones, sin parámetros booleanos.

## Relacionados

- `relacionado-con` → [folder-structure.md](folder-structure.md) — los sufijos de rol y el casing.
- `relacionado-con` → [../runtime/error-handling.md](../runtime/error-handling.md) — la validación como valores y el manejo de lo inesperado.
