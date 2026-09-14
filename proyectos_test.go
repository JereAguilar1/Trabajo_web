package main

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	db "tp2/db/sqlc"
)

func TestProyecto_CRUD(t *testing.T) {
	queries := getTestQueries(t)
	ctx := context.Background()

	// Un proyecto necesita un cliente ya existente (clave foránea id_cliente).
	cliente, err := queries.CreateCliente(ctx, "Cliente Para Proyecto")
	if err != nil {
		t.Fatalf("Error al crear cliente para proyecto: %v", err)
	}

	var createdProyecto db.Proyecto
	tiempoEstimado := time.Now().AddDate(0, 3, 0).Truncate(24 * time.Hour)

	t.Run("Crear Proyecto", func(t *testing.T) {
		params := db.CreateProyectoParams{
			Nombre:         "Sistema de Pagos",
			Descripcion:    "Plataforma de procesamiento de pagos online",
			Estado:         "Planificación",
			Gastos:         2500.50,
			Presupuesto:    10000.00,
			TiempoEstimado: tiempoEstimado,
			IDCliente:      cliente.IDCliente,
		}

		proyecto, err := queries.CreateProyecto(ctx, params)
		if err != nil {
			t.Fatalf("Error al crear proyecto: %v", err)
		}
		if proyecto.ID <= 0 {
			t.Fatalf("ID de proyecto inválido: %d", proyecto.ID)
		}
		if proyecto.Nombre != params.Nombre {
			t.Errorf("Nombre incorrecto: esperado %q, obtenido %q", params.Nombre, proyecto.Nombre)
		}
		if proyecto.IDCliente != cliente.IDCliente {
			t.Errorf("IDCliente incorrecto: esperado %d, obtenido %d", cliente.IDCliente, proyecto.IDCliente)
		}
		createdProyecto = proyecto
	})

	t.Run("Leer Proyecto por ID", func(t *testing.T) {
		proyecto, err := queries.GetProyecto(ctx, createdProyecto.ID)
		if err != nil {
			t.Fatalf("Error al obtener proyecto: %v", err)
		}
		if proyecto.ID != createdProyecto.ID {
			t.Errorf("ID mismatch: %d vs %d", proyecto.ID, createdProyecto.ID)
		}
		if proyecto.Nombre != createdProyecto.Nombre {
			t.Errorf("Nombre mismatch: %q vs %q", proyecto.Nombre, createdProyecto.Nombre)
		}
	})

	t.Run("Actualizar Proyecto", func(t *testing.T) {
		params := db.UpdateProyectoParams{
			ID:             createdProyecto.ID,
			Nombre:         "Sistema de Pagos y Facturación",
			Descripcion:    "Módulo integrado de pagos y facturación",
			Estado:         "En Desarrollo",
			Gastos:         4000.00,
			Presupuesto:    15000.00,
			TiempoEstimado: tiempoEstimado.AddDate(0, 1, 0),
			IDCliente:      cliente.IDCliente,
		}

		err := queries.UpdateProyecto(ctx, params)
		if err != nil {
			t.Fatalf("Error al actualizar proyecto: %v", err)
		}

		updated, err := queries.GetProyecto(ctx, createdProyecto.ID)
		if err != nil {
			t.Fatalf("Error al obtener proyecto actualizado: %v", err)
		}
		if updated.Nombre != params.Nombre {
			t.Errorf("Nombre no actualizado: esperado %q, obtenido %q", params.Nombre, updated.Nombre)
		}
		if updated.Estado != params.Estado {
			t.Errorf("Estado no actualizado: esperado %q, obtenido %q", params.Estado, updated.Estado)
		}
		createdProyecto = updated
	})

	t.Run("Listar Proyectos", func(t *testing.T) {
		proyectos, err := queries.ListProyectos(ctx)
		if err != nil {
			t.Fatalf("Error al listar proyectos: %v", err)
		}
		encontrado := false
		for _, p := range proyectos {
			if p.ID == createdProyecto.ID {
				encontrado = true
				if p.Nombre != createdProyecto.Nombre {
					t.Errorf("Nombre en lista mismatch: %q vs %q", p.Nombre, createdProyecto.Nombre)
				}
				break
			}
		}
		if !encontrado {
			t.Errorf("Proyecto %d no fue encontrado en el listado", createdProyecto.ID)
		}
	})

	t.Run("Eliminar Proyecto", func(t *testing.T) {
		err := queries.DeleteProyecto(ctx, createdProyecto.ID)
		if err != nil {
			t.Fatalf("Error al eliminar proyecto: %v", err)
		}

		_, err = queries.GetProyecto(ctx, createdProyecto.ID)
		if err == nil {
			t.Fatalf("Se esperaba error al consultar proyecto eliminado")
		}
		if !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("Se esperaba sql.ErrNoRows, obtenido: %v", err)
		}
	})
}
