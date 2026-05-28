---
name: deployment-topology
depth: light
domain: delivery
tipo: decisión
adr-output: deployment-topology
source: 12-factor (factor VI: processes, factor IX: disposability)
---

# Fase: Topología de deployment

## Qué decide

Dónde y cómo corre la aplicación en producción: unidad de ejecución, cantidad de instancias, y si el proceso es stateless.

## Preguntas

### 1. Unidad de ejecución

> Define el modelo operacional y condiciona cómo se escala, se actualiza y se observa la aplicación.

- **Container (Docker / OCI)** — *default para la mayoría de aplicaciones productivas; portabilidad y reproducibilidad.*
- Serverless (función: Lambda, Cloud Run, Azure Functions) — ideal para workloads event-driven o de baja frecuencia; sin gestión de instancias.
- VM / servidor bare-metal — más control, menos abstracción; común en entornos on-premise.
- PaaS gestionado (Heroku, Render, Railway, Fly.io) — sin infraestructura a gestionar; límites en configuración.
- No sé, recomendame.

### 2. Instancias y estado

> Una aplicación stateful con múltiples instancias requiere estrategias de sincronización; stateless escala horizontalmente sin fricción (12-factor factor VI).

- **Una instancia, stateless** — *default para proyectos que arrancan.*
- Múltiples instancias, stateless — escala horizontal; sesiones y cache deben ser externos.
- Una instancia, stateful — estado en memoria; sencillo pero sin alta disponibilidad.
- Múltiples instancias, stateful — requiere sticky sessions o sincronización explícita.
- No sé, recomendame.

## Qué materializar

ADR `deployment-topology` con: unidad de ejecución elegida, cantidad de instancias prevista, si el proceso es stateless o no, y la consecuencia directa sobre cómo se guarda estado de sesión o cache (referenciando `caching` si aplica).
