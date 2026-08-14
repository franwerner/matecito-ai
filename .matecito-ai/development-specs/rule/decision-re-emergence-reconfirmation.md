# Capability — Reconfirmación de decisiones que resurgen

- **Status:** Accepted
- **Date:** 2026-08-14
- **Components:** cli

## Propósito

Una decisión ratificada en un gate puede reaparecer en un gate posterior del mismo cambio — su premisa cayó, o una rama diferente la trae de nuevo. Sin un registro, resurge como item nuevo. Con el registro, se reconoce en una sola pregunta de sí/no: ¿sigue valiendo lo que ratificamos? Sólo a las decisiones re-identificadas se les aplica este flujo más corto; las que tienen contenido diferente se presentan en su forma completa como items nuevos.

## Actores

- **Orquestador**: escribe en `sdd/{change-name}/ratified-decisions` al cerrar un gate de ratificación, registrando cada item que el usuario decidió; en un gate posterior lee ese registro para reconocer items que regresan
- **Usuario**: ratifica decisiones en un gate y las ve aparecer en forma corta en gates posteriores si resurgen sin cambio
- **Fase**: suministra items ratificables; ignora el registro de ratificaciones

## Precondiciones

- Un primer gate de ratificación (gate de decisiones pendientes, INTAKE GATE, o gate de minería) ha cerrado en el cambio
- Cada item ratificado en ese gate lleva un `record:` token que identifica la decisión de forma única (un slug de dominio)
- Un gate posterior en el mismo cambio presenta un item cuyo `record` coincide con uno registrado

## Flujo principal

1. Orquestador cierra un gate de ratificación con items decididos
2. Por cada item ratificado, escribir una fila a `sdd/{change-name}/ratified-decisions`: `record` (el token identificador), `ratified_summary` (el contenido después de ajustes del usuario), `anchor` (la fuente), `gate` (cuál gate lo ratificó)
3. En un gate posterior, consultar el registro de ratificaciones
4. Para cada item que el siguiente gate presente:
   - Si su `record` coincide exactamente con uno registrado: comparar también el `summary` entrante contra `ratified_summary`
   - Si el summary es idéntico: presentar una sola pregunta de sí/no ("¿Sigue valiendo X?"), mostrando la fuente del registro
   - Si el summary difiere: presentar el item en su forma completa, como item nuevo sin invocar la reconfirmación corta
   - Si no hay coincidencia en `record`: presentar como item nuevo
5. Registrar el resultado de la reconfirmación (sí = ratificado nuevamente, no = rechazado, reabierto)

## Ramas / flujos alternativos

- **El usuario rechaza una reconfirmación corta** → Se abre el item en su forma completa; puede ser ratificado nuevamente o rechazado
- **No hay registro previo** → Todos los items se presentan en forma completa (el primer gate de este cambio)
- **Record coincide pero summary cambió** → Presentación completa, no reconfirmación corta (el usuario no vio este nuevo contenido antes)
- **No hay items a ratificar** → No se consulta el registro; el gate continúa sin ello

## Casos borde

- **Un item que resurge idéntico se ratifica de nuevo sin reconfirmación** → Tal cual vio, con resultado sí
- **Un item resurge con `record` igual pero premisa diferente** → Se presenta en forma completa (usuario nunca vio esta premisa)
- **La misma decisión la traen dos gates posteriores distintos** → En el primero se ofrece reconfirmación; una vez decidida, ese resultado se registra; en el segundo, se consulta el registro actualizado
- **Un cambio con registro de ratificaciones no entra en verify** → El registro se mantiene intacto; no se elimina al cerrar el cambio
- **Fase intenta leer el registro de ratificaciones** → No tiene permisos; nunca debe consultarlo. Su único canal es el prompt de despacho

## Reglas de negocio

- El registro es por cambio (`sdd/{change-name}/ratified-decisions`), no global ni por fase
- La coincidencia es exacta en `record` (el token que identifica la decisión), nunca fuzzy
- Matching es `record` solamente; el anchor y summary solo intervienen en la presentación una vez que hay match
- Si `record` coincide pero el `summary` cambió, la comparación es trivial: son diferentes, se presenta en forma completa
- Answering "sí" a una reconfirmación corta registra la ratificación nuevamente (el contenido registrado no cambia a menos que el usuario lo ajuste)
- Answering "no" abre el item en su forma completa para ratificación nueva o rechazo
- `sdd-apply` NUNCA debe consultar este registro; su único canal para una resolución es el prompt de despacho del orquestador
- El registro no se poda, no se borra: viaja con el cambio hasta que cierra (y persiste en Engram como audit trail)

## Entidades y estados

