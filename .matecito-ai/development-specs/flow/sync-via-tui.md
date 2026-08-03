# Capability — Sincronizar desde la interfaz interactiva

- **Status:** Accepted
- **Date:** 2026-08-02
- **Components:** cli

## Propósito

Ofrecer la misma sincronización del ecosistema desde la interfaz interactiva: mostrar el plan con su desglose de borrados, pedir confirmación y transmitir el progreso en vivo sin que el plan reaparezca en el log, sin que la interfaz se re-ejecute nunca a sí misma. Cuando la sincronización encuentra un auto-reemplazo, difiere el despliegue del payload y la reconfiguración del ecosistema al arranque siguiente, de modo que se completen con el ejecutable ya actualizado.

## Actores

- **La persona** que entra a la pantalla de sincronización desde el menú de la interfaz interactiva.

## Precondiciones

- La pantalla recibe las opciones de sincronización ya armadas por la aplicación (entre ellas, la versión propia y, cuando existe, una detección de estado previa).

## Flujo principal

### Pantalla de sincronización

1. La persona entra a la pantalla de sincronización.
2. La pantalla resuelve el plan: si ya hay una detección de estado hecha, la **reutiliza**; si no, detecta.
3. Muestra el plan —las acciones no omitidas, numeradas, con su verbo y, para el componente de payload, su origen, desglose por categoría (*nuevos*, *cambiados*, *a borrar*) y destinos a borrar— y espera confirmación `y/n`.
4. Con la confirmación afirmativa, ejecuta la sincronización inyectando los estados ya detectados, de modo que no se vuelva a detectar, activando la señal de que el plan ya fue mostrado.
5. Transmite la salida de la ejecución línea a línea a la vista, mostrando las últimas líneas mientras corre. El log de ejecución no vuelve a listar el plan.
6. Al terminar sin reemplazo del propio ejecutable, vuelve automáticamente al menú.
7. Si encuentra un auto-reemplazo, deja registrada una marca de sincronización pendiente, queda en el aviso y sigue con el arranque siguiente.

### Arranque siguiente

8. Antes de abrir la interfaz interactiva, el sistema consulta la marca de sincronización pendiente.
9. Si la marca está puesta, corre la sincronización con el ejecutable ya actualizado (reutilizando el mismo motor, sin volver a hacer una nueva detección).
10. Si esa sincronización termina bien, limpia la marca y abre la interfaz.
11. Si falla, deja la marca puesta, avisa por la salida estándar y abre la interfaz igual — un fallo nunca impide que la interfaz abra, y el arranque siguiente vuelve a intentarla.

## Ramas / flujos alternativos

- **La persona cancela en la confirmación** → vuelve al menú sin ejecutar nada.
- **El plan resuelto no tiene ninguna acción activa** → la pantalla lo informa y no ofrece ejecutar.
- **La resolución del plan falla** → la pantalla queda en estado terminado mostrando el error.
- **El propio ejecutable fue reemplazado durante la ejecución** → la pantalla **no** vuelve al menú: se queda mostrando que el ejecutable se actualizó y que la sincronización se completará sola al volver a abrirlo, deja registrada una marca de sincronización pendiente.
- **La persona sale mientras la ejecución está en curso** → vuelve al menú; la ejecución en marcha no es interrumpible.
- **El arranque siguiente encuentra una marca de sincronización pendiente pero la sincronización falla** → avisa por la salida estándar, deja la marca puesta para el próximo arranque, y abre la interfaz igual.

## Casos borde

- **La interfaz interactiva nunca dispara el traspaso al binario nuevo.** Ante un auto-reemplazo, el único desenlace es quedarse en el aviso; no emite ninguna acción que pudiera relanzar el proceso. El traspaso es exclusivo de la línea de comandos (ver [`resume-self-replace-run.md`](resume-self-replace-run.md)).
- **La confirmación interna del motor se da por hecha**: la pantalla ya preguntó, así que ejecuta sin volver a preguntar.
- **La salida se muestra acotada**: solo se mantienen visibles las últimas líneas del progreso.
- **Al registrar la marca de sincronización pendiente, si el estado persistido existe pero es ilegible o corrupto, no se escribe nada**: se prefiere perder la marca antes que sobrescribir un archivo que no se pudo leer.
- **El fallo de la sincronización diferida nunca impide la apertura de la interfaz**: se avisa por salida estándar, la marca queda puesta y la interfaz abre igual.

