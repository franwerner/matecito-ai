---
name: development-decisions-bootstrap
description: Entrevista interactiva por fases para capturar las decisiones de ingeniería de un proyecto (arquitectura, convenciones y políticas) al iniciarlo y materializarlas como EDRs organizados por dominio en .matecito-ai/edr/<dominio>/. Usá esta skill SIEMPRE que el usuario inicie un proyecto nuevo, clone un repo vacío, mencione "arrancar un proyecto / empezar un proyecto / setup inicial", pida ayuda con la arquitectura inicial, hable de "definir capas, estructura, convenciones", quiera revisar o actualizar decisiones de arquitectura existentes, o cuando detectes un repo sin .matecito-ai/edr/ ni CLAUDE.md y el usuario esté por escribir código que toca estructura. También dispará si el usuario menciona "EDR", "decisión arquitectónica", "convenciones del proyecto", "manejo de errores", "capas", "acoplamiento", "estructura de carpetas".
---

# Development Decisions Bootstrap

Entrevista al usuario para capturar las decisiones de ingeniería del proyecto —arquitectura, convenciones y políticas— y las materializa como EDRs estructurados, **organizados por dominio**, que Claude consultará en futuras sesiones vía `.matecito-ai/edr/INDEX.md`.

El objetivo es que las decisiones queden **registradas y verificables**, no implícitas en la cabeza del autor. Eso le permite a Claude (y a cualquier nuevo dev) trabajar respetando las convenciones sin volver a preguntarlas.

> **Nota sobre "EDR".** Usamos el término en sentido amplio: el catálogo cubre *decisiones* (con trade-offs reales), *convenciones* (acuerdos de estilo) y *políticas* (reglas verificables). El campo `type` de cada fase lo refleja. No todo lo que se captura es "arquitectura" en sentido estricto, pero todo merece quedar escrito, fechado y justificado — que es lo que aporta el formato EDR.

---

## Qué es (y qué no es) un EDR

El concepto canónico —qué cuenta como EDR y qué no, y la diferencia entre un EDR en borrador (inferido) y uno aceptado— vive en `~/.claude/references/edr/README.md` (referencia consultable, agnóstica de cualquier flujo o skill). bootstrap **aplica** ese concepto; no lo redefine. La *estructura* concreta del archivo EDR está en `~/.claude/references/edr/templates/edr.md`.

---

## Arquitectura de esta skill

La skill está partida en **motor** y **datos**:

- **`SKILL.md` (este archivo) = el motor.** Define CÓMO se trata cualquier fase: el flujo, las reglas de UX, cómo se materializa un EDR, el modo update. Es estable; casi no cambia.
- **`concerns/` = el catálogo de fases (los datos), organizado por dominio.** Cada fase vive en `concerns/<dominio>/<slug>.md` con qué decide, qué preguntar y qué materializar. Cada dominio tiene su `concerns/<dominio>/INDEX.md` (detalle de sus concerns + criterio de pertenencia). `concerns/INDEX.md` es el índice raíz: mapa de dominios + matriz de aplicabilidad por tipo de proyecto. Esto crece con el tiempo (ver "Ratchet").
- **Salida = EDRs en `.matecito-ai/edr/<dominio>/` del proyecto objetivo.** Misma taxonomía de dominios que el catálogo. No confundir con `concerns/`, que es el catálogo de la skill.

El motor lee `concerns/INDEX.md` **una sola vez** para armar la lista de fases relevantes, y recién después lee el archivo individual de cada fase que se va a tratar. Así no carga al contexto fases que no aplican.

---

## Dominios canónicos (taxonomía fija)

La skill usa una taxonomía de dominios **cerrada e impuesta por el motor** — idéntica en el catálogo interno (`concerns/<dominio>/`) y en la salida (`.matecito-ai/edr/<dominio>/`). Esto garantiza que todos los repos del equipo se vean igual y que ningún tema quede sin un casillero claro.

**Activos** (tienen concerns hoy):
`structure` · `runtime` · `data` · `observability` · `security` · `contracts` · `delivery` · `frontend` · `quality`

**Reservados** (casilleros válidos, sin concerns todavía — se pueblan vía ratchet):
`lifecycle` · `integration` · `privacy` · `release` · `domain-logic` · `compliance` · `ux-product`

El significado de cada dominio y su criterio de pertenencia están en `concerns/<dominio>/INDEX.md`. **No inventés dominios nuevos por proyecto**: si de verdad falta uno, es una decisión de catálogo (agregarlo al motor y a `concerns/INDEX.md`), nunca improvisado en un repo.

Cada fase declara su dominio en el frontmatter (`domain:`). En la salida, el EDR de una fase se escribe en `.matecito-ai/edr/<domain>/<name>.md`.

