package clientes

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

// POST /clientes
func Create(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Nombre string `json:"nombre"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			handlers.WriteError(w, http.StatusBadRequest, "cuerpo JSON inválido")
			return
		}

		if strings.TrimSpace(body.Nombre) == "" {
			handlers.WriteError(w, http.StatusBadRequest, "el nombre no puede estar vacío")
			return
		}

		cliente, err := q.CreateCliente(r.Context(), body.Nombre)
		if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, "no se pudo crear el cliente")
			return
		}

		handlers.WriteJSON(w, http.StatusCreated, cliente)
	}
}

//GET /clientes

func GetAll(q* db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		clientes, err := q.ListClientes(r.Context())
		
		if err!= nil{
			handlers.WriteError(w, http.StatusInternalServerError, "no se pudo devolver el listado de clientes")
			return
		}
		handlers.WriteJSON(w, http.StatusOK, clientes)
	}
}

//GET /cliente {id}

func GetByID(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
		if err != nil {
			handlers.WriteError(w, http.StatusBadRequest, "el id recibido no es válido")
			return
		}

		cliente, err := q.GetCliente(r.Context(), int32(id))
		if errors.Is(err, sql.ErrNoRows) {
			handlers.WriteError(w, http.StatusNotFound, "el cliente no existe")
			return
		}
		if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, "no se pudo obtener el cliente")
			return
		}
		handlers.WriteJSON(w, http.StatusOK, cliente)
	}
}

func Update(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
		if err != nil {
			handlers.WriteError(w, http.StatusBadRequest, "el id recibido no es válido")
			return
		}

		var body struct {
			Nombre string `json:"nombre"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			handlers.WriteError(w, http.StatusBadRequest, "cuerpo JSON inválido")
			return
		}
		if strings.TrimSpace(body.Nombre) == "" {
			handlers.WriteError(w, http.StatusBadRequest, "el nombre no puede estar vacío")
			return
		}

		if _, err := q.GetCliente(r.Context(), int32(id)); errors.Is(err, sql.ErrNoRows) {
			handlers.WriteError(w, http.StatusNotFound, "el cliente no existe")
			return
		} else if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, "no se pudo verificar el cliente")
			return
		}

		if err := q.UpdateCliente(r.Context(), db.UpdateClienteParams{
			IDCliente: int32(id),
			Nombre:    body.Nombre,
		}); err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, "el cliente no se pudo actualizar")
			return
		}

		cliente, err := q.GetCliente(r.Context(), int32(id))
		if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, "no se pudo obtener el cliente actualizado")
			return
		}
		handlers.WriteJSON(w, http.StatusOK, cliente)
	}
}

func Delete(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
		if err != nil {
			handlers.WriteError(w, http.StatusBadRequest, "el id recibido no es válido")
			return
		}

		if _, err := q.GetCliente(r.Context(), int32(id)); errors.Is(err, sql.ErrNoRows) {
			handlers.WriteError(w, http.StatusNotFound, "el cliente no existe")
			return
		} else if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, "no se pudo verificar el cliente")
			return
		}

		if err := q.DeleteCliente(r.Context(), int32(id)); err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, "el cliente no se pudo eliminar")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}