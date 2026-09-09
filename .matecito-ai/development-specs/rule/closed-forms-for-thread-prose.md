# Capability — Formas cerradas para la prosa del hilo

- **Status:** Accepted
- **Date:** 2026-09-09
- **Components:** cli

## Propósito

Que lo que el asistente dice en el hilo tenga una forma fijada por acto en lugar de componerse libre,
para que "sé breve" deje de ser una instrucción incomprobable. Cierra el hueco por el que el acto más
frecuente del orquestador —el reporte que emite entre una fase y la siguiente— no tenía forma a la que
atarse, y el que dejaba pasar un preámbulo propio por la puerta de la excepción de emisión.

## Actores

- **El orquestador** — único que emite el reporte entre-fases y único que despacha.
- **El asistente** — en cualquier turno del hilo, cuando presenta una decisión, explica un hallazgo o
  informa lo que hizo.
- **El usuario** — recibe el turno y no puede des-leer lo que se le ofreció sin pedirlo.

## Requisitos

### Requisito: Los actos comunicativos del hilo son un conjunto cerrado, y el reporte entre-fases es uno de ellos

El texto de gobernanza que el asistente carga siempre DEBE enumerar los actos comunicativos que cubren
lo que dice en el hilo como un **conjunto cerrado**, y ese conjunto DEBE incluir el reporte que el
orquestador emite entre una fase y la siguiente. Cada acto DEBE llevar su forma propia enunciada; un
acto nombrado sin forma NO cuenta como declarado, porque no hay nada contra lo que chequearlo.

Los actos que ya estaban declarados —presentar una decisión que el usuario decide, explicar un hallazgo
a pedido, informar lo que se hizo— DEBEN quedar con su forma sin cambios. El acto nuevo se suma; NO DEBE
reemplazar ni reformular a ninguno de los tres.

#### Scenario: el conjunto declarado son cuatro actos

- **verification:** `deferred → tighten-orchestrator-prose`
- **GIVEN** la sección de gobernanza que fija las formas cerradas
- **WHEN** se la lee
- **THEN** enumera cuatro actos comunicativos y uno de ellos es el reporte entre-fases
- **AND** cada uno de los cuatro lleva su forma enunciada

#### Scenario: el acto nuevo no desplaza a los tres previos

- **verification:** `deferred → tighten-orchestrator-prose`
- **GIVEN** los tres actos que la sección ya declaraba
- **WHEN** se la lee después del cambio
- **THEN** los tres siguen declarados con su forma sin alterar
- **AND** el cuarto aparece como agregado, no como reemplazo de ninguno

#### Scenario: un acto sin forma no está declarado

- **verification:** `deferred → tighten-orchestrator-prose`
- **GIVEN** un acto nombrado en la enumeración sin una forma propia enunciada
- **WHEN** se evalúa si el conjunto está declarado
- **THEN** ese acto no cuenta como declarado, porque no hay forma contra la que chequear un turno

### Requisito: La forma del reporte entre-fases es ancla más una línea, sin narrativa ni acciones de ratificación

El reporte entre-fases DEBE tomar, por item, un ancla a la fuente concreta sobre la que el item es y
**una sola línea** de resumen. NO DEBE llevar narrativa, preámbulo, recapitulación del contexto ni
justificación previa al item.

NO DEBE ofrecer acciones de ratificación. El reporte es informativo: lo que se ratifica se ratifica en
un gate, que es otro momento con su propia forma, y ofrecer confirmar en un reporte informativo produce
exactamente la fatiga de confirmación que el gate existe para evitar.

La forma DEBE tener **un solo dueño declarado**. La regla del disparo de gates enuncia qué aparece en el
reporte cuando un item no dispara; NO DEBE enunciar además cómo se ve el reporte. Dos enunciados
parciales de la misma forma es la divergencia que este requisito cierra.

#### Scenario: un reporte de dos items no lleva narrativa

- **verification:** `deferred → tighten-orchestrator-prose`
- **GIVEN** una fase que retornó dos items que no disparan
- **WHEN** el orquestador emite el reporte entre-fases
- **THEN** cada item aparece como su ancla más una línea de resumen
- **AND** el turno no lleva preámbulo, recapitulación del contexto ni cierre que reenuncie lo dicho

#### Scenario: el reporte informativo no ofrece ratificar

- **verification:** `deferred → tighten-orchestrator-prose`
- **GIVEN** un reporte entre-fases con al menos un item
- **WHEN** se lo presenta
- **THEN** no ofrece confirmar, corregir ni elegir sobre ningún item
- **AND** nada del flujo queda esperando una respuesta a ese turno

#### Scenario: la forma tiene un solo dueño

- **verification:** `deferred → tighten-orchestrator-prose`
- **GIVEN** la regla del disparo de gates y la regla de las formas cerradas
- **WHEN** se busca en las dos cómo se ve el reporte entre-fases
- **THEN** solo la de las formas cerradas lo enuncia
- **AND** la del disparo enuncia únicamente qué del item aparece ahí

### Requisito: El alcance de la excepción de emisión es cerrado y nombrado

