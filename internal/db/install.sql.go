// Code generated by sqlc. DO NOT EDIT.
// source: install.sql

package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const createInstall = `-- name: CreateInstall :one
INSERT INTO install (nicename, ssh_username, ssh_port, ssh_hostname) VALUES ($1, $2, $3, $4) RETURNING install_id, nicename, ssh_username, ssh_port, ssh_hostname, is_multisite
`

type CreateInstallParams struct {
	Nicename    string `json:"nicename"`
	SshUsername string `json:"ssh_username"`
	SshPort     int32  `json:"ssh_port"`
	SshHostname string `json:"ssh_hostname"`
}

func (q *Queries) CreateInstall(ctx context.Context, arg CreateInstallParams) (Install, error) {
	row := q.db.QueryRowContext(ctx, createInstall,
		arg.Nicename,
		arg.SshUsername,
		arg.SshPort,
		arg.SshHostname,
	)
	var i Install
	err := row.Scan(
		&i.InstallID,
		&i.Nicename,
		&i.SshUsername,
		&i.SshPort,
		&i.SshHostname,
		&i.IsMultisite,
	)
	return i, err
}

const createReport = `-- name: CreateReport :one
INSERT INTO install_report (install_id, snapshot_id, created_at) VALUES ($1, $2, $3) RETURNING report_id, install_id, snapshot_id, created_at
`

