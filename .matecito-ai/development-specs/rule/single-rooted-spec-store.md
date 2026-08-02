# Capability — El store de specs nunca se parte por app

- **Status:** Accepted
- **Date:** 2026-08-02
- **Components:** cli

## Propósito

Garantizar que `.matecito-ai/development-specs/` vive **siempre en el root del repositorio**, no en carpetas de aplicaciones individuales (`apps/api`, `apps/ui`, etc.), porque el comportamiento del sistema **no respeta el límite de una sola surface** y partir el store crearía duplicación o arbitrariedad en su dueño.

## Reglas de negocio

- El store es **proyecto-level**, no app-level o surface-level — incluso en un repositorio con múltiples superficies o aplicaciones (`apps/api`, `apps/ui`, etc.). Si el autor de specs (bootstrap) pide crear el store dentro de una app (e.g., `apps/api/development-specs/`), el sistema lo materializa igual en el root del repositorio (`.matecito-ai/development-specs/`), explica el porqué y ofrece el eje `components` como alternativa para discriminar qué superficies participan en cada spec.
- El comportamiento es **intencionado y multisuperficie** — la regla de un flujo aplica a todas las apps que la usan, o a ninguna. Partir el store no resuelve eso; lo esconde.
- **La asimetría con los EDRs es explícita:** las decisiones de ingeniería (`../edr/`) son **cómo está construida una pieza**, univaluada por sub-app (e.g., `apps/api/.matecito-ai/edr/`). Los specs son **qué hace el sistema**, multivaluado — la misma regla puede gobernar comportamiento en la API, la UI y el CLI. Un contrato duplicado diverge y deja de ser contrato.
- **Store fuera del root ya existe (CRITICAL):** el validador de specs lo detecta y lo reporta como CRITICAL. Se ofrece consolidar (mover los specs al root), pero **nunca es automático** — mover specs es destructivo y puede haber contenido divergente; el usuario decide si mover, mergear o aceptar el riesgo.
- El consumidor (Claude) siempre consulta un único store centralizado.

## Escenarios

### Scenario: se pide un store por app

- **GIVEN** un pedido de crear el store dentro de una app (e.g., `apps/api/development-specs/`)
- **WHEN** bootstrap lo procesa
- **THEN** materializa en el root, explica el porqué y ofrece el eje `components` para discriminar por superficie

### Scenario: ya existe un store fuera del root (CRITICAL)

- **GIVEN** un repo con `development-specs/` detectado fuera del root (e.g., en `apps/api/`)
- **WHEN** corre el pre-flight de bootstrap o la validación
- **THEN** se reporta como CRITICAL y se ofrece consolidar sin hacer nada automático — el usuario decide si mover, mergear o aceptar el riesgo

## Referencias

- **Rule** → [`./repo-component-declaration.md`](./repo-component-declaration.md) — la alternativa para discriminar qué componentes participan en cada spec.
- **Referencia** → `~/.claude/references/spec/README.md` — el concepto de capability-spec como fuente de verdad del comportamiento del sistema.
