# Roadmap — Matecito UI

Estado: la base está cerrada (cambio SDD `ui-base-components`, archivado): catálogo tech instalado, tokens light, 13 primitives + sonner + `cn()` en `src/shared/`, gates ESLint/Prettier/tsc/husky operativos. Lo que sigue, en orden:

## 1. Componentes de dominio

Construir sobre los primitives, siguiendo la referencia visual (`docs/visual/`):

- [x] AgentNode
- [x] ArtifactCard
- [x] DecisionCard
- [x] DiffView
- [x] SeverityTag
- [x] CanvasEdge
- [x] LabeledEdge
- [x] TimelineScrubber

Van bajo `src/features/` según `structure/architecture-style.md`. Sin pantallas todavía.

## 2. Pantallas compuestas

- [ ] Shell del cockpit (layout, rail, header) — referencia `Matecito Cockpit.dc.html`
- [ ] Pantalla de componentes / galería — referencia `Matecito Components.dc.html`
- [ ] Routing con TanStack Router (ya instalado)

## 3. Follow-ups diferidos

Anotados en el archive-report de `ui-base-components`; ninguno bloquea:

- [ ] Tema dark + toggle (`frontend/styling.md` queda parcialmente implementado: hoy solo `:root` light)
- [ ] Primitive Command (⌘K) — `cmdk` ya está instalado, falta el componente
- [ ] Enforcement ESLint de la regla closed-set/named-constant (`structure/code-conventions.md`, hoy chequeo manual)
- [ ] Revisar contraste de `--border`/`--input` (~1.2:1 sobre card blanca, sugerencia de verify)
- [ ] Gate de sync Kubb/OpenAPI — espera a que exista el spec OpenAPI del broker (`apps/api`)
