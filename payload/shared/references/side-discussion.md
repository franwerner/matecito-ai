# Side discussion — one mechanism, read by both ends

<!-- matecito-ai: the mechanism for a user-opened side discussion in a separate, interactive Claude
     Code session. Lives in the shared tier, not a domain fragment, because it has two readers and one
     of them belongs to no domain: the orchestrator, and the side session itself, which loads this file
     as the first line of its own seed prompt to learn its protocol. The kernel's `### Side Discussion
     (opt-in)` section states the policy (what this is, that only the user opens one, what each type
     means for the main thread); this file states the mechanism only — neither restates the other. -->

This file is read by two different actors, at two different moments:

- **The orchestrator**, when the user opens a side discussion, to compose the handoff and open the
  session.
- **The side session itself**, as the first thing it reads after it starts, to learn its own protocol
  before it reads the handoff or reasons about anything.

## Opening a side discussion

The orchestrator writes a **handoff** to the artifact store, under topic_key
`side-discussion/{slug}/handoff`, where `{slug}` is a kebab-case slug of the topic that the
orchestrator fixes at this point — never chosen by the side session, never reused across unrelated
discussions.

The handoff carries one header line above five fixed sections, all filled in on every handoff — a
section with nothing to say is stated empty, never dropped:

```markdown
Read `~/.claude/references/side-discussion.md` first — it is your protocol. Then follow this handoff.

- **Type:** blocking | consultive

## Topic
{one line naming what is under discussion}

## What I already read
- {path}[:line] — {what it settled or contributed}

## References
- {deployed path or artifact key} — {why the side session must read it}

## Open question
{the exact question to answer}

## Return
Write your conclusion to Engram under topic_key `side-discussion/{slug}/conclusion`.
That write is your only output: read and reason, change nothing in the repo, commit nothing.
```

The `Type` line is not a sixth section — it governs what the *main thread* does while the discussion is
open (see "Pickup" below), not what the side session has to read. It sits in the handoff, not only in
the orchestrator's own head, so a resumed or compacted side session can recover it without asking again.

## Opening the session — automatic, tested by five properties, no tool named

**The orchestrator opens the session itself — automatically, with no action required from the user.**
Handing the user a command to run does not satisfy this: the mechanism opens the session, it does not
delegate opening it.

The launch must satisfy five properties. None of them names a tool, in this file or anywhere else in the
payload — not as an example, not as a default, not as "the one that works today":

1. **A new terminal** — a session separate from this one, not this thread continuing.
2. **Interactive** — the user can read what it says and type into it.
3. **Opened without the user doing anything** — the orchestrator opens it.
4. **Its start is known — and starting means working on the handoff, not existing** — the orchestrator
   learns that the session actually received the seed prompt and began acting on it, rather than
   assuming it did. An opened terminal is not a started session, and neither is one holding the seed
   prompt in its input: text nobody submitted is a session that never began. Whatever opens it must
   therefore wait until the session is ready to receive input before delivering the prompt, and deliver
   it as a submitted message rather than typed text left standing.
5. **It opens on the tree the launching session is working on, and creates none** — the change workspace
   when change-level isolation is active for the current change, the ordinary working tree otherwise. A
   launch that produces its own worktree or checkout does not satisfy this.

**How** those five are achieved is resolved at the moment of use, against whatever the environment
offers — it is fixed in no file. This mirrors the rule this ecosystem already applies to its exploration
index (`payload/domains/development/CLAUDE.md:49`): reference the capability, never hardcode the tool
that provides it, resolve it at use time.

**Why inheriting the directory is safe, instead of assumed.** Two sessions standing in one working tree
would be a real problem if both wrote — that is exactly the collision the change-workspace mechanism
exists to prevent. It is not a problem here because the side session only discusses: it reads and
reasons, and its single output is an Engram key. Inheriting is therefore not a shortcut around
isolation — it is what the read-only boundary buys. It also gives the discussion something a separate
checkout would have taken away: the side session sees the main thread's uncommitted work, because it is
standing in the same tree.

Whatever opens the session is used **only** to open it. If it also offers ways to read or drive the
session's terminal, the orchestrator does not use them: the conclusion is retrieved from the artifact
store, never by reading over the discussion's shoulder (see "Pickup" below).

**In the same act, the orchestrator starts a watch that wakes it when the conclusion lands.** A
detached command that waits for `side-discussion/{slug}/conclusion` to exist and exits once it does:
the harness re-invokes the orchestrator when a detached command finishes, so the arrival of the
conclusion becomes the wake-up, and the user never has to announce it. Poll at an interval matched to a
conversation, not to a machine — this waits on people talking, so tens of seconds, never one. If the
environment offers no way to run a detached command that outlives the turn, the watch is simply absent:
the mechanism still works, retrieval falls back to the two consulting moments in "Pickup", and the
orchestrator says once that it will not notice on its own.

The seed prompt carried into the new session is the only part of the launch fixed here, and it carries
**one** thing — where the handoff is:

```
Read the Engram observation with topic_key side-discussion/{slug}/handoff and follow it.
```

Nothing else goes in it: not which files to read, not the protocol's location, not a summary of the
topic. Everything the session needs is in the handoff, which opens by pointing at this file. A prompt
that names a second thing has started carrying context outside the artifact, and the next thing it
carries is the discussion itself.

No permission mode and no tool restriction are part of the launch — see "The write boundary" below for
why that is deliberate, not an oversight.

**When nothing in the environment resolves into a launch meeting all five**, the mechanism is
unavailable: the orchestrator says so in one line, offers to take the question into this thread instead —
naming the cost, that doing so spends the context the side discussion existed to save — and hands over
**no command**. It states once, in that same line, that a launcher would restore the mechanism, without
naming one, and does not proceed as if the discussion had happened.

## The side session's protocol

1. Read the handoff at the topic_key the seed prompt names — it opens by telling you to read this
   file, which is the whole of your protocol; nothing else configures you.
2. Read whatever `## References` points at, and whatever `## What I already read` says was already
   settled, before reasoning about `## Open question`.
