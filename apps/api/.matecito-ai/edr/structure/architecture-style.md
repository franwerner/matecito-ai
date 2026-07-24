# EDR — Estilo de arquitectura

- **Status:** Accepted
- **Type:** decision
- **Date:** 2026-07-23

## Contexto

El broker es un daemon backend local chico (Matecito UI / cockpit): recibe JSON estructurado del MCP por HTTP, lo persiste en una SQLite como event-log y lo sirve a una UI React por WebSocket. Es single-dev, greenfield y local-first. Para un componente de este tamaño, una arquitectura con capas ceremoniales (application-services, puertos por todos lados) agrega fricción sin beneficio; pero los bordes de I/O sí necesitan aislarse para testear con implementaciones reales.

## Decisión

Adoptamos un estilo layered pragmático organizado por componente, con acoplamiento pragmático. Definimos interfaces SOLO en los bordes de I/O (el store de persistencia, el tail del filesystem, el hub de WebSocket); el resto de la lógica interna es acoplamiento directo, sin capas intermedias ni una capa de servicios de aplicación ceremonial. La cohesión se logra agrupando por componente, no por sufijos técnicos.

## Alcance

- `apps/api/internal/**` — paquetes del broker organizados por componente.

## Reglas verificables

- **[manual]** solo hay interfaces en los bordes de I/O declarados (store, tail, hub); la lógica interna no define interfaces de abstracción.
- **[manual]** no existe una capa de application-service ceremonial: el broker es chico y se acopla directo.

## Alternativas consideradas

- **Arquitectura hexagonal / puertos y adaptadores completa:** descartada por sobredimensionada para un daemon chico single-dev; los puertos ceremoniales agregan indirección sin pagar su costo.

## Consecuencias

- Menos indirección y boilerplate; el código se lee de corrido.
- Los bordes de I/O quedan testeables y sustituibles sin arrastrar el resto.
- Trade-off: el acoplamiento directo interno exige disciplina manual para no filtrar detalles de I/O hacia la lógica; se apoya en las reglas verificables de este EDR.

## Relacionados

- `relacionado-con` → [folder-structure.md](folder-structure.md) — el layout paquete-por-componente materializa este estilo.
