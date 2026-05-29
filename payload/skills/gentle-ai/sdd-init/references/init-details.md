# SDD Init Details

## Testing Capability Checklist

- Test runner: `package.json` scripts/deps, `pyproject.toml`, `pytest.ini`, `go.mod`, `Cargo.toml`, `Makefile`.
- Test layers: unit runner; integration libraries (`testing-library`, `httpx`, `httptest`, `WebApplicationFactory`); E2E tools (`playwright`, `cypress`, `selenium`, `chromedp`).
- Coverage: `vitest --coverage`, `jest --coverage`, `c8`, `pytest-cov`, `go test -cover`, `coverlet`.
- Quality: linter, type checker, formatter commands.

## Skill Loading (matecito-ai)

<!-- matecito-ai: el escaneo y construcción del skill-registry fue removido. sdd-init ya no arma .atl/skill-registry.md. Las fases cargan su propia SKILL.md y leen las convenciones del proyecto directamente desde .matecito-ai/adr/, CLAUDE.md y config.yaml. -->

sdd-init no construye ningún registry. Cada fase SDD carga su propia `SKILL.md`. Las convenciones del proyecto se leen directamente de sus archivos: `.matecito-ai/adr/` (decisiones de arquitectura), `CLAUDE.md`, y `config.yaml`.
- For index files such as `AGENTS.md`, extract referenced file paths and include both the index and referenced files in the registry.

## LLM-First Skill Criteria

- Treat skills as runtime instruction contracts, not human documentation.
- Required structure: frontmatter, Activation Contract, Hard Rules, Decision Gates, Execution Steps, Output Contract, References.
- Keep `description` quoted, one physical line, trigger-first, and no longer than 250 characters.
- Target 180-450 body tokens; move examples, schemas, edge cases, and background into local `references/` or `assets/`.
- References must be local files and stable relative to the skill directory when possible.
- Quality gates: hard rules are observable, decision gates cover real forks, output contract states exactly what to return, and references resolve locally.

## Engram Saves

```text
mem_save title/topic_key: sdd-init/{project}
type: architecture
content: detected project context markdown
capture_prompt: false when available

mem_save title/topic_key: sdd/{project}/testing-capabilities
type: config
content: testing capabilities markdown
capture_prompt: false when available
<!-- matecito-ai: bloque mem_save del skill-registry removido -->
```

<!-- matecito-ai: sección OpenSpec Skeleton removida (engram-only) -->


## Testing Capabilities Format

```markdown
## Testing Capabilities

**Strict TDD Mode**: {enabled/disabled}
**Detected**: {date}

### Test Runner

- Command: `{command}`
- Framework: {name}

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
```

## Output Templates

For each mode, include project, stack, persistence, Strict TDD Mode, Testing Capabilities table, artifacts created/saved, limitations where relevant, and next steps. Engram mode must mention local/non-shareable limitations; none mode must recommend enabling persistence.
