# EDR — Dos sitios fuera del Scope cerrado se barren igual, por veredicto propio

- **Status:** Accepted
- **Date:** 2026-09-08

## Contexto
El cambio `sdd/drop-apps-subtree` declara un `## Scope` cerrado — "Nada fuera de esta lista está en alcance". Leyendo el store se encontraron dos sitios que ese Scope no nombra pero que el cambio igual falsifica: `.matecito-ai/.markdown_vault_mcp/state.json`, un índice versionado en git cuyo mapa `indexed` trae 59 entradas file→hash, 14 de las cuales nombran specs que este cambio borra; y `REFACTOR.md:233` (§6, ítem 5), cuya frase "arrancando por el piloto" apunta a un piloto que la reescritura de §2 (tarea 1.1) ya no propone. Dos EDRs Accepted tiran en direcciones opuestas sobre qué hacer con un sitio así: `structure/retired-vocabulary-in-record-stores.md` deja un sustantivo incidental dentro de una regla que sigue vigente, y `structure/gating-vocabulary-sweep-scope.md` barre un puntero que ya no resuelve. Sólo una lectura sitio-por-sitio resuelve los dos sin aplicar ninguno a ciegas — la misma forma que `structure/lane-vocabulary-sweep-in-record-stores.md` ya ratificó antes: un registro, veredictos IN/OUT por sitio, cada uno con su razón.

## Decisión
Ambos sitios entran al barrido, cada uno con su propio veredicto: **Sitio 1 — `state.json` (IN)** — es un índice, no una regla con un sustantivo incidental; sus 14 entradas retiradas son punteros a archivos que van a dejar de existir, exactamente lo que `gating-vocabulary-sweep-scope.md` corrige. El mecanismo fijado por el usuario en el gate de diseño es editar el archivo directamente (borrar las entradas que nombran specs retirados), nunca re-correr la tool markdown-vault — así el trabajo no depende de que esa tool esté disponible donde corre `sdd-apply`. **Sitio 2 — `REFACTOR.md:233` (IN)** — "arrancando por el piloto" no es un sustantivo incidental dentro de una regla vigente: es un puntero a un mecanismo (el piloto de qmd) que la propia reescritura de §2 retira en esta misma tanda, dejando la frase apuntando a algo que no existe. Ninguno de los dos es el residuo que `retired-vocabulary-in-record-stores.md` protege — ese EDR protege un sustantivo que sobrevive como ejemplo dentro de una regla intacta, no un puntero que ya no resuelve.

## Reglas verificables
- **[manual]** `.matecito-ai/.markdown_vault_mcp/state.json`'s `indexed` map no contiene ninguna entrada que nombre uno de los 14 capability-specs retirados por `sdd/drop-apps-subtree`.
- **[manual]** El pruning de `state.json` se hizo por edición directa del archivo, nunca re-corriendo la tool markdown-vault.
- **[manual]** `REFACTOR.md:233` (§6, ítem 5) no usa la frase "arrancando por el piloto" ni ninguna otra que apunte a un piloto que §2 ya no propone.
- **[manual]** Ningún otro sitio fuera del `## Scope` cerrado del delta se toca a partir de este EDR — el barrido queda acotado a estos dos, no se generaliza a un tercer sitio no ratificado.

## Alternativas consideradas
(a) Tratar el `## Scope` cerrado como exhaustivo y dejar los dos sitios como residuo aceptado — descartada porque uno es un índice que va a apuntar a archivos inexistentes (no un sustantivo incidental) y el otro es un puntero roto, ambos el caso que `gating-vocabulary-sweep-scope.md` corrige, no el que `retired-vocabulary-in-record-stores.md` protege. (b) Extender el barrido a todo sitio que el cambio falsifique, sin veredicto por sitio — descartada porque generaliza sin la lectura caso-por-caso que dos EDRs en tensión requieren, y es exactamente lo que `structure/lane-vocabulary-sweep-in-record-stores.md` ya rechazó como forma.

## Consecuencias
El repo no queda con un índice tracked apuntando a 14 archivos borrados ni con una frase de planificación citando un piloto retirado. El precedente de un registro único con veredictos IN/OUT por sitio, cada uno con su razón, queda reutilizable para el próximo barrido que enfrente la misma tensión entre los dos EDRs de vocabulario retirado.
