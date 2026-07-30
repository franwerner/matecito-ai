<!-- matecito-ai: this file used to describe ONE block with two audiences — `sdd-spec` authored a fully
     executable YAML (route + concrete locators) and `sdd-verify` ran it. That design was unworkable in the
     default lane and it contradicted a rule the ecosystem already states elsewhere.

     Unworkable: the schema required `url` and role+name/CSS targets, the intake brief carries neither, and
     for a NEW feature the control does not exist yet — so reading the frontend did not help either. The
     spec's only legal move was `blocked`, asking the user for routes and accessible names. Every
     `ui-test: needed` change stopped at spec.

     Contradictory: the capability-spec vocabulary rule forbids volatile internal identifiers in a spec —
     and `/login`, `role=button name="Submit"`, `#password` are exactly that. The WHAT belongs to the spec;
     the literal HOW belongs to the code.

     The split below follows who actually knows what, and when: the spec knows what must be true (the user
     confirmed it); `sdd-apply` knows the concrete locators BECAUSE IT JUST WROTE THEM; `sdd-verify` needs
     exact targets to stay deterministic. Resolving locators semantically at verify time was the other
     candidate and was rejected: it moves the verdict from a reproducible check to an agent's judgment,
     which is what every other guard here works to avoid. Deployed path: `~/.claude/references/ui-scenarios-schema.md`. -->

# ui-scenarios Schema

UI verification is defined in **two halves**, authored by different phases, bound by the scenario `name`.

| Audience | Phase | What it takes from this file |
|----------|-------|------------------------------|
| **Behavior author** | `sdd-spec` | Part 1 — how to write the behavioral scenario: what must be true, in domain language. It writes this only when the intake brief carries `ui-test: needed`. |
| **Executable author** | `sdd-apply` | Part 2 — how to write the executable counterpart of each behavioral scenario: the real route and the real locators of what it just built. |
| **Executor** | `sdd-verify` | Part 3 — how to run the counterparts, how to check they cover every behavioral scenario, and what to do when one is missing. |

Everything below is normative for all three. No phase authors what the next cannot use, and none of them
re-derives what a previous one already fixed.

**Why the halves.** The spec cannot know the route or the accessible name of a control that does not
exist yet, and it must not carry them anyway — they are implementation detail. `sdd-apply` knows both
without guessing, because it authored them. Splitting keeps the spec free of volatile identifiers, stops
it from blocking, and still hands `sdd-verify` exact targets.

**The binding key is `name`.** Not a new coupling: `name` was already the key `sdd-verify` uses to map a
scenario to its row in the verdict table. Both halves of a scenario carry the same `name`, verbatim.

---

## Part 1 — The behavioral scenario (authored by `sdd-spec`)

Lives inside the **spec artifact**, under a `ui-scenarios:` key. It is the visual-surface expression of
the Given/When/Then scenarios the spec already wrote — derived from them, never a second set of
requirements.

| Field  | Type   | Required | Description |
|--------|--------|----------|-------------|
| `name` | string | yes      | Scenario identifier. The binding key: `sdd-apply` reuses it verbatim, `sdd-verify` maps it 1:1 to a STATE-verdict row |
| `given`| string | yes      | The starting situation, in domain language |
| `when` | string | yes      | The interaction the user performs, in domain language |
| `then` | list   | yes      | What must be true afterwards — one item per observable condition |

**Domain language only — no implementation detail.** No routes, no CSS selectors, no accessible names, no
component names. Name what the user sees and does, the way the requirements above already name it. If you
find yourself writing `/login` or `role=button name="Submit"`, you are authoring Part 2, which is not
yours.

```yaml
ui-scenarios:
  - name: login renders and submits
    given: an unauthenticated visitor on the sign-in screen
    when: they submit valid credentials
    then:
      - the app greets them by showing a welcome heading
      - no error message is shown
```

A capability with no visual surface produces no entry. Never introduce behavior here that the
requirements above do not already state.

---

## Part 2 — The executable counterpart (authored by `sdd-apply`)

Lives inside the **`apply-progress` artifact**, under an `### UI Scenario Counterparts` section. One
counterpart per behavioral scenario, **cumulative across batches** like every other section of that
artifact.

| Field   | Type   | Required | Description |
|---------|--------|----------|-------------|
| `name`  | string | yes      | The behavioral scenario's `name`, **verbatim**. A name that does not match one is a counterpart nobody runs |
| `url`   | string | yes      | The real route you implemented, appended to the dev-server base URL (e.g. `/login`) |
| `steps` | list   | yes      | Ordered primitives that reach the state the `when` describes |
| `expect`| object | yes      | The `then` conditions expressed as executable assertions |

You are not designing behavior here: the behavioral scenario fixes what must be true, and you translate
it against the surface you actually built. If you cannot translate a `then` into an assertion, that is a
finding — say so; do not silently weaken it.

### Step primitives

