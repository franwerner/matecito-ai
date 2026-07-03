<!-- Canonical template: SECCIÓN de `CLAUDE.md` raíz del proyecto que apunta al índice de capability-specs. La agrega development-spec-bootstrap en la materialización, de forma IDEMPOTENTE: si el CLAUDE.md no existe, se crea con esta sección; si existe, se agrega/actualiza SOLO esta sección (ancla: el heading de abajo) sin tocar el resto (p. ej. la sección de ADRs). NO duplicar. -->

## Comportamiento del sistema (capability-specs)

El comportamiento de este proyecto —qué hace ante cada situación— está en `.matecito-ai/development-specs/`, **organizado por tipo** (`flow` / `rule` / `lifecycle` / `process`).

**Antes de escribir código o tests que implementen o modifiquen un comportamiento (un flujo, una regla de negocio, un ciclo de vida, un proceso):**
1. Abrí `.matecito-ai/development-specs/INDEX.md` (índice raíz) e identificá el **tipo** relevante a tu tarea.
2. Abrí `.matecito-ai/development-specs/<type>/INDEX.md` y leé el capability-spec de lo que vas a tocar — es el contrato del *qué hace*, con sus escenarios verificables.
3. Si hay contradicción entre tu plan y un spec `Accepted`: pará y preguntale al usuario.

Los specs dicen *qué hace* el sistema; el *por qué* de cada elección técnica vive en `.matecito-ai/adr/` y el *cómo* literal en el código. Para definir, actualizar o validar comportamiento, usá `development-spec-bootstrap` (y `development-spec-validate` para chequear coherencia).
