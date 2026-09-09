# EDR — El guard de decisiones no mantiene su propia lista de secciones gateables; referencia la tabla canónica

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
La lista de secciones gateables (qué sección de qué fase dispara, y bajo qué valor de `gates:`) vivía
duplicada en el propio `### Unresolved Decisions Guard` y en el contrato canónico
(`_shared/sdd-phase-common.md`, Sección D.3). Cada edición de una tabla corría el riesgo de desalinear
la otra copia, sin que nada lo detectara.

## Decisión
El guard deja de mantener su propia lista de secciones gateables. Las secciones candidatas son
exactamente las que la tabla canónica de la Sección D.3 marca `contested`, `always` o `muted` — el guard
lee esa lista ahí y no mantiene una copia paralela.

## Reglas verificables
- **[manual]** `payload/domains/development/CLAUDE.md`, sección `### Unresolved Decisions Guard (MANDATORY)`, declara que las secciones candidatas son las que la Sección D.3 marca `contested`, `always` o `muted`, y que el guard no mantiene una copia paralela de esa lista.

## Alternativas consideradas
Mantener la lista de secciones gateables duplicada en el guard, sincronizada a mano en cada cambio —
descartada: es exactamente el defecto que motivó esta decisión, dos copias que divergen sin que nada lo
note.

## Consecuencias
Cualquier cambio al conjunto de secciones gateables se hace una sola vez, en la Sección D.3, y el guard
lo hereda por referencia.
