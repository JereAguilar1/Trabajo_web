package web

import (
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	db "tp2/db/sqlc"
	"tp2/views"
)

// parseMonto convierte un campo numérico del formulario; rechaza texto,
// NaN, infinitos y negativos.
func parseMonto(s string) (float64, bool) {
	v, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) || v < 0 {
		return 0, false
	}
	return v, true
}

func estadoValido(estado string) bool {
	for _, e := range views.Estados {
		if e.Valor == estado {
			return true
		}
	}
	return false
}

// POST /proyectos
func CreateProyecto(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Datos inválidos", http.StatusBadRequest)
			return
		}

		nombre := strings.TrimSpace(r.PostFormValue("nombre"))
		descripcion := strings.TrimSpace(r.PostFormValue("descripcion"))
		estado := r.PostFormValue("estado")
		presupuesto, okPresupuesto := parseMonto(r.PostFormValue("presupuesto"))
		gastos, okGastos := parseMonto(r.PostFormValue("gastos"))
		plazo, errPlazo := time.Parse("2006-01-02", r.PostFormValue("tiempo_estimado"))
		idCliente, errCliente := strconv.ParseInt(r.PostFormValue("id_cliente"), 10, 32)

		if nombre == "" || descripcion == "" || !estadoValido(estado) ||
			!okPresupuesto || !okGastos || errPlazo != nil || errCliente != nil {
			http.Error(w, "Datos inválidos", http.StatusBadRequest)
			return
		}

		_, err := q.CreateProyecto(r.Context(), db.CreateProyectoParams{
			Nombre:         nombre,
			Descripcion:    descripcion,
			Estado:         estado,
			Gastos:         gastos,
			Presupuesto:    presupuesto,
			TiempoEstimado: plazo,
			IDCliente:      int32(idCliente),
		})
		if err != nil {
			if pgCode(err) == fkViolation {
				http.Error(w, "El cliente no existe", http.StatusBadRequest)
				return
			}
			fail500(w, err)
			return
		}

		filas, err := cargarProyectosFilas(r.Context(), q)
		if err != nil {
			fail500(w, err)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := views.ProyectoList(filas).Render(r.Context(), w); err != nil {
			fail500(w, err)
		}
	}
}

// DELETE /proyectos/{id}
func DeleteProyecto(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		if err := q.DeleteProyecto(r.Context(), int32(id)); err != nil {
			if pgCode(err) == fkViolation {
				http.Error(w, "No se puede eliminar: el proyecto tiene desarrolladores asignados.", http.StatusConflict)
				return
			}
			fail500(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// POST /proyectos/{id}/desarrolladores
func AssignDesarrollador(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		idProyecto, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			http.Error(w, "ID de proyecto inválido", http.StatusBadRequest)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Formulario inválido", http.StatusBadRequest)
			return
		}

		devIdStr := r.PostFormValue("id_desarrollador")
		idDesarrollador, err := strconv.ParseInt(devIdStr, 10, 32)
		if err != nil {
			http.Error(w, "ID de desarrollador inválido", http.StatusBadRequest)
			return
		}

		_, err = q.CreateDesarrolladorEnProyecto(r.Context(), db.CreateDesarrolladorEnProyectoParams{
			IDDesarrollador: int32(idDesarrollador),
			IDProyecto:      int32(idProyecto),
		})
		if err != nil {
			fail500(w, err)
			return
		}

		filas, err := cargarProyectosFilas(r.Context(), q)
		if err != nil {
			fail500(w, err)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := views.ProyectoList(filas).Render(r.Context(), w); err != nil {
			fail500(w, err)
		}
	}
}

// DELETE /proyectos/{id}/desarrolladores/{devId}
func RemoveDesarrollador(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		idProyecto, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			http.Error(w, "ID de proyecto inválido", http.StatusBadRequest)
			return
		}

		devIdStr := r.PathValue("devId")
		idDesarrollador, err := strconv.ParseInt(devIdStr, 10, 32)
		if err != nil {
			http.Error(w, "ID de desarrollador inválido", http.StatusBadRequest)
			return
		}

		err = q.DeleteDesarrolladorEnProyecto(r.Context(), db.DeleteDesarrolladorEnProyectoParams{
			IDDesarrollador: int32(idDesarrollador),
			IDProyecto:      int32(idProyecto),
		})
		if err != nil {
			fail500(w, err)
			return
		}

		filas, err := cargarProyectosFilas(r.Context(), q)
		if err != nil {
			fail500(w, err)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := views.ProyectoList(filas).Render(r.Context(), w); err != nil {
			fail500(w, err)
		}
	}
}
