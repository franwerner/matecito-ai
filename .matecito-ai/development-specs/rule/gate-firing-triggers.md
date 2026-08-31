# Capability — Cuándo un gate detiene el flujo

- **Status:** Accepted
- **Date:** 2026-08-31
- **Components:** cli

## Propósito

**Cuándo** un gate dispara sobre los items que una fase retorna, y qué llega al usuario cuando no lo hace. Hoy la decisión pertenece a la sección: cualquier contenido en una sección marcada Tier 1 detiene el flujo. Tras este cambio la decisión pertenece al item, y se reduce a dos triggers. `flow/ratify-gate-items` cubre **cómo** un gate presenta una vez que disparó; `rule/contract-shape-proposal` cubre el contenido del trigger (b). Ninguno cubre cuándo.

## Los cuatro valores de `gates:`

El vocabulario fijado por este spec. Léase como una escalera, más fuerte a más débil — cada valor dice qué una sección hace con un item que declara **ningún** trigger:

| Valor | Un item que declara un trigger | Un item que declara ninguno | Secciones que lo toman |
|---|---|---|---|
| `always` | detiene | n/a — todo item detiene por definición | `### Contract Shapes Proposed` (en cada uno de `sdd-propose`, `sdd-spec`, `sdd-design`, `sdd-apply`) |
| `contested` | detiene | aparece en el resumen entre-fases | `### Scope and approach (unconfirmed)`, `### Derived capabilities (unconfirmed)`, ambas `### New Decisions`, `### Tasks not traceable to spec/design`, `### Unmandated Forks` |
| `reported` | aparece | aparece | `sdd-verify`'s `## Decision Gaps` y `## UI Verdict` — las filas que D.3 marcó `—`, que nunca detuvieron y siempre reportaron |
| `muted` | aparece | **nada llega al usuario** | `### Open Questions`, `### Mandated Departures` |

`contested` y `always` mantienen los significados que la proposal fijó. Lo que la proposal llamaba `never` se divide en `reported` (nunca detiene, siempre reporta) y `muted` (nunca detiene, reporta solo un item disparado). Ninguna sección toma más de un valor, y toda sección en la tabla D.3 toma exactamente uno.

## Requisitos

### Requisito: El criterio de disparo son dos triggers cerrados, evaluados por item

Un gate sobre un retorno de fase DEBE disparar para un item **si y solo si** uno de dos triggers se cumple:

- **(a) contested finding** — el veredicto `contested` del item es algo distinto de `none`, incluyendo ausente u hedgeado;
- **(b) contract shape** — el item es una forma de contrato propuesta.

Un item que no cumple ninguno DEBE proceder sin turno de usuario. La decisión de disparo DEBE NO ser tomada de la sección donde el item aparece. Una sección cuyo valor `gates:` es `reported` o `muted` nunca dispara, sea lo que declare sus items.

#### Scenario: Un retorno limpio se despacha silenciosamente

- **GIVEN** un retorno de fase cuyo cada item ratificable declara `contested: none`, y que no lleva `### Contract Shapes Proposed`
- **WHEN** el orquestador aplica el Guard de Decisiones sin Resolver
- **THEN** no abre gate, no corre walkthrough, y la siguiente fase se despacha
- **AND** el mecanismo no se menciona al usuario en absoluto

#### Scenario: Un item contested entre muchos abre el gate solo para ese

- **GIVEN** un retorno con cuatro items ratificables, de los cuales uno declara `contested: contradicts-record`
- **WHEN** el guard corre
- **THEN** el índice del gate cuenta uno y el walkthrough presenta ese solo
- **AND** los otros tres no toman resultado

#### Scenario: Un contrato dispara sin importar ningún otro item

- **GIVEN** un retorno llevando un item `### Contract Shapes Proposed` y ningún veredicto `contested` distinto de `none` en otro lado
- **WHEN** el guard corre
- **THEN** el gate abre para el item de contrato
- **AND** esa sección declara `gates: always` en su `.yaml` de fase y no lleva token `contested`

#### Scenario: Un item que no dispara no se borra silenciosamente

