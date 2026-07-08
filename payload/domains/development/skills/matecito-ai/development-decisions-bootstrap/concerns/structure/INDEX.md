# Dominio: `structure`

Cómo está organizado el código: el patrón macro, las capas, las reglas de dependencia entre ellas y la disposición física de archivos. Define la *forma* del sistema.

## Criterio de pertenencia

Un concern nuevo va en `structure` si trata sobre la organización estática del código, las reglas de quién-puede-depender-de-quién, o cómo se escribe el código a nivel de estilo. Si trata sobre comportamiento en ejecución, va en `runtime`; si trata sobre la lógica del negocio en sí, va en `domain-logic`.

**Frontera de naming (para no duplicar):** el *sufijo de rol* de un archivo y su ubicación → `folder-structure`; el *casing/estilo* de los nombres → `code-conventions`.

## Concerns en este dominio

| Concern | Prof. | Type | Qué decide |
|---|---|---|---|
| [architecture-style](architecture-style.md) | deep | decision | El patrón macro de organización del código y el nivel de desacople entre componentes. Es la decisión que más condiciona las fases siguientes (capas, DI, test... |
| [code-conventions](code-conventions.md) | deep | convention | Convenciones estilísticas de cómo se escribe el código a nivel micro: enum vs magic string, ausencia, inmutabilidad, estrictez de tipos, literales, casing de nombres, forma de función, igualdad, iteración, estilo de imports. Idiomáticas del lenguaje, mayormente lintables. |
| [folder-structure](folder-structure.md) | light | convention | Cómo se organiza el código dentro de cada capa (por feature vs por tipo técnico) y las convenciones de nombres por tipo de artefacto. Decisión de baja entrop... |
| [inter-layer-communication](inter-layer-communication.md) | deep | decision | Cómo fluyen los datos entre capas: si se usan DTOs o entidades en los bordes, si la comunicación es síncrona o con eventos, dónde se declaran las interfaces,... |
| [layers-and-dependencies](layers-and-dependencies.md) | deep | decision | Los nombres concretos de cada capa del sistema y las reglas de dependencia entre ellas, escritas de forma verificable. Es la decisión más importante para man... |
