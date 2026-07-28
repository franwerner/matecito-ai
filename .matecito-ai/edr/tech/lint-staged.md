# lint-staged

- **Category:** Other
- **Version:** sin pinear
- **Status:** Accepted
- **Decided in phase:** delivery
- **Date:** 2026-07-27

## Por qué la elegimos

Corre los checks solo sobre los archivos staged, y en monorepo **rutea por la config más cercana a cada archivo**: lo staged bajo `apps/ui` usa el bloque lint-staged de la UI (eslint + prettier), los `.go` usan el bloque del `package.json` root (`gofmt -w`). Un solo hook, checks por sub-app.

## Alternativas descartadas

- Script propio de ruteo por path en el hook: reimplementa lo que lint-staged ya hace nativo.
- `golangci-lint` en pre-commit: no es apto por-archivo (opera por paquete); queda como gate de PR.

## Notas

Usada en: `.husky/pre-commit`. Configs: `package.json` root (`*.go`) y `apps/ui/package.json` (TS/CSS/JSON/MD). Registrada acá (root) porque es transversal.
