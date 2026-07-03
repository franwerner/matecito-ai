---
name: development-spec-bootstrap
description: Entrevista interactiva por capability para definir el COMPORTAMIENTO de un sistema (qué hace, no cómo está construido) y materializarlo como capability-specs organizados por tipo en .matecito-ai/development-specs/<type>/. Usá esta skill SIEMPRE que el usuario quiera "definir el sistema completo antes de empezar", "especificar el comportamiento", "escribir los flujos", "definir las reglas de negocio", "qué hace el sistema ante X", pida ayuda con casos borde de un flujo, o mencione "capability spec", "spec de comportamiento", "definir capabilities". También si detectás un repo con .matecito-ai/adr/ pero sin .matecito-ai/development-specs/ y el usuario está por definir comportamiento (flujos, reglas, estados) todavía no escrito. NO es para decisiones técnicas (eso es development-decisions-bootstrap → ADRs).
---

# Development Spec Bootstrap

Entrevista al usuario para capturar el **comportamiento** del sistema —qué hace ante cada situación— y lo materializa como **capability-specs** estructurados, **organizados por tipo**, que Claude consultará en futuras sesiones vía `.matecito-ai/development-specs/INDEX.md`.

El objetivo es que el comportamiento quede **definido y verificable** antes de construir, no implícito en el PRD o en la cabeza del autor. Eso le permite a Claude (y a cualquier dev) implementar y testear contra un contrato claro del *qué hace*.

Es la contraparte de `development-decisions-bootstrap`: esa captura *qué se eligió y por qué* (ADRs); esta captura *qué hace* (specs). Ver la distinción en `~/.claude/references/spec/README.md`.

---

## Qué es (y qué no es) un capability-spec

El concepto canónico —qué cuenta como capability-spec, la separación con el ADR (qué-hace vs por-qué), los tipos y los estados— vive en `~/.claude/references/spec/README.md` (referencia consultable, agnóstica de flujo). bootstrap **aplica** ese concepto; no lo redefine. La *estructura* concreta del archivo está en `~/.claude/references/spec/templates/capability.md`.

**Lo esencial que el motor asume:**
- El spec describe **comportamiento observable y verificable**, en **idioma de dominio + contrato público**. NUNCA identificadores internos volátiles (clases, métodos, columnas, errores internos, rutas de archivo). El *cómo* es del código; el *por qué* es del ADR.
- Cada regla/flujo/borde importante lleva al menos un **escenario Given/When/Then** — es lo que lo vuelve testeable.

---

## Arquitectura de esta skill

A diferencia de `development-decisions-bootstrap`, esta skill **no tiene catálogo de datos**: las capabilities no son una taxonomía fija, se **descubren por proyecto**. Lo único fijo es la **taxonomía de tipos**, que vive acá mismo (abajo). Todo el motor es este `SKILL.md`.

- **Entrada:** la descripción del sistema + el PRD/proposal si existe + los ADRs ya definidos (para linkear, no repetir).
- **Salida:** capability-specs en `.matecito-ai/development-specs/<type>/<capability>.md`, con índice raíz (enruta por tipo) + un `INDEX.md` por tipo.

---

## Tipos de capability (taxonomía fija)

Taxonomía **cerrada**: toda capability es de uno de estos 4 tipos, que definen la subcarpeta (`<type>/<capability>.md`) y qué secciones de la plantilla son su **esqueleto**. No inventés tipos nuevos por proyecto.

| `type` | Qué captura | Ejemplo | Secciones esqueleto |
|---|---|---|---|
| `flow` | Operación de cara a un actor, con pasos y ramas | `connect-meta-account`, `send-message` | Flujo principal · Ramas · Casos borde · Errores de cara al actor |
| `rule` | Regla de negocio transversal, sin flujo | ventana de 24h, deduplicación | Reglas de negocio · Escenarios |
| `lifecycle` | Máquina de estados de una entidad | `message` (pending→sent→failed) | Entidades y estados · Escenarios |
| `process` | Comportamiento reactivo/de fondo, disparado por evento (no por un actor) | reconciliación, jobs | Flujo principal (del proceso) · Casos borde · Reglas de negocio |

Todas usan la **misma plantilla**; cada tipo enfatiza su esqueleto y omite lo que no aplica.