---

## Cuándo correr esta skill

- Proyecto nuevo (greenfield) sin `.matecito-ai/edr/` ni `CLAUDE.md`
- Repo existente que el usuario quiere "ordenar"
- El usuario pide explícitamente revisar/actualizar decisiones de arquitectura
- Detectás que vas a tocar estructura/capas/auth/errores y no hay convenciones documentadas

Si `.matecito-ai/edr/INDEX.md` ya existe con contenido: **NO rehagas todo**. Andá al modo `update` (final del documento).

---

## Reglas del motor (aplican a TODAS las fases)

Estas reglas son la diferencia entre una skill que la gente usa y una que abandona en el tercer turno.

**Una pregunta por turno.** Nunca dumpees una lista de 8 preguntas. Hacé una, esperá la respuesta, leéla, recién ahí formulá la siguiente.

**Opciones concretas con default sugerido.** Mal: "¿cómo manejás errores?". Bien: "¿Excepciones (default para Python/Java/Node), Result types (default para Rust/Go), o mix pragmático? Para tu stack te recomiendo X porque Y." Los defaults y las opciones de cada fase están en su archivo `concerns/<dominio>/<slug>.md`.

**Siempre incluí "no sé, recomendame".** Mucha gente no tiene opinión formada todavía. Cuando elijan esa opción, proponé con justificación de 2 líneas y pedí confirmación.

**Una línea de "por qué importa" antes de cada pregunta.** Usá el "Qué decide" del archivo de la fase. Sin sermones.

**Adaptate al contexto.** Si la descripción del proyecto o el pre-flight ya respondieron algo, no lo vuelvas a preguntar. Saltá lo que no aplique, pero **nunca omitas en silencio** (ver siguiente regla).

**Nunca omitas en silencio.** Si una fase se salta — por elección del usuario, por atajo, o porque "no aplica" — queda igualmente registrada con una **razón breve de 1-2 líneas**. Una omisión sin justificación es una decisión perdida. El registro depende del status: `Not Applicable` se anota como fila en el INDEX de su dominio (no genera archivo); `Pending` y `Deferred` se materializan como EDR-archivo (llevan trigger/condición y el modo update los rellena). Ver "Cómo manejar fases omitidas".

**Permití aplazar explícitamente.** Cualquier fase puede quedar `Pending` con la razón ("lo definimos cuando llegue el feature de pagos"). Mejor un EDR honesto con "pendiente + por qué" que una decisión inventada.

**Registrá tecnologías a medida que aparecen.** Cada vez que el usuario menciona o confirma una tecnología concreta, creá su mini-EDR en `.matecito-ai/edr/tech/<nombre>.md` en el momento. No esperes a la materialización final. Ver "Catálogo de tecnologías".

---

## Migración de legacy (`.matecito-ai/adr/` → `.matecito-ai/edr/`)

El artefacto se llamaba **ADR** y vivía en `.matecito-ai/adr/`; ahora es **EDR** en `.matecito-ai/edr/`. Si el proyecto viene de una versión anterior, detectalo **antes del pre-flight** y **ofrecé** migrar (nunca en silencio — son records del usuario):

```bash
if [ -d .matecito-ai/adr ] && [ ! -d .matecito-ai/edr ]; then echo "LEGACY: existe .matecito-ai/adr/ (formato ADR viejo) — ofrecer migración"; fi
```

Si existe, proponé al usuario y confirmá antes de tocar:
1. Renombrar la carpeta: `git mv .matecito-ai/adr .matecito-ai/edr` (o `mv` si no es git).
2. Actualizar el término dentro de los archivos migrados (`# ADR —` → `# EDR —`, "ADR"→"EDR" en los INDEX y records) y el pointer del `CLAUDE.md` del proyecto (`.matecito-ai/adr/` → `.matecito-ai/edr/`).

Tras migrar (o si el usuario prefiere no migrar todavía), seguí con el pre-flight. Si NO migra, el resto de la skill no encontrará records bajo `.matecito-ai/edr/` y tratará el proyecto como sin decisiones previas.

---

## Pre-flight (siempre primero)

Antes de la primera pregunta, inspeccioná el repo objetivo:

```bash
ls -la
test -f CLAUDE.md && echo "--- CLAUDE.md existe ---" && cat CLAUDE.md
test -d .matecito-ai/edr && echo "--- EDRs existentes (por dominio) ---" && find .matecito-ai/edr -name '*.md' | sort
test -d .matecito-ai/edr/tech && echo "--- Tech ya registrada ---" && ls .matecito-ai/edr/tech/
test -f package.json && echo "--- package.json ---" && head -50 package.json
test -f pyproject.toml && echo "--- pyproject.toml ---" && head -50 pyproject.toml
test -f go.mod && echo "--- go.mod ---" && cat go.mod
test -f Cargo.toml && echo "--- Cargo.toml ---" && head -30 Cargo.toml
test -f composer.json && echo "--- composer.json ---" && head -30 composer.json
test -f Gemfile && echo "--- Gemfile ---" && cat Gemfile
```

