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
			redirectError(w, r, "invalido")
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
			redirectError(w, r, "invalido")
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
		switch {
		case err == nil:
			redirectOK(w, r)
		case pgCode(err) == fkViolation: // el cliente no existe
			redirectError(w, r, "no_existe")
		default:
			fail(w, r, err)
		}
	}
}
