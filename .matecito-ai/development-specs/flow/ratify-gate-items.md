# Capability — Ratificar items de gate uno a uno con evidencia anclada

- **Status:** Accepted
- **Date:** 2026-08-13
- **Components:** cli

## Propósito

Un gate entrega todo a la vez, y cada item llega sin nombrar sobre qué trata, así que el usuario no puede juzgar un item ni pedir detalle. Tras este cambio un gate abre con un índice de lo que ese retorno trae, camina los items uno a uno, y cada item nombra su fuente. La brevedad deja de ser una petición y se convierte en un límite medido.

## Actores

- **Orquestador**: abre cada gate con un índice único, presenta items uno a uno, y espera un resultado explícito de cada uno antes de mostrar el siguiente
- **Usuario**: confirma, ajusta o rechaza cada item; puede saltar al final con "confirmar el resto"
- **Fase SDD**: suministra cada item con resumen, ancla y razonamiento completo

## Precondiciones

- El retorno de la fase lleva items ratificables (Tier 1) en una o más secciones
- Cada item que declara split lleva tres partes: `summary`, `anchor` y `rationale`
- El `summary` no supera la longitud máxima declarada en el contrato de la sección

## Flujo principal

1. Orquestador recibe el retorno de fase o materializa un momento de decisión propio (Discovery Gate, Uncommitted-Work Gate, Review Workload Guard, validadores, o risks)
2. Contar los items ratificables (Tier 1) a presentar; si el conteo es **0**, no mostrar nada; si es **exactamente 1**, presentar el template solamente sin índice; si es **2 o más**, construir un índice único mencionando cuántos hay y agrupados por sección
3. Presentar el primer item con su summary y anchor solamente
4. Esperar resultado: confirmado, ajustado o rechazado
5. Si ajustado, registrar el texto corregido como el ratificado
6. Presentar el siguiente item sin mostrar al usuario el razonamiento completo a menos que lo pida
7. En cualquier momento (o antes del primer item), el usuario puede pedir "confirmar el resto" para aceptar todos los pendientes
8. Una vez que todos los items tienen resultado, el gate cierra

## Ramas / flujos alternativos

- **Sin items ratificables** → No mostrar índice ni walkthrough; continuar silenciosamente al siguiente paso
- **Usuario pide el razonamiento completo de un item** → Reproducir la línea `· rationale:` verbatim del bloque persistido; volver al mismo item sin marcar resultado
- **Usuario pide ver la fuente anclada** → Recuperar el material en esa ubicación y mostrarlo; volver al mismo item sin marcar resultado
- **Usuario rechaza un item** → Registrar el rechazo; continuar con el siguiente
- **Usuario invoca "confirmar el resto" antes del primer item** → Marcar todos los pendientes como confirmados; cerrar el gate sin presentar items individuales

## Casos borde

- **Un solo item ratificable** → El índice sigue mostrándose (estructura uniforme); la opción "confirmar el resto" sigue disponible
- **Items de varias secciones en el mismo retorno** → Un índice único contando ambas; el walkthrough es plano, sin agrupar por sección durante la presentación
- **El usuario ajusta un item a texto idéntico** → El ajuste se registra como tal (no se compara con el original)
- **El usuario pide detalle de un item ya decidido** → La decisión se mantiene; se devuelve el detalle sin cambiar el resultado

## Reglas de negocio

