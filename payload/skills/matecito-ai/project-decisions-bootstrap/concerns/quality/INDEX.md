# Dominio: `quality`

Atributos de calidad cuantificables (NFRs): performance, escalabilidad, internacionalización. Definen umbrales medibles, no funcionalidad.

## Criterio de pertenencia

Un concern nuevo va en `quality` si define un *objetivo medible de calidad* (latencia, throughput, capacidad, cobertura de idiomas) transversal al sistema. Si es una decisión de diseño puntual, va en su dominio funcional.

## Concerns en este dominio

| Concern | Prof. | Tipo | Qué decide |
|---|---|---|---|
| [i18n-l10n](i18n-l10n.md) | light | decisión | Si la aplicación soporta múltiples idiomas o configuraciones regionales, qué librería gestiona las traducciones, y dónde viven los strings. |
| [nfr-performance](nfr-performance.md) | light | decisión | Los objetivos cuantitativos de tiempo de respuesta y throughput del sistema. Sin números acordados, "performance aceptable" es subjetivo y no testeable. |
| [scalability](scalability.md) | light | decisión | El modelo de escalado esperado (vertical u horizontal) y si la arquitectura lo soporta desde el inicio. ISO/IEC 25010 lo define bajo "capacity": grado en que... |