Con eso ya sabés:
- Si hay decisiones previas (`.matecito-ai/edr/INDEX.md` existe → modo update)
- Stack y framework principal (para inferir tipo y defaults)
- Si el repo es greenfield o tiene código existente

---

## El flujo

### 1. Descripción del proyecto (entrada conversacional)

Una sola pregunta abierta, sin interrogatorio:

> Contame a grandes rasgos de qué trata el proyecto: qué hace, qué tan importante es la seguridad, qué convenciones te importan, y qué stack pensás usar. No hace falta detalle, solo una idea para arrancar.

### 2. Inferencia + recomendación

De la descripción + el pre-flight, inferí (sin re-preguntar lo que ya quedó claro):

- **Tipo de proyecto** → mapealo a un token de `concerns/INDEX.md` (`api-rest`, `api-graphql`, `cli`, `libreria`, `web-spa`, `web-ssr`, `microservicio`, `monolito-modular`, `script`). Si es ambiguo entre dos, hacé UNA pregunta puntual.
- **Stack** → del pre-flight o de lo que mencionó. Si no se detectó y es crítico, preguntá.
- **Tamaño de equipo y punto de partida** → inferilos de la descripción y el pre-flight (greenfield si el repo está vacío; código existente / migración si ya hay código). Escalan la formalidad: solo / equipo chico → mantené lo mínimo; equipo grande → subí las fases de convención (`folder-structure`, `ci-quality-gates`, `arch-enforcement`).
- **"Knobs" de intensidad** del lenguaje de la descripción:
  - Menciona seguridad alta / datos sensibles / usuarios externos → subí las fases del dominio `security` (`authorization`, `input-validation`, `rate-limiting`, `cors`, `secrets-management`, `dependency-scanning`) a Requerido.
  - Menciona MVP / prototipo / "rápido" → mantené lo esencial, ofrecé el resto como opcional.
  - Menciona convenciones estrictas / equipo grande → incluí `folder-structure`, `ci-quality-gates`, `arch-enforcement`.

**Atajo para scripts mínimos.** Si es un `script` de una sola persona, proponé un set mínimo —saltá `architecture-style`, `layers-and-dependencies` e `inter-layer-communication`, y andá directo a `folder-structure` + lo esencial de `delivery`. Pedí permiso explícito para el atajo; las fases salteadas quedan registradas como `Not Applicable` con su razón.

Leé `concerns/INDEX.md` UNA vez. Armá el set: **Requerido(token) + Recomendado(token)**, ajustado por los knobs. Presentalo **agrupado por dominio** (el mismo orden del índice raíz):

> Por lo que contás, esto parece un **[tipo]**. Te recomiendo estas fases, agrupadas por dominio:
> - **`structure`:** `architecture-style` (por qué), `folder-structure` (por qué)…
> - **`runtime`:** `error-handling` (por qué)…
> - **`security`:** …
> - …
>
> Quedan afuera por ahora (no parecen aplicar): `fase-x`, `fase-y`.
> Dominios sin fases en este proyecto: `frontend`, `quality`…
>
> ¿Confirmás el tipo y el stack que detecté? [solo si quedó alguna duda]

### 3. Ajuste del set

> ¿Querés sacar alguna de estas, o agregar otras del catálogo?

Mostrá qué más hay disponible para sumar (las fases del catálogo no incluidas). Permití también **fase custom** (ver "Fase custom"). Lo que el usuario saque del set recomendado queda igual registrado — `Not Applicable` como fila en el INDEX del dominio, `Pending`/`Deferred` como EDR-archivo con razón — nunca hueco silencioso.

### 4. Recorrido de fases

Por cada fase del set final, seguí el procedimiento de "Cómo tratar una fase". Intercalá el registro de tecnologías cuando aparezcan.

### 5. Materialización

Cuando se recorrieron todas, materializá (ver "Materialización").

### 6. Validación (recomendada)

Al cerrar, ofrecé correr el validador `development-decisions-validate` en **contexto fresco** (como sub-agente), pasándole el tipo de proyecto y la lista de fases relevantes. Chequea coherencia entre EDRs, completitud y verificabilidad, y reporta con severidad. No modifica nada — los hallazgos los resuelve el usuario vía modo update. Es opcional pero recomendado: ojos frescos atrapan contradicciones que el flujo de la entrevista no ve.

