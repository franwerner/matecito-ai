# Capability — Reclamar el turno de la rama compartida

- **Status:** Accepted
- **Date:** 2026-09-10
- **Components:** cli

## Propósito

Que dos sesiones concurrentes en un mismo repositorio no escriban a la vez sobre la rama compartida. Una sesión reclama el turno, hace una unidad de trabajo y lo devuelve con el recibo que el reclamo le dio. Nada intercepta un comando y nada le pregunta al llamador quién es.

## Actores

- **La sesión que va a escribir** en una rama que no posee
- **El run de consolidación** de un batch de implementación — único punto cableado hoy
- **La persona** que consulta el estado o decide qué hacer ante un turno tomado

## Precondiciones

- El repositorio está en estado válido
- Quien va a escribir está listo para reclamar antes de escribir

## Flujo principal

1. Quien va a escribir invoca el comando de reclamo, encadenado delante de su escritura.
2. Si el turno está libre, el reclamo lo toma, termina en cero y devuelve un recibo en su propio renglón.
3. Si el turno está tomado, el reclamo registra la solicitud en la cola, pausa brevemente y reintenta exactamente una vez.
4. Si el reintento tiene éxito, elimina su propia entrada de cola y devuelve el recibo.
5. Si sigue tomado, el reclamo termina distinto de cero y reporta quién sostiene el turno, desde qué árbol, en qué momento, hace cuánto y qué hay en cola; la escritura encadenada no llega a correr.
6. Quien sostiene el turno hace su unidad de trabajo y lo libera invocando el comando de liberación con el recibo que el reclamo le dio.
7. La liberación borra el turno sólo si todavía lleva ese recibo, y termina en cero en cualquiera de sus dos resultados.

## Ramas / flujos alternativos

- **Nadie arbitró** → el reclamo termina en cero pero deja un mensaje en el flujo de error nombrando que no se arbitró y por qué; el trabajo procede igual.
- **El turno se movió por debajo** → la liberación no borra nada, reporta la discrepancia en el flujo de error, y termina en cero igual.
- **Consulta de estado** → puede invocarse sin reclamar nada; siempre termina en cero y nunca modifica el turno ni la cola.

## Casos borde

- **Una sesión muere sosteniendo el turno** → queda tomado; nada lo libera en su nombre.
- **Evidencia fuerte de que quien lo sostiene desapareció** → no autoriza romperlo; la decisión llega al usuario.
- **Entradas de cola abandonadas** → se acumulan tal cual, reportadas como solicitudes registradas, nunca barridas.
- **Un segundo reclamo de quien ya sostiene el turno** → es rechazado igual que si lo sostuviera cualquier otro.

## Reglas de negocio

- El turno se reclama y se libera por llamada explícita; nada intercepta, inspecciona ni clasifica un comando de shell para decidirlo.
- El mecanismo es asesor: funciona porque cada sesión se compromete a consultarlo; nada impide técnicamente escribir sin reclamar.
- Un reclamo tiene exactamente tres resultados distinguibles: tuyo, nadie arbitró, o tomado — nunca silenciosos entre sí.
- La liberación se decide sólo por el recibo; sin recibo es error de uso, y ningún campo del registro se compara ni se testea, sólo se imprime.
- No hay liberación automática, y romper un turno ajeno es siempre una decisión humana.
- El run de consolidación reclama una sola vez por ronda, antes del primer cherry-pick, y libera una sola vez después del último, sin condición y cualquiera sea el resultado del loop.

## Entidades y estados

- **Turno** — un ref único; estados: libre → tomado (con recibo) → libre (liberado o movido). Sin identidad de quién lo sostiene más allá de lo informativo.
- **Cola** — entradas de solicitudes pendientes; se acumulan, nunca se barren automáticamente.
- **Recibo** — el valor devuelto por un reclamo exitoso; es lo único que autoriza una liberación.

## Errores de cara al actor

