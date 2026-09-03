# EDR — Un cambio a un script del payload envía tests de regresión en dev-tests/, aunque el spec no los liste

- **Status:** Accepted
- **Date:** 2026-09-03

## Contexto
La tabla de Scope de este cambio — la fuente que `sdd-apply` tiene para saber qué tocar — no lista ningún archivo bajo `dev-tests/`. Leída al pie de la letra, escribir tests de regresión para el cambio a `render-return.js` queda fuera. Pero el cambio toca un script del que dependen las diez fases del pipeline a través de su contrato de retorno, y dejarlo sin regresión es exactamente el tipo de cambio silencioso que más cuesta detectar después. `structure/installer-step-ships-with-tests.md` (Accepted) ya resolvió el mismo fork para un installer step en Go — su alcance es explícitamente ese, así que no cubre un script JS del payload, pero su razonamiento (todo hermano existente tiene cobertura, Strict TDD está apagado y no fuerza la respuesta, un cambio sin test es el que más regresiona en silencio) aplica igual de bien acá.

## Decisión
El cambio a `render-return.js` envía sus tests de regresión en `payload/domains/development/dev-tests/render-return.test.js`, aunque la tabla de Scope del spec no lo liste como archivo a tocar. Es el gemelo de `installer-step-ships-with-tests.md` para scripts JS del payload en vez de pasos de instalación en Go.

## Reglas verificables
- **[auto]** Todo script de `payload/domains/development/scripts/` que gana o cambia una rama de comportamiento observable gana su cobertura correspondiente en `payload/domains/development/dev-tests/`, corrida vía `node --test`.
- **[manual]** La ausencia de un archivo de test en la tabla de Scope del spec de un cambio no exime a ese cambio de escribir tests de regresión para el script que toca — el spec fija ALCANCE de producto, no la lista completa de artefactos de soporte.

## Alternativas consideradas
No escribir tests, apoyado en que la tabla de Scope del spec no los lista: descartada, por el mismo razonamiento que ya cerró `installer-step-ships-with-tests.md` — todo script hermano en este árbol tiene cobertura, Strict TDD apagado no fuerza la respuesta por sí solo, y un cambio de comportamiento compartido por diez fases sin test de regresión es exactamente el caso donde más cuesta notar una rotura.

## Consecuencias
Todo cambio futuro a un script de `payload/domains/development/scripts/` que introduzca o modifique una rama de comportamiento observable se espera que envíe cobertura equivalente en `dev-tests/`, salvo que un cambio posterior ratifique explícitamente omitirla — absorber ese hueco en silencio reabre el mismo fork que este registro cierra.

## Relacionados
- `relacionado-con` → [installer-step-ships-with-tests.md](installer-step-ships-with-tests.md) — el mismo fork, resuelto antes para un installer step en Go — éste es su gemelo para scripts JS del payload
