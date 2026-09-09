# EDR — El mine gate del kernel se desactiva por presencia cuando el dominio activo declara su propio mecanismo

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
`development` declaró su propio mecanismo de captura de decisiones in-flow (proposal → ratificar una vez → materializar en apply) en vez del mine gate genérico post-verify del kernel. Borrar esta sección del kernel hubiera roto a `design` silenciosamente, porque `design` sigue dependiendo del mine gate genérico. Hacer el mecanismo opt-in en el otro sentido hubiera exigido editar `payload/domains/design/`, fuera de alcance de cualquier cambio que no sea del propio dominio design.

## Decisión
El kernel no nombra qué dominios tienen mecanismo propio. Antes de evaluar el mine gate, chequea por presencia si el fragmento del dominio activo declara su propio mecanismo de captura de decisiones. Si lo declara, el gate entero no corre para ese dominio — ni trigger check, ni dispatch, ni mención — el mecanismo propio del dominio es lo que rige de punta a punta. Un dominio cuyo fragmento no declara mecanismo propio hereda este gate sin cambios.

## Reglas verificables
- **[manual]** `payload/core/CLAUDE.md`, sección `### Decision-Gap Capture (mine gate)`, evalúa por presencia si el fragmento del dominio activo declara su propio mecanismo, y si lo declara, ninguno de los pasos del gate corre para ese dominio.
- **[manual]** El kernel no enumera por nombre qué dominios tienen mecanismo propio.
- **[manual]** `payload/domains/development/CLAUDE.md`, sección `### In-Flow Decision Capture (MANDATORY)`, es el lado que invoca esa cláusula de override desde el fragmento de `development`: declara su propio mecanismo (proposal · materializar straight-through en apply) en vez del mine gate genérico del kernel, y apunta al mecanismo completo en `in-flow-capture.md` sin restatearlo.

## Alternativas consideradas
Borrar la sección del kernel y mover el mecanismo genérico al fragmento de `design` — descartada: rompe a `design` en el mismo cambio que no lo toca. Hacer el mecanismo de `development` opt-in mediante una bandera en el fragmento de `design` — descartada: edita un fragmento fuera del alcance ratificado de este cambio.

## Consecuencias
El kernel queda con una cláusula de override genérica por presencia, en vez de una lista de excepciones nombradas que hay que mantener sincronizada cada vez que un dominio nuevo declara su propio mecanismo.
