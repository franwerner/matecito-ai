// Harness for `../scripts/render-return.js`'s `renderItems()` — TODAY's engine, unmodified.
//
// This suite must go green against the current, unchanged script BEFORE any new render code is
// written (Phase 3 lifts `renderItems()`'s per-item logic — cap, tokens, rationale — into a shared
// helper the table and labeled-lists renderers also call). It is the regression net that proves the
// refactor didn't change observable behavior; it never touches `render-return.js` itself.
//
// The `renderTable`/`renderLabeledLists` block below covers Phase 3's addition: a table or
// labeled-lists section that also declares `items` gets the same per-item shaping `renderItems()`
// always did, keyed by `items.key` for a table (a table row cannot carry continuation lines the way a
// bullet can, so its adornments print as a detail list below the table instead).
//
// Run with: node --test payload/domains/development/dev-tests/
'use strict';

const assert = require('node:assert/strict');
const { test, before, after } = require('node:test');

const { loadRenderReturnInternals, cleanupRenderReturnInternals } = require('./support.js');

let renderItems;
let renderTable;
let renderLabeledLists;

before(() => {
  ({ renderItems, renderTable, renderLabeledLists } = loadRenderReturnInternals());
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

// `renderTable()` — a section without `items` stays exactly as before; one with `items.key` gets a
// keyed detail block underneath the table (and its footer, when there is one).

test('a table with no `items` declared renders exactly as before — no detail block', () => {
  const section = {
    field: 'rows',
    title: '### Completeness',
    render: 'table',
    columns: [{ key: 'metric', label: 'Metric' }, { key: 'value', label: 'Value' }],
  };
  const out = renderTable(section, { rows: [{ metric: 'Tasks', value: '3/3' }] });
  assert.equal(out, '| Metric | Value |\n|---|---|\n| Tasks | 3/3 |');
});

test('a table with `items.key` emits one keyed detail-block bullet per row, below the table', () => {
  const section = {
    field: 'gaps',
    title: '### Decision Gaps',
    render: 'table',
    columns: [{ key: 'record', label: 'Record' }, { key: 'task', label: 'Task' }],
    items: { key: 'record', text: 'summary', tokens: [{ name: 'anchor', field: 'anchor' }] },
  };
  const out = renderTable(section, {
    gaps: [{ record: 'contracts/foo', task: '3.2', summary: 'a gap', anchor: '.matecito-ai/edr/foo.md' }],
  });
  assert.equal(
    out,
    '| Record | Task |\n|---|---|\n| contracts/foo | 3.2 |\n\n- contracts/foo — a gap\n  · anchor: .matecito-ai/edr/foo.md'
  );
});

test('a table declaring `items` with no `items.key` fails the render — the detail block has nothing to key by', () => {
  const section = {
    field: 'gaps',
    title: '### Decision Gaps',
    render: 'table',
    columns: [{ key: 'record', label: 'Record' }],
    items: { text: 'summary' },
  };
  assert.throws(
    () => renderTable(section, { gaps: [{ record: 'contracts/foo', summary: 'a gap' }] }),
    /no `items\.key`/
  );
});

test('a table row missing its own `items.key` column value fails the render', () => {
  const section = {
    field: 'gaps',
    title: '### Decision Gaps',
    render: 'table',
    columns: [{ key: 'task', label: 'Task' }],
    items: { key: 'record', text: 'summary' },
  };
  assert.throws(
    () => renderTable(section, { gaps: [{ task: '3.2', summary: 'a gap' }] }),
    /is missing `record`/
  );
});

// `renderLabeledLists()` — a list without `items` stays exactly as before; one with `items` shapes
// every entry across every list the same way `renderItems()` does.

test('a labeled-lists section with no `items` declared renders exactly as before — bare bullets', () => {
  const section = { title: '### Issues Found', render: 'labeled-lists', lists: [{ field: 'critical', label: 'Critical' }] };
  const out = renderLabeledLists(section, { critical: ['a bare string entry'] });
  assert.equal(out, '**Critical**:\n- a bare string entry');
});

test('a labeled-lists section with `items` shapes every entry, across every list', () => {
  const section = {
    title: '### Issues Found',
    render: 'labeled-lists',
    lists: [{ field: 'critical', label: 'Critical' }, { field: 'warning', label: 'Warning' }],
    items: { tokens: [{ name: 'anchor', field: 'anchor' }], rationale: 'rationale' },
  };
  const out = renderLabeledLists(section, {
    critical: [{ text: 'bad thing', anchor: 'foo.md', rationale: 'why it matters' }],
    warning: [{ text: 'lesser thing', anchor: 'bar.md', rationale: 'why it matters less' }],
  });
  assert.equal(
    out,
    '**Critical**:\n- bad thing\n  · anchor: foo.md\n  · rationale: why it matters\n' +
      '**Warning**:\n- lesser thing\n  · anchor: bar.md\n  · rationale: why it matters less'
  );
});
