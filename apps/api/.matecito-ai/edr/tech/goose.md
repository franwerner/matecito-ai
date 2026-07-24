# goose

- **Category:** Migrations
- **Version:** sin pinear
- **Status:** Accepted
- **Decided in phase:** data
- **Date:** 2026-07-23

## Por qué la elegimos

Aplica migraciones versionadas en modo librería al arranque del broker; con las migraciones embebidas en el binario, el arranque mantiene el schema al día sin herramientas externas y el binario queda autocontenido.

## Alternativas descartadas

- **Migraciones aplicadas por una CLI externa:** rompen el objetivo de binario autocontenido que se migra solo al arrancar.

## Notas

Usada en: data/data-access. Las migraciones se embeben en el binario.
