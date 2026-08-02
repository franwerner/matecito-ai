# Capability — Sincronizar desde la interfaz interactiva

- **Status:** Accepted
- **Date:** 2026-08-02
- **Components:** cli

## Propósito

Ofrecer la misma sincronización del ecosistema desde la interfaz interactiva: mostrar el plan con su desglose de borrados, pedir confirmación y transmitir el progreso en vivo sin que el plan reaparezca en el log, sin que la interfaz se re-ejecute nunca a sí misma.

## Actores

- **La persona** que entra a la pantalla de sincronización desde el menú de la interfaz interactiva.

## Precondiciones

- La pantalla recibe las opciones de sincronización ya armadas por la aplicación (entre ellas, la versión propia y, cuando existe, una detección de estado previa).

## Flujo principal

1. La persona entra a la pantalla de sincronización.
2. La pantalla resuelve el plan: si ya hay una detección de estado hecha, la **reutiliza**; si no, detecta.
3. Muestra el plan —las acciones no omitidas, numeradas, con su verbo y, para el componente de payload, su origen, desglose por categoría (*nuevos*, *cambiados*, *a borrar*) y destinos a borrar— y espera confirmación `y/n`.
4. Con la confirmación afirmativa, ejecuta la sincronización inyectando los estados ya detectados, de modo que no se vuelva a detectar, activando la señal de que el plan ya fue mostrado.
5. Transmite la salida de la ejecución línea a línea a la vista, mostrando las últimas líneas mientras corre. El log de ejecución no vuelve a listar el plan.
6. Al terminar sin reemplazo del propio ejecutable, vuelve automáticamente al menú.

## Ramas / flujos alternativos

- **La persona cancela en la confirmación** → vuelve al menú sin ejecutar nada.
- **El plan resuelto no tiene ninguna acción activa** → la pantalla lo informa y no ofrece ejecutar.
- **La resolución del plan falla** → la pantalla queda en estado terminado mostrando el error.
- **El propio ejecutable fue reemplazado durante la ejecución** → la pantalla **no** vuelve al menú: se queda mostrando que el ejecutable se actualizó y que hay que volver a correrlo para usar la versión nueva.
- **La persona sale mientras la ejecución está en curso** → vuelve al menú; la ejecución en marcha no es interrumpible.

## Casos borde

- **La interfaz interactiva nunca dispara el traspaso al binario nuevo.** Ante un auto-reemplazo, el único desenlace es quedarse en el aviso de reinicio manual; no emite ninguna acción que pudiera relanzar el proceso. El traspaso es exclusivo de la línea de comandos (ver [`resume-self-replace-run.md`](resume-self-replace-run.md)).
- **La confirmación interna del motor se da por hecha**: la pantalla ya preguntó, así que ejecuta sin volver a preguntar.
- **La salida se muestra acotada**: solo se mantienen visibles las últimas líneas del progreso.

## Reglas de negocio

- El plan que se muestra excluye las acciones *omitir*: solo se listan las que realmente se van a ejecutar.
- La detección de estado se hace **una sola vez** por entrada a la pantalla y se reutiliza al ejecutar.
- Confirmar es `y` o *enter*; cancelar es `n` o *esc*.
- Tras un auto-reemplazo, la pantalla queda en el aviso de reinicio y no emite ninguna acción de navegación.
- El plan ya mostrado en la pantalla NO reaparece entre las líneas de ejecución que se transmiten: una señal intacta garantiza que el motor no lo repita.

## Escenarios

### Scenario: auto-reemplazo deja la pantalla en el aviso de reinicio

- **GIVEN** una ejecución que terminó reportando que el propio ejecutable fue reemplazado
- **WHEN** la pantalla procesa el fin de la ejecución
- **THEN** queda marcada como terminada, muestra el aviso de reinicio manual y **no** emite ninguna acción que pueda relanzar el proceso

### Scenario: sin auto-reemplazo vuelve al menú

- **GIVEN** una ejecución que terminó sin reemplazo del propio ejecutable
- **WHEN** la pantalla procesa el fin de la ejecución
- **THEN** queda marcada como terminada y navega de vuelta al menú

### Scenario: se reutiliza la detección previa

- **GIVEN** una entrada a la pantalla con una detección de estado ya disponible
- **WHEN** la pantalla resuelve el plan
- **THEN** usa esos estados en lugar de volver a detectar, y los inyecta también en la ejecución

### Scenario: cancelar en la confirmación

- **GIVEN** la pantalla mostrando el plan y esperando confirmación
- **WHEN** la persona responde `n` o presiona *esc*
- **THEN** vuelve al menú sin ejecutar ninguna acción

### Scenario: plan vacío

- **GIVEN** un plan resuelto sin ninguna acción activa
- **WHEN** la pantalla lo presenta
- **THEN** informa que no hay nada para hacer en lugar de ofrecer la ejecución

### Scenario: progreso en vivo

- **GIVEN** una ejecución confirmada y en curso
- **WHEN** la sincronización emite salida
- **THEN** la pantalla la va mostrando línea a línea, conservando visibles las últimas

### Scenario: el log de ejecución no repite el plan

- **GIVEN** una pantalla que mostró el plan y la persona confirmó
- **WHEN** la ejecución transmite su salida
- **THEN** esas líneas no vuelven a listar el plan

### Scenario: la pantalla anticipa los borrados

- **GIVEN** un plan con destinos a borrar
- **WHEN** la pantalla lo presenta
- **THEN** muestra el desglose y los destinos a borrar antes del prompt de confirmación

## Referencias

- **Flow** → [`resume-self-replace-run.md`](resume-self-replace-run.md) — por qué el traspaso al binario nuevo es exclusivo de la línea de comandos.
- **Flow** → [`update-ecosystem.md`](update-ecosystem.md) — el mismo motor de detección, plan y ejecución expuesto como comando.
