# huma

- **Category:** Framework web
- **Version:** sin pinear
- **Status:** Accepted
- **Decided in phase:** contracts
- **Date:** 2026-07-23

## Por qué la elegimos

Framework Go que genera OpenAPI 3.1 y la validación de request/response a partir de structs, corriendo sobre `net/http` (stdlib, ServeMux). Da un contrato Go-first como fuente de verdad que la UI consume para generar sus tipos y schemas (vía Kubb), sin escribir el spec a mano.

## Alternativas descartadas

- **swaggo:** genera OpenAPI 2.0, versión más vieja del spec.
- **fuego:** descartado frente a huma.

## Notas

Usada en: contracts/api-contract (superficie HTTP-in del MCP → broker). Corre sobre `net/http` de la stdlib.
