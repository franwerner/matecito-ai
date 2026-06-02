---
name: layers-and-dependencies
depth: deep
domain: structure
type: decision
source: práctica clásica de arquitectura en capas · arc42 §5 (vista de building blocks)
---

# Fase: Capas y reglas de dependencia

## Qué decide

Los nombres concretos de cada capa del sistema y las reglas de dependencia entre ellas, escritas de forma verificable. Es la decisión más importante para mantener la arquitectura íntegra con el tiempo.

## Preguntas

Una por turno. Para cada una: línea de "por qué importa", opciones con default.

### 1. Nombres de capas

> Los nombres concretos son los que van al código, al linter y a las revisiones. Nombres vagos producen reglas vagas.

Basándote en el patrón elegido en architecture-style, proponé las capas con nombres de carpeta concretos y pedí confirmación. Ejemplo para Clean Architecture:

> Te propongo estas capas:
> - `domain/` — entidades, value objects, lógica pura de negocio
> - `application/` — casos de uso, orquestación, interfaces de repos
> - `infrastructure/` — implementaciones concretas (DB, HTTP clients, filesystem)
> - `presentation/` — controllers, CLI handlers, schemas/DTOs de API
>
> ¿Confirmás, querés renombrar, o agregar/quitar alguna?

Ejemplo para Layered N-tier:

> Te propongo:
> - `controllers/` — puntos de entrada HTTP/CLI
> - `services/` — lógica de negocio y orquestación
> - `repositories/` — acceso a datos
> - `models/` — entidades / ORM models
>
> ¿Confirmás o ajustás?

Ejemplo para Vertical Slice:

> Te propongo organizar por feature en lugar de por capa técnica:
> - `features/<nombre-del-feature>/` — todo lo del feature junto (handler, service, repo, DTOs, tests)
> - `shared/` — código compartido entre features (errors, utils, auth middleware)
>
> ¿Confirmás o ajustás?

### 2. Reglas de dependencia

> Las reglas escritas como paths/globs son las que se pueden verificar con un linter. Sin eso, son solo intenciones.

Mostrá las reglas como lista en formato path/glob verificable y pedí confirmación. Ejemplo para Clean Architecture con src/ como raíz:

> Reglas propuestas:
> - `src/domain/**` solo puede importar de `src/domain/**`
> - `src/application/**` puede importar de `src/domain/**` y `src/application/**`
> - `src/infrastructure/**` puede importar de cualquier capa (implementa interfaces)
> - `src/presentation/**` puede importar de `src/application/**` y `src/domain/**`
> - **Prohibido:** `src/presentation/**` → `src/infrastructure/**` directo
> - **Prohibido:** `src/domain/**` → cualquier módulo de framework externo
>
> ¿Ajustás alguna?

## Notas de lógica (para el motor)

- **Proponer capas según el estilo elegido en architecture-style:** no pedir al usuario que defina las capas desde cero. Derivar la propuesta del patrón y pedir confirmación/ajuste.
- **Escribir las reglas de dependencia como globs/paths verificables:** formato `src/<capa>/**` o `<módulo>/**` según la convención del proyecto. Las reglas vagas ("no acoplar capas") no van al ADR.
- **Mencionar enforcement al cerrar la fase:** una vez confirmadas las reglas, informar al usuario que se pueden enforcar con:
  - `import-linter` (Python) — archivo `.importlinter` con secciones `[importlinter:contract:...]`
  - `dependency-cruiser` (JS/TS) — archivo `.dependency-cruiser.js` con reglas `forbidden`
  - `ArchUnit` (Java) — tests de arquitectura en el suite de unit tests
  - No forzar que lo configuren ahora; solo mencionarlo para que sepan que la decisión es accionable.
- **Si eligió Vertical Slice:** las "reglas de dependencia" cambian: la restricción es que los features no se importen entre sí directamente; solo comparten a través de `shared/`. Adaptar la propuesta y el ADR.
- **Si eligió Sin patrón formal:** esta fase es Not Applicable. Crear el ADR con ese status y el motivo.

## Qué materializar

ADR `layers-and-dependencies` materializado según `../../templates/adr.md`. Debe contener:

- **Contexto** y **Decisión**: lista de capas con nombre de carpeta y responsabilidad de cada una.
- **Reglas verificables**: las reglas de dependencia escritas como paths/globs (qué puede importar qué, qué está explícitamente prohibido), cada una como aserción con su mecanismo al inicio según la herramienta de enforcement elegida o disponible. Ej: `[tool: dependency-cruiser]` ningún import desde `src/presentation/**` hacia `src/infrastructure/**`; `[tool: import-linter]` `src/domain/**` solo importa de `src/domain/**`; `[tool: ArchUnit]` para Java. Usá `[manual]` solo si no hay herramienta disponible todavía. Nombrá la herramienta concreta (`import-linter` / `dependency-cruiser` / `ArchUnit` / `deptrac`) en cada `[tool: ...]` para que en sesiones futuras se sepa cómo verificar las reglas.
- **Alcance**: como decisión estructural, incluí los globs **a nivel convención** que delimitan cada capa (ej: `src/domain/**`, `src/application/**`, `src/infrastructure/**`, `src/presentation/**`, o `features/<feature>/**` y `shared/**` en Vertical Slice). Patrones estables, no archivos concretos.
- **Relacionados** (opcional): vinculá con `architecture-style` (decisión de la que depende) y con `arch-enforcement` (que ejecuta estas reglas en CI).
