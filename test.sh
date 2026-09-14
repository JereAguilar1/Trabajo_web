#!/usr/bin/env bash

# --- Tareas previas ---

sqlc generate

go build ./...

docker compose down -v --remove-orphans
docker compose up -d

# espera a que postgres este listo
MAX_ATTEMPTS=30
ATTEMPT=0
until docker compose exec -T db pg_isready -U postgres -d tp2_db >/dev/null 2>&1; do
  ATTEMPT=$((ATTEMPT + 1))
  if [ "$ATTEMPT" -ge "$MAX_ATTEMPTS" ]; then
    echo "Error: PostgreSQL no respondió a tiempo."
    docker compose logs
    docker compose down -v --remove-orphans
    exit 1
  fi
  sleep 1
done

# --- Ejecución de tests ---

go test ./...
# guarda codigo de salida
TEST_RESULT=$?

# --- Tareas posteriores ---

docker compose down -v --remove-orphans

exit $TEST_RESULT
