---
name: setup-record-search
version: 1.1.0
description: Configura la búsqueda por significado sobre los records durables de un proyecto — crea una colección por store, declara sus exclusiones y las credenciales de embeddings, e indexa. Usá esta skill cuando la búsqueda semántica no esté disponible y quieras habilitarla, cuando el proyecto sume o quite un componente y las colecciones queden desalineadas, cuando `find-records` avise que trabajó sin la vía semántica, o cuando el usuario pida "configurar la búsqueda de records" o "que encuentre los EDR por significado". Es setup: corre una vez por proyecto, no en cada búsqueda.
license: MIT
metadata: {"hermes":{"tags":["records","search","setup","qmd","mcp"],"category":"cross-domain","related_skills":["find-records"]},"author":"matecito-ai","version":"1.1.0"}
---

# Setup Record Search

Deja lista la vía de búsqueda por significado que `find-records` usa cuando está
disponible. Sin esto, esa skill trabaja con el índice y la búsqueda literal, y
pierde el único camino que encuentra un record cuando no compartís vocabulario
con él.

**Esto es setup.** Corre una vez por proyecto, y otra vez cuando cambian los
stores. No es parte de buscar.

## Paso 1 — Ver qué stores hay

```sh
node ~/.claude/skills/find-records/scripts/map-records.js --json
```

Devuelve, por cada store: su `kind` (`decisions` o `behavior`), su `path`, el
`component` que lo gobierna (`null` cuando gobierna todo el repo), y la
`collection` con la que hay que darlo de alta.

## Paso 2 — Una colección por store

**Una colección por store, sin excepciones.** No agrupes por tipo ni por
componente: el filtro por store es lo que hace que buscar comportamiento no
devuelva decisiones, y agrupar lo pierde.

**Usá el nombre que devolvió el paso 1**, tal cual: `<proyecto>-<kind>-<component>`,
con `root` cuando el store no pertenece a ningún componente. Esa forma no depende
de cómo se llame la carpeta —`edr` en development, `ddr` en design— así que sirve
igual en cualquier dominio.

El prefijo del proyecto no es decoración. Un solo servidor sirve a todos tus
repos desde un único archivo de configuración, y ahí el nombre es la clave: si un
segundo proyecto declara `decisions-root`, no se suma — reemplaza la ruta del
primero, y esa colección queda apuntando al repo equivocado sin que nada falle ni
avise.

Para cada store del paso 1:

```sh
qmd collection add <path absoluto del store> --name <collection del paso 1> --mask '**/*.md'
```

**Excluí los `INDEX.md`.** Son tablas de navegación: matchean con todo y ensucian
el primer resultado. El comando no expone exclusiones, así que se declaran en
`~/.config/qmd/index.yml`, agregando a cada colección:

```yaml
    ignore:
      - "**/INDEX.md"
```

Un repo llamado `tienda`, con un store por componente más los de la raíz, queda
con `tienda-decisions-root`, `tienda-behavior-root`, `tienda-decisions-api` y
`tienda-decisions-ui`.

## Paso 3 — Credenciales

El servidor necesita un proveedor de embeddings. Se configura por entorno:

```sh
QMD_OPENAI_BASE_URL=<endpoint compatible con la API de OpenAI>
QMD_OPENAI_API_KEY=<la key>
QMD_EMBED_MODEL=<el modelo de embeddings>
```

**Nunca escribas la key en `.matecito-ai/config.json`**: ese archivo está
versionado y la estarías publicando. Va en el entorno, o en la configuración
global del usuario (`~/.matecito-ai/config.json`), que no viaja con el repo.

Sin credenciales no sigas: indexar sin proveedor falla a mitad y deja el índice
a medio construir. Decilo y pará.

## Paso 4 — Indexar

```sh
qmd update
qmd embed
```

El primero recorre los archivos, el segundo genera los vectores. Sobre un store
de unos cien records esto tarda menos de un minuto y cuesta fracciones de
centavo; sobre uno grande, proporcionalmente más.

Con esto el setup está completo: las colecciones, sus exclusiones, las
credenciales y el índice.

## Lo que esta skill no hace

**Levantar el servidor no es parte de este setup.** El MCP se declara junto a
los demás, con su comando, y el cliente lo arranca al abrir la sesión —
declaración que escribe el despliegue de matecito-ai, no esta skill. Por eso acá
no hay ningún proceso que arrancar, ningún puerto que leer y nada que quede
corriendo de fondo entre sesiones.

Si `find-records` sigue reportando que la vía semántica no está disponible
después de este setup, el problema está en esa declaración o en el permiso
(`mcp__qmd__*`: sin él el agente declara la herramienta pero no puede llamarla,
y el síntoma es idéntico) — no en la configuración de colecciones.

## Cuándo volver a correr esto

Cuando el proyecto **suma o quita un componente con store propio**: las
colecciones quedan desalineadas y los records de ese componente dejan de
aparecer, en silencio. El paso 1 lo detecta — si devuelve un store sin colección
que le corresponda, hay que crearla.
