# EDR — Topología de despliegue

- **Status:** Accepted
- **Date:** 2026-07-23

## Contexto

Matecito UI acompaña el flujo SDD de todos los proyectos del usuario en una máquina. Correr un daemon y una base por proyecto multiplicaría procesos y fragmentaría el estado. Un único daemon global por máquina, con una sola persistencia compartida y todo scopeado por proyecto, sirve un cockpit global y simplifica la operación. El empaquetado debe dar un artefacto único y autocontenido que incluya la UI ya construida.

## Decisión

Adoptamos una topología de **instancia única global**: un único daemon broker por máquina (no uno por proyecto), servicio de larga vida y stateful, con una única persistencia compartida en el home (`~/.matecito-ai/`) para todos los proyectos, sirviendo un cockpit global. El empaquetado produce un binario único: el broker va embebido en el binario del CLI de la herramienta (vinculado al módulo propio del broker) y la UI (el build de React) va embebida en el binario. Para el multi-proyecto, cada proyecto se registra por su ruta (vía MCP o UI); todas las tablas se scopean por proyecto y la identidad de un change es la combinación de proyecto y nombre de change/rama. Como consecuencia de build, la release construye la UI antes del build de Go, mediante un hook previo que produce el bundle de la UI y lo embebe.

## Reglas verificables

- **[manual]** corre un único daemon broker global por máquina, no una instancia por proyecto.
- **[manual]** el estado vive solo en la persistencia única bajo el home compartido; no hay estado de sesión en memoria que no sea derivable de la persistencia.
- **[manual]** el binario de release embebe la UI ya construida (el build de la UI ocurre antes del build de Go).
- **[manual]** cada proyecto se identifica por su ruta registrada y las filas se scopean por proyecto.

## Alternativas consideradas

- **Un daemon y una base por proyecto:** descartado; multiplica procesos y fragmenta el estado, sin beneficio para un uso local single-user.
- **UI servida por separado del binario:** descartada; embeberla da un artefacto único y autocontenido.

## Consecuencias

- Un solo proceso y una sola base para operar y respaldar.
- Un artefacto de release autocontenido (binario con broker + UI).
- El scope por proyecto habilita el cockpit global sobre estado compartido.
- Trade-off: la release acopla el build de la UI al del binario (la UI debe construirse primero), lo que agrega un paso al pipeline.

## Relacionados

- `relacionado-con` → [configuration.md](configuration.md) — la ubicación de la persistencia en el home compartido deriva de esta topología.
- `relacionado-con` → [../data/data-modeling.md](../data/data-modeling.md) — el scope por proyecto en todas las tablas se fija ahí.
- `relacionado-con` → [../contracts/api-contract.md](../contracts/api-contract.md) — la superficie de registro de proyecto deriva de esta topología.
