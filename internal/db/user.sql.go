// Code generated by sqlc. DO NOT EDIT.
// source: user.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createAccessToken = `-- name: CreateAccessToken :exec
INSERT INTO access_token (token, user_id, created_at, expires_at) VALUES ($1, $2, $3, $4)
`

type CreateAccessTokenParams struct {
	Token     string    `json:"token"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (q *Queries) CreateAccessToken(ctx context.Context, arg CreateAccessTokenParams) error {
	_, err := q.db.ExecContext(ctx, createAccessToken,
		arg.Token,
		arg.UserID,
		arg.CreatedAt,
		arg.ExpiresAt,
	)
	return err
}

const createAppToken = `-- name: CreateAppToken :one
INSERT INTO app_token (created_at, user_id) VALUES ($1, $2) RETURNING token_id, created_at, user_id
`

type CreateAppTokenParams struct {
	CreatedAt time.Time `json:"created_at"`
	UserID    uuid.UUID `json:"user_id"`
}

func (q *Queries) CreateAppToken(ctx context.Context, arg CreateAppTokenParams) (AppToken, error) {
	row := q.db.QueryRowContext(ctx, createAppToken, arg.CreatedAt, arg.UserID)
	var i AppToken
	err := row.Scan(&i.TokenID, &i.CreatedAt, &i.UserID)
	return i, err
}

const createUserAccount = `-- name: CreateUserAccount :one
INSERT INTO user_account (created_at, google_id, nicename, email, photo) VALUES ($1, $2, $3, $4, $5) RETURNING user_id, created_at, google_id, nicename, email, photo
`

type CreateUserAccountParams struct {
	CreatedAt time.Time `json:"created_at"`
	GoogleID  string    `json:"google_id"`
	Nicename  string    `json:"nicename"`
	Email     string    `json:"email"`
	Photo     string    `json:"photo"`
}

func (q *Queries) CreateUserAccount(ctx context.Context, arg CreateUserAccountParams) (UserAccount, error) {
	row := q.db.QueryRowContext(ctx, createUserAccount,
		arg.CreatedAt,
		arg.GoogleID,
		arg.Nicename,
		arg.Email,
		arg.Photo,
	)
	var i UserAccount
	err := row.Scan(
		&i.UserID,
		&i.CreatedAt,
		&i.GoogleID,
		&i.Nicename,
		&i.Email,
		&i.Photo,
	)
	return i, err
}

const getAccessToken = `-- name: GetAccessToken :one
SELECT token, user_id, created_at, expires_at FROM access_token WHERE token = $1
`

func (q *Queries) GetAccessToken(ctx context.Context, token string) (AccessToken, error) {
	row := q.db.QueryRowContext(ctx, getAccessToken, token)
	var i AccessToken
	err := row.Scan(
		&i.Token,
		&i.UserID,
		&i.CreatedAt,
		&i.ExpiresAt,
	)
	return i, err
}

const getAppToken = `-- name: GetAppToken :one
SELECT token_id, created_at, user_id FROM app_token WHERE token_id = $1
`

func (q *Queries) GetAppToken(ctx context.Context, tokenID uuid.UUID) (AppToken, error) {
	row := q.db.QueryRowContext(ctx, getAppToken, tokenID)
	var i AppToken
	err := row.Scan(&i.TokenID, &i.CreatedAt, &i.UserID)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT user_id, created_at, google_id, nicename, email, photo FROM user_account WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (UserAccount, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i UserAccount
	err := row.Scan(
		&i.UserID,
		&i.CreatedAt,
		&i.GoogleID,
		&i.Nicename,
		&i.Email,
		&i.Photo,
	)
	return i, err
}

const getUserByGoogleId = `-- name: GetUserByGoogleId :one
SELECT user_id, created_at, google_id, nicename, email, photo FROM user_account WHERE google_id = $1
`

func (q *Queries) GetUserByGoogleId(ctx context.Context, googleID string) (UserAccount, error) {
	row := q.db.QueryRowContext(ctx, getUserByGoogleId, googleID)
	var i UserAccount
	err := row.Scan(
		&i.UserID,
		&i.CreatedAt,
		&i.GoogleID,
		&i.Nicename,
		&i.Email,
		&i.Photo,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT user_id, created_at, google_id, nicename, email, photo FROM user_account WHERE user_id = $1
`

func (q *Queries) GetUserByID(ctx context.Context, userID uuid.UUID) (UserAccount, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, userID)
	var i UserAccount
	err := row.Scan(
		&i.UserID,
		&i.CreatedAt,
		&i.GoogleID,
		&i.Nicename,
		&i.Email,
		&i.Photo,
	)
	return i, err
}

const getUserIDForToken = `-- name: GetUserIDForToken :one
SELECT user_id FROM access_token WHERE token = $1
`

func (q *Queries) GetUserIDForToken(ctx context.Context, token string) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, getUserIDForToken, token)
	var user_id uuid.UUID
	err := row.Scan(&user_id)
	return user_id, err
}