type CreateReportParams struct {
	InstallID  uuid.UUID `json:"install_id"`
	SnapshotID uuid.UUID `json:"snapshot_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func (q *Queries) CreateReport(ctx context.Context, arg CreateReportParams) (InstallReport, error) {
	row := q.db.QueryRowContext(ctx, createReport, arg.InstallID, arg.SnapshotID, arg.CreatedAt)
	var i InstallReport
	err := row.Scan(
		&i.ReportID,
		&i.InstallID,
		&i.SnapshotID,
		&i.CreatedAt,
	)
	return i, err
}

const createReportEntry = `-- name: CreateReportEntry :one
INSERT INTO install_report_entry (report_id, filepath, snapshot_hash, current_hash) VALUES ($1, $2, $3, $4) RETURNING entry_id, report_id, filepath, snapshot_hash, current_hash
`

type CreateReportEntryParams struct {
	ReportID     uuid.UUID `json:"report_id"`
	Filepath     string    `json:"filepath"`
	SnapshotHash string    `json:"snapshot_hash"`
	CurrentHash  string    `json:"current_hash"`
}

func (q *Queries) CreateReportEntry(ctx context.Context, arg CreateReportEntryParams) (InstallReportEntry, error) {
	row := q.db.QueryRowContext(ctx, createReportEntry,
		arg.ReportID,
		arg.Filepath,
		arg.SnapshotHash,
		arg.CurrentHash,
	)
	var i InstallReportEntry
	err := row.Scan(
		&i.EntryID,
		&i.ReportID,
		&i.Filepath,
		&i.SnapshotHash,
		&i.CurrentHash,
	)
	return i, err
}

const createReportPlugin = `-- name: CreateReportPlugin :one
INSERT INTO install_report_entry_plugin (report_id, slug, issue) VALUES ($1, $2, $3) RETURNING entry_id, report_id, slug, issue
`

type CreateReportPluginParams struct {
	ReportID uuid.UUID       `json:"report_id"`
	Slug     string          `json:"slug"`
	Issue    json.RawMessage `json:"issue"`
}

func (q *Queries) CreateReportPlugin(ctx context.Context, arg CreateReportPluginParams) (InstallReportEntryPlugin, error) {
	row := q.db.QueryRowContext(ctx, createReportPlugin, arg.ReportID, arg.Slug, arg.Issue)
	var i InstallReportEntryPlugin
	err := row.Scan(
		&i.EntryID,
		&i.ReportID,
		&i.Slug,
		&i.Issue,
	)
	return i, err
}

const createSnapshot = `-- name: CreateSnapshot :one
INSERT INTO install_snapshot (created_at, created_by, install_id) VALUES ($1, $2, $3) RETURNING snapshot_id, install_id, created_at, created_by
`

type CreateSnapshotParams struct {
	CreatedAt time.Time `json:"created_at"`
	CreatedBy uuid.UUID `json:"created_by"`
	InstallID uuid.UUID `json:"install_id"`
}

func (q *Queries) CreateSnapshot(ctx context.Context, arg CreateSnapshotParams) (InstallSnapshot, error) {
	row := q.db.QueryRowContext(ctx, createSnapshot, arg.CreatedAt, arg.CreatedBy, arg.InstallID)
	var i InstallSnapshot
	err := row.Scan(
		&i.SnapshotID,
		&i.InstallID,
		&i.CreatedAt,
		&i.CreatedBy,
	)
	return i, err
}

const createSnapshotHash = `-- name: CreateSnapshotHash :one
INSERT INTO install_snapshot_hash_tree (snapshot_id, filepath, filehash) VALUES ($1, $2, $3) RETURNING hash_id, snapshot_id, filepath, filehash
`

type CreateSnapshotHashParams struct {
	SnapshotID uuid.UUID `json:"snapshot_id"`
	Filepath   string    `json:"filepath"`
	Filehash   string    `json:"filehash"`
}

func (q *Queries) CreateSnapshotHash(ctx context.Context, arg CreateSnapshotHashParams) (InstallSnapshotHashTree, error) {
	row := q.db.QueryRowContext(ctx, createSnapshotHash, arg.SnapshotID, arg.Filepath, arg.Filehash)
	var i InstallSnapshotHashTree
	err := row.Scan(
		&i.HashID,
		&i.SnapshotID,
		&i.Filepath,
		&i.Filehash,
	)
	return i, err
}

const createSnapshotPlugin = `-- name: CreateSnapshotPlugin :one
INSERT INTO install_snapshot_plugins (snapshot_id, nicename, slug, version, status) VALUES ($1, $2, $3, $4, $5) RETURNING plugin_id, snapshot_id, nicename, slug, version, status
`

type CreateSnapshotPluginParams struct {
	SnapshotID uuid.UUID `json:"snapshot_id"`
	Nicename   string    `json:"nicename"`
	Slug       string    `json:"slug"`
	Version    string    `json:"version"`
	Status     string    `json:"status"`
}

func (q *Queries) CreateSnapshotPlugin(ctx context.Context, arg CreateSnapshotPluginParams) (InstallSnapshotPlugin, error) {
	row := q.db.QueryRowContext(ctx, createSnapshotPlugin,
		arg.SnapshotID,
		arg.Nicename,
		arg.Slug,
		arg.Version,
		arg.Status,
	)
	var i InstallSnapshotPlugin
	err := row.Scan(
		&i.PluginID,
		&i.SnapshotID,
		&i.Nicename,
		&i.Slug,
		&i.Version,
		&i.Status,
	)
	return i, err
}

const getAllInstallIdentifiers = `-- name: GetAllInstallIdentifiers :many
SELECT install_id, ssh_username FROM install
`

type GetAllInstallIdentifiersRow struct {
	InstallID   uuid.UUID `json:"install_id"`
	SshUsername string    `json:"ssh_username"`
}

func (q *Queries) GetAllInstallIdentifiers(ctx context.Context) ([]GetAllInstallIdentifiersRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllInstallIdentifiers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllInstallIdentifiersRow
	for rows.Next() {
		var i GetAllInstallIdentifiersRow
		if err := rows.Scan(&i.InstallID, &i.SshUsername); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllInstalls = `-- name: GetAllInstalls :many
SELECT install_id, nicename, ssh_username, ssh_port, ssh_hostname, is_multisite FROM install
`

func (q *Queries) GetAllInstalls(ctx context.Context) ([]Install, error) {
	rows, err := q.db.QueryContext(ctx, getAllInstalls)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Install
	for rows.Next() {
		var i Install
		if err := rows.Scan(
			&i.InstallID,
			&i.Nicename,
			&i.SshUsername,
			&i.SshPort,
			&i.SshHostname,
			&i.IsMultisite,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getInstallByID = `-- name: GetInstallByID :one
SELECT install_id, nicename, ssh_username, ssh_port, ssh_hostname, is_multisite FROM install WHERE install_id = $1
`

func (q *Queries) GetInstallByID(ctx context.Context, installID uuid.UUID) (Install, error) {
	row := q.db.QueryRowContext(ctx, getInstallByID, installID)
	var i Install
	err := row.Scan(
		&i.InstallID,
		&i.Nicename,
		&i.SshUsername,
		&i.SshPort,
		&i.SshHostname,
		&i.IsMultisite,
	)
	return i, err
}

const getInstallByIdentifier = `-- name: GetInstallByIdentifier :one
SELECT install_id, nicename, ssh_username, ssh_port, ssh_hostname, is_multisite FROM install WHERE ssh_username = $1
`

func (q *Queries) GetInstallByIdentifier(ctx context.Context, sshUsername string) (Install, error) {
	row := q.db.QueryRowContext(ctx, getInstallByIdentifier, sshUsername)
	var i Install
	err := row.Scan(
		&i.InstallID,
		&i.Nicename,
		&i.SshUsername,
		&i.SshPort,
		&i.SshHostname,
		&i.IsMultisite,
	)
	return i, err
}

const getInstallIDForReportID = `-- name: GetInstallIDForReportID :one
SELECT install_id FROM install_report WHERE report_id = $1
`

func (q *Queries) GetInstallIDForReportID(ctx context.Context, reportID uuid.UUID) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, getInstallIDForReportID, reportID)
	var install_id uuid.UUID
	err := row.Scan(&install_id)
	return install_id, err
}

const getLatestSnapshotForInstallID = `-- name: GetLatestSnapshotForInstallID :one
SELECT snapshot_id, install_id, created_at, created_by FROM install_snapshot WHERE install_id = $1 ORDER BY created_at DESC LIMIT 1
`

func (q *Queries) GetLatestSnapshotForInstallID(ctx context.Context, installID uuid.UUID) (InstallSnapshot, error) {
	row := q.db.QueryRowContext(ctx, getLatestSnapshotForInstallID, installID)
	var i InstallSnapshot
	err := row.Scan(
		&i.SnapshotID,
		&i.InstallID,
		&i.CreatedAt,
		&i.CreatedBy,
	)
	return i, err
}

const getReportEntriesForReport = `-- name: GetReportEntriesForReport :many
SELECT entry_id, report_id, filepath, snapshot_hash, current_hash FROM install_report_entry WHERE report_id = $1
`

func (q *Queries) GetReportEntriesForReport(ctx context.Context, reportID uuid.UUID) ([]InstallReportEntry, error) {
	rows, err := q.db.QueryContext(ctx, getReportEntriesForReport, reportID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []InstallReportEntry
	for rows.Next() {
		var i InstallReportEntry
		if err := rows.Scan(
			&i.EntryID,
			&i.ReportID,
			&i.Filepath,
			&i.SnapshotHash,
			&i.CurrentHash,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSnapshotHashes = `-- name: GetSnapshotHashes :many
SELECT hash_id, snapshot_id, filepath, filehash FROM install_snapshot_hash_tree WHERE snapshot_id = $1
`

func (q *Queries) GetSnapshotHashes(ctx context.Context, snapshotID uuid.UUID) ([]InstallSnapshotHashTree, error) {
	rows, err := q.db.QueryContext(ctx, getSnapshotHashes, snapshotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []InstallSnapshotHashTree
	for rows.Next() {
		var i InstallSnapshotHashTree
		if err := rows.Scan(
			&i.HashID,
			&i.SnapshotID,
			&i.Filepath,
			&i.Filehash,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSnapshotPluginsForReport = `-- name: GetSnapshotPluginsForReport :many
SELECT plugin_id, snapshot_id, nicename, slug, version, status FROM install_snapshot_plugins WHERE snapshot_id = $1
`

func (q *Queries) GetSnapshotPluginsForReport(ctx context.Context, snapshotID uuid.UUID) ([]InstallSnapshotPlugin, error) {
	rows, err := q.db.QueryContext(ctx, getSnapshotPluginsForReport, snapshotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []InstallSnapshotPlugin
	for rows.Next() {
		var i InstallSnapshotPlugin
		if err := rows.Scan(
			&i.PluginID,
			&i.SnapshotID,
			&i.Nicename,
			&i.Slug,
			&i.Version,
			&i.Status,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}