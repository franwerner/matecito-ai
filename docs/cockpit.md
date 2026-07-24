# Matecito UI — contrato de salida estructurado y vista del flujo SDD

> **Estado:** borrador conceptual · 2026-07-23 · **Matecito UI**, dentro de este mismo repo (`matecito-ai`).
> Este doc captura **qué vamos a hacer y por qué**, a nivel visión. El diseño técnico fino
> (forma exacta de los schemas, contratos entre capas) lo produce después la fase de diseño del flujo.

## 0. TL;DR

Empezó como "una UI para ver los agentes del flujo SDD". Al discutirlo, el problema de raíz resultó ser otro: **el flujo SDD emite artefactos en lenguaje natural poco estructurado**, y por eso todo lo que quieras construir encima (observabilidad, trazabilidad, cumplimiento, UI) es difícil de extraer.

El reencuadre: **esto no es "una UI". Es darle al flujo SDD un contrato de salida estructurado.** La UI pasa a ser un consumidor río abajo que cae casi gratis una vez que los datos están estructurados.

Dos piezas, en orden:

1. **Ahora — la fundación (lectura):** un **contrato de salida estructurado** por fase (JSON como fuente de verdad) + una **vista de lectura/trazabilidad** de una corrida del flujo.
2. **Después — el norte (escritura):** operar desde la UI — aprobar gates, reiterar un pedazo puntual, co-modelar entidades/contratos.

---

## 1. El problema y el reencuadre

Hoy cada fase del flujo (`intake → … → archive`) deja su artefacto como **prosa/markdown**. Eso tiene dos costos:

- **No es legible por máquina.** Para mostrar "qué decidió cada agente", "qué código tocó", "qué MCP ejecutó" hay que *parsear prosa* — frágil y ambiguo.
- **No hay estándar de salida.** Cada fase se expresa a su manera; no hay una forma garantizada de cada artefacto.

El keystone de todo este trabajo:

> **Si estructurás en el origen, todo lo de río abajo se vuelve trivial.** La UI, la trazabilidad, la capa de cumplimiento y los diffs dejan de ser "extracción desde texto" y pasan a ser "render de datos ya estructurados". Por eso lo primero que se construye **no es la UI** — es el contrato estructurado.

## 2. Principios

### 2.1 La inversión: JSON es la fuente de verdad

Hoy el markdown es canónico (los EDRs *son* archivos `.md`). Lo invertimos:

```
        ANTES                              AHORA
   markdown canónico          JSON estructurado = fuente de verdad
        │                                 │
        └─ (la UI parsea prosa)           └─ el markdown es una PROYECCIÓN renderizada
```

La AI llena un JSON tipado; el markdown legible es una **vista renderizada** de ese JSON, no al revés.

### 2.2 Estructurar *alrededor* de la prosa, no *en lugar* de ella

La trampa de "estructurar todo" es matar el matiz. El punto dulce: cada artefacto es un **esqueleto tipado con bolsillos de prosa**.

- **Slots tipados** para lo que es máquina: refs, `status`, `%`, listas, ids.
- **Texto libre** para lo que es narrativa: el rationale, la explicación del tradeoff, el porqué de un desvío.

No se estructura el "por qué" — se estructura *envuelto*.

### 2.3 El MCP es dependencia dura del flujo

Decisión (Q9, opción A): el MCP **debe estar corriendo**. El flujo estructurado pasa por él y **no hay fallback** — sin MCP, el flujo SDD no corre.

> **Consecuencia asumida (revierte una invariante).** Esto invierte el `offer, don't impose` que este documento sostenía antes: el MCP pasa a ser **punto de dependencia**, no opcional. Se gana simplicidad (un solo camino de render, cero lógica de degradación); se paga que el flujo **requiere el MCP prendido**, y que las corridas **headless/CI** deben garantizar su disponibilidad. Es un cambio de principio central del ecosistema → **requiere su propio EDR**, no se desliza en silencio.

