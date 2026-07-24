# EDR — Escaneo de dependencias

- **Status:** Accepted
- **Type:** policy
- **Date:** 2026-07-23

## Contexto

Proyecto npm que conviene mantener con las deps escaneadas por vulnerabilidades, con un gate que se endurezca cuando exista CI.

## Decisión

**Dependabot** (servicio de GitHub) para las deps npm: PRs automáticos de seguridad. **`pnpm audit`** como gate cuando exista CI. La política ante una vulnerabilidad crítica (bloquear o abrir PR) se fija junto con el CI.

## Reglas verificables

- **[tool: dependabot]** las deps npm tienen escaneo de vulnerabilidades con PRs automáticos.
- **[tool: pnpm audit]** corre como gate en CI cuando exista.

## Relacionados

- `relacionado-con` → [../delivery/ci-quality-gates.md](../delivery/ci-quality-gates.md) — el gate se endurece con el CI.
