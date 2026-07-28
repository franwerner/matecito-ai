# husky

- **Category:** Other
- **Version:** sin pinear
- **Status:** Superseded
- **Decided in phase:** delivery
- **Date:** 2026-07-23

**Reemplazado por:** [`../../../../../.matecito-ai/edr/tech/husky.md`](../../../../../.matecito-ai/edr/tech/husky.md) — husky pasó a nivel root del monorepo (los git hooks son por repo; el montaje anidado en la UI requería un wrapper).

## Por qué la elegimos

Hooks de git para correr los mismos checks de calidad en pre-commit, en conjunto con lint-staged.

## Alternativas descartadas

- Ninguna evaluada formalmente.

## Notas

Usada en: delivery/ci-quality-gates (pre-commit).
