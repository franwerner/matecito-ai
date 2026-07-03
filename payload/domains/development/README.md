# Dominio: Development

Asistencia al desarrollo de software guiada por **SDD (Spec-Driven Development)**, sobre un repositorio de código. Es el dominio que lleva cada cambio desde un pedido en lenguaje natural hasta el código, capturando las decisiones de arquitectura una vez y respetándolas a lo largo del tiempo.

Este dominio es un plugin sobre el kernel agnóstico (`../../core/`). Acá viven su vocabulario concreto, sus fases, sus skills y sus guards. El contrato general de área está en [`../README.md`](../README.md).

## Vocabulario

El kernel define slots genéricos; el dominio los ata a términos concretos.

| Slot del kernel | Binding en development |
| --- | --- |
| Flow | **SDD** (Spec-Driven Development) |
| Artefacto de alineación | `spec` |
| Decision record | **ADR**, en `.matecito-ai/adr/` |
| Comportamiento durable | **capability-spec**, en `.matecito-ai/development-specs/` |
| Catálogo canónico | `design-patterns` (`~/.claude/references/design-patterns/`) |
| Exploration | **codegraph** (`mcp__codegraph__*`), activo cuando existe `.codegraph/` |
| Guards | `strict-tdd` + `review-workload` |
| Workspace | repositorio |
| Agentes de fase | `sdd-*` (`sdd-intake`, `sdd-explore`, …, `sdd-archive`) |

## Pipeline de fases

Nueve fases:

```
intake → explore → propose → spec → design → tasks → apply → verify → archive
```

No todo cambio recorre las nueve. El flujo es una **base inmutable** más **add-ons opcionales**:

- **Base (siempre corre):** `intake → spec → apply → verify → archive`.
- **Add-ons (se activan según el tamaño del cambio):** `explore`, `propose`, `design`, `tasks`.

`intake` es la fase de entrada: hace 2-4 preguntas para estructurar el pedido, lo clasifica y produce un brief. El orquestador **siempre muestra ese brief y espera tu confirmación** (gate humano) antes de seguir. Un fix trivial va directo; un cambio grande activa todos los add-ons.

| Fase | Lee | Escribe |
| --- | --- | --- |
| `sdd-intake` | pedido crudo | `intake` |
| `sdd-explore` | intake (brief) | `explore` |
| `sdd-propose` | exploration (opcional) | `proposal` |
| `sdd-spec` | proposal (requerido) + **capability-spec durable** (para Modified Capabilities) | `spec` |
| `sdd-design` | proposal + **ADRs** + **capability-specs durables** (requerido) | `design` |
| `sdd-tasks` | spec + design + **capability-specs tocados** (requerido) | `tasks` |
| `sdd-apply` | tasks + spec + design + apply-progress | `apply-progress` |
| `sdd-verify` | spec + tasks + apply-progress + **ADRs tocados** + **capability-specs tocados** | `verify-report` |
| `sdd-archive` | todos los artefactos | `archive-report` + **merge en capability-specs** |

En lanes `reduced`/`custom` algunas fases upstream no corren; cada fase lee el upstream disponible más cercano (`sdd-spec` cae al brief de intake cuando no hay proposal; `sdd-apply` toma `spec` como piso y saltea `tasks`/`design` si faltan). Los **capability-specs durables** se leen solo si existe `.matecito-ai/development-specs/` (gate por presencia, igual que los ADRs).

## Capability-specs — el comportamiento durable

El dominio mantiene **dos capas de conocimiento durable** en `.matecito-ai/`, versionadas en git y **nunca** en Engram:

- **ADRs** (`adr/<dominio>/`) — *qué se eligió y por qué*: la decisión técnica (tecnología, patrón, política) y su justificación.
- **capability-specs** (`development-specs/<tipo>/`) — *qué hace* el sistema: el comportamiento observable de cada capacidad, con escenarios verificables.

Tres capas, no dos: **capability-spec** (qué hace) · **ADR** (qué se eligió y por qué) · **código** (cómo se implementa). El spec habla en **idioma de dominio + contrato público**; nunca nombra identificadores internos (eso es del código) ni justifica elecciones (eso es del ADR).

