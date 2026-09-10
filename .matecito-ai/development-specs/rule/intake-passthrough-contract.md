# Capability — El contrato passthrough de intake

- **Status:** Accepted
- **Date:** 2026-08-31
- **Components:** cli

## Propósito

El intake actual corre un ciclo de dos pasos (Pass 1 formula discovery, Pass 2 investiga/clasifica/triage). El cambio de esta spec lo reduce a un passthrough: recibe el pedido, decide tres flags sin investigar, y emite el brief en un solo paso.

## Actores

- **Intake**: decide (nunca ratifica) los tres flags; corre el early guard de EDRs; emite el brief
- **Orquestador**: reporta los valores de flags en una línea de aviso; no espera en ellos
- **Usuarios downstream** (sdd-design, sdd-spec, sdd-verify, etc.): leen los flags decididos sin re-preguntar

## Precondiciones

- El pedido ha llegado a intake ya estructurado (request, domain, etc.)
- Intake aún tiene acceso a su early guard de EDRs (si el store de records está activo)

## Flujo principal

1. Intake recibe el pedido
2. Decide los tres flags (diagram, ui-test, components) sin investigar
3. Si el store de EDRs está activo, corre el early guard (Aceptados que conflictúan con el pedido)
4. Emite el brief con esos tres flags en `### Classification`
5. Orquestador reporta los valores en una línea de aviso
6. Próxima fase se despacha

## Ramas / flujos alternativos

- **Conflicto con un EDR Aceptado** → early guard lo levanta; es un blocker de flow
- **Store de EDRs ausente** → early guard se salta silenciosamente
- **Pedido ambiguo** → intake decide lo mejor que puede; no pregunta; no hay Pass 2

## Casos borde

- **Cambio trivial** → intake sigue siendo passthrough; no le importa el tamaño
- **Cambio grande** → intake sigue siendo passthrough; sigue sin clasificar
- **Cambio sobre dos superficies (cli + api)** → intake decide lo mejor que puede sobre componentes; la decisión viaja reportada, no ratificada

## Reglas de negocio

- Intake NUNCA investigates, NUNCA clasificará tamaño, NUNCA recomendará lane
- Los tres flags son decididos por intake, reportados por orquestador, actuados por sus lectores
- Si un flag es decidido mal, se corrige corrigiendo el brief, una vez, antes de que nada se despache
- Intake no retorna `needs-input` NUNCA bajo este contrato
- El brief tiene exactamente cuatro secciones (Request, Classification, Early Guard, Next) y nada más

## Entidades y estados

- **Brief de intake** — contiene el request estructurado y los tres flags decididos; estado: emitido una sola vez, nunca revisitado
- **Tres flags** — (diagram, ui-test, components); cada uno: decidido por intake, reportado por orquestador, actuado sin confirmación

## Errores de cara al actor

- **Flag decidido incorrectamente**: se corrige corrigiendo el brief
- **Conflicto con un EDR**: early guard lo levanta, brief no se emite
- **Pedido imposible de entender**: intake decide lo mejor que puede; no pregunta

## Requisitos

### Requisito: The intake phase is a passthrough

The intake phase MUST receive the raw request, decide the three decision flags, run the decision-record early guard when that store is active, and emit the brief — in one pass, one block. It MUST NOT investigate, MUST NOT classify a size, MUST NOT triage or recommend a lane, and MUST NOT return `needs-input`. A re-dispatch carrying the original request plus a user's correction is a **fresh pass over new input**, not a second pass with memory: the contract above is unchanged by it, and the re-run overwrites the same artifact key.

#### Scenario: One pass, one block

- GIVEN a raw request
- WHEN the intake phase runs
- THEN it returns once, with the brief, and never asks for anything first

#### Scenario: The removed classifications are gone

- GIVEN the brief after the change
- WHEN it is read
- THEN it carries no size, no lane, and no triage
- AND nothing downstream looks for them

#### Scenario: The early guard survives

- GIVEN a project whose decision-record store is active and a request that conflicts with an Accepted record
- WHEN the intake phase runs
- THEN it raises the conflict exactly as before this change

#### Scenario: The early guard stays silent with no store

- GIVEN a project with no decision-record store
- WHEN the intake phase runs
- THEN it skips the guard silently and mentions nothing

#### Scenario: A re-dispatch with a correction is a fresh pass

- GIVEN a brief the user corrected, and a re-dispatch carrying the original request plus those words
- WHEN the intake phase runs again
- THEN it runs one pass over that input, returns once, and still never returns `needs-input`
- AND it carries nothing over from the previous run beyond what its input states

### Requisito: The three decision flags are decided by intake and travel inside the brief

The intake phase MUST write all three flags into the brief's classification block. The orchestrator MUST report their values and MUST NOT wait on any one of them: no flag is offered, walked or confirmed as an item of its own. Each reader MUST act on the decided value without re-asking, and a reader that finds its flag absent MUST treat it as not-needed and close silently. The flags reach their readers inside a brief the user has confirmed as a whole, so the governing text MUST NOT state that a wrongly-decided flag reaches its reader with nobody having checked it.

Every place the governing text states the count MUST say three, and every enumeration MUST name those three and none other. No shipped text MUST still name a workspace-isolation flag, whether in the flags table, in intake's decision instructions, in the brief template, or in the onboarding material.

#### Scenario: The flags are reported, and the flow does not stop

- GIVEN a brief carrying all three flags
- WHEN it returns
- THEN their values are reported and nothing waits on any individual flag
- AND the only stop is the one question asked over the brief as a whole

#### Scenario: A reader acts on the decided value

- GIVEN a flag whose value the intake phase decided
- WHEN its reader runs
- THEN it acts on that value and does not re-ask or re-derive it

#### Scenario: The accepted cost is written down

- GIVEN the governing text for the flags after the change
- WHEN it is read
- THEN it no longer states that a wrongly-decided flag passes unchecked
- AND it states instead that a wrong value is corrected by correcting the brief, once, before anything is dispatched

#### Scenario: The count says three everywhere it is stated

- GIVEN the shipped text stating how many flags intake decides
- WHEN it is read at each place the count appears
- THEN all say three
- AND none is left at four

#### Scenario: No text still names the retired flag

- GIVEN the payload after the change
- WHEN it is searched for the isolation flag's name, in the flags table, in intake's instructions, in the brief template, and in the onboarding material
- THEN there are no matches

#### Scenario: Intake no longer decides the retired flag

- GIVEN the instructions intake follows to emit the brief
- WHEN they are read end to end
- THEN there is no step deciding an isolation value
- AND the brief it emits carries no such line
