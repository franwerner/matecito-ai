# Catálogo de tecnologías — root (transversales del monorepo)

Registro vivo de las tecnologías elegidas **a nivel repo** (afectan a todos los sub-apps). Las tecnologías propias de cada sub-app viven en su catálogo: `apps/api/.matecito-ai/edr/tech/` y `apps/ui/.matecito-ai/edr/tech/`.

**Para Claude:** consultá esta tabla antes de sugerir una herramienta repo-wide. Si pisa con algo ya elegido, **no la agregues sin preguntar**.

## Por categoría

### Otros

| Tech                          | Versión    | Por qué (resumen)                                                                    |
| ----------------------------- | ---------- | ------------------------------------------------------------------------------------ |
| [husky](husky.md)             | sin pinear | Hooks de git a nivel root del monorepo; un solo pre-commit para todos los sub-apps.  |
| [lint-staged](lint-staged.md) | sin pinear | Rutea los checks de pre-commit por config más cercana: UI → eslint/prettier, Go → gofmt. |

## Mantenimiento

- **Agregar tech:** crear `<nombre>.md`, sumar fila en la categoría.
- **Reemplazar tech:** borrar el archivo viejo y su fila del INDEX, crear el nuevo — git conserva el historial del archivo borrado.
- **Actualizar versión:** editar el archivo, anotar en Notas si hay breaking changes.
