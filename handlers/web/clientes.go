package web

import (
	"net/http"
	"strings"

	db "tp2/db/sqlc"
)

// POST /clientes
func CreateCliente(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			redirectError(w, r, "invalido")
			return
		}

		nombre := strings.TrimSpace(r.PostFormValue("nombre"))
		if nombre == "" {
			redirectError(w, r, "invalido")
			return
		}

		if _, err := q.CreateCliente(r.Context(), nombre); err != nil {
			fail(w, r, err)
			return
		}
		redirectOK(w, r)
	}
}
