#!/usr/bin/env node
/**
 * Mapea los stores de records de un proyecto.
 *
 * Existe porque la relación entre un componente y su store no es uniforme y
 * nadie la declara: un componente puede tener store propio bajo su path, o
 * gobernarse por el store de la raíz. Inferir eso leyendo el árbol sale bien
 * casi siempre, y ese "casi" es un record que nadie consultó.
 *
 * Solo mapea. No busca, no lee records, no decide relevancia — esas decisiones
 * viven en el SKILL.md, donde se pueden leer y discutir.
 *
 * Uso:
 *   node map-records.js [--root <dir>] [--component <a,b>] [--json]
 *
 * `--component` acepta varios separados por coma, o repetido. Un cambio que
 * toca dos componentes se rige por los stores de ambos MÁS los de la raíz, así
 * que el filtro suma en vez de elegir.
 *
 * Sin --root usa el directorio actual y sube hasta encontrar `.matecito-ai/`.
 */

const fs = require("fs");
const path = require("path");

/** Subdirectorios de `.matecito-ai/` que son stores, y qué gobierna cada uno. */
const STORE_KINDS = {
  edr: { kind: "decisions", label: "decisiones de ingeniería" },
  ddr: { kind: "decisions", label: "decisiones de diseño" },
  "development-specs": { kind: "behavior", label: "comportamiento del sistema" },
};

/** Minúsculas y guiones: el nombre es clave de YAML y argumento de línea de comandos. */
function slug(s) {
  return String(s)
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "");
}

/**
 * Nombre corto del proyecto, que prefija cada colección.
 *
 * Las colecciones de todos los repos viven en UN archivo de configuración, en el
 * home del usuario, y ahí el nombre es la clave: si dos proyectos declaran
 * `decisions-root`, el segundo no se suma — reemplaza la ruta del primero, y la
 * colección queda apuntando al repo equivocado sin que nada falle ni avise.
 */
function projectName(root, config) {
  return slug(config?.repo?.name || path.basename(root));
}

function findRoot(start) {
  let dir = path.resolve(start);
  while (true) {
    if (fs.existsSync(path.join(dir, ".matecito-ai"))) return dir;
    const parent = path.dirname(dir);
    if (parent === dir) return null;
    dir = parent;
  }
}

function readConfig(root) {
  const file = path.join(root, ".matecito-ai", "config.json");
  if (!fs.existsSync(file)) return {};
  try {
    return JSON.parse(fs.readFileSync(file, "utf8"));
  } catch (err) {
    return { _error: `config.json ilegible: ${err.message}` };
  }
}

function countRecords(dir) {
  let n = 0;
  const walk = (d) => {
    for (const entry of fs.readdirSync(d, { withFileTypes: true })) {
      const full = path.join(d, entry.name);
      if (entry.isDirectory()) walk(full);
      else if (entry.name.endsWith(".md") && entry.name !== "INDEX.md") n++;
    }
  };
  try {
    walk(dir);
  } catch {
    /* un store ilegible se reporta con 0, no rompe el mapa */
  }
  return n;
}

/** Todo `.matecito-ai/` del árbol, saltando lo que no es del proyecto. */
function findStoreRoots(root) {
  const skip = new Set(["node_modules", ".git", "dist", "build", "vendor", ".venv"]);
  const found = [];
  const walk = (dir, depth) => {
    if (depth > 4) return;
    let entries;
    try {
      entries = fs.readdirSync(dir, { withFileTypes: true });
    } catch {
      return;
    }
    for (const entry of entries) {
      if (!entry.isDirectory() || skip.has(entry.name)) continue;
      const full = path.join(dir, entry.name);
      if (entry.name === ".matecito-ai") found.push(full);
      else walk(full, depth + 1);
    }
  };
  walk(root, 0);
  return found;
}

/**
 * El componente dueño de un store, o null si es de la raíz y ningún
 * componente se declaró dueño de ella.
 * Un store bajo `apps/api/.matecito-ai/` pertenece al componente cuyo path lo
 * contiene. El de la raíz pertenece al componente que se declaró `root: true`
 * si hay uno — entonces se filtra como cualquier peer — o, si ninguno lo
 * declaró, gobierna a todos exactamente como siempre.
 */
function ownerOf(storeRoot, root, components) {
  const rel = path.relative(root, path.dirname(storeRoot));
  if (rel === "") {
    const rootComponent = components.find((c) => c.root);
    return rootComponent ? rootComponent.name : null;
  }
  for (const comp of components) {
    for (const p of comp.paths || []) {
      if (rel === p || rel.startsWith(p + path.sep)) return comp.name;
    }
  }
  return undefined; // hay store pero ningún componente declarado lo cubre
}

