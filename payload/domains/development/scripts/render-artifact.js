#!/usr/bin/env node
'use strict';

// Builds a durable artifact's body (EDR, capability-spec) — and separately its INDEX entries —
// from data, against the contract in `<type>.yaml` beside the artifact's canonical template.
//
// Mirrors render-return.js one level out: the caller never types a heading or a bullet line by
// hand, so a section cannot be dropped, mis-ordered, or filled in when the contract says a status
// forces it empty (the Inferred EDR invariant). The renderer never writes to disk — stdout only.
//
// Usage:
//   render-artifact.js --type edr --schema
//   render-artifact.js --type edr --data data.json                  # stdout: the body
//   render-artifact.js --type edr --data data.json --index-entries  # stdout: JSON entries
//
// Exit 0 renders to stdout, 1 on invalid data (with the offending field), 2 when it cannot run.

const fs = require('fs');
const path = require('path');
const yaml = require('./lib/yaml');

const TYPES = {
  edr: ['edr', 'edr.yaml'],
  'capability-spec': ['spec', 'capability.yaml'],
};

function get(data, dotted) {
  return dotted.split('.').reduce((o, k) => (o == null ? undefined : o[k]), data);
}

function fail(msg) {
  process.stderr.write(`${msg}\n`);
  process.exit(1);
}

function bail(msg) {
  process.stderr.write(`${msg}\n`);
  process.exit(2);
}

function present(v) {
  if (v === undefined || v === null || v === '') return false;
  if (Array.isArray(v)) return v.length > 0;
  return true;
}

// Fills a `{key}` template from a flat object — used for header/section formats (per-item) and for
// the title/path/index-row templates (against the whole data object).
function fmt(template, obj) {
  return template.replace(/\{(\w+)\}/g, (_, k) => {
    if (obj[k] === undefined) fail(`missing \`${k}\` for format "${template}"`);
    return String(obj[k]);
  });
}

// --- rendering per kind -----------------------------------------------------

function renderProse(section, data) {
  const lead = section.lead ? fmt(section.lead, data) : null;
  const body = String(get(data, section.field) || '');
  if (!lead) return body;
  return body ? `${lead}\n\n${body}` : lead;
}

function renderBullets(section, data) {
  const items = get(data, section.field) || [];
  return items.map((it, i) => {
    if (!section.format) return `- ${it}`;
    if (typeof it !== 'object' || it === null) {
      fail(`\`${section.field}[${i}]\` must be an object — ${section.title} renders it as "${section.format}"`);
    }
    return fmt(section.format, it);
  }).join('\n');
}

function renderNumbered(section, data) {
  const items = get(data, section.field) || [];
  return items.map((it, i) => `${i + 1}. ${it}`).join('\n');
}

// A fixed set of labeled lines rendered once — not a list field. Each bullet fails on its own
// missing value unless it declares `emitted: when-present`, which is how `## Evidencia (inferida)`
// keeps `prevalencia` out when the evidence `kind` does not carry one (config, ausencia).
function renderLabeledBullets(section, data) {
  const out = [];
  for (const b of section.bullets || []) {
    const v = get(data, b.field);
    if (!present(v)) {
      if (b.emitted === 'when-present') continue;
      fail(`missing required field \`${b.field}\` (renders "${b.label}" under ${section.title})`);
    }
    out.push(`- **${b.label}:** ${v}`);
  }
  return out.join('\n');
}

// matecito-ai: `verification` is optional and leads the scenario, before GIVEN — its ABSENCE carries
// meaning (`in-scope`), so a scenario without the line is not an unmarked scenario. Unknown fields
// abort instead of being dropped: this renderer writes durable capability-specs, and content
// disappearing from one without an error is the exact failure the token was introduced to close.
const SCENARIO_FIELDS = ['name', 'given', 'when', 'then', 'verification'];

function renderScenarios(section, data) {
  const items = get(data, section.field) || [];
  return items.map((s, i) => {
    for (const k of ['name', 'given', 'when', 'then']) {
      if (!present(s[k])) fail(`\`${section.field}[${i}]\` is missing \`${k}\` — every scenario needs name/given/when/then`);
    }
    for (const k of Object.keys(s)) {
      if (!SCENARIO_FIELDS.includes(k)) {
        fail(`\`${section.field}[${i}]\` carries unknown field \`${k}\` — this renderer emits only ${SCENARIO_FIELDS.join(', ')}, so it would be dropped without a trace`);
      }
    }
    const lines = [];
    if (present(s.verification)) lines.push(`- **verification:** ${s.verification}`);
    lines.push(`- **GIVEN** ${s.given}`, `- **WHEN** ${s.when}`, `- **THEN** ${s.then}`);
    return `### Scenario: ${s.name}\n\n${lines.join('\n')}`;
  }).join('\n\n');
}

