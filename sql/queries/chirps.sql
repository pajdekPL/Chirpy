-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(), NOW(), NOW(), $1, $2
)
RETURNING *;

-- name: GetChirps :many
SELECT c.id, c.created_at, c.updated_at, c.body, c.user_id, u.user_name as author_name
FROM chirps c
JOIN users u ON c.user_id = u.id
ORDER BY c.created_at ASC;

-- name: GetChirp :one
SELECT id, created_at, updated_at, body, user_id FROM chirps
WHERE id = $1;

-- name: GetChirpsByUser :many
SELECT c.id, c.created_at, c.updated_at, c.body, c.user_id, u.user_name as author_name
FROM chirps c
JOIN users u ON c.user_id = u.id
WHERE c.user_id = $1
ORDER BY c.created_at ASC;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1;