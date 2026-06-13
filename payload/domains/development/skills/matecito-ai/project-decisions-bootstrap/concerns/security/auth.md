---
name: auth
depth: deep
domain: security
type: decision
source: OWASP ASVS §2-3 (autenticación y gestión de sesiones) · arc42 §8 (conceptos transversales)
---

# Fase: Autenticación y autorización

## Qué decide

El mecanismo de autenticación, el modelo de permisos, dónde se valida, y la política de tokens/sesiones. Es una de las decisiones más difíciles de cambiar después — una vez que hay usuarios en producción, migrar el mecanismo de auth es costoso.

## Preguntas

Una por turno. Para cada una: línea de "por qué importa", opciones con default, y siempre la opción "no sé, recomendame".

### 1. Mecanismo de autenticación

> Define cómo el sistema verifica la identidad de quien hace una request. Determina la complejidad de setup y las dependencias externas.

- ***JWT (JSON Web Tokens) — default para APIs REST sin estado.*** Access token de corta duración, refresh token opcional.
- **Sessions con cookie** — *default para apps web con SSR.* Más simple, requiere almacén de sesiones (Redis o DB).
- **OAuth 2.0 / OIDC con proveedor externo** (Auth0, Keycloak, Cognito, Google) — delega la autenticación; menos código propio, dependencia externa.
- **API keys** — *default para integraciones B2B o APIs sin usuario humano.*
- **Mix** (ej: sessions para web, JWT para API pública, API keys para integración) — justificado en sistemas con múltiples tipos de clientes.
- No sé, recomendame.

### 2. Modelo de permisos

> Define qué tan granular es el control de acceso y qué tan complejo de mantener es.

- ***RBAC (Role-Based Access Control) — default razonable para la mayoría de sistemas.*** Roles predefinidos asignados a usuarios.
- **ABAC (Attribute-Based Access Control)** — permisos basados en atributos del usuario, el recurso, y el contexto. Más flexible, más complejo.
- **Simple booleano** (`is_admin` / `is_authenticated`) — ok para sistemas con dos o tres niveles de acceso sin necesidad de granularidad.
- **Sin modelo de permisos** — todos los usuarios autenticados tienen el mismo acceso.
- No sé, recomendame.

### 3. Dónde se valida

> Si la validación está dispersa, es fácil que un endpoint quede sin proteger. El patrón de validación centralizada reduce la superficie de olvido.

- ***Middleware / guard centralizado — default recomendado.*** Un solo lugar que protege todas las rutas; las excepciones se declaran explícitamente.
- **Decorator por endpoint / controller** — más explícito, más verboso, más fácil de olvidar.
- **Manual en cada handler** — sin framework de auth; máxima flexibilidad, máximo riesgo de olvido.
- No sé, recomendame.

### 4. Refresh tokens y política de expiración

> Una expiración muy larga es un riesgo de seguridad; una muy corta degrada la UX. Esta pregunta solo aplica si en (1) se eligió JWT o sessions con cookie.

- ***Access token: 15 min — Refresh token: 7 días con rotación — default OWASP para JWT.***
- Access token: 1 hora — sin refresh token (más simple, menor seguridad).
- Sessions: expiración por inactividad (ej: 30 min sin actividad).
- Definir según el nivel de sensibilidad de los datos del proyecto.
- No sé, recomendame.

## Notas de lógica (para el motor)

- **Si en Fase 0 el tipo de proyecto es `cli`, `libreria` o `script`:** esta fase probablemente es Not Applicable. Crear el ADR con ese status y motivo antes de preguntar.
- **Si eligió OAuth 2.0 / OIDC:** la pregunta 4 (expiración) la maneja el proveedor; aclarar eso y saltear o simplificar la pregunta.
- **Si eligió "API keys":** la pregunta 4 no aplica; en cambio, preguntar si las keys tienen fecha de expiración y cómo se rotan.
- **Si eligió "Sin modelo de permisos":** la pregunta 2 y parte de la 3 quedan reducidas; documentar el motivo en el ADR.

## Tech a registrar

Librería de auth si se usa una específica (`passport.md`, `authlib.md`, `next-auth.md`, `django-allauth.md`, `guardian.md`), proveedor de identidad externo si aplica (`auth0.md`, `keycloak.md`, `cognito.md`).

## Qué materializar

ADR `auth` materializado según `~/.claude/references/adr/templates/adr.md`. Debe contener:

- **Contexto**: tipo de clientes (web/API/integración), nivel de sensibilidad de los datos, y por qué este mecanismo es difícil de migrar una vez que hay usuarios en producción.
- **Decisión**: mecanismo de autenticación elegido y sus valores; modelo de permisos con descripción de los roles si aplica; dónde y cómo se valida (middleware/guard centralizado, decorator por endpoint, o manual); política de tokens/sesiones con las duraciones escritas como valores concretos, no como "corta duración".
- **Reglas verificables** (cada una con su mecanismo):
  - `[tool: test]` el access token expira a la duración decidida (ej. 15 min).
  - `[manual]` el refresh token rota en cada uso y caduca a la duración decidida (ej. 7 días).
  - `[manual]` toda ruta está protegida por el punto de validación elegido salvo las excepciones declaradas explícitamente.
  - Para API keys: `[manual]` las keys tienen fecha de expiración y un procedimiento de rotación documentado.
- **Alternativas consideradas**: los otros mecanismos evaluados (sessions, OAuth/OIDC, mix) y por qué no se eligieron.
- **Consecuencias**: dependencias externas introducidas (proveedor de identidad, almacén de sesiones) y trade-offs de seguridad/UX de las duraciones elegidas.
- **Relacionados** (si aplica): `relacionado-con` → `authorization.md` cuando el modelo de permisos se detalla allí.
