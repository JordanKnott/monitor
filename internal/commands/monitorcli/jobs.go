package monitorcli

import (
	"github.com/spf13/cobra"
)

func newJobsCmd() *cobra.Command {
	cc := &cobra.Command{
		Use:   "jobs",
		Short: "run the migrations",
		Long:  "Run the migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			/*
				appConfig, err := config.GetAppConfig()
				if err != nil {
					return err
				}

				fmt.Println(appConfig.Database.GetDatabaseConnectionUri())
				db, err := sql.Open("postgres", appConfig.Database.GetDatabaseConnectionUri())
				if err != nil {
					return err
				}
				return goose.Up(db, "migrations")
			*/
			return nil
		},
	}
	cc.AddCommand(newJobsSnapshotCmd(), newJobsCheckCmd(), newJobsUploadCmd())
	return cc
}
