---
name: data-access
depth: deep
domain: data
type: decision
source: práctica clásica de patrones de acceso a datos · arc42 §8 (conceptos transversales)
---

# Fase: Acceso a datos

## Qué decide

Cómo el sistema lee y escribe datos persistentes: qué nivel de abstracción se usa sobre la DB, si hay patrón Repository, cómo se manejan las migraciones, y dónde se inician las transacciones.

## Preguntas

Una por turno. Para cada una: línea de "por qué importa", opciones con default, y siempre la opción "no sé, recomendame".

### 1. Nivel de abstracción sobre la DB

> Define cuánto SQL escribís a mano y qué tan ligado está el código al motor de DB elegido.

- **ORM completo** (SQLAlchemy, Prisma, Hibernate, Entity Framework, ActiveRecord) — menos SQL manual, más magia, migraciones integradas en algunos.
- ***Query builder — default para proyectos donde se quiere control sin ORM full.*** (Knex, Drizzle, jOOQ, Dapper). SQL explícito con ayuda de tipos.
- **Raw SQL** — máximo control, mínima abstracción. Justificado para performance crítica o queries muy complejas.
- **Mix: ORM para CRUD simple, raw para queries complejas** — pragmático, pero requiere criterio claro de cuándo usar cada uno.
- No sé, recomendame.

### 2. Patrón Repository

> Define si la lógica de acceso a datos queda encapsulada detrás de una interfaz o si los casos de uso/servicios consultan directamente.

- ***Sí, Repository por entidad de dominio — default cuando hay patrón arquitectónico (Clean, Layered).*** Facilita mocking y testabilidad.
- **No, acceso directo** — el servicio/caso de uso llama al ORM o DB directamente. Más simple, más acoplado.
- **Solo para las entidades agregadas** (DDD light) — Repository solo donde la complejidad lo justifica.
- No sé, recomendame.

### 3. Herramienta de migraciones

> Las migraciones sin herramienta formal producen esquemas inconsistentes entre entornos.

- **Integrada en el ORM** (Alembic para SQLAlchemy, Prisma Migrate, Flyway/Liquibase para Java, EF Core Migrations) — *default cuando ya se usa ORM que las incluye.*
- **Herramienta independiente** (Flyway, Liquibase, golang-migrate, dbmate) — útil cuando se usa query builder o raw SQL.
- **Sin migraciones formales** — solo para proyectos sin DB relacional o DB gestionada externamente.
- No sé, recomendame.

### 4. Transacciones — dónde se inician

> Si las transacciones se inician en el lugar incorrecto, la lógica de negocio queda atada a la infraestructura o las transacciones son demasiado largas.

- **En el caso de uso / application service** — *default para Clean Architecture.* El caso de uso define la unidad de trabajo; infra solo ejecuta.
- **En el service de dominio** — solo si el dominio tiene lógica transaccional propia (raro).
- **En el controller** — ok para apps CRUD simples sin casos de uso diferenciados.
- **Nunca (sin transacciones)** — solo si la DB no las soporta o cada operación es atómica por diseño (ej: Firestore).
- No sé, recomendame.

## Tech a registrar

Motor de DB (`postgresql.md`, `mongodb.md`, `mysql.md`, `sqlite.md`), ORM o query builder elegido (`sqlalchemy.md`, `prisma.md`, `drizzle.md`, `typeorm.md`), herramienta de migraciones si es separada del ORM (`alembic.md`, `flyway.md`, `golang-migrate.md`).

## Qué materializar

ADR `data-access` materializado según el template `../../templates/adr.md`. La **Decisión** captura: nivel de abstracción sobre la DB (ORM / query builder / raw SQL / mix), si hay patrón Repository y para qué entidades, la herramienta de migraciones, y dónde se inician las transacciones. Si hay mix (ej: ORM para CRUD + raw para reportes), documentar el criterio concreto para elegir cuándo usar cada uno. Registrar como tech el motor de DB, el ORM/query builder y la herramienta de migraciones si es separada.

**Reglas verificables** (cada una con su mecanismo al inicio):

- **[manual]** si se eligió Repository: el acceso a datos de las entidades cubiertas pasa por su Repository; los servicios/casos de uso no llaman al ORM ni a la DB directamente.
- **[manual]** las transacciones se inician en la capa decidida (ej: caso de uso / application service); las capas inferiores solo ejecutan dentro de la unidad de trabajo.
- **[manual]** en un esquema mix, el criterio documentado decide ORM vs raw; no hay raw SQL en rutas que el criterio asigna al ORM (y viceversa).
- **[manual]** todo cambio de esquema pasa por la herramienta de migraciones; no hay DDL aplicado a mano fuera de migraciones versionadas.

**Relacionados:** vincular con `data-modeling` (las migraciones materializan las convenciones de IDs, timestamps y borrado) y con un eventual ADR de capas/dependencias si el Repository es parte de la arquitectura.
