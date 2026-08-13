#!/usr/bin/env node
'use strict';

// Builds a phase's return block from data, against the contract in `<phase>.yaml`.
//
// The agent never types a heading, so a section cannot be dropped, re-levelled or left without its
// sentinel. Derived values are computed here, which is what makes a summary that contradicts its own
// body impossible to express rather than merely detectable afterwards.
//
// Usage:
//   render-return.js --phase sdd-design --schema
//   render-return.js --phase sdd-design --data data.json [--contracts DIR]
//
// Exit 0 renders to stdout, 1 on invalid data (with the offending field), 2 when it cannot run.

const fs = require('fs');
const path = require('path');
const yaml = require('./lib/yaml');

function get(data, dotted) {
  return dotted.split('.').reduce((o, k) => (o == null ? undefined : o[k]), data);
}

function fail(msg) {
  process.stderr.write(`${msg}\n`);
  process.exit(1);
}

function fmt(template, value, field) {
  if (value && typeof value === 'object') {
    // Report every missing key at once: one round-trip per key is how a two-object contract turns
    // into four failed runs.
    const missing = keysOf(template).filter((k) => value[k] === undefined);
    if (missing.length) fail(`\`${field}\` is missing ${missing.map((k) => `\`${k}\``).join(', ')} — it needs { ${keysOf(template).join(', ')} }`);
    return template.replace(/\{(\w+)\}/g, (_, k) => String(value[k]));
  }
  return template.replace(/\{n\}/g, String(value));
}

function derive(spec, data) {
  const [kind, arg] = spec.split(':');
  if (kind === 'count') return (get(data, arg) || []).length;
  fail(`unknown derived value \`${spec}\``);
}

// A section with a single `token` field desugars to a one-element `tokens` list here, so the render
// and validate loops have exactly one shape to walk regardless of how many tokens an item carries.
// `sdd-design`'s single-token contract goes through this path unchanged — same output either way.
function tokensOf(items) {
  if (items.tokens) return items.tokens;
  if (items.token) return [{ name: items.token, field: items.token_field, values: items.values, passing: items.passing }];
  return [];
}

// --- rendering ------------------------------------------------------------

function renderLabeledBullets(section, data) {
  const out = [];
  for (const b of section.bullets || []) {
    // Same shape as the section-level gate below: required, never defaulted to "off" when omitted.
    if (b.emitted === 'conditional') {
      const gate = data[b.when];
      if (gate === undefined) fail(`missing required field \`${b.when}\` — it decides whether "${b.label}" is emitted, and guessing it is how a bullet goes missing`);
      if (!gate) continue;
    }
    let value;
    if (b.derived) value = derive(b.derived, data);
    else {
      value = get(data, b.field);
      if (value === undefined || value === null || value === '') {
        fail(`missing required field \`${b.field}\` (renders "${b.label}" under ${section.title})`);
      }
    }
    out.push(`- **${b.label}**: ${b.format ? fmt(b.format, value, b.field || b.derived) : value}`);
  }
  return out.join('\n');
}

function renderFields(section, data) {
  const out = [];
  const lead = get(data, section.lead);
  if (lead === undefined || lead === '') fail(`missing required field \`${section.lead}\` (opens ${section.title})`);
  out.push(String(lead), '');
  for (const f of section.fields || []) {
    const v = get(data, f.field);
    if (v === undefined || v === '') fail(`missing required field \`${f.field}\` (renders "${f.label}" under ${section.title})`);
    out.push(`**${f.label}**: ${v}`, '');
  }
  return out.join('\n').trimEnd();
}

