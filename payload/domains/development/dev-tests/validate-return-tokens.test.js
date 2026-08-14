// Regression coverage for `../scripts/validate-return.js`'s `checkItems` — specifically the
// free-form-token patch (decision-gap-capture-rework): a token declared WITHOUT `values` (e.g. the
// mailbox item's `record: <domain>/<slug>`) is free-form, and a token that DOES declare `values`
// keeps its exact prior behavior, byte for byte.
//
// Run with: node --test payload/domains/development/dev-tests/
'use strict';

const assert = require('node:assert/strict');
const { test, before, after } = require('node:test');

const {
  loadValidateReturnInternals, cleanupValidateReturnInternals,
  loadRenderReturnInternals, cleanupRenderReturnInternals,
} = require('./support.js');

let checkItems;
let renderItems;

before(() => {
  ({ checkItems } = loadValidateReturnInternals());
  ({ renderItems } = loadRenderReturnInternals());
});

after(() => {
  cleanupValidateReturnInternals();
  cleanupRenderReturnInternals();
});

// `checkItems(section, seenTitle, body)` walks the `- ` items of one section's body against
// `section.items.tokens`. Fixtures below mirror the shape the real `.yaml` contracts declare.

test('a token declared without `values` passes with any present, non-null value (free-form)', () => {
  const section = { items: { tokens: [{ name: 'record', field: 'record' }] } };
  const body = ['- item one', '  · record: contracts/some-slug'].join('\n');
  assert.deepEqual(checkItems(section, '### New Decisions', body), []);
});

test('a token declared without `values` still requires presence — TOKEN-MISSING, not silently optional', () => {
  const section = { items: { tokens: [{ name: 'record', field: 'record' }] } };
  const body = ['- item one', '  · rationale: because'].join('\n');
  const problems = checkItems(section, '### New Decisions', body);
  assert.equal(problems.length, 1);
  assert.equal(problems[0].code, 'TOKEN-MISSING');
});

test('a token that DOES declare `values` still rejects an illegal value — TOKEN-ILLEGAL, unchanged', () => {
  const section = {
    items: {
      tokens: [{ name: 'blocking-test', field: 'blocking_test', values: ['none', 'infra'], passing: ['none'] }],
    },
  };
  const body = ['- item one', '  · blocking-test: bogus'].join('\n');
  const problems = checkItems(section, '### New Decisions', body);
  assert.equal(problems.length, 1);
  assert.equal(problems[0].code, 'TOKEN-ILLEGAL');
});

test('a token that DOES declare `values` still rejects a legal-but-non-passing value — TOKEN-WRONG-MAILBOX, unchanged', () => {
  const section = {
    items: {
      tokens: [{ name: 'mandate', field: 'mandate', values: ['covered', 'forced', 'chosen'], passing: ['chosen'] }],
    },
  };
  const body = ['- item one', '  · mandate: covered'].join('\n');
  const problems = checkItems(section, '### Unmandated Forks', body);
  assert.equal(problems.length, 1);
  assert.equal(problems[0].code, 'TOKEN-WRONG-MAILBOX');
});

test('a mix of a free-form token and a values-bearing token on the same item — both checked independently', () => {
  const section = {
    items: {
      tokens: [
        { name: 'blocking-test', field: 'blocking_test', values: ['none', 'infra'], passing: ['none'] },
        { name: 'record', field: 'record' },
      ],
    },
  };
  const passing = ['- item one', '  · blocking-test: none', '  · record: contracts/some-slug'].join('\n');
  assert.deepEqual(checkItems(section, '### New Decisions', passing), []);

  const failing = ['- item one', '  · blocking-test: none'].join('\n');
  const problems = checkItems(section, '### New Decisions', failing);
  assert.equal(problems.length, 1);
  assert.equal(problems[0].code, 'TOKEN-MISSING');
});

// gate-coverage-gaps (task 3.3): `checkItems()` is render-agnostic (it gates on `s.items`, never on
// `s.render`), so a table's rendered detail block — the `| ... |` rows it cannot parse as items,
// followed by the `- {key} — {text}` bullets `render-return.js` now emits underneath it — walks
// through exactly the same loop as a plain items-list. A row missing its declared `anchor` token
// still raises `TOKEN-MISSING`, unchanged.
test('a table\'s detail block missing its declared `anchor` token raises TOKEN-MISSING, same as any items list', () => {
  const section = { items: { tokens: [{ name: 'anchor', field: 'anchor' }] } };
  const body = [
    '| Record | Task |',
    '|---|---|',
    '| contracts/foo | 3.2 |',
    '',
    '- contracts/foo — a gap',
  ].join('\n');
  const problems = checkItems(section, '### Decision Gaps', body);
  assert.equal(problems.length, 1);
  assert.equal(problems[0].code, 'TOKEN-MISSING');
});

// `items.fields` — the compound-item grammar (contract-shape-gate, Phase 1): `checkItems()` learns the
// same `· field: {a} — {b} — {c}` continuation line `render-return.js` emits, splitting on the first
// N-1 occurrences of the declared separator so a description containing the separator survives intact.

test('a malformed `· field:` line missing a part raises FIELD-MALFORMED, naming the item and the part', () => {
  const section = {
    items: { fields: { key: 'field', parts: ['name', 'type', 'description'], separator: ' — ' } },
  };
  const body = ['- a contract needs shaping', '  · field: id —  — the primary key'].join('\n');
  const problems = checkItems(section, '### Contract Shapes Proposed', body);
  assert.equal(problems.length, 1);
  assert.equal(problems[0].code, 'FIELD-MALFORMED');
});

test('an item declaring fields but carrying no `· field:` line raises FIELDS-MISSING', () => {
  const section = {
    items: { fields: { key: 'field', parts: ['name', 'type', 'description'], separator: ' — ' } },
  };
  const body = ['- a contract needs shaping', '  · rationale: because it matters'].join('\n');
  const problems = checkItems(section, '### Contract Shapes Proposed', body);
  assert.equal(problems.length, 1);
  assert.equal(problems[0].code, 'FIELDS-MISSING');
});

// The round-trip: the exact grammars must agree with each other, not merely by reading both sides —
// a compound item rendered by `render-return.js` MUST read back clean through `validate-return.js`'s
// `checkItems()`, with no `SECTION-UNPARSEABLE` false positive and no field line skipped in silence.
test('a rendered compound item reads back clean through checkItems — no SECTION-UNPARSEABLE, no skipped field', () => {
  const section = {
    field: 'contract_proposals',
    title: '### Contract Shapes Proposed',
    items: {
      tokens: [{ name: 'anchor', field: 'anchor' }],
      rationale: 'rationale',
      fields: { key: 'field', parts: ['name', 'type', 'description'], separator: ' — ', field_max: 160 },
    },
  };
  const rendered = renderItems(section, {
    contract_proposals: [{
      text: 'a contract needs shaping',
      anchor: 'scripts/render-return.js:98',
      field: [
        { name: 'id', type: 'uuid', description: 'the primary key' },
        { name: 'email', type: 'string', description: 'the user email — includes an em dash — mid-sentence' },
      ],
      rationale: 'a wrong guess propagates to the DB schema',
    }],
  });
  const problems = checkItems(section, '### Contract Shapes Proposed', rendered);
  assert.deepEqual(problems, []);
});
