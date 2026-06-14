---
name: design-decisions-bootstrap
description: Entrevista interactiva por fases para capturar las decisiones de DISEÑO de una pieza o sistema (foundation, components, layout, brand, accessibility) al iniciarlo y materializarlas como DDRs organizados por surface en .matecito-ai/ddr/<surface>/. Usá esta skill SIEMPRE que el usuario arranque un sistema de diseño nuevo, una identidad de marca, una pieza con sistema visual propio, mencione "arrancar el diseño / definir el sistema visual / setup de marca", pida ayuda con la paleta/tipografía/grid/tokens iniciales, quiera revisar o actualizar decisiones de diseño existentes, ratificar DDRs Inferred que dejó la minería de Figma, o cuando detectes una pieza sin .matecito-ai/ddr/ y el usuario esté por producir assets que fijan el sistema visual. También dispará si el usuario menciona "DDR", "decisión de diseño", "design tokens", "paleta", "escala tipográfica", "contraste", "biblioteca de componentes", "breakpoints", "tono de voz".
---

# Design Decisions Bootstrap

Entrevista al usuario para capturar las decisiones de diseño de la pieza o sistema —foundation, components, layout, brand, accessibility— y las materializa como DDRs estructurados, **organizados por surface**, que Claude consultará en futuras sesiones vía `.matecito-ai/ddr/INDEX.md`.

El objetivo es que las decisiones queden **registradas y verificables contra Figma**, no implícitas en la cabeza del diseñador. Eso le permite a Claude (y a cualquier nuevo diseñador) trabajar respetando el sistema sin volver a preguntarlo.

> **Nota sobre "DDR".** Usamos el término en sentido amplio: el catálogo cubre *decisiones* (con trade-offs reales), *convenciones* (acuerdos de estilo) y *políticas* (reglas verificables). El campo `type` de cada fase lo refleja. No todo lo que se captura es "arquitectura visual" en sentido estricto, pero todo merece quedar escrito, fechado y justificado — que es lo que aporta el formato DDR.

---

## Qué es (y qué no es) un DDR

El concepto canónico —qué cuenta como DDR y qué no, y la diferencia entre un DDR en borrador (inferido) y uno aceptado— vive en `~/.claude/references/ddr/README.md` (referencia consultable, agnóstica de cualquier flujo o skill). bootstrap **aplica** ese concepto; no lo redefine. La *estructura* concreta del archivo DDR está en `~/.claude/references/ddr/templates/ddr.md`.

---

## Arquitectura de esta skill

La skill está partida en **motor** y **datos**:

- **`SKILL.md` (este archivo) = el motor.** Define CÓMO se trata cualquier fase: el flujo, las reglas de UX, cómo se materializa un DDR, el modo update. Es estable; casi no cambia.
- **`concerns/` = el catálogo de fases (los datos), organizado por surface.** Cada fase vive en `concerns/<surface>/<slug>.md` con qué decide, qué preguntar y qué materializar. Cada surface tiene su `concerns/<surface>/INDEX.md` (detalle de sus concerns + criterio de pertenencia). `concerns/INDEX.md` es el índice raíz: mapa de surfaces + matriz de aplicabilidad por tipo de pieza. Esto crece con el tiempo (ver "Ratchet").
- **Salida = DDRs en `.matecito-ai/ddr/<surface>/` de la pieza objetivo.** Misma taxonomía de surfaces que el catálogo. No confundir con `concerns/`, que es el catálogo de la skill.

El motor lee `concerns/INDEX.md` **una sola vez** para armar la lista de fases relevantes, y recién después lee el archivo individual de cada fase que se va a tratar. Así no carga al contexto fases que no aplican.

---

## Surfaces canónicas (taxonomía fija)

La skill usa una taxonomía de surfaces **cerrada e impuesta por el motor** — idéntica en el catálogo interno (`concerns/<surface>/`) y en la salida (`.matecito-ai/ddr/<surface>/`). Esto garantiza que todas las piezas del equipo se vean igual y que ningún tema quede sin un casillero claro.

