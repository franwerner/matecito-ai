# EDR — Estilos y theming

- **Status:** Accepted
- **Type:** decision
- **Date:** 2026-07-23

## Contexto

El cockpit tiene tema dual (default light) y una tensión visual deliberada: un chrome cálido frente a un área de datos limpia. Hace falta fijar el mecanismo de estilos y theming sin invadir el territorio de los valores de token, que es del track de diseño.

## Decisión

**Tailwind CSS + shadcn/ui**: los componentes de Radix se copian al component base compartido y quedan editables. Theming dual por **CSS variables** siguiendo la convención de shadcn (`:root` + `.dark`) con un toggle; **default: claro (light)**. La tensión chrome cálido (radios grandes, tipografía display) vs. data limpia (radios chicos/rectos, mono) se resuelve con variantes y utilidades de Tailwind.

Los **valores finales de los tokens** (paleta exacta, escala tipográfica) son territorio del track de diseño (DDRs / Visual.md), **no** de este EDR de ingeniería: acá se fija el mecanismo.

## Reglas verificables

- **[manual]** los estilos se escriben con Tailwind + shadcn.
- **[manual]** el theming dual se implementa por CSS variables con toggle, default light.
- **[manual]** los valores finales de los tokens viven en el track de diseño, no en este EDR.

## Relacionados

- `relacionado-con` → [accessibility.md](accessibility.md) — el contraste AA se valida sobre los pares reales de cada tema.
- `relacionado-con` → [../delivery/configuration.md](../delivery/configuration.md) — la superficie de configuración de la app.
