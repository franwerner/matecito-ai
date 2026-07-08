---
name: api-contract
depth: light
domain: contracts
type: decision
source: arc42 §8 (cross-cutting concepts) · Richardson Maturity Model
---

# Fase: Contrato de API

## Qué decide

Cómo se versiona la API, cómo se pagina, el formato de respuesta estándar, e idempotencia de operaciones de escritura.

## Preguntas

### 1. Versionado y schema

> Sin estrategia de versionado, cualquier cambio breaking rompe clientes existentes sin aviso. El enfoque schema-first o code-first afecta quién es la fuente de verdad del contrato.

- **Versión en URL (`/v1/...`) + schema-first (OpenAPI / GraphQL SDL)** — *default recomendado; el schema es la fuente de verdad.*
- Versión en header (`Accept: application/vnd.api+json;version=1`).
- Sin versionado explícito — solo para APIs internas con un único consumidor controlado.
- No sé, recomendame.

### 2. Paginación e idempotencia

> La paginación sin estándar produce implementaciones inconsistentes entre endpoints. La idempotencia en escrituras permite reintentos seguros ante fallos de red.

- **Cursor-based para listas + idempotency key en POST/PUT** — *default recomendado para APIs que escalan.*
- Offset/limit — más simple; suficiente para datasets pequeños sin crecimiento esperado.
- No aplica — la API no tiene listas ni operaciones de escritura que necesiten reintentos.

## Notas de lógica (para el motor)

- Para `api-graphql`, el versionado por URL no aplica; adaptar la pregunta 1 a evolución de schema (deprecation de campos vs tipos nuevos).
- Si el proyecto es `microservicio`, mencionar que el contrato es también interfaz entre servicios internos.

## Qué materializar

EDR `api-contract` materializado según `~/.claude/references/edr/templates/edr.md`. Debe contener:

- **Contexto**: por qué sin estrategia de versionado cualquier cambio breaking rompe clientes sin aviso, y qué consumidores tiene la API.
- **Decisión**: estrategia de versionado (URL, header, o sin versionado), enfoque schema-first vs code-first, formato de respuesta estándar (estructura de éxito y de error), mecanismo de paginación (cursor-based, offset/limit), política de idempotencia, y si hay un schema publicado (URL o ubicación en el repo).
- **Reglas verificables** (cada una con su mecanismo):
  - `[tool: schema-validation]` toda respuesta valida contra el schema publicado (estructura de éxito y de error).
  - `[manual]` las operaciones de escritura (POST/PUT) aceptan idempotency key y reintentos seguros producen el mismo efecto.
  - `[manual]` las listas se paginan con el mecanismo decidido de forma consistente entre endpoints.
  - `[manual]` los cambios breaking incrementan la versión según la estrategia elegida.
- **Alternativas consideradas**: los otros mecanismos de versionado y paginación evaluados y por qué no se eligieron.
- **Consecuencias**: quién es la fuente de verdad del contrato (schema vs código) y el costo de mantener el versionado.
