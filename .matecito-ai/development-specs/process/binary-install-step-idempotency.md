# Idempotencia de paso de instalación de binario

- **Status:** Accepted
- **Date:** 2026-09-02
- **Components:** cli

## Propósito

Garantizar que el sistema decide la presencia de un binario instalado observándolo donde su paso lo deja, en ambas superficies (instalación y reporte de estado del entorno), para que una máquina no reporte contradicciones entre ambas ni repita descargas de binarios ya presentes y sanos.

## Actores

- La persona que corre `matecito-ai install` — dispara la planificación e instalación del ecosistema
- La persona que corre `matecito-ai verify` — solicita el estado del entorno
- El paso de instalación de cada binario — deja el ejecutable en su lugar de instalación
- El chequeo de reporte de estado — observa la presencia del binario en su lugar de instalación

## Requisitos

### Requisito: La presencia de un binario se observa donde su paso lo instala

El sistema DEBE decidir si el binario de un paso está instalado observándolo **en el lugar donde ese paso lo deja** —que exista ahí y que corra—, y NO DEBE decidirlo por si la sesión en curso alcanza un programa con ese nombre. La regla rige **cada vez que el sistema decide esa presencia**: tanto al planificar la instalación del ecosistema como al reportarle a la persona el estado del entorno. Que ese lugar todavía no esté incorporado al entorno de la sesión en curso NO DEBE, por sí solo, hacer que el binario se dé por faltante: lo que una corrida deja persistido para sesiones futuras no rige sobre la sesión que la ejecutó. Simétricamente, la mera alcanzabilidad tampoco alcanza para darlo por instalado: un binario ausente de su lugar, o presente ahí pero incapaz de correr, DEBE reportarse como faltante en cualquiera de las dos superficies.

#### Scenario: el lugar donde el paso instala todavía no está en el entorno, y la persona instala

- **GIVEN** una máquina recién estrenada donde la corrida anterior instaló el binario, y este corre desde el lugar donde su paso lo deja
- **AND** la sesión en curso todavía no incorporó ese lugar a su entorno
- **WHEN** la persona corre la instalación del ecosistema inmediatamente después
- **THEN** el plan no lista ese binario, no se pide confirmación por él, y nada se descarga ni se reinstala

#### Scenario: el lugar donde el paso instala todavía no está en el entorno, y la persona pide el estado

- **GIVEN** una máquina recién estrenada donde la corrida anterior instaló el binario, y este corre desde el lugar donde su paso lo deja
- **AND** la sesión en curso todavía no incorporó ese lugar a su entorno
- **WHEN** la persona pide el estado del entorno
- **THEN** el reporte da ese binario por correcto, con su versión detectada, y no lo cuenta como chequeo faltante

#### Scenario: el binario no está en su lugar

- **GIVEN** una máquina donde el binario no está en el lugar donde su paso lo deja
- **WHEN** la persona corre la instalación del ecosistema, o pide el estado del entorno
- **THEN** la instalación lista ese binario como pendiente
- **AND** el reporte de estado lo marca como faltante, con su pista de remediación

#### Scenario: el binario está en su lugar pero no corre

- **GIVEN** una máquina donde el binario está en el lugar donde su paso lo deja pero falla al ejecutarse
- **WHEN** la persona corre la instalación del ecosistema, o pide el estado del entorno
- **THEN** la instalación lista ese binario como pendiente
- **AND** el reporte de estado lo marca como faltante

#### Scenario: alcanzable desde el entorno pero ausente de su lugar

- **GIVEN** una máquina donde la sesión alcanza un programa con el nombre del binario, y el lugar donde su paso lo deja está vacío
- **WHEN** la persona corre la instalación del ecosistema, o pide el estado del entorno
- **THEN** la instalación lista ese binario como pendiente
- **AND** el reporte de estado lo marca como faltante

### Requisito: Instalar y reportar el estado no se contradicen sobre la misma máquina

Sobre la misma máquina, y sin que el entorno cambie entre una consulta y la otra, la respuesta de la instalación y la del reporte de estado sobre si un binario está instalado DEBEN coincidir. NO DEBE existir un estado en que la instalación reporte que no hay nada pendiente para ese binario mientras el reporte de estado lo da por faltante, ni el inverso. Qué piezas son críticas y cuáles opcionales no cambia por esto: la coincidencia es sobre si el binario está instalado, no sobre el peso que cada superficie le da.

#### Scenario: instalar y pedir el estado enseguida, en la misma sesión

- **GIVEN** una máquina limpia donde la persona acaba de correr la instalación del ecosistema con éxito
- **AND** la sesión no incorporó ningún lugar nuevo a su entorno desde entonces
- **WHEN** la persona pide el estado del entorno a continuación
- **THEN** el reporte da por correctos los mismos binarios que la instalación acaba de dejar instalados
- **AND** ninguno de ellos aparece como chequeo faltante