**Organización por tipo** (la ruta lleva el tipo, como el dominio en un ADR): `flow` (operación de cara a un actor, con pasos/ramas), `rule` (regla de negocio transversal), `lifecycle` (máquina de estados de una entidad), `process` (reactivo/de fondo, por evento). Dos niveles de índice: raíz que enruta por tipo + un `INDEX.md` por tipo. Concepto y plantillas en `~/.claude/references/spec/`.

**Dos "spec" que no hay que confundir:**
- el artefacto de fase `spec` (efímero, en Engram, `sdd/{change}/spec`) es el **delta** de un cambio (qué AGREGA / MODIFICA / QUITA);
- el **capability-spec** durable es el **estado acumulado** del comportamiento. Al archivar, `sdd-archive` **mergea** el delta en el capability-spec (anclado en escenarios, no destructivo); `sdd-spec` lee el durable para construir el delta de una capability existente.

**Skills:** `development-spec-bootstrap` define el comportamiento *upfront* (entrevista por capability, por tipo) y asienta el pointer en el `CLAUDE.md` del proyecto; `development-spec-validate` chequea coherencia entre capabilities. **Gate por presencia:** si el proyecto no tiene `development-specs/`, todas las fases lo saltean en silencio (igual que con los ADRs).

## Componentes

Las piezas específicas de desarrollo del ecosistema:

| Capa | Pieza | Rol |
| --- | --- | --- |
| **Flujo** | Fork SDD | Fases intake → … → archive, con base inmutable + add-ons opcionales. Modelo por agente y Strict TDD configurables. |
| **Skill** | `development-decisions-bootstrap` | Entrevista por fases que captura decisiones de ingeniería y las materializa como ADRs por dominio. |
| **Skill** | `development-decisions-validate` | Validador consultivo: coherencia, completitud y verificabilidad de los ADRs. |
| **Skill** | `development-decisions-mine` | Mina decisiones desde el código de un repo existente y las propone como ADRs `Inferred` (borradores) para que un humano las ratifique vía bootstrap. |
| **Skill** | `development-spec-bootstrap` | Entrevista por capability que captura el comportamiento del sistema (qué hace) y lo materializa como capability-specs por tipo en `.matecito-ai/development-specs/`. Contraparte de decisions-bootstrap (qué-hace vs por-qué). |
| **Skill** | `development-spec-validate` | Validador consultivo: coherencia entre capabilities, completitud, verificabilidad y referencias de los capability-specs. |
| **Referencia** | `adr` | Definición canónica de qué es (y qué no es) un ADR + plantillas de estructura. Consultable y agnóstica de flujo. |
| **Referencia** | `spec` | Definición canónica de qué es (y qué no es) un capability-spec (comportamiento) + plantillas. Consultable y agnóstica de flujo. |
| **Catálogo** | `design-patterns` | Catálogo canónico de patrones de diseño. Los ADRs lo citan por nombre; `sdd-design` respeta la definición cuando un ADR declara `Patrón aplicado`. |
| **MCP** | `engram` | Memoria persistente (mecanismo del núcleo): artefactos del SDD entre fases + descubrimientos/fixes entre sesiones. |
| **MCP** | `context7` | Documentación de librerías al día (contra APIs alucinadas). Se engancha en `apply`. |
| **MCP** | `codegraph` | Grafo de código pre-indexado (tree-sitter + SQLite) para explorar por estructura. |
| **MCP** | `drawio` _(next-ai-draw-io)_ | Render de diagramas de arquitectura on-demand y **efímeros**: el thread principal renderiza en vivo el `<mxGraphModel>` en el paso de `design` (preview en la URL que reporta `start_session`; el puerto es dinámico). El **vocabulario** (formas, iconos, estilos, layout) lo aporta la skill `drawio`. No se exporta ningún archivo al repo. |
| **MCP** | `debugger` _(mcp-debugger)_ | Step-through DAP on-demand, nunca automático. Hogar principal: `sdd-apply` (diagnostica + corrige en el mismo contexto). En `sdd-verify`: solo diagnóstico, sin fixes. Se omite en silencio si el toolchain de debug del lenguaje no está disponible. |
| **CLI** | `proofshot` | Verificación visual de UI: graba el browser y valida los scenarios. `sdd-verify` la corre cuando el cambio toca UI y proofshot está disponible. |
| **Hook** | `git-commit-validator` | Hook `PreToolUse` sobre `git commit` (subcomando Go `matecito-ai hook git-commit-validate`): bloquea atribución a IA (`Co-Authored-By`, `Claude`, …), avisa si el formato no es Conventional Commits, y deja pasar lo que no puede inspeccionar. |
| **Agentes** | `sdd-*` | Un sub-agente por fase, con contexto propio. |

