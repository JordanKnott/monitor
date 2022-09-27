-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE user_account (
    user_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at timestamptz NOT NULL,
    google_id text NOT NULL UNIQUE,
    nicename text NOT NULL,
    email text NOT NULL UNIQUE,
    photo text NOT NULL
);

CREATE TABLE access_token (
    token text NOT NULL,
    user_id uuid NOT NULL REFERENCES user_account (user_id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL,
    expires_at timestamptz NOT NULL
);

CREATE TABLE install (
    install_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    nicename text NOT NULL,
    ssh_username text NOT NULL UNIQUE,
    ssh_port int NOT NULL,
    ssh_hostname text NOT NULL,
    is_multisite boolean DEFAULT NULL
);


CREATE TABLE install_snapshot (
    snapshot_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    install_id uuid NOT NULL REFERENCES install(install_id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL,
    created_by uuid NOT NULL REFERENCES user_account(user_id) ON DELETE CASCADE
);


CREATE TABLE install_snapshot_plugins (
    plugin_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    snapshot_id uuid NOT NULL REFERENCES install_snapshot(snapshot_id) ON DELETE CASCADE,
    nicename text NOT NULL,
    slug text NOT NULL,
    version text NOT NULL,
    status text NOT NULL
);

CREATE TABLE install_snapshot_hash_tree (
    hash_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    snapshot_id uuid NOT NULL REFERENCES install_snapshot(snapshot_id) ON DELETE CASCADE,
    filepath text NOT NULL,
    filehash text NOT NULL
);


CREATE TABLE install_report (
    report_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    install_id uuid NOT NULL REFERENCES install(install_id) ON DELETE CASCADE,
    snapshot_id uuid NOT NULL REFERENCES install_snapshot(snapshot_id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL
);

CREATE TABLE install_report_entry (
    entry_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    report_id uuid NOT NULL REFERENCES install_report(report_id) ON DELETE CASCADE,
    filepath text NOT NULL,
    snapshot_hash text NOT NULL,
    current_hash text NOT NULL
);

CREATE TABLE install_report_entry_plugin (
    entry_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    report_id uuid NOT NULL REFERENCES install_report(report_id) ON DELETE CASCADE,
    slug  text NOT NULL,
    issue json NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
