-- +goose Up
-- +goose StatementBegin
CREATE TABLE install_metadata (
    metadata_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    install_id uuid NOT NULL REFERENCES install(install_id) ON DELETE CASCADE,
    is_multisite boolean NOT NULL DEFAULT false,
    primary_url text NOT NULL,
    core_version text NOT NULL,
);

CREATE TABLE install_user (
    user_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    install_id uuid NOT NULL REFERENCES install(install_id) ON DELETE CASCADE,
    wordpress_id text NOT NULL,
    nicename text NOT NULL,
    username text NOT NULL,
    email text NOT NULL,
    roles text NOT NULL
);

CREATE TABLE install_plugin (
    plugin_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    install_id uuid NOT NULL REFERENCES install(install_id) ON DELETE CASCADE,
    title text NOT NULL,
    name text NOT NULL,
    version text NOT NULL,
    status text NOT NULL
);

CREATE TABLE install_site (
    site_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    install_id uuid NOT NULL REFERENCES install(install_id) ON DELETE CASCADE,
    url text NOT NULL
);

CREATE TABLE install_site_user (
    user_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    site_id uuid NOT NULL REFERENCES install_site(site_id) ON DELETE CASCADE,
    nicename text NOT NULL,
    username text NOT NULL,
    email text NOT NULL,
    roles text NOT NULL
);

CREATE TABLE install_site_plugin (
    plugin_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    site_id uuid NOT NULL REFERENCES install_site(site_id) ON DELETE CASCADE,
    title text NOT NULL,
    name text NOT NULL,
    version text NOT NULL,
    status text NOT NULL
);

CREATE TABLE install_ssl_cert (
    cert_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    install_id uuid NOT NULL REFERENCES install(install_id) ON DELETE CASCADE,
    site_id uuid NOT NULL REFERENCES install_site(site_id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL,
    expires_at timestamptz NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
