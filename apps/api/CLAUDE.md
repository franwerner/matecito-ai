# Project Conventions for Claude

Las decisiones de ingeniería de este broker/MCP (arquitectura, convenciones y políticas) están en `.matecito-ai/edr/`, **organizadas por dominio**.

**Antes de escribir código que toque arquitectura, capas, errores, datos, transporte o convenciones:**
1. Abrí `.matecito-ai/edr/INDEX.md` (índice raíz) e identificá el **dominio** relevante a tu tarea.
2. Abrí `.matecito-ai/edr/<dominio>/INDEX.md` y leé los EDRs de ese dominio antes de escribir código.
3. Si hay contradicción entre tu plan y un EDR: pará y preguntale al usuario.

**Antes de instalar/sugerir cualquier dependencia nueva (lib, framework, herramienta), leé `.matecito-ai/edr/tech/INDEX.md`** para ver qué tecnologías ya están elegidas. Si tu sugerencia pisa con algo ya registrado, no la introduzcas sin preguntar.

**Cuando un EDR declara `Applied pattern: X`,** la definición canónica del patrón está en `~/.claude/references/design-patterns/patterns/<x>.md`. Consultá ese archivo antes de implementar para entender el contrato del patrón. Si vas a desviarte de la definición canónica, justificalo en el EDR — no improvises una variante.

Si una decisión no está documentada o algo no queda claro, **preguntá al usuario antes de inventar una convención**. Las decisiones se registran como EDR, no se improvisan.

Para crear, actualizar o revisar decisiones de ingeniería (incluyendo agregar/cambiar tecnologías del catálogo), usá la skill `development-decisions-bootstrap`.

## Comportamiento del sistema (capability-specs)

El comportamiento del sistema vive a **nivel de proyecto** (root), no en este sub-app: `<repo-root>/.matecito-ai/development-specs/`, organizado por tipo (`flow` / `rule` / `lifecycle` / `process`). El broker implementa varias de esas capabilities; su contrato del *qué hace* se lee desde ahí.

**Antes de implementar o modificar un comportamiento** (un flujo, una regla, un ciclo de vida, un proceso): leé `<repo-root>/.matecito-ai/development-specs/INDEX.md` y el spec correspondiente. Si tu plan contradice un spec `Accepted`, pará y preguntale al usuario. Los specs dicen *qué hace*; el *por qué* técnico vive en los EDRs de este sub-app (`.matecito-ai/edr/`) y el *cómo* en el código. Para definir/validar comportamiento: `development-spec-bootstrap` / `development-spec-validate`.
