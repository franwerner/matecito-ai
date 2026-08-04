#!/usr/bin/env node
'use strict';

// Checks a durable artifact store (EDR, capability-spec) for structural conformance: per-file shape
// against the same contract render-artifact.js used to produce it (`edr.yaml` / `capability.yaml`),
// plus the mechanical checks declared in `references/artifact-checks/checks.yaml` (index sync,
// dangling links, folder taxonomy, and more — see that file's README for the closed kind set).
// Read-only — it never writes.
//
// Usage:
//   validate-artifact.js --type <edr|capability-spec> (--file <path> | --store <dir>)
//                         [--root <dir>] [--refs <dir>] [--checks <file>]
//   validate-artifact.js --self-check [--refs <dir>] [--checks <file>]
//
// `--file` and `--store` are mutually exclusive; exactly one is required. `--file` runs only the
// per-file shape check. `--store` sweeps every artifact under the given directory and additionally
// runs the checks.yaml rows for this type; a non-conformant file never aborts the sweep.
//
// `--self-check` is a third scope, orthogonal to `--file`/`--store`: it checks the SHIPPED contracts,
// templates and `checks.yaml` against each other instead of a written artifact, so it needs none of
// `--type`, `--file`, `--store` or `--root`. It covers the two contract↔template pairs (edr, capability-
// spec), the checks.yaml rows against what their templates and taxonomy actually declare, and the
// `../` depth of the cross-store/same-store example links, derived from each contract's own `path:`.
//
// `--root` has no default. If any resolved check declares `requires: [root]` or `[config]` and
// `--root` is absent, this exits 2 naming those checks — it never silently degrades.
//
// Exit 0 when clean (or only warning/nota findings), 1 when any `error` finding exists, 2 when it
// could not run at all (bad args, unreadable contract/target, or a missing required `--root`).
// Findings are emitted as a single JSON record on stdout; on exit 2, stdout is empty and the reason
// is on stderr.

const fs = require('fs');
const path = require('path');
const yaml = require('./lib/yaml');
const libStore = require('./lib/store');

const TYPES = {
  edr: ['edr', 'edr.yaml'],
  'capability-spec': ['spec', 'capability.yaml'],
};

// Field labels used by `dangling-link` rows that read a header line instead of a section — only
// `Reemplazado por` needs this today, so a full generic mapping would be speculative.
const FIELD_LABELS = { replaced_by: 'Reemplazado por' };

const DEFAULT_IGNORE = new Set(['.git', 'node_modules']);

// Matches exactly two hashes, so a `### Scenario: ...` inside `## Escenarios` stays part of that
// section's body instead of being read as a sibling section.
const SECTION_HEADING = /^##\s+(.*)$/;
const HEADER_BULLET = /^-\s+\*\*([^*]+):\*\*\s*(.*)$/;

function parseArtifact(text) {
  const lines = text.split('\n');
  const header = {};
  let i = 0;
  while (i < lines.length && !SECTION_HEADING.test(lines[i])) {
    const m = lines[i].match(HEADER_BULLET);
    if (m) header[m[1].trim()] = m[2].trim();
    i++;
  }
  const sections = [];
  let current = null;
  for (; i < lines.length; i++) {
    const m = lines[i].match(SECTION_HEADING);
    if (m) {
      current = { title: lines[i].trim(), line: i + 1, body: [] };
      sections.push(current);
      continue;
    }
    if (current) current.body.push(lines[i]);
  }
  return { header, sections: sections.map((s) => ({ ...s, body: s.body.join('\n').trim() })) };
}

// --- per-file shape checks (edr.yaml / capability.yaml) --------------------

function checkHeader(contract, header, add) {
  const status = header['Status'];
  if (status === undefined) {
    add('error', 'STATUS-MISSING', 'no `- **Status:** ...` line found in the header');
  } else if (!(contract.statuses || []).includes(status)) {
    add('error', 'STATUS-ILLEGAL', `\`Status\` is "${status}", not one of ${JSON.stringify(contract.statuses)}`);
  }

  for (const h of contract.header || []) {
    const seen = header[h.label] !== undefined;
    const wantedByStatus = h.emitted !== 'on-status' || (status !== undefined && (h.statuses || []).includes(status));
    if (h.emitted !== 'when-present' && wantedByStatus && !seen) {
      add('error', 'HEADER-MISSING', `header \`${h.label}\` is required${h.emitted === 'on-status' ? ` on status \`${status}\`` : ''} and absent`);
    }
    if (h.emitted === 'on-status' && !wantedByStatus && seen) {
      add('error', 'HEADER-UNEXPECTED', `header \`${h.label}\` is only for status ${JSON.stringify(h.statuses)}, but \`Status\` is \`${status}\``);
    }
  }
  return status;
}

