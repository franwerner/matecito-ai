---
name: feature-flags
depth: light
domain: delivery
type: decision
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

- Si elige "Ninguno por ahora", no hacer la pregunta 2. Materializar el EDR con `Status: Pending` y motivo "sin necesidad identificada aún".

## Tech a registrar

Si se elige una librería o servicio de feature flags, registrarlo en `tech/`.

## Qué materializar

EDR `feature-flags` materializado según `~/.claude/references/edr/templates/edr.md`. Debe contener:

- **Contexto** y **Decisión**: mecanismo elegido (o decisión explícita de no usar, con `Status: Pending` y motivo si es "ninguno por ahora"), la convención de naming, y quién tiene permiso de cambiar un flag en producción.
- **Reglas verificables**: las convenciones de naming y ciclo de vida como aserciones con su mecanismo al inicio. Ej: `[manual]` todo flag sigue el prefijo `feat_<ticket>_<nombre>` / `exp_<ticket>_<nombre>`; `[manual]` todo flag se elimina en el sprint siguiente a su activación definitiva; `[manual]` solo <rol> puede cambiar un flag en producción. Usá `[tool: <linter/test>]` si el naming o la expiración es chequeable automáticamente. Conservá los valores concretos del prefijo y la regla de expiración.
