---
name: caching
depth: light
domain: runtime
tipo: decisión
adr-output: caching
source: checklists de production-readiness / SRE
---

# Fase: Caching

## Qué decide

Qué se cachea, dónde, y cómo se invalida. Mal hecho sirve datos viejos; bien hecho define latencia y costo.

## Preguntas

Una o dos, según haga falta.

### 1. Capa de cache

> Dónde vive el cache.

- **Ninguno por ahora** — *default honesto si no hay un problema de performance medido.*
- In-memory en proceso (ej: LRU) — simple, no compartido entre instancias.
- Cache distribuido (Redis / Memcached) — compartido, sobrevive reinicios.
- CDN / cache HTTP en el borde — para respuestas públicas.
- No sé, recomendame.

### 2. Estrategia de invalidación

> El problema difícil del caching. **Solo si en la 1 eligió cachear algo.**

- TTL fijo (expira por tiempo).
- Invalidación por evento (al escribir, se purga).
- Mix (TTL + purga en escrituras críticas).

## Notas de lógica (para el motor)

- Si en la pregunta 1 eligió "Ninguno por ahora", no hagas la pregunta 2. Materializá el ADR con `Status: Pending` y razón ("sin problema de performance medido todavía"), no como decisión hueca.

## Tech a registrar

Si se elige Redis / Memcached u otra herramienta, registrala en el catálogo `tech/`.

## Qué materializar

ADR `caching` con: qué se cachea y qué no, la capa elegida, la estrategia de invalidación, y TTLs concretos si se definieron.
