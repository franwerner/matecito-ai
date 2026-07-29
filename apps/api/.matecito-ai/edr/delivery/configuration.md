# EDR — Configuración

- **Status:** Accepted
- **Date:** 2026-07-23

## Contexto

El broker es un daemon local con pocos parámetros (puerto, ubicación de la persistencia). No justifica una librería de config ni un archivo de configuración. Debe arrancar solo con defaults sensatos y permitir override puntual al ejecutarlo. Además debe ser configurable de forma independiente para poder embeberse después en el CLI, que le pasaría los mismos parámetros. En el modelo de instancia única global, la persistencia default vive en el home compartido, no en el repo.

## Decisión

Los defaults viven en el código y se overridean con flags al ejecutar el comando (puerto, ubicación de la persistencia, etc.) y variables de entorno, usando solo la librería estándar, sin librería de config ni archivo. La configuración es independiente del broker: si más adelante se embebe en el CLI, el CLI le pasa los flags. Al arranque se validan los valores (puerto válido, ubicación de la persistencia escribible) y se aborta con un mensaje claro si alguno es inválido. Por el modelo global, la ubicación default de la persistencia es el home compartido (`~/.matecito-ai/`), no un directorio dentro del repo.

Lo que el daemon crea en esa ubicación queda **restringido a su dueño**: el directorio de persistencia sin acceso para terceros, y el identificador de máquina legible y escribible solo por el dueño. El sistema es single-user por diseño, así que no hay caso de uso para que otro usuario de la misma máquina lea el trabajo indexado; el default permisivo del lenguaje sería una concesión sin contraparte.

## Reglas verificables

- **[manual]** el broker arranca con defaults, sin config externa.
- **[manual]** los flags al ejecutar el comando overridean los defaults (puerto, ubicación de la persistencia).
- **[manual]** la config se valida al startup y aborta con mensaje claro si un valor es inválido.
- **[manual]** la ubicación default de la persistencia es el home compartido (`~/.matecito-ai/`).
- **[auto]** el directorio de persistencia se crea sin acceso para terceros y el identificador de máquina sin acceso para grupo ni terceros; aflojar cualquiera de los dos rompe el gate de lint.

## Alternativas consideradas

- **Librería de config + archivo de configuración:** descartada; para un puñado de parámetros, flags y env vars con la stdlib alcanzan.

## Consecuencias

- Arranque sin fricción y override puntual cuando hace falta.
- La configuración independiente permite embeber el broker en el CLI sin acoplar la carga de config.
- Trade-off: sin archivo de config, una configuración con muchos overrides se vuelve una línea de comando larga.

## Relacionados

- `relacionado-con` → [deployment-topology.md](deployment-topology.md) — la ubicación en el home compartido deriva del modelo de instancia única global.
