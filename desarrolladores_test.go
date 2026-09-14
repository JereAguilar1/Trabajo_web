package main

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	db "tp2/db/sqlc"
)

func TestDesarrollador_CRUD(t *testing.T) {
	queries := getTestQueries(t)
	ctx := context.Background()

	var createdDesarrollador db.Desarrolladore

	t.Run("Crear Desarrollador", func(t *testing.T) {
		nombre := "Juan"
		apellido := "Dev"
		desarrollador, err := queries.CreateDesarrollador(ctx, db.CreateDesarrolladorParams{
			Nombre:   nombre,
			Apellido: apellido,
		})
		if err != nil {
			t.Fatalf("Error al crear desarrollador: %v", err)
		}
		if desarrollador.IDDesarrollador <= 0 {
			t.Fatalf("ID de desarrollador invalido: %d", desarrollador.IDDesarrollador)
		}
		if desarrollador.Nombre != nombre {
			t.Errorf("Nombre incorrecto: esperado %q, obtenido %q", nombre, desarrollador.Nombre)
		}
		if desarrollador.Apellido != apellido {
			t.Errorf("Apellido incorrecto: esperado %q, obtenido %q", apellido, desarrollador.Apellido)
		}
		createdDesarrollador = desarrollador
	})

	t.Run("Leer Desarrollador por ID", func(t *testing.T) {
		desarrollador, err := queries.GetDesarrollador(ctx, createdDesarrollador.IDDesarrollador)
		if err != nil {
			t.Fatalf("Error al obtener desarrollador: %v", err)
		}
		if desarrollador.IDDesarrollador != createdDesarrollador.IDDesarrollador {
			t.Errorf("ID mismatch: %d vs %d", desarrollador.IDDesarrollador, createdDesarrollador.IDDesarrollador)
		}
		if desarrollador.Nombre != createdDesarrollador.Nombre {
			t.Errorf("Nombre mismatch: %q vs %q", desarrollador.Nombre, createdDesarrollador.Nombre)
		}
		if desarrollador.Apellido != createdDesarrollador.Apellido {
			t.Errorf("Apellido mismatch: %q vs %q", desarrollador.Apellido, createdDesarrollador.Apellido)
		}
	})

	t.Run("Actualizar Desarrollador", func(t *testing.T) {
		params := db.UpdateDesarrolladorParams{
			IDDesarrollador: createdDesarrollador.IDDesarrollador,
			Nombre:          "Juan Carlos",
			Apellido:        "Devoto",
		}
		err := queries.UpdateDesarrollador(ctx, params)
		if err != nil {
			t.Fatalf("Error al actualizar desarrollador: %v", err)
		}

		updated, err := queries.GetDesarrollador(ctx, createdDesarrollador.IDDesarrollador)
		if err != nil {
			t.Fatalf("Error al obtener desarrollador actualizado: %v", err)
		}
		if updated.Nombre != params.Nombre {
			t.Errorf("Nombre no actualizado: esperado %q, obtenido %q", params.Nombre, updated.Nombre)
		}
		if updated.Apellido != params.Apellido {
			t.Errorf("Apellido no actualizado: esperado %q, obtenido %q", params.Apellido, updated.Apellido)
		}
		createdDesarrollador = updated
	})

	t.Run("Listar Desarrolladores", func(t *testing.T) {
		desarrolladores, err := queries.ListDesarrolladores(ctx)
		if err != nil {
			t.Fatalf("Error al listar desarrolladores: %v", err)
		}
		encontrado := false
		for _, d := range desarrolladores {
			if d.IDDesarrollador == createdDesarrollador.IDDesarrollador {
				encontrado = true
				if d.Nombre != createdDesarrollador.Nombre {
					t.Errorf("Nombre en lista mismatch: %q vs %q", d.Nombre, createdDesarrollador.Nombre)
				}
				break
			}
		}
		if !encontrado {
			t.Errorf("Desarrollador %d no fue encontrado en el listado", createdDesarrollador.IDDesarrollador)
		}
	})

	t.Run("Eliminar Desarrollador", func(t *testing.T) {
		err := queries.DeleteDesarrollador(ctx, createdDesarrollador.IDDesarrollador)
		if err != nil {
			t.Fatalf("Error al eliminar desarrollador: %v", err)
		}

		_, err = queries.GetDesarrollador(ctx, createdDesarrollador.IDDesarrollador)
		if err == nil {
			t.Fatalf("Se esperaba error al consultar desarrollador eliminado")
		}
		if !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("Se esperaba sql.ErrNoRows, obtenido: %v", err)
		}
	})
}
