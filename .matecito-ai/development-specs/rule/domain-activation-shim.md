# Capability — Resolución de los dominios activos

- **Status:** Inferred
- **Date:** 2026-07-29
- **Components:** cli

## Propósito

Definir qué dominios se consideran activos en una máquina, para que todo lo que depende de esa respuesta —qué se despliega, qué binarios se instalan, qué integraciones se registran, qué se chequea— la derive de un solo lugar y de la misma forma.

## Reglas de negocio

- **Descubrimiento:** un dominio existe si el payload trae, bajo su carpeta, su declaración de dominio. Una carpeta de dominio sin esa declaración **no** se descubre. El descubrimiento devuelve los identificadores en orden alfabético, para que el resultado sea determinista.
- **Conjunto vacío o no configurado = todos.** Si el config no declara ningún dominio, los activos son **todos los descubiertos en el payload**. Es un shim de compatibilidad: preserva el comportamiento previo a que el ecosistema tuviera varios dominios, sin obligar a nadie a configurar nada.
- **Conjunto explícito = filtro.** Si el config declara dominios, los activos son los declarados **filtrados contra los descubiertos**, conservando el orden en que fueron declarados.
- **Un identificador declarado sin declaración de dominio en el payload se descarta en silencio**: no activa nada ni hace fallar la resolución.
- Lo que agregan los dominios activos —integraciones a registrar, binarios a instalar— se **deduplica** entre dominios, conservando el orden de primera aparición al recorrer los dominios ya resueltos.
- **Degradación ante error de resolución:** cuando la resolución de los dominios activos falla, los consumidores no dejan de funcionar. Cada uno cae a su comportamiento amplio: el despliegue trata el conjunto como vacío (todos), la detección de binarios los detecta todos, y el registro de integraciones cae a su conjunto de seguridad. Una configuración rota nunca deja el entorno a medio instalar en silencio.

## Escenarios

### Scenario: carpeta de dominio sin declaración se ignora

- **GIVEN** un payload con dos dominios que declaran su contrato y una carpeta de dominio que no lo declara
- **WHEN** el sistema descubre los dominios
- **THEN** devuelve solo los dos que lo declaran, en orden alfabético

### Scenario: sin configurar, todos activos

- **GIVEN** un config que no declara ningún dominio
- **WHEN** el sistema resuelve los activos
- **THEN** devuelve todos los dominios descubiertos en el payload

### Scenario: conjunto explícito filtra

- **GIVEN** un config que declara un dominio existente y otro que no existe en el payload
- **WHEN** el sistema resuelve los activos
- **THEN** devuelve solo el existente; el inexistente se descarta sin error

### Scenario: deduplicación entre dominios

- **GIVEN** dos dominios activos que declaran un mismo binario, y uno de ellos además otros dos
- **WHEN** el sistema resuelve los binarios de los activos
- **THEN** devuelve la unión sin repetidos, en orden de primera aparición

### Scenario: error de resolución no bloquea

- **GIVEN** un entorno donde la resolución de dominios activos falla
- **WHEN** el sistema detecta binarios o planifica el despliegue
- **THEN** cae a su comportamiento amplio (detectar todos / desplegar todos) en lugar de omitir componentes en silencio

## Referencias

- **Process** → [`../process/deploy-payload-to-host.md`](../process/deploy-payload-to-host.md) — cómo el conjunto activo filtra el índice, los fragmentos y los componentes desplegados (y por qué los compartidos no se filtran).
- **Rule** → [`hook-registry-domain-filtering.md`](hook-registry-domain-filtering.md) — el mismo conjunto activo aplicado al registro de hooks, con la excepción del dominio compartido.
- **Rule** → [`mcp-permission-auto-approval.md`](mcp-permission-auto-approval.md) — cómo se derivan del conjunto activo los patrones de auto-aprobación.
