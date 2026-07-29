# EDR — Acceso a datos con entity framework

- **Status:** Accepted
- **Date:** 2026-07-28
- **Applied pattern:** Repository — aislar el borde de persistencia detrás de una interfaz angosta, para que el código generado por el framework quede como detalle del store y los tests corran contra una base real sin acoplar la lógica.

## Contexto

El borde de persistencia se había resuelto asumiendo que el store era esencialmente un log de eventos: appends y estado derivado por query, sin entidades CRUD, por lo que un ORM "pelearía contra el diseño". Ese supuesto quedó desactualizado al cerrarse el modelo concreto: de las once tablas del store, **una sola** es un log append-only; las otras diez son entidades con ciclo de vida propio —estados, borrado lógico, versionado copy-on-write—. Además, el subconjunto de records forma un grafo que se recorre, no un conjunto de lecturas planas.

Sostener ese modelo con SQL de migraciones escrito a mano *más* un conjunto de queries escrito a mano obliga a mantener dos definiciones del mismo schema sincronizadas manualmente, y el costo crece con cada relación del grafo. Lo que hace falta es una única definición del schema en código, de la que las migraciones se **deriven**.

Los condicionantes que no cambian: el binario debe seguir siendo autocontenido (migraciones embebidas, versionadas y aplicadas al arranque en modo librería), la cross-compilación debe seguir siendo trivial (driver de base embebida en Go puro, sin cgo), el borde de persistencia debe seguir siendo angosto, la frontera transaccional vive en el único goroutine escritor, y los tests corren contra base real embebida, sin mocks.

## Decisión

Definimos el schema una sola vez, en código, con el entity framework **Ent**, y **derivamos** las migraciones de esa definición con **Atlas**. Las migraciones son un artefacto generado, versionado y revisado; no se mantiene un juego de migraciones escrito a mano en paralelo a la definición del schema.

Las migraciones generadas se siguen embebiendo en el binario y aplicándose al arranque en modo librería: el binario sigue siendo autocontenido. El driver de la base embebida sigue siendo Go puro, sin cgo — el framework se apoya en la librería estándar de acceso a base, así que la elección de driver no cambia.

El borde de persistencia sigue expuesto como una interfaz angosta (un lado de escritura y uno de lectura), **no** como repositorios por entidad: el código generado es un detalle interno del store, no la superficie que consume el resto del sistema. La frontera transaccional sigue viviendo en el único goroutine escritor: cada comando de escritura atómico se envuelve en una transacción ahí.

Los constraints de integridad —unicidad de las claves naturales compuestas y validaciones declarativas de columna— se declaran en la definición del schema, no editando a mano las migraciones generadas, y deben materializarse como constraints reales en la base: la integridad no puede depender de que toda escritura pase por este binario.

## Alcance

- `apps/api/internal/store/**` — el borde de persistencia: definición del schema, código generado, migraciones derivadas y embebidas.

## Reglas verificables

- **[manual]** el schema se define en código y las migraciones se derivan de esa definición; no se mantienen migraciones escritas a mano en paralelo.
- **[manual]** las migraciones están embebidas en el binario, versionadas, y se aplican al arranque.
- **[manual]** el driver de la base embebida es Go puro (sin cgo).
- **[manual]** el store expone una interfaz angosta de escritura/lectura, no repositorios por entidad.
- **[manual]** las escrituras atómicas se envuelven en una transacción dentro del único escritor.
- **[manual]** los constraints de integridad (unicidad de claves naturales compuestas y validaciones declarativas de columna) se declaran en la definición del schema, no editando las migraciones generadas.
- **[manual]** esos constraints existen en la base, no solo en el código generado: una escritura cruda que los viole es rechazada por el motor.

## Alternativas consideradas

- **Queries type-safe generadas desde SQL escrito a mano, con migraciones versionadas mantenidas aparte (la decisión anterior):** descartada por obligar a mantener dos definiciones del mismo schema sincronizadas a mano; el costo crece con cada relación del grafo de records.
- **Migración automática por diferencia contra la base al arrancar:** descartada por no ser versionada ni auditable — no deja artefacto que revisar en el cambio ni al que volver.
- **Mapeo objeto-relacional basado en reflection:** descartado por perder la verificación en tiempo de compilación que sí da la generación de código.

## Consecuencias

- Una sola definición del schema: las migraciones dejan de ser una segunda fuente de verdad que hay que mantener en sync a mano.
- El grafo de records se recorre con travesías generadas en vez de una lectura escrita a mano por cada camino.
- La unicidad de las claves naturales y las validaciones declarativas quedan expresadas donde se lee el modelo, en vez de sobrevivir solo dentro de las migraciones.
- Trade-off: el framework reserva la clave primaria para su propio identificador, así que las entidades cuya identidad natural es compuesta llevan un identificador sustituto y expresan esa identidad como unicidad. La garantía es equivalente; el modelo pierde la propiedad de que la clave natural *sea* la clave primaria.
- El binario sigue autocontenido y cross-compilable sin toolchain nativo; el borde angosto y la frontera transaccional en el único escritor no cambian.
- Trade-off: mayor compromiso con el framework — define el schema, genera el código de acceso y produce las migraciones; salir de él implicaría rehacer las tres cosas, no solo una capa.
- Trade-off: el framework todavía no alcanzó su versión 1.0, así que su superficie puede cambiar entre versiones menores; por eso la versión queda pineada en el catálogo de tecnologías.

## Relacionados

- `depende-de` → [../runtime/concurrency-async.md](../runtime/concurrency-async.md) — la frontera transaccional se apoya en el único goroutine escritor.
- `relacionado-con` → [data-modeling.md](data-modeling.md) — el schema concreto se define ahí.
