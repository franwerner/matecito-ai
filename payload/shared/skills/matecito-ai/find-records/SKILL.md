---
name: find-records
version: 1.0.0
description: Localiza los records durables que gobiernan una parte del proyecto — decisiones de ingeniería (EDR), de diseño (DDR) y capability-specs de comportamiento — antes de escribir código, proponer un diseño o verificar contra ellos. Usá esta skill SIEMPRE que estés por tocar arquitectura, capas, contratos, datos, transporte, convenciones o comportamiento, y necesites saber qué records aplican; cuando una fase del flujo tenga que leer "los EDRs que este cambio toca" o "los capability-specs tocados"; o cuando quieras saber qué stores tiene un repo y qué gobierna cada uno. Combina el índice del store, la búsqueda literal y la búsqueda por significado, y devuelve un conjunto de candidatos para decidir sobre él — nunca un único resultado.
license: MIT
metadata: {"hermes":{"tags":["records","edr","ddr","specs","search","discovery"],"category":"cross-domain","related_skills":[]},"author":"matecito-ai","version":"1.0.0"}
---

# Find Records

Encuentra los records durables que gobiernan lo que estás por tocar.

Un proyecto matecito-ai guarda dos clases de conocimiento durable: **decisiones**
(por qué se eligió cada cosa — `edr/` en development, `ddr/` en design) y
**comportamiento** (qué hace el sistema — `development-specs/`). Escribir código
sin haber leído el record que lo gobierna es el modo de falla que esta skill
existe para cerrar.

**El problema no es buscar: es no saber que había algo que buscar.** Un record
que nadie consultó no deja rastro. Por eso el procedimiento no termina en "lo
encontré" sino en "miré el conjunto completo de candidatos y decidí".

## Paso 1 — Mapear los stores (siempre)

```sh
node ~/.claude/skills/find-records/scripts/map-records.js
node ~/.claude/skills/find-records/scripts/map-records.js --component api,ui --json
```

Devuelve qué stores existen, qué gobierna cada uno, cuántos records tiene, dónde
está su índice y **con qué nombre buscar en él**.

Ese nombre no lo deduzcas ni lo escribas de memoria: usá el que devolvió el
script. Un nombre equivocado no da error — devuelve nada, y eso se lee igual que
"no hay records que apliquen".

**Un cambio que toca varios componentes se rige por la unión de sus stores más
los de la raíz**, que gobiernan a todos — salvo que el proyecto haya declarado
un componente raíz (`root: true`), en cuyo caso el store de la raíz pertenece
a ese componente y se filtra como cualquier peer: para incluirlo hace falta
sumar `root` a la lista. `--component` acepta varios separados
por coma (o repetido) y suma en vez de elegir: filtrar por uno solo cuando el
cambio toca dos deja afuera records que aplican. Un componente que no existe en
la configuración se ignora con un aviso, nunca en silencio.

No hace falta adivinar qué componentes toca el cambio: cuando el flujo está
activo, el brief de intake ya lo trae resuelto y confirmado por el usuario en su
línea `Components:`. Usá esa lista.

**Esto no se deduce mirando el árbol.** Un componente puede tener store propio
bajo su path o regirse por el de la raíz, y la diferencia no está declarada en
ningún lado. El script la resuelve igual siempre; el ojo, no. Su salida incluye
`componentsWithoutOwnStore` justamente para que nadie concluya que un componente
sin carpeta propia no tiene records que lo gobiernen — los tiene, en la raíz,
**mientras el proyecto no haya declarado un componente raíz**. Si lo declaró,
el store de la raíz pertenece a ese componente (`root`) y deja de cubrir
automáticamente a los demás: para un componente sin store propio, sumá `root`
a la búsqueda en vez de asumir que la raíz igual lo gobierna.

Corré esto **antes** de decidir qué leer. No asumas rutas ni supongas que el
único store es el de la raíz.

## Paso 2 — Reunir candidatos por tres vías

