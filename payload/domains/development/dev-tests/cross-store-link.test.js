// Regression coverage for the store-closure engine in `../scripts/validate-artifact.js`
// (`resolveStoreLink`, `fallsUnderStoreRoot`, and the `cross-store-link` ENGINE entry).
//
// Run with: node --test payload/domains/development/dev-tests/
//
// Own fixtures throughout (see support.js) — never the repo's real `.matecito-ai/` store.
'use strict';

const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const { test, before, after } = require('node:test');

const { loadInternals, cleanupInternals, makeFixtureRepo, removeFixtureRepo, makeExternalSymlink } = require('./support.js');

let internals;
let repoRoot;
let siblingRoot;

before(() => {
  internals = loadInternals();
  repoRoot = makeFixtureRepo();
  // A second, unrelated fixture "repo" — simulates a sibling project's own `.matecito-ai/...` store
  // living elsewhere on the same machine (the shape of this batch's CRITICAL finding).
  siblingRoot = makeFixtureRepo();
});

after(() => {
  removeFixtureRepo(repoRoot);
  removeFixtureRepo(siblingRoot);
  cleanupInternals();
});

function capabilitySpecCtx(root) {
  const storeDir = path.join(root, '.matecito-ai', 'development-specs');
  return {
    root,
    ownBase: path.basename(storeDir),
    commonParent: path.dirname(storeDir),
    stores: internals.checksConfig.stores,
  };
}

// --- resolveStoreLink: the 4 verdicts ---------------------------------------------------------

test('resolveStoreLink: intra-store link resolves ok', () => {
  const ctx = capabilitySpecCtx(repoRoot);
  const artifact = { relPath: 'rule/test-rule.md' };
  const verdict = internals.resolveStoreLink(ctx, artifact, '../flow/other-flow.md', 'capability-spec');
  assert.equal(verdict, 'ok');
});

test('resolveStoreLink: link with no file at the resolved path is missing', () => {
  const ctx = capabilitySpecCtx(repoRoot);
  const artifact = { relPath: 'rule/test-rule.md' };
  const verdict = internals.resolveStoreLink(ctx, artifact, '../flow/does-not-exist.md', 'capability-spec');
  assert.equal(verdict, 'missing');
});

test('resolveStoreLink: link resolving to a real file outside the declared store is out-of-store', () => {
  const ctx = capabilitySpecCtx(repoRoot);
  const artifact = { relPath: 'rule/test-rule.md' };
  const verdict = internals.resolveStoreLink(ctx, artifact, '../../edr/security/auth.md', 'capability-spec');
  assert.equal(verdict, 'out-of-store');
});

test('resolveStoreLink: external href short-circuits to ok, never evaluated against any store', () => {
  const ctx = capabilitySpecCtx(repoRoot);
  const artifact = { relPath: 'rule/test-rule.md' };
  const verdict = internals.resolveStoreLink(ctx, artifact, 'https://example.com', 'capability-spec');
  assert.equal(verdict, 'ok');
});

// --- The CRITICAL this batch fixed: a link escaping into another repo entirely -----------------

test('CRITICAL regression: a link resolving into a SIBLING repo\'s own store is out-of-store, not ok', () => {
  // `commonParent` mimics the real CLI shape (dirname of `--store`), but `root` is THIS fixture repo
  // — the link below climbs enough `../` to land inside `siblingRoot`'s own `.matecito-ai/...` store,
  // a REAL file that exists on disk, just not under `repoRoot`. `href` is derived with `path.relative`
  // (never hand-counted `../`), so the test stays correct regardless of where the OS puts temp dirs.
  const storeDir = path.join(repoRoot, '.matecito-ai', 'development-specs');
  const ctx = {
    root: repoRoot,
    ownBase: path.basename(storeDir),
    commonParent: path.dirname(storeDir),
    stores: internals.checksConfig.stores,
  };
  const artifact = { relPath: 'rule/test-rule.md' };
  const fromDirAbs = path.join(storeDir, 'rule'); // the real directory holding rule/test-rule.md
  const targetAbs = path.join(siblingRoot, '.matecito-ai', 'development-specs', 'flow', 'other-flow.md');
  const href = path.relative(fromDirAbs, targetAbs).split(path.sep).join('/');

  const verdict = internals.resolveStoreLink(ctx, artifact, href, 'capability-spec');
  assert.equal(verdict, 'out-of-store');
});

// --- Symlink escape: the narrower instance of the same guarantee found on the 2nd verify run --------

test('CRITICAL regression #2: a symlink whose real target lives outside the repo is out-of-store', (t) => {
  // `fallsUnderStoreRoot` is purely textual (`path.resolve`/`path.relative`) unless it also resolves
  // through symlinks. A symlink placed INSIDE the store, pointing at a REAL file OUTSIDE the repo,
  // textually resolves under the store root even though its dereferenced target does not — the exact
  // gap `fs.realpathSync` closes. Own fixture symlink, not anything from the real filesystem.
  const created = makeExternalSymlink(repoRoot);
  if (!created) {
    t.skip('platform refused to create a symlink (fs.symlinkSync threw) — skipping, not failing');
    return;
  }
  const { externalDir, linkPath } = created;
  try {
    const storeDir = path.join(repoRoot, '.matecito-ai', 'development-specs');
    const ctx = {
      root: repoRoot,
      ownBase: path.basename(storeDir),
      commonParent: path.dirname(storeDir),
      stores: internals.checksConfig.stores,
    };
    const artifact = { relPath: 'rule/test-rule.md' };

    // Sanity check: the symlink's own textual path DOES land under the declared store — proving this
    // test actually exercises the realpath resolution, not just a link that was already out-of-store
    // on the text alone.
    assert.ok(linkPath.startsWith(storeDir), 'the symlink must textually resolve under the store root');
    assert.ok(fs.existsSync(linkPath), 'the symlink must dereference to a real, existing file');

    const verdict = internals.resolveStoreLink(ctx, artifact, '../flow/symlinked-out.md', 'capability-spec');
    assert.equal(verdict, 'out-of-store');
  } finally {
    removeFixtureRepo(externalDir);
  }
});

