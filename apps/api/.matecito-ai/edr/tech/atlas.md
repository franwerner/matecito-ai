# atlas

- **Category:** Migrations
- **Version:** v1.2.3
- **Status:** Accepted
- **Decided in phase:** data
- **Date:** 2026-07-28

## Por qué la elegimos

Deriva las migraciones de la definición del schema en código, así que no hay un juego de migraciones escrito a mano que sostener en paralelo. Las migraciones versionadas son un flujo de primera clase y producen archivos SQL, que se siguen embebiendo en el binario y aplicando al arranque en modo librería: el binario queda autocontenido, igual que antes.

## Alternativas descartadas

- **Migraciones versionadas escritas a mano (`goose`):** deja el schema definido dos veces y obliga a sincronizarlo manualmente con la definición en código.
- **Migración automática por diferencia contra la base al arrancar:** no es versionada ni auditable — no deja artefacto que revisar en el cambio ni al que volver.
- **Migraciones aplicadas por una CLI externa:** rompen el objetivo de binario autocontenido que se migra solo al arrancar.

## Notas

Usada en: data/data-access-entity-framework. Las migraciones generadas se embeben en el binario y se aplican al arranque.