**La regla de clasificación** (orden de decisión + el tie-breaker del solapamiento por reuso) es canónica y vive en `~/.claude/references/spec/README.md` → «Cómo clasificar el tipo». Aplicala al enumerar; no la redefinas acá. La usa también `development-spec-validate` para chequear que el tipo declarado sea el correcto.

---

## Cuándo correr esta skill

- El usuario quiere **definir el sistema completo antes de construir** (su necesidad típica).
- Hay ADRs (`.matecito-ai/adr/`) pero el comportamiento no está escrito en ningún lado durable.
- El usuario pide especificar un flujo, una regla, un ciclo de vida o un proceso con detalle.

Si `.matecito-ai/development-specs/INDEX.md` ya existe con contenido: **NO rehagas todo**. Andá al modo `update` (final del documento).

---

## Reglas del motor (aplican a TODA capability)

**Una pregunta por turno.** Nunca dumpees una lista de preguntas. Una, esperá, leé, seguí.

**Preguntá por comportamiento, no por implementación.** Mal: "¿en qué clase va esto?". Bien: "¿qué pasa si el usuario cancela a mitad del flujo?". Si el usuario responde en términos de código, reformulá la respuesta a comportamiento antes de escribirla.

**Una línea de "por qué importa" antes de cada sección.** Sin sermones.

**Empujá los casos borde.** El happy path el usuario lo tiene claro; el valor está en las ramas y los bordes. Preguntá activamente: "¿y si…?" (falla el sistema externo, llega dos veces, el actor no tiene permiso, el recurso no existe, se vence una ventana).

**Todo comportamiento afirmado necesita su escenario.** Antes de cerrar una capability, cada regla/rama/borde relevante tiene un Given/When/Then. Si no se puede escribir el escenario, la afirmación es vaga: refinala.

**Permití `Draft`.** Una capability puede quedar `Draft` (le faltan escenarios o bordes) con una nota de qué falta. Mejor un spec honesto en Draft que uno inventado en Accepted.

**Linkeá al ADR, no lo repitas.** Cuando un comportamiento está gobernado por una decisión técnica ya escrita en un ADR, referencialo en "Referencias"; no copies el cómo/por qué al spec.

---

## Pre-flight (siempre primero)

Antes de la primera pregunta, inspeccioná el repo:

```bash
ls -la
test -d .matecito-ai/development-specs && echo "--- specs existentes (por tipo) ---" && find .matecito-ai/development-specs -name '*.md' | sort
test -d .matecito-ai/adr && echo "--- ADRs a linkear ---" && find .matecito-ai/adr -name '*.md' | sort
test -f CLAUDE.md && echo "--- CLAUDE.md ---" && cat CLAUDE.md
find . -maxdepth 2 -iname 'PRD*' -o -iname 'README*' -o -iname '*proposal*' 2>/dev/null | head
```

Con eso sabés: si hay specs previos (→ modo update), qué ADRs existen (para linkear), y si hay un PRD del que derivar las capabilities.

---

## El flujo

### 1. Descripción del sistema

Una pregunta abierta:

> Contame a grandes rasgos qué hace el sistema: qué operaciones ofrece, con qué actores interactúa, y qué reglas de negocio importan. Si tenés un PRD, pasámelo o decime dónde está.

### 2. Enumeración de capabilities

De la descripción + el PRD, **derivá la lista de capabilities** del sistema y **clasificá cada una por tipo**. Este es el paso propio de esta skill (no hay catálogo que la dé). Presentala agrupada por tipo:

> Por lo que contás, el sistema tiene estas capabilities:
> - **`flow`:** `connect-meta-account`, `send-message`, `receive-message`
> - **`rule`:** `messaging-window-24h`, `message-deduplication`
> - **`lifecycle`:** `message`
> - **`process`:** `outbound-echo-reconciliation`
>
> ¿Falta alguna, sobra alguna, o querés renombrar?

Slugs en **inglés kebab-case**. Clasificá cada capability con la regla de `~/.claude/references/spec/README.md` → «Cómo clasificar el tipo» (orden de decisión + tie-breaker por reuso: una `rule`/`lifecycle` es capability aparte SOLO si la comparten varios flujos; si es exclusiva de uno, vive dentro de ese flujo). Ante duda genuina, confirmá con el usuario.

### 3. Ajuste del set

