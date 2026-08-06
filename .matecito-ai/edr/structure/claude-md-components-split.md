# EDR — Partición del restatement de `components` en la guía del dominio

- **Status:** Accepted
- **Date:** 2026-08-05

## Contexto
La guía del dominio resumía el concepto repo-level de `components` y su proyección por-capability en un único bloque, anidado bajo la sección del capability-spec. Ese encuadre invertía el orden real: la declaración repo-level no depende de ninguna proyección, pero vivía subordinada a una de ellas.

## Decisión
El restatement se divide en dos secciones, en el orden que espeja la relación real: una sección propia para el concepto repo-level, ubicada antes de la sección del capability-spec, y la proyección por-capability se queda dentro de esa sección, citando el concepto repo-level en vez de repetirlo.

## Reglas verificables
- **[manual]** En la guía del dominio, la sección del concepto repo-level de `components` precede a la sección de capability-specs.

## Alternativas consideradas
Mantener un solo bloque con una cita a cada documento se descartó: un bloque único no puede expresar que la declaración precede y no depende de la proyección — solo separarlo en dos secciones, en ese orden, deja esa relación clara para quien lee la guía.

## Consecuencias
La guía del dominio queda ordenada de declaración a proyección, consistente con el resto del cambio. Cada mitad cita el documento que efectivamente la contiene, en vez de una cita compuesta apuntando a documentos distintos desde un solo párrafo.
