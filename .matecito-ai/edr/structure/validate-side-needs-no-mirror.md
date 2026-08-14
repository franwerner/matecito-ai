# EDR — El lado de validación no necesita un helper espejo del lado de renderizado

- **Status:** Accepted
- **Date:** 2026-08-14

## Contexto
Cuando el lado de renderizado gana un helper nuevo para producir el adorno de un ítem en más de una forma de sección (tabla y lista), surge la pregunta de si el lado de validación necesita un cambio equivalente para poder chequear esa misma forma.

## Decisión
El lado de validación no gana un helper espejo. Ya chequea la declaración de ítems por su presencia, sin importar qué forma de renderizado la acompañe; su única edición es ampliar el filtro que decide qué secciones entran al autochequeo de contrato, para que también cubra la nueva forma de sección.

## Reglas verificables
- **[auto]** El validador de retornos chequea la declaración de ítems por su presencia, no por la forma de renderizado de la sección que la contiene.
- **[manual]** Ningún helper de adorno se agrega al lado de validación; su único cambio es el filtro de qué secciones entran al autochequeo.

## Alternativas consideradas
Agregar un helper espejo en el lado de validación — descartada porque el chequeo ya es agnóstico a la forma de renderizado; agregar un espejo sería una segunda copia de una lógica que el validador no necesita para su propio trabajo.

## Consecuencias
El validador se mantiene desacoplado de cómo se renderiza una sección, cambiando sólo lo estrictamente necesario — a costa de que el criterio de qué entra al autochequeo ahora cubre una forma de sección más, y hay que recordar seguir ampliándolo si aparece una tercera.

## Relacionados
- `relacionado-con` → [item-shaping-helper-seam.md](item-shaping-helper-seam.md) — el helper del que este EDR explica por qué no necesita espejo.
- `relacionado-con` → [two-scripts-render-and-validate.md](two-scripts-render-and-validate.md) — la separación general entre renderizar y validar de la que este EDR es un caso.
