-- name: CreateUser :one
INSERT INTO users (
    id,
    username,
    hashed_password,
    created_at,
    updated_at
) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY id
LIMIT $1;

-- name: UpdateUserPassword :exec
UPDATE users SET hashed_password = $1
WHERE id = $2;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;