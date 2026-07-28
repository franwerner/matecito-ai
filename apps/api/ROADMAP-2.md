# ROADMAP-2 — Fase 2: Store (schema + borde de persistencia)

Desglose detallado de la **Fase 2** de [`ROADMAP.md`](ROADMAP.md). Se construye en **4 batches sobre la misma rama**, sin PRs: al cerrar cada batch el usuario revisa el diff acumulado antes de arrancar el siguiente.

**Contratos que gobiernan esta fase.** `edr/data/data-modeling` · `edr/data/data-access-entity-framework` · `edr/data/storage-sync-model` · `edr/delivery/testing-strategy` · `edr/observability/health-checks` · `edr/tech/*`. Comportamiento: `<repo-root>/.matecito-ai/development-specs/` — `lifecycle/record-version`, `process/index-decision-records`, `flow/iteration-notes`, `rule/event-scoping`, `lifecycle/change`.

**Convención de uso.** Marcá cada checkbox (`[ ]` → `[x]`) al completar y verificar el ítem, en el mismo batch que lo completa. No marques nada no verificado. Si un ítem revela una decisión no cubierta por el anexo de contratos de abajo: pará y preguntá al usuario. Al cerrar el batch 4, marcá los checkboxes de la Fase 2 en `ROADMAP.md`.

---

## Anexo — Contratos de schema confirmados

Decididos con el usuario antes de implementar. **No re-inferir ni improvisar sobre estos puntos**; cualquier desvío se consulta.

### Convenciones transversales

- **IDs:** `TEXT` con UUID v7 canónico con guiones (legible en base y en logs de correlación).
- **Timestamps** (`created_at`, `updated_at`, `occurred_at`, `frozen_at`): `TEXT` RFC 3339 UTC con milisegundos (`2026-07-27T18:06:33.123Z`) — ordena lexicográficamente igual que cronológicamente.
- **Enums de estado:** `TEXT` con `CHECK (col IN (...))`. **Valores siempre en inglés**, aunque el capability-spec que los describe esté en castellano.
- **Flags booleanos:** `INTEGER NOT NULL DEFAULT 0`.
- **`created_at` / `updated_at` NOT NULL en todas las tablas.** Sin DELETE físico; única excepción: notas de iteración en estado `pending`.
- **`project_id` denormalizado** en toda tabla de contenido que de otro modo llegaría al proyecto solo por join — cumple literal la regla verificable 32 de `data-modeling` y hace imposible cruzar proyectos por olvidar un join.

### Por tabla

