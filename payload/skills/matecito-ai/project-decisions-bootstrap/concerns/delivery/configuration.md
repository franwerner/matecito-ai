---
name: configuration
depth: light
domain: delivery
type: convention
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

ADR `configuration` materializado según `../../templates/adr.md`. Debe contener:

- **Contexto** y **Decisión**: mecanismo de configuración elegido (env vars puras / `.env` + env vars / archivos por entorno), qué archivos se commitean y cuáles no, si hay validación al startup, y la librería de schema tipado si aplica (`pydantic-settings`, `zod` + `dotenv`, `viper`, etc.).
- **Reglas verificables**: cada política como aserción con su mecanismo al inicio. Ej: `[tool: <gitignore/CI check>]` `.env` está en `.gitignore` y nunca se commitea; `[tool: <schema lib>]` la app valida y tipa la config al arranque y aborta con mensaje claro si falta una variable; `[manual]` las variables por entorno viven en env vars como fuente de verdad, `.env` solo para desarrollo local. Nombrá la librería elegida en el `[tool: ...]` de la validación.
- **Relacionados** (opcional): vinculá con `secrets-management` (secretos excluidos de este ADR) y `deployment-topology` si la config depende del entorno de ejecución.
