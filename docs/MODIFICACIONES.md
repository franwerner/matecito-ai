# matecito-ai — Modificaciones al SDD del Gentleman

> *Mientras la IA trabaja por vos, te tomás unos ricos mates.*

Este es el SDD del Gentleman (`payload/agents/` + `payload/skills/`) **modificado directamente** (fork) para integrarlo con el ecosistema matecito-ai: ADRs, codegraph y context7.

## Cómo está organizado

- **`vendor-original/`** — los archivos del Gentleman SIN tocar (snapshot de referencia). No editar. Sirven para hacer diff cuando salga un update upstream.
- **`payload/agents/` y `payload/skills/`** — la copia de trabajo, con las modificaciones de matecito-ai. Todo lo deployable a `~/.claude/` vive bajo `payload/`.

Todos los cambios están marcados con comentarios `<!-- matecito-ai: ... -->` (o `# matecito-ai:` en frontmatter), así son rastreables y reaplicables tras un update.

## Qué se modificó (4 fases, 8 archivos)

| Fase | Archivos | Cambio |
|------|----------|--------|
| **explore** | `payload/agents/sdd-explore.md`, `payload/skills/gentle-ai/sdd-explore/SKILL.md` | Step 3: política codegraph-first (estructura/relaciones) con grep como fallback (texto literal, archivos no indexados, o cuando codegraph no resuelve). Chequea `.codegraph/`. Tools codegraph agregadas al agente. |
| **design** | `payload/agents/sdd-design.md`, `payload/skills/gentle-ai/sdd-design/SKILL.md` | Step 2a: lee los ADRs vigentes de `.claude/adr/` antes de diseñar. Respeta los `Accepted`, frena (blocker) si los contradice, y flaggea decisiones nuevas no cubiertas para captura vía project-decisions-bootstrap. Secciones nuevas en el design doc: ADR Alignment / New Decisions / ADR Conflicts. |
| **apply** | `payload/agents/sdd-apply.md`, `payload/skills/gentle-ai/sdd-apply/SKILL.md` | Step 2 (5 y 6): respeta los ADRs aplicables como restricciones duras; usa context7 para docs de librerías y `codegraph_impact` antes de cambiar símbolos existentes. Tools codegraph + context7 agregadas al agente. |
| **verify** | `payload/agents/sdd-verify.md`, `payload/skills/gentle-ai/sdd-verify/SKILL.md` | Step 6b: verifica que el código del cambio respete los ADRs que tocó (acotado a este cambio, NO audita todo el catálogo). Violación → CRITICAL `ADR-VIOLATION`. |

## Qué NO se modificó (pero sigue presente)

propose, spec, tasks, archive, onboard. Operan sobre artefactos del SDD, no sobre código ni decisiones de arquitectura. (Sus 8 agentes sí tuvieron un retoque menor: el campo `skill_resolution` del envelope, por la remoción del registry — ver abajo.)

## Componentes ELIMINADOS (mecanismo de inyección)

matecito-ai usa **fork directo, no inyección de reglas**. Por eso se removió todo el mecanismo de skill-registry:

- **Borrados:** `skills/skill-creator/`, `skills/skill-registry/`, `skills/_shared/skill-resolver.md`.
- **Referencias colgantes reescritas para matecito-ai:**
  - `_shared/sdd-phase-common.md` — Sección A (carga de skills) reescrita: cada fase carga su propia `SKILL.md` y lee las convenciones del proyecto (`.claude/adr/`, `CLAUDE.md`, `config.yaml`). Sin bloque `Project Standards` inyectado. Campo `skill_resolution` simplificado a `phase-skill | none`.
  - `_shared/persistence-contract.md` — sección "Skill Registry" reemplazada por "Skill Loading".
  - `sdd-init` (SKILL + `references/init-details.md`) — ya no construye `.atl/skill-registry.md`: se quitaron el paso, las reglas de escaneo y el `mem_save` del registry.
  - Los 8 agentes — `skill_resolution: injected` → `phase-skill | none`.

Las convenciones del proyecto viven en los archivos del proyecto (ADRs, CLAUDE.md), no en un registry intermedio.

## Modos de persistencia ELIMINADOS (openspec / hybrid)

matecito-ai es **engram-only**. Se removieron los modos de persistencia basados en archivos (`openspec` e `hybrid`). Modos válidos ahora: **`engram | none`**.

- **Borrado:** `skills/_shared/openspec-convention.md`.
- **persistence-contract.md** reescrito a solo `engram | none` (tablas, roles y reglas simplificadas). Incluye la regla explícita "NEVER create `openspec/`".
- **Las 8 fases** (explore, propose, spec, design, tasks, apply, verify, archive) — se quitaron todas las ramas `IF mode is openspec/hybrid`, las creaciones de carpetas `openspec/`, las refs a `openspec/config.yaml` y `openspec/specs/`, y las líneas "Location/Archived to" que mencionaban rutas de archivo. Quedó solo la lógica engram/none.
- **sdd-archive** simplificado fuerte: ya no mueve carpetas ni mergea specs en filesystem; su trabajo en engram-only es escribir el reporte de archivo y marcar el estado.
- **sdd-init** — quitada la sección "OpenSpec Skeleton" y refs a config.yaml de openspec.
- **engram-convention.md, strict-tdd.md, agentes** — refs sueltas a openspec/hybrid limpiadas.

Esto es máxima divergencia del upstream (el SDD original es multi-modo). Mantener ante updates será más trabajo; los bloques propios están marcados `matecito-ai`.

## Pendiente de verificar antes de usar

1. **Nombres de las tools MCP.** Agregué las tools asumiendo los prefijos `mcp__codegraph__*` y `mcp__context7__*`. **Confirmá que coincidan** con cómo están registrados tus MCP servers (mirá tu `~/.claude.json` o `claude mcp list`). Si difieren, ajustá los nombres en el frontmatter de `payload/agents/sdd-explore.md` y `payload/agents/sdd-apply.md`. Están marcados con `# matecito-ai: ... VERIFY`.
2. **context7: nombres reales de sus tools.** Usé `resolve_library_id` y `query` como placeholders típicos; verificá los reales de tu instalación.

## Mantenimiento ante updates del Gentleman

1. Traé la versión nueva del SDD a una carpeta aparte.
2. `diff vendor-original/ <nueva>/` → mirá qué cambió el Gentleman.
3. Si un archivo que modificaste cambió upstream, reaplicá tus bloques `matecito-ai` (son fáciles de ubicar por el marcador).
4. Actualizá `vendor-original/` al nuevo snapshot.

## El ecosistema matecito-ai (contexto)

```
SKILLS    project-decisions (ADRs) + issue-brief + SDD (este fork)
MCP       codegraph + context7
AGENTES   sub-agentes del SDD (este fork)
ENGRAM    memoria de sesión (standalone)
```

- **ADRs** (`project-decisions`) → decisiones de arquitectura: qué/por qué, verificable, por dominio.
- **Engram** → memoria de sesión: descubrimientos, contexto, fixes. (Sin solapar con ADRs.)
- **codegraph** → explorar estructura del código (eficiente en tokens y tool calls).
- **context7** → docs de librerías al día.