- **Turno tomado** → código de salida distinto de cero; reporta quién lo sostiene, su árbol, su momento de escritura, hace cuánto, y la cola.
- **Liberación sin recibo** → falla como error de uso; no se borra nada.
- **El turno se movió por debajo** → se reporta la discrepancia; nada se borra; termina en cero.

## Requisitos

### Requisito: El turno se reclama y se libera por llamada explícita, nunca por interceptación

El turno DEBE tomarse invocando el comando de reclamo y devolverse invocando el de liberación. Nada DEBE interceptar, inspeccionar, parsear ni clasificar un comando de shell para decidir si hace falta un turno, y ningún manejador registrado en el host DEBE pertenecer a este mecanismo. El enforcement vive en el código de salida del reclamo, que quien escribe encadena delante de su escritura.

El mecanismo es **asesor**: funciona porque cada sesión se compromete a consultarlo. Nada impide técnicamente escribir a la rama compartida sin reclamar, y eso DEBE quedar enunciado donde el mecanismo se documenta.

#### Scenario: nada lee un comando

- **GIVEN** la implementación del mecanismo
- **WHEN** se la examina
- **THEN** ninguna parte inspecciona, parsea ni clasifica un comando de shell
- **AND** ninguno de los manejadores que el producto registra en el host pertenece a este mecanismo

#### Scenario: una escritura encadenada no corre detrás de un reclamo rechazado

- **GIVEN** un reclamo encadenado por éxito delante de una escritura
- **WHEN** el reclamo termina con código distinto de cero
- **THEN** la escritura no llega a correr

#### Scenario: el carácter asesor está enunciado

- **GIVEN** la documentación del mecanismo
- **WHEN** se la lee
- **THEN** dice que funciona porque cada sesión se compromete a consultarlo, y que nada impide escribir sin reclamar

### Requisito: Un reclamo tiene exactamente tres resultados distinguibles

El reclamo DEBE producir uno de tres resultados, distinguibles por quien lo encadena:

- **El turno es tuyo** — termina en cero, sin nada en el flujo de error, y devuelve el recibo en su propio renglón.
- **Nadie arbitró** — termina en cero, así que el trabajo procede, pero DEBE dejar un mensaje en el flujo de error nombrando que no se arbitró y por qué. Ningún camino de no-arbitración DEBE ser silencioso ni indistinguible de un reclamo exitoso.
- **El turno está tomado** — termina distinto de cero, que es lo que detiene la escritura encadenada, y DEBE reportar qué cambio lo sostiene, desde qué árbol, en qué momento de escritura, hace cuánto, y qué hay en cola.

Ante un turno tomado, el reclamo DEBE registrar la solicitud en la cola, pausar brevemente y reintentar **exactamente una vez** dentro de la misma llamada. Si el reintento toma el turno, DEBE eliminar su propia entrada de cola —condicionada al valor que escribió— y devolver el recibo. El comando NO DEBE preguntar nada ni esperar a nadie.

#### Scenario: un turno libre se reclama y devuelve un recibo

- **GIVEN** ninguna sesión sostiene el turno
- **WHEN** una sesión lo reclama
- **THEN** el turno queda tomado, la llamada termina en cero y nada va al flujo de error
- **AND** el recibo llega en su propio renglón, sin ambigüedad

#### Scenario: reclamar deja el árbol de trabajo intacto

- **GIVEN** una sesión con árbol de trabajo limpio
- **WHEN** reclama el turno
- **THEN** no se produce ningún commit, rama ni cambio de archivo por el acto de reclamar
- **AND** el árbol se sigue reportando limpio

#### Scenario: el turno tomado se encola, pausa y reintenta una sola vez

- **GIVEN** el turno está tomado cuando una sesión lo reclama
- **WHEN** corre el reclamo
- **THEN** la solicitud pendiente queda registrada en la cola
- **AND** el reclamo se reintenta exactamente una vez, tras una pausa breve, antes de rendirse

#### Scenario: el reintento toma el turno y limpia su propia entrada

