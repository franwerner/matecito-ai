# ROADMAP — Broker/MCP (apps/api)

Plan de construcción del broker contra sus contratos: los **capability-specs** (`<repo-root>/.matecito-ai/development-specs/` — el *qué hace*) y los **EDRs** (`.matecito-ai/edr/` — el *cómo/por qué*). Antes de implementar cada ítem, leé el spec/EDR que se referencia en su fase.

**Convención de uso (para el agente):** marcá cada checkbox (`[ ]` → `[x]`) al completar y verificar el ítem, en el mismo commit/batch que lo completa. No marques nada no verificado. Si un ítem revela una decisión no cubierta por EDR/spec: pará y preguntá al usuario (no improvises — regla del proyecto). Las fases van en orden de dependencia; no arranques una fase con la anterior a medias salvo pedido explícito.

---

## Fase 1 — Fundaciones del daemon

Contratos: `edr/structure/*`, `edr/runtime/error-handling`, `edr/delivery/configuration`, `edr/observability/*`, `edr/delivery/ci-quality-gates`.

- [x] Estructura de carpetas según `edr/structure/folder-structure` y estilo según `edr/structure/architecture-style`.
- [x] Carga de configuración (flags/env) según `edr/delivery/configuration`.
- [x] Logging estructurado con correlación proyecto + change + agente (`edr/observability/logging`). Nota: helper de correlación construido y listo para conectar; los valores reales de proyecto/change llegan en Fase 3, el de agente en Fase 4.
- [x] Modelo de errores interno y traducción al contrato `{error, code}` sin internals (`edr/runtime/error-handling`).
- [x] Lifecycle del daemon: arranque, shutdown limpio, verificación de disponibilidad del MCP al arrancar (falla con mensaje claro si falta — `edr/contracts/mcp-server`). Nota: el check de disponibilidad MCP es un stub no-op (siempre éxito); la verificación real contra la superficie MCP se cablea en Fase 3.
- [x] Health check (`edr/observability/health-checks`). Nota: el campo de estado de persistencia es un flag stub siempre `ok`; se cablea contra el store real en Fase 2.
- [x] CI en GitHub Actions sobre PR: `golangci-lint` + `go test ./...` + build del binario (`edr/delivery/ci-quality-gates`).

## Fase 2 — Store (schema + borde de persistencia)

Contratos: `edr/data/data-modeling`, `edr/data/data-access-entity-framework`, `edr/data/storage-sync-model`. Tech: modernc sqlite + Ent + Atlas. Desglose detallado en [`ROADMAP-2.md`](ROADMAP-2.md).

- [ ] Migraciones goose con el schema inicial:
  - [ ] `projects` (UUID v7, identidad por `project-id`, lookup de paths por máquina, estado activo/inactivo).
  - [ ] `changes` (por proyecto, mapeo rama↔nombre, estado `active`/`closed`).
  - [ ] `events` (envelope `{type, payload}`, posición monótona por change, `occurred_at`, agente, idempotencia por hash del contenido canónico).
  - [ ] Store de contenido content-addressable (dedup por hash) compartido por versiones de records y fotos de código.
  - [ ] Versiones de records + pins de eventos a versiones (`lifecycle/record-version`).
  - [ ] Índice/proyección de records por `(proyecto, rama)` con owning-root y soft-delete.
  - [ ] Notas de iteración (estados `pendiente`/`entregada`/`resuelta`, anclaje change/evento/archivo+rango).
- [ ] Timestamps `created_at`/`updated_at` en todas las tablas; sin DELETE físico (excepción: notas `pendientes`).
- [ ] Borde de persistencia (Repository) según `edr/data/data-access-entity-framework`; schema definido en código y migraciones derivadas de él.
- [ ] Tests de integración con SQLite real (temporal/memoria), sin mocks (`edr/delivery/testing-strategy`).

## Fase 3 — Identidad y registro de proyectos

Contrato: `specs/flow/register-project`, `specs/rule/event-scoping`, `edr/contracts/mcp-server`.

