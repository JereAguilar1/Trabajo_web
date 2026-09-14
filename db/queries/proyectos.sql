-- name: CreateProyecto :one
INSERT INTO proyectos (
    nombre,
    descripcion,
    estado,
    gastos,
    presupuesto,
    tiempo_estimado,
    id_cliente
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetProyecto :one
SELECT * FROM proyectos
WHERE id = $1;

-- name: ListProyectos :many
SELECT * FROM proyectos
ORDER BY id;

-- name: UpdateProyecto :exec
UPDATE proyectos
SET
    nombre = $2,
    descripcion = $3,
    estado = $4,
    gastos = $5,
    presupuesto = $6,
    tiempo_estimado = $7,
    id_cliente = $8
WHERE id = $1;

-- name: DeleteProyecto :exec
DELETE FROM proyectos
WHERE id = $1;
