# EDR — El Discovery Gate vive en `sdd-explore`, no en `sdd-intake`, porque necesita el código ya leído

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
Con las dos lanes fijas, `sdd-intake` pasó a ser un passthrough sin rol de discovery. El ciclo de dos
pasadas del Discovery Gate — leer el código afectado y sólo entonces formular el formulario de preguntas
— sólo tiene sentido una vez que el cambio ya fue leído: una pregunta anclada en lo que `sdd-explore`
encontró, no una inventada antes de abrir un archivo.

## Decisión
El Discovery Gate corre en `sdd-explore`, después de que la fase lee el código afectado, no en
`sdd-intake`. `sdd-explore` retorna `status: needs-input` cuando tiene preguntas reales, ancladas a lo
que ya leyó.

## Reglas verificables
- **[manual]** `payload/domains/development/CLAUDE.md`, sección `### Discovery Gate (MANDATORY)`,
  declara que `sdd-explore` lee el código afectado PRIMERO y sólo entonces formula el formulario de
  discovery.

## Alternativas consideradas
Dejar el discovery en `sdd-intake`, antes de leer código — descartada: una pregunta formulada antes de
abrir un archivo no está anclada a nada concreto del cambio, y es exactamente el tipo de pregunta
inventada que este mecanismo existe para evitar.

## Consecuencias
El formulario de discovery sólo existe después de que `sdd-explore` corrió su primera pasada de lectura;
`sdd-intake` no tiene preguntas propias que hacer.
