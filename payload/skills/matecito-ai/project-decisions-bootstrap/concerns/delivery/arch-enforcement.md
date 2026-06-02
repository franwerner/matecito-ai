---
name: arch-enforcement
depth: light
domain: delivery
type: policy
source: práctica de architecture-as-code / CI quality gates
---

# Fase: Enforcement de arquitectura

## Qué decide

Cómo convertir las reglas de dependencia definidas en `layers-and-dependencies` en configuración ejecutable de un linter que corra en CI y falle el build si se viola una regla.

## Preguntas

### 1. Herramienta de enforcement según stack

> Sin enforcement automatizado, las reglas de capas son recomendaciones que se degradan con el tiempo. El linter las vuelve verificables en cada PR.

- **`import-linter`** — *default para Python; reglas en `setup.cfg` o `.importlinter`.*
- **`dependency-cruiser`** — *default para JS / TS; config en `.dependency-cruiser.js`.*
- **`ArchUnit`** — *default para Java / Kotlin; tests en JUnit.*
- **`deptrac`** — alternativa para PHP.
- Sin enforcement por ahora — solo convención documentada.
- No sé, recomendame.

### 2. Integración en CI

> Que el linter exista en local pero no bloquee el merge no da garantías reales.

- **Sí, bloquea el merge en CI** — *default recomendado.*
- Solo corre localmente (pre-commit hook).
- No por ahora.

## Notas de lógica (para el motor)

- Esta fase requiere que `layers-and-dependencies` esté en status `Accepted` con reglas definidas. Si no hay capas definidas o están `Pending`, marcar esta fase como `Pending` con motivo "sin reglas de dependencia para enforcer".
- Las reglas a enforcer (paths/globs tipo "X solo importa de Y", "prohibido X → Z") ya están acordadas en el ADR `layers-and-dependencies`. No hay que volver a pedirlas — se leen de ese ADR.
- El mapeo de herramienta por stack (para la pregunta 1): Python → `import-linter`; JS/TS → `dependency-cruiser`; Java/Kotlin → `ArchUnit`; PHP → `deptrac`. Mostralo como default según el lenguaje detectado.
- Si el repo ya tiene código, después de implementar el enforcement ofrecé correr el linter una vez para confirmar que la configuración es válida y que el código actual respeta las reglas. Si es greenfield, queda listo para cuando haya código.

## Tech a registrar

La herramienta elegida (ej: `import-linter.md`, `dependency-cruiser.md`, `archunit.md`, `deptrac.md`).

## Qué materializar

> **Nota sobre el artefacto.** Esta fase es de las pocas cuyo output incluye un archivo de configuración ejecutable, no solo una decisión escrita. El ADR documenta *qué* se decidió; la **config concreta del linter la escribe el agente en el repo** al implementar, traduciendo las reglas reales del ADR `layers-and-dependencies` a la sintaxis de la herramienta elegida. El concern NO trae plantillas de config hardcodeadas — eso es trabajo de implementación, guiado por el ADR.

ADR `arch-enforcement` materializado según `../../templates/adr.md`. Debe contener:

- **Contexto** y **Decisión**: la herramienta de enforcement elegida y por qué (normalmente el default del stack), si corre en CI bloqueando el merge / solo localmente / todavía no, y la referencia al ADR `layers-and-dependencies` como origen de las reglas que el linter traduce (no duplicar las reglas acá).
- **Reglas verificables**: expresá las garantías que da esta decisión como aserciones con su mecanismo al inicio, nombrando la herramienta elegida. Ej: `[tool: dependency-cruiser]` el step de arch-lint corre en CI y bloquea el merge ante cualquier violación; `[tool: import-linter]` la config existe en su ubicación esperada y el comando definido la ejecuta. Usá `[manual]` solo si por ahora es convención documentada sin check.
- **Alcance**: como decisión estructural, incluí la **ubicación esperada de la config** y los globs **a nivel convención** que el enforcement cubre (ej: `.importlinter`, `.dependency-cruiser.js`, test de ArchUnit, `deptrac.yaml`; y `src/**` como superficie analizada). Indicá el **comando** que la ejecuta para que quien implemente sepa qué archivo crear y qué step agregar al CI. La config de CI concreta (GitHub Actions, GitLab, etc.) depende del proyecto.
- **Relacionados**: vinculá con `layers-and-dependencies` (fuente de las reglas) y `ci-quality-gates` (donde este check se integra como gate).
