# Capability — Localizar los registros duraderos que gobiernan el trabajo

- **Status:** Accepted
- **Date:** 2026-09-02
- **Components:** cli

## Propósito

Localizar los registros duraderos —decisiones y comportamiento— que gobiernan un trabajo antes de escribir código, proponer un diseño o verificar contra ellos. La falla que cierra es silenciosa: un registro que existe, gobierna el trabajo y usa vocabulario que el lector no comparte deja sin rastro cuando nadie lo consulta.

## Actores

- **Cualquier agente** a punto de tocar arquitectura, capas, contratos, datos, transporte, convenciones o comportamiento.
- **Las fases que deben leer "los registros que este cambio toca"** — diseño y verificación.

## Requisitos

### Requisito: Mapear los stores antes de decidir qué leer

El sistema DEBE primero determinar qué stores de registros el proyecto tiene, cuál gobierna cada uno, cuántos registros cada uno sostiene, y el nombre bajo el que cada uno se busca. Ese mapeo DEBE venir de la declaración propia del proyecto, nunca de inspeccionar el árbol de directorios — un componente puede tener su propio store o estar gobernado por el store raíz, y la diferencia no es visible en el árbol. Un cambio que toca varios componentes DEBE estar gobernado por la unión de sus stores más el store raíz; cuando el proyecto declara un componente raíz, el store raíz le pertenece y DEBE ser solicitado explícitamente como cualquier otro.

#### Scenario: un cambio a través de dos componentes

- **GIVEN** un cambio cuyo alcance confirmado nombra dos componentes
- **WHEN** el sistema mapea los stores
- **THEN** el conjunto candidato cubre los stores de ambos componentes y el store raíz, no solo uno

#### Scenario: un componente sin store propio

- **GIVEN** un componente sin store bajo su propia ruta, en un proyecto que no declara componente raíz
- **WHEN** el sistema mapea los stores
- **THEN** ese componente es reportado como gobernado por el store raíz, nunca como sin registros

#### Scenario: un nombre de componente que no existe

- **GIVEN** un nombre de componente ausente de la declaración del proyecto
- **WHEN** el sistema mapea los stores
- **THEN** es ignorado con un aviso, nunca silenciosamente

### Requisito: Reunir candidatos por tres caminos acumulativos

El sistema DEBE reunir candidatos por el índice del store, búsqueda literal y búsqueda semántica, y DEBE acumular a través de ellas en lugar de elegir una — cada una cubre un punto ciego que las otras no tienen. El índice del store DEBE ser leído siempre: es el único camino que enumera.

#### Scenario: los tres caminos disponibles

- **GIVEN** un proyecto cuyo camino semántico está configurado
- **WHEN** el sistema reúne candidatos
- **THEN** el índice del store es leído completo, y ambas búsquedas contribuyen a un conjunto candidato acumulado

### Requisito: Degradar por presencia y decirlo

Cuando el camino semántico no está disponible, el procedimiento DEBE continuar en los otros dos caminos sin error y sin mencionar la integración faltante por nombre en sus instrucciones. Debe, sin embargo, decirlo al presentar los candidatos — para que un lector pueda diferenciar "se buscó de tres formas" de "se buscó de dos, algo puede faltar".

#### Scenario: camino semántico no disponible

- **GIVEN** una sesión donde la integración de búsqueda semántica no es alcanzable
- **WHEN** el sistema presenta sus candidatos
- **THEN** los candidatos se presentan normalmente
- **AND** la presentación dice que no se usó el camino semántico

#### Scenario: camino semántico disponible

- **GIVEN** una sesión donde la integración de búsqueda semántica responde
- **WHEN** el sistema presenta sus candidatos
- **THEN** no se reporta ausencia alguna

### Requisito: Decidir sobre el conjunto completo, nunca sobre el top hit

La decisión DEBE ser tomada leyendo el conjunto candidato acumulado, no tomando el resultado de rango más alto. Los resultados semánticos DEBEN ser solicitados a una amplitud mínima de cinco y DEBEN estar alcanzados a los stores que aplican; los índices de navegación DEBEN estar excluidos de lo que el camino semántico devuelve, porque hacen match con todo y desplazan registros reales de la parte alta.

