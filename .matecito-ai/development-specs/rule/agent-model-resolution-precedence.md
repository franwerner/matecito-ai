# Capability — Precedencia al resolver el modelo de un agente

- **Status:** Inferred
- **Date:** 2026-07-29
- **Components:** cli

## Propósito

Fijar de dónde sale el modelo con el que corre un agente de un dominio, y qué significa que nadie lo haya fijado — para que la decisión más específica gane y la ausencia de decisión no se sustituya por una inventada.

## Reglas de negocio

- La resolución es **por dominio y por agente**: se pregunta siempre por el par (dominio, agente), nunca por el agente suelto.
- **Precedencia:** override del **proyecto** → override **global** → **default**.
- Un override cuenta solo si tiene valor. Una entrada vacía es equivalente a no tener entrada y deja pasar al siguiente nivel.
- **Ausencia en ambos scopes no es un valor:** la resolución devuelve el marcador *default* y **ningún modelo**. El consumidor debe **omitir el parámetro de modelo** en la invocación, para que rija el default que el propio agente declara. Nunca se sustituye por el modelo de la conversación en curso ni por un valor fijo elegido por la herramienta.
- La resolución informa además **de dónde salió** el valor (proyecto, global o default), para que la interfaz pueda mostrar el origen sin volver a resolver.
- **Un agente desconocido resuelve siempre a default**: preguntar por un agente que ningún scope declara no es un error, devuelve el mismo resultado que la ausencia.
- La ausencia del archivo de config del proyecto es indistinguible, a efectos de precedencia, de un config de proyecto sin override para ese agente: en ambos casos se pasa al global.
- Un config con la forma antigua —overrides al tope, sin dominio— resuelve igual, porque se pliega al dominio por defecto antes de consultarse. Ver [`legacy-config-migration.md`](legacy-config-migration.md).

## Escenarios

### Scenario: el proyecto gana

- **GIVEN** un agente con override global y override de proyecto, distintos entre sí
- **WHEN** el sistema resuelve su modelo
- **THEN** devuelve el del proyecto e informa que el origen es el proyecto

### Scenario: sin config de proyecto, gana el global

- **GIVEN** un agente con override global y sin config de proyecto
- **WHEN** el sistema resuelve su modelo
- **THEN** devuelve el global e informa que el origen es el global

### Scenario: sin override en ningún scope

- **GIVEN** un agente sin override ni en el proyecto ni en el global
- **WHEN** el sistema resuelve su modelo
- **THEN** no devuelve modelo e informa que el origen es el default — el consumidor omite el parámetro

### Scenario: agente desconocido

- **GIVEN** configs que declaran overrides para otros agentes
- **WHEN** el sistema resuelve el modelo de un agente que ninguno declara
- **THEN** no devuelve modelo e informa que el origen es el default

### Scenario: ambos configs ausentes

- **GIVEN** ni config global ni config de proyecto
- **WHEN** el sistema resuelve el modelo de cualquier agente
- **THEN** no devuelve modelo e informa que el origen es el default

### Scenario: config con la forma antigua

- **GIVEN** un config global con overrides al tope, sin dominio
- **WHEN** el sistema lo normaliza y resuelve el modelo de uno de esos agentes en el dominio por defecto
- **THEN** devuelve el override e informa que el origen es el global

## Referencias

- **Rule** → [`strict-tdd-resolution-precedence.md`](strict-tdd-resolution-precedence.md) — la misma precedencia aplicada al guard de un dominio.
- **Rule** → [`legacy-config-migration.md`](legacy-config-migration.md) — cómo un config de la forma antigua llega a ser consultable con esta precedencia.
- **Flow** → [`../flow/configure-agent-model.md`](../flow/configure-agent-model.md) — la pantalla donde la persona fija estos overrides, y por qué no persiste los que igualan al default.
