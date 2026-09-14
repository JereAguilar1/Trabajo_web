-- name: CreateDesarrollador :one
INSERT INTO desarrolladores (
    nombre,
    apellido
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetDesarrollador :one
SELECT * FROM desarrolladores
WHERE id_desarrollador = $1;

-- name: ListDesarrolladores :many
SELECT * FROM desarrolladores
ORDER BY id_desarrollador;

-- name: UpdateDesarrollador :exec
UPDATE desarrolladores
SET
    nombre = $2,
    apellido = $3
WHERE id_desarrollador = $1;

-- name: DeleteDesarrollador :exec
DELETE FROM desarrolladores
WHERE id_desarrollador = $1;
