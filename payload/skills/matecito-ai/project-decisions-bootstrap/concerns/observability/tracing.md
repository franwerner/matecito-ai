---
name: tracing
depth: light
domain: observability
type: decision
source: SRE · OpenTelemetry
---

# Fase: Tracing distribuido

## Qué decide

Si el proyecto instrumenta trazas distribuidas, cómo propaga el contexto entre servicios, y con qué backend. Sin propagación de contexto, los logs y métricas de múltiples servicios no se pueden correlacionar.

## Preguntas

Una o dos, según haga falta.

### 1. Nivel de tracing

> El tracing distribuido es imprescindible cuando una request cruza más de un proceso. En un monolito puede ser útil pero no crítico. Instrumentar después es costoso.

- **Sin tracing por ahora** — *default honesto para monolitos sin dependencias externas.*
- Tracing en proceso (spans internos, sin propagación cross-service) — útil para monolitos con queries lentas.
- Tracing distribuido completo con propagación de contexto — *default para microservicios y sistemas con múltiples servicios.*
- No sé, recomendame.

### 2. Implementación y backend

> OpenTelemetry es el estándar abierto; el backend es intercambiable. Elegir un backend vendor-locked dificulta migrar después. **Solo si en la 1 eligió algún nivel de tracing.**

- **OpenTelemetry SDK** (OTLP) con backend libre (Jaeger, Tempo, Honeycomb, Datadog) — *default recomendado.*
- SDK propietario del proveedor (Datadog APM, AWS X-Ray) — si ya hay un contrato con el proveedor y no se prevé migrar.
- No sé, recomendame.

## Notas de lógica (para el motor)

- Si eligió "Sin tracing por ahora", materializá con `Status: Pending` y razón ("monolito sin dependencias cross-service todavía; revisar si el sistema crece").
- Si el usuario ya eligió OTel en la Fase `metrics`, confirmá que usarán el mismo SDK para unificar.

## Tech a registrar

SDK de tracing si se elige uno concreto (ej: `opentelemetry-sdk.md`, `jaeger.md`).

## Qué materializar

ADR `tracing` materializado según el template `../../templates/adr.md`. La **Decisión** captura: nivel de tracing elegido (sin tracing / en proceso / distribuido con propagación), el protocolo de propagación de contexto (W3C TraceContext si OTel), el SDK, el backend de almacenamiento, y la regla de sampling si se definió.

**Reglas verificables** (cada una con su mecanismo al inicio):

- **[manual]** si hay tracing distribuido: el contexto se propaga entre servicios con el protocolo decidido (ej: W3C TraceContext), sin tramos donde se pierda el `trace-id`.
- **[manual]** el sampling sigue la regla decidida (ej: 100% en desarrollo, 10% en producción salvo errores, que siempre se muestrean).
- **[manual]** la instrumentación usa el SDK decidido (ej: OpenTelemetry), no SDKs mezclados que dificulten cambiar de backend.

Si se eligió "Sin tracing por ahora", el ADR va con `Status: Pending` y la razón concreta ("monolito sin dependencias cross-service todavía; revisar si el sistema crece"); en ese caso no lleva Reglas verificables.

**Relacionados:** vincular con `logging` (la correlación de logs puede reusar el `trace-id`) y con `metrics` cuando comparten el SDK OTLP.
