---
name: architecture-style
depth: deep
domain: structure
type: decision
source: práctica clásica de patrones arquitectónicos · arc42 §4 (vista de solución)
---

# Fase: Estilo arquitectónico y acoplamiento

## Qué decide

El patrón macro de organización del código y el nivel de desacople entre componentes. Es la decisión que más condiciona las fases siguientes (capas, DI, testing).

## Preguntas

Una por turno. Para cada una: línea de "por qué importa", opciones con default, y siempre la opción "no sé, recomendame".

### 1. Patrón arquitectónico

> Define cómo se organiza el sistema a gran escala y qué tan independiente es la lógica de negocio de la infraestructura.

- **Clean / Hexagonal / Ports & Adapters** — cuando la lógica de negocio es rica y querés que sea testeable e independiente de framework y DB.
- **Layered N-tier clásica** (Controller → Service → Repository) — pragmática, fácil de entender, *default razonable para CRUD-heavy.*
- **Vertical Slice / Feature folders** — cuando los features son independientes y querés cohesión por feature en lugar de por capa técnica.
- **MVC simple** — apps web tradicionales, prototipos.
- **Event-driven / CQRS** — sistemas con muchos side effects, alta escala de lectura, o auditoría fuerte.
- **Sin patrón formal** (lo más simple posible) — scripts, prototipos, librerías chicas.
- No sé, recomendame.

### 2. Nivel de acoplamiento buscado

> Define qué tan fácil es reemplazar una pieza (DB, HTTP client, framework) sin tocar la lógica de negocio.

- **Alto desacople** — interfaces y DI por todos lados. Mejor para test-driven, peor para velocidad inicial.
- ***Pragmático — default recomendado.*** Interfaces solo en bordes I/O (DB, HTTP externo, filesystem). Resto concreto. Buen balance.
- **Acoplamiento directo** — sin interfaces. Código más simple, menos testeable, ok para prototipos.
- No sé, recomendame.

## Notas de lógica (para el motor)

- **Default según tipo de proyecto:** si en Fase 0 eligió `script` o `libreria chica`, proponé "Sin patrón formal" en la pregunta 1 y "Acoplamiento directo" en la 2. Mostrá el default y pedí confirmación.
- **Si eligió Clean / Hexagonal:** las fases layers-and-dependencies e inter-layer-communication son requeridas; marcalas como tales antes de avanzar.
- **Si eligió Vertical Slice:** la fase layers-and-dependencies cambia de naturaleza (no hay capas horizontales, hay módulos por feature); avisale al usuario y ajustá el cuestionario de esa fase.
- **Si eligió Event-driven / CQRS:** la fase inter-layer-communication incluye preguntas adicionales sobre message bus; avisale antes de entrar.

## Qué materializar

ADR `architecture-style` materializado según `../../templates/adr.md`. Debe contener:

- **Contexto** y **Decisión**: patrón arquitectónico elegido, nivel de acoplamiento buscado, y justificación concreta (por qué este patrón para este tipo de proyecto y equipo).
- **Reglas verificables**: si el patrón implica restricciones específicas, enumeralas como aserciones chequeables con su mecanismo al inicio. Ej: `[tool: dependency-cruiser]` ningún módulo de `domain/**` importa framework externo; `[manual]` los bordes I/O exponen interfaces, el resto es concreto (si no hay check automático). Estas reglas son la base de layers-and-dependencies.
- **Alcance**: como decisión estructural, incluí los globs **a nivel convención** que el patrón gobierna (ej: `src/domain/**`, `src/application/**`, `src/infrastructure/**` para Clean; `features/<feature>/**` para Vertical Slice). Patrones estables, no archivos concretos.
- **Relacionados** (opcional): vinculá con `layers-and-dependencies` e `inter-layer-communication`, que refinan esta decisión.
