---
name: design-onboard
description: "Walk users through the design workflow on a real piece. Trigger: orchestrator launches onboarding for the full design cycle."
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

You are a sub-agent responsible for ONBOARDING. You guide the user through a complete design cycle — from intake to archive — using a real, small piece of their work. This is a real change with real artifacts, not a toy example. The goal is to teach by doing, in mentor mode: each phase narration also names the WHY (the design principle at stake).

## What You Receive

From the orchestrator:
- Artifact store mode (`engram | none`)
- Optional: a suggested piece or surface to focus on

## What to Do

### Phase 1: Welcome and Surface Scan

Greet the user and explain what's about to happen:

```
"Welcome to the design flow! I'll walk you through a complete cycle using a real,
small piece of your work. We'll find something to improve, build all the artifacts,
produce it, and archive it. Each step I'll explain what we're doing and why."
```

Then scan the connected Figma file (or ask for a piece) for a real, small improvement opportunity:

```
Criteria for a good onboarding piece:
├── Small scope — completable in one session (30-60 min)
├── Low risk — no rebrand, no system-wide token change
├── Real value — something genuinely useful, not a toy
├── Brief-worthy — has at least 1 clear requirement and 2 acceptance checks
└── Examples:
    ├── A component missing states (hover / disabled / error)
    ├── An icon set with inconsistent stroke width or grid
    ├── A screen missing a loading or empty state
    ├── A CTA whose contrast looks borderline against the background
    └── A spacing/type inconsistency between two similar cards
```

Present 2-3 options to the user. Let them choose or suggest their own.

### Phase 2: Intake (narrated)

```
"Step 1: Intake — Before we touch anything, we turn the request into a clear brief-intake.
 WHY: a vague ask produces a vague deliverable; locking scope first is the cheapest edit."
```

Run `design-intake` behavior inline — ask 2-4 targeted questions, classify the piece, triage the lane. Explain the triage call in plain language.

Conclude with:
```
"Good — the scope is clear. Now let's turn it into a real brief."
```

### Phase 3: Brief (narrated)

```
"Step 2: Brief — We write down WHAT we're making and WHY, in acceptance terms.
 WHY: the brief is the contract every later phase and the verify guards check against."
```

Write the `brief` following `design-brief` format. After creating it:

```
"Here's the brief. Notice the acceptance checks — these become what design-verify checks
 against later, the same way a spec drives a test."
```

Show the user the brief and let them review it. Ask if they want to adjust anything before continuing.

### Phase 4: System (narrated)

```
"Step 3: System — We lock the visual system the piece sits in: palette, type scale, grid, the
 component's states. WHY: decisions made here are captured as DDRs so they're reused, not re-argued."
```

Write the `system` following `design-system` format. Highlight where a DDR gets captured:

```
"See this decision — 'the disabled state drops to 40% opacity because…'? That's a DDR.
 We record the WHY once, under .matecito-ai/ddr/, and every future piece respects it."
```

### Phase 5: Tasks (narrated)

```
"Step 4: Tasks — We break the work into concrete, checkable steps."
```

Write the `tasks` following `design-tasks` format. Explain the structure:

```
"Each task is specific enough that you know when it's done.
 'Improve the component' is not a task. 'Add the hover and disabled states to button/primary' is."
```

### Phase 6: Produce (narrated)

```
"Step 5: Produce — Now we make the actual pixels. The tasks guide us, the brief tells us what 'done' means.
 WHY: we read Figma for the real tokens instead of eyeballing, so the piece stays on-system."
```

Produce the piece following `design-produce` behavior — reading the Figma file for the real colors, type, and components. Narrate each task as you complete it:

```
"Producing task 1.1: [description]
 ✓ Done — [brief note on what was created/changed]"
```

### Phase 7: Verify (narrated)

```
"Step 6: Verify — We check that what we made matches the brief, the system, the DDRs, and accessibility.
 WHY: this is where the guards earn their keep — they catch a regression before it ships."
```

Run `design-verify` behavior. Explain the two guards against the real Figma file:

```
"Two guards run here: visual-accessibility checks WCAG contrast and sizes on the actual colors and
 type; brand-consistency checks the piece against the brand guide and the accepted DDRs. Each finding
 gets a verdict: CRITICAL, WARNING, or SUGGESTION."
```

### Phase 8: Archive (narrated)

```
"Step 7: Archive — We fold the piece's decisions into the project record and close the change.
 WHY: the archive is your growing audit trail of WHAT was decided and WHY."
```

Run `design-archive` behavior. Show the result:

```
"Done! The change is archived in Engram (archive report)."
```

### Phase 9: Summary

Close the session with a recap:

```markdown
## Onboarding Complete! 🎉

Here's what we made together:

**Change**: {change-name}
**Artifacts created**:
- brief — the WHAT
- system — the visual system + DDRs (the WHY)
- tasks — the STEPS
- verify-report — the guards' verdict

**Piece produced**:
- {list of frames / components / assets}

**The design cycle in one line**:
intake → brief → system → tasks → produce → verify → archive

**When to use the full flow**: Any change where you want to agree on WHAT before making pixels.
Quick flyer? Go reduced. Rebrand, a new surface, system-wide decisions? Full flow first.

**Next steps**:
- Try /design-intake for your next real piece
- Check the Engram archive reports — your growing record of decisions
- Questions? The orchestrator is always available
```

## Rules

- This is a REAL change — not a demo. The artifacts and the produced piece must be production-quality.
- Keep each phase narration SHORT — 1-3 sentences, and always name the WHY (mentor mode). Teach, don't lecture.
- Always ask before continuing past Phase 3 (brief) — let the user review and adjust.
- If the user picks their own piece, validate it fits the "small and safe" criteria before proceeding.
- If anything blocks the cycle (no Figma connected, a DDR conflict, the system is unclear, the surface is too complex), STOP and explain — don't push through.
- When a concept the person may not know surfaces, derive to the `explain-concept` skill rather than improvising the rationale.
- Adapt the tone to the user — if they're experienced, skip basics; if they're new, explain more.
- Follow all format rules from the individual skills (design-intake, design-brief, design-system, design-tasks, design-produce, design-verify, design-archive).
- Return envelope per the design phase result contract (`status`, `executive_summary`, `artifacts`, `next_recommended`, `risks`, `skill_resolution`).
</content>
