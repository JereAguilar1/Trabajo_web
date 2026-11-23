package web

import (
	"context"
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
		proyectosFilas, err := cargarProyectosFilas(ctx, q)
		if err != nil {
			fail500(w, err)
			return
		}

		data := views.PageData{
			Clientes:        clientes,
			Desarrolladores: desarrolladores,
			Proyectos:       proyectosFilas,
			Error:           mensajes[r.URL.Query().Get("error")], // "" si el código no existe
		}

		if err := views.IndexPage("Sistema de Gestión de Proyectos", data).Render(ctx, w); err != nil {
			fail500(w, err)
		}
	}
}

func fail500(w http.ResponseWriter, err error) {
	log.Printf("web: %v", err)
	http.Error(w, "internal server error", http.StatusInternalServerError)
}

// cargarProyectosFilas recupera proyectos, clientes, desarrolladores y sus asignaciones,
// y retorna las filas listas para renderizar en views.ProyectoList.
func cargarProyectosFilas(ctx context.Context, q *db.Queries) ([]views.ProyectoRow, error) {
	proyectos, err := q.ListProyectos(ctx)
	if err != nil {
		return nil, err
	}
	clientes, err := q.ListClientes(ctx)
	if err != nil {
		return nil, err
	}
	desarrolladores, err := q.ListDesarrolladores(ctx)
	if err != nil {
		return nil, err
	}
	asignaciones, err := q.ListDesarrolladoresEnProyectos(ctx)
	if err != nil {
		return nil, err
	}

	return armarFilas(proyectos, clientes, desarrolladores, asignaciones), nil
}

// armarFilas combina cada proyecto con su cliente y sus desarrolladores asignados y disponibles.
func armarFilas(
	proyectos []db.Proyecto,
	clientes []db.Cliente,
	desarrolladores []db.Desarrolladore,
	asignaciones []db.DesarrolladoresEnProyecto,
) []views.ProyectoRow {
	nombreCliente := make(map[int32]string, len(clientes))
	for _, c := range clientes {
		nombreCliente[c.IDCliente] = c.Nombre
	}

	devsByID := make(map[int32]db.Desarrolladore, len(desarrolladores))
	for _, d := range desarrolladores {
		devsByID[d.IDDesarrollador] = d
	}

	devsPorProyecto := make(map[int32][]int32)
	for _, a := range asignaciones {
		devsPorProyecto[a.IDProyecto] = append(devsPorProyecto[a.IDProyecto], a.IDDesarrollador)
	}

	filas := make([]views.ProyectoRow, 0, len(proyectos))
	for _, p := range proyectos {
		var asignados []db.Desarrolladore
		asignadosSet := make(map[int32]bool)
		for _, devID := range devsPorProyecto[p.ID] {
			if d, ok := devsByID[devID]; ok {
				asignados = append(asignados, d)
				asignadosSet[devID] = true
			}
		}

		var disponibles []db.Desarrolladore
		for _, d := range desarrolladores {
			if !asignadosSet[d.IDDesarrollador] {
				disponibles = append(disponibles, d)
			}
		}

		filas = append(filas, views.ProyectoRow{
			Proyecto:        p,
			Cliente:         nombreCliente[p.IDCliente],
			Desarrolladores: asignados,
			Disponibles:     disponibles,
		})
	}
	return filas
}
