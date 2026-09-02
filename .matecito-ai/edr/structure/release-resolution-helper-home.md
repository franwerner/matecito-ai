# EDR — The release-resolution helper lives in the installer package, not in releasedl

- **Status:** Accepted
- **Date:** 2026-09-01

## Contexto
releasedl's stated contract is goreleaser-style per-OS/arch assets plus a checksums.txt, and it errors out when either is missing (releasedl.go:138-164) — both would fire for qmd, whose release publishes neither. Its GitHub plumbing is also not liftable: githubToken() is unexported (releasedl.go:168), so nothing outside the package can reuse it.

## Decisión
The code that asks GitHub for the latest release lives in the installer package (internal/setup/install) as a small private helper, rather than being added as a new mode of the existing releasedl download helper.

## Reglas verificables
- **[auto]** qmdLatestTarballURL lives in internal/setup/install, not in internal/setup/releasedl.
- **[auto]** The helper accepts a variadic apiBaseURL ...string parameter, mirroring releasedl.LatestReleaseWithTimeout's own test seam (releasedl.go:88-96).

## Alternativas consideradas
Adding an assets-only, no-checksum mode to releasedl. Discarded: releasedl's whole point is mandatory integrity verification (goreleaser assets + checksums.txt); widening it to also serve a caller that structurally cannot have any checksum would make that guarantee optional for every caller, not just qmd.

## Consecuencias
The two helpers stay independent: releasedl keeps its mandatory-checksum contract intact for its existing callers (engram, matecito-ai self-update, codegraph binary), and the new qmd helper owns its own, checksum-less, resolution logic with no shared state.
