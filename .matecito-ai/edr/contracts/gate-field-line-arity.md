# EDR — La aridad de la línea de campo de un item compuesto la decide el tipo de item, no una forma única

- **Status:** Accepted
- **Date:** 2026-09-01

## Contexto
Un item presentado por el walkthrough compartido de gates puede llevar contenido compuesto por varios campos en vez de una sola línea. Hasta ahora el único caso conocido era la propuesta de contrato, con su línea de continuación de tres partes (nombre — tipo — descripción). El brief de intake se modela también como item compuesto en "La cuarta forma" de `gate-presentation.md`, imprimiendo cada flag decidido como su propia línea de campo — pero un flag de decisión no tiene dimensión de tipo: solo nombre y valor decidido.

## Decisión
La aridad de la línea de campo la decide el tipo de item, no una forma única para todo item compuesto: tres partes (nombre — tipo — descripción) para la propuesta de contrato, dos (nombre — valor) para el brief de intake. Ambas aridades quedan como líneas de continuación del mismo item, dentro de la misma familia compartida ya fijada por `contracts/nested-field-continuation-line.md`, nunca como viñetas anidadas ni como una tabla dentro del item.

## Reglas verificables
- **[manual]** El brief de intake imprime cada flag decidido como una línea de continuación de dos partes (`· field: {name} — {value}`), nunca de tres.
- **[manual]** La propuesta de contrato sigue imprimiendo su línea de continuación de tres partes (`· field: {name} — {type} — {description}`), sin cambios por esta decisión.
- **[manual]** "La cuarta forma" de `gate-presentation.md` nombra los dos casos — contrato y brief de intake — sin enumerar los flags concretos de ningún dominio.

## Alternativas consideradas
Reusar la línea de tres partes del contrato para el brief, con el slot de tipo llevando un valor dummy o quedando vacío. Descartada porque un flag de decisión tiene nombre y valor decidido, y ninguna dimensión de tipo; cualquiera de las dos formas de llenar el slot medio (dummy o vacío) haría que el formato afirmara tres partes significativas donde solo hay dos.

## Consecuencias
El costo de sostener dos aridades en vez de una se paga explícitamente, fijado en `gate-presentation.md` y registrado acá, en vez de dejarlo para que un lector futuro lo infiera de los ejemplos. Un tercer tipo de item compuesto que aparezca más adelante deberá declarar su propia aridad del mismo modo, en vez de asumir que hereda una de las dos ya existentes.

## Relacionados
- `relacionado-con` → [nested-field-continuation-line.md](nested-field-continuation-line.md) — ambas aridades siguen siendo líneas de continuación del mismo item, nunca viñetas anidadas.
- `relacionado-con` → [per-field-description-cap.md](per-field-description-cap.md) — el tope por campo sigue aplicando a las líneas de valor del brief, no solo a las del contrato.
