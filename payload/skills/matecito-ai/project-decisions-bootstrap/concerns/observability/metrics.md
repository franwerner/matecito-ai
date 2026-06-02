---
name: metrics
depth: light
domain: observability
type: decision
source: SRE (RED method · USE method)
---

# Fase: Métricas

## Qué decide

Qué se mide (qué señales), en qué formato se exponen, y dónde se almacenan. Sin métricas no hay SLOs; sin SLOs no hay alertas accionables.

## Preguntas

Una o dos, según haga falta.

### 1. Qué se mide y cómo se expone

> **RED** (Rate, Errors, Duration) aplica a servicios orientados a requests. **USE** (Utilization, Saturation, Errors) aplica a recursos (CPU, memoria, conexiones de DB). Elegir el modelo antes de instrumentar evita métricas ad-hoc difíciles de mantener.

- **RED para endpoints + USE para recursos de infra** — *default recomendado para APIs y microservicios.*
- Solo RED — suficiente si la infra está administrada (cloud-managed DB, serverless).
- Solo métricas de negocio (pedidos por minuto, conversiones) — válido si la infra la monitorea otro equipo.
- Ninguna por ahora — *solo si el proyecto es interno o sin SLA.*
- No sé, recomendame.

### 2. Formato y destino

> Prometheus/OTel son los estándares de facto; la elección define qué herramienta de visualización y alertas se puede usar. **Solo si en la 1 eligió medir algo.**

- **Prometheus** (scraping pull, `/metrics` endpoint) + Grafana — *default para proyectos auto-hospedados o Kubernetes.*
- OpenTelemetry (OTLP push) — vendor-neutral; compatible con Datadog, Honeycomb, Grafana Cloud, etc.
- Métricas cloud-native (CloudWatch, Google Cloud Monitoring, Azure Monitor) — si el proyecto vive en un solo proveedor cloud.
- No sé, recomendame.

## Notas de lógica (para el motor)

- Si eligió "Ninguna por ahora", no hagas la pregunta 2. Materializá con `Status: Pending` y razón.
- Si el usuario ya eligió OpenTelemetry en la Fase `tracing`, sugerí OTLP como default en la pregunta 2 para unificar el stack de observabilidad.

## Tech a registrar

Librería de instrumentación si se elige una concreta (ej: `prometheus-client.md`, `opentelemetry-sdk.md`).

## Qué materializar

ADR `metrics` materializado según el template `../../templates/adr.md`. La **Decisión** captura: modelo elegido (RED / USE / negocio / ninguno), formato de exposición (Prometheus pull / OTLP push / cloud-native), destino de almacenamiento, y las métricas concretas iniciales si se definieron (ej: `http_requests_total`, `http_request_duration_seconds`, `db_pool_connections_active`).

**Reglas verificables** (cada una con su mecanismo al inicio):

- **[manual]** si se eligió RED: cada endpoint expone Rate, Errors y Duration; no hay endpoints instrumentados ad-hoc fuera del modelo.
- **[manual]** las métricas se exponen en el formato decidido (ej: endpoint `/metrics` Prometheus u OTLP), no en formatos mezclados.
- **[manual]** las métricas iniciales enumeradas en la Decisión existen y se emiten con los nombres acordados.

Si se eligió "Ninguna por ahora", el ADR va con `Status: Pending` y la razón concreta; en ese caso no lleva Reglas verificables.

**Relacionados:** vincular con `tracing` cuando ambos comparten el SDK de OpenTelemetry/OTLP para unificar el stack de observabilidad.
