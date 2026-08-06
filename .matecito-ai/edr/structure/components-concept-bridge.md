# EDR — Forma del puente al mudar el concepto `components`

- **Status:** Accepted
- **Date:** 2026-08-05

## Contexto
Al mudar la mitad repo-level del concepto a su hogar propio, el lugar de origen —la documentación del concepto de capability-spec— necesitaba una forma de transición: qué queda ahí, y cómo se re-apuntan las citas existentes que hoy nombran el ancla original.

## Decisión
El lugar de origen conserva únicamente su propia proyección del concepto y renombra su encabezado, dejando el ancla original exclusiva del hogar nuevo. Las citas existentes se re-apuntan por ruta y subsección precisa, nunca por el ancla genérica que dejó de existir en el lugar de origen.

## Reglas verificables
- **[manual]** Ninguna cita del concepto `components` en el payload usa el ancla genérica del lugar de origen; cada una apunta a una ruta y una subsección precisas.

## Alternativas consideradas
Conservar el encabezado original en el lugar de origen, sin renombrarlo, como forma de puente que no requeriría tocar las citas que lo nombran, se descartó: dejar ese rótulo en el lugar de origen mantiene ahí el mismo encuadre que el cambio existe para cerrar, y vuelve inverificable por búsqueda mecánica que el ancla quedó exclusiva del hogar nuevo.

## Consecuencias
El criterio de cierre queda mecánicamente verificable: una búsqueda del ancla original en el payload devuelve un único resultado, el título del hogar nuevo. El costo es re-apuntar cada cita existente a su sección precisa en vez de dejarlas apuntando al ancla genérica.
