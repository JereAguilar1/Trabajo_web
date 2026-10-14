# Trabajo Práctico 3: Exponiendo la Lógica como API

API REST para el Sistema de Gestión de Proyectos del TP1, construida sobre la capa de persistencia del TP2 (`sqlc` + PostgreSQL).

## Requisitos

- [Docker](https://www.docker.com/) y Docker Compose.
- [Go](https://go.dev/) 1.22+, [sqlc](https://sqlc.dev/) y `curl` — los usa `./test.sh` en tu máquina (generar código, correr los tests de persistencia y esperar a que la API responda).

`db/sqlc/` (código generado) y `.env` no se versionan: `sqlc generate` y la copia de `.env.example` a `.env` los hace `./test.sh` (y el `Dockerfile` para el código generado).

## Cómo ejecutar todo (tests + API)

```bash
git clone https://github.com/JereAguilar1/Trabajo_web.git
cd Trabajo_web
git checkout tp3
./test.sh
```

`test.sh` hace todo en orden:

1. **Tareas previas**: crea `.env` desde `.env.example` si no existe, corre `sqlc generate` y `go build`, borra contenedores y volúmenes, levanta Postgres y espera a que esté listo, levanta la aplicación (`docker compose up -d --build app`) y espera a que responda en `http://localhost:8090`.
2. **Pruebas**: `go test ./...` (persistencia) y `docker compose run --rm hurl` (cada endpoint de la API, con `requests.hurl`).
3. **Tareas posteriores**: borra contenedores y volúmenes, y sale con código distinto de 0 si algo falló.

## Configuración (`.env`)

Las credenciales y el puerto de Postgres se leen de `.env`. `./test.sh` lo crea a partir de `.env.example` si no existe; para correr `docker compose` a mano, copialo primero (`cp .env.example .env`). Si el puerto 5432 de tu máquina está ocupado, cambiá `DB_HOST_PORT` (por defecto `5433`).

## Levantar solo el servidor

```bash
cp .env.example .env
docker compose up --build
```

Levanta **`db`** (PostgreSQL 16 con el esquema de `db/schema/schema.sql` aplicado al iniciar) y **`app`** (el servidor HTTP, que espera a que `db` esté saludable y escucha en `http://localhost:8090`). Docker compila el binario dentro de la imagen (`Dockerfile`, multi-stage) y ejecuta `sqlc generate` en el build.

## Pruebas de la API

`requests.hurl` ejercita el ciclo de vida completo de las entidades del dominio: crear, listar, obtener por ID, actualizar y eliminar, incluyendo casos de validación (`400`), de recurso inexistente (`404`) y el alta/baja/listado de la relación N:M entre desarrolladores y proyectos. `./test.sh` ya lo corre; con el servidor levantado también se puede correr aparte:

```bash
docker compose run --rm hurl
```

El servicio `hurl` usa `network_mode: "host"` para llegar a `localhost:8090` y tiene `profiles: ["test"]`, por lo que no se levanta con un `docker compose up` normal.

## Tests de persistencia (heredado del TP2)

`go test ./...` ejercita `db/sqlc` directamente, sin pasar por la API. Necesitan Postgres corriendo, por eso se corren a través de `./test.sh`.

## Documentación

### Arquitectura de la API

Siguiendo el modelo de tres capas (Presentación → Lógica de negocio → Datos):

- **Handlers HTTP** (`handlers/`): un paquete por recurso — `handlers/clientes`, `handlers/desarrolladores`, `handlers/proyectos` y `handlers/desarrolladores_en_proyectos` (cada uno en su `handlers.go`). Parsean la petición, validan los datos de entrada y traducen el resultado a códigos de estado HTTP. Exponen `Create`, `GetAll`, `GetByID`, `Update` y `Delete` (el de la relación N:M expone `Assign`, `List` y `Remove`).
- **Repositorio** (`db/sqlc/`, del TP2): ejecuta las consultas SQL. Los handlers lo reciben como parámetro (inyección de dependencias, `main.go` lo crea una sola vez) en vez de crear su propia conexión.
- **Enrutamiento** (`main.go`): usa el `http.ServeMux` de Go 1.22+, que soporta método + wildcards en el patrón (`"GET /clientes/{id}"`), sin necesidad de parsear el path a mano.
- **Helpers HTTP** (`handlers/http_helpers.go`): `WriteJSON`/`WriteError` (respuesta consistente en JSON) y `ParseID` (parseo de parámetros de ruta a `int32`), reutilizados por todos los paquetes de handlers.

### Rutas expuestas

| Recurso | Método y ruta | Descripción |
|---|---|---|
| Clientes | `POST /clientes` | Crea un cliente (`nombre`). `201` + objeto creado, o `400` si el nombre está vacío. |
| | `GET /clientes` | Lista todos los clientes. |
| | `GET /clientes/{id}` | Obtiene un cliente. `404` si no existe. |
| | `PUT /clientes/{id}` | Actualiza un cliente. `404` si no existe. |
| | `DELETE /clientes/{id}` | Elimina un cliente. `204`, o `404` si no existe. |
| Desarrolladores | `POST /desarrolladores` | Crea un desarrollador (`nombre`, `apellido`). |
| | `GET /desarrolladores` | Lista todos los desarrolladores. |
| | `GET /desarrolladores/{id}` | Obtiene un desarrollador. `404` si no existe. |
| | `PUT /desarrolladores/{id}` | Actualiza un desarrollador. `404` si no existe. |
| | `DELETE /desarrolladores/{id}` | Elimina un desarrollador. `204`, o `404` si no existe. |
| Proyectos | `POST /proyectos` | Crea un proyecto (`nombre`, `descripcion`, `estado`, `gastos`, `presupuesto`, `tiempo_estimado` en RFC3339, `id_cliente`). |
| | `GET /proyectos` | Lista todos los proyectos. |
| | `GET /proyectos/{id}` | Obtiene un proyecto. `404` si no existe. |
| | `PUT /proyectos/{id}` | Actualiza un proyecto. `404` si no existe. |
| | `DELETE /proyectos/{id}` | Elimina un proyecto. `204`, o `404` si no existe. |
| Desarrolladores en proyectos (N:M) | `POST /proyectos/{id}/desarrolladores` | Asigna un desarrollador (`id_desarrollador` en el body) a un proyecto. `404` si el proyecto o el desarrollador no existen. |
| | `GET /proyectos/{id}/desarrolladores` | Lista los desarrolladores asignados a un proyecto. |
| | `DELETE /proyectos/{id}/desarrolladores/{devId}` | Quita a un desarrollador de un proyecto. `404` si esa asignación no existe. |

La relación `desarrolladores_en_proyectos` no tiene un ID propio (su clave es el par `id_desarrollador` + `id_proyecto`), por eso se modela como sub-recurso de `proyectos` en vez de como un CRUD independiente con `{id}` — no hay `PUT` (no tiene más campos que editar) ni un `GET` de una sola asignación.

### Validación y códigos de estado

Cada `POST`/`PUT` valida los campos obligatorios de su entidad (que los campos de texto no estén vacíos, que `presupuesto`/`gastos` no sean negativos, que `id_cliente`/`id_desarrollador` estén presentes) **antes** de tocar la base de datos, devolviendo `400` si la validación falla. En `PUT`/`DELETE`/el `POST` de la relación N:M, antes de modificar se verifica que el recurso exista (con un `Get` interno), porque las operaciones `UPDATE`/`DELETE` generadas por `sqlc` no informan si no afectaron ninguna fila, y una foreign key inválida daría un `500` de Postgres en vez de un `404` claro.

| Situación | Código |
|---|---|
| Creación exitosa | `201 Created` |
| Lectura/actualización exitosa | `200 OK` |
| Borrado exitoso | `204 No Content` |
| Dato inválido | `400 Bad Request` |
| Recurso no encontrado | `404 Not Found` |
| Falla inesperada de la base | `500 Internal Server Error` |

### Persistencia (TP2)

- **Motor**: PostgreSQL 16, corriendo en un contenedor Docker.
- **Esquema** (`db/schema/schema.sql`): define las cuatro tablas del dominio.
  - `clientes.id_cliente` y `desarrolladores.id_desarrollador` son `SERIAL` (autoincremental) y clave primaria de su tabla.
  - `proyectos.id` es clave primaria, y `proyectos.id_cliente` es clave foránea hacia `clientes`, modelando la relación 1:N cliente → proyectos.
  - `desarrolladores_en_proyectos` es una tabla intermedia con clave primaria compuesta (`id_desarrollador`, `id_proyecto`) y claves foráneas hacia ambas tablas, modelando la relación N:M.
  - Todas las tablas tienen `created_at` con default `CURRENT_TIMESTAMP` para trazabilidad de alta.
  - Al arrancar el contenedor de Postgres, este archivo se monta en `/docker-entrypoint-initdb.d/` y se ejecuta automáticamente una única vez (cuando el volumen de datos está vacío), creando el esquema desde cero.
- **Consultas** (`db/queries/*.sql`): una consulta anotada por operación CRUD de cada entidad, siguiendo la convención de `sqlc` (`-- name: NombreQuery :one|:many|:exec`).
- **Código generado** (`db/sqlc/`): a partir del esquema y las queries, `sqlc` genera structs Go (`models.go`) y métodos tipados que envuelven `database/sql`. No se edita a mano; se regenera con `sqlc generate` (ver `sqlc.yaml`).
- **Conexión**: `main.go` arma la cadena de conexión a partir de las variables de entorno `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD` y `DB_NAME`. Dentro de `docker-compose.yml`, el servicio `app` se conecta a `db:5432` (nombre de servicio + puerto interno de Postgres, no el puerto mapeado hacia el host).

### Tests

- **De persistencia** (`clientes_test.go`, `desarrolladores_test.go`, `proyectos_test.go`, `desarrolladores_en_proyectos_test.go`, `testdb_test.go`): usan el paquete `testing` de Go y ejercitan `db/sqlc` directamente, sin pasar por HTTP. Se corren con `./test.sh`.
- **De la API** (`requests.hurl`): pruebas de caja negra sobre los endpoints HTTP, con el servidor corriendo. Se corren con `docker compose run --rm hurl` (o `hurl --test requests.hurl` si tenés Hurl instalado localmente).
