# Capability — Línea `Components:` en el header del spec

- **Status:** Accepted
- **Date:** 2026-08-02
- **Components:** cli

## Propósito

Registrar, en cada capability-spec, qué componentes (superficies) implementan el comportamiento descrito, de modo que un desarrollador en un monorepo pueda ver rápidamente qué specs aplican a su parte del código.

## Reglas de negocio

- La línea solo aplica cuando el eje `components` está activo —hay `repo.components` declarado en el config— y sobre specs que existen en `.matecito-ai/development-specs/<type>/<capability>.md`.
- Al autorizar o actualizar un spec, el autor de specs agrega la línea `- **Components:** api, ui` en el header, junto a `Status` y `Date`; cada valor listado es una superficie declarada en `repo.components`.
- La línea se propone al minar y se ratifica junto con el spec; el validador de specs chequea que cada valor nombrado pertenece al set declarado.
- La línea es **multivaluada** — cada spec puede implementarse en múltiples componentes.
- Los valores son exactos: un componente es un `name` del set, no un path ni un prefijo.
- Si el eje está apagado (sin `repo.components`), la línea **no existe en ningún spec** — es presence-based.
- Un spec sin la línea (cuando el eje está activo) se reporta como WARNING durante la validación; un spec con un componente no declarado también se reporta como WARNING.
- La validación **nunca modifica** los specs; solo reporta — el consumidor lee la línea al minar o validar, nunca la escribe unilateralmente.
- La línea es la proyección **por-capability** del set: su valor describe esa capability y sólo se escribe al autorizar o actualizar **ese** spec. MUST NOT derivarse, copiarse ni ensancharse a partir del valor de otra proyección —en particular, del valor por-cambio ratificado en la INTAKE GATE—, ni al aplicar ni al archivar un cambio.
- La **forma** de la línea no cambia con esta regla: misma línea multivaluada, mismos valores exactos del set, mismos WARNINGs, mismo gate.

## Escenarios

### Scenario: multivaluada, y un valor no declarado se reporta

- **verification:** `standing → development-spec-validate`
- **GIVEN** un spec que lista dos superficies declaradas y una no declarada
- **WHEN** se valida el store
- **THEN** las declaradas son válidas y la otra se reporta como WARNING

### Scenario: specs previos al eje (WARNING, sin modificación)

- **verification:** `standing → development-spec-validate`
- **GIVEN** specs sin la línea y un `repo.components` recién declarado
- **WHEN** se valida el store
- **THEN** cada spec sin la línea se reporta como WARNING y ningún archivo se modifica

### Scenario: archivar un cambio no escribe la línea

- **verification:** `deferred → intake-components-flag`
- **GIVEN** un cambio con un valor de componentes ratificado a nivel de cambio y una capability tocada por él que no tiene línea `Components:`
- **WHEN** el cambio se archiva
- **THEN** la capability sigue sin la línea y la validación la reporta como WARNING, igual que antes del cambio

### Scenario: el contrato de la línea es el mismo que antes del cambio

- **GIVEN** un spec cuya línea `Components:` era válida antes de este cambio
- **WHEN** se valida después del cambio
- **THEN** sigue siendo válida, con la misma forma y los mismos findings

## Referencias

- **Rule** → [`./repo-component-declaration.md`](./repo-component-declaration.md) — la declaración del set de componentes.
- **Rule** → [`./component-inference-ratification.md`](./component-inference-ratification.md) — cómo se proponen y ratifican los componentes de un spec.
