-- name: GetInstallByID :one
SELECT * FROM install WHERE install_id = $1;

-- name: GetInstallByIdentifier :one
SELECT * FROM install WHERE ssh_username = $1;

-- name: GetAllInstallIdentifiers :many
SELECT install_id, ssh_username FROM install;

-- name: GetAllInstalls :many
SELECT * FROM install;

-- name: CreateInstall :one
INSERT INTO install (nicename, ssh_username, ssh_port, ssh_hostname) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: CreateSnapshot :one
INSERT INTO install_snapshot (created_at, created_by, install_id) VALUES ($1, $2, $3) RETURNING *;

-- name: CreateSnapshotHash :one
INSERT INTO install_snapshot_hash_tree (snapshot_id, filepath, filehash) VALUES ($1, $2, $3) RETURNING *;

-- name: GetSnapshotHashes :many
SELECT * FROM install_snapshot_hash_tree WHERE snapshot_id = $1;

-- name: GetLatestSnapshotForInstallID :one
SELECT * FROM install_snapshot WHERE install_id = $1 ORDER BY created_at DESC LIMIT 1;

-- name: CreateReport :one
INSERT INTO install_report (install_id, snapshot_id, created_at) VALUES ($1, $2, $3) RETURNING *;

-- name: CreateReportEntry :one
INSERT INTO install_report_entry (report_id, filepath, snapshot_hash, current_hash) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetReportEntriesForReport :many
SELECT * FROM install_report_entry WHERE report_id = $1;

-- name: CreateSnapshotPlugin :one
INSERT INTO install_snapshot_plugins (snapshot_id, nicename, slug, version, status) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetSnapshotPluginsForReport :many
SELECT * FROM install_snapshot_plugins WHERE snapshot_id = $1;

-- name: CreateReportPlugin :one
INSERT INTO install_report_entry_plugin (report_id, slug, issue) VALUES ($1, $2, $3) RETURNING *;

-- name: GetInstallIDForReportID :one
SELECT install_id FROM install_report WHERE report_id = $1;