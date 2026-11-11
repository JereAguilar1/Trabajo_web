package web

import (
	"log"
	"net/http"

	db "tp2/db/sqlc"
	"tp2/views"
)

// GET /
// Carga clientes, desarrolladores y proyectos, y renderiza la
// página completa.
func Index(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		clientes, err := q.ListClientes(ctx)
		if err != nil {
			fail500(w, err)
			return
		}
		desarrolladores, err := q.ListDesarrolladores(ctx)
		if err != nil {
			fail500(w, err)
			return
		}
		proyectos, err := q.ListProyectos(ctx)
		if err != nil {
			fail500(w, err)
			return
		}

		data := views.PageData{
			Clientes:        clientes,
			Desarrolladores: desarrolladores,
			Proyectos:       armarFilas(proyectos, clientes),
			Error:           mensajes[r.URL.Query().Get("error")], // "" si el código no existe
		}

		if err := views.IndexPage("Sistema de Gestión de Proyectos", data).Render(ctx, w); err != nil {
			fail500(w, err)
		}
	}
}

func fail500(w http.ResponseWriter, err error) {
	log.Printf("GET /: %v", err)
	http.Error(w, "internal server error", http.StatusInternalServerError)
}

// armarFilas combina cada proyecto con el nombre de su cliente.
func armarFilas(proyectos []db.Proyecto, clientes []db.Cliente) []views.ProyectoRow {
	nombreCliente := make(map[int32]string, len(clientes))
	for _, c := range clientes {
		nombreCliente[c.IDCliente] = c.Nombre
	}

	filas := make([]views.ProyectoRow, 0, len(proyectos))
	for _, p := range proyectos {
		filas = append(filas, views.ProyectoRow{Proyecto: p, Cliente: nombreCliente[p.IDCliente]})
	}
	return filas
}
