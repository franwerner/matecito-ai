# go-arch-lint

- **Category:** Other
- **Version:** v1.16.0
- **Status:** Accepted
- **Decided in phase:** arch-enforcement
- **Date:** 2026-07-28

## Por qué la elegimos

Vuelve ejecutable el grafo de dependencias entre componentes que hasta ahora solo vivía escrito en los EDRs de estructura: se declaran los componentes y qué puede depender de qué, y todo lo no declarado queda prohibido por default. También gobierna qué librería externa puede usar cada componente, que es lo que mantiene los detalles de persistencia encerrados en su borde.

## Alternativas descartadas

- **depguard solo (viene en golangci-lint):** razona por archivo e import, sin noción de componente ni de grafo; un paquete nuevo no declarado no rompe nada hasta que alguien escriba su regla. Se mantiene igual, pero acotado a lo que go-arch-lint no ve.
- **Reglas solo `[manual]` en los EDRs:** es lo que había; depende de que el reviewer se acuerde.

## Notas

No gobierna la librería estándar — solo componentes internos y dependencias externas declaradas. Las restricciones sobre paquetes de la stdlib se cubren con depguard.

El binario no se instala: se invoca con la versión pineada desde el propio módulo, así CI y las máquinas locales corren exactamente la misma.