const RENDERERS = {
  prose: renderProse,
  bullets: renderBullets,
  numbered: renderNumbered,
  'labeled-bullets': renderLabeledBullets,
  scenarios: renderScenarios,
};

// --- body -------------------------------------------------------------------

function renderHeader(contract, data) {
  const status = data.status;
  const out = [];
  for (const h of contract.header || []) {
    if (h.emitted === 'on-status' && !(h.statuses || []).includes(status)) continue;
    const src = get(data, h.field);
    const v = h.format ? (present(src) ? fmt(h.format, src) : undefined) : src;
    if (!present(v)) {
      if (h.emitted === 'when-present') continue;
      fail(`missing required field \`${h.field}\` (header "${h.label}"${h.emitted === 'on-status' ? `, required on status ${JSON.stringify(h.statuses)}` : ''})`);
    }
    out.push(`- **${h.label}:** ${Array.isArray(v) ? v.join(', ') : v}`);
  }
  return out;
}

// A section's presence is decided here — by the contract and the data — never by the caller. The
// three levers: `always` (must be there), `on-status` (there iff Status matches), `when-present`
// (there iff the field has content). `empty_on` sits across all three: when Status is in that list
// the heading still renders, but the body is forced empty regardless of what the caller supplied —
// that is what makes the Inferred-EDR invariant a property of the renderer, not of good behavior.
function sectionBody(section, data, status) {
  if (section.emitted === 'on-status' && !(section.statuses || []).includes(status)) return null;
  if ((section.empty_on || []).includes(status)) return '';

  const renderer = RENDERERS[section.render];
  if (!renderer) fail(`contract error: section "${section.title}" declares unknown render "${section.render}"`);

  if (section.render === 'labeled-bullets') return renderer(section, data);

  const value = get(data, section.field);
  const hasContent = Array.isArray(value) ? value.length > 0 : present(value);
  if (!hasContent) {
    if (section.emitted === 'when-present') return null;
    fail(`missing required field \`${section.field}\` (renders ${section.title}${section.emitted === 'on-status' ? `; required on status ${JSON.stringify(section.statuses)}` : ''})`);
  }
  return renderer(section, data);
}

function render(contract, data) {
  const status = data.status;
  if (!present(status)) fail('missing required field `status`');
  if (!(contract.statuses || []).includes(status)) {
    fail(`\`status\` is "${status}", not one of ${JSON.stringify(contract.statuses)}`);
  }
  for (const loc of contract.location || []) {
    if (!present(data[loc])) fail(`missing required field \`${loc}\` — part of the artifact's location`);
  }

  const out = [fmt(contract.title, data), '', ...renderHeader(contract, data), ''];
  for (const section of contract.sections || []) {
    const body = sectionBody(section, data, status);
    if (body === null) continue;
    out.push(section.title, body, '');
  }
  return out.join('\n').replace(/\n{3,}/g, '\n\n').trimEnd() + '\n';
}

// --- index entries ------------------------------------------------------------

// `path`, `file` and every `row` are derived from --data, never supplied — the caller cannot drift
// the index from the body because both come out of the same object.
function indexEntries(contract, data) {
  return (contract.index_entries || []).map((e) => {
    const keyValue = get(data, e.key);
    if (!present(keyValue)) fail(`missing \`${e.key}\` — needed to key the "${e.index}" INDEX row`);
    return {
      index: e.index,
      file: fmt(e.file, data),
      key: keyValue,
      section: e.section,
      scaffold: e.scaffold,
      row: fmt(e.row, data),
    };
  });
}

// --- schema ---------------------------------------------------------------

