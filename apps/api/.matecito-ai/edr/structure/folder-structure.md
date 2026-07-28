# EDR — Estructura de carpetas

- **Status:** Accepted
- **Type:** convention
- **Date:** 2026-07-23

## Contexto

El broker vive dentro del monorepo bajo `apps/api` pero es un componente autónomo con su propio ciclo de build (se compila como binario Go). Necesita un layout que lo mantenga como módulo Go propio, separe el entrypoint del daemon de su lógica interna y agrupe cada componente en su propio paquete siguiendo el idiom de Go (organización por paquete, no por sufijos de archivo).

## Decisión

El broker es un módulo Go propio: tiene su `go.mod` bajo `apps/api`. El layout sigue el convencional de Go: un entrypoint del daemon separado de la lógica interna, y la lógica interna dividida en paquetes, uno por componente (transporte, store, ingesta, config). La organización es paquete-por-componente; las interfaces se declaran en los bordes de I/O de cada paquete.

## Alcance

- `apps/api/cmd/broker/**` — entrypoint del daemon.
- `apps/api/internal/transport/**` — borde HTTP-in + WS-out.
- `apps/api/internal/store/**` — acceso a SQLite.
- `apps/api/internal/ingest/**` — hooks + tail del transcript.
- `apps/api/internal/config/**` — carga y validación de config.

## Reglas verificables

- **[tool: go-arch-lint]** cada componente vive en su propio paquete bajo `internal/` y está declarado en el grafo; un paquete que no pertenece a ningún componente declarado rompe el gate.
- **[manual]** el entrypoint del daemon vive en `cmd/broker`.
- **[tool: go-arch-lint]** el entrypoint es el único componente que puede depender de los tres componentes de primer nivel a la vez; ningún componente de primer nivel depende de otro salvo hacia el borde de persistencia.

## Alternativas consideradas

- **Organización por capa técnica (handlers/, services/, models/):** descartada por no ser idiomática en Go y por dispersar un mismo componente en varias carpetas.

## Consecuencias

- Cada componente se encuentra en un solo lugar; agregar uno nuevo es agregar un paquete.
- El módulo propio permite compilar y versionar el broker por separado.
- Trade-off: mantener `go.mod` propio dentro del monorepo agrega un paso de coordinación (go.work) para el binario que lo embebe.

## Relacionados

- `relacionado-con` → [architecture-style.md](architecture-style.md) — este layout materializa el estilo layered pragmático por componente.
