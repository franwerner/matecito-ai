# Capability — Ciclo de desarrollo en un solo comando

- **Status:** Accepted
- **Date:** 2026-08-03
- **Components:** cli

## Propósito

Ofrecer a la persona que desarrolla un solo comando para compilar el binario del árbol y dejarlo operativo en el entorno local sin que la instalación del ecosistema lo pise con la release publicada.

## Actores

- **La persona que desarrolla** que corre `make dev-install` desde el root del repo.

## Flujo principal

1. La persona corre `make dev-install`.
2. El sistema compila el binario del árbol con la configuración por defecto (`go build`).
3. Si la compilación falla, el comando termina en error sin copiar ni instalar nada.
4. Copia el binario compilado a `~/.local/bin/matecito-ai`.
5. Ejecuta ese binario exacto con `install -y`, sin pedir confirmación.
6. Si algún paso falla, el comando termina en error.

## Casos borde

- **Compilación fallida** → el comando termina en error; no se copia ni se instala nada.
- **Copia fallida** → el comando termina en error; no se ejecuta `install`.
- **Instalación fallida** → el comando termina en error; el binario dev queda en `~/.local/bin/` (la copia precedente de `install` si ocurrió no fue a revertirse).
- **El payload existente se va a reemplazar** → sin destino separado ni reversión; el deploy pisa lo que haya en el host destino.

## Reglas de negocio

- Cada paso de la cadena (`go build`, `cp`, `install -y`) es una invocación de shell independiente dentro del Makefile, por lo que el fracaso de cualquiera aborta la cadena y evita los pasos siguientes (make's default abort-on-error).
- No se busca aislamiento del binario dev: el payload del árbol reemplaza el existente sin destinos alternativos.
- El skip de `matecito-ai` durante `install` (por ser dev build) no detiene la instalación de los demás componentes; ver [`../flow/install-ecosystem.md`](../flow/install-ecosystem.md).

## Escenarios

### Scenario: ciclo completo en un comando

- **GIVEN** un árbol de trabajo con cambios sin publicar
- **WHEN** la persona corre `make dev-install`
- **THEN** el binario compilado queda en `~/.local/bin/`, el ecosistema se instala desde él sin confirmación, y `matecito-ai` no se reemplaza por la release

### Scenario: un paso que falla corta el ciclo

- **GIVEN** un árbol que no compila
- **WHEN** la persona corre `make dev-install`
- **THEN** el comando termina en error sin copiar ni instalar nada

### Scenario: el despliegue pisa lo que haya

- **GIVEN** un host con el payload de una release desplegado
- **WHEN** la persona corre `make dev-install`
- **THEN** el payload del árbol lo reemplaza, sin destino separado ni reversión

## Referencias

- **Flow** → [`../flow/install-ecosystem.md`](../flow/install-ecosystem.md) — el motor de instalación que este proceso invoca; nota especialmente su regla sobre dev builds.
