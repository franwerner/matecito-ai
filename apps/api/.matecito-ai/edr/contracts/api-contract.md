# EDR — Contrato de API

- **Status:** Accepted
- **Type:** decision
- **Date:** 2026-07-23

## Contexto

La "API" del broker no es una API pública tradicional: tiene dos superficies para consumidores controlados y locales. Por un lado el MCP postea JSON estructurado (HTTP-in); por otro la UI se suscribe y recibe eventos (WS-out). Como es un daemon local con consumidores conocidos, no hace falta negociación de versión en runtime, pero sí una fuente de verdad versionada del contrato. Además, en el modelo de instancia única global el MCP apunta a un broker global compartido, así que hace falta una superficie para registrar cada proyecto por su ruta.

## Decisión

Definimos el contrato sobre dos superficies:

1. **Versionado:** schema-first, sin versión en la URL. El schema del contrato JSON es la fuente de verdad, vive en el repo y se versiona con git; no hay negociación de versión en runtime porque los consumidores son controlados y locales.
2. **HTTP-in (MCP → broker):** el MCP postea el JSON estructurado; en éxito responde un ack/eco, en error responde el error público con su código (reusa la política de manejo de errores). Se implementa con `net/http` (stdlib, ServeMux) más `huma`, que produce el **OpenAPI 3.1** de forma Go-first: el contrato queda expresado como OpenAPI, que la UI consume para generar sus tipos y schemas vía Kubb.
3. **WS-out (broker → UI):** la UI se suscribe y recibe eventos con un envelope `{type, payload}`.
4. **Idempotencia:** las escrituras se identifican por la identidad del evento; el append es idempotente, de modo que un reintento produce el mismo efecto.
5. **Paginación del read/history:** Pending (será cursor-based cuando se construya); la lectura inicial es push por WebSocket.
6. **Registro de proyecto (modelo global):** el MCP apunta a un broker único global; se suma una superficie para registrar un proyecto por su ruta, ya sea vía MCP (por ejemplo, auto-registro por el directorio de la sesión) o vía UI. La identidad de un change es la combinación del proyecto y el nombre de change/rama.

## Reglas verificables

- **[manual]** el schema del contrato JSON es la fuente de verdad y se versiona en git; no hay versión en la URL.
- **[manual]** HTTP-in responde un ack en éxito y el error público con su código en caso de error.
- **[manual]** el HTTP-in se implementa con net/http + huma y produce el spec OpenAPI que la UI consume.
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

- `relacionado-con` → [../data/data-modeling.md](../data/data-modeling.md) — el shape del envelope y la clave de idempotencia se fijan ahí.
- `relacionado-con` → [../delivery/deployment-topology.md](../delivery/deployment-topology.md) — el modelo de instancia única global condiciona el registro de proyecto.

> Nota (cross-store): la UI consume este contrato como OpenAPI generándose tipos y schemas con Kubb. La decisión del lado consumidor vive en el store de EDRs de la UI (`apps/ui/.matecito-ai/edr/security/input-validation.md`).