- **GIVEN** un item en una sección cuyo valor es `contested`, declarando `contested: none`
- **WHEN** se presenta el resumen entre-fases de la fase
- **THEN** el `summary` del item aparece ahí con su anchor
- **AND** nada aguarda en él

### Requisito: El token `contested` y su conjunto cerrado de valores

Toda sección que dispara sobre el trigger (a) DEBE declarar un token por-item llamado `contested`, sobre el field `contested`, con exactamente estos valores:

`none` · `contradicts-statement` · `contradicts-record` · `unverified-assumption`

La declaración DEBE NO llevar una lista `passing`. Cada valor distinto de `none` DEBE nombrar una contraparte concreta — una declaración explícita del usuario u orquestador, un record Aceptado, una asunción que la fase misma no pudo verificar — y el conjunto DEBE NO contener valor alguno expresando incertidumbre. Una sección cuyo valor es `muted` DEBE declarar el mismo token, porque el surfacing de sus items depende de él.

#### Scenario: Una contradicción honesta es renderizable

- **GIVEN** una fase escribiendo un item que declara `contested: contradicts-record`
- **WHEN** renderiza su retorno
- **THEN** el render sucede y el bloque lleva la línea `· contested: contradicts-record`

#### Scenario: Declarar `passing` haría un veredicto honesto no-renderizable

- **GIVEN** la declaración del token `contested` en cualquiera de los cinco archivos `.yaml` de fase
- **WHEN** se lee
- **THEN** declara `values` y no `passing`
- **AND** la regla del renderer que un valor distinto de `passing` falla el render por tanto nunca alcanza este token

#### Scenario: Un valor fuera del conjunto se rechaza mecánicamente

- **GIVEN** un item declarando `contested: maybe`
- **WHEN** se valida el retorno
- **THEN** la validación reporta `TOKEN-ILLEGAL`, nombrando el item y los cuatro valores legales

#### Scenario: No hay escondite para la incertidumbre

- **GIVEN** los cuatro valores declarados
- **WHEN** una fase busca un valor que signifique "inseguro"
- **THEN** ninguno existe, y la fase debe elegir entre `none` y una contraparte nombrada

### Requisito: Ausencia y hedging disparan — solo un `none` afirmativo las salta

Un veredicto `contested` ausente u hedgeado DEBE leerse como disparando. El silencio DEBE NO saltarse un gate, y en una sección `muted` DEBE NO saltarse el surfacing tampoco.

#### Scenario: Un token faltante se atrapa antes del guard

- **GIVEN** un item ratificable sin línea `contested:`
- **WHEN** se valida el retorno
- **THEN** se reporta `TOKEN-MISSING` para ese item

#### Scenario: Un token faltante que alcanza el guard sigue disparando

- **GIVEN** un item que alcanza el orquestador sin veredicto `contested` legible
- **WHEN** el guard lo clasifica
- **THEN** dispara, bajo la lectura estricta
- **AND** el orquestador no suministra un veredicto en nombre de la fase

### Requisito: El orquestador clasifica en el token solo

El orquestador DEBE clasificar cada item en su token `contested` y DEBE NO re-derivar el veredicto leyendo la prosa del item. La clasificación DEBE estar enunciada como una tabla en el Guard de Decisiones sin Resolver, en la misma forma que las tablas `blocking-test` y `design-conflict` existentes.

#### Scenario: Prosa que se lee como una contradicción no anula el token

- **GIVEN** un item cuyo `summary` se lee como contradictorio a un record Aceptado, pero que declara `contested: none`
- **WHEN** el guard corre
- **THEN** el item no dispara
- **AND** el orquestador no re-lee la prosa para cuestionarse el veredicto

#### Scenario: La tabla coincide con la forma de las otras dos

- **GIVEN** la nueva tabla "Reading the `contested` token"
- **WHEN** se compara con las tablas `blocking-test` y `design-conflict` en el mismo archivo
- **THEN** tiene las mismas tres columnas y la misma fila ausente-o-hedgeada

### Requisito: El costo de la auto-reporte se enuncia llano y nunca presentado como mitigado

