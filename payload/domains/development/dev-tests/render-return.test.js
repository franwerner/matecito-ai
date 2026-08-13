// Harness for `../scripts/render-return.js`'s `renderItems()` — TODAY's engine, unmodified.
//
// This suite must go green against the current, unchanged script BEFORE any new render code is
// written (Phase 3 lifts `renderItems()`'s per-item logic — cap, tokens, rationale — into a shared
// helper the table and labeled-lists renderers also call). It is the regression net that proves the
// refactor didn't change observable behavior; it never touches `render-return.js` itself.
//
// Run with: node --test payload/domains/development/dev-tests/
'use strict';

const assert = require('node:assert/strict');
const { test, before, after } = require('node:test');

const { loadRenderReturnInternals, cleanupRenderReturnInternals } = require('./support.js');

let renderItems;

before(() => {
  ({ renderItems } = loadRenderReturnInternals());
});

after(() => {
  cleanupRenderReturnInternals();
});

// `renderItems(section, data)` — `section.field` names the list in `data`; `section.items` carries
// the optional `summary_max` cap, `tokens`/`token`, and `rationale` key. Fixtures below mirror the
// shapes the real `.yaml` contracts declare (see `sdd-verify.yaml`'s `## Decision Gaps`, etc.).

test('a plain item under the summary_max cap renders as a bare bullet', () => {
  const section = { field: 'items', title: '### T', items: { summary_max: 10 } };
  const out = renderItems(section, { items: ['short'] });
  assert.equal(out, '- short');
});

test('an item exactly AT the summary_max cap is accepted, not rejected', () => {
  const section = { field: 'items', title: '### T', items: { summary_max: 10 } };
  // Exactly 10 characters — the check is `length > summary_max`, so the boundary itself must pass.
  const out = renderItems(section, { items: ['1234567890'] });
  assert.equal(out, '- 1234567890');
});

test('an item OVER the summary_max cap fails the render, naming the section and the cap', () => {
  const section = { field: 'items', title: '### T', items: { summary_max: 10 } };
  assert.throws(
    () => renderItems(section, { items: ['12345678901'] }),
    /over the 10-char cap/
  );
});

test('a declared token missing from the item fails the render, naming the field', () => {
  const section = {
    field: 'items',
    title: '### Decision Gaps',
    items: { tokens: [{ name: 'anchor', field: 'anchor' }] },
  };
  assert.throws(
    () => renderItems(section, { items: [{ text: 'gap one' }] }),
    /has no `anchor`/
  );
});

test('a present token renders its `· name: value` line after the bullet', () => {
  const section = {
    field: 'items',
    title: '### Decision Gaps',
    items: { tokens: [{ name: 'anchor', field: 'anchor' }] },
  };
  const out = renderItems(section, { items: [{ text: 'gap one', anchor: '.matecito-ai/edr/foo.md' }] });
  assert.equal(out, '- gap one\n  · anchor: .matecito-ai/edr/foo.md');
});

test('a declaring section with the rationale field absent fails the render', () => {
  const section = { field: 'items', title: '### Unmandated Forks', items: { rationale: 'rationale' } };
  assert.throws(
    () => renderItems(section, { items: [{ text: 'fork one' }] }),
    /has no `rationale`/
  );
});

test('a rationale spanning multiple lines fails the render', () => {
  const section = { field: 'items', title: '### Unmandated Forks', items: { rationale: 'rationale' } };
  assert.throws(
    () => renderItems(section, { items: [{ text: 'fork one', rationale: 'line one\nline two' }] }),
    /spans multiple lines/
  );
});

test('an empty list renders the "None." sentinel, never an empty string', () => {
  const section = { field: 'items', title: '### T', items: {} };
  assert.equal(renderItems(section, { items: [] }), 'None.');
});

test('the rendered bullet count is derived from the list length, one per entry', () => {
  const section = { field: 'items', title: '### T', items: {} };
  const out = renderItems(section, { items: ['one', 'two', 'three'] });
  const bulletCount = out.split('\n').filter((l) => l.startsWith('- ')).length;
  assert.equal(bulletCount, 3);
});

test('an item carrying both a token and a rationale renders both adornment lines, in order', () => {
  const section = {
    field: 'items',
    title: '### Decision Gaps',
    items: { tokens: [{ name: 'anchor', field: 'anchor' }], rationale: 'rationale' },
  };
  const out = renderItems(section, {
    items: [{ text: 'gap one', anchor: '.matecito-ai/edr/foo.md', rationale: 'because it matters' }],
  });
  assert.equal(out, '- gap one\n  · anchor: .matecito-ai/edr/foo.md\n  · rationale: because it matters');
});
