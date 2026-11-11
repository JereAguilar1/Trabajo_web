#!/usr/bin/env bash
# TP5 - Renderizado del lado del servidor con templ
# Uso: ./run-app.sh     (equivale a `make run-app`)
set -euo pipefail
cd "$(dirname "$0")"

# --- Tareas posteriores: se ejecutan siempre al terminar (Ctrl+C, error o fin) ---
cleanup() {
  echo
  echo "Borrando contenedores y volúmenes..."
  docker compose down -v --remove-orphans
}
trap cleanup EXIT
trap 'exit 130' INT TERM

# --- Tareas previas ---

# crea .env a partir de .env.example (con contraseña aleatoria) si no existe
if [ ! -f .env ]; then
  cp .env.example .env
  pass=$(head -c 24 /dev/urandom | base64 | tr -dc A-Za-z0-9)
  sed -i "s/^POSTGRES_PASSWORD=.*/POSTGRES_PASSWORD=$pass/" .env
  echo "Creado .env a partir de .env.example (con contraseña aleatoria)"
fi

# la app lee DB_*; se arman desde las variables de .env
set -a; . ./.env; set +a
export DB_HOST=127.0.0.1 DB_USER="$POSTGRES_USER" DB_PASSWORD="$POSTGRES_PASSWORD" \
       DB_NAME="$POSTGRES_DB" DB_PORT="$DB_HOST_PORT"

# genera el código: sqlc (db/sqlc) y templ (views/*_templ.go)
sqlc generate
go run github.com/a-h/templ/cmd/templ@v0.3.1100 generate

# compila
go build -o bin/app .

# borra contenedores y volúmenes de ejecuciones anteriores, levanta la BD y
# espera a que esté lista (healthcheck)
docker compose down -v --remove-orphans
docker compose up -d --wait db

# --- Ejecución ---
echo "Aplicación en http://localhost:8090 (Ctrl+C para detener)"
./bin/app
