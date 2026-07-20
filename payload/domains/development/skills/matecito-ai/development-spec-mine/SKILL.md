---
name: development-spec-mine
description: Minería de comportamiento as-built desde el código fuente de un repo existente. Produce candidatos a capability-spec con Status Inferred a partir de evidencia de route-handlers, validaciones, máquinas de estado y event-handlers, corroborada (o no) por tests. Un solo modo — Mode A (scan brownfield invocado por la skill, o disparado por el flag flagSpecMine tras sdd-init). Usá esta skill cuando el usuario pida "minear comportamiento", "encontrar specs que faltan", "¿qué hace el sistema que no está documentado?", o cuando el flag flagSpecMine dispare el scan automático.
---

# Development Spec Mine

Minería consultativa de comportamiento as-built implícito en el código: escanea el repo, construye candidatos a capability-spec con evidencia observable, y los presenta para confirmación humana antes de escribir cualquier cosa. Nunca escribe capability-specs sin confirmación explícita.

> **Arquitectura de esta skill — motor/datos igual que bootstrap.** `SKILL.md` (este archivo) es el motor: pipeline, modelo de evidencia, reglas de confianza, flujo Mode A, invariantes. La taxonomía de tipos (fija: `flow`/`rule`/`lifecycle`/`process`) es de `development-spec-bootstrap` (READ-ONLY); las plantillas de estructura del capability-spec viven en la referencia `~/.claude/references/spec/templates/` (READ-ONLY). La salida son capability-specs con `Status: Inferred` en `.matecito-ai/development-specs/<type>/<capability>.md`, misma taxonomía.

> **Concepto de capability-spec — referencia canónica.** Qué ES y qué NO ES un capability-spec, y la diferencia borrador(inferido)/aceptado, están en `~/.claude/references/spec/README.md` (referencia agnóstica de flujo/skill). mine NO redefine el concepto — lo aplica. Si dudás si un hallazgo "es un capability-spec", esa referencia es la regla.

---

## Arquitectura de esta skill

La skill está partida en **motor** y **executor**:

- **`SKILL.md` (este archivo) = el motor.** Define el pipeline, el modelo de evidencia, las reglas de confianza (router + label), el flujo Mode A, y los invariantes. Es el contrato que el executor debe seguir.
- **`payload/agents/development-spec-mine.md` = el executor.** Agente de contexto fresco que hace el trabajo pesado de scan/discovery y RETORNA un bloque `candidates[]`. No escribe capability-specs.
- **El thread principal = gate + materialize.** Recibe `candidates[]` del executor, renderiza la tabla de confirmación, y ejecuta la materialización solo después del confirm explícito del usuario.

---

## Invariantes (no negociables)

Estas dos invariantes son la columna vertebral de la skill. Cualquier desviación las rompe.

**Invariante 1 — capability-specs Inferred nunca obligatorios; el ejecutor recibe un scope, no decide nada.** El ejecutor NO lee config, NO resuelve el flag `flagSpecMine`, NO se ramifica por modo: recibe un **scope** de su caller y lo procesa. La intención (si correr) vive en el caller:

- invocación directa de la skill → scope = repo completo (o el área que pidió el usuario).
- el orquestador, con `flagSpecMine` resuelto en `true` tras `sdd-init`/inicio de sesión sobre un repo con store ausente/disperso → scope = repo completo.

Con flag off el orquestador no invoca → silencio total. A diferencia de `development-decisions-mine`, spec-mine **no tiene Mode B / gap list**: siempre escanea comportamiento as-built completo, nunca un subconjunto in-flow. La dependencia dura del ejecutor es la **taxonomía de tipos** (siempre presente, fija) — NO un catálogo de concerns (los capability-specs no tienen catálogo de datos). La ausencia de `.matecito-ai/development-specs/` NO bloquea: "nada especificado todavía" → cada candidato es un hueco y mine puede bootstrapear (crear la carpeta + los primeros specs) tras el confirm. La presencia de un spec se chequea **por-candidato** (dedup), nunca como gate global.

