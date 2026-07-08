---
name: folder-structure
depth: light
domain: structure
type: convention
source: práctica clásica de convenciones de proyecto · arc42 §8 (conceptos transversales)
---

# Fase: Estructura de carpetas y naming

## Qué decide

Cómo se organiza el código dentro de cada capa (por feature vs por tipo técnico) y qué **sufijo de rol** identifica cada tipo de artefacto. Decisión de baja entropía pero alta frecuencia: se aplica en cada archivo nuevo. (El *casing* de los nombres NO se decide acá — es de `code-conventions`.)

## Preguntas

Una o dos, según haga falta.

### 1. Organización dentro de cada capa

> Impacta qué tan fácil es encontrar todo lo relacionado con un feature vs todo lo relacionado con un tipo técnico.

- ***Por tipo técnico dentro de cada capa — default para proyectos con pocas features bien delimitadas.*** Ej: `services/user_service.py`, `services/order_service.py`.
- **Por feature dentro de cada capa** — ej: `users/service.py`, `orders/service.py`. Mejor para proyectos con muchos features independientes.
- **Híbrido: por feature en `application/` e `infrastructure/`, por tipo en `domain/`** — combina cohesión de feature con pureza del dominio.
- No sé, recomendame.

### 2. Sufijos de rol por tipo de artefacto

> El sufijo de un archivo (`.controller`, `.use-case`, `.repository`) marca su **rol y su capa** — es un ancla estructural verificable, no solo un nombre.

Definí qué **sufijo** identifica cada tipo de artefacto relevante al stack, y si es **obligatorio**:

- **Servicios:** `*.service` · **Repositorios:** `*.repository` · **Casos de uso / handlers:** `*.use-case` / `*.handler` · **DTOs / schemas:** `*.dto` / `*.schema` · **Entidades de dominio:** `*.entity` · **Controllers / handlers HTTP:** `*.controller` / `*.routes`

Confirmar: ¿sufijo obligatorio por tipo o sin sufijo? El **casing** de esos nombres (kebab / camel / Pascal) NO se decide acá — es una convención global de `code-conventions`. Acá va solo *qué sufijo = qué rol*.

## Notas de lógica (para el motor)

- Si en architecture-style eligió "Vertical Slice": la pregunta 1 ya está respondida implícitamente (por feature es el default de Vertical Slice). Documentarlo en el EDR como consecuencia del estilo elegido y saltar la pregunta.
- El **casing** de nombres (archivos, tipos, variables) NO se decide en esta fase: es una convención global de `code-conventions`. Acá solo se fija la taxonomía de sufijos de rol; el EDR de esta fase hereda el casing de `code-conventions`. Si `code-conventions` no se va a tratar, remití el casing a esa fase (no lo dupliques acá).

## Qué materializar

EDR `folder-structure` materializado según `~/.claude/references/edr/templates/edr.md`. Debe contener:

- **Contexto** y **Decisión**: criterio de organización dentro de cada capa (por tipo técnico vs por feature vs híbrido), las convenciones de nombres por tipo de artefacto (clase y archivo), y si los sufijos son obligatorios o no. Conservá los ejemplos concretos de paths: `src/application/users/create_user.py`, `src/domain/user.py`, `src/infrastructure/db/user_repository.py`.
- **Reglas verificables**: cada **sufijo de rol** como aserción chequeable con su mecanismo al inicio. Ej: `[tool: <linter/naming check>]` sufijo `*.routes.ts` obligatorio en handlers HTTP; `[tool: <linter>]` cada tipo de artefacto lleva su sufijo. El **casing** NO va acá — es de `code-conventions`.
- **Alcance**: como decisión estructural, incluí los globs **a nivel convención** que la decisión gobierna (ej: `src/**/<feature>/`, `src/**/*.routes.ts`, `src/domain/**`). Patrones estables, no archivos concretos.
- **Relacionados** (opcional): vinculá con `architecture-style` y `layers-and-dependencies` como decisiones de las que esta depende, y `relacionado-con` → `code-conventions` (que fija el casing de los nombres cuyos sufijos define esta fase).
