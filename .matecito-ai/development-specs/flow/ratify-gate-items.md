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

1. Orquestador recibe el retorno de fase con items ratificables
2. Si hay items, construir un índice único mencionando cuántos hay y agrupados por sección
3. Presentar el primer item con su summary y anchor solamente
4. Esperar resultado: confirmado, ajustado o rechazado
5. Si ajustado, registrar el texto corregido como el ratificado
6. Presentar el siguiente item sin mostrar al usuario el razonamiento completo a menos que lo pida
7. En cualquier momento, el usuario puede pedir "confirmar el resto" para aceptar todos los pendientes
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
- Item by item es el default; "confirmar el resto" es el único atajo global
- El template fijo (summary + anchor + acciones) no tiene slot para narrativa libre
- Lo que el usuario ve en el gate es solo lo que el contrato declara para esa sección (summary, anchor, tokens)
- La rationale completa viaja siempre en el bloque persistido; nunca se imprime por defecto
- El ancla DEBE ser suministrado por la fase; nunca se deriva

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

## Referencias

- **Contrato compartido** → [`../../shared/references/gate-presentation.md`](../../shared/references/gate-presentation.md) — El walkthrough (índice → uno a uno → confirmar-el-resto a mitad) y el template fijo de slots (summary, anchor, acciones, sin narrativa)
- **Guard de orquestador** → [`../../payload/domains/development/CLAUDE.md`](../../payload/domains/development/CLAUDE.md) — El Unresolved Decisions Guard cita el archivo compartido en lugar de enunciar el batching mechanic
- **Anchoring criterion** → [`../../payload/domains/development/skills/gentle-ai/_shared/sdd-phase-common.md`](../../payload/domains/development/skills/gentle-ai/_shared/sdd-phase-common.md) Section D.3 — Formas legales de anchor (`<repo-path>[:line]` | `<engram-key>`), start-line-only, regla target-not-yet-written