- El índice es único por retorno, nunca acumulativo across el flow
- La forma del gate es decidida por conteo (0/1/≥2) Y por si un item es compuesto. Un item cuyo contenido es un conjunto de fields tipados (un contrato) se presenta dentro del mismo template compartido que cualquier otro, como UN item, con sus fields mostrados debajo de su summary y encima de las acciones. Cuenta como un item en el índice y toma exactamente un resultado
- Item by item es el default cuando hay múltiples; "confirmar el resto" es el único atajo global, EXCEPTO: cuando hay dos o más contratos (items compuestos), el gate DEBE ofrecer ritmo (uno-a-uno o todos-a-la-vez) antes de mostrar el primero; esta oferta es la única excepción nombrada a "no hay otro bulk action"
- El template fijo (summary + anchor + acciones + opcionalmente fields) no tiene slot para narrativa libre
- Los slots del template se separan por línea en blanco; imprimirlos como líneas consecutivas es una violación del template, no una variante de estilo
- Cuando la respuesta que el item necesita es una elección entre alternativas cerradas, se presenta por el control de preguntas del host (summary como pregunta, cada alternativa como opción con una línea de su costo) en vez de prosa terminada en una lista entre paréntesis; el ancla viaja igual, dentro del texto de la pregunta. Aplica a TODO item cuya respuesta sea una elección, no sólo a ratificaciones. Cuando la respuesta es algo que el usuario redacta (una pregunta de discovery abierta, una corrección, un ajuste), no hay alternativas que enumerar y la forma es el template en prosa
- Lo que el usuario ve en el gate es solo lo que el contrato declara para esa sección (summary, anchor, tokens, fields si el item es compuesto)
- La rationale completa viaja siempre en el bloque persistido; nunca se imprime por defecto
- El ancla DEBE ser suministrado por la fase; nunca se deriva
- El Discovery Gate es la única excepción al requisito de anchor (porque sus items son preguntas sobre un pedido, no sobre artefactos)
- Risks se presentan informativaentes: llevan anchor y forma de items, pero no bloquean el gate ni contabilizan en el índice de items ratificables
- Seis momentos del orquestador (Discovery Gate, Uncommitted-Work Gate, Review Workload Guard, `blocked` returns, findings de validador, risks) presentan través de este mismo template y walkthrough, sin redacción propia

## Entidades y estados

- **Índice del gate** — Resumen de items ratificables en el retorno, agrupado por sección, conteo total. Estados: no generado → presentado → confirmado
- **Item** — Entrada de una sección de gate, con summary, anchor, rationale y tokens opcionales. Estados: pendiente → confirmado → ajustado → rechazado
- **Walkthrough** — La secuencia uno a uno de presentación de items. Estados: no iniciado → en curso (N items pendientes) → completado

## Errores de cara al actor

- **Item sin ancla**: la renderización de retorno falla antes de llegar al gate; exit 1
- **Summary sobre límite**: la renderización falla antes de llegar al gate; exit 1
- **El usuario no responde a un item**: el gate sigue esperando; no hay timeout silencioso

## Escenarios

### Scenario: Un índice por retorno, nunca acumulado en el flow

- **GIVEN** una sección de retorno con items ratificables y luego otro retorno de fase con más items
- **WHEN** se abre el first gate
- **THEN** el índice cuenta solo los items del primer retorno
- **AND** cuando se abre el segundo gate, su índice cuenta solo los suyos, sin repetir lo ya decidido

### Scenario: Item by item es el default y bloquea

- **GIVEN** un índice con tres items ratificables
- **WHEN** el walkthrough comienza
- **THEN** se presenta el primer item solo y el segundo no se muestra hasta que el primero tiene resultado
- **AND** nada downstream se despacha mientras quede un item sin resultado

### Scenario: Confirm the rest antes del first item

- **GIVEN** un índice que acaba de mostrarse
- **WHEN** el usuario pide confirmar todo
- **THEN** cada item se registra como confirmado y no se presenta ninguno individualmente

### Scenario: Confirm the rest a mitad del walkthrough

- **GIVEN** un walkthrough de siete items donde los dos primeros fueron decididos y se presenta el tercero
- **WHEN** el usuario pide confirmar el resto
- **THEN** los items tres a siete se registran como confirmados, el primero y segundo mantienen sus resultados, y no se presentan más items

### Scenario: El razonamiento se pide sin responder el item

- **GIVEN** un item siendo presentado
- **WHEN** el usuario pide por qué
- **THEN** se reproduce la línea `· rationale:` verbatim del bloque, y el mismo item sigue awaiting decisión

### Scenario: La fuente anclada se pide sin responder el item

