# I6 — El costo de una corrida: el lane por fase y el largo de `apply`

**Precondición:** ninguna dura. **Riesgo:** medio — tocar los disparadores de lane cambia lo que se te
recomienda en cada cambio, y el sesgo de lane mínimo está ahí por una razón. **Tamaño:** mediano,
divisible: **una corrida por bloque**.

## La medición

`python3 scripts/sdd-cost-breakdown.py` sobre las corridas registradas. Total: **$187.24 en 2038
requests**, 7 agentes.

| agente | runs | $ | $/run | reqs/run | modelo |
|---|---|---|---|---|---|
| `sdd-apply` | 4 | **80.49** | 20.12 | 196 | sonnet |
| `sdd-verify` | 15 | **44.57** | 2.97 | 48 | sonnet |
| `sdd-design` | 3 | 31.09 | 10.36 | 66 | opus |
| `sdd-spec` | 3 | 17.47 | 5.82 | 40 | opus |
| `sdd-intake` | 4 | 8.70 | 2.18 | 22 | sonnet |
| `sdd-explore` | 2 | 3.81 | 1.91 | 33 | sonnet |
| `sdd-archive` | 1 | 1.12 | 1.12 | 69 | haiku |

`apply` + `verify` = **67% del total**. Los 15 runs de verify son instancias del fan-out de
sub-verificadores (6 grupos por despacho), no 15 despachos.

## Hallazgo 1 — el costo lo manda el arrastre de contexto, no el trabajo

La corrida grande de `sdd-apply` (487 requests, $54.67), por cuartos:

| cuarto | reqs | costo | prompt promedio |
|---|---|---|---|
| 1º | 121 | $9.67 | 185.675 |
| 2º | 121 | $10.88 | 272.788 |
| 3º | 121 | $14.94 | 379.611 |
| 4º | 124 | **$19.18** | 483.870 |

Mismo número de requests, **el doble de precio**. El prompt arrancó en 50k y terminó en 541k.

Se confirma entre corridas: **6,2× los requests cuesta 8,2×** — el precio por request sube 32% con el
largo de la corrida. No es lineal: partir `apply` en tandas más cortas no reparte el costo, lo **baja**.

## Hallazgo 2 — `tasks` no dispara nunca, y por eso `apply` corre en una sola tanda

`sdd-propose` y `sdd-tasks` tienen **0 corridas**. No es un bug: es mecánico. La regla del lane fork
dice *"escalate only for a **concrete, named reason** (an architectural decision, a large surface, or an
unclear area)"* — y esos tres disparadores mapean a **dos** add-ons: decisión arquitectónica → `design`,
área poco clara → `explore`.

**`propose` y `tasks` no tienen disparador propio.** Y como `custom` es *"base + sólo el add-on que ese
disparador necesita"*, un add-on sin disparador nunca entra en `custom`. Le queda `full`, que dice
textual *"do not recommend by default"*.

**La consecuencia de costo:** el fan-out paralelo de `sdd-apply` se dispara desde las marcas
`· parallel-group:` **del artefacto de tasks**. Sin fase `tasks` no hay marcas, y sin marcas `apply` no
puede paralelizar nunca. Por eso la corrida más cara fueron 487 requests en un solo batch.

**El sesgo de lane mínimo ahorra una fase de `tasks` y lo paga multiplicado en `apply`.**

## Hallazgo 3 — el campo que decide el lane no tiene ningún criterio

