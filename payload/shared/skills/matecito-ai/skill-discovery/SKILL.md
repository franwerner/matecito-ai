---
name: skill-discovery
version: 1.1.0
description: Use ONLY on-demand, when the user explicitly asks to find or discover skills for the current project (e.g. "find skills for this repo", "search skillsmp", "what skills exist for X"). NEVER trigger automatically or as part of any flow phase. Drives the official `skillsmp` MCP (`mcp__skillsmp__*`) to search the SkillsMP catalog, ranks results, dedups against already-available capability, then INSPECTS each candidate's real `SKILL.md` and audits it for foreign-project contamination, mechanism/version conflicts, trigger collisions and executable content before proposing anything at a gate. Nothing is written anywhere without an explicit per-item verdict.
license: MIT
metadata: {"hermes":{"tags":["skillsmp","skills","discovery","catalog","audit"],"category":"cross-domain","related_skills":[]},"author":"matecito-ai","version":"1.1.0"}
---

# Skill Discovery

**On-demand only — never automatic, never part of any flow phase.** Cross-domain, always available once matecito-ai is installed.

## The `skillsmp` MCP — official, read-only

This skill drives the `skillsmp` MCP, registered by the matecito-ai installer as a **remote HTTP server** at `https://skillsmp.com/mcp`, **without an API key** — the default is the anonymous tier.

This skill knows the MCP **exists**; it does not hardcode what the MCP **exposes** beyond its documented, read-only surface. The concrete tool signatures are owned by the MCP server and are learned at runtime from the loaded `mcp__skillsmp__*` tools — treat their live descriptions as authoritative over this table:

| Tool | Use for |
|------|---------|
| `mcp__skillsmp__search_skills` | Search the catalog by keyword/topic |
| `mcp__skillsmp__get_skill` | Fetch details for one specific skill |
| `mcp__skillsmp__list_categories` | List the catalog's categories |

These tools are **read-only**. There is no scan, install, or uninstall tool on the MCP side — search/browse is all the MCP does. Everything after that (ranking, dedup, inspection, gating, and the actual install) is this skill's responsibility.

## The risk this skill exists to manage

**The catalog is not curated.** It indexes `SKILL.md` files from any public GitHub repo, so the majority of results are skills written **for another project**, encoding that project's paths, business rules, tooling and stack assumptions. The catalog has no vetted tier and no human review; `search_skills` returns metadata only.

**Catalog metadata does not reveal this.** A description is routinely accurate about the topic and completely silent about the contamination — the foreign rules live in the body of the file. A skill can also be redundant with capability the project already has, or contradict the project's own recorded decisions.

Therefore the load-bearing rule of this skill: **never propose a candidate from metadata alone.** Read the real file first (Step 3).

## Step 1 — Project context and redundancy (before searching)

Prefer a project init/context artifact in Engram if one exists — matecito domains store it under `<domain>-init/{project}` — for stack/structure signals that bias ranking (e.g. a Go repo ranks Go-specific skills higher; a visual/design-oriented project ranks design skills higher).

If no init artifact exists, do a **lightweight self-detect only**: read top-level manifest/lockfiles (`go.mod`, `package.json`, `pyproject.toml`, `Cargo.toml`, etc.) and top-level directory names. This is a cheap hint, not a repo scan. In a monorepo, detect **per app** — the stack that matters is the one of the app the skill would serve.

**Some projects keep richer structured context than the manifest reveals** — recorded technology/architecture decisions the project maintains. A manifest can be nearly empty (e.g. a scaffolded frontend) while such records already capture the real framework/state-management/testing choices. When the project has them, read them alongside (not instead of) the manifest/self-detect signals before ranking.

### Redundancy check — three sources, before spending a search

A candidate that duplicates capability the project already has is not a find; it is a conflict risk. Check, in order:

1. **Already-available skills** — `<available_skills>` for this session, plus existing `SKILL.md` files under `~/.claude/skills/` and `./.claude/skills/`.
2. **The project's own recorded decisions and specs** — if the project already fixes an architecture, a skill that "teaches" that architecture is redundant *and* liable to contradict it. A skill must never be allowed to silently override an accepted project decision.
3. **`context7`** — if the topic is essentially **documentation for a library or API**, context7 covers it better: it serves current, version-pinned, field-level reference on demand, and it does not go stale the way a vendored docs index does. **Do not propose a docs-index skill for a library context7 already carries.**

