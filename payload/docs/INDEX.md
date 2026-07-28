# Guías del payload

Cómo implementar sobre `payload/` — lo que matecito-ai despliega en la máquina
del usuario. Estas guías documentan el repo; **no se despliegan** ni se leen en
runtime.

| Guía | Consultala cuando… |
| --- | --- |
| [Armar un dominio](build-a-domain.md) | Creás un dominio nuevo, agregás o movés un componente dentro de uno (`agents/`, `skills/`, `references/`), nombrás una skill, tocás `manifest.json`, o trabajás sobre el tier `shared/`. |
| [Referenciar rutas desplegadas](reference-deployed-paths.md) | Escribís una instrucción que le dice a un agente o skill qué archivo leer, declarás `skills:` en el frontmatter de un agente, o citás un template/reference. |

## Regla de entrada

Si tu cambio agrega o mueve **archivos** dentro del payload → empezá por *Armar
un dominio*. Si agrega o cambia **texto que otro agente va a leer y seguir** →
empezá por *Referenciar rutas desplegadas*. Muchos cambios tocan las dos.

## Por qué existen

Ninguna de las dos reglas se deduce leyendo el código: el deploy aplana
directorios (la ruta desplegada no se parece a la del repo) y `payload/` no
existe en el proyecto del usuario. Peor, **probar dentro de este repo esconde
ambos errores**, porque acá las rutas del repo sí resuelven.
