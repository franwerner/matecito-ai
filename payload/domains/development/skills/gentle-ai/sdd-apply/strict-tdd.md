# Strict TDD Module — Apply Phase

> **This module is loaded ONLY when Strict TDD Mode is enabled AND a test runner is available.**
> If you are reading this, the orchestrator already verified both conditions. Follow every instruction.

## TDD Philosophy

TDD is not testing. TDD is **software design driven by tests**. You write a test that describes what the code SHOULD do, then write the minimum code to make it real. The tests design the API, the contracts, the behavior. Code is a side effect of tests.

### The Three Laws

1. **Do NOT write production code** until you have a failing test
2. **Do NOT write more test** than is necessary to fail
3. **Do NOT write more code** than is necessary to pass the test

## TDD Implementation Cycle

For EVERY task assigned to you, follow this cycle strictly:

```
FOR EACH TASK:
├── 0. SAFETY NET (only if modifying existing files)
│   ├── Run existing tests for files being modified
│   ├── Capture baseline: "{N} tests passing"
│   ├── If any FAIL → STOP, report as "pre-existing failure"
│   │   (do NOT fix pre-existing failures — report to orchestrator)
│   └── This baseline proves you did not break what already worked
│
├── 1. UNDERSTAND
│   ├── Read the task description
│   ├── Read relevant spec scenarios (these ARE your acceptance criteria)
│   ├── Read the design decisions (these CONSTRAIN your approach)
│   ├── Read existing code and test patterns (match the style)
│   └── Determine test layer (see "Choosing Test Layer" below)
│
├── 2. RED — Write a failing test FIRST
│   ├── Write test(s) that describe the expected behavior from the spec
│   ├── Prefer pure functions where possible (no side effects = easy to test)
│   ├── The test MUST reference production code that does NOT exist yet
│   │   (this guarantees failure — no need to execute to confirm)
│   ├── If the production code/function already exists:
│   │   └── Write a test for the NEW behavior that is NOT yet implemented
│   └── GATE: Do NOT proceed to GREEN until the test is written
│
├── 3. GREEN — Write the MINIMUM code to pass
│   ├── Implement ONLY what the failing test needs
│   ├── Fake It is VALID here (hardcoded return values are OK)
│   ├── EXECUTE tests → must PASS
│   │   ├── ✅ Passed → proceed to TRIANGULATE or REFACTOR
│   │   └── ❌ Failed → fix the implementation, NOT the test
│   └── GATE: Do NOT proceed until GREEN is confirmed by execution
│
├── 4. TRIANGULATE (MANDATORY for most tasks)
│   ├── DEFAULT: triangulation is REQUIRED. You need a compelling reason to skip it.
│   ├── Add a second test case with DIFFERENT inputs/expected outputs
│   ├── EXECUTE tests → if Fake It breaks (hardcoded no longer works):
│   │   └── Generalize to real logic (this is the whole point)
│   ├── Repeat until ALL spec scenarios for this task are covered
│   ├── Each triangulation pass: write test → run → fix implementation
│   ├── MINIMUM: at least 2 test cases per behavior (happy path + one edge case)
│   │   ├── One test with data that produces a NON-EMPTY/NON-TRIVIAL result
│   │   └── One test with data that exercises a DIFFERENT code path
│   ├── WATCH OUT for GREEN that passes trivially:
│   │   ├── If your test passes because the component/element isn't rendered → NOT a real GREEN
│   │   ├── If your test passes because a loop iterates 0 times → NOT a real GREEN
│   │   ├── If your test passes because the setup doesn't trigger the code path → NOT a real GREEN
│   │   └── A real GREEN means: production code RAN and produced the expected output
│   ├── Skip triangulation ONLY when ALL of these are true:
│   │   ├── The task is purely structural (config file, constant definition, type export)
│   │   ├── There is literally ONE possible output (no branching, no logic)
│   │   └── You explicitly note "Triangulation skipped: {reason}" in the evidence table
│   └── GATE: All spec scenarios for this task must have tests before REFACTOR
│
├── 5. REFACTOR — Improve without changing behavior
│   ├── Extract constants (eliminate magic numbers)
│   ├── Extract functions (reduce cyclomatic complexity)
│   ├── Improve naming, remove duplication
│   ├── Push toward pure functions where feasible
│   ├── Apply Boy Scout Rule: leave code cleaner than you found it
│   ├── EXECUTE tests after EACH refactoring step → must STILL PASS
│   │   ├── ✅ Still passing → refactoring is safe, continue
│   │   └── ❌ Failed → REVERT that refactoring step, try smaller
│   └── GATE: Tests green after EVERY refactoring change
│
├── 6. Mark task complete [x]
└── 7. Note any deviations or issues discovered
```