| Primitive     | Syntax                                      | agent-browser command            |
|---------------|---------------------------------------------|----------------------------------|
| `open`        | `open: <url>`                               | `open <url>`                     |
| `snapshot`    | `snapshot`                                  | `snapshot` (accessibility tree)  |
| `fill`        | `fill: { target: <locator>, value: <text> }`| `fill <target> "<text>"`         |
| `click`       | `click: <locator>`                          | `click <target>`                 |
| `screenshot`  | `screenshot: <label>`                       | `screenshot [<label>]`           |
| `wait`        | `wait: { selector: <locator>, timeout_ms: <n> }` | wait until locator resolves or timeout |

Use `wait` for async UI — navigation delays, data loading, animation. `selector` is a locator that must
become present before proceeding; `timeout_ms` is a positive integer and the scenario FAILS if the
element does not appear in that window. Always place `wait` before assertions that depend on async state.

### Target / locator rules

Every target (in `fill`, `click`, `wait.selector`, `expect.visible`) MUST be one of:

1. **role+name locator** (primary): `role=<role> name="<accessible-name>"` — e.g. `role=button name="Submit"`
2. **CSS selector** (accepted alternative): e.g. `#password`, `.submit-btn`

**Forbidden: runtime snapshot refs (`@eN`).** `@e1`, `@e2` are ephemeral accessibility-tree refs that
agent-browser assigns **per snapshot**. They look correct when written and go stale on any DOM change.
`sdd-verify` rejects any target matching `@e\d+` with a CRITICAL failure that fails the scenario. This
prohibition now applies to YOUR output — you are the one authoring targets.

### `expect` block

```yaml
expect:
  visible: [<locator>, ...]          # STATE assertion — evaluated LIVE after steps
  text_contains: [<string>, ...]     # STATE assertion — evaluated LIVE after steps
  no_console_errors: true            # session-level error gate (deferred to SUMMARY.md)
  no_server_errors: true             # session-level error gate (deferred to SUMMARY.md)
```

| Class | Evaluation point | Source |
|-------|------------------|--------|
| `visible` / `text_contains` | Live, per-scenario | agent-browser snapshot after the steps |
| `no_console_errors` / `no_server_errors` | Session-level gate | `SUMMARY.md` aggregates after `proofshot stop` |

STATE assertions are evaluated against the live snapshot taken right after the last step; they are always
per-scenario. Error-absence assertions defer to the `consoleErrorCount` / `serverErrorCount` scalars in
`SUMMARY.md` — session-wide aggregates with **no** per-scenario breakdown, so they are never attributed to
an individual scenario.

### Counterpart of the example above

```yaml
### UI Scenario Counterparts

- name: login renders and submits
  url: /login
  steps:
    - open: /login
    - snapshot
    - fill: { target: 'role=textbox name="Email"', value: "a@b.com" }
    - fill: { target: "#password", value: "secret" }
    - click: 'role=button name="Submit"'
    - wait: { selector: 'role=heading name="Welcome"', timeout_ms: 3000 }
    - snapshot
    - screenshot: post-submit
  expect:
    visible: ['role=heading name="Welcome"']
    text_contains: ["Welcome"]
    no_console_errors: true
    no_server_errors: true
```

---

## Part 3 — Execution and coverage (performed by `sdd-verify`)

Read the behavioral scenarios from the **spec** and the counterparts from **`apply-progress`**, then pair
them by `name`.

**Coverage check — before running anything:**

| Situation | Verdict |
|---|---|
| Every behavioral scenario has a counterpart | Proceed to execute |
| A behavioral scenario has **no** counterpart | That scenario is `UNTESTED`, **CRITICAL**. Apply must not be able to drop one quietly |
| A counterpart's `name` matches no behavioral scenario | **WARNING** — orphan counterpart; report it, do not run it |
| **No** `### UI Scenario Counterparts` section at all, while the spec declares behavioral scenarios | An explicit **CRITICAL** finding: the UI check could not run and the report says so. **Never a silent skip** — silence here is indistinguishable from "no UI work was needed" |

That last row is the whole point of the coverage check. The UI gate closes silently by design when
`ui-test != needed` or `uiTest.available = ❌`; it must NOT close silently when the check was warranted
and the inputs were missing.

**Then execute**, per the ProofShot session lifecycle in this phase's own steps: static validation
(reject `@e\d+`), one session for the whole run, drive each counterpart's steps, snapshot, evaluate its
STATE assertions, stop, read the session aggregates.

`ui-verdict = (all per-scenario STATE assertions PASS) AND (session-level ERROR GATE PASS)`. Any FAIL is
CRITICAL and blocks archive.

---

## Validation rules

Authored by whoever owns the field:

1. `name` is a non-empty string, and identical in both halves — it is the verdict-table key.
2. `given` / `when` are non-empty and carry no route, selector or accessible name (Part 1).
3. `then` has at least one condition (Part 1).
4. `url` begins with `/` (Part 2).
5. Every step with a `target` uses role+name or CSS — never `@eN` (Part 2).
6. `wait.timeout_ms` is a positive integer (Part 2).
7. `expect` contains at least one assertion (Part 2).
8. `no_console_errors` / `no_server_errors` are booleans; only `true` is meaningful (omit to skip).
