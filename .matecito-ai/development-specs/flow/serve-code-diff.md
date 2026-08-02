# Capability — Servir las vistas de diff de código

- **Status:** Accepted
- **Date:** 2026-07-27
- **Components:** api, ui

## Propósito

Servirle a la UI los diffs del código de un change —qué cambió cada batch, el neto del change y la comparación contra el archivo real— derivados on demand de la cadena de fotos, para que el cockpit muestre los cambios sin depender de git.

## Actores

- **La UI (el cockpit)** — pide las vistas de diff de un change, en nombre del usuario.

## Precondiciones

- El proyecto está registrado.
- El change existe (en cualquier estado: `active` o `closed` — el historial de diffs se sirve igual).
- Sin permisos ni auth (daemon local, localhost).
- Existen fotos capturadas para el change (archivos que ningún apply tocó no tienen diff que servir).

## Flujo principal

1. La UI pide una vista de diff de un change: de un archivo puntual o del conjunto de archivos tocados.
2. El sistema deriva la vista **on demand** comparando fotos de la cadena (o una foto contra el disco), según la vista pedida:
   - **Individual (por batch):** la foto del antes vs la foto del después de ese batch — muestra solo lo que ese batch cambió, lo previo se cancela por estar en ambos lados.
   - **Colapsada (neto del change):** la base del archivo en el change vs su última foto — todos los batches fundidos en un solo diff.
   - **Contra el archivo real:** la última foto vs el contenido actual en disco — evidencia el drift posterior al último apply.
3. Si entre la foto del después de un batch y la del antes del siguiente hay diferencia, el sistema la señala como bloque de **cambio fuera del flujo** (una edición que nadie fotografió mientras ocurría), delimitado entre esos dos puntos.
4. La UI muestra el diff con sus líneas cambiadas y contexto.

## Ramas / flujos alternativos

- **Archivo sin fotos en el change** (ningún apply lo tocó) → no hay diff que derivar para ese archivo; no es un error.
- **Archivo real inaccesible** (repo movido, checkout ausente) → la vista contra el archivo real no está disponible y se informa; las vistas derivadas de fotos (individual, colapsada) se sirven igual, porque son self-contained.
- **Archivo creado o borrado por el apply** → la ausencia en la base o en el después se deriva como archivo nuevo entero o como borrado.

## Reglas de negocio

- **Los diffs nunca se almacenan**: toda vista es derivada al momento, comparando dos estados (foto vs foto, o foto vs disco). Lo almacenado son solo las fotos.
- **Ninguna vista depende de git**: ni de commits, ni de ramas, ni de la historia del repo. La fuente es la cadena de fotos del change.
- La lectura no muta nada: servir diffs no altera fotos, eventos ni el estado del change.

## Errores de cara al actor

- **Change inexistente** → `change_not_found`.
- **Proyecto no registrado** → `project_not_registered`.
- La forma de error es `{error, code}`, sin internals, como en el resto del sistema.

## Escenarios

### Scenario: diff individual muestra solo lo del batch

- **GIVEN** un archivo tocado por los batches N y N+1
- **WHEN** la UI pide el diff individual del batch N+1
- **THEN** el diff se deriva entre el antes y el después de ese batch y muestra solo sus cambios; lo del batch N no aparece

### Scenario: diff colapsado muestra el neto

- **GIVEN** un archivo tocado por varios batches del change
- **WHEN** la UI pide la vista colapsada
- **THEN** el diff se deriva entre la base del archivo en el change y su última foto, con todos los batches fundidos en el resultado neto

### Scenario: drift contra el archivo real

- **GIVEN** un archivo editado a mano después del último apply que lo tocó
- **WHEN** la UI pide la vista contra el archivo real
- **THEN** el diff entre la última foto y el disco muestra exactamente el cambio manual

### Scenario: cambio fuera del flujo señalado

- **GIVEN** un archivo cuya foto del después del batch N difiere de la foto del antes del batch N+1
- **WHEN** la UI pide las vistas del change
- **THEN** esa diferencia se señala como bloque de cambio fuera del flujo, delimitado entre ambos batches

### Scenario: sin checkout siguen las vistas de fotos

- **GIVEN** un change cuyas fotos están en la base pero cuyo archivo real es inaccesible
- **WHEN** la UI pide las vistas
- **THEN** la vista contra el archivo real se informa como no disponible y las vistas individual y colapsada se sirven igual

### Scenario: archivo creado se deriva entero

- **GIVEN** un archivo cuya base en el change es la ausencia
- **WHEN** la UI pide cualquier vista que lo incluya
- **THEN** el diff lo muestra como archivo nuevo entero

### Scenario: change inexistente

- **GIVEN** un identificador de change que no existe
- **WHEN** la UI pide sus diffs
- **THEN** el sistema responde `change_not_found` con la forma `{error, code}`

## Referencias

- **Process** → [`../process/capture-code-snapshots.md`](../process/capture-code-snapshots.md) — la captura de las fotos de las que se derivan todas las vistas.
- **Flow** → [`serve-change-state.md`](serve-change-state.md) — el lado lectura general del change (snapshot + push) que estas vistas complementan.
- **EDR** → [`../../../apps/api/.matecito-ai/edr/contracts/api-contract.md`](../../../apps/api/.matecito-ai/edr/contracts/api-contract.md) — la superficie de lectura de la UI por la que se sirven estas vistas.
- **EDR** → [`../../../apps/api/.matecito-ai/edr/runtime/error-handling.md`](../../../apps/api/.matecito-ai/edr/runtime/error-handling.md) — la forma de error `{error, code}`.
