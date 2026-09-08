# Catálogo de tecnologías — root (transversales del monorepo)

Registro vivo de las tecnologías elegidas para el repo.

**Para Claude:** consultá esta tabla antes de sugerir una herramienta repo-wide. Si pisa con algo ya elegido, **no la agregues sin preguntar**.

## Por categoría

### Otros

| Tech                          | Versión    | Por qué (resumen)                                                                    |
| ----------------------------- | ---------- | ------------------------------------------------------------------------------------ |
| [husky](husky.md)             | sin pinear | Hooks de git a nivel root del monorepo; un solo pre-commit para todos los sub-apps.  |
| [lint-staged](lint-staged.md) | sin pinear | Rutea los checks de pre-commit por config más cercana: UI → eslint/prettier, Go → gofmt. |

## EDRs en este dominio

(Registros con el formato EDR canónico — Contexto/Decisión/Reglas verificables — para tecnologías cuya
inclusión acá no fue obvia y necesitó su propio razonamiento, a diferencia de la tabla de arriba.)

| EDR | Status | Consultá cuando... |
|---|---|---|
| [qmd.md](qmd.md) | Accepted | Vas a instalar o tocar algo relacionado con qmd, el MCP de búsqueda semántica de records. |

## Mantenimiento

- **Agregar tech:** crear `<nombre>.md`, sumar fila en la categoría.
- **Reemplazar tech:** borrar el archivo viejo y su fila del INDEX, crear el nuevo — git conserva el historial del archivo borrado.
- **Actualizar versión:** editar el archivo, anotar en Notas si hay breaking changes.
