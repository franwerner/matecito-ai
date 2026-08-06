# Componentes del repo — Eje `components`

Referencia canónica del **concepto repo-level** de `components`: qué es, cómo se declara, y el gate que lo activa. Cualquier consumidor —la documentación del capability-spec, un EDR, o una herramienta futura— cita esta página en vez de redefinir el concepto.

## Qué es un componente

Un **componente** es una **superficie que el consumidor del producto reconoce** — `api`, `cli`, `ui` — no todo paquete o carpeta interna. Es un eje **opcional** sobre el store de capability-specs: le agrega a cada spec qué superficie(s) del repo participan del comportamiento que describe.

## Declaración

El set de componentes se declara **una sola vez**, en el config **del proyecto** (nunca del global), en un bloque `repo` al tope, hermano de `domains`/`domainConfig`:

```json
"repo": {
  "components": [
    { "name": "cli", "paths": ["cmd", "internal"] },
    { "name": "api", "paths": ["apps/api"] },
    { "name": "ui",  "paths": ["apps/ui"] }
  ]
}
```

Array de objetos: `name` (la superficie) + `paths` (una o más carpetas del repo, relativas a la raíz, donde vive). Es el **único** lugar del ecosistema donde se escribe la carpeta de un componente — el capability-spec solo referencia el `name` en su línea de header; no repite las `paths`.

## Gate presence-based

Sin `repo.components` declarado, el eje **no existe**: ningún spec lleva la línea `Components:`, y ningún chequeo de validación se dispara por su ausencia. El eje es enteramente opt-in — un repo se comporta exactamente igual que antes de que este eje existiera hasta que alguien declara el set.

## Proyecciones

Una **proyección** es el uso que un consumidor hace del set declarado, con su propia granularidad, su propia vida y su propio momento de ratificación. Hoy hay dos: la línea `Components:` de cada capability-spec (por-capability, durable, ratificada spec por spec) y el valor de un cambio del flujo SDD (por-cambio, efímero, ratificado una sola vez en la confirmación de alcance). Toda proyección nueva declara esas tres cosas: su granularidad, su vida, y cuándo se ratifica.

Las proyecciones comparten exactamente dos cosas: el **set declarado** y el **gate presence-based** de arriba. No comparten el **valor**: una ratificación hecha en una granularidad no se escribe, copia ni propaga a otra — tampoco en batch automático desde una granularidad más gruesa hacia una más fina. Cuando un cambio toca capabilities cuyos componentes difieren entre sí, el valor por-cambio no decide los componentes de cada capability: cada una conserva el suyo, ratificado spec por spec.
