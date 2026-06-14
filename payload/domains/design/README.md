# Dominio: Design

Asistencia al diseño visual full-spectrum: **marca + UI/UX + prototipos + guías de marca**, guiada por **SDD (Spec-Driven Design)**. Este dominio **no toca código** — el handoff diseño → código pertenece al dominio `development`.

Es un plugin sobre el kernel agnóstico (`../../core/`). Acá viven su vocabulario, sus fases, sus agentes, sus skills y sus guards. El contrato general de área está en [`../README.md`](../README.md).

## Vocabulario

El kernel define slots genéricos; el dominio los ata a términos concretos. Es el espejo de development sobre un sustrato visual.

| Slot del kernel | Binding en design |
| --- | --- |
| Flow | **SDD** (Spec-Driven **Design**) |
| Artefacto de alineación | `brief` (el "spec" de design) |
| Decision record | **DDR** (Design Decision Record), en `.matecito-ai/ddr/` |
| Catálogo canónico | `design-principles` (`~/.claude/references/design-principles/`) |
| Exploration | **Figma** (MCP `figma`), activo cuando hay un archivo conectado |
| Guards | `visual-accessibility` + `brand-consistency` |
| Workspace | carpeta |
| Agentes de fase | `design-*` (`design-intake`, `design-explore`, …, `design-archive`) |

## Pipeline de fases

Nueve fases:

```
intake → explore → propose → brief → system → tasks → produce → verify → archive
```

Espejo del pipeline de development: `brief` es el `spec`; `system` es la fase `design` (fija el sistema visual —paleta, escala tipográfica, grilla, componentes— y lee/escribe DDRs); `produce` es el `apply`.

No todo trabajo recorre las nueve. El flujo es una **base inmutable** más **add-ons opcionales**:

- **Base (siempre corre):** `intake → brief → produce → verify → archive`.
- **Add-ons:** `explore`, `propose`, `system`, `tasks`.

`intake` es la fase de entrada y produce el brief; el orquestador **muestra direcciones y espera que elijas** (gate humano) antes de seguir. Un flyer rápido va `reduced` (solo base); un rebrand va `full`.

| Fase | Lee | Escribe |
| --- | --- | --- |
| `design-intake` | pedido crudo | `intake` (brief-intake) |
| `design-explore` | intake | `explore` |
| `design-propose` | exploration (opcional) | `proposal` (direcciones) |
| `design-brief` | proposal (requerido) | `brief` |
| `design-system` | brief + **DDRs** | `system` |
| `design-tasks` | brief + system | `tasks` |
| `design-produce` | tasks + brief + system + produce-progress | `produce-progress` |
| `design-verify` | brief + system + produce-progress + **DDRs tocados** | `verify-report` |
| `design-archive` | todos los artefactos | `archive-report` |

## Agentes

`design-*` — un sub-agente por fase, con contexto propio. Son el espejo de los `sdd-*` de development:

```
design-intake  design-explore  design-propose  design-brief  design-system
design-tasks  design-produce  design-verify  design-archive
```

## Skills

La diferencia clave: una skill **de fase** es la *receta* de esa fase del pipeline; una skill **de capacidad** es una *técnica* reutilizable, no atada a una fase.

**Por fase** (`design-phases/design-*`): la receta de cada fase.

```
design-intake  design-explore  design-propose  design-brief  design-system
design-tasks  design-produce  design-verify  design-archive
```

**De capacidad** (`design-core/*`): técnicas que las fases invocan o que se usan ad hoc.

```
brand-from-references  # derivar marca desde referencias visuales
brand-guide            # producir/actualizar la guía de marca
explore-variations     # generar variantes/direcciones de una pieza
generate-assets        # producir los assets finales
consistency-audit      # auditar consistencia de un set de piezas
visual-accessibility   # chequeo de contraste/tamaños/jerarquía (WCAG)
figma-hygiene          # higiene del archivo Figma (capas, nombres, estructura)
explain-concept        # el motor mentor: explica el porqué de un concepto
design-review          # revisión crítica de una pieza
```

