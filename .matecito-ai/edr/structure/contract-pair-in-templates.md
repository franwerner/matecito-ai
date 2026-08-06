# EDR — El contrato máquina vive junto al template legible

- **Status:** Accepted
- **Date:** 2026-08-06

## Contexto

Cada artefacto durable tiene dos caras que describen la misma estructura: una legible, que una persona abre para entender qué va en cada sección, y una declarativa, que las herramientas leen para renderizar y validar. Las dos dicen lo mismo en idiomas distintos, y sólo sirven si están sincronizadas.

Cuando dos archivos deben moverse juntos, la distancia entre ellos es lo que decide si el drift se nota o no.

## Decisión

El contrato máquina vive **en el mismo directorio que el template legible**, como par adyacente, y no en un directorio de esquemas aparte.

## Alcance

- `payload/domains/development/references/*/templates/*.yaml` — el lado declarativo del par.
- `payload/domains/development/references/*/templates/*.md` — el lado legible, que vive junto al anterior.

## Reglas verificables

- **[manual]** Todo artefacto durable con contrato máquina tiene su par declarativo en el mismo directorio que su template legible, no en un árbol de esquemas separado.
- **[auto]** El auto-chequeo de cada herramienta compara el contrato del repo contra el desplegado y falla si divergieron.
- **[manual]** El template legible ilustra la forma; el contrato declarativo la impone. Ninguno de los dos redefine lo que el otro ya fija.

## Alternativas consideradas

- **Un directorio de esquemas aparte.** Descartado: separa dos archivos que sólo tienen sentido sincronizados, y hace que editar uno sin el otro sea invisible hasta que algo falla en runtime.
- **Un solo archivo que sea a la vez contrato y documentación.** Descartado: el template legible tiene valor propio como material que una persona lee de corrido, y comprimirlo dentro de un formato declarativo lo vuelve ilegible sin ganar nada.

## Consecuencias

- El drift entre las dos caras se nota al abrir el directorio: están una al lado de la otra.
- Un cambio de estructura obliga a tocar dos archivos, siempre. Es fricción deliberada: es lo que hace que la sincronía sea un acto consciente.
- El directorio de templates deja de ser sólo material de lectura y pasa a contener el contrato ejecutable, lo que exige que quien edite ahí sepa cuál de los dos gobierna.

## Relacionados

- `relacionado-con` → [two-scripts-render-and-validate.md](two-scripts-render-and-validate.md) — los dos ejecutables que leen este par.
- `relacionado-con` → [../contracts/data-contract-derived-and-producer-neutral.md](../contracts/data-contract-derived-and-producer-neutral.md) — qué expresa el lado declarativo del par.