**Invariante 2 — Nunca auto-materializar.** El executor NO tiene capacidad de escritura de capability-specs en la fase discover. La materialización es un paso separado, alcanzable solo después de un confirm explícito del usuario en el gate del thread principal.

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
│ [type/capability]  kind · observado · prevalencia · confidence │
│ ✅ alto-signal-candidate  →  [aceptar / editar / saltar]    │
│ ⬜ otro-candidate         →  [aceptar / editar / saltar]    │
└─────────────────────────────────────────────────────────────┘

Preguntas abiertas (→ development-spec-bootstrap, no Inferred):
- [pregunta 1]: ¿es esto comportamiento verificable?
- ...
```

Opciones de confirm: `accept-all` / por ítem / `none`.

### Paso 3: Materialize (thread principal — solo post-confirm)

Para cada candidato aceptado/editado, el thread principal escribe el capability-spec con `Status: Inferred` usando `~/.claude/references/spec/templates/capability.md` (READ-ONLY). Ver "Materialización" más abajo.

---

## Modelo de evidencia — keyed by `kind`

El `kind` de cada candidato determina qué se lee durante el scan y qué sección esqueleto del capability-spec se llena.

| kind | qué se lee | llena en el capability-spec | proposedType |
|---|---|---|---|
| `route-handler` | codegraph (routing table + handler callgraph) ▸ grep | `Flujo principal` | `flow` |
| `validation` | codegraph (guard/validation sites) ▸ grep | `Reglas de negocio` | `rule` |
| `state-machine` | codegraph (entity enum + transition edges) ▸ grep | `Entidades y estados` | `lifecycle` |
| `event-handler` | codegraph (subscribers/webhooks/jobs) ▸ grep | `Flujo principal` (de proceso) | `process` |
| `test-assertion` | tests (Given/When/Then arrange/act/assert) | `Escenarios` — y sube la confidence del kind que corrobora | — (oráculo) |

Reglas de llenado:

- El candidato lleva `observado` (el QUÉ estructural/de comportamiento visto, sin el porqué) y `prevalencia` (cuando aplica, mostrada verbatim).
- `test-assertion` NO llena una sección de evidencia propia: alimenta `Escenarios` (vía `proposedScenarios`) y decide el ruteo de confidence del kind que corrobora.
- `Propósito`, `Actores`, `Precondiciones`, `Ramas`, `Casos borde`, `Errores de cara al actor` y `Referencias` quedan **sin completar u omitidos** en capability-specs `Inferred` salvo que la evidencia observada los sustente directamente — mine no inventa contexto que no vio (ver "Materialización").

---

## Reglas de confianza — router + visible label (NO banda numérica)

La confianza hace exactamente dos cosas: (a) **routear** cada candidato (draft Inferred vs pregunta abierta); (b) **etiquetar + ordenar** en el gate para que la lista sea de alta señal y barata de revisar. La corrección la garantiza el gate humano, no un umbral.

### Reglas por kind

**route-handler / validation / state-machine / event-handler (shape-kinds):**
- Corroborado por evidencia `test-assertion` (un test cuyo Given/When/Then cubre el comportamiento observado) → `high` → candidato Inferred firme; el escenario derivado del test llena `proposedScenarios`.
- Solo evidencia de código crudo, sin `test-assertion` que lo corrobore → `low` → pregunta abierta, NUNCA un candidato firme.
- La prevalencia (cuando aplica) es una *señal mostrada verbatim* (ej: `8/8 handlers cubiertos`), NO un umbral — no decide el ruteo, solo informa al humano en el gate.

**test-assertion:**
- No genera candidato propio. Es el oráculo que corrobora (o no) a los shape-kinds de arriba.

**Low / claramente-no-es-comportamiento-verificable:**
→ pregunta abierta. NUNCA un Inferred vacío.

---

## Mode A — Brownfield scan (único modo)

A diferencia de `development-decisions-mine`, spec-mine **no tiene Mode B**: corre siempre como scan brownfield, invocado por la skill directamente O disparado por el flag `flagSpecMine` (mirror tipado de `flagDecisionGaps`) tras `sdd-init`/inicio de sesión, cuando el repo tiene código pero `.matecito-ai/development-specs/` está ausente o disperso. NO requiere que el store exista: si no existe, cada candidato es un hueco nuevo y mine puede bootstrapear la carpeta + los primeros specs tras el confirm.

### Flujo

1. **Referencia:** cargar la taxonomía de tipos de `development-spec-bootstrap` (READ-ONLY, fija: `flow`/`rule`/`lifecycle`/`process`) + la regla de clasificación de `~/.claude/references/spec/README.md` → «Cómo clasificar el tipo». `.matecito-ai/development-specs/` puede no existir — su ausencia significa "nada especificado todavía", no es un guard de salida.
2. **Preflight codegraph:** ¿existe `.codegraph/`?
   - Sí → codegraph como fuente primaria (routing tables, guard sites, entity enums, event subscribers).
   - No → grep como fallback.
3. **Executor escanea** → construye `candidates[]` con evidencia por kind, corroborada o no por `test-assertion`. **Dedup por-candidato:** listar `.matecito-ai/development-specs/<type>/*.md` existentes y chequear si el candidato ya está cubierto — si existe → saltar (o drift-check en re-run); si no existe → es un hueco → Inferred.
4. **Discovery report** (ANTES de cualquier escritura):
   - Candidatos agrupados por `proposedType`, con `kind · observado · prevalencia · confidence`.
   - Resumen: "Se draftearían Inferred: N / Preguntas abiertas: M".
5. **Gate → confirm → materialize** (ver pipeline arriba).
6. **Re-run sobre capability-specs existentes = detección de drift:**
   - Si el comportamiento descrito en las secciones esqueleto de un spec `Inferred` dejó de observarse en el código → drift spec-vs-código.
   - Reportar drift por separado; NO es un error de coherencia interna (ese es trabajo de `development-spec-validate`).

---

## Trigger del flag (referencia — implementación en Fase 5)

El flag `flagSpecMine` (mirror tipado de `flagDecisionGaps`) dispara Mode A cuando resuelve `true`: el orquestador, post-`sdd-init` (o primer comando de flow de la sesión) sobre un repo con código pero store ausente/disperso, despacha spec-mine sobre el repo completo → gate obligatorio. Flag off → silencio total, comportamiento idéntico a antes de que el flag existiera. La resolución del flag y su wiring en config/TUI son responsabilidad de la Fase 5 de `development-spec-mine`; esta skill solo documenta el punto de disparo — el ejecutor NUNCA lee el flag (Invariante 1).

---

## Invariante de taxonomía cerrada

Mine NUNCA inventa tipos nuevos. Toda capability cae en uno de los 4 tipos fijos (`flow`/`rule`/`lifecycle`/`process`), determinado por el `kind` de evidencia observado (ver tabla del Paso 3 del executor). A diferencia de `development-decisions-mine`, no hay concepto de "concern fuera de catálogo" ni de `catalog_gap_flags` — la taxonomía de tipos no tiene concerns internos que puedan faltar.

---

## Escalado — muchos candidatos (contrato del orquestador / thread principal)

El ejecutor es `scope → candidates[]` y NO escala por sí mismo; el escalado lo maneja el caller:

- **Fan-out:** si el repo es grande, el orquestador puede partir el scan en batches por área y despachar **varios ejecutores en paralelo**, cada uno con un slice del scope. Ninguno sabe que es uno de varios.
- **Merge + dedup:** el orquestador mergea los `candidates[]` y **deduplica por `type/capability`** antes del gate.
- **Gate con bulk:** el cuello de botella es la confirmación humana, no el scan. El thread principal ordena por confianza, agrupa por tipo, muestra un resumen primero ("N high / M preguntas"), y ofrece acciones bulk (`accept-all` de high, accept por-tipo, o por-ítem). Si son muchos, presentar en rondas por tipo.
- **Materialize en batch:** escribir los confirmados y actualizar `.matecito-ai/development-specs/INDEX.md` (raíz y de tipo) **una sola vez al final**, no por cada spec.

---

## Materialización de un capability-spec Inferred

Después del confirm en el gate, para cada candidato aceptado:

1. Leer `~/.claude/references/spec/templates/capability.md` (READ-ONLY).
2. Completar el header: `Status: Inferred`, `Date: <hoy>`.
3. Llenar la sección esqueleto de su `proposedType` (según la tabla del Paso 3 del executor) con `observado` + la evidencia observada.
4. Si hubo `test-assertion` que corroboró el candidato: llenar `## Escenarios` con `proposedScenarios` (el Given/When/Then parseado del test).
5. Dejar `## Propósito`, `## Actores`, `## Precondiciones`, y las secciones no-esqueleto de su tipo (`Ramas`, `Casos borde`, `Errores de cara al actor`, `Referencias`) **sin completar u omitidas** salvo que la evidencia observada las sustente directamente — mine no inventa contexto que no vio. El humano las completa/corrige al ratificar (`development-spec-bootstrap` modo update, caso "Ratificar un Inferred").
6. Escribir en `.matecito-ai/development-specs/<proposedType>/<proposedCapability>.md`.
7. Si la carpeta del tipo no existía: `mkdir -p .matecito-ai/development-specs/<proposedType>` y crear `INDEX.md` mínimo para el tipo.
8. Actualizar `.matecito-ai/development-specs/INDEX.md` (raíz) y `.matecito-ai/development-specs/<proposedType>/INDEX.md` con la entrada del nuevo spec.

---

## Shape del candidate (lo que retorna el executor)

```json
{
  "kind": "route-handler | validation | state-machine | event-handler",
  "observado": "descripción del QUÉ visto — sin el porqué",
  "prevalencia": "8/8 handlers | null",
  "confidence": "high | low",
  "proposedType": "flow | rule | lifecycle | process",
  "proposedCapability": "<kebab-case>",
  "testEvidence": "referencia al test que corrobora | null",
  "proposedScenarios": ["GIVEN ... WHEN ... THEN ..."],
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

Lista de observaciones de baja señal — comportamiento observado en código crudo sin test que lo corrobore.
```

(No hay `catalog_gap_flags` — los capability-specs no tienen catálogo de concerns.)

---

## Re-run y detección de drift

Cuando mine corre sobre un repo que ya tiene capability-specs `Inferred`:

1. Para cada spec Inferred existente, volver a buscar la evidencia original según su tipo:
   - ¿El comportamiento descrito en sus secciones esqueleto sigue siendo observable en el código?
2. Si divergió → reportar como **drift spec-vs-código** (distinto del chequeo de coherencia interna de `development-spec-validate`).
3. El reporte de drift va en una sección separada `## Drift detectado` antes del discovery report de nuevos candidatos.

---

## Anti-patterns

- No escribir capability-specs antes del confirm → invariante 2.
- No inventar un Mode B / gap list → spec-mine es Mode A únicamente; la ausencia del store NO es razón para no correr: mine bootstrapea los primeros.
- No inventar el contexto de una capability (`Propósito`/`Actores`/`Precondiciones`/`Ramas`/`Casos borde`/`Errores de cara al actor`) que no fue observado directamente → esas secciones quedan sin completar hasta la ratificación humana.
- No volcar identificadores internos volátiles (clases, métodos, columnas, rutas de archivo, errores internos) en `observado` ni en ninguna sección materializada → siempre idioma de dominio + contrato público (ver `~/.claude/references/spec/README.md` → "No es el cómo").
- No crear Inferred con confidence low → los low van a `open_questions`.
- No marcar un capability-spec Inferred como Accepted → eso lo hace el humano vía `development-spec-bootstrap` modo update (caso "Ratificar un Inferred").
- No inventar tipos nuevos → taxonomía cerrada (`flow`/`rule`/`lifecycle`/`process`).
- No reportar drift como error de coherencia interna → es drift spec-vs-código, categoría que resuelve `development-spec-validate` por separado.
