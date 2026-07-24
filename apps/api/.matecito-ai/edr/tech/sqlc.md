# sqlc

- **Category:** ORM
- **Version:** sin pinear
- **Status:** Accepted
- **Decided in phase:** data
- **Date:** 2026-07-23

## Por qué la elegimos

Genera código de acceso a datos type-safe a partir de SQL escrito a mano: da queries tipadas sin un ORM, lo que encaja con el modelo event-log (appends + derivación por query) en vez de pelear contra él.

## Alternativas descartadas

- **ORM:** pelea contra el event-log, que no es un modelo de entidades CRUD.
- **SQL a mano con escaneo manual de filas:** más verboso y sin la verificación de tipos en compilación que da la generación.

## Notas

Usada en: data/data-access. Se apoya en la librería estándar de acceso a base para ejecutar las queries generadas.
