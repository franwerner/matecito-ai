# EDR — El puntero del kernel a la excepción de fan-out de sdd-verify no ofrece el patrón a ninguna otra fase

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
`### Sub-Agent Launch Pattern` menciona que `sdd-verify` es la excepción al despacho de un-agente-por-fase. Sin nombrarlo explícitamente, esa mención podía leerse como la apertura de un patrón general — "despachar como más de un agente" — disponible para cualquier fase futura, en vez de una excepción única y nombrada cuya partición vive en el fragmento de dominio.

## Decisión
La línea que menciona la excepción de `sdd-verify` en el kernel declara explícitamente que es un puntero único, no una regla generalizada: `sdd-verify` es LA excepción nombrada al patrón de un-agente-por-fase en este pipeline, su partición vive en el fragmento de dominio de development, y esta línea no ofrece el patrón a ninguna otra fase.

## Reglas verificables
- **[manual]** `payload/core/CLAUDE.md`, sección `### Sub-Agent Launch Pattern`, declara que la mención a `sdd-verify` es un puntero único y no un patrón general ofrecido a otra fase.

## Alternativas consideradas
Dejar la mención sin la cláusula explícita de exclusividad, confiando en que el fragmento de dominio ya la tiene — descartada: el kernel es el archivo que siempre está en contexto, y una mención ambigua ahí es la que un agente lee primero.

## Consecuencias
La línea del kernel queda un poco más larga, pero cierra la lectura de que el fan-out es un patrón disponible para cualquier fase que quiera adoptarlo.

## Relacionados
- `relacionado-con` → [phase-fanout-two-cases.md](phase-fanout-two-cases.md) — el mismo principio de exclusividad, aplicado en el fragmento de dominio a los dos casos declarados de fan-out.
