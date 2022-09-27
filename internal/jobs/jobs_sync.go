package jobs

import (
	"context"

	"github.com/google/uuid"
	"github.com/jordanknott/monitor/internal/wordpress"
)

func (t *JobTasks) SyncInstall(installIDEncoded string) (bool, error) {
	installID, err := uuid.Parse(installIDEncoded)
	if err != nil {
		return false, err
	}
	install, err := t.Data.GetInstallByID(context.Background(), installID)
	if err != nil {
		return false, err
	}

	_, err = wordpress.RunSync(install)
	if err != nil {
		return false, err
	}
	return true, err
}
