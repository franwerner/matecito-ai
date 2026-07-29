# EDR — Modelado de datos

- **Status:** Accepted
- **Date:** 2026-07-28

## Contexto

El store del broker es una única base global por máquina (SQLite embebida en el daemon), self-contained para poder deployarse o compartirse, con la UI leyendo de la base y el AI de los archivos. El modelo debe cubrir las entidades que los capability-specs ya fijaron: proyectos, changes, el event-log por change, las versiones content-addressable (de records y de fotos de código), los pins de eventos a versiones, las notas de iteración y el índice de decisiones. Las convenciones de bajo nivel —identidad, borrado, tiempo, idempotencia— atraviesan todas las tablas y son costosas de migrar después; se fijan acá, antes de construir el store.

## Decisión

- **IDs: UUID v7** como clave primaria de toda entidad referenciable desde afuera. Ordenable por tiempo (índices baratos, natural para logs), no expone secuencias, y no colisiona entre instancias si la base se comparte o deploya. El `project-id` que porta el archivo de identidad commiteado usa el mismo formato. Esta decisión fija el **formato de los IDs** de las entidades que existan (proyectos, changes, eventos o notas son ejemplos ilustrativos, no una enumeración normativa): **no** determina qué entidades deben existir ni obliga a promover a tabla propia un atributo que hoy es solo una etiqueta — eso lo deciden el modelo concreto y los capability-specs.
- **Posición del log, además del ID:** cada event-log de change lleva una secuencia monótona propia — es el cursor de resume de la lectura en vivo, no una identidad; conviven ambas.
- **Contenido versionado: keyeado por hash** (content-addressable, dedup) — ya fijado por el modelo de storage; las versiones de records y las fotos de código comparten ese store de contenido.
- **Borrado: sin DELETE físico en todo el store.** La desaparición se modela como estado: un change se cierra (nunca se borra), un record queda `deleted` por rama, un proyecto se desregistra marcándose inactivo con su historia intacta. Única excepción: una nota de iteración `pendiente` puede borrarla su autor desde la UI antes de que el AI la lea; entregada o resuelta ya es historia y queda.
- **Timestamps:** toda entidad lleva sus timestamps de creación y actualización. Los eventos llevan además su **timestamp de ocurrencia** como campo propio: el orden semántico es cuándo ocurrió, no cuándo llegó (la reconciliación del spool entrega tarde).
- **Idempotencia: por hash del contenido canónico.** La identidad de un envío es su contenido: para un submit, el hash del artefacto canónico más su anclaje (change, fase), excluyendo campos que estampa el servidor; para un ítem de ingesta, ídem incluyendo su timestamp original. Un reintento o re-flush produce el mismo hash → mismo evento, sin duplicar y sin que ningún productor genere ni arrastre claves.
- **Envelope de eventos:** todo evento se persiste y se sirve como `{type, payload}` tipado por su clase (submit de fase, mecánico, nota) más su metadata (change, posición, timestamp de ocurrencia, agente); el payload es la forma canónica de su clase.
- **Sin multitenancy:** instancia única global single-user; el aislamiento es el scoping por proyecto — toda tabla de contenido se scopea por el ID del proyecto (y lo que corresponde, además por rama, según el modelo de storage).

## Alcance

- `apps/api/internal/**` — el schema del store, sus migraciones y el código que lo consume.

## Reglas verificables

- **[manual]** toda tabla de entidad referenciable usa UUID v7 como clave primaria; no se mezclan tipos de ID sin justificación en este EDR.
- **[manual]** el event-log de cada change mantiene una secuencia monótona propia (la posición), independiente del ID del evento.
- **[manual]** toda tabla lleva `created_at` y `updated_at` NOT NULL; la tabla de eventos lleva además `occurred_at`.
- **[manual]** no existe DELETE físico sobre proyectos, changes, eventos, versiones ni records; la única operación de borrado real permitida es sobre notas de iteración en estado `pendiente`, por su autor.
- **[manual]** toda tabla de contenido incluye la referencia al proyecto; ninguna query cruza proyectos sin filtrarla.
- **[manual]** la deduplicación de escrituras usa el hash del contenido canónico como clave de idempotencia; ningún productor genera idempotency-keys propias.
- **[manual]** los eventos se persisten como envelope `{type, payload}` con su clase y metadata; la UI consume ese shape.

## Alternativas consideradas

- **Autoincrement/rowid como PK:** lo más simple en SQLite, pero acopla la identidad a la instancia local — rompe el norte de base self-contained compartible/deployable.
- **ULID:** equivalente funcional a UUID v7; se descartó por no sumar sobre v7 en este stack.
- **Idempotency-key generada por el productor:** estándar HTTP, pero obliga a cada productor (agente, hook, spool) a generar y persistir claves entre retries — infraestructura extra para lo que el hash del contenido ya garantiza.
- **DELETE físico de proyectos/changes:** libera espacio pero rompe la trazabilidad y los pins de versiones que los eventos referencian; descartado.

## Consecuencias

- La identidad es portable: mover la base o compartirla no re-keyea nada; el `project-id` viaja con el repo y matchea el formato de la base.
- El costo del "sin borrado físico" es crecimiento monótono del store; aceptable para el volumen esperado (texto y metadata) y mitigable a futuro con compactación de contenido no referenciado — que sería una decisión nueva, no un DELETE ad-hoc.
- La idempotencia por hash hace la deduplicación transparente para todos los productores, pero exige definir con precisión la forma canónica de cada payload (qué campos entran al hash); esa definición vive con el contrato de cada clase de evento en el código.
- Trade-off: dos identificadores por evento (ID + posición) — mínima redundancia a cambio de separar identidad estable de cursor de lectura.

## Relacionados

- `depende-de` → [storage-sync-model.md](storage-sync-model.md) — el modelo de storage (base self-contained, content-addressable, scoping por proyecto/rama) que este EDR aterriza en convenciones de schema.
- `relacionado-con` → [data-access-entity-framework.md](data-access-entity-framework.md) — las migraciones y el borde de persistencia materializan estas convenciones.
- `relacionado-con` → [../contracts/mcp-server.md](../contracts/mcp-server.md) — la idempotencia de las tools de escritura se apoya en la clave por hash.
- `relacionado-con` → [../contracts/api-contract.md](../contracts/api-contract.md) — el envelope `{type, payload}` es el shape que consume la superficie de lectura de la UI.