Cuando otra regla del ecosistema ordena emitir contenido, esa regla gana y el contenido se emite
entero: eso NO cambia. Lo que cambia es el alcance de lo que el presupuesto sí acota, que DEBE quedar
**nombrado**.

La excepción NO DEBE enunciarse como "la prosa alrededor de ese material". "Alrededor" no es un
alcance: se lee como permiso para un preámbulo propio que las formas cerradas ya prohíben, y así la
excepción termina autorizando justo lo que la sección prohíbe tres párrafos antes. El texto DEBE nombrar
qué puede cubrir la prosa propia del asistente cuando una forma cerrada aplica.

#### Scenario: el alcance se nombra en lugar de quedar en "alrededor"

- **verification:** `deferred → tighten-orchestrator-prose`
- **GIVEN** el párrafo que declara la excepción de emisión
- **WHEN** se lo lee
- **THEN** nombra qué puede cubrir la prosa propia del asistente cuando una forma cerrada aplica
- **AND** no deja el alcance enunciado como "la prosa alrededor de ese material"

#### Scenario: el contenido mandado sigue emitiéndose entero

- **verification:** `deferred → tighten-orchestrator-prose`
- **GIVEN** una sección de retorno cuyo contrato parte cada item en resumen y razonamiento
- **WHEN** la fase la emite
- **THEN** las dos partes van completas al bloque
- **AND** el alcance recortado de la excepción no alcanza a ninguna de las dos

### Requisito: Narrar el prompt de despacho está prohibido por nombre

La lista de prohibiciones DEBE llevar, **nombrada**, la narración de las instrucciones con que se
despachó una fase: qué se le pidió, qué artefactos se le pasaron, cómo se le acotó el alcance. Es la
mayor fuente de contexto reenunciado en un turno entre fases, y hasta ahora nada la prohibía: cada forma
cerrada la excluía por implicación, que es exactamente el tipo de regla que este ecosistema no da por
cumplida.

La prohibición DEBE estar redactada para atar al orquestador. Una fase no tiene prompt de despacho
propio que narrar, así que una redacción genérica no le agrega obligación a nadie y sí deja al
orquestador con la lectura de que no lo alcanza.

#### Scenario: el turno de despacho no narra lo que se pidió

- **verification:** `deferred → tighten-orchestrator-prose`
- **GIVEN** el orquestador que acaba de despachar una fase
- **WHEN** informa al usuario
- **THEN** nombra la fase despachada y nada más de las instrucciones que le dio
- **AND** no reproduce el pedido, ni los artefactos que le pasó, ni el alcance con que la acotó

#### Scenario: la prohibición nombra al orquestador

- **verification:** `deferred → tighten-orchestrator-prose`
- **GIVEN** la lista de prohibiciones después del cambio
- **WHEN** se la lee
- **THEN** una de sus entradas nombra la narración del prompt de despacho
- **AND** está redactada de modo que ata al orquestador, que es el único que despacha

### Requisito: Las cuatro formas se auto-reportan; lo verificable es que la regla exista

Nada en el ecosistema observa el texto que el asistente escribe en el hilo: los mecanismos de chequeo se
atan al ciclo de vida de las herramientas, no al mensaje. Las cuatro formas, el alcance cerrado de la
excepción y la prohibición de narrar el despacho son por tanto **auto-reportadas**, exactamente como ya
lo eran las tres formas previas.

El texto de gobernanza DEBE enunciar ese costo llano y NO DEBE presentar ninguna verificación como
existente. Lo verificable de esta capacidad es que **la regla esté escrita**, nunca que la conducta la
siga: un chequeo de presencia de texto DEBE reportarse como que la regla existe, y NO DEBE reportarse
como que el asistente la cumple. Confundir las dos convierte una garantía nula en una garantía aparente,
que es peor que ninguna.

#### Scenario: el costo se enuncia, no se implica

- **verification:** `deferred → tighten-orchestrator-prose`
- **GIVEN** la sección que declara las cuatro formas
- **WHEN** se la lee
- **THEN** enuncia que su cumplimiento es auto-reportado y que nada lo observa
- **AND** no ofrece ningún mecanismo del ecosistema como detección de su incumplimiento

#### Scenario: un chequeo de presencia no se cuenta como verificación de conducta

- **verification:** `deferred → tighten-orchestrator-prose`
- **GIVEN** un chequeo que confirma que las cuatro formas y las dos prohibiciones están escritas
- **WHEN** se reporta su resultado
- **THEN** se reporta como que la regla existe en el texto
- **AND** no se reporta como que el turno del asistente la respetó

## Escenarios

El spec no declara escenarios adicionales más allá de los definidos en sus cinco requisitos
anteriormente.

## Referencias

- **Conceptualmente relacionado**: [`gate-firing-triggers.md`](gate-firing-triggers.md) — enuncia qué
  del item aparece en el reporte entre-fases cuando no dispara; este spec enuncia cómo se ve ese
  reporte.
- **Conceptualmente relacionado**: [`../flow/ratify-gate-items.md`](../flow/ratify-gate-items.md) —
  cómo se presenta un gate una vez que disparó, un momento distinto con su propia forma.
