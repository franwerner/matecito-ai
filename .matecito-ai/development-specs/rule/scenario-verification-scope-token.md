# Capability — Token `verification:` por escenario

- **Status:** Accepted
- **Date:** 2026-08-06
- **Components:** cli

## Propósito

Un cambio sólo promete lo que puede cumplir. Un escenario cuya verificación depende de trabajo excluido del alcance, o de comportamiento preexistente que el delta no altera, no es cobertura pendiente de ese cambio: es cobertura **de otro dueño**. Esta capability fija cómo se declara ese dueño —un token por escenario—, qué hace `sdd-verify` ante cada valor, y cómo el token sobrevive al archivado.

## Actores

- `sdd-verify` — fase de verificación que lee cada escenario y clasifica la cobertura del cambio.
- `sdd-archive` — fase de archivado que materializa los escenarios en el capability-spec durable.
- Lector del spec durable — desarrollador que consulta qué comportamiento es verificable por quién.

## Reglas de negocio

- El token es **por escenario**, exactamente una línea, primer bullet del bloque, **antes del `GIVEN`**, con la forma `- **verification:** ` seguida del valor.
- Los valores legales son exactamente tres y el conjunto es cerrado: `in-scope` (por ausencia), `deferred → <cambio>`, `standing → <dueño>`.
- La ausencia de la línea es `in-scope` por default deliberado: los escenarios ya verificados no llevan marca.
- Dos escenarios del mismo requisito pueden llevar valores distintos; el alcance del requisito **no** se deriva de ningún token.
- Un token **nunca rebaja el contrato**: lo diferido es quién verifica y cuándo, no si la regla vale.
- `sdd-verify` respeta el token: un escenario `deferred` o `standing` no se clasifica `UNTESTED`, no emite un CRITICAL por falta de cobertura, y no empuja el veredicto a FAIL por ese escenario.
- `sdd-archive` conserva el token **verbatim** al plegar el delta al capability-spec durable: mismo texto, mismo lugar, sin agregar/reescribir/normalizar/quitar.
- El enforcement es **documental** (prosa en skills y references): sin check mecánico nuevo en `checks.yaml` ni `validate-artifact.js`.
- Que el cambio nombrado en `deferred → <cambio>` cubra efectivamente esos escenarios queda **fuera de este alcance** (diferido).
- Por conservarse el token, un capability-spec durable pasa a llevar **además del contrato ratificado, el estado de verificación de cada escenario**. Es intencional.

## Entidades y estados

- **Escenario** — descripción Given/When/Then de un comportamiento observable. Estados: `in-scope` (verificable por este cambio) → `deferred` (verificable por otro cambio) → `standing` (ya verificado por otro flujo). Las transiciones no son temporales: un escenario lleva su token en el spec durable y lo conserva.

## Escenarios

### Scenario: la definición se encuentra completa sin abrir ninguna skill

- **GIVEN** un lector que necesita saber qué es el token, qué valores admite y qué significa su ausencia
- **WHEN** busca la definición
- **THEN** la encuentra completa en la doc del concepto de capability-spec, sin abrir la skill de ninguna fase

### Scenario: no aparece superficie de referencia nueva

- **GIVEN** el árbol de references del dominio antes y después del cambio
- **WHEN** se comparan
- **THEN** no hay archivo ni directorio de reference nuevo: la definición vive en el archivo que `sdd-verify` y `sdd-archive` ya citaban

### Scenario: las skills citan, no redefinen

- **GIVEN** la prosa que este cambio agrega a `sdd-verify` y a `sdd-archive`
- **WHEN** se busca en ella la forma del token o la semántica de sus valores
- **THEN** cada skill remite al hogar canónico por su ruta desplegada (nunca una ruta `payload/…`) y enuncia sólo su propia reacción a cada valor

### Scenario: un `deferred` en un spec `Accepted` se lee como intención

- **GIVEN** un lector que se topa con un escenario `deferred` dentro de un capability-spec durable `Accepted`
- **WHEN** consulta la definición para interpretarlo
- **THEN** la definición declara que el spec durable lleva también el estado de verificación de cada escenario, y que es intencional

### Scenario: la ausencia de la línea es `in-scope`

- **GIVEN** un escenario cuyo bloque arranca directamente en `GIVEN`
- **WHEN** se clasifica su alcance de verificación
- **THEN** cuenta como `in-scope`, sin marca y sin ambigüedad

### Scenario: dos escenarios del mismo requisito con alcances distintos

- **GIVEN** un requisito con dos escenarios, uno con `deferred → X` y otro sin línea
- **WHEN** se leen
- **THEN** cada uno lleva su propio alcance y el del requisito no se deriva de ninguno de los dos

