---
name: context
depth: deep
domain: context
type: decision
source: práctica clásica de inception / arc42 §1-3 (contexto y alcance)
---

# Fase: Contexto del proyecto

## Qué decide

El tipo de sistema, stack principal, tamaño de equipo y punto de partida. Es la entrada de todas las fases siguientes: sin este contexto, los defaults de cada fase no se pueden calcular.

## Preguntas

Una por turno. Para cada una: línea de "por qué importa", opciones con default, y siempre la opción "no sé, recomendame" donde aplica.

### 1. Tipo de proyecto

> Define qué fases son requeridas y cuáles se pueden saltear o simplificar.

- API REST
- API GraphQL
- CLI
- Librería
- App web SPA
- App web SSR
- Microservicio
- Monolito modular
- Script / automatización
- Otro (describí brevemente)

### 2. Stack principal

> Permite inferir defaults de lenguaje, idioma de errores, herramientas naturales del ecosistema.

Si el pre-flight detectó lenguaje y framework, mostrá lo detectado y pedí confirmación. Sino, preguntá:

- Lenguaje (y versión si la conocés)
- Framework principal (si aplica)
- No sé, recomendame.

### 3. Tamaño del equipo esperado

> El nivel de formalidad en convenciones escala con el equipo; no tiene sentido sobrediseñar para un proyecto de una sola persona.

- Solo
- *2-3 personas — default razonable para mayoría de proyectos nuevos.*
- 4-10 personas
- 10+ personas

### 4. Greenfield o sobre código existente

> Si hay código existente, los defaults pueden quedar obsoletos y hay que leer el estado actual antes de proponer convenciones.

- Greenfield (repo vacío o sin código propio relevante).
- *Sobre código existente — default cuando el repo tiene archivos propios.*
- Migración (el código existe y se está reescribiendo).

## Notas de lógica (para el motor)

- **Atajo para scripts:** si en (1) dijo "script / automatización" y en (3) dijo "solo", proponé saltar las Fases 1-3 (architecture-style, layers-and-dependencies, inter-layer-communication) e ir directo a folder-structure + subset mínimo de la Fase 5. Pedí permiso explícito para el atajo.
- **Stack detectado en pre-flight:** si el manifest (package.json, pyproject.toml, go.mod, Cargo.toml, etc.) ya reveló el stack, mostralo como default en la pregunta 2 en lugar de preguntar desde cero. Cambia la pregunta a confirmación.
- **Tech a registrar en esta fase:** lenguaje y framework principal (ej: `python.md`, `fastapi.md`, `node.md`, `nestjs.md`). Pedí versión (intentá detectarla del manifest) y "por qué" en 1-2 líneas.

## Tech a registrar

Lenguaje principal y framework web/CLI si aplica. Son las primeras entradas de `tech/INDEX.md`.

## Qué materializar

EDR `context` materializado según el template `~/.claude/references/edr/templates/edr.md`. Esta fase es mayormente descriptiva (caracteriza el proyecto, no impone una restricción chequeable), así que el cuerpo se concentra en `Contexto` y `Decisión`; **no se inventan `Reglas verificables`** salvo que alguna respuesta derive en una restricción real (en cuyo caso se marca con su mecanismo `[tool: ...]` o `[manual]`).

- **Contexto / Decisión** deben capturar los campos descriptivos: tipo de proyecto, lenguaje, framework, versión detectada o declarada, tamaño de equipo, y greenfield vs. código existente vs. migración.
- Si se usó el **atajo de script** (script + solo), documentar en `Decisión` qué fases se saltaron (architecture-style, layers-and-dependencies, inter-layer-communication) y por qué.
- `Relacionados`: este EDR es la entrada de las demás fases; si querés, podés enlazar con `relacionado-con` los EDRs que derivan sus defaults de este contexto. No es obligatorio.
