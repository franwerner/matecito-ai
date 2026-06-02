---
name: nfr-performance
depth: light
domain: quality
type: decision
source: ISO/IEC 25010 (característica: performance efficiency — time behaviour)
---

# Fase: Performance y latencia

## Qué decide

Los objetivos cuantitativos de tiempo de respuesta y throughput del sistema. Sin números acordados, "performance aceptable" es subjetivo y no testeable.

## Preguntas

### 1. Objetivos de latencia y throughput

> ISO/IEC 25010 define "time behaviour" como el tiempo de respuesta y tasa de procesamiento bajo condiciones definidas. Sin un número acordado, cualquier medición es ambigua.

Definí al menos uno de estos (dejá en blanco los que no apliquen):

- **Latencia P99 en endpoint crítico:** `___` ms (ej: 200 ms para login, 500 ms para búsqueda).
- **Latencia P50 (mediana):** `___` ms.
- **Throughput mínimo sostenido:** `___` req/s (ej: 100 req/s en carga normal).
- Sin objetivos formales por ahora — se mide cuando haya usuarios reales.
- No sé, recomendame un punto de partida razonable.

### 2. Presupuesto de respuesta y alertas

> **Solo si definió al menos un objetivo en la pregunta 1.** Un objetivo sin monitoreo ni alerta es decorativo.

- **Alerta cuando P99 supera el umbral** — requiere métricas (ver `metrics`).
- Revisión periódica manual (mensual, por release).
- Test de carga en CI en cada release.
- Sin mecanismo de alerta formal por ahora.

## Notas de lógica (para el motor)

- Si elige "Sin objetivos formales por ahora", no hacer la pregunta 2. Materializar el ADR con `Status: Pending` y motivo.

## Qué materializar

ADR `nfr-performance` materializado según el template `../../templates/adr.md`.

- **Contexto:** por qué se fijan objetivos de performance (expectativa de carga, criticidad del flujo) y bajo qué condiciones de medición se entienden.
- **Decisión:** los objetivos cuantitativos acordados, el endpoint o flujo de referencia para medirlos, y el mecanismo de detección de regresión elegido.
- **Reglas verificables:** cada objetivo se reformula como una aserción con valor concreto y su mecanismo. Ejemplos:
  - **[tool: test de carga]** la latencia P99 del endpoint `<ref>` se mantiene ≤ `___` ms bajo `___` req/s.
  - **[tool: test de carga]** la latencia P50 (mediana) del flujo `<ref>` se mantiene ≤ `___` ms.
  - **[tool: test de carga]** el throughput sostenido es ≥ `___` req/s en carga normal.
  - **[tool: alertas/métricas]** se dispara alerta cuando P99 supera el umbral (depende de `metrics`).
  - **[manual]** revisión periódica de las métricas (mensual / por release) cuando no hay alerta automática.

  Si se eligió "sin objetivos formales por ahora", materializar con `Status: Pending` y el motivo/trigger esperado, sin inventar números ni reglas.
- `Relacionados`: `depende-de` → `metrics` si la detección de regresión requiere instrumentación de métricas.
