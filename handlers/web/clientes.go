package web

import (
	"database/sql"
	"errors"
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

// editParam lee ?edit=ID de la URL (0 si falta o no es un número válido).
func editParam(r *http.Request) int32 {
	id, err := strconv.ParseInt(r.URL.Query().Get("edit"), 10, 32)
	if err != nil {
		return 0
	}
	return int32(id)
}

// GET /clientes[?edit=ID]
// Devuelve la tabla de clientes; con ?edit=ID esa fila sale como formulario.
func ListClientes(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientes, err := q.ListClientes(r.Context())
		if err != nil {
			fail500(w, err)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := views.ClienteListEdit(clientes, editParam(r)).Render(r.Context(), w); err != nil {
			fail500(w, err)
		}
	}
}

// PUT /clientes/{id}
func UpdateCliente(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Datos inválidos", http.StatusBadRequest)
			return
		}
		nombre := strings.TrimSpace(r.PostFormValue("nombre"))
		if nombre == "" {
			http.Error(w, "Datos inválidos", http.StatusBadRequest)
			return
		}

		// UPDATE no falla si el id no existe (0 filas afectadas): se verifica antes.
		if _, err := q.GetCliente(r.Context(), int32(id)); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "El cliente no existe", http.StatusNotFound)
				return
			}
			fail500(w, err)
			return
		}
		if err := q.UpdateCliente(r.Context(), db.UpdateClienteParams{IDCliente: int32(id), Nombre: nombre}); err != nil {
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
		// el nombre aparece en el <select> de proyectos y en su tabla
		_ = views.ProyectoFormOOB(clientes).Render(r.Context(), w)
		if filas, err := cargarProyectosFilas(r.Context(), q); err == nil {
			_ = views.ProyectoListOOB(filas).Render(r.Context(), w)
		}
	}
}