**Activas** (tienen concerns hoy):
`foundation` · `components` · `layout` · `brand` · `accessibility`

El significado de cada surface y su criterio de pertenencia están en `concerns/<surface>/INDEX.md`. **No inventés surfaces nuevas por pieza**: si de verdad falta una, es una decisión de catálogo (agregarla al motor y a `concerns/INDEX.md`), nunca improvisada en un repo.

Cada fase declara su surface en el frontmatter (`domain:` — el campo se llama así por compatibilidad del formato de concern; su valor es una surface). En la salida, el DDR de una fase se escribe en `.matecito-ai/ddr/<surface>/<name>.md`.

---

## Cuándo correr esta skill

- Sistema de diseño / marca nuevo (greenfield) sin `.matecito-ai/ddr/`
- Pieza existente cuyo sistema visual el usuario quiere "ordenar"
- El usuario pide explícitamente revisar/actualizar decisiones de diseño
- Hay DDRs `Inferred` que dejó `design-decisions-mine` y el usuario quiere ratificarlos
- Detectás que vas a fijar paleta/tipografía/grid/tokens y no hay sistema documentado

Si `.matecito-ai/ddr/INDEX.md` ya existe con contenido: **NO rehagas todo**. Andá al modo `update` (final del documento).

---

## Reglas del motor (aplican a TODAS las fases)

Estas reglas son la diferencia entre una skill que la gente usa y una que abandona en el tercer turno.

**Una pregunta por turno.** Nunca dumpees una lista de 8 preguntas. Hacé una, esperá la respuesta, leéla, recién ahí formulá la siguiente.

**Opciones concretas con default sugerido.** Mal: "¿qué paleta querés?". Bien: "¿Primario + neutros + semánticos (default para `app-ui`), solo marca + neutros (default para `marketing-asset`), o roles completos con secundario? Para tu tipo de pieza te recomiendo X porque Y." Los defaults y las opciones de cada fase están en su archivo `concerns/<surface>/<slug>.md`.

**Siempre incluí "no sé, recomendame".** Mucha gente no tiene opinión formada todavía. Cuando elijan esa opción, proponé con justificación de 2 líneas y pedí confirmación.

**Una línea de "por qué importa" antes de cada pregunta.** Usá el "Qué decide" del archivo de la fase. Sin sermones.

**Adaptate al contexto.** Si la descripción de la pieza, el archivo Figma conectado o el pre-flight ya respondieron algo, no lo vuelvas a preguntar. Saltá lo que no aplique, pero **nunca omitas en silencio** (ver siguiente regla).

**Nunca omitas en silencio.** Si una fase se salta — por elección del usuario, por atajo, o porque "no aplica" — queda igualmente registrada con una **razón breve de 1-2 líneas**. Una omisión sin justificación es una decisión perdida. El registro depende del status: `Not Applicable` se anota como fila en el INDEX de su surface (no genera archivo); `Pending` y `Deferred` se materializan como DDR-archivo (llevan trigger/condición y el modo update los rellena). Ver "Cómo manejar fases omitidas".

**Permití aplazar explícitamente.** Cualquier fase puede quedar `Pending` con la razón ("definimos el dark mode cuando llegue el rediseño de la app"). Mejor un DDR honesto con "pendiente + por qué" que una decisión inventada.

**Las reglas verificables son valores concretos contra Figma.** Cada DDR `Accepted` lleva reglas chequeables con valores —hex, ratio, escala, px, nombres de tokens— no adjetivos vagos. El mecanismo va al inicio de cada regla (`[tool: figma]`, `[tool: contrast]`, `[manual]`). Esto es lo que hace el sistema verificable por `design-decisions-mine` y por `design-verify`.

---

## Pre-flight (siempre primero)

