# ui-scenarios Block Schema

The `ui-scenarios` block is a YAML structure authored inside the spec artifact. sdd-verify executes it deterministically to drive a browser and evaluate assertions.

## Per-Scenario Fields

| Field   | Type   | Required | Description |
|---------|--------|----------|-------------|
| `name`  | string | yes      | Human-readable scenario identifier; maps 1:1 to a STATE-verdict row in the verify-report |
| `url`   | string | yes      | URL path appended to the dev-server base URL (e.g. `/login`) |
| `steps` | list   | yes      | Ordered action primitives executed sequentially by the browser agent |
| `expect`| object | yes      | Pass conditions; evaluated after steps complete |

## Step Primitives

Each step maps directly to an agent-browser command:

| Primitive     | Syntax                                      | agent-browser command            |
|---------------|---------------------------------------------|----------------------------------|
| `open`        | `open: <url>`                               | `open <url>`                     |
| `snapshot`    | `snapshot`                                  | `snapshot` (accessibility tree)  |
| `fill`        | `fill: { target: <locator>, value: <text> }`| `fill <target> "<text>"`         |
| `click`       | `click: <locator>`                          | `click <target>`                 |
| `screenshot`  | `screenshot: <label>`                       | `screenshot [<label>]`           |
| `wait`        | `wait: { selector: <locator>, timeout_ms: <n> }` | wait until locator resolves or timeout |

### `wait` / Timeout Primitive

Use `wait` to handle async UI: navigation delays, data-loading states, animation completion.

```yaml
wait: { selector: 'role=status name="Loading"', timeout_ms: 3000 }
```

- `selector`: a role+name locator or CSS selector that must become visible/present before proceeding.
- `timeout_ms`: maximum wait in milliseconds; the scenario FAILS if the element does not appear within this window.
- Always place `wait` steps before assertions that depend on async state.

## Target / Locator Rules

Targets (used in `fill`, `click`, `wait.selector`, and `expect.visible`) MUST be authored as one of:

1. **role+name locator** (primary): `role=<role> name="<accessible-name>"` — e.g. `role=button name="Submit"`, `role=textbox name="Email"`
2. **CSS selector** (accepted alternative): standard CSS string — e.g. `#password`, `.submit-btn`

### Forbidden: runtime snapshot refs (`@eN`)

`@e1`, `@e2`, etc. are ephemeral accessibility-tree refs assigned by agent-browser **per snapshot**. They go stale on any DOM change and are NEVER valid in authored scenarios. Static validation in sdd-verify rejects any step target matching `@e\d+` with a CRITICAL failure.

## `expect` Block

```yaml
expect:
  visible: [<locator>, ...]          # STATE assertion — evaluated LIVE after steps
  text_contains: [<string>, ...]     # STATE assertion — evaluated LIVE after steps
  no_console_errors: true            # session-level error gate (deferred to SUMMARY.md)
  no_server_errors: true             # session-level error gate (deferred to SUMMARY.md)
```

### Assertion Classes

| Class                              | Evaluation point        | Source                         |
|------------------------------------|-------------------------|--------------------------------|
| `visible` / `text_contains`        | Live, per-scenario      | agent-browser snapshot after steps |
| `no_console_errors` / `no_server_errors` | Session-level gate | SUMMARY.md aggregates after `proofshot stop` |

STATE assertions (`visible`, `text_contains`) are evaluated against the live snapshot taken immediately after the scenario's last step. They CANNOT be attributed across sessions — they are always per-scenario.

Error-absence assertions (`no_console_errors`, `no_server_errors`) defer to `consoleErrorCount`/`serverErrorCount` scalars in `SUMMARY.md`. These are session-wide aggregates with no per-scenario breakdown; they cannot be attributed to any individual scenario.

## Canonical Example

```yaml
ui-scenarios:
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

## Validation Rules

1. `name` must be a non-empty string — used as the verdict-table key.
2. `url` must begin with `/`.
3. Every step with a `target` field must use role+name or CSS — never `@eN`.
4. `wait.timeout_ms` must be a positive integer.
5. `expect` must contain at least one assertion.
6. `no_console_errors` / `no_server_errors` are boolean; only `true` is meaningful (omit to skip).
