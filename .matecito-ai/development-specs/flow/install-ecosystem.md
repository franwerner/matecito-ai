# Capability — Instalar el ecosistema en la máquina del usuario

- **Status:** Accepted
- **Date:** 2026-08-02
- **Components:** cli

## Propósito

Dejar la máquina de la persona con todo el ecosistema puesto y al día en un solo comando: los binarios que declaran los dominios activos, el payload desplegado en el host destino y la configuración del host reconciliada — mostrando el plan exactamente una sola vez.

## Actores

- **La persona** que corre `matecito-ai install` desde una terminal.

## Precondiciones

- Se puede resolver la carpeta de respaldo de la corrida; si no, el comando falla antes de tocar nada.

## Flujo principal

1. La persona corre `matecito-ai install`.
2. El sistema detecta **una sola vez** el estado observado de cada componente gestionado: el propio ejecutable, los binarios que declaran los dominios activos, el payload desplegado en el host destino y la configuración del host.
3. Deriva de ese estado el plan: por componente, *instalar*, *actualizar* u *omitir*.
4. Si ninguna acción queda activa, informa que no hay nada para hacer y termina sin tocar nada.
5. Muestra el plan numerado: por componente, el verbo que le corresponde (*instalar* o *actualizar*), su origen de payload si procede, y para el componente de payload el desglose por categoría (*nuevos*, *cambiados*, *a borrar*) y qué destinos se van a borrar.
6. Pide confirmación (`y/N`) y ejecuta las acciones en orden, continuando ante el fallo de cualquiera de ellas y mostrando el resultado de cada una.
7. Durante la ejecución del componente de payload, informa aquellos destinos que se han preservado por ser legacy sin coincidencia de hash, con la acción sugerida si la persona desea eliminarlos, y aquellos cuyas ediciones previas fueron pisadas al borrar, informando dónde quedó el respaldo.
8. Si el propio ejecutable fue reemplazado durante la corrida, entrega la ejecución al binario nuevo (ver [`resume-self-replace-run.md`](resume-self-replace-run.md)).
9. Termina en error **si y solo si al menos un componente falló**, nombrando cuáles y conservando sus causas; si no, termina OK.

## Ramas / flujos alternativos

- **`--dry-run`** → se muestra el plan y el comando termina sin ejecutar ninguna acción ni pedir confirmación.
- **`--yes`** → el plan se muestra igual, pero no se pide confirmación: se ejecuta directo.
- **Respuesta no afirmativa en la confirmación** → se informa que se canceló y no se ejecuta ninguna acción.
- **Corrida relanzada tras un auto-reemplazo** → se omiten por completo el plan, el `--dry-run` y la confirmación propios de este comando: esa corrida ya fue confirmada una vez. Ver [`resume-self-replace-run.md`](resume-self-replace-run.md).

## Casos borde

- **Un componente falla** → se registra su error, se informa el fallo y la corrida **sigue** con los componentes restantes (continue-on-error). No hay rollback.
- **Terminada la corrida, cualquier componente fallado la lleva a error**: ningún fallo se traga. El código de salida es binario —hubo fallos o no— y son el mensaje y la salida por pantalla los que dicen cuáles.
- **No hay nada desactualizado ni faltante** → el comando informa "nada para hacer" y no imprime plan ni pide confirmación.

## Reglas de negocio

- La detección del estado se hace **una sola vez** por corrida y se reutiliza para ejecutar: no se vuelve a consultar el estado remoto entre el plan y la ejecución.
- El verbo del plan lo determina el estado del componente: ausente → *instalar*; presente pero desactualizado, con payload cambiado o con reconciliación pendiente → *actualizar*.
- La confirmación acepta como afirmativas `y`, `yes`, `s`, `si` y `sí`, sin distinguir mayúsculas; **cualquier otra respuesta, o el cierre de la entrada, cancela**.
- Una vez confirmado en este comando, el motor de sincronización no vuelve a preguntar: la confirmación ocurre una sola vez por corrida.
- El error de salida se propaga **envuelto**, conservando la causa original de cada componente fallado y nombrándolos a todos. El detalle por componente ya se mostró durante la ejecución; el error final dice cuáles fallaron, no vuelve a explicarlos.

## Errores de cara al actor

