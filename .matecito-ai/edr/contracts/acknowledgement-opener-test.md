# EDR — La prohibición de aperturas vacías pasa de lista de frases a test verificable

- **Status:** Accepted
- **Date:** 2026-08-12

## Contexto
"no 'sure / great question'" prohibía frases puntuales, pero una lista cerrada de frases se esquiva con una frase nueva que cumple la misma función: abrir sin aportar información.

## Decisión
La regla pasa de una lista a un test: una línea de apertura que no aporta información se elimina; una línea que dice qué lectura se tomó, qué archivo se tocó o qué elección se hizo es contenido, no una apertura, y se conserva.

## Reglas verificables
- **[manual]** Toda respuesta se revisa por si su primera línea aporta información (la lectura tomada, el archivo tocado, la elección hecha); si no aporta nada, se elimina antes de enviar.

## Alternativas consideradas
Ampliar la lista de frases prohibidas. Descartado: una lista cerrada nunca cubre la frase siguiente que cumple la misma función vacía; el test ataca la función, no la forma superficial.

## Consecuencias
Una frase nueva que abre sin aportar nada sigue prohibida aunque no esté en ninguna lista; una línea que sí aporta lectura o elección deja de leerse como "apertura sospechosa" y se conserva.

## Relacionados
- `relacionado-con` → [closing-line-only-for-a-pending-decision.md](closing-line-only-for-a-pending-decision.md) — la otra prohibición de la misma revisión que pasa de lista a test/excepción con límite.
