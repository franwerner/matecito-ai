# REFACTOR — matecito-ai

Documento de trabajo. Reúne cuatro frentes de refactor con su alcance medido sobre el repo.
Nada de esto está ejecutado todavía: es el material para decidir qué se hace y en qué orden.

**Estado:** propuesto · **Fecha:** 2026-08-14 · **Rama:** `master`

---

## 1. Eliminar el dominio `design`

El dominio deja de mantenerse. Se retira del payload, del núcleo, de la documentación y de los tests.

### Alcance medido

| Qué | Dónde | Tamaño |
|---|---|---|
| El dominio completo | `payload/domains/design/` | 70 archivos, 552 KB |
| Cláusula de override del mine gate | `payload/core/CLAUDE.md` (≈599-602) | 1 bloque + su comentario |
| Portada y tabla de dominios | `README.md` líneas 17, 38, 58, 76, 83, 121, 160 | 7 menciones |
| Guía de auto-mine | `docs/guide/05-auto-mine.md` líneas 15-17, 24 | la sección `design` entera |
| Guía de configuración | `docs/guide/08-configuracion.md` | menciones sueltas |
| Comentario del manifest | `internal/manifest/manifest.go:20` (nombra DDR/design) | 1 línea |
| Tests del motor Go | `internal/manifest/manifest_test.go`, `internal/setup/deploy/deploy_test.go`, `internal/setup/deploy/registry_test.go` | 3 archivos |

Se van con él: los 10 agentes `design-*`, las skills de fase (`design-phases/*`), las dos skills de
decisiones (`design-decisions-mine`, `design-decisions-validate`), el store de DDR con sus templates,
`design-principles`, y las dependencias que el dominio declaraba en su `manifest.json` — los MCP
`figma` y `canva`, que ningún otro dominio pide.

### Consecuencia de segundo orden: el mine gate del núcleo queda huérfano

`payload/core/CLAUDE.md` documenta, en su propio comentario de mantenimiento, que `development`
declara su mecanismo de captura de decisiones in-flow *precisamente* para que `design` siguiera usando
el mine gate genérico del núcleo — y que borrar esa sección "habría roto `design` en silencio".

Sin `design`, la **Decision-Gap Capture (mine gate)** del núcleo se queda sin un solo consumidor:
`development` la anula por su cláusula de override. Son unas 25 líneas de kernel más su sección en
`docs/guide/05-auto-mine.md` que pasan a ser código muerto.

**No lo doy por decidido:** retirarlo es un cambio de núcleo con su propio alcance, y podría querer
conservarse como punto de extensión para un dominio futuro. Queda listado como decisión abierta (§5).

### Riesgo

Bajo y contenido. El dominio es un plugin: el microkernel está diseñado para que un dominio ausente no
afecte a los demás. El único punto de contacto real con el núcleo es la cláusula de override citada
arriba; el resto es documentación y tests.

---

## 2. QMD como motor de búsqueda

