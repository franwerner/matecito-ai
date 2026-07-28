# ent

- **Category:** ORM
- **Version:** v0.14.6
- **Status:** Accepted
- **Decided in phase:** data
- **Date:** 2026-07-28

## Por qué la elegimos

Permite definir el schema una sola vez en código y generar desde ahí el acceso a datos type-safe, en vez de mantener migraciones y queries como dos definiciones sincronizadas a mano. El modelo del broker es mayormente entidades con ciclo de vida (solo una de las once tablas es un log append-only) y un grafo de records que se recorre; la generación de código cubre esas travesías mejor que escribir una lectura por camino. Soporta claves primarias compuestas sin identificador sustituto y constraints declarativos de columna en la propia definición del schema.

## Alternativas descartadas

- **Queries type-safe generadas desde SQL escrito a mano (`sqlc`):** obliga a mantener el schema definido dos veces —en las migraciones y en las queries— y a sincronizarlo manualmente.
- **Mapeo objeto-relacional basado en reflection:** pierde la verificación en tiempo de compilación que da la generación de código.

## Notas

Usada en: data/data-access-entity-framework. Se apoya en la librería estándar de acceso a base, así que no cambia el driver de SQLite en Go puro. Versión pineada: el proyecto todavía no alcanzó su 1.0 y su superficie puede cambiar entre versiones menores.
