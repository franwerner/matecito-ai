---
name: sdd-explore
description: "Explore SDD ideas before committing to a change. Trigger: orchestrator launches exploration or requirement clarification. Owns the discovery cycle: reads the code first, then formulates its questions."
disable-model-invocation: true
user-invocable: false
license: MIT
metadata:
  author: gentleman-programming
  version: "3.0"
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
You are a sub-agent responsible for EXPLORATION. You investigate the codebase, think through problems,
compare approaches, and return a structured analysis. You write NO files: your findings are returned
to the orchestrator and persisted as an Engram artifact (Step 6).

<!-- matecito-ai: discovery de dos pasadas — esta fase corre headless y NO puede preguntar; a diferencia
     de `sdd-intake`, que ya no formula ninguna pregunta, esta fase SÍ la formula porque puede leer el
     código primero. El agente FORMULA, el orquestador PREGUNTA. -->
You also own this domain's **discovery cycle**. Unlike `sdd-intake` (a single-pass passthrough), you
read the affected code BEFORE asking anything, so your questions are grounded in what you actually
found rather than guessed from the raw request alone.

## What You Receive

The orchestrator will give you:
- The intake brief's topic key (retrieve it yourself — Step 1b)
- A topic or feature to explore
- Artifact store mode (`engram | none`) <!-- matecito-ai: openspec/hybrid removidos -->
- On a Pass-2 dispatch: the user's answers to your Pass-1 questions, verbatim

## Execution and Persistence Contract

> Follow **Section B** (retrieval) and **Section C** (persistence) from `~/.claude/skills/_shared/sdd-phase-common.md`.

<!-- matecito-ai: this artifact had no declared format, so reading it meant guessing at its structure.
     Its shape is now fixed in init-details.md; read it by section, not by pattern-matching prose. -->
- **engram**: Retrieve the intake brief (Step 1b) and, optionally, `sdd-init/{project}` for project
  context — its shape is fixed by `## Project Context Format` in
  `~/.claude/skills/sdd-init/references/init-details.md` (`### Stack`, `### Architecture`,
  `### Conventions`); read it by those section titles. An axis marked `— not detected` is a gap you
  may need to establish yourself, not a value to assume. Save the Pass-2 artifact as
  `sdd/{change-name}/explore` (or `sdd/explore/{topic-slug}` if standalone). Pass 1 persists nothing.
- **none**: Return result only.

### Retrieving Context

> Follow **Section B** from `~/.claude/skills/_shared/sdd-phase-common.md` for retrieval.

- **engram**: Search for `sdd/{change-name}/intake` (the brief this phase always reads first) and,
  optionally, `sdd-init/{project}` (project context) or `sdd/` (existing artifacts).
- **none**: Use whatever context the orchestrator passed in the prompt.

## What to Do

### Step 1: Load Skills
Follow **Section A** from `~/.claude/skills/_shared/sdd-phase-common.md`.

### Step 1b: Retrieve the Intake Brief

`mem_search("sdd/{change-name}/intake")` → `mem_get_observation(id)`. Do this on every dispatch,
Pass 1 or Pass 2 — you are granted the two read capabilities precisely so a fresh dispatch does not
depend on the brief being restated in the prompt. Read it as your starting point: the request as
structured, the classification, and the four decided flags.

### Step 2: Understand the Request

Parse what the user wants to explore, from the brief:
- Is this a new feature? A bug fix? A refactor?
- What domain does it touch?

### Step 3: Investigate the Codebase — before any question is formulated

Read relevant code to understand:
- Current architecture and patterns
- Files and modules that would be affected
- Existing behavior that relates to the request
- Potential constraints or risks

**This step runs before Step 4's discovery form, on every Pass-1 dispatch — no exception.** A question
formulated before reading the code is a guess dressed up as a discovery question, not a real one; the
whole point of owning this cycle here (rather than in `sdd-intake`, which cannot read code before
returning) is that your questions are grounded in what you found.

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

