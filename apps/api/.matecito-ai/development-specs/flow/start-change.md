# Capability — Iniciar o continuar un change

- **Status:** Accepted
- **Date:** 2026-07-24

## Propósito

Dar de alta o continuar un change sobre la rama actual y registrar el mapeo rama↔nombre-de-change, de modo que una sesión nueva sobre la misma rama continúe el mismo change en lugar de crear otro.

## Actores

- **El orquestador** (thread principal de Claude Code) que llama la tool `start_change` una vez al arrancar el flujo.

## Precondiciones

- El proyecto está registrado.
- Si hay una rama git, la key del change es (proyecto, rama); si no hay rama, se requiere un `change-name` explícito como fallback.

## Flujo principal

1. El orquestador llama `start_change(change-name, [branch])`.
2. El sistema identifica el proyecto a partir del contexto de la request (header de proyecto).
3. Resuelve la identidad del change: con rama, la key es (proyecto, rama); sin rama, la key es (proyecto) y el `change-name` pasa a ser el "change actual del proyecto".
4. Busca si el change ya existe para esa key: si existe, lo continúa (no lo duplica); si no existe, lo da de alta con el nombre y registra el mapeo.
5. Devuelve el change (su identidad y su estado).

## Ramas / flujos alternativos

- **El change ya existe para la misma key** → lo continúa; no crea uno nuevo ni duplica.
- **El proyecto no está registrado** → se rechaza con `project_not_registered`.
- **No hay rama ni `change-name`** → se rechaza con `change_name_required`.
- **Conflicto de nombre** (la rama ya está mapeada a otro `change-name`) → se rechaza con `change_name_conflict` y no cambia nada.
- **El change target está `closed`** → se rechaza con `change_closed` (guard del ciclo de vida del change).

## Casos borde

- **Rama distinta = change distinto** (modelo M1): dos ramas del mismo proyecto resuelven a dos changes separados, sin cruce.
- **Continuación cross-session**: una sesión nueva sobre la misma rama resuelve a la misma key y continúa el mismo change, no crea otro.

## Reglas de negocio

- Un change se identifica por la key (proyecto, rama) cuando hay rama; por (proyecto) con el `change-name` como change actual cuando no la hay (fallback).
- El mapeo rama↔nombre-de-change se registra una sola vez, al alta; un intento posterior de mapear la misma rama a un nombre distinto es un conflicto y se rechaza.
- Invariante de identidad: un change = un worktree = una rama (modelo M1).

## Errores de cara al actor

- **Proyecto no registrado** → `project_not_registered`.
- **Rama ya mapeada a otro nombre** → `change_name_conflict`.
- **Sin rama ni change-name** → `change_name_required`.
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

### Scenario: alta sin rama (fallback)

- **GIVEN** un proyecto registrado sin rama git y un `change-name` explícito
- **WHEN** el orquestador llama `start_change(change-name)`
- **THEN** el sistema da de alta el change con la key (proyecto) y el `change-name` como change actual del proyecto

### Scenario: conflicto de nombre

- **GIVEN** una rama ya mapeada a un `change-name`
- **WHEN** el orquestador llama `start_change` con un `change-name` distinto para esa misma rama
- **THEN** el sistema responde `change_name_conflict` y no cambia nada

### Scenario: sin rama ni nombre

- **GIVEN** un proyecto registrado sin rama git y sin `change-name`
- **WHEN** el orquestador llama `start_change`
- **THEN** el sistema responde `change_name_required`

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

- **EDR** → [`../../edr/contracts/mcp-server.md`](../../edr/contracts/mcp-server.md) — la tool `start_change`, la identidad por request y el registro del mapeo rama↔nombre.
- **EDR** → [`../../edr/delivery/deployment-topology.md`](../../edr/delivery/deployment-topology.md) — el daemon global único y el invariante change↔worktree↔rama.
- **Rule** → [`../rule/event-scoping.md`](../rule/event-scoping.md) — cómo se deriva el change de la rama (modelo M1) y el fallback sin rama.
- **Lifecycle** → [`../lifecycle/change.md`](../lifecycle/change.md) — el guard `change_closed`.
