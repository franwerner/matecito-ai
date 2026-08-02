# Capability — Flags tipados por dominio y su precedencia

- **Status:** Inferred
- **Date:** 2026-07-29
- **Components:** cli

## Propósito

Definir cómo se guardan y se leen los flags de comportamiento que un dominio declara (por ejemplo, los que habilitan la minería automática de decisiones o de comportamiento), de modo que "no lo decidí" siga siendo distinguible de "lo decidí en false" y la decisión de un dominio no se filtre a otro.

## Reglas de negocio

- **Tri-estado obligatorio.** Un flag tipado tiene tres estados observables: **sin valor** (la clave no está), **false explícito** y **true explícito**. La ausencia **no** se colapsa a `false` al guardarse ni al leerse.
- **Aislamiento por dominio.** Los flags se guardan bajo el dominio al que pertenecen. Fijar un flag en un dominio **no** cambia lo que se lee para otro: un dominio sin ese flag sigue leyendo *sin valor*.
- **Precedencia:** valor del **proyecto** (si la clave está seteada) → valor **global** (si la clave está seteada) → **false**. Es la misma cadena que la del guard de TDD estricto: el tri-estado es lo que la hace posible.
- **La puerta es la intención, no la presencia del store.** Un flag en `true` habilita su comportamiento aunque el store que ese comportamiento alimenta todavía no exista; un flag en `false` o sin valor lo apaga por completo, sin dejar rastro.
- **Los flags que un dominio ofrece son los que su declaración de dominio enumera.** Cada uno se muestra y se persiste con el tipo que declara; los que el sistema conoce de forma tipada se guardan como tales, y el resto se guardan como ajustes genéricos del dominio — de modo que un dominio pueda sumar configuración sin cambiar el modelo de config.
- **Compatibilidad hacia atrás:** un flag guardado con la forma antigua —al tope, sin dominio— se pliega al dominio por defecto al leerse, y deja de existir en la forma antigua. Ver [`legacy-config-migration.md`](legacy-config-migration.md).

## Escenarios

### Scenario: sin valor es distinto de false

- **GIVEN** un config sin el flag declarado para un dominio
- **WHEN** el sistema lo lee
- **THEN** obtiene *sin valor* — no `false`

### Scenario: false explícito se conserva

- **GIVEN** un config donde el flag se fijó explícitamente en `false` para un dominio
- **WHEN** el sistema lo lee
- **THEN** obtiene `false` explícito, distinguible de la ausencia

### Scenario: true explícito se conserva

- **GIVEN** un config donde el flag se fijó en `true` para un dominio
- **WHEN** el sistema lo lee
- **THEN** obtiene `true`

### Scenario: un dominio no contamina a otro

- **GIVEN** un config con el flag en `true` para un dominio
- **WHEN** el sistema lee ese mismo flag para otro dominio
- **THEN** obtiene *sin valor*

### Scenario: la ausencia en el proyecto deja pasar al global

- **GIVEN** un config de proyecto sin el flag y un global con el flag en `true`
- **WHEN** el sistema resuelve el flag efectivo
- **THEN** queda en `true` — la ausencia en el proyecto no se lee como `false`

### Scenario: un flag guardado con la forma antigua se pliega al dominio por defecto

- **GIVEN** un config que declara el flag al tope, sin dominio
- **WHEN** el sistema lo carga
- **THEN** el flag queda disponible bajo el dominio por defecto y la clave al tope deja de existir

## Referencias

- **Rule** → [`strict-tdd-resolution-precedence.md`](strict-tdd-resolution-precedence.md) — la misma cadena de precedencia, sobre el guard del dominio.
- **Rule** → [`legacy-config-migration.md`](legacy-config-migration.md) — el plegado de las claves al tope hacia el dominio por defecto.
- **Rule** → [`domain-activation-shim.md`](domain-activation-shim.md) — de dónde sale la lista de dominios cuya configuración se ofrece.
