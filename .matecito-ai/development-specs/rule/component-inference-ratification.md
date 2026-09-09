# Capability — Inferencia y ratificación de componentes

- **Status:** Accepted
- **Date:** 2026-08-02
- **Components:** cli

## Propósito

Permitir que las herramientas de autoría (`development-spec-bootstrap` y `sdd-intake`) propongan componentes basándose en la estructura del repo o en el alcance del cambio, pero garantizar que **ninguna propuesta se aplica sin confirmación explícita del usuario**.

## Reglas de negocio

- **development-spec-bootstrap** infiere el set de componentes del repo a partir de manifiestos (`package.json`, `go.mod`, etc.), detectando múltiples manifiestos de proyecto y proponiendo directorios como componentes iniciales; **sdd-intake** infiere componentes a partir del alcance del cambio descrito en el brief.
- **Ninguna propuesta se aplica sin ratificación explícita del usuario.** El sistema propone; el usuario ratifica el set del repo (propuesta editada, aceptada una sola vez), los componentes de cada spec (propuesta aceptada spec por spec o descartada, al autorizarlo), y el valor del cambio (propuesta aceptada una sola vez por `sdd-intake` y reportado).
- La **inferencia del set ocurre una sola vez** en bootstrap — requiere manifiestos múltiples de proyecto; sin ellos, no se propone set alguno y la declaración queda manual. La propuesta es editable porque un manifiesto (e.g., `package.json`) puede corresponder a un componente cuya superficie real (`paths`) el usuario debe splitear — un manifiesto ambiguo (`package.json` en la raíz de un monorepo) no determina las fronteras de cada componente, así que bootstrap deja los `paths` iniciales pero editables. Una vez ratificado, el set no se re-infiere.
- La **inferencia per-spec** ocurre cada vez que se autoriza un spec en bootstrap — requiere que el eje `components` esté activo (`repo.components` en el config). Es una proposición independiente con ratificación spec por spec, nunca batch automático: el sistema propone la línea `- **Components:**` y espera confirmación del usuario.
- La **inferencia por-cambio** ocurre una sola vez en `sdd-intake` — requiere el eje activo y una descripción del alcance que mapee a componentes declarados. La ratificación ocurre por `sdd-intake` y reportado, y el valor muere con el cambio (nunca se propaga a `sdd-spec` ni `sdd-archive`, que conservan las líneas por-capability existentes sin alterar).
- Las **vías de ratificación son enumeradas y disjuntas**, y cada una ratifica su propia proyección y ninguna otra: (a) el **set del repo** — una sola vez, editable, al declararlo; (b) los **componentes de un spec** — spec por spec, al autorizarlo; (c) el **valor por-cambio** — una sola vez por cambio, por `sdd-intake` y reportado.
- Una ratificación de la vía (c) MUST NOT contar como ratificación de la vía (b): no autoriza escribir la línea de ningún spec, ni de uno solo ni en batch. La prohibición de batch automático cubre **cualquier** propagación desde una granularidad más gruesa hacia una más fina.

## Escenarios

### Scenario: sin manifiestos múltiples no se propone nada

- **verification:** `standing → development-spec-bootstrap`
- **GIVEN** un repo con un solo manifiesto de proyecto
- **WHEN** se ofrece declarar los componentes
- **THEN** no se propone set alguno y la declaración queda manual

### Scenario: nada se escribe sin ratificar

- **verification:** `standing → development-spec-bootstrap`
- **GIVEN** una inferencia de componentes, del set o de un spec
- **WHEN** el usuario no la confirma
- **THEN** no se escribe ni el bloque `repo` ni la línea del header

### Scenario: la ratificación por-cambio no ratifica ningún spec

- **verification:** `deferred → intake-components-flag`
- **GIVEN** un valor por-cambio ratificado una sola vez por `sdd-intake` y reportado
- **WHEN** se autoriza o se archiva un spec tocado por ese cambio
- **THEN** la línea de ese spec sigue requiriendo su propia confirmación spec por spec

### Scenario: ratificar el set no ratifica las líneas de los specs

- **verification:** `standing → development-spec-bootstrap`
- **GIVEN** un set del repo recién ratificado
- **WHEN** se recorre el store
- **THEN** ningún spec recibe línea `Components:` por eso — cada uno se ratifica al autorizarlo

## Referencias

- **Rule** → [`./repo-component-declaration.md`](./repo-component-declaration.md) — la declaración del set.
- **Rule** → [`./capability-spec-components-axis.md`](./capability-spec-components-axis.md) — la línea en el spec.
- **Proceso** → bootstrap es un artefacto del payload que aplica este comportamiento; no vive en capability-specs.
