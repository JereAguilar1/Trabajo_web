package proyectos

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

func validateProyecto(nombre, descripcion, estado string, gastos, presupuesto float64, idCliente int32) error {
	if strings.TrimSpace(nombre) == "" {
		return errors.New("el nombre no puede estar vacío")
	}
	if strings.TrimSpace(descripcion) == "" {
		return errors.New("la descripción no puede estar vacía")
	}
	if strings.TrimSpace(estado) == "" {
		return errors.New("el estado no puede estar vacío")
	}
	if gastos < 0 {
		return errors.New("los gastos no pueden ser negativos")
	}
	if presupuesto < 0 {
		return errors.New("el presupuesto no puede ser negativo")
	}
	if idCliente <= 0 {
		return errors.New("id_cliente es obligatorio")
	}
	return nil
}

// POST /proyectos
func Create(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body db.CreateProyectoParams
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			handlers.WriteError(w, http.StatusBadRequest, "cuerpo JSON inválido")
			return
		}
		if err := validateProyecto(body.Nombre, body.Descripcion, body.Estado, body.Gastos, body.Presupuesto, body.IDCliente); err != nil {
			handlers.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		proyecto, err := q.CreateProyecto(r.Context(), body)
		if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, "no se pudo crear el proyecto")
			return
		}
		handlers.WriteJSON(w, http.StatusCreated, proyecto)
	}
}

// GET /proyectos
func GetAll(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proyectos, err := q.ListProyectos(r.Context())
		if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, "no se pudo devolver el listado de proyectos")
			return
		}
		handlers.WriteJSON(w, http.StatusOK, proyectos)
	}
}

// GET /proyectos/{id}
func GetByID(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
		if err != nil {
			handlers.WriteError(w, http.StatusBadRequest, "el id recibido no es válido")
			return
		}

		proyecto, err := q.GetProyecto(r.Context(), int32(id))
		if errors.Is(err, sql.ErrNoRows) {
			handlers.WriteError(w, http.StatusNotFound, "el proyecto no existe")
			return
		}
		if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, "no se pudo obtener el proyecto")
			return
		}
		handlers.WriteJSON(w, http.StatusOK, proyecto)
	}
}

// PUT /proyectos/{id}
func Update(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
		if err != nil {
			handlers.WriteError(w, http.StatusBadRequest, "el id recibido no es válido")
			return
		}

		var body db.UpdateProyectoParams
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			handlers.WriteError(w, http.StatusBadRequest, "cuerpo JSON inválido")
			return
		}
		if err := validateProyecto(body.Nombre, body.Descripcion, body.Estado, body.Gastos, body.Presupuesto, body.IDCliente); err != nil {
			handlers.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		body.ID = int32(id)

		if _, err := q.GetProyecto(r.Context(), int32(id)); errors.Is(err, sql.ErrNoRows) {
			handlers.WriteError(w, http.StatusNotFound, "el proyecto no existe")
			return
		} else if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, "no se pudo verificar el proyecto")
			return
		}

		if err := q.UpdateProyecto(r.Context(), body); err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, "el proyecto no se pudo actualizar")
			return
		}

		proyecto, err := q.GetProyecto(r.Context(), int32(id))
		if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, "no se pudo obtener el proyecto actualizado")
			return
		}
		handlers.WriteJSON(w, http.StatusOK, proyecto)
	}
}

// DELETE /proyectos/{id}
func Delete(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
		if err != nil {
			handlers.WriteError(w, http.StatusBadRequest, "el id recibido no es válido")
			return
		}

		if _, err := q.GetProyecto(r.Context(), int32(id)); errors.Is(err, sql.ErrNoRows) {
			handlers.WriteError(w, http.StatusNotFound, "el proyecto no existe")
			return
		} else if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, "no se pudo verificar el proyecto")
			return
		}

		if err := q.DeleteProyecto(r.Context(), int32(id)); err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, "el proyecto no se pudo eliminar")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