#### Scenario: el registro que gobierna no fue ranqueado primero

- **GIVEN** un conjunto candidato donde el registro que gobierna el trabajo está ranqueado tercero
- **WHEN** el sistema decide qué registros aplican
- **THEN** lee el conjunto completo e identifica ese registro, en lugar de frenar en el primero

#### Scenario: una búsqueda sin alcance cruzaría proyectos

- **GIVEN** un índice de búsqueda compartido que sirve varios proyectos
- **WHEN** el sistema corre el camino semántico
- **THEN** alcanza la consulta a los stores que aplican a este proyecto, para que ningún registro de otro proyecto entre el conjunto

### Requisito: El ranking no es cobertura

En verificación, los resultados de los caminos de búsqueda DEBEN ser tratados como candidatos a leer, nunca como prueba de que nada más aplica. Solo el índice del store lleva una garantía de enumeración. Un registro que nunca salió del radar no deja rastro de su ausencia.

#### Scenario: la verificación no reclama completitud desde la búsqueda

- **GIVEN** un paso de verificación que usó los caminos de búsqueda
- **WHEN** reporta qué registros verificó
- **THEN** no presenta la ausencia de un hit de búsqueda como evidencia de que no registra nada aplica

### Requisito: Disponible para las fases que leen registros, sin nombrar la integración

Toda fase que deba leer registros durables — hoy especificación, diseño, tareas, verificación e implementación, las cinco que la tabla de lectura del dominio marca leyendo decisiones o capability-specs — DEBE empezar con esta capacidad ya disponible para ella, y con la integración de búsqueda entre sus tools. Sus instrucciones DEBEN nombrar la capacidad y decir cuándo recurrir a ella, y NO DEBEN nombrar la integración que la implementa: la capacidad la conoce, las fases conocen la capacidad.

El criterio manda y la enumeración lo sigue. Cuando la fila de lectura de una fase pasa a incluir decisiones o capability-specs, esa fase DEBE quedar cableada y la enumeración DEBE nombrarla: una enumeración que se queda atrás del criterio no reduce la obligación, sólo la esconde.

#### Scenario: las cinco fases empiezan cableadas

- **GIVEN** las fases de especificación, diseño, tareas, verificación e implementación
- **WHEN** cualquiera de ellas arranca
- **THEN** esta capacidad está precargada y la integración de búsqueda está entre sus tools
- **AND** ninguna de las instrucciones de la fase menciona la integración por nombre

#### Scenario: la enumeración sigue al criterio

- **GIVEN** una fase cuya fila de lectura pasó a incluir capability-specs sin que la enumeración la nombrara
- **WHEN** se compara la enumeración con el criterio
- **THEN** la enumeración la incluye y la fase está cableada
- **AND** el criterio queda sin cambios: es la enumeración la que se corrige

#### Scenario: las instrucciones nombran la capacidad, no lo que la implementa

- **GIVEN** el cuerpo de instrucciones de cualquiera de esas fases
- **WHEN** se lo lee buscando cómo localizar los registros que la fase debe leer
- **THEN** nombra la capacidad y dice cuándo recurrir a ella
- **AND** no aparece el nombre de la integración de búsqueda en ninguna parte

#### Scenario: una fase que no lee registros no se cablea

- **GIVEN** una fase cuya fila de lectura no incluye decisiones ni capability-specs
- **WHEN** arranca
- **THEN** ni la capacidad ni la integración están entre lo que se le precarga

#### Scenario: el helper se resuelve donde la capacidad lo dice

- **GIVEN** las instrucciones propias de la capacidad, que nombran una ruta desplegada para su helper de mapeo de stores
- **WHEN** la capacidad se ejecuta en una máquina donde el ecosistema está desplegado
- **THEN** el helper es alcanzable exactamente en esa ruta, sin edit de esas instrucciones

## Referencias

- **Process** → [`process/configure-record-search.md`](../process/configure-record-search.md) — la configuración única que activa el camino semántico.
- **Process** → [`process/install-mcp-from-release-tarball.md`](../process/install-mcp-from-release-tarball.md) — cómo llega la integración de búsqueda a la máquina.