**The guardrail exception to (3).** A skill *does* add what context7 cannot when it encodes a **negative constraint** — e.g. "do NOT emit v3 syntax". context7 answers *when asked*; it does not stop a model from confidently writing stale syntax from training data without looking anything up. An anti-staleness guardrail is a behavioral constraint, not documentation, and is worth proposing. Validate it per Step 3c.

## Step 2 — Search and rank

1. **Search.** Call `search_skills` for the user's topic; use `list_categories` first if the topic is broad and category-shaped, then narrow with `search_skills`. Use `get_skill` to pull full detail for any promising candidate. If a query/category returns no candidate that genuinely fits the project's context, do NOT force a weak match — report it explicitly in the summary ("searched X, 0 relevant fit") instead of padding the list.
2. **Rank** the results by relevance to the user's stated topic and the Step 1 project context. **`stars` is a repo-host metric, not a skill-quality signal** — the catalog's `stars` field reflects the GitHub repo hosting the skill, which is very often a large monorepo/aggregator with many unrelated skills vendored in it. Never use `stars` as a primary ranking factor; at most use it as a light tiebreaker between two otherwise-equal candidates.
3. **Dedup against already-available skills** — per Step 1's redundancy check. Drop or mark anything already covered.
4. **Dedup candidates against each other (inter-repo clones).** The same skill is commonly re-vendored across dozens of repos (monorepo aggregators, forks, copy-pasted skill packs). Group raw candidates by `name` plus description similarity, and collapse each group into a single listed entry. Within a group, pick the **canonical** repo using this precedence: (1) the skill author's own/original repo (name/org matches the skill's stated author, or the repo is dedicated to that skill specifically) over (2) a generic monorepo/aggregator that merely vendors many unrelated skills; if still tied, prefer the repo with the clearer/more complete `get_skill` detail. Never list the same skill N times just because it appears in N repos.

## Step 3 — Inspect the real content (MANDATORY before the gate)

Metadata is insufficient to judge a candidate. For every candidate that survives Step 2:

1. **Clone shallow into a scratch directory OUTSIDE the project.** Never clone into the repo, never leave the clone behind, never `git add` it. **Cloning for inspection is not installing** — nothing enters the project until Step 5.
2. **Enumerate every file in the skill's folder**, then read its `SKILL.md` in full — frontmatter *and* body.

Then run the five audits below and record the findings **per candidate**. A candidate whose audits were not run is never presented as clean.

### 3a — Contamination audit (foreign-project residue)

Scan the body for anything that belongs to the origin project rather than to the topic:

- **Paths that do not exist in this repo** — a hardcoded helper file, a directory layout from another stack.
- **Business rules with proper nouns** — named approval gates, flags, or policies from the origin company/product. These are the most dangerous: they read as authoritative project rules and are entirely fictional here.
- **Tools or MCPs unavailable in this environment** — instructions to call a fetch/search tool this setup does not have.
- **Named env vars, tables, services, or endpoints** of the origin system.
- **References to sibling skills you do not have** — harmless as soft pointers, but note them.
- **Stack assumptions incompatible with the project** — a framework, ORM, or directory convention other than the project's.

Report each hit with its line. A candidate can still be worth installing **cleaned** — say exactly what would be removed before anything is written.

### 3b — Mechanism contradiction

Read the skill's own **"NOT for" / "Do NOT use when" / anti-examples** section and cross it against the project's actual toolchain. If the skill **excludes the very mechanism the project uses**, it is advisory reference at best — present it that way, never as a workflow guide, and say plainly which part does not apply.

### 3c — Version validation (bidirectional)

Any skill asserting a version constraint must be checked against the version **actually installed in the manifest**. The verdict can invert on the major version: a skill forbidding a framework's previous-generation syntax is a correct guardrail on the new major and **actively harmful** on the old one, where it would forbid the correct syntax. Never accept the skill's own version claim as the project's — read the manifest.

### 3d — Trigger collision

Check the frontmatter `name` and any declared triggers/keywords against installed skills **and** against the project's own flow phase names and domain vocabulary. A generic name can hijack activation away from the project's own process. On collision, **propose a rename at the gate** rather than dropping an otherwise-good candidate.

### 3e — Executable content (an actual review, not a reminder)

"Review the code before installing" is not a check. Perform one:

