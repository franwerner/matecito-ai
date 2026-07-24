# EDR — Contrato de API

- **Status:** Accepted
- **Type:** decision
- **Date:** 2026-07-23

## Contexto

La "API" del broker no es una API pública tradicional: tiene dos superficies distintas para consumidores controlados y locales, y conviene no confundirlas. Una es la **superficie de escritura cara a Claude**: Claude Code emite las fases del flujo contra el MCP, que no es un poster HTTP externo sino código in-process en el mismo binario que el store; su contrato se expresa con schemas de tool MCP (los schemas de artefacto por fase), no con REST/OpenAPI, y su detalle vive en su propio EDR. La otra es la **superficie de lectura cara a la UI**: un snapshot HTTP (REST) que la UI consulta más el WS-out por el que recibe eventos en vivo; esta es la que se describe como OpenAPI Go-first y de la que la UI genera sus tipos y schemas. Como es un daemon local con consumidores conocidos, no hace falta negociación de versión en runtime, pero sí una fuente de verdad versionada del contrato. Además, en el modelo de instancia única global el MCP apunta a un broker global compartido, así que hace falta una superficie para registrar cada proyecto por su ruta.

## Decisión

Definimos el contrato sobre dos superficies:

1. **Versionado:** schema-first, sin versión en la URL. El schema del contrato JSON es la fuente de verdad, vive en el repo y se versiona con git; no hay negociación de versión en runtime porque los consumidores son controlados y locales.
2. **Superficie de escritura (cara a Claude):** el MCP recibe las fases del flujo como salida estructurada y las escribe al store **in-process** — no es un poster REST hacia el broker, es código en el mismo binario. Su contrato se describe con **schemas de tool MCP** (los schemas de artefacto por fase), no con OpenAPI/REST; en error propaga el error público con su código (reusa la política de manejo de errores). El diseño de esta superficie —proceso, transporte, identidad por request y tools— vive en su propio EDR.
3. **Superficie de lectura (cara a la UI):** un snapshot HTTP (REST) que la UI consulta, implementado con `net/http` (stdlib, ServeMux) más `huma`, que produce el **OpenAPI 3.1** de forma Go-first. **Este** OpenAPI describe la lectura de la UI (no la escritura del MCP), y es de lo que la UI genera sus tipos y schemas vía Kubb.
4. **WS-out (broker → UI):** parte de la superficie de lectura; la UI se suscribe y recibe eventos con un envelope `{type, payload}`.
5. **Schemas compartidos:** los schemas de artefacto y de evento son structs Go compartidos (una sola fuente), expuestos de dos formas — como tool-schema MCP en la escritura y como componente OpenAPI en la lectura. No se duplican por superficie.
6. **Idempotencia:** las escrituras se identifican por la identidad del evento; el append es idempotente, de modo que un reintento produce el mismo efecto.
7. **Paginación del read/history:** Pending (será cursor-based cuando se construya); la lectura inicial es push por WebSocket.
8. **Registro de proyecto (modelo global):** el MCP apunta a un broker único global; se suma una superficie para registrar un proyecto por su ruta, ya sea vía MCP (por ejemplo, auto-registro por el directorio de la sesión) o vía UI. La identidad de un change es la combinación del proyecto y el nombre de change/rama.

## Reglas verificables

- **[manual]** el schema del contrato JSON es la fuente de verdad y se versiona en git; no hay versión en la URL.
- **[manual]** la escritura del MCP es in-process (no un POST REST al broker) y se describe con schemas de tool MCP, no con OpenAPI.
- **[manual]** la superficie de lectura de la UI (snapshot REST) se implementa con net/http + huma y produce el spec OpenAPI que la UI consume vía Kubb.
- **[manual]** los schemas de artefacto/evento son structs Go compartidos, expuestos como tool-schema MCP (escritura) y como componente OpenAPI (lectura), sin duplicar.
- **[manual]** los eventos WS-out usan el envelope `{type, payload}`.
- **[manual]** las escrituras son idempotentes por identidad del evento.
- **[manual]** existe una superficie para registrar un proyecto por su ruta (vía MCP o UI) y un change se identifica por la combinación de proyecto y nombre de change/rama.

## Alternativas consideradas

- **Versión en la URL / negociación de versión en runtime:** descartada; los consumidores son locales y controlados, el schema versionado en git alcanza.
- **Un solo canal (todo HTTP con polling):** descartado; el WS-out con push es lo que da la actualización en vivo del cockpit.

## Consecuencias

- El contrato es una sola fuente de verdad versionada, sin lógica de negociación.
- La idempotencia por identidad del evento hace seguros los reintentos.
- La superficie de registro habilita el modelo multi-proyecto sobre una instancia global.
- Trade-off: al no versionar en runtime, un cambio incompatible del schema exige coordinar el redeploy de los consumidores locales.

## Relacionados

- `relacionado-con` → [mcp-server.md](mcp-server.md) — el detalle de la superficie de escritura (proceso in-process, transporte, identidad por request y tools de submit por fase).
- `relacionado-con` → [../data/data-modeling.md](../data/data-modeling.md) — el shape del envelope y la clave de idempotencia se fijan ahí.
- `relacionado-con` → [../delivery/deployment-topology.md](../delivery/deployment-topology.md) — el modelo de instancia única global condiciona el registro de proyecto.

> Nota (cross-store): la UI consume este contrato como OpenAPI generándose tipos y schemas con Kubb. La decisión del lado consumidor vive en el store de EDRs de la UI (`apps/ui/.matecito-ai/edr/security/input-validation.md`).