---

## Cómo tratar una fase

Este es el procedimiento genérico del motor. Vale para cualquier fase, sea del catálogo o custom:

1. **Leé `concerns/<dominio>/<slug>.md`** (solo cuando vas a tratar esa fase; el dominio sale de la matriz del índice raíz o del campo `domain` del frontmatter).
2. Mostrá su **"Qué decide"** como la línea de "por qué importa".
3. Hacé sus **preguntas, una por turno**, en el orden del archivo. Para cada una: ofrecé las opciones con el default marcado e incluí "no sé, recomendame".
4. Si el archivo tiene **"Notas de lógica (para el motor)"**, aplicalas: defaults según stack, preguntas condicionales, propuestas según respuestas de fases previas.
5. **Confirmá** la decisión antes de seguir.
6. Si el archivo tiene **"Tech a registrar"** y se eligió una tecnología concreta, creá su mini-EDR en `tech/` en el momento (ver "Catálogo de tecnologías").
6b. **Si la decisión corresponde a un patrón canónico** del catálogo en `~/.claude/references/design-patterns/` (típicamente fases de los dominios `structure`, `runtime`, `data`), preguntá UNA vez cuál patrón aplica y registralo en el EDR como `**Applied pattern:** <Nombre> — <1 línea de por qué>`. No fuerces: si la decisión no es un patrón (ej. una convención de naming, una política de rate limiting), omití este paso. El catálogo se consulta por nombre, sin link en el EDR — el pointer a la ubicación está en el `CLAUDE.md` del proyecto.
7. **Materializá el EDR** en `.matecito-ai/edr/<dominio>/<name>.md`, con el `type` del frontmatter en su encabezado, según la sección "Qué materializar" del archivo.

Si la fase estaba recomendada pero el usuario la sacó, o no aplica: no la trates, pero dejá su registro — `Not Applicable` como fila en el INDEX del dominio; `Pending`/`Deferred` como EDR-archivo con su razón. Ver "Cómo manejar fases omitidas".

---

## Cómo manejar fases omitidas

Cuando el usuario saca una fase, dice "no aplica" o "lo decidimos después": **no la trates, pero dejá registrado el motivo.** Dónde queda el registro depende del status (ver "Dónde se registra cada status").

No preguntes por concern. Juntá todas las fases del set recomendado que quedaron afuera y clasificá los motivos **en una sola pasada, en bloque**: para las que salen por tipo de proyecto proponé la razón templada ("no aplica a un {tipo}") y confirmá en conjunto; abrí pregunta puntual solo por las que el usuario saca contra la recomendación de la matriz. Los concerns que la matriz nunca recomendó para este tipo no se enumeran uno por uno — su "no aplica" ya está implícito en la matriz.

Clasificá el motivo:

1. **No aplica al tipo de proyecto** → `Not Applicable`. Ej: "Es un script CLI sin red, no necesita auth."
2. **Lo decidimos después** → `Pending` (con trigger esperado opcional). Ej: "Definimos auth cuando llegue el milestone de usuarios públicos."
3. **No me interesa documentarlo / ad-hoc** → `Not Applicable` con motivo honesto.
4. **Otra razón** → el status que aplique, motivo libre.

### Status posibles

Conjunto cerrado, así el INDEX y las revisiones futuras son consistentes:

- **`Accepted`** — Decisión tomada y vigente.
- **`Pending`** — Sabemos que hay que decidirlo, todavía no es el momento. Incluye trigger ("cuando…") si se conoce.
- **`Not Applicable`** — Decisión consciente de que este tema no aplica. Lleva razón obligatoria.
- **`Deferred`** — Postergado deliberadamente con fecha o condición de revisión.
- **`Superseded`** — Reemplazado por otro EDR. Lleva referencia al que lo sustituye.

### Dónde se registra cada status

- **`Not Applicable`** → fila en el INDEX del dominio (`.matecito-ai/edr/<dominio>/INDEX.md`, sección "No aplican"). **No genera archivo propio.** Si el dominio entero queda sin ningún EDR-archivo (`Accepted`/`Pending`/`Deferred`), no se crea carpeta: el dominio se lista en el INDEX raíz, sección "Dominios sin uso" (ver Materialización).
- **`Pending` / `Deferred`** → EDR-archivo individual en la carpeta del dominio, con su trigger o condición de revisión. El modo update los resuelve rellenando contenido y pasándolos a `Accepted`.
- **`Accepted`** → EDR-archivo individual con contenido completo.

---

## Fase custom

Si el usuario quiere un tema que no está en el catálogo:

