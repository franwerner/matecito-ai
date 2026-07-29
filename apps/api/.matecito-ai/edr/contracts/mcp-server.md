# EDR — MCP server (superficie cara a Claude)

- **Status:** Accepted
- **Date:** 2026-07-27

## Contexto

El cockpit necesita una superficie de escritura que Claude Code consuma para materializar cada fase del flujo SDD como salida estructurada. Esa superficie es un MCP server, y su diseño está condicionado por dos hechos: el modelo de despliegue del broker (un único daemon global por máquina, con el store embebido) y los límites reales de identidad que Claude Code expone al servidor MCP. En particular, Claude Code no reenvía al MCP el identificador de sesión del protocolo, e interpola variables de entorno en los headers de su configuración de servidores MCP; el scoping por request tiene que apoyarse en lo que sí viaja, no en lo que se desearía tener. Además el flujo depende del MCP: sin él prendido, no hay dónde escribir las fases.

## Decisión

Exponemos el MCP como una superficie **in-process** dentro del mismo binario que aloja el broker y el store: un único daemon global por máquina lo sirve, y todas las sesiones de Claude Code de esa máquina consumen del mismo proceso. No hay un subproceso MCP por sesión.

- **Transporte:** streamable HTTP (el HTTP+SSE anterior está deprecado en la spec MCP; en streamable HTTP el SSE es un modo de respuesta opcional dentro de un endpoint único). Se implementa sobre el SDK Go oficial de MCP.
- **Identidad y scoping por request**, diseñado alrededor de los límites de Claude Code:
  - El **proyecto** llega en cada request por un header que transporta la ruta del directorio de proyecto, inyectado en la configuración de servidores MCP de Claude Code mediante la interpolación de la variable de entorno del directorio de proyecto. Viaja por request, por proceso.
  - El **change** se deriva server-side de la rama de git del proyecto: el daemon lee la rama del directorio de proyecto. Los sub-agentes no precisan conocer el nombre del change; el orquestador registra el mapeo rama↔nombre una sola vez al arrancar el flujo.
  - La **sesión** es un atributo del evento, no una clave de scoping. El MCP no la recibe (Claude Code no reenvía el header de sesión del protocolo); la identidad de sesión entra por el canal de hooks, no por el MCP.
  - **Invariante de identidad:** un change se corresponde con un working directory (un worktree) y una rama. Features en paralelo son worktrees distintos, es decir directorios de proyecto distintos. El mismo directorio y la misma rama son el mismo change; el atributo de sesión distingue corridas dentro de ese change.
- **Tools expuestas:** `register_project` para dar de alta (o linkear) un proyecto por su ruta, `find_project(name)` para buscar un proyecto por nombre en la base —el lookup que permite linkear un checkout sin `.id` a un proyecto ya existente—, `start_change` para dar de alta un change registrando el mapeo rama↔nombre, `change_status` que transiciona el estado del change entre `active` y `closed` (la única vía de cambio de estado, explícita y nunca efecto lateral de otra tool), **una tool de submit por cada fase del flujo** — intake, propuesta, spec, diseño, tareas, apply, verify y archivo — cada una validando su carga contra el schema de artefacto de esa fase, y las dos tools de **notas de iteración**: `list_pending_notes` (lee las notas pendientes del change — solo tras la confirmación del usuario en el chat; la señal de que existen llega por el canal de hooks, nunca con contenido) y `resolve_note` (marca una nota atendida, dejando el vínculo trazado en el event-log).
- **Identidad del proyecto:** la identidad estable de un proyecto es el `project-id` de `<git-toplevel>/.matecito-ai/.id` —un archivo commiteado que viaja con el repo—, **no el path local**. El path que llega en el header es un lookup por-máquina que el daemon resuelve a ese `project-id`; la base keyea por `project-id`, no por ruta.
- **Gates bloqueantes:** no existen en esta superficie. Los gates del flujo (intake gate, checkpoints, gates de candidatos) se resuelven siempre en la conversación del chat; la UI es visibilidad —renderiza artefactos y event-log— y a lo sumo inicia una iteración que se conversa en el chat. No hay mecanismo de retención de request abierta ni rendezvous UI↔MCP, ni previsto a futuro.
- **Dependencia dura:** el flujo SDD asume el MCP prendido. El arranque del daemon verifica su disponibilidad y falla con un mensaje claro si falta.

