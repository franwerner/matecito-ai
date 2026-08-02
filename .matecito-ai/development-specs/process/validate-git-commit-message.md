# Capability — Validar el mensaje de un commit de git

- **Status:** Inferred
- **Date:** 2026-07-29
- **Components:** cli

## Propósito

Impedir que un commit salga con atribución de IA en su mensaje, y avisar —sin bloquear— cuando el mensaje no sigue Conventional Commits. La política de commits del ecosistema se hace cumplir en el momento en que el commit se intenta, no después.

## Actores

- **El host**, que intercepta el intento de commit antes de ejecutarlo y entrega la carga del evento al validador.
- **El agente** que intentó el commit, que recibe el bloqueo o el aviso.

## Precondiciones

- El validador se registra para el evento previo al uso de herramienta, acotado a los comandos de shell que son un `git commit`.

## Flujo principal

1. El host entrega la carga del evento con el comando de shell que está por ejecutarse.
2. El validador extrae de ese comando el mensaje del primer argumento `-m`.
3. **Regla de bloqueo:** si el mensaje contiene atribución de IA, responde bloqueando con código `2` y un mensaje que nombra el patrón encontrado y pide quitar la atribución.
4. **Regla de aviso:** si no hay atribución pero el mensaje no matchea Conventional Commits, responde con código `0` (no bloquea) y un aviso sugiriendo revisarlo.
5. Si el mensaje pasa ambas reglas, responde sin código de bloqueo y sin mensaje.

## Reglas de negocio

- **Bloqueo (código `2`)** cuando el mensaje contiene alguno de los patrones de atribución: la marca de co-autoría, el nombre del asistente, la frase de generación asistida o el emoji de robot. La comparación **no distingue mayúsculas**; el emoji se compara además sobre el texto original, porque es multi-byte.
- **Aviso no bloqueante (código `0` con mensaje)** cuando el mensaje no matchea Conventional Commits. La forma esperada es un tipo de la lista permitida, un scope opcional entre paréntesis en minúsculas, dos puntos y espacio, y un asunto que empieza en minúscula y no termina en punto.
- **La regla de bloqueo tiene precedencia sobre la de formato**: un mensaje que además de tener atribución tenga mal formato se bloquea, no se avisa.
- **Fail-open** — el validador nunca bloquea cuando no puede afirmar nada: carga vacía, carga que no se puede interpretar, carga sin comando, o comando sin argumento `-m`.
- El mensaje se toma del **primer** argumento `-m`; se aceptan valores entre comillas dobles, entre comillas simples o sin comillas (hasta el primer espacio). Un guion que no sea exactamente el argumento `-m` —por ejemplo, una opción larga— no se confunde con él.
- El aviso se escribe en la salida de error del validador; el código de salida es lo que el host interpreta como permitir (`0`) o bloquear (`2`).

## Casos borde

- **Commit interactivo** (sin `-m`) → fail-open: no se valida nada.
- **Commit con el mensaje en un archivo** (sin `-m`) → fail-open.
- **Carga malformada o vacía** → fail-open.
- **Carga sin el comando** → fail-open.
- **Enmienda sin `-m`** → fail-open.

## Errores de cara al actor

- **Atribución de IA en el mensaje** → se bloquea el commit con código `2` y un mensaje que empieza por `BLOCKED`, nombra el patrón detectado y pide quitar la atribución.
- **Mensaje fuera de Conventional Commits** → no se bloquea: se emite un mensaje que empieza por `WARN` sugiriendo revisarlo, con un ejemplo de la forma esperada.

## Escenarios

### Scenario: mensaje convencional válido pasa

- **GIVEN** un `git commit -m "feat(hooks): add validator"`
- **WHEN** el validador lo procesa
- **THEN** responde sin bloquear y sin mensaje

### Scenario: mensaje convencional sin scope pasa

- **GIVEN** un `git commit -m "fix: handle nil pointer"`
- **WHEN** el validador lo procesa
- **THEN** responde sin bloquear y sin mensaje

### Scenario: marca de co-autoría bloquea

- **GIVEN** un mensaje que incluye la marca de co-autoría del asistente
- **WHEN** el validador lo procesa
- **THEN** bloquea con código `2` y un mensaje que empieza por `BLOCKED`

### Scenario: nombre del asistente en el cuerpo bloquea

- **GIVEN** un mensaje que menciona el nombre del asistente en minúsculas dentro del asunto
- **WHEN** el validador lo procesa
- **THEN** bloquea con código `2` — la comparación no distingue mayúsculas

### Scenario: emoji de robot bloquea

- **GIVEN** un mensaje que contiene el emoji de robot
- **WHEN** el validador lo procesa
- **THEN** bloquea con código `2`

### Scenario: frase de generación asistida bloquea

- **GIVEN** un mensaje que contiene la frase de generación asistida
- **WHEN** el validador lo procesa
- **THEN** bloquea con código `2`

### Scenario: formato incorrecto solo avisa

- **GIVEN** un mensaje capitalizado y terminado en punto, sin atribución
- **WHEN** el validador lo procesa
- **THEN** responde con código `0` y un mensaje que empieza por `WARN` — el commit no se bloquea

### Scenario: commit sin `-m` no se valida

- **GIVEN** un `git commit` interactivo, o uno que toma el mensaje de un archivo
- **WHEN** el validador lo procesa
- **THEN** responde sin bloquear y sin mensaje

### Scenario: carga inválida no bloquea

- **GIVEN** una carga vacía, malformada, o sin el comando
- **WHEN** el validador la procesa
- **THEN** responde sin bloquear y sin mensaje

## Referencias

- **Rule** → [`../rule/hook-registry-domain-filtering.md`](../rule/hook-registry-domain-filtering.md) — cuándo un validador de este tipo se considera activo según los dominios habilitados.
- **Rule** → [`../rule/hook-reconciliation-by-matecito-id.md`](../rule/hook-reconciliation-by-matecito-id.md) — cómo se registra y se mantiene al día su handler en el archivo de settings del host.