<!-- matecito-ai: este módulo REEMPLAZA el Step 4 de la skill, así que la regla de desvíos no le
     llegaba: en modo Strict TDD la conducta vieja ("anotá y seguí") era la única instrucción presente. -->
**Deviations rule (applies inside this cycle, same as in Standard Mode).** Step 7 is for recording
what you resolved, NOT a licence to resolve it. If you find the design wrong or incomplete **in
something the task you are about to write needs**, STOP and return with the gap and the concrete
options (`blocked` or `partial`, and persist first — see "Stopping mid-batch" below) — before writing
the RED test, because the test encodes the decision. Noting it
afterwards is not authorization. If the gap does not affect what you are writing, note it and
continue. The cut is between **execution detail** (an internal name, guard ordering, how you split a
private helper) and **decision** (a contract, a new dependency, which layer the logic lives in,
changing the design's approach); the canonical criterion is in `~/.claude/references/edr/README.md`.
Each deviation you do record must state whether `sdd-verify` will check it against the design.

<!-- matecito-ai: mismo agujero que la regla de desvíos — este módulo reemplaza el Step 4, pero
     marcar tareas y persistir apply-progress son los Steps 5-6 de la skill, POSTERIORES al ciclo.
     En Strict TDD el corte es aún más temprano (antes del test RED), así que sin esta regla un stop
     en la tarea 3 de 5 dejaba el código de las tareas 1-2 escrito en el repo y sin registrar. -->
**Stopping mid-batch — completed work is ALWAYS persisted (same as in Standard Mode).** Every early
exit from this cycle is a stop **inside** the batch, not an escape from it: the design gap above, a
pre-existing failure in step 0 SAFETY NET, an infrastructure failure of the test runner, or any other
unexpected blocker. Before returning control you ALWAYS, in this order:

1. Mark every task whose cycle you completed as `[x]` — Step 5 of `SKILL.md` (this module does NOT
   replace it).
2. Persist `apply-progress` with the REAL state — tasks done, task stopped on and at which stage,
   tasks untouched — Step 6 of `SKILL.md`, MANDATORY on every exit path.
3. Only then return.

**Never return control with completed work unpersisted.** The code and tests of the finished tasks
are already in the repo; leaving them unrecorded makes the next batch re-implement them.

Which status to return — same rule as Standard Mode, resolved top down (`SKILL.md` → "Stopping
Mid-Batch"; the full blocks are in `~/.claude/references/phase-returns/sdd-apply.md`):

- **`blocked`** — the blocker also stops the rest of the batch: you cannot keep going. An
  infrastructure failure of the test runner is always this case: without a runner no further task
  can pass the GREEN execution gate.
- **`partial`** — the phase is not finished: tasks of this change remain. The ordinary continuation
  batch and a stop that left one task untouched are both `partial`; it claims nothing about whether
  a blocker exists.
- **`done`** — nothing remains.

The blocker, when there is one, is reported in ONE place: the return template's `### Blocker`
section, with the gap and the concrete options — not `### Issues Found`, not `risks`. Carry the
envelope per **Section D** of `~/.claude/skills/_shared/sdd-phase-common.md`.

The **TDD Cycle Evidence** table still ships on an early exit: one row per task whose cycle you
completed, plus a row for the task you stopped on showing the last stage it reached (e.g. RED written,
GREEN never run). Do not omit the table because the batch was cut short.

## Choosing Test Layer

Based on the testing capabilities cached in Engram (`sdd/{project}/testing-capabilities`), choose the appropriate test layer for each task:

```
Determine test layer by WHAT the task does:
├── Pure logic, utility function, calculation, data transformation
│   └── Unit test (always available if test runner exists)
│
├── Component rendering, user interaction, state changes
│   ├── IF integration tools available → Integration test
│   └── IF NOT → Unit test with mocks (degrade gracefully)
│
├── Multi-component flow, API interaction, context/provider behavior
│   ├── IF integration tools available → Integration test
│   └── IF NOT → Unit test with mocks
│
├── Critical business flow, full user journey, cross-page navigation
│   ├── IF E2E tools available → E2E test
│   ├── IF NOT but integration available → Integration test
│   └── IF neither → Unit test (degrade gracefully)
│
└── Default: Unit test (always the fallback)
```

**Key rule**: Use the HIGHEST available layer that fits the task. But NEVER skip a task because a layer is unavailable — degrade to the next available layer.

## Test Execution

Detect the test runner from the cached testing capabilities:

```
Read test command from:
├── Cached capabilities → `test_runner.command` (fastest — already detected)
│   └── literal key, under `### Test Runner` in `sdd/{project}/testing-capabilities`
├── project test config → test_command (override)
└── Fallback: detect from package.json/pyproject.toml/go.mod

When executing tests during TDD:
├── Run ONLY the relevant test file, not the entire suite
│   ├── JS/TS: {runner} {test-file-path} (e.g., pnpm vitest run src/utils/tax.test.ts)
│   ├── Python: pytest {test-file-path}
│   ├── Go: go test ./{package}/... -run {TestName}
│   └── Adapt to the runner's CLI
├── This keeps the cycle FAST
└── Full suite runs happen in sdd-verify, not here
```

## Pure Function Preference

When writing production code in GREEN/TRIANGULATE steps, prefer pure functions:

```
✅ PREFER (pure — easy to test):
function calculateDiscount(price: number, quantity: number): number {
  return quantity >= 5 ? price * quantity * 0.1 : 0
}

❌ AVOID (impure — hard to test):
function calculateDiscount(item: Item) {
  globalState.lastDiscount = item.price * 0.1  // side effect
  updateDOM()                                   // side effect
  return globalState.lastDiscount
}
```

**Why**: Pure functions are deterministic (same input → same output), have no side effects, and are trivially testable. TDD naturally pushes you toward pure functions — embrace it.

## Approval Testing (for refactoring existing code)

When a task involves REFACTORING existing code (not writing new code):

```
BEFORE touching production code:
├── 1. Identify existing behavior to preserve
├── 2. Write "approval tests" that capture current behavior:
│   ├── Call the function with known inputs
│   ├── Assert the CURRENT outputs (even if ugly or wrong)
│   └── These tests document what the code does NOW
├── 3. Run approval tests → must PASS (they describe current reality)
├── 4. NOW refactor the production code
├── 5. Run approval tests again → must STILL PASS
│   ├── ✅ Passing → refactoring preserved behavior
│   └── ❌ Failing → refactoring broke something, revert
└── 6. If the spec says behavior should CHANGE:
    ├── Update the approval test to reflect NEW expected behavior
    ├── Run → test FAILS (RED — new behavior not implemented yet)
    └── Implement new behavior → GREEN
```

## Return Summary Extension

<!-- matecito-ai: este módulo reemplaza el Step 4 de la skill, no el formato del retorno. La forma de
     las dos secciones vive en el template; acá sólo lo que significan sus celdas. -->

This module does NOT define its own return format. The return is the one in
`~/.claude/references/phase-returns/sdd-apply.md`, followed literally, for whichever of `done`,
`partial` or `blocked` you resolved. What Strict TDD Mode adds is that its two **conditional**
sections — `### TDD Cycle Evidence` and `### Test Summary` — become mandatory: in this mode they are
part of the return, and their absence is a broken return, not "nothing to report". Emit the block's
`**Mode**: Strict TDD` header line so the orchestrator can tell the condition holds.

Fill those two sections as the template shows. What each column means is this module's business:

**Column definitions**:
- **Safety Net**: Pre-existing tests run before modifying files. "N/A (new)" for new files.
- **RED**: Test written first, referencing code that doesn't exist yet. Always "✅ Written".
- **GREEN**: Tests executed and passing after minimal implementation. Must show execution result.
- **TRIANGULATE**: Additional test cases added to force real logic. "➖ Single" if spec has only one scenario.
- **REFACTOR**: Code improved with tests still passing. "➖ None needed" if code was already clean.

## Assertion Quality Rules (MANDATORY)

**Every assertion must verify REAL behavior.** A test that passes without exercising production logic is worse than no test — it gives false confidence.

### Banned Assertion Patterns (NEVER write these)

```
# TRIVIAL ASSERTIONS — test proves nothing
expect(true).toBe(true)              # ❌ Tautology
expect(false).toBe(false)            # ❌ Tautology
expect(1).toBe(1)                    # ❌ Tautology — no production code involved
assert True                          # ❌ Always passes
assert 1 == 1                        # ❌ Always passes

# EMPTY COLLECTION ASSERTIONS without setup context
expect(result).toEqual([])           # ❌ ONLY valid if you set up conditions for empty
expect(result).toHaveLength(0)       # ❌ Same — why is it empty? Did production code run?
assert len(result) == 0              # ❌ Same — prove the emptiness comes from real logic
assert result == []                  # ❌ Same

# TYPE-ONLY ASSERTIONS — proves existence, not behavior
expect(result).toBeDefined()         # ❌ Alone is useless — WHAT is the value?
expect(result).not.toBeNull()        # ❌ Alone is useless — assert the actual value
expect(typeof result).toBe('object') # ❌ Alone is useless — what does the object contain?
assert result is not None            # ❌ Alone — assert what result actually IS

# GHOST LOOP — assertion inside a loop that iterates 0 times
const items = screen.queryAllByTestId("item");  // returns []
for (const item of items) {
  expect(item).toHaveTextContent("value");       # ❌ NEVER EXECUTES — loop body is dead code
}
# FIX: assert the collection is non-empty FIRST, or set up data so it IS non-empty:
expect(items).toHaveLength(3);                   # ✅ Proves items exist
for (const item of items) { ... }                # ✅ Now the loop actually runs

# INCOMPLETE TDD CYCLE — GREEN without TRIANGULATE
# If your GREEN test passes because the setup doesn't exercise the code path,
# you are NOT done. You MUST triangulate with a setup that DOES exercise it.
# Example: testing "search doesn't update until Enter" but the component
# that receives the search is never rendered → the test proves nothing.
# FIX: add a test where the component IS rendered and verify the behavior.
```

### What Makes a REAL Assertion

Every test assertion must satisfy ALL of these:
1. **Calls production code** — the test invokes a function, method, or component from the implementation
2. **Asserts a specific output** — compares against a concrete expected value derived from the spec
3. **Would FAIL if the production code were wrong** — if you change the implementation logic, THIS test breaks

```
# ✅ REAL assertions — production code determines the result
expect(calculateDiscount(100, 10)).toBe(10)       # Real input → real output
expect(screen.getByText('Welcome, John')).toBeInTheDocument()  # Rendered from data
assert result[0].status == "FAIL"                  # Specific finding from check execution
assert response.status_code == 403                 # Real HTTP response from the endpoint
expect(result).toHaveLength(3)                     # AND you set up exactly 3 items
```

### Empty Collection Rule

`expect(result).toEqual([])` or `assert len(result) == 0` is ONLY valid when:
1. You set up a specific precondition that SHOULD produce an empty result (e.g., no matching records)
2. The production code actually ran and filtered/processed data to arrive at empty
3. A companion test with different setup produces a NON-EMPTY result (triangulation)

If you cannot explain WHY the result is empty based on setup → the assertion is trivial.

### Smoke Test Rule

A test that only renders a component without asserting any output is NOT a valid test:

```
# ❌ SMOKE TEST ONLY — proves nothing about behavior
render(<MyComponent data={mockData} />);
expect(screen.getByTestId("wrapper")).toBeInTheDocument();  # Just proves it rendered

# ✅ BEHAVIORAL TEST — proves what the component DOES with the data
render(<MyComponent data={mockData} />);
expect(screen.getByText("Expected Title")).toBeInTheDocument();  # Verifies output from data
expect(screen.getByRole("button")).toHaveTextContent("Submit");  # Verifies real content
```

"Renders without crash" is a smoke test. It is NOT a unit test, NOT an integration test, and it does NOT count toward TDD coverage. If you need a smoke test, it must be accompanied by real behavioral assertions.

### Mock Hygiene Rules

**If you need more mocks than assertions, you are testing at the WRONG level.**

```
Mock/assertion ratio guide:
├── ≤ 3 mocks for a test file → ✅ Healthy — focused test
├── 4–6 mocks → ⚠️ Consider extracting logic to a pure function
├── 7+ mocks → ❌ STOP — you are testing at the wrong layer
│   ├── Extract the logic under test to a PURE FUNCTION and test it without mocks
│   ├── OR move the test to integration/E2E layer where real dependencies exist
│   └── NEVER write 10+ mocks to verify a one-line transformation
```

**Extract-Before-Mock Rule**: If the behavior you want to test is a data transformation, mapping, filtering, or conditional logic (e.g., `MUTED → FAIL` status conversion), EXTRACT it to a pure function FIRST, then test the pure function directly. No mocks needed.

```
# ❌ BAD: 15 mocks to test a one-line status conversion
vi.mock("next/navigation", ...);
vi.mock("next/link", ...);
vi.mock("@/components/shadcn", ...);
// ... 12 more mocks ...
render(<StatusCell row={mutedRow} />);
expect(screen.getByText("FAIL")).toBeInTheDocument();

# ✅ GOOD: extract and test the logic directly
// In production code:
export function resolveDisplayStatus(status: string, isMuted: boolean): string {
  return status === "MUTED" ? "FAIL" : status;
}

// In test — ZERO mocks needed:
expect(resolveDisplayStatus("MUTED", true)).toBe("FAIL");
expect(resolveDisplayStatus("PASS", false)).toBe("PASS");
```

### Implementation Detail Coupling Rule

Tests must assert **behavior visible to the user**, not internal implementation details:

```
# ❌ COUPLED TO IMPLEMENTATION — breaks on any style refactor
expect(element.className).toContain("text-xs");
expect(element.className).toContain("-mt-2.5");
expect(element.className).toContain("border-border-error-primary");
expect(element.style.color).toBe("red");

# ❌ COUPLED TO INTERNALS — breaks when implementation changes
expect(mockService.mock.calls.length).toBe(3);  # Why 3? Brittle.
expect(component.state.isLoading).toBe(true);    # Internal state, not behavior.

# ✅ BEHAVIORAL — survives refactors, tests what users see
expect(screen.getByText("Error: Payment failed")).toBeInTheDocument();
expect(screen.getByRole("alert")).toHaveTextContent("Risk:");
expect(screen.getByRole("button")).toBeDisabled();
```

**CSS class assertions are NEVER valid test assertions.** If you need to verify visual styling:
1. Test the **semantic outcome** (e.g., element has `role="alert"`, text is visible, button is disabled)
2. OR use a visual regression tool / E2E screenshot comparison
3. NEVER assert specific Tailwind/CSS class names — they are implementation details

## Rules (Strict TDD specific)

- NEVER write production code before writing its test — this is the ONE rule that cannot be broken
- NEVER skip the GREEN execution gate — you MUST run tests and confirm they pass
- NEVER skip triangulation when the spec defines multiple scenarios — hardcoded Fake It must be forced out
- NEVER write trivial assertions (see Banned Assertion Patterns above) — they are WORSE than no test
- ALWAYS verify that every assertion CALLS production code and asserts a SPECIFIC expected value
- ALWAYS run the Safety Net before modifying existing files — protect what already works
- ALWAYS report the TDD Cycle Evidence table — the verify phase will check it
- ALWAYS persist completed work before returning control — a stop mid-cycle (`blocked` or `partial`) still marks the finished tasks and saves `apply-progress` first (see "Stopping mid-batch" above)
<!-- matecito-ai: esta regla decía "reportá Blocked y seguí con la tarea siguiente", en contradicción
     directa con "Stopping mid-batch" (que trata el fallo de infraestructura como salida temprana que
     corta el batch). Vale la de arriba: sin runner no hay compuerta GREEN, así que "seguir" sólo
     puede producir código sin test verde — exactamente lo que este módulo existe para impedir. -->
- If a test runner execution fails for infrastructure reasons (not test failures), that is an early exit from the cycle — do NOT continue to the next task. Mark what you finished, persist `apply-progress`, and return `blocked` with the failure stated in `### Blocker` (see "Stopping mid-batch" above). Without a working runner no further task can pass the GREEN execution gate, so continuing can only produce production code that was never proven green
- Prefer pure functions — but don't force it where it doesn't fit (e.g., React components with state)
- For refactoring tasks, ALWAYS write approval tests before touching code
- Run ONLY the relevant test file during the cycle, not the full suite