## 3. Arquitectura

Tres capas, cada una con un dueño claro:

```
┌──────────────────────────────────────────────────────────────────────┐
│  CANÓNICO EXTERNO     .matecito-ai/* (EDRs, specs en .md) · git         │
│                       Engram (memoria semántica, cross-proyecto)       │
├──────────────────────────────────────────────────────────────────────┤
│  STORE ESTRUCTURADO   SQLite autocontenida (event-log) — persistente   │
│  (Matecito UI)        eventos + artefactos de flujo + índice del       │
│                       change/agentes/código/decisiones · queryable     │
├──────────────────────────────────────────────────────────────────────┤
│  BROKER               mini-API — único escritor de la SQLite            │
│                       sirve a la UI por WebSocket, rutea acciones       │
└──────────────────────────────────────────────────────────────────────┘
```

- **Canónico externo:** los EDRs/specs siguen en `.md`, el código en git, la **memoria semántica en Engram**. Matecito UI **no reemplaza** nada de esto.
- **Store estructurado (SQLite):** el corazón de Matecito UI. Autocontenida (un archivo, sin server aparte, sin infra pesada). El **event-log** (Q1) es una tabla `events`; el estado derivado son queries/vistas → timeline, diff v1→v2 y `supersedes` salen gratis, y ahora **consultables** por la UI. **Persistente** (no efímero): un change sobrevive entre sesiones, **identificado por su `change-name`/rama** — una sesión nueva continúa en la misma fila (Q10). **Indexa (no mueve)** los EDR/spec `.md` — mismo patrón que `codegraph` indexa código que vive en archivos: `.md` canónico, índice consultable en la DB; si divergen, gana el `.md` (se re-indexa).
- **Broker (mini-API):** único escritor de la SQLite (evita el split-brain AI+UI; SQLite en WAL da 1 escritor + N lectores). La AI escribe vía MCP→API; la UI vía HTTP→API.

> **Actualización (2026-07-23) — store global, no per-proyecto.** El store pasó a **instancia única global**: **un daemon broker único por máquina** + **una SQLite única en `~/.matecito-ai/`** compartida entre todos los proyectos (revierte el "per-proyecto" de Q10/Q11). Cada proyecto se **registra por su ruta** (vía MCP o UI) y toda la data se scopea por `project_id`; la identidad del change pasa a `project_id + change-name`/rama. El invariante único-escritor se mantiene porque hay **un solo proceso** escribiendo. Las decisiones de arquitectura del broker viven como EDRs en `apps/api/.matecito-ai/edr/`.

> **Actualización (2026-07-24) — contenido versionado, no solo índice.** El broker no solo **indexa** la proyección de los EDR/spec: guarda su **contenido versionado** (copy-on-write lazy). Un evento pinea la versión **vigente** al persistir; una versión solo se **congela/forkea** cuando el `.md` cambia y ya estaba referenciada (si no, update-in-place). Así la UI muestra un EDR/spec en la **versión exacta que aplicó cada evento**, aun si después se editó o borró — sin versiones de más. El `.md` sigue canónico para editar (git). Detalle en los capability-specs `index-decision-records` / `submit-phase-artifact`.

> **SQLite y Engram conviven — no compiten.** Engram es memoria **semántica, difusa, cross-proyecto** ("¿cómo resolvimos auth la otra vez?"); la SQLite es estado **estructurado, exacto** ("los eventos y artefactos del change #42"), ahora **global y compartido** entre proyectos (scopeado por `project_id`). El overlap es superficial. Reemplazar Engram sería una decisión **aparte**, con impacto en todo el ecosistema — no acoplada a Matecito UI.

### 3.1 El rol del MCP

```
CLAUDE  ──(JSON estructurado)──►  MCP  ──►  SERVER/mini-API  ──►  UI
   ▲                              │
   └────(markdown renderizado)────┘
```

El MCP hace dos cosas:

