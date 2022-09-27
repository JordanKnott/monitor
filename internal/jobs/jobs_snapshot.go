package jobs

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jordanknott/monitor/internal/db"
	"github.com/jordanknott/monitor/internal/wordpress"
)

func (t *JobTasks) CreateSnapshot(installIDEncoded string) (bool, error) {
	installID, err := uuid.Parse(installIDEncoded)
	if err != nil {
		return false, err
	}
	install, err := t.Data.GetInstallByID(context.Background(), installID)
	if err != nil {
		return false, err
	}
	ctx := context.Background()
	now := time.Now().UTC()

	user, err := t.Data.GetUserByEmail(ctx, "jordan@drivendigital.us")
	if err != nil {
		return false, err
	}
	snapshot, err := t.Data.CreateSnapshot(ctx, db.CreateSnapshotParams{InstallID: installID, CreatedAt: now, CreatedBy: user.UserID})
	if err != nil {
		return false, err
	}
	response, err := wordpress.RunScan(install)
	if err != nil {
		return false, err
	}
	for filename, hash := range response.Hashes {
		_, err := t.Data.CreateSnapshotHash(ctx, db.CreateSnapshotHashParams{
			SnapshotID: snapshot.SnapshotID,
			Filepath:   filename,
			Filehash:   hash,
		})
		if err != nil {
			return false, err
		}

	}
	for _, plugin := range response.Plugins {
		_, err := t.Data.CreateSnapshotPlugin(ctx, db.CreateSnapshotPluginParams{
			SnapshotID: snapshot.SnapshotID,
			Version:    plugin.Version,
			Status:     plugin.Status,
			Nicename:   plugin.Name,
			Slug:       plugin.Name,
		})
		if err != nil {
			return false, err
		}
	}
	return true, err
}