1. Tratalo con el procedimiento genérico, haciéndole 2-3 preguntas para extraer qué decide, opciones y qué materializar.
2. **Asignale un dominio canónico.** Mirá el "criterio de pertenencia" en cada `concerns/<dominio>/INDEX.md` para decidir dónde encaja (incluí los reservados: `lifecycle`, `integration`, `privacy`, `release`, `domain-logic`, `compliance`, `ux-product`). No inventés un dominio nuevo. Si genuinamente no encaja en ninguno, es señal de que falta un dominio en la taxonomía — eso es una decisión de catálogo, avisале al usuario, no lo resuelvas en el repo.
3. **Asignale un `type`** (`decision` / `convention` / `policy`).
4. Materializá el EDR en `.matecito-ai/edr/<dominio>/<slug>.md`. Una fase custom es **siempre solo para este proyecto**: no toques el catálogo `concerns/` (es read-only, se deploya desde el repo matecito-ai). Si el concern merece sumarse al catálogo para todos los proyectos, eso se hace editando `payload/skills/.../concerns/` en el repo matecito-ai (ver "Ratchet"), no desde acá.

---

## Catálogo de tecnologías (transversal a todas las fases)

Registro paralelo que se construye intercalado con la conversación. Cada vez que el usuario menciona o confirma una tecnología concreta, creás su mini-EDR.

### Cuándo crear un mini-EDR de tecnología

- El usuario nombra una lib/framework/herramienta ("usemos Postgres", "para tests pytest").
- El usuario elige una opción que implica una tecnología ("ORM" → preguntar cuál → registrar).
- Vos recomendás algo y lo acepta.
- Lo detectaste en pre-flight y el usuario confirma seguir usándola.

**No registres** versiones internas del lenguaje, dependencias transitivas, ni herramientas de build estándar (npm, pip) salvo que el usuario las haya elegido explícitamente sobre otra (ej: pnpm sobre npm sí).

### Flujo al detectar una tecnología

Tres preguntas rápidas (pueden ir en un turno):

1. **Versión.** Si el manifest la tiene, mostrala como default.
2. **Por qué (1-2 líneas).** Si no tiene una razón clara, sugerí una y pedí confirmación.
3. **Alternativas descartadas (1 línea).** 1-3 que se consideraron, o "ninguna evaluada" (información honesta).

Escribí el archivo y seguí con la fase. No detengas el flujo principal por esto.

### Estructura

```
.matecito-ai/edr/tech/
├── INDEX.md                  # tabla por categoría
├── python.md
├── fastapi.md
├── postgresql.md
└── ...
```

Naming: `<nombre-en-kebab-case>.md`, sin prefijos numéricos.

### Templates de tech

Los templates de estructura de EDR viven en la referencia canónica `~/.claude/references/edr/templates/`, un archivo por artefacto:

- Mini-EDR de tecnología → `~/.claude/references/edr/templates/tech-edr.md`
- INDEX del catálogo (`.matecito-ai/edr/tech/INDEX.md`) → `~/.claude/references/edr/templates/tech-index.md`

Las categorías sin filas en el INDEX se dejan vacías para que se vean los huecos.

---

## Materialización

### Paso 1: Resumir y confirmar

Antes de escribir nada, mostrá un resumen completo de todas las decisiones, agrupadas por fase, con su status. Pedí confirmación final. Permití editar cualquier respuesta.

### Paso 2: Estructura de archivos a generar

Los EDRs de salida son **slug-based** (sin prefijos numéricos) y van **agrupados por dominio en subcarpetas**, con un índice por dominio más un índice raíz. Misma taxonomía de dominios que el catálogo interno.

```
<root>/
├── CLAUDE.md                              # mínimo, apunta al índice raíz
└── .matecito-ai/
    └── edr/
        ├── INDEX.md                       # índice RAÍZ: explica cada dominio + cómo navegar
        ├── <dominio>/                      # una carpeta por dominio CON al menos un EDR
        │   ├── INDEX.md                    # índice del DOMINIO: lista sus EDRs + criterio
        │   ├── <slug>.md                   # un EDR por fase tratada (ej: error-handling.md)
        │   └── ...
        ├── ...                             # otros dominios
        └── tech/
            ├── INDEX.md                    # catálogo de tecnologías
            └── <una tech>.md
```

Reglas de la estructura:

- **Solo se crean carpetas de dominios que tienen al menos un EDR-archivo** (`Accepted`, `Pending` o `Deferred`). No repliques los 17 dominios en cada proyecto — la salida refleja lo que el proyecto realmente definió. (El catálogo interno sí tiene los 17; la salida solo los usados.)
- **Qué genera archivo y qué no:**
  - `Accepted` → EDR-archivo con contenido completo.
  - `Pending` / `Deferred` → EDR-archivo con su trigger/condición (el modo update los rellena).
  - `Not Applicable` → **fila en el INDEX del dominio** (sección "No aplican"), sin archivo propio.
  - El nombre de archivo es el `name` del concern (el slug); el dominio es el campo `domain` del frontmatter.
