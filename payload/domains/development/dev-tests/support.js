// Shared test support for payload/domains/development/dev-tests/.
//
// Lives in `dev-tests/`, a sibling of `scripts/` inside the domain folder that deploy never walks —
// `internal/setup/deploy/deploy.go`'s `domainComponents` only globs 5 literal names per component root
// (`agents`, `skills`, `references`, `scripts`, `matecito-ai`; see `componentRoots` at deploy.go:781),
// so a directory named anything else is invisible to `buildMappings` and never reaches `~/.claude/`.
// Confirmed empirically (not just assumed from the name) via `matecito-ai install --dry-run` and a
// direct `deploy.Plan()` probe — see the merged `sdd/stores-cerrados-bajo-referencias/apply-progress`
// artifact in Engram for both runs; no Go test lives in this repo for that check, since it is a
// property of `internal/setup/deploy`, not of this JS suite.
//
// `validate-artifact.js` ends in an unconditional `main()` call (no `require.main === module` guard),
// so it cannot be `require()`-d directly without running the CLI and calling `process.exit()`. This
// builds a requireable copy in a throwaway temp directory instead of adding exports to the shipped
// script — the same approach used to verify the fix this suite guards against a regression of.
'use strict';

const fs = require('fs');
const os = require('os');
const path = require('path');

const SCRIPTS_DIR = path.join(__dirname, '..', 'scripts');
const SRC = path.join(SCRIPTS_DIR, 'validate-artifact.js');
const LIB_DIR = path.join(SCRIPTS_DIR, 'lib');
const CHECKS_FILE = path.join(__dirname, '..', 'references', 'artifact-checks', 'checks.yaml');

let cached = null;
let cachedLoadDir = null;

// Returns `{ resolveStoreLink, fallsUnderStoreRoot, ENGINE, checksConfig }` from the shipped script,
// without touching it and without adding exports to production code. Memoized per process.
function loadInternals() {
  if (cached) return cached;

  const yaml = require(path.join(LIB_DIR, 'yaml.js'));
  const checksConfig = yaml.parse(fs.readFileSync(CHECKS_FILE, 'utf8'));

  let src = fs.readFileSync(SRC, 'utf8');
  src = src.replace("require('./lib/yaml')", `require(${JSON.stringify(path.join(LIB_DIR, 'yaml.js'))})`);
  src = src.replace("require('./lib/store')", `require(${JSON.stringify(path.join(LIB_DIR, 'store.js'))})`);
  src = src.replace(/\nmain\(\);\n$/, '\n');
  src += '\nmodule.exports = { resolveStoreLink, fallsUnderStoreRoot, ENGINE, mkFinding };\n';

  cachedLoadDir = fs.mkdtempSync(path.join(os.tmpdir(), 'matecito-dev-tests-load-'));
  const tmpFile = path.join(cachedLoadDir, 'validate-artifact.testable.js');
  fs.writeFileSync(tmpFile, src);

  const internals = require(tmpFile);
  cached = { ...internals, checksConfig };
  return cached;
}

// Removes the throwaway `require()`-able copy `loadInternals` wrote to the OS temp dir. Call once,
// after all tests in the process have run (a global `after` hook, not a per-test one).
function cleanupInternals() {
  if (cachedLoadDir) removeFixtureRepo(cachedLoadDir);
  cachedLoadDir = null;
  cached = null;
}

let cachedRenderReturn = null;
let cachedRenderReturnLoadDir = null;
const RENDER_RETURN_SRC = path.join(SCRIPTS_DIR, 'render-return.js');

