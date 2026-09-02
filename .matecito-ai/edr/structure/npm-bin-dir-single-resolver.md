# EDR — install.UserNpmBinDir() as the single exported npm-bin-dir resolver

- **Status:** Accepted
- **Date:** 2026-09-02

## Contexto
Both the install surface's codegraph/proofshot detection and the verify surface's equivalent checks need the folder where npm puts global executables, to probe a binary there instead of resolving it again on their own. internal/setup/install already had this logic as an unexported resolveUserNpmBinDir (install.go:495), built for qmdMCPStep's own Check.

## Decisión
Export UserNpmBinDir() (string, error) as a thin wrapper over the existing unexported resolveUserNpmBinDir — resolveUserNpmBinDir itself is neither renamed nor removed, since structure/mcp-step-guards-both-artifacts.md's Consecuencias names it verbatim as the resolver qmdMCPStep's Check and Run share. install.CodegraphBinaryPath and install.ProofshotBinaryPath are both built on top of UserNpmBinDir().

## Reglas verificables
- **[auto]** install.UserNpmBinDir() returns byte-identical output to a direct call to resolveUserNpmBinDir(), which still exists under its current name — exercised indirectly through internal/setup/sync/binary_detection_test.go's real-resolver subtests, which resolve through UserNpmBinDir() and assert the joined canonical path is what gets probed.

## Alternativas consideradas
(a) A new shared package exporting the npm-bin-dir resolver — churn a bug fix does not justify; the design's own first pass already rejected this for the same reason. (b) Renaming resolveUserNpmBinDir to UserNpmBinDir directly instead of wrapping it — would require editing mcp-step-guards-both-artifacts.md's Consecuencias, which names resolveUserNpmBinDir verbatim, turning a mechanical export into an EDR edit unrelated to this change's own scope. (c) Each of CodegraphBinaryPath/ProofshotBinaryPath calling resolveUserNpmBinDir directly (unexported, package-internal) since both already live in internal/setup/install — technically possible, but the exported UserNpmBinDir() is also what the verify-surface accessors and any future in-package caller reach for, so a single named entry point reads clearer than two callers using different visibility for the same fact.

## Consecuencias
internal/setup/install gains one more exported name; resolveUserNpmBinDir keeps its current signature, callers and tests untouched. A future accessor for another npm-installed binary's canonical path is built on UserNpmBinDir() the same way CodegraphBinaryPath and ProofshotBinaryPath are, rather than re-deriving the npm prefix.