**De decisiones** (`matecito-ai/*`): el ciclo de vida de los DDRs.

```
design-decisions-bootstrap  # captura interactiva de decisiones de diseño → DDRs; modo update ratifica Inferred→Accepted
design-decisions-mine       # mina decisiones desde el archivo Figma (styles/components) → DDRs Inferred
design-decisions-validate   # auditor consultivo: coherencia, completitud y verificabilidad de los DDRs
```

**Setup y onboarding** (meta, no son fases del pipeline):

```
design-init      # setup inicial: detecta contexto de diseño (Figma/Canva/marca) y capabilities, bootstrapea persistencia
design-onboard   # recorrido guiado del flujo de diseño de punta a punta, enseñando por hacer
```

## Guards

Gates de verificación que corren en `design-verify`.

- **visual-accessibility** — ratios de contraste WCAG, tamaños mínimos y jerarquía, chequeados contra los colores y la tipografía reales del archivo Figma. Flaggea lo que esté por debajo de AA.
- **brand-consistency** — cada pieza producida se chequea contra la guía de marca y los DDRs aceptados. Flaggea cualquier pieza que contradiga un decision record.

## Modo mentor

Regla transversal a todas las fases y skills: explicar el **porqué** detrás de cada decisión o hallazgo en 1-2 líneas — el principio de diseño subyacente, no solo el qué. Es el motor de aprendizaje: el equipo se vuelve más eficiente y aprende en el camino. Cuando aparece un concepto que la persona puede no conocer, deriva a la skill `explain-concept`. Cita el catálogo canónico `design-principles` en vez de improvisar la justificación.

## MCP

- **Engram (`engram`)** — memoria persistente (mecanismo del núcleo): artefactos del flujo entre fases + descubrimientos entre sesiones. design lo declara como dependencia propia (`mcp` + `binaries`), igual que cualquier otra.
- **Figma (`figma`)** — lo registra `install` (read-only): deja al agente **leer** el archivo Figma (revisar, auditar, extraer marca). OAuth una vez por persona vía `/mcp`.
- **Canva (`canva`)** — lo registra `install` (`claude mcp add --transport http canva https://mcp.canva.com/mcp`), el MCP oficial hosteado de Canva: deja al agente crear/editar piezas on-brand. OAuth una vez por persona vía `/mcp`. No usar el `@canva/cli ... mcp` de los tutoriales — ese es para construir apps de Canva, no para diseñar.

> **Dependencias declaradas (manifest).** `mcp: [engram, figma, canva]` · `binaries: [engram]`. Nada se instala global: el ecosistema instala esto solo cuando design está activo, y deriva de `mcp` los permisos de Claude Code (`mcp__<name>__*`). design **no** instala codegraph, drawio, context7 ni proofshot.

## Config del dominio

Lo que aparece en la pantalla de configuración del dominio (resuelto por-proyecto → global → default):

- **Models per agent** (`models`) — qué modelo usa cada fase (`design-intake`, `design-brief`, `design-system`, …). Sin valor configurado, cada agente usa su default curado.
- **Auto-mine DDR** (`flagDecisionGaps`, default `false`) — opt-in. Con el flag on: `design-tasks` marca cada decisión de marca como `· ddr: <surface>/<slug>` (exista o no el DDR), `design-verify` confirma cuáles se implementaron (sección `## Decision Gaps`), y al cerrar el orquestador dispara `design-decisions-mine` —que lee el archivo Figma vivo (styles/components)— y ofrece materializar las decisiones como DDRs `Inferred` (siempre con tu confirmación). Ratificás `Inferred → Accepted` con `design-decisions-bootstrap` en modo update. Aparte, `design-verify` compara los DDRs `Accepted` contra el estado real de Figma (drift): cualquier divergencia es un `DDR-VIOLATION`. Canva queda fuera (sin tokens legibles).

## Ver también

- [README raíz del ecosistema](../../../README.md) — visión general de matecito-ai.
- [Contrato de área](../README.md) — cómo se estructura un dominio.
