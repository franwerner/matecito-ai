# Capability — Materializar el delta de comportamiento de un cambio en el store durable

- **Status:** Accepted
- **Date:** 2026-09-08
- **Components:** cli

## Propósito

que el comportamiento que un cambio agrega o modifica quede escrito en el store durable una sola vez, por la misma fase que ya materializa lo demás que el cambio deja en disco, y que el cierre del cambio quede como una operación de registro que no escribe nada.

## Actores

- la fase que implementa el cambio —dispara el pliegue al terminar su trabajo—
- la fase que archiva —cierra el cambio sin escribir en el store—
- la fase que verifica —lee el store durable después del pliegue—
- el lector del store durable

## Flujo principal

1. La fase que implementa el cambio llega al despacho en que no quedan tareas pendientes.
2. Evalúa si el delta del cambio toca alguna capacidad.
3. Si no hay capacidades, termina sin escribir.
4. Si hay capacidades, para cada una: decide si es Nueva, ADDED a una existente, o MODIFIED en una existente.
5. Para capacidades Nuevas: crea el archivo durable y lo puebla con los requisitos y escenarios del delta.
6. Para ADDED/MODIFIED: pliega el delta no-destructivamente en el archivo durable.
7. Actualiza el índice de tipo y el índice raíz.
8. La fase que archiva, después, cierra el cambio sin escribir en el store durable.

## Reglas de negocio

- El pliegue DEBE dispararse en el despacho de implementación que da el trabajo del cambio por terminado —el que no deja tareas pendientes—, y DEBE abarcar el delta completo del cambio, todas sus capacidades a la vez.
- Un despacho de continuación, o una ronda de trabajo en paralelo que deja tareas pendientes, NO DEBE plegar nada: el store nunca queda a medio plegar entre despachos.
- El pliegue NO DEBE hacerse por tarea ni por capacidad: nada en el cambio asocia una tarea con la capacidad que implementa, y plegar antes de que todas las tareas de una capacidad hayan aterrizado documentaría comportamiento que todavía no está construido.
- El pliegue es no destructivo: preserva todo lo que el spec durable ya tenía y el delta no menciona, con sus tokens de verificación si los tenían.
- Lo que el delta ya escribió se copia tal como está, no se recompone; sólo se sintetiza lo que el delta no trae y el spec durable necesita para estar completo.
- Un pliegue que dejaría afuera contenido que el delta nunca mencionó detiene la fase en vez de escribir; ese bloqueo viaja por el único lugar de bloqueo que la fase declara, junto a sus otras causas, y NO se repite en ninguna otra sección del retorno.

## Escenarios

### Scenario: el cierre del cambio no toca el store durable

- **verification:** `deferred → rebind-spec-referencias`
- **GIVEN** un cambio cuyo delta declara comportamiento nuevo y comportamiento modificado
- **WHEN** corre la fase que archiva
- **THEN** no se crea ni se modifica ningún archivo del store durable de comportamiento
- **AND** su única salida persistida es su reporte, en la memoria del flujo

### Scenario: el cierre del cambio sigue leyendo el delta

- **verification:** `deferred → rebind-spec-referencias`
- **GIVEN** la fase que archiva, que necesita nombrar los artefactos del cambio en su reporte
- **WHEN** los reúne
- **THEN** el delta del cambio sigue entre lo que lee
- **AND** esa lectura no deriva en ninguna escritura al store durable

### Scenario: un solo dueño declarado

- **verification:** `deferred → rebind-spec-referencias`
- **GIVEN** la tabla que declara qué lee y qué escribe cada fase del flujo
- **WHEN** se busca quién escribe el store durable de comportamiento
- **THEN** aparece únicamente la fase que implementa, con la lectura del store que va a plegar
- **AND** la fase que archiva ya no figura como escritora de ese store

### Scenario: un despacho que deja tareas pendientes no pliega

- **verification:** `deferred → rebind-spec-referencias`
- **GIVEN** un despacho de implementación que termina con tareas del cambio todavía pendientes
- **WHEN** el despacho cierra
- **THEN** el store durable queda exactamente como estaba
- **AND** el pliegue queda pendiente para el despacho que cierre el cambio

### Scenario: el despacho que cierra pliega el delta entero

- **verification:** `deferred → rebind-spec-referencias`
- **GIVEN** un despacho de implementación que termina sin tareas pendientes, en un cambio cuyo delta toca dos capacidades
- **WHEN** el despacho cierra
- **THEN** las dos capacidades quedan plegadas en el store durable en esa misma corrida
- **AND** ninguna quedó plegada antes que la otra

