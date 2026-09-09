# EDR — La clave del Init Guard se lee de la declaración del dominio, nunca se deriva de su id

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
El Init Guard construía la clave de búsqueda por convención a partir del id del dominio. `development` declara `sdd-init/{project}`, pero la convención por id producía `development-init/{project}` — la búsqueda nunca coincidía, y el guard re-despachaba init en CADA comando de flujo. `design` coincidía por casualidad, por eso el defecto no se notó antes.

## Decisión
El Init Guard lee la clave que el dominio declara — la fila "Init topic key" de su tabla de vocabulario — sustituye `{project}`, y busca exactamente esa clave. Nunca deriva la clave del id del dominio. Todo dominio DEBE declarar esa fila; un fragmento sin ella deja al guard sin nada que leer, y eso es un defecto del fragmento, no licencia para volver a la convención.

## Reglas verificables
- **[manual]** `payload/core/CLAUDE.md`, sección `### Init Guard (MANDATORY)`, lee la clave de la fila "Init topic key" que el dominio declara, y nunca la deriva del id del dominio.

## Alternativas consideradas
Corregir sólo la fórmula de derivación por convención (por ejemplo, usar el nombre del primer agente en vez del id) — descartada: cualquier fórmula derivada vuelve a divergir la próxima vez que un dominio declare su namespace distinto a lo que la fórmula asume; leer la declaración explícita no tiene ese modo de falla.

## Consecuencias
Un dominio que no declara su fila de "Init topic key" deja al guard sin clave que buscar — comportamiento a propósito, tratado como defecto del fragmento en vez de resuelto con un default silencioso.