## Alcance

- `apps/api/internal/transport/**` — el borde que aloja la superficie MCP/HTTP.

## Reglas verificables

- **[manual]** un único daemon global expone el MCP in-process; las sesiones consumen del mismo proceso, sin subproceso MCP por sesión.
- **[manual]** el transporte es streamable HTTP.
- **[manual]** el scoping por request usa el header `X-Project-Path` (proyecto) más la rama de git derivada server-side (change); nunca depende del identificador de sesión del protocolo MCP.
- **[manual]** los sub-agentes no reciben el nombre del change; se deriva de la rama, y el orquestador registra el mapeo rama↔nombre una vez vía la tool de alta de change.
- **[manual]** hay una tool de submit por fase, y cada una valida el schema de artefacto de su fase.
- **[manual]** las notas de iteración se leen solo vía `list_pending_notes` (tras confirmación del usuario) y se resuelven solo vía `resolve_note`; la señal de pendientes viaja por el canal de hooks sin contenido.
- **[manual]** el estado de un change (`active` ↔ `closed`) solo se transiciona vía la tool `change_status`; ninguna otra tool lo cambia como efecto lateral.
- **[manual]** cualquier tool que opere sobre un change `closed` (salvo `change_status`, que puede reabrirlo) se rechaza con `change_closed`.
- **[manual]** una escritura se rechaza si el proyecto no fue registrado previamente.
- **[manual]** la identidad del proyecto es el `project-id` de `<git-toplevel>/.matecito-ai/.id`; el path del header se resuelve a ese `project-id` y la base keyea por `project-id`, nunca por el path local.
- **[manual]** `find_project(name)` resuelve un proyecto por nombre en la base, para linkear un checkout sin `.id` a un `project-id` existente.
- **[manual]** `register_project` rechaza con `not_a_git_repo` una ruta que no resuelve a un toplevel git; no hay registración fuera de un repo git.
- **[manual]** toda tool de escritura —`change_status` incluida— se rechaza con `branch_required` cuando el proyecto no está sobre una rama (detached HEAD); la rama se deriva con `git symbolic-ref --short HEAD`, que devuelve el nombre incluso en un repo sin commits (rama usable) y falla en detached HEAD.
- **[manual]** el arranque del daemon verifica la disponibilidad del MCP y falla con un mensaje claro si no está.

## Alternativas consideradas

- **Un subproceso MCP por sesión de Claude Code:** descartado; multiplicaría procesos contra un store único y global, y rompería el modelo de daemon único. El scoping por request sobre un proceso compartido cubre el aislamiento necesario.
- **Transporte HTTP+SSE clásico:** descartado; está deprecado en la spec MCP a favor de streamable HTTP.
- **Usar el identificador de sesión del protocolo como clave de scoping:** inviable; Claude Code no lo reenvía al MCP. Por eso el scoping se apoya en el header de proyecto y la rama derivada.

## Consecuencias

- Un solo proceso concentra broker, store y superficie MCP: menos partes móviles y una única fuente de verdad por máquina.
- Los sub-agentes quedan con cero awareness del change: el mapeo vive una vez en el orquestador y todo lo demás se deriva de la rama.
- Trade-off: al derivar el change de la rama, el invariante un-change↔un-worktree↔una-rama es obligatorio; trabajar dos changes sobre el mismo directorio y rama los colapsaría en uno.
- Los gates del flujo viven exclusivamente en el chat: la superficie MCP no retiene requests ni media confirmaciones interactivas, y la UI nunca resuelve un gate.

## Relacionados

- `relacionado-con` → [api-contract.md](api-contract.md) — separa esta superficie de escritura (MCP) de la superficie de lectura de la UI.
- `depende-de` → [../delivery/deployment-topology.md](../delivery/deployment-topology.md) — el daemon global único es la premisa de que el MCP sea in-process y compartido.
- `relacionado-con` → [../observability/logging.md](../observability/logging.md) — la correlación por proyecto + change + agente se apoya en esta identidad por request.