3. **Discuss it with the user.** This is the step the whole mechanism exists for, and it is not a
   preamble to writing something: you put your reading to them, they push back, and the answer is what
   comes out of that exchange. You read and you reason. You do **not** edit, create, or delete any file
   in the repo, and you do **not** commit anything — not as a shortcut, not as a demonstration, not
   even when the answer would be obvious to write down as code.
4. **Write nothing until the user says what the conclusion is.** Reaching an answer on your own and
   writing it down is the failure this mechanism exists to avoid — it makes the session a delegated
   analysis, which needs no separate session at all. Your own reading, however well argued, is an input
   to the discussion, never its outcome.
5. Once the user states the conclusion — or approves one you offered — write it to
   `side-discussion/{slug}/conclusion` (see "The conclusion" below for its shape), as **they** settled
   it, not as you would have. That write is your only output, and it is the last thing you do.

## The write boundary

The launch carries no permission mode and no tool restriction. The "only discusses" limit lives
**only as prose**, in the handoff's `## Return` and in this file's own instructions — nothing in this
ecosystem enforces it mechanically here.

This is deliberate, and it is not the same bet this ecosystem makes everywhere else it uses prose rules:
a side discussion is **interactive**, and the user can read and type into it. If it starts writing or
reaches for a shell command, whoever is watching the terminal sees it happen. That visibility is what
makes a prose-only boundary defensible here.

**The cost, stated plainly, and now larger than a manually-launched session's.** If the user leaves the
session unattended, nothing stops it from writing or committing. Because the side terminal inherits the
launching session's working directory instead of a checkout of its own, a stray write lands **in the tree
the main thread is working in**, not in a discardable copy. And because the launch is automatic, the user
never has to touch the terminal for the session to start — "attended" is an expectation the mechanism
does not itself demonstrate. The boundary is only as strong as the attention behind it.

