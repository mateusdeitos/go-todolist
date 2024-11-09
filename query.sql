-- name: CreateTodo :one
INSERT INTO todos (
	name, complete
) VALUES (
	$1, $2
)
RETURNING *;

-- name: UpdateTodo :exec
UPDATE todos
SET name = $1, complete = $2
WHERE id = $3;

-- name: DeleteTodo :exec
DELETE FROM todos
WHERE id = $1;

-- name: GetTodo :one
SELECT * FROM todos
WHERE id = $1 LIMIT 1;

-- name: ListTodos :many
SELECT * FROM todos
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: CountTodos :one
SELECT COUNT(*) FROM todos;
