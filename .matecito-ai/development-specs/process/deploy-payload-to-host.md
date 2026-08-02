# Capability — Desplegar el payload en el host destino

- **Status:** Accepted
- **Date:** 2026-08-02
- **Components:** cli

## Propósito

Llevar lo que el ecosistema declara como fuente —los componentes de cada dominio activo, los compartidos y el archivo de instrucciones raíz— al directorio de configuración del host destino, resolviendo antes qué archivo es nuevo, cuál cambió, cuál ya está igual y cuál debe borrarse, para escribir solo lo necesario, respaldar lo que se pisa y limpiar lo obsoleto.

## Actores

- **Los flujos de instalación, actualización y sincronización interactiva** — este proceso no tiene punto de entrada propio: es un motor que ellos invocan como uno de sus componentes.

## Precondiciones

- Se puede resolver el payload de origen: se prefiere un payload local si existe en el directorio actual o en alguno de sus padres (modo desarrollo); si no, se usa el payload embebido en el ejecutable.
- Se puede resolver el directorio de configuración del host destino (hoy, `~/.claude/`).
- El payload contiene el núcleo del archivo de instrucciones raíz y un directorio de dominios; sin ellos la planificación falla.

## Flujo principal

1. El proceso resuelve el payload de origen y el directorio del host destino.
2. **Compone el archivo de instrucciones raíz**: el núcleo común del payload, seguido de un **índice generado** de los dominios activos. Ese índice se genera en memoria — no es un archivo del payload.
3. **Despliega el fragmento de cada dominio activo como archivo suelto**, en la ubicación que el índice le indica al agente que lea bajo demanda.
4. **Mapea los componentes de cada dominio activo** —agentes, skills y referencias— a su destino plano en el host. Un dominio que no trae alguno de esos componentes simplemente lo saltea.
5. **Suma los componentes compartidos** con el mismo mapeo, sin filtrarlos por dominio activo.
6. **Detecta clashes de destino**: si dos orígenes distintos resuelven al mismo archivo destino, la planificación falla sin escribir nada.
7. **Calcula el estado de cada archivo** contra lo que ya hay en el host: *nuevo*, *cambiado*, *igual* o *a borrar*.
8. Al aplicar, salta los *iguales*; para los *cambiados* respalda primero la copia existente en la carpeta de respaldo de la corrida (los *nuevos* no tienen qué respaldar); escribe el contenido tal cual, byte a byte, sin transformarlo; y para los *a borrar* respalda y elimina según su procedencia.

## Casos borde

- **Conjunto de dominios activos vacío** → se despliegan todos los dominios presentes en el payload (ver [`../rule/domain-activation-shim.md`](../rule/domain-activation-shim.md)), pero los componentes **compartidos se despliegan igual**, con o sin dominios activos: no están sujetos a esa puerta.
- **Un dominio inactivo** no aporta ni su fragmento, ni sus componentes, ni una fila en el índice generado; sus archivos **que sigan siendo registrados de una corrida anterior se marcan *a borrar***.
- **La capa de agrupación bajo las skills es organizativa y se descarta**: la ruta desplegada de una skill nunca incluye ni el dominio ni el grupo. Esto es lo que hace posible el clash entre dominios.
- **Archivos marcadores de directorio vacío** (los que existen solo para que el repositorio conserve una carpeta) nunca se despliegan; un archivo real que conviva con uno de ellos en el mismo subárbol sí se despliega.
- **Un dominio activo sin fragmento de instrucciones** no aparece en el índice generado ni deja un archivo suelto.
- **La carpeta de respaldo se crea perezosamente**: solo si llegó a haber al menos un archivo *cambiado* o *a borrar* que respaldar.
- **El cuerpo del fragmento de un dominio nunca se incrusta** en el archivo de instrucciones raíz: ahí solo va el núcleo más el índice.
- **Sin registro de deployments anteriores, la corrida es una migración**: no borra nada por ausencia de registro; solo se evalúa la lista embebida de huérfanos legacy conocidos para los que se tenga prueba de propiedad.
- **Un destino registrado —o del catálogo legacy— que ya no existe en disco no se planifica como *a borrar***: no aparece en el plan, no se cuenta en el desglose y no se lista entre los destinos a borrar. Al reescribirse el registro con lo que quedó vigente, la entrada desaparece sola. Anunciar un borrado que no va a ocurrir infla el preview justo en la corrida donde la persona decide si confía en que el proceso borre.
- **Un directorio queda vacío tras borrar un archivo** — el proceso elimina también los directorios que quedaron sin contenido, subiendo por sus padres hasta (pero nunca incluyendo) la raíz del componente bajo el directorio de configuración del host (`skills/`, `agents/`, `references/`, `scripts/`, `matecito-ai/`), que nunca se borra ni aunque quede completamente vacía. La limpieza aplica para borrados de las dos fuentes: registro y catálogo legacy.