// Same trick as `loadInternals` above, for `render-return.js` — it also ends in an unconditional
// `main()` call, so it cannot be `require()`-d directly either. Additionally, `render-return.js`'s
// `fail()` calls `process.exit(1)` from WITHIN `renderItems()` itself (not only from `main()`), so
// calling it directly with bad data would kill the `node --test` process instead of failing one
// assertion. `fail()`'s body is replaced here with a `throw` so a rejected render surfaces as a
// catchable exception (`assert.throws`) — everything else in the module is untouched. Returns
// `{ renderItems }`, memoized per process.
function loadRenderReturnInternals() {
  if (cachedRenderReturn) return cachedRenderReturn;

  let src = fs.readFileSync(RENDER_RETURN_SRC, 'utf8');
  src = src.replace("require('./lib/yaml')", `require(${JSON.stringify(path.join(LIB_DIR, 'yaml.js'))})`);
  src = src.replace(
    'function fail(msg) {\n  process.stderr.write(`${msg}\\n`);\n  process.exit(1);\n}',
    'function fail(msg) {\n  throw new Error(msg);\n}'
  );
  src = src.replace(/\nmain\(\);\n$/, '\n');
  src += '\nmodule.exports = { renderItems, renderTable, renderLabeledLists, tokensOf, derive, fail };\n';

  cachedRenderReturnLoadDir = fs.mkdtempSync(path.join(os.tmpdir(), 'matecito-dev-tests-load-'));
  const tmpFile = path.join(cachedRenderReturnLoadDir, 'render-return.testable.js');
  fs.writeFileSync(tmpFile, src);

  cachedRenderReturn = require(tmpFile);
  return cachedRenderReturn;
}

// Removes the throwaway `require()`-able copy `loadRenderReturnInternals` wrote to the OS temp dir.
// Call once, after all tests in the process have run (a global `after` hook, not a per-test one).
function cleanupRenderReturnInternals() {
  if (cachedRenderReturnLoadDir) removeFixtureRepo(cachedRenderReturnLoadDir);
  cachedRenderReturnLoadDir = null;
  cachedRenderReturn = null;
}

let cachedValidateReturn = null;
let cachedValidateReturnLoadDir = null;
const VALIDATE_RETURN_SRC = path.join(SCRIPTS_DIR, 'validate-return.js');

// Same trick as `loadInternals` above, for `validate-return.js` — it also ends in an unconditional
// `main()` call, so it cannot be `require()`-d directly either. Returns `{ validate, checkItems,
// tokensOf, parseReturn, isSentinel }`, memoized per process.
function loadValidateReturnInternals() {
  if (cachedValidateReturn) return cachedValidateReturn;

  let src = fs.readFileSync(VALIDATE_RETURN_SRC, 'utf8');
  src = src.replace("require('./lib/yaml')", `require(${JSON.stringify(path.join(LIB_DIR, 'yaml.js'))})`);
  src = src.replace(/\nmain\(\);\n$/, '\n');
  src += '\nmodule.exports = { validate, checkItems, tokensOf, parseReturn, isSentinel };\n';

  cachedValidateReturnLoadDir = fs.mkdtempSync(path.join(os.tmpdir(), 'matecito-dev-tests-load-'));
  const tmpFile = path.join(cachedValidateReturnLoadDir, 'validate-return.testable.js');
  fs.writeFileSync(tmpFile, src);

  cachedValidateReturn = require(tmpFile);
  return cachedValidateReturn;
}

// Removes the throwaway `require()`-able copy `loadValidateReturnInternals` wrote to the OS temp dir.
// Call once, after all tests in the process have run (a global `after` hook, not a per-test one).
function cleanupValidateReturnInternals() {
  if (cachedValidateReturnLoadDir) removeFixtureRepo(cachedValidateReturnLoadDir);
  cachedValidateReturnLoadDir = null;
  cachedValidateReturn = null;
}

