# Trabajo Práctico 2: Persistiendo el Dominio

Capa de persistencia (esquema, consultas y código de acceso a datos generado con `sqlc`) para el dominio del TP1: un Sistema de Gestión de Proyectos.

## Requisitos

- [Go](https://go.dev/) 1.22+
- [sqlc](https://sqlc.dev/)
- [Docker](https://www.docker.com/) y Docker Compose
- Puerto **5432** libre
## Cómo ejecutar

```bash
git clone https://github.com/JereAguilar1/Trabajo_web.git
cd Trabajo_web
git checkout tp2
./test.sh
```

`test.sh` realiza lo siguiente:

1. **Tareas previas**: regenera el código con `sqlc generate`, compila el proyecto (`go build ./...`), destruye cualquier contenedor/volumen previo, levanta un Postgres limpio con Docker Compose y espera (haciendo polling con `pg_isready`, hasta 30 intentos de 1 segundo) a que esté aceptando conexiones.
2. **Ejecución de tests**: corre `go test ./...`, que ejercita el paquete estándar `testing` de Go contra la base real.
3. **Tareas posteriores**: destruye el contenedor y el volumen de Postgres, se hayan pasado los tests o no, dejando el entorno limpio.

El código de salida del script es el mismo que devuelve `go test` (0 si todo pasó, distinto de 0 si algo falló).

## Documentación

### Persistencia

- **Motor**: PostgreSQL 16, corriendo en un contenedor Docker efímero (se recrea en cada corrida de `test.sh`, sin datos persistentes entre ejecuciones).
- **Esquema** (`db/schema/schema.sql`): define las cuatro tablas del dominio.
  - `clientes.id_cliente` y `desarrolladores.id_desarrollador` son `SERIAL` (autoincremental) y clave primaria de su tabla.
  - `proyectos.id` es clave primaria, y `proyectos.id_cliente` es clave foránea hacia `clientes`, modelando la relación 1:N cliente → proyectos.
  - `desarrolladores_en_proyectos` es una tabla intermedia con clave primaria compuesta (`id_desarrollador`, `id_proyecto`) y claves foráneas hacia ambas tablas, modelando la relación N:M.
  - Todas las tablas tienen `created_at` con default `CURRENT_TIMESTAMP` para trazabilidad de alta.
  - Al arrancar el contenedor de Postgres, este archivo se monta en `/docker-entrypoint-initdb.d/` y se ejecuta automáticamente una única vez (cuando el volumen de datos está vacío), creando el esquema desde cero.
- **Consultas** (`db/queries/*.sql`): una consulta anotada por operación CRUD de cada entidad, siguiendo la convención de `sqlc` (`-- name: NombreQuery :one|:many|:exec`).
- **Código generado** (`db/sqlc/`): a partir del esquema y las queries, `sqlc` genera structs Go (`models.go`) y métodos tipados (`clientes.sql.go`, `desarrolladores.sql.go`, `proyectos.sql.go`, `desarrolladores_en_proyectos.sql.go`) que envuelven `database/sql`. No se edita a mano; se regenera con `sqlc generate` (ver `sqlc.yaml`).
- **Conexión**: tanto `main.go` como los tests (`testdb_test.go`) arman la cadena de conexión a partir de las variables de entorno `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD` y `DB_NAME`, con defaults que apuntan al Postgres levantado por `docker-compose.yml` (`127.0.0.1:5432`, `postgres`/`postgres`, `tp2_db`).

### Tests

- `testdb_test.go` — helper compartido (`getTestQueries`) que conecta a la base de test usando el paquete `testing`.
- `clientes_test.go`, `desarrolladores_test.go`, `proyectos_test.go`, `desarrolladores_en_proyectos_test.go` — tests cubriendo el CRUD completo de cada entidad contra una base Postgres real.