- **GIVEN** un reclamo cuyo primer intento falló y cuyo reintento tiene éxito
- **WHEN** la llamada retorna
- **THEN** elimina su entrada de cola, condicionada al valor que escribió
- **AND** termina en cero con el recibo del intento exitoso

#### Scenario: sigue tomado después del reintento

- **GIVEN** un reclamo cuyo reintento también encuentra el turno tomado
- **WHEN** la llamada retorna
- **THEN** termina con código distinto de cero
- **AND** su reporte nombra el cambio que lo sostiene, su árbol, su momento de escritura, hace cuánto, y qué hay en cola

#### Scenario: ningún camino de no-arbitración es silencioso

- **GIVEN** cada razón por la que el mecanismo puede no llegar a arbitrar, ejercida por separado
- **WHEN** se inspecciona cada resultado
- **THEN** cada uno termina en cero y lleva un mensaje que nombra la razón
- **AND** ninguno es indistinguible de un reclamo que tomó el turno

#### Scenario: un segundo reclamo dentro del mismo bracket es rechazado

- **GIVEN** quien ya sostiene el turno
- **WHEN** reclama otra vez
- **THEN** la llamada termina distinto de cero, igual que si el turno lo sostuviera cualquier otro
- **AND** el reporte nombra su propio cambio y árbol, lo que vuelve el error evidente para una persona, sin que el mecanismo haya detectado nada

#### Scenario: el comando reporta y quien lo llama conversa

- **GIVEN** un reclamo que terminó distinto de cero
- **WHEN** quien lo llamó lo maneja
- **THEN** pone los hechos reportados ante el usuario y no elige ninguna opción por su cuenta
- **AND** el comando mismo no pidió nada ni esperó por nadie

### Requisito: La liberación se decide sólo por el recibo

La liberación DEBE borrar el turno únicamente si éste todavía lleva exactamente el recibo que se le pasa. Sin recibo, DEBE fallar como error de uso, y NO DEBE existir ninguna ruta alternativa que identifique al llamador de otra forma ni que elija qué borrar comparando un campo registrado.

Sus dos resultados —**liberado** y **el turno se movió**— DEBEN terminar ambos en cero: una liberación es limpieza, no un veredicto sobre el trabajo. Ante un turno que ya no lleva el recibo, NO DEBE borrarse nada, y la discrepancia DEBE reportarse en el flujo de error en lugar de forzar.

Los campos del registro del turno —cambio, árbol, rama destino, momento, base y marca de tiempo UTC— DEBEN ser informativos: se imprimen en un reporte que lee una persona, y ninguno DEBE compararse, testearse ni usarse en ninguna decisión.

#### Scenario: el recibo coincide

- **GIVEN** quien sostiene el recibo que su reclamo devolvió, y un turno que aún lleva ese valor
- **WHEN** libera
- **THEN** el turno deja de existir y la llamada termina en cero

#### Scenario: el turno se movió por debajo

- **GIVEN** un turno que ya no lleva el valor del recibo, se haya movido o ya liberado
- **WHEN** quien lo sostenía libera
- **THEN** no se borra nada, la discrepancia se reporta en el flujo de error, y la llamada termina en cero igual

#### Scenario: una liberación sin recibo es error de uso

- **GIVEN** un llamador que no proporciona recibo
- **WHEN** llama a liberar
- **THEN** la llamada falla como error de uso
- **AND** no se borra nada y no se ofrece ninguna ruta alternativa de liberación

#### Scenario: todos los campos difieren y la liberación funciona igual

- **GIVEN** un turno cuyos campos de registro difieren todos de quien lo libera
- **WHEN** ese llamador libera con el recibo correcto
- **THEN** el turno se libera

#### Scenario: todos los campos coinciden y la liberación falla igual

- **GIVEN** un turno cuyos campos son idénticos a los de quien libera, pero cuyo valor no es el recibo que sostiene
- **WHEN** ese llamador libera
- **THEN** no se borra nada

#### Scenario: los campos se imprimen, nunca se consultan

