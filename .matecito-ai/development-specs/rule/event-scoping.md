# Capability — Scoping de eventos a (proyecto, change)

- **Status:** Accepted
- **Date:** 2026-07-27

## Propósito

Definir cómo cada evento —venga de una tool MCP o del canal de hooks— se asocia a un (proyecto, change), para que las dos fuentes converjan sobre la misma corrida y la UI pueda mostrarla unificada.

## Actores

- **La superficie MCP** (canal semántico) — cada request de tool aporta un evento a scopear.
- **El canal de hooks** (canal mecánico) — cada hook aporta un evento a scopear, con el session id de Claude Code como atributo.

## Reglas de negocio

- **Proyecto** ← se identifica por el **`project-id` estable** de `<git-toplevel>/.matecito-ai/.id` (commiteado, machine-independent: sobrevive a mover o renombrar el repo y sirve deployado), **no por el path local**. El header `X-Project-Path: ${CLAUDE_PROJECT_DIR}` de la request MCP (y el `cwd` para los hooks) es el **lookup por-máquina** que el daemon **resuelve a ese `project-id`**.
- **Change** (modelo M1) ← se deriva de la rama git → key (proyecto, rama), **siempre**. No hay otra key: un change está sujeto a una rama. El sub-agente nunca provee el change.
- **Rama obligatoria:** sin rama no hay change posible, y **toda tool de escritura se rechaza con `branch_required`** (incluida `change_status`). El único caso real "sin rama" es un **detached HEAD**; un repo git recién inicializado **sin commits sí tiene rama usable** (su rama existe como nombre aunque todavía no tenga ref) y se usa normalmente como key. El guard **no** alcanza a las otras superficies: la lectura de la UI sigue sirviendo historial y la ingesta mecánica sigue descartando en silencio.
- **Session** ← el session id de Claude Code es un **atributo** del evento (llega por los hooks), no forma parte de la clave de scoping.
- **Invariante:** un change = un worktree = una rama. La concurrencia se resuelve por worktrees, que son directorios de proyecto distintos.
- Los dos canales —hooks (mecánico) y MCP (semántico)— se scopean con la misma (proyecto, change) y se juntan en la vista de la corrida.
- El scoping nunca depende del `Mcp-Session-Id` del protocolo MCP: Claude Code no lo reenvía al servidor.
- **Owning-root de los records:** dentro de un proyecto (repo), los records EDR/spec se **sub-identifican por su `owning-root`** (la carpeta padre directa de su `.matecito-ai/`), porque un monorepo puede tener varios `.matecito-ai/`. Ver [`../process/index-decision-records.md`](../process/index-decision-records.md).
- **Scoping por rama de todo lo que refleja el repo:** todo lo que refleja contenido del repo (que varía por rama) se scopea por `(proyecto, rama)` — los eventos y changes, y los records EDR/spec (su índice, status, `active`/`deleted` y versión vigente). El **proyecto es el contenedor branch-independiente**: solo su registración y el watch son project-level. El **contenido de versiones** es content-addressable, **compartido entre ramas por hash** (no se scopea por rama). Ver [`../../../apps/api/.matecito-ai/edr/data/storage-sync-model.md`](../../../apps/api/.matecito-ai/edr/data/storage-sync-model.md).

## Escenarios

### Scenario: proyecto por header

- **GIVEN** una request MCP con el header `X-Project-Path` (y un hook con su `cwd`)
- **WHEN** el sistema scopea el evento
- **THEN** el daemon **resuelve** ese path al `project-id` estable del `.id` y scopea el evento a ese proyecto — el path es el lookup por-máquina, no la identidad

### Scenario: change por rama (M1)

- **GIVEN** un proyecto con una rama git activa
- **WHEN** el sistema scopea un evento
- **THEN** el change se deriva de la rama, resolviendo a la key (proyecto, rama)

### Scenario: sin rama (detached HEAD) → rechazo

- **GIVEN** un proyecto en detached HEAD, es decir sin rama activa
- **WHEN** una tool de escritura pide scopear un evento
- **THEN** el sistema no resuelve ningún change y rechaza la escritura con `branch_required`

### Scenario: repo git sin commits (unborn branch)

- **GIVEN** un repo git recién inicializado, sin commits todavía, cuya rama existe solo como nombre
- **WHEN** el sistema scopea un evento
- **THEN** usa ese nombre de rama como key (proyecto, rama) con normalidad — no es un caso de rechazo

### Scenario: session como atributo

- **GIVEN** dos sesiones de Claude Code sobre el mismo directorio y la misma rama
- **WHEN** cada una emite eventos
- **THEN** ambos se scopean al mismo change, y se distinguen entre sí por el session id (atributo), no por la clave

### Scenario: concurrencia por worktrees

- **GIVEN** dos ramas del mismo proyecto en worktrees distintos
- **WHEN** cada una emite eventos
- **THEN** el sistema los scopea a dos changes distintos, sin cruce

### Scenario: los canales se juntan por (proyecto, change)

- **GIVEN** eventos del canal de hooks y del canal MCP de la misma corrida
- **WHEN** el sistema los scopea
- **THEN** ambos comparten la misma (proyecto, change) y se unen en la vista de la corrida

### Scenario: el `Mcp-Session-Id` no afecta el scoping

- **GIVEN** una request MCP que incluye un header `Mcp-Session-Id`
- **WHEN** el sistema scopea el evento
- **THEN** ignora el `Mcp-Session-Id` y resuelve el change por (proyecto, rama) — el session id del protocolo MCP nunca participa de la clave

## Referencias

- **EDR** → [`../../../apps/api/.matecito-ai/edr/contracts/mcp-server.md`](../../../apps/api/.matecito-ai/edr/contracts/mcp-server.md) — el modelo de identidad por request: header de proyecto, change derivado de la rama, session como atributo, e independencia del `Mcp-Session-Id`.
- **Process** → [`../process/index-decision-records.md`](../process/index-decision-records.md) — la identidad de los records por slug + owning-root en monorepos con varios `.matecito-ai/`.
- **EDR** → [`../../../apps/api/.matecito-ai/edr/data/storage-sync-model.md`](../../../apps/api/.matecito-ai/edr/data/storage-sync-model.md) — el modelo de almacenamiento y sincronización: archivo canónico para editar, base espejo con versiones content-addressable, y el scoping por `(proyecto, rama)` de lo que refleja el repo.