- Enumerate every file and **flag anything that is not `.md`** — scripts, hooks, `bin/`, package manifests, binaries.
- Read what `allowed-tools` declares. A reference skill has no reason to claim `Write`/`Edit`/`Bash`.
- Scan the body for network calls, install commands, credential/env access, and instructions that would execute something.
- **State the result explicitly per candidate.** "Single `SKILL.md`, no executables, no tool claims" is a finding worth recording — not an assumption to skip.

There is no automatic security scan upstream: the catalog returns metadata, not a safety guarantee.

## Step 4 — Gate (mandatory, not skippable)

Present the surviving candidates **with their Step 3 findings attached**:

- **Summary first**: "N found / M already covered (deduped) / K duplicate repos collapsed / J flagged by audit".
- **Grouped by category** (from `list_categories`/`get_skill`, or a heuristic grouping if a result has no category field).
- **Per-item verdict — not a binary accept/skip.** Offer:
  - **install as-is** — audits clean, fits the stack and version.
  - **install cleaned** — sound core, foreign residue removed; show exactly what will be stripped or adapted *before* writing.
  - **install renamed** — trigger collision; propose the new name.
  - **reference only** — content is sound but the mechanism or version does not match this project; it must not drive workflow.
  - **discard** — contaminated beyond its value, or redundant with existing capability / context7.
- **State the destination per item.** In a monorepo, name the target app — a UI skill does not belong in the API's skills dir.
- **No silent caps.** If a search response was truncated by a result limit, state the exact cap and how many more results likely exist; offer to fetch more rather than presenting a partial list as complete.
- **No silent zero-fit.** If a searched query/category returned nothing that genuinely fits, say so explicitly ("searched WebSocket libraries — 0 relevant fit") instead of omitting it or padding with a loosely-related candidate.
- **Nothing is written anywhere without its own explicit verdict** — this holds even in an automatic/non-interactive session. Bulk-accept is only ever offered for candidates whose audits came back clean.

## Step 5 — Install and record

Only for items with an explicit verdict from Step 4.

1. **Destination.** `./.claude/skills/<skill-name>/` for the project — in a monorepo, the **specific app's** skills dir, never defaulting to the repo root. `~/.claude/skills/<skill-name>/` only if the user explicitly wants it global.
2. **Apply the verdict.** For *cleaned* or *renamed*, make exactly the changes shown at the gate — no additional edits.
3. **Record provenance** alongside the installed skill: source repo URL, the path within that repo, and the commit inspected. Without it the skill cannot be re-audited or updated later, and its origin becomes unknowable.
4. **Remove the scratch clone.**
5. **Note the reload requirement** — skills are picked up at session start, so a newly installed skill is not active until the session restarts.

## Rate limit / optional key (honest wording — do not overstate)

The anonymous tier is **50 requests/day, 10/minute** — enough for occasional on-demand discovery. If the user wants a higher limit, they can register a free key from **skillsmp.com/docs/api** and re-register the MCP themselves:

```
claude mcp add --transport http --header "Authorization: Bearer sk_live_..." skillsmp https://skillsmp.com/mcp
```

Present this as **optional — it MAY raise the limit**. Do not promise a specific higher number (e.g. 500/day): SkillsMP does not document the remote MCP's auth behavior, so the actual effect of the header is unconfirmed. This skill never applies or asks for a key itself.

## Anti-patterns

- **Proposing from catalog metadata alone**, without reading the real `SKILL.md`. The contamination is never in the description.
- Presenting a candidate as clean without having run the Step 3 audits.
- Cloning into the project directory, or leaving a scratch clone behind.
- Installing anything without its own explicit per-item verdict.
- Skipping the Step 4 gate, or bulk-accepting candidates that carry audit findings.
- Installing a **docs-index skill** for a library `context7` already covers.
- Accepting the skill's **version claim** instead of checking the version installed in the manifest.
- Treating a skill whose own "NOT for" excludes the project's mechanism as a workflow guide.
- Dropping a good candidate over a **name collision** instead of proposing a rename.
- Installing into the **wrong app's** skills dir in a monorepo, or defaulting to the repo root.
- Letting a vendored skill override an **accepted project decision** instead of surfacing the conflict.
- Silently truncating or capping the candidate list without saying so.
- Listing the same skill multiple times because it is vendored in multiple repos — collapse to the canonical repo (Step 2.4).
- Using `stars` as a quality signal — it measures the host repo, not the skill.
- Forcing a weak/loosely-related candidate into a category just to avoid reporting "0 fit".
- Promising a specific rate limit increase from the optional key — say "may raise" only.
- Assuming the MCP exposes scan/install/uninstall tools — it does not; it is search/browse only.
