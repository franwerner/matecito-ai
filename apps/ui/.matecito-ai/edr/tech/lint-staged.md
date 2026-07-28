# lint-staged

- **Category:** Other
- **Version:** sin pinear
- **Status:** Superseded
- **Decided in phase:** delivery
- **Date:** 2026-07-23

**Reemplazado por:** [`../../../../../.matecito-ai/edr/tech/lint-staged.md`](../../../../../.matecito-ai/edr/tech/lint-staged.md) — lint-staged pasó a nivel root del monorepo; la config de la UI en su `package.json` sigue viva, ruteada por la config más cercana.

## Por qué la elegimos

Corre los checks (ESLint, Prettier, tipos) solo sobre los archivos staged en el pre-commit, junto con husky.

## Alternativas descartadas

- Ninguna evaluada formalmente.

## Notas

Usada en: delivery/ci-quality-gates (pre-commit sobre staged).
