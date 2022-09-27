package monitorcli

import "github.com/spf13/cobra"

func newInstallsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "installs",
	}
	cmd.AddCommand(newInstallsListCmd(), newInstallsSyncCmd())
	return cmd
}
