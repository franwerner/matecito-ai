---
name: caching
depth: light
domain: runtime
type: decision
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

- Si en la pregunta 1 eligió "Ninguno por ahora", no hagas la pregunta 2. Materializá el EDR con `Status: Pending` y razón ("sin problema de performance medido todavía"), no como decisión hueca.

## Tech a registrar

Si se elige Redis / Memcached u otra herramienta, registrala en el catálogo `tech/`.

## Qué materializar

EDR `caching` materializado según el template `~/.claude/references/edr/templates/edr.md`. La **Decisión** captura: qué se cachea y qué no, la capa elegida (in-memory / distribuido / CDN-HTTP), la estrategia de invalidación (TTL fijo / por evento / mix) y los TTLs concretos si se definieron.

**Reglas verificables** (cada una con su mecanismo al inicio):

- **[manual]** solo se cachea lo enumerado en la Decisión; cachear datos fuera de esa lista requiere actualizar el EDR.
- **[manual]** cada entrada cacheada tiene un TTL explícito o una regla de purga por evento; ninguna entrada queda sin política de expiración.
- **[manual]** en escrituras críticas la entrada correspondiente se purga o se reescribe en el mismo flujo, para no servir datos viejos.

Si se eligió "Ninguno por ahora", el EDR va con `Status: Pending` y la razón concreta ("sin problema de performance medido todavía"); en ese caso no lleva Reglas verificables.
