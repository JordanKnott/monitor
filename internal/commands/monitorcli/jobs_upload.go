package monitorcli

import (
	"github.com/spf13/cobra"
)

func newJobsUploadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upload <installID>",
		Short: "Create a new migration file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			/*
				appConfig, err := config.GetAppConfig()
				if err != nil {
					return err
				}

				dbConn, err := sql.Open("postgres", appConfig.Database.GetDatabaseConnectionUri())
				if err != nil {
					return err
				}
				data := db.NewData(dbConn)
				jobConfig := appConfig.Job.GetJobConfig()
				server, err := machinery.NewServer(&jobConfig)
				if err != nil {
					// do something with the error
				}
				queueLog.Set(&jobs.MachineryLogger{})
				redisClient, err := appConfig.MessageQueue.GetMessageQueueClient()
				if err != nil {
					return err
				}
				if args[0] == "all" {
					jobs.RegisterTasks(server, *data, appConfig, redisClient)
					check := &tasks.Signature{
						Name: "uploadAll",
					}
					_, err = server.SendTask(check)
					if err != nil {
						return err
					}
					return nil
				}
				install, err := data.GetInstallByIdentifier(context.Background(), args[0])
				if err != nil {
					return err
				}
				jobs.RegisterTasks(server, *data, appConfig, redisClient)
				check := &tasks.Signature{
					Name: "upload",
					Args: []tasks.Arg{
						{
							Type:  "string",
							Value: install.InstallID,
						},
					},
				}
				_, err = server.SendTask(check)
				if err != nil {
					return err
				}
			*/
			return nil

		},
	}
}