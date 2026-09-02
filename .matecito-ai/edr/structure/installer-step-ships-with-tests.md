# EDR — The new installer code ships with tests, matching what its six existing installer-step siblings already have

- **Status:** Accepted
- **Date:** 2026-09-01

## Contexto
The change's spec Scope table introduces itself as 'the only source sdd-apply has for what to touch' and lists no test file for the qmd installer step, so writing one is outside it on a literal reading. Every existing MCP step in this package has tests (install_skillsmp_test.go, install_debugger_test.go, install_codegraph_test.go, install_proofshot_test.go, install_mcp_test.go), and the qmd step was the only one in the repo with no test precedent at all — precisely the case where shipping untested is worst. Strict TDD is off for this project (config.json strictTdd: false), so no guard forced the answer either way, and the gap also changed the line estimate sdd-tasks forecast.

## Decisión
The qmd MCP installer step ships with tests in the same shape its six siblings use: install_qmd_test.go covers the step's Check semantics (registration absent, registration present but binary missing, both present) and every Run failure path (npm absent, claude absent, binary absent after a successful npm install, the release asset absent, a non-200 GitHub API response), plus the happy path reaching the npm install stage. install_mcp_test.go's registry and permission-pattern tests are extended with a qmd case, mirroring the existing skillsmp coverage.

## Reglas verificables
- **[auto]** install_qmd_test.go exists and covers both Check branches (registration absent; registration present but binary missing) plus every Run failure path (npm absent, claude absent, binary absent after npm success, release asset absent, non-200 API response) — enforced by `go test ./internal/setup/install/...`.
- **[auto]** install_mcp_test.go's TestPermissionPattern_ConventionAndOverride includes a "qmd" case, and a TestMCPRegistry_Qmd asserts the registry entry exists and that defaultMCP excludes "qmd", mirroring TestMCPRegistry_Skillsmp — enforced by `go test ./internal/setup/install/...`.

## Alternativas consideradas
Shipping the qmd installer step with no tests, on the literal reading that the spec's Scope table names no test file and therefore adding one is out of scope. Discarded: every sibling MCP step already has test coverage, Strict TDD being off removes the one guard that would have forced a decision either way, and an untested network-touching install step is exactly the shape most likely to regress silently.

## Consecuencias
Every future installer step that reaches the network or a second durable artifact is expected to ship equivalent Check/Run coverage unless a later change explicitly ratifies skipping it — absorbing that gap silently again reopens the same fork this record closes.
