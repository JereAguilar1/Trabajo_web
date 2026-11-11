package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"
	db "tp2/db/sqlc"
	"tp2/handlers/web"
)

func main() {
	// La conexión se configura solo con variables de entorno (sin valores
	// por defecto): no hay credenciales en el código. Ver .env.example.
	cfg := map[string]string{}
	var faltan []string
	for _, k := range []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME"} {
		cfg[k] = os.Getenv(k)
		if cfg[k] == "" {
			faltan = append(faltan, k)
		}
	}
	if len(faltan) > 0 {
		log.Fatalf("Faltan variables de entorno: %s (ver .env.example)", strings.Join(faltan, ", "))
	}
	host, port, dbname := cfg["DB_HOST"], cfg["DB_PORT"], cfg["DB_NAME"]

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, cfg["DB_USER"], cfg["DB_PASSWORD"], dbname)

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

	// Vistas renderizadas en el servidor (templ). Reemplazan a la API JSON
	// del TP3 y al frontend JS del TP4. Los formularios HTML solo soportan
	// GET y POST.
	// {$} hace que "/" matchee solo la raíz y no cualquier ruta desconocida.
	mux.HandleFunc("GET /{$}", web.Index(queries))

	mux.HandleFunc("POST /clientes", web.CreateCliente(queries))

	mux.HandleFunc("POST /desarrolladores", web.CreateDesarrollador(queries))

	mux.HandleFunc("POST /proyectos", web.CreateProyecto(queries))

	// Uso el 8090 porque el 8080 lo tengo ocupado
	http.ListenAndServe(":8090", mux)
}
