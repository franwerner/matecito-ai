# SDD Init Details

## Testing Capability Checklist

- Test runner: `package.json` scripts/deps, `pyproject.toml`, `pytest.ini`, `go.mod`, `Cargo.toml`, `Makefile`.
- Test layers: unit runner; integration libraries (`testing-library`, `httpx`, `httptest`, `WebApplicationFactory`); E2E tools (`playwright`, `cypress`, `selenium`, `chromedp`).
- Coverage: `vitest --coverage`, `jest --coverage`, `c8`, `pytest-cov`, `go test -cover`, `coverlet`.
- Quality: linter, type checker, formatter commands.

## Skill Loading (matecito-ai)

<!-- matecito-ai: el escaneo y construcción del skill-registry fue removido. sdd-init ya no arma .atl/skill-registry.md. Las fases cargan su propia SKILL.md y leen las convenciones del proyecto directamente desde .matecito-ai/edr/, CLAUDE.md y config.yaml. -->

<!-- matecito-ai: this paragraph was in Spanish inside an English body. Translated without changing the
     instruction. -->
`sdd-init` builds no registry. Each SDD phase loads its own `SKILL.md`. The project's conventions are read straight from its own files: `.matecito-ai/edr/` (architecture decisions), `CLAUDE.md`, and `config.yaml`.

<!-- matecito-ai: two leftovers of the removed skill-registry were dropped here. (1) A dangling bullet
     about `AGENTS.md` that ended in "include both the index and referenced files in the registry" — a
     registry that no longer exists. (2) The whole `## LLM-First Skill Criteria` rubric, which was worse
     than orphaned: it declared a REQUIRED structure (Activation Contract / Hard Rules / Decision Gates /
     Execution Steps / Output Contract / References) that half the skills in this payload do not use —
     `sdd-intake`, `sdd-design`, `sdd-explore` and others are built on Purpose / What to Do / Rules. An
     executor taking the rubric literally concludes they are malformed. A rubric with no enforcer that
     invalidates the thing it describes is not documentation, it is a trap. -->

## Engram Saves

```text
mem_save title/topic_key: sdd-init/{project}
type: architecture
content: project context markdown, per `## Project Context Format` below
capture_prompt: false when available

mem_save title/topic_key: sdd/{project}/testing-capabilities
type: config
content: testing capabilities markdown
capture_prompt: false when available
<!-- matecito-ai: bloque mem_save del skill-registry removido -->
```

<!-- matecito-ai: sección OpenSpec Skeleton removida (engram-only) -->


<!-- matecito-ai: this artifact had NO declared format anywhere. The only production instruction was
     "content: detected project context markdown" — while its sibling, `testing-capabilities`, has a full
     format right below. So every init executor invented its own structure and the two consumers
     (`sdd-explore`, `sdd-propose`) went looking for whatever they imagined. The three axes below are the
     ones the return template already names in `### Project`; this is their full artifact version, of
     which that return carries only the one-line summary. -->
## Project Context Format

The shape of the `sdd-init/{project}` artifact. Read by `sdd-explore` and `sdd-propose` as project
context (optional for both), and by any later phase that needs the ground truth. Section titles are
fixed — a consumer looks for them literally.

This is the FULL detail. The `### Project` section of `~/.claude/references/phase-returns/sdd-init/sdd-init.md`
is the one-line-per-axis summary of these same three axes, for the orchestrator; the two never restate
each other.

```markdown
## Project Context

**Project**: {name}
**Detected**: {date}

### Stack

- Primary language: {the one language that governs the project}
- Languages: {all detected, with versions when the manifest pins them}
- Frameworks / runtimes: {name + version}
- Package manager: `{command}`
- Manifests read: {the files this detection came from}

### Architecture

- Layout: {single project | monorepo — and the apps/packages it contains}
- Entry points: {the concrete paths}
- Layers / boundaries: {how the code is divided, in the project's own vocabulary}

### Conventions

- What the project's own files fix: {`CLAUDE.md`, lint/format/type-check config, commit convention —
  cite the file for each, so a later phase can read the rule rather than trust this summary}
- Naming and structure patterns observed: {only what the code actually shows}
```

Detection rules for this artifact:

- **Primary language is one value, and other things depend on it.** `debugger.language` in the Testing
  Capabilities artifact derives from it, so it is recorded here explicitly and not left implicit in the
  languages list.
- **Never guess.** An axis you could not establish is reported as `— not detected`, never filled with a
  plausible value. Evidence that does not resolve to one answer is the `blocked` case of the return
  template, not a default you pick.
- **Cite the source.** Each axis names the files the detection came from. A later phase that disagrees
  needs to know what was read, and a stale entry needs to be diagnosable.
- **Ground truth only, never gate state.** Do NOT record whether the decision-record or capability-spec
  stores exist: those are presence-based gates that every phase evaluates for itself at use time, and a
  snapshot taken at init is stale the moment a store is created.

## Testing Capabilities Format

