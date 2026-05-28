# Dominio: `security`

Cómo se protege el sistema de uso indebido: autenticación, autorización, validación de input, secretos, CORS, rate limiting, escaneo de dependencias.

## Criterio de pertenencia

Un concern nuevo va en `security` si trata sobre *proteger* el sistema o sus datos de un actor hostil. Si trata sobre el derecho a tener/usar datos personales, va en `privacy`; si trata sobre cumplir una regulación, va en `compliance`.

## Concerns en este dominio

| Concern | Prof. | Tipo | Qué decide |
|---|---|---|---|
| [auth](auth.md) | deep | decisión | El mecanismo de autenticación, el modelo de permisos, dónde se valida, y la política de tokens/sesiones. Es una de las decisiones más difíciles de cambiar de... |
| [authorization](authorization.md) | light | decisión | Qué modelo de permisos rige el acceso a recursos o acciones, y dónde se evalúa esa lógica en el sistema. |
| [cors](cors.md) | light | política | Qué orígenes pueden hacer requests cross-origin a la API y si se permiten credenciales (cookies, Authorization header). |
| [dependency-scanning](dependency-scanning.md) | light | política | Qué herramienta detecta vulnerabilidades conocidas en dependencias y si ese escaneo corre de forma automática en el pipeline de CI. |
| [input-validation](input-validation.md) | light | política | Dónde y cómo se valida y sanitiza el input externo antes de que entre al sistema, y qué se hace cuando falla la validación. |
| [rate-limiting](rate-limiting.md) | light | decisión | Si se limita la cantidad de requests, con qué granularidad (IP / usuario / API key) y dónde se aplica el límite. |
| [secrets-management](secrets-management.md) | light | política | Dónde se almacenan los secretos del sistema (credenciales, tokens, claves), cómo se rotan, y qué está explícitamente prohibido commitear o loggear. |
