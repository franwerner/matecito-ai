# Capability — Reconciliación de handlers por marcador de identidad

- **Status:** Inferred
- **Date:** 2026-07-29
- **Components:** cli

## Propósito

Mantener al día en el archivo de settings del host los handlers que el ecosistema declara —reemplazando los propios que quedaron viejos y borrando los que ya no se declaran— sin tocar jamás los que puso la persona.

## Reglas de negocio

- **La propiedad se declara con un marcador de identidad.** El ecosistema es dueño de todo handler que lleve un marcador de identidad propio no vacío; los handlers **sin** ese marcador son de la persona.
- **Nunca se toca un handler ajeno.** Un handler sin marcador se conserva intacto, esté donde esté, aunque comparta grupo con uno propio.
- **Un handler propio que ya coincide con lo declarado se conserva.** Coincidir significa coincidir en todo lo observable: mismo marcador, mismo evento, mismo matcher, mismo comando, misma condición y mismo tiempo límite.
- **Un handler propio que no coincide se reemplaza**: se quita el viejo y se agrega el declarado. Es lo que hace que cambiar el comando de un hook no deje dos versiones conviviendo.
- **Un handler propio que ya no está declarado se elimina.** Retirar un hook del ecosistema lo retira también del host.
- **Los grupos que quedan vacíos se descartan**, y si un evento se queda sin ningún grupo, el evento desaparece: no quedan estructuras vacías como resto.
- **Un handler declarado que no existe se agrega** al grupo de su evento y matcher; si ese grupo no existe, se crea.
- **Idempotencia:** si el conjunto declarado ya está exactamente presente, la reconciliación **no reporta cambio** y no toca el documento. Correr la instalación dos veces no ensucia el archivo.
- El marcador de identidad se deriva del dominio y del nombre del handler cuando el hook no declara uno propio, así que cada hook tiene identidad estable sin que haya que asignarla a mano.

## Escenarios

### Scenario: idempotente cuando ya está todo

- **GIVEN** un archivo de settings que ya contiene exactamente el handler declarado, con el mismo marcador, evento, matcher y comando
- **WHEN** el sistema reconcilia
- **THEN** no reporta cambio y el documento queda igual

### Scenario: reemplazo de un handler propio obsoleto

- **GIVEN** un handler propio con el mismo marcador pero con un comando anterior
- **WHEN** el sistema reconcilia contra el comando nuevo
- **THEN** reporta cambio y queda un único handler, con el comando nuevo y el mismo marcador

### Scenario: el handler de la persona se preserva

- **GIVEN** un grupo que contiene un handler sin marcador y el handler propio ya correcto
- **WHEN** el sistema reconcilia
- **THEN** no reporta cambio y ambos handlers siguen presentes

### Scenario: se elimina un handler que ya no se declara

- **GIVEN** un archivo de settings con un handler propio y un conjunto declarado vacío
- **WHEN** el sistema reconcilia
- **THEN** reporta cambio y no queda ningún handler propio

### Scenario: el grupo vacío se descarta

- **GIVEN** un evento con un único grupo cuyo único handler era propio, y un conjunto declarado vacío
- **WHEN** el sistema reconcilia
- **THEN** el evento queda sin grupos y su entrada desaparece del documento

### Scenario: se agrega un handler declarado que falta

- **GIVEN** un archivo de settings sin ningún handler para el evento y matcher declarados
- **WHEN** el sistema reconcilia
- **THEN** reporta cambio y crea el grupo con el handler declarado y su marcador

## Referencias

- **Rule** → [`hook-registry-domain-filtering.md`](hook-registry-domain-filtering.md) — qué handlers entran en el conjunto declarado según los dominios activos.
- **Process** → [`../process/validate-git-commit-message.md`](../process/validate-git-commit-message.md) — el handler que hoy se reconcilia con este mecanismo.
- **Lifecycle** → [`../lifecycle/component-check-status.md`](../lifecycle/component-check-status.md) — cómo se reporta que un handler declarado falta en el archivo de settings.