1. **Conecta** los componentes (Claude ↔ server ↔ UI).
2. **Es el registry + renderer del contrato:** recibe el JSON que la AI llena, y dada una **plantilla de texto** rellena las partes y devuelve el artefacto renderizado. La AI nunca improvisa el formato del artefacto — llena slots tipados, el MCP garantiza la forma canónica.

### 3.2 El split: hooks para lo mecánico, MCP para lo semántico

No todo pasa por el MCP. Se parte por naturaleza:

| Naturaleza | Ejemplos | Canal |
|---|---|---|
| **Mecánica / observabilidad** | árbol de agentes, MCP calls hechas, diffs | **hooks + transcript** (disparan solos; la AI no se narra) |
| **Semántica / bidireccional** | el artefacto que la AI produce, análisis de cumplimiento, respuestas de gate, "reiterá esto" desde la UI | **MCP** (regla del flujo) |

> **Por qué el split.** Si "la AI debe llamar al MCP para reportar su estado" fuera regla, la fidelidad del registro quedaría a merced de que el agente se acuerde. Los hooks nunca se olvidan. El MCP se reserva para lo que **solo** la AI puede producir o para lo que necesita respuesta de vuelta.

## 4. El contrato de salida estructurado

Cada tipo de artefacto (intake brief, propuesta, spec, EDR, tasks, verify report…) tiene:

- un **schema** (los slots tipados + los bolsillos de prosa), y
- una **plantilla** de render a markdown.

Ambos son propiedad del contrato (con fallback de render local, §2.3). La AI piensa y llena el JSON; el contrato valida y proyecta. Esto es lo que convierte "el flujo" en algo con **estándar de salida**.

## 5. La vista de lectura (trazabilidad) — la primera pieza

**Read-only.** Solo mirar y navegar; nada de aprobar, editar ni materializar (eso es el norte, §7).

**Unidad:** el *change* (una corrida del flujo). Puede haber un índice de changes por encima.

**El objeto central: un grafo** con tres tipos de nodo —fases/agentes, decisiones/specs, código— y las llamadas MCP colgando de los agentes. No son tres pantallas: es el mismo grafo mirado desde **tres puertas**.

| Puerta | Pregunta | Qué ves |
|---|---|---|
| **Por el flujo** | ¿cómo se construyó? | el árbol orquestador → fases → agentes, en orden |
| **Por el código** | ¿qué se tocó y de dónde salió? | entrás por un archivo, ves qué fase/agente/decisión lo tocó (`git blame` por flujo) |
| **Por la decisión** | ¿qué se decidió? | entrás por la lista de EDRs, ves qué código/spec justificó cada una |

**Qué muestra cada nodo:**

- **Agente/fase:** tipo, el mandato recibido, duración, estado, **los MCP que ejecutó**, lo que produjo, su result contract.
- **Decisión (EDR) / spec:** contenido, estado (`Inferred`/`Accepted`), el porqué, qué justificó.
- **Código:** referencia de archivo + **diff (post-apply, no en vivo)**, qué agente lo escribió, a qué spec/EDR traza.
- **Llamada MCP:** qué tool, qué consultó — evidencia de "qué miró el agente antes de decidir".

**Dimensión temporal:** como el estado es event-log, hay un **scrubber** que reproduce el orden real de la corrida. En vivo es un stream; congelado es un mapa navegable. **Es el mismo objeto.**

**Principio de DX:** es una superficie de lectura, así que la única regla es **todo es remontable** — cada cosa tangible linkea a su origen, nunca hay punto muerto para "¿de dónde salió esto?". Cero prosa suelta, todo referencias.

### 5.1 El flujo de agentes en vivo

La puerta "por el flujo" es el corazón de lo que motivó todo esto: **ver la bifurcación que va tomando el orquestador y qué hace cada agente en el camino.** Se muestra en tres capas, de más garantizada a más rica:

