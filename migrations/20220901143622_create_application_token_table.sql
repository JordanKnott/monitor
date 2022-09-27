-- +goose Up
-- +goose StatementBegin
INSERT INTO user_account (created_at, email, google_id, nicename, photo) VALUES (NOW(), 'jordan@drivendigital.us', '', '', '') ON CONFLICT DO NOTHING;

CREATE TABLE app_token (
    token_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at timestamptz NOT NULL,
    user_id uuid NOT NULL REFERENCES user_account(user_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
