# Dominio: `observability`

Cómo se hace visible el estado del sistema en producción: logs, métricas, trazas distribuidas y health checks.

## Criterio de pertenencia

Un concern nuevo va en `observability` si trata sobre *medir o exponer* el estado del sistema corriendo. Si trata sobre proteger el sistema, va en `security`.

## Concerns en este dominio

| Concern | Prof. | Tipo | Qué decide |
|---|---|---|---|
| [health-checks](health-checks.md) | light | decisión | Qué endpoints de salud expone el servicio, qué chequea cada uno, y cómo los usa el orquestador (Kubernetes, ECS, load balancer). Sin esta distinción, un rest... |
| [logging](logging.md) | light | política | Formato de logs, niveles disponibles, correlación de requests, y librería según stack. Decisión transversal: afecta debugging, observabilidad, y compliance. |
| [metrics](metrics.md) | light | decisión | Qué se mide (qué señales), en qué formato se exponen, y dónde se almacenan. Sin métricas no hay SLOs; sin SLOs no hay alertas accionables. |
| [tracing](tracing.md) | light | decisión | Si el proyecto instrumenta trazas distribuidas, cómo propaga el contexto entre servicios, y con qué backend. Sin propagación de contexto, los logs y métricas... |
