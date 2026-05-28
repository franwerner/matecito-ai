---
name: input-validation
depth: light
domain: security
tipo: política
adr-output: input-validation
source: OWASP ASVS v4 §5 (Validation, Sanitization and Encoding)
---

# Fase: Validación de input

## Qué decide

Dónde y cómo se valida y sanitiza el input externo antes de que entre al sistema, y qué se hace cuando falla la validación.

## Preguntas

### 1. Capa y herramienta de validación

> El input no confiable validado tarde o en múltiples lugares crea superficies de inyección. Centralizarlo en el borde del sistema es la defensa principal (OWASP ASVS 5.1).

- **Validación en el borde (controller / handler) con schema declarativo** — *default recomendado; ej: Zod, Pydantic, Joi, class-validator.*
- Validación manual en cada endpoint — propenso a inconsistencias.
- Mix: schema en el borde + validaciones de negocio en dominio.
- No sé, recomendame.

### 2. Respuesta ante input inválido

> Define qué ve el cliente cuando el input no pasa validación. Importante para no filtrar internals.

- **400 Bad Request con descripción del campo que falló, sin stack trace** — *default recomendado.*
- 422 Unprocessable Entity (semánticamente más preciso para REST).
- No sé, recomendame.

## Tech a registrar

Si se elige una librería de validación (Zod, Pydantic, Joi, class-validator, etc.), registrarla en el catálogo `tech/`.

## Qué materializar

ADR `input-validation` con: capa donde se valida, herramienta elegida, política de respuesta ante falla, y regla explícita sobre qué NUNCA se retorna al cliente (stack traces, mensajes de error internos, detalles de base de datos).
