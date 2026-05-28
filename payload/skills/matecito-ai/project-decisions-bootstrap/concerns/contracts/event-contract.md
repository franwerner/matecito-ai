---
name: event-contract
depth: light
domain: contracts
tipo: decisión
adr-output: event-contract
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

ADR `event-contract` con: formato de schema y herramienta de validación, convención de naming de tipos de evento, estrategia de versionado (versión en el tipo vs campo de versión en el payload), política de idempotencia del consumidor, y política de backward compatibility (campos nuevos opcionales, campos eliminados solo en versión major).
