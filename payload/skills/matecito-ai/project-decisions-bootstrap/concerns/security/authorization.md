---
name: authorization
depth: light
domain: security
type: decision
source: OWASP ASVS v4 §4 (Access Control)
---

# Fase: Autorización

## Qué decide

Qué modelo de permisos rige el acceso a recursos o acciones, y dónde se evalúa esa lógica en el sistema.

## Preguntas

### 1. Modelo de permisos

> La elección define la complejidad del sistema de control de acceso y dónde vive esa lógica.

- **Boolean simple (puede / no puede)** — *default para proyectos pequeños con pocos roles.*
- RBAC (Role-Based): permisos asignados a roles, roles asignados a usuarios.
- ABAC (Attribute-Based): permisos según atributos del usuario, el recurso y el contexto.
- No sé, recomendame.

### 2. Punto de evaluación

> Dónde se chequean los permisos. Centralizar evita que cada endpoint lo reimplemente (y que alguno se olvide).

- **Middleware / guard global antes de llegar al handler** — *default recomendado.*
- En cada use case o service de dominio.
- Mix: guard para autenticación, dominio para reglas de negocio finas.
- No sé, recomendame.

## Notas de lógica (para el motor)

- Si el proyecto no tiene usuarios ni autenticación, marcar esta fase como `Not Applicable`.
- Si se elige ABAC, advertir al usuario que la complejidad de implementación es significativamente mayor que RBAC; pedir confirmación.

## Qué materializar

ADR `authorization` materializado según `~/.claude/references/adr/templates/adr.md`. Debe contener:

- **Contexto**: cuántos roles o niveles de acceso existen y por qué se centraliza la evaluación de permisos.
- **Decisión**: modelo elegido (boolean simple, RBAC o ABAC), descripción de los roles o atributos relevantes si aplica, y punto de evaluación (middleware/guard global, use case/service de dominio, o mix).
- **Reglas verificables** (cada una con su mecanismo):
  - `[manual]` cuando el check de permisos falla, el sistema responde con el status decidido (ej. 403 vs 404 vs redirección) de forma consistente en todos los endpoints.
  - `[manual]` la evaluación de permisos ocurre en el punto único elegido; ningún handler reimplementa la lógica de acceso.
- **Alternativas consideradas**: los otros modelos evaluados y, si se eligió ABAC, por qué se asumió su mayor complejidad.
- **Consecuencias**: costo de mantenimiento del modelo elegido y trade-off de granularidad vs complejidad.
- **Relacionados** (si aplica): `relacionado-con` → `auth.md` cuando el mecanismo de autenticación se decide allí.