- **Capa 1 — El roster.** Cada agente que lanza el orquestador aparece como una tarjeta: tipo (`sdd-apply`, `sdd-explore`…), fase, el **mandato** que recibió, inicio, estado (`corriendo`/`listo`). Al terminar, su **result contract** (status, summary, artifacts, risks). Sale de los hooks `SubagentStart/Stop` (§6).
- **Capa 2 — El feed de actividad por agente.** Qué está haciendo cada uno: qué archivo lee, qué edita, qué **MCP ejecuta**, su razonamiento. Sale de tailear el transcript del agente (§6) — casi en vivo.
- **Capa 3 — El grafo con aristas etiquetadas.** Contesta *qué información pasa entre los agentes*. En cada arista, el payload real:
  - **hacia el agente:** el mandato + las topic-keys de Engram que leyó + el modelo + la fase.
  - **de vuelta:** el result contract.
  - **efectos laterales:** lo que escribió en Engram / `.matecito-ai/*` / el código.

**El camino completo del change.** Uniendo las tres capas ves la corrida de punta a punta: qué fases se activaron, cuántos agentes se lanzaron en cada una, la bifurcación que tomó el orquestador y —por cada agente— **qué parte del código tocó y qué decisiones (EDRs) gobiernan esa parte** (con su cumplimiento, §5.2). Del código remontás a la decisión, y de la decisión al agente que la aplicó.

**Paralelismo legible.** El orquestador puede abrir varios agentes a la vez; en el chat eso es un choclo intercalado ilegible. En la UI van lado a lado, cada uno con su feed. Acá la UI le pasa el trapo al chat.

### 5.2 Capa de cumplimiento (spec + EDR)

En `verify`, la AI **marca las partes del código que tocó** con una etiqueta a revisar —`CRITICAL` / `WARNING` / `SUGGESTION`— y el **problema encontrado**, contra sus **dos contratos**: el **spec** (comportamiento) y el **EDR** (decisión). Es lo que `verify` ya juzga hoy (*Decision Gaps*, severidades); la exigencia nueva es emitirlo **estructurado** y **anclado al pedazo de código**, no en prosa suelta.

> **Sin porcentaje.** Se descartó el % de cumplimiento — juicio de la AI inherentemente difuso, que mentía con confianza. Solo **etiqueta + problema**.

Se apoya en que **cada artefacto de código implementado ya lleva su justificación** (de `apply`): el nodo de código en el grafo reúne ruta + diff + justificación + —si `verify` lo marcó— etiqueta + problema. Del hallazgo remontás a la justificación y a la decisión/spec que debía cumplir.

## 6. Qué es capturable (y qué no)

Verificado contra la documentación oficial de Claude Code (hooks + sub-agents):

- **Roster (garantizado):** los hooks `SubagentStart` / `SubagentStop` traen `agent_id`, `agent_type` y (en Stop) el `last_assistant_message`. Marcan inicio/fin + result contract.
- **Feed de actividad por agente (destrabado vía transcript):** cada sub-agente escribe su propio transcript en
  `~/.claude/projects/{project}/{sessionId}/subagents/agent-{agentId}.jsonl`.
  La mini-API **tailea** ese archivo → atribución por construcción (es *su* archivo). El transcript incluye acciones (reads/edits/tool-calls, **incluidas las MCP calls**) y el texto de razonamiento, escrito a medida que trabaja.

**Plumbing:** `SubagentStart` (da `agent_id`) → aparece la tarjeta + arranca el tail de su `.jsonl` → feed; `SubagentStop` → cierra + `last_assistant_message`. Encaja sin tocar la arquitectura: la mini-API ya es file-watcher (event-log); suma el watch de `subagents/*.jsonl`.

> **El límite y su corrección.** No hay stream token-a-token *empujado* — pero tailear el transcript da texto + acciones **casi en vivo**, mucho más rico que "solo el contrato final". Cuán fluido se ve depende de la cadencia de flush del transcript — **medido (Q5): casi en vivo, a nivel de evento** (cada tool-call/mensaje se escribe con ~1-2 s de latencia, incremental durante toda la corrida, nunca buffered al final). Alcanza de sobra para un feed en vivo por acciones. El tail debe ser **incremental** (el archivo llega a ~150KB), no re-leer entero.

