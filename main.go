package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	db "tp2/db/sqlc"
)

func main() {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
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
		fmt.Printf("No se pudo conectar a PostgreSQL en %s:%s (%v). Ejecute ./test.sh.\n", host, port, err)
		return
	}

	fmt.Printf("Conexión establecida con PostgreSQL en %s:%s (DB: %s)\n", host, port, dbname)

	queries := db.New(database)
	_ = queries
}