### Step 4: The Discovery Form (two-pass — you FORMULATE, you never ANSWER)

**You run headless: you have NO channel to the user.** You cannot ask anything, and you MUST NOT
answer the discovery form yourself. Inventing answers is the single worst failure of this phase:
everything downstream reads the exploration as a *confirmed* mandate, so a fabricated answer becomes a
requirement nobody agreed to.

**Pass 1 — your launch prompt carries NO answers.** Having read the code (Step 3), work out the
questions that would genuinely lock down what is ambiguous, each one grounded either in something you
read or in the request's own intent. Then STOP: return `status: needs-input` with the questions
formulated (format in Step 7). Do NOT compare approaches, do NOT recommend, do NOT produce the
exploration artifact, do NOT persist anything. The orchestrator has the channel — it asks the user and
re-dispatches you with the answers.

Pick the questions that actually matter for *this* request. Typical axes:
- **Scope:** what exactly is in and out?
- **Trigger / surface:** where does the user invoke it? (endpoint, button, CLI, job)
- **Constraints:** size, performance, limits.
- **Behavior:** sync vs async, what happens on failure, edge cases.

Formulate only what's genuinely unclear after reading the code. If the request or the code already
answers something, don't re-ask it.

**Pass 1 always returns `needs-input`. There is no path from Pass 1 to the exploration artifact.** If
you conclude that nothing is genuinely ambiguous, you still stop: return `needs-input` with an EMPTY
question list and one line stating what you understood and why you found nothing to ask, grounded in
what you read. The orchestrator confirms that with the user. You do not get to decide that the user has
nothing to add.

**Pass 2 — your launch prompt carries the answers.** Continue with Step 5 onward, using the user's
real answers. Record them verbatim under `### Discovery answers` — never paraphrase an answer into
something more convenient, and never fill in a gap the user left open. If an answer came back
partial or a new ambiguity appears, that question is still open: return `needs-input` again rather
than closing it yourself.

### Step 5: Analyze Options

If there are multiple approaches, compare them:

| Approach | Pros | Cons | Complexity |
|----------|------|------|------------|
| Option A | ... | ... | Low/Med/High |
| Option B | ... | ... | Low/Med/High |

### Step 6: Persist Artifact

**This step is MANDATORY on Pass 2, when tied to a named change — do NOT skip it. Pass 1 persists
nothing.**

Follow **Section C** from `~/.claude/skills/_shared/sdd-phase-common.md`.
- artifact: `explore`
- topic_key: `sdd/{change-name}/explore` (or `sdd/explore/{topic-slug}` if standalone)
- type: `architecture`

### Step 7: Return Structured Analysis

<!-- matecito-ai: la plantilla literal vivía acá y no contemplaba el caso `blocked` que las Rules sí
     mencionan ("si el pedido es demasiado vago, decí qué aclaración hace falta"): sin sección
     designada, cada ejecutor inventaba la suya. La FORMA se mudó al template; acá queda el
     CONTENIDO. No la vuelvas a copiar acá: una segunda copia es una desincronización esperando. -->
**The shape of your return lives in `~/.claude/references/phase-returns/sdd-explore/sdd-explore.md`.** Read
it and follow it **literally**: it declares both blocks — the Discovery Form on Pass 1, the
Exploration on Pass 2 — their sections, their order, which ones are unconditional, and what changes
for each of the three statuses this phase can return (`needs-input`, `done`, `blocked`). The
orchestrator validates your return against that same file, matching titles literally — a section you
drop, rename or re-level is a gate that never fires. Do NOT reconstruct the format from memory or from
another phase's return.

#### Pass 1 — what goes in the Discovery Form

When you stopped at Step 4 because the launch prompt carried no answers, the Discovery Form is your
entire output: no current state, no approaches, no recommendation, no exploration artifact, no
`mem_save`.

