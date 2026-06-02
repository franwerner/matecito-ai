---
name: concurrency-async
depth: light
domain: runtime
type: decision
source: 12-factor (VIII: concurrency)
---

# Fase: Modelo de concurrencia / async

## Qué decide

Cómo el proyecto maneja operaciones que no son estrictamente secuenciales: async nativo del lenguaje, threads, workers de proceso, o simplemente síncrono directo.

## Preguntas

Una o dos, según haga falta.

### 1. Modelo de concurrencia

> Afecta librerías disponibles, testing, y cómo se diseña el código I/O-bound (DB, HTTP externo, filesystem).

- **Síncrono directo** — *default honesto para proyectos CRUD simples sin latencia crítica.*
- Async nativo del lenguaje (async/await — asyncio, tokio, Node event loop) — para I/O-bound con baja latencia esperada.
- Threads (threading estándar) — para CPU-bound ligero o compatibilidad con librerías bloqueantes.
- Workers de proceso (fork/spawn, multiprocessing) — para CPU-bound pesado o aislamiento por proceso (12-factor: scale out por tipo).
- No sé, recomendame.

### 2. Política de mezcla sync/async

> Mezclar sync y async en el mismo proceso genera deadlocks y bugs difíciles de diagnosticar. **Solo si en la 1 eligió async nativo.**

- Async puro — ningún llamado bloqueante en el event loop; lo sync se corre en executor/threadpool.
- Pragmático — solo las rutas críticas son async; el resto sync (con riesgo controlado).

## Notas de lógica (para el motor)

- Si el stack detectado en Fase 0 es Node.js: sugerí async como default (el event loop es la base del ecosistema).
- Si el stack es Python con FastAPI: sugerí async. Con Django o Flask sin ASGI: sugerí sync directo salvo que el usuario tenga un caso concreto.
- Si eligió "Síncrono directo", no hagas la pregunta 2.

## Tech a registrar

Si elige un runtime async con librería concreta (ej: `asyncio`, `tokio`, `trio`, `uvloop`), registrala en el catálogo `tech/`.

## Qué materializar

ADR `concurrency-async` materializado según el template `../../templates/adr.md`. La **Decisión** captura: modelo elegido (síncrono directo / async nativo / threads / workers de proceso), la razón basada en el tipo de carga esperada (I/O-bound vs CPU-bound), la política de mezcla sync/async si aplica (async puro vs pragmático) y la tech concreta si se registró.

**Reglas verificables** (cada una con su mecanismo al inicio):

- **[manual]** si se eligió async puro: ningún llamado bloqueante corre directo en el event loop; todo lo síncrono se delega a un executor/threadpool.
- **[tool: linter]** si el ecosistema tiene reglas de async sin esperar (ej: `no-floating-promises` / `await` faltante), están activas y sin excepciones sin justificar.
- **[manual]** el modelo elegido es consistente con el tipo de carga declarado (I/O-bound → async/threads; CPU-bound pesado → workers de proceso).

**Relacionados:** vincular con `background-jobs` cuando los workers de proceso materializan también el modelo de jobs.