## Reglas de negocio

- El archivo de instrucciones raíz es **generado**, no copiado: núcleo del payload + índice de dominios activos. Su contenido no existe como archivo en el payload.
- El índice lista los dominios en orden de directorio (determinista) y toma de cada uno su etiqueta y su resumen; si un dominio no los declara, se cae al identificador del dominio como etiqueta.
- **Convención de nombres única entre dominios:** dos dominios activos no pueden exponer una skill con el mismo nombre de carpeta. El sistema **detecta** el choque y falla; no lo resuelve automáticamente prefijando. Cuando los dos orígenes en conflicto son de dominios distintos, el mensaje nombra ambos dominios y, si el destino es una skill, también la skill — para que la violación de convención sea evidente.
- Un choque entre un componente compartido y uno de dominio también es un clash y también falla.
- Los archivos del payload se copian **byte a byte, sin ninguna transformación**.
- El respaldo se hace **solo** para archivos que ya existían y cambiaron, y para archivos que van a ser borrados; los nuevos no se respaldan y los iguales no se tocan.
- El estado *igual* se determina comparando el contenido a escribir con el contenido en disco; si el destino no se puede leer, el archivo cuenta como *nuevo*.
- El estado *a borrar* se clasifica en **dos fuentes** con criterios distintos:
  - **Entrada del registro de deployments**: se respalda y se borra **siempre**, aunque el contenido en disco haya cambiado desde que se escribió. La pertenencia al registro es la prueba de propiedad. Si el contenido no coincide con el hash registrado, se reporta la edición pisada para que la persona sepa qué se perdió.
  - **Entrada del catálogo embebido de huérfanos legacy**: se borra **solo si el contenido coincide** con uno de sus hashes históricos. Si no coincide, la entrada se preserva intacta y se reporta. El hash es la única prueba de propiedad cuando no existe registro.

## Escenarios

### Scenario: copia byte a byte

- **GIVEN** un payload con un archivo de agente y un host destino vacío
- **WHEN** el proceso planifica y aplica
- **THEN** el archivo en el host queda idéntico byte a byte al del payload

### Scenario: el archivo raíz es núcleo más índice, sin cuerpos de dominio

- **GIVEN** un payload con núcleo de instrucciones y un dominio que declara etiqueta, resumen y fragmento
- **WHEN** el proceso planifica y aplica
- **THEN** el archivo de instrucciones raíz contiene el núcleo y una entrada de índice con la etiqueta, el resumen y la ubicación del fragmento — y **no** contiene el cuerpo del fragmento

### Scenario: el fragmento del dominio se despliega suelto

- **GIVEN** el mismo payload del escenario anterior
- **WHEN** el proceso aplica
- **THEN** el cuerpo del fragmento queda como archivo suelto en el host, byte a byte igual al del payload

### Scenario: un conjunto activo explícito filtra dominios

- **GIVEN** un payload con dos dominios que traen fragmento y skills, y un conjunto activo que nombra solo a uno
- **WHEN** el proceso planifica y aplica
- **THEN** el índice lista solo el dominio activo, solo su fragmento queda suelto y solo su skill se despliega; nada del dominio inactivo llega al host

### Scenario: clash de skills entre dominios

- **GIVEN** dos dominios que exponen una skill con el mismo nombre de carpeta
- **WHEN** el proceso planifica
- **THEN** falla sin escribir nada, con un mensaje que nombra los dos dominios y la skill en conflicto

### Scenario: clash entre un compartido y un dominio

- **GIVEN** un componente compartido y uno de dominio que resuelven al mismo destino
- **WHEN** el proceso planifica
- **THEN** falla sin escribir nada

### Scenario: los compartidos se despliegan sin dominios activos

- **GIVEN** un payload con una skill compartida y un conjunto de dominios activos vacío
- **WHEN** el proceso planifica
- **THEN** la skill compartida está entre los archivos a desplegar

### Scenario: los compartidos conviven con los de dominio

- **GIVEN** un payload con una skill compartida y un agente de un dominio activo
- **WHEN** el proceso planifica
- **THEN** ambos están entre los archivos a desplegar

### Scenario: la skill compartida también se aplana

- **GIVEN** una skill compartida ubicada bajo una capa de agrupación
- **WHEN** el proceso planifica
- **THEN** su destino no incluye la capa de agrupación

