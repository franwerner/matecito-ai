// Regression coverage for `../scripts/validate-return.js`'s `checkItems` — specifically the
// free-form-token patch (decision-gap-capture-rework): a token declared WITHOUT `values` (e.g. the
// mailbox item's `record: <domain>/<slug>`) is free-form, and a token that DOES declare `values`
// keeps its exact prior behavior, byte for byte.
//
// Run with: node --test payload/domains/development/dev-tests/
'use strict';

const assert = require('node:assert/strict');
const { test, before, after } = require('node:test');

const { loadValidateReturnInternals, cleanupValidateReturnInternals } = require('./support.js');

let checkItems;

before(() => {
  ({ checkItems } = loadValidateReturnInternals());
});

after(() => {
  cleanupValidateReturnInternals();
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