// One item's adornment: the text (capped, single-line) plus its `· token:`/`· rationale:` lines —
// everything an item carries AFTER its own first bullet line. Every render form keeps its own first
// line (`- {text}` for a plain items list, `- {key} — {text}` for a table's detail block), because a
// table row cannot carry continuation lines the way a bullet can; this is the part they share.
// `fieldName` lets a caller with no single `section.field` (labeled-lists, one list per entry) name
// the right path in an error instead of blaming `undefined`.
function shapeItem(section, item, i, fieldName) {
  const field = fieldName || section.field;
  const spec = section.items || {};
  const textField = spec.text || 'text';
  const text = typeof item === 'string' ? item : item[textField];
  if (!text) fail(`\`${field}[${i}]\` has no \`${textField}\``);
  if (String(text).includes('\n')) fail(`\`${field}[${i}].${textField}\` spans multiple lines — each part must be a single line`);
  // `items.summary_max` caps `text` in characters, not lines — a terminal width is invisible to a
  // script, a character count is not. Enforced here, while the phase is still alive to shorten it;
  // `validate-return.js` stays unchanged (see the design's "Where the cap runs").
  if (spec.summary_max && String(text).length > spec.summary_max) {
    fail(`\`${field}[${i}].${textField}\` is ${String(text).length} chars, over the ${spec.summary_max}-char cap — shorten it before it reaches the gate`);
  }

  const lines = [];
  // One or more tokens, in declared order — a single `token` desugars to one via `tokensOf`.
  for (const t of tokensOf(spec)) {
    const value = typeof item === 'object' ? item[t.field] : undefined;
    if (value === undefined) fail(`\`${field}[${i}]\` has no \`${t.field}\` — a declared token is required, never inferred`);
    if (t.values && !t.values.includes(value)) {
      fail(`\`${field}[${i}].${t.field}\` is "${value}", not one of ${JSON.stringify(t.values)}`);
    }
    // matecito-ai: el mensaje nombraba la contradicción y no la salida. El ejecutor que declara
    // honestamente un valor no-`passing` tiene a dónde ir —su propia skill se lo dice— pero acá
    // llegaba un error que sólo le decía que estaba mal, en el momento en que más barato es
    // recordarle dónde va el ítem.
    if (t.passing && !t.passing.includes(value)) {
      fail(`\`${field}[${i}].${t.field}\` is "${value}", which contradicts ${section.title} — an item with this token does not belong in that section: return \`blocked\` and file it under \`### Blocker\` instead`);
    }
    lines.push(`  · ${t.name}: ${value}`);
  }

  // `items.rationale` IS the opt-in for the summary/rationale split (see the phase-return `.md`).
  // Both parts are required and single-line — a section with this key never emits a half item.
  if (spec.rationale) {
    const rationale = typeof item === 'object' ? item[spec.rationale] : undefined;
    if (!rationale) fail(`\`${field}[${i}]\` has no \`${spec.rationale}\` — a declaring section requires both parts`);
    if (String(rationale).includes('\n')) fail(`\`${field}[${i}].${spec.rationale}\` spans multiple lines — each part must be a single line`);
    lines.push(`  · rationale: ${rationale}`);
  }

  return { text, lines };
}

function renderItems(section, data) {
  const list = get(data, section.field);
  if (list === undefined) fail(`missing required field \`${section.field}\` (renders ${section.title}; use [] when there is nothing)`);
  if (!Array.isArray(list)) fail(`\`${section.field}\` must be a list`);
  // The sentinel is emitted because the list is empty, never because the agent remembered to.
  if (list.length === 0) return 'None.';

  const out = [];
  for (const [i, item] of list.entries()) {
    const { text, lines } = shapeItem(section, item, i);
    out.push(`- ${text}`, ...lines);
  }
  return out.join('\n');
}

// A table the agent supplies as rows of objects cannot come out misaligned, and a column it forgets is
// named at render time instead of surviving as a ragged pipe nobody notices.
//
// A row that also declares `section.items` gets a keyed detail block below the table (and its
// footer): a table row is one markdown line and cannot carry the continuation lines an adornment
// needs, so the same `shapeItem()` output that a plain items-list would print inline instead prints as
// its own `- {key} — {text}` entry, keyed by the column `items.key` names.
function renderTable(section, data) {
  const rows = get(data, section.field);
  if (rows === undefined) fail(`missing required field \`${section.field}\` (renders ${section.title}; use [] when there is nothing)`);
  if (!Array.isArray(rows)) fail(`\`${section.field}\` must be a list of rows`);
  if (rows.length === 0 && section.sentinel) return 'None.';

  if (section.items && !section.items.key) {
    fail(`\`${section.title}\` declares \`items\` on a table but no \`items.key\` — the detail block needs the column that labels each row`);
  }

  const cols = section.columns || [];
  const out = [
    `| ${cols.map((c) => c.label).join(' | ')} |`,
    `|${cols.map(() => '---').join('|')}|`,
  ];
  const details = [];
  for (const [i, row] of rows.entries()) {
    const missing = cols.filter((c) => row[c.key] === undefined).map((c) => c.key);
    if (missing.length) fail(`\`${section.field}[${i}]\` is missing ${missing.map((k) => `\`${k}\``).join(', ')} — a row needs { ${cols.map((c) => c.key).join(', ')} }`);
    out.push(`| ${cols.map((c) => String(row[c.key]).replace(/\|/g, '\\|')).join(' | ')} |`);

    if (section.items) {
      const keyValue = row[section.items.key];
      if (keyValue === undefined) fail(`\`${section.field}[${i}]\` is missing \`${section.items.key}\` — \`items.key\` names the column that labels its detail block`);
      const { text, lines } = shapeItem(section, row, i);
      details.push(`- ${keyValue} — ${text}`, ...lines);
    }
  }
  if (section.footer) {
    const value = section.footer.derived ? derive(section.footer.derived, data) : get(data, section.footer.field);
    if (value === undefined) fail(`missing required field \`${section.footer.field}\` (renders the footer of ${section.title})`);
    out.push('', `**${section.footer.label}**: ${section.footer.format ? fmt(section.footer.format, value, section.footer.field) : value}`);
  }
  if (details.length) out.push('', ...details);
  return out.join('\n');
}

