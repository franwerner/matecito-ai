# EDR — Los catorce capability-specs retirados se borran, no se marcan Deprecated

- **Status:** Accepted
- **Date:** 2026-09-08

## Contexto
El cambio `sdd/drop-apps-subtree` retira catorce capability-specs cockpit-only (`flow`×7, `lifecycle`×2, `process`×3, `rule`×2) junto con el código que documentan. `.matecito-ai/development-specs/INDEX.md:43` dice, sin condiciones: "Retirar una capacidad: marcá el spec `Deprecated` con link a su reemplazo; no borres el archivo" — esa redacción queda tal cual está, ratificada, así que este cambio deja catorce ausencias que contradicen una convención que sigue escrita en la página. El precedente para cerrar ese hueco es `structure/retired-vocabulary-in-record-stores.md`, que existe justamente para que un barrido futuro no vuelva a marcar como defecto una excepción deliberada.

## Decisión
Los catorce specs se borran como archivo, no se marcan `Deprecated`, porque ninguno tiene un reemplazo al que enlazar — el brief confirmado y ratificado en el gate de propose manda el borrado para este cambio puntual. La convención general de `INDEX.md:43` no se toca ni se reescribe: sigue rigiendo para cualquier retiro futuro de una sola capability con reemplazo. Este EDR es el registro de que la tanda de catorce es una excepción decidida, no una ruptura accidental de esa convención.

## Reglas verificables
- **[manual]** Ninguno de los catorce specs retirados por `sdd/drop-apps-subtree` existe como archivo tras el cambio, y ninguno quedó marcado `Deprecated`.
- **[manual]** `.matecito-ai/development-specs/INDEX.md:43` conserva la redacción "marcá el spec Deprecated... no borres el archivo" sin cambios — la convención general no se toca por esta excepción.
- **[manual]** Un futuro barrido o auditoría sobre `INDEX.md:43` encuentra esta excepción documentada acá, en vez de leer las catorce ausencias como una ruptura accidental de la convención.

## Alternativas consideradas
Confiar sólo en la línea de mantenimiento que agrega `sdd-archive` al cerrar el cambio. Descartada: esa línea registra qué pasó, no que fue una excepción deliberada, y vive en una lista cronológica que nadie relee al chequear una convención vigente.

## Consecuencias
La convención de `INDEX.md:43` ("no borres, marcá Deprecated") sigue intacta para cualquier retiro futuro de una sola capability con reemplazo. Este EDR es el punto que explica por qué esta tanda de catorce fue la excepción, sin que la convención general tenga que cambiar de redacción para dar cuenta de ella.
