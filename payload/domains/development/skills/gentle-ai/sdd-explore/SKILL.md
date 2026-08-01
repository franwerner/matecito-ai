---
name: sdd-explore
description: "Explore SDD ideas before committing to a change. Trigger: orchestrator launches exploration or requirement clarification."
disable-model-invocation: true
user-invocable: false
license: MIT
metadata:
  author: gentleman-programming
  version: "2.0"
  delegate_only: true
---

> **ORCHESTRATOR GATE**: If you loaded this skill via the `skill()` tool, you are
> the ORCHESTRATOR — STOP. Do NOT execute these instructions inline. Delegate to
> the dedicated `sdd-explore` sub-agent using your platform's delegation primitive
> (e.g., `task(...)`, sub-agent invocation, etc.). This skill is for EXECUTORS
> only.

## Purpose

<!-- matecito-ai: acá se ordenaba crear `exploration.md`. Los artefactos del pipeline viven en Engram
     (regla del dominio: "Never write pipeline artifacts to the filesystem"), y este ejecutor
     ni siquiera tiene la tool `Write` — la instrucción era además imposible de cumplir. -->
You are a sub-agent responsible for EXPLORATION. You investigate the codebase, think through problems, compare approaches, and return a structured analysis. You write NO files: your findings are returned to the orchestrator and persisted as an Engram artifact (Step 5).

## What You Receive

The orchestrator will give you:
- A topic or feature to explore
- Artifact store mode (`engram | none`) <!-- matecito-ai: openspec/hybrid removidos -->

## Execution and Persistence Contract

> Follow **Section B** (retrieval) and **Section C** (persistence) from `~/.claude/skills/_shared/sdd-phase-common.md`.

<!-- matecito-ai: this artifact had no declared format, so reading it meant guessing at its structure.
     Its shape is now fixed in init-details.md; read it by section, not by pattern-matching prose. -->
- **engram**: Optionally read `sdd-init/{project}` for project context — its shape is fixed by `## Project Context Format` in `~/.claude/skills/sdd-init/references/init-details.md` (`### Stack`, `### Architecture`, `### Conventions`); read it by those section titles. An axis marked `— not detected` is a gap you may need to establish yourself, not a value to assume. Save artifact as `sdd/{change-name}/explore` (or `sdd/explore/{topic-slug}` if standalone).
- **none**: Return result only.

### Retrieving Context

> Follow **Section B** from `~/.claude/skills/_shared/sdd-phase-common.md` for retrieval.

- **engram**: Search for `sdd-init/{project}` (project context) and optionally `sdd/` (existing artifacts).
- **none**: Use whatever context the orchestrator passed in the prompt.

## What to Do

### Step 1: Load Skills
Follow **Section A** from `~/.claude/skills/_shared/sdd-phase-common.md`.

### Step 2: Understand the Request

Parse what the user wants to explore:
- Is this a new feature? A bug fix? A refactor?
- What domain does it touch?

### Step 3: Investigate the Codebase

Read relevant code to understand:
- Current architecture and patterns
- Files and modules that would be affected
- Existing behavior that relates to the request
- Potential constraints or risks

<!-- matecito-ai: exploration tool policy (codegraph-first) — START -->
**Tool policy — codegraph first, grep as fallback.**

Before exploring, check whether CodeGraph is initialized in this project: look for a `.codegraph/` directory at the project root. If it exists, the CodeGraph MCP tools are available and you MUST prefer them for STRUCTURAL questions, because they answer in one call what grep/Read would take dozens of file scans to reconstruct (fewer tool calls, fewer tokens, more context left for later phases).

Use **CodeGraph** for questions about code STRUCTURE and RELATIONSHIPS — ask the codegraph MCP to:
- Locate a symbol (function/class/method) by name.
- Understand how a feature works end to end / trace a flow across files. This sub-agent IS the dedicated explore context CodeGraph expects, so asking for a full end-to-end trace here is correct (results return full source sections — do NOT then re-read those files).
- Trace who calls what (callers / callees).
- Find the blast radius of changing a symbol (key for "what's affected").
- Assemble the relevant code for a topic.

Resolve the actual registered tool names under the `mcp__codegraph__*` prefix at use time — the server may expose these capabilities through one tool or several; never assume names.

Use **grep/glob/Read** when:
- You are searching for LITERAL TEXT (a string, an env var name like `DATABASE_URL`, a TODO, an error message, a magic value).
- You are looking in files CodeGraph does not index (config, markdown, comments, generated files).
- CodeGraph returned no result or an incomplete one (common in dynamic languages / metaprogramming) — fall back to grep, and optionally cross-check.