- **GIVEN** un reporte producido para un turno tomado
- **WHEN** se examina el código que lo produjo
- **THEN** cada campo llegó al reporte como texto
- **AND** ninguno fue comparado, coincidido ni testeado en el camino

### Requisito: La consulta de estado es de sólo lectura y siempre tiene éxito

La consulta de estado DEBE poder invocarse sin reclamar nada, DEBE reportar el turno como libre o describir a quien lo sostiene, hace cuánto y qué hay en cola, DEBE terminar siempre en cero y NO DEBE modificar el turno ni escribir ninguna entrada de cola.

#### Scenario: el turno está libre

- **GIVEN** ningún turno sostenido
- **WHEN** se consulta el estado
- **THEN** se reporta libre y la llamada tiene éxito

#### Scenario: el turno está tomado, con cola detrás

- **GIVEN** un turno tomado con entradas en cola
- **WHEN** se consulta el estado
- **THEN** se reporta quién lo sostiene, hace cuánto y qué hay en cola, y la llamada tiene éxito

#### Scenario: consultar no cambia nada

- **GIVEN** cualquier estado del turno
- **WHEN** se consulta el estado
- **THEN** el turno queda exactamente como estaba y no se escribió ninguna entrada de cola

### Requisito: No hay liberación automática, y romper un turno ajeno es decisión humana

Una sesión que termina sosteniendo el turno DEBE dejarlo tomado: nada DEBE liberarlo en su nombre. Liberar un turno que uno no sostiene DEBE seguir siendo una decisión humana, incluso ante evidencia fuerte de que quien lo sostiene desapareció. Las entradas de cola abandonadas DEBEN acumularse tal cual, reportadas como solicitudes registradas y nunca como evidencia de que alguien sigue trabajando; nada DEBE barrerlas.

#### Scenario: una sesión muere sosteniendo el turno

- **GIVEN** una sesión que termina mientras sostiene el turno
- **WHEN** se inspecciona el turno después
- **THEN** sigue tomado y nada lo liberó en su nombre

#### Scenario: evidencia fuerte no autoriza romperlo

- **GIVEN** evidencia de que quien sostiene el turno desapareció
- **WHEN** quien espera considera liberarlo
- **THEN** no lo hace, y la decisión llega al usuario en su lugar

#### Scenario: las entradas de cola abandonadas se acumulan

- **GIVEN** un llamador que registró una solicitud pendiente y después se rindió
- **WHEN** se lee la cola más tarde
- **THEN** la entrada sigue ahí, reportada como solicitud registrada y no como trabajo en curso
- **AND** nada la barre

### Requisito: El run de consolidación reclama una vez por ronda, sin condición

El run de consolidación DEBE reclamar el turno **una sola vez, antes del primer cherry-pick de la ronda**, y liberarlo **una sola vez después del último, antes de la pasada de limpieza**. NO DEBE reclamar por cherry-pick, y DEBE liberar cualquiera sea el resultado del loop: una ronda que no integró nada libera igual.

El reclamo DEBE ser **incondicional**. La prosa que lo cablea NO DEBE contener ninguna prueba sobre qué contenedor tiene la ronda, ni ninguna rama en la que no se reclame: hay un solo destino de integración posible, la rama de trabajo, y por lo tanto una sola respuesta.

Ante un reclamo rechazado, la ronda NO DEBE integrar nada, DEBE reportarse bloqueada con los hechos que el reclamo reportó, y DEBE dejar intacto el trabajo que sus corridas aisladas ya produjeron. Quien realiza la escritura es quien reclama y quien libera.

#### Scenario: un batch de consolidación es un bracket

- **GIVEN** una ronda a punto de hacer cherry-pick de varios commits
- **WHEN** corre
- **THEN** reclama una vez antes del primer cherry-pick y libera una vez después del último, antes de la pasada de limpieza
- **AND** libera cualquiera sea el resultado del loop, incluida una ronda que no integró nada

#### Scenario: el reclamo del batch no tiene condición

