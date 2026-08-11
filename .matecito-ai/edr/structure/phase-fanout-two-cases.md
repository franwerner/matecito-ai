# EDR — La prosa de fan-out nombra dos casos, cada uno con su propia definición

- **Status:** Accepted
- **Date:** 2026-08-10

## Contexto

Hasta esta decisión, la única fase del pipeline que se despachaba como más de un agente era la de
verificación, y la prosa que lo documentaba lo decía explícitamente: excepción única, nadie más fanout-ea.
La fase de implementación gana ahora su propio mecanismo de despacho concurrente (batch paralelo con
corridas aisladas más una consolidación) — un segundo caso real, con su propia forma. Documentarlo como
una excepción aparte, sin tocar la cláusula existente, dejaría dos prosas que se contradicen: una dice
"la única excepción soy yo", la otra dice "yo también fanout-eo".

## Decisión

La sección de fan-out nombra **los dos casos declarados** —verificación e implementación—, cada uno con
su propio archivo de definición, y conserva la cláusula de exclusividad **re-apuntada a ambos**: sigue
sin ser un patrón general disponible para cualquier fase, y un tercer caso necesita su propio cambio.

## Reglas verificables

- **[manual]** La sección de fan-out nombra exactamente dos casos, cada uno con su archivo de definición propio.
- **[manual]** La cláusula de exclusividad sigue presente y re-apuntada a ambos casos: ningún texto la deja leer como "patrón general" ni como "excepción de una sola fase".

## Alternativas consideradas

- **Dejar la prosa acotada a verificación y documentar el batch paralelo aparte, sin tocarla.** Descartado: produce dos afirmaciones contradictorias sobre cuántas fases fanout-ean, y quien lea solo una de las dos queda con información falsa.
- **Soltar la cláusula de exclusividad al agregar el segundo caso.** Descartado explícitamente: la cláusula es lo que impide que un tercer caso se cuele sin su propio cambio; perderla convierte el fan-out en un patrón ofrecido a cualquier fase, que es justo lo que este mecanismo no es.

## Consecuencias

- Cada caso mantiene su propia forma y su propio archivo de definición (`subverifier-groups.md`, `parallel-batch.md`); la prosa compartida es solo el nombre de la sección y la cláusula de alcance, no el mecanismo.
- Agregar un tercer caso de fan-out en el futuro requiere tocar esta decisión explícitamente, no solo agregar prosa suelta en la fase que lo quiera.

## Relacionados

- `relacionado-con` → [dispatch-batch-bound-integration.md](dispatch-batch-bound-integration.md) — el mecanismo que esta prosa nombra como el segundo caso.
- `relacionado-con` → [consolidation-run-is-the-integrator.md](consolidation-run-is-the-integrator.md) — quién ejecuta ese segundo caso.
