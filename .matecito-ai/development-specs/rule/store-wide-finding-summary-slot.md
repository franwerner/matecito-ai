# Capability — Slot de resumen de hallazgos de alcance store

- **Status:** Accepted
- **Date:** 2026-09-03
- **Components:** cli

## Propósito

Que un hallazgo de alcance store —el resultado de validar un store completo— tenga un lugar propio y declarado en el retorno de la fase, separado de la columna que describe la divergencia de una sola fila, de modo que lo heredado y lo que abrió el cambio queden contables por separado.

## Actores

- **Fase verify**: emite un store-wide summary line en su retorno cuando la validación de una capacidad-spec completó
- **Sub-verifier**: proporciona el resultado de haber validado el store entero
- **Renderizador**: publica el slot y lo emite en su lugar declarado, incluso cuando la tabla está vacía
- **Orquestador**: captura el resultado del store-wide summary y lo persiste en el reporte

## Precondiciones

- Una sección de retorno de fase declara un slot de pie de tabla (`footer`) con una etiqueta y una clave
- La sección tiene datos tipados (tabla o labeled-list)
- El sub-verifier que produce la sección suministra la clave en su resultado

## Flujo principal

1. El contrato de retorno de una sección declara un `footer` con `label` y `field`
2. El sub-verifier produce su salida incluyendo el valor para esa clave
3. El renderizador valida que la clave está presente en los datos
4. El renderizador emite el pie de tabla debajo de la tabla, con su valor
5. El orquestador persiste el resumen de alcance store en el bloque completo

## Ramas / flujos alternativos

- **La tabla está vacía pero el footer tiene `on_sentinel: true`** → el renderizador emite `None.` seguido del pie de tabla, para distinguir "se validó y salió limpio" de "no se validó"
- **La sección no está presente** → todo lo que gatea sobre su presencia (incluyendo el footer) se omite

## Casos borde

- **Fragmento que no trae la clave del footer** → renderizador falla nombrando la clave faltante; exit 1, stdout vacío
- **Footer con `on_sentinel: true` pero tabla vacía** → se emite el sentinel `None.` seguido de una línea en blanco y luego el pie de tabla
- **Dos secciones con footers, uno con `on_sentinel` y otro sin** → solo el primero sobrevive al `None.` case; el segundo sigue la regla anterior
- **La tabla tiene filas pero el footer está ausente de los datos** → falla de renderización nombrando el campo

## Reglas de negocio

- El slot de pie de tabla es un campo tipado que el sub-verifier debe suministrar cada vez que emite la sección, sea cual sea el contenido de la tabla
- El campo es emitido siempre, incluso cuando la tabla está vacía: su ausencia hace indistinguible "se validó y salió limpio" de "no se validó"
- El pie de tabla emite una línea de alcance store que cuenta solo lo que la validación del store reporta, y NO se mezcla con hallazgos de la sección de items (que describen solo hallazgos que abrió el cambio)
- Una sección que no declara el footer no produce un campo correspondiente; el sub-verifier que la produce no lo suministra
- Cuando la tabla está vacía y la sección tiene `sentinel: true`, el pie de tabla sobrevive la línea `None.` solo si su declaración incluye `on_sentinel: true`; sin esa bandera, la tabla vacía se emite como bare `'None.'`
- El pie de tabla es renderizado por el mismo renderizador que renderiza las filas de la tabla; no hay renderizado diferenciado por tipo de footer
- La clave del footer es suministrada por el productor, no derivada

## Entidades y estados

- **Sección de retorno** — sección que emite datos tipados (tabla, labeled-list). Puede opcionalmente declarar un `footer`. Puede tener `sentinel: true` para emitir `None.` cuando está vacía.
- **Footer** — declaración de un pie de tabla con `label` (texto a imprimir), `field` (clave de dónde leer el valor), y opcionalmente `on_sentinel` (booleano para que sobreviva el `None.`). Emitida por el renderizador debajo de la tabla, con `label` seguida de `:` y el valor.

## Errores de cara al actor

- **Campo del footer ausente en los datos** → no se emite bloque, exit 1, stderr nombra la sección y el campo faltante
- **Tabla vacía con footer declarado pero sin `on_sentinel`** → se emite bare `None.` sin el pie de tabla

## Escenarios

### Requisito: Un hallazgo de alcance store tiene su propio lugar declarado

Cuando una fase debe reportar el resultado de una validación que abarca un store completo —no un elemento de la tabla—, el contrato de esa sección DEBE declarar un lugar propio para esa línea, debajo de la tabla y con su etiqueta, y la fase DEBE emitirla ahí. La columna de detalle por fila NO DEBE transportarla: esa columna está reservada a la divergencia de un único elemento, y usarla para un resumen del store entero hace que ninguna de las dos cosas se pueda leer.

