# Qué es (y qué no es) un capability-spec

Referencia canónica del **concepto** de capability-spec. Es la fuente de verdad de la *idea*; cualquier skill o agente que trabaje con specs apunta acá en vez de redefinirla. La *estructura/plantilla* concreta se define por separado (`templates/capability.md`); esto define qué cuenta como capability-spec y qué no.

Un capability-spec es la contraparte del EDR. El **EDR** captura *qué se eligió y por qué* (la decisión técnica y su justificación). El **capability-spec** captura *qué hace el sistema* (el comportamiento). Son ortogonales y se referencian, no se solapan.

## Qué ES un capability-spec

Un capability-spec describe el **comportamiento intencionado** de una capacidad del sistema (un flujo, una operación de cara a un actor). Tres rasgos lo definen:

1. **Describe qué hace** — el comportamiento observable ante cada situación: el flujo, sus ramas, sus casos borde, sus reglas. No cómo está implementado: qué hace.
2. **Es verificable** — cada regla y cada flujo se puede chequear con un escenario concreto (Given/When/Then). Si una afirmación no se puede testear, no es spec: es opinión.
3. **Perdura como fuente de verdad del comportamiento** — el código, los tests y los cambios futuros se validan contra él. Es durable y acumulado, no un delta de un cambio puntual.

Responde: *"qué debe hacer el sistema ante cada situación, y cómo se verifica"*.

## Qué NO es un capability-spec

- **No es el "por qué".** La justificación de una elección técnica —el trade-off, la alternativa descartada— es un **EDR**. El spec no argumenta; especifica.
- **No es el "cómo".** El detalle de implementación vive en el código. El spec no nombra **identificadores internos volátiles** (clases, métodos, columnas de base de datos, errores internos, rutas de archivo). Se escribe en el **idioma del dominio** más el **contrato público** de cara al actor (endpoints públicos, códigos de error expuestos). Un identificador interno en el spec lo vuelve un calco del código que se pudre con el primer rename.
- **No es una tarea ni un plan.** Una unidad de trabajo se ejecuta y se termina; el spec perdura y gobierna. El plan de cómo llegar al comportamiento es otra cosa.
- **No es el delta de un cambio.** El artefacto efímero que produce la fase `sdd-spec` (en Engram, `sdd/{change}/spec`) describe lo que un cambio AGREGA/MODIFICA/QUITA; el capability-spec es el **estado acumulado y durable** resultante. El delta se materializa en el capability-spec al archivar el cambio.
- **No es un modelo de datos.** Las entidades aparecen por su **semántica de dominio y sus invariantes de comportamiento**, no por su forma de persistencia (esa es una decisión de datos → EDR).

## Tipos y organización

Los capability-specs se organizan por **tipo** en subcarpetas de `.matecito-ai/development-specs/`, con dos niveles de índice (un índice raíz que enruta por tipo + un `INDEX.md` dentro de cada tipo). El tipo lo lleva la **ruta** (`<type>/<capability>.md`), no una propiedad del header —igual que el dominio de un EDR. Cuatro tipos, cada uno con las secciones que son su "esqueleto":

| Tipo | Qué captura | Secciones esqueleto |
|---|---|---|
| `flow` | Operación de cara a un actor, con pasos y ramas | Flujo principal · Ramas · Casos borde · Errores de cara al actor |
| `rule` | Regla de negocio transversal, sin flujo | Reglas de negocio · Escenarios |
| `lifecycle` | Máquina de estados de una entidad | Entidades y estados · Escenarios |
| `process` | Comportamiento reactivo/de fondo, disparado por evento (no por un actor) | Flujo principal (del proceso) · Casos borde · Reglas de negocio |

Todos usan la misma plantilla (`templates/capability.md`): cada tipo enfatiza sus secciones esqueleto y omite las que no le aplican. Un tipo sin ningún spec-archivo no crea carpeta (se deja constancia en el índice raíz).

### Cómo clasificar el tipo

Toda capability es de **un solo** tipo. Decidilo en este orden:

1. ¿La dispara un **actor** (usuario, consumidor) con una petición, y tiene pasos? → **`flow`**.
2. ¿La dispara un **evento o el sistema** (webhook, job, timer), no un actor, y tiene pasos? → **`process`**.
3. ¿Es una **restricción o invariante sin pasos**, que aplica transversalmente? → **`rule`**.
4. ¿Es el **ciclo de vida de una entidad** (sus estados y transiciones)? → **`lifecycle`**.

`flow` vs `process` se decide por **una** cosa: el disparador (actor vs evento/sistema).

**Tie-breaker del solapamiento (por reuso).** Casi todo `flow` además chequea reglas y cambia estados de entidades; eso NO los vuelve capabilities aparte por sí solo:

> Una regla o un ciclo de vida es una capability **separada** (`rule` / `lifecycle`) SOLO si lo **comparten varios flujos/procesos**. Si es exclusivo de una capability, vive **dentro** de su spec (en la sección "Reglas de negocio" / "Entidades y estados"), no como archivo aparte.

Ejemplos:
- La ventana de 24h aplica a `send-message` y a otros envíos → `rule` propio (`rule/messaging-window-24h`).
- El ciclo de `message` (pending→sent→failed) lo tocan varios flujos → `lifecycle` propio (`lifecycle/message`).
- Una validación de formato que solo usa `connect-meta-account` → vive DENTRO de ese `flow`, no como `rule`.