El texto de gobernanza DEBE enunciar que `unverified-assumption` es auto-reportado y que ningún script puede probar el valor cierto; que esta es una garantía **estrictamente más débil** que la regla de identidad de sección que reemplaza, porque una fase puede ahora optarse a sí misma de un gate escribiendo una palabra; y DEBE NO presentar el strict default o el resumen entre-fases como detección o mitigación. En una sección `muted` el costo se compone — un item limpio erróneamente no es meramente desbloqueante, es invisible — y eso DEBE enunciarse donde el valor se define. Los dos scripts DEBEN mantenerse intactos.

#### Scenario: El costo se escribe, no se implica

- **GIVEN** el Guard de Decisiones sin Resolver reescrito
- **WHEN** se lee
- **THEN** nombra la falla — una fase aseverando `none` sobre una asunción real sin verificar — y enuncia que nada en el ecosistema lo detecta
- **AND** ni el strict default ni el transcript se llama mitigación

#### Scenario: El costo compuesto de `muted` se nombra

- **GIVEN** la definición del valor `muted`
- **WHEN** se lee
- **THEN** enuncia que un item erróneamente limpio en tal sección llega al usuario en ningún lado
- **AND** no ofrece el transcript como fallback, porque no hay entrada de transcript a la que caer

#### Scenario: Ningún script aprende a hacer valer el veredicto

- **GIVEN** el diff del cambio
- **WHEN** se inspecciona `validate-return.js` y `render-return.js`
- **THEN** ninguno aparece en él
- **AND** un grep por `tier` o `gates` en cualquiera de los dos scripts no retorna nada

### Requisito: `### Unmandated Forks` dispara por contrato, no por auto-reporte

Esa sección DEBE declarar `contested` con un **único** valor legal, `unverified-assumption`. Cualquier otro valor DEBE fallar el render. Ninguna regla de orquestador se agrega para ella: el conjunto cerrado hace el trabajo, porque el token `mandate:` de la sección es `chosen` por construcción — más de una resolución era válida, nada la respaldó, y nada fue aplicado.

#### Scenario: Un valor distinto de disparo es imposible de redactar

- **GIVEN** `sdd-apply` escribiendo un item `### Unmandated Forks` con `contested: none`
- **WHEN** renderiza su retorno
- **THEN** el render falla nombrando el field y el valor único legal, y nada llega a stdout

#### Scenario: Cada fork sigue deteniendo el flujo

- **GIVEN** un retorno llevando un item `### Unmandated Forks` y nada más contested
- **WHEN** el guard corre
- **THEN** el gate abre para ese item
- **AND** la debilidad de auto-reporte del token `contested` no alcanza esta sección

### Requisito: Una fuente de verdad de qué dispara; todo copia la cita

Sección D.3 de `_shared/sdd-phase-common.md` DEBE ser la declaración única de cuál sección toma cuál valor. En los cinco contratos `.yaml` de fase el field `tier: 1|2` DEBE **renombrarse** a `gates: always | contested | reported | muted` — renombrado, no redefinido en su lugar, así una copia stale se lee como stale en lugar de como autoridad. La tabla de D.3 DEBE asignar uno de los cuatro a **toda** fila, incluyendo las dos filas de `sdd-verify` que previamente marcó `—`. Toda otra copia de la clasificación DEBE citar D.3 y el token en lugar de reenstatar la regla.

#### Scenario: El field se renombra en todos los cinco contratos

- **GIVEN** cada uno de los cinco archivos `.yaml` de fase
- **WHEN** se lee tras el cambio
- **THEN** ninguna clave `tier:` subsiste, y toda sección que llevaba una declara exactamente uno de los cuatro valores

#### Scenario: Los cuatro valores mapean sin ambigüedad a secciones

- **GIVEN** la tabla D.3 tras el cambio
- **WHEN** se lee cada fila
- **THEN** lleva exactamente uno de `always`, `contested`, `reported`, `muted`
- **AND** ninguna fila queda sin valor, incluyendo las dos filas de `sdd-verify` que previamente llevaban `—`

#### Scenario: El vocabulario de tier se fue del payload

