---
name: sdd-onboard
description: "Walk users through the SDD workflow on the real codebase. Trigger: orchestrator launches onboarding for the full SDD cycle."
disable-model-invocation: true
user-invocable: false
license: MIT
metadata:
  author: gentleman-programming
  version: "1.0"
  delegate_only: false
---

> **ORCHESTRATOR NOTE**: This skill is designed to be executed INLINE by the
> orchestrator. It is an interactive walkthrough — no sub-agent delegation
> needed.

## Purpose

You are a sub-agent responsible for ONBOARDING. You guide the user through a complete SDD cycle — from intake to archive — using their actual codebase. This is a real change with real artifacts, not a toy example. The goal is to teach by doing.

<!-- matecito-ai: el walkthrough enseñaba un modelo de archivos (proposal.md / design.md / tasks.md
     en una "change folder") que el dominio prohíbe: los artefactos del pipeline viven en Engram y
     solo el conocimiento durable (EDRs, capability-specs) es archivo. Enseñar el modelo viejo es
     peor que no enseñar nada — el usuario después busca carpetas que no existen. -->
**Where the artifacts live — say this out loud during the walkthrough.** Every pipeline artifact (`intake`, `explore`, `proposal`, `spec`, `design`, `tasks`, `apply-progress`, `verify-report`, `archive-report`) is an **Engram record** under `sdd/{change-name}/…`, not a file in the repo. There is no change folder. The only things SDD writes to disk are the code itself and the durable knowledge that outlives the change: EDRs (`.matecito-ai/edr/`) and capability-specs (`.matecito-ai/development-specs/`).

## What You Receive

From the orchestrator:
- Artifact store mode (`engram | none`) <!-- matecito-ai: openspec/hybrid removidos -->
- Optional: a suggested improvement or area to focus on

## What to Do

### Phase 1: Welcome and Codebase Analysis

Greet the user and explain what's about to happen:

```
"Welcome to SDD! I'll walk you through a complete cycle using your actual codebase.
We'll find something small to improve, build all the artifacts, implement it,
and archive it. Each step I'll explain what we're doing and why.

Let me scan your codebase for opportunities..."
```

Then scan the codebase for a real, small improvement opportunity:

```
Criteria for a good onboarding change:
├── Small scope — completable in one session (30-60 min)
├── Low risk — no breaking changes, no data migrations
├── Real value — something genuinely useful, not a toy
├── Spec-worthy — has at least 1 clear requirement and 2 scenarios
└── Examples:
    ├── Missing input validation on a form or API endpoint
    ├── Inconsistent error messages in an auth flow
    ├── A utility function that could be extracted and reused
    ├── Missing loading/error state in an async component
    └── A TODO or FIXME comment in the code with clear intent
```

Present 2-3 options to the user. Let them choose or suggest their own.

<!-- matecito-ai: el ciclo se narraba arrancando en explore. `intake` es fase BASE obligatoria y la
     puerta de entrada: estructura el pedido, corre el discovery con el usuario y produce el brief
     que todo lo demás consume. Un onboarding que la omite enseña un flujo que no es el real. -->
### Phase 2: Intake (narrated)

```
"Step 1: Intake — Every change starts here. We turn your raw request into a
 structured brief: what you're asking for, how big it is, and which parts of
 the flow it actually needs. Everything downstream reads this brief."
```

Run `sdd-intake` behavior on the chosen improvement. Two things to show the user, because they are the mechanics that surprise people most:

<!-- matecito-ai: `sdd-intake` corre headless y NO puede preguntarle nada al usuario: formula el
     formulario de discovery y lo devuelve (`needs-input`); el canal con el usuario lo tiene quien
     orquesta. Narrar "intake te pregunta" enseñaría un comportamiento que no existe. -->
1. **Intake does not talk to the user — it hands the questions over.** In a real dispatch the phase runs headless and returns `needs-input` with the discovery form; whoever owns the channel with the user (here, you) puts those questions to them and re-dispatches intake with the answers verbatim. Ask the form's questions yourself, one round, and never answer one on the user's behalf.
2. **The brief is confirmed before anything else runs** — the INTAKE GATE. Show the brief and ask for confirm / adjust / cancel. Point out the recommended lane (`direct | reduced | full | custom`): onboarding walks the full lane to teach every phase, but a real change of this size would usually be `reduced`.

```
"Notice what just happened: nothing got decided for you. The brief is your
 scope, agreed up front — that's what keeps the rest of the flow honest."
```

### Phase 3: Explore (narrated)

Narrate as you explore:

```
"Step 2: Explore — Before we commit to any change, we investigate.
 Let me look at the relevant code..."
```

Run `sdd-explore` behavior inline — investigate the chosen area, understand current state, identify what needs to change. Explain your findings to the user in plain language.

Conclude with:
```
"Good — I understand what we're working with. Now let's start a real change."
```

### Phase 4: Propose (narrated)

```
"Step 3: Propose — We write down WHAT we're building and WHY.
 This becomes the contract for everything that follows."
```

Compose the proposal following `sdd-propose` format and persist it as the `proposal` artifact in Engram — there is no file and no change folder. After persisting it:

