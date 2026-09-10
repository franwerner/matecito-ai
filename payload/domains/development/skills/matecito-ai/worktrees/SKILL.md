---
name: worktrees
description: Read before `git merge`, `git cherry-pick`, or any command that lands commits produced in another worktree onto a branch this session does not own. The test - does this write move work produced in another tree onto a branch you do not own? Yes, claim the turn first and release it after. No (a commit or push in the tree you are already in), this skill does not apply. USE THIS SKILL to claim and release that turn with `matecito-ai turn claim`/`release` - two commands around one native git ref, no interception, no identity - and what a refused claim means.
---

# Worktrees — the shared-branch turn

Concurrent orchestrator sessions, each driving its own change in the same repository, can reach the same
shared branch at the same time. This skill is the entry point to a turn-taking mechanism that keeps two
such writes from landing at once: two commands, `matecito-ai turn claim` and `matecito-ai turn release`,
around one git ref. Nothing intercepts a command, nothing inspects a command string, and nothing anywhere
asks who you are — see `## One claim, one release, per unit of work` for the one thing that makes that
safe.

## The scope criterion

Apply this test to the write in front of you, verbatim:

> **Does this write move work produced in another tree onto a branch this session does not own?**

Yes → the turn matters for this write. No (a commit or push in the tree you are already in) → no turn,
and nothing about this mechanism applies. The test is answered by comparing two things you already know
when you reach the write — the tree the work came from, and whether this session owns the destination
branch — never by whether the work is finished, which is not what the test asks.

## One claim, one release, per unit of work

