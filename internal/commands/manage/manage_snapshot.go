package manage

import (
	"encoding/json"
	"fmt"

	"github.com/jordanknott/monitor/internal/wordpress"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newManageSnapshotCmd() *cobra.Command {
	return &cobra.Command{
		Use: "snapshot",
		RunE: func(cmd *cobra.Command, args []string) error {
			plugins, err := wordpress.GetPlugins()
			if err != nil {
				logrus.WithError(err).Error("while getting plugins")
				return err
			}
			hashtree := wordpress.MakeHashTree([]string{"public/wp-content/themes", "public/wp-includes"})
			raw, err := json.Marshal(wordpress.SnapshotResponse{
				Hashes:  hashtree,
				Plugins: plugins,
			})
			if err != nil {
				logrus.WithError(err).Error("while marshal response")
				return err

			}
			fmt.Println(string(raw))
			return nil
		},
	}

}
