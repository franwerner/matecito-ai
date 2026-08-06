# EDR — Render y validación como ejecutables separados

- **Status:** Accepted
- **Date:** 2026-08-06

## Contexto

Los artefactos durables del ecosistema se producen y se auditan en momentos distintos y por actores distintos. Producir es un acto de escritura que ocurre dentro de una fase, sobre un artefacto por vez, con datos que esa fase acaba de generar. Auditar es una lectura que puede correr sobre un store entero, sin dueño ni fase, y que también corre sobre material que nadie produjo con estas herramientas.

La garantía que más importa del lado de la auditoría es negativa: **el auditor no escribe**. Una garantía negativa que depende de disciplina se rompe en la primera excepción conveniente.

## Decisión

El renderizado y la validación viven en **ejecutables separados**, no como subcomandos de una herramienta única. Cada uno lee el mismo contrato, y ninguno importa al otro.

## Alcance

- `payload/domains/development/scripts/*.js` — los ejecutables del ecosistema; la separación aplica al par render/validate de cada familia de artefactos.
- `payload/domains/development/references/*/templates/*.yaml` — el contrato que ambos leen, y que es lo único compartido entre ellos.

## Reglas verificables

- **[manual]** Ningún ejecutable de validación escribe, mueve ni borra archivos del store que audita.
- **[manual]** El renderizador y el validador de una misma familia de artefactos no se importan entre sí; lo único que comparten es el contrato declarativo.
- **[auto]** Cada validador soporta un modo de auto-chequeo que compara el contrato del repo contra el desplegado y sale sin escribir.

## Alternativas consideradas

- **Un solo ejecutable con subcomandos.** Descartado: comparten proceso, y entonces "esto no escribe" pasa a ser una propiedad de la disciplina de quien edita y no de la estructura. Un validador que puede escribir es un validador en el que hay que confiar en vez de uno que se puede correr sobre cualquier cosa.
- **Un ejecutable que valide y ofrezca corregir.** Descartado por la misma razón, agravada: el auto-fix convierte cada hallazgo en una mutación silenciosa, que es el modo de falla que estas herramientas existen para cerrar.

## Consecuencias

- El validador se puede apuntar a cualquier store, propio o ajeno, sin evaluar el riesgo de que mute algo.
- Cuesta duplicar el parseo del contrato en dos lugares. El auto-chequeo de cada uno acota ese costo: si el contrato cambia y una de las dos lecturas queda vieja, el modo de auto-chequeo lo reporta.
- La invocación es más verbosa: producir y auditar son dos llamadas, nunca una.

## Relacionados

- `relacionado-con` → [contract-pair-in-templates.md](contract-pair-in-templates.md) — el contrato único que ambos ejecutables leen.
- `relacionado-con` → [../contracts/three-level-checker-severity.md](../contracts/three-level-checker-severity.md) — cómo reporta el validador lo que encuentra.
