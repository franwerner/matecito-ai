# Go

- **Category:** Lenguaje
- **Version:** sin pinear
- **Status:** Accepted
- **Decided in phase:** structure
- **Date:** 2026-07-23

## Por qué la elegimos

Daemon local de larga vida que necesita concurrencia idiomática (goroutines/channels/context), un binario único autocontenido y cross-compilación trivial. Trae de base formateo y análisis estático canónicos.

## Alternativas descartadas

- Ninguna evaluada formalmente: elección de base para el broker.

## Notas

Usada en: structure/architecture-style, structure/folder-structure, structure/code-conventions y, de forma transversal, todos los EDRs del broker.