function checkSections(contract, sections, status, add) {
  const byTitle = new Map();
  for (const s of sections) {
    if (byTitle.has(s.title)) {
      add('warning', 'SECTION-DUPLICATE', `\`${s.title}\` appears more than once (line ${s.line})`, s.line);
      continue;
    }
    byTitle.set(s.title, s);
  }

  const declaredTitles = new Set((contract.sections || []).map((s) => s.title));

  for (const s of contract.sections || []) {
    const seen = byTitle.get(s.title);
    const wanted = s.emitted !== 'on-status' || (s.statuses || []).includes(status);

    if (s.emitted !== 'when-present' && wanted && !seen) {
      // `always` is unconditional — its absence is a real structural defect, `error`. `on-status`
      // is the renderer's production instruction ("emit only when this status applies"), not a
      // conformity requirement on a hand-written file — the original rubric's closest checks for
      // this case (section-set-symmetry, section-minimum) migrated at `warning`/nota, never `error`.
      const severity = s.emitted === 'on-status' ? 'warning' : 'error';
      add(severity, 'SECTION-MISSING', `\`${s.title}\` is required${s.emitted === 'on-status' ? ` on status \`${status}\`` : ''} and absent`);
      continue;
    }
    if (s.emitted === 'on-status' && !wanted && seen) {
      add('error', 'SECTION-UNEXPECTED', `\`${s.title}\` is only for status ${JSON.stringify(s.statuses)}, but \`Status\` is \`${status}\``, seen.line);
      continue;
    }
    if (seen && (s.empty_on || []).includes(status) && seen.body !== '') {
      add('error', 'SECTION-FORBIDDEN-CONTENT', `\`${s.title}\` carries content, but status \`${status}\` forbids reasoning there`, seen.line);
    }
  }

  // Order — reported, not failed, and reported as a `nota`: the rule is newer than the store (every
  // capability-spec written before it exists is free to disagree), and a `lifecycle` skeleton
  // legitimately opens with a section this store's declared order puts second. One nota is enough
  // signal; a cascade of follow-on mismatches after the first would just be noise from the same root
  // cause.
  const declaredOrder = (contract.sections || []).map((s) => s.title);
  const foundDeclared = sections.filter((s) => declaredTitles.has(s.title)).map((s) => s.title);
  const expected = declaredOrder.filter((t) => foundDeclared.includes(t));
  for (let i = 0; i < expected.length; i++) {
    if (expected[i] !== foundDeclared[i]) {
      add('nota', 'SECTION-ORDER', `\`${foundDeclared[i]}\` appears out of the declared order (expected \`${expected[i]}\` at this position)`, byTitle.get(foundDeclared[i]).line);
      break;
    }
  }

  // Unknown headings — a human added something the contract does not know about. Informational
  // only: the durable file is theirs to extend.
  for (const s of sections) {
    if (!declaredTitles.has(s.title)) {
      add('nota', 'SECTION-UNKNOWN', `\`${s.title}\` (line ${s.line}) is not declared by the contract`, s.line);
    }
  }
}

function validateArtifactText(contract, text) {
  const problems = [];
  const add = (severity, code, message, line = null) => problems.push({ severity, code, message, line });
  const { header, sections } = parseArtifact(text);
  const status = checkHeader(contract, header, add);
  if (status !== undefined && (contract.statuses || []).includes(status)) {
    checkSections(contract, sections, status, add);
  }
  return problems;
}

function validateFileShape(contract, text, filePath) {
  return validateArtifactText(contract, text).map((p) => ({ ...p, file: filePath }));
}

// --- checks.yaml engine -----------------------------------------------------

function fillTemplate(message, vars) {
  return message.replace(/\{(\w+)\}/g, (m, k) => (k in vars && vars[k] !== undefined ? String(vars[k]) : m));
}

function mkFinding(row, vars) {
  return {
    severity: row.severity,
    code: row.code,
    message: fillTemplate(row.message, vars),
    file: vars.file,
    line: vars.line !== undefined ? vars.line : null,
    related: vars.related || [],
  };
}

function isExternalHref(href) {
  return /^[a-z][a-z0-9+.-]*:\/\//i.test(href);
}

function folderOf(relPath) {
  return relPath.includes('/') ? relPath.split('/')[0] : null;
}

