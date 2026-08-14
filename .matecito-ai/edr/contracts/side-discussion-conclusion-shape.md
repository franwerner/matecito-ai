# EDR — La conclusión de discusión lateral lleva cinco secciones fijas, cortadas para calzar con un ítem de buzón

- **Status:** Accepted
- **Date:** 2026-08-14

## Contexto
La conclusión es lo único que la sesión lateral devuelve. Cuando esa conclusión zanja algo arquitectónico, tiene que poder entrar al mecanismo de captura in-flow existente (propone → ratifica → materializa) sin que alguien tenga que reescribirla primero.

## Decisión
La conclusión lleva cinco secciones — `## Question` / `## Conclusion` / `## Why` / `## What it rests on` / `## Still open` — cortadas así a propósito: `## Conclusion` es directamente el `summary` de un ítem de buzón, `## Why` es su `rationale`, y `## What it rests on` es su `anchor`, sin transformación.

## Reglas verificables
- **[manual]** La conclusión lleva siempre las cinco secciones `## Question`, `## Conclusion`, `## Why`, `## What it rests on` y `## Still open`; una sección sin contenido se declara vacía (p. ej. `## Still open` → "None."), nunca se omite.
- **[manual]** Cuando una conclusión zanja una decisión arquitectónica, `## Conclusion` se usa directamente como `summary`, `## Why` como `rationale` y `## What it rests on` como `anchor` del ítem de buzón que la propone — sin reescritura.

## Alternativas consideradas
Cuerpo libre — descartado, mismo motivo que el traspaso: nada garantiza que la conclusión traiga lo que la captura in-flow necesita. Una forma calcada del template de EDR — descartada porque haría que la conclusión se vea como un registro ya ratificado, que es exactamente lo que no es (ver la decisión sobre que la conclusión no es un registro).

## Consecuencias
Una conclusión que zanja algo arquitectónico entra al camino de captura existente sin transformación. El costo es que la sesión lateral tiene que ceñirse a esta forma incluso cuando la conclusión no es una decisión arquitectónica — el corte igual sirve, sólo que `## What it rests on` puede no alimentar ningún buzón.

## Relacionados
- `relacionado-con` → [side-discussion-handoff-shape.md](side-discussion-handoff-shape.md) — la otra mitad del intercambio, con el mismo criterio de secciones fijas.
- `relacionado-con` → [side-discussion-conclusion-is-not-a-record.md](side-discussion-conclusion-is-not-a-record.md) — fija qué pasa con esta conclusión cuando zanja algo arquitectónico.
