package jobs

import (
	"context"

	"github.com/google/uuid"
	"github.com/jordanknott/monitor/internal/wordpress"
)

func (t *JobTasks) UploadManage(installIDEncoded string) (bool, error) {
	installID, err := uuid.Parse(installIDEncoded)
	if err != nil {
		return false, err
	}
	install, err := t.Data.GetInstallByID(context.Background(), installID)
	if err != nil {
		return false, err
	}
	err = wordpress.UploadFile(install, "./dist/manage", "./manage")
	if err != nil {
		return false, err
	}
	return true, err
}
