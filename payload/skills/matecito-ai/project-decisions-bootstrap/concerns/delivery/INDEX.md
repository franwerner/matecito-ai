# Dominio: `delivery`

Cómo se construye, prueba, configura y despliega el sistema: CI/CD, testing, DI, configuración, deployment, feature flags, documentación.

## Criterio de pertenencia

Un concern nuevo va en `delivery` si trata sobre el proceso de llevar el código a producción de forma confiable. Si trata sobre la estrategia de versionado del *producto* (changelogs, breaking changes, soporte de versiones), va en `release`.

## Concerns en este dominio

| Concern | Prof. | Type | Qué decide |
|---|---|---|---|
| [arch-enforcement](arch-enforcement.md) | light | policy | Cómo convertir las reglas de dependencia definidas en `layers-and-dependencies` en configuración ejecutable de un linter que corra en CI y falle el build si ... |
| [ci-quality-gates](ci-quality-gates.md) | light | policy | Qué checks corren automáticamente en CI y cuáles bloquean el merge. Incluye pre-commit hooks para fallas rápidas en local. |
| [configuration](configuration.md) | light | convention | Cómo la aplicación lee su configuración (no secretos) según el entorno, y si esa configuración se valida y tipiea al arranque. |
| [dependency-injection](dependency-injection.md) | deep | decision | Cómo se conectan e instancian las dependencias del sistema: si se hace manualmente en un composition root, con un container de DI, o si lo maneja el framewor... |
| [deployment-topology](deployment-topology.md) | light | decision | Dónde y cómo corre la aplicación en producción: unidad de ejecución, cantidad de instancias, y si el proceso es stateless. |
| [documentation](documentation.md) | light | convention | Qué se documenta, dónde vive cada tipo de doc, y qué convenciones se siguen. El objetivo es que la documentación sea mantenible, no exhaustiva. |
| [feature-flags](feature-flags.md) | light | decision | Si el proyecto usa flags de features, con qué mecanismo, y cómo se nombran y eliminan para evitar deuda técnica acumulada. |
| [testing-strategy](testing-strategy.md) | deep | decision | La pirámide objetivo de tests, política de mocks vs reales, si TDD es obligatorio, y la cobertura mínima si se mide. Define qué tan costoso es cambiar el sis... |
