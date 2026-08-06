# EDR — Hogar canónico del concepto `components` (repo-level)

- **Status:** Accepted
- **Date:** 2026-08-05

## Contexto
El concepto repo-level de `components` —qué es un componente, cómo se declara el set, y el gate que lo activa— vivía subordinado a la documentación del concepto de capability-spec, que es apenas uno de sus consumidores. Un consumidor nuevo de ese concepto tenía que abrir la documentación de capability-spec para encontrarlo, una dependencia invertida: el concepto no tenía hogar propio dentro del catálogo de referencias del dominio.

## Decisión
El concepto repo-level se muda a una reference propia del dominio development, independiente de cualquier consumidor puntual. Cada consumidor —incluida la documentación del capability-spec— cita esa reference en vez de redefinir o alojar el concepto.

## Reglas verificables
- **[manual]** El concepto repo-level de `components` (qué es un componente, cómo se declara el set, el gate) vive únicamente en su hogar canónico; ningún otro documento del payload lo redefine.

## Alternativas consideradas
Nombrar la reference nueva a secas, sin calificarla, se descartó: el sustantivo ya es de primera clase en otro dominio del ecosistema (diseño), y las references no se namespacean por dominio en el deploy — una reference así nombrada colisiona en el árbol compartido. Alojarla en el tier compartido entre dominios también se descartó: no cumple el criterio de admisión a ese tier, porque el concepto es específico de development, no transversal a todos los dominios activos.

## Consecuencias
El mecanismo de deploy detecta colisiones de reference de forma temprana y sin excepciones por tier, así que una colisión futura con esta reference sería ruidosa, no silenciosa. El nombre elegido espeja la clave de configuración que declara el set.
