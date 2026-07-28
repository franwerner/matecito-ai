# github.com/google/uuid

- **Category:** Otros
- **Version:** v1.6.0
- **Status:** Accepted
- **Decided in phase:** data
- **Date:** 2026-07-28

## Por qué la elegimos

Genera los identificadores UUID v7 que el modelado de datos fija como clave primaria de toda entidad del store. La librería estándar no ofrece UUID, y la versión 7 exige un layout concreto (timestamp de milisegundos en los bits altos más relleno aleatorio) que hay que implementar exacto para que el orden por tiempo se sostenga.

## Alternativas descartadas

- **Implementarlo a mano:** son pocas líneas, pero un error en el layout de bits rompe el ordenamiento temporal de forma silenciosa y recién se nota con volumen; no vale el riesgo para evitar una dependencia chica y estable.
- **Otras librerías de UUID del ecosistema:** equivalentes en funcionalidad; se eligió la más difundida y con menos superficie.

## Notas

Usada en: data/data-modeling (la convención de identificadores) y data/data-access-entity-framework (el store que los genera).
