---
name: development-decisions-mine
description: Minería de decisiones de ingeniería desde el código fuente de un repo existente. Produce candidatos a ADR con Status Inferred a partir de evidencia estructural, de configuración, de patrones o de ausencia. Dos modos — Mode A (scan brownfield invocado por la skill) y Mode B (in-flow, opt-in via flagDecisionGaps). Usá esta skill cuando el usuario pida "minear decisiones", "encontrar ADRs que faltan", "¿qué decisiones hay implícitas en el código?", o cuando el orchestrator dispare el boundary dispatch post-verify con gaps implementados.
---

# Development Decisions Mine

Minería consultativa de decisiones de ingeniería implícitas en el código: escanea el repo, construye candidatos a ADR con evidencia observable, y los presenta para confirmación humana antes de escribir cualquier cosa. Nunca escribe ADRs sin confirmación explícita.

> **Arquitectura de esta skill — motor/datos igual que bootstrap.** `SKILL.md` (este archivo) es el motor: pipeline, modelo de evidencia, reglas de confianza, flujos Mode A y Mode B, invariantes. El catálogo de dominios y la taxonomía de concerns son de `development-decisions-bootstrap` (READ-ONLY); las plantillas de estructura de ADR (incl. `adr.md`) viven en la referencia `~/.claude/references/adr/templates/` (READ-ONLY). La salida son ADRs con `Status: Inferred` en `.matecito-ai/adr/<dominio>/<slug>.md`, misma taxonomía.

> **Concepto de ADR — referencia canónica.** Qué ES y qué NO ES un ADR, y la diferencia borrador(inferido)/aceptado, están en `~/.claude/references/adr/README.md` (referencia agnóstica de flujo/skill). mine NO redefine el concepto — lo aplica. Si dudás si un hallazgo "es un ADR", esa referencia es la regla.

---

## Arquitectura de esta skill

La skill está partida en **motor** y **executor**:

- **`SKILL.md` (este archivo) = el motor.** Define el pipeline, el modelo de evidencia, las reglas de confianza (router + label), los flujos de Mode A y Mode B, y los invariantes. Es el contrato que el executor debe seguir.
- **`payload/agents/development-decisions-mine.md` = el executor.** Agente de contexto fresco que hace el trabajo pesado de scan/discovery y RETORNA un bloque `candidates[]`. No escribe ADRs. El executor sirve tanto Mode A como Mode B.
- **El thread principal = gate + materialize.** Recibe `candidates[]` del executor, renderiza la tabla de confirmación, y ejecuta la materialización solo después del confirm explícito del usuario.

---

## Invariantes (no negociables)

Estas dos invariantes son la columna vertebral de la skill. Cualquier desviación las rompe.

**Invariante 1 — ADRs nunca obligatorios; el ejecutor recibe un scope, no decide nada.** El ejecutor NO lee config, NO resuelve el flag, NO se ramifica por modo: recibe un **scope** de su caller y lo procesa. La intención (si correr) vive en el caller:

- invocación directa de la skill → scope = repo completo.
- el orquestador (con `flagDecisionGaps` resuelto en true + gaps implementados) → scope = la gap list.

Con flag off el orquestador no invoca → silencio total. La dependencia dura del ejecutor es el **catálogo de concerns** (siempre presente), de donde sale la taxonomía de dominios — NO los ADR generados. La ausencia de `.matecito-ai/adr/` NO bloquea: "nada decidido todavía" → cada candidato es un hueco y mine puede bootstrapear (crear la carpeta + los primeros ADR) tras el confirm. La presencia de un ADR se chequea **por-candidato** (dedup), nunca como gate global.

**Invariante 2 — Nunca auto-materializar.** El executor NO tiene capacidad de escritura de ADRs en la fase discover. La materialización es un paso separado, alcanzable solo después de un confirm explícito del usuario en el gate del thread principal.

---

## Pipeline runtime: `discover → confirm → materialize`

El pipeline corre en tres contextos distintos por diseño:

```
[executor — contexto fresco]          [thread principal]
discover → draft candidates[]   →→→   confirm (gate) → materialize
```

### Paso 1: Discover (executor)

El executor escanea el repo, construye candidatos con evidencia, y RETORNA un bloque `candidates[]` estructurado. No escribe nada.

### Paso 2: Confirm (thread principal — gate)

El thread principal recibe los candidatos y renderiza la tabla de confirmación:

```
Candidatos ordenados por confidence label (high primero):
┌─────────────────────────────────────────────────────────────┐
│ [dominio/slug]  kind · observado · prevalencia · confidence │
│ ✅ alto-signal-candidate  →  [aceptar / editar / saltar]    │
│ ⬜ otro-candidate         →  [aceptar / editar / saltar]    │
└─────────────────────────────────────────────────────────────┘

Preguntas abiertas (→ bootstrap, no Inferred):
- [pregunta 1]: ¿es esto una decisión?
- ...

Posibles gaps del catálogo (advisory):
- concern X no encaja en ningún dominio canónico
```

Opciones de confirm: `accept-all` / por ítem / `none`.

### Paso 3: Materialize (thread principal — solo post-confirm)

Para cada candidato aceptado/editado, el thread principal escribe el ADR con `Status: Inferred` usando `~/.claude/references/adr/templates/adr.md` (READ-ONLY). Ver "Materialización" más abajo.

---

## Modelo de evidencia — keyed by `kind`

El `kind` de cada candidato determina qué se lee durante el scan y qué secciones del ADR se llenan.

| kind | qué se lee | llena `## Evidencia (inferida)` | llena `## Alcance` |
|---|---|---|---|
| `estructural` | codegraph (symbols/edges) ▸ grep | kind + observado + prevalencia | globs (el locator estructural) |
| `patrón` | codegraph (recurring shape) ▸ grep | kind + observado + prevalencia | globs |
| `config` | manifest entry (package.json / pyproject.toml / go.mod / CI yaml) | kind + observado (sin prevalencia) | omitido |
| `ausencia` | grep probando ausencia en sitios esperados | kind + observado (sin prevalencia) | omitido (no hay glob) |

Reglas de llenado:

- `## Evidencia (inferida)` contiene solo la metadata de la inferencia: `kind`, `observado`, `prevalencia` (cuando aplica). El locator estructural NO va acá.
- `## Alcance` lleva los globs para `estructural`/`patrón` — los mismos que el validador usa como ancla de drift.
- `## Contexto`, `## Decisión`, `## Consecuencias`, `## Alternativas consideradas` quedan **vacíos** en ADRs `Inferred`. Mine NUNCA infiere el porqué.
- `## Reglas verificables` se omite para `Inferred` (solo aplica a `Accepted`).

---

## Reglas de confianza — router + visible label (NO banda numérica)

La confianza hace exactamente dos cosas: (a) **routear** cada candidato (draft Inferred vs pregunta abierta → bootstrap); (b) **etiquetar + ordenar** en el gate para que la lista sea de alta señal y barata de revisar. La corrección la garantiza el gate humano, no un umbral.

### Reglas por kind

**config:**
- Entrada en el manifest existe → `high`. Un solo dato es autoritativo (es un hecho, no una inferencia).
- Siempre → candidato Inferred.

**estructural / patrón:**
- La prevalencia es una *señal mostrada verbatim* (ej: `40/42 handlers`), NO un umbral.
- Patrón domina sus sitios aplicables → `high` → candidato Inferred.
- Claramente marginal, o patrón compitiendo a peso similar → `low` → pregunta abierta.
- Genuinamente ambiguo → preferir mostrar como candidato con label `low` en el gate (el humano filtra barato) en lugar de pre-rutear.

**ausencia:**
- Ausencia uniforme en todos los sitios esperados → `high`.
- Cualquier presencia parcial → `low`.
- `high` → candidato Inferred; `low` → pregunta abierta.

**Low / claramente-no-es-una-decisión:**
→ pregunta abierta ruteada a bootstrap. NUNCA un Inferred vacío.

---

## Mode A — Brownfield scan (invocado por la skill)

Mode A corre cuando el usuario pide explícitamente minear decisiones del repo. Es independiente del flag `flagDecisionGaps`. NO requiere que `.matecito-ai/adr/` exista: si no existe, cada candidato es un hueco nuevo y mine puede bootstrapear la carpeta + los primeros ADR tras el confirm.

### Flujo

1. **Referencia:** cargar el catálogo de concerns de bootstrap (READ-ONLY) como taxonomía de clasificación. `.matecito-ai/adr/` puede no existir — su ausencia significa "nada decidido todavía" (todo será hueco), no es un guard de salida.
2. **Preflight codegraph:** ¿existe `.codegraph/`?
   - Sí → codegraph como fuente primaria (symbols, edges, recurring shapes).
   - No → grep como fallback.
3. **Executor escanea** → construye `candidates[]` con evidencia por kind. **Dedup por-candidato:** mapear cada hallazgo a un concern y chequear si ya existe `.matecito-ai/adr/<dominio>/<slug>.md` — si existe → saltar (o drift-check en re-run); si no existe → es un hueco → Inferred.
4. **Discovery report** (ANTES de cualquier escritura):
   - Candidatos agrupados por `proposedDomain`, con `kind · observado · prevalencia · confidence`.
   - Resumen: "Se draftearían Inferred: N / Preguntas abiertas → bootstrap: M / Posibles gaps del catálogo: K".
