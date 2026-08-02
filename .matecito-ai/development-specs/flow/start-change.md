# Capability — Iniciar o continuar un change

- **Status:** Accepted
- **Date:** 2026-07-27
- **Components:** api

## Propósito

Dar de alta o continuar un change sobre la rama actual y registrar el mapeo rama↔nombre-de-change, de modo que una sesión nueva sobre la misma rama continúe el mismo change en lugar de crear otro.

## Actores

- **El orquestador** (thread principal de Claude Code) que llama la tool `start_change` una vez al arrancar el flujo.

## Precondiciones

- El proyecto está registrado.
- El proyecto está sobre una rama git: la key del change es siempre (proyecto, rama).

## Flujo principal

1. El orquestador llama `start_change(change-name, [branch])`.
2. El sistema identifica el proyecto a partir del contexto de la request (header de proyecto).
3. Resuelve la identidad del change: la key es (proyecto, rama).
4. Busca si el change ya existe para esa key: si existe, lo continúa (no lo duplica); si no existe, lo da de alta con el nombre y registra el mapeo.
5. Devuelve el change (su identidad y su estado).

## Ramas / flujos alternativos

- **El change ya existe para la misma key** → lo continúa; no crea uno nuevo ni duplica.
- **El proyecto no está registrado** → se rechaza con `project_not_registered`.
- **El proyecto no está sobre una rama** (detached HEAD) → se rechaza con `branch_required`.
- **Conflicto de nombre** (la rama ya está mapeada a otro `change-name`) → se rechaza con `change_name_conflict` y no cambia nada.
- **El change target está `closed`** → se rechaza con `change_closed` (guard del ciclo de vida del change).

## Casos borde

- **Rama distinta = change distinto** (modelo M1): dos ramas del mismo proyecto resuelven a dos changes separados, sin cruce.
- **Continuación cross-session**: una sesión nueva sobre la misma rama resuelve a la misma key y continúa el mismo change, no crea otro.

## Reglas de negocio

- Un change se identifica siempre por la key (proyecto, rama); sin rama no hay change posible.
- El mapeo rama↔nombre-de-change se registra una sola vez, al alta; un intento posterior de mapear la misma rama a un nombre distinto es un conflicto y se rechaza.
- Invariante de identidad: un change = un worktree = una rama (modelo M1).

## Errores de cara al actor

- **Proyecto no registrado** → `project_not_registered`.
- **Rama ya mapeada a otro nombre** → `change_name_conflict`.
- **Proyecto sin rama (detached HEAD)** → `branch_required`.
- **Change cerrado** → `change_closed`.
- La forma de error es `{error, code}`, estricta, sin internals.

## Escenarios

### Scenario: alta con rama

- **GIVEN** un proyecto registrado y una rama git sin change previo
- **WHEN** el orquestador llama `start_change(change-name, branch)`
- **THEN** el sistema da de alta el change con la key (proyecto, rama), registra el mapeo rama↔nombre y devuelve su identidad y estado

### Scenario: continúa cross-session

- **GIVEN** un change ya existente para la key (proyecto, rama)
- **WHEN** una sesión nueva sobre la misma rama llama `start_change`
- **THEN** el sistema continúa el mismo change y no crea ni duplica otro

### Scenario: conflicto de nombre

- **GIVEN** una rama ya mapeada a un `change-name`
- **WHEN** el orquestador llama `start_change` con un `change-name` distinto para esa misma rama
- **THEN** el sistema responde `change_name_conflict` y no cambia nada

### Scenario: proyecto sin rama (detached HEAD)

- **GIVEN** un proyecto registrado en detached HEAD, es decir sin rama activa
- **WHEN** el orquestador llama `start_change(change-name)`
- **THEN** el sistema responde `branch_required` y no da de alta ningún change

### Scenario: proyecto no registrado

- **GIVEN** un contexto cuyo proyecto no fue registrado previamente
- **WHEN** el orquestador llama `start_change`
- **THEN** el sistema responde `project_not_registered`

### Scenario: rama distinta = change distinto (M1)

- **GIVEN** un change existente para (proyecto, rama A)
- **WHEN** el orquestador llama `start_change` sobre la rama B del mismo proyecto
- **THEN** el sistema resuelve a un change distinto, sin cruzarse con el de la rama A

### Scenario: sobre change cerrado

- **GIVEN** un change en estado `closed`
- **WHEN** el orquestador llama `start_change` sobre esa key
- **THEN** el sistema responde `change_closed`

## Referencias

- **EDR** → [`../../../apps/api/.matecito-ai/edr/contracts/mcp-server.md`](../../../apps/api/.matecito-ai/edr/contracts/mcp-server.md) — la tool `start_change`, la identidad por request y el registro del mapeo rama↔nombre.
- **EDR** → [`../../../apps/api/.matecito-ai/edr/delivery/deployment-topology.md`](../../../apps/api/.matecito-ai/edr/delivery/deployment-topology.md) — el daemon global único y el invariante change↔worktree↔rama.
- **Rule** → [`../rule/event-scoping.md`](../rule/event-scoping.md) — cómo se deriva el change de la rama (modelo M1) y el rechazo `branch_required` sin rama.
- **Lifecycle** → [`../lifecycle/change.md`](../lifecycle/change.md) — el guard `change_closed`.
