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

## Hard Rules

- Read proposal, spec, design, and tasks before judging implementation.
- Execute relevant tests; static analysis alone is never verification.
- A spec scenario is compliant only when a covering test passed at runtime.
- Compare specs first, design second, task completion third.
<!-- matecito-ai: also verify the changed code respects the ADRs it touched (scoped to this change, not a full catalog audit) — see Execution Step 6b -->
- Verify the changed code respects the ADRs it touched (scoped to this change); a violation is CRITICAL `ADR-VIOLATION`.
- Do not fix issues; report them for the orchestrator/user.
- Persist `verify-report` to Engram (engram mode) or inline-only (none mode). <!-- matecito-ai: engram-only -->
- If Strict TDD is active, load `strict-tdd-verify.md` from this skill directory; if inactive, never load it.
- When `ui-test == needed` and `uiTest.available = ✅`, run the ProofShot UI step (one session per run); otherwise skip silently.
- Return the Section D envelope from `../_shared/sdd-phase-common.md`.

## Decision Gates

| Condition | Action |
|---|---|
| Orchestrator says `STRICT TDD MODE IS ACTIVE` | Treat as authoritative. |
| Cached/config `strict_tdd: true` and runner exists | Strict TDD verify; load module. |
| Strict TDD false or no runner | Standard verify; skip TDD checks. |
| Task incomplete | CRITICAL for core task, WARNING for cleanup task. |
| Test command exits non-zero | CRITICAL. |
| Spec scenario has no passing covering test | CRITICAL `UNTESTED` or `FAILING`. |
| Design deviation exists | WARNING unless it breaks a spec. |
<!-- matecito-ai: ADR violation in the changed code is CRITICAL -->
| Changed code violates an ADR it touched | CRITICAL `ADR-VIOLATION` (cite the ADR). |
| `ui-test != needed` OR `uiTest.available` absent or ❌ | Skip UI step silently — no mention, no UI Verdict section. |
| Scenario step target matches `@e\d+` | CRITICAL — reject authored runtime ref; scenario FAILS static validation. |
| Any per-scenario STATE assertion FAIL | CRITICAL — blocks archive. |
| Session-level error gate FAIL (consoleErrorCount or serverErrorCount > 0) | CRITICAL — blocks archive. |

## Execution Steps

1. Load relevant skills via shared SDD Section A.
2. Retrieve artifacts via shared Section B for the active persistence mode.
3. Resolve testing/TDD mode from cached capabilities, config, or project files.
3b. UI-test gate: read `ui-test` from the spec artifact and `uiTest.available` from `sdd/{project}/testing-capabilities`. If `ui-test != needed` OR `uiTest.available = ❌` OR either is absent → silently skip all UI steps (3c–3e, 3f). No mention, no UI Verdict section.
3c. Static validation (gate passed): for every scenario in `ui-scenarios`, reject any step target matching `@e\d+`. A matched target is CRITICAL — fail that scenario immediately.
3d. ProofShot session (gate and static validation passed): generate a collision-safe `outputDir` (`proofshot-artifacts/{change}-{timestamp}-{random}/`); `proofshot start --run "{devServer.command}" --port {port} --output {outputDir}`; for EACH scenario drive its steps then take a LIVE agent-browser `snapshot` and evaluate `visible`/`text_contains` STATE assertions against it; after ALL scenarios `proofshot stop`; read `SUMMARY.md` aggregates `consoleErrorCount`/`serverErrorCount` for the session-level ERROR GATE; delete `{outputDir}/session.webm` by default (retain only with explicit `retain-video` flag).
3e. SPLIT verdict: `ui-verdict = (all STATE assertions PASS) AND (error gate PASS)`; any FAIL → CRITICAL → blocks archive.
3f. Append `## UI Verdict` to the report: per-scenario STATE table (`Scenario | STATE | Failure Reason`), session-level ERROR GATE row (`consoleErrorCount`, `serverErrorCount`, PASS/FAIL), artifact path `proofshot-artifacts/{outputDir}/`.
4. Count completed and incomplete tasks.
5. Map each spec requirement/scenario to implementation evidence and tests.
6. Check design decisions against changed code.
<!-- matecito-ai: verify the change respects the ADRs it touched -->
6b. Check ADR compliance (scoped to THIS change). For each ADR listed in the design's "ADR Alignment" section (or, if absent, the ADRs in `.matecito-ai/adr/<domain>/` for the domains this change touched), confirm the implemented code actually honors that ADR's concrete rules (e.g. auth mechanism, error format, validation location, layer dependencies). This is scoped to the current change — do NOT audit the whole ADR catalog here. Report any violation as CRITICAL `ADR-VIOLATION` (cite the ADR). If `.matecito-ai/adr/` does not exist, skip this step.
7. Run test, build/type-check, and coverage commands when available.
8. Build the behavioral compliance matrix from actual test results.
9. Persist and return the verification report.

## Output Contract

Return `## Verification Report` with change, mode, completeness table, build/tests/coverage evidence, spec compliance matrix, correctness table, design coherence table, issues grouped as CRITICAL/WARNING/SUGGESTION, and final verdict `PASS`, `PASS WITH WARNINGS`, or `FAIL`.

## References

- [references/report-format.md](references/report-format.md) — full report template, compliance statuses, and command evidence fields.
- [references/ui-scenarios-schema.md](references/ui-scenarios-schema.md) — `ui-scenarios` block schema: field definitions, step primitives, target rules, `wait` primitive, assertion classes, validation rules.
- [strict-tdd-verify.md](strict-tdd-verify.md) — load only when Strict TDD is active.
- `../_shared/sdd-phase-common.md` — skill loading, retrieval, persistence, and return envelope.