**Incertidumbres no bloqueantes** (se resuelven empíricamente, no cambian el diseño):
- ~~Cadencia de flush del transcript~~ → **resuelto (Q5): casi en vivo, event-level ~1-2 s.**
- Atribución vía hooks anidados `PreToolUse/PostToolUse`: doc contradictoria → **no se usa** (se usa el transcript).
- ~~Payload del launch (`Agent`/`Task`) con `subagent_type` + prompt: no documentado~~ → **resuelto (Q6): el mandato está en la línea 1 del transcript**, disponible desde el arranque; se lee de ahí (no se depende del payload del hook).

## 7. El norte (escritura — fase posterior)

Explícitamente **futuro**. La vista de lectura es el **sustrato** sobre el que estas acciones van a parar; bien pensada ahora, el write después es casi gratis.

- **Gates en la UI:** aprobar el INTAKE GATE, materializar candidatos, ratificar `Inferred → Accepted` — desde la UI en vez del chat. La respuesta vuelve por el resultado de la tool MCP (la AI está bloqueada esperando ahí).
- **Reiterar puntual:** ver código que no cumple un EDR y **redisparar un agente sobre ese pedazo y solo ese** — posible únicamente porque la vista de lectura ya capturó la referencia tangible (diff ↔ EDR ↔ agente/task).
- **Canvas de co-modelado:** modelar entidades/contratos en un lienzo que la AI y el humano editan a la par; el output *es* el capability-spec/design que la fase consume. Encarna el principio **"todo lo que la AI infiere es un borrador con una manija"**: cada inferencia se agarra y se tuerce, en vez de argumentar en prosa.

> **El chat sigue siendo el centro de la conversación.** La UI es para lo estructurado (listas, diagramas, records); el chat para el razonamiento y el matiz. La AI reacciona a los cambios de la UI.

## 8. Decisiones abiertas / EDRs necesarios

Cambios de ecosistema que requieren su propio EDR y no deben deslizarse en silencio:

1. **Modelo de artefactos de flujo:** los artefactos de flujo pasan a **estructurados** (JSON → SQLite), markdown como proyección. **EDR/spec NO se migran** — siguen `.md` canónicos e **indexados** en SQLite (Q7/Q11).
2. **MCP como dependencia dura del flujo (Q9-A):** el flujo estructurado pasa por el MCP y **no corre sin él** → **revierte** la invariante `offer, don't impose` (§2.3). Toca el orquestador/skills de matecito-ai.
3. **Alcance de la re-plomería:** estructurar la salida de cada fase es un cambio profundo (todos los skills/agentes de fase). Se hace **fase por fase, empezando por `apply`** (Q8/Q8b). Territorio EDR + flujo completo.

## 9. Secuenciación

**Slice vertical fino primero** (Q8). No se construye la UI al principio ni se estructuran las nueve fases a ciegas:

```
FASE 0 — Slice vertical (validar el lazo) · fase = `apply`
  1. Contrato estructurado de UNA fase (schema + plantilla)
  2. La AI llena el JSON → el MCP valida + renderiza + persiste en SQLite
  3. Una UI mínima consulta la SQLite y lo muestra (lectura)
     └─ si el lazo cierra, el resto es repetir el patrón

FASE 1 — Escalar la lectura
  · Replicar el patrón al resto de las fases (cada una = schema + plantilla + índice)
  · Indexar los EDR/spec .md → se completa el grafo de trazabilidad (flujo · código · decisión)

FASE 2 — Cumplimiento (Q3)
  · verify emite etiquetas estructuradas (CRITICAL/WARNING/SUGGESTION + problema),
    ancladas al pedazo de código. Sin %.

FASE 3 — El norte: escritura (§7)
  · gates en la UI · reiterar puntual · canvas de co-modelado
```

