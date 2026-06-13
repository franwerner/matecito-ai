---
name: logging
depth: light
domain: observability
type: policy
source: 12-factor (XI: logs) · production-readiness
---

# Fase: Logging

## Qué decide

Formato de logs, niveles disponibles, correlación de requests, y librería según stack. Decisión transversal: afecta debugging, observabilidad, y compliance.

## Preguntas

Una o dos, según haga falta.

### 1. Formato y niveles

> Los logs estructurados en JSON son parseables por cualquier plataforma (Datadog, CloudWatch, Loki) sin configuración extra. Texto plano es más legible en desarrollo pero difícil de filtrar en producción.

- **Estructurado JSON en producción, texto en desarrollo** — *default recomendado para cualquier proyecto productivo.*
- Estructurado JSON siempre — si el stack es microservicios o hay un agregador centralizado.
- Texto plano siempre — solo para scripts o herramientas internas sin agregador.
- No sé, recomendame.

(Niveles a usar — confirmar o ajustar:)

- **`debug` / `info` / `warn` / `error`** — *default estándar en la mayoría de ecosistemas.*
- `trace` / `debug` / `info` / `warn` / `error` / `fatal` — si el stack y librería los soportan nativamente.

### 2. Correlación y librería

> Sin un `request-id` o `trace-id` propagado en cada log, correlacionar qué pasó durante una request en múltiples logs es imposible a escala.

- **`request-id` generado en el borde y propagado en todos los logs de esa request** — *default recomendado para cualquier API.*
- `trace-id` de OpenTelemetry si hay tracing distribuido — integra con la Fase de tracing.
- Sin correlación por ahora — solo si el proyecto no recibe tráfico concurrente.

(Librería sugerida según stack — confirmar:)

- Python → `structlog` (estructurado, contexto por request) o `logging` estándar con formatter JSON.
- Node.js / TypeScript → `pino` (performance) o `winston` (más configurable).
- Go → `zerolog` (performance) o `slog` (estándar desde Go 1.21).
- Java → `logback` + `logstash-logback-encoder` para JSON.
- Rust → `tracing` (integra con OpenTelemetry).
- No sé, recomendame.

## Notas de lógica (para el motor)

- La selección de librería puede ir en el mismo turno que la pregunta 2 si el stack ya está detectado — proponer el default y pedir confirmación.
- Si el usuario elige correlación por `trace-id` de OTel, registrar dependencia con la Fase `tracing`.

## Tech a registrar

Librería de logging elegida (ej: `structlog.md`, `pino.md`, `zerolog.md`, `winston.md`).

## Qué materializar

ADR `logging` materializado según el template `~/.claude/references/adr/templates/adr.md`. La **Decisión** captura: formato elegido (JSON estructurado / texto / mixto por entorno), niveles disponibles, la política de correlación (`request-id` o `trace-id`, dónde se genera y cómo se propaga) y la librería elegida (registrada también como tech).

**Reglas verificables** (cada una con su mecanismo al inicio):

- **[manual]** nunca se loggean passwords, tokens ni PII en ningún nivel.
- **[manual]** todo log de error incluye el stack trace completo.
- **[manual]** el nivel mínimo en producción es `info` (sin `debug` en prod).
- **[manual]** en producción los logs salen en el formato decidido (ej: JSON estructurado), no texto plano.
- **[manual]** todo log emitido durante una request lleva el identificador de correlación (`request-id`/`trace-id`) propagado desde el borde.

**Relacionados:** vincular con `tracing` si la correlación usa el `trace-id` de OpenTelemetry, y con `error-handling` (la política de qué/cómo se loggean los errores se decide allí).
