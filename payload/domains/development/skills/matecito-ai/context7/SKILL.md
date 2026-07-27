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

**Absent MCP → silent degradation.** If no context7 prefix is registered, skip this skill silently: no error, no mention. Resolve versions through the ecosystem's remote lookup alone (see below), or let the package manager resolve them (`npm install pkg` / `go get pkg@latest`) — never from a remembered version number.

## Hard triggers — consult BEFORE acting

| About to… | Consult context7 for |
| --- | --- |
| Add/install a dependency | Current stable (or LTS) version; deprecation status |
| Write a version into a manifest | The version you are about to pin — never from memory |
| Write config for a library/framework | Current config shape (files, keys, defaults) |
| Use a library API non-trivially | Current signature/usage — especially if the lib moves fast |
| Migrate a library across versions | Breaking changes, migration guide |
| Debug a library-specific error | Known issues, changed behavior between versions |

The first two are **mandatory**: no dependency is added and no version is pinned without resolving it through context7 first and then corroborating it (see below).

## Version choice rules

- **New dependency** → current stable; prefer LTS when the project targets long-term support runtimes.
- **Existing dependency** → the repo's manifest/lockfile wins for consistency. Do NOT silently upgrade an already-pinned version — an upgrade is a scope change: surface it and ask.
- **Peer/engine constraints** → check the manifest's engine/runtime pins (e.g. `engines.node`) before picking a version the runtime can't run.

## Corroborate with the ecosystem's remote lookup

context7's docs are a snapshot and can lag the registry, so a version taken from it alone may already be stale. After context7 gives you a candidate, **corroborate it against the ecosystem's own remote lookup** — then apply the resolution rule below.

**Identify the ecosystem from the manifest present**, and use its lookup command:

| Manifest | Remote lookup |
| --- | --- |
| `package.json` | `npm view <pkg> version` (or the project's package manager equivalent) |
| `go.mod` | `go list -m -versions <module>` |
| `pyproject.toml` / `requirements.txt` | `pip index versions <pkg>` |
| `Cargo.toml` | `cargo info <crate>` |
| `Gemfile` | `gem list -r -e <gem>` |
| `composer.json` | `composer show -a <pkg>` |

**The number NEVER comes from what is installed locally.** Installed state (`node_modules`, `npm ls`, `pip show`, `go.sum`, any lockfile) tells you what IS there, not what is current — it is not an answer to "which version should I pin". Only context7 or the ecosystem's **remote** lookup may produce the number.

**Resolution:**

- **Both agree** → write it.
- **They differ** → write the ecosystem's version **only if it is greater** than context7's; if it is lower or equal, keep context7's. Either way, note the discrepancy in one line.
- **context7 empty, ecosystem resolves it** → write the ecosystem's version and state that context7 did not have the library.
- **No lookup tool available** for that ecosystem (toolchain absent) → context7's version stands; say it was not corroborated.

## Not found in ANY source → report, NEVER guess

If neither context7 nor the ecosystem's remote lookup finds the requested library (or its specific version):

- **NEVER** write a version inferred from training memory, from a "similar" library, from a lockfile, or from what "sounds right". A guessed version is worse than no version.
- **NEVER** invent a third source to fill the gap.
- **STOP and report it explicitly**: state that nothing was found for that specific library, show what you queried in each source, and ask the user how to proceed.
- This is blocking: no manifest edit happens for that dependency until the user answers. The rest of the task may continue if independent.

## What context7 is NOT for

Business logic, refactors, general programming concepts, code review, debugging your own code. Consulting it there adds latency without signal. One focused query per real question — don't spray.

## Scope

Applies wherever code or manifests are written: `sdd-apply`, the `direct` lane, and any ad-hoc edit outside the flow. Headless phase agents with context7 access apply it the same way.

## Self-check (before touching a manifest or library API)

1. Am I about to write a **version number**? → came from context7 or the ecosystem's remote lookup? If from neither → stop.
2. Did I corroborate context7's candidate against the ecosystem's remote lookup, taking the greater of the two?
3. Am I reading the number off installed state (`node_modules`, `npm ls`, `pip show`, `go.sum`, a lockfile)? → that is not a source; go query remotely.
4. Did BOTH sources come up empty? → report what you queried in each and wait — never guess.
5. Am I writing **config/API code** for a library that moves fast? → queried current docs?
6. Is this business logic or a general concept? → context7 does not apply.
