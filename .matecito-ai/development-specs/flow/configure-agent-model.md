# Capability — Elegir el modelo de cada agente de un dominio

- **Status:** Inferred
- **Date:** 2026-07-29
- **Components:** cli

## Propósito

Dejar que la persona fije, por dominio y por agente, con qué modelo corre cada agente — distinguiendo lo que decide para todas sus máquinas de lo que decide solo para un proyecto, y sin escribir overrides que no aportan nada.

## Actores

- **La persona** que abre la pantalla de modelos por agente desde la interfaz interactiva.

## Precondiciones

- Hay un dominio activo seleccionado; sus agentes son las filas de la pantalla.
- Hay un scope activo: **global** (vale para todas las máquinas de la persona) o **proyecto** (vale solo para el repositorio actual).

## Flujo principal

1. La persona abre la pantalla para un dominio.
2. El sistema siembra una fila por agente del dominio, partiendo de los overrides explícitos que ya tenga el config del scope activo.
3. Para los agentes sin override explícito, completa el valor según el scope (ver *Reglas de negocio*).
4. La persona recorre las filas y cambia el modelo de la fila actual, ciclando entre las opciones o eligiendo una directamente por su posición.
5. Al salir guardando, el sistema persiste los valores editados en el config del scope activo, descartando los que no corresponde persistir.
6. Vuelve a la pantalla anterior.

## Ramas / flujos alternativos

- **Salir descartando** → se abandona la pantalla sin escribir nada en el config.
- **El guardado falla** → la pantalla **no** se cierra: muestra el error para que la persona lo vea, en lugar de volver al menú en silencio.

## Casos borde

- **Un agente sin override en scope proyecto** queda marcado con el valor especial de herencia, no con un modelo concreto: no se inventa una decisión de proyecto que la persona no tomó.
- **Un agente sin override en scope global** se siembra con el modelo que el propio payload declara por defecto para ese agente; solo si el payload no declara ninguno se cae a un valor de último recurso.
- **El valor especial de herencia solo existe durante la edición**: nunca se persiste como valor.
- **Al mostrar una fila heredada**, el sistema resuelve y muestra cuál es el modelo efectivo: el override global si lo hay, y si no el default del payload para ese agente.
- **Un agente que ya tenía override en el config del scope** conserva su valor al abrir la pantalla, aunque coincida o no con el default.

## Reglas de negocio

- El valor especial de herencia se ofrece **solo en scope proyecto**, y es la primera opción del ciclo: es la forma de *quitar* el override de proyecto de un agente.
- **Al guardar en scope proyecto**, cada agente marcado con el valor de herencia **limpia** su override en el config del proyecto.
- **Al guardar en scope global**, un valor que coincide con el default que el payload declara para ese agente **no se persiste**: escribirlo contaminaría la herencia de todos los proyectos, que dejarían de ver el default del payload y pasarían a ver un override global espurio.
- Los valores seleccionables son los del **conjunto de modelos válidos que declara el host**; en scope proyecto, ese conjunto se antepone con el valor de herencia.
- La herencia que se muestra nunca es un valor fijo en el código: se resuelve contra el config global y, en su defecto, contra el default declarado en el payload para ese agente.
- El dominio determina qué agentes se listan y bajo qué entrada del config se persisten: la elección de modelos es **por dominio**, no global al ecosistema.

## Escenarios

### Scenario: en proyecto, sin override propio, hereda el global

- **GIVEN** un dominio cuyo agente tiene un override en el config global y ninguno en el del proyecto, y un default de payload distinto
- **WHEN** la persona abre la pantalla en scope proyecto
- **THEN** la fila queda marcada con el valor de herencia y muestra como efectivo el modelo del override global, no el default del payload

### Scenario: en proyecto, sin override en ningún lado, hereda el default del payload

- **GIVEN** un agente sin override ni en el proyecto ni en el global
- **WHEN** la persona abre la pantalla en scope proyecto
- **THEN** la fila muestra como efectivo el default que el payload declara para ese agente

### Scenario: en global, se siembra con el default del payload

- **GIVEN** un dominio con dos agentes, uno con override global y otro sin ninguno
- **WHEN** la persona abre la pantalla en scope global
- **THEN** el primero conserva su override y el segundo se siembra con su default del payload — no con un modelo fijo elegido por la herramienta

### Scenario: al guardar en global no se persiste lo igual al default

- **GIVEN** dos agentes en scope global, uno cuyo valor coincide con su default de payload y otro cuyo valor difiere
- **WHEN** la persona sale guardando
- **THEN** el config guardado no contiene entrada para el que coincide con su default, y sí contiene la del que difiere

### Scenario: el valor de herencia limpia el override de proyecto

- **GIVEN** un agente con override en el config del proyecto
- **WHEN** la persona lo cambia al valor de herencia y sale guardando
- **THEN** el config del proyecto queda sin override para ese agente

### Scenario: el valor de herencia no se ofrece en scope global

- **GIVEN** la pantalla abierta en scope global
- **WHEN** la persona cicla las opciones de una fila
- **THEN** solo aparecen modelos del conjunto válido; el valor de herencia no está entre ellos

### Scenario: fallo al guardar mantiene la pantalla abierta

- **GIVEN** una edición pendiente y un config que no se puede escribir
- **WHEN** la persona sale guardando
- **THEN** la pantalla sigue abierta mostrando el error del guardado, sin volver al menú

## Referencias

- **Rule** → [`../rule/agent-model-resolution-precedence.md`](../rule/agent-model-resolution-precedence.md) — la precedencia proyecto → global → default con la que se resuelve el modelo efectivo.
- **Rule** → [`../rule/model-value-validation.md`](../rule/model-value-validation.md) — qué valores de modelo se aceptan y por qué los nombres de agente no están restringidos.
