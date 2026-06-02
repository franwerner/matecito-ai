---
name: input-validation
depth: light
domain: security
type: policy
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

ADR `input-validation` materializado según `../../templates/adr.md`. Esta es una decisión de tipo `policy`; sus reglas deben quedar especialmente accionables. Debe contener:

- **Contexto**: por qué el input no confiable validado tarde o en múltiples lugares crea superficies de inyección, y por qué la defensa principal es validar en el borde (OWASP ASVS 5.1).
- **Decisión**: capa donde se valida (borde con schema declarativo, manual por endpoint, o mix), herramienta elegida (ej. Zod, Pydantic, Joi, class-validator), y política de respuesta ante input inválido (400 con descripción del campo, o 422).
- **Reglas verificables** (cada una con su mecanismo):
  - `[tool: test]` el input inválido recibe el status decidido (400/422) con la descripción del campo que falló.
  - `[tool: test]` la respuesta de error NUNCA incluye stack traces, mensajes de error internos ni detalles de base de datos.
  - `[manual]` todo input externo pasa por el schema declarativo en la capa de borde antes de entrar al sistema.
- **Alternativas consideradas**: validación manual por endpoint y el otro código de status, con su trade-off de consistencia/precisión semántica.
- **Consecuencias**: dependencia de la librería de validación elegida y disciplina requerida para mantener los schemas en el borde.
