-- name: ListTasks :many
SELECT * FROM tasks;

-- name: CreateTask :exec
INSERT INTO tasks (
	id, name, description, parent_id, last_done_on, reoccurrence, type, index
) VALUES (
	$1, $2, $3, $4, $5, $6, $7, $8
);

-- name: CompleteTask :exec
UPDATE tasks
SET last_done_on = NOW()
WHERE id = $1;
