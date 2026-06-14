# Qué es (y qué no es) un DDR

Referencia canónica del **concepto** de DDR (Design Decision Record). Es la fuente de verdad de la *idea*; cualquier skill o agente que trabaje con DDRs apunta acá en vez de redefinirla. La *estructura/plantilla* concreta se define por separado; esto define qué cuenta como DDR y qué no.

## Qué ES un DDR

Un DDR captura una **decisión** de diseño: una elección deliberada entre alternativas, con su razón. Tres rasgos lo definen:

1. **Decide** algo — elige una opción frente a otras (una paleta, una escala tipográfica, una grilla, un set de tokens, un componente, un tono de marca, un objetivo de accesibilidad). No describe ni produce: elige.
2. **Tiene un porqué** — el contexto y el trade-off que justifican la elección.
3. **Perdura como restricción** — gobierna el diseño futuro; quien produzca piezas después la respeta, y se puede chequear contra Figma si una pieza se desvió de ella.

Cubre tres sabores, todos "decisiones" en sentido amplio:
- **decisión** — un trade-off real (p. ej. qué paleta primaria, qué familia tipográfica).
- **convención** — un acuerdo de estilo o estructura (p. ej. cómo se nombran los tokens, qué escala de espaciado).
- **política** — una regla verificable (p. ej. ratio de contraste mínimo, objetivo WCAG).

Responde: *"qué decidimos, por qué, y qué gobierna de acá en adelante"*.

## Qué NO es un DDR

- **No es una pieza ni un asset.** Un mockup, un export, un frame concreto se produce y se entrega; no decide nada. La mayoría del trabajo visual es ejecución y no merece DDR. Forzar un DDR por cada pieza es ruido.
- **No es una verificación.** Un criterio de aceptación o un check de contraste confirma que algo cumple ("style → ratio ≥ 4.5:1"). No tiene porqué ni alternativas, y no perdura como restricción. Verifica; no decide.
- **No es una señal ni un recordatorio.** Un "falta definir el dark mode", un TODO o una nota *apuntan* a una decisión ausente — son el dedo señalando, no la decisión.
- **No es un valor incidental.** Un color o un px que aparece pocas veces o sin intención clara —un valor suelto en un mockup sin token que lo ancle— no es una decisión, es ruido. Sin evidencia de elección deliberada, no hay DDR.
- **No es el "cómo".** El detalle de producción vive en el archivo Figma y en la pieza. El DDR captura el *porqué* de una elección, no el paso a paso de cómo se armó el frame.
- **No es un porqué adivinado.** Inferir la razón de una decisión sin que conste es inventar. El porqué lo aporta una persona.

## Estados (un DDR no es estático)

Un DDR tiene ciclo de vida. Dos estados importan para el concepto:

- **Borrador / inferido** — se registró el **QUÉ** (qué se eligió, observable como evidencia desde Figma) pero el **PORQUÉ está vacío**: es *evidencia de que se tomó una decisión*, sin la razón ratificada. Es un candidato, no una decisión cerrada.
- **Aceptado** — una persona lo ratificó y completó el porqué (contexto, alternativas, consecuencias). Recién ahí es un DDR pleno.

La promoción de borrador a aceptado siempre la hace una persona; nunca es automática. Un borrador todavía no "decide" — solo deja constancia de que hay algo para decidir.
