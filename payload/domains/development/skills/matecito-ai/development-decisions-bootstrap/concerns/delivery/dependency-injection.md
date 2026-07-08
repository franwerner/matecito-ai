---
name: dependency-injection
depth: deep
domain: delivery
type: decision
source: práctica clásica de IoC y composition root · arc42 §8 (conceptos transversales)
---

# Fase: Inyección de dependencias

## Qué decide

Cómo se conectan e instancian las dependencias del sistema: si se hace manualmente en un composition root, con un container de DI, o si lo maneja el framework. Impacta directamente la testabilidad y la complejidad de setup.

## Preguntas

Una por turno. Para cada una: línea de "por qué importa", opciones con default, y siempre la opción "no sé, recomendame".

### 1. Mecanismo de DI

> Define dónde y cómo se cablean las dependencias. Es la decisión central de esta fase; las siguientes son detalles de la elegida.

- **Manual (composition root / main)** — factories explícitas en un único punto de entrada. Máxima visibilidad, sin magia, más código de setup.
- ***Container de DI — default cuando hay muchas dependencias o el proyecto crece.*** Más automático, menos código manual, requiere aprender la librería.
- **Framework-provided** (NestJS modules, Spring IoC, ASP.NET DI, FastAPI Depends) — *default cuando el framework ya lo trae integrado y el equipo lo conoce.*
- No sé, recomendame.

### 2. Librería de container (solo si eligió "Container de DI")

> La elección de librería determina convenciones de registro, scopes y testabilidad.

- **Python:** `dependency-injector` (basado en providers explícitos), `injector` (basado en anotaciones).
- **JS/TS:** `awilix` (registro explícito, funcional), `tsyringe` (decoradores), `inversify` (decoradores, más maduro).
- **Java:** Spring IoC (de facto), Guice (alternativa más liviana).
- **C#:** Microsoft.Extensions.DependencyInjection (built-in), Autofac.
- **Go:** wire (generación de código), fx (runtime).
- No sé, recomendame.

### 3. Scopes de vida de los objetos

> Un scope incorrecto produce bugs sutiles (estado compartido inesperado) o performance innecesaria (objetos que se recrean de más).

- ***Singleton para servicios stateless, Scoped (por request) para unit of work y repos — default recomendado para apps web.***
- Singleton para todo — solo si el código es completamente stateless.
- Transient para todo — recrear siempre; seguro pero puede ser costoso.
- Definido caso por caso (sin regla general).
- No sé, recomendame.

## Notas de lógica (para el motor)

- **Si en architecture-style eligió "Sin patrón formal" o "acoplamiento directo":** DI probablemente es Not Applicable o se reduce a instanciación manual en main. Crear el EDR con ese status y motivo.
- **Si el framework ya provee DI** (NestJS, Spring, ASP.NET): la pregunta 2 (librería de container) no aplica — anotar que se usa el mecanismo del framework y por qué (integración nativa, menos dependencias).
- **Si eligió "Manual":** la pregunta 2 y 3 cambian de naturaleza — ya no hay container ni scopes formales. El EDR documenta el composition root y qué factories existen.

## Tech a registrar

Si se eligió un container de DI externo: `dependency-injector.md`, `awilix.md`, `tsyringe.md`, `inversify.md`, `wire.md`, `fx.md`, u otro. No registrar el DI del framework si ya está registrado como parte del framework principal.

## Qué materializar

EDR `dependency-injection` materializado según `~/.claude/references/edr/templates/edr.md`. Debe contener:

- **Contexto** y **Decisión**: mecanismo elegido (manual / container / framework), librería específica si aplica (`dependency-injector`, `awilix`, `tsyringe`, Spring IoC, etc.), scopes por tipo de componente (servicio, repository, unit of work), y — si es composition root manual — dónde vive y quién es responsable de armarlo.
- **Reglas verificables**: las reglas de scope como aserciones con su mecanismo al inicio, no como intención vaga. Ej: `[manual]` todo `Service` se registra como singleton; `[manual]` todo `Repository` y unit of work es scoped al request; `[tool: <linter/test>]` si la convención de registro es chequeable por la librería de DI o un test de wiring. Conservá los valores concretos de scope por tipo de componente.
- **Relacionados** (opcional): vinculá con `architecture-style` y `inter-layer-communication` (dónde se declaran las interfaces que el container cablea).
