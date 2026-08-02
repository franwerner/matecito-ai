# Capability — Registrar (o linkear) un proyecto

- **Status:** Accepted
- **Date:** 2026-07-27
- **Components:** api

## Propósito

Dar de alta —o linkear— un proyecto en el store por su identidad estable, para que sus changes, eventos y records queden asociados a él. La operación es idempotente: registrar dos veces el mismo proyecto no lo duplica.

## Actores

- **El orquestador / Claude Code** — llama `register_project(path)`; la ruta llega en el header `X-Project-Path` (que transporta el valor `${CLAUDE_PROJECT_DIR}`).
- **Un usuario vía la UI** — registra un proyecto por una ruta explícita.

## Precondiciones

- El daemon está corriendo.
- La ruta resuelve a un directorio existente.
- La ruta está dentro de un repo git (matecito requiere git).

## Flujo principal

1. El actor llama `register_project(path)`.
2. El sistema **normaliza la ruta**: la vuelve absoluta, resuelve symlinks y limpia los segmentos `.`/`..`. Luego resuelve la **raíz del repo git (su toplevel)** como raíz del proyecto; si la ruta no resuelve a un toplevel git, se rechaza con `not_a_git_repo`.
3. **Resuelve la identidad** del proyecto vía `<toplevel>/.matecito-ai/.id`, el archivo commiteado que porta el `project-id` estable:
   - **Hay `.id`** → ese es el `project-id` (un clon fresco ya lo trae, así que auto-linkea sin lookup).
   - **No hay `.id`** → no se genera un id a ciegas: con `find_project(name)` se busca en la base si ya existe un proyecto con ese nombre (por ejemplo, uno compartido o deployado). Si el usuario **confirma un match** → se **materializa `.id`** con ese `project-id` (LINK). Si **no hay match** o el proyecto es nuevo → se **genera un `project-id` nuevo** y se materializa `.id` (CREATE).
4. Si es un alta (primera vez) → el daemon **suma el proyecto a su watch**, **descubre recursivamente todos los `.matecito-ai/`** bajo el toplevel (cada record se sub-identifica por su owning-root, la ruta relativa de su store) y hace el **build inicial del índice**.
5. **Devuelve el proyecto**: su identidad estable, el `project-id`.

## Ramas / flujos alternativos

- **`.id` presente / proyecto ya registrado** → operación **idempotente**: devuelve el mismo `project-id`, no duplica.
- **Sin `.id` y con match confirmado** → materializa `.id` con el `project-id` existente y linkea el checkout (LINK).
- **Sin `.id` y sin match** → genera un `project-id` nuevo, materializa `.id` (CREATE).
- **Ruta inexistente o que no es un directorio** → se rechaza con `invalid_path`.
- **Ruta que no resuelve a un toplevel git** → se rechaza con `not_a_git_repo`.

## Casos borde

- **Ruta relativa, symlink u ortografía distinta** → la normalización canónica las colapsa a la misma raíz → misma identidad, un solo proyecto.
- **Registración anidada** (por ejemplo registrar `apps/api` con la raíz ya registrada) → el toplevel de git resuelve al **mismo proyecto**; no crea uno nuevo. El `.matecito-ai/` de `apps/api` es un owning-root dentro del proyecto, no un proyecto aparte.
- **Ruta válida sin `.matecito-ai/` todavía** → se registra igual; el watch queda armado e indexa cuando el `.matecito-ai/` aparezca.
- **Cambio de ruta** (mover o renombrar el repo) → el `.id` commiteado sobrevive al movimiento → la ruta nueva se re-mapea al mismo `project-id`, mismo proyecto.

## Reglas de negocio

- **La identidad del proyecto es el `project-id` del `.id`**, no el path local. El `.id` está commiteado, así que viaja con el repo: es machine-independent, sobrevive a mover o renombrar el checkout y sirve deployado o compartido. El **path local es solo un lookup por-máquina**, no la identidad.
- La base **keyea por `project-id`**, nunca por el path local.
- Un `project-id` no se genera nunca a ciegas cuando falta el `.id`: primero se busca un match por nombre (`find_project`) para poder linkear; solo si no hay match se crea uno nuevo.
- El descubrimiento de records es recursivo: un proyecto (repo) puede tener varios `.matecito-ai/`, y cada record se sub-identifica por su owning-root.
- Solo la registración del proyecto y el watch son a nivel proyecto (branch-independiente); el contenido que varía por rama se scopea aparte.

## Errores de cara al actor

