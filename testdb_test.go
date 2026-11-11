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

	// Sin valores por defecto: las credenciales vienen del entorno (.env).
	cfg := map[string]string{}
	for _, k := range []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME"} {
		cfg[k] = os.Getenv(k)
		if cfg[k] == "" {
			t.Skipf("falta la variable de entorno %s (ver .env.example); se omiten los tests de BD", k)
		}
	}

	connStr := "host=" + cfg["DB_HOST"] + " port=" + cfg["DB_PORT"] + " user=" + cfg["DB_USER"] +
		" password=" + cfg["DB_PASSWORD"] + " dbname=" + cfg["DB_NAME"] + " sslmode=disable"

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Error al abrir la conexión: %v", err)
	}
	t.Cleanup(func() {
		conn.Close()
	})

	if err := conn.Ping(); err != nil {
		t.Fatalf("Error al hacer ping a la base de datos (%s:%s): %v", cfg["DB_HOST"], cfg["DB_PORT"], err)
	}

	return db.New(conn)
}