function backtickName(cell) {
  const m = (cell || '').trim().match(/^`([^`]+)`$/);
  return m ? m[1] : null;
}

function linkSlug(cell) {
  const links = libStore.extractLinks(cell || '');
  if (!links.length) return null;
  return links[0].href.replace(/\.md$/, '').split('/').pop();
}

function componentNames(configJson) {
  if (!configJson || !configJson.repo || !Array.isArray(configJson.repo.components)) return [];
  return configJson.repo.components.map((c) => c.name);
}

// A link found in `artifact` resolves relative to the store's own parent directory (both stores in
// this ecosystem live as siblings under `.matecito-ai/`), so cross-store links (`../edr/...`) and
// same-store links (`slug.md`, `../<domain>/slug.md`) resolve through the same path, uniformly.
function resolveStoreLink(ctx, artifact, href) {
  if (isExternalHref(href)) return true;
  const from = path.posix.join(ctx.ownBase, artifact.relPath);
  const targetRel = libStore.resolveRelativeLink(from, href);
  return fs.existsSync(path.join(ctx.commonParent, targetRel));
}

function findDirsNamed(root, name) {
  const out = [];
  const recurse = (dir) => {
    let entries;
    try {
      entries = fs.readdirSync(dir, { withFileTypes: true });
    } catch (e) {
      return;
    }
    for (const entry of entries) {
      if (!entry.isDirectory() || DEFAULT_IGNORE.has(entry.name)) continue;
      const full = path.join(dir, entry.name);
      if (entry.name === name) out.push(full);
      recurse(full);
    }
  };
  if (fs.existsSync(root)) recurse(root);
  return out;
}

function indexSyncCheck(row, ctx) {
  const out = [];
  const isEdr = row.store === 'edr';
  const folders = ctx.storeConfig.folders || [];
  if (row.level === 'file') {
    for (const folder of folders) {
      const idxText = ctx.indexes[`${folder}/INDEX.md`];
      const heading = isEdr ? '## EDRs en este dominio' : '## Capacidades';
      const rows = idxText ? libStore.extractTableRows(idxText, heading) : [];
      const inIndex = new Set(rows.map((r) => (isEdr ? linkSlug(r[0]) : backtickName(r[0]))).filter(Boolean));
      const onDisk = new Set(ctx.artifacts.filter((a) => a.folder === folder).map((a) => a.slug));
      for (const slug of onDisk) {
        if (!inIndex.has(slug)) out.push(mkFinding(row, { slug, capability: slug, domain: folder, type: folder, file: `${folder}/${slug}.md` }));
      }
      for (const slug of inIndex) {
        if (!onDisk.has(slug)) out.push(mkFinding(row, { slug, capability: slug, domain: folder, type: folder, file: `${folder}/INDEX.md` }));
      }
    }
  } else {
    const rootText = ctx.indexes['INDEX.md'];
    const heading = isEdr ? '## Dominios de este proyecto' : '## Tipos de este proyecto';
    const rows = rootText ? libStore.extractTableRows(rootText, heading) : [];
    const declared = new Set(rows.map((r) => backtickName(r[0])).filter(Boolean));
    const present = new Set(folders.filter((f) => ctx.indexes[`${f}/INDEX.md`] !== undefined || ctx.artifacts.some((a) => a.folder === f)));
    for (const f of present) {
      if (!declared.has(f)) out.push(mkFinding(row, { domain: f, type: f, file: 'INDEX.md' }));
    }
    for (const f of declared) {
      if (!present.has(f)) out.push(mkFinding(row, { domain: f, type: f, file: 'INDEX.md' }));
    }
  }
  return out;
}

const ENGINE = {
  'section-non-empty': (rows, ctx) => {
    const out = [];
    for (const row of rows) {
      for (const a of ctx.artifacts) {
        if (row.when_status && !row.when_status.includes(a.header['Status'])) continue;
        const body = libStore.extractSection(a.text, row.section);
        if (body !== null && body.trim() === '') {
          out.push(mkFinding(row, { file: a.relPath, status: a.header['Status'] }));
        }
      }
    }
    return out;
  },

  'section-required-when': (rows, ctx) => {
    const out = [];
    for (const row of rows) {
      for (const a of ctx.artifacts) {
        if (row.when_domain_in && !row.when_domain_in.includes(a.folder)) continue;
        if (row.when_status && !row.when_status.includes(a.header['Status'])) continue;
        if (libStore.extractSection(a.text, row.section) === null) out.push(mkFinding(row, { file: a.relPath }));
      }
    }
    return out;
  },

  'section-set-symmetry': (rows, ctx) => {
    const out = [];
    for (const row of rows) {
      const [secA, secB] = row.sections || [];
      for (const a of ctx.artifacts) {
        if (row.when_status && !row.when_status.includes(a.header['Status'])) continue;
        const nonEmpty = (b) => b !== null && b.trim() !== '';
        const bodyA = libStore.extractSection(a.text, secA);
        const bodyB = libStore.extractSection(a.text, secB);
        if (nonEmpty(bodyA) !== nonEmpty(bodyB)) out.push(mkFinding(row, { file: a.relPath }));
      }
    }
    return out;
  },

  'section-minimum': (rows, ctx) => {
    const out = [];
    for (const row of rows) {
      for (const a of ctx.artifacts) {
        if (row.when_status && !row.when_status.includes(a.header['Status'])) continue;
        const hasAny = (row.required_any || []).some((sec) => {
          const b = libStore.extractSection(a.text, sec);
          return b !== null && b.trim() !== '';
        });
        if (!hasAny) out.push(mkFinding(row, { file: a.relPath }));
      }
    }
    return out;
  },

  'bullet-prefix': (rows, ctx) => {
    const out = [];
    for (const row of rows) {
      for (const a of ctx.artifacts) {
        const body = libStore.extractSection(a.text, row.section);
        if (body === null) continue;
        const items = body.split('\n').filter((l) => /^- /.test(l));
        const bad = items.some((it) => !(row.prefix_any || []).some((p) => it.slice(2).startsWith(p)));
        if (bad) out.push(mkFinding(row, { file: a.relPath }));
      }
    }
    return out;
  },

  'legacy-status': (rows, ctx) => {
    const out = [];
    for (const row of rows) {
      for (const a of ctx.artifacts) {
        if (a.header[row.header] === row.value) out.push(mkFinding(row, { file: a.relPath }));
      }
    }
    return out;
  },

  'scenario-shape': (rows, ctx) => {
    const out = [];
    for (const row of rows) {
      for (const a of ctx.artifacts) {
        const body = libStore.extractSection(a.text, row.section);
        if (body === null) continue;
        const blocks = body.split(/^###\s+/m).filter((b) => b.trim() !== '');
        const scan = blocks.length ? blocks : [body];
        for (const block of scan) {
          const missing = (row.requires_shape || []).some((tok) => !new RegExp(`\\*\\*${tok}\\*\\*`, 'i').test(block));
          if (missing) {
            out.push(mkFinding(row, { file: a.relPath }));
            break;
          }
        }
      }
    }
    return out;
  },

  'skeleton-sections': (rows, ctx) => {
    const out = [];
    for (const row of rows) {
      for (const a of ctx.artifacts) {
        if (a.folder !== row.type) continue;
        if (row.when_status && !row.when_status.includes(a.header['Status'])) continue;
        for (const sec of row.sections || []) {
          const body = libStore.extractSection(a.text, sec);
          if (body === null || body.trim() === '') out.push(mkFinding(row, { file: a.relPath, section: sec }));
        }
      }
    }
    return out;
  },

  'header-in-set': (rows, ctx) => {
    const out = [];
    const names = componentNames(ctx.configJson);
    for (const row of rows) {
      for (const a of ctx.artifacts) {
        const raw = a.header[row.header];
        if (raw === undefined) continue;
        for (const v of raw.split(',').map((s) => s.trim()).filter(Boolean)) {
          if (!names.includes(v)) out.push(mkFinding(row, { file: a.relPath, value: v }));
        }
      }
    }
    return out;
  },

  'header-required-when-source': (rows, ctx) => {
    const out = [];
    for (const row of rows) {
      for (const a of ctx.artifacts) {
        if (a.header[row.header] === undefined) out.push(mkFinding(row, { file: a.relPath }));
      }
    }
    return out;
  },

  'glob-drift': (rows, ctx) => {
    const out = [];
    const rootFiles = libStore.walk(ctx.root);
    for (const row of rows) {
      for (const a of ctx.artifacts) {
        const body = libStore.extractSection(a.text, row.section);
        if (body === null) continue;
        for (const line of body.split('\n')) {
          const m = line.match(/^-\s+`([^`]+)`/);
          if (!m) continue;
          const glob = m[1];
          const hit = rootFiles.some((f) => libStore.matchesGlob(glob, f));
          if (!hit) out.push(mkFinding(row, { file: a.relPath, glob }));
        }
      }
    }
    return out;
  },

  'folder-taxonomy': (rows, ctx) => {
    const out = [];
    let dirs;
    try {
      dirs = fs.readdirSync(ctx.storeDir, { withFileTypes: true }).filter((e) => e.isDirectory()).map((e) => e.name);
    } catch (e) {
      dirs = [];
    }
    const canonical = new Set(ctx.storeConfig.folders || []);
    const extra = new Set(ctx.storeConfig.extra_folders || []);
    for (const row of rows) {
      for (const dir of dirs) {
        if (DEFAULT_IGNORE.has(dir) || canonical.has(dir) || extra.has(dir)) continue;
        out.push(mkFinding(row, { folder: dir, file: dir }));
      }
    }
    return out;
  },

  'index-sync': (rows, ctx) => rows.flatMap((row) => indexSyncCheck(row, ctx)),

  'location-by-slug': (rows, ctx) => {
    const out = [];
    for (const row of rows) {
      for (const a of ctx.artifacts) {
        const expected = (ctx.storeConfig.slug_domains || {})[a.slug];
        if (expected && expected !== a.folder) {
          out.push(mkFinding(row, { file: a.relPath, slug: a.slug, expected_domain: expected, folder: a.folder }));
        }
      }
    }
    return out;
  },

  'dangling-link': (rows, ctx) => {
    const out = [];
    for (const row of rows) {
      for (const a of ctx.artifacts) {
        if (row.when_status && !row.when_status.includes(a.header['Status'])) continue;
        let links = [];
        if (row.field) {
          const raw = a.header[FIELD_LABELS[row.field] || row.field];
          if (raw) links = libStore.extractLinks(raw);
        } else if (row.section) {
          const body = libStore.extractSection(a.text, row.section);
          if (body) links = libStore.extractLinks(body);
        }
        for (const link of links) {
          if (!resolveStoreLink(ctx, a, link.href)) out.push(mkFinding(row, { file: a.relPath, href: link.href }));
        }
      }
    }
    return out;
  },

  'orphan-folder': (rows, ctx) => {
    const out = [];
    for (const row of rows) {
      for (const folder of ctx.storeConfig.folders || []) {
        const hasIndex = ctx.indexes[`${folder}/INDEX.md`] !== undefined;
        const hasArtifacts = ctx.artifacts.some((a) => a.folder === folder);
        if (hasIndex && !hasArtifacts) out.push(mkFinding(row, { domain: folder, type: folder, folder, file: `${folder}/INDEX.md` }));
      }
    }
    return out;
  },

  'inbound-reference-count': (rows, ctx) => {
    const out = [];
    const graph = libStore.buildInboundReferenceGraph(ctx.artifacts.map((a) => ({ path: a.relPath, text: a.text })));
    for (const row of rows) {
      for (const a of ctx.artifacts) {
        if (!(row.applies_to_types || []).includes(a.folder)) continue;
        const refs = graph.get(a.relPath) || [];
        const referrers = new Set(refs.filter((r) => (row.referenced_by_types || []).includes(folderOf(r.from))).map((r) => r.from));
        if (referrers.size >= 1 && referrers.size < (row.min_referrers || 2)) {
          out.push(mkFinding(row, { capability: a.slug, type: a.folder, referrer: [...referrers][0], file: a.relPath }));
        }
      }
    }
    return out;
  },

  'deprecated-referenced': (rows, ctx) => {
    const out = [];
    const graph = libStore.buildInboundReferenceGraph(ctx.artifacts.map((a) => ({ path: a.relPath, text: a.text })));
    const byPath = new Map(ctx.artifacts.map((a) => [a.relPath, a]));
    for (const row of rows) {
      for (const a of ctx.artifacts) {
        if (row.when_status && !row.when_status.includes(a.header['Status'])) continue;
        const refs = graph.get(a.relPath) || [];
        const liveReferrer = refs.find((r) => {
          const from = byPath.get(r.from);
          return from && from.header['Status'] !== 'Deprecated';
        });
        if (liveReferrer) out.push(mkFinding(row, { file: a.relPath, referrer: liveReferrer.from }));
      }
    }
    return out;
  },

  'store-uniqueness': (rows, ctx) => {
    const out = [];
    for (const row of rows) {
      const storeRoot = ctx.storeConfig.root || '';
      const basename = path.basename(storeRoot);
      const canonicalAbs = path.resolve(ctx.root, storeRoot);
      const found = findDirsNamed(ctx.root, basename).filter((p) => path.resolve(p) !== canonicalAbs);
      for (const p of found) out.push(mkFinding(row, { path: path.relative(ctx.root, p), file: path.relative(ctx.root, p) }));
    }
    return out;
  },
};

