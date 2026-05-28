---
name: cors
depth: light
domain: security
tipo: política
adr-output: cors
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

ADR `cors` con: lista de orígenes permitidos (o criterio de allowlist dinámico), si se permiten credenciales, métodos y headers habilitados, y la regla de dónde se configura (middleware de la app vs config del reverse proxy).
