# EDR — CI y quality gates

- **Status:** Accepted
- **Type:** policy
- **Date:** 2026-07-28

## Contexto

La estrategia de testing ya está decidida (suite stdlib con SQLite real, cobertura informativa sin umbral); falta fijar qué bloquea un merge y qué garantiza el release. El repo ya tiene un workflow de release por tags con goreleaser, y la topología de despliegue anticipa que el binario del broker embebe el bundle de la UI — el orden de build importa.

Con la persistencia adentro, el estilo dejó de ser lo único que vale la pena verificar de forma determinista: aparece un borde que no debe filtrar sus detalles hacia el resto del daemon, y ese límite no se sostiene solo con revisión humana. Un gate que nadie puede olvidar vale más que una convención que hay que recordar.

## Decisión

- **Gates que bloquean merge** (GitHub Actions sobre PR), en orden de costo creciente: coherencia del manifest de dependencias, grafo de dependencias entre componentes, agregador de linters, suite de tests y build del binario. Sin umbral de cobertura como gate (decisión de testing: informativa).
- **Toda herramienta de verificación corre con su versión pineada**, la misma en CI y en las máquinas locales: una verificación cuyo resultado depende de qué versión bajó ese día no es determinista.
- **Release** (workflow de goreleaser sobre tags): el bundle de la UI se construye **antes** del build de Go, para embeberse en el binario; el release falla si el bundle no construye.
- **Pre-commit vía husky a nivel root del monorepo:** los `.go` staged pasan por el formateador estricto en el pre-commit (lint-staged rutea por config más cercana). El agregador de linters no corre en pre-commit (no es apto por-archivo, opera por paquete); el gate duro vive en el PR.

## Reglas verificables

- **[tool: go-arch-lint]** el merge se bloquea si un componente importa algo que su grafo declarado no permite, o si usa una dependencia externa que no tiene concedida.
- **[tool: golangci-lint]** el merge se bloquea si el linter reporta errores.
- **[tool: ci]** el merge se bloquea si `go mod tidy -diff` detecta que el manifest de dependencias quedó desincronizado del código.
- **[tool: ci]** el merge se bloquea si `go test ./...` falla o el binario no buildea.
- **[tool: ci]** el pipeline de release construye el bundle de la UI antes del build de Go y falla si el bundle no construye.
- **[tool: ci]** el workflow toma la versión del lenguaje del manifest del módulo, nunca de un valor escrito aparte que pueda desincronizarse.
- **[manual]** no se agrega umbral de cobertura como gate.
- **[tool: husky/lint-staged]** los `.go` staged pasan por `gofumpt` en el pre-commit del root.
- **[manual]** cada herramienta de verificación se invoca con versión explícita; ninguna resuelve "la última".

## Alternativas consideradas

- **Umbral de cobertura como gate:** ya descartado en la estrategia de testing; la cobertura es informativa.
- **Sin pre-commit del lado Go (solo gate de PR):** considerado primero; se optó por el husky root para atrapar formato antes del push con costo local mínimo.
- **`golangci-lint` en pre-commit:** descartado; opera por paquete, no por archivo staged — vive en el PR.
- **Escáner de vulnerabilidades de dependencias como gate:** evaluado y dejado afuera por ahora; la superficie de dependencias es chica y todas son de primera línea.
- **Herramientas sin pinear:** descartado; el gate pasaría a depender de qué release existía el día de la corrida.

## Consecuencias

- El PR es el único gate duro: lo que mergea, respeta el grafo de dependencias, lintea, testea y buildea.
- El límite del borde de persistencia deja de depender de que el reviewer se acuerde: filtrar los tipos de la base hacia otro componente rompe el build del PR.
- El release nunca publica un binario sin la UI embebida o con una UI rota.
- El formato se corrige solo al commitear; el resto se entera en CI — el pre-commit no reemplaza al gate del PR.
- Trade-off: el grafo de componentes hay que mantenerlo al día. Un componente nuevo no compila el gate hasta que alguien lo declare — que es exactamente lo que se buscaba, pero es un paso más en cada cambio estructural.
- Trade-off: compilar las herramientas desde fuente en vez de bajar binarios agrega tiempo a cada corrida de CI, a cambio de que la versión sea reproducible.

## Relacionados

- `relacionado-con` → [testing-strategy.md](testing-strategy.md) — la suite y la política de cobertura que estos gates ejecutan.
- `relacionado-con` → [deployment-topology.md](deployment-topology.md) — el embebido del bundle de la UI que ordena el build del release.
- `relacionado-con` → [../structure/architecture-style.md](../structure/architecture-style.md) — el grafo de componentes que uno de estos gates hace cumplir.
- `relacionado-con` → [../../../ui/.matecito-ai/edr/delivery/ci-quality-gates.md](../../../ui/.matecito-ai/edr/delivery/ci-quality-gates.md) — los gates espejo de la UI.
