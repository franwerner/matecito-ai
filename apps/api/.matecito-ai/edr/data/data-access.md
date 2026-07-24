# EDR — Acceso a datos

- **Status:** Accepted
- **Type:** decision
- **Date:** 2026-07-23
- **Applied pattern:** Repository — aislar el borde de persistencia detrás de una interfaz angosta para testear con SQLite real sin acoplar la lógica a los detalles del store.

## Contexto

El corazón del broker es una SQLite que funciona como event-log: se le hacen appends de eventos y se deriva estado a partir de queries. No es un modelo de entidades CRUD, así que un ORM pelearía contra el diseño. Además se quiere mantener trivial la cross-compilación (binario autocontenido, sin dependencias nativas) y poder testear el acceso a datos contra una base real, no mockeada.

## Decisión

Usamos un driver de SQLite en Go puro (sin cgo) para que la cross-compilación siga siendo trivial. La capa de queries se construye sobre la librería estándar de acceso a base más un generador de queries type-safe a partir de SQL escrito a mano; sin ORM. El borde de persistencia se expone como una interfaz de Repository angosta (un lado de escritura y uno de lectura), NO como repositorios por entidad, porque es un event-log y no entidades. Las migraciones se embeben en el binario, están versionadas y se aplican al arranque con una herramienta de migración en modo librería, para mantener el binario autocontenido. La frontera transaccional vive en el escritor del store: cada comando de escritura atómico se envuelve en una transacción dentro del único goroutine escritor.

## Alcance

- `apps/api/internal/store/**` — el borde de persistencia (interfaz Repository, queries, migraciones embebidas).

## Reglas verificables

- **[manual]** el driver de SQLite es Go puro (sin cgo).
- **[manual]** las queries se definen en SQL y se generan type-safe; no se usa ORM.
- **[manual]** el store expone una interfaz angosta de escritura/lectura, no repositorios por entidad.
- **[manual]** las migraciones están embebidas en el binario y se aplican al arranque.
- **[manual]** las escrituras atómicas se envuelven en una transacción dentro del único escritor.

## Alternativas consideradas

- **ORM:** descartado; pelea contra el modelo event-log (appends + derivación por query), no contra entidades.
- **Driver SQLite basado en cgo:** descartado por complicar la cross-compilación; el driver Go puro la mantiene trivial.
- **Repositorios por entidad:** descartados; no hay entidades CRUD, hay un log de eventos con una interfaz angosta de escritura/lectura.

## Consecuencias

- Binario autocontenido y cross-compilable sin toolchain nativo.
- Queries type-safe con SQL explícito, alineado con el event-log.
- El único escritor concentra la frontera transaccional, lo que simplifica la atomicidad.
- Trade-off: escribir SQL a mano y regenerar el código type-safe es más trabajo por query que un ORM que lo autogenera.

## Relacionados

- `depende-de` → [../runtime/concurrency-async.md](../runtime/concurrency-async.md) — la frontera transaccional se apoya en el único goroutine escritor.
- `relacionado-con` → [data-modeling.md](data-modeling.md) — el schema concreto de tablas se define ahí.
