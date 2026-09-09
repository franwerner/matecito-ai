# EDR — El token `blocking-test` de `sdd-design` es auto-declarado; el guard clasifica, no re-corre el test

- **Status:** Accepted
- **Date:** 2026-09-09

## Contexto
El blocking test de `sdd-design` era pura auto-evaluación: el ejecutor lo corría en su cabeza y publicaba
sólo el veredicto. Una decisión cuyo propio texto decía "esto necesita una cola y un worker que el
proyecto no tiene hoy" (eje 1, literalmente) llegó a `New Decisions` con `status: done`, y el guard no
tenía forma mecánica de notarlo — notarlo hubiera significado leer la prosa de la decisión y re-correr el
test, exactamente la interpretación que este guard existe para no hacer.

## Decisión
Cada item bajo `### New Decisions` declara un token `· blocking-test: none | infra | contract |
data-model`. El guard clasifica en base a ese token solamente — nunca leyendo la decisión y re-corriendo
el test por su cuenta. Mismo arreglo que `verify-checks:` para las desviaciones de diseño, un nivel más
arriba: la fase declara, el guard clasifica.

## Reglas verificables
- **[manual]** `payload/domains/development/CLAUDE.md`, sección `### Unresolved Decisions Guard (MANDATORY)`, declara que el guard clasifica el token `blocking-test` sin re-correr el test ni re-leer la prosa de la decisión.

## Alternativas consideradas
Que el guard lea la prosa de cada decisión y re-derive si el blocking test debería haber dado positivo —
descartada: es exactamente la interpretación mecánica que este guard existe para evitar, y ya falló en la
práctica con una decisión cuyo propio texto contradecía su veredicto declarado.

## Consecuencias
El token es la única evidencia de que el test corrió. El costo — que es auto-reportado y nada lo
verifica mecánicamente — es el mismo costo que ya asume el resto de los tokens `contested`.