- **Dominio sin ningún EDR-archivo:** si todas sus fases quedaron `Not Applicable`, no se crea carpeta; el dominio se lista en el INDEX raíz (sección "Dominios sin uso") con una razón de 1 línea.
- **Qué se lista como N/A:** solo las fases que la matriz daba como Requerido/Recomendado para el tipo y se sacaron. Los concerns que la matriz nunca recomendó para este tipo no se enumeran — su "no aplica" ya está en la matriz.
- **Dos niveles de índice:** el raíz (`edr/INDEX.md`) enruta por dominio y lista los dominios sin uso; cada dominio (`edr/<dominio>/INDEX.md`) lista sus EDRs y sus N/A. Más `tech/INDEX.md` para tecnologías.

### Paso 3: Templates

Los templates de estructura de EDR son el **contrato canónico** y viven en la referencia `~/.claude/references/edr/templates/` (índice en `~/.claude/references/edr/templates/INDEX.md`). El template del `CLAUDE.md` raíz es propio de bootstrap y vive en `templates/claude-md.md`. No se duplican acá: antes de materializar, leé el template del artefacto que vas a escribir.

| Artefacto | Template |
|---|---|
| EDR individual (`<dominio>/<slug>.md`) | `~/.claude/references/edr/templates/edr.md` |
| Índice raíz (`edr/INDEX.md`) | `~/.claude/references/edr/templates/index-root.md` |
| Índice de dominio (`edr/<dominio>/INDEX.md`) | `~/.claude/references/edr/templates/index-domain.md` |
| Mini-EDR de tecnología (`edr/tech/<nombre>.md`) | `~/.claude/references/edr/templates/tech-edr.md` |
| Índice de tech (`edr/tech/INDEX.md`) | `~/.claude/references/edr/templates/tech-index.md` |
| `CLAUDE.md` raíz | [`templates/claude-md.md`](templates/claude-md.md) |

Notas del contrato del EDR (también en `~/.claude/references/edr/templates/edr.md`): **no hay sección `Historial`** (lo lleva git; la evolución se ve en la cadena de `Superseded`); **links entre EDRs** — dentro del mismo dominio ruta corta (`<slug>.md`), entre dominios ruta relativa (`../<dominio>/<slug>.md`).

### Paso 4: Escribir y reportar

1. Para cada dominio con al menos un EDR-archivo (`Accepted`/`Pending`/`Deferred`): `mkdir -p .matecito-ai/edr/<dominio>`. Además `mkdir -p .matecito-ai/edr/tech`.
2. Escribir `CLAUDE.md` (si no existe; si existe, **NO sobrescribir** — preguntar al usuario qué hacer)
3. Escribir `.matecito-ai/edr/INDEX.md` (índice raíz) listando los dominios con EDR-archivo y, en su sección "Dominios sin uso", los dominios cuyas fases quedaron todas `Not Applicable`.
4. Escribir `.matecito-ai/edr/<dominio>/INDEX.md` para cada dominio usado: la tabla de EDRs (Accepted/Pending/Deferred) y la sección "No aplican" con las fases `Not Applicable` del dominio y su razón.
5. Escribir los EDR-archivo de las fases con contenido: `Accepted` completo; `Pending`/`Deferred` con su trigger/condición. Los `Not Applicable` no generan archivo — quedan como fila en el INDEX del dominio (o del raíz si el dominio quedó sin uso).
6. Escribir `tech/INDEX.md` (los archivos individuales de tech ya se fueron creando intercalados).
7. Reportar al usuario:
   - Lista de archivos creados (path completo), **agrupada por dominio**
   - Resumen de 1 línea por EDR-archivo, con su status entre corchetes (`[Accepted]`, `[Pending]`, `[Deferred]`) y su `type`
   - Conteo de `Not Applicable` por dominio (viven en los INDEX), no uno por uno
   - Tecnologías registradas en `tech/`
   - **Lista separada de EDRs `Pending`/`Deferred` con su trigger**, así sabe qué quedó por decidir
   - Sugerencia de commitear estos archivos al repo
8. Ofrecer correr el validador `development-decisions-validate` en contexto fresco (ver flujo, paso 6) antes de dar por cerrado el bootstrap.

