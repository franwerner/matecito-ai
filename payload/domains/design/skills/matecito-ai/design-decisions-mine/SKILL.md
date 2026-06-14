---
name: design-decisions-mine
description: Minería de decisiones de DISEÑO desde el archivo Figma conectado de un proyecto. Produce candidatos a DDR con Status Inferred a partir de evidencia de tokens, componentes, patrones repetidos o ausencias. Dos modos — Mode A (scan del archivo Figma invocado por la skill) y Mode B (in-flow, opt-in via flagDecisionGaps). Usá esta skill cuando el usuario pida "minear decisiones de diseño", "encontrar DDRs que faltan", "¿qué decisiones de diseño hay implícitas en el Figma?", o cuando el orchestrator dispare el boundary dispatch post-verify con gaps implementados.
---

# Design Decisions Mine

Minería consultativa de decisiones de diseño implícitas en el archivo Figma: escanea el archivo, construye candidatos a DDR con evidencia observable (styles/components nombrados, patrones repetidos, ausencias), y los presenta para confirmación humana antes de escribir cualquier cosa. Nunca escribe DDRs sin confirmación explícita.

> **Arquitectura de esta skill — motor/datos igual que bootstrap.** `SKILL.md` (este archivo) es el motor: pipeline, modelo de evidencia, reglas de confianza, flujos Mode A y Mode B, invariantes. El catálogo de surfaces y la taxonomía de design concerns son de `design-decisions-bootstrap` (READ-ONLY); las plantillas de estructura de DDR (incl. `ddr.md`) viven en la referencia `~/.claude/references/ddr/templates/` (READ-ONLY). La salida son DDRs con `Status: Inferred` en `.matecito-ai/ddr/<surface>/<slug>.md`, misma taxonomía.

> **Concepto de DDR — referencia canónica.** Qué ES y qué NO ES un DDR, y la diferencia borrador(inferido)/aceptado, están en `~/.claude/references/ddr/README.md` (referencia consultable, agnóstica de flujo/skill). mine NO redefine el concepto — lo aplica. Si dudás si un hallazgo "es un DDR", esa referencia es la regla.

---

## Arquitectura de esta skill

La skill está partida en **motor** y **executor**:

- **`SKILL.md` (este archivo) = el motor.** Define el pipeline, el modelo de evidencia, las reglas de confianza (router + label), los flujos de Mode A y Mode B, y los invariantes. Es el contrato que el executor debe seguir.
- **`payload/agents/design-decisions-mine.md` = el executor.** Agente de contexto fresco que hace el trabajo pesado de scan/discovery sobre el archivo Figma (read-only) y RETORNA un bloque `candidates[]`. No escribe DDRs. El executor sirve tanto Mode A como Mode B.
- **El thread principal = gate + materialize.** Recibe `candidates[]` del executor, renderiza la tabla de confirmación, y ejecuta la materialización solo después del confirm explícito del usuario.

---

## Fuente de evidencia: Figma (no codegraph)

La evidencia de esta skill sale del **archivo Figma conectado, leído read-only** vía el MCP figma — NO de codegraph (eso es el dominio development). La evidencia FUERTE son los **styles y components nombrados**; un mockup con valores sueltos sin tokens es señal débil.

**Canva queda EXCLUIDO de la minería.** Canva no expone tokens legibles vía API, así que no produce evidencia fuerte de decisiones. Si una pieza vive solo en Canva, la skill lo reporta y no infiere decisiones a partir de ella.

---

## Invariantes (no negociables)

Estas dos invariantes son la columna vertebral de la skill. Cualquier desviación las rompe.

**Invariante 1 — DDRs nunca obligatorios; el ejecutor recibe un scope, no decide nada.** El ejecutor NO lee config, NO resuelve el flag, NO se ramifica por modo: recibe un **scope** de su caller y lo procesa. La intención (si correr) vive en el caller:

- invocación directa de la skill → scope = archivo Figma completo.
- el orquestador (con `flagDecisionGaps` resuelto en true + gaps implementados) → scope = la gap list.