- **Ruta inexistente o que no es un directorio** → `invalid_path`.
- **Ruta que no está en un repo git** → `not_a_git_repo`.
- La forma de error es `{error, code}`, estricta, sin internals. En caso de éxito, la operación devuelve el proyecto.

## Escenarios

### Scenario: alta nueva (create)

- **GIVEN** una ruta válida sin `.id` y sin match por nombre en la base
- **WHEN** el actor llama `register_project(path)`
- **THEN** el sistema genera un `project-id` nuevo, materializa `.id`, arranca el watch, descubre recursivamente los `.matecito-ai/`, hace el build inicial del índice y devuelve el proyecto

### Scenario: link a proyecto existente

- **GIVEN** una ruta sin `.id` cuyo nombre matchea un proyecto existente en la base
- **WHEN** el actor llama `register_project(path)` y confirma el match de `find_project(name)`
- **THEN** el sistema materializa `.id` con el `project-id` existente y linkea el checkout a ese proyecto

### Scenario: auto-link con `.id` presente

- **GIVEN** un clon con `.id` commiteado
- **WHEN** el actor llama `register_project(path)`
- **THEN** el sistema se registra directo a ese `project-id`, sin lookup por nombre

### Scenario: idempotente

- **GIVEN** un proyecto ya registrado
- **WHEN** el actor vuelve a llamar `register_project` sobre él
- **THEN** el sistema devuelve el mismo `project-id` y no duplica

### Scenario: ruta inexistente

- **GIVEN** una ruta que no existe o no es un directorio
- **WHEN** el actor llama `register_project(path)`
- **THEN** el sistema responde `invalid_path`

### Scenario: normalización canónica

- **GIVEN** dos spellings de la misma ruta (relativa, con symlink o con `.`/`..`)
- **WHEN** el actor registra por cada una
- **THEN** ambas resuelven a la misma raíz y a un solo proyecto

### Scenario: registración anidada = mismo proyecto

- **GIVEN** un proyecto ya registrado por su raíz de git
- **WHEN** el actor llama `register_project` sobre `apps/api` dentro de ese repo
- **THEN** el sistema resuelve al mismo toplevel y al mismo `project-id`; el `.matecito-ai/` de `apps/api` queda como un owning-root adentro, no como un proyecto nuevo

### Scenario: sin `.matecito-ai/` todavía

- **GIVEN** una ruta válida sin `.matecito-ai/`
- **WHEN** el actor llama `register_project(path)`
- **THEN** el sistema registra el proyecto y arma el watch, que indexa cuando el `.matecito-ai/` aparezca

### Scenario: cambio de ruta sobrevive

- **GIVEN** un proyecto registrado cuyo repo se mueve o renombra
- **WHEN** el actor registra la ruta nueva
- **THEN** el `.id` commiteado mantiene el `project-id` y la ruta nueva se re-mapea al mismo proyecto

### Scenario: alta en carpeta sin repo git → rechazo

- **GIVEN** una ruta válida que no pertenece a ningún repo git
- **WHEN** el actor llama `register_project(path)`
- **THEN** el sistema responde `not_a_git_repo` y no registra ningún proyecto ni materializa `.id`

### Scenario: descubrimiento recursivo multi-store

- **GIVEN** un repo con varios `.matecito-ai/` (la raíz, `apps/api`, `apps/ui`)
- **WHEN** el actor registra la raíz del repo
- **THEN** el sistema descubre y vigila todos los `.matecito-ai/`, cada uno sub-identificado por su owning-root

## Referencias

- **EDR** → [`../../../apps/api/.matecito-ai/edr/contracts/mcp-server.md`](../../../apps/api/.matecito-ai/edr/contracts/mcp-server.md) — las tools `register_project` y `find_project`, y que la identidad del proyecto es el `project-id` del `.id`, no el path.
- **EDR** → [`../../../apps/api/.matecito-ai/edr/data/storage-sync-model.md`](../../../apps/api/.matecito-ai/edr/data/storage-sync-model.md) — el `.id` commiteado como identidad portable y la base que keyea por `project-id`.
- **EDR** → [`../../../apps/api/.matecito-ai/edr/delivery/deployment-topology.md`](../../../apps/api/.matecito-ai/edr/delivery/deployment-topology.md) — el daemon global único que aloja el watch y el store.
- **Rule** → [`../rule/event-scoping.md`](../rule/event-scoping.md) — el proyecto identificado por `project-id` y el scoping branch-independiente de su registración.
- **Process** → [`../process/index-decision-records.md`](../process/index-decision-records.md) — el watch y el índice que arranca la registración, y el owning-root de los records.
