# Capability — Migración de la configuración heredada

- **Status:** Inferred
- **Date:** 2026-07-29
- **Components:** cli

## Propósito

Que una configuración escrita con una forma anterior siga funcionando y se ponga al día sola, sin que la persona tenga que editar nada a mano ni perder lo que ya había decidido.

## Reglas de negocio

Hay **dos migraciones distintas**, que ocurren en momentos distintos:

**1. Plegado de las claves al tope hacia el dominio por defecto (en cada lectura).**

- Un config que declara ajustes **al tope** —sin dominio— es de la forma anterior al modelo por dominio. Al cargarse, esos ajustes se pliegan **al dominio por defecto** y las claves al tope se limpian.
- **El plegado nunca pisa un valor presente:** si el dominio por defecto ya tiene ese ajuste, el valor heredado del tope se descarta. Lo explícito y moderno gana.
- El plegado es **idempotente**: aplicarlo de nuevo sobre un config ya plegado no cambia nada.
- Tras el plegado, todo consumidor ve **solo** la forma por dominio, y todo lo que se escriba de vuelta se escribe **solo** en esa forma. La forma anterior no se reescribe nunca.

**2. Migración del archivo de modelos heredado (una sola vez).**

- Cuando el archivo de configuración no existe pero **sí** existe, en el mismo directorio, el archivo heredado que solo contenía los modelos por agente, el sistema lo migra: arma la configuración con esos modelos, la escribe como archivo de configuración y **borra** el archivo heredado.
- La migración deja el guard de TDD estricto escrito como **`false` explícito**, no como clave ausente: la migración toma la decisión y la deja registrada.
- **Es de una sola vez.** Después de migrar, el archivo heredado ya no existe; la carga siguiente lee el archivo de configuración con normalidad y no vuelve a migrar nada.
- **Si el archivo de configuración existe, el heredado se ignora** y **no** se borra: no hay migración cuando ya hay configuración.
- Si no existe ninguno de los dos, no hay configuración y no hay error: la ausencia es un estado válido.
- Un archivo de configuración presente pero ilegible **sí** es un error; no se degrada a "sin configuración".

## Escenarios

### Scenario: los modelos al tope se pliegan al dominio por defecto

- **GIVEN** un config que declara overrides de modelo al tope, sin dominio
- **WHEN** el sistema lo carga
- **THEN** esos overrides quedan disponibles bajo el dominio por defecto y las claves al tope quedan limpias

### Scenario: un flag al tope se pliega al dominio por defecto

- **GIVEN** un config que declara un flag de comportamiento al tope
- **WHEN** el sistema lo carga
- **THEN** el flag queda disponible bajo el dominio por defecto y la clave al tope queda limpia

### Scenario: el plegado no pisa lo ya presente

- **GIVEN** un config que declara un ajuste al tope y también, para el dominio por defecto, ese mismo ajuste con otro valor
- **WHEN** el sistema lo carga
- **THEN** conserva el valor del dominio y descarta el del tope

### Scenario: migración del archivo de modelos heredado

- **GIVEN** un directorio sin archivo de configuración pero con el archivo de modelos heredado
- **WHEN** el sistema carga la configuración
- **THEN** obtiene los modelos bajo el dominio por defecto, con el guard de TDD estricto en `false` explícito; el archivo de configuración queda creado y el heredado, borrado

### Scenario: la migración no se repite

- **GIVEN** un directorio donde ya ocurrió la migración
- **WHEN** el sistema vuelve a cargar la configuración
- **THEN** lee el archivo de configuración, obtiene los mismos modelos y no recrea el archivo heredado

### Scenario: con configuración presente, el heredado se ignora

- **GIVEN** un directorio con archivo de configuración y también con el archivo de modelos heredado, con valores distintos
- **WHEN** el sistema carga la configuración
- **THEN** usa la del archivo de configuración y **no** borra el heredado

### Scenario: sin ninguno de los dos archivos

- **GIVEN** un directorio sin archivo de configuración ni heredado
- **WHEN** el sistema carga la configuración
- **THEN** no hay configuración y no hay error

### Scenario: archivo de configuración ilegible

- **GIVEN** un archivo de configuración con contenido inválido
- **WHEN** el sistema lo carga
- **THEN** falla con error y no devuelve configuración

## Referencias

- **Rule** → [`agent-model-resolution-precedence.md`](agent-model-resolution-precedence.md) · [`strict-tdd-resolution-precedence.md`](strict-tdd-resolution-precedence.md) · [`typed-flag-resolution-precedence.md`](typed-flag-resolution-precedence.md) — las tres resoluciones que dependen de que la configuración ya esté plegada a la forma por dominio.
