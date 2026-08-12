# EDR — Emitir en inglés y presentar traducido no alcanza a los avisos internos

- **Status:** Accepted
- **Date:** 2026-08-12

## Contexto
La nueva regla de idioma (payload y salida de fase en inglés, presentado en el idioma de la conversación) podía leerse como que reemplaza también la regla existente de avisos internos, que ya fija una línea corta en inglés.

## Decisión
La regla de emitir-en-inglés-presentar-traducido no alcanza a los avisos internos; conservan su propia regla de una única línea corta en inglés, sin traducir.

## Reglas verificables
- **[manual]** El bloque de idioma nombra explícitamente la regla de avisos internos como fuera de su alcance.
- **[manual]** La regla de avisos internos sigue fijando el aviso como una única línea corta en inglés.

## Alternativas consideradas
Dejar que la regla de idioma se aplique también a los avisos internos. Descartado: un aviso interno ("Loaded intake.", "Saved to Engram.") es una nota mecánica de una línea, no material gate-facing que un lector sin el vocabulario del ecosistema necesite entender traducido; traducirlo agrega una transformación a algo que ya es mínimo.

## Consecuencias
Los avisos internos no cambian de comportamiento con este cambio; sólo el material gate-facing (summary/rationale, resumen ejecutivo, el bloque de tono) queda gobernado por la regla de idioma.
