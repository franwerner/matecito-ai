---
name: data-modeling
depth: light
domain: data
type: decision
source: práctica clásica de modelado de datos
---

# Fase: Modelado de datos

## Qué decide

Convenciones de bajo nivel que afectan esquema de DB, APIs, y código de dominio: tipo de IDs, borrado lógico vs físico, timestamps estándar, y si el modelo soporta multitenancy.

## Preguntas

Una o dos, según haga falta.

### 1. Tipo de IDs y borrado

> El tipo de ID afecta performance (UUID vs serial en índices), privacidad (no exponer secuencias), y portabilidad. El soft delete afecta todas las queries futuras si no se define desde el inicio.

- **Autoincrement / serial** — *default simple para proyectos sin distribución ni privacidad de IDs.*
- UUID v4 — no secuencial, adecuado cuando los IDs se exponen externamente o hay múltiples fuentes de datos.
- UUID v7 — ordenable por tiempo, combina ventajas de serial y UUID v4 (recomendado para sistemas nuevos).
- ULID o CUID2 — alternativas a UUID v7 con mejor legibilidad.
- No sé, recomendame.

(Respuesta separada, mismo turno o siguiente:)

- **Borrado físico (DELETE real)** — *default honesto si no hay requisito de auditoría.*
- Soft delete (`deleted_at` nullable) — necesario para auditoría o recuperación; agrega complejidad a todas las queries.

### 2. Timestamps y multitenancy

> Timestamps estándar (`created_at`, `updated_at`) son fáciles de agregar ahora y muy costosos de migrar después. Multitenancy en DB es una decisión irreversible de esquema.

- **Timestamps estándar en todas las entidades** (`created_at`, `updated_at`) — *default recomendado siempre.*
- Sin timestamps — solo si el dominio lo justifica explícitamente.

(Respuesta separada:)

- **Sin multitenancy** — *default para la mayoría de proyectos.*
- Multitenancy por `tenant_id` en cada tabla — simple, un solo schema.
- Schema separado por tenant — aislamiento fuerte, complejidad operativa alta.
- Base de datos separada por tenant — aislamiento máximo, costo operativo muy alto.

## Notas de lógica (para el motor)

- Si el proyecto es `script` o `librería`, saltá esta fase completa con `Status: Not Applicable` y razón.
- La segunda parte de cada pregunta (borrado / multitenancy) puede hacerse en el mismo turno que la primera si son cortas.

## Qué materializar

ADR `data-modeling` materializado según el template `~/.claude/references/adr/templates/adr.md`. La **Decisión** captura: tipo de ID elegido y su justificación (autoincrement / UUID v4 / UUID v7 / ULID-CUID2), la política de borrado (físico vs soft delete con razón), la convención de timestamps, y la estrategia de multitenancy si aplica.

**Reglas verificables** (cada una con su mecanismo al inicio):

- **[manual]** toda tabla usa el tipo de ID decidido como clave primaria; no se mezclan tipos de ID entre entidades sin justificación en el ADR.
- **[manual]** si se eligieron timestamps estándar: toda tabla lleva `created_at` y `updated_at` NOT NULL.
- **[manual]** si se eligió soft delete: las tablas afectadas tienen columna `deleted_at` nullable y las queries por default excluyen las filas con `deleted_at IS NOT NULL`.
- **[manual]** si hay multitenancy por `tenant_id`: toda tabla multi-tenant incluye `tenant_id` y ninguna query cruza tenants sin filtrarlo.

Si el proyecto es `script` o `librería`, la fase se salta con `Status: Not Applicable` (vive como fila en el INDEX del dominio, sin Reglas verificables).

**Relacionados:** vincular con `data-access` (las migraciones materializan estas convenciones de esquema).
