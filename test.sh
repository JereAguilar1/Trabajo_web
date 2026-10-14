#!/usr/bin/env bash

# --- Tareas previas ---

# crea .env a partir de .env.example si no existe
if [ ! -f .env ]; then
  cp .env.example .env
  echo "Creado .env a partir de .env.example"
fi

# carga las variables y las traduce a las que usan los tests de Go
set -a; source .env; set +a
export DB_USER="$POSTGRES_USER" DB_PASSWORD="$POSTGRES_PASSWORD" DB_NAME="$POSTGRES_DB" DB_PORT="$DB_HOST_PORT"

sqlc generate

go build ./...

docker compose down -v --remove-orphans
docker compose up -d db

# espera a que postgres este listo
MAX_ATTEMPTS=30
ATTEMPT=0
until docker compose exec -T db pg_isready -U "$POSTGRES_USER" -d "$POSTGRES_DB" >/dev/null 2>&1; do
  ATTEMPT=$((ATTEMPT + 1))
  if [ "$ATTEMPT" -ge "$MAX_ATTEMPTS" ]; then
    echo "Error: PostgreSQL no respondió a tiempo."
    docker compose logs
    docker compose down -v --remove-orphans
    exit 1
  fi
  sleep 1
done

# levanta la aplicacion (compila la imagen) y espera a que responda
docker compose up -d --build app

ATTEMPT=0
until curl -s -o /dev/null http://localhost:8090/clientes; do
  ATTEMPT=$((ATTEMPT + 1))
  if [ "$ATTEMPT" -ge "$MAX_ATTEMPTS" ]; then
    echo "Error: la aplicación no respondió a tiempo."
    docker compose logs app
    docker compose down -v --remove-orphans
    exit 1
  fi
  sleep 1
done

# --- Ejecución de tests ---

# tests de persistencia
go test ./...
TEST_RESULT=$?

# pruebas de la API (cada endpoint, con Hurl)
docker compose run --rm hurl
HURL_RESULT=$?

# --- Tareas posteriores ---

docker compose down -v --remove-orphans

if [ "$TEST_RESULT" -ne 0 ] || [ "$HURL_RESULT" -ne 0 ]; then
  exit 1
fi
exit 0
