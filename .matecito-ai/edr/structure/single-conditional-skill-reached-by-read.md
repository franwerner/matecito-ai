# EDR — The two new phases reach resolve-library-docs by a directed Read of its deployed path

- **Status:** Accepted
- **Date:** 2026-09-08

## Contexto
sdd-propose and sdd-explore both need to reach the resolve-library-docs skill before writing a library fact (version, API, config option, support status, migration step). Three precedents for reaching a conditional skill already coexist in the repo: the Skill tool plus prose (sdd-apply, and what payload/docs/reference-deployed-paths.md prescribes for conditional skills), a directed Read of the deployed path (sdd-design, stated in its own frontmatter comment), and a skills: preload (development-decisions-mine, for this same skill). Nothing in the repo fixed which of the three a phase should take, or why.

## Decisión
Both phases reach resolve-library-docs via a directed Read of its deployed path (`~/.claude/skills/resolve-library-docs/SKILL.md`), triggered on demand. Neither gains the `Skill` tool, and neither carries a `skills:` preload naming resolve-library-docs. The reach mechanism is identical in both phases.

The criterion this decision fixes, so the next phase that needs to reach a conditional skill does not have to re-derive it: unconditional reach → `skills:` preload; more than one conditional skill needed → the `Skill` tool; exactly one conditional skill needed → a directed `Read` of its deployed path. This reads all three existing precedents as coherent rather than contradictory, and makes `reference-deployed-paths.md`'s conditional-skill rule the middle branch of the criterion, stated without its qualifier.

## Reglas verificables
- **[manual]** A phase agent frontmatter reaching exactly one conditional skill MUST NOT add the `Skill` tool nor a `skills:` entry for it — it reaches the skill via a directed `Read` of `~/.claude/skills/<skill>/SKILL.md`, stated in its own numbered instructions.
- **[manual]** A phase agent frontmatter reaching more than one conditional skill MAY add the `Skill` tool and name each skill in prose instead of a directed `Read` per skill.
- **[manual]** A phase agent frontmatter whose need for a skill is unconditional (every run of the phase reaches it) preloads it via `skills:` instead of either of the above.
- **[manual]** sdd-propose and sdd-explore reach resolve-library-docs by the identical mechanism — a directed Read, never diverging between the two.

## Alternativas consideradas
(A) Add `Skill` to `tools:` and name the skill in prose — what `sdd-apply` does, and what `payload/docs/reference-deployed-paths.md` prescribes for conditional skills generally. Rejected: `Skill` is an unbounded grant — it lets the agent invoke any installed skill, to reach one. It pays for itself in an agent that reaches several conditional skills (`sdd-apply` reaches `resolve-library-docs`, the debugger skill and others); it does not in an agent that reaches exactly one, which is the case for both `sdd-propose` and `sdd-explore`, and for `sdd-design` before them. Both frontmatters justify every grant they hold narrowly and even narrow `Bash`'s charter in prose; an unbounded grant sits badly there. `Read` is already held by both agents, so this option would add a grant where the chosen option adds none.
(C) A `skills:` preload naming resolve-library-docs — what `development-decisions-mine` does for this same skill. Rejected on `sdd-design`'s own stated reasoning: a preload loads the full skill body on every run of a phase whose trigger fires only when a library is named, most of the time for nothing.

## Consecuencias
sdd-propose and sdd-explore gain no new tool beyond mcp__context7 (bare, server-level) for this capability. Both phases' instructions carry an explicit reach line (a numbered step) naming resolve-library-docs and its deployed path, so the grant is never inert surface. Future phases that need to reach a conditional skill can resolve the mechanism from the three-way criterion this record fixes, instead of re-deriving it from scratch or picking whichever precedent they saw last.

## Relacionados
- `spec` → [../../development-specs/rule/library-facts-resolved-not-recalled.md](../../development-specs/rule/library-facts-resolved-not-recalled.md) — capability-spec this decision implements — requiere que ambas fases empiecen cableadas por el mismo mecanismo
