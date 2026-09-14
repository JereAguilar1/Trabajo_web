-- name: CreateCliente :one
INSERT INTO clientes (nombre) 
VALUES ($1)
RETURNING *;

-- name: GetCliente :one
SELECT * FROM clientes
WHERE id_cliente = $1;

-- name: ListClientes :many
SELECT * FROM clientes
ORDER BY id_cliente;

-- name: UpdateCliente :exec
UPDATE clientes
SET nombre = $2
WHERE id_cliente = $1;

-- name: DeleteCliente :exec
DELETE FROM clientes
WHERE id_cliente = $1;
