# lint-staged

- **Category:** Other
- **Version:** sin pinear
- **Status:** Accepted
- **Decided in phase:** delivery
- **Date:** 2026-07-27

## Por qué la elegimos

Corre los checks solo sobre los archivos staged. Con un único componente en el repo, el bloque `lint-staged` vive entero en el `package.json` root: los `.go` staged corren `gofumpt` (`go run mvdan.cc/gofumpt@v0.11.0 -w`), invocado por `.husky/pre-commit` (`pnpm exec lint-staged`). Un solo hook, un solo bloque de config.

## Alternativas descartadas

- Script propio de ruteo por path en el hook: reimplementa lo que lint-staged ya hace nativo.
- `golangci-lint` en pre-commit: no es apto por-archivo (opera por paquete); queda como gate de PR.

## Notas

Usada en: `.husky/pre-commit`. Config: `package.json` root (`*.go`). Registrada acá (root) porque es transversal.
