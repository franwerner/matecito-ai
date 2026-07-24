# Cockpit — preguntas abiertas

> Cabos pendientes de la discusión de [`cockpit.md`](cockpit.md). Los resolvemos **una a la vez**.
> Cada pregunta tiene un **Estado** (`Pendiente` / `Resuelta`) y su **Resolución** cuando se cierra.

## Micro-decisiones

### Q1 — Event-log vs snapshot para el JSON
El estado vivo de sesión, ¿es un **event-log** (appendeás eventos, la API deriva el estado; da timeline/diff/`supersedes` gratis) o un **snapshot** (un JSON = estado actual, se reescribe; más simple)?
- **Recomendación:** event-log. El doc ya lo asume.
- **Estado:** Resuelta
- **Resolución:** **Event-log.** Appendeás eventos y la mini-API deriva el estado; da timeline (scrubber), diff v1→v2 y `supersedes` gratis, y aguanta mejor la concurrencia. **Medio: SQLite autocontenido** (ver Q11) — la tabla `events` es el log; el estado derivado son queries.

### Q2 — Puerta de entrada: codebase permanente vs change en vuelo
¿El cockpit abre sobre el **codebase como mapa permanente** (con los changes adentro como historial, apoyado en codegraph) o solo sobre el **change en vuelo**?
- **Estado:** Resuelta
- **Resolución:** **Change en vuelo.** La vista se arma con los archivos que tocó *esta* corrida (referenciados por ruta + diff + procedencia al agente/EDR), no con el árbol entero del proyecto. El mapa permanente del codebase queda como evolución futura.

### Q3 — ¿La capa de cumplimiento por EDR entra en la primera pieza?
Meterla (§5.2) implica pedirle a `verify` que la emita **estructurada** — o sea toca el flujo ya en la primera pieza, no es solo UI. ¿Entra ahora o después?
- **Estado:** Resuelta
- **Resolución:** Entra, pero **sin el %** (juicio difuso, descartado). Queda: en `verify`, la AI marca las partes del código que tocó con `CRITICAL`/`WARNING`/`SUGGESTION` + el **problema encontrado**, contra sus dos contratos (**spec** y **EDR**). Es lo que `verify` ya juzga (Decision Gaps/severidades), ahora **estructurado y anclado al código**. Se apoya en que cada artefacto de código ya lleva su **justificación** (de `apply`). El orden de construcción lo define Q8.

### Q4 — Nombre del proyecto
"cockpit" es tentativo.
- **Estado:** Resuelta
- **Resolución:** **Matecito UI.** Vive **dentro de este mismo repo** (`matecito-ai`).

## A validar empíricamente (no bloquean, no cambian el diseño)

### Q5 — Cadencia de flush del transcript
¿Cada cuánto flushea el `.jsonl` del sub-agente? Define cuán "en vivo" se ve el feed. Experimento de ~10 min.
- **Estado:** Resuelta (medido)
- **Resolución:** **Casi en vivo, a nivel de evento.** Medido con un agente descartable (~47s) sondeando el transcript cada 1s: se escribe **incremental durante toda la corrida** (de a 1-3 líneas cada ~1-4s, siguiendo cada acción), **nunca buffered-al-final** ni token-a-token. Latencia del feed ~1-2s. Nota práctica: el archivo es grande (~150KB), el tail debe ser **incremental** (por offset/líneas nuevas), no re-leer entero.

### Q6 — De dónde se lee el mandato para la arista
¿`SubagentStart` trae el prompt/`subagent_type`, o se lee de las primeras líneas del transcript?
- **Estado:** Resuelta (verificado)
- **Resolución:** **Se lee de la línea 1 del transcript.** Confirmado: el mandato aparece en la línea 1 del `.jsonl`, presente desde `t=0` (el archivo ya tiene 6 líneas al lanzar). No se depende del payload del hook (que quedó no documentado).

## Decisiones grandes / EDRs (para cuando se formalice)

### Q7 — Migración markdown → JSON
Cómo pasan los EDRs/specs de markdown-canónico a JSON-fuente-de-verdad. ¿Coexisten? ¿Se migran? ¿Big-bang o gradual?
- **Estado:** Resuelta
- **Resolución:** **Greenfield + índice.** (1) Artefactos de flujo: sin back-migración — nacen estructurados de acá en adelante; los viejos en prosa quedan como historial. (2) **EDR/spec NO se migran ni se mueven a SQLite**: siguen **canónicos en `.md`** (git, PR-review, edición a mano, grep) y Matecito UI los **indexa** en SQLite (id/status/refs + relaciones) — mismo patrón que `codegraph` indexa código que vive en archivos. La SQLite es índice/proyección, no fuente de verdad de EDR/spec; si divergen, **gana el `.md`** (se re-indexa vía el file-watch de la mini-API).
- **Actualización (2026-07-24) — de índice a contenido versionado:** el broker no guarda **solo** la proyección (id/status/refs): guarda el **contenido de cada EDR/spec versionado** en la DB, con **copy-on-write lazy**. Un evento referencia siempre la versión **vigente** (implícita, el agente no manda versión); una versión solo se **congela y forkea** cuando el `.md` cambia *y* esa versión ya estaba referenciada por un evento (si no, se actualiza en el lugar → cero versiones de más). Motivo: la UI debe poder mostrar un EDR/spec en la **versión exacta que aplicó un evento pasado**, aun si el `.md` cambió o se borró. El `.md` sigue **canónico para editar** (git); el broker versiona para el pin/traza. El **contenido versionado vive en la DB self-contained** (content-addressable, dedup por hash) — **NO derivado de git**: git es local en cada máquina y bloquearía el futuro **deploy/compartir con compañeros** (un servidor no tiene el checkout de nadie). Modelo completo (dueños, write-through, watch, lectura del AI sobre archivos) en el EDR `data/storage-sync-model` + los specs `index-decision-records` / `submit-phase-artifact`.