- **GIVEN** una búsqueda por el vocabulario de tier (`Tier 1`, `Tier-1`, `Tier 2`, `Tier-2`, `tier: 1`, `tier: 2`) bajo `payload/`
- **WHEN** se corre tras el cambio
- **THEN** solo coincide `payload/domains/design/references/design-principles/README.md` y `payload/domains/development/references/design-patterns/README.md`

#### Scenario: Los mirrors apuntan en lugar de reenstatar

- **GIVEN** cada uno de los 23 archivos in-scope que llevaban el vocabulario de tier
- **WHEN** se lee tras el cambio
- **THEN** ninguno enuncia la regla de disparo en sus propias palabras; cada uno nombra el token y apunta a D.3

#### Scenario: El field se queda inerte

- **GIVEN** el field renombrado `gates:`
- **WHEN** se buscan en los scripts
- **THEN** ningún script lo lee, exactamente como ninguno leía `tier`

### Requisito: "Ratificado" significa que el gate no lo controló

Una proposal DEBE contar como ratificada cuando sea el usuario la confirmó en un gate, o el gate nunca disparó para ella. Un item declarando `contested: none` DEBE ser reenviado a `sdd-apply` marcado **ratificado**, verbatim como fue redactado, sin turno de usuario. El prompt de despacho DEBE ser byte-idéntico en ambos caminos, así el paso de materialización de `sdd-apply` y su regla "una resolución faltante retorna `blocked`" no cambian.

#### Scenario: Una proposal no contestada se materializa sin turno de usuario

- **GIVEN** un item `### New Decisions` declarando `contested: none`
- **WHEN** el guard corre
- **THEN** no abre gate, y el item se reenvía a `sdd-apply` marcado ratificado, verbatim
- **AND** `sdd-apply` lo materializa como un record `Accepted` en el mismo paso que implementa el trabajo que lo rige

#### Scenario: Los dos caminos son indistinguibles downstream

- **GIVEN** un item ratificado-por-usuario y uno ratificado-automáticamente
- **WHEN** se comparan sus prompts de despacho
- **THEN** la forma de reenvío es idéntica, y nada le dice a `sdd-apply` cuál camino la produjo

#### Scenario: Una resolución faltante sigue bloqueando

- **GIVEN** `sdd-apply` alcanzando una tarea gobernada por un item su prompt de despacho no menciona
- **WHEN** evalúa el item
- **THEN** retorna `blocked` nombrando el item y la resolución faltante, sin cambio por este cambio

#### Scenario: Lo que `Accepted` ahora garantiza es más angosto

- **GIVEN** una decisión auto-ratificada escrita a disco como `Accepted`
- **WHEN** un lector downstream la consulta
- **THEN** fue vista por el usuario solo como una línea en un resumen entre-fases
- **AND** todo lector que trata `Accepted` como intención ratificada sigue haciéndolo

### Requisito: El ledger de ratificación registra la ruta auto

El ledger `sdd/{change-name}/ratified-decisions` DEBE ganar una fila por item auto-ratificado, llevando `gate: auto` en lugar de un nombre de gate. El matching de re-emergence (string exacto en `record`) DEBE funcionar sobre esas filas.

#### Scenario: Un item auto-ratificado se registra

- **GIVEN** un item auto-ratificado porque su veredicto fue `none`
- **WHEN** el paso del gate cierra
- **THEN** existe una fila del ledger para él con `gate: auto` y su texto ratificado

#### Scenario: Un item re-emergente que sigue limpio se silencia

- **GIVEN** un item posterior coincidiendo una fila del ledger cuyo `gate` es `auto`, y cuyo propio veredicto es `none`
- **WHEN** el paso del gate posterior corre
- **THEN** se auto-ratifica de nuevo, silenciosamente, sin pregunta hecha

#### Scenario: Un item re-emergente que ahora es contested camina la forma completa

- **GIVEN** un item posterior coincidiendo esa misma fila `auto`, declarando un veredicto distinto de `none`
- **WHEN** el gate abre
- **THEN** se camina en forma completa, no a través de la forma corta de re-emergence
- **AND** la razón es que la forma corta pregunta si una decisión de usuario anterior sigue aplicando, y no hubo una

