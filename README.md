# Trabajo Práctico 6: Interfaces Dinámicas con HTMX

Sistema de Gestión de Proyectos (TP1): administración de proyectos, clientes y desarrolladores para un equipo, ahora con una interfaz **dinámica y reactiva utilizando HTMX** sobre la capa de Server-Side Rendering con [`templ`](https://templ.guide/) y la capa de persistencia del TP2 (`sqlc` + PostgreSQL). Reemplaza las recargas completas de página del TP5 por swaps parciales de HTML en el DOM sin escribir código JavaScript en el frontend.

## Ejecutar (TP6)

Requisitos: Docker + Docker Compose, Go 1.22+, [sqlc](https://sqlc.dev/) y `make`. Clonar, cambiar a la rama `tp6` y ejecutar:

```bash
git clone https://github.com/JereAguilar1/Trabajo_web.git
cd Trabajo_web
git checkout tp6
make run-app
```

Si no tenés `make`, `./run-app.sh` hace exactamente lo mismo.

`make run-app` (o `./run-app.sh`) hace todo en orden:

1. **Tareas previas**: crea `.env` desde `.env.example` si no existe, ejecuta `sqlc generate` y `templ generate`, compila (`go build`), borra contenedores y volúmenes anteriores, levanta PostgreSQL en Docker y espera a que esté listo (el esquema se carga al iniciar).
2. **Ejecuta la aplicación** en <http://localhost:8090>.
3. **Tareas posteriores**: al detenerla con `Ctrl+C`, borra los contenedores y volúmenes (los datos no se conservan).

**Configuración sin credenciales en el código:** la app lee `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD` y `DB_NAME` del entorno y no tiene valores por defecto: si falta alguna, termina con un error que las lista. `make run-app` las arma desde `.env` (`POSTGRES_*` y `DB_HOST_PORT`), y `docker-compose.yml` también las toma de ahí. `.env` está en `.gitignore`; el repo solo trae `.env.example`, sin contraseña. Al crear `.env`, `make run-app` genera una contraseña aleatoria.

Otros targets: `make test` (tests de persistencia contra una BD temporal), `make generate`, `make build`, `make down` (borra contenedores y volúmenes a mano) y `make docker` (BD y app en contenedores). El proyecto no incluye archivos JavaScript propios (solo la biblioteca HTMX vía CDN).

## Interfaz

Abrir <http://localhost:8090>. Es una sola página, maquetada con Pico CSS (solo unas pocas reglas propias para compactar los botones dentro de las tablas):

- **Navegación:** un `<nav>` arriba con enlaces a las secciones Proyectos, Clientes y Desarrolladores.
- **Secciones:** cada entidad es un `<article>` con su tabla (sin columna ID) y un botón "Nuevo …" que despliega el formulario de alta (`<details>`).
- **Datos legibles:** los montos se muestran como `$ 1500.00` y las fechas como `dd/mm/aaaa`.
- **Desarrolladores de un proyecto:** en cada fila, un desplegable "Desarrolladores (N)" muestra los asignados (con botón "Quitar") y un select para asignar otro.
- **Sin recargas:** crear (`POST`), borrar (`DELETE`) y asignar se hacen con peticiones asíncronas de HTMX, que actualizan solo la tabla afectada. Mientras se guarda se muestra "Guardando..." (`hx-indicator`); al terminar bien, el formulario se limpia y se cierra.
- **Errores:** si se intenta borrar algo con dependencias (por ejemplo un cliente con proyectos), el servidor responde `409 Conflict`, la interfaz conserva el registro y muestra el aviso.

## Rutas

| Ruta | Método | Acción |
|---|---|---|
| `GET /` | GET | Página completa inicial (formularios + listas) |
| `POST /clientes`, `/desarrolladores`, `/proyectos` | POST | Inserta la entidad y devuelve el fragmento HTML de la lista actualizada |
| `DELETE /clientes/{id}`, `/desarrolladores/{id}`, `/proyectos/{id}` | DELETE | Elimina la entidad y responde 200 OK vacío (HTMX remueve la fila) |
| `POST /proyectos/{id}/desarrolladores` | POST | Asigna un desarrollador al proyecto y actualiza la tabla de proyectos |
| `DELETE /proyectos/{id}/desarrolladores/{devId}` | DELETE | Desvincula al desarrollador del proyecto y actualiza la tabla de proyectos |

## Estructura del código

Código: `views/` (componentes `.templ`; los `_templ.go` son generados), `handlers/web/` (handlers HTML reactivos), `db/` (`sqlc`).
