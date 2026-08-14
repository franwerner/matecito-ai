# EDR — Tope de caracteres al resumen de un ítem de gate, aplicado al construirlo

- **Status:** Accepted
- **Date:** 2026-08-14

## Contexto
Un ítem de un mailbox de gate imprime una línea de resumen que el usuario lee antes de confirmar. Sin un tope, esa línea puede crecer sin control y volver el gate ilegible. El chequeo tiene que aplicarse en algún punto del ciclo de vida del retorno, y ese punto determina qué tan barato es corregir un exceso: mientras la fase todavía está viva y tiene el contexto cargado, o después, cuando ya terminó.

## Decisión
El tope se mide en caracteres del string —no en líneas de terminal, cuyo ancho ningún script conoce— y se declara por sección, con valor 250. Se aplica en el momento de construir el bloque de retorno, mientras la fase todavía está viva y corregir cuesta una edición y una re-corrida; no al validar un retorno ya armado, momento en el que la fase ya terminó y corregir exige re-despacharla entera.

## Reglas verificables
- **[auto]** Cada sección que declara ítems fija su propio tope de caracteres para el resumen, y el renderizador aborta sin emitir nada si algún resumen lo supera.
- **[manual]** El límite se mide en caracteres del string, nunca en líneas de terminal.

## Alternativas consideradas
Un tope en líneas de terminal — descartado porque el ancho de terminal es invisible para un script. Valores de 200 y 160, o un default único para toda sección — descartados frente a una medición sobre prosa real del repo (percentiles de longitud de líneas de índice existentes), que mostró que 250 recorta los casos más extremos sin afectar la prosa ordinaria. Aplicar el chequeo también (o solo) al validar el retorno ya construido — descartado porque en ese punto la fase ya terminó y corregir exige una re-corrida completa en vez de una edición inmediata.

## Consecuencias
Un exceso de longitud se detecta y se corrige en el momento más barato del ciclo, mientras la fase todavía tiene el contexto cargado — a costa de que el validador no vuelve a chequear el límite sobre un retorno ya emitido: un bloque escrito a mano en vez de renderizado no queda medido.

## Relacionados
- `relacionado-con` → [anchor-token-free-form.md](anchor-token-free-form.md) — mismo mecanismo de ítems de gate, mismo cambio que lo introdujo.
- `relacionado-con` → [per-field-description-cap.md](per-field-description-cap.md) — mismo criterio de tope por caracteres, aplicado a la descripción de un campo en vez de al resumen del ítem.
