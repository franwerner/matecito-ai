# EDR — El sweep de vocabulario de gating incluye el 24° archivo fuera del Scope original

- **Status:** Accepted
- **Date:** 2026-08-31

## Contexto
El cambio narrow-gating-triggers ratificó un Scope de 23 archivos que llevaban el vocabulario Tier 1/Tier 2 retirado. `payload/domains/development/skills/gentle-ai/sdd-apply/SKILL.md:493` decía "el tier de estos buzones y quién los consume están fijados en Sección D.3" y quedaba fuera de ese Scope — invisible al grep literal que el spec usa como aceptación, y, una vez completado el cambio, apuntando a una Sección D.3 que ya no tiene columna `Tier`.

## Decisión
El 24° archivo SÍ se barre: se reescribe el comentario de la línea 493 para que nombre `gates:` en vez del sustantivo de tier retirado — un solo comentario reescrito, sin cambio de comportamiento y sin abrir una clase de archivo nueva. Cae del allow-list del sweep final (tarea 3.4).

## Reglas verificables
- **[manual]** `payload/domains/development/skills/gentle-ai/sdd-apply/SKILL.md` no contiene el sustantivo `tier` en ninguna de sus líneas.
- **[manual]** La línea que antes decía "el tier de estos buzones" nombra `gates:` en su lugar, sin otro cambio en esa línea ni en las que la rodean.

## Alternativas consideradas
Dejar la línea como estaba, aceptando que un skill leído en runtime quede apuntando a un concepto retirado (Sección D.3 sin columna Tier). Descartada: un ejecutor que la lea después del cambio recibe una instrucción que ya no resuelve a nada verificable.

## Consecuencias
El sweep final (tarea 3.4) ya no necesita listar este archivo en su allow-list de residuos conocidos — el archivo queda limpio junto con los otros 23.