- **No se pudo resolver la carpeta de respaldo** → el comando falla antes de detectar o ejecutar nada.
- **Falló al menos un componente** → el comando termina en error, nombrando los componentes que fallaron y conservando sus causas originales.

## Entidades y estados

- **Plan** — representa el conjunto de acciones a ejecutar por componente. Estado: mostrado (una sola vez) → ejecutándose (engine ejecuta sin re-mostrar).

## Escenarios

### Scenario: plan combinado con verbo por componente

- **GIVEN** un entorno con un componente ausente y otro presente pero desactualizado
- **WHEN** la persona corre `matecito-ai install`
- **THEN** el sistema muestra un plan numerado que marca *instalar* para el ausente y *actualizar* para el desactualizado, y para el componente de payload agrega su origen

### Scenario: dry-run no ejecuta nada

- **GIVEN** un plan con al menos una acción activa
- **WHEN** la persona corre `matecito-ai install --dry-run`
- **THEN** el sistema muestra el plan, avisa que no se ejecutó nada y termina sin pedir confirmación ni modificar el entorno

### Scenario: confirmación negativa cancela

- **GIVEN** un plan con acciones activas y sin `--yes`
- **WHEN** la persona responde algo distinto de una afirmación en el prompt `y/N`
- **THEN** el sistema informa que se canceló y no ejecuta ninguna acción

### Scenario: `--yes` no pregunta

- **GIVEN** un plan con acciones activas
- **WHEN** la persona corre `matecito-ai install --yes`
- **THEN** el sistema muestra el plan y ejecuta las acciones sin pedir confirmación

### Scenario: nada para hacer

- **GIVEN** un entorno donde todo está presente y al día
- **WHEN** la persona corre `matecito-ai install`
- **THEN** el sistema informa que no hay nada para hacer y termina sin plan, sin confirmación y sin cambios

### Scenario: un componente falla y la corrida sigue

- **GIVEN** un plan con varias acciones activas donde la primera falla
- **WHEN** el sistema las ejecuta
- **THEN** informa el fallo de esa acción y continúa ejecutando las restantes, sin cancelar la corrida

### Scenario: el fallo de cualquier componente sale como error

- **GIVEN** una corrida donde falló el despliegue del payload y también un binario, y el componente de exploración estructural del código terminó bien
- **WHEN** el comando termina
- **THEN** devuelve un error que nombra los dos componentes fallados y envuelve sus causas originales, en lugar de terminar OK

### Scenario: sin componentes fallados termina OK

- **GIVEN** una corrida donde todos los componentes terminaron bien
- **WHEN** el comando termina
- **THEN** termina OK

### Scenario: corrida relanzada tras auto-reemplazo

- **GIVEN** una corrida relanzada porque el propio ejecutable fue reemplazado
- **WHEN** el comando arranca
- **THEN** no muestra plan, no evalúa `--dry-run` y no pide confirmación: pasa directo a ejecutar lo que queda

### Scenario: no se imprime el plan dos veces

- **GIVEN** una corrida de `install` que muestra el plan y después invoca al motor
- **WHEN** se ejecuta
- **THEN** el plan —con su desglose y sus destinos a borrar— aparece una sola vez en la salida

### Scenario: el plan anticipa los borrados y reporta lo que tocó

- **GIVEN** un payload con destinos a borrar: uno del registro que la persona editó, y un huérfano legacy que no coincide
- **WHEN** el sistema muestra el plan y ejecuta
- **THEN** lista el desglose y los destinos a borrar antes de pedir confirmación, y al terminar informa la edición pisada con su respaldo y el legacy preservado con su ruta **y la acción que la persona puede tomar sobre él**, sin que el comando falle

## Referencias

- **Flow** → [`update-ecosystem.md`](update-ecosystem.md) — el mismo motor de detección y ejecución, expuesto como comando de actualización.
- **Flow** → [`resume-self-replace-run.md`](resume-self-replace-run.md) — qué hace la corrida relanzada tras el auto-reemplazo del ejecutable.
- **Process** → [`../process/deploy-payload-to-host.md`](../process/deploy-payload-to-host.md) — el motor de despliegue del payload que este comando invoca como uno de sus componentes.
- **Rule** → [`../rule/domain-activation-shim.md`](../rule/domain-activation-shim.md) — qué dominios se consideran activos al detectar binarios y payload.
