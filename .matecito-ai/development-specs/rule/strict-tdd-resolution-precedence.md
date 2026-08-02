# Capability — Precedencia al resolver el guard de TDD estricto

- **Status:** Inferred
- **Date:** 2026-07-29
- **Components:** cli

## Propósito

Fijar cómo se decide si el guard de TDD estricto está activo para un dominio, distinguiendo "no lo decidí" de "lo decidí en false" — porque solo esa distinción permite que un proyecto herede lo global en vez de pisarlo.

## Reglas de negocio

- La resolución es **por dominio**: el guard se pregunta siempre para un dominio concreto.
- **Precedencia:** valor del **proyecto** (solo si la clave está seteada) → valor **global** (solo si la clave está seteada) → **false**.
- **Clave ausente ≠ false.** Un config que no declara la clave no decide nada y deja pasar al siguiente nivel; un config que la declara en `false` **sí** decide, y corta la cadena.
- La ausencia del archivo de config del proyecto equivale, a efectos de precedencia, a un config de proyecto sin la clave: se pasa al global.
- Cuando ningún scope declara la clave, el guard queda **inactivo** (`false`). Es el default seguro: nada se endurece sin una decisión explícita.
- Un config con la forma antigua —la clave al tope, sin dominio— resuelve igual, porque se pliega al dominio por defecto antes de consultarse. Ver [`legacy-config-migration.md`](legacy-config-migration.md).

## Escenarios

### Scenario: el proyecto gana

- **GIVEN** un dominio con el guard en `false` a nivel global y en `true` a nivel proyecto
- **WHEN** el sistema lo resuelve
- **THEN** el guard queda activo

### Scenario: sin config de proyecto, gana el global

- **GIVEN** un dominio con el guard en `true` a nivel global y sin config de proyecto
- **WHEN** el sistema lo resuelve
- **THEN** el guard queda activo

### Scenario: clave ausente en el proyecto deja pasar al global

- **GIVEN** un config de proyecto que **no** declara la clave y un global que la declara en `true`
- **WHEN** el sistema lo resuelve
- **THEN** el guard queda activo — la ausencia en el proyecto no se interpreta como `false`

### Scenario: clave ausente en ambos scopes

- **GIVEN** ni el proyecto ni el global declaran la clave
- **WHEN** el sistema lo resuelve
- **THEN** el guard queda inactivo

### Scenario: sin ningún config

- **GIVEN** un config global vacío y sin config de proyecto
- **WHEN** el sistema lo resuelve
- **THEN** el guard queda inactivo

## Referencias

- **Rule** → [`agent-model-resolution-precedence.md`](agent-model-resolution-precedence.md) — la misma cadena de precedencia aplicada al modelo por agente.
- **Rule** → [`typed-flag-resolution-precedence.md`](typed-flag-resolution-precedence.md) — la distinción ausente/false generalizada al resto de los flags por dominio.
- **Rule** → [`legacy-config-migration.md`](legacy-config-migration.md) — cómo un config de la forma antigua llega a ser consultable por dominio.
