# TP5 - Renderizado del lado del servidor con templ
#
#   make run-app   -> tareas previas + ejecuta la aplicación + tareas posteriores
#
# La configuración de la BD sale de .env (se crea desde .env.example si falta).
# Se carga dentro de cada receta con ". ./.env" para que funcione también la
# primera vez, cuando .env todavía no existía al leer este Makefile.

.PHONY: test run-app env generate build clean-containers db-up run down docker

SHELL := /bin/bash
BIN = bin/app

# Variables que lee main.go (DB_*), traducidas desde las de .env
LOAD_ENV = set -a; . ./.env; set +a; \
           export DB_HOST=127.0.0.1 DB_USER=$$POSTGRES_USER DB_PASSWORD=$$POSTGRES_PASSWORD \
                  DB_NAME=$$POSTGRES_DB DB_PORT=$$DB_HOST_PORT

# Entrega: todo en orden. La BD se borra al terminar (también con Ctrl+C).
run-app: env generate build clean-containers db-up run

# crea .env a partir de .env.example (con contraseña aleatoria) si no existe
env:
	@[ -f .env ] || (cp .env.example .env \
	  && sed -i "s/^POSTGRES_PASSWORD=.*/POSTGRES_PASSWORD=$$(head -c 24 /dev/urandom | base64 | tr -dc A-Za-z0-9)/" .env \
	  && echo "Creado .env a partir de .env.example (con contraseña aleatoria)")

# genera el código Go: sqlc (db/sqlc) y templ (views/*_templ.go)
generate:
	sqlc generate
	go run github.com/a-h/templ/cmd/templ@v0.3.1100 generate

# compila la aplicación
build:
	go build -o $(BIN) .

# borra contenedores y volúmenes de ejecuciones anteriores
clean-containers:
	docker compose down -v --remove-orphans

# levanta PostgreSQL y espera a que esté healthy (el schema se carga al iniciar)
db-up: env
	docker compose up -d --wait db

# corre el servidor (:8090) y, al terminar o con Ctrl+C, borra contenedores y volúmenes
run: env
	@$(LOAD_ENV); \
	trap 'echo; echo "Borrando contenedores y volúmenes..."; docker compose down -v --remove-orphans' EXIT; \
	trap 'exit 130' INT TERM; \
	echo "Aplicación en http://localhost:8090 (Ctrl+C para detener)"; \
	./$(BIN)

# alternativa: BD y app en contenedores (la app se compila en la imagen)
docker: env generate
	docker compose up -d --build db app

# tareas posteriores a mano: borra contenedores y volúmenes
down:
	docker compose down -v --remove-orphans

# tests de persistencia (go test) contra una BD temporal; borra contenedores y volúmenes al terminar
test: env generate clean-containers db-up
	@$(LOAD_ENV); \
	trap 'docker compose down -v --remove-orphans' EXIT; \
	go test ./...
