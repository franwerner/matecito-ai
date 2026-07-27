# Matecito UI — Dirección visual (UI/UX)

> Brief de diseño autocontenido: define **qué es el producto** y **cómo debe verse y sentirse**.
> Es la dirección; el sistema de tokens final se fija en la fase de sistema de diseño.

## 1. Qué es el producto

**Matecito UI** es un **cockpit** para supervisar un **flujo de trabajo orquestado por IA**. Dada una corrida de trabajo, muestra —de forma tangible y remontable— el **árbol de agentes** que la ejecutó, qué produjo cada uno (código, decisiones), qué herramientas usó y el estado de cada pieza. No reemplaza la conversación: **la conversación sigue siendo el centro**; la UI es para lo *estructurado* (grafos, diffs, registros, listas).

Quién lo usa: alguien que supervisa a la IA — de alta presencia (quiere ver el flujo, el código, las decisiones) y que valora **entender el porqué** de cada cosa sin tener que preguntarlo.

El lazo central: *abrís una corrida → ves cómo se construyó → remontás cualquier cosa a su origen.*

## 2. Principios de UX

- **Todo es remontable.** Cada cosa tangible (un diff, una etiqueta, un nodo) linkea a su origen. Nunca un punto muerto para "¿de dónde salió esto?".
- **Densidad calma.** Mucha información, pero jerarquía sutil, aire medido, transiciones suaves. Denso sin ser ruidoso.
- **Keyboard-first.** Navegación y acciones por teclado; el mouse es opcional. Paleta de comandos.
- **La conversación es el centro.** La UI no compite con ella: muestra lo estructurado y reacciona a lo que pasa.
- **Ofrecer, no dar codazos.** "Qué necesita de mí" es mirable de un vistazo; nunca insiste.
- **Presencia y ausencia.** Se puede mirar en vivo *o* volver a una cola de "esto pasó, esto te espera".

## 3. Dirección visual (las tres decisiones)

1. **Tema dual — oscuro por defecto + claro cálido.** Oscuro por defecto (mejor para data densa y sesiones largas), con un modo claro cálido. Toggle. La identidad vive en ambos.
2. **Cálido en el chrome, limpio en la data.** La personalidad (mascota, calidez, redondeces, violeta) vive en navegación, estados vacíos y acentos. Las zonas densas (grafo, diffs, código) van **limpias, neutras y legibles**.
3. **Densidad calma.** Jerarquía por contraste sutil y espaciado —no por bordes pesados—; microtransiciones suaves. Calmo, no cargado.

## 4. Sistema visual

### 4.1 Identidad
Acento de IA **violeta `#8B7BFF`** (grafo, nodos) · calabash marrón `#B26A3C` · yerba verde `#6B8E4E` · plata `#C7CAD0` · crema `#FBF6EE` · tinta cálida `#3A2A1F`. Display **Fredoka**, body **Nunito**, mascota kawaii (un mate con carita). Tono: cálido, argentino sin cliché, *tech sin frialdad* — "IA = acento, nunca domina".

### 4.2 Superficies por tema
Oscuras **cálidas**, nunca negro puro (mantener la calidez):

| Rol | Oscuro (aprox.) | Claro (aprox.) |
|---|---|---|
| Fondo base | `#1A1410` | `#FBF6EE` |
| Superficie / tarjeta | `#221B15` | `#FFFDF7` |
| Superficie elevada | `#2B231C` | `#FFF6E0` |
| Borde / divisor | `#3A2F26` | `#E8DCC8` |
| Texto principal | `#F5ECE0` | `#3A2A1F` |
| Texto secundario | `#B8A896` | `#6F5A4A` |

El **violeta de IA se mantiene** en ambos temas (pop sobre oscuro, acento sobre claro).

### 4.3 Paleta semántica de datos
Colores de *significado*, separados de los de identidad, para las zonas de data:

| Significado | Color | Uso |
|---|---|---|
| IA / agente / nodo | violeta `#8B7BFF` | el hilo de la IA en todo el cockpit |
| Éxito / cumple | yerba `#6B8E4E` | pasa, completado |
| Crítico | rojo cálido `#C4553D` | violación, hallazgo crítico |
| Advertencia | ámbar `#D9963A` | warning (distinto de calabash) |
| Sugerencia / info | azul calmo `#5B8DB8` | sugerencia, metadata |
| Neutral / metadata | plata `#C7CAD0` | ids, timestamps, refs secundarias |

