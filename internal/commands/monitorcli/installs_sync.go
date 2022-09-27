package monitorcli

import (
	"errors"
	"fmt"
	"os"

	"github.com/imroc/req/v3"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jordanknott/monitor/internal/api"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func newInstallsSyncCmd() *cobra.Command {
	return &cobra.Command{
		Use:  "sync <filepath>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := os.Stat(args[0]); err != nil {
				return errors.New(args[0] + " does not exist!")
			}
			f, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer f.Close()
			client := req.C()
			appConfig, err := GetAppConfig()
			if err != nil {
				return err
			}
			url, err := appConfig.GetURL("installs/sync/preview")
			if err != nil {
				return err
			}

			var result api.InstallSyncPreviewResponse
			r, err := client.R().SetFileReader("file", "installs.csv", f).SetResult(&result).SetHeader("Authorization", appConfig.Token).Post(url)
			if err != nil {
				return err
			}
			if r.StatusCode == 401 {
				return errors.New("authentication failed")
			} else if r.StatusCode != 200 {
				return fmt.Errorf("invalid status code %d", r.StatusCode)
			}
			fmt.Println(result)
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"Nicename", "Status"})
			for _, install := range result.ToBeCreated {
				t.AppendRow(table.Row{
					install.SshUsername,
					"Create",
				})
			}
			for _, install := range result.ToBeMerged {
				t.AppendRow(table.Row{
					install.SshUsername,
					"Merge",
				})
			}
			for _, install := range result.ToBeDeleted {
				t.AppendRow(table.Row{
					install.SshUsername,
					"Delete",
				})
			}
			total := len(result.ToBeCreated) + len(result.ToBeDeleted) + len(result.ToBeMerged)
			t.AppendFooter(table.Row{"Total Changes", total})
			t.AppendFooter(table.Row{"Will delete", len(result.ToBeDeleted)})
			t.Render()
			prompt := promptui.Prompt{
				Label:     "Sync installs",
				IsConfirm: true,
			}
			res, err := prompt.Run()
			if err != nil {
				return err
			}
			if res == "y" {
				f, err := os.Open(args[0])
				if err != nil {
					return err
				}
				defer f.Close()
				url, err = appConfig.GetURL("installs/sync")
				if err != nil {
					return err
				}
				var result api.InstallSyncPreviewResponse
				r, err := client.R().SetFileReader("file", "installs.csv", f).SetResult(&result).SetHeader("Authorization", appConfig.Token).Post(url)
				if err != nil {
					return err
				}
				if r.StatusCode != 200 {
					return fmt.Errorf("invalid status code %d", r.StatusCode)
				}
			}
			return nil
		},
	}
}