- [ ] Superficie MCP esqueleto: streamable HTTP in-process (SDK Go oficial), scoping por header `X-Project-Path`, derivación server-side del change por rama.
- [ ] Normalización canónica de rutas (absoluta, symlinks, `.`/`..`) y resolución del toplevel git; si la ruta no resuelve a un toplevel git, rechazo `not_a_git_repo`.
- [ ] Resolución de identidad: `.id` presente → auto-link; ausente → `find_project(name)` + confirmación → LINK o CREATE materializando `.id`.
- [ ] Tool `register_project` idempotente + errores `invalid_path` y `not_a_git_repo`.
- [ ] Descubrimiento recursivo de `.matecito-ai/` (owning-roots) y alta del watch.
- [ ] Tests: escenarios del spec (create, link, auto-link, idempotente, anidada, cambio de ruta, rechazo no-git).

## Fase 4 — Changes y event-log (escritura MCP)

Contratos: `specs/flow/start-change`, `specs/flow/submit-phase-artifact`, `specs/lifecycle/change`.

- [ ] Tool `start_change` (alta/continuación, mapeo rama↔nombre).
- [ ] Tool `change_status` (`active` ↔ `closed`, única vía de transición).
- [ ] Guard `change_closed` en toda tool de escritura salvo `change_status`.
- [ ] Guard `branch_required` en toda tool de escritura, `change_status` incluida (rama derivada con `git symbolic-ref --short HEAD`; repo sin commits = rama usable, detached HEAD = rechazo).
- [ ] Tools `submit_<fase>` (×8) con validación de schema por fase; rechazo `validation_failed` todo-o-nada.
- [ ] Append al event-log con posición monótona, idempotencia por hash, sin validación de orden de fases.
- [ ] Pin implícito de la versión vigente de EDR/spec referenciados (congela la versión — `lifecycle/record-version`).
- [ ] Render del artefacto a forma canónica y respuesta al agente.
- [ ] Tests: escenarios de los specs (válido, inválido, duplicado, fuera de orden, closed, pin).

## Fase 5 — Indexado y versionado de records

Contratos: `specs/process/index-decision-records`, `specs/lifecycle/record-version`, `edr/data/storage-sync-model`.

- [ ] Build inicial al arrancar (recorre los `.md` bajo los owning-roots).
- [ ] File-watcher por proyecto + comando de sync/reindex on-demand (gana el archivo).
- [ ] Write-through: escritura del AI por el daemon = archivo + base en la misma operación.
- [ ] Extracción de proyección (slug, tipo, status, refs) + contenido; `.md` malformado se saltea y flaggea; refs colgadas se marcan sin romper.
- [ ] Versionado copy-on-write lazy: update-in-place si no referenciada, fork si congelada.
- [ ] Estado por `(proyecto, rama)`: soft-delete/reaparición, contenido compartido por hash entre ramas.
- [ ] Tests: escenarios del spec (build, write-through, checkout, borrado/reaparición, owning-root, dedup).

## Fase 6 — Lectura de la UI (OpenAPI + WebSocket)

Contratos: `specs/flow/serve-change-state`, `edr/contracts/api-contract`. Tech: huma (OpenAPI 3.1 Go-first) + coder/websocket.

- [ ] Endpoints OpenAPI de lectura: snapshot de change (estado + event-log con posición), listado de proyectos y de changes por proyecto.
- [ ] Hub de suscripciones WS por change (coder/websocket + golang-x-sync), push de eventos nuevos.
- [ ] Resume por posición del log: suscripción desde la posición declarada (sin huecos snapshot↔suscripción), reconexión sin re-snapshot, rechazo `invalid_position`.
- [ ] Eventos servidos por clase (fase con artefacto / mecánicos sin) en el envelope `{type, payload}`.
- [ ] Errores `change_not_found` / `project_not_registered` con `{error, code}`.
- [ ] Export del OpenAPI para Kubb (gate de sync de tipos en el CI de la UI).
- [ ] Tests: escenarios del spec (snapshot, push, huecos, reconexión, closed, multi-suscripción, scoping).