### 4.4 Tipografía
- **Fredoka** (display): títulos, wordmark, estados vacíos — *chrome*, no data.
- **Nunito** (body/UI): navegación, etiquetas, texto de tarjetas.
- **Monoespaciada** (una familia mono): **código, diffs, refs, ids, datos técnicos** — para legibilidad y alineación.

### 4.5 Espaciado, grid, radios
- Grid **8pt**, densidad calma (aire medido, ni apretado ni disperso).
- Radios **más contenidos en la data**: tarjetas de chrome redondeadas (`~14px`); zonas densas (tablas, diffs, grafo) con radios chicos o rectas — la redondez es personalidad de *chrome*, no de data.

## 5. Layout — canvas de workflow

Unidad = **una corrida**. La pieza central es un **canvas** (infinito, pan/zoom) donde el flujo se ve como un **workflow de ramificaciones** —el estándar actual para visualizar flujos— pero **calmo**: nodos limpios, aristas sutiles, buen aire. *Canvas sí, maraña no.*

```
┌──────────────────────────────────────────────┬───────────────┐
│  CANVAS (home) — workflow de ramificaciones   │   INSPECTOR   │
│                                               │   (drawer, al │
│      ┌ orquestador ┐                          │    click de   │
│      │             │  aristas = qué info pasa  │    un nodo)   │
│   [ nodo ]─[ nodo ]─┬─[ nodo ]─[ nodo ]        │               │
│                     └─[ nodo ]                 │   • feed      │
│   nodos = fases/agentes · violeta = IA         │   • diffs     │
│                                               │   • decisiones│
├──────────────────────────────────────────────┴───────────────┤
│  TIMELINE / scrubber — reproduce los nodos apareciendo         │
└────────────────────────────────────────────────────────────────┘
```

- **Canvas = el mapa (home).** Orquestador → agentes ramificando; nodos = fases/agentes, aristas = qué información pasa. Pan/zoom.
- **Inspector (drawer lateral, al click de un nodo).** Despliega **todo lo que ese nodo generó y usó**, denso y remontable:
  - **diffs de código** que escribió (con su justificación y etiquetas de severidad),
  - **decisiones** que aplicó (con su estado),
  - **specs** que usó,
  - **herramientas** que ejecutó,
  - su **feed de actividad** (en vivo) y su **resultado**.

  Desde un diff se remonta a la decisión/spec que lo gobierna. El canvas queda de mapa; la lectura pesada pasa al inspector.
- **Lentes.** "Código" y "decisión" se pueden mirar como filtros/resaltados sobre el mismo canvas (por archivo tocado, por decisión), o vistas que abren desde un nodo.
- **Timeline** reproduce los nodos apareciendo y ramificando.
- **Chrome cálido / canvas limpio.** Nav, header y timeline con personalidad; canvas e inspector limpios, mono donde hay técnica. **Nodos de IA en violeta.**

## 6. Componentes clave

- **Nodo de agente** — tipo, fase, mandato, estado, herramientas usadas; feed de actividad en vivo al expandir. Chrome cálido; contenido técnico en mono.
- **Aristas del canvas** — etiquetadas con qué info pasa entre nodos; sutiles, no ruidosas.
- **Vista de diff** — código en mono, adiciones/quitas legibles, con su justificación y etiquetas de revisión ancladas.
- **Tarjeta de decisión** — título, estado (borrador / ratificada), relaciones (qué código cumple).
- **Etiquetas de severidad** — crítico / advertencia / sugerencia con la paleta semántica (§4.3) + el problema; siempre con texto/ícono, no solo color.
- **Scrubber de timeline** — reproduce el orden de la corrida.
- **Tarjeta de artefacto** — cada artefacto renderizado como tarjeta con campos y refs, no prosa suelta.

## 7. Estados y personalidad

- **Estados vacíos** — acá vive la mascota kawaii: "todavía no hay una corrida activa", "sin decisiones tocadas". Calidez donde no hay data que estorbar.
- **En vivo** — el feed de actividad pulsa sutil (no parpadeos agresivos); casi en tiempo real.
- **La mascota guía, no decora** — aparece en onboarding, vacíos y momentos de éxito; nunca encima de la data densa.

## 8. Accesibilidad

- **WCAG AA** en ambos temas, sobre los pares reales de texto/fondo (el violeta y el ámbar sobre oscuro se validan, no se asumen).
- Tamaños mínimos legibles para data densa; foco de teclado visible (keyboard-first lo exige).
- El color nunca es el único portador de significado (las severidades llevan texto/ícono, no solo color).

---

> Esta dirección alimenta el sistema de diseño (tokens finales, componentes) y la producción de las pantallas.
