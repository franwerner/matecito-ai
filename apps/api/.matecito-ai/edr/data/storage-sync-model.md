# EDR — Modelo de almacenamiento y sincronización de records

- **Status:** Accepted
- **Date:** 2026-07-24

## Contexto

El broker persiste dos cosas de naturaleza distinta: el estado runtime del cockpit (eventos, changes, actividad de agentes) y una copia de los records de decisión y comportamiento (los EDR/spec del repo). Había que decidir quién es canónico, cómo se sincronizan las dos representaciones de los records, y dónde vive el versionado. El norte que condiciona todo es poder **deployar y compartir con compañeros** desde un servidor: features hechos, eventos e historia disponibles sin depender del entorno local de ninguna persona.

## Decisión

El contenido versionado de los records vive en la base de datos, **content-addressable con dedup por hash** (no se duplica el contenido repetido). Se descarta apoyarse en git como store de versiones: git vive local en cada máquina y un servidor compartido no tiene el checkout de nadie, así que usar git como store bloquearía el deploy/compartir. La base de datos es el artefacto **portable y deployable**, self-contained.

Hay tres tipos de dato, cada uno con un único dueño:

- **Records (los `.md` de decisión y comportamiento):** el **archivo (git) es canónico para editar**; la base de datos es **espejo más versiones**.
- **Eventos y estado runtime:** **viven solo en la base y la base es canónica** (nunca fueron archivos).
- **Fotos de código (los diffs del cockpit):** el mismo principio de no depender de git se extiende al código que un change toca — las vistas de diff de la UI **no** se derivan de git (ni commits, ni ramas, ni disciplina de commit exigida al flujo): el broker captura **fotos del contenido completo** de cada archivo en los bordes de cada batch de apply, en el mismo store content-addressable (dedup por hash), y los diffs se **derivan on demand** comparando fotos (o una foto contra el disco). Se descartó almacenar solo deltas/líneas cambiadas: un patch guardado se desalinea con el drift, no renderiza contexto y las cadenas de patches son frágiles; los estados completos hacen toda vista siempre computable. Las fotos son self-contained: sobreviven a rebase, ramas borradas y sirven deployado.

La sincronización tiene dos caminos:

- **Escrituras del AI → write-through por el daemon:** el daemon escribe el archivo **y** la base en la misma operación, sin ventana de desincronización.
- **Cambios out-of-band** (edición a mano, checkout o merge de git) → el **watch** re-espeja a la base; **si el archivo y la base divergen, gana el archivo**. La ventana asíncrona del watch se puede cerrar on-demand con un comando de sync/reindex (o un hook post-checkout de git).

El reparto lectura/escritura queda así: el **AI lee de los archivos** (el contenido actual, local; el versionado es transparente para su lectura), el **AI escribe por el daemon**, la **UI lee de la base** (queries, grafo, historia, versiones), y **deploy/compartir se sirve de la base**.

Todo lo que refleja contenido del repo se scopea por `(proyecto, rama)`; el **proyecto es el contenedor branch-independiente** (solo su registración y el watch son project-level). El contenido de versiones es compartido entre ramas por hash.

La **identidad del proyecto es el `project-id` de `<git-toplevel>/.matecito-ai/.id`**, un archivo commiteado que viaja con el repo: es machine-independent, sobrevive a mover o renombrar el checkout y funciona deployado o compartido. La **base keyea por `project-id`, no por el path local** —que sería inválido en un servidor que no tiene el checkout de nadie—. Esto refuerza que la base sea self-contained y deployable.

## Reglas verificables

- **[manual]** el contenido versionado de los records vive en la base, content-addressable con dedup por hash; no se deriva de git.
- **[manual]** para editar un record, el archivo (git) es canónico; si el archivo y la copia de la base divergen, gana el archivo (se re-espeja).
- **[manual]** las escrituras del AI pasan por el daemon (write-through: archivo más base en la misma operación); los cambios out-of-band se re-espejan por el watch, cerrables on-demand con un comando de sync.
- **[manual]** el AI lee el contenido de los archivos, no de la base; la UI y el deploy leen de la base.
- **[manual]** la base es self-contained y deployable: un servidor puede servir features, eventos e historia sin git ni los archivos locales.
- **[manual]** la identidad del proyecto es el `project-id` de `<git-toplevel>/.matecito-ai/.id` (commiteado, viaja con el repo); la base keyea por `project-id`, nunca por el path local.
- **[manual]** los eventos y el estado runtime viven solo en la base.
- **[manual]** las vistas de diff de código se derivan de las fotos del store (o de una foto contra el disco), nunca de git; los diffs no se almacenan.
- **[manual]** las fotos de código son contenido completo por archivo, en los bordes de batch, en el store content-addressable (dedup por hash).

## Alternativas consideradas

- **Git como store de versiones:** descartado; git es local a cada máquina y un servidor compartido no tiene el checkout de nadie, así que bloquearía el deploy/compartir.
- **Base canónica total, con los archivos como export puro:** descartado; rompe editar los records a mano, la operación git-native y la conciencia de rama.
- **Diffs de código derivados de git (merge-base/commits por batch):** descartado; exige disciplina de commits por batch, se rompe con rebase/ramas borradas y no funciona deployado sin checkout. Las fotos content-addressable dan lo mismo sin esas dependencias.
- **Almacenar solo deltas/líneas cambiadas de código:** descartado; frágil ante drift (los patches se desalinean), sin contexto para renderizar, y las vistas colapsadas dependen de cadenas de patches que no pueden fallar.

## Consecuencias

- La base es un artefacto deployable y compartible: un servidor sirve features, eventos e historia por sí solo.
- El record sigue siendo git-native localmente: se edita como archivo y se versiona con git.
- Trade-off: hay más que construir que apoyarse en git como store, a cambio de la portabilidad del artefacto.

## Relacionados

- `relacionado-con` → [data-access-entity-framework.md](data-access-entity-framework.md) — el borde de persistencia donde viven la escritura del espejo y de las versiones.
- `relacionado-con` → [data-modeling.md](data-modeling.md) — las convenciones de schema (IDs, borrado, idempotencia por hash, envelope) que aterrizan este modelo.
- `relacionado-con` → [../contracts/api-contract.md](../contracts/api-contract.md) — la superficie de lectura de la UI que consume la base (queries, grafo, historia, versiones).
- `relacionado-con` → [../contracts/mcp-server.md](../contracts/mcp-server.md) — la identidad por request y el scoping por proyecto/rama de los eventos.
