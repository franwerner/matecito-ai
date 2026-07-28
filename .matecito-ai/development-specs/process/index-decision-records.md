# Capability — Indexar y versionar los EDR/spec records

- **Status:** Accepted
- **Date:** 2026-07-27

## Propósito

Mantener, por proyecto, un índice consultable del estado actual de los records de decisión y comportamiento (los `.md` de `edr/` y `development-specs/`) junto con su **contenido versionado guardado en la base**, para poder navegar el estado vigente y mostrar la versión exacta que un evento aplicó. La base es self-contained: guarda tanto el índice como el contenido de cada versión, de modo que la UI puede consultarlo y el broker se puede deployar/compartir **sin apoyarse en git como store de versiones**. El `.md` es la fuente canónica para editar (vive en git); el sistema versiona su contenido en la base para sostener el pin y la traza.

## Actores

- **La escritura del AI vía daemon** (write-through) — cuando el AI crea o edita un record, el daemon escribe el `.md` **y** actualiza la base en la misma operación, sin ventana de desincronización.
- **El file-watcher del daemon** — observa el árbol de records del proyecto y dispara el (re)indexado ante cada cambio out-of-band (edición a mano, `git checkout`/`merge`).
- **El comando de sync/reindex** — cierra on-demand la ventana asíncrona del watch (o lo dispara un hook post-checkout de git), re-espejando el estado a la base.
- **El proceso de indexado** — lee el `.md`, extrae su proyección y contenido, aplica el versionado y actualiza el índice.

## Precondiciones

- El proyecto está registrado (el índice y las versiones se scopean por proyecto).
- Existen los directorios de records observados del proyecto.

## Flujo principal (del proceso)

1. El (re)indexado se dispara por tres caminos, todos convergen en el mismo flujo:
   - **Build inicial:** al arrancar, el daemon recorre los `.md` existentes bajo los directorios de records del proyecto y los indexa.
   - **Write-through:** cuando el AI escribe un record por el daemon, la escritura del `.md` y la actualización de la base ocurren en la misma operación (archivo y base consistentes al instante, sin pasar por la ventana del watch).
   - **Watch (out-of-band) + sync:** ante un cambio hecho por fuera del daemon (edición a mano, `git checkout`/`merge`), el file-watch re-espeja el `.md` a la base; su ventana asíncrona se puede cerrar on-demand con el comando de sync/reindex.
2. El sistema lee el `.md` disparado.
3. Extrae la **proyección** (id/slug, tipo, status, refs/relaciones a otros records) y el **contenido**.
4. Aplica el versionado **copy-on-write lazy**, según el estado de la versión actual (gobernado por el ciclo de vida de la versión de record, que este proceso aplica pero no define): update-in-place solo si esa versión no está congelada **y** ninguna otra rama la apunta; en cualquier otro caso —congelada por un evento, o apuntada por otra rama— fork a una nueva versión actual para la rama que cambió.
5. Actualiza el índice (upsert) apuntando a la versión actual del record, con su proyección al día.

## Casos borde

- **Un `.md` indexado se borra** → el record pasa a estado **`deleted`** (soft-delete) **en esa `(proyecto, rama)`**: su entrada **permanece** en el índice (marcada `deleted`), no se elimina; las versiones **congeladas** (referenciadas por eventos) sobreviven y los eventos que las apuntan las siguen resolviendo. El `deleted` es por `(proyecto, rama)`: borrarlo en una rama **no afecta** su estado en otra. Si el `.md` reaparece, el record vuelve a **`active`**.
- **El AI escribe un record por el daemon** (write-through) → el `.md` y la base quedan consistentes al instante, en la misma operación; no se pasa por la ventana asíncrona del watch.
- **Cambio de rama con `git checkout`** → el watch (o el comando de sync on-demand) re-espeja a la base el estado de records de la nueva rama; el índice/estado por `(proyecto, rama)` refleja la rama activa sin arrastrar el de la anterior.
- **Ref colgada** (una relación entre records apunta a un slug inexistente) → se marca la relación como colgada (flag) y el indexado del record sigue; no rompe el índice.
- **`.md` malformado** (no se puede parsear su status o sus refs) → ese record se saltea y se flaggea; el resto del árbol se indexa igual.
- **`.md` nuevo aparece** → se indexa como su primera versión (que pasa a ser la versión actual).

