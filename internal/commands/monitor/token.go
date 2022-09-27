package monitor

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jordanknott/monitor/internal/db"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newTokenCmd() *cobra.Command {
	cc := &cobra.Command{
		Use:   "token",
		Short: "Creates an application token",
		Args:  cobra.ExactArgs(1),
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
			dbConn, err := sql.Open("postgres", uri)
			if err != nil {
				logrus.Panic(err)
			}
			dbConn.SetMaxOpenConns(25)
			dbConn.SetMaxIdleConns(25)
			dbConn.SetConnMaxLifetime(5 * time.Minute)
			defer dbConn.Close()
			err = dbConn.Ping()
			if err != nil {
				logrus.WithError(err).Error("while pinging database")
				return err
			}

			data := db.NewData(dbConn)
			ctx := context.Background()
			user, err := data.GetUserByEmail(ctx, args[0])
			if err != nil {
				return err
			}
			now := time.Now().UTC()
			fmt.Printf("Creating token for %s\n", user.UserID)
			token, err := data.CreateAppToken(ctx, db.CreateAppTokenParams{CreatedAt: now, UserID: user.UserID})
			if err != nil {
				return err
			}
			fmt.Printf("Token: %s\n", token.TokenID)
			return nil

		},
	}
	return cc
}