// --- store loading ----------------------------------------------------------

function loadStore(storeDir) {
  const artifacts = [];
  const indexes = {};
  for (const rel of libStore.walk(storeDir)) {
    if (!rel.endsWith('.md')) continue;
    const full = path.join(storeDir, ...rel.split('/'));
    const text = fs.readFileSync(full, 'utf8');
    if (path.posix.basename(rel) === 'INDEX.md') {
      indexes[rel] = text;
      continue;
    }
    const segs = rel.split('/');
    const folder = segs.length > 1 ? segs[0] : null;
    const slug = path.posix.basename(rel).replace(/\.md$/, '');
    const { header } = parseArtifact(text);
    artifacts.push({ relPath: rel, folder, slug, text, header });
  }
  return { artifacts, indexes };
}

// --- cli ---------------------------------------------------------------------

function arg(name, fallback) {
  const i = process.argv.indexOf(`--${name}`);
  return i === -1 ? fallback : process.argv[i + 1];
}

function contractPath(type) {
  const explicit = arg('refs');
  const [dir, file] = TYPES[type];
  if (explicit) return path.join(explicit, dir, 'templates', file);
  const home = path.join(require('os').homedir(), '.claude', 'references');
  if (fs.existsSync(path.join(home, dir, 'templates', file))) return path.join(home, dir, 'templates', file);
  return path.resolve(__dirname, '..', 'references', dir, 'templates', file);
}

