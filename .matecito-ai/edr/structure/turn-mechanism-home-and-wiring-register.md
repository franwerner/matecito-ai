# EDR — The turn mechanism lives entirely in the `worktrees` skill, and every wiring point opens on the git command

- **Status:** Accepted
- **Date:** 2026-08-15

## Contexto
The turn-taking mechanism shipped as a discoverable skill plus a separate reference file, wired with a paragraph that opened by naming the acting role and the moment in the cycle. Two real Claude Code sessions were offered the skill, both read the wiring paragraph end to end, and neither claimed the turn — merging `--no-ff` 23 seconds apart without either consulting the mechanism. In the same file, instructions written as concrete git mechanics (`git worktree add <dir> -b matecito-ai/<change-name> <original-branch>`; confirming the `.gitignore` line) were followed by both sessions, unprompted. Two Accepted records — `structure/shared-branch-turn-prose-homes.md` and `structure/shared-branch-turn-wiring-points.md` — fix the old split (a `shared-branch-turn` skill and reference pair, wiring paragraphs with an unconstrained register); this design contradicts both, so this record replaces them rather than editing either — the route chosen at the intake gate and re-confirmed as a spec requirement.

## Decisión
The turn mechanism gets one home and one name — the `worktrees` skill; no separate definition file ships, and each wiring point opens on the git command being run instead of on who runs it. This record replaces both `shared-branch-turn-prose-homes.md` and `structure/shared-branch-turn-wiring-points.md`; neither is edited. The folder name `worktrees` was checked unique against all 48 skill folders the two domains ship — `workspace` was rejected as the flow's own word (`change workspace`), `merge-lock` as naming a merge, when the mechanism's only remaining write-moment is a cherry-pick loop, which is no merge. The criterion lives in the skill body in full and in the frontmatter `description` in compact form — no third home.

## Alcance
- `payload/domains/development/skills/matecito-ai/worktrees/**` — The renamed skill — the mechanism's only home.

## Reglas verificables
- **[manual]** Exactly one folder ships the mechanism, named `worktrees`, globally unique across every installed domain's skill folders, and matching its frontmatter `name`.
- **[manual]** No separate definition file ships for this mechanism, and no shipped instruction cites one — the decision-record store is excluded, since its citations of the old paths are historical record.
- **[manual]** Each wiring point's opening sentence states what to do around the git command being typed, never who the reader is or which phase they are in.
- **[manual]** The scope criterion appears in the skill body in full, with its yes/no consequences, and in compact form in the frontmatter `description` — nowhere else.
- **[manual]** Worktree, nesting, consolidation and cleanup mechanics are cited from their owners (the domain fragment's Shared-branch turn pointer, `parallel-batch.md`) and never restated inside the skill.

## Alternativas consideradas
`workspace` as the base name — rejected: it is the flow's own word (`change workspace`) and would read as the ceremony noun that already failed to be recognized. `merge-lock` — rejected, with a strengthened reason now that only one write-moment remains: the mechanism's one write-moment is a cherry-pick loop, which is no merge, so the name would describe something that no longer happens. Keeping `shared-branch-turn` and rewording only the wiring paragraphs — rejected: `shared-branch-turn` names a ceremony, and the register of that name is the same register that failed.

## Consecuencias
The mechanism's whole definition lives in one place a reader already has in front of them when a skill listing is shown, and the compact criterion in the frontmatter description lets an agent decide whether it needs the skill without opening it. The evidence this answers to is stated, not assumed: every rule above is a documentary check, checkable by reading the payload with no run at all, and every one of them can pass while the mechanism still never fires at a real integration — a live run (two concurrent consolidation runs reaching the same cherry-pick loop at once) is the only confirmation that recognition actually improved.