- **GIVEN** un item anclado a una ubicación en el repositorio
- **WHEN** el usuario pide verlo
- **THEN** el material en esa ubicación se recupera y muestra, y el mismo item sigue awaiting decisión

### Scenario: Material fuera de los slots fijos no se presenta

- **GIVEN** un retorno que lleva prosa beyond the declared sections del template
- **WHEN** su gate presenta un item
- **THEN** solo el summary, anchor y acciones disponibles se muestran
- **AND** la narrativa libre no llega al gate

### Scenario: Dos gates, la misma presentación

- **GIVEN** dos gates distintos cada uno presentando un item ratificable
- **WHEN** cada uno abre
- **THEN** los slots ofrecidos y acciones disponibles son idénticos en ambos, y ninguno enuncia su propia redacción

### Scenario: El modo unattended sigue caminando

- **GIVEN** un cambio corriendo en modo de ejecución unattended
- **WHEN** un gate abre con items ratificables
- **THEN** el índice y el walkthrough item-by-item ocurren exactamente como en modo interactive

### Scenario: Un ajuste reemplaza el contenido del item

- **GIVEN** un item que el usuario corrige en lugar de confirmar tal cual ofrecido
- **WHEN** el walkthrough avanza
- **THEN** el contenido ajustado es lo registrado como ratificado, y el texto originalmente ofrecido no se registra

### Scenario: Items informativos quedan fuera

- **GIVEN** un retorno con dos items ratificables y uno informativo
- **WHEN** su gate abre
- **THEN** el índice cuenta dos, el walkthrough presenta dos, y el item informativo aparece solo en el resumen entre-fases

### Scenario: Items informativos también anclan

- **GIVEN** un item informativo en el resumen entre-fases
- **WHEN** se muestra
- **THEN** nombra la fuente sobre la que es, y nada lo aguarda

### Scenario: Tres gates ratifican por este walkthrough

- **GIVEN** el gate de decisiones pendientes, el gate de scope-confirmation, y ambos gates de confirmación minada
- **WHEN** cada uno abre
- **THEN** todos presentan items través del mismo template compartido con índice único y walkthrough uno-a-uno

### Scenario: El orquestador mismo presenta momentos de decisión

- **GIVEN** un momento donde el orquestador elige presentar algo (Discovery Gate, Uncommitted-Work Gate, Review Workload Guard, validador que falla, risks)
- **WHEN** ese momento abre
- **THEN** presenta através del mismo template compartido (index + walkthrough si ≥2 items, template solo si 1, silencio si 0)
- **AND** el índice y walkthrough tienen la misma forma que un gate de fase

### Scenario: Cero items ratificables, silencio total

- **GIVEN** un gate cuyo material está vacío o es solo informativo
- **WHEN** se intenta presentar
- **THEN** nada se muestra y el mechanism se omite
- **AND** el siguiente paso continúa sin mención del gate

### Scenario: Un solo item, template sin índice

- **GIVEN** un gate con exactamente un item (un blocker a fase, un resultado de validación, o un risk)
- **WHEN** se abre
- **THEN** se muestra el item través del template (summary, anchor, acciones)
- **AND** no aparece ningún índice

### Scenario: El mismo momento en dos cambios produce distinto conteo

- **GIVEN** el Discovery Gate en dos changes distintos, uno con 3 preguntas y otro con 0
- **WHEN** cada uno abre
- **THEN** el primero muestra índice + walkthrough; el segundo muestra una sola línea confirmando la lectura del pedido
- **AND** nada en el código clasificó "Discovery" como multi-item: la forma cambió solo por el conteo

### Scenario: Discovery Gate no declara anchor

- **GIVEN** el Discovery Gate presentando preguntas sobre un pedido sin artefactos aún
- **WHEN** se redactan sus items
- **THEN** ninguno lleva anchor (porque no hay ubicación recuperable que apuntar)
- **AND** el archivo de presentación declara esta excepción explícitamente, no por defecto

### Scenario: Findings de validador, mucho o poco

