# Capability — Capturar fotos de código

- **Status:** Accepted
- **Date:** 2026-07-27

## Propósito

Mantener, por change, la cadena de fotos del contenido de cada archivo que el apply toca —el estado antes y después de cada batch— para que las vistas de diff se deriven comparando fotos, sin depender de git ni de disciplina de commits.

## Actores

- **Fuente de captura pre-edición** — el mecanismo del runtime del coder que intercepta **sincrónicamente** cada edición de archivo antes de que ocurra y reporta el path a la ingesta del broker; la edición no procede hasta que la captura termina. Hoy: los hooks pre-tool de Claude Code, vía el mismo cliente mínimo del canal mecánico.
- **El submit del batch de apply** — al persistirse el artefacto de un batch, el broker captura el estado posterior de los archivos tocados.
- **El broker** — lee los archivos, decide si fotografiar y sostiene la cadena de fotos.

## Precondiciones

- La fuente de captura está configurada en el runtime del coder (la gestiona el installer de matecito).
- Para que se capture: el proyecto está registrado y hay un change `active` para el contexto — mismo gate que la ingesta mecánica; sin change activo, no se fotografía nada.

## Flujo principal

1. El agente va a editar un archivo; la fuente de captura intercepta la edición **antes** de que ocurra y reporta el path.
2. Si es el **primer toque de ese archivo desde el último batch** (o desde el inicio del change), el broker lee el archivo del disco —todavía intacto— y guarda la foto de su contenido: la primera de la cadena es la **base** del archivo en este change. Toques subsiguientes dentro del mismo batch no fotografían.
3. La edición procede.
4. Al submit del batch de apply, el broker lee el estado final de los archivos tocados y guarda la foto del **después** del batch.
5. La cadena por archivo queda: base → después del batch 1 → (antes del batch 2) → después del batch 2 → … Cada evento de apply referencia, por archivo, su foto antes → después.

## Casos borde

- **Broker caído al momento de una foto pre-edición** → aplica el spool de ingesta con contenido: el cliente lee el archivo él mismo y guarda la foto completa en el spool local; se reconcilia al volver el contacto. Ninguna base se pierde.
- **El apply crea un archivo nuevo** → su base es la **ausencia** ("no existía"): las vistas lo derivan como archivo nuevo entero.
- **El apply borra un archivo** → la foto del después registra la **ausencia**: las vistas lo derivan como borrado.
- **Edición que no pasa por las tools de edición del runtime** (un comando arbitrario, una edición a mano) → no dispara captura y no genera foto propia. Queda evidenciada igual: si un batch posterior toca el archivo, su foto de antes ya la incluye (el hueco entre fotos la delata como cambio fuera del flujo); si nada la vuelve a tocar, la delata la comparación contra el archivo real.
- **Pre-edición cuyo contenido es idéntico a la última foto de la cadena** → no se guarda nada nuevo (dedup por hash).

## Reglas de negocio

- **Las fotos son estados completos, no cambios**: cada foto es el contenido íntegro del archivo en ese instante, con todo lo acumulado. Los diffs son siempre derivados comparando fotos (o una foto contra el disco); nunca se almacena un diff.
- **Se fotografía solo en los bordes de batch**: la primera edición que toca el archivo en un batch (el antes) y el submit del batch (el después). Las ediciones intermedias dentro del batch no generan fotos.
- La captura pre-edición es **sincrónica**: la edición no procede hasta que la foto está tomada (o spooleada); no hay carrera entre fotografiar y editar.
- El contenido de las fotos es **content-addressable con dedup por hash** — la misma maquinaria de versionado que los records; archivos idénticos no se re-almacenan.
- La cadena y su base se scopean **por change**, no por sesión: un change que abarca varias sesiones continúa su cadena; la base es la primera foto del archivo dentro del change.

## Escenarios

### Scenario: base capturada en el primer toque

- **GIVEN** un change activo y un archivo aún no tocado por este change
- **WHEN** el agente va a editarlo y la fuente intercepta la edición
- **THEN** el broker fotografía el contenido intacto ANTES de que la edición ocurra y esa foto queda como base del archivo en la cadena del change

### Scenario: el después se captura al submit

- **GIVEN** un batch de apply que editó archivos
- **WHEN** se persiste el artefacto del batch
- **THEN** el broker fotografía el estado final de cada archivo tocado como el después de ese batch, y el evento de apply referencia por archivo su antes → después

### Scenario: ediciones intermedias no fotografían

- **GIVEN** un archivo ya fotografiado como antes en el batch en curso
- **WHEN** el mismo batch lo vuelve a editar
- **THEN** no se toma ninguna foto adicional hasta el submit del batch

### Scenario: edición manual encerrada entre fotos

- **GIVEN** un archivo con foto del después del batch N, editado a mano fuera del runtime
- **WHEN** el batch N+1 lo vuelve a tocar y se captura su antes
- **THEN** esa foto ya incluye la edición manual, y el hueco entre "después de N" y "antes de N+1" evidencia el cambio fuera del flujo

### Scenario: broker caído no pierde la base

- **GIVEN** el broker caído en el momento de una captura pre-edición
- **WHEN** la fuente intercepta la edición
- **THEN** el cliente lee el archivo, guarda la foto completa en el spool local, la edición procede, y al volver el contacto la foto se reconcilia con su timestamp original

### Scenario: archivo creado

- **GIVEN** un batch que crea un archivo que no existía
- **WHEN** se captura la cadena de ese archivo
- **THEN** su base es la ausencia y las vistas lo derivan como archivo nuevo entero

### Scenario: archivo borrado

- **GIVEN** un batch que borra un archivo fotografiado
- **WHEN** se captura el después del batch
- **THEN** la foto registra la ausencia y las vistas lo derivan como borrado

### Scenario: sin change activo no se fotografía

- **GIVEN** un proyecto registrado sin change `active` para el contexto (o con el change `closed`)
- **WHEN** el agente edita archivos
- **THEN** no se captura ninguna foto; el gate es el mismo que el de la ingesta mecánica

### Scenario: contenido idéntico no se re-almacena

- **GIVEN** una captura cuyo contenido tiene el mismo hash que la última foto de la cadena
- **WHEN** el broker la procesa
- **THEN** no almacena contenido nuevo (dedup por hash)

## Referencias

- **Rule** → [`../rule/ingestion-spool.md`](../rule/ingestion-spool.md) — el spool con contenido que garantiza que ninguna foto se pierda.
- **Flow** → [`../flow/serve-code-diff.md`](../flow/serve-code-diff.md) — las vistas de diff que se derivan de estas fotos.
- **Process** → [`ingest-mechanical-events.md`](ingest-mechanical-events.md) — el canal hermano de ingesta mecánica, con el mismo gate de change activo y el mismo cliente.
- **EDR** → [`../../../apps/api/.matecito-ai/edr/data/storage-sync-model.md`](../../../apps/api/.matecito-ai/edr/data/storage-sync-model.md) — el store content-addressable con dedup por hash que estas fotos reutilizan.
- **EDR** → [`../../../apps/api/.matecito-ai/edr/data/data-modeling.md`](../../../apps/api/.matecito-ai/edr/data/data-modeling.md) — (Pending) el modelado del store de fotos y sus referencias desde los eventos de apply.