- **Ledger de ratificaciones** — tabla en `sdd/{change-name}/ratified-decisions`: cuatro columnas (`record`, `ratified_summary`, `anchor`, `gate`), una fila por item ratificado. Estados: vacío → poblado → leído por gates posteriores
- **Item ratificado** — entrada de una sección de gate, después de la decisión. Estados: nuevo → registrado → re-emerge → reconfirmado o reabierto
- **Coincidencia** — prueba exacta de `record` entre un item nuevo y el registro. Estados: sin registro → coincidencia hallada → diferencia de contenido detectada

## Errores de cara al actor

- **Record duplicado en el registro** — Error de lógica del orquestador: solo un row por `record` por gate. Falla en auditoria
- **Usuario elige rechazar reconfirmación** → Item se abre en forma completa; sin falla
- **`sdd-apply` intenta leer el registro** → La guía explícita lo prohíbe; `sdd-apply` recibe resoluciones solo en su prompt

## Escenarios

### Scenario: La misma decisión sin cambio, reconfirmación corta

- **GIVEN** una decisión ratificada en el primer gate de ratificación de este cambio
- **WHEN** el mismo item resurge en un gate posterior idéntico en contenido
- **THEN** se presenta como un sí/no nombrando qué fue ratificado y dónde, sin reabrirlo en forma completa

### Scenario: El item ajustado en el gate anterior

- **GIVEN** un item que el usuario corrigió en lugar de confirmar tal cual se ofreció
- **WHEN** se registra al cerrar el gate
- **THEN** el texto ajustado es lo que se registra, no la versión original

### Scenario: Ratificación pendiente cuando resurge

- **GIVEN** una decisión ratificada en un gate cuya implementación no ha corrido aún
- **WHEN** un gate posterior abre en el mismo cambio
- **THEN** la ratificación ya está registrada y disponible para reconfirmación

### Scenario: La premisa cayó

- **GIVEN** una decisión ratificada "archivo X debe vivir en Y"
- **WHEN** un gate posterior presenta la misma decisión, pero el contenido ahora es "archivo X debe vivir en Z porque Y ya no existe"
- **THEN** el summary difiere, se presenta en forma completa
- **AND** si el usuario rechaza, reabre; si confirma, esa nueva versión se registra

### Scenario: Matching exacto en `record`, idéntico en summary

- **GIVEN** un item cuyo `record` y `summary` coinciden exactamente con lo registrado
- **WHEN** se presenta
- **THEN** se muestra como reconfirmación de una sola línea, no como item nuevo

### Scenario: Una decisión, dos gates posteriores

- **GIVEN** un item ratificado en el gate de Unresolved Decisions, y el mismo `record` resurge en el INTAKE GATE después
- **WHEN** el primero se presenta
- **THEN** se ofrece reconfirmación; se registra el resultado
- **WHEN** el segundo lo presenta
- **THEN** consulta el registro actualizado del primer gate

### Scenario: Record diferente, misma fuente

- **GIVEN** dos items que ambos anclan a la misma ubicación pero `record` tokens distintos
- **WHEN** ambos aparecen en el mismo gate
- **THEN** se presentan como dos items separados
- **AND** nada confunde uno con otro por la fuente común

### Scenario: Fase intenta consultar el registro

- **GIVEN** una fase que recibe `sdd/{change-name}/ratified-decisions` en su prompt de despacho
- **WHEN** intenta leerlo
- **THEN** debe fallar — tiene permisos solo sobre el contenido de su propio artefacto, nunca este registro
- **AND** la guía lo explícita en el dispatch CLAUDDE.md

### Scenario: El registro viaja con el cambio

- **GIVEN** un cambio que ha ratificado tres decisiones en su first gate
- **WHEN** el cambio entra en verify o se detiene por otra razón
- **THEN** el registro se mantiene intacto en Engram
- **AND** no se poda, no se borra

### Scenario: Reconfirmación rechazada

- **GIVEN** un item en reconfirmación corta que el usuario rechaza
- **WHEN** responde "no"
- **THEN** el item se abre en su forma completa, ofreciendo todas las acciones (confirmar, ajustar, rechazar)

## Referencias

- **Contrato compartido** → [`../../shared/references/gate-presentation.md`](../../shared/references/gate-presentation.md) — La forma corta de reconfirmación (una sola pregunta de sí/no, mostrando summary y anchor registrados)
- **Ledger persistencia** → [`../../payload/domains/development/CLAUDE.md`](../../payload/domains/development/CLAUDE.md) — El orquestador escribe el ledger al cerrar cada gate, Engram key `sdd/{change-name}/ratified-decisions`, cuatro campos (`record`, `ratified_summary`, `anchor`, `gate`)
- **Prohibición en apply** → Mismo archivo — `sdd-apply` MUST NOT read this key; its only channel is the dispatch prompt
- **Identificación de decisión** → [`../rule/scenario-verification-scope-token.md`](../rule/scenario-verification-scope-token.md) — El token `record:` es el identificador estable de una decisión dentro de una fase
