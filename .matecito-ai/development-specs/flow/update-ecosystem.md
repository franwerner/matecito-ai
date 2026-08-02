# Capability — Actualizar los componentes del ecosistema

- **Status:** Accepted
- **Date:** 2026-08-02
- **Components:** cli

## Propósito

Reconciliar en un solo comando el estado de los componentes ya instalados —binarios, payload desplegado en el host destino y configuración del host— contra lo que el ecosistema declara hoy, anticipando qué destinos se van a borrar e informando al final si alguno quedó sin sincronizar.

## Actores

- **La persona** que corre `matecito-ai update` desde una terminal.

## Precondiciones

- Se puede resolver la carpeta de respaldo de la corrida; si no, el comando falla antes de tocar nada.

## Flujo principal

1. La persona corre `matecito-ai update`.
2. El sistema detecta el estado observado de cada componente gestionado y deriva el plan (*instalar* / *actualizar* / *omitir*).
3. Si ninguna acción queda activa, informa que no hay nada para hacer y termina.
4. Muestra el plan numerado con el verbo por componente y, para el componente de payload, su origen, desglose por categoría (*nuevos*, *cambiados*, *a borrar*) y destinos a borrar.
5. Pide confirmación (`y/N`) salvo que se haya pasado `--yes`.
6. Ejecuta las acciones en orden, continuando ante el fallo de cualquiera y mostrando el resultado de cada una.
7. Durante la ejecución del componente de payload, informa aquellos destinos que se han preservado por ser legacy sin coincidencia de hash, con la acción sugerida, y aquellos cuyas ediciones previas fueron pisadas al borrar, informando dónde quedó el respaldo. Ninguno de estos cuenta como componente fallado.
8. Si el propio ejecutable fue reemplazado, entrega la ejecución al binario nuevo (ver [`resume-self-replace-run.md`](resume-self-replace-run.md)).
9. Termina en error **si y solo si al menos un componente falló**; si no, termina OK.

## Ramas / flujos alternativos

- **`--dry-run`** → se muestra el plan y el comando termina sin ejecutar nada.
- **`--yes`** → se omite la confirmación interactiva.
- **Respuesta no afirmativa en la confirmación** → se informa la cancelación y no se ejecuta ninguna acción.
- **Corrida relanzada tras un auto-reemplazo** → se omiten plan y confirmación. Ver [`resume-self-replace-run.md`](resume-self-replace-run.md).

## Casos borde

- **La reconciliación de la configuración del host va última**: se ejecuta después de los binarios y del despliegue del payload, para que al registrarse las integraciones y ajustarse los permisos el entorno ya esté completo. Esto es lo que hace que sumar una integración nueva quede consistente con `update` y no solo con `install`.
- **No se pudo consultar la versión más reciente de un componente** (por ejemplo, sin red) → ese componente se omite en vez de forzar una reinstalación.
- **Un componente falla** → se registra su error, la corrida sigue con los demás, y el comando termina en error al final.

## Reglas de negocio

- El plan se deriva del estado observado con esta precedencia: versión más reciente desconocida → *omitir*; componente ausente → *instalar*; payload cambiado → *actualizar*; reconciliación pendiente → *actualizar*; sin versión de referencia para comparar → *omitir*; versión actual distinta de la más reciente (normalizando el prefijo `v`) → *actualizar*; en cualquier otro caso → *omitir*.
- La ejecución es **continue-on-error**: el fallo de un componente nunca cancela los demás.
- El resultado de salida es binario: hubo al menos un componente fallado, o no hubo ninguno. El comando no distingue *cuál* falló en su código de salida.
- La confirmación acepta como afirmativas `y`, `yes`, `s`, `si` y `sí`, sin distinguir mayúsculas; cualquier otra respuesta cancela.

## Errores de cara al actor

- **No se pudo resolver la carpeta de respaldo** → el comando falla antes de detectar o ejecutar nada.
- **Al menos un componente falló** → el comando termina en error indicando que uno o más componentes fallaron durante la sincronización; el detalle por componente ya se mostró durante la ejecución.

## Escenarios

### Scenario: todo al día

- **GIVEN** un entorno donde ningún componente está ausente, desactualizado ni pendiente
- **WHEN** la persona corre `matecito-ai update`
- **THEN** el sistema informa que no hay nada para hacer y termina OK

### Scenario: versión remota desconocida se omite

- **GIVEN** un componente presente cuya versión más reciente no se pudo consultar
- **WHEN** el sistema deriva el plan
- **THEN** ese componente queda *omitido*, sin forzar reinstalación

### Scenario: diferencia de versión solo por el prefijo

- **GIVEN** un componente cuya versión instalada y cuya versión más reciente difieren únicamente en el prefijo `v`
- **WHEN** el sistema deriva el plan
- **THEN** las considera equivalentes y omite el componente

### Scenario: un componente falla → salida en error

- **GIVEN** un plan con varias acciones donde una falla
- **WHEN** la corrida termina
- **THEN** el sistema ejecutó igualmente las acciones restantes y el comando termina en error indicando que uno o más componentes fallaron

### Scenario: ningún componente falla → salida OK

- **GIVEN** un plan cuyas acciones se ejecutaron todas sin error
- **WHEN** la corrida termina
- **THEN** el comando termina OK

### Scenario: la configuración del host se reconcilia al final

- **GIVEN** un plan que incluye binarios, payload y configuración del host
- **WHEN** el sistema ejecuta las acciones
- **THEN** la configuración del host se reconcilia después de los binarios y del despliegue del payload

### Scenario: ni preservar ni pisar es fallar

- **GIVEN** una corrida cuyo payload preservó un huérfano legacy y borró una entrada del registro que estaba editada
- **WHEN** termina
- **THEN** mostró el desglose y los destinos a borrar, informa el preservado con su ruta **y su acción sugerida** y la edición pisada con su respaldo, y el comando termina OK

## Referencias

- **Flow** → [`install-ecosystem.md`](install-ecosystem.md) — el comando hermano, con la misma regla de salida y el mismo motor, más su propio preview del plan.
- **Flow** → [`resume-self-replace-run.md`](resume-self-replace-run.md) — la corrida relanzada tras el auto-reemplazo del ejecutable.
- **Rule** → [`../rule/mcp-permission-auto-approval.md`](../rule/mcp-permission-auto-approval.md) — qué patrones de auto-aprobación se reconcilian como parte de la configuración del host.