- **`projects`** — identidad por `project-id` (UUID v7, el mismo del archivo `.id`); `name` **sin unicidad** (`find_project` puede devolver varios y el usuario confirma); `status` `'active'|'inactive'`.
- **`project_paths`** — tabla aparte: N paths por (proyecto, máquina), para worktrees y clones múltiples. `machine_id` = UUID v7 generado una vez y persistido en el persist-dir del daemon (`~/.matecito-ai/`), **nunca el hostname** (cambia al renombrar, colisiona entre máquinas). `root_path` canónico (absoluto, symlinks resueltos, sin `.`/`..`). `UNIQUE (machine_id, root_path)`. Repo movido → la fila vieja pasa a `'stale'` y se inserta la nueva.
- **`changes`** — `branch TEXT NOT NULL` con un único `UNIQUE (project_id, branch)`. **No existe el fallback sin rama** (enmienda del 2026-07-27: un change requiere siempre rama git). `UNIQUE (project_id, name)` — el nombre es único por proyecto, para no colisionar con los topic keys de Engram `sdd/{change-name}/…`. `status` `'active'|'closed'`.
- **`events`** — envelope `{type, payload}`; `payload TEXT` con **JSON canónico, exactamente los bytes que se hashean, nunca re-serializados**; `position` monótona por change; `occurred_at` propio; `agent_kind` (`phase|orchestrator|subagent`), `agent_name` y `session_id` como **texto nullable, sin tabla `agents`** — un agente es hoy solo una etiqueta, sin atributos ni ciclo de vida, y normalizarla costaría un upsert dentro de la transacción de cada append; si gana atributos, se promueve a tabla y se backfillea desde los eventos. Sus valores reales llegan en Fase 3-4; `UNIQUE (change_id, position)`; `UNIQUE (idempotency_hash)` global — seguro porque el anclaje (change, fase) entra al hash per `data-modeling`.
- **`content_objects`** — `id` propio como PK más `UNIQUE (project_id, hash)` (ver la nota de claves compuestas al final del anexo); el dedup es **entre ramas dentro del proyecto**, no entre proyectos. `hash TEXT` con prefijo de algoritmo inline: `'sha256:<hex minúscula>'`, para poder migrar de algoritmo sin ambigüedad. SHA-256 de la stdlib (sin dependencia nueva). `bytes BLOB` crudo, **sin compresión** (fuera de alcance; se agregaría con una columna `encoding` en un `ALTER`). Store único compartido por versiones de records y fotos de código, **sin columna discriminadora de tipo** (rompería el dedup entre ambos).
- **`records`** — identidad = slug + `owning_root` derivado del path: `UNIQUE (project_id, owning_root, slug)`. Branch-independiente.
- **`record_versions`** — **no scopeada por rama**. `status` `'current'|'frozen'`. Lleva la **proyección** (`kind`, `doc_status`, …), porque `record-version.md` define la versión como "una foto del contenido **y la proyección**" y por lo tanto se congela con ella. **NO existe un índice único de "una sola `current` por record"**: dos ramas con contenido divergente tienen dos versiones vigentes simultáneas del mismo record; la vigencia por rama la expresa `record_branch_state.current_version_id`.
- **`record_version_refs`** — una fila por línea de `## Relacionados` / `## Referencias`: `relation`, `target_slug`, `target_owning_root`, `dangling`. El destino se guarda como **texto, sin FK al record destino** — los destinos se borran, reaparecen y se renombran, y con texto ese movimiento solo recalcula el flag en vez de re-apuntar punteros.
- **`record_branch_state`** — el estado por `(record, rama)`, con `id` propio como PK más `UNIQUE (record_id, branch)`: `status` `'active'|'deleted'` (soft-delete), `current_version_id`, y el flag **`malformed`** de `.md` que no se pudo parsear (vive acá porque es la única fila por rama: un merge conflictivo rompe el archivo en una rama concreta, y la entrada sigue apuntando a la última versión sana).
- **`event_record_pins`** — `id` propio como PK más `UNIQUE (event_id, record_version_id)`.
- **`iteration_notes`** — `status` `'pending'|'delivered'|'resolved'`, sin vuelta atrás. **Sin columna `author`**: el sistema es single-user por EDR y el "su autor" del contrato contrasta usuario-vs-AI, no persona-vs-persona; la regla real es que solo las `pending` admiten borrado físico. CHECKs de anclaje: rango de líneas solo si hay `file_path`; `file_path` solo si hay `event_id`; `line_from <= line_to`.

### Nota — claves compuestas

Tres tablas (`content_objects`, `record_branch_state`, `event_record_pins`) se identifican naturalmente por una clave compuesta, pero **llevan un `id` propio como clave primaria y expresan la clave natural como índice único**. El motivo es una limitación verificada de la herramienta: en Ent v0.14.6 la clave primaria compuesta solo funciona en tablas intermedias de relaciones muchos-a-muchos, y en una entidad común rompe la generación de código.

**La garantía de integridad es la misma** — el índice único impide los duplicados exactamente igual que lo haría la clave primaria—, así que el dedup por contenido, el estado por rama y el pineo de versiones no cambian en nada. El costo es una columna por fila en esas tres tablas.

### Reglas de comportamiento que el store debe respetar

