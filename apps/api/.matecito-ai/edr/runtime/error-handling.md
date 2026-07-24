# EDR — Manejo de errores

- **Status:** Accepted
- **Type:** policy
- **Date:** 2026-07-23

## Contexto

El broker recibe payloads del MCP y sirve a la UI; sus fallos van desde entradas inválidas hasta problemas de persistencia o de tail del filesystem. Necesita una política uniforme de propagación y traducción de errores que dé una respuesta pública estable al consumidor, y que además no exponga en logs el contenido sensible de los artefactos ni del transcript.

## Decisión

Usamos el manejo idiomático de Go: valores de error explícitos, envueltos con wrapping para preservar la cadena, y sentinels por paquete para los casos que el llamador necesita distinguir. El boundary de manejo de errores está centralizado en el borde de transporte, que traduce cualquier error interno a la respuesta pública de error con su código. El logging de errores es estructurado y nunca vuelca payloads de artefactos ni contenido de transcript.

## Reglas verificables

- **[manual]** los errores se propagan como valores con wrapping, no como panics, salvo bugs irrecuperables.
- **[manual]** el boundary que traduce al error público (con su código) vive en el borde de transporte.
- **[manual]** cada paquete expone sus propios sentinels para los errores que el llamador debe distinguir.
- **[manual]** al loggear un error nunca se vuelca el payload del artefacto ni el contenido del transcript.

## Alternativas consideradas

- **Panic/recover como flujo de control:** descartada; en Go los errores esperables son valores, no excepciones.

## Consecuencias

- La respuesta pública de error es estable e independiente de la estructura interna.
- La cadena de wrapping conserva el diagnóstico sin filtrar contenido sensible.
- Trade-off: el wrapping explícito y los sentinels por paquete son más verbosos que un mecanismo de excepciones.

## Relacionados

- `relacionado-con` → [../observability/logging.md](../observability/logging.md) — el logging estructurado de errores y la exclusión de payloads/transcript se definen ahí.
