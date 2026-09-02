# EDR — Run fails naming a qmd it did not install when one still shadows it on PATH

- **Status:** Accepted
- **Date:** 2026-09-02

## Contexto
qmdMCPStep.Run installs its own qmd via `npm install -g <tarball>` into the canonical npm bin dir and registers with the host when the registration is absent. Before this change, a qmd already winning on PATH from somewhere else (e.g. a bun-installed copy, or one under a root-owned /usr/local/bin) was never distinguished from the one this step installs — Run would reinstall unconditionally and return nil, reporting success while the host kept launching the foreign, possibly broken, copy.

## Decisión
qmdMCPStep.Run captures `exec.LookPath("qmd")` BEFORE ensureUserNpmPrefix runs — the load-bearing ordering detail, since ensureUserNpmPrefix prepends the npm bin dir to the process PATH (install.go, ensureUserNpmPrefix), so any LookPath performed after it always resolves to this step's own qmd and the check below would never fire. Run then proceeds exactly as before: resolve the tarball, `npm install -g`, verify the binary landed on PATH, register with the host only when the registration is still absent. Only after all of that, if the path captured before the mutation resolved outside the canonical npm bin dir (resolveUserNpmBinDir, computed fresh so it reflects any prefix reconfiguration ensureUserNpmPrefix just performed), Run returns an error naming that exact foreign absolute path. The machine is left maximally repaired: this step's own qmd is installed and registered, and the one fact the person cannot derive themselves (which of possibly several stacked qmd installs wins on PATH) is named for them.

## Reglas verificables
- **[auto]** on a machine with a foreign qmd shadowing the install, Run still installs and registers its own copy, then returns an error naming the foreign absolute path — TestQmdMCPStep_Run_FailsNamingShadowingExecutable (internal/setup/install/install_qmd_test.go)
- **[auto]** the error names the foreign path, never the canonical one — TestQmdMCPStep_Run_FailsNamingShadowingExecutable (internal/setup/install/install_qmd_test.go)
- **[auto]** on a clean machine (no foreign qmd on PATH before the install mutates it) Run succeeds with no error — TestQmdMCPStep_Run_AssetFound_ReachesInstall and TestQmdMCPStep_Check_NotPendingRightAfterRun (internal/setup/install/install_qmd_test.go)
- **[auto]** a repair run over a registration already present and a pre-existing qmd already in the canonical dir still skips a duplicate `claude mcp add` — TestQmdMCPStep_Run_AlreadyRegistered_SkipsClaudeMCPAdd (internal/setup/install/install_qmd_test.go)

## Alternativas consideradas
Reinstalling on top and returning nil was rejected: `npm install -g` writes only into this step's own prefix (e.g. ~/.npm-global/bin) and cannot remove a root-owned /usr/local/bin/qmd, PATH order is unchanged, so the next Check reports pending again and every future run redownloads while reporting success — the exact permanent loop the incident exposed. Warning and returning nil was rejected for the same loop, minus the download: the run still reports success while the integration stays effectively unusable, which is what the governing spec's new requirement forbids in substance even though a warning technically names the problem. Failing *before* installing and registering was rejected: it leaves the machine strictly worse — no good qmd on disk once the person removes the shadowing one, and no registration — while installing first costs one npm run and turns the remediation into a single `rm`.

## Consecuencias
A qmd install that used to always end in `nil` can now end in an error even though the step's own artifacts (binary + registration) are both correctly in place — the surrounding capability already contracts for a step failing without aborting the run (process/install-mcp-from-release-tarball, "Una falla de este paso no aborta la corrida"), so the overall install run continues and reports this one component in error. The error message is the only user-facing surface of this decision: it must keep naming the exact foreign absolute path, since that is the one piece of information the person cannot derive themselves.

## Relacionados
- `relacionado-con` → [mcp-step-guards-both-artifacts.md](mcp-step-guards-both-artifacts.md) — Check's provenance clause es lo que dispara esta decisión de Run — este record fija qué hace Run una vez que dispara.
