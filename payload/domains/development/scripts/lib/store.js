'use strict';

// Store-level primitives shared by the artifact checks engine: walking a store directory, matching
// the globs an EDR's `## Alcance` declares, and reading the two things store-wide checks need out of
// markdown (INDEX table rows, links). No external deps — this runs from the user's `~/.claude/scripts`
// with zero installs, same constraint as `lib/yaml.js`.

const fs = require('fs');
const path = require('path');

const DEFAULT_IGNORE = new Set(['.git', 'node_modules']);

// Real files under `dir`, recursive, as POSIX paths relative to `dir`. A missing `dir` is a store
// that has not been materialized yet, not an error — the caller (a store sweep) treats it as empty.
function walk(dir, opts = {}) {
  const ignore = opts.ignore ? new Set(opts.ignore) : DEFAULT_IGNORE;
  const out = [];
  const recurse = (current) => {
    for (const entry of fs.readdirSync(current, { withFileTypes: true })) {
      if (ignore.has(entry.name)) continue;
      const full = path.join(current, entry.name);
      if (entry.isDirectory()) recurse(full);
      else if (entry.isFile()) out.push(path.relative(dir, full).split(path.sep).join('/'));
    }
  };
  if (fs.existsSync(dir)) recurse(dir);
  return out.sort();
}

const REGEXP_SPECIAL = /[.+^${}()|[\]\\]/;

// `**` crosses path segments (matches `/`), `*` and `?` stay within one segment — the same split every
// `## Alcance` glob in this repo already assumes. No brace expansion, no character classes: nothing in
// the real corpus needs them, so the matcher does not carry them either.
function globToRegExp(pattern) {
  let re = '';
  for (let i = 0; i < pattern.length; i++) {
    const c = pattern[i];
    if (c === '*' && pattern[i + 1] === '*') {
      re += '.*';
      i++;
    } else if (c === '*') {
      re += '[^/]*';
    } else if (c === '?') {
      re += '[^/]';
    } else if (REGEXP_SPECIAL.test(c)) {
      re += `\\${c}`;
    } else {
      re += c;
    }
  }
  return new RegExp(`^${re}$`);
}

function matchesGlob(pattern, filePath) {
  return globToRegExp(pattern).test(filePath);
}

// The body lines of the first `heading` found, up to the next heading at the same or shallower level
// (`##` closes on the next `##` or `#`, never on a `###` inside it). `null` when the heading is absent.
function extractSection(markdown, heading) {
  const lines = markdown.split('\n');
  const level = (heading.match(/^#+/) || ['#'])[0].length;
  const start = lines.findIndex((l) => l.trim() === heading.trim());
  if (start === -1) return null;
  let end = lines.length;
  for (let i = start + 1; i < lines.length; i++) {
    const m = lines[i].match(/^(#+)\s/);
    if (m && m[1].length <= level) {
      end = i;
      break;
    }
  }
  return lines.slice(start + 1, end).join('\n');
}

function isSeparatorRow(cells) {
  return cells.every((c) => /^:?-+:?$/.test(c));
}

// Data rows of the first markdown table under `heading` — the header row and the `|---|---|`
// separator are consumed, never returned. `[]` when the heading or its table is absent.
function extractTableRows(markdown, heading) {
  const body = extractSection(markdown, heading);
  if (body === null) return [];
  const tableLines = body
    .split('\n')
    .map((l) => l.trim())
    .filter((l) => l.startsWith('|'));
  const rows = tableLines
    .map((l) => l.slice(1, l.endsWith('|') ? -1 : undefined).split('|').map((c) => c.trim()))
    .filter((cells) => !isSeparatorRow(cells));
  return rows.slice(1); // drop the header row
}

// `[text](href)` occurrences, in order. Cheap on purpose: no title-attribute, no reference-style
// `[text][ref]` — neither form appears in any durable template this reads.
function extractLinks(text) {
  const out = [];
  const re = /\[([^\]]*)\]\(([^)]+)\)/g;
  let m;
  while ((m = re.exec(text))) out.push({ text: m[1], href: m[2] });
  return out;
}

function isExternalHref(href) {
  return /^[a-z][a-z0-9+.-]*:\/\//i.test(href);
}

// Resolves a link found in `fromPath` against the directory it lives in, POSIX-style — the same
// resolution a browser or a markdown viewer would do for a relative link.
function resolveRelativeLink(fromPath, href) {
  const dir = path.posix.dirname(fromPath);
  return path.posix.normalize(path.posix.join(dir, href));
}

// `files`: `[{ path, text }]`, `path` POSIX-relative to the store's common root. Returns a `Map` from
// resolved target path to the list of `{ from, text, href }` links that point at it — one entry per
// distinct target, so both "does anything link here" (dangling-link) and "how many things link here"
// (inbound-reference-count) read off the same graph.
function buildInboundReferenceGraph(files) {
  const graph = new Map();
  for (const file of files) {
    for (const link of extractLinks(file.text)) {
      if (isExternalHref(link.href)) continue;
      const target = resolveRelativeLink(file.path, link.href);
      if (!graph.has(target)) graph.set(target, []);
      graph.get(target).push({ from: file.path, text: link.text, href: link.href });
    }
  }
  return graph;
}

module.exports = {
  walk,
  globToRegExp,
  matchesGlob,
  extractSection,
  extractTableRows,
  extractLinks,
  resolveRelativeLink,
  buildInboundReferenceGraph,
};
