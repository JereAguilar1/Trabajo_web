package web

import (
	"net/http"
	"strings"

	db "tp2/db/sqlc"
)

// POST /desarrolladores
func CreateDesarrollador(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			redirectError(w, r, "invalido")
			return
		}

		nombre := strings.TrimSpace(r.PostFormValue("nombre"))
		apellido := strings.TrimSpace(r.PostFormValue("apellido"))
		if nombre == "" || apellido == "" {
			redirectError(w, r, "invalido")
			return
		}

		_, err := q.CreateDesarrollador(r.Context(), db.CreateDesarrolladorParams{
			Nombre:   nombre,
			Apellido: apellido,
		})
		if err != nil {
			fail(w, r, err)
			return
		}
		redirectOK(w, r)
	}
}
