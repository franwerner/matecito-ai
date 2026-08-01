'use strict';

// Hand-rolled because these scripts must run with zero installs in the user's ~/.claude. Covers only
// what the phase contracts use: scalars, inline arrays, nested maps, arrays of maps.

// Section titles are `"### Summary"`, so a naive comment stripper would eat them at the first `#`.
function stripComment(line) {
  let quote = null;
  for (let i = 0; i < line.length; i++) {
    const c = line[i];
    if (quote) {
      if (c === quote) quote = null;
      continue;
    }
    if (c === '"' || c === "'") quote = c;
    else if (c === '#') return line.slice(0, i).replace(/\s+$/, '');
  }
  return line.replace(/\s+$/, '');
}

function parseScalar(raw) {
  const v = raw.trim();
  if (v === '') return null;
  if (v === 'true') return true;
  if (v === 'false') return false;
  if (v === 'null') return null;
  if (/^-?\d+$/.test(v)) return Number(v);
  if ((v.startsWith('"') && v.endsWith('"')) || (v.startsWith("'") && v.endsWith("'"))) {
    return v.slice(1, -1);
  }
  if (v.startsWith('[') && v.endsWith(']')) {
    const inner = v.slice(1, -1).trim();
    return inner === '' ? [] : splitInline(inner).map(parseScalar);
  }
  if (v.startsWith('{') && v.endsWith('}')) {
    const inner = v.slice(1, -1).trim();
    if (inner === '') return {};
    const out = {};
    for (const pair of splitInline(inner)) {
      const sep = keyEnd(pair);
      if (sep === -1) continue;
      out[pair.slice(0, sep).trim()] = parseScalar(pair.slice(sep + 1));
    }
    return out;
  }
  return v;
}

// Commas inside quotes are data, not separators.
function splitInline(s) {
  const out = [];
  let cur = '';
  let quote = null;
  for (const c of s) {
    if (quote) {
      cur += c;
      if (c === quote) quote = null;
      continue;
    }
    if (c === '"' || c === "'") { quote = c; cur += c; continue; }
    if (c === ',') { out.push(cur); cur = ''; continue; }
    cur += c;
  }
  out.push(cur);
  return out.map((x) => x.trim()).filter((x) => x !== '');
}

function keyEnd(text) {
  let quote = null;
  for (let i = 0; i < text.length; i++) {
    const c = text[i];
    if (quote) { if (c === quote) quote = null; continue; }
    if (c === '"' || c === "'") { quote = c; continue; }
    if (c === ':') return i;
  }
  return -1;
}

function toLines(text) {
  return text
    .split('\n')
    .map(stripComment)
    .map((l) => ({ indent: l.search(/\S/), text: l.trim() }))
    .filter((l) => l.text !== '');
}

function parseNode(lines, start, indent) {
  if (start >= lines.length) return [null, start];
  if (lines[start].text.startsWith('- ')) return parseList(lines, start, indent);
  return parseMap(lines, start, indent);
}

function parseList(lines, start, indent) {
  const out = [];
  let i = start;
  while (i < lines.length && lines[i].indent === indent && lines[i].text.startsWith('- ')) {
    const head = lines[i].text.slice(2).trim();
    // `- { key: a }` and `- [a, b]` are whole values, not the first pair of a block-style map.
    const inlineCollection = head.startsWith('{') || head.startsWith('[');
    if (!inlineCollection && keyEnd(head) !== -1 && !head.startsWith('"')) {
      // An item that opens a map: its first pair sits on the dash line, the rest below it.
      const synthetic = [{ indent: indent + 2, text: head }];
      let j = i + 1;
      while (j < lines.length && lines[j].indent > indent) synthetic.push(lines[j++]);
      out.push(parseMap(synthetic, 0, indent + 2)[0]);
      i = j;
    } else {
      out.push(parseScalar(head));
      i += 1;
    }
  }
  return [out, i];
}

function parseMap(lines, start, indent) {
  const out = {};
  let i = start;
  while (i < lines.length && lines[i].indent === indent) {
    const sep = keyEnd(lines[i].text);
    if (sep === -1) break;
    const key = lines[i].text.slice(0, sep).trim();
    const rest = lines[i].text.slice(sep + 1).trim();
    if (rest !== '') {
      out[key] = parseScalar(rest);
      i += 1;
      continue;
    }
    const childIndent = i + 1 < lines.length ? lines[i + 1].indent : indent;
    const [value, next] = parseNode(lines, i + 1, childIndent);
    out[key] = value;
    i = next;
  }
  return [out, i];
}

function parse(text) {
  const lines = toLines(text);
  if (lines.length === 0) return {};
  return parseNode(lines, 0, lines[0].indent)[0];
}

module.exports = { parse };
