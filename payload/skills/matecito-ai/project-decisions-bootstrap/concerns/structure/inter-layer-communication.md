---
name: inter-layer-communication
depth: deep
domain: structure
tipo: decisión
adr-output: inter-layer-communication
source: práctica clásica de arquitectura en capas · arc42 §8 (conceptos transversales)
---

# Fase: Comunicación entre capas

## Qué decide

Cómo fluyen los datos entre capas: si se usan DTOs o entidades en los bordes, si la comunicación es síncrona o con eventos, dónde se declaran las interfaces, y dónde vive la validación. Estas cuatro decisiones se intersectan y condicionan el testing.

## Preguntas

Una por turno. Para cada una: línea de "por qué importa", opciones con default, y siempre la opción "no sé, recomendame".

### 1. DTOs vs entidades en los bordes

> Define si la capa externa (controller, infra) trabaja con las entidades de dominio directamente o con objetos de transferencia propios. Impacta el nivel de acoplamiento y la cantidad de código de mapeo.

- **DTOs siempre en bordes** — *default para Clean Architecture estricta.* Más código de mapeo, pero el dominio queda aislado de cambios en la API/DB.
- **Entidades pueden cruzar** — más simple, ok para apps chicas o CRUD sin lógica de negocio real.
- **Mix: entrada con DTOs, salida puede ser entidad** — reduce mapeo de salida sin exponer la entidad al input externo.
- No sé, recomendame.

### 2. Sync vs async + eventos de dominio

> Determina si los efectos secundarios de un caso de uso se disparan de forma directa o desacoplada. Impacta testabilidad y complejidad operativa.

- ***Toda comunicación directa síncrona — default simple y suficiente para la mayoría de los casos.***
- **Eventos de dominio in-process** (mediator pattern) — desacople entre casos de uso sin agregar infraestructura externa.
- **Message bus externo** (Kafka / RabbitMQ / SQS) — desacople real entre servicios, complejidad operativa alta.
- No sé, recomendame.

### 3. Dirección de dependencias — dónde se declaran las interfaces

> Dónde viven las interfaces de repositorios y servicios externos define el nivel de pureza de la arquitectura. Impacta qué capa depende de qué.

- **En `domain/`** — Hexagonal estricto. El dominio define los contratos; infra los implementa.
- ***En `application/` — Clean clásica, default recomendado.*** Los casos de uso definen qué necesitan; dominio no sabe nada de IO.
- **No usamos interfaces** — depende de implementaciones concretas. Ok para proyectos sin requisito de testabilidad alta.
- No sé, recomendame.

### 4. Validación — dónde vive

> Duplicar validación es costoso; no duplicarla es un riesgo. El balance depende del nivel de confianza en los datos de entrada.

- **Solo en el borde** (DTO/schema con pydantic / zod / joi / bean validation) — simple, pero el dominio confía en que quien lo llame ya validó.
- **Solo en el dominio** (constructores estrictos, value objects) — más robusto, más código en el dominio.
- ***Ambos — defensa en profundidad, default robusto.*** El borde valida formato/tipo; el dominio valida reglas de negocio.
- No sé, recomendame.

## Notas de lógica (para el motor)

- **Si en architecture-style eligió "Sin patrón formal":** esta fase es Not Applicable. Crear el ADR con ese status y motivo.
- **Si eligió "Vertical Slice":** la pregunta 3 cambia — no hay interfaces de repositorio en `domain/` ni `application/`; en Vertical Slice las interfaces viven dentro del feature. Adaptar la opción y el ADR.
- **Si eligió "Message bus externo" en la pregunta 2:** es una decisión de infraestructura que implica tech adicional (Kafka, RabbitMQ, SQS, etc.). Registrala en `tech/` y anotá que la fase de data-access puede necesitar revisión.

## Tech a registrar

Si en la pregunta 4 se menciona una librería de validación concreta (pydantic, zod, joi, bean validation, FluentValidation, etc.), registrarla en `tech/`.

## Qué materializar

ADR `inter-layer-communication` con: política de DTOs vs entidades (con nombres concretos de las clases de mapeo si los hay), estilo de comunicación sync/async, dónde se declaran las interfaces de repositorios y servicios externos, y la política de validación (qué valida cada capa, con qué herramienta). Todo como reglas verificables, no como intenciones vagas.