### Requisito: La ruta de rechazo mantiene su forma exacta

`### Rejected Proposals Checked`, el token `design-conflict: none | conflicts`, la regla que `conflicts` fuerza `status: blocked`, y el conjunto esperado sostenido por el orquestador de rechazos reenviados DEBEN ser sin cambios. Los rechazos se hacen más raros — uno puede ahora originarse solo en un gate que disparó — pero el chequeo sobrevive intacto.

#### Scenario: Un rechazo sigue debiendo un veredicto

- **GIVEN** una proposal rechazada reenviada en un prompt de despacho cuya tarea gobernada el run alcanzó
- **WHEN** `sdd-apply` retorna
- **THEN** el retorno lleva un item para ese `record:` con un veredicto `design-conflict`
- **AND** un rechazo reenviado sin item coincidente es un retorno incorrecto

#### Scenario: `conflicts` sigue forzando un stop

- **GIVEN** un item declarando `design-conflict: conflicts`
- **WHEN** el orquestador lo lee
- **THEN** el retorno debe ser `blocked`, y ambas versiones se presentan

#### Scenario: Nada fue reenviado como rechazado

- **GIVEN** un cambio donde ninguna proposal alcanzó un gate
- **WHEN** `sdd-apply` retorna
- **THEN** la sección es legítimamente ausente, per la regla de emisión condicional sin cambios

### Requisito: La escalación `blocking-test` no se toca

Los valores de `blocking-test`, su tabla de clasificación y su escalada nombrada por eje DEBEN permanecer byte-idénticos. Un eje nombrado es una escalada **pasado** el walkthrough, no un caso de él, así que estrechar cuáles items entran al walkthrough no lo alcanza.

#### Scenario: Un eje nombrado sigue escalando sobre un veredicto limpio

- **GIVEN** un item `### New Decisions` declarando `blocking-test: data-model` y `contested: none`
- **WHEN** el guard lo clasifica
- **THEN** se surfacea como un stop de estilo `blocked`, citando el token
- **AND** el veredicto `contested` no lo degrada

### Requisito: Solo uno de los nueve momentos de cita cambia, y nueve no es el inventario completo

De los nueve momentos que citan el walkthrough compartido, solo el gate de decisiones pendientes DEBE cambiar. Los otros ocho DEBEN ser sin cambios: el INTAKE GATE, el gate de decisión-minada (que nunca dispara en este dominio), el gate de spec-minada, el Discovery Gate, el Uncommitted-Work Gate, el Review Workload Guard, un retorno `blocked` de fase, y `risks` / Hallazgos del Validador. El texto de gobernanza DEBE NO reclamar que los nueve son el inventario completo de puntos de interrupción.

#### Scenario: Los ocho quedan solos

- **GIVEN** el diff del cambio
- **WHEN** se inspeccionan las definiciones de los otros ocho momentos
- **THEN** ninguno de sus disparadores cambió
- **AND** ninguno de ellos adquirió un token `contested`

#### Scenario: Dos puntos de interrupción fuera de los nueve se nombran, no se absorben

- **GIVEN** el reclamo calificado en `gate-presentation.md`
- **WHEN** se lee
- **THEN** nombra el reporte de merge-conflict del Change Workspace y el STOP de atomicidad de commit como puntos de interrupción fuera de los nueve
- **AND** enuncia que ninguno de los triggers lo alcanza, y ninguno se modifica

## Escenarios

El spec no declara escenarios adicionales más allá de los definidos en sus doce requisitos anteriormente.

## Referencias

- **Conceptualmente relacionado**: `flow/ratify-gate-items.md` — define cómo un gate presenta una vez que disparó
- **Conceptualmente relacionado**: `rule/contract-shape-proposal.md` — define el contenido del trigger (b)
- **Conceptualmente relacionado**: `rule/decision-re-emergence-reconfirmation.md` — describe el comportamiento de re-emergence adyacente; tras archive ambas coexisten sin fusión
