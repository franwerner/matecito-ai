# EDR — El sub-verificador lee la regla de una sola línea en el envelope que ya lee primero

- **Status:** Accepted
- **Date:** 2026-09-03

## Contexto
Un sub-verificador de `sdd-verify` no le pide al renderer la forma de sus datos: arma su fragmento leyendo la prosa de su propio método y la definición de su grupo. Su propio archivo de agente (`agents/sdd-verify.md`) le dice explícitamente que la Sección D del protocolo compartido NO lo alcanza — así que apuntarlo ahí es señalarle una autoridad que se le dijo que ignore. La restricción de una sola línea (para `summary` y `rationale`) no estaba enunciada en ningún lugar que un sub-verificador efectivamente lea, mientras que el tope de caracteres sí lo estaba en dos notas cortas — una restricción presente y la otra ausente es peor que ninguna, porque leer una completa se lee como haber leído todas.

## Decisión
La regla completa se escribe en extenso una sola vez, en el envelope del Sub-Report (`subverifier-groups.md`) — el documento que todo sub-verificador ya lee antes de correr. Las dos notas cortas que ya existían junto al tope de caracteres (SKILL.md, pasos 3f y 6c) ganan la forma conjunta más un puntero a ese envelope, en vez de una tercera copia en extenso.

## Reglas verificables
- **[manual]** El envelope del Sub-Report (`subverifier-groups.md`) enuncia, en extenso, que `summary` y `rationale` son de una sola línea, con la consecuencia concreta (falla el render consolidado, salida 1, sin stdout) y el momento en que esa falla se descubre (recién al merge, después de que todos los grupos ya terminaron).
- **[manual]** Las dos notas de la SKILL (pasos 3f y 6c) que ya mencionaban el tope de caracteres pasan a mencionar también la regla de una sola línea, junto con un puntero al envelope por su ruta desplegada (`~/.claude/…`).
- **[manual]** Ningún sub-verificador es apuntado a la Sección D del protocolo compartido para esta regla — su propio agente le dice que esa sección no lo alcanza.

## Alternativas consideradas
Apuntar a `~/.claude/skills/_shared/sdd-phase-common.md` D.2: descartada, el agente del sub-verificador dice explícitamente que esa sección no lo alcanza. Sólo la SKILL (3f, 6c): descartada, deja al documento que todo sub-verificador lee primero — el envelope — sin ninguna mención de una restricción de la forma que declara. Tres copias en extenso: descartada, tres lugares para que diverjan, contra el patrón de un-solo-hogar que ya fijan dos EDRs Accepted.

## Consecuencias
Quedan dos menciones cortas (SKILL 3f, 6c) que son punteros, no autoridades — no fijan valor propio, sólo reenuncian lo que la renderización ya impone y dicen dónde está la versión completa. Un sub-verificador que escribe un `rationale` multi-línea recién lo descubre al fallar el merge consolidado, después de que todos los grupos ya entregaron su Sub-Report — esta prosa es la única barrera previa que existe antes de ese punto.

## Relacionados
- `relacionado-con` → [../contracts/subverifier-item-shape-single-declaration.md](../contracts/subverifier-item-shape-single-declaration.md) — fija que la forma de un finding se declara una sola vez, en el mismo envelope
- `relacionado-con` → [prose-register-single-home.md](prose-register-single-home.md) — el patrón general de un hogar en extenso y punteros en todos lados más — aplicado acá a esta regla concreta