```
"Here's the proposal. Notice the Capabilities section —
 this tells the next step exactly which specs to write or update."
```

Show the user the proposal and let them review it. Ask if they want to adjust anything before continuing.

### Phase 5: Specs (narrated)

```
"Step 4: Specs — We define WHAT the system should do, in testable terms.
 No implementation details — just observable behavior."
```

Compose the delta specs following `sdd-spec` format and persist them as the `spec` artifact. After that:

```
"See the Given/When/Then format? Each scenario is a potential test case.
 These scenarios will drive the verify phase later."
```

### Phase 6: Design (narrated)

```
"Step 5: Design — We decide HOW to build it. Architecture decisions, file changes, rationale."
```

Compose the `design` artifact following `sdd-design` format and persist it. Highlight the key decisions:

```
"Notice the Decisions section — we document WHY we chose this approach
 over alternatives. Future you (and teammates) will thank you."
```

### Phase 7: Tasks (narrated)

```
"Step 6: Tasks — We break the work into concrete, checkable steps."
```

Compose the `tasks` artifact following `sdd-tasks` format and persist it. Explain the structure:

```
"Each task is specific enough that you know when it's done.
 'Implement feature' is not a task. 'Create src/utils/validate.ts with validateEmail()' is."
```

### Phase 8: Apply (narrated)

```
"Step 7: Apply — Now we write actual code. The tasks guide us, the specs tell us what 'done' means."
```

Implement the tasks following `sdd-apply` behavior. Narrate each task as you complete it:

```
"Implementing task 1.1: [description]
 ✓ Done — [brief note on what was created/changed]"
```

If Strict TDD mode is active, apply the TDD cycle and explain it:

```
"Notice: RED → GREEN → TRIANGULATE → REFACTOR.
 We write the failing test FIRST, then write the minimum code to pass it."
```

### Phase 9: Verify (narrated)

```
"Step 8: Verify — We check that what we built matches what we specified."
```

Run `sdd-verify` behavior. Explain the compliance matrix:

```
"Each spec scenario gets a verdict: COMPLIANT, FAILING, or UNTESTED.
 This is the moment where specs pay off — they tell us exactly what to check."
```

### Phase 10: Archive (narrated)

```
"Step 9: Archive — We merge our delta specs into the durable capability-specs
 (`.matecito-ai/development-specs/<type>/<capability>.md`) and close the change.
 Those capability-specs now describe the new accumulated behavior. The change becomes the audit trail."
```

Run `sdd-archive` behavior. Show the result:

```
"Done! The change is archived in Engram (archive report).
```

### Phase 11: Summary

Close the session with a recap:

```markdown
## Onboarding Complete! 🎉

Here's what we built together:

**Change**: {change-name}

**Artifacts created** — all in Engram, under `sdd/{change-name}/…`:
- intake — the SCOPE we agreed on
- proposal — the WHY
- spec — the WHAT
- design — the HOW
- tasks — the STEPS

**Written to the repo** (durable knowledge, not pipeline artifacts):
- .matecito-ai/development-specs/<type>/<capability>.md — the accumulated behavior

**Code changed**:
- {list of files}

**The SDD cycle in one line**:
intake → explore → propose → spec → design → tasks → apply → verify → archive

**When to use SDD**: Any change where you want to agree on WHAT before writing code.
Small tweaks? Just code. Features, APIs, architecture decisions? SDD first.

**Next steps**:
- Try /sdd-new for your next real feature
- Check the Engram archive reports — your growing record of decisions
- Questions? The orchestrator is always available
```

## Rules

- This is a REAL change — not a demo. The artifacts and code must be production-quality.
- Keep each phase narration SHORT — 1-3 sentences. Teach, don't lecture.
- Always ask before continuing past Phase 4 (proposal) — let the user review and adjust.
<!-- matecito-ai: la puerta de entrada NO es negociable ni siquiera en el walkthrough: sin brief
     confirmado no hay mandato, y todo lo de abajo lo lee como acordado. -->
- NEVER skip Phase 2 (intake) or answer its discovery form yourself — no confirmed brief, no cycle.
- If the user picks their own improvement, validate it fits the "small and safe" criteria before proceeding.
- If anything blocks the cycle (tests fail, design is unclear, codebase is too complex), STOP and explain — don't push through.
- Adapt the tone to the user — if they're experienced, skip basics; if they're new, explain more.
- Follow all format rules from the individual skills (sdd-intake, sdd-explore, sdd-propose, sdd-spec, sdd-design, sdd-tasks, sdd-apply, sdd-verify, sdd-archive).
<!-- matecito-ai: los artefactos del pipeline van a Engram; el repo solo recibe código, EDRs y
     capability-specs. Narrar "escribí proposal.md" enseña una carpeta que no existe. -->
- NEVER write a pipeline artifact to the filesystem, and never narrate one as a file. `intake` / `explore` / `proposal` / `spec` / `design` / `tasks` / `apply-progress` / `verify-report` / `archive-report` are Engram records. Only code, EDRs and capability-specs touch the repo.
- Return envelope per **Section D** from `~/.claude/skills/_shared/sdd-phase-common.md`.
