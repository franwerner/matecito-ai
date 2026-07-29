# EDR — Gates de calidad y CI

- **Status:** Accepted
- **Date:** 2026-07-23

## Contexto

No hay suite de tests (read-first, single-dev), así que los gates que bloquean merge se apoyan en lint, formato y chequeo de tipos, más el sync de los tipos generados desde el OpenAPI del broker.

## Decisión

Gates que bloquean merge: **ESLint + Prettier + `tsc --noEmit`** (sin tests, no hay suite). Pre-commit con **husky + lint-staged a nivel root del monorepo** (los mismos checks sobre lo staged de la UI; el hook vive en el root del repo y rutea por config más cercana — la config de la UI sigue en su `package.json`). CI en **GitHub Actions** sobre PR. Gate extra: los **tipos generados por Kubb deben estar en sync con el OpenAPI** del broker (falla si quedaron desactualizados).

## Reglas verificables

- **[auto]** el merge se bloquea si el linter reporta errores.
- **[auto]** se bloquea si hay diff de formato.
- **[auto]** se bloquea ante errores de tipos.
- **[auto]** los mismos checks corren en local antes del commit.
- **[auto]** el PR falla si los tipos generados por Kubb quedaron desincronizados del OpenAPI.

## Relacionados

- `relacionado-con` → [../security/dependency-scanning.md](../security/dependency-scanning.md) — `pnpm audit` como gate cuando exista CI.