// A throwaway two-store fixture tree, isolated per test — never the repo's real `.matecito-ai/`.
// Layout:
//   <root>/.matecito-ai/development-specs/rule/test-rule.md   (in-store, links out via `../edr/...`)
//   <root>/.matecito-ai/development-specs/flow/other-flow.md  (in-store link target)
//   <root>/.matecito-ai/edr/security/auth.md                  (cross-store link target)
//   <root>/apps/api/.matecito-ai/edr/data/nested.md            (per-sub-app nested EDR store)
//   <root>/apps/api/.matecito-ai/edr/data/nested-sibling.md    (intra-nested-store link target)
function makeFixtureRepo() {
  const root = fs.mkdtempSync(path.join(os.tmpdir(), 'matecito-dev-tests-repo-'));

  const specRule = path.join(root, '.matecito-ai', 'development-specs', 'rule');
  const specFlow = path.join(root, '.matecito-ai', 'development-specs', 'flow');
  const edrSecurity = path.join(root, '.matecito-ai', 'edr', 'security');
  const nestedEdrData = path.join(root, 'apps', 'api', '.matecito-ai', 'edr', 'data');

  fs.mkdirSync(specRule, { recursive: true });
  fs.mkdirSync(specFlow, { recursive: true });
  fs.mkdirSync(edrSecurity, { recursive: true });
  fs.mkdirSync(nestedEdrData, { recursive: true });

  fs.writeFileSync(
    path.join(specFlow, 'other-flow.md'),
    '# Flow — other\n\nFixture target for an intra-store link.\n'
  );
  fs.writeFileSync(
    path.join(edrSecurity, 'auth.md'),
    '# EDR — auth\n\nFixture target for a cross-store link.\n'
  );
  fs.writeFileSync(
    path.join(nestedEdrData, 'nested-sibling.md'),
    '# EDR — nested sibling\n\nFixture target for an intra-nested-store link.\n'
  );

  // The real `rule/event-scoping.md:26` shape: an inline cross-store link in the BODY, outside any
  // section `dangling-link`/`scan: section` would inspect, plus a `## Referencias` section that only
  // ever carries intra-store/dangling/external links (the store is closed — see spec/README.md).
  fs.writeFileSync(
    path.join(specRule, 'test-rule.md'),
    [
      '# Rule — test',
      '',
      '## Reglas de negocio',
      '',
      'Cross-store link inline in the body: [auth](../../edr/security/auth.md).',
      '',
      '## Referencias',
      '',
      '- intra-store → [other-flow](../flow/other-flow.md)',
      '- colgado → [missing](../flow/does-not-exist.md)',
      '- externo → [example](https://example.com)',
      '',
    ].join('\n')
  );

  return root;
}

function removeFixtureRepo(root) {
  fs.rmSync(root, { recursive: true, force: true });
}

// Creates a symlink at `<repo>/.matecito-ai/development-specs/flow/symlinked-out.md` whose REAL
// target is a file that lives OUTSIDE `repo` entirely (a separate temp dir) — the shape of a link
// that textually resolves under the store's own root but physically points elsewhere, which is what
// exposed the purely-textual `path.resolve`/`path.relative` comparison this fixture guards against.
//
// Returns `{ externalDir, linkPath }` on success. Returns `null` — never throws — when the platform
// refuses to create a symlink (e.g. missing privilege); callers MUST check for `null` and skip their
// test with a visible reason instead of treating it as a failure.
function makeExternalSymlink(repo) {
  const externalDir = fs.mkdtempSync(path.join(os.tmpdir(), 'matecito-dev-tests-external-'));
  const externalTarget = path.join(externalDir, 'outside.md');
  fs.writeFileSync(externalTarget, '# Outside\n\nReal file living OUTSIDE the fixture repo entirely.\n');

  const linkPath = path.join(repo, '.matecito-ai', 'development-specs', 'flow', 'symlinked-out.md');
  try {
    fs.symlinkSync(externalTarget, linkPath, 'file');
  } catch (e) {
    removeFixtureRepo(externalDir);
    return null;
  }
  return { externalDir, linkPath };
}

module.exports = {
  loadInternals,
  cleanupInternals,
  loadRenderReturnInternals,
  cleanupRenderReturnInternals,
  loadValidateReturnInternals,
  cleanupValidateReturnInternals,
  makeFixtureRepo,
  removeFixtureRepo,
  makeExternalSymlink,
};
