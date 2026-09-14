package main

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	db "tp2/db/sqlc"
)

func TestDesarrolladoresEnProyectos_CRUD(t *testing.T) {
	queries := getTestQueries(t)
	ctx := context.Background()

	// Entidades necesarias para poder asignar: un cliente, un proyecto y dos desarrolladores.
	cliente, err := queries.CreateCliente(ctx, "Cliente Relaciones")
	if err != nil {
		t.Fatalf("Error al crear cliente: %v", err)
	}

	proyecto, err := queries.CreateProyecto(ctx, db.CreateProyectoParams{
		Nombre:         "Proyecto Cloud",
		Descripcion:    "Migración a la nube",
		Estado:         "Activo",
		Gastos:         100.0,
		Presupuesto:    1000.0,
		TiempoEstimado: time.Now().AddDate(0, 2, 0),
		IDCliente:      cliente.IDCliente,
	})
	if err != nil {
		t.Fatalf("Error al crear proyecto: %v", err)
	}

	dev1, err := queries.CreateDesarrollador(ctx, db.CreateDesarrolladorParams{
		Nombre:   "Ada",
		Apellido: "Lovelace",
	})
	if err != nil {
		t.Fatalf("Error al crear dev1: %v", err)
	}

	dev2, err := queries.CreateDesarrollador(ctx, db.CreateDesarrolladorParams{
		Nombre:   "Margaret",
		Apellido: "Hamilton",
	})
	if err != nil {
		t.Fatalf("Error al crear dev2: %v", err)
	}

	t.Run("Asignar Desarrolladores", func(t *testing.T) {
		asig1, err := queries.CreateDesarrolladorEnProyecto(ctx, db.CreateDesarrolladorEnProyectoParams{
			IDDesarrollador: dev1.IDDesarrollador,
			IDProyecto:      proyecto.ID,
		})
		if err != nil {
			t.Fatalf("Error al asignar dev1 al proyecto: %v", err)
		}
		if asig1.IDDesarrollador != dev1.IDDesarrollador || asig1.IDProyecto != proyecto.ID {
			t.Errorf("Datos de asignación erróneos: %+v", asig1)
		}

		_, err = queries.CreateDesarrolladorEnProyecto(ctx, db.CreateDesarrolladorEnProyectoParams{
			IDDesarrollador: dev2.IDDesarrollador,
			IDProyecto:      proyecto.ID,
		})
		if err != nil {
			t.Fatalf("Error al asignar dev2 al proyecto: %v", err)
		}
	})

	t.Run("Obtener Asignación", func(t *testing.T) {
		asig, err := queries.GetDesarrolladorEnProyecto(ctx, db.GetDesarrolladorEnProyectoParams{
			IDDesarrollador: dev1.IDDesarrollador,
			IDProyecto:      proyecto.ID,
		})
		if err != nil {
			t.Fatalf("Error al obtener asignación: %v", err)
		}
		if asig.IDDesarrollador != dev1.IDDesarrollador || asig.IDProyecto != proyecto.ID {
			t.Errorf("Asignación incorrecta: %+v", asig)
		}
	})

	t.Run("Listar Desarrolladores por Proyecto", func(t *testing.T) {
		devsEnProy, err := queries.ListDesarrolladoresByProyecto(ctx, proyecto.ID)
		if err != nil {
			t.Fatalf("Error al listar desarrolladores del proyecto: %v", err)
		}
		if len(devsEnProy) != 2 {
			t.Errorf("Se esperaban 2 desarrolladores asignados, obtenidos: %d", len(devsEnProy))
		}
	})

	t.Run("Listar Proyectos por Desarrollador", func(t *testing.T) {
		proysDeDev, err := queries.ListProyectosByDesarrollador(ctx, dev1.IDDesarrollador)
		if err != nil {
			t.Fatalf("Error al listar proyectos de dev1: %v", err)
		}
		if len(proysDeDev) != 1 {
			t.Errorf("Se esperaba 1 proyecto para dev1, obtenidos: %d", len(proysDeDev))
		}
		if proysDeDev[0].IDProyecto != proyecto.ID {
			t.Errorf("Proyecto incorrecto para dev1: %d vs %d", proysDeDev[0].IDProyecto, proyecto.ID)
		}
	})

	t.Run("Listar Todas las Asignaciones", func(t *testing.T) {
		todas, err := queries.ListDesarrolladoresEnProyectos(ctx)
		if err != nil {
			t.Fatalf("Error al listar todas las asignaciones: %v", err)
		}
		if len(todas) < 2 {
			t.Errorf("Se esperaban al menos 2 asignaciones globales, obtenidos: %d", len(todas))
		}
	})

	t.Run("Desasignar Desarrollador de Proyecto", func(t *testing.T) {
		err := queries.DeleteDesarrolladorEnProyecto(ctx, db.DeleteDesarrolladorEnProyectoParams{
			IDDesarrollador: dev1.IDDesarrollador,
			IDProyecto:      proyecto.ID,
		})
		if err != nil {
			t.Fatalf("Error al desasignar dev1: %v", err)
		}

		_, err = queries.GetDesarrolladorEnProyecto(ctx, db.GetDesarrolladorEnProyectoParams{
			IDDesarrollador: dev1.IDDesarrollador,
			IDProyecto:      proyecto.ID,
		})
		if err == nil {
			t.Fatalf("Se esperaba error al consultar asignación eliminada")
		}
		if !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("Se esperaba sql.ErrNoRows, obtenido: %v", err)
		}

		devsRestantes, err := queries.ListDesarrolladoresByProyecto(ctx, proyecto.ID)
		if err != nil {
			t.Fatalf("Error al listar desarrolladores restantes: %v", err)
		}
		if len(devsRestantes) != 1 {
			t.Errorf("Se esperaba 1 desarrollador restante, obtenidos: %d", len(devsRestantes))
		}
	})
}