Antes de la primera pregunta, inspeccioná la pieza objetivo:

```bash
ls -la
test -d .matecito-ai/ddr && echo "--- DDRs existentes (por surface) ---" && find .matecito-ai/ddr -name '*.md' | sort
```

Y si hay un archivo Figma conectado vía el MCP figma, leé sus styles / components existentes (`mcp__figma__*`) como punto de partida — un sistema ya empezado en Figma responde varias preguntas sin tener que hacerlas.

Con eso ya sabés:
- Si hay decisiones previas (`.matecito-ai/ddr/INDEX.md` existe → modo update; DDRs `Inferred` presentes → hay borradores para ratificar)
- Si la pieza es greenfield o ya tiene un sistema visual empezado en Figma
- Qué tokens / componentes ya existen (para inferir tipo de pieza y defaults)

> **Canva queda fuera.** Canva no expone tokens legibles, así que no sirve como fuente de evidencia. Si la pieza vive solo en Canva, las preguntas se contestan a mano; no hay reglas `[tool: figma]` que chequear.

---

## El flujo

### 1. Descripción de la pieza (entrada conversacional)

Una sola pregunta abierta, sin interrogatorio:

> Contame a grandes rasgos qué estás diseñando: qué tipo de pieza es (una app, una landing, un sistema de marca, un asset de marketing), qué tan importante es la consistencia y la accesibilidad, y si ya hay algo armado en Figma. No hace falta detalle, solo una idea para arrancar.

### 2. Inferencia + recomendación

De la descripción + el pre-flight, inferí (sin re-preguntar lo que ya quedó claro):

- **Tipo de pieza** → mapealo a un token de `concerns/INDEX.md` (`landing`, `app-ui`, `brand-system`, `marketing-asset`). Si es ambiguo entre dos, hacé UNA pregunta puntual.
- **Estado del Figma** → del pre-flight: greenfield vs. sistema empezado. Lo empezado provee defaults.
- **"Knobs" de intensidad** del lenguaje de la descripción:
  - Menciona accesibilidad / público amplio / requisitos WCAG → subí `contrast-target` a Requerido.
  - Menciona MVP / "rápido" / one-off → mantené lo esencial (foundation básico), ofrecé el resto como opcional.
  - Menciona sistema reutilizable / equipo / escala → incluí `design-tokens`, `component-library`.

Leé `concerns/INDEX.md` UNA vez. Armá el set: **Requerido(token) + Recomendado(token)**, ajustado por los knobs. Presentalo **agrupado por surface** (el mismo orden del índice raíz):

> Por lo que contás, esto parece un **[tipo de pieza]**. Te recomiendo estas fases, agrupadas por surface:
> - **`foundation`:** `color-palette` (por qué), `type-scale` (por qué)…
> - **`components`:** `component-library` (por qué)…
> - **`accessibility`:** …
> - …
>
> Quedan afuera por ahora (no parecen aplicar): `fase-x`, `fase-y`.
> Surfaces sin fases en esta pieza: `layout`…
>
> ¿Confirmás el tipo de pieza que detecté? [solo si quedó alguna duda]

### 3. Ajuste del set

> ¿Querés sacar alguna de estas, o agregar otras del catálogo?

Mostrá qué más hay disponible para sumar (las fases del catálogo no incluidas). Permití también **fase custom** (ver "Fase custom"). Lo que el usuario saque del set recomendado queda igual registrado — `Not Applicable` como fila en el INDEX de la surface, `Pending`/`Deferred` como DDR-archivo con razón — nunca hueco silencioso.

### 4. Recorrido de fases

Por cada fase del set final, seguí el procedimiento de "Cómo tratar una fase".

### 5. Materialización

Cuando se recorrieron todas, materializá (ver "Materialización").

---

## Cómo tratar una fase

Este es el procedimiento genérico del motor. Vale para cualquier fase, sea del catálogo o custom:

1. **Leé `concerns/<surface>/<slug>.md`** (solo cuando vas a tratar esa fase; la surface sale de la matriz del índice raíz o del campo `domain` del frontmatter).
2. Mostrá su **"Qué decide"** como la línea de "por qué importa".
3. Hacé sus **preguntas, una por turno**, en el orden del archivo. Para cada una: ofrecé las opciones con el default marcado e incluí "no sé, recomendame".
4. Si el archivo tiene **"Notas de lógica (para el motor)"**, aplicalas: defaults según tipo de pieza, preguntas condicionales, propuestas según respuestas de fases previas.
5. **Confirmá** la decisión antes de seguir.
5b. **Si la decisión corresponde a un principle canónico** del catálogo en `~/.claude/references/design-principles/` (típicamente fases de las surfaces `foundation`, `components`, `layout`), preguntá UNA vez cuál principle aplica y registralo en el DDR como `**Applied principle:** <Nombre> — <1 línea de por qué>`. No fuerces: si la decisión no mapea a un principle (ej. una convención de naming de tokens, una política de tono de voz), omití este paso. El catálogo se consulta por nombre, sin link en el DDR.
6. **Materializá el DDR** en `.matecito-ai/ddr/<surface>/<name>.md`, con el `type` del frontmatter en su encabezado, según la sección "Qué materializar" del archivo. Las **Reglas verificables** salen con valores concretos contra Figma y su mecanismo (`[tool: figma]` / `[tool: contrast]` / `[manual]`).

Si la fase estaba recomendada pero el usuario la sacó, o no aplica: no la trates, pero dejá su registro — `Not Applicable` como fila en el INDEX de la surface; `Pending`/`Deferred` como DDR-archivo con su razón. Ver "Cómo manejar fases omitidas".

---

## Cómo manejar fases omitidas

Cuando el usuario saca una fase, dice "no aplica" o "lo decidimos después": **no la trates, pero dejá registrado el motivo.** Dónde queda el registro depende del status (ver "Dónde se registra cada status").

No preguntes por concern. Juntá todas las fases del set recomendado que quedaron afuera y clasificá los motivos **en una sola pasada, en bloque**: para las que salen por tipo de pieza proponé la razón templada ("no aplica a un {tipo de pieza}") y confirmá en conjunto; abrí pregunta puntual solo por las que el usuario saca contra la recomendación de la matriz. Los concerns que la matriz nunca recomendó para este tipo de pieza no se enumeran uno por uno — su "no aplica" ya está implícito en la matriz.

Clasificá el motivo:

1. **No aplica al tipo de pieza** → `Not Applicable`. Ej: "Es un asset de marketing one-off, no necesita biblioteca de componentes."
2. **Lo decidimos después** → `Pending` (con trigger esperado opcional). Ej: "Definimos dark mode cuando llegue el rediseño."
3. **No me interesa documentarlo / ad-hoc** → `Not Applicable` con motivo honesto.
4. **Otra razón** → el status que aplique, motivo libre.

### Status posibles

Conjunto cerrado, así el INDEX y las revisiones futuras son consistentes:

- **`Accepted`** — Decisión tomada y vigente.
- **`Pending`** — Sabemos que hay que decidirlo, todavía no es el momento. Incluye trigger ("cuando…") si se conoce.
- **`Not Applicable`** — Decisión consciente de que este tema no aplica. Lleva razón obligatoria.
- **`Deferred`** — Postergado deliberadamente con fecha o condición de revisión.
- **`Superseded`** — Reemplazado por otro DDR. Lleva referencia al que lo sustituye.
- **`Inferred`** — Borrador minado del archivo Figma por `design-decisions-mine`: tiene el QUÉ y la evidencia, pero el PORQUÉ vacío. No es una decisión ratificada hasta que el modo update lo entrevista y lo pasa a `Accepted`.

### Dónde se registra cada status

