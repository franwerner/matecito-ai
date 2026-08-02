# Capability — Filtrado de hooks por dominio activo

- **Status:** Inferred
- **Date:** 2026-07-29
- **Components:** cli

## Propósito

Decidir qué hooks del ecosistema corresponden a una máquina según los dominios que tenga habilitados, distinguiendo los que pertenecen a un dominio de los que valen siempre.

## Reglas de negocio

- Cada hook declara **a qué dominio pertenece**, junto con su evento, su matcher y su implementación. La declaración y la implementación viven juntas.
- **Un hook de un dominio se incluye solo si ese dominio está activo.** Si el dominio no está entre los activos, el hook no se incluye ni se registra.
- **Existe un dominio compartido**, que no es un dominio de trabajo sino un sello de "vale siempre": un hook marcado así **se incluye para cualquier conjunto activo, incluido el conjunto vacío**.
- **Identidad estable:** cuando un hook no declara su marcador de identidad, se deriva de su dominio y su nombre de handler. Ese marcador es el que después permite reconciliarlo sin tocar lo ajeno.
- **El conjunto activo se resuelve del entorno** con la misma regla que el resto del ecosistema (ver [`domain-activation-shim.md`](domain-activation-shim.md)); si esa resolución falla, la resolución de hooks activos también falla en lugar de adivinar.
- **Ningún dominio activo con hooks** significa que la sección de hooks no existe: no se reconcilia ni se reporta nada.

## Escenarios

### Scenario: el hook compartido entra con conjunto activo vacío

- **GIVEN** un hook registrado bajo el dominio compartido
- **WHEN** el sistema filtra los hooks con un conjunto de dominios activos vacío
- **THEN** ese hook está incluido

### Scenario: el hook compartido convive con el de un dominio activo

- **GIVEN** un hook del dominio compartido y otro de un dominio concreto
- **WHEN** el sistema filtra con ese dominio activo
- **THEN** ambos están incluidos

### Scenario: un hook de un dominio inactivo se excluye

- **GIVEN** un hook de un dominio que no es el compartido
- **WHEN** el sistema filtra con un conjunto activo vacío, y también con un conjunto que no incluye ese dominio
- **THEN** en los dos casos el hook queda fuera

### Scenario: sin hooks activos no hay sección

- **GIVEN** un entorno donde ningún dominio activo declara hooks
- **WHEN** el sistema chequea el estado del entorno
- **THEN** no reporta ningún resultado de hooks

## Referencias

- **Rule** → [`domain-activation-shim.md`](domain-activation-shim.md) — cómo se resuelve el conjunto de dominios activos.
- **Rule** → [`hook-reconciliation-by-matecito-id.md`](hook-reconciliation-by-matecito-id.md) — qué se hace con los hooks que resultan incluidos.
- **Process** → [`../process/validate-git-commit-message.md`](../process/validate-git-commit-message.md) — un hook de dominio concreto y su comportamiento.