function schema(contract) {
  const out = [`Data shape for \`${contract.artifact}\` — write this as JSON and pass it with --data.`, ''];
  out.push(`status: one of ${JSON.stringify(contract.statuses)}`);
  for (const loc of contract.location || []) out.push(`${loc}: string   (location — part of the derived path)`);
  out.push(`title: string   (fills "${contract.title}")`);
  out.push('', `(derived, do NOT supply: \`path\` — built from ${JSON.stringify(contract.location)} via "${contract.path}")`);

  for (const h of contract.header || []) {
    let when = '';
    if (h.emitted === 'when-present') when = ' [optional — omitted, not emptied, when absent]';
    if (h.emitted === 'on-status') when = ` [only on status ${JSON.stringify(h.statuses)}]`;
    const shape = h.format ? `{ ${[...new Set([...h.format.matchAll(/\{(\w+)\}/g)].map((m) => m[1]))].join(', ')} }` : 'string';
    out.push(`${h.field}: ${shape}${when}   (header "${h.label}")`);
  }

  for (const s of contract.sections || []) {
    let when = '';
    if (s.emitted === 'on-status') when = ` [only on status ${JSON.stringify(s.statuses)}]`;
    if (s.emitted === 'when-present') when = ' [optional — omitted, not emptied, when absent]';
    if (s.empty_on) when += ` (forced empty on status ${JSON.stringify(s.empty_on)}, regardless of any value supplied)`;
    out.push('', `${s.title}${when}`);
    if (s.render === 'labeled-bullets') {
      for (const b of s.bullets || []) out.push(`  ${b.field}: string${b.emitted === 'when-present' ? ' [optional]' : ''}   ("${b.label}")`);
    } else if (s.render === 'scenarios') {
      out.push(`  ${s.field}: [{ name, given, when, then, verification? }]   (verification is optional — its absence means \`in-scope\`; any other field aborts the render)`);
    } else if (s.format) {
      const keys = [...new Set([...s.format.matchAll(/\{(\w+)\}/g)].map((m) => m[1]))];
      out.push(`  ${s.field}: [{ ${keys.join(', ')} }]   (each item rendered as "${s.format}")`);
    } else if (s.render === 'prose') {
      out.push(`  ${s.field}: string`);
    } else {
      out.push(`  ${s.field}: [string]`);
    }
  }

  // matecito-ai: este bloque enumeraba index/key/file/section/scaffold y nunca los placeholders del
  // `row:`. Los campos que sólo el row consume (`blurb`, `trigger`, …) no aparecían en ninguna parte
  // del schema, así que un agente que lo seguía al pie renderizaba bien el artefacto y moría en la
  // segunda llamada — y tres skills instruyen exactamente esa secuencia de dos invocaciones.
  const declared = new Set([
    ...(contract.location || []),
    'title', 'status', 'path',
    ...(contract.header || []).map((h) => h.field),
    ...(contract.sections || []).map((s) => s.field),
  ]);
  out.push('', '--index-entries (pass the same --data; a second invocation, no writes):');
  for (const e of contract.index_entries || []) {
    out.push(`  ${e.index}: ${e.key} keys a row into ${e.file} under "${e.section}" (scaffold: ${e.scaffold} when the file is absent)`);
    const rowFields = [...new Set([...String(e.row || '').matchAll(/\{(\w+)\}/g)].map((m) => m[1]))];
    const extra = rowFields.filter((f) => !declared.has(f));
    if (extra.length) {
      out.push(`    row needs ${extra.map((f) => `\`${f}\``).join(', ')} — supplied for this row only, absent from the artifact body above`);
    }
  }
  return out.join('\n') + '\n';
}

// --- cli ------------------------------------------------------------------

function arg(name, fallback) {
  const i = process.argv.indexOf(`--${name}`);
  return i === -1 ? fallback : process.argv[i + 1];
}

// Deployed, the contracts sit beside their canonical template under the user's ~/.claude; inside the
// payload repo they sit next to this script under `references/<x>/templates/`. Same probe pattern as
// render-return.js: the contract's own file, not the directory (a deployed dir can exist without it).
function contractPath(type) {
  const explicit = arg('refs');
  const [dir, file] = TYPES[type];
  if (explicit) return path.join(explicit, dir, 'templates', file);
  const home = path.join(require('os').homedir(), '.claude', 'references');
  if (fs.existsSync(path.join(home, dir, 'templates', file))) return path.join(home, dir, 'templates', file);
  return path.resolve(__dirname, '..', 'references', dir, 'templates', file);
}

function main() {
  const type = arg('type');
  if (!type || !TYPES[type]) {
    process.stderr.write(`usage: render-artifact.js --type <edr|capability-spec> [--schema | --data <file.json> [--index-entries]] [--refs <dir>]\n`);
    process.exit(2);
  }
  const file = contractPath(type);

  let contract;
  try {
    contract = yaml.parse(fs.readFileSync(file, 'utf8'));
  } catch (e) {
    bail(`cannot read the contract for \`${type}\` at ${file}: ${e.message}`);
  }

  if (process.argv.includes('--schema')) {
    process.stdout.write(schema(contract));
    return;
  }

  const dataFile = arg('data');
  let data;
  try {
    data = JSON.parse(dataFile ? fs.readFileSync(dataFile, 'utf8') : fs.readFileSync(0, 'utf8'));
  } catch (e) {
    bail(`cannot read the data: ${e.message}`);
  }

  if (process.argv.includes('--index-entries')) {
    process.stdout.write(`${JSON.stringify(indexEntries(contract, data), null, 2)}\n`);
    return;
  }
  process.stdout.write(render(contract, data));
}

main();