- **GIVEN** un validador que falla emitiendo hallazgos
- **WHEN** se presentan
- **THEN** si hay múltiples, abren con índice + walkthrough; si hay uno, solo el template; si cero, nada
- **AND** nada cambió en el validador o el script: el conteo decidió la forma

### Scenario: Risks no bloquean pero sí anclan

- **GIVEN** un retorno de fase con dos items ratificables y tres risks
- **WHEN** su gate abre
- **THEN** el índice cuenta dos (solo ratificables), el walkthrough presenta dos, y los risks se muestran junto al resumen sin abrir un nuevo walkthrough
- **AND** cada risk nombra su fuente (un archivo, un artefacto, la fuente que lo motivó)

### Scenario: Un retorno es solo risks

- **GIVEN** un retorno cuyo único contenido reportable son sus risks
- **WHEN** se presenta
- **THEN** no abre índice ni walkthrough
- **AND** los risks se muestran con sus anchors

### Scenario: Un contrato con muchos fields

- **GIVEN** un gate presentando un contrato que lleva nueve fields
- **WHEN** se camina
- **THEN** es un item tomando un resultado, y ningún field toma resultado de su propio

### Scenario: Items compuesto e ordinarios en el mismo retorno

- **GIVEN** un retorno llevando dos contratos y tres items ordinarios
- **WHEN** su gate abre
- **THEN** un índice cuenta cinco y el walkthrough es plano, cada contrato mostrado con sus propios fields

### Scenario: La cuarta forma vive con las otras tres

- **GIVEN** la declaración única de cómo un gate presenta items
- **WHEN** la forma compuesta se busca
- **THEN** se declara ahí, junto a las formas decididas por conteo, y ninguna otra declaración de ella existe

### Scenario: Tres contratos llegan a un gate

- **GIVEN** un retorno llevando tres formas de contrato
- **WHEN** el gate abre
- **THEN** declara cuántas hay y ofrece uno-a-uno o todos-a-la-vez antes de mostrar el primero

### Scenario: La excepción está escrita donde la regla que exceptúa está

- **GIVEN** la declaración que ningún otro bulk action existe
- **WHEN** se lee tras este cambio
- **THEN** la oferta de ritmo de contrato se nombra ahí como su excepción

### Scenario: Confirmar el resto sigue disponible durante el walk de contratos

- **GIVEN** un walkthrough de formas de contrato con dos ya decididas
- **WHEN** el usuario pide confirmar el resto
- **THEN** las indecididas se registran como confirmadas y las decididas mantienen sus resultados

### Scenario: Dos declaraciones, dos unidades

- **GIVEN** la declaración que items se caminan uno a la vez y la declaración que un contrato nunca se propone field-por-field
- **WHEN** se leen juntas
- **THEN** cada una nombra la unidad que gobierna y ninguna se lee como excepción a la otra

## Referencias

- **Contrato compartido** → [`../../shared/references/gate-presentation.md`](../../shared/references/gate-presentation.md) — El walkthrough (índice → uno a uno → confirmar-el-resto), el template fijo de slots (summary, anchor, acciones, sin narrativa), los nueve momentos (tres gates de fase + seis momentos de orquestador), la regla de conteo (0/1/≥2), y la excepción de anchor del Discovery Gate
- **Guard de orquestador** → [`../../payload/domains/development/CLAUDE.md`](../../payload/domains/development/CLAUDE.md) — Los seis momentos de orquestador (Discovery Gate, Uncommitted-Work Gate, Review Workload Guard, `blocked` returns, findings de validadores, risks) citan el archivo compartido
- **Anchoring criterion** → [`../../payload/domains/development/skills/gentle-ai/_shared/sdd-phase-common.md`](../../payload/domains/development/skills/gentle-ai/_shared/sdd-phase-common.md) Section D.3 — Formas legales de anchor (`<repo-path>[:line]` | `<engram-key>`), start-line-only, regla target-not-yet-written; Discovery Gate es la excepción explícita
