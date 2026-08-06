# Intake Brief — línea Components

- **Status:** Accepted
- **Date:** 2026-08-06
- **Components:** cli

## Propósito

El Intake Brief declara **qué superficies del repo toca un cambio**. Es la proyección **por-cambio** del set declarado en `repo.components`: granularidad de cambio, vida del cambio, ratificada una sola vez en la INTAKE GATE. Esta capability fija su contrato — cómo se llama, dónde vive, qué se emite cuando nada matchea, cómo se apaga, quién la infiere, quién la ratifica y, sobre todo, que **nadie la consume**: es metadata que lee una persona.

## Reglas de negocio

- La línea se llama exactamente **`Components`**, alineada con `repo.components[].name`. Cada valor es el `name` exacto de un componente del set declarado; nunca paths, prefijos ni etiquetas inventadas.
- La línea es **multivaluada**: nombra todos los componentes cuyas `paths` cubren el alcance del pedido, separados por coma.
- **Gate presence-based**: sin `repo.components` declarado en el config del proyecto, el campo no existe en el brief y no se menciona. Con el set declarado, la línea se emite siempre (su ausencia es anomalía, como la ausencia de `Diagram`).
- El valor no se hereda del config global; sólo cuenta el config de proyecto.
- Con el eje activo y ninguna `paths` cubriendo el alcance, la línea se emite con el valor único `unassigned`. Si al menos una `paths` cubre el alcance, la línea nombra sólo esos componentes sin agregar `unassigned`.
- `sdd-intake` infiere el valor mapeando el alcance del pedido contra `repo.components[].paths`. Es una propuesta que el usuario ratifica en la INTAKE GATE — su única ratificación. El orquestador la surfacea por nombre porque está declarada en la tabla de campos-del-brief del fragmento del dominio.
- El valor es **metadata del cambio para un lector humano**. Ninguna fase lo lee para cambiar su comportamiento, ninguna lo consulta en `sdd-verify`/`sdd-tasks`, y ninguna lo escribe o extiende hacia la línea `Components:` de los capability-specs durables. La proyección por-capability conserva su propia ratificación.
- El renderizador del bloque de retorno soporta un gate **por bullet** dentro de una sección etiquetada, con el mismo shape que el gate de sección: el bullet declara su condición y el ejecutor suministra el booleano. Resuelto en falso, el bullet se omite sin fallar. Resuelto en verdadero, el bullet exige su valor. El booleano es obligatorio: su ausencia falla el render nombrando el campo, nunca se interpreta como "apagado" ni omite en silencio. El impresor del esquema declara el gate de cada bullet que lo lleve, como ya lo hace para las secciones.

## Escenarios

### Scenario: un cambio acotado a una sola superficie

- **GIVEN** un repo con el set declarado y un pedido cuyo alcance cae dentro de las `paths` de un solo componente
- **WHEN** el brief se produce
- **THEN** su clasificación lleva la línea `Components` con ese único `name`

### Scenario: un cambio que cruza dos superficies

- **GIVEN** un pedido cuyo alcance toca `paths` de dos componentes declarados
- **WHEN** el brief se produce
- **THEN** la línea nombra a los dos, separados por coma — no elige uno "principal" ni parte el brief

### Scenario: los valores son nombres del set, nunca rutas

- **GIVEN** un alcance que se describió por carpetas del repo
- **WHEN** se escribe la línea
- **THEN** lleva los `name` exactos de los componentes que esas carpetas resuelven, nunca las carpetas mismas ni un prefijo

### Scenario: el gate apagado no deja línea en el brief

- **GIVEN** un repo sin `repo.components` en el config de proyecto
- **WHEN** se produce el brief y se muestra en la INTAKE GATE
- **THEN** no hay línea `Components` en ninguna parte del brief y nada menciona componentes

### Scenario: el gate encendido y la línea ausente es un defecto

- **GIVEN** un repo con el set declarado y un brief de Pass 2 cuya clasificación no lleva la línea
- **WHEN** se lee ese brief
- **THEN** se trata como campo faltante —el mismo tratamiento que un `Diagram` faltante— y no como "el eje no aplica acá"

### Scenario: el global no enciende el gate

- **GIVEN** un config global con `repo.components` y uno de proyecto sin él
- **WHEN** se produce el brief
- **THEN** el gate sigue apagado y no hay línea: el set no se hereda del global

### Scenario: el gate resuelto en falso omite la línea sin fallar

- **GIVEN** un bullet que declara su gate y un ejecutor que suministra el booleano en falso
- **WHEN** se renderiza el bloque
- **THEN** el bullet no aparece y el render termina bien

### Scenario: el valor olvidado con el gate encendido falla el render

- **GIVEN** un bullet con el gate resuelto en verdadero y sin valor suministrado
- **WHEN** se renderiza
- **THEN** el render falla nombrando el campo, igual que cualquier bullet requerido — no se omite la línea

### Scenario: el booleano del gate omitido falla el render

- **GIVEN** un bullet que declara su gate y un ejecutor que no suministra el booleano
- **WHEN** se renderiza
- **THEN** el render falla nombrando el campo del gate, y no asume "apagado"

### Scenario: los bullets existentes no cambian de comportamiento

- **GIVEN** cualquier retorno de fase vigente, cuyos bullets no declaran gate
- **WHEN** se renderiza antes y después del cambio
- **THEN** la salida es idéntica y los campos faltantes siguen fallando igual

