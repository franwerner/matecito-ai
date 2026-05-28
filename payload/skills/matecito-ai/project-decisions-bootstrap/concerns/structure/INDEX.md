# Dominio: `structure`

Cómo está organizado el código: el patrón macro, las capas, las reglas de dependencia entre ellas y la disposición física de archivos. Define la *forma* del sistema.

## Criterio de pertenencia

Un concern nuevo va en `structure` si trata sobre la organización estática del código y las reglas de quién-puede-depender-de-quién. Si trata sobre comportamiento en ejecución, va en `runtime`; si trata sobre la lógica del negocio en sí, va en `domain-logic`.

## Concerns en este dominio

| Concern | Prof. | Tipo | Qué decide |
|---|---|---|---|
| [architecture-style](architecture-style.md) | deep | decisión | El patrón macro de organización del código y el nivel de desacople entre componentes. Es la decisión que más condiciona las fases siguientes (capas, DI, test... |
| [folder-structure](folder-structure.md) | light | convención | Cómo se organiza el código dentro de cada capa (por feature vs por tipo técnico) y las convenciones de nombres por tipo de artefacto. Decisión de baja entrop... |
| [inter-layer-communication](inter-layer-communication.md) | deep | decisión | Cómo fluyen los datos entre capas: si se usan DTOs o entidades en los bordes, si la comunicación es síncrona o con eventos, dónde se declaran las interfaces,... |
| [layers-and-dependencies](layers-and-dependencies.md) | deep | decisión | Los nombres concretos de cada capa del sistema y las reglas de dependencia entre ellas, escritas de forma verificable. Es la decisión más importante para man... |