Con flag off el orquestador no invoca → silencio total. La dependencia dura del ejecutor es el **catálogo de design concerns** (siempre presente), de donde sale la taxonomía de surfaces — NO los DDR generados. La ausencia de `.matecito-ai/ddr/` NO bloquea: "nada decidido todavía" → cada candidato es un hueco y mine puede bootstrapear (crear la carpeta + los primeros DDR) tras el confirm. La presencia de un DDR se chequea **por-candidato** (dedup), nunca como gate global.

**Invariante 2 — Nunca auto-materializar.** El executor NO tiene capacidad de escritura de DDRs en la fase discover. La materialización es un paso separado, alcanzable solo después de un confirm explícito del usuario en el gate del thread principal.

---

## Pipeline runtime: `discover → confirm → materialize`

El pipeline corre en tres contextos distintos por diseño:

```
[executor — contexto fresco]          [thread principal]
discover → draft candidates[]   →→→   confirm (gate) → materialize
```

### Paso 1: Discover (executor)

El executor escanea el archivo Figma, construye candidatos con evidencia, y RETORNA un bloque `candidates[]` estructurado. No escribe nada.

### Paso 2: Confirm (thread principal — gate)

El thread principal recibe los candidatos y renderiza la tabla de confirmación:

```
Candidatos ordenados por confidence label (high primero):
┌─────────────────────────────────────────────────────────────┐
│ [surface/slug]  kind · observado · prevalencia · confidence │
│ ✅ alto-signal-candidate  →  [aceptar / editar / saltar]    │
│ ⬜ otro-candidate         →  [aceptar / editar / saltar]    │
└─────────────────────────────────────────────────────────────┘

Preguntas abiertas (→ bootstrap, no Inferred):
- [pregunta 1]: ¿es esto una decisión?
- ...

Posibles gaps del catálogo (advisory):
- concern X no encaja en ninguna surface canónica
```

Opciones de confirm: `accept-all` / por ítem / `none`.

### Paso 3: Materialize (thread principal — solo post-confirm)

Para cada candidato aceptado/editado, el thread principal escribe el DDR con `Status: Inferred` usando `~/.claude/references/ddr/templates/ddr.md` (READ-ONLY). Ver "Materialización" más abajo.

---

## Modelo de evidencia — keyed by `kind`

El `kind` de cada candidato determina qué se lee del archivo Figma durante el scan y qué secciones del DDR se llenan.

| kind | qué se lee | llena `## Evidencia (inferida)` | llena `## Alcance` |
|---|---|---|---|
| `token` | `get_styles` / variables del `get_file` — un style o variable **nombrado** (`Primary/500`, `Heading/H1`) | kind + observado + prevalencia | la lista de tokens / frames que gobierna |
| `component` | `get_components` — un componente o set de variantes | kind + observado + prevalencia | el set de componentes / frames que lo instancian |
| `pattern` | `get_file` / `get_node` — un valor o estilo **repetido sin nombrar** (mismo hex/px a mano) | kind + observado + prevalencia | los frames donde se repite |
| `absence` | scan probando ausencia en sitios esperados (no hay color styles, no hay text styles, etc.) | kind + observado (sin prevalencia) | omitido (no hay locator) |

Reglas de llenado:

- `## Evidencia (inferida)` contiene solo la metadata de la inferencia: `kind`, `observado`, `prevalencia` (cuando aplica). El locator NO va acá.
- `## Alcance` lleva el locator (lista de tokens / frames / set de componentes) para `token`/`component`/`pattern` — el mismo que el validador usa como ancla de drift.
- `## Contexto`, `## Decisión`, `## Consecuencias`, `## Alternativas consideradas` quedan **vacíos** en DDRs `Inferred`. Mine NUNCA infiere el porqué.
- `## Reglas verificables` se omite para `Inferred` (solo aplica a `Accepted`).

---

## Reglas de confianza — router + visible label (NO banda numérica)

