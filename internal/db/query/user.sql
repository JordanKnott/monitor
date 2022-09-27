-- name: CreateAccessToken :exec
INSERT INTO access_token (token, user_id, created_at, expires_at) VALUES ($1, $2, $3, $4);

-- name: GetAccessToken :one
SELECT * FROM access_token WHERE token = $1;

-- name: GetUserIDForToken :one
SELECT user_id FROM access_token WHERE token = $1;

-- name: GetUserByEmail :one
SELECT * FROM user_account WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM user_account WHERE user_id = $1;

-- name: GetUserByGoogleId :one
SELECT * FROM user_account WHERE google_id = $1;

-- name: CreateUserAccount :one
INSERT INTO user_account (created_at, google_id, nicename, email, photo) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: CreateAppToken :one
INSERT INTO app_token (created_at, user_id) VALUES ($1, $2) RETURNING *;

-- name: GetAppToken :one
SELECT * FROM app_token WHERE token_id = $1;