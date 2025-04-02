-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id, expiration_datetime)
VALUES (
    gen_random_uuid(), NOW(), NOW(), $1, $2, $3
)
RETURNING *;

-- name: GetChirps :many
SELECT c.id, c.created_at, c.updated_at, c.body, c.user_id, c.expiration_datetime, u.user_name as author_name
FROM chirps c
JOIN users u ON c.user_id = u.id
WHERE c.expiration_datetime > NOW()
ORDER BY c.created_at ASC;

-- name: GetExpiredChirps :many
SELECT c.id, c.created_at, c.updated_at, c.body, c.user_id, c.expiration_datetime, u.user_name as author_name
FROM chirps c
JOIN users u ON c.user_id = u.id
WHERE c.expiration_datetime < NOW()
ORDER BY c.created_at ASC;

-- name: GetChirp :one
SELECT id, created_at, updated_at, body, user_id, expiration_datetime FROM chirps
WHERE id = $1;

-- name: GetChirpsByUser :many
SELECT c.id, c.created_at, c.updated_at, c.body, c.user_id, c.expiration_datetime, u.user_name as author_name
FROM chirps c
JOIN users u ON c.user_id = u.id
WHERE c.user_id = $1 AND c.expiration_datetime > NOW()
ORDER BY c.created_at ASC;

-- name: GetExpiredChirpsByUser :many
SELECT c.id, c.created_at, c.updated_at, c.body, c.user_id, c.expiration_datetime, u.user_name as author_name
FROM chirps c
JOIN users u ON c.user_id = u.id
WHERE c.user_id = $1 AND c.expiration_datetime < NOW()
ORDER BY c.created_at ASC;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1;