Las tres cubren puntos ciegos distintas y ninguna sola alcanza. Usá las que
estén disponibles y **acumulá** — no elijas una.

**El índice del store.** `INDEX.md` enumera: si el record existe, está la fila.
Es la única vía con cobertura garantizada, y la única que te dice qué hay que
*no* buscaste. Leelo siempre. Con varios stores, leé el de cada uno que aplique
según el paso 1.

**Búsqueda literal.** `grep -rn` sobre las rutas que devolvió el paso 1.
Encuentra el término exacto cuando lo recordás: es instantánea y precisa, y
falla entera cuando el record usa otras palabras que vos — o está escrito en
otro idioma que tu consulta.

**Búsqueda por significado — el MCP `qmd`.** Cruza el vocabulario: encuentra el
record por lo que trata aunque no compartas una sola palabra con él. Es la única
vía que resuelve el caso "no me acuerdo cómo se llama ni cómo está escrito".

Llamá su tool `query` con estos parámetros, que no son preferencias sino lo que
salió de medir sobre un store real:

- `searches: [{ type: "vec", query: "<tu pregunta en lenguaje natural>" }]`
- `limit: 5` como piso — con 3 se pierde uno de cada cuatro records.
- `collections: [...]` con los nombres que devolvió el paso 1 para los stores
  que aplican. Sin acotar, los records de decisiones le ganan el ranking a los
  de comportamiento — y además el índice es compartido entre proyectos, así que
  una búsqueda sin acotar puede devolver records de otro repo.

Los resultados traen `file` y `score`. El `score` ordena, no decide: leé los
cinco.

## Paso 3 — Decidir sobre el conjunto, no sobre el primero

Reunidos los candidatos, **leelos** y decidí cuáles aplican. Reglas que salen de
medir esto sobre un store real, no de la intuición:

1. **Pedile cinco resultados o más a la búsqueda, nunca tres.** Fue la mejora
   más grande de todas las probadas y no cuesta nada: el record correcto casi
   siempre está en la lista, lo que falla es suponer que está primero.
2. **Acotá la búsqueda al store que corresponde.** Sin filtro, los records de
   decisiones ganan el ranking incluso cuando buscás comportamiento.
3. **Excluí los `INDEX.md` de la búsqueda.** Son tablas de navegación: matchean
   con todo y ensucian el primer puesto. Al índice se lo lee entero, no se lo
   busca.
4. **El primer resultado no es la respuesta.** En la medición, el record exacto
   quedó primero en dos de cada tres consultas sobre decisiones y menos de la
   mitad sobre comportamiento — pero apareció entre los tres primeros en cuatro
   de cada cinco. La lista sirve; el primer puesto no.

## Cuando `qmd` no está disponible

Es **opcional por presencia**: si no tenés su tool entre las tuyas, o el
servidor no responde, no lo menciones y trabajá con el índice y grep. Todo el
procedimiento sigue en pie — pierde la vía que cruza vocabulario, no la
cobertura.

Lo que **no** es aceptable es callar que faltó: cuando la vía semántica no
estuvo, decilo al presentar los candidatos. Quien lea el resultado tiene que
poder distinguir "busqué por las tres vías y esto es lo que hay" de "busqué por
dos y puede haber quedado algo con otro vocabulario".

## Qué NO hace esta skill

**No reemplaza al índice.** El índice enumera y no omite; la búsqueda rankea, y
rankear admite omisiones que nadie ve. Si lo que necesitás es garantía de haber
visto todo lo que aplica, la garantía la da el índice.

**No decide por vos.** Devuelve candidatos. Cuál gobierna tu caso lo decidís
leyendo, y cuando un record `Accepted` contradice tu plan, la regla del
ecosistema es parar y preguntarle al usuario — no elegir en silencio cuál gana.

**No sirve como veredicto en verificación.** Verificar exige cobertura, y una
tasa de acierto alta no es cobertura: el record que no apareció no deja rastro
de su ausencia. En verificación los resultados son candidatos a leer, nunca
prueba de que no había nada más.