## Reglas de negocio

- El `.md` es la fuente de verdad para editar. Si el índice o el contenido versionado divergen del `.md`, gana el `.md`: se re-indexa (con posible nueva versión según la regla de bump).
- **El AI lee de los archivos** (el contenido actual, local; el versionado es transparente para su lectura) y escribe por el daemon; **la UI lee de la base** (índice, versiones, historia). El contenido versionado se guarda en la base, no se deriva de git.
- El versionado de las versiones —cuándo se congela una versión, update-in-place vs fork, la inmutabilidad y supervivencia de las congeladas— lo gobierna el **ciclo de vida de la versión de record** ([`../lifecycle/record-version.md`](../lifecycle/record-version.md)); este proceso lo aplica, no lo define.
- El **índice/estado** de un record (que existe, su status, `active`/`deleted` y su versión vigente) se scopea por **`(proyecto, rama)`** — cada rama guarda su estado en la base, no se deriva de git, así el `active`/`deleted` **no flip-flopea** al cambiar de rama. El **contenido de versiones** es content-addressable y se **comparte entre ramas por hash** (no se scopea por rama). El proyecto es el contenedor branch-independiente: solo su registración y el watch son project-level.
- **Identidad del record en monorepo (owning-root):** un monorepo puede contener varios `.matecito-ai/` (uno por app). El **ID de un record** es su slug + su **`owning-root`**, donde `owning-root` = el nombre de la **carpeta padre directa** que contiene ese `.matecito-ai/`. El daemon **deriva el owning-root del path** del archivo al indexar. Ejemplos: `.matecito-ai/` en la raíz del repo → owning-root = la raíz del proyecto; `apps/api/.matecito-ai/` → owning-root = `api`; `apps/ui/.matecito-ai/` → owning-root = `ui`. Así dos records con el mismo slug bajo `.matecito-ai/` distintos **no colisionan**.

## Entidades y estados

- **Record (EDR/spec)** — el record identificado por su slug + su `owning-root`. Estados: **`active`** (el `.md` está presente en la estructura de carpetas) → **`deleted`** (el `.md` se borró de la estructura), y de vuelta **`deleted` → `active`** si el `.md` reaparece. Es un **soft-delete**: la entrada del record no se elimina, se marca `deleted`. Las versiones **congeladas** del record sobreviven en cualquier estado.
- **Versión de un record** — definida en su lifecycle propio: [`../lifecycle/record-version.md`](../lifecycle/record-version.md) (`actual` → `congelada`).
- **Entrada de índice** — el puntero del proyecto a la versión actual de cada record; se hace upsert al indexar y **permanece** cuando el `.md` se borra, marcada `deleted` (sin afectar versiones congeladas).

## Escenarios

### Scenario: build inicial al arrancar

- **GIVEN** un proyecto con `.md` de records ya presentes en su árbol
- **WHEN** arranca el daemon y corre el build inicial
- **THEN** cada `.md` existente queda indexado con su proyección y su primera versión como versión actual

### Scenario: el índice apunta a la versión que resulta del re-indexado

- **GIVEN** un record cuyo `.md` cambia
- **WHEN** se re-indexa
- **THEN** el índice apunta a la versión actual resultante — la misma actualizada en el lugar, o la nueva si hubo fork (según el ciclo de vida de la versión de record)

### Scenario: borrado con versión congelada superviviente (soft-delete)

- **GIVEN** un record indexado con una versión congelada que un evento referencia
- **WHEN** su `.md` se borra
- **THEN** el record pasa a estado `deleted` (su entrada permanece, marcada `deleted`), la versión congelada sobrevive y el evento que la referenciaba la sigue resolviendo

### Scenario: reaparición de un record borrado