Incorporá lo que el usuario agregue/saque/renombre. Una capability que el usuario decide **no especificar ahora** no se inventa: se puede dejar como `Draft` con una nota, o directamente no crearla (a diferencia de los ADRs, un comportamiento sin definir no necesita constancia formal — no hay "Not Applicable" para specs).

### 4. Recorrido por capability

Por cada capability del set final, seguí "Cómo tratar una capability".

### 5. Materialización

Cuando se recorrieron todas, materializá (ver "Materialización").

### 6. Validación (recomendada)

Al cerrar, ofrecé correr `development-spec-validate` en **contexto fresco** (sub-agente): chequea coherencia entre specs (contradicciones, entidades/estados no definidos, referencias colgadas a ADRs), completitud y verificabilidad. No modifica nada — los hallazgos los resuelve el usuario vía modo update.

---

## Cómo tratar una capability

Procedimiento genérico del motor, para cualquier tipo:

1. **Identificá su tipo** y mirá sus **secciones esqueleto** (tabla de tipos). Esas son las que sí o sí trabajás; las demás secciones de la plantilla son opcionales según el caso.
2. Mostrá el **propósito** en una línea y confirmalo ("Esta capability logra X para Y, ¿sí?").
3. Recorré las secciones **una por turno**, en orden de la plantilla, saltando las que no aplican al tipo:
   - **Actores / Precondiciones** — quién la dispara y qué debe ser cierto antes.
   - **Flujo principal** (flow/process) — pasos observables, en lenguaje de dominio.
   - **Ramas / Casos borde** — empujá los "¿y si…?".
   - **Reglas de negocio** (rule/process) — invariantes con valores concretos.
   - **Entidades y estados** (lifecycle) — transiciones y qué las dispara.
   - **Errores de cara al actor** — el contrato de error observable.
4. **Escenarios.** Por cada regla/rama/borde relevante, escribí un Given/When/Then. Es el paso que vuelve verificable la capability; no lo saltees.
5. **Referencias.** Si algún comportamiento está gobernado por un ADR existente (lo viste en pre-flight), linkealo. Si notás que falta un ADR (una decisión técnica no tomada), anotalo como pregunta abierta para `development-decisions-bootstrap` — no lo resuelvas acá.
6. **Self-check de vocabulario** antes de escribir: ningún identificador interno volátil en ninguna sección. Reformulá a idioma de dominio / contrato público, o mové el ancla técnica a un link en "Referencias".
7. **Materializá** el spec en `.matecito-ai/development-specs/<type>/<capability>.md`.

---

## Materialización

### Estructura de archivos a generar

```
.matecito-ai/development-specs/
├── INDEX.md                        # índice RAÍZ: enruta por tipo
├── <type>/                         # una carpeta por tipo CON al menos un spec
│   ├── INDEX.md                    # índice del TIPO: sus capabilities + cuándo consultar
│   ├── <capability>.md             # un spec por capability
│   └── ...
└── ...
```

Reglas:
- **Solo se crean carpetas de tipos que tienen al menos un spec-archivo.** Los tipos sin uso se listan en el índice raíz (sección "Tipos sin uso"), sin carpeta.
- Nombre de archivo = el slug de la capability (kebab-case, inglés), sin prefijos numéricos.
- Dos niveles de índice: raíz (`development-specs/INDEX.md`) enruta por tipo; cada tipo (`<type>/INDEX.md`) lista sus capabilities y su criterio de "cuándo consultar".

### Templates

Los templates son el **contrato canónico** y viven en `~/.claude/references/spec/templates/`. Antes de materializar, leé el del artefacto que vas a escribir. No los dupliques acá.

| Artefacto | Template |
|---|---|
| Capability-spec (`<type>/<capability>.md`) | `~/.claude/references/spec/templates/capability.md` |
| Índice raíz (`development-specs/INDEX.md`) | `~/.claude/references/spec/templates/index-root.md` |
| Índice de tipo (`<type>/INDEX.md`) | `~/.claude/references/spec/templates/index-type.md` |
| Sección de `CLAUDE.md` (pointer a specs) | [`templates/claude-md-spec.md`](templates/claude-md-spec.md) (propio de esta skill) |

### Escribir y reportar

