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

### Requisito: La segunda corrida no reporta nada pendiente

Correr la instalación de nuevo sobre una integración ya registrada DEBE reportar nada pendiente para ella y NO DEBE descargar, reinstalar ni re-registrar.

#### Scenario: segunda corrida idempotente

- **GIVEN** una máquina donde la integración ya está instalada y registrada
- **WHEN** la persona corre la instalación del ecosistema de nuevo
- **THEN** el plan no lista esta integración y nada se descarga ni re-registra

### Requisito: Una falla de este paso no aborta la corrida

Una falla en cualquier punto de este paso DEBE ser reportada nombrando el componente, y la corrida de instalación que la rodea DEBE continuar con sus componentes restantes — el contrato de continue-on-error existente, que este paso no amplía ni escapa.

#### Scenario: la descarga falla a mitad de la corrida

- **GIVEN** una corrida de instalación con varios componentes activos, donde la descarga del artefacto falla
- **WHEN** se ejecuta la corrida
- **THEN** reporta la falla de este componente, continúa con los restantes y termina en error nombrando el que falló

## Referencias

- **Flow** → [`install-ecosystem.md`](../flow/install-ecosystem.md) — el plan / confirmación / contrato de continue-on-error en el que este paso se ejecuta.
- **Rule** → [`rule/mcp-permission-auto-approval.md`](../rule/mcp-permission-auto-approval.md) — de dónde viene el patrón de permisos convencional y por qué no se declara override acá.