Cuando la sección se emite, la línea de alcance store se emite siempre, incluso cuando la validación no encontró nada. Un chequeo cuyo único rastro sería un hallazgo es indistinguible de un chequeo que nunca corrió.

#### Scenario: el store validado tiene hallazgos

- **GIVEN** una sección de coherencia cuyo store fue validado y arrojó hallazgos
- **WHEN** la fase emite su bloque de retorno
- **THEN** debajo de la tabla aparece una línea etiquetada con el resultado de alcance store y su conteo
- **AND** ninguna celda de detalle por fila lleva ese conteo

#### Scenario: el store validado no tiene hallazgos

- **GIVEN** una sección de coherencia cuyo store fue validado y salió limpio
- **WHEN** la fase emite su bloque de retorno
- **THEN** la línea de alcance store igual aparece, reportando cero
- **AND** el lector puede distinguir "se validó y no había nada" de "no se validó"

#### Scenario: una fila con divergencia puntual

- **GIVEN** un elemento del store que diverge de lo que el cambio implementó
- **WHEN** la fase emite su bloque de retorno
- **THEN** la celda de detalle de esa fila describe únicamente esa divergencia
- **AND** el resumen del store no aparece en ninguna celda

#### Scenario: el store no está presente

- **GIVEN** un proyecto sin ese store
- **WHEN** la fase emite su bloque de retorno
- **THEN** la sección entera no se emite, y con ella tampoco la línea de alcance store

### Requisito: Lo heredado se cuenta aparte de lo que abrió el cambio

La línea de alcance store DEBE contar sólo lo que la validación del store reporta, y NO DEBE mezclarse con los hallazgos que produjo el cambio en curso, que viajan en la sección de hallazgos del mismo bloque. Un lector DEBE poder decir, sin abrir nada más, cuántos hallazgos son heredados y cuántos los abrió este cambio.

#### Scenario: hallazgos heredados y hallazgos propios en el mismo retorno

- **GIVEN** un store con hallazgos previos y un cambio que además abrió hallazgos propios
- **WHEN** la fase emite su bloque de retorno
- **THEN** la línea de alcance store reporta sólo los heredados
- **AND** la sección de hallazgos reporta sólo los que abrió el cambio
- **AND** los dos conteos se leen por separado

#### Scenario: store limpio y cambio con hallazgos propios

- **GIVEN** un store sin hallazgos previos y un cambio que abrió hallazgos propios
- **WHEN** la fase emite su bloque de retorno
- **THEN** la línea de alcance store reporta cero
- **AND** los hallazgos del cambio siguen apareciendo íntegros en su propia sección

### Requisito: Quien produce la sección conoce el lugar que debe llenar

Cuando el contrato agrega un lugar nuevo que la fase debe llenar, la definición del productor que arma esa sección DEBE enumerar la clave nueva entre las que devuelve, y el método de la fase DEBE apuntar a ese lugar en vez de a la columna de detalle. Las dos cosas DEBEN cambiar junto con el contrato, no después: un productor que no sabe de la clave hace fallar la construcción del bloque final, cuando todos los productores ya terminaron su trabajo.

Un fragmento que no trae la clave DEBE hacer fallar la construcción del bloque nombrando la clave faltante, sin emitir bloque.

#### Scenario: el productor lee su propia definición

- **GIVEN** el productor que arma la sección de coherencia
- **WHEN** consulta qué claves debe devolver
- **THEN** la clave del resumen de alcance store está enumerada entre ellas

#### Scenario: el método de la fase apunta al lugar nuevo

- **GIVEN** el paso del método que ordena reportar el resultado de validar el store
- **WHEN** se lo lee
- **THEN** indica emitir la línea en el lugar declarado para ella, y no en la columna de detalle de una fila

#### Scenario: un fragmento que no trae la clave

- **GIVEN** un fragmento de la sección que omite la clave del resumen de alcance store
- **WHEN** se construye el bloque de retorno
- **THEN** la construcción falla nombrando la clave faltante; no se emite bloque; salida 1

## Referencias

- **Rule** → [`mailbox-item-summary-rationale-split.md`](mailbox-item-summary-rationale-split.md) — el mismo bloque de retorno transporta items con split; ese contrato gobierna las partes de cada item, este gobierna el resumen de alcance store que va debajo de la tabla. No se solapan: uno es por item, el otro es por store.