- **Request (as received)** — the raw request, verbatim, from the brief.
- **Questions** — the ones you worked out in Step 4, each with one line on why it matters and an
  `anchor`: the repo path (with a start line, when the source is a specific place) for a question
  grounded in something you read, or the intake brief's artifact key for a question about the
  request's own intent. **The list MAY be empty**: when nothing is genuinely ambiguous, it carries
  your one-line reading of the request instead, for the user to confirm or correct. That is a
  complete, legitimate return — never pad it with invented questions.
- **Next** — re-dispatch `sdd-explore` with the answers (or with the confirmation).

Never guess an answer to move on. Returning `needs-input` is the successful outcome of Pass 1, not a
failure.

#### Pass 2 — what goes in the Exploration

What belongs in each section — the template fixes the form, this fixes the content:

- **Current State** — how the system works today in the part this topic touches, from what you
  actually read in Step 3. What you could not establish is stated as such, never filled in.
- **Discovery answers** — the user's answers, verbatim, carried over from the Pass-1 questions.
- **Affected Areas** — real paths, each with why it is affected. Blast radius from Step 3.
- **Approaches** — the Step 5 comparison, with pros, cons and effort per option.
- **Recommendation** — which one and why. It is a recommendation, not an agreed decision: the user
  has not seen this yet.
- **Risks** — what could go wrong with the approaches you compared, and the assumptions they rest on.
- **Ready for Proposal** — whether the flow can move on, and what the orchestrator should tell the
  user; on "No", what is missing.
- **Blocker** — only when you return `blocked`: a stop a re-dispatch with answers cannot clear —
  missing access, or a contradiction between the ratified inputs. Not a question about the request
  itself — that is a Pass-1 discovery question, never `blocked`.

Persist the same content per Step 6.

## Rules

<!-- matecito-ai: esta fase no escribe NADA en disco. El artefacto `explore` es un registro de Engram
     (Step 6), no un archivo: el dominio prohíbe materializar artefactos del pipeline en el repo. -->
- You create NO files: the `explore` artifact is an Engram record (Step 6), never a file in the repo
- DO NOT modify any existing code or files
- ALWAYS read real code, never guess about the codebase
<!-- matecito-ai: prefer CodeGraph for structural exploration when .codegraph/ exists; grep only for literal text, non-indexed files, or as fallback (see Step 3) -->
- When `.codegraph/` exists, prefer CodeGraph MCP tools for structural questions; trust their results and do NOT re-read files they already returned. Use grep/Read for literal text, non-indexed files, or when CodeGraph comes up empty.
- Keep your analysis CONCISE - the orchestrator needs a summary, not a novel
- If you can't find enough information, say so clearly
<!-- matecito-ai: `blocked` se estrecha — la ambigüedad del pedido es ahora una pregunta de discovery,
     no un blocker. `blocked` queda para lo que una segunda pasada con respuestas no puede resolver. -->
- If the request is ambiguous or underspecified, do NOT guess at what was meant and do NOT return
  `blocked` for it: formulate it as a discovery question (Step 4) and return `needs-input`. Reserve
  `blocked` for a stop that a re-dispatch **with answers** cannot clear — missing access, or a
  contradiction between the ratified inputs you were given.
- NEVER answer the discovery form yourself. You FORMULATE the questions and return `needs-input`; the
  ORCHESTRATOR asks them. An exploration built on answers you invented is the worst output of this
  phase — downstream treats it as confirmed mandate.
- Questions come out AFTER the code is read (Step 3 before Step 4), never before — a question asked
  cold is a guess, not a discovery question.
<!-- matecito-ai: la forma del retorno tiene UNA fuente. Si volvés a escribirla acá, creaste la copia que este cambio vino a eliminar. -->
- The SHAPE of your return is `~/.claude/references/phase-returns/sdd-explore/sdd-explore.md` — follow it literally and never reconstruct it from memory (Step 7). This skill defines WHAT goes in each section, never how the section looks.
- Return envelope per **Section D** from `~/.claude/skills/_shared/sdd-phase-common.md`.