[`tobi/qmd`](https://github.com/tobi/qmd) — MIT, SQLite FTS5 + sqlite-vec, BM25 + vectorial +
re-ranking por LLM local, con servidor MCP (`query`, `get`, `multi_get`, `status`) — está **adoptado**
como buscador sobre el corpus markdown del proyecto. No es una propuesta: `process/configure-record-
search.md` y `flow/find-durable-records.md` (ambos `Accepted`, 2026-09-02) ya fijan cuándo se
configura y cuándo se usa.

### Costos asumidos

- ~1,9 GB de modelos GGUF en `~/.cache/qmd/models/`, descargados una vez.
- El modo híbrido carga tres modelos por invocación salvo que corra como daemon HTTP (`qmd mcp --http --daemon`).
- El índice se desactualiza al editar un `.md`; el disparador de reindexado es el que
  `configure-record-search.md` ya declara: correr de nuevo cuando cambian los stores del proyecto.

### El punto a resolver — cerrado

El pedido original incluía *quitar el indexado que existe hoy para la búsqueda*. Al medirlo en su
momento, ese indexado no hacía búsqueda: `.matecito-ai/development-specs/process/index-decision-
records.md` definía un índice de **estado** más un **versionado copy-on-write lazy** — pin de versión
exacta por evento, contenido content-addressable, soft-delete por `(proyecto, rama)` — sin relación
con capacidades de búsqueda, y dependía de cuatro EDRs y tres specs del broker.

Esa pregunta no se respondió: **se cerró**. `index-decision-records` y el broker que lo implementaba se
retiraron junto con el resto del cockpit (`sdd/drop-apps-subtree`), así que la contraparte que este
punto discutía ya no existe. No queda nada que argumentar acá.

---

## 3. Ordenar los templates del payload

Cuatro defectos concretos, de más de fondo a más cosmético.

### 3.1 El idioma se contradice consigo mismo

`payload/core/CLAUDE.md:115` fija: *"Everything authored in the payload, and everything a phase emits,
is written in **English**"*. Pero `capability.md` y `edr.md` clavan los títulos de sección en español
(`## Propósito`, `## Casos borde`, `## Decisión`, `## Alcance`, `## Reglas verificables`) y toda su
prosa guía también. Los `.yaml` agravan la mezcla: comentados en inglés, fijan los mismos encabezados
en español.

O los records durables son una excepción declarada a esa regla, o los templates la incumplen. Hoy no
está escrito cuál de las dos, y cada material nuevo hereda la ambigüedad. **Es el que más desordena**,
porque decide el contenido de todo lo que se escriba después.

### 3.2 `tech-edr` quedó fuera del motor

`validate-artifact.js` reconoce exactamente dos artefactos: `edr → edr.yaml` y
`capability-spec → capability.yaml` (líneas 39-40). `tech-edr.md` tiene forma propia — `Category`,
`Version`, `Decided in phase`, y `Status: Accepted` literal en vez de placeholder — y es el único
artefacto de development con template pero sin `.yaml`.

Nadie lo renderiza ni valida su forma: los EDR de `tech/` se escriben a mano y se confía en la
disciplina, que es justo lo que el resto del sistema decidió no hacer. **Es el de mayor riesgo
silencioso.**

### 3.3 Dos escaleras de Status que no comparten casi nada

`capability.yaml` declara `[Inferred, Draft, Accepted, Deprecated]`; `edr.yaml`,
`[Inferred, Accepted, Pending, Deferred]`. Coinciden en dos de cuatro.

La asimetría es defendible — una decisión se aplaza, un comportamiento no; un EDR se edita en el lugar
en vez de deprecarse — pero ese razonamiento no está escrito en ningún lado, así que se lee como
descuido y el próximo artefacto no sabe cuál copiar. El arreglo es escribirlo, no unificar las
escaleras.

### 3.4 Un residuo mal pegado

La línea 1 de `capability.md` cierra un comentario y abre otro pegado en la misma línea, describiendo
un bug de rutas ya arreglado. La convención de notas `<!-- matecito-ai: -->` es deliberada y
sistemática — hay más de veinte en `payload/core/CLAUDE.md` — pero siempre en línea propia. Esta es un
artefacto de edición, y es lo primero que se ve al abrir el template.

### No son defectos

- Que el DDR no tenga `.yaml` es coherente: `payload/domains/design/` no tiene `scripts/` ni motor.
  Con el §1 ejecutado, el punto desaparece solo.
- Todas las rutas `~/.claude/references/...` citadas en el payload resuelven correctamente.

---

## 4. La asimetría entre silencio y emisión

El síntoma observado: el sistema carga contexto, lee Engram, guarda, y no devuelve nada al chat. Tiene
una forma medible sobre `payload/core/CLAUDE.md` (684 líneas).

**El silencio está especificado por mecanismo; la emisión, una sola vez y en abstracto.** Hay doce
mandatos de callar, repartidos por todo el kernel y pegados cada uno a su mecanismo concreto:

| Línea | Qué manda callar |
|---|---|
| 123 | piezas del ecosistema ausentes |
| 153, 649 | el gate de decision records, cuando el store está inactivo |
| 216, 396 | el spec-mine, según el flag |
| 364 | el guard temprano del intake |
| 385 | la corrida de init |
| 604 | el override de dominio |
| 625 | el mine gate cuando no dispara |

Enfrente, casi todo el lenguaje de emisión vive apiñado en **dos líneas** — 112 y 115 — enunciado en
general y sin nombrar un solo mecanismo. Cuando el agente está ejecutando el paso, la instrucción
cercana le dice `silently` y la que dice "salvo que haya una regla de emisión" está cuatrocientas
líneas más lejos y no habla de su caso. Gana la local, siempre.

**Y lo poco que sale es de baja utilidad.** La línea 122 fija los avisos internos como *"a single short
line in English"* — el ejemplo literal es `"Saved to Engram."` — y la 117 los exime a propósito de la
regla de idioma de la 115. En una conversación en español, lo único que devuelve el sistema es un
muñón en inglés que no dice qué guardó ni bajo qué clave.

**Encima, el contenido tiene prohibido cruzar el chat.** La línea 635 hace que el orquestador le pase a
los sub-agentes *"topic-key references, not content"*, y la 89 manda los deliverables a archivos y
prohíbe resumir después de un cambio salvo pedido explícito.

Son tres capas independientes empujando en la misma dirección, sin nada que las contrapese en el punto
donde se toma la decisión. Se cruza con §3.1: es el mismo eje — qué se emite y en qué idioma —
asomando por dos superficies distintas.

---

## 5. Decisiones abiertas

Bloquean el frente que nombran; el resto puede avanzar sin ellas.

1. **§1 — qué pasa con el mine gate del núcleo** una vez que `design`, su único consumidor, no existe:
   se retira, o se conserva como punto de extensión.

2. **§3.1 — cuál de las dos reglas de idioma gana**: los records durables son excepción declarada, o
   los templates se pasan a inglés.

---

## 6. Orden sugerido

1. **§3.4 y §3.3** — media hora entre los dos, cero riesgo, y dejan los templates presentables.
2. **§1** — la eliminación de `design`, que además resuelve solo el falso pendiente del DDR sin `.yaml`.
3. **§3.1 y §4** — juntos, porque son el mismo eje; deciden lo que se va a ver de todo lo demás.
4. **§3.2** — el contrato de `tech-edr`, una vez que el criterio de idioma esté fijado.
5. **§2** — al final; la adopción ya está cerrada, no quedan pasos pendientes que ordenar.
