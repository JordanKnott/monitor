package jobs

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jordanknott/monitor/internal/db"
	"github.com/jordanknott/monitor/internal/wordpress"
)

func GetPlugin(entries []wordpress.PluginEntry, target string) (bool, wordpress.PluginEntry) {
	for _, p := range entries {
		if p.Name == target {
			return true, p
		}
	}
	return false, wordpress.PluginEntry{}
}

type PluginStatus string

var (
	StatusRemoved  PluginStatus = "removed"
	StatusAdded    PluginStatus = "added"
	StatusUpgraded PluginStatus = "upgraded"
)

type VersionChange struct {
	Snapshot string
	Current  string
}

type PluginIssue struct {
	Status  PluginStatus
	Version *VersionChange
}

func (t *JobTasks) CheckInstall(installIDEncoded string) (string, error) {
	installID, err := uuid.Parse(installIDEncoded)
	if err != nil {
		return "", err
	}
	install, err := t.Data.GetInstallByID(context.Background(), installID)
	if err != nil {
		return "", err
	}
	ctx := context.Background()
	now := time.Now().UTC()

	snapshot, err := t.Data.GetLatestSnapshotForInstallID(ctx, installID)
	if err != nil {
		return "", err
	}
	snapshotHashes := make(map[string]string)
	hashes, err := t.Data.GetSnapshotHashes(ctx, snapshot.SnapshotID)
	if err != nil {
		return "", err
	}
	for _, hash := range hashes {
		snapshotHashes[hash.Filepath] = hash.Filehash
	}
	report, err := t.Data.CreateReport(ctx, db.CreateReportParams{InstallID: installID, CreatedAt: now, SnapshotID: snapshot.SnapshotID})
	if err != nil {
		return "", err
	}
	response, err := wordpress.RunScan(install)
	if err != nil {
		return "", err
	}
	snapshotPlugins, err := t.Data.GetSnapshotPluginsForReport(ctx, report.SnapshotID)
	if err != nil {
		return "", err
	}
	slugs := []string{}
	for _, sPlugin := range snapshotPlugins {
		slugs = append(slugs, sPlugin.Slug)
		if ok, p := GetPlugin(response.Plugins, sPlugin.Slug); ok {
			if p.Version != sPlugin.Version {
				raw, err := json.Marshal(PluginIssue{Status: StatusUpgraded,
					Version: &VersionChange{Snapshot: sPlugin.Version, Current: p.Version}})
				if err != nil {
					return "", err
				}
				t.Data.CreateReportPlugin(ctx, db.CreateReportPluginParams{
					ReportID: report.ReportID,
					Slug:     sPlugin.Slug,
					Issue:    raw,
				})
			}
		} else {
			raw, err := json.Marshal(PluginIssue{Status: StatusRemoved})
			if err != nil {
				return "", err
			}
			t.Data.CreateReportPlugin(ctx, db.CreateReportPluginParams{
				ReportID: report.ReportID,
				Slug:     sPlugin.Slug,
				Issue:    raw,
			})
		}
	}
	for _, plugin := range response.Plugins {
		if !DoesExist(slugs, plugin.Name) {
			raw, err := json.Marshal(PluginIssue{Status: StatusAdded})
			if err != nil {
				return "", err
			}
			t.Data.CreateReportPlugin(ctx, db.CreateReportPluginParams{
				ReportID: report.ReportID,
				Slug:     plugin.Name,
				Issue:    raw,
			})
		}
	}

	for filename, hash := range response.Hashes {
		if snapHash, ok := snapshotHashes[filename]; ok {
			if snapHash != hash {
				_, err := t.Data.CreateReportEntry(ctx, db.CreateReportEntryParams{
					ReportID:     report.ReportID,
					Filepath:     filename,
					SnapshotHash: snapHash,
					CurrentHash:  hash,
				})
				if err != nil {
					return "", err
				}
			}
		} else {
			_, err := t.Data.CreateReportEntry(ctx, db.CreateReportEntryParams{
				ReportID:     report.ReportID,
				Filepath:     filename,
				SnapshotHash: "",
				CurrentHash:  hash,
			})
			if err != nil {
				return "", err
			}

		}
	}
	return report.ReportID.String(), err
}
