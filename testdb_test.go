package main

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
	db "tp2/db/sqlc"
)

// getTestQueries abre una conexión a la base de prueba (usando las mismas
// variables de entorno que test.sh) y devuelve un *db.Queries listo para usar.
func getTestQueries(t *testing.T) *db.Queries {
	t.Helper()

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

	connStr := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Error al abrir la conexión: %v", err)
	}
	t.Cleanup(func() {
		conn.Close()
	})

	if err := conn.Ping(); err != nil {
		t.Fatalf("Error al hacer ping a la base de datos (%s): %v", connStr, err)
	}

	return db.New(conn)
}
