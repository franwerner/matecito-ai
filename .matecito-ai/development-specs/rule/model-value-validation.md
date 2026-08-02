# Capability — Validación del valor de modelo configurado

- **Status:** Inferred
- **Date:** 2026-07-29
- **Components:** cli

## Propósito

Rechazar temprano un modelo que el host no reconoce, sin restringir de paso qué agentes pueden configurarse — porque los agentes los aportan los dominios y cambian, mientras que los modelos válidos los declara el host.

## Reglas de negocio

- **El valor sí se valida.** Un override cuyo modelo no pertenece al **conjunto de modelos válidos que declara el host** hace fallar la validación, con un error que nombra el valor rechazado, el agente al que pertenecía y el conjunto aceptado. (Hoy ese conjunto son los alias `fable`, `opus`, `sonnet` y `haiku`.)
- **El nombre del agente NO se valida.** Los agentes de cada dominio se **descubren del payload**, no de una lista fija: validarlos contra un roster fijo rechazaría agentes nuevos o renombrados. Cualquier nombre de agente es aceptable como clave.
- La validación recorre **todos los dominios** presentes en la configuración, no solo el dominio por defecto.
- La validación **normaliza primero**: un config con la forma anterior se pliega al dominio por defecto antes de validarse, así que sus overrides también se validan. Ver [`legacy-config-migration.md`](legacy-config-migration.md).
- **El guard del dominio no está restringido**: cualquier valor booleano —incluida la ausencia— es válido.
- **Una configuración ausente es válida**: validar "nada" no es un error.
- **La validez del alias no se prueba contra el host.** El sistema no averigua qué alias soporta la instalación concreta del host; solo comprueba pertenencia al conjunto declarado. Un alias válido pero no soportado por esa instalación se degrada, en el momento de invocar al agente, al default que el agente declara — no acá.

## Escenarios

### Scenario: valor fuera del conjunto válido

- **GIVEN** una configuración con un override cuyo modelo no pertenece al conjunto declarado por el host
- **WHEN** el sistema la valida
- **THEN** falla con un error que menciona el valor rechazado

### Scenario: valor dentro del conjunto válido

- **GIVEN** una configuración con un override cuyo modelo pertenece al conjunto y con el guard del dominio fijado
- **WHEN** el sistema la valida
- **THEN** pasa sin error

### Scenario: el nombre del agente no se restringe

- **GIVEN** una configuración con un override para un agente que no está en ninguna lista fija, con un modelo válido
- **WHEN** el sistema la valida
- **THEN** pasa sin error

### Scenario: configuración ausente

- **GIVEN** ninguna configuración
- **WHEN** el sistema la valida
- **THEN** pasa sin error y sin fallar

### Scenario: se valida también lo que venía en la forma anterior

- **GIVEN** una configuración con overrides al tope, sin dominio, uno de ellos con un modelo inválido
- **WHEN** el sistema la valida
- **THEN** los pliega al dominio por defecto y falla por el modelo inválido

## Referencias

- **Rule** → [`agent-model-resolution-precedence.md`](agent-model-resolution-precedence.md) — cómo se elige el override efectivo una vez validado.
- **Rule** → [`legacy-config-migration.md`](legacy-config-migration.md) — el plegado previo a la validación.
- **Flow** → [`../flow/configure-agent-model.md`](../flow/configure-agent-model.md) — la pantalla que solo ofrece valores del conjunto declarado.