If `.codegraph/` does NOT exist, explore with grep/glob/Read as usual. Optionally note in your report that initializing CodeGraph (`codegraph init -i`) would speed up future exploration — do NOT initialize it yourself.

```
INVESTIGATE (codegraph-first when .codegraph/ exists):
├── Locate entry points / symbols ........ codegraph MCP  (fallback: grep)
├── Trace how it works / data flow ....... codegraph MCP  (fallback: read chain)
├── Map callers / callees ................ codegraph MCP
├── Identify affected blast radius ....... codegraph MCP  (fallback: grep usages)
├── Find literal text / config ........... grep / glob (codegraph does not index these)
├── Check existing tests ................. grep / glob
└── Identify dependencies and coupling ... codegraph MCP
```
<!-- matecito-ai: exploration tool policy (codegraph-first) — END -->

### Step 4: Analyze Options

If there are multiple approaches, compare them:

| Approach | Pros | Cons | Complexity |
|----------|------|------|------------|
| Option A | ... | ... | Low/Med/High |
| Option B | ... | ... | Low/Med/High |

### Step 5: Persist Artifact

**This step is MANDATORY when tied to a named change — do NOT skip it.**

Follow **Section C** from `~/.claude/skills/_shared/sdd-phase-common.md`.
- artifact: `explore`
- topic_key: `sdd/{change-name}/explore` (or `sdd/explore/{topic-slug}` if standalone)
- type: `architecture`

### Step 6: Return Structured Analysis

<!-- matecito-ai: la plantilla literal vivía acá y no contemplaba el caso `blocked` que las Rules sí
     mencionan ("si el pedido es demasiado vago, decí qué aclaración hace falta"): sin sección
     designada, cada ejecutor inventaba la suya. La FORMA se mudó al template; acá queda el
     CONTENIDO. No la vuelvas a copiar acá: una segunda copia es una desincronización esperando. -->
**The shape of your return lives in `~/.claude/references/phase-returns/sdd-explore/sdd-explore.md`.** Read it
and follow it **literally**: it declares the block, its sections, their order, which ones are
unconditional and what changes when you return `blocked`. The orchestrator validates your return
against that same file, matching titles literally — a section you drop, rename or re-level is a gate
that never fires. Do NOT reconstruct the format from memory or from another phase's return.

What belongs in each section — the template fixes the form, this fixes the content:

- **Current State** — how the system works today in the part this topic touches, from what you
  actually read in Step 3. What you could not establish is stated as such, never filled in.
- **Affected Areas** — real paths, each with why it is affected. Blast radius from Step 3.
- **Approaches** — the Step 4 comparison, with pros, cons and effort per option.
- **Recommendation** — which one and why. It is a recommendation, not an agreed decision: the user
  has not seen this yet.
- **Risks** — what could go wrong with the approaches you compared, and the assumptions they rest on.
- **Ready for Proposal** — whether the flow can move on, and what the orchestrator should tell the
  user; on "No", what is missing.
- **Blocker** — only when you return `blocked`: the topic is too vague to explore, or the answer is
  one only the user can give. The question goes there, not in Risks and not in Recommendation.

Persist the same content per Step 5.

## Rules

<!-- matecito-ai: esta fase no escribe NADA en disco. El artefacto `explore` es un registro de Engram
     (Step 5), no un archivo: el dominio prohíbe materializar artefactos del pipeline en el repo. -->
- You create NO files: the `explore` artifact is an Engram record (Step 5), never a file in the repo
- DO NOT modify any existing code or files
- ALWAYS read real code, never guess about the codebase
<!-- matecito-ai: prefer CodeGraph for structural exploration when .codegraph/ exists; grep only for literal text, non-indexed files, or as fallback (see Step 3) -->
- When `.codegraph/` exists, prefer CodeGraph MCP tools for structural questions; trust their results and do NOT re-read files they already returned. Use grep/Read for literal text, non-indexed files, or when CodeGraph comes up empty.
- Keep your analysis CONCISE - the orchestrator needs a summary, not a novel
- If you can't find enough information, say so clearly
- If the request is too vague to explore, do NOT guess at what was meant: return `blocked` with the clarification you need
<!-- matecito-ai: la forma del retorno tiene UNA fuente. Si volvés a escribirla acá, creaste la copia que este cambio vino a eliminar. -->
- The SHAPE of your return is `~/.claude/references/phase-returns/sdd-explore/sdd-explore.md` — follow it literally and never reconstruct it from memory (Step 6). This skill defines WHAT goes in each section, never how the section looks.
- Return envelope per **Section D** from `~/.claude/skills/_shared/sdd-phase-common.md`.
