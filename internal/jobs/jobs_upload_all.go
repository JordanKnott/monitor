package jobs

import (
	"context"

	"github.com/RichardKnop/machinery/v1/tasks"
)

func DoesExist(a []string, t string) bool {
	for _, b := range a {
		if b == t {
			return true
		}
	}
	return false
}

func (t *JobTasks) UploadAllInstalls() (bool, error) {
	installs, err := t.Data.GetAllInstalls(context.Background())
	if err != nil {
		return false, err
	}
	whitelist := []string{"drivendigitalrebuild", "lcdv2", "glacierchocolate"}
	for _, install := range installs {
		if DoesExist(whitelist, install.SshUsername) {
			_, err = t.Server.SendTask(&tasks.Signature{
				Name: "upload",
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
	return false, err
}
