package desarrolladores

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"tp2/handlers"
	db "tp2/db/sqlc"
)

func validateDesarrollador(nombre, apellido string) error {
	if strings.TrimSpace(nombre) == "" {
		return errors.New("el nombre no puede estar vacío")
	}
	if strings.TrimSpace(apellido) == "" {
		return errors.New("el apellido no puede estar vacío")
	}
	return nil
}

// POST /desarrolladores
func Create(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body db.CreateDesarrolladorParams
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			handlers.WriteError(w, http.StatusBadRequest, "cuerpo JSON inválido")
			return
		}
		if err := validateDesarrollador(body.Nombre, body.Apellido); err != nil {
			handlers.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		desarrollador, err := q.CreateDesarrollador(r.Context(), body)
		if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, "no se pudo crear el desarrollador")
			return
		}
		handlers.WriteJSON(w, http.StatusCreated, desarrollador)
	}
}

// GET /desarrolladores
func GetAll(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		desarrolladores, err := q.ListDesarrolladores(r.Context())
		if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, "no se pudo devolver el listado de desarrolladores")
			return
		}
		handlers.WriteJSON(w, http.StatusOK, desarrolladores)
	}
}

// GET /desarrolladores/{id}
func GetByID(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
		if err != nil {
			handlers.WriteError(w, http.StatusBadRequest, "el id recibido no es válido")
			return
		}

		desarrollador, err := q.GetDesarrollador(r.Context(), int32(id))
		if errors.Is(err, sql.ErrNoRows) {
			handlers.WriteError(w, http.StatusNotFound, "el desarrollador no existe")
			return
		}
		if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, "no se pudo obtener el desarrollador")
			return
		}
		handlers.WriteJSON(w, http.StatusOK, desarrollador)
	}
}

// PUT /desarrolladores/{id}
func Update(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
		if err != nil {
			handlers.WriteError(w, http.StatusBadRequest, "el id recibido no es válido")
			return
		}

		var body db.UpdateDesarrolladorParams
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			handlers.WriteError(w, http.StatusBadRequest, "cuerpo JSON inválido")
			return
		}
		if err := validateDesarrollador(body.Nombre, body.Apellido); err != nil {
			handlers.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		body.IDDesarrollador = int32(id)

		if _, err := q.GetDesarrollador(r.Context(), int32(id)); errors.Is(err, sql.ErrNoRows) {
			handlers.WriteError(w, http.StatusNotFound, "el desarrollador no existe")
			return
		} else if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, "no se pudo verificar el desarrollador")
			return
		}

		if err := q.UpdateDesarrollador(r.Context(), body); err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, "el desarrollador no se pudo actualizar")
			return
		}

		desarrollador, err := q.GetDesarrollador(r.Context(), int32(id))
		if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, "no se pudo obtener el desarrollador actualizado")
			return
		}
		handlers.WriteJSON(w, http.StatusOK, desarrollador)
	}
}

// DELETE /desarrolladores/{id}
func Delete(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
		if err != nil {
			handlers.WriteError(w, http.StatusBadRequest, "el id recibido no es válido")
			return
		}

		if _, err := q.GetDesarrollador(r.Context(), int32(id)); errors.Is(err, sql.ErrNoRows) {
			handlers.WriteError(w, http.StatusNotFound, "el desarrollador no existe")
			return
		} else if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, "no se pudo verificar el desarrollador")
			return
		}

		if err := q.DeleteDesarrollador(r.Context(), int32(id)); err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, "el desarrollador no se pudo eliminar")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