#### Scenario: un binario ausente se ve igual desde las dos superficies

- **GIVEN** una máquina donde uno de los binarios no está en el lugar donde su paso lo deja
- **WHEN** la persona corre la instalación del ecosistema y luego pide el estado del entorno
- **THEN** la instalación lo lista como pendiente y el reporte lo marca como faltante — las dos superficies dicen lo mismo sobre él

### Requisito: Un binario presente y sano no se vuelve a descargar ni a reinstalar

Cuando el binario de un paso está presente y sano, la corrida de instalación NO DEBE descargar su artefacto de distribución ni reinstalarlo. El verbo del plan lo determina el estado observado del binario: ausente → *instalar*; presente pero con una versión distinta de la vigente → *actualizar*; presente y al día → nada. Un binario observado como presente NO DEBE quedar clasificado para *instalar*. Este requisito rige sobre la superficie que instala; el reporte de estado no descarga ni instala nada.

#### Scenario: la segunda corrida no descarga nada por un binario ya instalado

- **GIVEN** una máquina donde el binario está instalado, al día y corre desde el lugar donde su paso lo deja
- **WHEN** la persona corre la instalación del ecosistema de nuevo
- **THEN** el plan no lista ese binario
- **AND** no se descarga ningún artefacto de distribución para él

#### Scenario: un binario desactualizado se actualiza, no se reinstala

- **GIVEN** una máquina donde el binario está presente en su lugar y corre, con una versión distinta de la vigente
- **WHEN** la persona corre la instalación del ecosistema
- **THEN** el plan lo lista con el verbo *actualizar*, no con el verbo *instalar*

#### Scenario: primera corrida en una máquina limpia

- **GIVEN** una máquina donde ninguno de los binarios está instalado
- **WHEN** la persona corre la instalación del ecosistema
- **THEN** el plan lista cada uno de esos binarios con el verbo *instalar*
- **AND** al terminar la corrida cada uno corre desde el lugar donde su paso lo deja

### Requisito: El lugar de instalación es propio de cada paso

El lugar donde un paso deja su ejecutable es propio de ese paso: el sistema DEBE resolverlo paso por paso y NO DEBE asumir un único lugar común a todos los binarios. Esto rige igual en las dos superficies. El veredicto de un paso NO DEBE quedar condicionado por el de otro. Cuando el lugar de un paso no se puede resolver, ese paso DEBE reportarse como faltante nombrando su binario —no puede establecer dónde viviría lo suyo, y darlo por instalado sería afirmar lo que no observó—, y los demás binarios DEBEN conservar su veredicto.

#### Scenario: dos pasos con lugares de instalación distintos

- **GIVEN** una máquina donde dos binarios se instalan en lugares distintos, uno presente y sano en el suyo y el otro ausente del suyo
- **WHEN** la persona corre la instalación del ecosistema
- **THEN** el plan lista únicamente el ausente
- **AND** el presente no se reporta pendiente por el veredicto del otro

#### Scenario: el lugar de un paso no se puede resolver, al instalar

- **GIVEN** una máquina donde el lugar de instalación de uno de los binarios no se puede resolver, y los demás binarios están presentes y sanos en los suyos
- **WHEN** la persona corre la instalación del ecosistema
- **THEN** el plan lista como pendiente el binario cuyo lugar no se pudo resolver
- **AND** los demás no se listan

#### Scenario: el lugar de un paso no se puede resolver, al reportar el estado

- **GIVEN** una máquina donde el lugar de instalación de uno de los binarios no se puede resolver, y los demás binarios están presentes y sanos en los suyos
- **WHEN** la persona pide el estado del entorno
- **THEN** el reporte marca como faltante el binario cuyo lugar no se pudo resolver, nombrándolo y con su pista de remediación
- **AND** los demás se reportan correctos

## Referencias

- **Flow** → [`../flow/install-ecosystem.md`](../flow/install-ecosystem.md) — el plan, la confirmación única y el contrato de continue-on-error dentro de los cuales la superficie de instalación se ejecuta.
- **Process** → [`install-mcp-from-release-tarball.md`](install-mcp-from-release-tarball.md) — el paso que además registra la integración con el host: cumple este piso y suma requisitos propios que este spec no relaja.
- **Lifecycle** → [`../lifecycle/component-check-status.md`](../lifecycle/component-check-status.md) — el vocabulario (`OK` / `Missing` / `Outdated`), la criticidad y la pista de remediación con que la superficie de reporte expresa el veredicto que este spec fija. Ese contrato no cambia: acá cambia qué se observa, no cómo se lo nombra.