test('fallsUnderStoreRoot: exact 3 cases the verify report demonstrated', () => {
  const { fallsUnderStoreRoot } = internals;
  const sibling = path.join(siblingRoot, '.matecito-ai', 'development-specs', 'rule', 'test-rule.md');
  const own = path.join(repoRoot, '.matecito-ai', 'development-specs', 'rule', 'test-rule.md');
  const nested = path.join(repoRoot, 'apps', 'api', '.matecito-ai', 'edr', 'data', 'nested.md');

  assert.equal(fallsUnderStoreRoot(sibling, '.matecito-ai/development-specs', repoRoot), false, 'sibling repo -> out-of-store');
  assert.equal(fallsUnderStoreRoot(own, '.matecito-ai/development-specs', repoRoot), true, 'own repo root store -> ok');
  assert.equal(fallsUnderStoreRoot(nested, '.matecito-ai/edr', repoRoot), true, 'nested sub-app EDR -> ok');
});

// --- The nested sub-app case: the reason the segment-subsequence match exists at all -----------

test('nested per-sub-app EDR store still resolves ok (the original mandated departure)', () => {
  const storeDir = path.join(repoRoot, 'apps', 'api', '.matecito-ai', 'edr');
  const ctx = {
    root: repoRoot,
    ownBase: path.basename(storeDir),
    commonParent: path.dirname(storeDir),
    stores: internals.checksConfig.stores,
  };
  // A link from data/nested.md (the source need not exist — resolveStoreLink only checks the target)
  // to its real sibling in the SAME nested store, one directory level down from the repo root.
  const artifact = { relPath: 'data/nested.md' };
  const verdict = internals.resolveStoreLink(ctx, artifact, 'nested-sibling.md', 'edr');
  assert.equal(verdict, 'ok');
});

// --- scan: body engine — inside vs. outside the designated section -----------------------------

test('cross-store-link (scan: body): catches an inline cross-store link outside any declared section', () => {
  const storeDir = path.join(repoRoot, '.matecito-ai', 'development-specs');
  const relPath = 'rule/test-rule.md';
  const text = fs.readFileSync(path.join(storeDir, relPath), 'utf8');
  const artifact = { relPath, folder: 'rule', slug: 'test-rule', text, header: { Status: 'Accepted' } };
  const ctx = {
    root: repoRoot,
    ownBase: path.basename(storeDir),
    commonParent: path.dirname(storeDir),
    stores: internals.checksConfig.stores,
    artifacts: [artifact],
  };

  const row = internals.checksConfig.checks.find(
    (r) => r.kind === 'cross-store-link' && r.store === 'capability-spec'
  );
  assert.ok(row, 'checks.yaml must still declare the capability-spec cross-store-link row');
  assert.equal(row.scan, 'body', 'the row this suite guards must stay scan: body');

  const findings = internals.ENGINE['cross-store-link']([row], ctx);
  assert.equal(findings.length, 1);
  assert.match(findings[0].message, /\.\.\/\.\.\/edr\/security\/auth\.md/);
  assert.equal(findings[0].severity, 'warning');
});

test('negative control: a scan: section sweep of Referencias alone never sees the inline body link', () => {
  const storeDir = path.join(repoRoot, '.matecito-ai', 'development-specs');
  const text = fs.readFileSync(path.join(storeDir, 'rule', 'test-rule.md'), 'utf8');
  const libStore = require(path.join(__dirname, '..', 'scripts', 'lib', 'store.js'));
  const referenciasBody = libStore.extractSection(text, '## Referencias');
  const referenciasLinks = libStore.extractLinks(referenciasBody || '');
  assert.ok(!referenciasLinks.some((l) => l.href === '../../edr/security/auth.md'));
});

// --- EDR -> spec direction: no positive test existed for this row before this batch -------------

test('cross-store-link (edr side): an EDR linking a capability-spec is flagged', () => {
  const storeDir = path.join(repoRoot, '.matecito-ai', 'edr');
  const relPath = 'security/fake-edr.md';
  const text = [
    '# EDR — fake',
    '',
    '## Relacionados',
    '',
    '- `relacionado-con` → [../security/auth.md](../security/auth.md)',
    '- cita cruzada → [../../development-specs/rule/test-rule.md](../../development-specs/rule/test-rule.md)',
    '',
  ].join('\n');
  fs.writeFileSync(path.join(storeDir, relPath), text);

  try {
    const artifact = { relPath, folder: 'security', slug: 'fake-edr', text, header: { Status: 'Accepted' } };
    const ctx = {
      root: repoRoot,
      ownBase: path.basename(storeDir),
      commonParent: path.dirname(storeDir),
      stores: internals.checksConfig.stores,
      artifacts: [artifact],
    };
    const row = internals.checksConfig.checks.find((r) => r.kind === 'cross-store-link' && r.store === 'edr');
    assert.ok(row, 'checks.yaml must still declare the edr cross-store-link row');

    const findings = internals.ENGINE['cross-store-link']([row], ctx);
    assert.equal(findings.length, 1);
    assert.match(findings[0].message, /test-rule\.md/);

    const intraStoreOnly = findings.every((f) => !f.message.includes('auth.md'));
    assert.ok(intraStoreOnly, 'the intra-store Relacionados link must not also be flagged');
  } finally {
    fs.unlinkSync(path.join(storeDir, relPath));
  }
});
