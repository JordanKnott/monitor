package monitorcli

import (
	"context"
	"net/http"
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jordanknott/monitor/internal/graph/client"
	"github.com/spf13/cobra"
)

func newInstallsListCmd() *cobra.Command {
	return &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			appConfig, err := GetAppConfig()
			if err != nil {
				return err
			}
			url, err := appConfig.GetURL("graphql")
			if err != nil {
				return err
			}

			c := client.NewClient(url, http.DefaultClient, appConfig.Token)

			r, _, err := client.GetAllInstalls(ctx, c)
			if err != nil {
				return err
			}

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"ID", "Nicename", "Username", "Hostname", "Port"})
			for _, install := range r.Installs {
				t.AppendRow(table.Row{
					install.ID,
					install.Nicename,
					install.Username,
					install.Hostname,
					install.Port,
				})
			}
			t.AppendFooter(table.Row{"", "", "", "Total", len(r.Installs)})
			t.Render()
			return nil
		},
	}
}
