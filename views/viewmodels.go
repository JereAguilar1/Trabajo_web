package views

import (
	"fmt"
	"time"

	db "tp2/db/sqlc"
)

// ProyectoRow es un proyecto listo para mostrar: con el nombre de su cliente.
type ProyectoRow struct {
	Proyecto db.Proyecto
	Cliente  string
}

// PageData agrupa todo lo que necesita la página principal.
type PageData struct {
	Clientes        []db.Cliente
	Desarrolladores []db.Desarrolladore
	Proyectos       []ProyectoRow
	Error           string // mensaje ya resuelto por el servidor (nunca texto del usuario)
}

// Estados válidos de un proyecto (valor guardado -> etiqueta mostrada).
var Estados = []struct{ Valor, Etiqueta string }{
	{"nuevo", "Nuevo"},
	{"en_progreso", "En progreso"},
	{"completado", "Completado"},
}

func estadoLabel(valor string) string {
	for _, e := range Estados {
		if e.Valor == valor {
			return e.Etiqueta
		}
	}
	return valor
}

func money(v float64) string { return fmt.Sprintf("%.2f", v) }

func fecha(t time.Time) string { return t.Format("2006-01-02") }
