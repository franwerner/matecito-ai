---
name: skill-discovery
version: 1.0.0
description: Use ONLY on-demand, when the user explicitly asks to find or discover skills for the current project (e.g. "find skills for this repo", "search skillsmp", "what skills exist for X"). NEVER trigger automatically or as part of any flow phase. Drives the official `skillsmp` MCP (`mcp__skillsmp__*`) to search the SkillsMP catalog, ranks results by relevance to the current project, dedups against already-available/installed skills, and proposes candidates at a gate. Installation is always MANUAL — this skill guides the user to the skill's GitHub repo, it never writes files or installs anything itself.
license: MIT
metadata: {"hermes":{"tags":["skillsmp","skills","discovery","catalog"],"category":"cross-domain","related_skills":[]},"author":"matecito-ai","version":"1.0.0"}
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

These tools are **read-only**. There is no scan, install, or uninstall tool on the MCP side — search/browse is all the MCP does. Everything after that (ranking, dedup, gating, and the actual install) is this skill's responsibility, performed by a human following manual guidance — never automated.

## Step 1 — Project context (cheap, before searching)

Prefer a project init/context artifact in Engram if one exists — matecito domains store it under `<domain>-init/{project}` — for stack/structure signals that bias ranking (e.g. a Go repo ranks Go-specific skills higher; a visual/design-oriented project ranks design skills higher).

If no init artifact exists, do a **lightweight self-detect only**: read top-level manifest/lockfiles (`go.mod`, `package.json`, `pyproject.toml`, `Cargo.toml`, etc.) and top-level directory names. This is a cheap hint, not a repo scan.

**Some projects keep richer structured context than the manifest reveals** — recorded technology/architecture decisions the project maintains. A manifest can be nearly empty (e.g. a scaffolded frontend) while such records already capture the real framework/state-management/testing choices. When the project has them, read them alongside (not instead of) the manifest/self-detect signals before ranking.

## Step 2 — Search and rank

1. **Search.** Call `search_skills` for the user's topic; use `list_categories` first if the topic is broad and category-shaped, then narrow with `search_skills`. Use `get_skill` to pull full detail for any promising candidate before ranking or presenting it. If a query/category returns no candidate that genuinely fits the project's context, do NOT force a weak match — report it explicitly in the summary ("searched X, 0 relevant fit") instead of padding the list.
2. **Rank** the results by relevance to the user's stated topic and the Step 1 project context. **`stars` is a repo-host metric, not a skill-quality signal** — the catalog's `stars` field reflects the GitHub repo hosting the skill, which is very often a large monorepo/aggregator with many unrelated skills vendored in it. Never use `stars` as a primary ranking factor; at most use it as a light tiebreaker between two otherwise-equal candidates.
3. **Dedup against already-available skills.** Drop or mark any candidate whose skill name already matches an entry in `<available_skills>` (this session's loaded skill list) or an existing `SKILL.md` under `~/.claude/skills/` (global) or `./.claude/skills/` (project).
4. **Dedup candidates against each other (inter-repo clones).** The same skill is commonly re-vendored across dozens of repos (monorepo aggregators, forks, copy-pasted skill packs). Group raw candidates by `name` plus description similarity, and collapse each group into a single listed entry. Within a group, pick the **canonical** repo to present using this precedence: (1) the skill author's own/original repo (name/org matches the skill's stated author, or the repo is dedicated to that skill specifically) over (2) a generic monorepo/aggregator that merely vendors many unrelated skills; if still tied, prefer the repo with the clearer/more complete `get_skill` detail. Never list the same skill N times just because it appears in N repos.

## Step 3 — Gate (mandatory, not skippable)

Present the ranked, deduped candidates:

- **Summary first**: "N found / M already available (deduped) / K duplicate repos collapsed".
- **Grouped by category** (from `list_categories`/`get_skill`, or a heuristic grouping if a catalog result has no category field).
- **Bulk actions**: accept-all, accept per-category, accept/skip per-item.
- **No silent caps.** If a search response was truncated by a result limit, state the exact cap and how many more results likely exist; offer to fetch more rather than quietly presenting a partial list as complete.
- **No silent zero-fit.** If a searched query/category came back with no candidate that genuinely fits the project, say so explicitly in the summary (e.g. "searched WebSocket libraries — 0 relevant fit") instead of omitting it or forcing a loosely-related candidate to fill the slot.
- **Nothing is installed from this gate or ever by this skill** — this skill only proposes; the user installs manually (Step 4). This holds even in an automatic/non-interactive session.

## Step 4 — Manual install guidance (per confirmed skill)

For each skill the user wants, this skill NEVER writes files itself. Instead it gives the user:

1. The skill's **GitHub repository URL** (from the catalog entry).
2. How to copy it into the current project: clone/download the skill's folder and place it at `./.claude/skills/<skill-name>/` (or `~/.claude/skills/<skill-name>/` if the user explicitly wants it global).
3. A clear reminder: **review the skill's code before installing it** — there is no automatic security scan; the catalog only returns metadata, not a vetted-safe guarantee.

## Rate limit / optional key (honest wording — do not overstate)

The anonymous tier is **50 requests/day, 10/minute** — enough for occasional on-demand discovery. If the user wants a higher limit, they can register a free key from **skillsmp.com/docs/api** and re-register the MCP themselves:

```
claude mcp add --transport http --header "Authorization: Bearer sk_live_..." skillsmp https://skillsmp.com/mcp
```

Present this as **optional — it MAY raise the limit**. Do not promise a specific higher number (e.g. 500/day): SkillsMP does not document the remote MCP's auth behavior, so the actual effect of the header is unconfirmed. This skill never applies or asks for a key itself.

## Anti-patterns

- Installing anything automatically — this skill only proposes; the human always installs manually.
- Skipping the Step 3 gate or bulk-accepting without the user's explicit confirmation.
- Silently truncating or capping the candidate list without saying so.
- Listing the same skill multiple times because it's vendored in multiple repos — collapse to the canonical repo (see Step 2's inter-repo dedup).
- Using `stars` as a quality signal — it measures the host repo, not the skill.
- Forcing a weak/loosely-related candidate into a category just to avoid reporting "0 fit".
- Promising a specific rate limit increase from the optional key — say "may raise" only.
- Assuming the MCP exposes scan/install/uninstall tools — it does not; it is search/browse only.
