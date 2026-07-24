# EDR — Estructura de carpetas

- **Status:** Accepted
- **Type:** convention
- **Date:** 2026-07-23

## Contexto

Consecuencia directa de la organización por vertical slice: hace falta un layout que agrupe cada feature con todo lo suyo y separe lo transversal.

## Decisión

Layout raíz `src/{app, routes, features, shared}` más el punto de entrada. `routes` contiene las rutas file-based (finas, importan de `features`; el enrutador se decide en routing). Cada feature contiene lo suyo: sus componentes, sus hooks, su store y sus queries. Lo transversal vive en `shared`:
- component base (shadcn) como capa de UI compartida,
- el cliente WS y el setup de TanStack Query como borde de datos compartido,
- utilidades y hooks genéricos.

**Sufijos de rol** (recomendados; el casing lo fija code-conventions): componentes en `PascalCase.tsx`, hooks `use*.ts`, store `*.store.ts`, queries `*.queries.ts`, tipos `*.types.ts`.

## Alcance

- `apps/ui/src/app/**` — composición raíz de la app.
- `apps/ui/src/routes/**` — rutas file-based finas que importan de `features/`.
- `apps/ui/src/features/**` — un directorio por feature con sus componentes, hooks, store y queries.
- `apps/ui/src/shared/**` — component base, borde de datos y utilidades transversales.

## Reglas verificables

- **[manual]** cada feature vive en su propia carpeta bajo `features/`.
- **[manual]** lo transversal vive bajo `shared/`.
- **[tool: eslint]** cada tipo de artefacto lleva su sufijo de rol.

## Relacionados

- `relacionado-con` → [architecture-style.md](architecture-style.md) — el layout deriva del estilo vertical slice.
- `relacionado-con` → [code-conventions.md](code-conventions.md) — el casing y las convenciones de nombre.