<!-- matecito-ai: consumers name dotted keys (`uiTest.available`, `test_runner.command`,
     `debugger.available`) and this format emitted prose labels (`- Command:`, an `**available**` row).
     The literal keys existed only in init's RETURN, which nobody retrieves: a reader looking for
     `uiTest.available` in the artifact found nothing and its gate closed silently. The artifact now
     emits the literal keys; this list is the contract between the two sides. -->
### Canonical keys — the contract with the consumers

Downstream phases cite these keys **literally**; the artefact below emits them literally. Neither
side paraphrases: a key that only exists under a prose label is a key nobody finds.

| Key | Emitted in | Read by |
| --- | --- | --- |
| `test_runner.command` | `### Test Runner` | `sdd-apply` (Strict TDD test execution), the orchestrator's Strict TDD forwarding |
| `test_runner.framework` | `### Test Runner` | `sdd-apply`, as context |
| `uiTest.proofshot` | `### UI Test` | `sdd-init` itself, to derive `uiTest.available` |
| `uiTest.devServer.command` | `### UI Test` | `sdd-verify` (the `proofshot start --run` argument) |
| `uiTest.available` | `### UI Test` | `sdd-verify` (the UI-test gate) |
| `debugger.language` | `### Debugger` | the `debugger` skill, to name the toolchain |
| `debugger.toolchain` | `### Debugger` | the `debugger` skill, so a ❌ is diagnosable |
| `debugger.available` | `### Debugger` | `sdd-apply` and `sdd-verify`, to decide whether a session can be opened |

```markdown
## Testing Capabilities

**Detected**: {date}

### Test Runner

- `test_runner.command`: `{command}`
- `test_runner.framework`: {name}

### Test Layers

| Layer       | Available | Tool        |
| ----------- | --------- | ----------- |
| Unit        | ✅ / ❌   | {tool or —} |
| Integration | ✅ / ❌   | {tool or —} |
| E2E         | ✅ / ❌   | {tool or —} |

### Coverage

- Available: ✅ / ❌
- Command: `{command or —}`

### Quality Tools

| Tool         | Available | Command        |
| ------------ | --------- | -------------- |
| Linter       | ✅ / ❌   | {command or —} |
| Type checker | ✅ / ❌   | {command or —} |
| Formatter    | ✅ / ❌   | {command or —} |

### UI Test

| Key                        | Available | Detail                                           |
| -------------------------- | --------- | ------------------------------------------------ |
| `uiTest.proofshot`         | ✅ / ❌   | `proofshot` binary found on PATH / not found     |
| `uiTest.devServer.command` | ✅ / ❌   | `{the resolved run command}` / —                 |
| `uiTest.available`         | ✅ / ❌   | proofshot AND devServer both ✅                  |

Detection notes:
- `uiTest.proofshot`: check `proofshot` on PATH (equivalent to `exec.LookPath("proofshot")`). If not on PATH at init time, detected as ❌ even if installed elsewhere. Document the limitation: proofshot installed outside PATH → detected absent.
- `uiTest.devServer.command`: inspect `package.json` scripts for `dev`, `start`, or `serve` keys (in that priority order); fall back to framework config (e.g. `vite.config.*`, `next.config.*`). Record the resolved command string in the Detail cell — `sdd-verify` passes it verbatim to `proofshot start --run`.
- `uiTest.available` = proofshot ✅ AND devServer ✅; both must be present for sdd-verify to run the UI step.

### Debugger

| Key                    | Available | Detail                                                    |
| ---------------------- | --------- | --------------------------------------------------------- |
| `debugger.language`    | —         | the project's primary language, as resolved in step 1      |
| `debugger.toolchain`   | ✅ / ❌   | the debug binary looked for, and whether it was found      |
| `debugger.available`   | ✅ / ❌   | the toolchain binary for that language is installed        |

Detection notes:
- **Procedural, never a lookup table.** Resolve the standard debug toolchain for the language you actually detected and check whether its binary is installed — `dlv` for Go, `debugpy` for Python, the Node inspector, and so on. Record WHICH binary you looked for, so a ❌ is diagnosable instead of mysterious.
- **Adapter present ≠ binary installed.** The MCP may report a language as supported while its debug binary is absent from the machine; what matters here is the binary.
- **❌ is a normal outcome, not a blocker.** `sdd-apply` and `sdd-verify` skip debugger usage silently when `available` is ❌; the `debugger` skill surfaces the install command if someone later needs a session.
- **Init detects, it does not provision.** Never install the toolchain here.
```

## Output Templates

<!-- matecito-ai: esto era una TERCERA descripción del retorno (además del template y de la skill).
     El formato vive en un solo lugar; acá queda sólo lo que el template no sabe: lo que cambia
     según el modo de persistencia. -->
The shape of the return lives in `~/.claude/references/phase-returns/sdd-init/sdd-init.md` — that file owns
the block and its sections; do not restate them here.

What this file adds, because it is mode-specific: **engram** mode must mention that the persisted
context is local and non-shareable; **none** mode must recommend enabling persistence.