**Claim once, do one unit of work, release on every path out — and never claim a second time inside that
bracket.** With no identity anywhere, the mechanism cannot tell a caller that already holds the turn from
any other caller: a second claim inside your own bracket is refused exactly like anyone else's would be.
That refusal is legible rather than decided — the report names your own change, tree and moment, which
makes the mistake obvious to a person reading it — but the mechanism itself never acts on that
coincidence. Both write-moments this mechanism guards are single units: the change integration is one
merge (claim once, then the whole merge sequence including its rebase-and-retry path, then release on
whichever path it ends on), and a consolidation round is one cherry-pick loop (claim once before the
first cherry-pick, release once after the last, before the cleanup pass, whatever the loop's outcome).

## Claiming the turn

```
matecito-ai turn claim --change <name> --destination <branch> --moment <change-integration|batch-consolidation>
```

`--tree` defaults to the current working tree; every flag is descriptive — none of it is ever compared
against anything (see `contracts/turn-release-keyed-to-the-claim-token.md`). Three outcomes:

- **Free — claimed.** Exit 0, and the command prints `token: <sha>` on its own line, after `turn
  claimed`. That sha is your receipt: hold onto it, you need it to release. Nothing else is printed.
- **Held.** Exit 1, with a report naming who holds it, since when, and what else is queued. The command
  already queued your request, paused, and retried once before reporting this — there is nothing left
  for you to do by hand. See `## When the turn is taken` below.
- **Not arbitrated.** Exit 0, with a message on stderr saying the mechanism could not decide (no
  repository, an unreadable record, a claim that failed for some other reason). The work proceeds; the
  failure is never silent.

Exit 1 is the enforcement — chain the claim ahead of the write with `&&`, so the write never runs behind
a refused claim:

```
matecito-ai turn claim --change my-change --destination main --moment change-integration \
  && git merge --ff-only matecito-ai/my-change
```

## Releasing the turn

```
matecito-ai turn release --token <sha>
```

Release on **every** path out of the unit of work the claim started — a clean write and a failed one
alike. `--token` is mandatory; there is deliberately no other way to identify yourself to `release`.

- **Matches.** The turn no longer exists, exit 0.
- **Mismatch.** The turn no longer carries that value — moved, or already released. Exit 0, with the
  mismatch reported on stderr; nothing is deleted, nothing is forced.
- **No token.** Exit 1, a usage error. There is no fallback that identifies you any other way.

## When the turn is taken

A claim that exits 1 already performed the queue-and-retry
`structure/guard-queues-and-retries-agent-asks.md` requires, inside that one call. What is left to you is
the conversation: read the three facts the report names — holder, age, queue — and put them to the user.
Offer a recommendation and choose none of the options yourself: wait and retry, leave it, or release the
turn under the user's explicit instruction. **No caller breaks another caller's turn on its own
initiative, whatever the evidence** — not after a wait, not on a stale timestamp, not on an absent
holder. A headless run has no one to ask: it integrates nothing and reports itself `status: blocked` with
the same three facts, leaving its work untouched for a later attempt.

## Reading the state

```
matecito-ai turn status
```

Always succeeds — free, or the holder with its age and its queue — and never claims or releases anything.
Use it whenever you want to look without attempting a claim: a human diagnosing, an orchestrator asked
why something is stuck. Under the hood this reads the same two things a raw `git` inspection would:
`refs/matecito-ai/` holds the lock (if any) and every queued entry, one `git for-each-ref` away; a single
entry's fields are `git cat-file -p <ref>`. A queue entry is a claimed pending request, never
confirmed-live — its owner may already be gone, and nothing checks. Report it as **claimed**, never as
evidence the owner is still working.

## When matecito-ai is not installed

Everything below applies only on a machine where the `matecito-ai` binary is not available — the two
commands above cover every other case. Do not write to the shared branch, do not delete or overwrite the
lock ref, and do not force a way past it by any other route.

**Claiming by hand.**

1. Build the record as plain text, one field per line: `change` · `tree` (the absolute path of the
   working tree this write comes from) · `destination` · `moment` (`change-integration` or
   `batch-consolidation`) · `base` (the sha this write applies to) · `at` (an ISO-8601 UTC timestamp).
2. `git hash-object -w --stdin` over that record → a blob sha.
3. `git update-ref refs/matecito-ai/merge-lock <blob-sha> ""` — the empty old value makes this a
   create-only-if-absent write.

Exit 0: the turn is yours — proceed to the write. Exit 128, `fatal: … reference already exists`: the
turn is held — go to "If it is already taken" below. **This record carries no receipt the tool issued.**
It is exactly as valid as a record `matecito-ai turn claim` would have written, but nothing distinguishes
it from one written by mistake except that a human built it by hand — the same trust the mechanism
withdrew everywhere else once it stopped depending on anyone's word. Weigh that before choosing this
path over installing the tool.

**Releasing by hand**, once the write is over, on every path out:

```
git update-ref -d refs/matecito-ai/merge-lock <blob-sha>
```

Exit 0: released. A failure means the ref no longer carries `<blob-sha>` — report the mismatch and stop;
do not retry with a different value, do not delete unconditionally, and do not treat the failure as
harmless.

**If it is already taken:**

1. Build your own record as above (`moment`/`base` describe the write you are waiting to make), write it
   as a blob, and point a queue ref at it: `git update-ref refs/matecito-ai/merge-queue/<change-name>
   <blob-sha>`.
2. Pause briefly, then re-attempt the claim **once**.
3. Still taken → stop trying. Read the state (`## Reading the state` above, or the raw refs) and report
   who holds it, since when, and what else is queued — same conversation as `## When the turn is taken`
   above.
4. Once your claim succeeds, delete your own queue ref by compare-and-swap on the value you wrote:
   `git update-ref -d refs/matecito-ai/merge-queue/<change-name> <blob-sha>`.

**No session breaks another session's turn on its own initiative, whatever the evidence.** The only path
by which a turn is released by someone other than its holder is the user saying so explicitly. Accepted
cost: a crashed session blocks the queue until a human intervenes.

## Worktree mechanics live elsewhere

This skill owns the turn — what it is, how to claim and release it, why a claim comes back refused, and
how to do either by hand where the tool is not installed. It does not own how the per-task worktree a
parallel batch opens is created, how the consolidation run's cherry-pick loop runs, or how either is
cleaned up afterward. For that:

- `~/.claude/references/phase-returns/sdd-apply/parallel-batch.md` — per-task worktrees, the base
  handshake, the consolidation run, and its cleanup pass.
