---
name: cors
depth: light
domain: security
type: policy
source: OWASP ASVS v4 §14.5 (HTTP Request Header Validation)
---

# Fase: Política CORS

## Qué decide

Qué orígenes pueden hacer requests cross-origin a la API y si se permiten credenciales (cookies, Authorization header).

## Preguntas

### 1. Orígenes permitidos y credenciales

> Una política `*` con credenciales es inválida por spec y una política laxa sin credenciales puede exponer datos a orígenes no esperados. OWASP ASVS 14.5.3 requiere validación explícita del origen.

- **Lista explícita de orígenes (ej: `https://app.midominio.com`)** — *default recomendado; credenciales solo si es necesario.*
- Origen dinámico reflejado del request — solo aceptable si se valida contra un allowlist en código.
- `*` sin credenciales — aceptable solo para APIs públicas de solo lectura.
- No aplica — la API no es consumida por browsers.
- No sé, recomendame.

## Notas de lógica (para el motor)

- Si el proyecto es `cli`, `libreria` o `script`, marcar esta fase como `Not Applicable`.
- Si el usuario elige `*` con credenciales, advertir que viola la spec de CORS y no funcionará en browsers.

## Qué materializar

ADR `cors` materializado según `~/.claude/references/adr/templates/adr.md`. Esta es una decisión de tipo `policy`; sus reglas deben quedar especialmente accionables. Debe contener:

- **Contexto**: si la API es consumida por browsers, y por qué OWASP ASVS 14.5.3 exige validación explícita del origen.
- **Decisión**: lista de orígenes permitidos (o el criterio de allowlist dinámico), si se permiten credenciales, métodos y headers habilitados, y dónde se configura (middleware de la app vs config del reverse proxy).
- **Reglas verificables** (cada una con su mecanismo):
  - `[manual]` el header `Access-Control-Allow-Origin` solo refleja orígenes presentes en la lista explícita (o en el allowlist validado en código); nunca `*` cuando se permiten credenciales.
  - `[tool: test]` una request cross-origin desde un origen fuera del allowlist no recibe headers CORS permisivos.
  - `[manual]` `Access-Control-Allow-Credentials: true` solo coexiste con orígenes explícitos, nunca con `*`.
- **Alternativas consideradas**: origen dinámico reflejado, `*` sin credenciales, y por qué se descartaron o limitaron.
- **Consecuencias**: orígenes que quedan habilitados y el riesgo de exposición si el allowlist crece sin control.