// Deployed, the checks contract sits in the user's ~/.claude; inside the payload repo it sits next
// to this script. Falling back to the sibling keeps `--checks` genuinely optional in both places —
// same pattern as `contractPath` above and `validate-return.js`'s `contractsDir`.
function checksPath() {
  const explicit = arg('checks');
  if (explicit) return explicit;
  const home = path.join(require('os').homedir(), '.claude', 'references', 'artifact-checks', 'checks.yaml');
  if (fs.existsSync(home)) return home;
  return path.resolve(__dirname, '..', 'references', 'artifact-checks', 'checks.yaml');
}

// Same precedence as `contractPath`/`checksPath` (explicit `--refs`, then the deployed home, then the
// repo-local sibling), collapsed to the base directory — surfaced on stderr so someone editing the
// payload in-repo notices when a bare `--self-check` is actually reading the deployed copy instead.
function refsRootPath() {
  const explicit = arg('refs');
  if (explicit) return explicit;
  const home = path.join(require('os').homedir(), '.claude', 'references');
  if (fs.existsSync(home)) return home;
  return path.resolve(__dirname, '..', 'references');
}

const USAGE = 'usage: validate-artifact.js --type <edr|capability-spec> (--file <path> | --store <dir>) [--root <dir>] [--refs <dir>] [--checks <file>]\n       validate-artifact.js --self-check [--refs <dir>] [--checks <file>]\n';

// --- self-check -------------------------------------------------------------
// The contracts (`edr.yaml`, `capability.yaml`) are only safe as a second source of truth for
// `render-artifact.js` / `validate-artifact.js` while this proves they still agree with their
// templates and with `checks.yaml`, which references both by literal text.

// A contract's `title:`/`path:` carry placeholder slots (`{title}`, `{domain}`) that a template
// renders as its own placeholder syntax (`<título>`); comparing them literally would always fail, so
// a placeholder slot becomes a wildcard before matching against the template's lines.
function templateHasPattern(pattern, md) {
  const escaped = pattern.replace(/[.*+?^${}()|[\]\\]/g, '\\$&').replace(/\\\{\w+\\\}/g, '.+');
  const re = new RegExp(`^${escaped}$`);
  return md.split('\n').some((l) => re.test(l.trim()));
}