5. **Gate → confirm → materialize** (ver pipeline arriba).
6. **Re-run sobre ADRs existentes = detección de drift:**
   - Si los globs de `## Alcance` dejaron de matchear → drift ADR-vs-código.
   - Si `observado` ya no coincide con la evidencia actual → divergencia.
   - Reportar drift por separado; NO es un error de coherencia interna (ese es trabajo de validate).

---

## Mode B — In-flow gap detection (opt-in, flagDecisionGaps)

Mode B es **un contexto de invocación**, no una lógica distinta del ejecutor: el orquestador, cuando `flagDecisionGaps` resuelve a `true` y hay gaps implementados, invoca al ejecutor con scope = la gap list. **El ejecutor NO lee config ni resuelve el flag** — eso lo hace el orquestador (ver CLAUDE.md, regla "SDD agent launch" + "Decision-Gap Capture"), igual que con `strictTdd`/modelo. Con flag off el orquestador no invoca: silencio total. La presencia de `.matecito-ai/adr/` NO condiciona nada — sin ADRs los huecos se bootstrapean; la existencia de un ADR se chequea por-candidato (dedup).

### Hooks en el flow SDD (contratos mínimos)

**sdd-tasks (hook de detección):**
- Un gap = una task que lleva `· adr: <dominio>/<slug>` cuyo ADR destino NO existe bajo `.matecito-ai/adr/`.
- Las referencias `· adr:` colgadas son la lista de gaps — se cargan verbatim en el artifact `tasks` (ningún campo nuevo).
- Solo cuando flag-on; si no → omitir silenciosamente. NO depende de que existan ADRs: con cero ADRs, toda `· adr:` que toque una decisión apunta a algo inexistente, así que todo es hueco (bootstrap de los primeros).

**sdd-verify (hook de confirmación):**
- Para cada gap detectado en tasks (referencia `· adr:` colgada), confirmar que: (a) la task correspondiente está completa, (b) su `criteria:` pasó en el código embarcado.
- Emitir una sección `## Decision Gaps` en el verify-report: `| dominio/slug | task | implemented? |`.
- Solo cuando flag-on.

**Boundary dispatch (verify→archive — dispatch condicional del orchestrator):**
- Trigger: `flagDecisionGaps resuelto-true` AND el verify-report lista ≥1 gap implementado.
- El orchestrator despacha el executor de mine Mode B: le pasa `change-name`, la lista de gaps implementados (`dominio/slug` + hint de `Alcance`), y el root del repo.
- El executor minea el código **real embarcado** (evidencia fuerte post-ship).
- Retorna `candidates[]`; el thread principal presenta el gate; confirmed → materializar ADRs Inferred.
- Si no hay gaps implementados → no se despacha (silencio).

**sdd-archive:**
- NO tiene hook. Los ADR (cualquier status, incl. Inferred) viven SOLO en sus `.md` bajo `.matecito-ai/adr/` — nunca se registran en Engram ni en el archive-report.

---

## Invariante de taxonomía cerrada

Mine NUNCA inventa concerns en el catálogo compartido. El `kind` clasifica la evidencia; el **concern** (del catálogo) clasifica a qué decisión pertenece. Cuando un hallazgo con evidencia real NO matchea ningún concern del catálogo:

1. **Igual se le asigna un dominio canónico** (todo ADR cae en un dominio; los dominios son cerrados, los concerns dentro de un dominio son lo que puede crecer).
2. Se materializa como **ADR custom project-local** (`.matecito-ai/adr/<dominio>/<slug>.md`) — como la "fase custom" de bootstrap. NUNCA edita el catálogo compartido.
3. Se **flagea como catalog-gap** (advisory: candidato a Ratchet; promover al catálogo es un acto manual en el repo matecito-ai).
4. **Confianza conservadora:** el catálogo ES la definición curada de "qué cuenta como decisión que merece ADR". Sin un concern que lo ancle, mine está menos seguro de que sea una decisión → inclinar hacia mostrarlo como candidato flaggeado / pregunta (el humano juzga en el gate), no auto-draftear en silencio.

El candidato lleva `concern: <slug> | null`. `null` → custom local + catalog-gap + confianza conservadora.

Mine PUEDE señalar "este concern no está en el catálogo"; NUNCA puede agregar concerns al catálogo.

---

## Escalado — muchos gaps (contrato del orquestador / thread principal)

El ejecutor es `scope → candidates[]` y NO escala por sí mismo; el escalado lo maneja el caller:

- **Fan-out:** si la gap list es grande, el orquestador la parte en batches y despacha **varios ejecutores en paralelo**, cada uno con un slice del scope. Ninguno sabe que es uno de varios.
- **Merge + dedup:** el orquestador mergea los `candidates[]` y **deduplica por `dominio/slug`** (gaps distintos o ejecutores paralelos pueden proponer el mismo) antes del gate.
- **Gate con bulk:** el cuello de botella es la confirmación humana, no el scan. El thread principal ordena por confianza, agrupa por dominio, muestra un resumen primero ("N high / M a revisar / K preguntas"), y ofrece acciones bulk (`accept-all` de high, accept por-dominio, o por-ítem). Si son muchos, presentar en rondas por dominio.
- **Materialize en batch:** escribir los confirmados y actualizar `.matecito-ai/adr/INDEX.md` **una sola vez al final**, no por cada ADR.

(Muchos gaps en un solo change suele ser un olor: la tarea fue demasiado grande, o el repo recién arranca ADRs — ahí Mode A de entrada encaja mejor que gotear in-flow.)

---

## Materialización de un ADR Inferred

Después del confirm en el gate, para cada candidato aceptado:

1. Leer `~/.claude/references/adr/templates/adr.md` (READ-ONLY).
2. Completar el header: `Status: Inferred`, `Type: <proposedType>`, `Date: <hoy>`.
3. Dejar `## Contexto`, `## Decisión`, `## Consecuencias`, `## Alternativas consideradas` **vacíos** (el humano los completa al promover a Accepted).
4. Llenar `## Evidencia (inferida)` con `kind`, `observado`, `prevalencia` (si aplica para el kind).
5. Para `estructural`/`patrón`: llenar `## Alcance` con `proposedAlcanceGlobs`.
6. Omitir `## Reglas verificables` (solo Accepted).
7. Escribir en `.matecito-ai/adr/<proposedDomain>/<proposedSlug>.md`.
8. Si la carpeta del dominio no existía: `mkdir -p .matecito-ai/adr/<proposedDomain>` y crear `INDEX.md` mínimo para el dominio.
9. Actualizar `.matecito-ai/adr/INDEX.md` con la entrada del nuevo ADR.

---

## Shape del candidate (lo que retorna el executor)

```json
{
  "kind": "estructural | config | patrón | ausencia",
  "observado": "descripción del QUÉ visto — sin el porqué",
  "prevalencia": "40/42 handlers | null",
  "confidence": "high | low",
  "concern": "<concern-slug del catálogo> | null",
  "proposedDomain": "<dominio canónico>",
  "proposedSlug": "<kebab-case>",
  "proposedAlcanceGlobs": ["src/**/*.routes.ts"],
  "proposedType": "decision | convention | policy",
  "lowSignalReason": "descripción si confidence es low | null"
}
```

El executor retorna un bloque markdown con un array JSON bajo el header `## candidates`:

```markdown
## candidates

\`\`\`json
[
  { ... },
  { ... }
]
\`\`\`

## open_questions

Lista de observaciones de baja señal sugeridas a bootstrap como preguntas abiertas.

## catalog_gap_flags

Advisory: concerns que no encajan en ningún dominio canónico (no es acción bloqueante).
```

---

## Re-run y detección de drift

Cuando mine corre sobre un repo que ya tiene ADRs `Inferred`:

1. Para cada ADR Inferred existente, volver a buscar la evidencia original:
   - ¿Los globs de `## Alcance` siguen matcheando algo?
   - ¿El `observado` en `## Evidencia (inferida)` sigue siendo verdad?
2. Si los globs dejaron de matchear o `observado` divergió → reportar como **drift ADR-vs-código** (distinto del chequeo de coherencia interna de validate).
3. El reporte de drift va en una sección separada `## Drift detectado` antes del discovery report de nuevos candidatos.

---

## Anti-patterns

- No escribir ADRs antes del confirm → invariante 2.
- No correr Mode B cuando el flag es false → invariante 1. La ausencia de `.matecito-ai/adr/` NO es razón para no correr: mine bootstrapea los primeros.
- No inventar el porqué de una decisión → WHY siempre vacío en Inferred.
- No usar `path:line` como locator → siempre globs estructurales o entradas de manifest.
- No crear Inferred con confidence low → los low van a bootstrap como pregunta abierta.
- No marcar un ADR Inferred como Accepted → eso lo hace el humano vía bootstrap modo update.
- No inventar dominios nuevos ni concerns nuevos → taxonomía cerrada, solo flagear gaps.
- No reportar drift como error de coherencia interna → es drift ADR-vs-código, categoría distinta.