El `size` del cambio (`trivial` | `small` | `medium` | `large`) gobierna la escalación: habilita
`direct`, y es el **único disparador que por sí solo empuja a `full`** (*"several triggers or `large`
size → full"*).

Busqué en las 238 líneas de `sdd-intake/SKILL.md` cómo se establece. El paso 3 dice *"From the request
+ answers, classify"* y después lista las cuatro etiquetas. **No hay ningún test.** Más abajo hay una
exhortación — *"Be honest about size"* — que no es un criterio.

O sea: **el tamaño lo estima el modelo a partir de cómo el usuario describió el pedido**, sin mirar el
código. Es el campo más flojo del brief siendo el más consecuente. Y contamina el disparador de `tasks`:
"cuando el cambio da para varias tandas de apply" **es un juicio de tamaño**, así que se apoyaría en el
mismo vacío.

Hay un proxy de superficie que ya se calcula y no alimenta nada: el mapeo de `components` recorre el
alcance del pedido contra los `paths` declarados, y el propio skill dice que **ninguna fase lo lee**.

## Hallazgo 4 — el 42% de los requests no llama ninguna herramienta

205 de 487 en la corrida grande, y se repite: 42%, 41%, 44% en las otras tres. Son requests de puro
texto, y a 484k de prompt promedio son los más caros de la corrida. No producen nada verificable.

Sin diagnóstico: no está medido si es razonamiento necesario, narración, o reintentos. **Averiguarlo es
la primera tarea del bloque 3**, no una conclusión de este documento.

## Hallazgo 5 — dos add-ons sin ninguna corrida real

Que `propose` y `tasks` nunca hayan corrido significa que su contrato de retorno, sus buzones y los
guards que los leen **nunca se ejercitaron con trabajo de verdad**. Es deuda silenciosa: no falla nada
hasta que alguien pide `full`.

## Lo que NO es un problema

`sdd-init` con 0 corridas está bien. La clave `sdd-init/matecito-ai` existe en Engram (detectada
2026-08-19), el Init Guard la encuentra y no re-despacha. El bug conocido de esa guarda daba lo
contrario —re-despachar en cada comando—, así que cero es la señal buena.

---

## Bloque 1 — El lane se resuelve por fase, no de una en el intake

**El cambio de enfoque.** Hoy `sdd-intake` fija el lane entero al principio, apoyado en un `size` que
no tiene criterio (hallazgo 3). La alternativa: **intake clasifica el pedido a partir de la charla, y
cada fase habilita el add-on que sigue** — porque cada fase es el primer punto donde esa pregunta es
contestable con evidencia. `sdd-design` ya produjo el diseño y **sabe** si el trabajo tiene piezas
independientes; intake sólo tenía la descripción del pedido.

El mecanismo existe a medias: cada fase ya devuelve **`next_recommended`** en su contrato de retorno.
Hoy es orientativo; esto lo vuelve portante.

**El límite estructural, que define hasta dónde llega la idea.** Los add-ons tienen posición fija, y
una fase sólo puede habilitar uno que todavía no pasó:

```
intake → [explore] → [propose] → spec → [design] → [tasks] → apply → verify → archive
```

| add-on | quién puede decidirlo | hoy tiene disparador |
|---|---|---|
| `explore` | sólo intake — no hay nada aguas arriba | sí (área poco clara) |
| `propose` | intake, o `explore` si corrió | **no** |
| `design` | **`sdd-spec`**, que acaba de escribir el spec | sí (decisión arquitectónica) |
| `tasks` | **`sdd-design`**, que acaba de producir lo que se desglosaría | **no** |

**Resuelve `tasks`** — el que no tiene disparador y el que más plata cuesta. **No hace nada por
`propose`**, que sigue atascado en intake sin criterio, lo cual refuerza que puede sobrar.

**La fase DECIDE y REPORTA — no pregunta.** Si cada fase abre un gate "¿sumamos tasks?", se cambia una
pregunta al principio por N en el medio, que es lo contrario de todo lo demás que se está haciendo.
Bajo el criterio de I1 no hay nada que ratificar: habilitar un add-on no contradice ningún artefacto,
no es ambiguo, y la fase **puede** inferirlo — con más evidencia de la que tiene intake hoy. Va como
una línea del retorno.

**Lo que se pierde:** saber al principio cuánto va a costar el cambio. Aunque hoy esa previsión se
apoya en un tamaño adivinado, así que es previsibilidad de algo que no se está midiendo.

### Tareas

- [ ] **Leer `payload/docs/INDEX.md` antes de tocar archivos.**
- [ ] **Decidir el alcance**: si la resolución progresiva aplica a `design` y `tasks`, o si se deja
      `design` en intake (ya tiene disparador) y sólo `tasks` se mueve.
- [ ] **Escribir el criterio de `tasks` en `sdd-design`**: cuándo el diseño que acaba de producir da
      para varias tandas de `apply` o para trabajo paralelizable. **Es el punto del batch** — sin él,
      el fan-out paralelo sigue siendo una capacidad inalcanzable.
- [ ] **Volver `next_recommended` portante**, no orientativo: definir qué hace el orquestador con él y
      cómo se distingue de una recomendación que se ignora.
- [ ] **Reescribir el `Lane fork`** de `payload/CLAUDE.md` y el paso 4 de `sdd-intake/SKILL.md`: intake
      deja de fijar el lane completo y pasa a clasificar el pedido y habilitar sólo lo que le
      corresponde (`explore`, y `propose` si sobrevive).
- [ ] **Resolver `propose`**: darle un disparador o darlo de baja. Su candidato —más de un enfoque
      viable— se solapa con el de `design`, y ese solapamiento es parte de por qué nunca dispara.
- [ ] **Y sobre el `size`**: decidir si se le escribe un test, si se apoya en el mapeo de `components`
      que ya se calcula y nadie lee, o si con la resolución progresiva deja de importar tanto como para
      justificar el trabajo.

## Bloque 2 — Acotar el largo de `apply`

- [ ] **Tratar el largo de la tanda como decisión de costo, no sólo de alcance.** La capacidad de tanda
      de continuación ya existe; lo que falta es un criterio de cuándo cortar.
- [ ] **Definir dónde vive ese criterio**: en el desglose de `tasks` (cuántas tareas por tanda), en el
      orquestador al formar la ronda, o en `sdd-apply` al reportar `partial`.
- [ ] **Revisar el `Review Workload Guard`**, que hoy ya razona sobre tamaño de PR y tandas encadenadas
      — puede ser el lugar donde esto encaje sin inventar un guard nuevo.

## Bloque 3 — Los requests sin herramienta

- [ ] **Medir qué son.** Muestrear los requests sin tool-call de una corrida y clasificarlos:
      razonamiento necesario, narración, o reintento. **Sin esto no hay nada que decidir.**
- [ ] **Recién con el diagnóstico**, decidir si hay algo que cambiar y dónde.

## Bloque 4 — Ejercitar lo que nunca corrió

- [ ] **Una corrida real en lane `full`** sobre un cambio chico, sólo para que `propose` y `tasks`
      pasen por su contrato una vez. No es una prueba de humo: es la única forma de saber si sus
      buzones y sus guards funcionan.

---

## Lo que este batch corrige de lo ya planeado

**I2 no ataca este gasto.** El canal por Engram y el recorte de output bajan el contexto del
**orquestador**, y el 67% del costo
está adentro de los sub-agentes, cuyo prompt crece por su propio trabajo. Son $125 de $187 que el plan
actual no toca. Conviene tenerlo escrito para no atribuirles un ahorro que no van a producir.

**Lo único que sí toca algo es I2-C**: al recortar los `.md` de `phase-returns` —`sdd-apply.md` son 472
líneas— baja el
prompt de entrada, que hoy arranca en 47-50k antes de que la fase haga nada.

## Cierre

Una corrida donde `tasks` se recomiende por su propio disparador, `apply` corra en tandas acotadas, y
el costo por request no suba con el largo de la corrida. Y el desglose corriendo de nuevo, para
comparar contra los números de arriba en vez de contra una impresión.