> **Dependencias declaradas (manifest).** `mcp: [engram, context7, codegraph, drawio, debugger]` · `binaries: [engram, codegraph, proofshot]`. Nada se instala global: el ecosistema instala esto solo cuando development está activo, y deriva de `mcp` los permisos de Claude Code (`mcp__<name>__*`).

> **Hooks (transparencia).** Cuando development está activo, el instalador **escribe el hook en tu `~/.claude/settings.json`** y este corre **automáticamente** (Claude Code no pide aprobación para hooks). El hook se identifica con una marca `matecitoId` y se reconcilia por identidad en cada `install`/`update` (reemplaza el suyo, nunca toca tus hooks). Su declaración e implementación viven compiladas dentro del binario `matecito-ai` (subcomando `hook git-commit-validate`) — no hay archivo de hook en el payload ni script alguno depositado en tu `~/.claude`.

## Skills

**Por fase** (`gentle-ai/sdd-*`): una skill por fase del pipeline.

```
sdd-explore  sdd-propose  sdd-spec  sdd-design  sdd-tasks  sdd-apply  sdd-verify  sdd-archive
sdd-init  sdd-onboard
```

**De capacidad** (`matecito-ai/*`): técnicas transversales, no atadas a una fase.

```
git                          # formato de commits (Conventional Commits), atomicidad, atribución
development-decisions-bootstrap  # captura interactiva de decisiones → ADRs
development-decisions-validate   # validación consultiva de ADRs
development-decisions-mine       # minería de decisiones desde el código → ADRs Inferred
development-spec-bootstrap    # captura interactiva de comportamiento → capability-specs por tipo
development-spec-validate     # validación consultiva de coherencia entre capability-specs
sdd-intake                   # estructura el pedido crudo y produce el brief de entrada
debugger                     # preflight por lenguaje, loop de debug DAP y install helper via mcp__debugger__*
```

## Guards

Gates de verificación que el flujo corre.

- **strict-tdd** — si está activo (opt-in), `apply` y `verify` siguen el ciclo test-first: el test se escribe antes que la implementación. El test runner sale de `sdd/{project}/testing-capabilities` en Engram.
- **review-workload** — presupuesto de PR. Después de `tasks` y antes de `apply`, inspecciona el `Review Workload Forecast`; si recomienda PRs encadenados o el presupuesto de ~400 líneas está en riesgo, frena y pregunta (PRs encadenados vs. `size:exception`).

## Config del dominio

Lo que aparece en la pantalla de configuración del dominio (resuelto por-proyecto → global → default):

- **Models per agent** (`models`) — qué modelo usa cada fase (`sdd-intake`, `sdd-spec`, `sdd-design`, …). Sin valor configurado, cada agente usa su default curado.
- **Strict TDD** (`strictTdd`, default `false`) — si está activo, `apply` y `verify` siguen test-first.
- **Auto-mine ADR** (`flagDecisionGaps`, default `false`) — opt-in. Activa la detección de decisiones implementadas sin ADR durante el flujo; al cerrar, ofrece minarlas como ADRs `Inferred` (siempre con tu confirmación).

## Ver también

- [Guía profunda del flujo SDD](../../../docs/guide/README.md) — cómo funciona todo de punta a punta: fases, herramientas y la capa de decisiones (bootstrap / validate / mine).
- [README raíz del ecosistema](../../../README.md) — visión general de matecito-ai.
- [Contrato de área](../README.md) — cómo se estructura un dominio.
