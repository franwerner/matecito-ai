# EDR — Validación de entrada

- **Status:** Accepted
- **Type:** policy
- **Date:** 2026-07-23

## Contexto

Todo lo que entra a la UI proviene del broker (envelope WS `{type, payload}` y snapshot HTTP). El broker es schema-first (Go-first vía huma), así que su contrato es la fuente de verdad; la UI no debe escribir schemas a mano ni dejar que un payload inválido se cuele al estado.

## Decisión

**Zod + tipos TS autogenerados por `Kubb`** desde el OpenAPI del broker. El borde de datos parsea todo lo entrante (envelope WS y snapshot HTTP) contra esos schemas generados **antes** de que entre al store o la cache. Un payload inválido es un error controlado: no se cuela al estado, no tumba la app y se surfacea vía el manejo de errores. Los schemas Zod generados espejan el contrato schema-first del broker. Cuando llegue la fase de escritura, la validación de forms también será con Zod.

## Alcance

- `apps/ui/src/shared/api/**` — el borde donde se parsea lo entrante contra los schemas generados.

## Reglas verificables

- **[manual]** todo dato entrante del broker pasa por su schema Zod (generado por Kubb) en el borde de datos antes de entrar al estado.
- **[manual]** un payload inválido nunca se cuela al estado ni tumba la app.
- **[tool: kubb]** los tipos y schemas se generan del OpenAPI del broker, no se escriben a mano.

## Relacionados

- `relacionado-con` → [../structure/architecture-style.md](../structure/architecture-style.md) — el borde de datos como punto único de entrada.