- **`Not Applicable`** → fila en el INDEX de la surface (`.matecito-ai/ddr/<surface>/INDEX.md`, sección "No aplican"). **No genera archivo propio.** Si la surface entera queda sin ningún DDR-archivo (`Accepted`/`Pending`/`Deferred`/`Inferred`), no se crea carpeta: la surface se lista en el INDEX raíz, sección "Surfaces sin uso" (ver Materialización).
- **`Pending` / `Deferred`** → DDR-archivo individual en la carpeta de la surface, con su trigger o condición de revisión. El modo update los resuelve rellenando contenido y pasándolos a `Accepted`.
- **`Accepted`** → DDR-archivo individual con contenido completo, incluidas las reglas verificables.
- **`Inferred`** → DDR-archivo individual con header + `## Evidencia (inferida)` + `## Alcance` (lo escribe `design-decisions-mine`, no bootstrap). El modo update lo ratifica.

---

## Fase custom

Si el usuario quiere un tema que no está en el catálogo:

1. Tratalo con el procedimiento genérico, haciéndole 2-3 preguntas para extraer qué decide, opciones y qué materializar.
2. **Asignale una surface canónica.** Mirá el "criterio de pertenencia" en cada `concerns/<surface>/INDEX.md` para decidir dónde encaja. No inventés una surface nueva. Si genuinamente no encaja en ninguna, es señal de que falta una surface en la taxonomía — eso es una decisión de catálogo, avisале al usuario, no lo resuelvas en el repo.
3. **Asignale un `type`** (`decision` / `convention` / `policy`).
4. Materializá el DDR en `.matecito-ai/ddr/<surface>/<slug>.md`. Una fase custom es **siempre solo para esta pieza**: no toques el catálogo `concerns/` (es read-only, se deploya desde el repo matecito-ai). Si el concern merece sumarse al catálogo para todas las piezas, eso se hace editando `payload/domains/design/skills/.../concerns/` en el repo matecito-ai (ver "Ratchet"), no desde acá.

---

## Materialización

### Paso 1: Resumir y confirmar

Antes de escribir nada, mostrá un resumen completo de todas las decisiones, agrupadas por fase, con su status. Pedí confirmación final. Permití editar cualquier respuesta.

### Paso 2: Estructura de archivos a generar

Los DDRs de salida son **slug-based** (sin prefijos numéricos) y van **agrupados por surface en subcarpetas**, con un índice por surface más un índice raíz. Misma taxonomía de surfaces que el catálogo interno.

```
<root>/
└── .matecito-ai/
    └── ddr/
        ├── INDEX.md                       # índice RAÍZ: explica cada surface + cómo navegar
        ├── <surface>/                     # una carpeta por surface CON al menos un DDR
        │   ├── INDEX.md                   # índice de la SURFACE: lista sus DDRs + criterio
        │   ├── <slug>.md                  # un DDR por fase tratada (ej: color-palette.md)
        │   └── ...
        └── ...                            # otras surfaces
```

Reglas de la estructura:

- **Solo se crean carpetas de surfaces que tienen al menos un DDR-archivo** (`Accepted`, `Pending`, `Deferred` o `Inferred`). No repliques las 5 surfaces en cada pieza — la salida refleja lo que la pieza realmente definió.
- **Qué genera archivo y qué no:**
  - `Accepted` → DDR-archivo con contenido completo (incluidas reglas verificables).
  - `Pending` / `Deferred` → DDR-archivo con su trigger/condición (el modo update los rellena).
  - `Not Applicable` → **fila en el INDEX de la surface** (sección "No aplican"), sin archivo propio.
  - El nombre de archivo es el `name` del concern (el slug); la surface es el campo `domain` del frontmatter.
