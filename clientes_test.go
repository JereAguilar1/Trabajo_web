package main

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	db "tp2/db/sqlc"
)

func TestCliente_CRUD(t *testing.T) {
	queries := getTestQueries(t)
	ctx := context.Background()

	var createdCliente db.Cliente

	t.Run("Crear Cliente", func(t *testing.T) {
		nombre := "Empresa de Software ABC"
		cliente, err := queries.CreateCliente(ctx, nombre)
		if err != nil {
			t.Fatalf("Error al crear cliente: %v", err)
		}
		if cliente.IDCliente <= 0 {
			t.Fatalf("ID de cliente inválido: %d", cliente.IDCliente)
		}
		if cliente.Nombre != nombre {
			t.Errorf("Nombre incorrecto: esperado %q, obtenido %q", nombre, cliente.Nombre)
		}
		createdCliente = cliente
	})

	t.Run("Leer Cliente por ID", func(t *testing.T) {
		cliente, err := queries.GetCliente(ctx, createdCliente.IDCliente)
		if err != nil {
			t.Fatalf("Error al obtener cliente: %v", err)
		}
		if cliente.IDCliente != createdCliente.IDCliente {
			t.Errorf("ID mismatch: %d vs %d", cliente.IDCliente, createdCliente.IDCliente)
		}
		if cliente.Nombre != createdCliente.Nombre {
			t.Errorf("Nombre mismatch: %q vs %q", cliente.Nombre, createdCliente.Nombre)
		}
	})

	t.Run("Actualizar Cliente", func(t *testing.T) {
		params := db.UpdateClienteParams{
			IDCliente: createdCliente.IDCliente,
			Nombre:    "Empresa ABC Internacional",
		}
		err := queries.UpdateCliente(ctx, params)
		if err != nil {
			t.Fatalf("Error al actualizar cliente: %v", err)
		}

		updated, err := queries.GetCliente(ctx, createdCliente.IDCliente)
		if err != nil {
			t.Fatalf("Error al obtener cliente actualizado: %v", err)
		}
		if updated.Nombre != params.Nombre {
			t.Errorf("Nombre no actualizado: esperado %q, obtenido %q", params.Nombre, updated.Nombre)
		}
		createdCliente = updated
	})

	t.Run("Listar Clientes", func(t *testing.T) {
		clientes, err := queries.ListClientes(ctx)
		if err != nil {
			t.Fatalf("Error al listar clientes: %v", err)
		}
		encontrado := false
		for _, c := range clientes {
			if c.IDCliente == createdCliente.IDCliente {
				encontrado = true
				if c.Nombre != createdCliente.Nombre {
					t.Errorf("Nombre en lista mismatch: %q vs %q", c.Nombre, createdCliente.Nombre)
				}
				break
			}
		}
		if !encontrado {
			t.Errorf("Cliente %d no fue encontrado en el listado", createdCliente.IDCliente)
		}
	})

	t.Run("Eliminar Cliente", func(t *testing.T) {
		err := queries.DeleteCliente(ctx, createdCliente.IDCliente)
		if err != nil {
			t.Fatalf("Error al eliminar cliente: %v", err)
		}

		_, err = queries.GetCliente(ctx, createdCliente.IDCliente)
		if err == nil {
			t.Fatalf("Se esperaba error al consultar cliente eliminado")
		}
		if !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("Se esperaba sql.ErrNoRows, obtenido: %v", err)
		}
	})
}
