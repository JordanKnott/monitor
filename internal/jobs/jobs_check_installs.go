package jobs

import (
	"context"

	"github.com/RichardKnop/machinery/v1/tasks"
)

func (t *JobTasks) CheckAllInstalls() (bool, error) {
	installs, err := t.Data.GetAllInstalls(context.Background())
	if err != nil {
		return false, err
	}
	whitelist := []string{"drivendigitalrebuild", "lcdv2", "glacierchocolate"}
	for _, install := range installs {
		if DoesExist(whitelist, install.SshUsername) {
			_, err = t.Server.SendTask(&tasks.Signature{
				Name: "checkInstall",
				Args: []tasks.Arg{
					{
						Type:  "string",
						Value: install.InstallID,
					},
				},
			})
			if err != nil {
				return false, err
			}
		}
	}
	return true, nil
}
