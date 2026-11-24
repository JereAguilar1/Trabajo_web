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

// POST /desarrolladores
func CreateDesarrollador(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Datos inválidos", http.StatusBadRequest)
			return
		}

		nombre := strings.TrimSpace(r.PostFormValue("nombre"))
		apellido := strings.TrimSpace(r.PostFormValue("apellido"))
		if nombre == "" || apellido == "" {
			http.Error(w, "Datos inválidos", http.StatusBadRequest)
			return
		}

		_, err := q.CreateDesarrollador(r.Context(), db.CreateDesarrolladorParams{
			Nombre:   nombre,
			Apellido: apellido,
		})
		if err != nil {
			fail500(w, err)
			return
		}

		desarrolladores, err := q.ListDesarrolladores(r.Context())
		if err != nil {
			fail500(w, err)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := views.DesarrolladorList(desarrolladores).Render(r.Context(), w); err != nil {
			fail500(w, err)
			return
		}
		if filas, err := cargarProyectosFilas(r.Context(), q); err == nil {
			_ = views.ProyectoListOOB(filas).Render(r.Context(), w)
		}
	}
}

// DELETE /desarrolladores/{id}
func DeleteDesarrollador(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		if err := q.DeleteDesarrollador(r.Context(), int32(id)); err != nil {
			if pgCode(err) == fkViolation {
				http.Error(w, "No se puede eliminar: el desarrollador está asignado a proyectos.", http.StatusConflict)
				return
			}
			fail500(w, err)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if filas, err := cargarProyectosFilas(r.Context(), q); err == nil {
			_ = views.ProyectoListOOB(filas).Render(r.Context(), w)
		}
	}
}

// GET /desarrolladores[?edit=ID]
func ListDesarrolladores(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		desarrolladores, err := q.ListDesarrolladores(r.Context())
		if err != nil {
			fail500(w, err)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := views.DesarrolladorListEdit(desarrolladores, editParam(r)).Render(r.Context(), w); err != nil {
			fail500(w, err)
		}
	}
}

// PUT /desarrolladores/{id}
func UpdateDesarrollador(q *db.Queries) http.HandlerFunc {
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
		apellido := strings.TrimSpace(r.PostFormValue("apellido"))
		if nombre == "" || apellido == "" {
			http.Error(w, "Datos inválidos", http.StatusBadRequest)
			return
		}

		if _, err := q.GetDesarrollador(r.Context(), int32(id)); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "El desarrollador no existe", http.StatusNotFound)
				return
			}
			fail500(w, err)
			return
		}
		err = q.UpdateDesarrollador(r.Context(), db.UpdateDesarrolladorParams{
			IDDesarrollador: int32(id),
			Nombre:          nombre,
			Apellido:        apellido,
		})
		if err != nil {
			fail500(w, err)
			return
		}

		desarrolladores, err := q.ListDesarrolladores(r.Context())
		if err != nil {
			fail500(w, err)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := views.DesarrolladorList(desarrolladores).Render(r.Context(), w); err != nil {
			fail500(w, err)
			return
		}
		// el nombre aparece en la tabla de proyectos (asignados y disponibles)
		if filas, err := cargarProyectosFilas(r.Context(), q); err == nil {
			_ = views.ProyectoListOOB(filas).Render(r.Context(), w)
		}
	}
}
