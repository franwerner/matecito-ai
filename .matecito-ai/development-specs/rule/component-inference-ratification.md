# Capability — Inferencia y ratificación de componentes

- **Status:** Accepted
- **Date:** 2026-08-02
- **Components:** cli

## Propósito

Permitir que las herramientas de autoría (`development-spec-bootstrap` y `development-spec-mine`) propongan los componentes de un spec basándose en la estructura del repo o en la ruta del archivo, pero garantizar que **ninguna propuesta se aplica sin confirmación explícita del usuario**.

## Reglas de negocio

- **development-spec-bootstrap** infiere el set de componentes del repo a partir de manifiestos (`package.json`, `go.mod`, etc.), detectando múltiples manifiestos de proyecto y proponiendo directorios como componentes iniciales; **development-spec-mine** propone componentes de cada spec escaneado basándose en el prefijo de path más largo que coincida con un componente declarado.
- **Ninguna propuesta se aplica sin ratificación explícita del usuario.** El sistema propone; el usuario ratifica el set del repo (propuesta editada, aceptada una sola vez) y los componentes de cada spec (propuesta aceptada spec por spec o descartada).
- La **inferencia del set ocurre una sola vez** en bootstrap — requiere manifiestos múltiples de proyecto; sin ellos, no se propone set alguno y la declaración queda manual. La propuesta es editable porque un manifiesto (e.g., `package.json`) puede corresponder a un componente cuya superficie real (`paths`) el usuario debe splitear — un manifiesto ambiguo (`package.json` en la raíz de un monorepo) no determina las fronteras de cada componente, así que bootstrap deja los `paths` iniciales pero editables. Una vez ratificado, el set no se re-infiere.
- La **inferencia per-spec** (bootstrap por autoría; mine por longest-prefix matching de la ruta) ocurre cada vez que se autoriza o mina un spec — requiere que el eje `components` esté activo (`repo.components` en el config) y, para mine, paths declarados con los que hacer matching. Es una proposición independiente con ratificación spec por spec, nunca batch automático: el sistema propone la línea `- **Components:**` y espera confirmación del usuario.
- El criterio de mine es el **prefijo declarado más largo** — la especificidad que el autor expresó al declarar la ruta más profunda; si hay overlap (e.g., `internal` vs `internal/tui`), gana el más profundo. Si ningún `paths` coincide con la ruta escaneada, mine **no propone componente** para ese spec — la validación lo reportará como WARNING, no una adivinanza.

## Escenarios

### Scenario: sin manifiestos múltiples no se propone nada

- **GIVEN** un repo con un solo manifiesto de proyecto
- **WHEN** se ofrece declarar los componentes
- **THEN** no se propone set alguno y la declaración queda manual

### Scenario: al minar, sin prefijo que matchee no se propone

- **GIVEN** una ruta escaneada que ningún `paths` declarado prefija
- **WHEN** se arma el candidato
- **THEN** viene sin componentes, no con una adivinanza

### Scenario: nada se escribe sin ratificar

- **GIVEN** una inferencia de componentes, del set o de un spec
- **WHEN** el usuario no la confirma
- **THEN** no se escribe ni el bloque `repo` ni la línea del header

## Referencias

- **Rule** → [`./repo-component-declaration.md`](./repo-component-declaration.md) — la declaración del set.
- **Rule** → [`./capability-spec-components-axis.md`](./capability-spec-components-axis.md) — la línea en el spec.
- **Proceso** → bootstrap y mine son artefactos del payload que aplican este comportamiento; no viven en capability-specs.
