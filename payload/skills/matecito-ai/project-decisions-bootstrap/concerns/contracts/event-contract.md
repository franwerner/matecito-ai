---
name: event-contract
depth: light
domain: contracts
type: decision
source: arc42 §8 (cross-cutting concepts) · CloudEvents spec v1.0
---

# Fase: Contrato de eventos

## Qué decide

Cómo se estructuran y versionan los eventos publicados, la convención de naming, e idempotencia del lado del consumidor.

## Preguntas

### 1. Schema y versionado de eventos

> Un evento sin schema versionado es un contrato implícito — los consumidores se rompen sin aviso cuando el productor cambia. CloudEvents provee un envelope estándar que simplifica routing, logging y tracing entre sistemas.

- **Schema declarativo (JSON Schema / Avro / Protobuf) + versión en el tipo del evento** — *default recomendado; ej: `order.created.v1`.*
- Envelope CloudEvents (spec estándar) con schema registrado.
- Sin schema formal — solo para sistemas internos con productor y consumidor bajo control total.
- No sé, recomendame.

### 2. Idempotencia del consumidor

> Los brokers de mensajes garantizan at-least-once en la mayoría de los casos. Sin idempotencia en el consumidor, un evento duplicado produce efectos dobles.

- **Deduplicación por `eventId` — el consumidor registra IDs procesados y descarta duplicados** — *default recomendado.*
- Operaciones naturalmente idempotentes — el handler ya es safe ante repetición por diseño.
- No aplica — garantía exactly-once provista por el broker.

## Notas de lógica (para el motor)

- Si el proyecto no publica ni consume eventos asincrónicos, marcar esta fase como `Not Applicable`.

## Tech a registrar

Si se elige un schema registry (Confluent Schema Registry, AWS Glue, etc.) o un formato de serialización (Avro, Protobuf), registrarlo en el catálogo `tech/`.

## Qué materializar

ADR `event-contract` materializado según `../../templates/adr.md`. Debe contener:

- **Contexto**: por qué un evento sin schema versionado es un contrato implícito que rompe consumidores sin aviso, y qué aporta un envelope estándar (CloudEvents) para routing, logging y tracing.
- **Decisión**: formato de schema y herramienta de validación (JSON Schema, Avro, Protobuf, CloudEvents), convención de naming de tipos de evento (ej. `order.created.v1`), estrategia de versionado (versión en el tipo vs campo de versión en el payload), política de idempotencia del consumidor (deduplicación por `eventId`, operaciones naturalmente idempotentes, o exactly-once del broker), y política de backward compatibility.
- **Reglas verificables** (cada una con su mecanismo):
  - `[tool: schema-validation]` todo evento publicado valida contra su schema registrado.
  - `[manual]` los tipos de evento siguen la convención de naming decidida, con la versión incluida.
  - `[manual]` un evento duplicado no produce efectos dobles, según la estrategia de idempotencia elegida.
  - `[manual]` los cambios solo agregan campos opcionales; los campos eliminados o renombrados solo ocurren en una versión major.
- **Alternativas consideradas**: los otros formatos de schema y estrategias de idempotencia evaluados y por qué no se eligieron.
- **Consecuencias**: dependencia de un schema registry si aplica y el contrato de compatibilidad que productores y consumidores deben respetar.
