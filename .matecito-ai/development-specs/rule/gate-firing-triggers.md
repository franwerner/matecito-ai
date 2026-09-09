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
| `contested` | detiene | aparece en el resumen entre-fases | `### Scope and approach (unconfirmed)`, `### Derived capabilities (unconfirmed)`, `### Tasks not traceable to spec/design`, `### Unmandated Forks` |
| `reported` | aparece | aparece | `sdd-verify`'s `## Decision Gaps` y `## UI Verdict` — las filas que D.3 marcó `—`, que nunca detuvieron y siempre reportaron —, ambas `### New Decisions` |
| `muted` | aparece | **nada llega al usuario** | `### Open Questions`, `### Mandated Departures` |

`contested` y `always` mantienen los significados que la proposal fijó. Lo que la proposal llamaba `never` se divide en `reported` (nunca detiene, siempre reporta) y `muted` (nunca detiene, reporta solo un item disparado). Ninguna sección toma más de un valor, y toda sección en la tabla D.3 toma exactamente uno.

## Requisitos

### Requisito: El criterio de disparo son dos triggers cerrados, evaluados por item

Un gate sobre un retorno de fase DEBE disparar para un item **si y solo si** uno de dos triggers se cumple:

- **(a) contested finding** — el veredicto `contested` del item es algo distinto de `none`, incluyendo ausente u hedgeado;
- **(b) contract shape** — el item es una forma de contrato propuesta.

Un item que no cumple ninguno DEBE proceder sin turno de usuario. La decisión de disparo DEBE NO ser tomada de la sección donde el item aparece. Una sección cuyo valor `gates:` es `reported` o `muted` nunca dispara, sea lo que declare sus items.

Este requisito fija **qué** aparece en el resumen entre-fases cuando un item no dispara —el `summary` del item con su anchor— y DEBE NO fijar la **forma** de ese reporte. La forma es de la regla de formas cerradas de la prosa del hilo, que la declara una sola vez; enunciarla también acá deja dos enunciados parciales de la misma cosa, y ninguno de los dos es el dueño.

#### Scenario: Un retorno limpio se despacha silenciosamente

- **verification:** `deferred → tighten-orchestrator-prose`
- **GIVEN** un retorno de fase cuyo cada item ratificable declara `contested: none`, y que no lleva `### Contract Shapes Proposed`
- **WHEN** el orquestador aplica el Guard de Decisiones sin Resolver
- **THEN** no abre gate, no corre walkthrough, y la siguiente fase se despacha
- **AND** el mecanismo no se menciona al usuario en absoluto

#### Scenario: Un item contested entre muchos abre el gate solo para ese

- **verification:** `deferred → tighten-orchestrator-prose`
- **GIVEN** un retorno con cuatro items ratificables, de los cuales uno declara `contested: contradicts-record`
- **WHEN** el guard corre
- **THEN** el índice del gate cuenta uno y el walkthrough presenta ese solo
- **AND** los otros tres no toman resultado

#### Scenario: Un contrato dispara sin importar ningún otro item

- **verification:** `deferred → tighten-orchestrator-prose`
- **GIVEN** un retorno llevando un item `### Contract Shapes Proposed` y ningún veredicto `contested` distinto de `none` en otro lado
- **WHEN** el guard corre
- **THEN** el gate abre para el item de contrato
- **AND** esa sección declara `gates: always` en su `.yaml` de fase y no lleva token `contested`

#### Scenario: Un item que no dispara no se borra silenciosamente

- **verification:** `deferred → tighten-orchestrator-prose`
- **GIVEN** un item en una sección cuyo valor es `contested`, declarando `contested: none`
- **WHEN** se presenta el resumen entre-fases de la fase
- **THEN** el `summary` del item aparece ahí con su anchor
- **AND** nada aguarda en él

#### Scenario: La forma del reporte no se enuncia acá

- **verification:** `deferred → tighten-orchestrator-prose`
- **GIVEN** el texto de este requisito después del cambio
- **WHEN** se lo lee buscando cómo se ve el resumen entre-fases
- **THEN** dice qué del item aparece ahí y remite la forma a la regla de formas cerradas de la prosa del hilo
- **AND** no enuncia por su cuenta la forma del reporte

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

El orquestador DEBE clasificar cada item en su token `contested` y DEBE NO re-derivar el veredicto leyendo la prosa del item. La clasificación DEBE estar enunciada como una tabla en el Guard de Decisiones sin Resolver, en la misma forma que la tabla `blocking-test` existente.

#### Scenario: Prosa que se lee como una contradicción no anula el token

