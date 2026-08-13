#!/usr/bin/env node
'use strict';

// Validates a phase return against the contract declared in `<phase>.yaml`, so the orchestrator's
// Return Contract Check stops being a human reading and comparing titles by eye.
//
// Usage:
//   validate-return.js --phase sdd-design [--status done] [--file return.md] [--contracts DIR]
//   validate-return.js --phase sdd-design --self-check [--contracts DIR]
//
// Exit 0 when clean, 1 on violations, 2 when it could not run at all. A validator that fails
// silently is worse than no validator: the orchestrator must surface exit 2, never swallow it.

const fs = require('fs');
const path = require('path');
const yaml = require('./lib/yaml');

const SENTINEL = /^(none|ninguna|ninguno|nada)\b/i;
const HEADING = /^(#{2,6})\s+(.*)$/;

// --- return parsing -------------------------------------------------------

// Returns the headings of the report block plus each one's body, keyed by literal title.
// A block title may carry placeholders (`## Discovery Form: {title}`), so it is matched as a pattern
// rather than compared literally.
function blockMatcher(blockTitle) {
  const escaped = blockTitle.replace(/[.*+?^${}()|[\]\\]/g, '\\$&').replace(/\\\{\w+\\\}/g, '.+');
  return new RegExp(`^${escaped}$`);
}

function parseReturn(text, blockTitle, declared) {
  const lines = text.split('\n');
  const matcher = blockMatcher(blockTitle);
  let start = lines.findIndex((l) => matcher.test(l.trim()));
  if (start === -1) return { blockFound: false, headings: [], bodies: {}, relevelled: [] };

  const blockLevel = blockTitle.match(HEADING)[1].length;
  const wanted = new Map();
  for (const s of declared || []) {
    for (const t of [s.title].concat(s.accepts || [])) wanted.set(t.replace(HEADING, '$2'), s.title);
  }

  const headings = [];
  const relevelled = [];
  const bodies = {};
  let current = null;

  for (let i = start + 1; i < lines.length; i++) {
    const m = lines[i].trim().match(HEADING);
    if (m && m[1].length <= blockLevel) {
      // Re-levelling a section reads as "the block ended here", which would blame every section
      // below it instead of the one heading that actually broke.
      const owner = wanted.get(m[2].trim());
      if (!owner) break;
      // A section the contract itself declares at the block's level is not re-levelled, it is correct:
      // `sdd-verify` puts `## UI Verdict` and `## Decision Gaps` there on purpose, because the gates
      // that read them match at that level.
      const found = lines[i].trim();
      if (found !== owner) relevelled.push({ found, expected: owner, line: i + 1 });
      current = owner;
      headings.push({ title: owner, line: i + 1 });
      bodies[owner] = [];
      continue;
    }
    if (m) {
      current = lines[i].trim();
      headings.push({ title: current, line: i + 1 });
      bodies[current] = [];
      continue;
    }
    if (current) bodies[current].push(lines[i]);
  }
  return { blockFound: true, headings, bodies, relevelled, full: text };
}

function bodyText(bodies, title) {
  return (bodies[title] || []).join('\n').trim();
}

function isSentinel(body) {
  const first = body.split('\n').map((l) => l.trim()).filter((l) => l !== '')[0];
  return first !== undefined && SENTINEL.test(first);
}

// --- checks ---------------------------------------------------------------

// A phase whose block changes by status (sdd-intake) declares `variants`; pick the one in play.
function resolveVariant(contract, status) {
  if (!contract.variants) return contract;
  const v = contract.variants.find((x) => (x.statuses || []).includes(status));
  return v ? { ...contract, ...v } : contract;
}

function validate(rawContract, text, status) {
  const problems = [];
  const add = (severity, code, message) => problems.push({ severity, code, message });

  const contract = resolveVariant(rawContract, status);
  if (rawContract.variants && status === null) {
    add('notice', 'STATUS-UNKNOWN', 'this phase returns a different block per status; pass --status to validate the right one');
  }
  const declared = contract.sections || [];
  const parsed = parseReturn(text, contract.block, declared);
  if (!parsed.blockFound) {
    add('error', 'BLOCK-MISSING', `the block \`${contract.block}\` is not in the return`);
    return problems;
  }

  for (const r of parsed.relevelled) {
    add('error', 'SECTION-RELEVELLED', `\`${r.found}\` (line ${r.line}) should be \`${r.expected}\` — the level is part of the title the guard matches`);
  }

  const titlesOf = (s) => [s.title].concat(s.accepts || []);
  const present = new Map();
  for (const h of parsed.headings) {
    const owner = declared.find((s) => titlesOf(s).includes(h.title));
    if (!owner) {
      add('error', 'SECTION-UNKNOWN', `\`${h.title}\` (line ${h.line}) is not declared by the contract — a renamed or re-levelled section is a gate that never fires`);
      continue;
    }
    if (present.has(owner.title)) {
      add('error', 'SECTION-DUPLICATE', `\`${owner.title}\` appears more than once`);
      continue;
    }
    present.set(owner.title, h.title);
  }

  for (const s of declared) {
    const seen = present.get(s.title);
    if (s.emitted === 'always' && !seen) {
      add('error', 'SECTION-MISSING', `\`${s.title}\` is unconditional and absent — it ships with a \`None…\` sentinel when there is nothing to report`);
    }
    if (s.emitted === 'on-status') {
      if (status === null) {
        add('notice', 'STATUS-UNKNOWN', `cannot check \`${s.title}\`: no \`**Status**:\` line found and no --status given`);
      } else if ((s.statuses || []).includes(status) && !seen) {
        add('error', 'SECTION-MISSING', `\`${s.title}\` is required on status \`${status}\` and is absent`);
      } else if (!(s.statuses || []).includes(status) && seen) {
        add('error', 'SECTION-UNEXPECTED', `\`${s.title}\` is only for status ${JSON.stringify(s.statuses)}, but the return is \`${status}\``);
      }
    }
  }

  for (const s of declared) {
    const seen = present.get(s.title);
    if (!seen || !s.items || !(s.items.token || s.items.tokens || s.items.rationale)) continue;
    problems.push(...checkItems(s, seen, bodyText(parsed.bodies, seen)));
  }

  for (const claim of contract.summary_claims || []) {
    const owner = declared.find((s) => s.title === claim.section);
    const seen = owner && present.get(owner.title);
    if (!seen) continue;
    if (!isSentinel(bodyText(parsed.bodies, seen))) continue;
    const hit = claimsContent(text, claim.pattern);
    if (hit) {
      add('error', 'SUMMARY-CONTRADICTS-BODY', `\`${seen}\` is closed with an empty sentinel, but the summary claims "${hit}"`);
    }
  }

  return problems;
}

// One walk per item, checking whichever of the independent markers the section declares: its
// token(s) (e.g. `blocking-test`, or `mandate` + `verify-checks`) and the `rationale` half of the
// summary/rationale split. A section may declare tokens, rationale, both, or neither — neither is
// filtered out by the caller. `tokensOf` desugars a single `token` to a one-element list, so a
// section with one token or several walks the same loop below.
function checkItems(section, seenTitle, body) {
  const out = [];
  if (isSentinel(body)) return out;

  const tokens = tokensOf(section.items).map((t) => ({ ...t, re: new RegExp(`^[\\s·*-]*${t.name}\\s*:\\s*(.+)$`, 'i') }));
  const wantsRationale = !!section.items.rationale;
  const lines = body.split('\n');
  const rationaleLine = /^[\s·*-]*rationale\s*:\s*(.+)$/i;

  let items = 0;
  for (let i = 0; i < lines.length; i++) {
    if (!/^\s*[-*]\s+\S/.test(lines[i])) continue;
    items += 1;
    const values = tokens.map(() => null);
    let rationaleValue = null;
    for (let j = i + 1; j < lines.length; j++) {
      if (/^\s*[-*]\s+\S/.test(lines[j])) break;
      const trimmed = lines[j].trim();
      for (const [k, t] of tokens.entries()) {
        if (values[k] !== null) continue;
        const m = trimmed.match(t.re);
        if (m) { values[k] = m[1].trim().replace(/[.`]+$/, ''); break; }
      }
      if (wantsRationale && rationaleValue === null) {
        const m = trimmed.match(rationaleLine);
        if (m) rationaleValue = m[1].trim();
      }
    }
    const item = lines[i].trim().slice(0, 60);
    for (const [k, t] of tokens.entries()) {
      const value = values[k];
      if (value === null) {
        out.push({ severity: 'error', code: 'TOKEN-MISSING', message: `"${item}" carries no \`${t.name}:\` — an omission gets the strict reading, not the benefit of the doubt` });
        continue;
      }
      // A token declared without `values` is free-form (e.g. `record: <domain>/<slug>`): any present,
      // non-null value passes — mirrors `render-return.js`'s `if (t.values && ...)` guard on the render
      // side. A token that DOES declare `values` keeps the exact legal/passing checks below, unchanged.
      if (!t.values) continue;
      const legal = t.values;
      const passing = t.passing || legal;
      if (!legal.includes(value)) {
        out.push({ severity: 'error', code: 'TOKEN-ILLEGAL', message: `"${item}" declares \`${t.name}: ${value}\`, which is not one of ${JSON.stringify(legal)}` });
      } else if (!passing.includes(value)) {
        out.push({ severity: 'error', code: 'TOKEN-WRONG-MAILBOX', message: `"${item}" declares \`${t.name}: ${value}\` — that value contradicts \`${seenTitle}\`` });
      }
    }
    if (wantsRationale && rationaleValue === null) {
      out.push({ severity: 'error', code: 'RATIONALE-MISSING', message: `"${item}" carries no \`· rationale:\` — a declaring section requires both parts` });
    }
  }

  if (items === 0) {
    const expected = tokens.length ? tokens.map((t) => `\`${t.name}:\``).join(' and ') : '`· rationale:`';
    out.push({ severity: 'error', code: 'SECTION-UNPARSEABLE', message: `\`${seenTitle}\` has content but no items to check — expected \`- \` bullets, each with its ${expected}` });
  }
  return out;
}

// Mirrors `render-return.js`'s desugar: a single `token` field becomes a one-element `tokens` list,
// so the walk below has one shape regardless of how many tokens an item declares.
function tokensOf(items) {
  if (items.tokens) return items.tokens;
  if (items.token) return [{ name: items.token, field: items.token_field, values: items.values, passing: items.passing }];
  return [];
}

// A claim only contradicts an empty sentinel when it asserts a non-zero count.
function claimsContent(text, pattern) {
  const re = new RegExp(`^.*${pattern}.*$`, 'gim');
  for (const line of text.match(re) || []) {
    const n = line.match(/(\d+)/);
    if (n && Number(n[1]) !== 0) return line.trim();
  }
  return null;
}

// --- self-check -----------------------------------------------------------

// The `.yaml` is only safe as a second file while this proves it still agrees with the `.md`.
function selfCheck(contract, md) {
  const problems = [];
  const headings = new Set(
    md.split('\n').map((l) => l.trim()).filter((l) => HEADING.test(l))
  );
  // Several templates document the phase's ARTIFACT alongside its return, in their own fenced block.
  // Only the fences that actually open with a declared return block count — otherwise an artifact-only
  // section (`### Tasks`, `### UI Scenario Counterparts` in sdd-apply) reads as an undeclared return
  // section, which is a false alarm about the one distinction these contracts exist to keep straight.
  const declaredBlocks = (contract.variants ? contract.variants.map((v) => v.block) : [contract.block])
    .filter(Boolean)
    .map(blockMatcher);
  const inFence = new Set();
  // Headings, inside a return fence, under which a literal `· rationale:` line appears — feeds the
  // bidirectional check below. Tracked alongside `inFence` rather than as a second pass over `md`.
  const rationaleUnderHeading = new Set();
  const RATIONALE_MARKER = /·\s*rationale\s*:/i;
  let fence = null;
  let currentHeading = null;
  // `sdd-verify.md` deliberately opens its return template with a `~~~` fence so it can nest literal
  // `` ``` `` blocks (build/test output) inside without closing early — the only template that does.
  // A fence line only opens or closes the CURRENTLY open fence when its delimiter matches; a different
  // delimiter while one is already open is nested content and is skipped, not treated as a toggle.
  for (const raw of md.split('\n')) {
    const l = raw.trim();
    const fenceMatch = l.match(/^(`{3,}|~{3,})/);
    if (fenceMatch) {
      const marker = fenceMatch[1][0];
      if (fence && marker !== fence.marker) continue;
      if (fence) {
        if (fence.isReturn) for (const h of fence.headings) inFence.add(h);
        fence = null;
        currentHeading = null;
      } else {
        fence = { isReturn: false, headings: [], marker };
      }
      continue;
    }
    if (!fence) continue;
    if (HEADING.test(l)) {
      if (declaredBlocks.some((m) => m.test(l))) fence.isReturn = true;
      fence.headings.push(l);
      currentHeading = l;
      continue;
    }
    if (fence.isReturn && currentHeading && RATIONALE_MARKER.test(l)) {
      rationaleUnderHeading.add(currentHeading);
    }
  }

  const allSections = contract.variants
    ? contract.variants.flatMap((v) => v.sections || [])
    : (contract.sections || []);
  const allBlocks = contract.variants ? contract.variants.map((v) => v.block) : [contract.block];

  for (const s of allSections) {
    if (!inFence.has(s.title) && !headings.has(s.title)) {
      problems.push({ severity: 'error', code: 'DRIFT', message: `\`${s.title}\` is declared in the contract but appears nowhere in the template` });
    }
    // An accepted variant is an alternative title, not the canonical one the template renders, so a
    // literal mention in prose is enough to prove the two files still agree on it.
    for (const variant of s.accepts || []) {
      if (!inFence.has(variant) && !headings.has(variant) && !md.includes(variant)) {
        problems.push({ severity: 'error', code: 'DRIFT', message: `the accepted variant \`${variant}\` is declared in the contract but appears nowhere in the template` });
      }
    }
  }
  for (const b of allBlocks) {
    const m = blockMatcher(b);
    const seen = [...inFence, ...headings].some((h) => m.test(h));
    if (!seen) problems.push({ severity: 'error', code: 'DRIFT', message: `the block \`${b}\` appears nowhere in the template` });
  }
  for (const title of inFence) {
    const level = title.match(HEADING)[1].length;
    if (level < 3) continue;
    const known = allSections.some((s) => [s.title].concat(s.accepts || []).includes(title));
    if (!known) {
      problems.push({ severity: 'error', code: 'DRIFT', message: `\`${title}\` appears in the template's example blocks but the contract does not declare it` });
    }
  }

  // Bidirectional rationale-marker rule, narrow to sections that can carry `items.rationale` — an
  // `items`-render section always (rationale or not, to also catch the "shows but doesn't declare"
  // case), plus a table or labeled-lists section only once it actually declares `items.rationale`
  // (most of them declare no `items` at all, and have nothing here to check). A generic "every
  // declared item marker must appear" rule would also fail `sdd-apply.md`, whose `verify-checks`
  // token is documented as inline prose rather than its own marker line — a defect this change
  // neither caused nor was asked to fix.
  for (const s of allSections) {
    if (s.render !== 'items' && !(s.items && s.items.rationale)) continue;
    const declares = !!(s.items && s.items.rationale);
    const titles = [s.title].concat(s.accepts || []);
    const shows = titles.some((t) => rationaleUnderHeading.has(t));
    if (declares && !shows) {
      problems.push({ severity: 'error', code: 'DRIFT', message: `\`${s.title}\` declares \`items.rationale\` but its template shows no \`· rationale:\` line under that heading` });
    }
    if (!declares && shows) {
      problems.push({ severity: 'error', code: 'DRIFT', message: `\`${s.title}\` shows a \`· rationale:\` line in its template but its contract does not declare \`items.rationale\`` });
    }
  }
  return problems;
}

// --- cli ------------------------------------------------------------------

function arg(name, fallback) {
  const i = process.argv.indexOf(`--${name}`);
  return i === -1 ? fallback : process.argv[i + 1];
}

// Deployed, the contracts sit in the user's ~/.claude; inside the payload repo they sit next to this
// script. Falling back to the sibling keeps `--contracts` genuinely optional in both places.
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
    process.stderr.write(`usage: validate-return.js --phase <name> [--status <s>] [--file <path>] [--contracts <dir>] [--self-check]\n  --contracts defaults to ${contractsDir(phase)}\n`);
    process.exit(2);
  }
  const dir = contractsDir(phase);
  const yamlPath = path.join(dir, phase, `${phase}.yaml`);
  const mdPath = path.join(dir, phase, `${phase}.md`);

  let contract;
  try {
    contract = yaml.parse(fs.readFileSync(yamlPath, 'utf8'));
  } catch (e) {
    process.stderr.write(`cannot read the contract for \`${phase}\` at ${yamlPath}: ${e.message}\n`);
    process.exit(2);
  }

  let problems;
  if (process.argv.includes('--self-check')) {
    let md;
    try {
      md = fs.readFileSync(mdPath, 'utf8');
    } catch (e) {
      process.stderr.write(`cannot read the template at ${mdPath}: ${e.message}\n`);
      process.exit(2);
    }
    problems = selfCheck(contract, md);
  } else {
    const file = arg('file');
    let text;
    try {
      text = file ? fs.readFileSync(file, 'utf8') : fs.readFileSync(0, 'utf8');
    } catch (e) {
      process.stderr.write(`cannot read the return: ${e.message}\n`);
      process.exit(2);
    }
    const explicit = arg('status');
    const found = text.match(/\*\*Status\*\*:\s*([a-z-]+)/i);
    problems = validate(contract, text, explicit || (found ? found[1].toLowerCase() : null));
  }

  const errors = problems.filter((p) => p.severity === 'error');
  for (const p of problems) process.stdout.write(`${p.severity.toUpperCase()} ${p.code}: ${p.message}\n`);
  if (errors.length === 0) process.stdout.write(`OK ${phase}: contract satisfied\n`);
  process.exit(errors.length === 0 ? 0 : 1);
}

main();