La confianza hace exactamente dos cosas: (a) **routear** cada candidato (draft Inferred vs pregunta abierta → bootstrap); (b) **etiquetar + ordenar** en el gate para que la lista sea de alta señal y barata de revisar. La corrección la garantiza el gate humano, no un umbral.

### Reglas por kind

**token / component:**
- Un style/variable o componente **nombrado** existe → `high`. Un style nombrado es autoritativo (es un hecho del archivo, no una inferencia).
- Siempre → candidato Inferred.

**pattern:**
- La prevalencia es una *señal mostrada verbatim* (ej: `usado en 12/14 frames`), NO un umbral.
- Patrón domina sus sitios aplicables → `high` → candidato Inferred.
- Claramente marginal, o patrón compitiendo a peso similar → `low` → pregunta abierta.
- Un valor suelto sin token que lo ancle (mockup sin styles) → `low` → pregunta abierta. NUNCA inventar una decisión a partir de un mockup sin tokens.
- Genuinamente ambiguo → preferir mostrar como candidato con label `low` en el gate (el humano filtra barato) en lugar de pre-rutear.

**absence:**
- Ausencia uniforme en todos los sitios esperados → `high`.
- Cualquier presencia parcial → `low`.
- `high` → candidato Inferred; `low` → pregunta abierta.

**Low / claramente-no-es-una-decisión:**
→ pregunta abierta ruteada a bootstrap. NUNCA un Inferred vacío.

---

## Mode A — Scan del archivo Figma (invocado por la skill)

Mode A corre cuando el usuario pide explícitamente minear decisiones de diseño del archivo Figma. Es independiente del flag `flagDecisionGaps`. NO requiere que `.matecito-ai/ddr/` exista: si no existe, cada candidato es un hueco nuevo y mine puede bootstrapear la carpeta + los primeros DDR tras el confirm.

### Flujo

1. **Referencia:** cargar el catálogo de design concerns de bootstrap (READ-ONLY) como taxonomía de clasificación. `.matecito-ai/ddr/` puede no existir — su ausencia significa "nada decidido todavía" (todo será hueco), no es un guard de salida.
2. **Preflight Figma:** ¿hay un archivo Figma conectado vía el MCP figma?
   - Sí → Figma como fuente primaria (`get_styles`, `get_components`, `get_file`/`get_node`).
   - No → no hay fuente de evidencia fuerte → `status: silenced`. Canva no cuenta (no expone tokens).
3. **Executor escanea** → construye `candidates[]` con evidencia por kind. **Dedup por-candidato:** mapear cada hallazgo a un concern y chequear si ya existe `.matecito-ai/ddr/<surface>/<slug>.md` — si existe → saltar (o drift-check en re-run); si no existe → es un hueco → Inferred.
4. **Discovery report** (ANTES de cualquier escritura):
   - Candidatos agrupados por `proposedSurface`, con `kind · observado · prevalencia · confidence`.
   - Resumen: "Se draftearían Inferred: N / Preguntas abiertas → bootstrap: M / Posibles gaps del catálogo: K".
5. **Gate → confirm → materialize** (ver pipeline arriba).
6. **Re-run sobre DDRs existentes = detección de drift:**
   - Si `observado` ya no coincide con la evidencia actual del archivo Figma (el style nombrado desapareció, el componente se borró) → divergencia.
   - Reportar drift por separado; NO es un error de coherencia interna (ese es trabajo de validate).

---

## Mode B — In-flow gap detection (opt-in, flagDecisionGaps)

Mode B es **un contexto de invocación**, no una lógica distinta del ejecutor: el orquestador, cuando `flagDecisionGaps` resuelve a `true` y hay gaps implementados, invoca al ejecutor con scope = la gap list. **El ejecutor NO lee config ni resuelve el flag** — eso lo hace el orquestador (ver CLAUDE.md, regla "design agent launch" + "Decision-Gap Capture"), igual que con el modelo. Con flag off el orquestador no invoca: silencio total. La presencia de `.matecito-ai/ddr/` NO condiciona nada — sin DDRs los huecos se bootstrapean; la existencia de un DDR se chequea por-candidato (dedup).