// Drift between one contract and its template: a declared `title:`, `header:` label or
// `sections[].title` that appears nowhere in the template; a `##` heading the contract does not
// declare; or a `statuses:` value the template's inline status enum omits, in either direction.
function pairDrift(contract, md, yamlPath, mdPath) {
  const problems = [];
  const add = (file, code, message) => problems.push({ severity: 'error', code, message, file });

  if (!templateHasPattern(contract.title, md)) {
    add(yamlPath, 'DRIFT-TITLE', `the artifact title \`${contract.title}\` appears nowhere in the template`);
  }

  for (const h of contract.header || []) {
    if (!md.includes(`**${h.label}`)) {
      add(yamlPath, 'DRIFT-HEADER', `header label \`${h.label}\` is declared in the contract but appears nowhere in the template`);
    }
  }

  const headings = new Set(md.split('\n').map((l) => l.trim()).filter((l) => /^##\s+/.test(l)));
  for (const s of contract.sections || []) {
    if (!headings.has(s.title)) {
      add(yamlPath, 'DRIFT-SECTION', `section \`${s.title}\` is declared in the contract but appears nowhere in the template`);
    }
  }
  const declaredTitles = new Set((contract.sections || []).map((s) => s.title));
  for (const h of headings) {
    if (!declaredTitles.has(h)) {
      add(mdPath, 'DRIFT-SECTION', `template heading \`${h}\` is not declared by the contract's \`sections:\``);
    }
  }

  // The status enum is prose (`**Status:** <A | B | C>`), not a rendered section — the fragile half
  // of this check, since nothing forces it to stay literal the way a `##` heading does.
  const enumMatch = md.match(/\*\*Status:\*\*\s*<([^>]+)>/);
  if (!enumMatch) {
    add(mdPath, 'DRIFT-STATUS', 'no inline `**Status:** <...>` enum found in the template to compare against `statuses:`');
  } else {
    const templateStatuses = enumMatch[1].split('|').map((s) => s.trim());
    const declared = contract.statuses || [];
    for (const s of declared) {
      if (!templateStatuses.includes(s)) add(yamlPath, 'DRIFT-STATUS', `\`statuses:\` declares \`${s}\`, which the template's status enum omits`);
    }
    for (const s of templateStatuses) {
      if (!declared.includes(s)) add(mdPath, 'DRIFT-STATUS', `the template's status enum declares \`${s}\`, which \`statuses:\` omits`);
    }
  }

  return problems;
}

// checks.yaml rows name section titles and taxonomy folders by literal text (never through the
// contract), so a rename on either side leaves the row referencing something that no longer exists —
// and a row that never matches never fires, silently.
function checksDrift(checksConfig, contracts, checksFile) {
  const problems = [];
  const add = (code, message) => problems.push({ severity: 'error', code, message, file: checksFile });

  const sectionTitlesOf = (store) => new Set(((contracts[store] || {}).sections || []).map((s) => s.title));
  const foldersOf = (store) => {
    const cfg = (checksConfig.stores || {})[store] || {};
    return new Set([...(cfg.folders || []), ...(cfg.extra_folders || [])]);
  };

  for (const row of checksConfig.checks || []) {
    const sections = sectionTitlesOf(row.store);
    const sectionRefs = [].concat(row.section || [], row.sections || [], row.required_any || []);
    for (const title of sectionRefs) {
      if (!sections.has(title)) add('DRIFT-CHECKS-SECTION', `row \`${row.code}\` references section \`${title}\`, which \`${row.store}\`'s contract does not declare`);
    }

    const folders = foldersOf(row.store);
    const folderRefs = [].concat(row.type || [], row.applies_to_types || [], row.referenced_by_types || [], row.when_domain_in || []);
    for (const name of folderRefs) {
      if (!folders.has(name)) add('DRIFT-CHECKS-FOLDER', `row \`${row.code}\` references folder \`${name}\`, absent from \`${row.store}\`'s taxonomy`);
    }
  }

  // Store-level taxonomy metadata that no single row owns: the per-type skeleton section lists and
  // the slug->domain catalog.
  for (const [store, cfg] of Object.entries(checksConfig.stores || {})) {
    const folders = foldersOf(store);
    const declaredSections = sectionTitlesOf(store);
    for (const [key, sections] of Object.entries(cfg.skeleton || {})) {
      if (!folders.has(key)) add('DRIFT-CHECKS-FOLDER', `stores.${store}.skeleton declares key \`${key}\`, absent from \`${store}\`'s taxonomy`);
      for (const title of sections) {
        if (!declaredSections.has(title)) add('DRIFT-CHECKS-SECTION', `stores.${store}.skeleton.${key} references section \`${title}\`, which \`${store}\`'s contract does not declare`);
      }
    }
    for (const domain of Object.values(cfg.slug_domains || {})) {
      if (!folders.has(domain)) add('DRIFT-CHECKS-FOLDER', `stores.${store}.slug_domains references folder \`${domain}\`, absent from \`${store}\`'s taxonomy`);
    }
  }

  return problems;
}

// Directory segments a contract's `path:` puts between `.matecito-ai/` and the filename, e.g.
// ".matecito-ai/edr/{domain}/{slug}.md" -> ["edr", "{domain}"].
function pathDirs(contract) {
  return contract.path.split('/').slice(1, -1);
}

function isPlaceholderSegment(seg) {
  return /^\{[\w-]+\}$/.test(seg);
}

// How many `../` a link from `source` to `target` must carry: climb from `source`'s own directory
// past whatever it does NOT literally share with `target`, then descend into `target`. A placeholder
// segment is never "shared", even when the two contracts spell it the same way (`{domain}` ==
// `{domain}`) — it stands for a value that varies per artifact, so two files of the very same contract
// still need to climb past it to reach each other.
function expectedDepth(source, target) {
  const sourceDirs = pathDirs(source);
  const targetDirs = pathDirs(target);
  let shared = 0;
  while (
    shared < sourceDirs.length && shared < targetDirs.length &&
    !isPlaceholderSegment(sourceDirs[shared]) && !isPlaceholderSegment(targetDirs[shared]) &&
    sourceDirs[shared] === targetDirs[shared]
  ) shared += 1;
  return sourceDirs.length - shared;
}

// Example links are placeholders, never resolved against the filesystem — only their `../` depth is
// checked, derived from the `path:` values the source and target contracts already declare. Driven by
// checks.yaml's own `dangling-link` rows (the ones with a `section:`), so the two link-carrying
// sections this covers are discovered, not hardcoded twice.
function linkDepthDrift(checksConfig, contracts, templatePaths, templates) {
  const problems = [];
  const rows = (checksConfig.checks || []).filter((r) => r.kind === 'dangling-link' && r.section);
  for (const row of rows) {
    const source = contracts[row.store];
    const target = contracts[row.target_store];
    const md = templates[row.store];
    if (!source || !target || md === undefined) continue;
    const body = libStore.extractSection(md, row.section);
    if (body === null) continue;
    for (const link of libStore.extractLinks(body)) {
      const up = link.href.match(/^(\.\.\/)*/)[0];
      const depthFound = up.length / 3;
      const remainder = link.href.slice(up.length).split('/');
      if (remainder.length <= 1) continue; // same-directory example — no location to derive
      const want = expectedDepth(source, target);
      if (depthFound !== want) {
        problems.push({
          severity: 'error',
          code: 'DRIFT-LINK-DEPTH',
          message: `\`${row.section}\` links to \`${link.href}\` (${depthFound} level(s) up), but the declared \`path:\` values put it ${want} level(s) up`,
          file: templatePaths[row.store],
        });
      }
    }
  }
  return problems;
}

function runSelfCheck(checksConfig, checksFile) {
  const contracts = {};
  const templatePaths = {};
  const templates = {};
  for (const type of Object.keys(TYPES)) {
    const yamlPath = contractPath(type);
    try {
      contracts[type] = yaml.parse(fs.readFileSync(yamlPath, 'utf8'));
    } catch (e) {
      process.stderr.write(`cannot read the contract for \`${type}\` at ${yamlPath}: ${e.message}\n`);
      process.exit(2);
    }
    const mdPath = yamlPath.replace(/\.yaml$/, '.md');
    templatePaths[type] = mdPath;
    try {
      templates[type] = fs.readFileSync(mdPath, 'utf8');
    } catch (e) {
      process.stderr.write(`cannot read the template for \`${type}\` at ${mdPath}: ${e.message}\n`);
      process.exit(2);
    }
  }

  const problems = [];
  for (const type of Object.keys(TYPES)) {
    problems.push(...pairDrift(contracts[type], templates[type], contractPath(type), templatePaths[type]));
  }
  problems.push(...checksDrift(checksConfig, contracts, checksFile));
  problems.push(...linkDepthDrift(checksConfig, contracts, templatePaths, templates));
  return problems;
}

function main() {
  if (process.argv.includes('--self-check')) {
    process.stderr.write(`self-check: comparing contracts against ${refsRootPath()}\n`);
    const checksFile = checksPath();
    let checksConfig;
    try {
      checksConfig = yaml.parse(fs.readFileSync(checksFile, 'utf8'));
    } catch (e) {
      process.stderr.write(`cannot read the checks contract at ${checksFile}: ${e.message}\n`);
      process.exit(2);
    }
    const severityMap = checksConfig.severity_map || {};

    const findings = runSelfCheck(checksConfig, checksFile).map((f) => ({
      severity: f.severity,
      display: severityMap[f.severity] || f.severity.toUpperCase(),
      code: f.code,
      file: f.file,
      line: f.line === undefined ? null : f.line,
      related: f.related || [],
      message: f.message,
    }));
    const counts = { error: 0, warning: 0, nota: 0 };
    for (const f of findings) counts[f.severity] = (counts[f.severity] || 0) + 1;

    const record = {
      tool: 'validate-artifact',
      type: null,
      mode: 'self-check',
      root: null,
      target: null,
      scanned: Object.keys(TYPES).length,
      findings,
      counts,
      skipped: [],
    };
    process.stdout.write(`${JSON.stringify(record)}\n`);
    process.exit(counts.error === 0 ? 0 : 1);
  }

  const type = arg('type');
  if (!type || !TYPES[type]) {
    process.stderr.write(USAGE);
    process.exit(2);
  }

  const fileArg = arg('file');
  const storeArg = arg('store');
  if ((fileArg && storeArg) || (!fileArg && !storeArg)) {
    process.stderr.write(`${USAGE}--file and --store are mutually exclusive; exactly one is required\n`);
    process.exit(2);
  }

  let artifactContract;
  const artifactContractFile = contractPath(type);
  try {
    artifactContract = yaml.parse(fs.readFileSync(artifactContractFile, 'utf8'));
  } catch (e) {
    process.stderr.write(`cannot read the contract for \`${type}\` at ${artifactContractFile}: ${e.message}\n`);
    process.exit(2);
  }

  let checksConfig;
  const checksFile = checksPath();
  try {
    checksConfig = yaml.parse(fs.readFileSync(checksFile, 'utf8'));
  } catch (e) {
    process.stderr.write(`cannot read the checks contract at ${checksFile}: ${e.message}\n`);
    process.exit(2);
  }

  const severityMap = checksConfig.severity_map || {};
  const storeConfig = (checksConfig.stores || {})[type] || {};

  let findings = [];
  const skipped = [];
  let scanned = 0;
  let mode;
  let target;
  let rootOut = null;

  if (fileArg) {
    mode = 'file';
    target = fileArg;
    let text;
    try {
      text = fs.readFileSync(fileArg, 'utf8');
    } catch (e) {
      process.stderr.write(`cannot read \`${fileArg}\`: ${e.message}\n`);
      process.exit(2);
    }
    findings = validateFileShape(artifactContract, text, fileArg);
    scanned = 1;
  } else {
    mode = 'store';
    target = storeArg;
    const rootArg = arg('root');
    const resolvedRows = (checksConfig.checks || []).filter((r) => r.store === type);
    const needingRoot = resolvedRows.filter((r) => (r.requires || []).some((x) => x === 'root' || x === 'config'));
    if (needingRoot.length && !rootArg) {
      process.stderr.write(`cannot run without --root — required by: ${needingRoot.map((r) => r.code).join(', ')}\n`);
      process.exit(2);
    }
    rootOut = needingRoot.length ? rootArg : null;

    let configJson = null;
    let configReadError = null;
    if (rootArg) {
      try {
        configJson = JSON.parse(fs.readFileSync(path.join(rootArg, '.matecito-ai', 'config.json'), 'utf8'));
      } catch (e) {
        configReadError = e;
      }
    }
    const configOk = !!(configJson && configJson.repo && Array.isArray(configJson.repo.components) && configJson.repo.components.length);

    if (!fs.existsSync(storeArg)) {
      process.stderr.write(`cannot read the store at ${storeArg}\n`);
      process.exit(2);
    }

    const { artifacts, indexes } = loadStore(storeArg);
    scanned = artifacts.length;

    // `extra_folders` (e.g. `tech/` in the EDR store) are legal store contents with no shape
    // contract of their own — `edr.yaml` describes a domain EDR, not a tech mini-EDR — and no check
    // here treats them as one of the canonical domains/types (see checks.yaml's `extra_folders`
    // comment). Shape-checking them against the wrong contract would flag every real tech EDR as
    // broken, so both the per-file shape check and the checks.yaml engine skip them.
    const extraFolders = new Set(storeConfig.extra_folders || []);
    const canonicalArtifacts = artifacts.filter((a) => !extraFolders.has(a.folder));
    for (const a of canonicalArtifacts) findings.push(...validateFileShape(artifactContract, a.text, a.relPath));

    const ctx = {
      type,
      storeDir: storeArg,
      root: rootArg,
      storeConfig,
      artifacts: canonicalArtifacts,
      indexes,
      ownBase: path.basename(storeArg),
      commonParent: path.dirname(storeArg),
      configJson,
    };

    const byKind = {};
    for (const row of resolvedRows) (byKind[row.kind] = byKind[row.kind] || []).push(row);

    for (const kind of Object.keys(byKind)) {
      const engineFn = ENGINE[kind];
      if (!engineFn) continue; // closed kind set — a row of an unknown kind is a checks.yaml defect, not this script's job to guess at
      const rows = byKind[kind];
      const configRows = rows.filter((r) => (r.requires || []).includes('config'));
      const otherRows = rows.filter((r) => !(r.requires || []).includes('config'));
      if (configRows.length) {
        if (configOk) {
          findings.push(...engineFn(configRows, ctx));
        } else {
          const reason = configReadError ? '.matecito-ai/config.json no existe bajo --root' : '.matecito-ai/config.json no declara repo.components';
          for (const r of configRows) skipped.push({ check: r.code, reason });
        }
      }
      if (otherRows.length) findings.push(...engineFn(otherRows, ctx));
    }
  }

  findings = findings.map((f) => ({
    severity: f.severity,
    display: severityMap[f.severity] || f.severity.toUpperCase(),
    code: f.code,
    file: f.file,
    line: f.line === undefined ? null : f.line,
    related: f.related || [],
    message: f.message,
  }));

  const counts = { error: 0, warning: 0, nota: 0 };
  for (const f of findings) counts[f.severity] = (counts[f.severity] || 0) + 1;

  const record = {
    tool: 'validate-artifact',
    type,
    mode,
    root: rootOut,
    target,
    scanned,
    findings,
    counts,
    skipped,
  };

  process.stdout.write(`${JSON.stringify(record)}\n`);
  process.exit(counts.error === 0 ? 0 : 1);
}

main();
