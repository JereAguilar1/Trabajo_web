package desarrolladores_en_proyectos

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"tp2/handlers"
	db "tp2/db/sqlc"
)

// POST /proyectos/{id}/desarrolladores  { "id_desarrollador": N }
// Asigna un desarrollador existente a un proyecto existente.
func Assign(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idProyecto, err := handlers.ParseID(r, "id")
		if err != nil {
			handlers.WriteError(w, http.StatusBadRequest, "el id de proyecto no es válido")
			return
		}

		var body struct {
			IDDesarrollador int32 `json:"id_desarrollador"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			handlers.WriteError(w, http.StatusBadRequest, "cuerpo JSON inválido")
			return
		}
		if body.IDDesarrollador <= 0 {
			handlers.WriteError(w, http.StatusBadRequest, "id_desarrollador es obligatorio")
			return
		}

		if _, err := q.GetProyecto(r.Context(), idProyecto); errors.Is(err, sql.ErrNoRows) {
			handlers.WriteError(w, http.StatusNotFound, "el proyecto no existe")
			return
		} else if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, "no se pudo verificar el proyecto")
			return
		}

		if _, err := q.GetDesarrollador(r.Context(), body.IDDesarrollador); errors.Is(err, sql.ErrNoRows) {
			handlers.WriteError(w, http.StatusNotFound, "el desarrollador no existe")
			return
		} else if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, "no se pudo verificar el desarrollador")
			return
		}

		relacion, err := q.CreateDesarrolladorEnProyecto(r.Context(), db.CreateDesarrolladorEnProyectoParams{
			IDDesarrollador: body.IDDesarrollador,
			IDProyecto:      idProyecto,
		})
		if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, "no se pudo asignar el desarrollador al proyecto")
			return
		}
		handlers.WriteJSON(w, http.StatusCreated, relacion)
	}
}

// GET /proyectos/{id}/desarrolladores
// Lista los desarrolladores asignados a un proyecto.
func List(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idProyecto, err := handlers.ParseID(r, "id")
		if err != nil {
			handlers.WriteError(w, http.StatusBadRequest, "el id de proyecto no es válido")
			return
		}

		if _, err := q.GetProyecto(r.Context(), idProyecto); errors.Is(err, sql.ErrNoRows) {
			handlers.WriteError(w, http.StatusNotFound, "el proyecto no existe")
			return
		} else if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, "no se pudo verificar el proyecto")
			return
		}

		asignaciones, err := q.ListDesarrolladoresByProyecto(r.Context(), idProyecto)
		if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, "no se pudo listar los desarrolladores del proyecto")
			return
		}
		handlers.WriteJSON(w, http.StatusOK, asignaciones)
	}
}

// DELETE /proyectos/{id}/desarrolladores/{devId}
// Quita a un desarrollador de un proyecto.
func Remove(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idProyecto, err := handlers.ParseID(r, "id")
		if err != nil {
			handlers.WriteError(w, http.StatusBadRequest, "el id de proyecto no es válido")
			return
		}
		idDesarrollador, err := handlers.ParseID(r, "devId")
		if err != nil {
			handlers.WriteError(w, http.StatusBadRequest, "el id de desarrollador no es válido")
			return
		}

		params := db.GetDesarrolladorEnProyectoParams{
			IDDesarrollador: idDesarrollador,
			IDProyecto:      idProyecto,
		}
		if _, err := q.GetDesarrolladorEnProyecto(r.Context(), params); errors.Is(err, sql.ErrNoRows) {
			handlers.WriteError(w, http.StatusNotFound, "el desarrollador no está asignado a ese proyecto")
			return
		} else if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, "no se pudo verificar la asignación")
			return
		}

		if err := q.DeleteDesarrolladorEnProyecto(r.Context(), db.DeleteDesarrolladorEnProyectoParams{
			IDDesarrollador: idDesarrollador,
			IDProyecto:      idProyecto,
		}); err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, "no se pudo quitar al desarrollador del proyecto")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