- **GIVEN** un item cuyo `summary` se lee como contradictorio a un record Aceptado, pero que declara `contested: none`
- **WHEN** el guard corre
- **THEN** el item no dispara
- **AND** el orquestador no re-lee la prosa para cuestionarse el veredicto

#### Scenario: La tabla coincide con la forma del único par sobreviviente

- **GIVEN** la nueva tabla "Reading the `contested` token"
- **WHEN** se compara con la tabla `blocking-test` en el mismo archivo
- **THEN** tiene las mismas tres columnas y la misma fila ausente-o-hedgeada
- **AND** ninguna comparación nombra una tabla `design-conflict`, porque ninguna existe

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

### Requisito: La escalación `blocking-test` no se toca

Los valores de `blocking-test`, su tabla de clasificación y su escalada nombrada por eje DEBEN permanecer byte-idénticos. Un eje nombrado es una escalada **pasado** el walkthrough, no un caso de él, así que estrechar cuáles items entran al walkthrough no lo alcanza.

### Requisito: Los momentos que citan el walkthrough compartido son ocho

La aritmética, explícita, para que nadie la componga sobre la línea equivocada: **antes de este cambio eran nueve; este cambio los lleva a ocho**.

Los ocho son: el gate de decisiones pendientes · el gate de confirmación del brief · el Discovery Gate · el gate de decisión-minada · el Uncommitted-Work Gate · el Review Workload Guard · un retorno `blocked` de fase · `risks` / Hallazgos del Validador (una fila, dos fuentes de anchor). El que este cambio remueve es el gate de spec-minada, que deja de existir junto con su ejecutor; el subgrupo que lo pierde es el de gates ratificantes. Toda declaración del conteo en el repo DEBE leer **ocho**, y donde una declaración parte el total en subgrupos, el subgrupo que pierde miembro por este cambio es el de gates ratificantes, no el de momentos de orquestador. Ninguno de los ocho sobrevivientes cambia por este cambio: no cambian sus disparadores y ninguno adquiere un token `contested`. El texto de gobernanza DEBE NO reclamar que los ocho son el inventario completo de puntos de interrupción.

#### Scenario: Toda declaración del conteo lee ocho

- **GIVEN** el repo después del cambio
- **WHEN** se recolecta cada declaración de cuántos momentos citan el walkthrough compartido
- **THEN** cada una lee ocho
- **AND** ninguna declaración lee nueve

#### Scenario: El momento removido es el gate de spec-minada y solo él

- **GIVEN** la enumeración después del cambio
- **WHEN** se compara item por item contra los nueve previos al cambio
- **THEN** exactamente una entrada está ausente, el gate de confirmación de spec-minada
- **AND** los siete restantes, más el gate de decisión-minada, están presentes con disparadores sin cambios

#### Scenario: Un split de sub-conteo nombra qué subgrupo se achicó

- **GIVEN** una declaración que parte el total en subgrupos
- **WHEN** se lee después del cambio
- **THEN** el subgrupo de gates ratificantes es el que perdió un miembro
- **AND** el subgrupo de momentos de orquestador queda sin cambios

#### Scenario: Dos puntos de interrupción fuera de los ocho se nombran, no se absorben

- **GIVEN** el reclamo calificado en el contrato compartido de presentación
- **WHEN** se lee
- **THEN** nombra el reporte de merge-conflict del Change Workspace y el STOP de atomicidad de commit como puntos de interrupción fuera de los ocho
- **AND** enuncia que ninguno de los triggers lo alcanza, y ninguno se modifica

### Requisito: Ninguna declaración del conteo queda desalineada

Cada archivo que declara el conteo de momentos que citan el walkthrough compartido DEBE leer ocho después del cambio, y ninguna declaración DEBE quedar en nueve. Las declaraciones son tres: el contrato compartido de presentación (dos oraciones más el encabezado de su tabla de momentos), este requisito (cuerpo y ambos títulos de escenario), y la fila de Referencias de `flow/ratify-gate-items`.

#### Scenario: El barrido del conteo no deja un archivo atrás

- **GIVEN** el repo después del cambio
- **WHEN** se busca cada declaración del conteo de momentos que citan el walkthrough
- **THEN** las tres leen ocho, y ninguna lee nueve

## Escenarios

El spec no declara escenarios adicionales más allá de los definidos en sus diez requisitos anteriormente.

## Referencias

- **Conceptualmente relacionado**: `flow/ratify-gate-items.md` — define cómo un gate presenta una vez que disparó
- **Conceptualmente relacionado**: `rule/contract-shape-proposal.md` — define el contenido del trigger (b)
- **Conceptualmente relacionado**: [`closed-forms-for-thread-prose.md`](closed-forms-for-thread-prose.md) — dueño de la forma del reporte entre-fases; este spec sólo fija qué del item aparece ahí