- **GIVEN** la prosa que cablea el reclamo en la consolidación, después del cambio
- **WHEN** se la lee de punta a punta
- **THEN** ordena reclamar sin ninguna prueba previa sobre el contenedor de la ronda
- **AND** no queda ninguna rama enunciada en la que la ronda no reclame

#### Scenario: una ronda rechazada no integra y se reporta bloqueada

- **GIVEN** un run de consolidación cuyo reclamo terminó distinto de cero
- **WHEN** maneja el resultado
- **THEN** no integra nada esa ronda y se reporta bloqueado con los hechos que el reclamo reportó
- **AND** el trabajo ya producido por las corridas aisladas queda intacto para un intento posterior

#### Scenario: quien escribe es quien reclama

- **GIVEN** cualquier momento de escritura que necesite el turno
- **WHEN** se realiza
- **THEN** el turno lo reclama y lo libera quien realiza esa escritura, nunca otro participante en su nombre

### Requisito: Del origen no viajan la condición, el doble contenedor ni el esquema de citas ajeno

El cableado de la consolidación es un **empalme** sobre la prosa que este producto ya tiene, nunca una copia de la del origen. Tres cosas del origen NO DEBEN aparecer en ningún archivo embarcado después del cambio, porque cada una reintroduce lo que este cambio retira:

1. **La prueba condicional del origen** sobre si el contenedor de la ronda es la rama de trabajo. La oración se **borra**, no se inlinea con su respuesta.
2. **Toda formulación de contenedor doble** — la forma "el workspace de cambio cuando el aislamiento está activo, la rama de trabajo cuando no" y sus variantes. Cada aparición colapsa a un solo valor escrito como una sola cláusula, nunca como una oración de dos ramas con una tachada.
3. **El esquema de citas por ruta de importación del origen**, que es su convención de despliegue y no la de este producto. Toda cita DEBE quedar en la forma de ruta desplegada que este producto usa.

Además, ningún archivo embarcado DEBE referirse al workspace de nivel de cambio: ni por el nombre de su flag, ni por su directorio, ni por su patrón de rama. La única excepción es la línea de ignorados del repositorio, conservada a propósito.

#### Scenario: la condición del origen no viajó

- **GIVEN** el payload después del cambio
- **WHEN** se busca una prueba sobre cuál es el contenedor de la ronda antes de reclamar
- **THEN** no hay ninguna, y el reclamo aparece enunciado sin condición

#### Scenario: ninguna formulación de contenedor doble sobrevive

- **GIVEN** el payload después del cambio
- **WHEN** se lo busca por formulaciones que nombren dos contenedores posibles para una misma decisión
- **THEN** no hay ninguna: cada punto que antes ramificaba enuncia un solo valor en una sola cláusula
- **AND** ninguna quedó como oración de dos ramas con una de ellas tachada o marcada como inaplicable

#### Scenario: las citas quedan en el esquema de este producto

- **GIVEN** la skill portada y la prosa de consolidación después del cambio
- **WHEN** se examinan sus citas
- **THEN** todas están en la forma de ruta desplegada de este producto
- **AND** ninguna usa el esquema de citas por ruta de importación del origen

#### Scenario: nada embarcado nombra ya el workspace de nivel de cambio

- **GIVEN** el payload después del cambio
- **WHEN** se lo busca por el nombre del flag retirado, por el directorio del workspace y por su patrón de rama
- **THEN** no hay coincidencias en ningún archivo embarcado
- **AND** la línea de ignorados del repositorio queda excluida de esa búsqueda, conservada a propósito

#### Scenario: la cita sin destino se elimina en vez de re-apuntarse

- **GIVEN** la skill portada, que en el origen citaba dos documentos
- **WHEN** se la lee después del cambio
- **THEN** cita el documento que sí tiene destino en este producto
- **AND** la cita cuyo contenido se retira no quedó re-apuntada a ningún otro documento: se eliminó

### Requisito: El mecanismo tiene un solo hogar, invocable y sin colisión