Esta regla la aplican por igual `development-spec-bootstrap` (al clasificar) y `development-spec-validate` (al chequear que el tipo declarado sea el correcto).

## Estados

Un capability-spec tiene ciclo de vida ligero:

- **Inferred** — borrador no-confiable derivado del código as-built por `development-spec-mine`; candidato pendiente de ratificación humana, todavía no es fuente de verdad. Ver la salvaguarda abajo.
- **Draft** — se está escribiendo o le faltan escenarios; todavía no es fuente de verdad.
- **Accepted** — ratificado por una persona; el comportamiento descrito es el intencionado y el código se valida contra él.
- **Deprecated** — la capacidad se retiró o se reemplazó; se conserva por trazabilidad (el reemplazo se linkea).

Al igual que el EDR, un capability-spec **puede** tener estado `Inferred`, con una salvaguarda estricta: un `Inferred` es un **borrador no-confiable** derivado del código as-built, NO la intención ratificada. No se infiere comportamiento *como si fuera la verdad* —eso calcaría la implementación actual, bugs incluidos—; se infiere como **candidato pendiente de ratificación humana**. Mientras esté `Inferred`, `sdd-verify` lo **ignora** (no es contrato): ese guardarraíl vuelve seguro tenerlo en el store. La intención la fija una persona al promoverlo a `Accepted` (vía `development-spec-bootstrap` modo update).

## Relación con el EDR

Tres capas, no dos:

- **capability-spec** → *qué hace* el sistema (comportamiento).
- **EDR** → *qué se eligió y por qué* (la decisión técnica: tecnología, patrón, política, y su justificación en prosa).
- **código** → el *cómo* literal (la implementación paso a paso).

El spec **referencia** los EDRs que gobiernan cómo se implementa un comportamiento (sección "Referencias"), pero no repite su contenido: el spec dice *la regla*, el EDR dice *cómo se implementa esa regla y por qué así*. Ejemplo: el spec de `outbound-message` dice "no se puede enviar fuera de la ventana de 24h; al vencer, el sistema responde con error X"; el EDR dice "modelamos esa condición como una jerarquía de errores de dominio mapeada a RFC 7807, porque…".

## Eje `components`

Un **componente** es una **superficie que el consumidor del producto reconoce** — `api`, `cli`, `ui` — no todo paquete o carpeta interna. Es un eje **opcional** sobre el store de capability-specs: le agrega a cada spec qué superficie(s) del repo participan del comportamiento que describe.

### Declaración

El set de componentes se declara **una sola vez**, en el config **del proyecto** (nunca del global), en un bloque `repo` al tope, hermano de `domains`/`domainConfig`:

```json
"repo": {
  "components": [
    { "name": "cli", "paths": ["cmd", "internal"] },
    { "name": "api", "paths": ["apps/api"] },
    { "name": "ui",  "paths": ["apps/ui"] }
  ]
}
```

Array de objetos: `name` (la superficie) + `paths` (una o más carpetas del repo, relativas a la raíz, donde vive). Es el **único** lugar del ecosistema donde se escribe la carpeta de un componente — el capability-spec solo referencia el `name` en su línea de header; no repite las `paths`.

### La línea del header

Con el eje declarado, cada capability-spec suma junto a `Status` y `Date` una línea multivaluada:

```
- **Components:** api, ui
```

Cada valor listado debe estar dentro del set declarado en `repo.components`. Es lo único que el eje agrega al store: los `INDEX.md` (raíz y de tipo) no cambian.

### Gate presence-based

Sin `repo.components` declarado, el eje **no existe**: ningún spec lleva la línea `Components:`, y ningún chequeo de validación se dispara por su ausencia. El eje es enteramente opt-in — un repo se comporta exactamente igual que antes de que este eje existiera hasta que alguien declara el set.

### Por qué el store nunca se parte por app

El store de capability-specs vive siempre en la raíz del repo (`.matecito-ai/development-specs/`), nunca dentro de una app o componente — a diferencia del EDR, que sí vive por sub-app (`apps/api/.matecito-ai/edr/`, `apps/ui/.matecito-ai/edr/`).

La razón no es arbitraria: son ejes distintos, no una inconsistencia entre el spec y el EDR.

- El **EDR** captura **cómo está construida** una pieza — es una decisión que, por naturaleza, pertenece a una sola pieza de código. Univaluado: vive donde vive el código que gobierna.
- El **capability-spec** captura **qué hace el sistema** — y un comportamiento casi nunca respeta el límite de una sola app: un flujo que arranca en la `ui`, pasa por la `api` y termina escribiendo algo que el `cli` también puede disparar es **una sola capacidad**, no tres. Multivaluado: el eje `components` existe precisamente para expresar esa lista de superficies sin partir el spec.

Partir el store por app fuerza a elegir un dueño arbitrario para cada spec que cruza superficies, o a duplicar el spec en cada app que toca — y un spec duplicado diverge con el primer cambio que se aplique de un solo lado, dejando de ser un contrato. Por eso, ante un pedido de crear el store dentro de una app, la autoría no lo hace en silencio ni lo discute: explica este porqué y ofrece declarar el eje `components` en su lugar.