### Hooks en el flow design (contratos mínimos)

**design-tasks (hook de detección):**
- Un gap = una task que lleva `· ddr: <surface>/<slug>` cuyo DDR destino NO existe bajo `.matecito-ai/ddr/`.
- Las referencias `· ddr:` colgadas son la lista de gaps — se cargan verbatim en el artifact `tasks` (ningún campo nuevo).
- Solo cuando flag-on; si no → omitir silenciosamente. NO depende de que existan DDRs: con cero DDRs, toda `· ddr:` que toque una decisión apunta a algo inexistente, así que todo es hueco (bootstrap de los primeros).

**design-verify (hook de confirmación):**
- Para cada gap detectado en tasks (referencia `· ddr:` colgada), confirmar que: (a) la task correspondiente está completa, (b) su `criteria:` pasó contra el archivo Figma embarcado.
- Emitir una sección `## Decision Gaps` en el verify-report: `| surface/slug | task | implemented? |`.
- Solo cuando flag-on.

**Boundary dispatch (verify→archive — dispatch condicional del orchestrator):**
- Trigger: `flagDecisionGaps resuelto-true` AND el verify-report lista ≥1 gap implementado.
- El orchestrator despacha el executor de mine Mode B: le pasa `change-name`, la lista de gaps implementados (`surface/slug` + hint de `Alcance`), y el root del repo.
- El executor minea el **archivo Figma real embarcado** (evidencia fuerte post-ship: los styles/components quedaron definidos).
- Retorna `candidates[]`; el thread principal presenta el gate; confirmed → materializar DDRs Inferred.
- Si no hay gaps implementados → no se despacha (silencio).

**design-archive:**
- NO tiene hook. Los DDR (cualquier status, incl. Inferred) viven SOLO en sus `.md` bajo `.matecito-ai/ddr/` — nunca se registran en Engram ni en el archive-report.

---

## Invariante de taxonomía cerrada

Mine NUNCA inventa concerns en el catálogo compartido. El `kind` clasifica la evidencia; el **concern** (del catálogo) clasifica a qué decisión pertenece. Cuando un hallazgo con evidencia real NO matchea ningún concern del catálogo:

1. **Igual se le asigna una surface canónica** (todo DDR cae en una surface; las surfaces son cerradas, los concerns dentro de una surface son lo que puede crecer).
2. Se materializa como **DDR custom project-local** (`.matecito-ai/ddr/<surface>/<slug>.md`) — como la "fase custom" de bootstrap. NUNCA edita el catálogo compartido.
3. Se **flagea como catalog-gap** (advisory: candidato a Ratchet; promover al catálogo es un acto manual en el repo matecito-ai).
4. **Confianza conservadora:** el catálogo ES la definición curada de "qué cuenta como decisión que merece DDR". Sin un concern que lo ancle, mine está menos seguro de que sea una decisión → inclinar hacia mostrarlo como candidato flaggeado / pregunta (el humano juzga en el gate), no auto-draftear en silencio.

El candidato lleva `concern: <slug> | null`. `null` → custom local + catalog-gap + confianza conservadora.

Mine PUEDE señalar "este concern no está en el catálogo"; NUNCA puede agregar concerns al catálogo.

---

## Escalado — muchos gaps (contrato del orquestador / thread principal)

El ejecutor es `scope → candidates[]` y NO escala por sí mismo; el escalado lo maneja el caller:

