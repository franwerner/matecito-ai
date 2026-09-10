# EDR — Un record portado a otro store se adapta cuando su premisa no existe en el destino, en vez de copiarse verbatim

- **Status:** Accepted
- **Date:** 2026-09-10

## Contexto
Portar records de un repo origen a este store no es siempre una copia textual. Al portar `contracts/turn-claimed-by-explicit-call.md` y `structure/turn-mechanism-home-and-wiring-register.md` desde v1.50.0, ambos enmarcaban el mecanismo alrededor de **dos** write-moments — el merge de integración de cambio y el loop de cherry-pick de la consolidación. El primero no existe en este repo: sólo el bracket de cherry-pick de la corrida de consolidación sobrevive. Ambos registros también cargaban cinco links de `## Relacionados` a records de la fuente ratificados como nunca-portados: `turn-released-on-both-session-end-events`, `wiring-paragraph-defers-to-the-guard`, `merge-guard-registration-and-matching`, `shared-branch-turn-prose-homes`, `shared-branch-turn-wiring-points`.

## Decisión
Los dos records se portan adaptados, no verbatim: el framing de dos write-moments colapsa al único que sobrevive (el loop de cherry-pick de la consolidación); el rechazo del nombre `merge-lock` — que sólo `structure/turn-mechanism-home-and-wiring-register.md` cargaba — se mantiene ahí con la razón reforzada, ya que el único write-moment que queda *es* un loop de cherry-pick, que no es un merge; `contracts/turn-claimed-by-explicit-call.md` no menciona `merge-lock` en ningún punto, antes ni después del port; los cinco links de `## Relacionados` a records nunca-portados se eliminan y su sustantivo se pliega como prosa (mención en backticks, no link) en el Contexto de cada registro, así el rastro del razonamiento sobrevive sin un link que resuelva a nada. Los otros ocho records portados no llevan path de payload, ni link de Relacionados, ni cláusula de nivel de cambio — para ellos el port es textualmente mecánico (verbatim, sólo paths traducidos).

## Reglas verificables
- **[manual]** Ambos registros adaptados (`contracts/turn-claimed-by-explicit-call.md`, `structure/turn-mechanism-home-and-wiring-register.md`) enmarcan el mecanismo alrededor de un solo write-moment — el loop de cherry-pick de la consolidación — nunca dos.
- **[manual]** El rechazo del nombre `merge-lock` sobrevive en `structure/turn-mechanism-home-and-wiring-register.md`, con la razón reforzada: el único write-moment que queda es un loop de cherry-pick, que no es un merge. `contracts/turn-claimed-by-explicit-call.md` no lo menciona, porque nunca lo tuvo.
- **[manual]** Los cinco links a records nunca-portados (turn-released-on-both-session-end-events, wiring-paragraph-defers-to-the-guard, merge-guard-registration-and-matching, shared-branch-turn-prose-homes, shared-branch-turn-wiring-points) no aparecen como link resoluble en ningún registro portado; su sustantivo sobrevive como prosa en Contexto.
- **[manual]** Los otros ocho records portados no llevan ni path de payload, ni link de Relacionados, ni cláusula de nivel de cambio — su port es textualmente mecánico.

## Alternativas consideradas
Copiar los diez records verbatim, dejando que las citas a records nunca-portados queden colgando y el framing de dos write-moments describa una operación que no existe en este repo. Descartada: un link colgante y una cláusula que describe un movimiento inexistente son exactamente lo que el chequeo de tres partes de `sdd-verify` castiga, y dejarlo así traslada el defecto al siguiente sweep en vez de resolverlo en el port mismo.

## Consecuencias
Dos de los diez records portados difieren textualmente de su fuente, más allá de la traducción de paths — quien compare byte a byte contra el repo origen encuentra una diferencia real, documentada acá y no en un diff sin contexto. El rastro de las cinco relaciones retiradas sobrevive como prosa, no como link: un lector que siga esos nombres los encuentra mencionados pero no puede navegarlos — costo aceptado, igual que el resto de los links a records retirados en este cambio.

## Relacionados
- `relacionado-con` → [../contracts/turn-claimed-by-explicit-call.md](../contracts/turn-claimed-by-explicit-call.md) — uno de los dos records adaptados por esta decisión.
- `relacionado-con` → [turn-mechanism-home-and-wiring-register.md](turn-mechanism-home-and-wiring-register.md) — el otro record adaptado por esta decisión.
