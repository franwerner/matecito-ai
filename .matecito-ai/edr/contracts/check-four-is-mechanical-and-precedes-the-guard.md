# EDR — Check 4 compara dos afirmaciones explícitas, mecánicamente, y vive antes del guard que calla

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
Un retorno con `### New Decisions: None.` en el cuerpo y `Summary: "Key Decisions: 2 documented"` pasaba
los primeros tres chequeos del Return Contract Check — presencia, títulos y forma por status — y
despachaba en silencio. El Unresolved Decisions Guard lee el sentinel de la sección y se queda callado
cuando está vacía, así que nadie ponía las dos afirmaciones lado a lado. El chequeo va acá, no en el
guard, porque este check corre ANTES de él y el guard es precisamente el que enmudece ante una sección
vacía.

## Decisión
El cuarto chequeo compara dos afirmaciones explícitas sobre la misma sección — lo que el `Summary` del
envelope dice y lo que el cuerpo de la sección realmente trae — nunca si la prosa del `Summary` es una
descripción fiel del cuerpo. Una sección con sólo el sentinel vacío se declara EMPTY; si el `Summary`
afirma una cuenta distinta de cero o la presencia de contenido para esa misma sección, eso es una
contradicción, y lo mismo al revés (filas reales con el `Summary` declarándolo en cero). Ningún otro
desacuerdo entre `Summary` y cuerpo dispara este chequeo.

## Reglas verificables
- **[manual]** `payload/domains/development/CLAUDE.md`, sección `### Return Contract Check (MANDATORY)`,
  item 4, declara que el chequeo compara dos afirmaciones explícitas sobre la misma sección — nunca si la
  prosa es fiel al cuerpo — y que ningún otro desacuerdo lo dispara.
- **[manual]** La misma sección declara que este chequeo corre antes del Unresolved Decisions Guard, que
  es el que se queda mudo ante una sección vacía.

## Alternativas consideradas
Delegar esta comparación al Unresolved Decisions Guard, ya que ese guard es quien lee las secciones de
todos modos — descartada: ese guard corre después de este check y precisamente calla cuando la sección
está vacía, que es el caso exacto que había que cazar.
Definir el chequeo como "¿el `Summary` es una descripción fiel del cuerpo?" — descartada: eso es
interpretación, lo único que este chequeo existe para no hacer; se define en cambio como comparar dos
afirmaciones explícitas y mecánicas.

## Consecuencias
Un retorno cuyo `Summary` contradice su propio cuerpo se detiene acá, mecánicamente, antes de que el
guard tenga oportunidad de leer una sección vacía y seguir de largo.
