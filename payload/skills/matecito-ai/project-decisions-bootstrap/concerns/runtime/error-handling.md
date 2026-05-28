---
name: error-handling
depth: deep
domain: runtime
tipo: decisión
adr-output: error-handling
source: práctica clásica de manejo de errores · arc42 §8 (conceptos transversales)
---

# Fase: Manejo de errores

## Qué decide

Cómo se representan, propagan y responden los errores en todo el sistema. Es de las decisiones más transversales: toca dominio, infraestructura y bordes.

## Preguntas

Una por turno. Para cada una: línea de "por qué importa", opciones con default, y siempre la opción "no sé, recomendame".

### 1. Estilo de errores

> Define la forma base de comunicar fallas en todo el código.

- **Excepciones** — *default para Python / Java / Node / C#.*
- **Result / Either types** — *default para Rust / Go; opcional en otros con librería.*
- **Mix pragmático** — excepciones para fallas inesperadas, Result para flujos de error esperados.
- No sé, recomendame.

### 2. Boundary handling

> Dónde se atrapan los errores que escapan de las capas internas.

- **Middleware / interceptor global** — *default para frameworks web.*
- Cada controller individual.
- Mix: middleware para crashes, controllers para errores de negocio.
- No sé, recomendame.

### 3. Errores de dominio custom

> Si hay una jerarquía propia de errores de negocio.

- Sí, jerarquía completa con base class de dominio.
- Solo para errores de negocio importantes.
- No, usamos los nativos del lenguaje.

### 4. Formato de respuesta de error

> Cómo ve el cliente un error. **Solo si el tipo de proyecto es `api-rest`, `api-graphql` o `microservicio`** (si no aplica, omitir esta pregunta).

- **RFC 7807 Problem Details** — *default recomendado para REST.*
- Formato custom JSON (`{error, code, details}`).
- Texto plano — *no recomendado para APIs serias.*

### 5. Política de logging de errores

> Qué se registra y qué nunca.

- Qué se loggea: todos / solo 5xx / sin PII.
- Nivel por tipo: error vs warn vs info.
- Qué NUNCA se loggea: passwords, tokens, datos personales, payloads completos.

## Notas de lógica (para el motor)

- **Default según stack:** mirá el lenguaje detectado en Fase 0. Si es Rust o Go → proponé `Result` en la pregunta 1. Si es Python / Java / Node / C# → proponé `Excepciones`. Mostrá el default y pedí confirmación.
- **Pregunta 4 condicional:** salteala si el proyecto no expone una API. No la registres como omisión; simplemente no aplica a la decisión.

## Tech a registrar

Si se elige una librería específica de Result/errores (ej: `returns` en Python, `neverthrow` en TS), registrala en el catálogo `tech/` con el flujo de tecnologías del motor.

## Qué materializar

ADR `error-handling` con: estilo elegido, dónde se hace boundary handling, jerarquía de errores (con nombres concretos si aplica: `UserNotFoundError`, `InsufficientFundsError`), formato de respuesta, y la política de logging escrita como reglas verificables (qué se loggea, en qué nivel, qué nunca).