- **Surface sin ningún DDR-archivo:** si todas sus fases quedaron `Not Applicable`, no se crea carpeta; la surface se lista en el INDEX raíz (sección "Surfaces sin uso") con una razón de 1 línea.
- **Qué se lista como N/A:** solo las fases que la matriz daba como Requerido/Recomendado para el tipo de pieza y se sacaron. Los concerns que la matriz nunca recomendó para este tipo no se enumeran — su "no aplica" ya está en la matriz.
- **Dos niveles de índice:** el raíz (`ddr/INDEX.md`) enruta por surface y lista las surfaces sin uso; cada surface (`ddr/<surface>/INDEX.md`) lista sus DDRs y sus N/A.

### Paso 3: Templates

Los templates de estructura de DDR son el **contrato canónico** y viven en la referencia `~/.claude/references/ddr/templates/` (índice en `~/.claude/references/ddr/templates/INDEX.md`). No se duplican acá: antes de materializar, leé el template del artefacto que vas a escribir.

| Artefacto | Template |
|---|---|
| DDR individual (`<surface>/<slug>.md`) | `~/.claude/references/ddr/templates/ddr.md` |
| Índice raíz (`ddr/INDEX.md`) | `~/.claude/references/ddr/templates/index-root.md` |
| Índice de surface (`ddr/<surface>/INDEX.md`) | `~/.claude/references/ddr/templates/index-surface.md` |

Notas del contrato del DDR (también en `~/.claude/references/ddr/templates/ddr.md`): **no hay sección `Historial`** (lo lleva git; la evolución se ve en la cadena de `Superseded`); **links entre DDRs** — dentro de la misma surface ruta corta (`<slug>.md`), entre surfaces ruta relativa (`../<surface>/<slug>.md`).

> A diferencia de development, design NO tiene catálogo de tecnologías (`tech/`) ni escribe un `CLAUDE.md` raíz.

### Paso 4: Escribir y reportar

1. Para cada surface con al menos un DDR-archivo (`Accepted`/`Pending`/`Deferred`): `mkdir -p .matecito-ai/ddr/<surface>`.
2. Escribir `.matecito-ai/ddr/INDEX.md` (índice raíz) listando las surfaces con DDR-archivo y, en su sección "Surfaces sin uso", las surfaces cuyas fases quedaron todas `Not Applicable`.
3. Escribir `.matecito-ai/ddr/<surface>/INDEX.md` para cada surface usada: la tabla de DDRs y la sección "No aplican" con las fases `Not Applicable` de la surface y su razón.
4. Escribir los DDR-archivo de las fases con contenido: `Accepted` completo (con reglas verificables); `Pending`/`Deferred` con su trigger/condición. Los `Not Applicable` no generan archivo — quedan como fila en el INDEX de la surface (o del raíz si la surface quedó sin uso).
5. Reportar al usuario:
   - Lista de archivos creados (path completo), **agrupada por surface**
   - Resumen de 1 línea por DDR-archivo, con su status entre corchetes (`[Accepted]`, `[Pending]`, `[Deferred]`) y su `type`
   - Conteo de `Not Applicable` por surface (viven en los INDEX), no uno por uno
   - **Lista separada de DDRs `Pending`/`Deferred` con su trigger**, así sabe qué quedó por decidir
   - Sugerencia de commitear estos archivos al repo

---

## Modo update (cuando `.matecito-ai/ddr/INDEX.md` ya existe)

