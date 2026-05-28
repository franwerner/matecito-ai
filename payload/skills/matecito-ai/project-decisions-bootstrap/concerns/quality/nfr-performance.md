---
name: nfr-performance
depth: light
domain: quality
tipo: decisión
adr-output: nfr-performance
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

ADR `nfr-performance` con: los objetivos cuantitativos acordados (latencia P50/P99, throughput), el endpoint o flujo de referencia para medirlos, y cómo se detecta una regresión (alerta, test de carga, revisión manual).
