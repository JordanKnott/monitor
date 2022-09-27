package monitor

import (
	"database/sql"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/RichardKnop/machinery/v1"
	queueLog "github.com/RichardKnop/machinery/v1/log"
	data "github.com/jordanknott/monitor/internal/db"
	"github.com/jordanknott/monitor/internal/jobs"
)

func newWorkerCmd() *cobra.Command {
	cc := &cobra.Command{
		Use:   "worker",
		Short: "Run the task queue worker",
		Long:  "Run the task queue worker",
		RunE: func(cmd *cobra.Command, args []string) error {
			Formatter := new(logrus.TextFormatter)
			Formatter.TimestampFormat = "02-01-2006 15:04:05"
			Formatter.FullTimestamp = true
			logrus.SetFormatter(Formatter)
			logrus.SetLevel(logrus.InfoLevel)

			appConfig, err := GetAppConfig()
			if err != nil {
				logrus.Panic(err)
			}
			uri := appConfig.Database.GetDatabaseConnectionUri()
			logrus.WithFields(logrus.Fields{
				"uri": uri,
			}).Info("connecting to database")
			db, err := sql.Open("postgres", uri)
			if err != nil {
				logrus.Panic(err)
			}
			db.SetMaxOpenConns(25)
			db.SetMaxIdleConns(25)
			db.SetConnMaxLifetime(5 * time.Minute)
			defer db.Close()
			err = db.Ping()
			if err != nil {
				logrus.WithError(err).Error("while pinging database")
				return err
			}

			logrus.Info("starting task queue server instance")
			jobConfig := appConfig.Job.GetJobConfig()
			logrus.WithField("broker", jobConfig.Broker).Info("connecting to redis")
			server, err := machinery.NewServer(&jobConfig)
			if err != nil {
				logrus.WithError(err).Error("while starting machinery server")
				return err
			}
			queueLog.Set(&jobs.MachineryLogger{})
			repo := *data.NewData(db)
			logrus.WithField("messageQueue", appConfig.MessageQueue.URI).Info("connecting to message queue")
			redisClient, err := appConfig.MessageQueue.GetMessageQueueClient()
			if err != nil {
				logrus.WithError(err).Error("while connecting message queue")
				return err
			}
			jobs.RegisterTasks(server, repo, appConfig, redisClient)
			jobs.RegisterPeriodicTasks(server)

			worker := server.NewWorker("taskcafe_worker", 3)
			logrus.Info("starting task queue worker")
			err = worker.Launch()
			if err != nil {
				logrus.WithError(err).Error("error while launching ")
				return err
			}
			return nil
		},
	}
	return cc
}
