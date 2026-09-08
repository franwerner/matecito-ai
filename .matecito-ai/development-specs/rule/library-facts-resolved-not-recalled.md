# Capability — Un dato de librería se resuelve, no se recuerda

- **Status:** Accepted
- **Date:** 2026-09-08
- **Components:** cli

## Propósito

Que ningún dato de una dependencia externa —una versión, una API, una opción de configuración, un estado de soporte, un paso de migración— entre en un artefacto del flujo por memoria. La falla que cierra es silenciosa: un dato viejo que suena plausible viaja por el flujo sin que nada lo contradiga, y recién se rompe cuando alguien lo ejecuta.

## Actores

- La fase que explora — compara enfoques y nombra dependencias candidatas.
- La fase que propone — fija el enfoque y con él sus dependencias.
- Las fases que ya estaban cableadas — diseño e implementación.
- El lector del artefacto, que no puede distinguir un dato resuelto de uno recordado.

## Requisitos

### Requisito: Un dato de librería se resuelve, no se recuerda

Toda fase a punto de afirmar en su artefacto la versión, la API, una opción de configuración, el estado de soporte o un paso de migración de una dependencia externa DEBE resolver ese dato a través de la capacidad de documentación de librerías antes de escribirlo. NO DEBE afirmarlo de memoria ni derivarlo del nombre de la librería.

La obligación la dispara la afirmación, no la fase: una fase cuyo artefacto no afirma ningún dato de ese tipo NO DEBE consultar nada, y esa ausencia NO DEBE reportarse como omisión.

#### Scenario: una fase fija una versión concreta

- **GIVEN** una fase a punto de escribir una versión concreta de una dependencia en su artefacto
- **WHEN** escribe esa línea
- **THEN** la versión salió de una resolución contra la documentación de la librería, no de memoria

#### Scenario: una comparación de dos librerías

- **GIVEN** una fase que compara dos librerías como alternativas de un mismo enfoque
- **WHEN** describe qué ofrece cada una
- **THEN** los datos de ambas salen de la capacidad
- **AND** ninguna de las dos queda descrita de memoria porque la otra sí se resolvió

#### Scenario: una fase que no nombra ninguna dependencia

- **GIVEN** una fase cuyo artefacto no afirma ningún dato de una dependencia externa
- **WHEN** termina
- **THEN** no consultó la capacidad
- **AND** eso no se cuenta como paso salteado

### Requisito: Las fases que comparan enfoques o nombran dependencias empiezan cableadas

Toda fase del flujo que compare enfoques o nombre dependencias externas —hoy exploración y propuesta, sumadas a diseño e implementación, que ya lo estaban— DEBE empezar con la capacidad de documentación de librerías alcanzable y con su integración entre sus tools. El grant DEBE ser a nivel de servidor y NO DEBE pinnear nombres de tools individuales, que cambian entre versiones del servidor.

Las dos fases que se cablean acá DEBEN alcanzar la capacidad por el mismo mecanismo: nada las distingue en este respecto, y dos mecanismos para la misma capacidad en fases hermanas es divergencia sin causa. Este requisito NO fija cuál mecanismo — exige que la capacidad sea alcanzable y que sea el mismo en las dos.

#### Scenario: exploración y propuesta arrancan cableadas

- **GIVEN** las fases de exploración y de propuesta
- **WHEN** cualquiera de ellas arranca
- **THEN** la integración de documentación de librerías está entre sus tools
- **AND** la capacidad es alcanzable desde sus instrucciones

#### Scenario: el grant es a nivel de servidor

- **GIVEN** la línea de tools de cualquiera de esas fases
- **WHEN** se lee el grant de documentación de librerías
- **THEN** nombra el servidor entero, no una tool suya

#### Scenario: el mismo mecanismo en las dos

- **GIVEN** las dos fases cableadas por este cambio
- **WHEN** se compara cómo alcanza cada una la capacidad
- **THEN** es el mismo mecanismo en las dos

### Requisito: El cableado incluye el disparador, no sólo el acceso

Un acceso concedido a una fase sin una instrucción que diga cuándo usarlo NO cuenta como cableado: la fase arranca con superficie que nada la lleva a usar. Las instrucciones que la fase sigue DEBEN nombrar la capacidad y decir en qué momento se recurre a ella.

Nombrar la integración en lugar de la capacidad NO alcanza: un lector que sólo ve el nombre del servidor no sabe qué procedimiento seguir. La instrucción nombra la capacidad; que además diga por qué integración está respaldada es permitido, no requerido.

#### Scenario: un acceso sin disparador no cuenta como cableado

- **GIVEN** una fase con el acceso a la documentación concedido y sin ninguna instrucción que diga cuándo recurrir a la capacidad
- **WHEN** se evalúa si está cableada
- **THEN** no lo está: el acceso quedó como superficie inerte

#### Scenario: el disparador vive en las instrucciones que la fase sigue

- **GIVEN** las instrucciones de método que la fase lee y ejecuta
- **WHEN** se busca en ellas cuándo consultar la documentación de una librería
- **THEN** aparece la capacidad nombrada y el momento en que se recurre a ella

## Escenarios

El spec no declara escenarios adicionales más allá de los definidos en sus tres requisitos anteriores.

## Referencias

- **Conceptualmente relacionado**: `flow/find-durable-records.md` — la misma forma de capacidad (grant a nivel de servidor + disparador en la instrucción de la fase), aplicada a otra integración.
