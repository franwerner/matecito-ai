# Project Conventions for Claude

## Comportamiento del sistema (capability-specs)

El comportamiento de **Matecito UI** —qué hace el sistema ante cada situación— está definido a nivel de proyecto en `.matecito-ai/development-specs/`, **organizado por tipo** (`flow` / `rule` / `lifecycle` / `process`). Vive en el root porque describe el producto completo, no un sub-app puntual.

**Antes de escribir código o tests que implementen o modifiquen un comportamiento (un flujo, una regla de negocio, un ciclo de vida, un proceso):**
1. Abrí `.matecito-ai/development-specs/INDEX.md` (índice raíz) e identificá el **tipo** relevante a tu tarea.
2. Abrí `.matecito-ai/development-specs/<type>/INDEX.md` y leé el capability-spec de lo que vas a tocar — es el contrato del *qué hace*, con sus escenarios verificables.
3. Si hay contradicción entre tu plan y un spec `Accepted`: pará y preguntale al usuario.

Los specs dicen *qué hace* el sistema; el *por qué* de cada elección técnica vive en los EDRs de cada sub-app (ej. `apps/api/.matecito-ai/edr/`) y el *cómo* literal en el código. Para definir, actualizar o validar comportamiento, usá `development-spec-bootstrap` (y `development-spec-validate` para chequear coherencia).

## Decisiones de ingeniería (EDRs)

Las decisiones técnicas viven **por sub-app**, junto al código que gobiernan: `apps/api/.matecito-ai/edr/` (broker/MCP), `apps/ui/.matecito-ai/edr/` (UI). Antes de tocar arquitectura, capas, errores, datos, transporte o convenciones de un sub-app, leé sus EDRs (empezando por su `edr/INDEX.md`) y su `CLAUDE.md`.