### Scenario: nada matchea y la línea se emite igual

- **GIVEN** un repo con el set declarado y un pedido cuyo alcance no cae en ninguna `paths` declarada
- **WHEN** el brief se produce
- **THEN** la línea se emite con el valor `unassigned`, no se omite

### Scenario: match parcial no arrastra `unassigned`

- **GIVEN** un alcance que toca las `paths` de un componente declarado y además un área que ningún componente declara
- **WHEN** se infiere la línea
- **THEN** nombra a ese componente y sólo a ése — `unassigned` no se agrega

### Scenario: `unassigned` y ausencia no son lo mismo

- **GIVEN** dos briefs: uno de un repo sin el set declarado y otro de un repo con el set donde nada matcheó
- **WHEN** se comparan
- **THEN** el primero no lleva línea y el segundo lleva `unassigned` — la distinción es legible sin abrir el config

### Scenario: la propuesta llega al gate junto al resto de la clasificación

- **GIVEN** un brief de Pass 2 con el eje activo
- **WHEN** el orquestrador abre la INTAKE GATE
- **THEN** el valor de `Components` se surfacea para confirmar o ajustar, junto con el lane y el resto de la clasificación

### Scenario: el orquestrador lo surfacea por nombre porque está declarado

- **GIVEN** el fragmento del dominio después del cambio
- **WHEN** el orquestrador busca qué campos del brief se confirman en la INTAKE GATE
- **THEN** encuentra `Components` declarado junto a `diagram` y `ui-test`, con su lector registrado como ninguno

### Scenario: el usuario corrige el valor inferido

- **GIVEN** una inferencia que nombra `cli` y un usuario que sabe que el cambio también toca `api`
- **WHEN** lo ajusta en la INTAKE GATE
- **THEN** el brief se actualiza y se re-muestra, exactamente como un ajuste de lane

### Scenario: ninguna fase posterior lo vuelve a preguntar

- **GIVEN** un valor ratificado en la INTAKE GATE
- **WHEN** corre cualquier fase posterior del lane
- **THEN** ninguna lo re-pregunta ni lo re-infiere: el gate fue su única confirmación

### Scenario: ninguna fase lee el valor

- **GIVEN** el payload después de este cambio
- **WHEN** se busca qué fase lee el valor por-cambio de componentes
- **THEN** ninguna lo lee: sólo `sdd-intake` lo escribe y sólo una persona lo consume

### Scenario: archivar no escribe ninguna línea de capability-spec

- **verification:** `standing → sdd-archive`
- **GIVEN** un cambio con valor por-cambio ratificado y capabilities tocadas sin línea `Components:`
- **WHEN** el cambio se archiva
- **THEN** ninguna línea `Components:` se crea, modifica ni ensancha — siguen requiriendo su ratificación spec por spec

### Scenario: el contrato declara la ausencia de consumidor

- **GIVEN** un lector que se pregunta para qué sirve el campo
- **WHEN** consulta su definición
- **THEN** encuentra declarado que no tiene consumidor de máquina y que eso es deliberado, no un pendiente

### Scenario: ninguna otra fase gana mención del campo

- **GIVEN** el payload antes y después del cambio
- **WHEN** se buscan menciones del valor por-cambio de componentes fuera de `sdd-intake` y de las dos piezas de mecanismo
- **THEN** no hay ninguna: ni prosa nueva en `sdd-spec`, ni en `sdd-archive`, ni en `sdd-tasks`, ni en `sdd-verify`

### Scenario: los contratos previos siguen válidos

- **GIVEN** un `repo.components`, un capability-spec con su línea `Components:` y los retornos de las otras nueve fases, todos válidos antes del cambio
- **WHEN** se cargan, se renderizan y se validan después
- **THEN** todo se comporta igual: sin errores nuevos, sin findings nuevos y sin migración

### Scenario: los escenarios diferidos no se reetiquetan acá

- **GIVEN** los escenarios que llevan `deferred → intake-components-flag` en capability-specs `Accepted`
- **WHEN** este cambio se cierra
- **THEN** sus tokens quedan como están: este cambio no reescribe specs `Accepted`, y reasignarles dueño queda para un cambio posterior

## Referencias

- **Rule** → [`./repo-component-declaration.md`](./repo-component-declaration.md) — la declaración del set del repo.
- **Rule** → [`./capability-spec-components-axis.md`](./capability-spec-components-axis.md) — la proyección **por-capability** (la otra proyección; comparten set y gate, nunca el valor).
- **Rule** → [`./component-inference-ratification.md`](./component-inference-ratification.md) — las tres vías de ratificación; ésta es la vía (c).
- **Rule** → [`./component-projection-independence.md`](./component-projection-independence.md) — por qué una ratificación por-cambio no escribe en la proyección durable.
- **Process** → [`../process/render-durable-artifact.md`](../process/render-durable-artifact.md) — el renderizador hermano, de donde sale el precedente de gate por bullet (`when-present`) del que este cambio se aparta a propósito.
- **Referencia de concepto** → [`../../../payload/domains/development/references/repo-components/README.md`](../../../payload/domains/development/references/repo-components/README.md) — el hogar canónico del concepto repo-level: qué es un componente, cómo se declara el set y el gate presence-based.
