# 07 — Herramientas

[← 06 Auto-mine de specs](06-auto-mine-spec.md) · [Índice](README.md) · Siguiente: [08 — Configuración →](08-configuracion.md)

Las herramientas (MCP y CLI) se enganchan en las fases donde aportan. El CLI `matecito-ai` se ocupa de instalarlas y mantener el entorno sano; cada una se usa con su propio binario/servidor.

> **Nada es global — todo lo declara el dominio.** No hay herramientas "base" que se instalen siempre: cada dominio lista en su `manifest.json` los `mcp` y `binaries` que usa, y solo eso se instala (la unión de los dominios activos). development declara **engram, context7, codegraph, drawio** como `mcp` y **engram, codegraph, proofshot** como `binaries`. **Engram** es el mecanismo de memoria del núcleo, pero —como todo— su instalación la declara el dominio. Otro dominio trae lo suyo (p. ej. **Figma** en design).

## Resumen

| Herramienta | Origen | Tipo | Rol | Engancha en |
|---|---|---|---|---|
| **Engram** | núcleo (memoria), declarado por dominio | binario standalone | memoria persistente: artefactos del SDD entre fases + descubrimientos/fixes entre sesiones | todas las fases (artifacts) |
| **context7** | development | MCP | documentación de librerías al día (contra APIs alucinadas) | `apply` |
| **drawio** | development | MCP | diagramas de arquitectura on-demand y **efímeros** (preview en vivo, sin archivo) | thread principal (en el paso de design) |
| **codegraph** | development | MCP | grafo de código pre-indexado (tree-sitter + SQLite) para explorar por estructura | `explore`, `apply`, `mine` |
| **proofshot** | development | CLI | verificación visual de UI (graba el browser, valida escenarios) | `verify` (si el cambio toca UI) |

## Engram — la memoria

Es el medio por el que **pasa la información entre fases**: cada fase guarda su artefacto bajo `sdd/<change>/<artefacto>` y la siguiente lo lee (ver [02](02-flujo-sdd.md#cómo-se-pasa-la-información-entre-fases)). Persiste entre sesiones, así que el contexto no se pierde al cerrar.
**No guarda EDRs** — esos viven solo en `.matecito-ai/edr/*.md`.

## codegraph — exploración por estructura

Cuando existe `.codegraph/`, las fases que necesitan entender el código (explore, apply, y mine al analizar) consultan el grafo en vez de escanear archivo por archivo. Para texto literal o archivos no indexados, se usa grep.

## context7 — docs al día

En `apply`, cuando se trabaja con una librería/framework, context7 trae la documentación actual para no codear contra APIs inventadas.

## drawio — diagramas on-demand

Los diagramas se generan **on-demand, nunca automáticamente**, y son **efímeros**: `intake` decide si el cambio amerita uno (complejidad estructural: varios componentes, flujo cruzando límites, etc.) y lo marca en el brief; se confirma en el INTAKE GATE. Cuando aplica, el **thread principal** lo renderiza en vivo en `localhost:6002` — **no se exporta ningún archivo** al repo. El sub-agente `design` no genera diagramas (es headless, no previsualiza); solo señala que conviene uno.

## proofshot — verificación visual de UI

`intake` decide si el cambio amerita `ui-test` según los escenarios. Cuando aplica y proofshot está disponible, `verify` conduce el browser según los escenarios del spec, valida el estado en vivo y chequea errores de consola/servidor. **Capability-gated**: si proofshot o el dev-server no están, se saltea en silencio.

## Referencias consultables (no son herramientas)

Además de las herramientas hay **references** — material pasivo que las skills consultan, en `~/.claude/references/`:
- **`edr/`** — el concepto de EDR + las plantillas de estructura. Ver [04](04-decisiones-edr.md).
- **`design-patterns/`** — catálogo canónico de patrones. Cuando un EDR declara `Patrón aplicado: X`, `design` respeta su definición.
