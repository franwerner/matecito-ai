<!-- Canonical template: `CLAUDE.md` raíz del proyecto objetivo (mínimo, apunta al índice raíz de ADRs). Consumido por la fase de Materialización de SKILL.md. NO sobrescribir un CLAUDE.md existente sin permiso. -->

# Project Conventions for Claude

Las decisiones de ingeniería de este proyecto (arquitectura, convenciones y políticas) están en `.matecito-ai/adr/`, **organizadas por dominio**.

**Antes de escribir código que toque arquitectura, capas, errores, auth, datos o convenciones:**
1. Abrí `.matecito-ai/adr/INDEX.md` (índice raíz) e identificá el **dominio** relevante a tu tarea.
2. Abrí `.matecito-ai/adr/<dominio>/INDEX.md` y leé los ADRs de ese dominio antes de escribir código.
3. Si hay contradicción entre tu plan y un ADR: pará y preguntale al usuario.

**Antes de instalar/sugerir cualquier dependencia nueva (lib, framework, herramienta, DB), leé `.matecito-ai/adr/tech/INDEX.md`** para ver qué tecnologías ya están elegidas. Si tu sugerencia pisa con algo ya registrado, no la introduzcas sin preguntar.

**Cuando un ADR declara `Applied pattern: X`,** la definición canónica del patrón está en `~/.claude/references/design-patterns/patterns/<x>.md`. Consultá ese archivo antes de implementar para entender el contrato del patrón. Si vas a desviarte de la definición canónica, justificalo en el ADR — no improvises una variante.

Si una decisión no está documentada o algo no queda claro, **preguntá al usuario antes de inventar una convención**. Las decisiones se registran como ADR, no se improvisan.

Para crear, actualizar o revisar decisiones de ingeniería (incluyendo agregar/cambiar tecnologías del catálogo), usá la skill `development-decisions-bootstrap`.