La definición del mecanismo —incluido el invariante de un reclamo y una liberación por unidad de trabajo, y el criterio que decide si una escritura necesita turno— DEBE vivir en un único lugar descubrible, sin un segundo archivo que la repita. Quien necesite referirse a la mecánica de worktrees o de consolidación DEBE citar al documento dueño de ese comportamiento en lugar de dar una segunda versión.

El criterio DEBE aparecer en exactamente dos lugares de ese hogar: en su descripción compacta —de modo que alcance para decidir si el mecanismo aplica sin abrir el cuerpo— y completo en el cuerpo, con las consecuencias de cada respuesta. El caso *no* NO DEBE seguir enunciado como "aterrizar en la rama del propio workspace de cambio": ese contenedor ya no existe, y el caso *no* es escribir en el árbol en el que uno ya está.

El nombre de la carpeta de la skill ES su comando de invocación y DEBE ser único entre todas las skills instaladas por todos los dominios; una colisión futura DEBE hacer fallar el despliegue de forma ruidosa, nunca sobrescribir en silencio. El nombre declarado en el encabezado DEBE coincidir con la carpeta.

Ese hogar DEBE conservar además el procedimiento manual de git para reclamar y liberar, para una máquina sin la herramienta, marcado explícitamente como más débil: un turno reclamado a mano no lleva un recibo emitido por la herramienta.

#### Scenario: la definición existe una sola vez

- **GIVEN** el punto de entrada descubrible del mecanismo
- **WHEN** se lo lee
- **THEN** lleva la definición completa, incluido el invariante de bracket, enunciado exactamente una vez
- **AND** no existe un segundo archivo embarcado que la contenga

#### Scenario: la mecánica ajena se cita, no se absorbe

- **GIVEN** el hogar del mecanismo, que necesita referirse a comportamiento de worktree o de consolidación
- **WHEN** lo hace
- **THEN** cita al documento dueño de ese comportamiento
- **AND** no agrega una segunda versión de nada de eso

#### Scenario: la descripción decide sola si el mecanismo aplica

- **GIVEN** sólo el nombre y la descripción compacta de la skill, como se ven en un listado
- **WHEN** se contrasta una escritura contra ellos
- **THEN** la descripción enuncia el criterio en una forma que decide si el mecanismo hace falta, sin abrir el cuerpo
- **AND** el cuerpo lleva el criterio completo con las consecuencias de cada respuesta

#### Scenario: el caso *no* del criterio ya no nombra un contenedor retirado

- **GIVEN** el criterio, en la descripción compacta y en el cuerpo
- **WHEN** se lee su caso negativo
- **THEN** dice que escribir en el árbol en el que uno ya está no necesita turno
- **AND** no menciona aterrizar en la rama de un workspace de cambio

#### Scenario: el nombre invocable no colisiona con ninguna skill instalada

- **GIVEN** las carpetas de skills que embarcan todos los dominios instalados
- **WHEN** se comparan sus nombres
- **THEN** el del mecanismo aparece exactamente una vez
- **AND** el nombre declarado en el encabezado de la skill coincide con esa carpeta

#### Scenario: una colisión futura hace fallar el despliegue

- **GIVEN** dos dominios que embarcaran una skill con el mismo nombre de carpeta
- **WHEN** se despliega
- **THEN** el despliegue falla nombrando la colisión
- **AND** ninguna de las dos sobrescribe a la otra en silencio

#### Scenario: el procedimiento manual queda, marcado como más débil

- **GIVEN** el hogar del mecanismo, en su sección para una máquina sin la herramienta
- **WHEN** se la lee
- **THEN** lleva el procedimiento de git para reclamar y liberar a mano
- **AND** dice que un turno reclamado así no lleva un recibo emitido por la herramienta, y por qué eso es más débil

## Referencias

- **Conceptualmente relacionado**: `process/isolate-change-workspace.md` — el run de consolidación que sostiene este turno es el mismo que integra el aislamiento por tarea