// Labeled entries whose value may be a command's real output, which has to survive verbatim.
function renderBlocks(section, data) {
  const out = [];
  for (const b of section.blocks || []) {
    const v = get(data, b.field);
    if (v === undefined || v === '') fail(`missing required field \`${b.field}\` (renders "${b.label}" under ${section.title})`);
    if (typeof v === 'object' && v.summary !== undefined) {
      out.push(`**${b.label}**: ${v.summary}`);
      if (v.output) out.push('```text', String(v.output).replace(/```/g, "'''"), '```');
    } else {
      out.push(`**${b.label}**: ${v}`);
    }
    out.push('');
  }
  return out.join('\n').trimEnd();
}

// Several named lists in one section, each with its own sentinel. `section.items` (declared once,
// the same shape `renderItems()` and `renderTable()` read) applies uniformly to every entry of every
// list here — a plain string entry only when the section declares no `items` at all.
function renderLabeledLists(section, data) {
  const out = [];
  for (const l of section.lists || []) {
    const items = get(data, l.field);
    if (items === undefined) fail(`missing required field \`${l.field}\` (renders "${l.label}" under ${section.title}; use [] when there is nothing)`);
    if (!Array.isArray(items)) fail(`\`${l.field}\` must be a list`);
    if (items.length === 0) {
      out.push(`**${l.label}**: None`);
      continue;
    }
    out.push(`**${l.label}**:`);
    for (const [i, it] of items.entries()) {
      if (section.items) {
        const { text, lines } = shapeItem(section, it, i, l.field);
        out.push(`- ${text}`, ...lines);
      } else {
        out.push(`- ${typeof it === 'string' ? it : it.text}`);
      }
    }
  }
  return out.join('\n');
}

const RENDERERS = {
  'labeled-bullets': renderLabeledBullets,
  'labeled-lists': renderLabeledLists,
  fields: renderFields,
  items: renderItems,
  table: renderTable,
  blocks: renderBlocks,
  line: (section, data) => {
    const v = get(data, section.field);
    if (v === undefined || v === '') fail(`missing required field \`${section.field}\` (renders ${section.title})`);
    return String(v);
  },
};

// A phase whose return changes block by status (sdd-intake ships a Discovery Form on Pass 1 and an
// Intake Brief on Pass 2) declares `variants`; everything else uses the top-level block and sections.
function resolveVariant(contract, status) {
  if (!contract.variants) return contract;
  const v = contract.variants.find((x) => (x.statuses || []).includes(status));
  if (!v) fail(`no variant of \`${contract.phase}\` covers status "${status}"`);
  return { ...contract, ...v };
}

function blockTitle(template, data) {
  return template.replace(/\{(\w+)\}/g, (_, k) => {
    if (data[k] === undefined || data[k] === '') fail(`missing required field \`${k}\` (it names the block: "${template}")`);
    return String(data[k]);
  });
}

function render(rawContract, data) {
  const status = data.status;
  if (!status) fail('missing required field `status`');
  if (!(rawContract.statuses || []).includes(status)) {
    fail(`\`status\` is "${status}", not one of ${JSON.stringify(rawContract.statuses)}`);
  }
  const contract = resolveVariant(rawContract, status);

  const out = [blockTitle(contract.block, data), ''];
  for (const h of contract.header || []) {
    const v = get(data, h.field);
    if (v === undefined || v === '') fail(`missing required field \`${h.field}\` (renders "${h.label}" in the header)`);
    out.push(`**${h.label}**: ${v}`);
  }
  if ((contract.header || []).length) out.push('');

  for (const section of contract.sections || []) {
    if (section.emitted === 'on-status' && !(section.statuses || []).includes(status)) continue;
    // A gate the agent evaluated (a store exists, a flag is on) travels as a boolean it supplies.
    // Requiring it explicitly is the point: an omitted gate would silently drop the section, which is
    // the failure these contracts exist to prevent.
    if (section.emitted === 'conditional') {
      const gate = data[section.when];
      if (gate === undefined) fail(`missing required field \`${section.when}\` — it decides whether ${section.title} is emitted, and guessing it is how a section goes missing`);
      if (!gate) continue;
    }

    let title = section.title;
    if (section.variant && data[section.variant.when]) title = section.variant.title;

    const renderer = RENDERERS[section.render];
    if (!renderer) fail(`contract error: section ${section.title} declares no known \`render\``);
    // Errors quote the title actually being rendered, not the base one, or they point at a section
    // the reader will not find in the output.
    out.push(title, renderer({ ...section, title }, data), '');
  }
  return out.join('\n').replace(/\n{3,}/g, '\n\n').trimEnd() + '\n';
}

// --- schema ---------------------------------------------------------------

function keysOf(format) {
  return [...format.matchAll(/\{(\w+)\}/g)].map((m) => m[1]);
}

// The field list an `items`-shaped entry carries, shared by every render form that can declare
// `items` (`items` itself, plus `table` and `labeled-lists` now that they wire into `shapeItem()`).
function itemsFields(spec) {
  const fields = [`${spec.text || 'text'}: string`];
  for (const t of tokensOf(spec)) fields.push(`${t.field}: ${JSON.stringify(t.passing || t.values)}`);
  if (spec.rationale) fields.push(`${spec.rationale}: string`);
  return fields;
}

// The agent asks the tool for the shape instead of memorizing it, so the two cannot drift.
function schema(contract) {
  const out = [`Data shape for \`${contract.phase}\` — write this as JSON and pass it with --data.`, ''];
  out.push(`status: one of ${JSON.stringify(contract.statuses)}`);
  if (contract.variants) {
    out.push('', 'This phase returns a DIFFERENT block per status:');
    for (const v of contract.variants) out.push(`  ${JSON.stringify(v.statuses)} → ${v.block}`);
    for (const v of contract.variants) {
      out.push('', `── ${v.block} ──`);
      out.push(...schema({ ...contract, ...v, variants: null, phase: contract.phase }).split('\n').slice(3));
    }
    return out.join('\n') + '\n';
  }
  for (const k of keysOf(contract.block || '')) out.push(`${k}: string   (names the block: "${contract.block}")`);
  for (const h of contract.header || []) out.push(`${h.field}: string   (header "${h.label}")`);

  for (const s of contract.sections || []) {
    let when = '';
    if (s.emitted === 'on-status') when = ` [only on status ${JSON.stringify(s.statuses)}]`;
    if (s.emitted === 'conditional') when = ` [emitted only when \`${s.when}\` is true]`;
    out.push('', `${s.title}${when}`);
    if (s.emitted === 'conditional') out.push(`  ${s.when}: boolean   (required — decides whether the section exists)`);
    if (s.render === 'table') {
      const cols = (s.columns || []).map((c) => c.key).join(', ');
      out.push(`  ${s.field}: [{ ${cols} }]   (one row per entry)`);
      if (s.sentinel) out.push('  (empty list renders the "None." sentinel — never omit the field)');
      if (s.items) out.push(`  (each row also carries { ${itemsFields(s.items).join(', ')} }, rendered as a detail list below the table keyed by its \`${s.items.key}\` column)`);
      if (s.footer && !s.footer.derived) {
        out.push(`  ${s.footer.field}: ${s.footer.format ? `{ ${keysOf(s.footer.format).join(', ')} }` : 'string'}   ("${s.footer.label}")`);
      } else if (s.footer) {
        out.push(`  (footer "${s.footer.label}" derived from ${s.footer.derived.split(':')[1]} — do NOT supply)`);
      }
    } else if (s.render === 'blocks') {
      for (const b of s.blocks || []) out.push(`  ${b.field}: string | { summary, output }   ("${b.label}"; output is fenced verbatim)`);
    } else if (s.render === 'labeled-lists') {
      for (const l of s.lists || []) {
        if (s.items) out.push(`  ${l.field}: [{ ${itemsFields(s.items).join(', ')} }]   ("${l.label}"; empty list renders "None")`);
        else out.push(`  ${l.field}: [string]   ("${l.label}"; empty list renders "None")`);
      }
    } else if (s.render === 'labeled-bullets') {
      for (const b of s.bullets || []) {
        const bWhen = b.emitted === 'conditional' ? ` [emitted only when \`${b.when}\` is true]` : '';
        if (b.derived) out.push(`  (derived from ${b.derived.split(':')[1]} — do NOT supply "${b.label}")${bWhen}`);
        else if (b.format) out.push(`  ${b.field}: { ${keysOf(b.format).join(', ')} }   (rendered as "${b.format}")${bWhen}`);
        else out.push(`  ${b.field}: string${bWhen}`);
        if (b.emitted === 'conditional') out.push(`  ${b.when}: boolean   (required — decides whether "${b.label}" bullet is emitted)`);
      }
    } else if (s.render === 'fields') {
      out.push(`  ${s.lead}: string`);
      for (const f of s.fields || []) out.push(`  ${f.field}: string   ("${f.label}")`);
    } else if (s.render === 'items') {
      const spec = s.items || {};
      out.push(`  ${s.field}: [{ ${itemsFields(spec).join(', ')} }]`);
      out.push(`  (empty list renders the "None." sentinel — never omit the field)`);
      if (spec.rationale) out.push(`  (${spec.text || 'text'} and ${spec.rationale} are both required, single-line, non-empty — ${spec.text || 'text'} prints at the gate, ${spec.rationale} never does by default)`);
      if (spec.summary_max) out.push(`  (${spec.text || 'text'} is capped at ${spec.summary_max} characters — an over-cap value fails the render, exit 1, no stdout)`);
    } else if (s.render === 'line') {
      out.push(`  ${s.field}: string`);
    }
    if (s.variant) out.push(`  ${s.variant.when}: boolean   (retitles to "${s.variant.title}")`);
  }
  return out.join('\n') + '\n';
}

// --- cli ------------------------------------------------------------------

function arg(name, fallback) {
  const i = process.argv.indexOf(`--${name}`);
  return i === -1 ? fallback : process.argv[i + 1];
}

// Deployed, the contracts sit in the user's ~/.claude; inside the payload repo they sit next to this
// script. The probe is the phase's own file, not the directory: the deployed directory can exist and
// still lack a contract that has not shipped yet.
function contractsDir(phase) {
  const explicit = arg('contracts');
  if (explicit) return explicit;
  const home = path.join(require('os').homedir(), '.claude', 'references', 'phase-returns');
  if (phase && fs.existsSync(path.join(home, phase, `${phase}.yaml`))) return home;
  return path.resolve(__dirname, '..', 'references', 'phase-returns');
}

function main() {
  const phase = arg('phase');
  if (!phase) {
    process.stderr.write(`usage: render-return.js --phase <name> [--schema | --data <file.json>] [--contracts <dir>]\n  --contracts defaults to ${contractsDir(phase)}\n`);
    process.exit(2);
  }
  const dir = contractsDir(phase);

  let contract;
  try {
    contract = yaml.parse(fs.readFileSync(path.join(dir, phase, `${phase}.yaml`), 'utf8'));
  } catch (e) {
    process.stderr.write(`cannot read the contract for \`${phase}\`: ${e.message}\n`);
    process.exit(2);
  }

  if (process.argv.includes('--schema')) {
    process.stdout.write(schema(contract));
    return;
  }

  const file = arg('data');
  let data;
  try {
    data = JSON.parse(file ? fs.readFileSync(file, 'utf8') : fs.readFileSync(0, 'utf8'));
  } catch (e) {
    process.stderr.write(`cannot read the data: ${e.message}\n`);
    process.exit(2);
  }
  process.stdout.write(render(contract, data));
}

main();