1. **Leé el índice raíz, los índices de surface y los DDRs** existentes (`find .matecito-ai/ddr -name '*.md'`).
2. **Mostrá un resumen agrupado por surface y, dentro de cada una, por status:** `Accepted`, `Pending` (con trigger), `Deferred`, `Not Applicable` (con razón), `Inferred` (borrador minado de Figma, sin porqué).
3. **Preguntá si algún `Pending` o `Deferred` ya está listo para resolverse.** Sin esto, los "lo decidimos después" se pierden.
3b. **Ratificá los `Inferred` (borradores minados de Figma).** Un `Inferred` tiene el QUÉ y la evidencia (`## Evidencia (inferida)` + `## Alcance`), pero el PORQUÉ vacío — es un candidato sin ratificar (ver `~/.claude/references/ddr/README.md`). Por cada uno, ofrecé ratificarlo: entrevistá por el porqué (Contexto, Decisión razonada, Alternativas, Consecuencias), llená esas secciones, agregá las **Reglas verificables** con valores concretos contra Figma, **descartá la sección `## Evidencia (inferida)`** (es transitoria), y cambiá `Status: Inferred → Accepted`. Si el usuario no quiere ratificarlo ahora, queda `Inferred` (no se pierde). NUNCA promuevas un `Inferred` sin la entrevista del porqué — eso es lo que lo convierte en decisión.
4. **Ratchet — barré el catálogo:** leé `concerns/INDEX.md`, listá las fases relevantes al tipo de pieza que **no tengan DDR todavía** (típicamente fases nuevas agregadas al catálogo desde la última corrida) y ofrecé tratarlas ahora. Mostralas agrupadas por surface, incluyendo si caen en una surface que la pieza todavía no usa (esa carpeta se crea recién al materializar el primer DDR). Esta es la forma de que los temas agregados al catálogo lleguen a piezas viejas.
5. **Después preguntá qué más quiere hacer:**
   - **Resolver un Pending/Deferred** → recorrer las preguntas de esa fase, cambiar Status a `Accepted`, llenar contenido.
   - **Ratificar un `Inferred`** → entrevistar por el porqué, llenar Contexto/Decisión/Alternativas/Consecuencias + Reglas verificables, descartar `## Evidencia (inferida)`, `Status → Accepted`.
   - **Actualizar una decisión (cambio menor)** → editar el DDR. Git lleva el historial.
   - **Cambiar una decisión (cambio de fondo)** → crear DDR nuevo en la misma surface, marcar el viejo `Superseded` con link al nuevo. No editar la decisión vieja en el lugar.
   - **Agregar una decisión nueva** no cubierta → crear DDR en su surface + fila en el índice de esa surface (y en el raíz si la surface es nueva en la pieza).
   - **Cambiar un `Not Applicable` a `Pending`/`Accepted`** → el contexto de la pieza cambió (ej: el asset one-off creció a sistema reutilizable y ahora sí hay biblioteca de componentes). Sacá la fila de la sección "No aplican" del INDEX de la surface (o "Surfaces sin uso" del raíz) y creá el DDR-archivo con el nuevo status y contenido; creá la carpeta de la surface si no existía.
   - **Rehacer todo desde cero** → confirmación doble. Antes de sobrescribir, mover el directorio a `.matecito-ai/ddr.old.<timestamp>/`.
6. Para actualizar/agregar, recorré solo las fases relevantes — no rehagas todo el cuestionario.
7. **Después de cualquier cambio, mantené los índices coherentes:** actualizá el índice de la surface afectada y, si agregaste o vaciaste una surface, el índice raíz.

---

## Ratchet: agregar una fase al catálogo

El valor de largo plazo de la skill es que **nunca se vuelva a olvidar un tema**. Cuando aparece un concern que no estaba:

1. Determiná a qué **surface canónica** pertenece (consultá el "criterio de pertenencia" en `concerns/<surface>/INDEX.md`).
2. Creá `concerns/<surface>/<slug>.md` con el formato estándar (mirá `concerns/foundation/color-palette.md` para una fase `deep` y `concerns/foundation/spacing-grid.md` para una `light`). Incluí en el frontmatter `domain` (la surface) y `type`.
3. Sumá la fila al `concerns/<surface>/INDEX.md` y a la matriz de `concerns/INDEX.md`, con la aplicabilidad por tipo de pieza (Requerido/Recomendado).

Desde ese momento, todo bootstrap futuro lo considera, y el modo update lo ofrece a piezas viejas (paso 4 de update). El catálogo se sembró de taxonomías externas (W3C Design Tokens, Atomic Design, WCAG 2.x, Material Design) para nacer casi completo y solo crecer.