**Vocabulario del EDR — self-check obligatorio antes de escribir cada EDR-archivo (`Accepted`):** el razonamiento (Contexto/Decisión/Consecuencias/Alternativas) va en **conceptos, patrones y límites**. Ningún identificador interno volátil (clase, método, columna, error interno, ruta de archivo) en esas secciones — reubicalo a `## Alcance` (glob estable) o a `## Reglas verificables` (ahí sí nombrás la clase, es el ancla que se chequea). No inventes subsecciones tipo "Forma hexagonal" / "Mecanismo concreto" que vuelquen la implementación en el porqué. Excepción: nombre de tecnología/librería (es la decisión) y contrato público (endpoint público, código de error expuesto). Ver `~/.claude/references/edr/README.md` → "No es el cómo".

---

## Modo update (cuando `.matecito-ai/edr/INDEX.md` ya existe)

1. **Leé el índice raíz, los índices de dominio y los EDRs** existentes (`find .matecito-ai/edr -name '*.md'`).
2. **Mostrá un resumen agrupado por dominio y, dentro de cada uno, por status:** `Accepted`, `Pending` (con trigger), `Deferred`, `Not Applicable` (con razón), `Inferred` (borrador minado del código, sin porqué).
3. **Preguntá si algún `Pending` o `Deferred` ya está listo para resolverse.** Es lo más importante del modo update — sin esto, los "lo decidimos después" se pierden.
3b. **Ratificá los `Inferred` (borradores minados del código).** Un `Inferred` tiene el QUÉ y la evidencia, pero el PORQUÉ vacío — es un candidato sin ratificar (ver `~/.claude/references/edr/README.md`). Por cada uno, ofrecé ratificarlo: entrevistá por el porqué (Contexto, Decisión razonada, Alternativas, Consecuencias), llená esas secciones, **descartá la sección `## Evidencia (inferida)`** (es transitoria), y cambiá `Status: Inferred → Accepted`. Si el usuario no quiere ratificarlo ahora, queda `Inferred` (no se pierde). NUNCA promuevas un `Inferred` sin la entrevista del porqué — eso es lo que lo convierte en decisión.
4. **Ratchet — barré el catálogo:** leé `concerns/INDEX.md`, listá las fases relevantes al tipo de proyecto que **no tengan EDR todavía** (típicamente fases nuevas agregadas al catálogo desde la última corrida) y ofrecé tratarlas ahora. Mostralas agrupadas por dominio, incluyendo si caen en un dominio que el proyecto todavía no usa (esa carpeta se crea recién al materializar el primer EDR). Esta es la forma de que los temas agregados al catálogo lleguen a proyectos viejos.
5. **Después preguntá qué más quiere hacer:**
   - **Resolver un Pending/Deferred** → recorrer las preguntas de esa fase, cambiar Status a `Accepted`, llenar contenido.
   - **Ratificar un `Inferred`** → entrevistar por el porqué, llenar Contexto/Decisión/Alternativas/Consecuencias, descartar `## Evidencia (inferida)`, `Status → Accepted`.
   - **Actualizar una decisión (cambio menor)** → editar el EDR. Git lleva el historial.
   - **Cambiar una decisión (cambio de fondo)** → crear EDR nuevo en el mismo dominio, marcar el viejo `Superseded` con link al nuevo. No editar la decisión vieja en el lugar.
   - **Agregar una decisión nueva** no cubierta → crear EDR en su dominio + fila en el índice de ese dominio (y en el raíz si el dominio es nuevo en el proyecto).
   - **Cambiar un `Not Applicable` a `Pending`/`Accepted`** → el contexto del proyecto cambió (ej: el script chico creció a app multiusuario y ahora sí hay auth). Sacá la fila de la sección "No aplican" del INDEX del dominio (o "Dominios sin uso" del raíz) y creá el EDR-archivo con el nuevo status y contenido; creá la carpeta del dominio si no existía.
   - **Agregar/cambiar/quitar una tecnología** → editar `tech/INDEX.md` y el archivo en `tech/<nombre>.md`. Si reemplazás, el viejo queda `Superseded` apuntando al nuevo.
   - **Rehacer todo desde cero** → confirmación doble. Antes de sobrescribir, mover el directorio a `.matecito-ai/edr.old.<timestamp>/`.
6. Para actualizar/agregar, recorré solo las fases relevantes — no rehagas todo el cuestionario.
7. **Después de cualquier cambio, mantené los índices coherentes:** actualizá el índice del dominio afectado y, si agregaste o vaciaste un dominio, el índice raíz.

---

## Ratchet: agregar una fase al catálogo

El valor de largo plazo de la skill es que **nunca se vuelva a olvidar un tema**. Cuando aparece un concern que no estaba:

1. Determiná a qué **dominio canónico** pertenece (consultá el "criterio de pertenencia" en `concerns/<dominio>/INDEX.md`). Si encaja en un dominio reservado, ese dominio pasa de reservado a activo.
2. Creá `concerns/<dominio>/<slug>.md` con el formato estándar (mirá `concerns/runtime/error-handling.md` para una fase `deep` y `concerns/runtime/caching.md` para una `light`). Incluí en el frontmatter `domain` y `type`.
3. Sumá la fila al `concerns/<dominio>/INDEX.md` y a la matriz de `concerns/INDEX.md`, con la aplicabilidad por tipo de proyecto (Requerido/Recomendado). Si el dominio pasó de reservado a activo, movelo de tabla en el índice raíz.

Desde ese momento, todo bootstrap futuro lo considera, y el modo update lo ofrece a proyectos viejos (paso 4 de update). El catálogo se sembró de taxonomías externas (ISO/IEC 25010, 12-factor, arc42, OWASP ASVS, production-readiness) para nacer casi completo y solo crecer.

---

## Anti-patterns que esta skill evita

- ❌ Tirar todas las preguntas en un solo turno → la gente abandona.
- ❌ Forzar Clean Architecture en un script de 200 líneas → adaptar el set de fases al tipo de proyecto.
- ❌ Saltar una fase sin documentar el motivo → siempre dejar registro: `Not Applicable` como fila en el INDEX del dominio; `Pending`/`Deferred` como EDR-archivo + razón.
- ❌ Crear un EDR-archivo por cada `Not Applicable` → los N/A viven como fila en el INDEX del dominio, no como archivo; solo se justifican las desviaciones de la matriz.
- ❌ Preguntar el motivo de cada N/A por separado → clasificar en bloque, una sola pasada, con razón templada por tipo de proyecto.
- ❌ Confundir "no decidido aún" (`Pending`) con "decidido que no aplica" (`Not Applicable`) → son status distintos.
- ❌ Editar una decisión de fondo en el lugar → para cambios de decisión, supersede (EDR nuevo + viejo `Superseded`). Cambios menores sí se editan; git lleva el historial.
- ❌ Mantener una tabla `Historial` manual → es redundante con git y se pudre.
- ❌ Inventar reglas no discutidas con el usuario en la materialización → todo lo que va al EDR fue confirmado.
- ❌ Reglas vagas tipo "tratá de no acoplar capas" → siempre verificable: paths, globs, ejemplos.
- ❌ Nombrar identificadores internos volátiles (clases, métodos, columnas, rutas de archivo, errores internos) en el razonamiento del EDR (Contexto/Decisión/Consecuencias/Alternativas) → van a `## Alcance` (glob estable) o `## Reglas verificables` (ancla chequeable), nunca al porqué. El razonamiento se escribe en conceptos/patrones/límites; se pudre con el primer rename si calca el código. Excepción: nombre de tecnología/librería y contrato público. Ver `~/.claude/references/edr/README.md`.
- ❌ Sobrescribir un `CLAUDE.md` existente sin permiso → preguntar y ofrecer merge.
- ❌ Asumir el stack en lugar de detectarlo en pre-flight → leer manifests primero.
- ❌ Leer todo el catálogo `concerns/` de una → leer `INDEX.md` (raíz) para seleccionar, y cada `concerns/<dominio>/<slug>.md` solo cuando se trata esa fase.
- ❌ Inventar un dominio nuevo en un repo → la taxonomía es fija e impuesta por el motor; un dominio nuevo es decisión de catálogo, no de proyecto.
- ❌ Replicar los 17 dominios en la salida → en `.matecito-ai/edr/` solo se crean las carpetas de dominios con al menos un EDR-archivo (`Accepted`/`Pending`/`Deferred`); los dominios solo-N/A se listan en el INDEX raíz.
- ❌ Dejar índices desincronizados tras un cambio → actualizá el índice del dominio afectado y el raíz si corresponde.
- ❌ Dejar el catálogo `tech/` vacío hasta el final → registrar intercalado, mientras la justificación está fresca.
- ❌ Agregar una dependencia en sesiones futuras sin consultar `tech/INDEX.md` → revisar primero si ya hay algo elegido.
- ❌ En modo update, no preguntar por los `Pending`/`Deferred` ni barrer el catálogo por fases nuevas → es como se pierden las decisiones aplazadas y las fases agregadas.

---

## Recordatorio final

El valor de esta skill no está en las preguntas — está en que las decisiones queden **escritas, accionables y mantenidas**. Si las preguntas son geniales pero los EDRs salen vagos, fallamos. Si los EDRs son específicos y verificables, Claude (y cualquier dev) puede trabajar respetando las convenciones sin volver a preguntar.

Escribí los EDRs con la misma claridad con la que le explicarías la convención a un dev nuevo el primer día.