### Scenario: el conjunto de valores es cerrado

- **GIVEN** la definición publicada
- **WHEN** se busca qué valores son legales
- **THEN** enumera exactamente tres y los declara cerrados: ninguna otra cadena es un token

### Scenario: el token no rebaja el contrato

- **GIVEN** un escenario con `deferred → X`
- **WHEN** se pregunta si la regla que enuncia sigue siendo obligatoria
- **THEN** sigue siéndolo: lo diferido es quién la verifica y cuándo, no si vale

### Scenario: un `deferred` sin cobertura no es cobertura faltante

- **GIVEN** un escenario con `deferred → X` y sin test que lo cubra
- **WHEN** corre la verificación del cambio
- **THEN** se reporta con su token y no como `UNTESTED`
- AND no aporta ningún CRITICAL ni empuja el veredicto a FAIL

### Scenario: un `standing` sin cobertura deja visible a su dueño

- **GIVEN** un escenario con `standing → <dueño>` y sin test que lo cubra
- **WHEN** corre la verificación del cambio
- **THEN** se reporta con su token, con el dueño legible en el reporte, y no como `UNTESTED`
- AND no aporta ningún CRITICAL ni empuja el veredicto a FAIL

### Scenario: los escenarios fuera de alcance no inflan ni diluyen el resumen

- **GIVEN** un cambio con escenarios `in-scope` y escenarios con token
- **WHEN** se arma el resumen de cumplimiento
- **THEN** el denominador cuenta sólo los `in-scope`, y los demás siguen visibles en la matriz con su token

### Scenario: un escenario `in-scope` sin cobertura sigue siendo CRITICAL

- **verification:** `standing → sdd-verify`
- **GIVEN** un escenario sin línea de token y sin test que pase cubriéndolo
- **WHEN** corre la verificación
- **THEN** se clasifica `UNTESTED` como CRITICAL, igual que antes de este cambio

### Scenario: un test que corrió y falló no queda exento

- **GIVEN** un escenario con `deferred → X` cuyo test sí corrió y falló
- **WHEN** se clasifica
- **THEN** el fallo se reporta como fallo: el token exime de la falta de cobertura, nunca de un resultado negativo

### Scenario: el token viaja al store durable

- **GIVEN** un escenario con `deferred → X` en el delta del cambio
- **WHEN** el cambio se archiva
- **THEN** el capability-spec durable lleva ese escenario con la misma línea, con el mismo texto y en el mismo lugar del bloque

### Scenario: un escenario sin token no gana uno al archivar

- **GIVEN** un escenario del delta sin línea de token
- **WHEN** el cambio se archiva
- **THEN** el escenario durable tampoco la lleva

### Scenario: en un MODIFIED manda el token del delta

- **GIVEN** un escenario MODIFIED cuya copia durable llevaba `standing → Y` y cuya copia en el delta lleva `deferred → X`
- **WHEN** el cambio se archiva
- **THEN** queda el del delta, sin fusión, sin normalización y sin conservar el anterior

### Scenario: el merge sigue siendo no destructivo

- **verification:** `standing → sdd-archive`
- **GIVEN** un capability-spec con escenarios que el delta no menciona
- **WHEN** el cambio se archiva
- **THEN** esos escenarios se preservan tal cual, con sus tokens si los tenían

### Scenario: la validación del store no cambia de resultado

- **GIVEN** un capability-spec cuyo escenario lleva la línea del token antes de su `GIVEN`
- **WHEN** se valida el store
- **THEN** los hallazgos son los mismos que sin la línea — en particular, el escenario sigue teniendo `GIVEN`, `WHEN` y `THEN` completos y no se agrega ningún check nuevo

### Scenario: nada obliga a la herencia

- **GIVEN** un escenario con `deferred → X`
- **WHEN** se planifica o se cierra el cambio X
- **THEN** nada de este cambio verifica ni exige que X cubra ese escenario: la herencia queda para un cambio posterior

## Referencias

- **Referencia de concepto** → [`~/.claude/references/spec/README.md`](../../../payload/domains/development/references/spec/README.md) — el hogar canónico donde se define el token, su forma, sus valores, su uso per-escenario y su consecuencia aceptada sobre el spec durable.
- **Skill: sdd-verify** → [`payload/domains/development/skills/gentle-ai/sdd-verify/SKILL.md`](../../../payload/domains/development/skills/gentle-ai/sdd-verify/SKILL.md) — cómo lee y respeta el token la fase de verificación.
- **Skill: sdd-archive** → [`payload/domains/development/skills/gentle-ai/sdd-archive/SKILL.md`](../../../payload/domains/development/skills/gentle-ai/sdd-archive/SKILL.md) — cómo conserva el token verbatim al archivar.