- **Fan-out:** si la gap list es grande, el orquestador la parte en batches y despacha **varios ejecutores en paralelo**, cada uno con un slice del scope. Ninguno sabe que es uno de varios.
- **Merge + dedup:** el orquestador mergea los `candidates[]` y **deduplica por `surface/slug`** (gaps distintos o ejecutores paralelos pueden proponer el mismo) antes del gate.
- **Gate con bulk:** el cuello de botella es la confirmación humana, no el scan. El thread principal ordena por confianza, agrupa por surface, muestra un resumen primero ("N high / M a revisar / K preguntas"), y ofrece acciones bulk (`accept-all` de high, accept por-surface, o por-ítem). Si son muchos, presentar en rondas por surface.
- **Materialize en batch:** escribir los confirmados y actualizar `.matecito-ai/ddr/INDEX.md` **una sola vez al final**, no por cada DDR.

(Muchos gaps en un solo change suele ser un olor: la tarea fue demasiado grande, o el repo recién arranca DDRs — ahí Mode A de entrada encaja mejor que gotear in-flow.)

---

## Materialización de un DDR Inferred

Después del confirm en el gate, para cada candidato aceptado:

1. Leer `~/.claude/references/ddr/templates/ddr.md` (READ-ONLY).
2. Completar el header: `Status: Inferred`, `Type: <proposedType>`, `Date: <hoy>`.
3. Dejar `## Contexto`, `## Decisión`, `## Consecuencias`, `## Alternativas consideradas` **vacíos** (el humano los completa al promover a Accepted).
4. Llenar `## Evidencia (inferida)` con `kind`, `observado`, `prevalencia` (si aplica para el kind).
5. Para `token`/`component`/`pattern`: llenar `## Alcance` con el locator (lista de tokens / frames / set de componentes).
6. Omitir `## Reglas verificables` (solo Accepted).
7. Escribir en `.matecito-ai/ddr/<proposedSurface>/<proposedSlug>.md`.
8. Si la carpeta de la surface no existía: `mkdir -p .matecito-ai/ddr/<proposedSurface>` y crear `INDEX.md` mínimo para la surface.
9. Actualizar `.matecito-ai/ddr/INDEX.md` con la entrada del nuevo DDR.

---

## Shape del candidate (lo que retorna el executor)

```json
{
  "kind": "token | component | pattern | absence",
  "observado": "descripción del QUÉ visto en Figma — sin el porqué",
  "prevalencia": "usado en 12/14 frames | null",
  "confidence": "high | low",
  "concern": "<concern-slug del catálogo> | null",
  "proposedSurface": "foundation | components | layout | brand | accessibility",
  "proposedSlug": "<kebab-case>",
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

Advisory: concerns que no encajan en ninguna surface canónica (no es acción bloqueante).
```

---

## Re-run y detección de drift

Cuando mine corre sobre un repo que ya tiene DDRs `Inferred`:

1. Para cada DDR Inferred existente, volver a buscar la evidencia original en el archivo Figma:
   - ¿El locator de `## Alcance` (tokens / frames / componentes) sigue existiendo?
   - ¿El `observado` en `## Evidencia (inferida)` sigue siendo verdad (el style nombrado/componente sigue ahí)?
2. Si el locator desapareció o `observado` divergió → reportar como **drift DDR-vs-Figma** (distinto del chequeo de coherencia interna de validate).
3. El reporte de drift va en una sección separada `## Drift detectado` antes del discovery report de nuevos candidatos.

---

## Anti-patterns

- No escribir DDRs antes del confirm → invariante 2.
- No correr Mode B cuando el flag es false → invariante 1. La ausencia de `.matecito-ai/ddr/` NO es razón para no correr: mine bootstrapea los primeros.
- No inventar el porqué de una decisión → WHY siempre vacío en Inferred.
- No inventar una decisión a partir de un mockup con valores sueltos sin tokens → eso es señal débil → pregunta abierta.
- No minar Canva → no expone tokens legibles; reportarlo y no inferir.
- No crear Inferred con confidence low → los low van a bootstrap como pregunta abierta.
- No marcar un DDR Inferred como Accepted → eso lo hace el humano vía bootstrap modo update.
- No inventar surfaces nuevas ni concerns nuevos → taxonomía cerrada, solo flagear gaps.
- No reportar drift como error de coherencia interna → es drift DDR-vs-Figma, categoría distinta.
