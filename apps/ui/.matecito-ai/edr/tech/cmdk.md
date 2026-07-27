# cmdk

- **Category:** Other
- **Version:** sin pinear
- **Status:** Accepted
- **Decided in phase:** frontend
- **Date:** 2026-07-27

## Por qué la elegimos

Primitive headless de comando/búsqueda filtrable (accesible, teclado-first) sobre el que se construye la paleta de comandos (⌘K) del cockpit. Aprobado en el gate de `ui-base-components` junto con `lucide-react`.

## Alternativas descartadas

- Ninguna evaluada formalmente.

## Notas

En este batch (`ui-base-components`) queda **instalado como dependencia únicamente**: el primitive Command que lo envuelve se difiere a un batch posterior. Se coordina con `radix-ui.md` (Dialog como shell del ⌘K) cuando se construya.
