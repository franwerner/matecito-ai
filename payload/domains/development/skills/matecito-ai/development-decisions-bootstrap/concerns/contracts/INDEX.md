# Dominio: `contracts`

Cómo se define la superficie pública que otros consumen: API REST/GraphQL, CLI, librería, eventos. Versionado, formato, compatibilidad.

## Criterio de pertenencia

Un concern nuevo va en `contracts` si define el contrato de *lo que este sistema expone*. Si trata sobre cómo este sistema *consume* servicios de terceros, va en `integration`.

## Concerns en este dominio

| Concern | Prof. | Type | Qué decide |
|---|---|---|---|
| [api-contract](api-contract.md) | light | decision | Cómo se versiona la API, cómo se pagina, el formato de respuesta estándar, e idempotencia de operaciones de escritura. |
| [cli-contract](cli-contract.md) | light | decision | Cómo se parsean los argumentos, qué exit codes se usan, qué va a stdout vs stderr, y el formato de output. |
| [event-contract](event-contract.md) | light | decision | Cómo se estructuran y versionan los eventos publicados, la convención de naming, e idempotencia del lado del consumidor. |
| [library-contract](library-contract.md) | light | decision | Qué expone la librería como superficie pública, cómo se versiona, y cuál es la política de backward compatibility y deprecación. |
