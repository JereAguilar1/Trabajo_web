-- name: CreateDesarrolladorEnProyecto :one
INSERT INTO desarrolladores_en_proyectos (
    id_desarrollador,
    id_proyecto
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetDesarrolladorEnProyecto :one
SELECT * FROM desarrolladores_en_proyectos
WHERE id_desarrollador = $1 AND id_proyecto = $2;

-- name: ListDesarrolladoresEnProyectos :many
SELECT * FROM desarrolladores_en_proyectos
ORDER BY id_proyecto, id_desarrollador;

-- name: ListDesarrolladoresByProyecto :many
SELECT * FROM desarrolladores_en_proyectos
WHERE id_proyecto = $1;

-- name: ListProyectosByDesarrollador :many
SELECT * FROM desarrolladores_en_proyectos
WHERE id_desarrollador = $1;

-- name: DeleteDesarrolladorEnProyecto :exec
DELETE FROM desarrolladores_en_proyectos
WHERE id_desarrollador = $1 AND id_proyecto = $2;
