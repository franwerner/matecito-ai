# EDR — El reporte entre-fases es el cuarto acto comunicativo, con su propia forma cerrada

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
Los actos comunicativos del hilo estaban declarados como un conjunto cerrado de tres: presentar una decisión, explicar un hallazgo, informar lo que se hizo. El reporte que el orquestador emite entre una fase y la siguiente no era ninguno de los tres, así que componía prosa libre — exactamente lo que la sección de formas cerradas existe para evitar. Peor: el párrafo de excepción de emisión, redactado en términos posicionales ("la prosa alrededor de ese material"), dejaba una puerta abierta a que ese reporte llevara preámbulo propio.

## Decisión
El reporte entre-fases se declara como el cuarto acto comunicativo, agregado dentro de `### Explaining — fixed forms, not free prose` — después del tercer acto y antes de la lista `Never` — sin reemplazar ni reformular a los tres existentes. Su forma es la del template fijo de ítem menos su slot de acciones: por ítem, el ancla más una sola línea de resumen; sin narrativa, preámbulo, recapitulación del contexto ni justificación previa al ítem; y sin ofrecer ninguna acción de ratificación, porque lo que se ratifica se ratifica en un gate — un momento distinto con su propia forma.

## Reglas verificables
- **[manual]** `payload/core/CLAUDE.md`, sección `### Explaining — fixed forms, not free prose`, enumera cuatro actos comunicativos; el reporte entre-fases es el cuarto y aparece después del tercer acto y antes de la lista `Never`.
- **[manual]** La forma del cuarto acto es, por ítem, el ancla más una línea de resumen, sin narrativa, preámbulo, recapitulación ni justificación previa, y declara explícitamente que no ofrece acciones de ratificación.
- **[manual]** Los tres actos ya declarados (presentar una decisión, explicar un hallazgo, informar lo que se hizo) quedan con su forma sin alterar.

## Alternativas consideradas
Una sección `###` propia para el reporte entre-fases — descartada: la sección abre con la oración "Three communicative acts cover nearly everything you say in the thread", un conjunto cerrado; un cuarto acto fuera de esa sección falsifica la oración que dice arreglar. Reordenar para que el acto nuevo encabece la lista, por ser el más frecuente — descartada: ningún requisito pide un reordenamiento y los tres actos existentes deben mantener su forma sin cambios.

## Consecuencias
El reporte entre-fases queda con una forma verificable en el mismo lugar donde ya se verifican los otros tres actos, en vez de quedar gobernado por partes en la regla del disparo de gates (que sólo dice qué aparece, nunca cómo se ve). El costo es que la sección de formas cerradas crece un párrafo más, pero deja de ser la única sección que no cubre el acto más frecuente del orquestador.

## Relacionados
- `relacionado-con` → [prose-register-single-home.md](prose-register-single-home.md) — mismo principio de un solo dueño para una forma de prosa, aplicado antes al registro de decisiones.
