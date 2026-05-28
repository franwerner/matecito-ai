---
name: feature-flags
depth: light
domain: delivery
tipo: decisión
adr-output: feature-flags
source: continuous delivery / trunk-based development
---

# Fase: Feature flags

## Qué decide

Si el proyecto usa flags de features, con qué mecanismo, y cómo se nombran y eliminan para evitar deuda técnica acumulada.

## Preguntas

### 1. Mecanismo de feature flags

> Los feature flags permiten deployar código inactivo y activarlo sin nuevo deploy. Si se usan sin convención, se acumulan y nadie sabe cuáles siguen activos.

- **Ninguno por ahora** — *default honesto si no hay necesidad identificada; se puede sumar después.*
- Variables de entorno / config — simple, sin UI; útil para flags de larga duración (ej: habilitar una integración).
- Librería in-process (ej: `unleash-client`, `flagsmith`, `flipt`, `LaunchDarkly SDK`) — evaluación local con reglas; requiere sincronización periódica.
- Servicio externo gestionado (LaunchDarkly, Flagsmith cloud, Statsig) — UI de gestión, targeting por usuario/segmento; costo operacional.
- No sé, recomendame.

### 2. Naming y ciclo de vida

> **Solo si eligió usar flags.** Un flag sin fecha de expiración es deuda técnica garantizada.

- Prefijo con tipo + ticket: `feat_<ticket>_<nombre>` / `exp_<ticket>_<nombre>` (feature / experiment).
- Convención libre con revisión periódica (sprint review, trimestral).
- Sin convención formal por ahora.

## Notas de lógica (para el motor)

- Si elige "Ninguno por ahora", no hacer la pregunta 2. Materializar el ADR con `Status: Pending` y motivo "sin necesidad identificada aún".

## Tech a registrar

Si se elige una librería o servicio de feature flags, registrarlo en `tech/`.

## Qué materializar

ADR `feature-flags` con: mecanismo elegido (o decisión explícita de no usar), convención de naming, regla de eliminación de flags (ej: "todo flag se elimina en el sprint siguiente a su activación definitiva"), y quién tiene permiso de cambiar un flag en producción.