### Q8 — Alcance de la re-plomería
¿Estructuramos la salida de **todas** las fases o arrancamos con **una sola** (slice vertical)?
- **Estado:** Resuelta (falta elegir la fase del slice)
- **Resolución:** **Vertical slice** — estructurar UNA fase de punta a punta (contrato → MCP → SQLite → UI mínima) para validar el lazo; el resto se replica. El **roadmap de "lo que viene después"** quedó asentado en `cockpit.md` §9 (FASE 0 slice → FASE 1 escalar lectura + indexar EDR/spec → FASE 2 cumplimiento en verify → FASE 3 escritura). **Slice (FASE 0) = `apply`** (el código: diffs + justificación + refs; el artefacto más valioso y el que ejercita las partes difíciles). `tasks` y el resto quedan para FASE 1.

### Q9 — Render local de fallback
Mecanismo para que el flujo renderice artefactos **sin** el MCP (preservar "offer, don't impose"). Decidido en principio, mecanismo sin definir.
- **Estado:** Resuelta
- **Resolución:** **Opción A — MCP dependencia dura de TODO el flujo. Sin fallback.** El MCP **debe estar corriendo**; sin él, el flujo SDD no corre. **Revierte conscientemente la invariante `offer, don't impose`** → **requiere su propio EDR**. Costo asumido: el flujo requiere el MCP prendido, y las corridas headless/CI deben garantizar su disponibilidad. Beneficio: un solo camino de render, cero lógica de degradación.

## Agujero de diseño abierto

### Q10 — Reconciliación multi-sesión
El JSON vivo es **por sesión**, pero un change puede cruzar varias. La verdad durable persiste, pero cómo se **reengancha el overlay vivo** al retomar un change en otra sesión no está resuelto.
- **Estado:** Resuelta
- **Resolución:** **Disuelto por Q11** — la SQLite ya es **persistente y per-proyecto**, así que un change que cruza sesiones sigue appendeando al mismo store; no hay overlay que reenganchar. Lo único que queda: una sesión nueva **identifica el change por lo que ya existe** (`change-name` / rama de git / estado de flujo en `.matecito-ai/`) y **continúa en su fila**; si no hay change activo, es uno nuevo. No se inventa un identificador nuevo.
- **Actualización (2026-07-23):** con el store global (ver Q11 actualizada), la SQLite deja de ser per-proyecto: es **única y compartida en `~/.matecito-ai/`**. La identidad del change pasa a **`project_id + change-name`/rama** (el proyecto se registra por su ruta). La reconciliación multi-sesión sigue válida, ahora scopeada por proyecto.

## Emergentes (surgieron durante la resolución)

### Q11 — Medio de persistencia del estado: SQLite autocontenido
Surgió al resolver Q5. En vez de archivos JSON, el estado de Matecito UI vive en una **SQLite autocontenida** (un archivo, sin server aparte, sin infra pesada).
- **Estado:** Resuelta
- **Resolución:** **SQLite para el store estructurado de Matecito UI** (event-log + artefactos de flujo + índice del change/agentes/código/decisiones). Persistente y **queryable**; la mini-API es el único escritor (SQLite en WAL: 1 escritor + N lectores). **Engram se queda** (memoria semántica, difusa, cross-proyecto) — **NO se reemplaza**; conviven, son datos distintos con trabajos distintos. **EDRs/specs siguen en `.md`.** Reemplazar Engram, si alguna vez, es una decisión **aparte**, no acoplada a Matecito UI. **Refinamiento (ver Q7):** los EDR/spec siguen `.md` canónicos y se **indexan** en SQLite (no se mueven) — patrón codegraph: archivo canónico, índice consultable en la DB.
- **Actualización (2026-07-23):** el store pasa a **instancia única global** — **un daemon broker único por máquina** + **una SQLite única en `~/.matecito-ai/`** compartida entre todos los proyectos (revierte el modelo per-proyecto de Q10/Q11). Cada proyecto se **registra por su ruta** (vía MCP o UI) y toda la data se scopea por `project_id`. Motivo: un cockpit y un almacenamiento único que sirvan a todos los proyectos. El invariante único-escritor se mantiene porque hay **un solo proceso** escribiendo. Ver EDR `deployment-topology` en `apps/api/.matecito-ai/edr/delivery/`.

---

## Próximo paso (después de resolver las 10)
Decidir si se **formaliza al flujo** (`intake`) para estructurar el cockpit como cambio real, y ahí elegir lane y modo.
