package web

import (
	"net/http"
	"strconv"
	"strings"

	db "tp2/db/sqlc"
	"tp2/views"
)

// POST /clientes
func CreateCliente(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Datos inválidos", http.StatusBadRequest)
			return
		}

		nombre := strings.TrimSpace(r.PostFormValue("nombre"))
		if nombre == "" {
			http.Error(w, "Datos inválidos", http.StatusBadRequest)
			return
		}

		if _, err := q.CreateCliente(r.Context(), nombre); err != nil {
			fail500(w, err)
			return
		}

		clientes, err := q.ListClientes(r.Context())
		if err != nil {
			fail500(w, err)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := views.ClienteList(clientes).Render(r.Context(), w); err != nil {
			fail500(w, err)
			return
		}
		_ = views.ProyectoFormOOB(clientes).Render(r.Context(), w)
	}
}

// DELETE /clientes/{id}
func DeleteCliente(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		if err := q.DeleteCliente(r.Context(), int32(id)); err != nil {
			if pgCode(err) == fkViolation {
				http.Error(w, "No se puede eliminar: el cliente tiene proyectos asociados.", http.StatusConflict)
				return
			}
			fail500(w, err)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if clientes, err := q.ListClientes(r.Context()); err == nil {
			_ = views.ProyectoFormOOB(clientes).Render(r.Context(), w)
		}
	}
}