- **Copy-on-write lazy con DOS condiciones de fork** (`lifecycle/record-version`, enmendado 2026-07-27): el update-in-place procede solo si la versión **no está congelada** Y **ninguna otra rama la apunta**. Si está congelada por un pin de evento **o** la comparte otra rama → fork a una versión nueva para la rama que cambió, dejando la previa intacta.
- **Idempotencia por hash del contenido canónico** (`data-modeling`): para un submit, el hash del artefacto canónico **más su anclaje (change, fase)**, excluyendo los campos que estampa el servidor. Ningún productor genera claves propias.
- **Health** (`observability/health-checks`): 200 **solo** con la persistencia abierta y migrada.

---

## Batch 1 — Tooling y schema completo

Deja la base creable desde cero y migrada al arrancar. Sin lógica de dominio todavía.

- [x] Instalar las dependencias con las versiones ya resueltas: `entgo.io/ent` v0.14.6, `ariga.io/atlas` v1.2.3 y `modernc.org/sqlite` v1.54.0 (driver Go puro sin cgo, embebe SQLite 3.53.3).
- [x] Configuración de la generación de código de Ent: ubicación del schema y del código generado según `edr/structure/folder-structure` (`internal/store/ent/schema/**`, código generado en `internal/store/ent/**`).
- [x] ~~Confirmar que las claves primarias compuestas funcionan en entidades comunes.~~ **Verificado: NO funcionan.** En Ent v0.14.6 la anotación de ID compuesto es válida solo en tablas intermedias de relaciones muchos-a-muchos; en una entidad común rompe la generación de código. Resuelto con `id` propio más índice único sobre la clave natural en las tres tablas afectadas (ver la nota de claves compuestas en el anexo).
- [x] Subir el `go` / `toolchain` de `go.mod` a la versión que exige `ariga.io/atlas` v1.2.3 (>= 1.26.4); `go.mod` quedó en `go 1.26.4` con línea `toolchain go1.26.4` explícita y coherente (evita el antecedente de la Fase 1).
- [x] Schema en código — identidad y flujo: `projects`, `project_paths`, `changes`.
- [x] Schema en código — event-log y contenido: `events`, `content_objects`.
- [x] Schema en código — records: `records`, `record_versions`, `record_version_refs`, `record_branch_state`, `event_record_pins`. `record_versions` referencia el contenido por `content_object_id → content_objects(id)` (relación normal, resuelto durante el batch — ver nota de claves compuestas).
- [x] Schema en código — notas: `iteration_notes`, con sus CHECKs de anclaje declarados en el propio schema.
- [x] Aplicar las convenciones transversales del anexo a **las 11 tablas** (tipos, timestamps, enums en inglés, flags, `project_id` denormalizado, FKs activas con `PRAGMA foreign_keys = ON` vía `_pragma=foreign_keys(1)` en el DSN de `modernc.org/sqlite`).
- [x] Generación y persistencia del `machine_id` (UUID v7) en el persist-dir del daemon; se crea al primer arranque y se reusa (`internal/store/machineid.go`, wireado en `cmd/broker/lifecycle.go`).
- [x] Generar con Atlas las **migraciones versionadas** derivadas del schema, embeberlas en el binario (`embed.FS`) y aplicarlas al arranque en modo librería. Las migraciones se derivan, **nunca se editan a mano**: si falta un constraint, se corrige en el schema y se regenera. (`internal/store/migrations/generate/main.go` deriva; `internal/store/migrations/embed.go` embebe; `internal/store/store.go` aplica vía el `migrate.Executor` puro-Go de Atlas — sin CLI, sin cgo.)
- [x] Apertura de la base en el arranque del daemon + ejecución de migraciones pendientes; fallo con mensaje claro si no puede. (`store.Open` en `cmd/broker/lifecycle.go`, antes del check de MCP y del bind del transporte.)
- [x] Test de integración: base efímera creada desde cero, migraciones aplicadas, schema verificado (tablas, índices y CHECKs presentes). (`internal/store/store_test.go` — schema completo, idempotencia de arranque, FKs enforced, y las 10 CHECKs de la tabla rechazando INSERT crudo + casos de control positivos.)

## Batch 2 — Borde de persistencia e identidad

Establece el patrón de Repository que los batches siguientes replican.

