# Capability — Corrida relanzada tras el auto-reemplazo del ejecutable

- **Status:** Inferred
- **Date:** 2026-07-29
- **Components:** cli

## Propósito

Terminar sin fricción una sincronización que se interrumpió porque el propio ejecutable se reemplazó a sí mismo: el binario nuevo retoma la corrida donde quedó, sin volver a preguntar lo que la persona ya confirmó y sin poder entrar en un ciclo de auto-actualizaciones.

## Actores

- **La corrida relanzada** — el ejecutable nuevo, arrancado automáticamente por el traspaso posterior al auto-reemplazo. No hay una persona en el medio: la confirmación ya se dio en la invocación original.

## Precondiciones

- Una corrida anterior reemplazó con éxito el ejecutable propio y disparó el traspaso al binario nuevo.
- El traspaso marca la corrida nueva como *relanzada*; esa marca es lo único que la distingue de una invocación normal.

## Flujo principal

1. La corrida original detecta que el componente del propio ejecutable fue reemplazado con éxito.
2. Solo la superficie de línea de comandos dispara el traspaso: el motor de sincronización nunca se re-ejecuta a sí mismo, únicamente informa que hubo reemplazo.
3. El traspaso relanza el mismo comando con la marca de *relanzada* incorporada al entorno de la corrida nueva.
4. La corrida relanzada **excluye del plan la acción sobre el propio ejecutable**, aunque el estado detectado la volviera a proponer.
5. Fuerza la confirmación como ya dada: no lee la entrada estándar en ningún momento.
6. Omite la impresión del plan: la persona ya lo vio y lo aceptó en la invocación original.
7. Ejecuta las acciones restantes emitiendo igual el progreso por acción, para que la persona vea qué se está haciendo.

## Ramas / flujos alternativos

- **No hubo reemplazo del ejecutable** → no hay traspaso ni corrida relanzada: el flujo termina normalmente.
- **El traspaso al binario nuevo falla** → no se trata como error de la corrida: se informa por pantalla que el ejecutable ya se actualizó y que hay que volver a correrlo a mano. El ejecutable en disco ya está reemplazado, así que fallar duro solo dejaría a la persona sin salida.
- **La superficie de interfaz interactiva** nunca dispara el traspaso; se queda mostrando el aviso de reinicio. Ver [`sync-via-tui.md`](sync-via-tui.md).

## Casos borde

- **Garantía de terminación:** excluir la acción sobre el propio ejecutable evita un segundo traspaso aunque la consulta de la versión más reciente falle o devuelva un valor inconsistente. Una corrida relanzada nunca vuelve a auto-actualizarse.
- **La marca de relanzada ya presente en el entorno** no se duplica al propagarse.
- **Marca con un valor distinto del esperado** → la corrida se considera normal, no relanzada.
- **Sin ninguna acción restante** además de la propia excluida → la corrida informa que no hay nada para hacer y termina.

## Reglas de negocio

- El motor de sincronización es puro respecto del traspaso: nunca se re-ejecuta a sí mismo; solo reporta que hubo auto-reemplazo y deja que el llamador decida.
- En una corrida relanzada, la confirmación se considera dada: la entrada estándar **no se lee**, ni siquiera cuando no se pasó la opción de omitir confirmación.
- En una corrida relanzada no se imprime el plan, pero **sí** el progreso por acción.
- La acción sobre el propio ejecutable se excluye del plan activo de la corrida relanzada, y por lo tanto esa corrida nunca reporta un auto-reemplazo propio.

## Escenarios

### Scenario: la acción propia se excluye

- **GIVEN** una corrida relanzada cuyo estado detectado propondría instalar el propio ejecutable y también otro componente
- **WHEN** el sistema arma el plan activo
- **THEN** la acción del propio ejecutable queda fuera, no se emite progreso para ella y la corrida no reporta auto-reemplazo

### Scenario: no se pide confirmación ni se lee la entrada

- **GIVEN** una corrida relanzada sin la opción de omitir confirmación
- **WHEN** el sistema va a ejecutar el plan
- **THEN** ejecuta sin preguntar y no lee la entrada estándar en ningún momento

### Scenario: no se imprime el plan pero sí el progreso

- **GIVEN** una corrida relanzada con al menos una acción distinta de la propia
- **WHEN** el sistema la ejecuta
- **THEN** no aparece la impresión del plan, pero sí aparece el progreso de esa acción

### Scenario: control — corrida normal sí muestra plan y pregunta

- **GIVEN** una corrida **no** relanzada con una acción activa y sin la opción de omitir confirmación
- **WHEN** la persona responde que no en el prompt
- **THEN** el sistema había impreso el plan y la corrida queda cancelada sin ejecutar nada

### Scenario: el traspaso falla y se degrada al aviso manual

- **GIVEN** una corrida donde el ejecutable se reemplazó con éxito pero el traspaso al binario nuevo falla
- **WHEN** el comando termina
- **THEN** no devuelve error: informa que el ejecutable se actualizó y que hay que volver a correrlo a mano

### Scenario: sin reemplazo no hay traspaso

- **GIVEN** una corrida donde el ejecutable propio no fue reemplazado
- **WHEN** el comando termina
- **THEN** no se dispara ningún traspaso ni se emite ningún aviso al respecto

## Referencias

- **Flow** → [`install-ecosystem.md`](install-ecosystem.md) · [`update-ecosystem.md`](update-ecosystem.md) — los dos comandos que disparan el traspaso.
- **Flow** → [`sync-via-tui.md`](sync-via-tui.md) — la superficie interactiva, que nunca dispara el traspaso y se queda en el aviso de reinicio.
