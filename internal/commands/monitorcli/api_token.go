package monitorcli

import (
	"github.com/spf13/cobra"
)

func newApiTokenCmd() *cobra.Command {
	cc := &cobra.Command{
		Use:   "token",
		Short: "run the migrations",
		Long:  "Run the migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			/*
				appConfig, err := config.GetAppConfig()
				if err != nil {
					return err
				}

				fmt.Println(appConfig.Database.GetDatabaseConnectionUri())
				dbConn, err := sql.Open("postgres", appConfig.Database.GetDatabaseConnectionUri())
				if err != nil {
					return err
				}

				prompt := promptui.Prompt{
					Label: "Email for Token",
				}
				email, err := prompt.Run()

				data := db.NewData(dbConn)
				user, err := data.GetUserByEmail(context.Background(), email)
				if err != nil {
					return err
				}

				token, err := utils.GenerateToken()
				if err != nil {
					return err
				}
				now := time.Now().UTC()
				expiresAt := now.AddDate(0, 0, 1)

				err = data.CreateAccessToken(context.Background(), db.CreateAccessTokenParams{UserID: user.UserID, Token: token, CreatedAt: now, ExpiresAt: expiresAt})
				if err != nil {
					return err
				}
				fmt.Println("Token: " + token)
			*/
			return nil
		},
	}
	return cc
}
