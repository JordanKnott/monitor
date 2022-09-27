package manage

import (
	"encoding/json"
	"fmt"

	"github.com/jordanknott/monitor/internal/wordpress"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newManageSyncCmd() *cobra.Command {
	return &cobra.Command{
		Use: "sync",
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := wordpress.GenerateResponse(wordpress.WPConfig{ProtectedPaths: []string{"public/wp-content/themes", "public/wp-includes"}})
			if err != nil {
				logrus.WithError(err).Error("while marshal response")
				return err
			}
			raw, err := json.Marshal(response)
			if err != nil {
				logrus.WithError(err).Error("while marshal response")
				return err
			}
			fmt.Println(string(raw))
			return nil
		},
	}

}
