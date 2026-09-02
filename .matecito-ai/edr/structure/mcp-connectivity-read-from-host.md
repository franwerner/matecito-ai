# EDR — A registered integration's connection state is read from the host, not assumed

- **Status:** Accepted
- **Date:** 2026-09-02

## Contexto
Find(needle) resolves presence from ~/.claude.json first (findInJSON) and falls back to `claude mcp list` (findViaCLI) only when the JSON lookup misses. Found carried a single `Connected bool`, and only findViaCLI ever set it — a JSON hit left it at its zero value, false, indistinguishable from "the host was asked and said no". Found.Describe() (mcp.go, then line 88) read that bool directly, so every JSON-sourced hit rendered as unconnected regardless of whether the integration was actually running.

## Decisión
Found.Connected bool becomes Found.Connection ConnectionState, a three-value type (ConnectionUnknown zero value, ConnectionUp, ConnectionDown). Find keeps findInJSON first for presence, untouched — a JSON hit's presence never depends on the CLI listing. On a JSON hit, Find additionally consults the memoized `claude mcp list` output (connectionFromCLI) and sets Connection from it: ConnectionUp when the matching line reports "✓ Connected", ConnectionDown when the line matches by name but does not, ConnectionUnknown when the CLI cannot be consulted at all or the listing has no line for that name. A CLI-sourced hit (findViaCLI) sets Connection directly from the same listing it already parsed to find the entry — never ConnectionUnknown, since the CLI was necessarily consulted to produce that hit. Describe() renders all three states, keyed off Connection rather than off Source.

## Reglas verificables
- **[auto]** a registered integration the host launches successfully is found and reports ConnectionUp — TestFindEnrichesJSONHitConnectionUp (internal/mcp/mcp_test.go)
- **[auto]** a registered integration the host fails to launch is found and reports ConnectionDown, not simply absent — TestFindEnrichesJSONHitConnectionDown (internal/mcp/mcp_test.go)
- **[auto]** an unregistered integration is reported not found, never as a hit with a negative connection state — TestFindUnregisteredNotConfusedWithDown (internal/mcp/mcp_test.go)
- **[auto]** enriching one integration's connectivity does not change whether every other registered or unregistered integration is found — TestFindOtherPresenceUnchangedByConnectivity (internal/mcp/mcp_test.go)
- **[auto]** a JSON hit invokes the CLI runner exactly once, for enrichment, without presence depending on it — TestFindJSONFirst (internal/mcp/mcp_test.go)

## Alternativas consideradas
Reordering Find to try the CLI first was rejected: it would make presence itself depend on parsing `claude mcp list`, which the governing spec's scenario 4 (presence of every other registered integration unchanged) explicitly forbids changing. Keeping Connected as a bool and just populating it correctly was rejected: a bool cannot represent "not consulted" as a state distinct from "consulted and negative", which the spec names as a distinction the reader must be able to make. A separate FindWithConnectivity entry point, leaving Find itself untouched, was rejected because every scenario in the requirement is written about the one query that also determines presence — splitting them would satisfy the letter of a two-function API and miss the requirement.

## Consecuencias
Connection is read in exactly one place today, Found.Describe() (mcp.go:88 before this change); no caller outside internal/mcp touches it — the five doctor checks (context7, engram, drawio, debugger, codegraph) and every MCP install step consume Describe() or discard the struct entirely, so the type change is contained to mcp.go and its tests. A JSON hit now always costs one additional read of the memoized CLI output (a cache hit after the first call in a process, never a second exec) — the enrichment is best-effort and never blocks or fails presence.
