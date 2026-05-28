---
name: configuration
depth: light
domain: delivery
tipo: convención
adr-output: configuration
source: 12-factor (factor III: config)
---

# Fase: Configuración

## Qué decide

Cómo la aplicación lee su configuración (no secretos) según el entorno, y si esa configuración se valida y tipiea al arranque.

## Preguntas

### 1. Mecanismo de configuración

> Determina de dónde lee la app sus valores de entorno y cómo se distribuyen entre entornos.

- **Variables de entorno puras** — *default 12-factor; portable y sin archivos en disco.*
- `.env` + variables de entorno — env vars como fuente de verdad, `.env` solo para desarrollo local (en `.gitignore`).
- Archivos de config por entorno (yaml / toml / json) — útil cuando hay mucha config estructurada; riesgo de filtrar valores si no se cuidan.
- No sé, recomendame.

### 2. Validación y tipado al startup

> Si la app arranca con config incompleta o malformada, el error aparece tarde y en producción. Validar al arranque lo convierte en un fallo rápido y obvio.

- **Sí, validación + schema tipado al startup** — *default recomendado; rompe temprano con mensaje claro.*
- Solo validación básica (variables presentes, sin tipos).
- Sin validación — los errores se detectan en runtime.
- No sé, recomendame.

## Tech a registrar

Si se elige una librería de config tipada (ej: `pydantic-settings` para Python, `zod` con `dotenv` para TS/JS, `viper` para Go), registrarla en el catálogo `tech/`.

## Qué materializar

ADR `configuration` con: mecanismo elegido, qué archivos se commitean y cuáles no, si hay validación al startup, librería de schema tipado (si aplica), y una regla explícita sobre dónde viven las variables por entorno. Secretos excluidos de este ADR — ver `secrets-management`.