## Reglas de negocio

- El plan que se muestra excluye las acciones *omitir*: solo se listan las que realmente se van a ejecutar.
- La detección de estado se hace **una sola vez** por entrada a la pantalla y se reutiliza al ejecutar.
- Confirmar es `y` o *enter*; cancelar es `n` o *esc*.
- Tras un auto-reemplazo, la pantalla queda en el aviso y no emite ninguna acción de navegación, pero deja registrada una marca de sincronización pendiente.
- El plan ya mostrado en la pantalla NO reaparece entre las líneas de ejecución que se transmiten: una señal intacta garantiza que el motor no lo repita.
- La sincronización diferida al arranque siguiente solo corre si la marca está puesta; es idempotente y puede reintentar sin riesgo.
- Cuando se registra la marca, si el estado es ilegible, el registro se omite antes de sobrescribir el archivo.

## Escenarios

### Scenario: auto-reemplazo deja la pantalla en el aviso de reinicio

- **GIVEN** una ejecución que terminó reportando que el propio ejecutable fue reemplazado
- **WHEN** la pantalla procesa el fin de la ejecución
- **THEN** queda marcada como terminada, muestra el aviso de que el ejecutable se actualizó y de que la sincronización se completará al volver a abrirlo, deja la marca de sincronización pendiente y **no** emite ninguna acción que relance el proceso

### Scenario: sin auto-reemplazo vuelve al menú

- **GIVEN** una ejecución que terminó sin reemplazo del propio ejecutable
- **WHEN** la pantalla procesa el fin de la ejecución
- **THEN** queda marcada como terminada, no deja ninguna marca de sincronización pendiente y navega de vuelta al menú

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

### Scenario: el auto-reemplazo marca la sincronización como pendiente

- **GIVEN** una ejecución de la interfaz que terminó reportando el reemplazo del propio ejecutable
- **WHEN** la pantalla procesa el fin de la ejecución
- **THEN** el estado persistido queda con la marca de sincronización pendiente puesta

### Scenario: sin auto-reemplazo no se marca nada

- **GIVEN** una ejecución de la interfaz que terminó sin reemplazo del propio ejecutable
- **WHEN** la pantalla procesa el fin de la ejecución
- **THEN** el estado persistido queda sin marca de sincronización pendiente

### Scenario: no hay estado persistido todavía

- **GIVEN** que aún no existe estado persistido
- **WHEN** se registra la marca de sincronización pendiente tras un auto-reemplazo
- **THEN** el estado se crea con la marca puesta

### Scenario: estado ilegible o corrupto — no se pisa

- **GIVEN** un estado persistido que existe pero no se puede leer ni interpretar
- **WHEN** se intenta registrar la marca de sincronización pendiente
- **THEN** no se escribe nada, la marca se pierde y el estado queda intacto
- AND la pantalla igual muestra su aviso

### Scenario: la marca dispara el sync antes de abrir la interfaz

- **GIVEN** un estado persistido con la marca de sincronización pendiente puesta
- **WHEN** se arranca la interfaz interactiva
- **THEN** la sincronización corre antes de que la interfaz se abra
- AND al terminar bien la marca queda limpia y la interfaz abre

### Scenario: sin marca no corre nada

- **GIVEN** un estado persistido sin marca de sincronización pendiente
- **WHEN** se arranca la interfaz interactiva
- **THEN** no se ejecuta ninguna sincronización y la interfaz abre directo

### Scenario: el sync diferido falla, avisa y se reintenta

- **GIVEN** un arranque con la marca de sincronización pendiente puesta cuya sincronización falla
- **WHEN** termina el intento
- **THEN** el fallo se avisa por la salida estándar, la marca queda puesta y la interfaz abre igual
- AND el arranque siguiente vuelve a intentarla

## Referencias

- **Flow** → [`resume-self-replace-run.md`](resume-self-replace-run.md) — por qué el traspaso al binario nuevo es exclusivo de la línea de comandos.
- **Flow** → [`update-ecosystem.md`](update-ecosystem.md) — el mismo motor de detección, plan y ejecución expuesto como comando.
