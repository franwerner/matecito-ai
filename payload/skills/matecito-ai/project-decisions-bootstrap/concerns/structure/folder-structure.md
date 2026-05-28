---
name: folder-structure
depth: light
domain: structure
tipo: convención
adr-output: folder-structure
source: práctica clásica de convenciones de proyecto · arc42 §8 (conceptos transversales)
---

# Fase: Estructura de carpetas y naming

## Qué decide

Cómo se organiza el código dentro de cada capa (por feature vs por tipo técnico) y las convenciones de nombres por tipo de artefacto. Decisión de baja entropía pero alta frecuencia: se aplica en cada archivo nuevo.

## Preguntas

Una o dos, según haga falta.

### 1. Organización dentro de cada capa

> Impacta qué tan fácil es encontrar todo lo relacionado con un feature vs todo lo relacionado con un tipo técnico.

- ***Por tipo técnico dentro de cada capa — default para proyectos con pocas features bien delimitadas.*** Ej: `services/user_service.py`, `services/order_service.py`.
- **Por feature dentro de cada capa** — ej: `users/service.py`, `orders/service.py`. Mejor para proyectos con muchos features independientes.
- **Híbrido: por feature en `application/` e `infrastructure/`, por tipo en `domain/`** — combina cohesión de feature con pureza del dominio.
- No sé, recomendame.

### 2. Convenciones de nombres por tipo de artefacto

> Sin convención explícita cada dev nombra diferente; la búsqueda por nombre deja de funcionar.

Preguntá y confirmá para cada tipo relevante al stack:

- **Servicios:** `UserService` / `user_service.py` / `user.service.ts`
- **Repositorios:** `UserRepository` / `user_repository.py` / `user.repository.ts`
- **Casos de uso / handlers:** `CreateUserUseCase` / `create_user.py` / `create-user.handler.ts`
- **DTOs / schemas:** `UserDTO` / `UserSchema` / `CreateUserRequest`
- **Entidades de dominio:** `User` / `user.entity.ts`
- **Controllers / handlers HTTP:** `UserController` / `user_router.py` / `user.controller.ts`

Confirmar: ¿snake_case, camelCase o PascalCase por tipo? ¿Con sufijo obligatorio (`Service`, `Repository`, `Controller`) o sin sufijo?

## Notas de lógica (para el motor)

- Si en architecture-style eligió "Vertical Slice": la pregunta 1 ya está respondida implícitamente (por feature es el default de Vertical Slice). Documentarlo en el ADR como consecuencia del estilo elegido y saltar la pregunta.
- El default de casing depende del lenguaje detectado en Fase 0: Python → snake_case para archivos y variables, PascalCase para clases; JS/TS → camelCase para variables, PascalCase para clases, kebab-case para archivos (con o sin sufijo según el framework); Java/C# → PascalCase para clases y archivos. Mostrar el default del lenguaje como propuesta y pedir confirmación.

## Qué materializar

ADR `folder-structure` con: criterio de organización dentro de cada capa (por tipo técnico vs por feature), tabla de convenciones de nombres por tipo de artefacto (clase y archivo), y si los sufijos son obligatorios o no. Preferiblemente con ejemplos concretos de paths: `src/application/users/create_user.py`, `src/domain/user.py`, `src/infrastructure/db/user_repository.py`.