---

## Anti-patterns que esta skill evita

- ❌ Tirar todas las preguntas en un solo turno → la gente abandona.
- ❌ Forzar una biblioteca de componentes completa en un asset de marketing one-off → adaptar el set de fases al tipo de pieza.
- ❌ Saltar una fase sin documentar el motivo → siempre dejar registro: `Not Applicable` como fila en el INDEX de la surface; `Pending`/`Deferred` como DDR-archivo + razón.
- ❌ Crear un DDR-archivo por cada `Not Applicable` → los N/A viven como fila en el INDEX de la surface, no como archivo; solo se justifican las desviaciones de la matriz.
- ❌ Preguntar el motivo de cada N/A por separado → clasificar en bloque, una sola pasada, con razón templada por tipo de pieza.
- ❌ Confundir "no decidido aún" (`Pending`) con "decidido que no aplica" (`Not Applicable`) → son status distintos.
- ❌ Editar una decisión de fondo en el lugar → para cambios de decisión, supersede (DDR nuevo + viejo `Superseded`). Cambios menores sí se editan; git lleva el historial.
- ❌ Promover un `Inferred` a `Accepted` sin entrevistar el porqué → eso es lo que lo convierte en decisión; un Inferred promovido a dedo es una decisión inventada.
- ❌ Dejar la sección `## Evidencia (inferida)` en un DDR ya ratificado → es transitoria; al pasar a `Accepted` se descarta (git conserva la traza).
- ❌ Mantener una tabla `Historial` manual → es redundante con git y se pudre.
- ❌ Inventar reglas no discutidas con el usuario en la materialización → todo lo que va al DDR fue confirmado.
- ❌ Reglas vagas tipo "que la paleta sea armoniosa" → siempre verificable: hex, ratio, escala, px, nombres de tokens, con su mecanismo (`[tool: figma]`/`[tool: contrast]`/`[manual]`).
- ❌ Inferir decisiones de un mockup con valores sueltos sin tokens → eso es señal débil; sin styles nombrados no hay regla `[tool: figma]` que chequear.
- ❌ Asumir el sistema en lugar de leerlo del Figma conectado en pre-flight → leer styles/components primero.
- ❌ Leer todo el catálogo `concerns/` de una → leer `INDEX.md` (raíz) para seleccionar, y cada `concerns/<surface>/<slug>.md` solo cuando se trata esa fase.
- ❌ Inventar una surface nueva en una pieza → la taxonomía es fija e impuesta por el motor; una surface nueva es decisión de catálogo, no de pieza.
- ❌ Replicar las 5 surfaces en la salida → en `.matecito-ai/ddr/` solo se crean las carpetas de surfaces con al menos un DDR-archivo; las surfaces solo-N/A se listan en el INDEX raíz.
- ❌ Dejar índices desincronizados tras un cambio → actualizá el índice de la surface afectada y el raíz si corresponde.
- ❌ Minar/ratificar decisiones de Canva → no expone tokens legibles; reportarlo y no inferir.
- ❌ En modo update, no preguntar por los `Pending`/`Deferred`, no ratificar los `Inferred`, ni barrer el catálogo por fases nuevas → es como se pierden las decisiones aplazadas, los borradores minados y las fases agregadas.

---

## Recordatorio final

El valor de esta skill no está en las preguntas — está en que las decisiones de diseño queden **escritas, verificables contra Figma y mantenidas**. Si las preguntas son geniales pero los DDRs salen vagos, fallamos. Si los DDRs son específicos y chequeables (hex, ratio, escala, nombres de tokens), Claude (y cualquier diseñador) puede trabajar respetando el sistema sin volver a preguntar.

Escribí los DDRs con la misma claridad con la que le explicarías el sistema visual a un diseñador nuevo el primer día.