function main() {
  const args = process.argv.slice(2);
  const get = (flag) => {
    const i = args.indexOf(flag);
    return i >= 0 ? args[i + 1] : undefined;
  };
  /** Todos los valores de un flag repetible, aplanando listas por coma. */
  const getAll = (flag) => {
    const out = [];
    args.forEach((a, i) => {
      if (a === flag && args[i + 1]) out.push(...args[i + 1].split(",").map((s) => s.trim()));
    });
    return out.filter(Boolean);
  };

  const root = findRoot(get("--root") || process.cwd());
  if (!root) {
    console.error("No encontré `.matecito-ai/` desde acá ni en ningún directorio padre.");
    process.exit(2);
  }

  const config = readConfig(root);
  const project = projectName(root, config);
  const components = config?.repo?.components || [];
  const filterComponents = getAll("--component");
  const warnings = [];
  const declared = new Set(components.map((c) => c.name));
  for (const name of filterComponents) {
    if (!declared.has(name)) {
      warnings.push(`Componente "${name}" no está declarado en config.json — se ignora en el filtro.`);
    }
  }
  if (config._error) warnings.push(config._error);
  if (components.length === 0) {
    warnings.push(
      "El proyecto no declara `repo.components`: no hay eje de componentes y todo store se trata como de raíz.",
    );
  }

  const stores = [];
  for (const storeRoot of findStoreRoots(root)) {
    const owner = ownerOf(storeRoot, root, components);
    if (owner === undefined) {
      warnings.push(
        `Store en ${path.relative(root, storeRoot)} sin componente declarado que lo cubra — se incluye igual, pero conviene declararlo en config.json.`,
      );
    }
    for (const [name, meta] of Object.entries(STORE_KINDS)) {
      const dir = path.join(storeRoot, name);
      if (!fs.existsSync(dir)) continue;
      const indexFile = path.join(dir, "INDEX.md");
      stores.push({
        kind: meta.kind,
        label: meta.label,
        // El nombre exacto con el que buscar en este store. Se devuelve en vez
        // de dejarlo deducir: escrito a mano es una cadena que nadie valida —
        // un nombre equivocado no falla, devuelve nada, y eso se lee igual que
        // "no hay records que apliquen".
        collection: `${project}-${meta.kind}-${owner ? slug(owner) : "root"}`,
        path: path.relative(root, dir),
        component: owner === undefined ? null : owner,
        // Después de reservar el nombre `root`, ningún otro campo distingue
        // las dos semánticas del store que hoy vive sin dueño en la raíz —
        // así que `scope` las auto-reporta (structure/root-store-scope-is-self-reported).
        scope:
          owner === "root"
            ? "componente root (raíz del repo)"
            : owner === null
              ? "todo el repo (sin componente raíz declarado)"
              : owner
                ? `componente ${owner}`
                : "todo el repo",
        index: fs.existsSync(indexFile) ? path.relative(root, indexFile) : null,
        records: countRecords(dir),
      });
    }
  }

  // Un componente sin store propio se gobierna por los de la raíz. Decirlo
  // explícito evita que alguien concluya que ese componente no tiene records.
  const covered = new Set(stores.map((s) => s.component).filter(Boolean));
  const inherit = components.map((c) => c.name).filter((n) => !covered.has(n));

  // Un cambio multi-componente se rige por la UNIÓN de sus stores más los de la
  // raíz, que gobiernan a todos. Filtrar por uno solo cuando toca dos deja
  // afuera records que aplican.
  const selected = filterComponents.length
    ? stores.filter((s) => s.component === null || filterComponents.includes(s.component))
    : stores;

  const result = {
    root,
    project,
    components: components.map((c) => c.name),
    componentsWithoutOwnStore: inherit,
    filteredBy: filterComponents.length ? filterComponents : null,
    stores: selected,
    warnings,
  };

  if (args.includes("--json")) {
    console.log(JSON.stringify(result, null, 2));
    return;
  }

  console.log(`Raíz del proyecto: ${root}`);
  console.log(`Proyecto: ${project}`);
  if (components.length) console.log(`Componentes: ${result.components.join(", ")}`);
  if (inherit.length) {
    console.log(`Sin store propio (se rigen por los de la raíz): ${inherit.join(", ")}`);
  }
  console.log("");
  for (const s of selected) {
    console.log(`${s.path}`);
    console.log(`  ${s.label} · ${s.scope} · ${s.records} records`);
    console.log(`  índice: ${s.index ?? "(sin INDEX.md)"}`);
    console.log(`  colección: ${s.collection}`);
  }
  if (warnings.length) {
    console.log("\nAvisos:");
    for (const w of warnings) console.log(`  - ${w}`);
  }
}

main();
