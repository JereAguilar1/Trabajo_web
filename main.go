package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"net/http"

	_ "github.com/lib/pq"
	db "tp2/db/sqlc"
	"tp2/handlers/clientes"
	"tp2/handlers/desarrolladores"
	"tp2/handlers/desarrolladores_en_proyectos"
	"tp2/handlers/proyectos"
)

func main() {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5433"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "tp2_db"
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error al inicializar driver de base de datos: %v", err)
	}
	defer database.Close()

	if err := database.Ping(); err != nil {
		log.Fatalf("No se pudo conectar a PostgreSQL en %s:%s (%v).", host, port, err)
	}

	fmt.Printf("Conexión establecida con PostgreSQL en %s:%s (DB: %s)\n", host, port, dbname)

	queries := db.New(database)

	// SERVER
	mux := http.NewServeMux()
	
	//Handlers cliente
	mux.HandleFunc("POST /clientes", clientes.Create(queries))
	mux.HandleFunc("GET /clientes", clientes.GetAll(queries))
	mux.HandleFunc("GET /clientes/{id}", clientes.GetByID(queries))
	mux.HandleFunc("PUT /clientes/{id}", clientes.Update(queries))
	mux.HandleFunc("DELETE /clientes/{id}", clientes.Delete(queries))

	//Handlers desarrolladores
	mux.HandleFunc("POST /desarrolladores", desarrolladores.Create(queries))
	mux.HandleFunc("GET /desarrolladores", desarrolladores.GetAll(queries))
	mux.HandleFunc("GET /desarrolladores/{id}", desarrolladores.GetByID(queries))
	mux.HandleFunc("PUT /desarrolladores/{id}", desarrolladores.Update(queries))
	mux.HandleFunc("DELETE /desarrolladores/{id}", desarrolladores.Delete(queries))

	//Handlers proyectos
	mux.HandleFunc("POST /proyectos", proyectos.Create(queries))
	mux.HandleFunc("GET /proyectos", proyectos.GetAll(queries))
	mux.HandleFunc("GET /proyectos/{id}", proyectos.GetByID(queries))
	mux.HandleFunc("PUT /proyectos/{id}", proyectos.Update(queries))
	mux.HandleFunc("DELETE /proyectos/{id}", proyectos.Delete(queries))

	//Handlers desarrolladores_en_proyectos (relación N:M, sub-recurso de proyectos)
	mux.HandleFunc("POST /proyectos/{id}/desarrolladores", desarrolladores_en_proyectos.Assign(queries))
	mux.HandleFunc("GET /proyectos/{id}/desarrolladores", desarrolladores_en_proyectos.List(queries))
	mux.HandleFunc("DELETE /proyectos/{id}/desarrolladores/{devId}", desarrolladores_en_proyectos.Remove(queries))


// Uso el 8090 porque el 8080 lo tengo ocupado
	http.ListenAndServe(":8090", mux)
}
