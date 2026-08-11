# EDR — Handshake de base en dos niveles contra un espacio de trabajo contaminado

- **Status:** Accepted
- **Date:** 2026-08-10

## Contexto

El aislamiento por espacio de trabajo que ofrece el harness no garantiza que el directorio entregado
parta realmente del estado actual de la rama de trabajo — un directorio reciclado o contaminado es una
base que nadie verificó. Confiar ciegamente en la garantía de aislamiento del harness deja abierta la
posibilidad de implementar sobre una base desconocida sin que nada lo note hasta que el conflicto (o
algo peor, un merge silencioso incorrecto) aparece más tarde.

## Decisión

Dos chequeos, uno por lado, cierran esa brecha:

- **Nivel 1 (corrida aislada, antes de escribir nada):** `git rev-parse HEAD` debe coincidir con la
  base que el orquestador capturó y pasó en el prompt de despacho, y `git status --porcelain` debe
  estar vacío. Si cualquiera de los dos falla, la corrida no implementa nada y reporta el motivo.
- **Nivel 2 (corrida de consolidación, antes de confiar en un commit reportado):** el padre del commit
  reportado debe coincidir con esa misma base.

## Reglas verificables

- **[auto]** Antes de escribir, la corrida aislada verifica `git rev-parse HEAD` contra la base recibida y `git status --porcelain` vacío; si cualquiera falla, no implementa nada y reporta `not-implemented / base-not-established`.
- **[auto]** La corrida de consolidación verifica que el padre del commit reportado coincide con la misma base antes de cherry-pickearlo.

## Alternativas consideradas

- **Confiar en la garantía de aislamiento del harness.** Descartado: es una promesa externa, no una verificación — la única defensa real contra un espacio de trabajo contaminado es un chequeo explícito hecho por la corrida misma, en los dos extremos del flujo (antes de escribir, y antes de integrar).

## Consecuencias

- Una base no establecida nunca se implementa sobre datos desconocidos: el peor resultado posible es una tarea reportada incompleta, nunca un commit sobre una base equivocada.
- El chequeo de nivel 2 da a la consolidación una segunda oportunidad de descartar un reporte cuya base resultó no confiable, sin depender de que el nivel 1 haya funcionado como se esperaba.
- Agrega dos verificaciones de git a cada corrida y a la consolidación — costo bajo, pagado una vez por tarea.

## Relacionados

- `relacionado-con` → [one-commit-per-isolated-run.md](one-commit-per-isolated-run.md) — el commit que este handshake protege en sus dos extremos.
- `relacionado-con` → [../structure/dispatch-batch-bound-integration.md](../structure/dispatch-batch-bound-integration.md) — el nivel 2 corre durante la integración que esa decisión ordena diferir.