### Scenario: los marcadores de directorio vacío no se despliegan

- **GIVEN** un subárbol del payload que contiene marcadores de directorio vacío y, en otro caso, un archivo real junto a uno de ellos
- **WHEN** el proceso planifica
- **THEN** ningún marcador aparece entre los archivos a desplegar, y el archivo real sí

### Scenario: fallback al payload embebido

- **GIVEN** un directorio de trabajo sin ningún payload local en él ni en sus padres
- **WHEN** el proceso resuelve el payload de origen
- **THEN** usa el payload embebido en el ejecutable, que contiene el núcleo de instrucciones

### Scenario: el registro refleja lo aplicado

- **GIVEN** una aplicación que escribió archivos y borró uno viejo
- **WHEN** termina
- **THEN** queda una entrada por archivo desplegado y ninguna del borrado

### Scenario: primera corrida con el binario nuevo

- **GIVEN** un host con archivos de versiones anteriores y sin registro
- **WHEN** corre la sincronización
- **THEN** se escribe el registro con los destinos del plan y **no se borra nada por ausencia de registro** — lo único que puede borrarse en esa corrida son los huérfanos del catálogo legacy con prueba de propiedad

### Scenario: un destino que ya no está en disco no se anuncia

- **GIVEN** un destino registrado que el plan de hoy no produce y que además ya no existe en el host
- **WHEN** el proceso planifica
- **THEN** no aparece entre los *a borrar* ni en el desglose, y al aplicar la entrada desaparece del registro sin error

### Scenario: la lista legacy se resuelve por hash

- **GIVEN** dos huérfanos legacy: uno cuyo contenido es alguno de sus hashes históricos, otro desplegado desde un working tree sucio y por eso sin hash histórico posible
- **WHEN** corre la migración
- **THEN** el primero se respalda y se borra; el segundo se preserva y entra al reporte de preservados

### Scenario: qué queda *a borrar*

- **GIVEN** un destino registrado que el plan de hoy no produce, y un dominio que dejó de estar activo con archivos en el registro
- **WHEN** el proceso planifica
- **THEN** ambos quedan *a borrar*, sin tratamiento especial

### Scenario: entrada del registro intacta

- **GIVEN** un *a borrar* del registro cuyo contenido coincide con su hash registrado
- **WHEN** el proceso aplica
- **THEN** lo respalda y lo elimina, y el reporte lo da como borrado limpio

### Scenario: entrada del registro que la persona editó

- **GIVEN** un *a borrar* del registro cuyo contenido NO coincide con su hash registrado
- **WHEN** el proceso aplica
- **THEN** lo respalda y lo elimina igual, y el reporte avisa que se pisó una edición y dónde quedó el respaldo

### Scenario: huérfano legacy que no coincide

- **GIVEN** un *a borrar* del catálogo legacy cuyo contenido no coincide con ninguno de sus hashes históricos
- **WHEN** el proceso aplica
- **THEN** no lo toca, lo deja en el host, y lo devuelve en el reporte de preservados; la corrida no falla

### Scenario: desactivar un dominio no deja carpetas fantasma

- **GIVEN** un dominio desactivado cuyos archivos ocupaban carpetas propias, algunas anidadas, dentro de `skills/`
- **WHEN** el proceso aplica los borrados
- **THEN** esas carpetas —incluidas las anidadas— dejan de existir, y `skills/` sigue existiendo

### Scenario: un directorio con algo ajeno queda intacto

- **GIVEN** un directorio del que se borra un archivo nuestro pero que conserva otro archivo que no es nuestro
- **WHEN** el proceso aplica el borrado
- **THEN** el directorio sigue existiendo, con lo ajeno adentro

### Scenario: la raíz del componente sobrevive vacía

- **GIVEN** un borrado que deja una raíz de componente sin ningún archivo adentro
- **WHEN** el proceso aplica
- **THEN** esa raíz sigue existiendo

## Referencias

- **Rule** → [`../rule/domain-activation-shim.md`](../rule/domain-activation-shim.md) — cómo se resuelve el conjunto de dominios activos y por qué el conjunto vacío significa "todos".
- **Flow** → [`../flow/install-ecosystem.md`](../flow/install-ecosystem.md) · [`../flow/update-ecosystem.md`](../flow/update-ecosystem.md) — los flujos que invocan este motor.
- **Contexto de repo** → [`../../../payload/docs/reference-deployed-paths.md`](../../../payload/docs/reference-deployed-paths.md) — el mapeo fuente → desplegado, y la regla de que ninguna instrucción de runtime nombre una ruta del payload.
