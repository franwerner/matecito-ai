---
name: context7
description: Defines when and how to consult the context7 MCP for current library documentation — version resolution before installing or pinning ANY dependency, library/framework configuration, API usage, version migrations, and library-specific debugging. USE THIS SKILL whenever you are about to add or install a dependency, write a version into a manifest (package.json, go.mod, requirements.txt, Cargo.toml, …), write configuration or API calls for a library/framework, migrate a library across versions, or debug an error specific to a library. Your training data is stale for versions and APIs — resolve them, never recall them.
---

# context7 — current library docs before you act

## The problem this solves

Training data lags reality. A dependency installed "from memory" lands one or two major versions behind; an API written "from memory" may target a shape that no longer exists. context7 is the MCP that serves **current** documentation — versions, APIs, configuration — straight from the source.

**Core rule: never pick a version or an API shape from training memory when context7 can answer it.** Resolving first costs one call; an outdated pin costs a migration.

## Tool resolution (no hardcoded names)

This skill deliberately names **capabilities, not literal tool names** — tool names change between server versions and registration styles. The context7 MCP exposes two capabilities under its registered prefix (typically `mcp__context7__*`):

1. **Library resolution** — turns a library name into a context7 library ID (and surfaces known versions).
2. **Docs query** — answers a focused question against that library's current docs.

Before first use in a session, confirm the actual registered tool names under the prefix (from the session's tool list, or via ToolSearch). Use whatever names are registered — if they differ from what you remember, the registration wins.

**Absent MCP → silent degradation.** If no context7 prefix is registered, skip this skill silently: no error, no mention. Fall back to the package manager's own resolution (install without hand-writing a version — `npm install pkg` / `go get pkg@latest` resolve current versions natively) rather than writing a remembered version number.

## Hard triggers — consult BEFORE acting

| About to… | Consult context7 for |
| --- | --- |
| Add/install a dependency | Current stable (or LTS) version; deprecation status |
| Write a version into a manifest | The version you are about to pin — never from memory |
| Write config for a library/framework | Current config shape (files, keys, defaults) |
| Use a library API non-trivially | Current signature/usage — especially if the lib moves fast |
| Migrate a library across versions | Breaking changes, migration guide |
| Debug a library-specific error | Known issues, changed behavior between versions |

The first two are **mandatory**: no dependency is added and no version is pinned without resolving it first (or falling back to native package-manager resolution when the MCP is absent).

## Version choice rules

- **New dependency** → current stable; prefer LTS when the project targets long-term support runtimes.
- **Existing dependency** → the repo's manifest/lockfile wins for consistency. Do NOT silently upgrade an already-pinned version — an upgrade is a scope change: surface it and ask.
- **Peer/engine constraints** → check the manifest's engine/runtime pins (e.g. `engines.node`) before picking a version the runtime can't run.

## Not found → report, NEVER guess

context7 is the ONLY version source — there is no secondary lookup (no registry queries, no "let the package manager decide", no docs sites). If context7 returns nothing relevant for the requested library, or resolves it but yields no concrete version:

- **NEVER** write a version inferred from training memory, from a "similar" library, or from what "sounds right". A guessed version is worse than no version.
- **NEVER** substitute another source on your own initiative.
- **STOP and report it explicitly**: state that nothing related to that specific library was found in context7, show what you queried, and ask the user how to proceed.
- This is blocking: no manifest edit happens for that dependency until the user answers. The rest of the task may continue if independent.

## What context7 is NOT for

Business logic, refactors, general programming concepts, code review, debugging your own code. Consulting it there adds latency without signal. One focused query per real question — don't spray.

## Scope

Applies wherever code or manifests are written: `sdd-apply`, the `direct` lane, and any ad-hoc edit outside the flow. Headless phase agents with context7 access apply it the same way.

## Self-check (before touching a manifest or library API)

1. Am I about to write a **version number**? → resolved via context7? If not → stop.
1b. Did context7 come up empty for this library or its version? → report "nothing found for <lib> in context7" and wait — never guess, never switch to another source.
2. Am I writing **config/API code** for a library that moves fast? → queried current docs?
3. Is the MCP absent? → degrade silently; never hand-write a remembered version.
4. Is this business logic or a general concept? → context7 does not apply.