- **GIVEN** un record en estado `deleted` (su `.md` se había borrado)
- **WHEN** su `.md` reaparece en la estructura y se re-indexa
- **THEN** el record vuelve a estado `active` sobre la misma entrada, y sus versiones congeladas siguen intactas

### Scenario: write-through al escribir el AI

- **GIVEN** el AI crea o edita un record por el daemon
- **WHEN** el daemon procesa la escritura
- **THEN** escribe el `.md` y actualiza la base en la misma operación, dejando archivo y base consistentes al instante, sin pasar por la ventana asíncrona del watch

### Scenario: sync on-demand tras cambiar de rama

- **GIVEN** un proyecto cuyos records difieren entre dos ramas y se hace `git checkout` a la otra rama
- **WHEN** el watch se dispara o se corre el comando de sync/reindex
- **THEN** la base re-espeja el estado de records de la rama activa, y el índice/estado por `(proyecto, rama)` refleja esa rama sin arrastrar el `active`/`deleted` de la anterior

### Scenario: mismo slug bajo `.matecito-ai/` distintos no colisiona (owning-root)

- **GIVEN** dos records con el mismo slug, uno bajo `apps/api/.matecito-ai/` y otro bajo `apps/ui/.matecito-ai/`
- **WHEN** el sistema los indexa
- **THEN** se identifican distinto por su owning-root (`api` vs `ui`), sin colisión

### Scenario: contenido compartido por hash, estado por rama

- **GIVEN** el mismo contenido de un record presente en dos ramas
- **WHEN** el sistema lo versiona
- **THEN** se guarda **una sola versión** (deduplicada por hash), compartida por ambas ramas y por cualquier evento que la referencie; y borrar el `.md` en una rama **no elimina ese contenido** (solo marca el record `deleted` en esa `(proyecto, rama)`; la versión congelada y el contenido siguen resolviéndose desde otra rama/evento)

### Scenario: ref colgada marcada sin romper el índice

- **GIVEN** un `.md` cuya relación apunta a un slug de record inexistente
- **WHEN** se indexa
- **THEN** la relación se marca como colgada (flag) y el record se indexa igual; el índice sigue consultable

### Scenario: `.md` malformado se saltea y flaggea

- **GIVEN** un árbol de records donde un `.md` no se puede parsear (status o refs ilegibles)
- **WHEN** corre el indexado
- **THEN** ese `.md` se saltea y se flaggea, y el resto de los records se indexan normalmente

## Referencias

- **EDR** → [`../../../apps/api/.matecito-ai/edr/data/storage-sync-model.md`](../../../apps/api/.matecito-ai/edr/data/storage-sync-model.md) — el modelo de almacenamiento y sincronización: base self-contained con el contenido versionado, archivo canónico para editar, write-through más watch, y el scoping por `(proyecto, rama)`.
- **EDR** → [`../../../apps/api/.matecito-ai/edr/data/data-access-entity-framework.md`](../../../apps/api/.matecito-ai/edr/data/data-access-entity-framework.md) — el borde de persistencia (Repository) donde viven la escritura del índice y de las versiones.
- **EDR** → [`../../../apps/api/.matecito-ai/edr/data/data-modeling.md`](../../../apps/api/.matecito-ai/edr/data/data-modeling.md) — el store de versiones content-addressable y el pin de la versión referenciada por un evento.
- **EDR** → [`../../../apps/api/.matecito-ai/edr/contracts/api-contract.md`](../../../apps/api/.matecito-ai/edr/contracts/api-contract.md) — la superficie de lectura de la UI que consume el índice y la versión exacta que aplicó un evento.
- **Lifecycle** → [`../lifecycle/record-version.md`](../lifecycle/record-version.md) — el ciclo de vida de la versión de record (`actual` → `congelada`) que este proceso aplica al versionar.
- **Rule** → [`../rule/event-scoping.md`](../rule/event-scoping.md) — cómo se scopea el índice y las versiones por proyecto.
- **Flow** → [`../flow/submit-phase-artifact.md`](../flow/submit-phase-artifact.md) — el flujo que pinea la versión vigente de un EDR/spec al persistir el evento.
