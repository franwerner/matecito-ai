# Capability — Instalar un servidor de integración desde un release tarball

- **Status:** Accepted
- **Date:** 2026-09-02
- **Components:** cli

## Propósito

Registrar un servidor de integración cuya única distribución es un tarball de un release en GitHub — sin entrada en un registro de paquetes — durante la corrida de instalación del ecosistema, idempotentemente y sin llevar el resto de la corrida hacia abajo.

## Actores

- **Los flujos de instalación y actualización del ecosistema** — este proceso no tiene punto de entrada propio: es un motor que ellos invocan como uno de sus pasos.

## Requisitos

### Requisito: Resolver el artefacto desde el release más reciente publicado

El sistema DEBE resolver el artefacto a instalar desde el **release más reciente publicado** de la forja de integración, en tiempo de instalación, y NO DEBE fijar una versión en ningún lado del repositorio. Cuando el release no publica un manifiesto de checksums, el sistema DEBE proceder sin verificación de integridad en lugar de fallar — el límite de confianza es el del transporte, el mismo que cada paso basado en registro ya está usando.

#### Scenario: el release más reciente es el que se instala

- **GIVEN** una forja de integración que ha publicado releases
- **WHEN** el paso de instalación se ejecuta
- **THEN** resuelve el artefacto del release más nuevo publicado e instala ese

#### Scenario: un release más nuevo se recoge sin edit del repositorio

- **GIVEN** una máquina que instaló la integración desde el release N
- **WHEN** se publica el release N+1 y la persona corre la instalación en una máquina limpia
- **THEN** se instala el release N+1, sin edición alguna del repositorio

#### Scenario: la ausencia de manifiesto de checksums no detiene la instalación

- **GIVEN** un release que no publica manifiesto de checksums
- **WHEN** el paso de instalación se ejecuta
- **THEN** instala el artefacto de todas formas, sin reportar falla de verificación

### Requisito: Instalar globalmente, luego registrar con el host

El sistema DEBE instalar el artefacto resuelto de modo que su ejecutable sea alcanzable en PATH, y luego registrar la integración con el host bajo el nombre que el manifest declara, para que el host la lance. El paso solo tiene éxito cuando ambas operaciones se han completado.

#### Scenario: entorno limpio

- **GIVEN** una máquina con la integración ni instalada ni registrada
- **WHEN** la persona corre la instalación del ecosistema
- **THEN** el ejecutable de la integración es llamable en PATH
- **AND** la integración está registrada con el host bajo su nombre declarado

#### Scenario: la auto-aprobación se deriva por convención

- **GIVEN** la integración declarada por un dominio activo, sin override explícito de permisos para ella
- **WHEN** el sistema deriva los patrones de permisos de tools a auto-aprobar
- **THEN** el patrón convencional para ese nombre está entre ellos

#### Scenario: falta una precondición

- **GIVEN** una máquina donde el CLI del host o el gestor de paquetes no está en PATH
- **WHEN** el paso de instalación se ejecuta
- **THEN** falla nombrando la precondición faltante, en lugar de reportar éxito o saltarse silenciosamente

### Requisito: El estado de conexión de una integración registrada se lee del host, no se asume

Cuando el sistema consulta al host por una integración registrada, el estado de conexión que obtiene DEBE ser el que el host reporta. El sistema NO DEBE devolver "no conectada" por defecto cuando no llegó a consultarlo — un valor no consultado y un valor consultado que dio negativo son cosas distintas, y quien lea el resultado DEBE poder distinguirlas. La determinación de **presencia** del registro NO DEBE cambiar por esto: toda integración registrada se sigue encontrando exactamente igual que antes.

#### Scenario: una integración registrada y viva se reporta conectada

- **GIVEN** una máquina donde la integración está registrada con el host y el host la lanza sin error
- **WHEN** el sistema consulta al host por esa integración
- **THEN** la encuentra
- **AND** el estado de conexión que obtiene dice que está conectada

#### Scenario: una integración registrada que el host no logra lanzar se reporta no conectada

- **GIVEN** una máquina donde la integración está registrada pero el host falla al lanzarla
- **WHEN** el sistema consulta al host por esa integración
- **THEN** la encuentra
- **AND** el estado de conexión que obtiene dice que no está conectada

#### Scenario: una integración no registrada no se confunde con una registrada y caída

