# Dominio: `data`

Cómo se modelan, leen y escriben los datos persistentes: nivel de abstracción sobre la DB, patrón de acceso, convenciones de esquema.

## Criterio de pertenencia

Un concern nuevo va en `data` si trata sobre el modelo o el acceso a datos *en uso normal*. Si trata sobre datos a lo largo del tiempo (migraciones, backups, retención, borrado), va en `lifecycle`.

## Concerns en este dominio

| Concern | Prof. | Tipo | Qué decide |
|---|---|---|---|
| [data-access](data-access.md) | deep | decisión | Cómo el sistema lee y escribe datos persistentes: qué nivel de abstracción se usa sobre la DB, si hay patrón Repository, cómo se manejan las migraciones, y d... |
| [data-modeling](data-modeling.md) | light | decisión | Convenciones de bajo nivel que afectan esquema de DB, APIs, y código de dominio: tipo de IDs, borrado lógico vs físico, timestamps estándar, y si el modelo s... |
