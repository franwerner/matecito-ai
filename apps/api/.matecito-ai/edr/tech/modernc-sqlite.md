# modernc.org/sqlite

- **Category:** DB
- **Version:** sin pinear
- **Status:** Accepted
- **Decided in phase:** data
- **Date:** 2026-07-23

## Por qué la elegimos

Driver de SQLite en Go puro (sin cgo): mantiene trivial la cross-compilación y el binario autocontenido, sin toolchain nativo.

## Alternativas descartadas

- **Driver SQLite basado en cgo:** complica la cross-compilación por requerir un toolchain de C.

## Notas

Usada en: data/data-access-entity-framework.
