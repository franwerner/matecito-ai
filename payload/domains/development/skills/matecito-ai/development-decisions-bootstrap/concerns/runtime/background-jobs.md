---
name: background-jobs
depth: light
domain: runtime
type: decision
source: 12-factor (VIII: concurrency · IX: disposability)
---

# Fase: Background jobs y tareas programadas

## Qué decide

Si el proyecto necesita procesar trabajo fuera del ciclo request/response, y con qué mecanismo: cola, scheduler, o ninguno.

## Preguntas

Una o dos, según haga falta.

### 1. Mecanismo de background jobs

> Un job que corre en el proceso web rompe 12-factor (VIII) y hace el proceso no-disposable (IX). Importante saber si hay trabajo diferido antes de diseñar la infra.

- **Ninguno por ahora** — *default honesto si no hay un caso de uso concreto identificado.*
- Cola de mensajes con workers separados (RabbitMQ, SQS, Redis Streams) — desacoplado, escalable horizontalmente.
- Queue in-process (Celery, BullMQ, Sidekiq, etc.) — más simple, mismo repo, worker separado.
- Scheduler (cron-like: APScheduler, node-cron, cron job de Kubernetes) — para tareas periódicas sin cola.
- No sé, recomendame.

### 2. Estrategia de reintentos y dead-letter

> Un job que falla sin reintento pierde trabajo silenciosamente. **Solo si en la 1 eligió cola o queue in-process.**

- Reintentos con backoff exponencial + dead-letter queue (DLQ) — *default recomendado.*
- Reintentos fijos sin DLQ — simple, riesgo de perder mensajes envenenados.
- Sin reintentos — solo si el job es idempotente y la pérdida es aceptable.

## Notas de lógica (para el motor)

- Si eligió "Ninguno por ahora", no hagas la pregunta 2. Materializá el EDR con `Status: Pending` y razón ("sin caso de uso de background identificado todavía").

## Tech a registrar

Si se elige una librería o broker concreto (ej: `celery.md`, `bullmq.md`, `rabbitmq.md`, `redis.md`), registrarla en el catálogo `tech/`.

## Qué materializar

EDR `background-jobs` materializado según el template `~/.claude/references/edr/templates/edr.md`. La **Decisión** captura: mecanismo elegido (ninguno / cola de mensajes / queue in-process / scheduler), broker o librería concreta si aplica, estrategia de reintentos y dead-letter, y si los workers corren como procesos separados del proceso web (12-factor VIII/IX).

**Reglas verificables** (cada una con su mecanismo al inicio):

- **[manual]** los background jobs no corren en el proceso web: hay un proceso/worker separado por tipo de carga (12-factor VIII).
- **[manual]** todo job encolado tiene política de reintento con backoff exponencial y, salvo idempotencia explícita, una dead-letter queue donde aterrizan los mensajes envenenados.
- **[tool: test]** los jobs sin reintento son idempotentes (re-ejecutarlos no duplica efectos), verificado con un test que los corre dos veces.

Si se eligió "Ninguno por ahora", el EDR va con `Status: Pending` y la razón concreta ("sin caso de uso de background identificado todavía"); en ese caso no lleva Reglas verificables.

**Relacionados:** vincular con `resilience` si los reintentos/backoff comparten política, y con `concurrency-async` cuando los workers definen el modelo de concurrencia.
