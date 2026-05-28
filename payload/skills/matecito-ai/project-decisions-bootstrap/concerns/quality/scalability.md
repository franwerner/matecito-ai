---
name: scalability
depth: light
domain: quality
tipo: decisión
adr-output: scalability
source: ISO/IEC 25010 (característica: performance efficiency — capacity)
---

# Fase: Escalabilidad

## Qué decide

El modelo de escalado esperado (vertical u horizontal) y si la arquitectura lo soporta desde el inicio. ISO/IEC 25010 lo define bajo "capacity": grado en que los límites del sistema cubren los requisitos.

## Preguntas

### 1. Modelo de escalado

> La diferencia entre escalar verticalmente (más CPU/RAM a la misma instancia) y horizontalmente (más instancias en paralelo) determina si el proceso debe ser stateless y cómo se gestiona la sesión y el cache.

- **Vertical por ahora, horizontal si se necesita** — *default para proyectos que arrancan; sin compromisos de arquitectura adicionales.*
- Horizontal desde el inicio — stateless obligatorio; sesiones y cache externos (ver `caching`, `deployment-topology`).
- Sin escalado previsto — instancia única, carga controlada.
- No sé, recomendame.

### 2. Proceso stateless

> **Solo si eligió escala horizontal.** Un proceso stateless puede reiniciarse y reemplazarse sin pérdida de datos; es el requisito base de escala horizontal (12-factor factor VI).

- **Sí, stateless — todo estado externo** (DB, cache distribuido, storage de sesión).
- Stateful con estado compartido explícito — sesiones en Redis u otro store compartido.

## Notas de lógica (para el motor)

- Si elige "Vertical por ahora" o "Sin escalado previsto", no hacer la pregunta 2. Registrar el ADR con el modelo elegido y la condición de revisión ("revisar si se supera X usuarios / Y req/s").
- Si ya se definió `deployment-topology` con "múltiples instancias", proponer "horizontal desde el inicio" como default en la pregunta 1.

## Qué materializar

ADR `scalability` con: modelo de escalado elegido, si el proceso es stateless o no, dónde vive el estado si hay múltiples instancias, y la condición que dispararía una revisión de esta decisión (umbral de carga, crecimiento de equipo, etc.).
