<!-- matecito-ai: DESIGN DOMAIN FRAGMENT.
     Appended to core/CLAUDE.md at deploy time. Binds the kernel's generic nouns
     to design's vocabulary and adds the design-specific rules (mentor mode,
     phases, guards) that must NOT live in the kernel. -->

# matecito-ai — Design domain

For a full-spectrum visual designer (brand + UI/UX + prototypes + brand guides).
This domain does NOT touch code — the design→code handoff belongs to the
`development` domain.

## Domain vocabulary (binds the kernel's generic slots)

| Kernel slot | Design binding |
| --- | --- |
| Structured flow name | SDD (Spec-Driven **Design**) — "spec" = "brief" |
| Phase pipeline | `intake → explore → propose → brief → system → tasks → produce → verify → archive` |
| Mandatory base phases | `intake → brief → produce → verify → archive` |
| Optional add-on phases | `explore`, `propose`, `system`, `tasks` |
| Phase agents | `design-*` (`design-intake`, `design-explore`, …, `design-archive`) |
| Alignment artifact | `brief` |
| Decision record | `DDR` (Design Decision Record), stored in `.matecito-ai/ddr/` |
| Canonical catalog | `design-principles` at `~/.claude/references/design-principles/` |
| Exploration index | Figma (`figma` MCP), active when a Figma file is connected |
| Guards | `visual-accessibility`, `brand-consistency` |
| Engram topic-key namespace | `design-init/{project}` · `design/{change-name}/{intake,explore,proposal,brief,system,tasks,produce-progress,verify-report,archive-report,state}` |

## Mentor mode (cross-cutting rule)

Every phase and skill explains the **why** behind each decision or finding in 1-2
lines — the underlying design principle, not just the what. This is the learning
engine: the person gets more efficient AND learns on the way. When a concept the
person may not know surfaces, derive to the `explain-concept` skill. Cite the
canonical `design-principles` catalog rather than improvising the rationale.

## Language

Design deliverables are visual (Figma files, brand guides, exported assets), not
code. Naming inside the work — layer names, component names, design-token names,
file names — uses English, kebab/Pascal per the file's existing convention.

## SDD Flow (design)

```
design-intake → design-explore → design-propose → design-brief → design-system → design-tasks → design-produce → design-verify → design-archive
                                                                        ^
                                                                 (system reads DDRs)
```

Mirrors development's pipeline: `brief` is design's `spec`; `system` is design's
`design` phase (locks the visual system — palette, type scale, grid, components —
and reads/writes DDRs); `produce` is design's `apply`. The lane fork (in the
kernel) governs which phases run: a quick flyer goes `reduced` (base only), a
rebrand goes `full`.

### Phase → agent + skills

| Phase | Agent | Skills it uses |
| --- | --- | --- |
| explore | `design-explore` | `brand-from-references`, `figma-hygiene` |
| propose | `design-propose` | `explore-variations` |
| brief | `design-brief` | — |
| system | `design-system` | `brand-guide` |
| produce | `design-produce` | `generate-assets` |
| verify | `design-verify` | `consistency-audit`, `visual-accessibility`, `design-review` |

Standalone skills (trigger by natural language, outside a flow): `explain-concept`
(the mentor engine), and any skill above used ad hoc.

## Guards

### visual-accessibility
WCAG contrast ratios, minimum sizes, and hierarchy checked against the real colors
and type in the Figma file. Run in `design-verify`; flag anything below AA.

### brand-consistency
Every produced piece is checked against the brand guide and the accepted DDRs.
Run in `design-verify`; flag a piece that contradicts a decision record.

## Decision records (DDR)

Brand and design decisions ("the palette is X because…", "CTAs look like Y
because…") are captured once as DDRs under `.matecito-ai/ddr/`, then **respected
and verified** by the `system` and `verify` phases. Same concept and gates as the
kernel's decision records — only the term (DDR) and store differ from development's
ADRs.

## MCP

- **Figma (`figma`)** — registered by install (`claude mcp add --transport http figma https://mcp.figma.com/mcp`). Lets the agent READ the Figma file (review, audit, extract brand). OAuth once per person via `/mcp`.
- **Canva (`canva`)** — registered by install (`claude mcp add --transport http canva https://mcp.canva.com/mcp`), Canva's official hosted MCP. Lets the agent create/edit on-brand designs. OAuth once per person via `/mcp`. Do NOT use the `@canva/cli ... mcp` from tutorials — that is for building Canva apps, not for designing.

## SDD Phase Read/Write (design)

| Phase | Reads | Writes |
| --- | --- | --- |
| `design-intake` | raw request | `intake` (brief-intake) |
| `design-explore` | intake | `explore` |
| `design-propose` | exploration (optional) | `proposal` (directions) |
| `design-brief` | proposal (required) | `brief` |
| `design-system` | brief + **DDRs** | `system` |
| `design-tasks` | brief + system | `tasks` |
| `design-produce` | tasks + brief + system + produce-progress | `produce-progress` |
| `design-verify` | brief + system + produce-progress + **DDRs touched** | `verify-report` |
| `design-archive` | all artifacts | `archive-report` |