## The conclusion

The side session writes its conclusion to `side-discussion/{slug}/conclusion`, in five fixed sections,
all filled in on every write:

```markdown
## Question
{the open question, restated so this reads alone}

## Conclusion
{the answer, in the terms the question asked}

## Why
{the reasoning, including the alternatives weighed and why each was discarded}

## What it rests on
- {path}[:line] or {artifact-key} — {what it contributed}

## Still open
{what this could not settle — "None." when nothing}
```

The section split is not arbitrary: `## Conclusion` is already a decision item's `summary`, `## Why` is
its `rationale`, and `## What it rests on` is its `anchor` — the shape a conclusion needs to enter a
domain's existing decision-capture path with no rewriting, when it turns out to settle one.

**A conclusion is working material, never a decision record.** It is the same kind of thing as any other
phase's reasoning artifact — the reasoning that leads to a decision, not the decision once ratified. When
a conclusion settles something that belongs in a durable decision record, it re-enters whatever path this
project already uses to propose and ratify one: it does not get written there directly, by either side.
The side session never writes a decision-record file itself, and the main thread never copies a
conclusion into one without going through that path's own ratification step.

## Pickup

A discussion is **blocking** or **consultive**, and the user says which when they open it — the
orchestrator asks if they do not say, and never picks one on its own. Whichever they pick is declared in
the handoff's `Type` line.

**What each type means for the main thread is stated once, in the kernel's `### Side Discussion
(opt-in)` section — read it there.** This file does not restate it: a definition written in two places
drifts the first time one of them is edited, and the two homes exist precisely so that each says a
different thing. What belongs here is the part the kernel does not carry — how the conclusion is
retrieved, which is the same for both types.

Both types are picked up the same way: **by consulting** `side-discussion/{slug}/conclusion`. Three
moments bring the main thread to consult it — the watch started at open time fires because the
conclusion landed, the user reports the discussion done, or the main thread reaches a point where it
needs the conclusion. The first makes the other two unnecessary in the ordinary case; they remain
because a watch can be absent (see "Opening the session") and because the user may want the conclusion
read before the discussion formally closes.

What none of the three is: a signal from the side session itself. There is no push, webhook, or
cross-session channel between the two sessions, and the watch is not one — it observes the artifact
store, not the discussion, and it cannot tell a conclusion that was written from a discussion that was
abandoned. It knows one thing: the key exists now. Everything else is still read from the artifact.

**When the conclusion is not there yet**, the main thread says so and offers ways forward instead of
guessing at one: wait for the user to finish it, re-open the discussion with the same handoff, or bring
the open question into the main thread directly.

**Reading a consultive conclusion carries one extra step, before the conclusion itself: the main thread
states what it advanced while the discussion was open.** The work it did, named — not a claim that none
of it was affected. It does not judge whether any of it is now invalid; it puts what it did and what
the discussion concluded side by side, and the user decides. A blocking discussion needs no such step,
because by definition the thread advanced nothing that depended on it. This repairs nothing, and is not
meant to: it makes visible the overlap the consultive type accepts by construction, per the kernel's
test for choosing between the two.

**There is no timeout**, and none is to be added: nothing here causes the main thread to proceed on its
own reading of the open question after any amount of elapsed time. The reason is the kernel's standing
one about silence, which this file does not restate.

## Abandonment

If a side discussion is opened and never finished — the session opens and the user never brings it to a
conclusion — nothing cleans it up and nothing times it out. The handoff and (if written) the conclusion
simply stay in the artifact store under their keys, available if the same discussion is ever picked back
up. There is no liveness check to build here: the whole mechanism is consult-only.

## Store mode

This mechanism requires both the artifact store and a launch that satisfies the five properties above.
With `artifact_store.mode = none` there is no return channel for a handoff or a conclusion, and the
mechanism is unavailable for that reason — same branch, same message shape as an unresolved launch.