- **GIVEN** una máquina donde la integración no está registrada con el host
- **WHEN** el sistema consulta al host por esa integración
- **THEN** reporta que no la encontró, en lugar de un hallazgo con estado de conexión negativo

#### Scenario: la presencia de las demás integraciones registradas no cambia

- **GIVEN** una máquina con varias integraciones registradas con el host y otras que no lo están
- **WHEN** el sistema consulta al host por cada una
- **THEN** encuentra todas las registradas y ninguna de las no registradas, igual que antes de este cambio

### Requisito: Un ejecutable que este paso no instaló nunca queda resuelto en silencio

Cuando la corrida encuentra, bajo el nombre de la integración, un ejecutable que este paso no instaló ni mantiene, la corrida NO DEBE terminar reportando que no había nada pendiente para ella y dejando la integración inutilizable. DEBE terminar en uno de dos estados: la integración instalada y funcionando, o el problema reportado nombrando la integración. Cuál de los dos —y si además la corrida falla por eso— es una elección de ingeniería fuera de este contrato; lo que el contrato prohíbe es el tercer resultado, el silencio.

#### Scenario: el ejecutable ajeno no se resuelve en silencio

- **GIVEN** una máquina con un ejecutable bajo el nombre de la integración que este paso no instaló
- **WHEN** la persona corre la instalación del ecosistema
- **THEN** al terminar, o la integración quedó instalada y su ejecutable corre, o la corrida reportó el problema nombrando la integración
- **AND** en ningún caso la corrida termina reportando que no había nada pendiente para esa integración

### Requisito: La segunda corrida no reporta nada pendiente

Correr la instalación de nuevo sobre una integración que este paso ya instaló y registró DEBE reportar nada pendiente para ella y NO DEBE descargar, reinstalar ni re-registrar. "Ya instalada" NO se satisface con la mera presencia de un ejecutable con ese nombre en PATH: el sistema DEBE distinguir el ejecutable que este paso instala y mantiene de otro que simplemente esté ahí, y uno que no es el suyo DEBE reportarse como pendiente. Qué señal concreta implementa esa distinción —dónde resuelve el ejecutable, de qué release es, si corre— es una elección de ingeniería y no parte de este contrato: lo que el contrato fija es que la distinción existe y es observable en lo que la corrida reporta.

#### Scenario: segunda corrida idempotente

- **GIVEN** una máquina donde la integración fue instalada por este paso, está registrada, y su ejecutable corre
- **WHEN** la persona corre la instalación del ecosistema de nuevo
- **THEN** el plan no lista esta integración y nada se descarga ni re-registra

#### Scenario: lo que este paso acaba de instalar no vuelve a reportarse pendiente

- **GIVEN** una máquina donde la corrida anterior instaló y registró la integración por este paso
- **WHEN** la persona corre la instalación del ecosistema inmediatamente después
- **THEN** el plan no lista esta integración y nada se descarga ni re-registra

#### Scenario: un ejecutable que este paso no instaló no cuenta como instalada

- **GIVEN** una máquina donde la integración está registrada con el host y hay en PATH un ejecutable con su nombre que este paso no instaló
- **WHEN** la persona corre la instalación del ecosistema
- **THEN** el plan lista esta integración como pendiente

#### Scenario: un ejecutable presente pero que no corre se reporta pendiente

- **GIVEN** una máquina donde la integración está registrada y hay un ejecutable con su nombre en PATH que falla al ejecutarse
- **WHEN** la persona corre la instalación del ecosistema
- **THEN** el plan lista esta integración como pendiente

### Requisito: Una falla de este paso no aborta la corrida

Una falla en cualquier punto de este paso DEBE ser reportada nombrando el componente, y la corrida de instalación que la rodea DEBE continuar con sus componentes restantes — el contrato de continue-on-error existente, que este paso no amplía ni escapa.

#### Scenario: la descarga falla a mitad de la corrida

- **GIVEN** una corrida de instalación con varios componentes activos, donde la descarga del artefacto falla
- **WHEN** se ejecuta la corrida
- **THEN** reporta la falla de este componente, continúa con los restantes y termina en error nombrando el que falló

## Referencias

- **Flow** → [`install-ecosystem.md`](../flow/install-ecosystem.md) — el plan / confirmación / contrato de continue-on-error en el que este paso se ejecuta.
- **Rule** → [`rule/mcp-permission-auto-approval.md`](../rule/mcp-permission-auto-approval.md) — de dónde viene el patrón de permisos convencional y por qué no se declara override acá.
