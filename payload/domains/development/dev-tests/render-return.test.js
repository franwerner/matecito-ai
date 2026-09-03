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
let schema;

before(() => {
  ({ renderItems, renderTable, renderLabeledLists, schema } = loadRenderReturnInternals());
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

// `items.fields` — the compound-item shape (contract-shape-gate, Phase 1): a repeated typed field,
// rendered as one `· field: {name} — {type} — {description}` continuation line per entry, after the
// tokens and before the rationale line. The field COUNT is never capped; only each entry's
// `description` part carries its own `field_max`.

test('a compound item renders its field lines after tokens and before rationale, in order', () => {
  const section = {
    field: 'contract_proposals',
    title: '### Contract Shapes Proposed',
    items: {
      tokens: [{ name: 'anchor', field: 'anchor' }],
      rationale: 'rationale',
      fields: { key: 'field', parts: ['name', 'type', 'description'], separator: ' — ', field_max: 160 },
    },
  };
  const out = renderItems(section, {
    contract_proposals: [{
      text: 'a contract needs shaping',
      anchor: 'scripts/render-return.js:98',
      field: [
        { name: 'id', type: 'uuid', description: 'the primary key' },
        { name: 'email', type: 'string', description: 'the user email' },
      ],
      rationale: 'a wrong guess propagates to the DB schema',
    }],
  });
  assert.equal(
    out,
    '- a contract needs shaping\n' +
      '  · anchor: scripts/render-return.js:98\n' +
      '  · field: id — uuid — the primary key\n' +
      '  · field: email — string — the user email\n' +
      '  · rationale: a wrong guess propagates to the DB schema'
  );
});

test('a field entry missing a part fails the render, naming the item, the field key and the part', () => {
  const section = {
    field: 'contract_proposals',
    title: '### Contract Shapes Proposed',
    items: { fields: { key: 'field', parts: ['name', 'type', 'description'], separator: ' — ' } },
  };
  assert.throws(
    () => renderItems(section, {
      contract_proposals: [{ text: 'a contract', field: [{ name: 'id', description: 'the key' }] }],
    }),
    /`contract_proposals\[0\]\.field\[0\]\.type` is empty/
  );
});

test('a field description over its field_max cap fails the render, naming the item and the cap', () => {
  const section = {
    field: 'contract_proposals',
    title: '### Contract Shapes Proposed',
    items: { fields: { key: 'field', parts: ['name', 'type', 'description'], separator: ' — ', field_max: 10 } },
  };
  assert.throws(
    () => renderItems(section, {
      contract_proposals: [{
        text: 'a contract',
        field: [{ name: 'id', type: 'uuid', description: 'this description is far too long for the cap' }],
      }],
    }),
    /over the 10-char cap/
  );
});

test('an item declaring fields but listing none fails the render, naming the item', () => {
  const section = {
    field: 'contract_proposals',
    title: '### Contract Shapes Proposed',
    items: { fields: { key: 'field', parts: ['name', 'type', 'description'], separator: ' — ' } },
  };
  assert.throws(
    () => renderItems(section, { contract_proposals: [{ text: 'a contract', field: [] }] }),
    /`contract_proposals\[0\]\.field` declares fields but carries none/
  );
});

test('a section declaring no `items.fields` renders byte-identical to before this change', () => {
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

// `schema()` — the shape-announcement path. Before this change the two constraint notes (single-line
// non-empty parts, the summary_max cap) were emitted only inside the `items` render branch; these
// cases prove every render form that can declare `items` (`table`, `labeled-lists`, `items` itself)
// announces the same two notes, and that a non-declaring section stays byte-identical.

test('a table section declaring items.rationale and summary_max announces both constraint notes', () => {
  const contract = {
    phase: 'test-phase',
    block: '## Test',
    statuses: ['done'],
    sections: [{
      title: '### T',
      emitted: 'always',
      render: 'table',
      field: 'rows',
      sentinel: true,
      columns: [{ key: 'id', label: 'Id' }],
      items: { key: 'id', text: 'summary', rationale: 'rationale', summary_max: 250 },
    }],
  };
  const out = schema(contract);
  assert.match(out, /summary and rationale are both required, single-line, non-empty/);
  assert.match(out, /summary is capped at 250 characters — an over-cap value fails the render, exit 1, no stdout/);
});

test('a labeled-lists section declaring items announces both notes once per section, not once per list', () => {
  const contract = {
    phase: 'test-phase',
    block: '## Test',
    statuses: ['done'],
    sections: [{
      title: '### Issues',
      emitted: 'always',
      render: 'labeled-lists',
      lists: [{ field: 'critical', label: 'Critical' }, { field: 'warning', label: 'Warning' }],
      items: { rationale: 'rationale', summary_max: 250 },
    }],
  };
  const out = schema(contract);
  const rationaleMatches = out.match(/are both required, single-line, non-empty/g) || [];
  const capMatches = out.match(/is capped at 250 characters/g) || [];
  assert.equal(rationaleMatches.length, 1);
  assert.equal(capMatches.length, 1);
});

test('an items-rendered declaring section announces exactly as before this change', () => {
  const contract = {
    phase: 'test-phase',
    block: '## Test',
    statuses: ['done'],
    sections: [{
      title: '### Decision Gaps',
      emitted: 'always',
      render: 'items',
      field: 'items',
      items: { rationale: 'rationale', summary_max: 250, tokens: [{ name: 'anchor', field: 'anchor' }] },
    }],
  };
  const out = schema(contract);
  assert.match(out, /items: \[\{ text: string, anchor: undefined, rationale: string \}\]/);
  assert.match(out, /\(empty list renders the "None\." sentinel — never omit the field\)/);
  assert.match(
    out,
    /\(text and rationale are both required, single-line, non-empty — text prints at the gate, rationale never does by default\)/
  );
  assert.match(
    out,
    /\(text is capped at 250 characters — an over-cap value fails the render, exit 1, no stdout\)/
  );
});

test('a non-declaring section states nothing about the split, unchanged by this change', () => {
  const contract = {
    phase: 'test-phase',
    block: '## Test',
    statuses: ['done'],
    sections: [{
      title: '### Completeness',
      emitted: 'always',
      render: 'table',
      field: 'rows',
      columns: [{ key: 'metric', label: 'Metric' }, { key: 'value', label: 'Value' }],
    }],
  };
  const out = schema(contract);
  assert.doesNotMatch(out, /required, single-line, non-empty/);
  assert.doesNotMatch(out, /is capped at/);
});

test('any phase whose contract has the same table+items combination inherits both notes, with no per-phase code', () => {
  const phaseAContract = {
    phase: 'phase-a',
    block: '## A',
    statuses: ['done'],
    sections: [{
      title: '### X',
      emitted: 'always',
      render: 'table',
      field: 'rows',
      sentinel: true,
      columns: [{ key: 'id', label: 'Id' }],
      items: { key: 'id', text: 'summary', rationale: 'rationale', summary_max: 250 },
    }],
  };
  const phaseBContract = {
    phase: 'phase-b',
    block: '## B',
    statuses: ['done'],
    sections: [{
      title: '### Y',
      emitted: 'always',
      render: 'table',
      field: 'rows',
      sentinel: true,
      columns: [{ key: 'id', label: 'Id' }],
      items: { key: 'id', text: 'summary', rationale: 'rationale', summary_max: 250 },
    }],
  };
  for (const out of [schema(phaseAContract), schema(phaseBContract)]) {
    assert.match(out, /are both required, single-line, non-empty/);
    assert.match(out, /is capped at 250 characters/);
  }
});

test('two declaring sections with different summary_max caps each announce only their own', () => {
  const contract = {
    phase: 'test-phase',
    block: '## Test',
    statuses: ['done'],
    sections: [
      { title: '### A', emitted: 'always', render: 'items', field: 'a_items', items: { rationale: 'rationale', summary_max: 250 } },
      { title: '### B', emitted: 'always', render: 'items', field: 'b_items', items: { rationale: 'rationale', summary_max: 500 } },
    ],
  };
  const out = schema(contract);
  const [sectionA, sectionB] = out.split('### B');
  assert.match(sectionA, /is capped at 250 characters/);
  assert.doesNotMatch(sectionA, /is capped at 500 characters/);
  assert.match(sectionB, /is capped at 500 characters/);
  assert.doesNotMatch(sectionB, /is capped at 250 characters/);
});

// `renderTable()`'s `footer.on_sentinel` opt-in — the store-wide line has to survive the empty-table
// `None.` return, or "clean" and "never ran" become indistinguishable exactly when the table is empty.

test('an empty table with footer.on_sentinel emits "None." plus the footer line', () => {
  const section = {
    field: 'rows',
    title: '### Coherence (Capability-Specs)',
    render: 'table',
    sentinel: true,
    columns: [{ key: 'spec', label: 'Capability-spec' }],
    footer: { field: 'spec_store_structure', label: 'Store structure (pre-existing)', on_sentinel: true },
  };
  const out = renderTable(section, { rows: [], spec_store_structure: '0 pre-existing findings' });
  assert.equal(out, 'None.\n\n**Store structure (pre-existing)**: 0 pre-existing findings');
});

test('an empty table without footer.on_sentinel stays the bare "None." sentinel', () => {
  const section = {
    field: 'rows',
    title: '### Coherence (Capability-Specs)',
    render: 'table',
    sentinel: true,
    columns: [{ key: 'spec', label: 'Capability-spec' }],
    footer: { field: 'spec_store_structure', label: 'Store structure (pre-existing)' },
  };
  const out = renderTable(section, { rows: [], spec_store_structure: '0 pre-existing findings' });
  assert.equal(out, 'None.');
});

test('an empty table with footer.on_sentinel and a missing footer field fails naming the field', () => {
  const section = {
    field: 'rows',
    title: '### Coherence (Capability-Specs)',
    render: 'table',
    sentinel: true,
    columns: [{ key: 'spec', label: 'Capability-spec' }],
    footer: { field: 'spec_store_structure', label: 'Store structure (pre-existing)', on_sentinel: true },
  };
  assert.throws(
    () => renderTable(section, { rows: [] }),
    /missing required field `spec_store_structure`/
  );
});

// Regression: the three pre-existing footers (`compliance_summary`, `breakdown.totals`, `error_gate`)
// render exactly as before `footerLine()` was extracted out of the non-empty-rows path.

test('the compliance_summary footer renders exactly as before this change', () => {
  const section = {
    field: 'spec_compliance',
    title: '### Spec Compliance Matrix',
    render: 'table',
    sentinel: true,
    columns: [{ key: 'requirement', label: 'Requirement' }],
    footer: { label: 'Compliance summary', field: 'compliance_summary', format: '{compliant}/{total} scenarios compliant' },
  };
  const out = renderTable(section, {
    spec_compliance: [{ requirement: 'REQ-01' }],
    compliance_summary: { compliant: 3, total: 4 },
  });
  assert.equal(out, '| Requirement |\n|---|\n| REQ-01 |\n\n**Compliance summary**: 3/4 scenarios compliant');
});

test('the breakdown.totals footer renders exactly as before this change', () => {
  const section = {
    field: 'breakdown.phases',
    title: '### Breakdown',
    render: 'table',
    sentinel: true,
    columns: [{ key: 'phase', label: 'Phase' }],
    footer: { label: 'Totals', field: 'breakdown.totals', format: '{tasks} tasks · {work_units} work units' },
  };
  const out = renderTable(section, {
    breakdown: { phases: [{ phase: 'Phase 1' }], totals: { tasks: 5, work_units: 2 } },
  });
  assert.equal(out, '| Phase |\n|---|\n| Phase 1 |\n\n**Totals**: 5 tasks · 2 work units');
});

test('the error_gate footer renders exactly as before this change', () => {
  const section = {
    field: 'ui_verdict',
    title: '## UI Verdict',
    render: 'table',
    columns: [{ key: 'scenario', label: 'Scenario' }],
    footer: { label: 'Error gate', field: 'error_gate', format: 'consoleErrorCount {console} / serverErrorCount {server} → {verdict}' },
  };
  const out = renderTable(section, {
    ui_verdict: [{ scenario: 'login' }],
    error_gate: { console: 0, server: 0, verdict: 'PASS' },
  });
  assert.equal(out, '| Scenario |\n|---|\n| login |\n\n**Error gate**: consoleErrorCount 0 / serverErrorCount 0 → PASS');
});