- [ ] Borde de persistencia (Repository) según `edr/data/data-access-entity-framework`: interfaz angosta, sin filtrar hacia adentro los tipos generados ni los de la base.
- [ ] Traducción de errores del store al modelo de errores interno de la Fase 1 (sin exponer internals; el borde HTTP ya los mapea a `{error, code}`).
- [ ] Manejo de transacción del único escritor (el daemon es global y único por máquina).
- [ ] Repository de `projects`: alta, lookup por `project-id`, búsqueda por nombre (varios resultados posibles), baja lógica a `inactive`.
- [ ] Repository de `project_paths`: alta, resolución de path canónico → `project-id` por máquina, transición de la fila vieja a `stale` al mover el repo.
- [ ] Repository de `changes`: alta/continuación por `(project_id, branch)`, transición `active` ↔ `closed`, rechazo de nombre duplicado en el proyecto.
- [ ] Tests de integración con SQLite real (sin mocks): alta y lookup de proyecto, path por máquina, repo movido, continuación de change sobre la misma rama, unicidad de nombre de change.

## Batch 3 — Event-log, contenido y notas

- [ ] Append de eventos con **posición monótona por change**, asignada server-side y sin huecos.
- [ ] Idempotencia por `idempotency_hash`: un reintento con el mismo contenido canónico no duplica el evento y devuelve el existente.
- [ ] Cálculo del hash canónico (artefacto + anclaje change/fase, excluyendo campos estampados por el servidor), con el `payload` persistido byte a byte tal como se hasheó.
- [ ] Lectura del event-log por change: snapshot completo y lectura desde una posición dada (el insumo del resume de la Fase 6).
- [ ] Repository de `content_objects`: `put` idempotente por `(project_id, hash)` que no reescribe si ya existe, y `get` por hash. Prefijo `sha256:` construido y validado en un solo lugar.
- [ ] Repository de `iteration_notes`: creación anclada, transiciones `pending → delivered → resolved`, borrado físico permitido **solo** en `pending`.
- [ ] Tests de integración: posición monótona bajo escrituras sucesivas, idempotencia (mismo contenido dos veces = un evento), dedup de contenido por hash, CHECKs de anclaje rechazando combinaciones inválidas, borrado de nota permitido solo en `pending`.

## Batch 4 — Records, versionado, health y cierre

- [ ] Repository de `records`: upsert por `(project_id, owning_root, slug)`.
- [ ] Repository de `record_versions` con la **proyección** y el contenido referenciado por hash.
- [ ] **Copy-on-write lazy con las dos condiciones de fork**: update-in-place solo si la versión no está congelada Y ninguna otra rama la apunta; en cualquier otro caso, fork. Es el ítem más delicado de la fase.
- [ ] Congelamiento (`current → frozen`) disparado por el pin de un evento, sin vuelta atrás.
- [ ] Repository de `record_branch_state`: estado `active`/`deleted` por rama, versión vigente por rama, flag `malformed`.
- [ ] Repository de `record_version_refs`: alta de refs de una versión y recálculo del flag `dangling` cuando un destino aparece o desaparece.
- [ ] Repository de `event_record_pins`: pin de un evento a la versión vigente al momento de persistir.
- [ ] **Stub C — health real**: `StoreStatus()` cableado contra el store. 200 con `{"process":"ok","store":"ok"}` solo con la persistencia abierta y migrada; **503** en otro caso con `store` en `unavailable` (base no abierta), `migrating` (migraciones corriendo) u `outdated` (pendientes o fallidas). Eliminar el stub siempre-`ok` de la Fase 1.
- [ ] Tests de integración: update-in-place sin referencias, fork por versión congelada, **fork por versión compartida con otra rama** (el escenario nuevo del spec enmendado), soft-delete y reaparición por rama, contenido compartido por hash entre ramas, ref colgada marcada sin romper el índice, `.md` malformado flaggeado por rama, health en sus cuatro estados.
- [ ] Pasada final: `golangci-lint` limpio, `go vet`, `gofmt`, suite completa verde con `-race`.
- [ ] Marcar los checkboxes de la **Fase 2** en [`ROADMAP.md`](ROADMAP.md).
