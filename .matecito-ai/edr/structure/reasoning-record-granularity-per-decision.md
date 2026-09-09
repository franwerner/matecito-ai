# EDR — Los comentarios de razonamiento reubicados se graban un record por decisión, no uno por comentario

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
Cerca de treinta comentarios `<!-- matecito-ai: -->` de clase 1 (narración de un defecto) y clase 2 (alternativa rechazada) salen de `payload/core/CLAUDE.md` y `payload/domains/development/CLAUDE.md` en este cambio. Sin una regla de granularidad, cada comentario removido podía terminar como su propio EDR, incluso cuando varios comentarios en sitios distintos documentaban exactamente la misma decisión — el mismo patrón de duplicación que ya forzó a centralizar otras reglas de este ecosistema.

## Decisión
Los comentarios se graban **un record por decisión**, no uno por comentario: varios comentarios que declaran la misma decisión en sitios distintos del payload colapsan en un único record, cuyo `## Reglas verificables` nombra cada archivo y sección que gobierna, verbatim — nunca en `## Alcance`, porque esa sección del template del EDR es para globs a nivel de convención, y un título de sección no es un glob. La forma naive que cada nota rechazaba se conserva en `## Alternativas consideradas`, porque eso es lo que evita que el defecto vuelva. El conteo final de records queda deliberadamente sin fijar de antemano.

## Reglas verificables
- **[manual]** Cada record que reubica el razonamiento de uno o más comentarios removidos nombra, en su propio `## Reglas verificables`, el archivo y la sección exactos que gobierna — nunca en `## Alcance`.
- **[manual]** Varios comentarios que declaran la misma decisión en sitios distintos del payload colapsan en un único record, con todos los sitios gobernados nombrados en ese mismo record.
- **[manual]** El `## Alternativas consideradas` de cada record conserva la forma naive que el comentario original rechazaba.

## Alternativas consideradas
Un record por comentario (1:1) — descartada: produce records duplicados que divergen en la próxima edición, el mismo defecto que la centralización de otras reglas de prosa de este ecosistema ya existe para evitar. Un único record omnibus para todo el sweep — descartada: una búsqueda literal por el título de una sección específica termina en un record cuya Decisión es sobre otra cosa.

## Consecuencias
El número de records que deja este sweep no está fijado por adelantado — lo fija el contenido: cuántas decisiones distintas hay detrás de los ~30 comentarios, no cuántos comentarios hay.

## Relacionados
- `relacionado-con` → [mirror-sites-join-the-prune-sweep.md](mirror-sites-join-the-prune-sweep.md) — precedente ya ratificado de la misma forma un-record-varios-sitios.
