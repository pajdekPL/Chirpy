-- name: GetUserByID :one
SELECT id, created_at, updated_at, email, user_name, hashed_password, is_chirpy_red FROM users
WHERE id = $1;


-- name: ChangeUserName :one
UPDATE users
SET user_name = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ChangeUserData :one
UPDATE users
SET email = $2, hashed_password = $3, user_name = $4, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, hashed_password, is_chirpy_red FROM users
WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password, user_name)
VALUES (
    gen_random_uuid(), NOW(), NOW(), $1, $2, $3
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: UpgradeUserToRed :one
UPDATE users
SET is_chirpy_red = TRUE, updated_at = NOW()
WHERE id = $1
RETURNING *;