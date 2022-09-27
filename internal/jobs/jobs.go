package jobs

import (
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/go-redis/redis/v8"

	"github.com/jordanknott/monitor/internal/config"
	"github.com/jordanknott/monitor/internal/db"
)

func RegisterPeriodicTasks(server *machinery.Server) error {
	signature := tasks.Signature{
		Name: "checkAllInstalls",
	}
	if !server.IsTaskRegistered("check-installs") {
		err := server.RegisterPeriodicTask("0 8 * * *", "check-installs", &signature)
		return err
	}
	return nil
}

func RegisterTasks(server *machinery.Server, data db.Data, appConfig config.AppConfig, messageQueue *redis.Client) {
	tasks := JobTasks{Data: data, Server: server, AppConfig: appConfig, MessageQueue: messageQueue}
	server.RegisterTasks(map[string]interface{}{
		"createSnapshot":    tasks.CreateSnapshot,
		"checkInstall":      tasks.CheckInstall,
		"upload":            tasks.UploadManage,
		"uploadAll":         tasks.UploadAllInstalls,
		"checkInstallEmail": tasks.CheckInstallEmailReport,
		"checkAllInstalls":  tasks.CheckAllInstalls,
	})
}

type JobTasks struct {
	AppConfig    config.AppConfig
	Data         db.Data
	Server       *machinery.Server
	MessageQueue *redis.Client
}
