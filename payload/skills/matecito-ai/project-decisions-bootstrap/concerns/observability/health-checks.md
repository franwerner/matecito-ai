---
name: health-checks
depth: light
domain: observability
type: decision
source: SRE · 12-factor (production-readiness)
---

# Fase: Health checks

## Qué decide

Qué endpoints de salud expone el servicio, qué chequea cada uno, y cómo los usa el orquestador (Kubernetes, ECS, load balancer). Sin esta distinción, un restart mal configurado puede matar instancias sanas o dejar tráfico en instancias rotas.

## Preguntas

Una o dos, según haga falta.

### 1. Endpoints y qué chequean

> **Liveness** responde "¿el proceso está vivo?". **Readiness** responde "¿puede aceptar tráfico ahora?". Confundirlos causa que Kubernetes reinicie instancias que solo están esperando la DB, o que el load balancer envíe tráfico a instancias que aún no terminaron de iniciar.

- **`/health/live` (liveness) + `/health/ready` (readiness)** — *default recomendado para cualquier servicio en Kubernetes o ECS.*
- Solo `/health` (liveness básico) — válido si no hay orquestador o el entorno es simple.
- Sin health checks — *solo para scripts o herramientas sin infraestructura que los consuma.*
- No sé, recomendame.

### 2. Qué chequea cada endpoint

> Si readiness chequea dependencias críticas (DB, cache), el orquestador saca la instancia del pool hasta que estén disponibles. Si chequea de más (servicios no críticos), un outage externo saca instancias sanas. **Solo si en la 1 eligió liveness + readiness.**

Liveness chequea:
- Solo que el proceso responde (HTTP 200) — *default; no depende de nada externo.*

Readiness chequea:
- **Conexión a DB + dependencias críticas** — *default recomendado.*
- Solo conexión a DB.
- Proceso listo (igual que liveness) — si no hay dependencias externas críticas.

## Notas de lógica (para el motor)

- Si el tipo de proyecto es `script`, `librería`, o `cli`, saltá esta fase con `Status: Not Applicable` y razón.
- Si eligió "Solo `/health`", no hagas la pregunta 2.

## Qué materializar

ADR `health-checks` materializado según el template `../../templates/adr.md`. La **Decisión** captura: endpoints expuestos (`/health/live` + `/health/ready` o solo `/health`), qué chequea cada uno con la lista concreta de dependencias si se definió, los timeouts de los checks, y cómo los consume el orquestador (config de probes de Kubernetes/ECS si aplica).

**Reglas verificables** (cada una con su mecanismo al inicio):

- **[manual]** liveness no llama a ninguna dependencia externa: responde solo si el proceso está vivo (HTTP 200).
- **[manual]** readiness falla (no-200) cuando una dependencia crítica declarada (ej: conexión a DB) no está disponible, sacando la instancia del pool.
- **[manual]** readiness no chequea servicios no críticos, para que un outage externo no saque instancias sanas.
- **[tool: test]** existe un test que verifica que `/health/live` responde 200 sin tocar dependencias y que `/health/ready` refleja el estado de las dependencias críticas.

Si el tipo de proyecto es `script`, `librería` o `cli`, la fase se salta con `Status: Not Applicable` (sin Reglas verificables).
