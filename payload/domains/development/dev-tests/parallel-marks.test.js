// Regression coverage for `../scripts/validate-parallel-marks.js` — the closed list of three
// malformed shapes it reports, and the two shapes it must NOT flag (unmarked tasks, a group of one).
//
// Run with: node --test payload/domains/development/dev-tests/
'use strict';

const assert = require('node:assert/strict');
const { test } = require('node:test');

const { validate } = require('../scripts/validate-parallel-marks.js');

test('a clean artifact with no marks produces no findings', () => {
  const text = [
    '- [ ] 1.1 Do a thing',
    '      criteria: it happened',
    '- [ ] 1.2 Do another thing',
    '      criteria: it also happened',
  ].join('\n');
  assert.deepEqual(validate(text), []);
});

test('a well-formed group of several tasks produces no findings', () => {
  const text = [
    '- [ ] 1.1 Alpha',
    '      criteria: alpha done',
    '      · parallel-group: a',
    '- [ ] 1.2 Beta',
    '      criteria: beta done',
    '      · parallel-group: a',
  ].join('\n');
  assert.deepEqual(validate(text), []);
});

test('a group of exactly one member is well-formed, not malformed', () => {
  const text = [
    '- [ ] 1.3 Solo',
    '      criteria: solo done',
    '      · parallel-group: b',
  ].join('\n');
  assert.deepEqual(validate(text), []);
});

test('an empty-id sub-line is reported and names the task', () => {
  const text = [
    '- [ ] 1.2 Beta',
    '      criteria: beta done',
    '      · parallel-group:',
  ].join('\n');
  const problems = validate(text);
  assert.equal(problems.length, 1);
  assert.equal(problems[0].code, 'MARK-EMPTY-ID');
  assert.equal(problems[0].task, '1.2');
});

test('a whitespace-only id is treated the same as an empty one', () => {
  const text = [
    '- [ ] 1.2 Beta',
    '      criteria: beta done',
    '      · parallel-group:    ',
  ].join('\n');
  const problems = validate(text);
  assert.equal(problems.length, 1);
  assert.equal(problems[0].code, 'MARK-EMPTY-ID');
});

test('two marks on one task is reported once, and does not also flag an empty id inside the set', () => {
  const text = [
    '- [ ] 1.4 Gamma',
    '      criteria: gamma done',
    '      · parallel-group: a',
    '      · parallel-group:',
  ].join('\n');
  const problems = validate(text);
  assert.equal(problems.length, 1);
  assert.equal(problems[0].code, 'MARK-DUPLICATE');
  assert.equal(problems[0].task, '1.4');
});

test('a mark on the checkbox line itself is reported, never read as a valid mark', () => {
  const text = [
    '- [ ] 1.5 Delta · parallel-group: a',
    '      criteria: delta done',
  ].join('\n');
  const problems = validate(text);
  assert.equal(problems.length, 1);
  assert.equal(problems[0].code, 'MARK-ON-CHECKBOX-LINE');
  assert.equal(problems[0].task, '1.5');
});

test('the check never inspects whether same-id tasks are genuinely independent', () => {
  // Two tasks sharing an id that (per their own descriptions) would collide on the same file — the
  // check has no opinion: independence is the producer's call, re-deriving it here is forbidden.
  const text = [
    '- [ ] 1.1 Rewrite shared.txt line 1',
    '      criteria: shared.txt updated',
    '      · parallel-group: a',
    '- [ ] 1.2 Rewrite shared.txt line 1 too',
    '      criteria: shared.txt updated',
    '      · parallel-group: a',
  ].join('\n');
  assert.deepEqual(validate(text), []);
});

test('findings across independent tasks are all reported, and a clean task is never implicated', () => {
  const text = [
    '- [ ] 1.1 Alpha',
    '      criteria: alpha done',
    '      · parallel-group: a',
    '- [ ] 1.2 Beta',
    '      criteria: beta done',
    '      · parallel-group:',
    '- [ ] 1.4 Gamma',
    '      criteria: gamma done',
    '      · parallel-group: a',
  ].join('\n');
  const problems = validate(text);
  assert.equal(problems.length, 1);
  assert.equal(problems[0].task, '1.2');
});

test('a heading between phases resets sub-line collection for the next task', () => {
  const text = [
    '## Phase 1: Foundation',
    '',
    '- [ ] 1.1 Alpha',
    '      criteria: alpha done',
    '      · parallel-group: a',
    '',
    '## Phase 2: Core',
    '',
    '- [ ] 2.1 Beta',
    '      criteria: beta done',
  ].join('\n');
  assert.deepEqual(validate(text), []);
});
