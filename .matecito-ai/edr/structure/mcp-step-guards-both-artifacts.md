# EDR — The qmd install step's Check guards both the binary and the registration, not the registration alone

- **Status:** Accepted
- **Date:** 2026-09-01

## Contexto
Every existing MCP step (context7, codegraph, drawio, debugger, skillsmp, figma, canva) registers a server the host re-fetches per launch via `npx -y <pkg>@latest`, so its only durable artifact is the registration and `Check: !mcp.Find(name)` covers it completely. qmd distributes as a globally installed npm package instead, so it has two durable artifacts: the binary and the registration. codegraph has the same two-artifact shape and solves it with a separate guardian, `codegraphBinaryStep` (install.go:463), contributed through the manifest's `binaries` array. qmd was ratified out of that array: `sync.Detect` resolves binaries through hand-written per-binary blocks (sync.go:228-280), and adding qmd there would need a fourth block plus a `fetchLatestQmd`, which is out of scope here.

## Decisión
qmdMCPStep's Check reports pending when the registration is absent, OR when the qmd executable is not on PATH. Run reinstalls the binary unconditionally and registers with the host only when the registration is still absent, so a repair run does not attempt a duplicate `claude mcp add`.

## Reglas verificables
- **[auto]** qmdMCPStep's Check returns true when the registration is present but the qmd binary is missing from PATH — `TestQmdMCPStep_Check_RegistrationPresentBinaryAbsent` (`internal/setup/install/install_qmd_test.go`).
- **[auto]** qmdMCPStep's Check returns false only when both the registration and the binary are present — `TestQmdMCPStep_Check_BothPresent` (`internal/setup/install/install_qmd_test.go`).
- **[auto]** a repair run (binary missing, registration present) reinstalls the binary and does not call `claude mcp add` again — `TestQmdMCPStep_Run_AlreadyRegistered_SkipsClaudeMCPAdd` (`internal/setup/install/install_qmd_test.go`).

## Alternativas consideradas
(a) Registration only — `!mcp.Find("qmd")`, matching all six existing MCP steps verbatim. Discarded: leaves a machine that loses the binary out of band reporting nothing pending, forever — the host would launch a command that is not there, and no install or update run would ever repair it. (b) Declaring qmd in the manifest's `binaries` array so `codegraphBinaryStep`'s pattern gives the binary its own guardian — the structurally cleanest answer, and exactly how codegraph solves it, but ratified out of scope: `sync.Detect` has no generic binary path today.

## Consecuencias
The step's own Check is now the qmd binary's only guardian. If this step's registration in mcpRegistry is ever removed while qmd remains installed via npm, that guardian disappears with it — unlike codegraph, qmd has no separate binaries-array watchdog.
