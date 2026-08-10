---
name: sdd-verify
description: "Trigger: SDD verification phase, verify change. Execute tests and prove implementation matches specs, design, and tasks."
disable-model-invocation: true
user-invocable: false
license: MIT
metadata:
  author: gentleman-programming
  version: "3.0"
  delegate_only: true
---

> **ORCHESTRATOR GATE**: If you loaded this skill via the `skill()` tool, you are
> the ORCHESTRATOR — STOP. Do NOT execute these instructions inline. Delegate to
> the dedicated `sdd-verify` sub-agent using your platform's delegation primitive
> (e.g., `task(...)`, sub-agent invocation, etc.). This skill is for EXECUTORS
> only.

## Activation Contract

Run when the orchestrator launches verification for an SDD change. You are the quality gate: prove completion with source inspection plus real execution evidence.

<!-- matecito-ai: verification fan-out. The orchestrator dispatches N instances of this skill in ONE
     message — one per group whose gate is active — instead of a single instance doing every check in
     series. `subverifier-groups.md` fixes the partition; this file's Hard Rules and Execution Steps
     stay the checks themselves, unchanged, and each instance runs only the subset its group owns. -->
### Group Scope (read before anything else)

You are dispatched for **one** of seven named groups (`execution`, `correctness`, `design-coherence`,
`edr-coherence`, `spec-coherence`, `ui`, `decision-gaps`) — the orchestrator names it. Read
`~/.claude/references/phase-returns/sdd-verify/subverifier-groups.md` first: it maps each Execution
Step below to the group that owns it, and defines the Sub-Report you return **instead of** the full
`## Verification Report` in "Output Contract" (that block is what the *orchestrator* assembles from
every dispatched group's Sub-Report — see "Output Contract" for where the split falls).

Run only the steps your group owns. A step owned by a different group is not yours: do not run it, do
not check it, and say nothing about it in your Sub-Report. Being dispatched for a group also means its
**gate is already active** — you do not re-check EDR-store presence, capability-spec-store presence,
`ui-test`, or `flagDecisionGaps`; the orchestrator resolved that before dispatching you.

## Hard Rules

<!-- matecito-ai: `apply-progress` faltaba en esta lista aunque el agente sí lo recupera: es donde viaja
     la declaración de desvíos de `sdd-apply` (paso 6), sin la cual no se puede clasificar un desvío.
     El brief entra por el flag `ui-test` (paso 3b), que nunca estuvo en el spec. -->
- Read the intake brief, proposal, spec, design, tasks, and apply-progress before judging implementation — but only the ones your group's steps actually need (see "Group Scope" above); do not fetch an artifact no step of yours reads.
<!-- matecito-ai: apply-progress is no longer only the deviations and the file list — it carries the
     executable half of the UI scenarios, without which the browser run has nothing to drive. -->
- UI verification pairs the spec's **behavioral** scenarios with the **executable counterparts** in apply-progress, by `name`. A behavioral scenario with no counterpart is `UNTESTED`/CRITICAL; a missing counterparts section while the spec declares scenarios is CRITICAL, never a silent skip.
- Execute relevant tests; static analysis alone is never verification.
- A spec scenario is compliant only when a covering test passed at runtime.
- Compare specs first, design second, task completion third.
<!-- matecito-ai: also verify the changed code respects the EDRs it touched (scoped to this change, not a full catalog audit) — see Execution Step 6b -->
- Verify the changed code respects the EDRs it touched (scoped to this change); a violation is CRITICAL `EDR-VIOLATION`.
<!-- matecito-ai: also validate the implemented behavior against the durable capability-specs (`.matecito-ai/development-specs/`) — accumulated behavior state, not only the change delta — see Execution Step 6d -->
- Validate the implemented behavior against the durable capability-specs of the capabilities this change touches (`.matecito-ai/development-specs/<type>/<capability>.md`, when present), **scoped to `Status: Accepted` only**; a divergence from the accumulated intended behavior of an `Accepted` spec is CRITICAL `SPEC-VIOLATION`. A spec with `Status: Inferred` is a non-authoritative draft — like `Draft`, it is skipped, never a contract, and can never fail verify.
<!-- matecito-ai: esta fase comparaba el spec durable contra el delta del cambio y nada más. La pregunta
     que nadie hacía es si el archivo respeta el contrato del store — y un merge que re-tipea la prosa
     no cambia el delta, así que tres pasadas seguidas dieron "sin derivas" mientras el archivo perdía
     las negritas de `- **GIVEN**` y escribía links que no resuelven a ningún lado. Lo encontró un
     `validate-artifact.js` corrido a mano, después de archivar. El motor ya existía; lo que faltaba
     era que esta fase lo invocara. -->
- When the change declares New or Modified capabilities, validate the durable store's **structure**, not only its behavior: a finding on a capability-spec this change touched is this change's own; a finding anywhere else is pre-existing and never pushes the verdict. See Execution Step 6e.
<!-- matecito-ai: verification-scope-token. Canonical definition lives at `~/.claude/references/spec/README.md`
     (deployed path) — this bullet cites it, it does not redefine the token's form or semantics. -->
- Read each scenario's `verification:` token (canonical definition: `~/.claude/references/spec/README.md`). A `deferred`/`standing` scenario is never classified `UNTESTED`, never raises a CRITICAL for missing coverage, and never pushes the verdict to FAIL for that reason; report it with its token so the owner of that verification stays visible, and exclude it from the compliance summary's denominator. The token exempts missing coverage, never a negative result: a `deferred`/`standing` scenario whose test ran and failed is still `FAILING`. A scenario with no `verification:` line is `in-scope` and classifies exactly as before this rule existed.
- Do not fix issues; report them for the orchestrator/user.
<!-- matecito-ai: persistence moved to the orchestrator, once, after the merge — see "Group Scope" above.
     A sub-verifier that persisted its own fragment would produce one `verify-report` write per group
     instead of one per run, and none of those partial writes would be the consolidated report. -->
- Never call `mem_save`. There is exactly one `verify-report` per run; the orchestrator persists it once, after merging every dispatched group's Sub-Report.
- If Strict TDD is active, load `strict-tdd-verify.md` from this skill directory; if inactive, never load it.
<!-- matecito-ai: el flag `ui-test` lo decide y lo escribe `sdd-intake` en su brief; el spec nunca lo
     llevó ni lo mencionó. Leerlo del spec devolvía siempre "ausente", y como el gate cierra en
     silencio ante un valor ausente, la verificación de UI no corría nunca y tampoco se quejaba. -->
- When `ui-test == needed` (read from the intake brief) and `uiTest.available = ✅`, run the ProofShot UI step (one session per run); otherwise skip silently.
<!-- matecito-ai: un desvío de diseño que sobrevive hasta acá ya no es "apply se desvió y lo anotó":
     `sdd-apply` debe devolver `blocked` cuando el hueco afecta lo que está por escribir. Lo que llega
     es o detalle de ejecución (no es asunto de esta fase) o una decisión tomada sin preguntar
     (violación de un hard-stop) — y sólo CRITICAL bloquea el archivado. -->
- A design deviation is CRITICAL unless `sdd-apply` declared it not verifiable against the design; see Execution Step 6.
<!-- matecito-ai: la forma del retorno dejó de vivir acá. Una copia inline más era una copia más para desincronizar — que es exactamente cómo `## Decision Gaps` y `## UI Verdict` terminaron sin lugar en el reporte. -->
<!-- matecito-ai: this used to say "return the Section D envelope" — that is now the ORCHESTRATOR's job,
     performed once after every dispatched group's Sub-Report is in. A group-scoped instance of this
     skill returns the Sub-Report envelope from `subverifier-groups.md` instead; see "Group Scope". -->
- Return the Sub-Report envelope `~/.claude/references/phase-returns/sdd-verify/subverifier-groups.md` defines — never the Section D phase envelope and never the full `## Verification Report` block.

## Decision Gates

| Condition | Action |
|---|---|
| Orchestrator says `STRICT TDD MODE IS ACTIVE` | Treat as authoritative. |
| Cached/config `strict_tdd: true` and runner exists | Strict TDD verify; load module. |
| Strict TDD false or no runner | Standard verify; skip TDD checks. |
| Task incomplete | CRITICAL for core task, WARNING for cleanup task. |
| Test command exits non-zero | CRITICAL. |
| Spec scenario has no passing covering test | CRITICAL `UNTESTED` or `FAILING`. |
<!-- matecito-ai: verification-scope-token — see the Hard Rules bullet above; this row and the next classify per that rule. -->
| Scenario carries `verification: deferred → <change>` or `standing → <owner>` | Report it with its token; not `UNTESTED`, no coverage CRITICAL, does not push the verdict to FAIL; excluded from the compliance summary's denominator. |
| A `deferred`/`standing` scenario's test ran and failed | Classify the failure normally (`FAILING`) — the token exempts missing coverage, never a negative result. |
<!-- matecito-ai: recalibrado. La escapatoria vieja ("WARNING unless it breaks a spec") sobra por dos
     lados: si rompe un spec, la Spec Compliance Matrix ya lo agarra sola como CRITICAL, y dejarla
     escrita sugería que un desvío que NO rompe un spec es inofensivo. La única rebaja legítima es la
     que declara `sdd-apply`: un desvío que no toca ningún requisito del spec ni ningún `criteria:` de
     tarea es detalle de ejecución. Sin declaración se asume verificable: el default seguro es el estricto. -->
<!-- matecito-ai: a deviation now lives under one of two apply-progress sections (`### Unmandated
     Forks` or `### Mandated Departures`) and carries a `mandate:` token alongside `verify-checks:` —
     but `mandate:` says WHO decided it, not whether verify audits it. Classification stays keyed on
     `verify-checks:` alone, from either section. -->
| Design deviation carrying `verify-checks: yes` — or no token at all, or a hedged one — found under either `### Unmandated Forks` or `### Mandated Departures` | CRITICAL `DESIGN-DEVIATION` (cite the deviation and the design decision it departs from). |
| Design deviation carrying `verify-checks: no` — found under either section | WARNING — execution detail, not this phase's business. |
<!-- matecito-ai: EDR violation in the changed code is CRITICAL -->
| Changed code violates an EDR it touched | CRITICAL `EDR-VIOLATION` (cite the EDR). |
| `ui-test != needed` OR `uiTest.available` absent or ❌ | Skip UI step silently — no mention, no UI Verdict section. |
| Behavioral scenario with no counterpart in apply-progress | CRITICAL `UNTESTED` for that scenario. |
| Counterpart declaring `covers: partial: …` | WARNING with what it leaves out — **even when it PASSES**. |
| Counterpart carrying no `covers` | WARNING, read as `partial: undeclared`. |
| Counterpart whose `name` matches no behavioral scenario | WARNING — orphan; report, do not run. |
| No `### UI Scenario Counterparts` at all, while the spec declares scenarios | CRITICAL — the UI check could not run; say so. NOT a silent skip. |
| Scenario step target matches `@e\d+` | CRITICAL — reject the runtime ref; scenario FAILS static validation. |
| Any per-scenario STATE assertion FAIL | CRITICAL — blocks archive. |
| Session-level error gate FAIL (consoleErrorCount or serverErrorCount > 0) | CRITICAL — blocks archive. |

## Execution Steps

<!-- matecito-ai: each step below still belongs to exactly one group, per `subverifier-groups.md` — this
     note is not repeated at every step; "Group Scope" above already says run only your own. -->
Run only the steps `subverifier-groups.md` maps to your assigned group; every other step is not yours.

1. Load relevant skills via shared SDD Section A.
2. Retrieve artifacts via shared Section B for the active persistence mode — scoped to what your group's own steps need (see "Group Scope").
3. Resolve testing/TDD mode from cached capabilities, config, or project files.
<!-- matecito-ai: el flag sale del brief de intake (`### Classification` → `UI test:`), NO del spec —
     el spec no lo contiene. El brief siempre existe porque intake es fase base, así que no hay fallback. -->
<!-- matecito-ai: this gate is now resolved by the ORCHESTRATOR, before dispatch — it decides whether
     the `ui` group is even launched. No sub-verifier runs this step; `subverifier-groups.md` starts the
     `ui` group's ownership at 3b-bis, deliberately after this one. Kept here as the record of what the
     orchestrator's gate check does, and because 3b-bis reads "the gate passed" as its precondition. -->
3b. UI-test gate (orchestrator, pre-dispatch — not run by any sub-verifier): read `ui-test` from the **intake brief** (`sdd/{change-name}/intake`, line `- UI test: {needed|not-needed}` under `### Classification`) and `uiTest.available` from `sdd/{project}/testing-capabilities` (literal key, in its `### UI Test` table — canonical-keys list in `~/.claude/skills/sdd-init/references/init-details.md`). If `ui-test != needed` OR `uiTest.available = ❌` OR either is absent → do not dispatch the `ui` group at all; the consolidated report carries no mention and no UI Verdict section. The scenarios come in two halves: the **behavioral** ones from the spec's `ui-scenarios` block and their **executable counterparts** from `### UI Scenario Counterparts` in apply-progress (step 3b-bis).
<!-- matecito-ai: two halves (schema Parts 1 and 2) because the spec cannot know a route or an accessible
     name before the code exists, and must not pin them anyway. Pairing is by `name` — the key this phase
     already used for the verdict table, not a new coupling. -->
3b-bis. Coverage check (gate passed, before executing anything): pair behavioral scenarios with counterparts **by `name`**, verbatim. No counterpart for a behavioral scenario → `UNTESTED`, **CRITICAL**. Orphan counterpart (name matches nothing) → WARNING, reported and not run. No `### UI Scenario Counterparts` section at all while the spec declares scenarios → **CRITICAL** finding, never a silent skip: the gate closes silently when the flag is off, never when the check was warranted and the input was missing. Read each counterpart's `covers`: `partial: …` (or absent, read as `partial: undeclared`) → WARNING naming what it leaves out, **reported even when the scenario PASSES** — a green row that proves less than its scenario claims is worse than a red one.
3c. Static validation (counterparts found): for every counterpart, reject any step target matching `@e\d+`. A matched target is CRITICAL — fail that scenario immediately.
3d. ProofShot session (gate and static validation passed): generate a collision-safe `outputDir` (`proofshot-artifacts/{change}-{timestamp}-{random}/`); `proofshot start --run "{uiTest.devServer.command}" --port {port} --output {outputDir}` (same `### UI Test` table); for EACH scenario drive its steps then take a LIVE agent-browser `snapshot` and evaluate `visible`/`text_contains` STATE assertions against it; after ALL scenarios `proofshot stop`; read `SUMMARY.md` aggregates `consoleErrorCount`/`serverErrorCount` for the session-level ERROR GATE; delete `{outputDir}/session.webm` by default (retain only with explicit `retain-video` flag).
3e. SPLIT verdict: `ui-verdict = (all STATE assertions PASS) AND (error gate PASS)`; any FAIL → CRITICAL → blocks archive.
3f. Return `ui_verdict` and `error_gate` in your Sub-Report — the `ui` group's data keys, per your row in `subverifier-groups.md`. They carry the per-scenario STATE results, the session-level ERROR GATE (`consoleErrorCount`, `serverErrorCount`, PASS/FAIL) and the artifact path `proofshot-artifacts/{outputDir}/`. The orchestrator renders them as `## UI Verdict` in the consolidated report, at the position and `##` level `~/.claude/references/phase-returns/sdd-verify/sdd-verify.md` fixes — you do not render that markdown section yourself.
4. Count completed and incomplete tasks.
5. Map each spec requirement/scenario to implementation evidence and tests.
<!-- matecito-ai: verification-scope-token. Definition owned by `~/.claude/references/spec/README.md`; this step only applies it while building the matrix. -->
5b. Before classifying a scenario in the Spec Compliance Matrix, read its `verification:` token — this matrix is scoped to the change's own delta requirements, so this step applies to the delta spec's scenarios. `deferred → <change>` or `standing → <owner>` → the `Result` cell reads `➖ OUT-OF-SCOPE — {token}` (the token inlined verbatim); exclude that scenario from the compliance summary's denominator. If its test ran and failed anyway, classify it `FAILING` instead — the token never exempts a negative result. No `verification:` line → `in-scope`; classify with the other four `Result` values exactly as before. The same token-reading habit carries over to any durable capability-spec loaded at step 6d — but 6d fills a different table (`### Coherence (Capability-Specs)`, columns `Followed?`/`Severity`/`Notes`) that has no `Result` cell at all; see 6d for its own vocabulary.
<!-- matecito-ai: `sdd-apply` ya declara, por cada desvío, si esta fase lo va a chequear contra el
     diseño (si toca un requisito del spec o un `criteria:` de tarea) — pero nadie leía esa declaración,
     así que la única rebaja legítima de severidad no tenía forma de aplicarse. Acá se lee. -->
6. Check design decisions against changed code. Read **both** `### Unmandated Forks` and `### Mandated Departures` in the apply-progress artifact: for each deviation recorded under either section, take its `verify-checks:` declaration of whether this phase will check it against the design (i.e. whether it touches a spec requirement or a task `criteria:`). Then classify per the Decision Gates table — declared verifiable, or no declaration at all, and the code does diverge from the design → CRITICAL `DESIGN-DEVIATION`; declared NOT verifiable → WARNING. A deviation with no declaration is treated as verifiable: the safe default is the strict one, and only apply has the context to downgrade it. **Classify by `verify-checks:` alone — a deviation's `mandate:` (`covered` / `forced` / `chosen`) says who decided it, never whether it is audited.** Record every deviation in `### Coherence (Design)` and list it under its severity in `### Issues Found`. If apply-progress has no such sections (or both are the `None` sentinel), there are no declared deviations — judge the code against the design as usual.
<!-- matecito-ai: verify the change respects the EDRs it touched -->
6b. Check EDR compliance (scoped to THIS change). For each EDR listed in the design's "EDR Alignment" section (or, if absent, the EDRs in `.matecito-ai/edr/<domain>/` for the domains this change touched), confirm the implemented code actually honors that EDR's concrete rules (e.g. auth mechanism, error format, validation location, layer dependencies). This is scoped to the current change — do NOT audit the whole EDR catalog here. Report any violation as CRITICAL `EDR-VIOLATION` (cite the EDR). If `.matecito-ai/edr/` does not exist, skip this step.
<!-- matecito-ai: this hard rule had no section of its own in the report — its only home was `### Issues
     Found`, which carries the finding only if there WAS a violation. With nothing found, no trace
     remained that the check ran, so "no violation" and "never checked" were indistinguishable. -->
   Record the outcome in `### Coherence (EDRs)` — one row per applicable EDR, emitted **whenever the store is active, including the all-clear case**, and omitted entirely (with no mention of EDRs anywhere) when the store is absent or empty. Its position, columns and conditionality are fixed by `~/.claude/references/phase-returns/sdd-verify/sdd-verify.md`; the meaning of its cells by [references/report-format.md](references/report-format.md). A violation still goes under CRITICAL in `### Issues Found` as well — the table is the evidence the check ran, not a replacement for the finding.
<!-- matecito-ai: validate implemented behavior against the durable capability-specs (accumulated behavior, not just the change delta), scoped to Accepted only -->
6d. Check capability-spec compliance (durable behavior), **scoped to `Status: Accepted` specs only**. For each capability this change touches that has a durable spec under `.matecito-ai/development-specs/<type>/<capability>.md` (type ∈ flow|rule|lifecycle|process), load it and check its `Status` at the load point: `Accepted` → for each scenario, first read its `verification:` token (canonical definition: `~/.claude/references/spec/README.md`) — a `deferred → <change>` or `standing → <owner>` scenario's absent or not-yet-existing behavior is never itself a divergence; note the token verbatim instead of flagging it. For every other scenario, confirm the implemented code honors the accumulated intended behavior that spec describes — not only the delta introduced by this change; report any divergence as CRITICAL `SPEC-VIOLATION` (cite the capability-spec). The token exempts missing/absent behavior, never a negative result: if a `deferred`/`standing` scenario's behavior WAS directly exercised and it failed anyway, that failure is still a divergence and is reported as `SPEC-VIOLATION` normally. `Inferred` (like `Draft`) → skip — it is a non-authoritative draft, not a contract, and can never trigger `SPEC-VIOLATION`. If `.matecito-ai/development-specs/` does not exist, skip this step entirely.
<!-- matecito-ai: same hole as 6b had — a hard rule whose only trace in the report was `### Issues Found`,
     so the all-clear case was indistinguishable from a check that never ran. Worse here, because a
     skipped `Inferred` spec ALSO looks like a clean check unless the row says so. -->
   Record the outcome in `### Coherence (Capability-Specs)` — one row per touched capability that has a durable spec, emitted **whenever the store is active, all-clear case included**, and absent entirely when the store is not. A skipped `Inferred`/`Draft` spec still gets its row, marked `➖ Not checked`: it says a spec exists but is not yet a contract, which is information, not an absence. When the only scenarios that would otherwise diverge carry a `deferred`/`standing` token, the row still reads `✅ Yes` — cite the token(s) verbatim in `Notes`, the same exemption as above, never covering a divergence actually observed. Shape fixed by `~/.claude/references/phase-returns/sdd-verify/sdd-verify.md`, cell meanings by [references/report-format.md](references/report-format.md). A divergence also goes under CRITICAL in `### Issues Found`.
<!-- matecito-ai: 6d y 6e preguntan cosas distintas sobre el mismo archivo, y por eso ninguna cubre a la
     otra: 6d pregunta si el código honra lo que el spec dice, 6e si el spec honra el contrato del
     store. Un merge que re-tipea la prosa —negritas perdidas en `- **GIVEN**`, un link con los `../`
     mal contados— pasa 6d intacto, porque el comportamiento descrito no cambió. El motor mecánico ya
     existe y otras skills ya lo invocan; acá sólo se lo llama y se particionan sus findings. -->
6e. Validate the durable store's STRUCTURE — **only when this change declares New or Modified capabilities**; otherwise skip entirely, with no mention. Run the mechanical engine rather than re-deriving its checks by hand:

   ```
   node ~/.claude/scripts/validate-artifact.js --type capability-spec --store .matecito-ai/development-specs --root <repo root>
   ```

   Partition its JSON findings **by file**. A finding on a capability-spec this change created or modified is this change's own: map the tool's severity (`error` → CRITICAL, `warning` → WARNING, `nota` → SUGGESTION) and list it under `### Issues Found`. Every other finding is pre-existing — a store carries defects from changes that are not yours, and letting them push the verdict makes every verify fail for reasons nobody in this change caused. Report those as a single count line in `### Coherence (Capability-Specs)` `Notes`, never as findings. If the store does not exist, skip.
<!-- matecito-ai: decision-gap confirmation hook
Active ONLY when flagDecisionGaps=true (does NOT depend on EDRs existing). When active: read the tasks artifact; collect all `· edr: <domain>/<slug>` whose file `.matecito-ai/edr/<domain>/<slug>.md` does NOT exist — these are the decision gaps. For each gap: (a) check the task is `[x]` (complete); (b) check its `criteria:` passes in the shipped code (static inspection or test result). If (a) and (b) → `implemented: yes`; otherwise `implemented: no`. La sección viaja DENTRO del bloque que devolvés (que es el mismo texto que persistís): el gate del orquestador sólo lee el retorno, así que una sección que quedara únicamente en la copia persistida no dispararía nada. Si tiene al menos un `yes`, el orquestador puede disparar el mine gate post-verify. When flag off: do NOT add the section, do NOT mention anything — byte-identical behavior to before. -->
6c. (Decision-gap confirmation) Being dispatched for the `decision-gaps` group already means `flagDecisionGaps=true`: from the tasks artifact collect all `· edr: <domain>/<slug>` refs whose target file does NOT exist → these are decision gaps. For each: confirm task is `[x]` AND `criteria:` passes in shipped code → mark `implemented: yes/no`. Return `decision_gaps` in your Sub-Report. The orchestrator renders it as `## Decision Gaps` in the consolidated report — position, columns and `##` level fixed by `~/.claude/references/phase-returns/sdd-verify/sdd-verify.md`, which the mine gate matches literally.
7. Run test, build/type-check, and coverage commands when available.
8. Build the behavioral compliance matrix from actual test results.
<!-- matecito-ai: this step used to persist AND return the full report — both are now the orchestrator's,
     once, after merging every dispatched group. See "Group Scope" and "Output Contract" below. -->
9. Return your Sub-Report (see "Output Contract" below). Do not persist anything.

## Output Contract

<!-- matecito-ai: el contrato dice QUÉ tiene que llevar el reporte; la FORMA vive una sola vez, en el template. Antes esta sección era una tercera descripción del mismo bloque (con la plantilla y con el contrato del ejecutor) y las tres divergían. -->
<!-- matecito-ai: this section used to say "return the block" — that instruction now belongs to the
     ORCHESTRATOR, once, after merging every dispatched group. What follows is still true of the
     CONSOLIDATED report (so "what the report must carry" stands), but what YOU personally return is
     the Sub-Report from "Group Scope" / `subverifier-groups.md`, carrying only your group's keys. -->
The `## Verification Report` block **exactly as `~/.claude/references/phase-returns/sdd-verify/sdd-verify.md` defines it** — sections, titles, order, which ones are conditional, and what changes per status — is what the **orchestrator** assembles from every dispatched group's Sub-Report, not what any single one of you returns. You return the Sub-Report; the orchestrator validates the *consolidated* block against that same file, matching titles literally, so a section any group drops, renames or re-levels in its own fragment is a gate that never fires downstream.

What the consolidated report must CARRY (the template fixes how it looks): change identity and mode; task completeness; real build/test/coverage execution evidence; the spec compliance matrix; correctness and design-coherence tables; issues grouped as CRITICAL/WARNING/SUGGESTION; and a final verdict `PASS`, `PASS WITH WARNINGS`, or `FAIL`. Plus the conditional sections, each only when its gate holds: `### Coherence (EDRs)` (step 6b, whenever the decision store is active), `### Coherence (Capability-Specs)` (step 6d, whenever that store is active), the Strict TDD extension, `## UI Verdict` (step 3f), and `## Decision Gaps` (step 6c). Your own contribution to each is exactly the data key(s) your row in `subverifier-groups.md` owns.

<!-- matecito-ai: la regla del veredicto tenía que seguir a la recalibración del desvío: si un desvío
     CRITICAL pudiera convivir con un PASS, subir la severidad no cambiaría nada aguas abajo. -->
Verdict rules for design coherence: a CRITICAL `DESIGN-DEVIATION` forbids `PASS` and `PASS WITH WARNINGS` — the verdict is `FAIL`, like any other CRITICAL. A deviation apply declared not verifiable is a WARNING: it caps the verdict at `PASS WITH WARNINGS` and never blocks archive.

## References

- `~/.claude/references/phase-returns/sdd-verify/subverifier-groups.md` — the group partition, the
  Sub-Report contract you actually return, and the merge the orchestrator runs over it. Read this
  **before** the return template below.
- `~/.claude/references/phase-returns/sdd-verify/sdd-verify.md` — **the** shape of the *consolidated*
  report the orchestrator assembles and persists. Read it to know where your data keys land, not as
  the thing you personally emit.
- [references/report-format.md](references/report-format.md) — compliance statuses and what counts as command evidence: the meaning of the cells the template lays out.
- `~/.claude/references/spec/README.md` — canonical definition of the `verification:` token (form, values, default by absence). This skill only cites it and applies its own reaction to each value; it never redefines the token.
<!-- matecito-ai: el schema dejó de ser privado de esta skill: ahora es contrato compartido entre quien
     PRODUCE el bloque (`sdd-spec`) y quien lo CONSUME (esta fase), así que vive en references/ y se
     despliega a `~/.claude/references/`. La ruta relativa vieja ya no resuelve. -->
- `~/.claude/references/ui-scenarios-schema.md` — the UI-scenario contract, shared by **three** phases: Part 1 the behavioral half (`sdd-spec`), Part 2 the executable counterparts (`sdd-apply`), **Part 3 execution and coverage — yours**. Read Part 3 before running the UI step, and Part 2 to know the shape of what you are executing.
- [strict-tdd-verify.md](strict-tdd-verify.md) — load only when Strict TDD is active.
- `~/.claude/skills/_shared/sdd-phase-common.md` — Sections A and B bind you (skill loading, retrieval); Section D does not — see "Group Scope".