## Fase 7 — Ingesta mecánica + cliente de hooks

Contratos: `specs/process/ingest-mechanical-events`, `specs/rule/ingestion-spool`, `edr/runtime/resilience`.

- [ ] Endpoint HTTP de ingesta (separado de la superficie MCP), idempotente por hash, orden por `occurred_at`.
- [ ] Gate de ingesta: proyecto registrado + change `active`; descarte silencioso (no registrado / sin change / `closed`).
- [ ] Ingesta tolerante: stop huérfano, start duplicado.
- [ ] Cliente mínimo de hooks (binario/subcomando): reporte SubagentStart/Stop con presupuesto de 500 ms, fire-and-forget.
- [ ] Spool local con contenido y timestamp original; flush en próximo contacto (pendientes antes que lo nuevo), sin proceso residente, sin límites.
- [ ] Registración de los hooks en la configuración del runtime vía el installer (tocará `internal/` del root — anunciar antes de tocar).
- [ ] Tests: escenarios de los specs (ingesta viva, descartes, spool/reconciliación, idempotencia, nunca bloquea).

## Fase 8 — Fotos de código y vistas de diff

Contratos: `specs/process/capture-code-snapshots`, `specs/flow/serve-code-diff`, `edr/data/storage-sync-model`.

- [ ] Captura pre-edición vía hook sincrónico (PreToolUse): foto en el primer toque por archivo desde el último batch; dedup por hash; fallback spool con contenido al vencer los 500 ms.
- [ ] Captura post al submit del batch de apply (estado final de los archivos tocados) + referencia antes→después por archivo en el evento.
- [ ] Representación de ausencia: archivo creado (base = no existía) y borrado (después = ausencia).
- [ ] Endpoints de vistas de diff derivadas on demand: individual por batch, colapsada (base→última), contra archivo real (drift), bloques fuera-de-flujo; nunca se almacenan diffs.
- [ ] Rama sin checkout: vistas de fotos disponibles, vista contra disco informada como no disponible.
- [ ] Tests: escenarios de los specs (base, bordes de batch, edición manual encerrada, creado/borrado, vistas, spool de fotos).

## Fase 9 — Notas de iteración

Contratos: `specs/flow/iteration-notes`, `edr/contracts/mcp-server`.

- [ ] Escritura de notas desde la UI (endpoint): anclaje change → evento → archivo+rango; borrado real solo de `pendientes` por su autor.
- [ ] Señal en el hook de prompt: cantidad de pendientes, nunca contenido; skip silencioso al vencer presupuesto o broker caído.
- [ ] Tool MCP `list_pending_notes` (pasa las notas a `entregada`) y `resolve_note` (a `resuelta`, con vínculo trazado en el event-log).
- [ ] Eventos de creación/resolución de nota en el event-log (llegan en vivo por la suscripción).
- [ ] Tests: escenarios del spec (creación anclada, señal, no-lee-sin-confirmar, entrega, resolución, anclaje estable, multi-sesión).

## Fase 10 — Endurecimiento y cierre

Contratos: `edr/delivery/*`, `edr/observability/*`, todos los specs (verificación final).

- [ ] Reconexión/backoff del lado servidor verificados contra `edr/runtime/resilience` (presupuestos, backoff 1 s→30 s).
- [ ] Pasada de verificación specs↔implementación: cada escenario `Accepted` tiene test o verificación manual anotada.
- [ ] Pasada de reglas verificables de los EDRs del broker.
- [ ] Release: bundle de la UI construido antes del build de Go y embebido en el binario; el release falla si el bundle no construye (`edr/delivery/ci-quality-gates` + goreleaser).
- [ ] `golangci-lint` limpio y suite verde en CI.
- [ ] Documentación mínima de operación del daemon (arranque, config, salud).
