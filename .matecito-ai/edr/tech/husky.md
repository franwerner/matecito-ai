# husky

- **Category:** Other
- **Version:** sin pinear
- **Status:** Accepted
- **Decided in phase:** delivery
- **Date:** 2026-07-27

## Por qué la elegimos

Hooks de git **a nivel root del monorepo**: un solo `pre-commit` intercepta todos los commits del repo (los git hooks son por repo, no por sub-app) y delega en lint-staged. Reemplaza el montaje anterior anidado en la UI, que requería un wrapper para resolver el git root.

## Alternativas descartadas

- Husky anidado en `apps/ui` (montaje anterior): funcionaba repo-wide igual, pero con un wrapper frágil y sin lugar natural para los checks de Go.
- Hooks de git a mano (`core.hooksPath` manual): sin gestión de instalación por `prepare`.

## Notas

Usada en: el pre-commit del repo (`.husky/pre-commit` → lint-staged). Instalación: `pnpm install` en el root (script `prepare`). Registrada acá (root) porque es transversal a todos los sub-apps.
