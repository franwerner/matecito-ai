# Capability — Estado de chequeo de un componente

- **Status:** Inferred
- **Date:** 2026-07-29
- **Components:** cli

## Propósito

Dar un vocabulario único para reportar cómo está cada pieza del entorno, de modo que el reporte se lea igual venga de donde venga —un binario, una integración, un archivo del host, un permiso o un handler— y que un problema siempre traiga cómo resolverlo.

## Entidades y estados

- **Resultado de chequeo** — el veredicto sobre una pieza del entorno en un momento dado. Lleva el nombre de la pieza, si es crítica, su estado, y —según el estado— la versión detectada, un detalle y una pista de remediación.

  Estados: **OK** · **Missing** · **Outdated**. Son excluyentes y se derivan de la observación en cada corrida: **no se persisten** ni transicionan en el tiempo. Volver a chequear vuelve a derivar el estado desde cero.

  - **OK** — la pieza está presente y satisface lo que se espera de ella. Puede portar la versión detectada.
  - **Missing** — la pieza no está, no se pudo alcanzar, o está pero no contiene lo que debía contener.
  - **Outdated** — la pieza está y se pudo leer, pero su versión no alcanza el mínimo requerido.

## Reglas de negocio

- **Un resultado `Missing` porta siempre una pista de remediación no vacía.** Es lo que vuelve accionable el reporte: no se le informa a la persona que falta algo sin decirle cómo resolverlo.
- **`Outdated` implica que la pieza fue encontrada**: no se puede estar desactualizado sin estar presente. Si la versión no se puede interpretar, el resultado también es `Outdated`, con el detalle explicando que no se pudo parsear.
- Una pieza que no está en el entorno, o cuya consulta de versión falla, es `Missing` — no `Outdated`.
- **Un resultado `Missing` de una pieza no crítica no hace fallar el reporte**: se muestra como advertencia y se marca como opcional. Solo las piezas **críticas** que no están en `OK` cuentan para el veredicto final.
- El veredicto global es binario: sin críticas fuera de `OK`, el estado es correcto; con al menos una, el reporte informa cuántos chequeos críticos faltan y el comando termina en fallo.
- Las secciones del reporte se muestran **solo para lo que corresponde a los dominios activos**; una sección sin piezas no se imprime.
- Un resultado `Missing` sin detalle se muestra con un detalle genérico, para que nunca quede una línea sin explicación.
- Las pistas de remediación se muestran solo cuando hay algo que remediar (`Missing` y `Outdated`); un `OK` no las muestra aunque las lleve.

## Escenarios

### Scenario: pieza presente y correcta

- **GIVEN** un binario presente en el entorno que responde a su consulta de versión
- **WHEN** el sistema lo chequea
- **THEN** el resultado queda en `OK` con la versión detectada

### Scenario: pieza ausente del entorno

- **GIVEN** un binario que no está en el entorno
- **WHEN** el sistema lo chequea
- **THEN** el resultado queda en `Missing`, con el detalle de que no se encontró y con una pista de remediación no vacía

### Scenario: consulta de versión que falla

- **GIVEN** un binario presente cuya consulta de versión falla
- **WHEN** el sistema lo chequea
- **THEN** el resultado queda en `Missing` —no en `Outdated`— con la salida del fallo como detalle y una pista de remediación

### Scenario: versión por debajo del mínimo

- **GIVEN** un binario presente cuya versión mayor está por debajo del mínimo requerido
- **WHEN** el sistema lo chequea
- **THEN** el resultado queda en `Outdated`, con el detalle de la versión encontrada y el mínimo, y con la pista de cómo actualizarlo

### Scenario: versión ilegible

- **GIVEN** un binario presente cuya versión no se puede interpretar
- **WHEN** el sistema lo chequea
- **THEN** el resultado queda en `Outdated` con el detalle de que no se pudo parsear

### Scenario: todo `Missing` trae remediación

- **GIVEN** cualquier conjunto de resultados de chequeo
- **WHEN** el sistema los produce
- **THEN** ninguno con estado `Missing` tiene la pista de remediación vacía

### Scenario: una pieza opcional faltante no rompe el veredicto

- **GIVEN** un resultado `Missing` de una pieza no crítica y ninguna pieza crítica fuera de `OK`
- **WHEN** el sistema resume el reporte
- **THEN** el veredicto global es correcto, y la pieza faltante se muestra como advertencia marcada como opcional

### Scenario: una pieza crítica faltante rompe el veredicto

- **GIVEN** al menos un resultado crítico que no está en `OK`
- **WHEN** el sistema resume el reporte
- **THEN** informa cuántos chequeos críticos faltan y el comando termina en fallo

### Scenario: los estados son excluyentes

- **GIVEN** cualquier resultado de chequeo
- **WHEN** se observa su estado
- **THEN** es exactamente uno de `OK`, `Missing` u `Outdated`

## Referencias

- **Rule** → [`../rule/mcp-permission-auto-approval.md`](../rule/mcp-permission-auto-approval.md) — el chequeo por patrón de auto-aprobación, que usa este vocabulario.
- **Rule** → [`../rule/hook-reconciliation-by-matecito-id.md`](../rule/hook-reconciliation-by-matecito-id.md) — el chequeo de presencia de cada handler declarado.
- **Rule** → [`../rule/domain-activation-shim.md`](../rule/domain-activation-shim.md) — qué secciones del reporte se muestran según los dominios activos.