> **La fase del slice (FASE 0) es la decisión inmediata.** El resto del roadmap ya está asentado arriba.

---

## 11. FASE 0 — desglose de trabajo

> Slice = `apply`. Todas las preguntas (`questions.md`) cerradas, sin TBDs. Este desglose *es* la fase `tasks` de un cambio real cuando esto se formalice al flujo (§9, opción b). Es el **qué**, no el **cómo** técnico fino — eso lo produce la fase de diseño.

**A · Contrato de `apply`**
- A1. Schema del artefacto de `apply` (por tarea: diff/archivo, justificación, refs a task/EDR/spec; slots tipados + bolsillos de prosa).
- A2. Plantilla de render (JSON → markdown) del artefacto.
- A3. Adaptar el skill/agente `sdd-apply` para que emita el JSON del contrato. *(toca el ecosistema → EDR)*

**B · MCP server "cockpit"**
- B1. Tool que recibe el JSON + lo valida contra el schema.
- B2. Render (plantilla) + persistencia (vía la mini-API).
- B3. Registrar el MCP en la config de Claude Code.
- B4. **Dependencia dura** (Q9-A): el flujo asume el MCP prendido; chequeo de arranque + mensaje claro si falta.

**C · SQLite (store/índice)**
- C1. Schema de tablas: `events` (log), `changes`, `agents`, `artifacts`, `code_refs`, `decisions_index` (EDR/spec) + relaciones.
- C2. Vistas/queries que derivan el estado para la lectura.
- C3. Re-indexado de EDR/spec `.md` (file-watch; gana el `.md`).
- C4. Identidad del change por `change-name`/rama; sesión nueva continúa en su fila (Q10).

**D · mini-API (broker)**
- D1. Único escritor de la SQLite (WAL); endpoints de escritura (MCP) y lectura (UI).
- D2. WebSocket para push a la UI.
- D3. Ingesta: hooks + tail del transcript → eventos.

**E · Observabilidad (ingesta)**
- E1. Hooks `SubagentStart/Stop` en `settings.json` → empujan a la mini-API.
- E2. Tail **incremental** del `.jsonl` (por offset) + parseo de acciones/MCP calls (Q5).
- E3. Mandato desde la línea 1 del transcript (Q6).

**F · UI mínima**
- F1. Shell de la app (dentro del repo).
- F2. Vista de un change: la **tarjeta del artefacto estructurado de `apply`** (diffs + justificación + refs).
- F3. Panel de flujo de agentes (roster + feed) — definir si entra en el mínimo (G3).
- F4. Consumo vía la mini-API (WebSocket).

**G · Transversal / setup**
- G1. Dónde vive la SQLite (¿`.matecito-ai/`?) — detalle de build.
- G2. El flujo asume MCP prendido (Q9-A) — chequeo de arranque.
- G3. Definir el "mínimo" del slice: ¿solo la tarjeta de `apply`, o también el panel de flujo?

---

## Glosario de decisiones tomadas en esta discusión

- El proyecto **no es "una UI"** — es un **contrato de salida estructurado** para el flujo SDD; la UI es consumidor río abajo. ✅
- **JSON = fuente de verdad**, markdown = proyección renderizada. ✅ (confirmado)
- Estructurar **alrededor de la prosa** (slots tipados + bolsillos de texto). ✅
- **MCP = validador/acelerador, no dependencia dura**; render local de fallback. ✅
- Tres capas: **verdad durable · estado vivo (JSON event-log) · broker (mini-API único escritor)**. ✅
- **Split hooks (mecánico) vs MCP (semántico/bidireccional)**. ✅
- Primera pieza = **vista de lectura / trazabilidad**, read-only. ✅
- Primera cosa a construir = **contrato estructurado de una fase**, no la UI. ✅
- **Write (gates, reiterar, canvas)** = norte, fase posterior. ✅
