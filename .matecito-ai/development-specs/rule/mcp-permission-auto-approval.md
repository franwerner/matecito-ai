# Capability — Auto-aprobación de las herramientas del ecosistema

- **Status:** Inferred
- **Date:** 2026-07-29
- **Components:** cli

## Propósito

Que la persona no tenga que aprobar a mano, una y otra vez, las herramientas que el propio ecosistema instaló: los patrones de auto-aprobación se derivan de lo que declaran los dominios activos, no de una lista escrita a mano.

## Reglas de negocio

- **Los patrones se derivan, no se enumeran.** Se produce un patrón por cada integración que declaran los dominios activos. No hay integraciones "de base": nada se registra ni se auto-aprueba si ningún dominio activo lo declara.
- **Convención de nombre por defecto.** El patrón de una integración se deriva de su nombre siguiendo la convención de patrón de herramientas del host (hoy, el prefijo de herramientas de integración seguido del nombre y un comodín). Un nombre que no está en el registro interno igual obtiene su patrón por convención.
- **Override explícito para las que rompen la convención.** Una integración cuyo nombre expuesto por el host no coincide con la convención declara su patrón explícitamente; ese override gana. Hoy el único caso es la integración de memoria persistente, que el host expone bajo su nombre de plugin y no bajo el nombre de la integración.
- **Un patrón fijo se incluye siempre**, además de los derivados: el que habilita la invocación de skills. No depende de ningún dominio.
- **Conjunto de seguridad ante error de resolución.** Si los dominios activos no se pueden resolver, se cae a un conjunto conocido en lugar de quedarse sin patrones — una configuración rota nunca deja el entorno sin auto-aprobación en silencio.
- **La reconciliación solo agrega.** Los patrones faltantes se suman a la lista de permitidos del archivo de settings del host; los que ya estaban se conservan y **nada se quita**, incluidos los que puso la persona.
- La reconciliación **no toca** el modo de permisos por defecto ni los permisos de las herramientas de shell, escritura o edición: se limita a la lista de permitidos.
- El chequeo del entorno reporta **un resultado por patrón esperado**, indicando si está o no en la lista de permitidos.

## Escenarios

### Scenario: patrón por convención

- **GIVEN** una integración cuyo nombre sigue la convención del host
- **WHEN** el sistema deriva su patrón de auto-aprobación
- **THEN** obtiene el patrón que la convención dicta para ese nombre

### Scenario: override para la integración que rompe la convención

- **GIVEN** la integración de memoria persistente, que el host expone bajo su nombre de plugin
- **WHEN** el sistema deriva su patrón
- **THEN** obtiene el patrón declarado explícitamente, no el que dictaría la convención sobre su nombre

### Scenario: nombre desconocido cae a la convención

- **GIVEN** un nombre de integración que no figura en el registro interno
- **WHEN** el sistema deriva su patrón
- **THEN** obtiene el patrón por convención sobre ese nombre, sin fallar

### Scenario: el patrón fijo siempre está

- **GIVEN** cualquier conjunto de dominios activos
- **WHEN** el sistema arma los patrones esperados
- **THEN** el patrón que habilita la invocación de skills está incluido, además de los derivados

### Scenario: la reconciliación solo agrega

- **GIVEN** una lista de permitidos que ya contiene un patrón esperado y le falta otro
- **WHEN** el sistema reconcilia
- **THEN** agrega el faltante, conserva el que estaba y no quita nada

### Scenario: nada que reconciliar

- **GIVEN** una lista de permitidos que ya contiene todos los patrones esperados
- **WHEN** el sistema reconcilia
- **THEN** no reporta cambios y no modifica el archivo

### Scenario: chequeo por patrón

- **GIVEN** un conjunto de patrones esperados donde uno falta en la lista de permitidos
- **WHEN** el sistema chequea el entorno
- **THEN** reporta ese patrón como faltante, con la indicación de cómo remediarlo, y el resto como presentes

## Referencias

- **Rule** → [`domain-activation-shim.md`](domain-activation-shim.md) — de dónde sale el conjunto de integraciones declaradas.
- **Flow** → [`../flow/update-ecosystem.md`](../flow/update-ecosystem.md) — el componente de configuración del host donde esta reconciliación se aplica.
- **Lifecycle** → [`../lifecycle/component-check-status.md`](../lifecycle/component-check-status.md) — la forma del resultado que reporta cada patrón.
