# Dominio: `context`

Punto de entrada: captura el propósito, alcance, stakeholders y stack del proyecto. Es la fase que alimenta los defaults de todas las demás.

## Criterio de pertenencia

Un concern nuevo va en `context` si responde *qué es y para quién es* el sistema, no *cómo se construye*. Si describe estructura, comportamiento u operación, va en otro dominio.

## Concerns en este dominio

| Concern | Prof. | Tipo | Qué decide |
|---|---|---|---|
| [context](context.md) | deep | decisión | El tipo de sistema, stack principal, tamaño de equipo y punto de partida. Es la entrada de todas las fases siguientes: sin este contexto, los defaults de cad... |