1. Para cada tipo con al menos un spec: `mkdir -p .matecito-ai/development-specs/<type>`.
2. Escribir cada `<type>/<capability>.md` (`Accepted` completo; `Draft` con nota de qué falta).
3. Escribir `<type>/INDEX.md` por tipo usado.
4. Escribir `.matecito-ai/development-specs/INDEX.md` (raíz), listando los tipos usados y, en "Tipos sin uso", los que no aplican.
5. **Asentar el pointer en el `CLAUDE.md` del proyecto** (idempotente, con la sección de `templates/claude-md-spec.md`): si `CLAUDE.md` no existe, crealo con esa sección; si existe, agregá o actualizá SOLO la sección `## Comportamiento del sistema (capability-specs)` (localizada por su heading) sin tocar el resto —en particular, NO pises la sección de ADRs si `development-decisions-bootstrap` ya la escribió—. Así Claude sabe que debe consultar los specs antes de implementar, igual que consulta los ADRs.
6. Reportar al usuario: lista de archivos creados **agrupada por tipo**, 1 línea por capability con su status entre corchetes, las capabilities que quedaron `Draft` con qué les falta, y si se creó/actualizó el `CLAUDE.md`. Sugerir commitear.
7. Ofrecer `development-spec-validate` en contexto fresco.

---

## Modo update (cuando `.matecito-ai/development-specs/INDEX.md` ya existe)

1. **Leé** el índice raíz, los índices de tipo y los specs existentes.
2. **Mostrá un resumen** agrupado por tipo y status (`Accepted`, `Draft` con qué falta, `Deprecated`).
3. **Preguntá si algún `Draft` está listo para completarse** — es lo más importante del update; sin esto los "lo definimos después" se pierden.
4. **Después preguntá qué más:**
   - **Completar un `Draft`** → recorrer sus secciones faltantes, `Status → Accepted`.
   - **Actualizar comportamiento (cambio menor)** → editar el spec. Git lleva el historial.
   - **Agregar una capability nueva** → tratarla con el procedimiento genérico + fila en el `INDEX.md` de su tipo (y en el raíz si el tipo es nuevo en el proyecto).
   - **Retirar una capability** → `Status → Deprecated` con link a su reemplazo; no borrar el archivo.
5. **Después de cualquier cambio, mantené los índices coherentes:** actualizá el índice del tipo afectado y, si agregaste o vaciaste un tipo, el raíz.
6. **Asegurá que el pointer del `CLAUDE.md` siga presente** (la sección `## Comportamiento del sistema (capability-specs)`); si falta, agregala desde `templates/claude-md-spec.md`.

> Nota: durante el flujo SDD, el comportamiento se actualiza **solo** vía el merge de `sdd-archive` (el delta del cambio → el capability-spec). El modo update es para autoría/mantenimiento manual fuera de un cambio.

---

## Anti-patterns que esta skill evita

- ❌ Nombrar identificadores internos volátiles (clases, métodos, columnas, rutas, errores internos) en cualquier sección → idioma de dominio + contrato público; el ancla técnica va como link al ADR en "Referencias".
- ❌ Escribir el *por qué* o el *cómo* en el spec → eso es ADR (por qué) o código (cómo). El spec es el *qué hace*.
- ❌ Afirmar comportamiento sin un escenario Given/When/Then que lo verifique → toda regla/rama/borde relevante lleva su escenario.
- ❌ Quedarse en el happy path → empujar los casos borde es el valor de la skill.
- ❌ Repetir el contenido de un ADR en el spec → linkearlo, no copiarlo.
- ❌ Inventar tipos nuevos por proyecto → la taxonomía de tipos es fija (`flow`/`rule`/`lifecycle`/`process`).
- ❌ Tirar todas las preguntas en un turno → una por turno.
- ❌ Dumpear la lista de capabilities sin clasificarlas por tipo → la enumeración siempre viene agrupada por tipo.
- ❌ Dejar índices desincronizados tras un cambio → actualizá el índice del tipo afectado y el raíz.
- ❌ Materializar comportamiento no confirmado por el usuario → todo lo que va al spec fue acordado.

---

## Recordatorio final

El valor de esta skill no está en las preguntas — está en que el comportamiento quede **escrito, verificable y mantenido**, y separado limpio del *por qué* (ADR) y del *cómo* (código). Si el spec sale vago o calca el código, fallamos. Si cada capability dice con precisión qué hace el sistema ante cada situación —con sus escenarios testeables— cualquiera puede implementar y verificar contra un contrato claro.

Escribí los specs con la misma claridad con la que le explicarías el comportamiento del sistema a un dev nuevo el primer día.
