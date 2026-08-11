#!/usr/bin/env node
'use strict';

// Validates the `· parallel-group:` marks of a tasks artifact, so the orchestrator's Eligibility
// check (`~/.claude/references/phase-returns/sdd-apply/parallel-batch.md` → `## Eligibility`) stops
// trusting the marks by eye before forming a parallel batch — same point, same manner as
// `validate-return.js`'s Return Contract Check, one level down (the artifact, not the return).
//
// Usage:
//   validate-parallel-marks.js --file tasks.md
//   validate-parallel-marks.js < tasks.md
//
// Exit 0 and silent stdout when every mark is well-formed — a clean artifact reports nothing (see
// parallel-batch.md's "a clean artifact says nothing" scenario). Exit 1 with one line per finding
// when a mark is malformed; the caller degrades only the named task to serial, it never halts the
// phase. Exit 2 when the check could not run at all (e.g. the file is unreadable) — never read as
// "clean", and never a licence to fall back on reading the marks by eye.
//
// Closed list of what this reports — nothing else, and in particular never whether two tasks
// sharing an id are genuinely independent (deriving that from files or descriptions is exactly the
// inference the eligibility rule forbids):
//   - a `· parallel-group:` sub-line with no id (empty or whitespace only)
//   - a task carrying more than one such sub-line
//   - a mark written on the task's `- [ ]` line instead of its own indented sub-line

const fs = require('fs');

const TASK_LINE = /^-\s*\[[ xX]\]\s*(\S+)/;
const MARK_LINE = /^\s*[·*-]?\s*parallel-group\s*:\s*(.*)$/i;
const INLINE_MARK = /parallel-group\s*:/i;
const HEADING = /^#{1,6}\s+/;

// One entry per task, in document order. `marks` collects every `· parallel-group:` sub-line found
// under it; `inlineMark` is true when the label also appears on the task's own `- [ ]` line.
function parseTasks(text) {
  const lines = text.split('\n');
  const tasks = [];
  let current = null;

  for (let i = 0; i < lines.length; i++) {
    const raw = lines[i];
    const taskMatch = raw.match(TASK_LINE);
    if (taskMatch) {
      current = { id: taskMatch[1], line: i + 1, marks: [], inlineMark: INLINE_MARK.test(raw) };
      tasks.push(current);
      continue;
    }
    if (!current) continue;
    // A heading or an unindented line ends the run of sub-lines belonging to the current task.
    if (HEADING.test(raw.trim()) || (raw.trim() !== '' && !/^\s/.test(raw))) {
      current = null;
      continue;
    }
    const markMatch = raw.match(MARK_LINE);
    if (markMatch) current.marks.push({ value: markMatch[1].trim(), line: i + 1 });
  }
  return tasks;
}

function validate(text) {
  const problems = [];
  const add = (code, taskId, line, message) => problems.push({ code, task: taskId, line, message });

  for (const task of parseTasks(text)) {
    if (task.inlineMark) {
      add('MARK-ON-CHECKBOX-LINE', task.id, task.line,
        `task ${task.id} — \`parallel-group\` appears on the \`- [ ]\` line, not on its own indented sub-line; not read as a mark, task runs serial`);
    }
    if (task.marks.length > 1) {
      // Multiple sub-lines make the task's intended group unrecoverable — report once, don't also
      // flag an empty id inside the same set as a second, independent problem.
      add('MARK-DUPLICATE', task.id, task.marks[task.marks.length - 1].line,
        `task ${task.id} — carries ${task.marks.length} \`· parallel-group:\` sub-lines; neither id is picked, task runs serial`);
      continue;
    }
    for (const mark of task.marks) {
      if (mark.value === '') {
        add('MARK-EMPTY-ID', task.id, mark.line,
          `task ${task.id} — \`· parallel-group:\` carries no id, task runs serial`);
      }
    }
  }

  return problems;
}

function arg(name) {
  const i = process.argv.indexOf(`--${name}`);
  return i === -1 ? undefined : process.argv[i + 1];
}

function main() {
  const file = arg('file');
  let text;
  try {
    text = file ? fs.readFileSync(file, 'utf8') : fs.readFileSync(0, 'utf8');
  } catch (e) {
    process.stderr.write(`cannot read the tasks artifact: ${e.message}\n`);
    process.exit(2);
  }

  const problems = validate(text);
  for (const p of problems) process.stdout.write(`ERROR ${p.code}: ${p.message}\n`);
  process.exit(problems.length === 0 ? 0 : 1);
}

if (require.main === module) main();

module.exports = { parseTasks, validate };