### Scenario: dos despachos no producen dos pliegues

- **verification:** `deferred → rebind-spec-referencias`
- **GIVEN** un cambio que necesitó un despacho de continuación antes de terminar
- **WHEN** el segundo despacho cierra el cambio
- **THEN** el delta se pliega una sola vez, en ese despacho
- **AND** el store no muestra rastro de un pliegue anterior parcial

### Scenario: lo que el delta no menciona sobrevive

- **verification:** `deferred → rebind-spec-referencias`
- **GIVEN** un spec durable con escenarios que el delta del cambio no menciona
- **WHEN** se pliega el delta
- **THEN** esos escenarios quedan tal cual, con sus tokens de verificación si los tenían

### Scenario: un pliegue que perdería contenido detiene la fase

- **verification:** `deferred → rebind-spec-referencias`
- **GIVEN** un pliegue que dejaría afuera secciones o escenarios que el delta nunca mencionó
- **WHEN** la fase llega a ese punto
- **THEN** no escribe nada en el store durable
- **AND** informa el bloqueo con la pregunta que necesita respuesta antes de seguir

### Scenario: el bloqueo por pliegue destructivo convive con las otras causas

- **verification:** `deferred → rebind-spec-referencias`
- **GIVEN** una fase que declara un único lugar para su bloqueo, con más de una causa posible
- **WHEN** el bloqueo lo causa un pliegue que perdería contenido
- **THEN** se informa en ese único lugar, nombrando esa causa
- **AND** no aparece repetido en la sección de hallazgos, en la de estado ni entre los riesgos

### Scenario: lo que el delta ya escribió se copia, no se recompone

- **verification:** `deferred → rebind-spec-referencias`
- **GIVEN** un escenario que el delta escribió completo
- **WHEN** se pliega
- **THEN** el spec durable lo lleva con el mismo texto, sin reescribirlo ni normalizarlo

### Scenario: una capacidad nueva queda enrutada desde los dos índices

- **verification:** `deferred → rebind-spec-referencias`
- **GIVEN** un delta que declara una capacidad que no existía en el store
- **WHEN** se pliega
- **THEN** el índice de su tipo lleva su fila
- **AND** el índice raíz enruta a ese tipo

### Scenario: un despacho que materializó lo reporta

- **GIVEN** un despacho que cerró el cambio y plegó dos capacidades
- **WHEN** emite su retorno
- **THEN** la sección de capacidades materializadas aparece, con las dos
- **AND** distingue la capacidad creada de la modificada

### Scenario: un despacho que no materializó no emite la sección

- **GIVEN** un despacho de continuación que dejó tareas pendientes
- **WHEN** emite su retorno
- **THEN** la sección de capacidades materializadas no aparece

### Scenario: la verificación encuentra el delta ya plegado

- **verification:** `deferred → rebind-spec-referencias`
- **GIVEN** un cambio que declaró una capacidad nueva y cuya implementación terminó
- **WHEN** corre la verificación y busca el spec durable de esa capacidad
- **THEN** lo encuentra en disco, con el delta ya plegado

### Scenario: hallazgos críticos no revierten el pliegue

- **verification:** `deferred → rebind-spec-referencias`
- **GIVEN** un cambio cuyo delta ya se plegó y cuya verificación reporta hallazgos críticos
- **WHEN** termina la verificación
- **THEN** el store durable conserva el delta plegado
- **AND** nada del flujo lo deshace por ese resultado

### Scenario: la prosa de la verificación nombra lo que el store ya contiene

- **verification:** `deferred → rebind-spec-referencias`
- **GIVEN** los pasos de verificación que leen el store durable de comportamiento
- **WHEN** se los lee
- **THEN** enuncian que el store ya lleva el delta de este cambio
- **AND** distinguen la parte circular del chequeo de la parte que sigue siendo regresión

## Referencias

- **Regla de verificación de escenarios** → [`rule/scenario-verification-scope-token.md`](../rule/scenario-verification-scope-token.md) — cómo el token de verificación rebaja o no el alcance de cobertura de un cambio.
- **Definición canónica del token** → [`~/.claude/references/spec/README.md`](../../../payload/domains/development/references/spec/README.md) — la forma del token, sus valores, su uso per-escenario y su consecuencia aceptada en el spec durable.
